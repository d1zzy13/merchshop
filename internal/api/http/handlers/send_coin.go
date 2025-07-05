package handlers

import (
	"encoding/json"
	"net/http"

	"merchshop/internal/api/http/middleware"
	"merchshop/internal/api/http/models"
)

// SendCoin godoc
// @Summary Отправить монеты другому пользователю
// @Tags default
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body models.SendCoinRequest true "Кому и сколько отправить"
// @Success 200 {object} models.InfoResponse "Успешный ответ"
// @Failure 400 {object} models.ErrorResponse "Неверный запрос"
// @Failure 401 {object} models.ErrorResponse "Неавторизован"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Router /sendCoin [post]
func (h *Handler) SendCoin(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Неавторизован")
		return
	}

	var req models.SendCoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Неверный запрос")
		return
	}

	// Получаем получателя
	receiver, err := h.userUseCase.GetByUsername(r.Context(), req.ToUser)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Неверный запрос")
		return
	}

	err = h.transactionUseCase.Transfer(r.Context(), userID, receiver.ID, req.Amount)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, "Успешно")
}
