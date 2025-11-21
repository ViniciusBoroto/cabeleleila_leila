package service

import (
	"testing"

	"github.com/ViniciusBoroto/cabeleleila_leila/internal/mocks"
	"github.com/ViniciusBoroto/cabeleleila_leila/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestServiceService_CreateService_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockServiceRepository(ctrl)
	svc := NewServiceService(mockRepo)

	service := models.Service{
		Name:            "Corte de Cabelo",
		Price:           50.0,
		DurationMinutes: 30,
	}

	mockRepo.EXPECT().Create(gomock.Any()).Return(models.Service{
		ID:              1,
		Name:            "Corte de Cabelo",
		Price:           50.0,
		DurationMinutes: 30,
	}, nil)

	result, err := svc.CreateService(service)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "Corte de Cabelo", result.Name)
}

func TestServiceService_CreateService_MissingName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockServiceRepository(ctrl)
	svc := NewServiceService(mockRepo)

	service := models.Service{
		Name:            "",
		Price:           50.0,
		DurationMinutes: 30,
	}

	_, err := svc.CreateService(service)
	assert.Error(t, err)
	assert.Equal(t, "service name is required", err.Error())
}

func TestServiceService_CreateService_NegativePrice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockServiceRepository(ctrl)
	svc := NewServiceService(mockRepo)

	service := models.Service{
		Name:            "Corte",
		Price:           -10.0,
		DurationMinutes: 30,
	}

	_, err := svc.CreateService(service)
	assert.Error(t, err)
	assert.Equal(t, "service price cannot be negative", err.Error())
}

func TestServiceService_CreateService_InvalidDuration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockServiceRepository(ctrl)
	svc := NewServiceService(mockRepo)

	service := models.Service{
		Name:            "Corte",
		Price:           50.0,
		DurationMinutes: 0,
	}

	_, err := svc.CreateService(service)
	assert.Error(t, err)
	assert.Equal(t, "service duration must be greater than 0", err.Error())
}

func TestServiceService_GetService_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockServiceRepository(ctrl)
	svc := NewServiceService(mockRepo)

	mockRepo.EXPECT().FindByID(uint(1)).Return(models.Service{
		ID:              1,
		Name:            "Escova",
		Price:           40.0,
		DurationMinutes: 45,
	}, nil)

	result, err := svc.GetService(1)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "Escova", result.Name)
}

func TestServiceService_GetService_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockServiceRepository(ctrl)
	svc := NewServiceService(mockRepo)

	_, err := svc.GetService(0)
	assert.Error(t, err)
	assert.Equal(t, "invalid service ID", err.Error())
}

func TestServiceService_ListServices_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockServiceRepository(ctrl)
	svc := NewServiceService(mockRepo)

	services := []models.Service{
		{ID: 1, Name: "Corte", Price: 50.0, DurationMinutes: 30},
		{ID: 2, Name: "Escova", Price: 40.0, DurationMinutes: 45},
	}

	mockRepo.EXPECT().FindAll().Return(services, nil)

	result, err := svc.ListServices()
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestServiceService_UpdateService_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockServiceRepository(ctrl)
	svc := NewServiceService(mockRepo)

	updated := models.Service{
		ID:              1,
		Name:            "Corte Premium",
		Price:           60.0,
		DurationMinutes: 40,
	}

	mockRepo.EXPECT().Update(gomock.Any()).Return(nil)
	mockRepo.EXPECT().FindByID(uint(1)).Return(updated, nil)

	result, err := svc.UpdateService(updated)
	assert.NoError(t, err)
	assert.Equal(t, "Corte Premium", result.Name)
	assert.Equal(t, 60.0, result.Price)
}

func TestServiceService_UpdateService_MissingID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockServiceRepository(ctrl)
	svc := NewServiceService(mockRepo)

	service := models.Service{
		ID:              0,
		Name:            "Corte",
		Price:           50.0,
		DurationMinutes: 30,
	}

	_, err := svc.UpdateService(service)
	assert.Error(t, err)
	assert.Equal(t, "service ID is required", err.Error())
}

func TestServiceService_DeleteService_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockServiceRepository(ctrl)
	svc := NewServiceService(mockRepo)

	mockRepo.EXPECT().FindByID(uint(1)).Return(models.Service{ID: 1}, nil)
	mockRepo.EXPECT().Delete(uint(1)).Return(nil)

	err := svc.DeleteService(1)
	assert.NoError(t, err)
}

func TestServiceService_DeleteService_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockServiceRepository(ctrl)
	svc := NewServiceService(mockRepo)

	err := svc.DeleteService(0)
	assert.Error(t, err)
	assert.Equal(t, "invalid service ID", err.Error())
}

func TestServiceService_DeleteService_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockServiceRepository(ctrl)
	svc := NewServiceService(mockRepo)

	mockRepo.EXPECT().FindByID(uint(999)).Return(models.Service{}, assert.AnError)

	err := svc.DeleteService(999)
	assert.Error(t, err)
}
