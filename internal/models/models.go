package models

import "time"

type ServiceItem struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	DurationMin int     `json:"duration_min"`
	Price       float64 `json:"price"`
	Status      string  `json:"status"` // ex: pending, done, cancelled
}

type Appointment struct {
	ID         string        `json:"id"`
	ClientName string        `json:"client_name"`
	ClientID   string        `json:"client_id"` // opcional identificador do cliente
	Date       time.Time     `json:"date"`
	Services   []ServiceItem `json:"services"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
	Confirmed  bool          `json:"confirmed"`
}

type AppointmentFilter struct {
	From time.Time `form:"from" time_format:"2006-01-02"`
	To   time.Time `form:"to" time_format:"2006-01-02"`
}
