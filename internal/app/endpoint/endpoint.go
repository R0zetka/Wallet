package endpoint

import (
	"Walet/internal/app/service"
	"Walet/internal/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
	"time"
)

func CreateWallet(w http.ResponseWriter, r *http.Request) {

	_, err := service.DB.Exec("INSERT INTO wallet (balance) VALUES (100)")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusNotFound)
		return
	}
}

func SendMoney(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	walletID := uuid.FromStringOrNil(params["walletId"])
	var transaction model.Transaction
	_ = json.NewDecoder(r.Body).Decode(&transaction)
	fmt.Println("!!!")
	walletFor, statusFor := service.SerchWallet(walletID)
	if statusFor == http.StatusOK {
		transaction.From = walletFor.ID
		walletTo, statusTO := service.SerchWallet(uuid.FromStringOrNil(transaction.To))
		if statusTO == http.StatusOK {
			var balansFor = walletFor.Balance
			var balansTo = walletTo.Balance
			if balansFor >= transaction.Amount {
				balansFor = balansFor - transaction.Amount
				balansTo = balansTo + transaction.Amount
				service.DB.QueryRow("UPDATE wallet SET balance = $1 WHERE id = $2", balansFor, transaction.From)
				service.DB.QueryRow("UPDATE wallet SET balance = $1 WHERE id = $2", balansTo, transaction.To)
				transaction.Time = time.Now()
				_, err := service.DB.Query("INSERT INTO transferhistory (timetransaction, fromwallet, towallet, amount) VALUES ($1, $2, $3, $4)", transaction.Time, uuid.FromStringOrNil(transaction.From), uuid.FromStringOrNil(transaction.To), transaction.Amount)
				if err != nil {
					log.Fatal(err)
				}
				w.WriteHeader(http.StatusOK)
				return
			} else {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		} else if statusTO == http.StatusNotFound || transaction.To == transaction.From {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	} else {
		w.WriteHeader(statusFor)
		return
	}

}

func GetWallet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	walletID := uuid.FromStringOrNil(params["walletId"])
	wallet, status := service.SerchWallet(walletID)
	if status == http.StatusOK {
		json.NewEncoder(w).Encode(wallet)
	}
	w.WriteHeader(status)

}
func TransactionHistory(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	walletID := uuid.FromStringOrNil(params["walletId"])
	wallet, status := service.SerchWallet(walletID)
	if status == http.StatusOK {
		var trans model.Transaction
		rows, _ := service.DB.Query("SELECT timetransaction, id,fromwallet,towallet,amount FROM transferhistory ")
		transactions := make([]model.Transaction, 0)

		for rows.Next() {
			rows.Scan(&trans.Time, &trans.ID, &trans.From, &trans.To, &trans.Amount)
			if trans.From == wallet.ID || trans.To == wallet.ID {
				transactions = append(transactions, trans)
			}

		}
		if len(transactions) != 0 {
			json.NewEncoder(w).Encode(transactions)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}

	} else {
		w.WriteHeader(http.StatusNotFound)
	}

}
