package takeout

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"tuneshift/internal/source"
)

// libraryEntry holds metadata from the music library songs CSV.
type libraryEntry struct {
	Title  string
	Album  string
	Artist string
}

// ParseZip extracts playlists from a Google Takeout ZIP.
// It joins playlist video IDs against the music library for metadata.
func ParseZip(data []byte) ([]source.Playlist, error) {
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to open zip: %w", err)
	}

	library := make(map[string]libraryEntry) // video ID -> metadata
	var playlistFiles []*zip.File
	var playlistIndex *zip.File

	for _, f := range r.File {
		name := filepath.ToSlash(f.Name)

		if strings.HasSuffix(name, "music library songs.csv") || strings.HasSuffix(name, "music-library-songs.csv") {
			entries, err := parseLibraryCSV(f)
			if err != nil {
				return nil, fmt.Errorf("failed to parse music library: %w", err)
			}
			for k, v := range entries {
				library[k] = v
			}
		} else if strings.HasSuffix(name, "-videos.csv") {
			playlistFiles = append(playlistFiles, f)
		} else if strings.HasSuffix(name, "playlists.csv") || strings.HasSuffix(name, "playlists/playlists.csv") {
			playlistIndex = f
		}
	}

	// Build playlist name lookup from the index CSV
	playlistNames := make(map[string]string) // playlist ID -> title
	if playlistIndex != nil {
		playlistNames, _ = parsePlaylistIndex(playlistIndex)
	}

	var playlists []source.Playlist
	for _, f := range playlistFiles {
		pl, err := parsePlaylistVideos(f, library, playlistNames)
		if err != nil {
			continue
		}
		if len(pl.Tracks) == 0 {
			continue
		}
		playlists = append(playlists, *pl)
	}

	if len(playlists) == 0 {
		return nil, fmt.Errorf("no playlists with tracks found in the ZIP — make sure to export both 'YouTube Music' playlists and music library")
	}

	return playlists, nil
}

// parseLibraryCSV reads the music library songs CSV.
// Format: Video ID, Song Title, Album Title, Artist Name 1
func parseLibraryCSV(f *zip.File) (map[string]libraryEntry, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	reader := csv.NewReader(rc)
	reader.LazyQuotes = true

	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	// Find column indices by header name (handles localization)
	colMap := make(map[string]int)
	for i, h := range header {
		colMap[strings.TrimSpace(strings.ToLower(h))] = i
	}

	videoIDCol := findCol(colMap, "video id")
	titleCol := findCol(colMap, "song title", "title")
	albumCol := findCol(colMap, "album title", "album")
	artistCol := findCol(colMap, "artist name 1", "artist")

	// Fallback: if headers not recognized, assume positional
	if videoIDCol < 0 {
		videoIDCol = 0
	}
	if titleCol < 0 {
		titleCol = 1
	}
	if albumCol < 0 {
		albumCol = 2
	}
	if artistCol < 0 {
		artistCol = 3
	}

	entries := make(map[string]libraryEntry)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		maxCol := max(videoIDCol, titleCol, albumCol, artistCol)
		if len(record) <= maxCol {
			continue
		}

		id := strings.TrimSpace(record[videoIDCol])
		if id == "" {
			continue
		}

		entries[id] = libraryEntry{
			Title:  strings.TrimSpace(record[titleCol]),
			Album:  strings.TrimSpace(record[albumCol]),
			Artist: strings.TrimSpace(record[artistCol]),
		}
	}

	return entries, nil
}

// parsePlaylistIndex reads the playlists.csv index to get playlist titles.
// Format: Playlist ID, ..., Playlist Title (Original), ...
func parsePlaylistIndex(f *zip.File) (map[string]string, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	reader := csv.NewReader(rc)
	reader.LazyQuotes = true

	header, err := reader.Read()
	if err != nil {
		return nil, err
	}

	colMap := make(map[string]int)
	for i, h := range header {
		colMap[strings.TrimSpace(strings.ToLower(h))] = i
	}

	idCol := findCol(colMap, "playlist id")
	titleCol := findCol(colMap, "playlist title (original)", "playlist title")

	if idCol < 0 {
		idCol = 0
	}
	if titleCol < 0 {
		titleCol = 3
	}

	names := make(map[string]string)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		maxCol := max(idCol, titleCol)
		if len(record) <= maxCol {
			continue
		}

		id := strings.TrimSpace(record[idCol])
		title := strings.TrimSpace(record[titleCol])
		if id != "" && title != "" {
			names[id] = title
		}
	}

	return names, nil
}

// parsePlaylistVideos reads a single playlist's video CSV.
// Filename format: "<playlist name>-videos.csv"
// CSV format: Video ID, Playlist Video Creation Timestamp
func parsePlaylistVideos(f *zip.File, library map[string]libraryEntry, playlistNames map[string]string) (*source.Playlist, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	reader := csv.NewReader(rc)
	reader.LazyQuotes = true

	// Skip header
	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	// Derive playlist name from filename: "test playlist-videos.csv" -> "test playlist"
	base := filepath.Base(f.Name)
	name := strings.TrimSuffix(base, "-videos.csv")

	playlist := &source.Playlist{
		Name: name,
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		if len(record) == 0 {
			continue
		}

		videoID := strings.TrimSpace(record[0])
		if videoID == "" {
			continue
		}

		entry, found := library[videoID]
		if !found {
			continue // track not in library, can't resolve metadata
		}

		playlist.Tracks = append(playlist.Tracks, source.Track{
			TrackName:   entry.Title,
			ArtistNames: entry.Artist,
			AlbumName:   entry.Album,
		})
	}

	return playlist, nil
}

// findCol looks up a column index by trying multiple header names.
func findCol(colMap map[string]int, names ...string) int {
	for _, name := range names {
		if idx, ok := colMap[name]; ok {
			return idx
		}
	}
	return -1
}
