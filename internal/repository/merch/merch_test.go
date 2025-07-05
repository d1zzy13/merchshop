package merch_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"merchshop/internal/entity"
	"merchshop/internal/repository/merch"
)

// Тест на успешное получение товара по имени
func TestMerch_GetByName_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := merch.NewMerchRepository(db)

	expected := &entity.Merchandise{
		Name:  "t-shirt",
		Price: 80,
	}

	rows := sqlmock.NewRows([]string{"name", "price"}).
		AddRow(expected.Name, expected.Price)

	mock.ExpectQuery(`SELECT name, price FROM merchandise WHERE name = \$1`).
		WithArgs(expected.Name).
		WillReturnRows(rows)

	ctx := context.Background()
	got, err := repo.GetByName(ctx, expected.Name)

	require.NoError(t, err)
	require.Equal(t, expected, got)

	require.NoError(t, mock.ExpectationsWereMet())
}

// Тест на случай, когда товар по имени не найден (возвращается ошибка)
func TestMerch_GetByName_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := merch.NewMerchRepository(db)

	mock.ExpectQuery(`SELECT name, price FROM merchandise WHERE name = \$1`).
		WithArgs("Unknown").
		WillReturnError(sql.ErrNoRows)

	ctx := context.Background()
	result, err := repo.GetByName(ctx, "Unknown")

	require.Error(t, err)
	require.Nil(t, result)

	require.NoError(t, mock.ExpectationsWereMet())
}

// Тест на успешный список всего мерча
func TestMerch_List_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := merch.NewMerchRepository(db)

	rows := sqlmock.NewRows([]string{"name", "price"}).
		AddRow("powerbank", 200).
		AddRow("t-shirt", 80)

	mock.ExpectQuery(`SELECT \* FROM merchandise`).WillReturnRows(rows)

	ctx := context.Background()
	items, err := repo.List(ctx)

	require.NoError(t, err)
	require.Len(t, items, 2)
	require.Equal(t, "powerbank", items[0].Name)
	require.Equal(t, 80, items[1].Price)

	require.NoError(t, mock.ExpectationsWereMet())
}

// Тест, если QueryContext вернул ошибку (например, проблемы с БД)
func TestMerch_List_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := merch.NewMerchRepository(db)

	mock.ExpectQuery(`SELECT \* FROM merchandise`).
		WillReturnError(sql.ErrConnDone)

	ctx := context.Background()
	items, err := repo.List(ctx)

	require.Error(t, err)
	require.Nil(t, items)

	require.NoError(t, mock.ExpectationsWereMet())
}

// Тест, если в строке ответа что-то не так (например, типы не совпадают)
func TestMerch_List_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := merch.NewMerchRepository(db)

	rows := sqlmock.NewRows([]string{"name", "price"}).
		AddRow("wallet", "not int")

	mock.ExpectQuery(`SELECT \* FROM merchandise`).WillReturnRows(rows)

	ctx := context.Background()
	items, err := repo.List(ctx)

	require.Error(t, err)
	require.Nil(t, items)

	require.NoError(t, mock.ExpectationsWereMet())
}
