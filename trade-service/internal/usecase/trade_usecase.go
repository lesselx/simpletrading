package usecase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"simpletrading/tradeservice/internal/config"
	"simpletrading/tradeservice/internal/repository/memory"
)

type TradeUsecase struct {
	repo *memory.TradeRepository
	cfg  *config.Config
}

func NewTradeUsecase(repo *memory.TradeRepository, cfg *config.Config) *TradeUsecase {
	return &TradeUsecase{repo: repo, cfg: cfg}
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (uc *TradeUsecase) PlaceTrade(amount float64) error {
	// Step 1: Get machine token from Auth Service
	token, err := GetMachineToken(uc.cfg.AuthUrl, uc.cfg.ClientId, uc.cfg.ClientSecret)
	if err != nil {
		return fmt.Errorf("auth failed: %v", err)
	}

	// Step 2: Request recent data from Data Service with Authorization
	req, err := http.NewRequest("GET", uc.cfg.DataUrl, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("data request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read and log the response body for debugging
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// Log the raw response body to see what we are getting
	fmt.Printf("Raw response body: %s\n", respBody)

	// Step 3: Try to decode the response into a struct
	var data []struct {
		Value     float64 `json:"value"`
		Timestamp string  `json:"timestamp"`
	}

	// Decode the JSON response
	if err := json.Unmarshal(respBody, &data); err != nil {
		return fmt.Errorf("data decode failed: %v", err)
	}

	// Step 4: Ensure there's some data to process
	if len(data) == 0 {
		return fmt.Errorf("no data available")
	}

	// Step 5: Find the lowest price in the last 24 hours
	min := data[0].Value
	for _, d := range data {
		if d.Value < min {
			min = d.Value
		}
	}

	// Step 6: Validate the trade price
	if amount < min/2 {
		return fmt.Errorf("trade price too low; must be at least %.2f", min/2)
	}

	fmt.Printf("Trade accepted: %.2f\n", amount)

	// TODO: Step 7: Save trade in database
	// return uc.repo.Insert(&Trade{Amount: amount, Timestamp: time.Now()})

	return nil
}

// func (uc *TradeUsecase) PlaceTrade(trade domain.Trade) error {
// 	// Step 1: Get machine token
// 	token, err := GetMachineToken(uc.cfg.AuthUrl, uc.cfg.ClientId, uc.cfg.ClientSecret)
// 	if err != nil {
// 		return err
// 	}

// 	// Step 2: Call Data Service
// 	req, err := http.NewRequest("GET", uc.cfg.DataUrl, nil)
// 	if err != nil {
// 		return err
// 	}
// 	req.Header.Set("Authorization", "Bearer "+token)

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	var data []struct {
// 		Value     float64 `json:"value"`
// 		Timestamp string  `json:"timestamp"`
// 	}
// 	err = json.NewDecoder(resp.Body).Decode(&data)
// 	if err != nil {
// 		return err
// 	}

// 	// Step 3: Calculate lowest price in last 24 hours
// 	min := data[0].Value
// 	for _, d := range data {
// 		if d.Value < min {
// 			min = d.Value
// 		}
// 	}

// 	// Step 4: Validate trade price
// 	if trade.Price < min/2 {
// 		return fmt.Errorf("trade price too low; must be at least %.2f", min/2)
// 	}

// 	// Step 5: Save trade
// 	return uc.repo.Insert(&trade)
// }

func GetMachineToken(authURL, clientID, clientSecret string) (string, error) {
	req, err := http.NewRequest("GET", authURL, nil)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(clientID, clientSecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get token: %s", resp.Status)
	}

	var tokenRes TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenRes)
	if err != nil {
		return "", err
	}

	log.Printf("Received token: %s", tokenRes.AccessToken)

	return tokenRes.AccessToken, nil
}

// func (uc *TradeUsecase) PlaceTrade(userID string, price float64) error {
// 	lowest, err := uc.repo.GetLowestPriceInLast24Hours()
// 	if err == nil && price < (lowest/2.0) {
// 		return errors.New("trade price must not be lower than half the lowest price in the last 24 hours")
// 	}

// 	trade := &domain.Trade{
// 		UserID: userID,
// 		Price:  price,
// 	}

// 	err = uc.repo.Insert(trade)
// 	if err != nil {
// 		log.Println("Error saving trade:", err)
// 		return err
// 	}

// 	return nil
// }
