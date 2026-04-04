package youtube

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const apiBase = "https://www.googleapis.com/youtube/v3"

type Client struct {
	accessToken string
	httpClient  *http.Client
}

func NewClient(accessToken string) *Client {
	return &Client{
		accessToken: accessToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) doGet(endpoint string, params url.Values, result interface{}) error {
	u := apiBase + "/" + endpoint + "?" + params.Encode()

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("youtube API error %d: %s", resp.StatusCode, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

func (c *Client) ListPlaylists() ([]PlaylistSummary, error) {
	var playlists []PlaylistSummary
	pageToken := ""

	for {
		params := url.Values{
			"part":       {"snippet,contentDetails"},
			"mine":       {"true"},
			"maxResults": {"50"},
		}
		if pageToken != "" {
			params.Set("pageToken", pageToken)
		}

		var resp struct {
			Items []struct {
				ID      string `json:"id"`
				Snippet struct {
					Title string `json:"title"`
				} `json:"snippet"`
				ContentDetails struct {
					ItemCount int `json:"itemCount"`
				} `json:"contentDetails"`
			} `json:"items"`
			NextPageToken string `json:"nextPageToken"`
		}

		if err := c.doGet("playlists", params, &resp); err != nil {
			return nil, fmt.Errorf("list playlists: %w", err)
		}

		for _, item := range resp.Items {
			playlists = append(playlists, PlaylistSummary{
				ID:         item.ID,
				Name:       item.Snippet.Title,
				TrackCount: item.ContentDetails.ItemCount,
			})
		}

		if resp.NextPageToken == "" {
			break
		}
		pageToken = resp.NextPageToken
	}

	// Add "Liked Music" playlist (special YouTube Music playlist)
	playlists = append(playlists, PlaylistSummary{
		ID:   "LM",
		Name: "Liked Music",
	})

	return playlists, nil
}

func (c *Client) ListPlaylistItems(playlistID string) ([]videoItem, error) {
	var items []videoItem
	pageToken := ""

	for {
		params := url.Values{
			"part":       {"snippet"},
			"playlistId": {playlistID},
			"maxResults": {"50"},
		}
		if pageToken != "" {
			params.Set("pageToken", pageToken)
		}

		var resp struct {
			Items []struct {
				Snippet struct {
					Title            string `json:"title"`
					VideoOwnerTitle  string `json:"videoOwnerChannelName"`
					ResourceID       struct {
						VideoID string `json:"videoId"`
					} `json:"resourceId"`
				} `json:"snippet"`
			} `json:"items"`
			NextPageToken string `json:"nextPageToken"`
		}

		if err := c.doGet("playlistItems", params, &resp); err != nil {
			return nil, fmt.Errorf("list playlist items: %w", err)
		}

		for _, item := range resp.Items {
			if item.Snippet.ResourceID.VideoID == "" {
				continue // skip deleted/private videos
			}
			items = append(items, videoItem{
				videoID:      item.Snippet.ResourceID.VideoID,
				title:        item.Snippet.Title,
				channelTitle: item.Snippet.VideoOwnerTitle,
			})
		}

		if resp.NextPageToken == "" {
			break
		}
		pageToken = resp.NextPageToken
	}

	return items, nil
}

func (c *Client) GetVideoDurations(videoIDs []string) (map[string]int, error) {
	durations := make(map[string]int)

	// Batch up to 50 IDs per request
	for i := 0; i < len(videoIDs); i += 50 {
		end := i + 50
		if end > len(videoIDs) {
			end = len(videoIDs)
		}
		batch := videoIDs[i:end]

		params := url.Values{
			"part": {"contentDetails"},
			"id":   {strings.Join(batch, ",")},
		}

		var resp struct {
			Items []struct {
				ID             string `json:"id"`
				ContentDetails struct {
					Duration string `json:"duration"`
				} `json:"contentDetails"`
			} `json:"items"`
		}

		if err := c.doGet("videos", params, &resp); err != nil {
			return nil, fmt.Errorf("get video durations: %w", err)
		}

		for _, item := range resp.Items {
			durations[item.ID] = parseISO8601Duration(item.ContentDetails.Duration)
		}
	}

	return durations, nil
}

// parseISO8601Duration parses "PT3M45S", "PT1H2M3S", etc. into milliseconds.
func parseISO8601Duration(s string) int {
	s = strings.TrimPrefix(s, "PT")
	var hours, minutes, seconds int

	if idx := strings.Index(s, "H"); idx >= 0 {
		fmt.Sscanf(s[:idx], "%d", &hours)
		s = s[idx+1:]
	}
	if idx := strings.Index(s, "M"); idx >= 0 {
		fmt.Sscanf(s[:idx], "%d", &minutes)
		s = s[idx+1:]
	}
	if idx := strings.Index(s, "S"); idx >= 0 {
		fmt.Sscanf(s[:idx], "%d", &seconds)
	}

	return (hours*3600 + minutes*60 + seconds) * 1000
}
