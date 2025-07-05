package handlers

import (
	"net/http"

	"merchshop/internal/api/http/middleware"
	_ "merchshop/internal/api/http/models"

	"github.com/gorilla/mux"
)

// Buy godoc
// @Summary Купить предмет из магазина
// @Tags default
// @Security BearerAuth
// @Produce json
// @Param item path string true "Название предмета"
// @Success 200 {object} models.InfoResponse "Успешный ответ"
// @Failure 400 {object} models.ErrorResponse "Неверный запрос"
// @Failure 401 {object} models.ErrorResponse "Неавторизован"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Router /buy/{item} [get]
func (h *Handler) Buy(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Неавторизован")
		return
	}

	vars := mux.Vars(r)
	merchName := vars["item"]

	merch, err := h.merchUseCase.GetByName(r.Context(), merchName)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Неверный запрос")
		return
	}

	err = h.purchaseUseCase.Purchase(r.Context(), userID, 1, merch.Name)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, "Успешно")
}
