package purchase_test

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"merchshop/internal/repository/purchase"
)

// Тест успешного создания покупки
func TestPurchase_Create_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := purchase.NewPurchaseRepository(db)

	const (
		userID    = 1
		merchName = "hoody"
		quantity  = 2
		price     = 300
	)

	totalPrice := price * quantity

	mock.ExpectBegin()

	rowPrice := sqlmock.NewRows([]string{"price"}).AddRow(price)

	mock.ExpectQuery(`SELECT price FROM merchandise WHERE name = \$1`).
		WithArgs(merchName).
		WillReturnRows(rowPrice)

	mock.ExpectExec(`UPDATE users SET balance = balance - \$1`).
		WithArgs(totalPrice, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(`INSERT INTO purchases`).
		WithArgs(userID, merchName, quantity, totalPrice).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	ctx := context.Background()
	err = repo.CreatePurchase(ctx, userID, merchName, quantity)
	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

// Тест при недостатке средств (обновление баланса не произошло)
func TestPurchase_Create_InsufficientFunds(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := purchase.NewPurchaseRepository(db)

	const (
		userID    = 1
		merchName = "hoody"
		quantity  = 1
		price     = 300
	)

	totalPrice := price * quantity

	mock.ExpectBegin()

	rowPrice := sqlmock.NewRows([]string{"price"}).AddRow(price)

	mock.ExpectQuery(`SELECT price FROM merchandise WHERE name = \$1`).
		WithArgs(merchName).
		WillReturnRows(rowPrice)

	mock.ExpectExec(`UPDATE users SET balance = balance - \$1`).
		WithArgs(totalPrice, userID).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 row affected

	mock.ExpectRollback()

	ctx := context.Background()
	err = repo.CreatePurchase(ctx, userID, merchName, quantity)
	require.ErrorContains(t, err, "insufficient funds")

	require.NoError(t, mock.ExpectationsWereMet())
}

// Тест на ошибку при получении цены товара.
func TestRepo_Purchase_MerchNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := purchase.NewPurchaseRepository(db)

	userID := 1
	merchName := "unknown-item"
	quantity := 2

	mock.ExpectBegin()

	mock.ExpectQuery(`SELECT price FROM merchandise WHERE name = \$1`).
		WithArgs(merchName).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectRollback()

	ctx := context.Background()
	err = repo.CreatePurchase(ctx, userID, merchName, quantity)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "get merchandise price")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// Тест получения покупок по пользователю
func TestPurchase_GetByUserId(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := purchase.NewPurchaseRepository(db)

	userID := 1
	now := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`
       SELECT id, user_id, merch_name, quantity, total_price, created_at
       FROM purchases
       WHERE user_id = $1
       ORDER BY created_at DESC`)).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "merch_name", "quantity", "total_price", "created_at",
		}).AddRow(1, userID, "hoody", 2, 600, now))

	ctx := context.Background()
	purchases, err := repo.GetByUserId(ctx, userID)

	require.NoError(t, err)
	require.Len(t, purchases, 1)
	require.Equal(t, "hoody", purchases[0].MerchName)

	require.NoError(t, mock.ExpectationsWereMet())
}
