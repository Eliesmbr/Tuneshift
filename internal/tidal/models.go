package tidal

type Track struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	ArtistNames []string `json:"artistNames"`
	Duration   int      `json:"duration"`
	ISRC       string   `json:"isrc"`
}
