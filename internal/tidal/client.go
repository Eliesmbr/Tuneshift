package tidal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const baseURL = "https://openapi.tidal.com/v2"

type Client struct {
	accessToken string
	userID      string
	httpClient  *http.Client
	countryCode string
}

func NewClient(accessToken, userID, countryCode string) *Client {
	if countryCode == "" {
		countryCode = "US"
	}
	return &Client{
		accessToken: accessToken,
		userID:      userID,
		countryCode: countryCode,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) doRequest(method, url string, body interface{}, result interface{}) error {
	var lastErr error
	for attempt := range 3 {
		lastErr = c.doRequestOnce(method, url, body, result)
		if lastErr == nil {
			return nil
		}

		if _, ok := lastErr.(*rateLimitError); ok {
			wait := time.Duration(2<<uint(attempt)) * time.Second // 2s, 4s, 8s
			log.Printf("Tidal rate limited, retrying in %s (attempt %d/3)", wait, attempt+1)
			time.Sleep(wait)
			continue
		}

		return lastErr
	}
	return lastErr
}

func (c *Client) doRequestOnce(method, url string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Accept", "application/vnd.api+json")
	if body != nil {
		req.Header.Set("Content-Type", "application/vnd.api+json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return &rateLimitError{}
	}

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("tidal API error %d: %s", resp.StatusCode, string(respBody))
	}

	if result != nil && resp.StatusCode != http.StatusNoContent {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

type rateLimitError struct{}

func (e *rateLimitError) Error() string { return "tidal: rate limited" }
