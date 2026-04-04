package youtube

type PlaylistSummary struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	TrackCount int    `json:"track_count"`
}

type videoItem struct {
	videoID      string
	title        string
	channelTitle string
}
