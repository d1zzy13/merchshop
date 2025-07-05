package transaction_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"merchshop/internal/repository/transaction"
)

// успешный перевод средств
func TestRepo_CreateTransaction_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := transaction.NewTransactionRepository(db)
	ctx := context.Background()

	const (
		senderID   = 1
		receiverID = 2
		amount     = 100
	)

	mock.ExpectBegin()

	// Мокаем списание у отправителя
	mock.ExpectExec(`UPDATE users SET balance = balance - \$1`).
		WithArgs(amount, senderID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Мокаем пополнение получателю
	mock.ExpectExec(`UPDATE users SET balance = balance \+ \$1`).
		WithArgs(amount, receiverID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Мокаем запись транзакции
	mock.ExpectExec(`INSERT INTO transactions`).
		WithArgs(senderID, receiverID, amount).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err = repo.CreateTransaction(ctx, senderID, receiverID, amount)
	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

// ошибка при недостатке средств
func TestRepo_CreateTransaction_InsufficientFunds(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := transaction.NewTransactionRepository(db)
	ctx := context.Background()

	const (
		senderID   = 1
		receiverID = 2
		amount     = 999
	)

	mock.ExpectBegin()

	// Списать средства не удалось (balance < amount)
	mock.ExpectExec(`UPDATE users SET balance = balance - \$1`).
		WithArgs(amount, senderID).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 row affected

	mock.ExpectRollback()

	err = repo.CreateTransaction(ctx, senderID, receiverID, amount)
	require.ErrorContains(t, err, "insufficient funds")

	require.NoError(t, mock.ExpectationsWereMet())
}

// получение транзакций по айди
func TestRepo_GetBySenderID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := transaction.NewTransactionRepository(db)
	ctx := context.Background()

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "sender_id", "receiver_id", "amount", "created_at", "sender_name", "receiver_name",
	}).AddRow(1, 1, 2, 100, now, "alice", "bob")

	mock.ExpectQuery(`SELECT t.id, t.sender_id, t.receiver_id, t.amount`).
		WithArgs(1).
		WillReturnRows(rows)

	txs, err := repo.GetBySenderID(ctx, 1)
	require.NoError(t, err)
	require.Len(t, txs, 1)

	tx := txs[0]
	require.Equal(t, 1, tx.ID)
	require.Equal(t, 1, tx.SenderID)
	require.Equal(t, 2, tx.ReceiverID)
	require.Equal(t, 100, tx.Amount)
	require.Equal(t, "alice", tx.SenderName)
	require.Equal(t, "bob", tx.ReceiverName)
}

// получение транзакции, которых нету у пользователя
func TestRepo_GetByReceiverID_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := transaction.NewTransactionRepository(db)
	ctx := context.Background()

	mock.ExpectQuery(`SELECT t.id, t.sender_id, t.receiver_id`).
		WithArgs(99).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "sender_id", "receiver_id", "amount", "created_at", "sender_name", "receiver_name",
		}))

	txs, err := repo.GetByReceiverID(ctx, 99)
	require.NoError(t, err)
	require.Empty(t, txs)
}
