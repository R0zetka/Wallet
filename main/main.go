package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
	"time"
)

type Wallet struct {
	ID      string `json:"id"`
	Balance int    `json:"balance"`
}

type Transaction struct {
	ID     string    `json:"id"`
	Time   time.Time `json:"timetransaction"`
	From   string    `json:"fromwallet"`
	To     string    `json:"towallet"`
	Amount int       `json:"amount"`
}

var db *sql.DB

func main() {

	db = connectBD()

	r := mux.NewRouter()

	r.HandleFunc("/api/v1/wallet", createWallet).Methods("POST")
	r.HandleFunc("/api/v1/wallet/{walletId}/send", sendMoney).Methods("POST")

	r.HandleFunc("/api/v1/wallet/{walletId}", getWallet).Methods("GET")
	r.HandleFunc("/api/v1/wallet/{walletId}/history", transactionHistory).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", r))
}

func connectBD() *sql.DB {
	connStr := "host=localhost port=8080 user=postgres password=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func createWallet(w http.ResponseWriter, r *http.Request) {

	_, err := db.Exec("INSERT INTO wallet (balance) VALUES (100)")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusNotFound)
		return
	}
}

func sendMoney(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	walletID := uuid.FromStringOrNil(params["walletId"])
	var transaction Transaction
	_ = json.NewDecoder(r.Body).Decode(&transaction)

	walletFor, statusFor := serchWallet(walletID)
	if statusFor == http.StatusOK {
		transaction.From = walletFor.ID
		walletTo, statusTO := serchWallet(uuid.FromStringOrNil(transaction.To))
		if statusTO == http.StatusOK {
			var balansFor = walletFor.Balance
			var balansTo = walletTo.Balance
			if balansFor >= transaction.Amount {
				balansFor = balansFor - transaction.Amount
				balansTo = balansTo + transaction.Amount
				db.QueryRow("UPDATE wallet SET balance = $1 WHERE id = $2", balansFor, transaction.From)
				db.QueryRow("UPDATE wallet SET balance = $1 WHERE id = $2", balansTo, transaction.To)
				transaction.Time = time.Now()
				_, err := db.Query("INSERT INTO transferhistory (timetransaction, fromwallet, towallet, amount) VALUES ($1, $2, $3, $4)", transaction.Time, uuid.FromStringOrNil(transaction.From), uuid.FromStringOrNil(transaction.To), transaction.Amount)
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

func getWallet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	walletID := uuid.FromStringOrNil(params["walletId"])
	wallet, status := serchWallet(walletID)
	if status == http.StatusOK {
		json.NewEncoder(w).Encode(wallet)
	}
	w.WriteHeader(status)

}
func transactionHistory(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	walletID := uuid.FromStringOrNil(params["walletId"])
	wallet, status := serchWallet(walletID)
	if status == http.StatusOK {
		var trans Transaction
		rows, _ := db.Query("SELECT timetransaction, id,fromwallet,towallet,amount FROM transferhistory ")
		transactions := make([]Transaction, 0)

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

func serchWallet(walletId uuid.UUID) (Wallet, int) {
	var wallet Wallet
	db.QueryRow("SELECT id, balance  FROM wallet WHERE id = $1", walletId).Scan(&wallet.ID, &wallet.Balance)
	if wallet.ID == "" {
		return wallet, http.StatusNotFound
	} else {
		return wallet, http.StatusOK
	}

}
