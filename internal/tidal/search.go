package tidal

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func (c *Client) SearchTracksByISRC(isrcs []string) (map[string]*Track, error) {
	result := make(map[string]*Track, len(isrcs))

	// Tidal allows max 20 ISRCs per request and returns one track per ISRC.
	for i := 0; i < len(isrcs); i += 20 {
		end := i + 20
		if end > len(isrcs) {
			end = len(isrcs)
		}

		params := url.Values{}
		params.Set("countryCode", c.countryCode)
		for _, isrc := range isrcs[i:end] {
			params.Add("filter[isrc]", strings.ToUpper(isrc))
		}

		u := baseURL + "/tracks?" + params.Encode()

		var resp struct {
			Data []searchItem `json:"data"`
		}
		if err := c.doRequest("GET", u, nil, &resp); err != nil {
			return nil, err
		}

		for _, item := range resp.Data {
			t := item.toTrack()
			if t.ISRC != "" {
				result[strings.ToUpper(t.ISRC)] = t
			}
		}
	}

	return result, nil
}

func (c *Client) SearchTrack(query string, limit int) ([]Track, error) {
	if limit == 0 {
		limit = 10
	}

	u := fmt.Sprintf("%s/searchResults/%s?countryCode=%s&include=tracks,tracks.artists",
		baseURL, url.PathEscape(query), c.countryCode)

	var result struct {
		Included []searchItem `json:"included"`
	}
	if err := c.doRequest("GET", u, nil, &result); err != nil {
		return nil, err
	}

	// Build artist name lookup from included artist objects
	artistNames := make(map[string]string)
	for _, item := range result.Included {
		if item.Type == "artists" || item.Type == "artist" {
			if item.Attributes.Name != "" {
				artistNames[item.ID] = item.Attributes.Name
			}
		}
	}

	var tracks []Track
	for _, item := range result.Included {
		if item.Type != "tracks" && item.Type != "track" {
			continue
		}
		t := item.toTrack()
		// Resolve artist names from relationships
		for _, rel := range item.Relationships.Artists.Data {
			if name, ok := artistNames[rel.ID]; ok {
				t.ArtistNames = append(t.ArtistNames, name)
			}
		}
		tracks = append(tracks, *t)
		if len(tracks) >= limit {
			break
		}
	}
	return tracks, nil
}

// searchItem handles the JSON:API format where duration can be string or int
type searchItem struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Title    string          `json:"title"`
		Name     string          `json:"name"` // for artist resources
		Duration json.RawMessage `json:"duration"`
		ISRC     string          `json:"isrc"`
	} `json:"attributes"`
	Relationships struct {
		Artists struct {
			Data []resourceIdentifier `json:"data"`
		} `json:"artists"`
	} `json:"relationships"`
}

type resourceIdentifier struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

func (s *searchItem) toTrack() *Track {
	dur := parseDuration(s.Attributes.Duration)
	return &Track{
		ID:       s.ID,
		Title:    s.Attributes.Title,
		Duration: dur,
		ISRC:     s.Attributes.ISRC,
	}
}

// parseDuration parses a duration value that can be an int (seconds),
// a numeric string, or an ISO 8601 duration like "PT2M58S".
// Returns seconds.
func parseDuration(raw json.RawMessage) int {
	if len(raw) == 0 {
		return 0
	}
	var i int
	if err := json.Unmarshal(raw, &i); err == nil {
		return i
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		if v, err := strconv.Atoi(s); err == nil {
			return v
		}
		return parseISO8601Duration(s)
	}
	return 0
}

func parseISO8601Duration(s string) int {
	s = strings.TrimPrefix(s, "PT")
	total := 0
	num := 0
	for _, c := range s {
		switch {
		case c >= '0' && c <= '9':
			num = num*10 + int(c-'0')
		case c == 'H':
			total += num * 3600
			num = 0
		case c == 'M':
			total += num * 60
			num = 0
		case c == 'S':
			total += num
			num = 0
		}
	}
	return total
}
