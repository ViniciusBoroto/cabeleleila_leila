package service

import (
	"testing"

	"github.com/ViniciusBoroto/cabeleleila_leila/internal/mocks"
	"github.com/ViniciusBoroto/cabeleleila_leila/internal/models"
	"github.com/golang/mock/gomock"
)

func TestCreateService_WithSuggestion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockAppointmentRepository(ctrl)
	apSrv := NewAppointmentService(mockRepo)

	mockRepo.EXPECT().Create(gomock.Any()).Return(models.Appointment{ID: 5}, nil)

}
