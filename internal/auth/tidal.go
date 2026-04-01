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
	tidalAuthURL  = "https://login.tidal.com/authorize"
	tidalTokenURL = "https://auth.tidal.com/v1/oauth2/token"
	tidalMeURL    = "https://api.tidal.com/v1/users"
	tidalScopes = "collection.read collection.write playlists.read playlists.write user.read"
)

type TidalAuth struct {
	clientID    string
	redirectURI string
}

func NewTidalAuth(clientID, redirectURI string) *TidalAuth {
	return &TidalAuth{
		clientID:    clientID,
		redirectURI: redirectURI,
	}
}

func (ta *TidalAuth) AuthURL(state, codeChallenge string) string {
	params := url.Values{
		"response_type":         {"code"},
		"client_id":             {ta.clientID},
		"redirect_uri":          {ta.redirectURI},
		"state":                 {state},
		"code_challenge_method": {"S256"},
		"code_challenge":        {codeChallenge},
	}
	params.Set("scope", tidalScopes)
	return tidalAuthURL + "?" + params.Encode()
}

type tidalTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	User         struct {
		UserID   int    `json:"userId"`
		Username string `json:"username"`
	} `json:"user"`
}

func (ta *TidalAuth) ExchangeCode(code, codeVerifier string) (*TokenData, error) {
	data := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {ta.redirectURI},
		"client_id":     {ta.clientID},
		"code_verifier": {codeVerifier},
	}

	req, err := http.NewRequest("POST", tidalTokenURL, strings.NewReader(data.Encode()))
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
		return nil, fmt.Errorf("tidal token exchange failed: %s", resp.Status)
	}

	var tokenResp tidalTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	token := &TokenData{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
		UserID:       fmt.Sprintf("%d", tokenResp.User.UserID),
		UserName:     tokenResp.User.Username,
	}

	// Fetch country code from user profile
	if cc, err := ta.fetchCountryCode(token.AccessToken, token.UserID); err == nil {
		token.CountryCode = cc
	}

	return token, nil
}

func (ta *TidalAuth) fetchCountryCode(accessToken, userID string) (string, error) {
	u := fmt.Sprintf("https://openapi.tidal.com/v2/users/%s?countryCode=US", userID)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.api+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch user: %s", resp.Status)
	}

	var result struct {
		Data struct {
			Attributes struct {
				Country string `json:"country"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.Data.Attributes.Country == "" {
		return "US", nil
	}
	return result.Data.Attributes.Country, nil
}

func (ta *TidalAuth) RefreshAccessToken(refreshToken string) (*TokenData, error) {
	data := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
		"client_id":     {ta.clientID},
	}

	req, err := http.NewRequest("POST", tidalTokenURL, strings.NewReader(data.Encode()))
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
		return nil, fmt.Errorf("tidal token refresh failed: %s", resp.Status)
	}

	var tokenResp tidalTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	newRefresh := tokenResp.RefreshToken
	if newRefresh == "" {
		newRefresh = refreshToken
	}

	return &TokenData{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: newRefresh,
		ExpiresAt:    time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
		UserID:       fmt.Sprintf("%d", tokenResp.User.UserID),
		UserName:     tokenResp.User.Username,
	}, nil
}
