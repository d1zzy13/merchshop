package transaction_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"merchshop/internal/entity"
	"merchshop/internal/usecase/transaction"
)

type mockRepos struct {
	GetByIDFunc           func(ctx context.Context, id int) (*entity.User, error)
	CreateTransactionFunc func(ctx context.Context, senderID, receiverID, amount int) error
	GetByUserIDFunc       func(ctx context.Context, userID int) ([]entity.Transaction, error)
	GetBySenderIDFunc     func(ctx context.Context, userID int) ([]entity.Transaction, error)
	GetByReceiverIDFunc   func(ctx context.Context, userID int) ([]entity.Transaction, error)
}

func (m *mockRepos) GetByID(ctx context.Context, id int) (*entity.User, error) {
	return m.GetByIDFunc(ctx, id)
}

func (m *mockRepos) CreateTransaction(ctx context.Context, senderID, receiverID, amount int) error {
	return m.CreateTransactionFunc(ctx, senderID, receiverID, amount)
}

func (m *mockRepos) GetByUserID(ctx context.Context, userID int) ([]entity.Transaction, error) {
	return m.GetByUserIDFunc(ctx, userID)
}

func (m *mockRepos) GetBySenderID(ctx context.Context, userID int) ([]entity.Transaction, error) {
	return m.GetBySenderIDFunc(ctx, userID)
}

func (m *mockRepos) GetByReceiverID(ctx context.Context, userID int) ([]entity.Transaction, error) {
	return m.GetByReceiverIDFunc(ctx, userID)
}

func (m *mockRepos) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	return nil, nil
}

func (m *mockRepos) CreateUser(ctx context.Context, username string, password string) (*entity.User, error) {
	return nil, nil
}
func TestTransfer_Success(t *testing.T) {
	mock := &mockRepos{
		GetByIDFunc: func(ctx context.Context, id int) (*entity.User, error) {
			return &entity.User{ID: id, Balance: 1000}, nil
		},
		CreateTransactionFunc: func(ctx context.Context, senderID, receiverID, amount int) error {
			return nil
		},
	}

	uc := transaction.NewUseCase(mock, mock)
	err := uc.Transfer(context.Background(), 1, 2, 500)

	assert.NoError(t, err)
}

func TestTransfer_InsufficientFunds(t *testing.T) {
	mock := &mockRepos{
		GetByIDFunc: func(ctx context.Context, id int) (*entity.User, error) {
			return &entity.User{ID: id, Balance: 100}, nil
		},
	}

	uc := transaction.NewUseCase(mock, mock)
	err := uc.Transfer(context.Background(), 1, 2, 200)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient funds")
}

func TestTransfer_InvalidAmount(t *testing.T) {
	mock := &mockRepos{
		GetByIDFunc: func(ctx context.Context, id int) (*entity.User, error) {
			return &entity.User{ID: id, Balance: 1000}, nil
		},
	}

	uc := transaction.NewUseCase(mock, mock)
	err := uc.Transfer(context.Background(), 1, 2, 0)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid amount")
}

func TestTransfer_SameUser(t *testing.T) {
	mock := &mockRepos{
		GetByIDFunc: func(ctx context.Context, id int) (*entity.User, error) {
			return &entity.User{ID: id, Balance: 1000}, nil
		},
	}

	uc := transaction.NewUseCase(mock, mock)
	err := uc.Transfer(context.Background(), 1, 1, 100)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sender and receiver are the same user")
}

func TestTransfer_UserNotFound(t *testing.T) {
	callCount := 0
	mock := &mockRepos{
		GetByIDFunc: func(ctx context.Context, id int) (*entity.User, error) {
			callCount++
			if callCount == 2 {
				return nil, errors.New("not found")
			}
			return &entity.User{ID: id, Balance: 1000}, nil
		},
	}

	uc := transaction.NewUseCase(mock, mock)
	err := uc.Transfer(context.Background(), 1, 2, 100)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get receiver")
}

func TestGetUserTransactions_Success(t *testing.T) {
	now := time.Now()
	mock := &mockRepos{
		GetByUserIDFunc: func(ctx context.Context, userID int) ([]entity.Transaction, error) {
			return []entity.Transaction{
				{ID: 1, SenderID: userID, ReceiverID: 2, Amount: 100, CreatedAt: now},
			}, nil
		},
	}

	uc := transaction.NewUseCase(mock, mock)
	txns, err := uc.GetUserTransactions(context.Background(), 1)

	assert.NoError(t, err)
	assert.Len(t, txns, 1)
	assert.Equal(t, 100, txns[0].Amount)
}

func TestGetSentTransactions_Success(t *testing.T) {
	mock := &mockRepos{
		GetBySenderIDFunc: func(ctx context.Context, userID int) ([]entity.Transaction, error) {
			return []entity.Transaction{{SenderID: userID}}, nil
		},
	}

	uc := transaction.NewUseCase(mock, mock)
	txns, err := uc.GetSentTransactions(context.Background(), 1)

	assert.NoError(t, err)
	assert.Len(t, txns, 1)
}

func TestGetReceivedTransactions_Success(t *testing.T) {
	mock := &mockRepos{
		GetByReceiverIDFunc: func(ctx context.Context, userID int) ([]entity.Transaction, error) {
			return []entity.Transaction{{ReceiverID: userID}}, nil
		},
	}

	uc := transaction.NewUseCase(mock, mock)
	txns, err := uc.GetReceivedTransactions(context.Background(), 1)

	assert.NoError(t, err)
	assert.Len(t, txns, 1)
}
