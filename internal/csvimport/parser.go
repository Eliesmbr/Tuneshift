package csvimport

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Track represents a track parsed from an Exportify CSV
type Track struct {
	TrackURI    string
	TrackName   string
	ArtistNames string
	AlbumName   string
	DurationMS  int
	ISRC        string
}

// Playlist represents a parsed playlist with its tracks
type Playlist struct {
	Name   string  `json:"name"`
	Tracks []Track `json:"tracks"`
}

// column indices we care about
type columnMap struct {
	trackName   int
	artistNames int
	albumName   int
	durationMS  int
	isrc        int
}

func findColumns(header []string) (*columnMap, error) {
	cm := &columnMap{
		trackName:   -1,
		artistNames: -1,
		albumName:   -1,
		durationMS:  -1,
		isrc:        -1,
	}

	for i, col := range header {
		switch strings.TrimSpace(col) {
		case "Track Name":
			cm.trackName = i
		case "Artist Name(s)":
			cm.artistNames = i
		case "Album Name":
			cm.albumName = i
		case "Track Duration (ms)":
			cm.durationMS = i
		case "ISRC":
			cm.isrc = i
		}
	}

	if cm.trackName == -1 {
		return nil, fmt.Errorf("missing required column 'Track Name' — is this an Exportify CSV?")
	}
	if cm.artistNames == -1 {
		return nil, fmt.Errorf("missing required column 'Artist Name(s)'")
	}

	return cm, nil
}

// ParseCSV parses an Exportify CSV and returns the tracks.
// playlistName is provided by the caller (derived from filename).
func ParseCSV(r io.Reader, playlistName string) (*Playlist, error) {
	reader := csv.NewReader(r)
	reader.LazyQuotes = true

	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	cols, err := findColumns(header)
	if err != nil {
		return nil, err
	}

	playlist := &Playlist{
		Name: playlistName,
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue // skip malformed rows
		}

		track := Track{}

		if cols.trackName >= 0 && cols.trackName < len(record) {
			track.TrackName = record[cols.trackName]
		}
		if cols.artistNames >= 0 && cols.artistNames < len(record) {
			track.ArtistNames = record[cols.artistNames]
		}
		if cols.albumName >= 0 && cols.albumName < len(record) {
			track.AlbumName = record[cols.albumName]
		}
		if cols.durationMS >= 0 && cols.durationMS < len(record) {
			track.DurationMS, _ = strconv.Atoi(record[cols.durationMS])
		}
		if cols.isrc >= 0 && cols.isrc < len(record) {
			track.ISRC = record[cols.isrc]
		}

		// Skip empty rows
		if track.TrackName == "" {
			continue
		}

		playlist.Tracks = append(playlist.Tracks, track)
	}

	return playlist, nil
}

// FirstArtist returns the first artist name from the comma-separated list
func (t Track) FirstArtist() string {
	parts := strings.SplitN(t.ArtistNames, ",", 2)
	if len(parts) > 0 {
		return strings.TrimSpace(parts[0])
	}
	return ""
}
