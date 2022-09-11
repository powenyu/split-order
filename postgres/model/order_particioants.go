package model

type OrderParticipant struct {
	ID      uint    `json:"id" db:"id"`
	OrderID uint    `json:"order_id" db:"order_id"`
	UserID  string  `json:"user_id" db:"user_id"`
	Price   float64 `json:"price" db:"price"`
}
