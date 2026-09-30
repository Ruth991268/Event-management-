package models

// Location stores event place information
type Location struct {
	ID      int64   `json:"id" db:"id"`
	Name    string  `json:"name" db:"name"`
	City    string  `json:"city,omitempty" db:"city"`
	Country string  `json:"country,omitempty" db:"country"`
	MapLat  float64 `json:"map_lat,omitempty" db:"map_lat"`
	MapLng  float64 `json:"map_lng,omitempty" db:"map_lng"`
}
