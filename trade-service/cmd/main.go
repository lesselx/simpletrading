package main

import (
	"log"
	"net/http"

	"simpletrading/tradeservice/internal/config"
	apphttp "simpletrading/tradeservice/internal/delivery/http"
	"simpletrading/tradeservice/internal/repository/memory"
	"simpletrading/tradeservice/internal/usecase"
)

func main() {

	db, cfg := config.Init()

	repo := memory.NewTradeRepository(db)
	uc := usecase.NewTradeUsecase(repo, cfg)
	handler := apphttp.NewHandler(uc)

	log.Println("Trade Service running on", cfg.Port)
	http.ListenAndServe(cfg.Port, handler.Router())
}
