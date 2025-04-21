package main

import (
	"log"
	"net/http"
	"simpletrading/authservice/internal/config"
	apphandler "simpletrading/authservice/internal/delivery/http"
	"simpletrading/authservice/internal/repository/memory"
	"simpletrading/authservice/internal/usecase"
)

func main() {

	// Initialize GORM SQLite database
	db, cfg := config.Init()

	// Example: Use GoogleOAuthConfig for Google login

	// Initialize the UserRepository using GORM DB instance
	userRepo := memory.NewUserRepository(db)

	// Initialize usecase with the repository
	uc := usecase.NewAuthUsecase(userRepo, *cfg)

	// Set up the HTTP handler
	handler := apphandler.NewHandler(uc)

	// Start the server
	log.Println("Auth Service running on", cfg.Port)
	http.ListenAndServe(cfg.Port, handler.Router())
}
