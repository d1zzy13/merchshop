package transaction

import (
	"context"
	"fmt"

	entities "merchshop/internal/entity"
	"merchshop/internal/repository/transaction"
	"merchshop/internal/repository/user"
)

type UseCase interface {
	Transfer(ctx context.Context, senderID, receiverID int, amount int) error
	GetUserTransactions(ctx context.Context, userID int) ([]entities.Transaction, error)
	GetReceivedTransactions(ctx context.Context, userID int) ([]entities.Transaction, error)
	GetSentTransactions(ctx context.Context, userID int) ([]entities.Transaction, error)
}

type useCase struct {
	transactionRepo transaction.Repository
	userRepo        user.Repository
}

func NewUseCase(transactionRepo transaction.Repository, userRepo user.Repository) UseCase {
	return &useCase{
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
	}
}

func (u *useCase) Transfer(ctx context.Context, senderID, receiverID, amount int) error {

	sender, err := u.userRepo.GetByID(ctx, senderID)
	if err != nil {
		return fmt.Errorf("failed to get sender %d: %w", senderID, err)
	}

	_, err = u.userRepo.GetByID(ctx, receiverID)
	if err != nil {
		return fmt.Errorf("failed to get receiver %d: %w", receiverID, err)
	}

	if senderID == receiverID {
		return fmt.Errorf("sender and receiver are the same user: %d", senderID)
	}

	if amount <= 0 {
		return fmt.Errorf("invalid amount: %d", amount)
	}

	if sender.Balance < amount {
		return fmt.Errorf("insufficient funds: have %d, need %d", sender.Balance, amount)
	}

	if err := u.transactionRepo.CreateTransaction(ctx, senderID, receiverID, amount); err != nil {
		return fmt.Errorf("failed to transfer money: %w", err)
	}

	return nil
}

func (u *useCase) GetUserTransactions(ctx context.Context, userID int) ([]entities.Transaction, error) {
	transactions, err := u.transactionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions for user %d: %w", userID, err)
	}

	return transactions, nil
}

func (u *useCase) GetSentTransactions(ctx context.Context, userID int) ([]entities.Transaction, error) {
	transactions, err := u.transactionRepo.GetBySenderID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sent transactions for user %d: %w", userID, err)
	}

	return transactions, nil
}

func (u *useCase) GetReceivedTransactions(ctx context.Context, userID int) ([]entities.Transaction, error) {
	transactions, err := u.transactionRepo.GetByReceiverID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get received transactions for user %d: %w", userID, err)
	}

	return transactions, nil
}
