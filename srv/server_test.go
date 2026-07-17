package srv

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func newTestServer(t *testing.T) *Server {
	t.Helper()
	tempDB := filepath.Join(t.TempDir(), "test.sqlite3")
	t.Cleanup(func() { os.Remove(tempDB) })

	server, err := New(tempDB, "test-hostname")
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}
	return server
}

func TestCreateAndListTodos(t *testing.T) {
	s := newTestServer(t)

	// Create a todo
	body := `{"item_type":"todo","title":"Buy milk","content":"From the store","priority":1}`
	req := httptest.NewRequest(http.MethodPost, "/api/items", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	s.handleCreateItem(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var created map[string]any
	if err := json.NewDecoder(w.Body).Decode(&created); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if created["title"] != "Buy milk" {
		t.Errorf("expected title 'Buy milk', got %v", created["title"])
	}

	// List todos
	req2 := httptest.NewRequest(http.MethodGet, "/api/items?type=todo", nil)
	w2 := httptest.NewRecorder()
	s.handleListItems(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w2.Code)
	}

	var items []map[string]any
	if err := json.NewDecoder(w2.Body).Decode(&items); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("expected 1 item, got %d", len(items))
	}
}

func TestToggleTodo(t *testing.T) {
	s := newTestServer(t)

	// Create a todo
	body := `{"item_type":"todo","title":"Test toggle"}`
	req := httptest.NewRequest(http.MethodPost, "/api/items", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	s.handleCreateItem(w, req)

	var created map[string]any
	_ = json.NewDecoder(w.Body).Decode(&created)

	// Toggle
	req2 := httptest.NewRequest(http.MethodPost, "/api/items/1/toggle", nil)
	req2.SetPathValue("id", "1")
	w2 := httptest.NewRecorder()
	s.handleToggleTodo(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w2.Code, w2.Body.String())
	}

	var toggled map[string]any
	_ = json.NewDecoder(w2.Body).Decode(&toggled)
	if toggled["done"].(float64) != 1 {
		t.Errorf("expected done=1 after toggle, got %v", toggled["done"])
	}
}

func TestReviewItem(t *testing.T) {
	s := newTestServer(t)

	// Create a note
	body := `{"item_type":"note","title":"Flashcard"}`
	req := httptest.NewRequest(http.MethodPost, "/api/items", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	s.handleCreateItem(w, req)

	// Review with rating 4
	reviewBody := `{"rating":4}`
	req2 := httptest.NewRequest(http.MethodPost, "/api/items/1/review", bytes.NewBufferString(reviewBody))
	req2.SetPathValue("id", "1")
	w2 := httptest.NewRecorder()
	s.handleReviewItem(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w2.Code, w2.Body.String())
	}

	var reviewed map[string]any
	_ = json.NewDecoder(w2.Body).Decode(&reviewed)
	if reviewed["repetitions"].(float64) != 1 {
		t.Errorf("expected repetitions=1, got %v", reviewed["repetitions"])
	}
	if reviewed["interval_days"].(float64) != 1 {
		t.Errorf("expected interval_days=1, got %v", reviewed["interval_days"])
	}
}

func TestSM2Algorithm(t *testing.T) {
	tests := []struct {
		name         string
		quality      int
		easeFactor   float64
		interval     float64
		repetitions  int
		wantEF       float64
		wantInterval float64
		wantReps     int
	}{
		{"fail resets", 2, 2.5, 6, 3, 2.5, 0, 0},
		{"first success", 4, 2.5, 0, 0, 2.5, 1, 1},
		{"second success", 4, 2.5, 1, 1, 2.5, 6, 2},
		{"third success", 4, 2.5, 6, 2, 2.5, 15, 3},
		{"perfect score", 5, 2.5, 6, 2, 2.6, 15, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ef, iv, reps := calculateSM2(tt.quality, tt.easeFactor, tt.interval, tt.repetitions)
			if reps != tt.wantReps {
				t.Errorf("reps: got %d, want %d", reps, tt.wantReps)
			}
			if iv != tt.wantInterval {
				t.Errorf("interval: got %f, want %f", iv, tt.wantInterval)
			}
			// Allow small float tolerance
			if ef < tt.wantEF-0.01 || ef > tt.wantEF+0.01 {
				t.Errorf("ease factor: got %f, want ~%f", ef, tt.wantEF)
			}
		})
	}
}

func TestArchiveItem(t *testing.T) {
	s := newTestServer(t)

	// Create
	body := `{"item_type":"todo","title":"To archive"}`
	req := httptest.NewRequest(http.MethodPost, "/api/items", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	s.handleCreateItem(w, req)

	// Archive
	req2 := httptest.NewRequest(http.MethodDelete, "/api/items/1", nil)
	req2.SetPathValue("id", "1")
	w2 := httptest.NewRecorder()
	s.handleArchiveItem(w2, req2)

	if w2.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w2.Code)
	}

	// List should be empty
	req3 := httptest.NewRequest(http.MethodGet, "/api/items?type=todo", nil)
	w3 := httptest.NewRecorder()
	s.handleListItems(w3, req3)

	var items []map[string]any
	_ = json.NewDecoder(w3.Body).Decode(&items)
	if len(items) != 0 {
		t.Errorf("expected 0 items after archive, got %d", len(items))
	}
}

func TestSearchItems(t *testing.T) {
	s := newTestServer(t)

	// Create items
	for _, title := range []string{"Go programming", "Rust language", "Go concurrency"} {
		body := `{"item_type":"note","title":"` + title + `"}`
		req := httptest.NewRequest(http.MethodPost, "/api/items", bytes.NewBufferString(body))
		w := httptest.NewRecorder()
		s.handleCreateItem(w, req)
	}

	// Search for "Go"
	req := httptest.NewRequest(http.MethodGet, "/api/search?q=Go", nil)
	w := httptest.NewRecorder()
	s.handleSearch(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var items []map[string]any
	_ = json.NewDecoder(w.Body).Decode(&items)
	if len(items) != 2 {
		t.Errorf("expected 2 results for 'Go', got %d", len(items))
	}
}

func TestListItemsRequiresType(t *testing.T) {
	s := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/items", nil)
	w := httptest.NewRecorder()
	s.handleListItems(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 without type param, got %d", w.Code)
	}
}
