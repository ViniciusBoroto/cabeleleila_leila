package service

import (
	"testing"
	"time"

	"github.com/ViniciusBoroto/cabeleleila_leila/internal/mocks"
	"github.com/ViniciusBoroto/cabeleleila_leila/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateService_WithSuggestion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockAppointmentRepository(ctrl)
	apSrv := NewAppointmentService(mockRepo)
	existentAp := models.Appointment{ID: 2, Date: time.Now().AddDate(0, 0, 2)}

	mockRepo.EXPECT().Create(gomock.Any()).Return(models.Appointment{ID: 5}, nil)
	mockRepo.EXPECT().FindUserAppointmentsInWeek(gomock.Any(), gomock.Any(), gomock.Any()).Return([]models.Appointment{existentAp}, nil)

	ap, suggestion, err := apSrv.CreateAppointment(1, []models.Service{{ID: 1, Name: "Corte"}}, time.Now().AddDate(0, 0, 3))
	assert.NoError(t, err)
	assert.NotNil(t, suggestion)
	assert.Equal(t, uint(5), ap.ID)
	assert.Equal(t, existentAp.Date, suggestion.Date)
}

func TestCreateService_NoSuggestion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockAppointmentRepository(ctrl)
	apSrv := NewAppointmentService(mockRepo)

	mockRepo.EXPECT().Create(gomock.Any()).Return(models.Appointment{ID: 5}, nil)
	mockRepo.EXPECT().FindUserAppointmentsInWeek(gomock.Any(), gomock.Any(), gomock.Any()).Return([]models.Appointment{}, nil)

	ap, suggestion, err := apSrv.CreateAppointment(1, []models.Service{{ID: 1, Name: "Corte"}}, time.Now().AddDate(0, 0, 3))
	assert.NoError(t, err)
	assert.Nil(t, suggestion)
	assert.Equal(t, uint(5), ap.ID)
}

// TestCreateAppointment_NoServices tests that appointment creation fails without services
func TestCreateAppointment_NoServices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockAppointmentRepository(ctrl)
	apSrv := NewAppointmentService(mockRepo)

	ap, suggestion, err := apSrv.CreateAppointment(1, []models.Service{}, time.Now().AddDate(0, 0, 3))

	assert.Error(t, err)
	assert.Nil(t, suggestion)
	assert.Equal(t, models.ErrAppointmentNoServices, err)
	assert.Equal(t, uint(0), ap.ID)
}

// TestCreateAppointment_WithServices tests that appointment is created with services
func TestCreateAppointment_WithServices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockAppointmentRepository(ctrl)
	apSrv := NewAppointmentService(mockRepo)

	services := []models.Service{
		{ID: 1, Name: "Corte", Price: 50.0},
		{ID: 2, Name: "Escova", Price: 40.0},
		{ID: 3, Name: "Coloração", Price: 100.0},
	}

	expectedAp := models.Appointment{
		ID:       5,
		UserID:   1,
		Services: services,
		Date:     time.Now().AddDate(0, 0, 3),
		Status:   models.StatusPending,
	}

	mockRepo.EXPECT().Create(gomock.Any()).Return(expectedAp, nil)
	mockRepo.EXPECT().FindUserAppointmentsInWeek(gomock.Any(), gomock.Any(), gomock.Any()).Return([]models.Appointment{}, nil)

	ap, suggestion, err := apSrv.CreateAppointment(1, services, time.Now().AddDate(0, 0, 3))

	assert.NoError(t, err)
	assert.Nil(t, suggestion)
	assert.Equal(t, uint(5), ap.ID)
	assert.Equal(t, uint(1), ap.UserID)
	assert.Equal(t, 3, len(ap.Services))
}

// TestCreateAppointment_UserLoaded tests that user is loaded in response
func TestCreateAppointment_UserLoaded(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockAppointmentRepository(ctrl)
	apSrv := NewAppointmentService(mockRepo)

	user := models.User{
		ID:       1,
		Email:    "test@example.com",
		Name:     "Test User",
		Phone:    "123456789",
		Role:     models.RoleCustomer,
		IsActive: true,
	}

	services := []models.Service{{ID: 1, Name: "Corte", Price: 50.0}}

	expectedAp := models.Appointment{
		ID:       5,
		User:     user,
		UserID:   1,
		Services: services,
		Date:     time.Now().AddDate(0, 0, 3),
		Status:   models.StatusPending,
	}

	mockRepo.EXPECT().Create(gomock.Any()).Return(expectedAp, nil)
	mockRepo.EXPECT().FindUserAppointmentsInWeek(gomock.Any(), gomock.Any(), gomock.Any()).Return([]models.Appointment{}, nil)

	ap, _, err := apSrv.CreateAppointment(1, services, time.Now().AddDate(0, 0, 3))

	assert.NoError(t, err)
	assert.NotEqual(t, uint(0), ap.User.ID)
	assert.Equal(t, "test@example.com", ap.User.Email)
	assert.Equal(t, "Test User", ap.User.Name)
}

// TestCreateAppointment_ServicesNotNil tests that services are not nil
func TestCreateAppointment_ServicesNotNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockAppointmentRepository(ctrl)
	apSrv := NewAppointmentService(mockRepo)

	services := []models.Service{
		{ID: 1, Name: "Corte", Price: 50.0},
		{ID: 2, Name: "Escova", Price: 40.0},
	}

	expectedAp := models.Appointment{
		ID:       5,
		UserID:   1,
		Services: services,
		Date:     time.Now().AddDate(0, 0, 3),
		Status:   models.StatusPending,
	}

	mockRepo.EXPECT().Create(gomock.Any()).Return(expectedAp, nil)
	mockRepo.EXPECT().FindUserAppointmentsInWeek(gomock.Any(), gomock.Any(), gomock.Any()).Return([]models.Appointment{}, nil)

	ap, _, err := apSrv.CreateAppointment(1, services, time.Now().AddDate(0, 0, 3))

	assert.NoError(t, err)
	assert.NotNil(t, ap.Services)
	assert.Equal(t, 2, len(ap.Services))
	assert.Equal(t, "Corte", ap.Services[0].Name)
	assert.Equal(t, "Escova", ap.Services[1].Name)
}
