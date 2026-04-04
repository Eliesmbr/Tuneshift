package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	googleAuthURL  = "https://accounts.google.com/o/oauth2/v2/auth"
	googleTokenURL = "https://oauth2.googleapis.com/token"
	googleScope    = "https://www.googleapis.com/auth/youtube.readonly"
)

type GoogleAuth struct {
	clientID     string
	clientSecret string
	redirectURI  string
}

func NewGoogleAuth(clientID, clientSecret, redirectURI string) *GoogleAuth {
	return &GoogleAuth{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
	}
}

func (ga *GoogleAuth) AuthURL(state string) string {
	params := url.Values{
		"response_type": {"code"},
		"client_id":     {ga.clientID},
		"redirect_uri":  {ga.redirectURI},
		"state":         {state},
		"scope":         {googleScope},
		"access_type":   {"offline"},
		"prompt":        {"consent"},
	}
	return googleAuthURL + "?" + params.Encode()
}

type googleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func (ga *GoogleAuth) ExchangeCode(code string) (*TokenData, error) {
	data := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {ga.redirectURI},
		"client_id":     {ga.clientID},
		"client_secret": {ga.clientSecret},
	}

	req, err := http.NewRequest("POST", googleTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google token exchange failed: %s", resp.Status)
	}

	var tokenResp googleTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	token := &TokenData{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
	}

	if name, err := ga.fetchChannelName(token.AccessToken); err == nil {
		token.UserName = name
	}

	return token, nil
}

func (ga *GoogleAuth) fetchChannelName(accessToken string) (string, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/youtube/v3/channels?part=snippet&mine=true", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch channel: %s", resp.Status)
	}

	var result struct {
		Items []struct {
			Snippet struct {
				Title string `json:"title"`
			} `json:"snippet"`
		} `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if len(result.Items) > 0 {
		return result.Items[0].Snippet.Title, nil
	}
	return "", nil
}
