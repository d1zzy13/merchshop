package repository

import (
	"database/sql"

	"merchshop/internal/repository/merch"
	"merchshop/internal/repository/purchase"
	"merchshop/internal/repository/transaction"
	"merchshop/internal/repository/user"
)

type Repositories struct {
	User        user.Repository
	Transaction transaction.Repository
	Purchase    purchase.Repository
	Merch       merch.Repository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		User:        user.NewUserRepository(db),
		Transaction: transaction.NewTransactionRepository(db),
		Purchase:    purchase.NewPurchaseRepository(db),
		Merch:       merch.NewMerchRepository(db),
	}
}
