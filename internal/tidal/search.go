package tidal

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func (c *Client) SearchTrackByISRC(isrc string) (*Track, error) {
	u := fmt.Sprintf("%s/tracks?countryCode=%s&filter[isrc]=%s",
		baseURL, c.countryCode, url.QueryEscape(strings.ToUpper(isrc)))

	var result struct {
		Data []searchItem `json:"data"`
	}
	if err := c.doRequest("GET", u, nil, &result); err != nil {
		return nil, err
	}
	if len(result.Data) == 0 {
		return nil, nil
	}

	return result.Data[0].toTrack(), nil
}

func (c *Client) SearchTrack(query string, limit int) ([]Track, error) {
	if limit == 0 {
		limit = 10
	}

	u := fmt.Sprintf("%s/searchResults/%s?countryCode=%s&include=tracks",
		baseURL, url.PathEscape(query), c.countryCode)

	var result struct {
		Included []searchItem `json:"included"`
	}
	if err := c.doRequest("GET", u, nil, &result); err != nil {
		return nil, err
	}

	var tracks []Track
	for _, item := range result.Included {
		if item.Type != "tracks" && item.Type != "track" {
			continue
		}
		tracks = append(tracks, *item.toTrack())
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
		Duration json.RawMessage `json:"duration"`
		ISRC     string          `json:"isrc"`
	} `json:"attributes"`
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

func parseDuration(raw json.RawMessage) int {
	if len(raw) == 0 {
		return 0
	}
	// Try as int first
	var i int
	if err := json.Unmarshal(raw, &i); err == nil {
		return i
	}
	// Try as string
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		if v, err := strconv.Atoi(s); err == nil {
			return v
		}
	}
	return 0
}
