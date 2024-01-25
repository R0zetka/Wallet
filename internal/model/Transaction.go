package model

import "time"

type Transaction struct {
	ID     string    `json:"id"`
	Time   time.Time `json:"timetransaction"`
	From   string    `json:"from"`
	To     string    `json:"to"`
	Amount int       `json:"amount"`
}
