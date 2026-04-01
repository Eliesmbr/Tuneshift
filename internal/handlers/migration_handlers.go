package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"tuneshift/internal/csvimport"
	"tuneshift/internal/migration"
	"tuneshift/internal/tidal"
)

type MigrateRequest struct {
	UploadSessionID string   `json:"upload_session_id"`
	Playlists       []string `json:"playlists"` // playlist names to migrate (empty = all)
}

func (h *Handler) StartMigration(w http.ResponseWriter, r *http.Request) {
	tidalToken, err := h.sessions.GetTokenCookie(r, "tuneshift_tidal")
	if err != nil {
		writeError(w, http.StatusUnauthorized, "not connected to Tidal")
		return
	}

	var req MigrateRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	playlists, ok := h.getUploadedPlaylists(req.UploadSessionID)
	if !ok || len(playlists) == 0 {
		writeError(w, http.StatusGone, "upload_expired")
		return
	}

	// Filter playlists if specific ones were selected
	if len(req.Playlists) > 0 {
		selected := make(map[string]bool, len(req.Playlists))
		for _, name := range req.Playlists {
			selected[name] = true
		}
		var filtered []csvimport.Playlist
		for _, pl := range playlists {
			if selected[pl.Name] {
				filtered = append(filtered, pl)
			}
		}
		playlists = filtered
	}

	if len(playlists) == 0 {
		writeError(w, http.StatusBadRequest, "no playlists selected")
		return
	}

	sessionID, err := generateSessionID()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate session ID")
		return
	}

	progress := migration.NewProgressReporter()

	h.mu.Lock()
	h.reporters[sessionID] = progress
	h.mu.Unlock()

	countryCode := tidalToken.CountryCode
	if countryCode == "" {
		countryCode = "US"
	}
	tidalClient := tidal.NewClient(tidalToken.AccessToken, tidalToken.UserID, countryCode)
	engine := migration.NewEngine(tidalClient, progress)

	go func() {
		defer func() {
			progress.Close()
			time.Sleep(30 * time.Second)
			h.mu.Lock()
			delete(h.reporters, sessionID)
			delete(h.uploads, req.UploadSessionID)
			h.mu.Unlock()
		}()

		ctx := context.Background()
		result, err := engine.RunFromCSV(ctx, playlists)
		if err != nil {
			progress.Send(migration.ProgressEvent{
				Type:    "error",
				Message: "Migration failed: " + err.Error(),
			})
			return
		}

		progress.Send(migration.ProgressEvent{
			Type:    "result",
			Message: fmt.Sprintf("Matched %d/%d tracks across %d playlists", result.MatchedTracks, result.TotalTracks, result.Playlists),
			Current: result.MatchedTracks,
			Total:   result.TotalTracks,
		})
	}()

	writeJSON(w, http.StatusOK, map[string]string{
		"session_id": sessionID,
	})
}

func (h *Handler) MigrationProgress(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		writeError(w, http.StatusBadRequest, "session_id required")
		return
	}

	h.mu.Lock()
	progress, ok := h.reporters[sessionID]
	h.mu.Unlock()

	if !ok {
		writeError(w, http.StatusNotFound, "migration session not found")
		return
	}

	progress.ServeSSE(w, r)
}
