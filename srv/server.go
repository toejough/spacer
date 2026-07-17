package srv

import (
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"runtime"

	"srv.exe.dev/db"
)

type Server struct {
	Hostname     string
	TemplatesDir string
	StaticDir    string
}

func New(dbPath, hostname string) (*Server, error) {
	_, thisFile, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(thisFile)
	srv := &Server{
		Hostname:     hostname,
		TemplatesDir: filepath.Join(baseDir, "templates"),
		StaticDir:    filepath.Join(baseDir, "static"),
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
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(s.StaticDir))))
	// Serve sw.js from root so it can control the whole origin
	mux.HandleFunc("GET /sw.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(s.StaticDir, "sw.js"))
	})
	slog.Info("starting server", "addr", addr)
	return http.ListenAndServe(addr, mux)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(s.TemplatesDir, "index.html")
	http.ServeFile(w, r, path)
}
