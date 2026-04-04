package handlers

import (
	"encoding/json"
	"net/http"

	"tuneshift/internal/youtube"
)

func (h *Handler) YouTubePlaylists(w http.ResponseWriter, r *http.Request) {
	token, err := h.sessions.GetTokenCookie(r, "tuneshift_google")
	if err != nil {
		writeError(w, http.StatusUnauthorized, "not connected to YouTube Music")
		return
	}

	client := youtube.NewClient(token.AccessToken)
	playlists, err := client.ListPlaylists()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch playlists: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"playlists": playlists,
	})
}

func (h *Handler) YouTubeFetch(w http.ResponseWriter, r *http.Request) {
	token, err := h.sessions.GetTokenCookie(r, "tuneshift_google")
	if err != nil {
		writeError(w, http.StatusUnauthorized, "not connected to YouTube Music")
		return
	}

	var req struct {
		Playlists []youtube.PlaylistSummary `json:"playlists"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if len(req.Playlists) == 0 {
		writeError(w, http.StatusBadRequest, "no playlists selected")
		return
	}

	client := youtube.NewClient(token.AccessToken)
	playlists, err := youtube.FetchPlaylists(client, req.Playlists)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch playlist tracks: "+err.Error())
		return
	}

	sessionID, err := h.storeUploadedPlaylists(playlists)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to store playlists")
		return
	}

	type playlistSummary struct {
		Name       string `json:"name"`
		TrackCount int    `json:"track_count"`
	}

	summaries := make([]playlistSummary, len(playlists))
	totalTracks := 0
	for i, pl := range playlists {
		summaries[i] = playlistSummary{
			Name:       pl.Name,
			TrackCount: len(pl.Tracks),
		}
		totalTracks += len(pl.Tracks)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"session_id":   sessionID,
		"playlists":    summaries,
		"total_tracks": totalTracks,
	})
}
