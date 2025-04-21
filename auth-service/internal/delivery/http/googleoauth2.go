package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"simpletrading/authservice/internal/config"

	"golang.org/x/oauth2"
)

func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := config.GoogleOAuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	code := r.URL.Query().Get("code")

	token, err := config.GoogleOAuthConfig.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Code exchange failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := config.GoogleOAuthConfig.Client(ctx, token)
	userInfoResp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer userInfoResp.Body.Close()

	body, err := io.ReadAll(userInfoResp.Body)
	if err != nil {
		http.Error(w, "Read error", http.StatusInternalServerError)
		return
	}

	var userInfo map[string]interface{}
	json.Unmarshal(body, &userInfo)

	// Example: get email from userInfo
	email, ok := userInfo["email"].(string)
	if !ok {
		http.Error(w, "Email not found in Google response", http.StatusInternalServerError)
		return
	}

	jwtToken, err := config.GenerateJWT(email)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Respond with JWT
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": jwtToken,
	})
}
