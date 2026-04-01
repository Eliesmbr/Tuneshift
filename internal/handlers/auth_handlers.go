package handlers

import (
	"log"
	"net/http"
	"net/url"

	"tuneshift/internal/auth"
)

// Tidal Auth

func (h *Handler) TidalLogin(w http.ResponseWriter, r *http.Request) {
	state, err := auth.GenerateState()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate state")
		return
	}
	verifier, err := auth.GenerateCodeVerifier()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate verifier")
		return
	}
	challenge := auth.GenerateCodeChallenge(verifier)

	if err := h.sessions.SetStateCookie(w, state, verifier, h.isSecure()); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to set state cookie")
		return
	}
	http.Redirect(w, r, h.tidalAuth.AuthURL(state, challenge), http.StatusTemporaryRedirect)
}

func (h *Handler) TidalCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	errParam := r.URL.Query().Get("error")

	if errParam != "" {
		http.Redirect(w, r, "/?error=tidal_"+url.QueryEscape(errParam), http.StatusTemporaryRedirect)
		return
	}

	savedState, verifier, err := h.sessions.GetStateCookie(r)
	if err != nil || savedState != state {
		writeError(w, http.StatusBadRequest, "invalid state parameter")
		return
	}

	token, err := h.tidalAuth.ExchangeCode(code, verifier)
	if err != nil {
		log.Printf("Tidal code exchange failed: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to connect to Tidal")
		return
	}

	if err := h.sessions.SetTokenCookie(w, "tuneshift_tidal", token, h.isSecure()); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save session")
		return
	}

	h.sessions.ClearCookie(w, "tuneshift_state")
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (h *Handler) TidalStatus(w http.ResponseWriter, r *http.Request) {
	token, err := h.sessions.GetTokenCookie(r, "tuneshift_tidal")
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"connected": false,
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"connected": true,
		"user": map[string]string{
			"id":   token.UserID,
			"name": token.UserName,
		},
	})
}

func (h *Handler) TidalLogout(w http.ResponseWriter, r *http.Request) {
	h.sessions.ClearCookie(w, "tuneshift_tidal")
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
