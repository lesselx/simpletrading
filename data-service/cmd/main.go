package main

import (
	"log"
	"net/http"

	"simpletrading/dataservice/internal/config"
	apphttp "simpletrading/dataservice/internal/delivery/http"
	"simpletrading/dataservice/internal/repository/memory"
	"simpletrading/dataservice/internal/usecase"
)

func main() {

	db, cfg := config.Init()

	repo := memory.NewDataRepo(db)
	uc := usecase.NewDataUsecase(*repo)
	handler := apphttp.NewHandler(uc)

	// Start data generation inside the handler
	handler.StartDataGeneration()

	log.Println("Data Service running on", cfg.Port)
	http.ListenAndServe(cfg.Port, handler.Router())
}
