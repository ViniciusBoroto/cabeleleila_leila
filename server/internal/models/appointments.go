package models

import (
	"errors"
	"time"
)

type AppointmentStatus string

const (
	StatusPending   AppointmentStatus = "PENDING"
	StatusConfirmed AppointmentStatus = "CONFIRMED"
	StatusDone      AppointmentStatus = "DONE"
	StatusCanceled  AppointmentStatus = "CANCELED"
)

var (
	ErrAppointmentNoServices = errors.New("appointment must have at least one service")
)

type Appointment struct {
	ID        uint              `gorm:"primaryKey" json:"id"`
	User      User              `gorm:"foreignKey:UserID" json:"user"`
	UserID    uint              `json:"user_id"`
	Services  []Service         `json:"services" gorm:"many2many:appointment_services;"`
	Date      time.Time         `json:"date"`
	Status    AppointmentStatus `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// Validate checks if the appointment is valid
func (a *Appointment) Validate() error {
	if len(a.Services) == 0 {
		return ErrAppointmentNoServices
	}
	return nil
}

type AppointmentFilter struct {
	UserID    *uint      `json:"user_id"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
}
