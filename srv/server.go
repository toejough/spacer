package srv

import (
	"fmt"
	"log/slog"
	"net/http"

	"srv.exe.dev/db"
)

const (
	StaticDir    = "/app/public"
	TemplatesDir = "/app/public"
)

type Server struct {
	Hostname string
	DBPath   string
}

func New(dbPath, hostname string) (*Server, error) {
	srv := &Server{
		Hostname: hostname,
		DBPath:   dbPath,
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

	// Serve static files from /app/public with no directory listing
	staticHandler := noStoreHeader(http.FileServer(http.Dir(StaticDir)))
	mux.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	// Serve sw.js from root so it can control the whole origin
	mux.HandleFunc("GET /sw.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		http.ServeFile(w, r, StaticDir+"/sw.js")
	})

	slog.Info("starting server", "addr", addr)
	return http.ListenAndServe(addr, mux)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	http.ServeFile(w, r, TemplatesDir+"/index.html")
}

// noStoreHeader wraps a handler to disable caching, so Cloudflare's edge
// cache doesn't hide deploys behind stale static assets.
func noStoreHeader(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		h.ServeHTTP(w, r)
	})
}
