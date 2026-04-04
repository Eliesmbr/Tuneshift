package youtube

import (
	"fmt"
	"log"

	"tuneshift/internal/source"
)

// FetchPlaylists fetches the selected playlists from YouTube and converts them to source.Playlist.
func FetchPlaylists(client *Client, selections []PlaylistSummary) ([]source.Playlist, error) {
	var playlists []source.Playlist

	for _, sel := range selections {
		log.Printf("Fetching YouTube playlist: %s (%s)", sel.Name, sel.ID)

		items, err := client.ListPlaylistItems(sel.ID)
		if err != nil {
			return nil, fmt.Errorf("fetch playlist %q: %w", sel.Name, err)
		}

		if len(items) == 0 {
			playlists = append(playlists, source.Playlist{Name: sel.Name})
			continue
		}

		// Collect video IDs for duration lookup
		videoIDs := make([]string, len(items))
		for i, item := range items {
			videoIDs[i] = item.videoID
		}

		durations, err := client.GetVideoDurations(videoIDs)
		if err != nil {
			log.Printf("Warning: failed to fetch durations for playlist %q: %v", sel.Name, err)
			// Continue without durations rather than failing entirely
		}

		tracks := make([]source.Track, 0, len(items))
		for _, item := range items {
			trackName, artistName := ParseVideoTitle(item.title, item.channelTitle)

			track := source.Track{
				TrackName:   trackName,
				ArtistNames: artistName,
			}
			if durations != nil {
				track.DurationMS = durations[item.videoID]
			}

			tracks = append(tracks, track)
		}

		playlists = append(playlists, source.Playlist{
			Name:   sel.Name,
			Tracks: tracks,
		})
	}

	return playlists, nil
}
