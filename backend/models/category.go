package models

// Category represents an event category (Music, Tech, etc.)
type Category struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
