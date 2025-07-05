package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"merchshop/internal/api/http/handlers"
	"merchshop/internal/api/http/middleware"
	"merchshop/internal/api/http/models"
	"merchshop/internal/entity"
)

type mockUserUseCase struct{ mock.Mock }
type mockPurchaseUseCase struct{ mock.Mock }
type mockTransactionUseCase struct{ mock.Mock }

func (m *mockUserUseCase) GetByID(ctx context.Context, userID int) (*entity.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *mockPurchaseUseCase) GetUserPurchases(ctx context.Context, userID int) ([]entity.Purchase, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]entity.Purchase), args.Error(1)
}

func (m *mockTransactionUseCase) GetSentTransactions(ctx context.Context, userID int) ([]entity.Transaction, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]entity.Transaction), args.Error(1)
}

func (m *mockTransactionUseCase) GetReceivedTransactions(ctx context.Context, userID int) ([]entity.Transaction, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]entity.Transaction), args.Error(1)
}

func (m *mockUserUseCase) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *mockUserUseCase) Register(ctx context.Context, username string, password string) (*entity.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *mockTransactionUseCase) GetUserTransactions(ctx context.Context, userID int) ([]entity.Transaction, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]entity.Transaction), args.Error(1)
}

func (m *mockTransactionUseCase) Transfer(ctx context.Context, senderID, receiverID int, amount int) error {
	args := m.Called(ctx, senderID, receiverID, amount)
	return args.Error(0)
}

func (m *mockPurchaseUseCase) Purchase(ctx context.Context, userID, quantity int, merchName string) error {
	args := m.Called(ctx, userID, quantity, merchName)
	return args.Error(0)
}

func TestInfo_Success(t *testing.T) {
	userUC := new(mockUserUseCase)
	purchaseUC := new(mockPurchaseUseCase)
	txUC := new(mockTransactionUseCase)

	userID := 1
	ctx := context.WithValue(context.Background(), middleware.UserIDKey, userID)

	userUC.On("GetByID", mock.Anything, userID).Return(&entity.User{ID: userID, Username: "test", Balance: 100}, nil)
	purchaseUC.On("GetUserPurchases", mock.Anything, userID).Return([]entity.Purchase{}, nil)
	txUC.On("GetSentTransactions", mock.Anything, userID).Return([]entity.Transaction{}, nil)
	txUC.On("GetReceivedTransactions", mock.Anything, userID).Return([]entity.Transaction{}, nil)

	h := handlers.NewHandler(userUC, txUC, purchaseUC, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/info", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	h.Info(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.InfoResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, 100, resp.Coins)
}

func TestInfo_Unauthorized(t *testing.T) {
	h := handlers.NewHandler(nil, nil, nil, nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/info", nil)
	w := httptest.NewRecorder()

	h.Info(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
