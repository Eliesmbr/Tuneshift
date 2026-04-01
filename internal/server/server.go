package server

import (
	"context"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"time"

	"tuneshift/internal/auth"
	"tuneshift/internal/config"
	"tuneshift/internal/handlers"
	"tuneshift/internal/middleware"
)

type Server struct {
	httpServer *http.Server
	cfg        *config.Config
}

func New(cfg *config.Config) (*Server, error) {
	sessionManager := auth.NewSessionManager(cfg.EncryptionKey)

	tidalAuth := auth.NewTidalAuth(
		cfg.TidalClientID,
		cfg.BaseURL+"/api/auth/tidal/callback",
	)

	h := handlers.New(sessionManager, tidalAuth, cfg)

	mux := http.NewServeMux()
	registerRoutes(mux, h)

	handler := middleware.Chain(
		mux,
		middleware.Logger,
		middleware.CORS(cfg.BaseURL),
		middleware.RateLimit(100, time.Minute),
	)

	return &Server{
		httpServer: &http.Server{
			Addr:         cfg.Address(),
			Handler:      handler,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 120 * time.Second, // Long for SSE
			IdleTimeout:  60 * time.Second,
		},
		cfg: cfg,
	}, nil
}

func registerRoutes(mux *http.ServeMux, h *handlers.Handler) {
	// Health
	mux.HandleFunc("GET /api/health", h.Health)

	// CSV Upload
	mux.HandleFunc("POST /api/upload", h.UploadCSV)

	// Tidal Auth
	mux.HandleFunc("GET /api/auth/tidal/login", h.TidalLogin)
	mux.HandleFunc("GET /api/auth/tidal/callback", h.TidalCallback)
	mux.HandleFunc("GET /api/auth/tidal/status", h.TidalStatus)
	mux.HandleFunc("POST /api/auth/tidal/logout", h.TidalLogout)

	// Migration
	mux.HandleFunc("POST /api/migrate", h.StartMigration)
	mux.HandleFunc("GET /api/migrate/progress", h.MigrationProgress)

	// SPA fallback — serve frontend
	mux.HandleFunc("/", spaHandler())
}

func spaHandler() http.HandlerFunc {
	distPath := "web/dist"
	if _, err := os.Stat(distPath); err != nil {
		return func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api/") {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(`<!DOCTYPE html><html><body><h1>Tuneshift</h1><p>Frontend not built. Run npm build in web/</p></body></html>`))
		}
	}

	fsys := os.DirFS(distPath)
	fileServer := http.FileServerFS(fsys)

	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		// Reject path traversal attempts
		if strings.Contains(path, "..") {
			http.NotFound(w, r)
			return
		}

		if _, err := fs.Stat(fsys, path); err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}

		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
