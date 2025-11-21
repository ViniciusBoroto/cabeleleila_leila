package models

type Service struct {
	ID              uint    `gorm:"primaryKey" json:"id"`
	Name            string  `json:"name"`
	Price           float64 `json:"price"`
	DurationMinutes int     `json:"duration_minutes"`
}
