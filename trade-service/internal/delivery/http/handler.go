package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"simpletrading/tradeservice/internal/usecase"
)

type Handler struct {
	uc *usecase.TradeUsecase
}

func NewHandler(uc *usecase.TradeUsecase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) Router() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/trade", JWTMiddleware(http.HandlerFunc(h.PlaceTrade)))
	return mux
}

type tradeRequest struct {
	Price float64 `json:"price"`
}

// PlaceTrade handles placing a new trade
func (h *Handler) PlaceTrade(w http.ResponseWriter, r *http.Request) {
	amountStr := r.URL.Query().Get("amount")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	err = h.uc.PlaceTrade(amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "trade accepted"})
}

// // getLowestPriceFromDataService fetches the lowest price from Data Service API
// func (h *Handler) getLowestPriceFromDataService() (float64, error) {
// 	// Request Data Service for the lowest price in the last 24 hours
// 	resp, err := http.Get("http://data-service:8081/data/lowest?since=24h")
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer resp.Body.Close()

// 	// Parse response from Data Service
// 	var result map[string]float64
// 	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
// 		return 0, err
// 	}

// 	// Return the lowest price
// 	return result["lowest"], nil
// }
