package handlers

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"tuneshift/internal/source"
	"tuneshift/internal/source/exportify"
	"tuneshift/internal/source/takeout"
)

const (
	maxUploadSize = 50 << 20 // 50MB
	maxFileCount  = 50
)

func (h *Handler) UploadCSV(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		writeError(w, http.StatusBadRequest, "file too large (max 50MB)")
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		writeError(w, http.StatusBadRequest, "no files uploaded")
		return
	}
	if len(files) > maxFileCount {
		writeError(w, http.StatusBadRequest, "too many files (max 50)")
		return
	}

	var playlists []source.Playlist

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to read file")
			return
		}
		defer file.Close()

		name := strings.TrimSuffix(fileHeader.Filename, filepath.Ext(fileHeader.Filename))

		playlist, err := exportify.ParseCSV(file, name)
		if err != nil {
			writeError(w, http.StatusBadRequest, "failed to parse CSV file: "+err.Error())
			return
		}

		playlists = append(playlists, *playlist)
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

	sessionID, err := h.storeUploadedPlaylists(playlists)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to store upload")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"session_id":   sessionID,
		"playlists":    summaries,
		"total_tracks": totalTracks,
	})
}

func (h *Handler) storeUploadedPlaylists(playlists []source.Playlist) (string, error) {
	id, err := generateSessionID()
	if err != nil {
		return "", err
	}

	h.mu.Lock()
	h.uploads[id] = uploadEntry{
		playlists: playlists,
		createdAt: time.Now(),
	}
	h.mu.Unlock()

	return id, nil
}

func (h *Handler) getUploadedPlaylists(sessionID string) ([]source.Playlist, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	entry, ok := h.uploads[sessionID]
	if !ok {
		return nil, false
	}
	return entry.playlists, true
}

func (h *Handler) UploadTakeout(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		writeError(w, http.StatusBadRequest, "file too large (max 50MB)")
		return
	}

	files := r.MultipartForm.File["file"]
	if len(files) == 0 {
		writeError(w, http.StatusBadRequest, "no file uploaded")
		return
	}

	fileHeader := files[0]
	if !strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".zip") {
		writeError(w, http.StatusBadRequest, "expected a .zip file")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read file")
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read file")
		return
	}

	playlists, err := takeout.ParseZip(data)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
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

	sessionID, err := h.storeUploadedPlaylists(playlists)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to store upload")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"session_id":   sessionID,
		"playlists":    summaries,
		"total_tracks": totalTracks,
	})
}
