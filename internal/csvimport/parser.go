package csvimport

import (
	"bufio"
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

// Exportify CSV column indices (fixed order regardless of language)
const (
	colTrackName   = 1
	colArtistNames = 3
	colAlbumName   = 5
	colDurationMS  = 12
	colISRC        = 16
	minColumns     = 17
)

// ParseCSV parses an Exportify CSV and returns the tracks.
// playlistName is provided by the caller (derived from filename).
func ParseCSV(r io.Reader, playlistName string) (*Playlist, error) {
	// Strip UTF-8 BOM if present (common when CSV is downloaded on Mac/iOS)
	br := bufio.NewReader(r)
	if bom, err := br.Peek(3); err == nil && len(bom) >= 3 && bom[0] == 0xEF && bom[1] == 0xBB && bom[2] == 0xBF {
		br.Discard(3)
	}

	reader := csv.NewReader(br)
	reader.LazyQuotes = true

	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	if len(header) < minColumns {
		return nil, fmt.Errorf("expected at least %d columns, got %d — is this an Exportify CSV?", minColumns, len(header))
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
		if len(record) < minColumns {
			continue
		}

		track := Track{
			TrackName:   record[colTrackName],
			ArtistNames: record[colArtistNames],
			AlbumName:   record[colAlbumName],
			ISRC:        record[colISRC],
		}
		track.DurationMS, _ = strconv.Atoi(record[colDurationMS])

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
