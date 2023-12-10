package mysql

import (
	"context"
	"io"
	"time"

	"github.com/jim-nnamdi/jinx/pkg/model"
)

//go:generate mockgen -destination=mocks/mock_database.go -package=mocks

type Database interface {
	io.Closer

	/* user interaction queries */
	CreateUser(ctx context.Context, username string, password string, email string, country string, phone string, sessionkey string, walletbalance float64) (bool, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CheckUser(ctx context.Context, email string, password string) (*model.User, error)
	GetBySessionKey(ctx context.Context, sessionkey string) (*model.User, error)
	GetUserPortfolio(ctx context.Context, user_email string) (*[]model.PortfolioOrder, error)

	/* transactions */
	GetUserTransactions(ctx context.Context, user_email string) (*[]model.Transaction, error)
	CreateNewTransaction(ctx context.Context, from_user int, from_user_email string, to_user int, to_user_email string, transactiontype string, created_at time.Time, updated_at time.Time, amount int, user_email string) (bool, error)
}
