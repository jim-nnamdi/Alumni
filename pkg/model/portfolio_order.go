package model

import "time"

type PortfolioOrder struct {
	Id        int       `json:"id"`         // ID of portfolio
	Type      string    `json:"type"`       // buy or sell
	Security  string    `json:"security"`   // symbol or stock bought
	Unit      int       `json:"unit"`       // units of type
	Status    string    `json:"status"`     // status of order
	Cancelled int       `json:"cancelled"`  // was order cancelled
	UserID    int       `json:"user_id"`    // userID of the order
	UserEmail string    `json:"user_email"` // userEmail of the order
	CreatedAt time.Time `json:"created_at"` // time order was created
	UpdatedAt time.Time `json:"updated_at"` // time order was updated
}

type PortfolioOrders struct {
	Orders []PortfolioOrder `json:"portfolio_orders"`
}
