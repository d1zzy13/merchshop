package purchase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"merchshop/internal/entity"
	"merchshop/internal/usecase/purchase"
)

type mockRepos struct {
	GetByIDFunc        func(ctx context.Context, id int) (*entity.User, error)
	GetByNameFunc      func(ctx context.Context, name string) (*entity.Merchandise, error)
	CreatePurchaseFunc func(ctx context.Context, userID int, merchName string, quantity int) error
	GetByUserIdFunc    func(ctx context.Context, userID int) ([]entity.Purchase, error)
}

func (m *mockRepos) GetByID(ctx context.Context, id int) (*entity.User, error) {
	return m.GetByIDFunc(ctx, id)
}

func (m *mockRepos) GetByName(ctx context.Context, name string) (*entity.Merchandise, error) {
	return m.GetByNameFunc(ctx, name)
}

func (m *mockRepos) CreatePurchase(ctx context.Context, userID int, merchName string, quantity int) error {
	return m.CreatePurchaseFunc(ctx, userID, merchName, quantity)
}

func (m *mockRepos) GetByUserId(ctx context.Context, userID int) ([]entity.Purchase, error) {
	return m.GetByUserIdFunc(ctx, userID)
}

func (m *mockRepos) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	return nil, nil
}

func (m *mockRepos) CreateUser(ctx context.Context, username string, password string) (*entity.User, error) {
	return nil, nil
}

func (m *mockRepos) List(ctx context.Context) ([]entity.Merchandise, error) {
	return nil, nil
}

func TestPurchase_Success(t *testing.T) {
	mock := &mockRepos{
		GetByIDFunc: func(ctx context.Context, id int) (*entity.User, error) {
			return &entity.User{ID: id, Balance: 1000}, nil
		},
		GetByNameFunc: func(ctx context.Context, name string) (*entity.Merchandise, error) {
			return &entity.Merchandise{Name: name, Price: 100}, nil
		},
		CreatePurchaseFunc: func(ctx context.Context, userID int, merchName string, quantity int) error {
			return nil
		},
	}

	useCase := purchase.NewUseCase(mock, mock, mock)
	err := useCase.Purchase(context.Background(), 1, 2, "hoody")

	assert.NoError(t, err)
}

func TestPurchase_InsufficientFunds(t *testing.T) {
	mock := &mockRepos{
		GetByIDFunc: func(ctx context.Context, id int) (*entity.User, error) {
			return &entity.User{ID: id, Balance: 100}, nil
		},
		GetByNameFunc: func(ctx context.Context, name string) (*entity.Merchandise, error) {
			return &entity.Merchandise{Name: name, Price: 100}, nil
		},
	}

	useCase := purchase.NewUseCase(mock, mock, mock)
	err := useCase.Purchase(context.Background(), 1, 2, "hoody")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient funds")
}

func TestPurchase_InvalidQuantity(t *testing.T) {
	mock := &mockRepos{
		GetByIDFunc: func(ctx context.Context, id int) (*entity.User, error) {
			return &entity.User{ID: id, Balance: 1000}, nil
		},
		GetByNameFunc: func(ctx context.Context, name string) (*entity.Merchandise, error) {
			return &entity.Merchandise{Name: name, Price: 100}, nil
		},
	}

	useCase := purchase.NewUseCase(mock, mock, mock)
	err := useCase.Purchase(context.Background(), 1, 0, "hoody")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid quantity")
}

func TestGetUserPurchases_Success(t *testing.T) {
	now := time.Now()

	mock := &mockRepos{
		GetByIDFunc: func(ctx context.Context, id int) (*entity.User, error) {
			return &entity.User{ID: id}, nil
		},
		GetByUserIdFunc: func(ctx context.Context, userID int) ([]entity.Purchase, error) {
			return []entity.Purchase{
				{
					ID:        1,
					UserID:    userID,
					MerchName: "hoodie",
					Quantity:  1,
					CreatedAt: now,
				},
			}, nil
		},
	}

	useCase := purchase.NewUseCase(mock, mock, mock)
	purchases, err := useCase.GetUserPurchases(context.Background(), 1)

	assert.NoError(t, err)
	assert.Len(t, purchases, 1)
	assert.Equal(t, "hoodie", purchases[0].MerchName)
}

func TestGetUserPurchases_UserNotFound(t *testing.T) {
	mock := &mockRepos{
		GetByIDFunc: func(ctx context.Context, id int) (*entity.User, error) {
			return nil, errors.New("not found")
		},
	}

	useCase := purchase.NewUseCase(mock, mock, mock)
	purchases, err := useCase.GetUserPurchases(context.Background(), 99)

	assert.Error(t, err)
	assert.Nil(t, purchases)
	assert.Contains(t, err.Error(), "failed to get user")
}
