package merch

import (
	"context"
	"database/sql"
	"fmt"

	entities "merchshop/internal/entity"
)

type Repository interface {
	GetByName(ctx context.Context, name string) (*entities.Merchandise, error)
	List(ctx context.Context) ([]entities.Merchandise, error)
}

type Repo struct {
	db *sql.DB
}

func NewMerchRepository(db *sql.DB) Repository {
	return &Repo{db: db}
}

func (r *Repo) GetByName(ctx context.Context, name string) (*entities.Merchandise, error) {
	const query = `
       SELECT name, price
       FROM merchandise
       WHERE name = $1`

	var merch entities.Merchandise
	err := r.db.QueryRowContext(ctx, query, name).
		Scan(&merch.Name, &merch.Price)

	if err != nil {
		return nil, fmt.Errorf("failed to get merchandise by name: %w", err)
	}

	return &merch, nil
}

func (r *Repo) List(ctx context.Context) ([]entities.Merchandise, error) {
	const query = `
		SELECT *
		FROM merchandise`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query merchandise: %w", err)
	}
	defer rows.Close()

	var merchItems []entities.Merchandise

	for rows.Next() {
		var merch entities.Merchandise

		if err := rows.Scan(&merch.Name, &merch.Price); err != nil {
			return nil, fmt.Errorf("failed to scan merchandise: %w", err)
		}

		merchItems = append(merchItems, merch)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning merchandise: %w", err)
	}

	return merchItems, err

}
