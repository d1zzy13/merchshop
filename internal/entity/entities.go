package entity

import "time"

type User struct {
	ID        int
	Username  string
	Password  string
	Balance   int
	CreatedAt time.Time
}

type Merchandise struct {
	Name  string
	Price int
}

type Transaction struct {
	ID           int
	SenderID     int
	ReceiverID   int
	SenderName   string
	ReceiverName string
	Amount       int
	CreatedAt    time.Time
}

type Purchase struct {
	ID         int
	UserID     int
	MerchName  string
	Quantity   int
	TotalPrice int
	CreatedAt  time.Time
}
