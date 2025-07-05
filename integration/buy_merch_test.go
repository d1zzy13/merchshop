package integration_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"merchshop/integration"
	"merchshop/internal/api/http/auth"
	"merchshop/internal/api/http/handlers"
	"merchshop/internal/api/http/middleware"
	"merchshop/internal/api/http/router"
	"merchshop/internal/repository"
	"merchshop/internal/usecase"

	"github.com/stretchr/testify/assert"
)

func TestBuyMerch(t *testing.T) {
	db := integration.SetupTestDB(t)
	defer db.Close()

	repo := repository.NewRepositories(db)

	useCases := usecase.NewUseCases(repo)

	tokenManager, err := auth.NewJWTManager("supersecret", 24*time.Hour)
	if err != nil {
		t.Fatalf("failed to initialize token manager: %v", err)
	}

	token, err := tokenManager.NewToken(1)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	handler := handlers.NewHandler(
		useCases.User,
		nil,
		useCases.Purchase,
		useCases.Merch,
		tokenManager,
	)

	r := router.NewRouter(handler, tokenManager)

	req := httptest.NewRequest(http.MethodGet, "/api/buy/t-shirt", http.NoBody)
	req.Header.Set("Authorization", "Bearer "+token)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var userBalance int

	err = db.QueryRow("SELECT balance FROM users WHERE id = $1", 1).Scan(&userBalance)
	assert.NoError(t, err)
	assert.Equal(t, 920, userBalance)

	var purchaseCount int

	err = db.QueryRow("SELECT COUNT(*) FROM purchases WHERE user_id = $1 AND merch_name = 't-shirt'", 1).Scan(&purchaseCount)
	assert.NoError(t, err)
	assert.Equal(t, 1, purchaseCount)
}
