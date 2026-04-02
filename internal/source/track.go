package source

import "strings"

type Track struct {
	TrackName   string
	ArtistNames string
	AlbumName   string
	DurationMS  int
	ISRC        string
}

type Playlist struct {
	Name   string  `json:"name"`
	Tracks []Track `json:"tracks"`
}

func (t Track) FirstArtist() string {
	parts := strings.SplitN(t.ArtistNames, ",", 2)
	if len(parts) > 0 {
		return strings.TrimSpace(parts[0])
	}
	return ""
}
