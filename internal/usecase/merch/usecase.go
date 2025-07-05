package merch

import (
	"context"
	"fmt"

	entities "merchshop/internal/entity"
	"merchshop/internal/repository/merch"
)

type UseCase interface {
	List(ctx context.Context) ([]entities.Merchandise, error)
	GetByName(ctx context.Context, name string) (*entities.Merchandise, error)
}

type useCase struct {
	merchRepo merch.Repository
}

func NewUseCase(mr merch.Repository) UseCase {
	return &useCase{
		merchRepo: mr,
	}
}

func (u *useCase) GetByName(ctx context.Context, name string) (*entities.Merchandise, error) {
	if name == "" {
		return nil, fmt.Errorf("empty merchandise name")
	}

	merch, err := u.merchRepo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get merchandise by name %s: %w", name, err)
	}

	return merch, nil
}

func (u *useCase) List(ctx context.Context) ([]entities.Merchandise, error) {
	merch, err := u.merchRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get merchandise list: %w", err)
	}

	return merch, nil
}
