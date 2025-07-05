package usecase

import (
	"merchshop/internal/repository"
	"merchshop/internal/usecase/merch"
	"merchshop/internal/usecase/purchase"
	"merchshop/internal/usecase/transaction"
	"merchshop/internal/usecase/user"
)

type UseCases struct {
	User        user.UseCase
	Transaction transaction.UseCase
	Purchase    purchase.UseCase
	Merch       merch.UseCase
}

func NewUseCases(repos *repository.Repositories) *UseCases {
	return &UseCases{
		User:        user.NewUseCase(repos.User),
		Transaction: transaction.NewUseCase(repos.Transaction, repos.User),
		Purchase:    purchase.NewUseCase(repos.Purchase, repos.User, repos.Merch),
		Merch:       merch.NewUseCase(repos.Merch),
	}
}
