package user_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"merchshop/internal/repository/user"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

// создание пользователя
func TestRepo_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := user.NewUserRepository(db)

	username := "testuser"
	password := "securepassword"

	createdAt := time.Now()

	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(username, password).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash", "balance", "created_at"}).
			AddRow(1, username, password, 1000, createdAt))

	ctx := context.Background()
	u, err := repo.CreateUser(ctx, username, password)

	require.NoError(t, err)
	require.Equal(t, 1, u.ID)
	require.Equal(t, username, u.Username)
	require.Equal(t, password, u.Password)
	require.Equal(t, 1000, u.Balance)
	require.WithinDuration(t, createdAt, u.CreatedAt, time.Second)

	require.NoError(t, mock.ExpectationsWereMet())
}

// ошибка создания пользователя
func TestRepo_FailedCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := user.NewUserRepository(db)

	username := "testuser"
	password := "securepassword"
	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash", "balance", "created_at"}))

	ctx := context.Background()
	u, err := repo.CreateUser(ctx, username, password)
	require.Nil(t, u)
	require.Error(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

// получение пользователя по айди
func TestRepo_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := user.NewUserRepository(db)

	createdAt := time.Now()

	mock.ExpectQuery(`SELECT id, username, password_hash, balance, created_at FROM users WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash", "balance", "created_at"}).
			AddRow(1, "user1", "hashpass", 800, createdAt))

	ctx := context.Background()
	u, err := repo.GetByID(ctx, 1)

	require.NoError(t, err)
	require.Equal(t, 1, u.ID)
	require.Equal(t, "user1", u.Username)
	require.Equal(t, "hashpass", u.Password)
	require.Equal(t, 800, u.Balance)
	require.WithinDuration(t, createdAt, u.CreatedAt, time.Second)

	require.NoError(t, mock.ExpectationsWereMet())
}

// получение пользователя по имени
func TestRepo_GetByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := user.NewUserRepository(db)

	createdAt := time.Now()

	mock.ExpectQuery(`SELECT id, username, password_hash, balance, created_at FROM users WHERE username = \$1`).
		WithArgs("user1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash", "balance", "created_at"}).
			AddRow(2, "user1", "pass123", 700, createdAt))

	ctx := context.Background()
	u, err := repo.GetByUsername(ctx, "user1")

	require.NoError(t, err)
	require.Equal(t, 2, u.ID)
	require.Equal(t, "user1", u.Username)
	require.Equal(t, "pass123", u.Password)
	require.Equal(t, 700, u.Balance)
	require.WithinDuration(t, createdAt, u.CreatedAt, time.Second)

	require.NoError(t, mock.ExpectationsWereMet())
}

// пользователь не найден
func TestRepo_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := user.NewUserRepository(db)

	mock.ExpectQuery(`SELECT id, username, password_hash, balance, created_at FROM users WHERE id = \$1`).
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	ctx := context.Background()
	u, err := repo.GetByID(ctx, 999)

	require.Error(t, err)
	require.Nil(t, u)

	require.True(t, errors.Is(err, sql.ErrNoRows) || u == nil)
	require.NoError(t, mock.ExpectationsWereMet())
}
