package tidal

type Track struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Duration int    `json:"duration"`
	ISRC     string `json:"isrc"`
}
