package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"tuneshift/internal/auth"
	"tuneshift/internal/config"
	"tuneshift/internal/migration"
	"tuneshift/internal/source"
)

type uploadEntry struct {
	playlists []source.Playlist
	createdAt time.Time
}

type Handler struct {
	sessions  *auth.SessionManager
	tidalAuth *auth.TidalAuth
	cfg       *config.Config

	mu        sync.Mutex
	reporters map[string]*migration.ProgressReporter
	uploads   map[string]uploadEntry
}

func New(sessions *auth.SessionManager, tidalAuth *auth.TidalAuth, cfg *config.Config) *Handler {
	h := &Handler{
		sessions:  sessions,
		tidalAuth: tidalAuth,
		cfg:       cfg,
		reporters: make(map[string]*migration.ProgressReporter),
		uploads:   make(map[string]uploadEntry),
	}

	// Cleanup expired uploads every 5 minutes
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			h.mu.Lock()
			for id, entry := range h.uploads {
				if time.Since(entry.createdAt) > 30*time.Minute {
					delete(h.uploads, id)
				}
			}
			h.mu.Unlock()
		}
	}()

	return h
}

func (h *Handler) isSecure() bool {
	return strings.HasPrefix(h.cfg.BaseURL, "https")
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func generateSessionID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
