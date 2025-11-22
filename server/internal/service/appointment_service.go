package service

import (
	"errors"
	"time"

	"github.com/ViniciusBoroto/cabeleleila_leila/internal/models"
	"github.com/ViniciusBoroto/cabeleleila_leila/internal/repository"
)

type AppointmentService interface {
	CreateAppointment(userID uint, services []models.Service, date time.Time) (created models.Appointment, suggestion *models.Appointment, res error)
	UpdateAppointment(id uint, newAp models.Appointment) (models.Appointment, error)
	ListHistory(start, end time.Time) ([]models.Appointment, error)
	ListAll() ([]models.Appointment, error)
	ConfirmAppointment(id uint) error
	GetWeeklyPerformance() (int, int, error)
}

type appointmentService struct {
	repo repository.AppointmentRepository
}

func NewAppointmentService(repo repository.AppointmentRepository) AppointmentService {
	return &appointmentService{repo}
}

func getWeekRange(date time.Time) (time.Time, time.Time) {
	weekday := int(date.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	start := date.AddDate(0, 0, -weekday+1)
	end := start.AddDate(0, 0, 6)
	return start, end
}

func (s *appointmentService) CreateAppointment(userID uint, services []models.Service, date time.Time) (created models.Appointment, suggestion *models.Appointment, err error) {
	// Validate that appointment has at least one service
	if len(services) == 0 {
		return models.Appointment{}, nil, models.ErrAppointmentNoServices
	}

	weekStart, weekEnd := getWeekRange(date)

	// Sugestão de unificação de datas
	existing, _ := s.repo.FindUserAppointmentsInWeek(userID, weekStart, weekEnd)
	if len(existing) > 0 {
		first := existing[0]
		suggestion = &first
	}

	ap := models.Appointment{
		UserID:   userID,
		Services: services,
		Date:     date,
		Status:   models.StatusPending,
	}

	created, err = s.repo.Create(ap)
	return
}

func (s *appointmentService) UpdateAppointment(id uint, newAp models.Appointment) (models.Appointment, error) {
	newAp.ID = id

	ap, err := s.repo.FindByID(id)
	if err != nil {
		return ap, err
	}

	diff := time.Until(ap.Date)
	if diff < (48 * time.Hour) {
		return ap, errors.New("alterações só podem ser feitas por telefone")
	}

	err = s.repo.Update(newAp)
	return ap, err
}

func (s *appointmentService) ListHistory(start, end time.Time) ([]models.Appointment, error) {
	return s.repo.ListByPeriod(start, end)
}

func (s *appointmentService) ListAll() ([]models.Appointment, error) {
	return s.repo.ListAll()
}

func (s *appointmentService) ConfirmAppointment(id uint) error {
	ap, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	ap.Status = models.StatusConfirmed
	return s.repo.Update(ap)
}

func (s *appointmentService) GetWeeklyPerformance() (int, int, error) {
	now := time.Now()
	start, end := getWeekRange(now)

	list, err := s.repo.ListByPeriod(start, end)
	if err != nil {
		return 0, 0, err
	}

	completed := 0
	for _, ap := range list {
		if ap.Status == models.StatusDone {
			completed++
		}
	}

	return len(list), completed, nil
}
