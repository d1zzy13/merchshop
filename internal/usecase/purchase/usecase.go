package purchase

import (
	"context"
	"fmt"

	entities "merchshop/internal/entity"
	"merchshop/internal/repository/merch"
	"merchshop/internal/repository/purchase"
	"merchshop/internal/repository/user"
)

type UseCase interface {
	Purchase(ctx context.Context, userID, quantity int, merchName string) error
	GetUserPurchases(ctx context.Context, userID int) ([]entities.Purchase, error)
}

type useCase struct {
	purchaseRepo purchase.Repository
	userRepo     user.Repository
	merchRepo    merch.Repository
}

func NewUseCase(purchaseRepo purchase.Repository, userRepo user.Repository, merchRepo merch.Repository) UseCase {
	return &useCase{
		purchaseRepo: purchaseRepo,
		userRepo:     userRepo,
		merchRepo:    merchRepo,
	}
}

func (u *useCase) GetUserPurchases(ctx context.Context, userID int) ([]entities.Purchase, error) {
	_, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user %d: %w", userID, err)
	}

	purchases, err := u.purchaseRepo.GetByUserId(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user purchases: %w", err)
	}

	return purchases, nil
}

func (u *useCase) Purchase(ctx context.Context, userID, quantity int, merchName string) error {

	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user %d: %w", userID, err)
	}

	merch, err := u.merchRepo.GetByName(ctx, merchName)
	if err != nil {
		return fmt.Errorf("failed to get merchandise %s: %w", merchName, err)
	}

	if quantity <= 0 {
		return fmt.Errorf("invalid quantity: %d", quantity)
	}

	totalPrice := merch.Price * quantity
	if user.Balance < totalPrice {
		return fmt.Errorf("insufficient funds: have %d, need %d", user.Balance, totalPrice)
	}

	if err := u.purchaseRepo.CreatePurchase(ctx, userID, merchName, quantity); err != nil {
		return fmt.Errorf("failed to process purchase: %w", err)
	}

	return nil
}
