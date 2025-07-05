package user

import (
	"context"
	"database/sql"
	"fmt"

	entities "merchshop/internal/entity"
)

type Repository interface {
	CreateUser(ctx context.Context, username string, password string) (*entities.User, error)
	GetByID(ctx context.Context, id int) (*entities.User, error)
	GetByUsername(ctx context.Context, username string) (*entities.User, error)
}

type Repo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) Repository {
	return &Repo{db: db}
}

func (r *Repo) CreateUser(ctx context.Context, username string, password string) (*entities.User, error) {
	const query = `
        INSERT INTO users (username, password_hash, balance)
        VALUES ($1, $2, 1000)
        RETURNING id, username, password_hash, balance, created_at`

	var user entities.User

	err := r.db.QueryRowContext(ctx, query, username, password).
		Scan(&user.ID, &user.Username, &user.Password, &user.Balance, &user.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func (r *Repo) GetByID(ctx context.Context, id int) (*entities.User, error) {
	const query = `
        SELECT id, username, password_hash, balance, created_at
        FROM users
        WHERE id = $1`

	var user entities.User
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&user.ID, &user.Username, &user.Password, &user.Balance, &user.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

func (r *Repo) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	const query = `
        SELECT id, username, password_hash, balance, created_at
        FROM users
        WHERE username = $1`

	var user entities.User

	err := r.db.QueryRowContext(ctx, query, username).
		Scan(&user.ID, &user.Username, &user.Password, &user.Balance, &user.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return &user, nil
}
