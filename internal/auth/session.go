package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

const (
	spotifyCookieName = "tuneshift_spotify"
	tidalCookieName   = "tuneshift_tidal"
	stateCookieName   = "tuneshift_state"
)

type SessionManager struct {
	gcm cipher.AEAD
}

type TokenData struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	UserID       string    `json:"user_id,omitempty"`
	UserName     string    `json:"user_name,omitempty"`
	UserImage    string    `json:"user_image,omitempty"`
	CountryCode  string    `json:"country_code,omitempty"`
}

func NewSessionManager(key []byte) *SessionManager {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic("invalid encryption key: " + err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic("failed to create GCM: " + err.Error())
	}
	return &SessionManager{gcm: gcm}
}

func (sm *SessionManager) encrypt(data []byte) (string, error) {
	nonce := make([]byte, sm.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := sm.gcm.Seal(nonce, nonce, data, nil)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (sm *SessionManager) decrypt(encoded string) ([]byte, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	nonceSize := sm.gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return sm.gcm.Open(nil, nonce, ciphertext, nil)
}

func (sm *SessionManager) SetTokenCookie(w http.ResponseWriter, name string, token *TokenData, secure bool) error {
	data, err := json.Marshal(token)
	if err != nil {
		return err
	}
	encrypted, err := sm.encrypt(data)
	if err != nil {
		return err
	}

	sameSite := http.SameSiteLaxMode
	if secure {
		sameSite = http.SameSiteStrictMode
	}

	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    encrypted,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		MaxAge:   3600, // 1 hour
	})
	return nil
}

func (sm *SessionManager) GetTokenCookie(r *http.Request, name string) (*TokenData, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return nil, err
	}
	data, err := sm.decrypt(cookie.Value)
	if err != nil {
		return nil, err
	}
	var token TokenData
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, err
	}
	return &token, nil
}

func (sm *SessionManager) ClearCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
}

func (sm *SessionManager) SetStateCookie(w http.ResponseWriter, state string, codeVerifier string, secure bool) error {
	data, err := json.Marshal(map[string]string{
		"state":         state,
		"code_verifier": codeVerifier,
	})
	if err != nil {
		return err
	}
	encrypted, err := sm.encrypt(data)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     stateCookieName,
		Value:    encrypted,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   600, // 10 minutes
	})
	return nil
}

func (sm *SessionManager) GetStateCookie(r *http.Request) (state string, codeVerifier string, err error) {
	cookie, err := r.Cookie(stateCookieName)
	if err != nil {
		return "", "", err
	}
	data, err := sm.decrypt(cookie.Value)
	if err != nil {
		return "", "", err
	}
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		return "", "", err
	}
	return m["state"], m["code_verifier"], nil
}
