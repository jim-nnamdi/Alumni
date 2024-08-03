package model

import "time"

type Chat struct {
	Id        int       `json:"int"`
	Sender    User      `json:"sender"`
	Recipient User      `json:"recipient"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
