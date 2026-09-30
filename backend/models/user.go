package models

import "time"

type User struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	
	// Do NOT return password in API responses
	Password  string    `json:"-" db:"password"`

	AvatarURL string    `json:"avatar_url,omitempty" db:"avatar_url"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
