package user

import (
	"context"
	"fmt"

	entities "merchshop/internal/entity"
	"merchshop/internal/repository/user"
)

type UseCase interface {
	Register(ctx context.Context, username string, password string) (*entities.User, error)
	GetByID(ctx context.Context, id int) (*entities.User, error)
	GetByUsername(ctx context.Context, username string) (*entities.User, error)
}

type useCase struct {
	userRepo user.Repository
}

func NewUseCase(userRepo user.Repository) UseCase {
	return &useCase{
		userRepo: userRepo,
	}
}

func (u *useCase) Register(ctx context.Context, username string, password string) (*entities.User, error) {
	user, err := u.userRepo.CreateUser(ctx, username, password)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (u *useCase) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	user, err := u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username %s: %w", username, err)
	}

	return user, nil
}

func (u *useCase) GetByID(ctx context.Context, id int) (*entities.User, error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id %d: %w", id, err)
	}

	return user, nil
}
