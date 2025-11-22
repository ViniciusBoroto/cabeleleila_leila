package models

import (
	"time"
)

type AppointmentStatus string

const (
	StatusPending   AppointmentStatus = "PENDING"
	StatusConfirmed AppointmentStatus = "CONFIRMED"
	StatusDone      AppointmentStatus = "DONE"
	StatusCanceled  AppointmentStatus = "CANCELED"
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

type AppointmentFilter struct {
	UserID    *uint      `json:"user_id"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
}
