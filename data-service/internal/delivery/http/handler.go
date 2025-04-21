package http

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"simpletrading/dataservice/internal/usecase"
)

type Handler struct {
	uc *usecase.DataUsecase
}

func NewHandler(uc *usecase.DataUsecase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) Router() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/data", JWTMiddleware(http.HandlerFunc(h.GetData)))
	mux.Handle("/data/lowest", JWTMiddleware(http.HandlerFunc(h.GetLowestPrice)))
	return mux
}

// GetData handles the GET /data endpoint
func (h *Handler) GetData(w http.ResponseWriter, r *http.Request) {
	email, ok := GetUserEmailFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Println("User:", email)
	limit := 10

	// Check if limit param is provided
	if queryLimit := r.URL.Query().Get("limit"); queryLimit != "" {
		if parsedLimit, err := strconv.Atoi(queryLimit); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	dataPoints, err := h.uc.GetRecentData(limit)
	if err != nil {
		http.Error(w, "Failed to get data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dataPoints)
}

// StartDataGeneration starts a goroutine to generate random data every minute
func (h *Handler) StartDataGeneration() {
	// Seed the random number generator once at the start of the application
	rand.Seed(time.Now().UnixNano())

	go func() {
		ticker := time.NewTicker(time.Minute) // Generate data every minute
		defer ticker.Stop()                   // Ensure the ticker is stopped when the function ends

		for {
			<-ticker.C
			randomValue := rand.Float64() * 10000 // Generates a random value between 0 and 100,000,000
			err := h.uc.GenerateData(randomValue)
			if err != nil {
				log.Println("Error generating data:", err)
			}
		}
	}()
}

// GetLowestPrice returns the lowest price in the last 24 hours
func (h *Handler) GetLowestPrice(w http.ResponseWriter, r *http.Request) {
	// Filter data for the last 24 hours
	lowest, err := h.uc.GetLowestPriceInLast24Hours()
	if err != nil {
		http.Error(w, "Error fetching lowest price: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the lowest price as a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"lowest": lowest})
}
