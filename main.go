package main

import (
	"context"
	"log"
	"time"

	"github.com/ViniciusBoroto/cabeleleila_leila/internal/handlers"
	"github.com/ViniciusBoroto/cabeleleila_leila/internal/models"
	"github.com/gin-gonic/gin"
	// "github.com/seuusuario/cabeleleila/internal/service/postgres" // sua implementação real
)

func main() {
	r := gin.Default()

	// aqui você injeta a implementação concreta do service
	// ex: svc := postgres.NewPostgresAppointmentService(db)
	var svc handlersNullService // temporário até implementar
	// substitua com a instância real:
	// svc := postgres.NewPostgresAppointmentService(db)

	h := handlers.NewAppointmentHandler(svc)
	api := r.Group("/api")
	h.RegisterRoutes(api)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

// handlersNullService é um stub para compilar o main enquanto você implementa o service.
// Você deve remover/substituir isso pela implementação real.
type handlersNullService struct{}

func (handlersNullService) CreateAppointment(_ context.Context, a *models.Appointment) (*models.Appointment, error) {
	return a, nil
}
func (handlersNullService) UpdateAppointment(_ context.Context, id string, a *models.Appointment) (*models.Appointment, error) {
	return a, nil
}
func (handlersNullService) GetAppointment(_ context.Context, id string) (*models.Appointment, error) {
	return &models.Appointment{ID: id}, nil
}
func (handlersNullService) ListAppointments(_ context.Context, from, to time.Time) ([]*models.Appointment, error) {
	return []*models.Appointment{}, nil
}
func (handlersNullService) SuggestSameWeekDate(_ context.Context, clientID string, weekOf time.Time) (*time.Time, error) {
	return nil, nil
}
func (handlersNullService) ListIncoming(_ context.Context) ([]*models.Appointment, error) {
	return []*models.Appointment{}, nil
}
func (handlersNullService) ConfirmAppointment(_ context.Context, id string) error {
	return nil
}
func (handlersNullService) UpdateServiceStatus(_ context.Context, appointmentID, serviceID, status string) error {
	return nil
}
func (handlersNullService) WeeklyPerformance(_ context.Context, weekStart time.Time) (map[string]any, error) {
	return map[string]any{}, nil
}
