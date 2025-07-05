package transaction

import (
	"context"
	"database/sql"
	"fmt"

	entities "merchshop/internal/entity"
)

type Repository interface {
	CreateTransaction(ctx context.Context, senderID, receiverID int, amount int) error
	GetByUserID(ctx context.Context, userID int) ([]entities.Transaction, error)
	GetBySenderID(ctx context.Context, senderID int) ([]entities.Transaction, error)
	GetByReceiverID(ctx context.Context, receiverID int) ([]entities.Transaction, error)
}

type Repo struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) Repository {
	return &Repo{db: db}
}

func (r *Repo) CreateTransaction(ctx context.Context, senderID, receiverID, amount int) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			fmt.Printf("rollback failed: %v\n", err)
		}
	}()

	const updateSender = `
        UPDATE users 
        SET balance = balance - $1 
        WHERE id = $2 AND balance >= $1`

	result, err := tx.ExecContext(ctx, updateSender, amount, senderID)
	if err != nil {
		return fmt.Errorf("update sender balance: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("insufficient funds")
	}

	const updateReceiver = `
        UPDATE users 
        SET balance = balance + $1 
        WHERE id = $2`

	if _, err = tx.ExecContext(ctx, updateReceiver, amount, receiverID); err != nil {
		return fmt.Errorf("update receiver balance: %w", err)
	}

	const insertTx = `
        INSERT INTO transactions (sender_id, receiver_id, amount) 
        VALUES ($1, $2, $3)`

	if _, err = tx.ExecContext(ctx, insertTx, senderID, receiverID, amount); err != nil {
		return fmt.Errorf("insert transaction: %w", err)
	}

	return tx.Commit()
}

func (r *Repo) queryTransactions(ctx context.Context, query string, args ...interface{}) ([]entities.Transaction, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []entities.Transaction

	for rows.Next() {
		var t entities.Transaction

		err := rows.Scan(
			&t.ID, &t.SenderID, &t.ReceiverID, &t.Amount, &t.CreatedAt,
			&t.SenderName, &t.ReceiverName,
		)

		if err != nil {
			return nil, fmt.Errorf("scan transaction: %w", err)
		}

		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return transactions, nil
}

func (r *Repo) GetByUserID(ctx context.Context, userID int) ([]entities.Transaction, error) {
	const query = `
        SELECT t.id, t.sender_id, t.receiver_id, t.amount, t.created_at,
               s.username as sender_name, r.username as receiver_name
        FROM transactions t
        JOIN users s ON t.sender_id = s.id
        JOIN users r ON t.receiver_id = r.id
        WHERE t.sender_id = $1 OR t.receiver_id = $1`

	return r.queryTransactions(ctx, query, userID)
}

func (r *Repo) GetBySenderID(ctx context.Context, senderID int) ([]entities.Transaction, error) {
	const query = `
        SELECT t.id, t.sender_id, t.receiver_id, t.amount, t.created_at,
               s.username as sender_name, r.username as receiver_name
        FROM transactions t
        JOIN users s ON t.sender_id = s.id
        JOIN users r ON t.receiver_id = r.id
        WHERE t.sender_id = $1`

	return r.queryTransactions(ctx, query, senderID)
}

func (r *Repo) GetByReceiverID(ctx context.Context, receiverID int) ([]entities.Transaction, error) {
	const query = `
        SELECT t.id, t.sender_id, t.receiver_id, t.amount, t.created_at,
               s.username as sender_name, r.username as receiver_name
        FROM transactions t
        JOIN users s ON t.sender_id = s.id
        JOIN users r ON t.receiver_id = r.id
        WHERE t.receiver_id = $1`

	return r.queryTransactions(ctx, query, receiverID)
}
