package service

import (
	"Walet/internal/model"
	"database/sql"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
)

var DB *sql.DB

func ConnectBD() *sql.DB {
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

func SerchWallet(walletId uuid.UUID) (model.Wallet, int) {
	var wallet model.Wallet
	DB.QueryRow("SELECT id, balance  FROM wallet WHERE id = $1", walletId).Scan(&wallet.ID, &wallet.Balance)
	fmt.Println("wa %s", walletId)
	fmt.Println(wallet.ID)
	if wallet.ID == "" {
		return wallet, http.StatusNotFound
	} else {
		return wallet, http.StatusOK
	}

}
