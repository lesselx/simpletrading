package http

import (
	"encoding/json"
	"log"
	"net/http"
	"simpletrading/authservice/internal/usecase"
)

type Handler struct {
	uc *usecase.AuthUsecase
}

func NewHandler(uc *usecase.AuthUsecase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) Router() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/auth/google", GoogleLogin)
	mux.HandleFunc("/auth/google/callback", GoogleCallback)
	mux.HandleFunc("/auth/token", h.GetToken)
	return mux
}

// GetToken handles token issuance for machines
func (h *Handler) GetToken(w http.ResponseWriter, r *http.Request) {
	// Assume the client credentials are passed in the Authorization header (Basic Auth)
	clientID, clientSecret, ok := r.BasicAuth()

	log.Printf("Client ID: %s, Client Secret: %s", clientID, clientSecret)
	if !ok {
		log.Println("Missing client credentials")
		http.Error(w, "Missing client credentials", http.StatusUnauthorized)
		return
	}

	log.Printf("Extracted client ID: %s", clientID)
	log.Printf("Extracted client secret: %s", clientSecret)

	// Validate the client credentials
	token, err := h.uc.GetToken(clientID, clientSecret)
	if err != nil {
		log.Printf("Error getting token: %v", err)
		http.Error(w, "Invalid client credentials", http.StatusUnauthorized)
		return
	}

	// Return the JWT token as a response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"access_token": token})
}
