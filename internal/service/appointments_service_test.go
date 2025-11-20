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
	mockRepo.EXPECT().FindCustomerAppointmentsInWeek(gomock.Any(), gomock.Any(), gomock.Any()).Return([]models.Appointment{existentAp}, nil)

	ap, suggestion, err := apSrv.CreateAppointment(1, []models.Service{{ID: 1, Name: "Corte"}}, time.Now().AddDate(0, 0, 3))
	assert.NoError(t, err)
	assert.NotNil(t, suggestion)
	assert.Equal(t, uint(5), ap.ID)
	assert.Equal(t, existentAp.Date, *suggestion)
}

func TestCreateService_NoSuggestion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockAppointmentRepository(ctrl)
	apSrv := NewAppointmentService(mockRepo)

	mockRepo.EXPECT().Create(gomock.Any()).Return(models.Appointment{ID: 5}, nil)
	mockRepo.EXPECT().FindCustomerAppointmentsInWeek(gomock.Any(), gomock.Any(), gomock.Any()).Return([]models.Appointment{}, nil)

	ap, suggestion, err := apSrv.CreateAppointment(1, []models.Service{{ID: 1, Name: "Corte"}}, time.Now().AddDate(0, 0, 3))
	assert.NoError(t, err)
	assert.Nil(t, suggestion)
	assert.Equal(t, uint(5), ap.ID)
}
