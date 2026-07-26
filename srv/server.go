package srv

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"

	"srv.exe.dev/db"
)

type Server struct {
	Hostname     string
	TemplatesDir string
	StaticDir    string
	DBPath       string
}

func New(dbPath, hostname string) (*Server, error) {
	_, thisFile, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(thisFile)
	srv := &Server{
		Hostname:     hostname,
		TemplatesDir: filepath.Join(baseDir, "templates"),
		StaticDir:    filepath.Join(baseDir, "static"),
		DBPath:       dbPath,
	}
	// Still open DB so migrations table exists (template requirement),
	// but we don't use it for items anymore.
	if wdb, err := db.Open(dbPath); err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	} else {
		if err := db.RunMigrations(wdb); err != nil {
			return nil, fmt.Errorf("failed to run migrations: %w", err)
		}
		wdb.Close()
	}
	return srv, nil
}

func (s *Server) Serve(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.handleIndex)
	mux.HandleFunc("/api/todos", s.handleTodos)
	mux.HandleFunc("/api/todos/", s.handleTodo)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(s.StaticDir))))
	// Serve sw.js from root so it can control the whole origin
	mux.HandleFunc("GET /sw.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		http.ServeFile(w, r, filepath.Join(s.StaticDir, "sw.js"))
	})
	slog.Info("starting server", "addr", addr)
	return http.ListenAndServe(addr, mux)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	path := filepath.Join(s.TemplatesDir, "index.html")
	http.ServeFile(w, r, path)
}

// handleTodos supports GET /api/todos
func (s *Server) handleTodos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	dbConn, err := db.Open(s.DBPath)
	if err != nil {
		http.Error(w, "failed to open db", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()
	rows, err := dbConn.Query("SELECT id, title, content, done, priority, due_date, tags, created_at, updated_at, archived, completed_at, archived_at FROM items WHERE item_type = 'todo'")
	if err != nil {
		http.Error(w, "db query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	type Todo struct {
		ID          int64   `json:"id"`
		Title       string  `json:"title"`
		Content     string  `json:"content"`
		Done        int     `json:"done"`
		Priority    int     `json:"priority"`
		DueDate     *string `json:"due_date"`
		Tags        string  `json:"tags"`
		CreatedAt   string  `json:"created_at"`
		UpdatedAt   string  `json:"updated_at"`
		Archived    int     `json:"archived"`
		CompletedAt *string `json:"completed_at"`
		ArchivedAt  *string `json:"archived_at"`
	}
	var todos []Todo
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Content, &t.Done, &t.Priority, &t.DueDate, &t.Tags, &t.CreatedAt, &t.UpdatedAt, &t.Archived, &t.CompletedAt, &t.ArchivedAt); err != nil {
			http.Error(w, "db scan failed", http.StatusInternalServerError)
			return
		}
		todos = append(todos, t)
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(todos)
}

// handleTodo supports PATCH /api/todos/:id
func (s *Server) handleTodo(w http.ResponseWriter, r *http.Request) {
	// Expect path /api/todos/{id}
	idStr := filepath.Base(r.URL.Path)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if r.Method != http.MethodPatch {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var payload struct {
		Status      string  `json:"status"`
		CompletedAt *string `json:"completed_at"`
		ArchivedAt  *string `json:"archived_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	dbConn, err := db.Open(s.DBPath)
	if err != nil {
		http.Error(w, "failed to open db", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()
	// map status to done/archived booleans and timestamps
	tx, err := dbConn.Begin()
	if err != nil {
		http.Error(w, "db tx failed", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	switch payload.Status {
	case "done":
		if payload.CompletedAt != nil {
			_, err = tx.Exec("UPDATE items SET done = 1, completed_at = ?, archived = 0, archived_at = NULL, updated_at = datetime('now') WHERE id = ?", *payload.CompletedAt, id)
		} else {
			_, err = tx.Exec("UPDATE items SET done = 1, completed_at = datetime('now'), archived = 0, archived_at = NULL, updated_at = datetime('now') WHERE id = ?", id)
		}
	case "archived", "abandoned":
		if payload.ArchivedAt != nil {
			_, err = tx.Exec("UPDATE items SET archived = 1, archived_at = ?, done = 0, completed_at = NULL, updated_at = datetime('now') WHERE id = ?", *payload.ArchivedAt, id)
		} else {
			_, err = tx.Exec("UPDATE items SET archived = 1, archived_at = datetime('now'), done = 0, completed_at = NULL, updated_at = datetime('now') WHERE id = ?", id)
		}
	case "open", "reopen":
		_, err = tx.Exec("UPDATE items SET done = 0, archived = 0, completed_at = NULL, archived_at = NULL, updated_at = datetime('now') WHERE id = ?", id)
	default:
		http.Error(w, "invalid status", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "db update failed", http.StatusInternalServerError)
		return
	}
	if err := tx.Commit(); err != nil {
		http.Error(w, "db commit failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
