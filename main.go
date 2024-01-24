package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Wallet struct {
	ID      string `json:"id"`
	Balance int    `json:"balance"`
}

type Transaction struct {
	Time   time.Time `json:"time"`
	From   string    `json:"from"`
	To     string    `json:"to"`
	Amount int       `json:"amount"`
}

var wallets []Wallet

var transactions []Transaction

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/wallet", createWallet).Methods("POST")
	r.HandleFunc("/api/v1/wallet/{walletId}/send", sendMoney).Methods("POST")

	r.HandleFunc("/api/v1/wallet/{walletId}", getWallet).Methods("GET")
	r.HandleFunc("/api/v1/wallet/{walletId}/history", transactionHistory).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", r))
}

func createWallet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var wallet Wallet
	wallet.ID = strconv.Itoa(rand.Intn(1000000))
	wallet.Balance = 100
	wallets = append(wallets, wallet)
	json.NewEncoder(w).Encode(wallets)
}

func sendMoney(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	parms := mux.Vars(r)
	var moneyTO Transaction
	_ = json.NewDecoder(r.Body).Decode(&moneyTO)
	walletFor, indexFor := serchWallet(parms["walletId"])
	if walletFor.ID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		walletTo, indexTo := serchWallet(moneyTO.To)
		if (walletTo.ID == walletFor.ID) || (walletTo.ID == "") {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			if moneyTO.Amount <= wallets[indexFor].Balance {
				wallets[indexFor].Balance = wallets[indexFor].Balance - moneyTO.Amount
				wallets[indexTo].Balance = wallets[indexTo].Balance + moneyTO.Amount
				moneyTO.From = wallets[indexFor].ID
				moneyTO.Time = time.Now()
				transactions = append(transactions, moneyTO)
				json.NewEncoder(w).Encode(moneyTO)
				w.WriteHeader(http.StatusOK)
				return
			} else {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	}
}

func getWallet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	parms := mux.Vars(r)
	item, _ := serchWallet(parms["walletId"])
	switch item.ID {
	case "":
		w.WriteHeader(http.StatusNotFound)
	default:
		json.NewEncoder(w).Encode(item)
	}
}
func transactionHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	parms := mux.Vars(r)
	item, _ := serchWallet(parms["walletId"])
	if item.ID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		var transactionsID []Transaction
		for _, transaction := range transactions {
			if transaction.From == item.ID || transaction.To == item.ID {
				transactionsID = append(transactionsID, transaction)
			}
		}
		json.NewEncoder(w).Encode(transactionsID)
	}
}

func serchWallet(walletId string) (Wallet, int) {
	for index, item := range wallets {
		if item.ID == walletId {
			return wallets[index], index
		}
	}
	return Wallet{}, -1
}
