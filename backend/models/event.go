package models

import "time"

// Event stores event information
type Event struct {
	ID          int64     `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description,omitempty" db:"description"`
	CategoryID  int64     `json:"category_id,omitempty" db:"category_id"`
	LocationID  int64     `json:"location_id,omitempty" db:"location_id"`
	Price       float64   `json:"price,omitempty" db:"price"`
	StartDate   time.Time `json:"start_date" db:"start_date"`
	EndDate     time.Time `json:"end_date" db:"end_date"`
	ImageURL    string    `json:"image_url,omitempty" db:"image_url"`
	CreatedBy   int64     `json:"created_by" db:"created_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
