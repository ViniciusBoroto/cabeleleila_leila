package service

import (
	"context"
	"time"

	"github.com/ViniciusBoroto/cabeleleila_leila/internal/models"
)

// AppointmentService é a interface que os handlers usam.
// Implementações concretas (DB, memory, etc) serão injetadas no handler.
type AppointmentService interface {
	CreateAppointment(ctx context.Context, a *models.Appointment) (*models.Appointment, error)
	UpdateAppointment(ctx context.Context, id string, a *models.Appointment) (*models.Appointment, error)
	GetAppointment(ctx context.Context, id string) (*models.Appointment, error)
	ListAppointments(ctx context.Context, from, to time.Time) ([]*models.Appointment, error)
	SuggestSameWeekDate(ctx context.Context, clientID string, weekOf time.Time) (*time.Time, error)

	// Operacional
	ListIncoming(ctx context.Context) ([]*models.Appointment, error)
	ConfirmAppointment(ctx context.Context, id string) error
	UpdateServiceStatus(ctx context.Context, appointmentID, serviceID, status string) error

	// Gerencial
	WeeklyPerformance(ctx context.Context, weekStart time.Time) (map[string]any, error)
}
