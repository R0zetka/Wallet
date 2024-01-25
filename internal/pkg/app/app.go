package app

import (
	"Walet/internal/app/endpoint"
	"Walet/internal/app/service"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func New() {

	service.DB = service.ConnectBD()

	r := mux.NewRouter()

	r.HandleFunc("/api/v1/wallet", endpoint.CreateWallet).Methods("POST")
	r.HandleFunc("/api/v1/wallet/{walletId}/send", endpoint.SendMoney).Methods("POST")

	r.HandleFunc("/api/v1/wallet/{walletId}", endpoint.GetWallet).Methods("GET")
	r.HandleFunc("/api/v1/wallet/{walletId}/history", endpoint.TransactionHistory).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", r))
}
