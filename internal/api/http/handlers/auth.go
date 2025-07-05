package handlers

import (
	"encoding/json"
	"net/http"

	"merchshop/internal/api/http/models"
	"merchshop/internal/config"
)

// Auth godoc
// @Summary Авторизация пользователя
// @Tags default
// @Accept json
// @Produce json
// @Param input body models.AuthRequest true "Данные авторизации"
// @Success 200 {object} models.InfoResponse "Успешный ответ"
// @Failure 400 {object} models.ErrorResponse "Неверный запрос"
// @Failure 401 {object} models.ErrorResponse "Неавторизован"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Router /auth [post]
func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {
	var req models.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Неверный запрос")
		return
	}

	user, err := h.userUseCase.GetByUsername(r.Context(), req.Username)
	if err != nil {
		hashedPassword, _ := config.HashPassword(req.Password)

		user, err = h.userUseCase.Register(r.Context(), req.Username, hashedPassword)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
			return
		}
	}

	if !(config.ComparePasswords(user.Password, req.Password)) {
		writeError(w, http.StatusUnauthorized, "Неавторизован")
		return
	}

	token, err := h.tokenManager.NewToken(user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		return
	}

	writeJSON(w, http.StatusOK, models.AuthResponse{Token: token})
}
