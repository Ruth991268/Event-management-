package models

import "time"

// Ticket represents a ticket purchase/reservation
type Ticket struct {
	ID          int64     `json:"id" db:"id"`
	UserID      int64     `json:"user_id" db:"user_id"`
	EventID     int64     `json:"event_id" db:"event_id"`
	Quantity    int       `json:"quantity" db:"quantity"`
	TotalPrice  float64   `json:"total_price" db:"total_price"`
	PurchasedAt time.Time `json:"purchased_at" db:"purchased_at"`
}
