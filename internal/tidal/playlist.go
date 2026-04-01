package tidal

import (
	"fmt"
)

type jsonAPIPlaylistCreate struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Privacy     string `json:"privacy"`
		} `json:"attributes"`
	} `json:"data"`
}

func (c *Client) CreatePlaylist(name, description string) (string, error) {
	body := jsonAPIPlaylistCreate{}
	body.Data.Type = "playlists"
	body.Data.Attributes.Name = name
	body.Data.Attributes.Description = description
	body.Data.Attributes.Privacy = "PUBLIC"

	var result struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := c.doRequest("POST", baseURL+"/playlists", body, &result); err != nil {
		return "", fmt.Errorf("create playlist failed: %w", err)
	}
	return result.Data.ID, nil
}

func (c *Client) AddTracksToPlaylist(playlistID string, trackIDs []string) error {
	// Add tracks in batches
	batchSize := 20
	for i := 0; i < len(trackIDs); i += batchSize {
		end := i + batchSize
		if end > len(trackIDs) {
			end = len(trackIDs)
		}

		batch := trackIDs[i:end]
		items := make([]map[string]string, len(batch))
		for j, id := range batch {
			items[j] = map[string]string{
				"type": "tracks",
				"id":   id,
			}
		}

		body := map[string]interface{}{
			"data": items,
		}

		u := fmt.Sprintf("%s/playlists/%s/relationships/items", baseURL, playlistID)
		if err := c.doRequest("POST", u, body, nil); err != nil {
			return fmt.Errorf("add tracks to playlist failed: %w", err)
		}
	}
	return nil
}
