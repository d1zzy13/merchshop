package models

// AuthRequest модель запроса авторизации
// swagger:model AuthRequest
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse модель ответа авторизации
// swagger:model AuthResponse
type AuthResponse struct {
	Token string `json:"token"`
}

// SendCoinRequest модель передачи коинов
// swagger:model SendCoinRequest
type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

// InfoResponse модель ответа информации
// swagger:model InfoResponse
type InfoResponse struct {
	Coins       int             `json:"coins"`
	Inventory   []InventoryItem `json:"inventory"`
	CoinHistory CoinHistoryInfo `json:"coinHistory"`
}

// InventoryItem элемент инвентаря
// swagger:model InventoryItem
type InventoryItem struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

// CoinHistoryInfo история коинов
// swagger:model CoinHistoryInfo
type CoinHistoryInfo struct {
	Received []CoinOperation `json:"received"`
	Sent     []CoinOperation `json:"sent"`
}

// CoinOperation операция с коинами
// swagger:model CoinOperation
type CoinOperation struct {
	FromUser string `json:"fromUser,omitempty"`
	ToUser   string `json:"toUser,omitempty"`
	Amount   int    `json:"amount"`
}

// ErrorResponse модель ошибок
// swagger:model ErrorResponse
type ErrorResponse struct {
	Errors string `json:"errors"`
}
