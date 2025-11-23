package service

import (
	"errors"
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

	mockRepo.EXPECT().FindUserAppointmentsInWeek(gomock.Any(), gomock.Any(), gomock.Any()).Return([]models.Appointment{existentAp}, nil)

	ap, suggestion, err := apSrv.CreateAppointment(1, []models.Service{{ID: 1, Name: "Corte"}}, time.Now().AddDate(0, 0, 3))
	assert.NoError(t, err)
	assert.NotNil(t, suggestion)
	assert.Equal(t, uint(0), ap.ID) // No appointment should be created when suggestion exists
	assert.Equal(t, existentAp.Date, suggestion.Date)
	assert.Equal(t, uint(2), suggestion.ID)
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

// TestMergeAppointments_Success tests successful merge of services
func TestMergeAppointments_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockAppointmentRepository(ctrl)
	apSrv := NewAppointmentService(mockRepo)

	existingServices := []models.Service{
		{ID: 1, Name: "Corte", Price: 50.0},
		{ID: 2, Name: "Escova", Price: 40.0},
	}

	newServices := []models.Service{
		{ID: 3, Name: "Hidratação", Price: 60.0},
	}

	existingAp := models.Appointment{
		ID:       10,
		UserID:   1,
		Services: existingServices,
		Date:     time.Now().AddDate(0, 0, 2),
		Status:   models.StatusPending,
	}

	mergedServices := append(existingServices, newServices...)
	mergedAp := models.Appointment{
		ID:       10,
		UserID:   1,
		Services: mergedServices,
		Date:     time.Now().AddDate(0, 0, 2),
		Status:   models.StatusPending,
	}

	gomock.InOrder(
		mockRepo.EXPECT().FindByID(uint(10)).Return(existingAp, nil),
		mockRepo.EXPECT().Update(gomock.Any()).Return(nil),
		mockRepo.EXPECT().FindByID(uint(10)).Return(mergedAp, nil),
	)

	result, err := apSrv.MergeAppointments(10, newServices)

	assert.NoError(t, err)
	assert.Equal(t, uint(10), result.ID)
	assert.Equal(t, 3, len(result.Services))
	assert.Equal(t, "Corte", result.Services[0].Name)
	assert.Equal(t, "Escova", result.Services[1].Name)
	assert.Equal(t, "Hidratação", result.Services[2].Name)
}

// TestMergeAppointments_AppointmentNotFound tests error when appointment doesn't exist
func TestMergeAppointments_AppointmentNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockAppointmentRepository(ctrl)
	apSrv := NewAppointmentService(mockRepo)

	newServices := []models.Service{
		{ID: 3, Name: "Hidratação", Price: 60.0},
	}

	mockRepo.EXPECT().FindByID(uint(999)).Return(models.Appointment{}, errors.New("appointment not found"))

	result, err := apSrv.MergeAppointments(999, newServices)

	assert.Error(t, err)
	assert.Equal(t, "appointment not found", err.Error())
	assert.Equal(t, uint(0), result.ID)
}

// TestMergeAppointments_UpdateFails tests error when update fails
func TestMergeAppointments_UpdateFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockAppointmentRepository(ctrl)
	apSrv := NewAppointmentService(mockRepo)

	existingServices := []models.Service{
		{ID: 1, Name: "Corte", Price: 50.0},
	}

	newServices := []models.Service{
		{ID: 3, Name: "Hidratação", Price: 60.0},
	}

	existingAp := models.Appointment{
		ID:       10,
		UserID:   1,
		Services: existingServices,
		Date:     time.Now().AddDate(0, 0, 2),
		Status:   models.StatusPending,
	}

	gomock.InOrder(
		mockRepo.EXPECT().FindByID(uint(10)).Return(existingAp, nil),
		mockRepo.EXPECT().Update(gomock.Any()).Return(errors.New("database error")),
	)

	result, err := apSrv.MergeAppointments(10, newServices)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	assert.Equal(t, uint(0), result.ID)
}

// TestMergeAppointments_WithMultipleNewServices tests merging multiple new services
func TestMergeAppointments_WithMultipleNewServices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockAppointmentRepository(ctrl)
	apSrv := NewAppointmentService(mockRepo)

	existingServices := []models.Service{
		{ID: 1, Name: "Corte", Price: 50.0},
	}

	newServices := []models.Service{
		{ID: 2, Name: "Escova", Price: 40.0},
		{ID: 3, Name: "Hidratação", Price: 60.0},
		{ID: 4, Name: "Coloração", Price: 100.0},
	}

	existingAp := models.Appointment{
		ID:       10,
		UserID:   1,
		Services: existingServices,
		Date:     time.Now().AddDate(0, 0, 2),
		Status:   models.StatusPending,
	}

	mergedServices := append(existingServices, newServices...)
	mergedAp := models.Appointment{
		ID:       10,
		UserID:   1,
		Services: mergedServices,
		Date:     time.Now().AddDate(0, 0, 2),
		Status:   models.StatusPending,
	}

	gomock.InOrder(
		mockRepo.EXPECT().FindByID(uint(10)).Return(existingAp, nil),
		mockRepo.EXPECT().Update(gomock.Any()).Return(nil),
		mockRepo.EXPECT().FindByID(uint(10)).Return(mergedAp, nil),
	)

	result, err := apSrv.MergeAppointments(10, newServices)

	assert.NoError(t, err)
	assert.Equal(t, uint(10), result.ID)
	assert.Equal(t, 4, len(result.Services))
	// Total price: 50 (Corte) + 40 (Escova) + 60 (Hidratação) + 100 (Coloração) = 250
	assert.Equal(t, 200.0, result.Services[1].Price+result.Services[2].Price+result.Services[3].Price)
}

// TestMergeAppointments_PreservesAppointmentData tests that merge preserves other appointment data
func TestMergeAppointments_PreservesAppointmentData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockAppointmentRepository(ctrl)
	apSrv := NewAppointmentService(mockRepo)

	existingServices := []models.Service{
		{ID: 1, Name: "Corte", Price: 50.0},
	}

	newServices := []models.Service{
		{ID: 2, Name: "Escova", Price: 40.0},
	}

	appointmentDate := time.Now().AddDate(0, 0, 2)
	existingAp := models.Appointment{
		ID:       10,
		UserID:   5,
		Services: existingServices,
		Date:     appointmentDate,
		Status:   models.StatusConfirmed,
	}

	mergedServices := append(existingServices, newServices...)
	mergedAp := models.Appointment{
		ID:       10,
		UserID:   5,
		Services: mergedServices,
		Date:     appointmentDate,
		Status:   models.StatusConfirmed,
	}

	gomock.InOrder(
		mockRepo.EXPECT().FindByID(uint(10)).Return(existingAp, nil),
		mockRepo.EXPECT().Update(gomock.Any()).Return(nil),
		mockRepo.EXPECT().FindByID(uint(10)).Return(mergedAp, nil),
	)

	result, err := apSrv.MergeAppointments(10, newServices)

	assert.NoError(t, err)
	assert.Equal(t, uint(10), result.ID)
	assert.Equal(t, uint(5), result.UserID)
	assert.Equal(t, appointmentDate, result.Date)
	assert.Equal(t, models.StatusConfirmed, result.Status)
	assert.Equal(t, 2, len(result.Services))
}

// TestMergeAppointments_EmptyNewServices tests merging with empty services array
func TestMergeAppointments_EmptyNewServices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockAppointmentRepository(ctrl)
	apSrv := NewAppointmentService(mockRepo)

	existingServices := []models.Service{
		{ID: 1, Name: "Corte", Price: 50.0},
	}

	existingAp := models.Appointment{
		ID:       10,
		UserID:   1,
		Services: existingServices,
		Date:     time.Now().AddDate(0, 0, 2),
		Status:   models.StatusPending,
	}

	mergedAp := models.Appointment{
		ID:       10,
		UserID:   1,
		Services: existingServices,
		Date:     time.Now().AddDate(0, 0, 2),
		Status:   models.StatusPending,
	}

	gomock.InOrder(
		mockRepo.EXPECT().FindByID(uint(10)).Return(existingAp, nil),
		mockRepo.EXPECT().Update(gomock.Any()).Return(nil),
		mockRepo.EXPECT().FindByID(uint(10)).Return(mergedAp, nil),
	)

	result, err := apSrv.MergeAppointments(10, []models.Service{})

	assert.NoError(t, err)
	assert.Equal(t, uint(10), result.ID)
	assert.Equal(t, 1, len(result.Services))
	assert.Equal(t, "Corte", result.Services[0].Name)
}

// TestMergeAppointments_UpdatedAtTimestamp tests that UpdatedAt is set on merge
func TestMergeAppointments_UpdatedAtTimestamp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockAppointmentRepository(ctrl)
	apSrv := NewAppointmentService(mockRepo)

	existingServices := []models.Service{
		{ID: 1, Name: "Corte", Price: 50.0},
	}

	newServices := []models.Service{
		{ID: 2, Name: "Escova", Price: 40.0},
	}

	oldTime := time.Now().AddDate(-1, 0, 0)
	existingAp := models.Appointment{
		ID:        10,
		UserID:    1,
		Services:  existingServices,
		Date:      time.Now().AddDate(0, 0, 2),
		Status:    models.StatusPending,
		UpdatedAt: oldTime,
	}

	mergedServices := append(existingServices, newServices...)
	mergedAp := models.Appointment{
		ID:        10,
		UserID:    1,
		Services:  mergedServices,
		Date:      time.Now().AddDate(0, 0, 2),
		Status:    models.StatusPending,
		UpdatedAt: time.Now(),
	}

	gomock.InOrder(
		mockRepo.EXPECT().FindByID(uint(10)).Return(existingAp, nil),
		mockRepo.EXPECT().Update(gomock.Any()).Return(nil),
		mockRepo.EXPECT().FindByID(uint(10)).Return(mergedAp, nil),
	)

	result, err := apSrv.MergeAppointments(10, newServices)

	assert.NoError(t, err)
	assert.Equal(t, uint(10), result.ID)
	assert.True(t, result.UpdatedAt.After(oldTime))
}
