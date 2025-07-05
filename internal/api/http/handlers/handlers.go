package handlers

import (
	"merchshop/internal/api/http/auth"
	"merchshop/internal/usecase/merch"
	"merchshop/internal/usecase/purchase"
	"merchshop/internal/usecase/transaction"
	"merchshop/internal/usecase/user"
)

type Handler struct {
	userUseCase        user.UseCase
	transactionUseCase transaction.UseCase
	purchaseUseCase    purchase.UseCase
	merchUseCase       merch.UseCase
	tokenManager       auth.TokenManager
}

func NewHandler(
	userUseCase user.UseCase,
	transactionUseCase transaction.UseCase,
	purchaseUseCase purchase.UseCase,
	merchUseCase merch.UseCase,
	tm auth.TokenManager,
) *Handler {
	return &Handler{
		userUseCase:        userUseCase,
		transactionUseCase: transactionUseCase,
		purchaseUseCase:    purchaseUseCase,
		merchUseCase:       merchUseCase,
		tokenManager:       tm,
	}
}
