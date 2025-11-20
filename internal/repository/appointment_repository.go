package repository

import (
	"time"

	"github.com/ViniciusBoroto/cabeleleila_leila/internal/models"
)

type AppointmentRepository interface {
	Create(ap models.Appointment) (models.Appointment, error)
	Update(ap models.Appointment) error
	FindByID(id uint) (models.Appointment, error)
	FindCustomerAppointmentsInWeek(customerID uint, weekStart, weekEnd time.Time) ([]models.Appointment, error)
	ListByPeriod(start, end time.Time) ([]models.Appointment, error)
	ListAll() ([]models.Appointment, error)
}
