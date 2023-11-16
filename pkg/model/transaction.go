package model

import "time"

type Transaction struct {
	Id              int       `json:"id"`
	FromUserID      int       `json:"from_user_id"`
	FromUserEmail   string    `json:"from_user_email"`
	ToUserID        int       `json:"to_user_id"`
	ToUserEmail     string    `json:"to_user_email"`
	TransactionType string    `json:"transaction_type"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Amount          int       `json:"amount"`
	UserEmail       string    `json:"user_email"`
}
