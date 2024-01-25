package model

import "time"

type Transaction struct {
	ID     string    `json:"id"`
	Time   time.Time `json:"timetransaction"`
	From   string    `json:"fromwallet"`
	To     string    `json:"towallet"`
	Amount int       `json:"amount"`
}
