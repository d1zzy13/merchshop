package handlers

import (
	"net/http"

	"merchshop/internal/api/http/middleware"
	"merchshop/internal/api/http/models"
)

// Info godoc
// @Summary Получить информацию о пользователе
// @Tags default
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.InfoResponse "Успешный ответ"
// @Failure 400 {object} models.ErrorResponse "Неверный запрос"
// @Failure 401 {object} models.ErrorResponse "Неавторизован"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Router /info [get]
func (h *Handler) Info(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Неавторизован")
		return
	}

	user, err := h.userUseCase.GetByID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		return
	}

	purchases, err := h.purchaseUseCase.GetUserPurchases(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		return
	}

	sentTx, err := h.transactionUseCase.GetSentTransactions(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		return
	}

	receivedTx, err := h.transactionUseCase.GetReceivedTransactions(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		return
	}

	resp := models.InfoResponse{
		Coins:     user.Balance,
		Inventory: mapInventory(purchases),
		CoinHistory: models.CoinHistoryInfo{
			Sent:     mapTransactions(sentTx, false),
			Received: mapTransactions(receivedTx, true),
		},
	}

	writeJSON(w, http.StatusOK, resp)
}
