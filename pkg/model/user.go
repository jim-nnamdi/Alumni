package model

import (
	"context"
	"time"
)

type User struct {
	Id            int       `json:"id"`
	Username      string    `json:"username"`
	Password      string    `json:"password"`
	Email         string    `json:"email"`
	Country       string    `json:"country"`
	Phone         string    `json:"phone"`
	SessionKey    string    `json:"session_key"`
	WalletBalance string    `json:"wallet_balance"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// key is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type key struct{}

// playerInfo is the key for model.PlayerInfo values in Contexts. It is
// unexported; clients use model.NewContext and model.FromContext
// instead of using this key directly.
var user key

// NewContext returns a new Context that carries value playerInfo.
func NewContext(ctx context.Context, pi *User) context.Context {
	return context.WithValue(ctx, user, pi)
}

// FromContext returns the User value stored in ctx, if any.
func FromContext(ctx context.Context) (*User, bool) {
	pi, ok := ctx.Value(user).(*User)
	return pi, ok
}
