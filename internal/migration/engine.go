package migration

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"tuneshift/internal/source"
	"tuneshift/internal/tidal"
)

type Result struct {
	TotalTracks   int            `json:"total_tracks"`
	MatchedTracks int            `json:"matched_tracks"`
	FailedTracks  int            `json:"failed_tracks"`
	Playlists     int            `json:"playlists_created"`
	NotFound      []NotFoundItem `json:"not_found"`
}

type NotFoundItem struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
	ISRC   string `json:"isrc,omitempty"`
}

type Engine struct {
	tidalClient *tidal.Client
	matcher     *Matcher
	progress    *ProgressReporter
}

func NewEngine(tidalClient *tidal.Client, progress *ProgressReporter) *Engine {
	return &Engine{
		tidalClient: tidalClient,
		matcher:     NewMatcher(tidalClient),
		progress:    progress,
	}
}

func (e *Engine) Run(ctx context.Context, playlists []source.Playlist) (*Result, error) {
	result := &Result{}

	// Check for existing playlists on Tidal
	existingNames := make(map[string]bool)
	if names, err := e.tidalClient.GetUserPlaylists(); err == nil {
		for _, name := range names {
			existingNames[name] = true
		}
	}

	for _, pl := range playlists {
		if ctx.Err() != nil {
			return result, ctx.Err()
		}

		if existingNames[pl.Name] {
			e.progress.Send(ProgressEvent{
				Type:    "duplicate",
				Message: fmt.Sprintf("Playlist '%s' already exists on Tidal - skipping", pl.Name),
			})
			continue
		}

		e.progress.Send(ProgressEvent{
			Type:    "phase",
			Message: fmt.Sprintf("Migrating playlist '%s' (%d tracks)...", pl.Name, len(pl.Tracks)),
		})

		tidalPlaylistUUID, err := e.tidalClient.CreatePlaylist(pl.Name, "Migrated from Spotify via Tuneshift")
		if err != nil {
			log.Printf("Failed to create playlist %q: %v", pl.Name, err)
			e.progress.Send(ProgressEvent{
				Type:    "error",
				Message: fmt.Sprintf("Failed to create playlist '%s': %s", pl.Name, err.Error()),
			})
			continue
		}

		matchedIDs := e.matchTracks(ctx, pl.Tracks, result)

		if len(matchedIDs) > 0 {
			if err := e.tidalClient.AddTracksToPlaylist(tidalPlaylistUUID, matchedIDs); err != nil {
				log.Printf("Failed to add tracks to playlist %q: %v", pl.Name, err)
			}
		}

		result.Playlists++
		e.progress.Send(ProgressEvent{
			Type:    "playlist",
			Message: fmt.Sprintf("Playlist '%s': %d/%d tracks matched", pl.Name, len(matchedIDs), len(pl.Tracks)),
		})
	}

	// Send not-found tracks
	for _, nf := range result.NotFound {
		artist := nf.Artist
		if artist != "" {
			artist = " by " + artist
		}
		e.progress.Send(ProgressEvent{
			Type:    "not_found",
			Message: nf.Name + artist,
		})
	}

	e.progress.Send(ProgressEvent{
		Type:    "complete",
		Message: "Migration complete!",
	})

	return result, nil
}

func (e *Engine) matchTracks(ctx context.Context, tracks []source.Track, result *Result) []string {
	total := len(tracks)
	matchedIDs := make([]string, 0, total)

	// Phase 1: Batch ISRC lookup
	isrcs := make([]string, 0, total)
	for _, t := range tracks {
		if t.ISRC != "" {
			isrcs = append(isrcs, t.ISRC)
		}
	}

	isrcMatches := make(map[string]*tidal.Track)
	if len(isrcs) > 0 {
		e.progress.Send(ProgressEvent{
			Type:    "phase",
			Message: fmt.Sprintf("Looking up %d tracks by ISRC...", len(isrcs)),
		})
		matches, err := e.tidalClient.SearchTracksByISRC(isrcs)
		if err != nil {
			log.Printf("Batch ISRC lookup failed: %v, falling back to individual search", err)
		} else {
			isrcMatches = matches
		}
	}

	// Phase 2: Process results, fuzzy search only for misses
	fuzzyCount := 0
	for i, track := range tracks {
		if ctx.Err() != nil {
			break
		}

		result.TotalTracks++
		var matched *tidal.Track

		if track.ISRC != "" {
			matched = isrcMatches[strings.ToUpper(track.ISRC)]
		}

		if matched == nil {
			if fuzzyCount > 0 {
				time.Sleep(300 * time.Millisecond)
			}
			fuzzyCount++
			t, err := e.matcher.FuzzyMatch(track)
			if err == nil {
				matched = t
			}
		}

		if matched != nil {
			matchedIDs = append(matchedIDs, matched.ID)
			result.MatchedTracks++
		} else {
			result.FailedTracks++
			result.NotFound = append(result.NotFound, NotFoundItem{
				Name:   track.TrackName,
				Artist: track.FirstArtist(),
				Album:  track.AlbumName,
				ISRC:   track.ISRC,
			})
		}

		if (i+1)%10 == 0 || i+1 == total {
			e.progress.Send(ProgressEvent{
				Type:    "progress",
				Message: fmt.Sprintf("Tracks: %d/%d matched", result.MatchedTracks, i+1),
				Current: i + 1,
				Total:   total,
			})
		}
	}

	return matchedIDs
}
