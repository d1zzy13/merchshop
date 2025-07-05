package handlers

import (
	"encoding/json"
	"net/http"

	"merchshop/internal/api/http/models"
	entities "merchshop/internal/entity"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, models.ErrorResponse{Errors: message})
}

func mapInventory(purchases []entities.Purchase) []models.InventoryItem {
	// Создаем map для группировки товаров
	inventory := make(map[string]int)

	// Группируем покупки по типу товара
	for _, purchase := range purchases {
		inventory[purchase.MerchName] += purchase.Quantity
	}

	// Преобразуем map в slice для ответа
	result := make([]models.InventoryItem, 0, len(inventory))
	for itemType, quantity := range inventory {
		result = append(result, models.InventoryItem{
			Type:     itemType,
			Quantity: quantity,
		})
	}

	return result
}

func mapTransactions(transactions []entities.Transaction, isReceived bool) []models.CoinOperation {
	result := make([]models.CoinOperation, len(transactions))

	for i, tx := range transactions {
		operation := models.CoinOperation{
			Amount: tx.Amount,
		}

		if isReceived {
			operation.FromUser = tx.SenderName
		} else {
			operation.ToUser = tx.ReceiverName
		}

		result[i] = operation
	}

	return result
}
