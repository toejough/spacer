package srv

import (
	"log/slog"
	"net/http"
)

const (
	StaticDir    = "/app/public"
	TemplatesDir = "/app/public"
)

type Server struct {
	Hostname string
}

func New(hostname string) (*Server, error) {
	return &Server{Hostname: hostname}, nil
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
