package srv

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"srv.exe.dev/db"
)

func TestTodosAPI(t *testing.T) {
	// create temp db file
	f, err := os.CreateTemp("", "testdb-*.sqlite")
	if err != nil {
		t.Fatalf("tmp db: %v", err)
	}
	f.Close()
	defer os.Remove(f.Name())
	// run migrations
	dbconn, err := db.Open(f.Name())
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.RunMigrations(dbconn); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	dbconn.Close()
	// insert a todo item
	dbconn, _ = db.Open(f.Name())
	_, err = dbconn.Exec("INSERT INTO items (item_type,title,content) VALUES ('todo','test','')")
	if err != nil {
		t.Fatalf("insert todo: %v", err)
	}
	dbconn.Close()
	// start server
	srv, err := New(f.Name(), "localhost")
	if err != nil {
		t.Fatalf("new srv: %v", err)
	}
	r := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.Serve(r.Context().Value("addr").(string))
	}))
	defer r.Close()
	// Instead of running the full server, directly call handler
	req := httptest.NewRequest("GET", "/api/todos", nil)
	w := httptest.NewRecorder()
	srv.handleTodos(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 got %d", res.StatusCode)
	}
	var todos []map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&todos); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(todos) != 1 {
		t.Fatalf("expected 1 todo got %d", len(todos))
	}
}

func TestPatchTodo(t *testing.T) {
	f, err := os.CreateTemp("", "testdb-*.sqlite")
	if err != nil {
		t.Fatalf("tmp db: %v", err)
	}
	f.Close()
	defer os.Remove(f.Name())
	dbconn, err := db.Open(f.Name())
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.RunMigrations(dbconn); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	// insert a todo
	res, err := dbconn.Exec("INSERT INTO items (item_type,title,content) VALUES ('todo','patch-test','')")
	if err != nil {
		t.Fatalf("insert: %v", err)
	}
	id, _ := res.LastInsertId()
	dbconn.Close()

	srv, err := New(f.Name(), "localhost")
	if err != nil {
		t.Fatalf("new srv: %v", err)
	}
	// PATCH to mark done
	body := `{"status":"done"}`
	req := httptest.NewRequest("PATCH", "/api/todos/"+fmt.Sprint(id), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	srv.handleTodo(w, req)
	if w.Result().StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204 got %d", w.Result().StatusCode)
	}
	// verify DB
	dbconn2, _ := db.Open(f.Name())
	var done int
	var completedAt sql.NullString
	row := dbconn2.QueryRow("SELECT done, completed_at FROM items WHERE id = ?", id)
	if err := row.Scan(&done, &completedAt); err != nil {
		t.Fatalf("scan: %v", err)
	}
	if done != 1 {
		t.Fatalf("expected done=1 got %d", done)
	}
	if !completedAt.Valid {
		t.Fatalf("expected completed_at to be set")
	}
}
