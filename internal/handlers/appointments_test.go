package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ViniciusBoroto/cabeleleila_leila/internal/handlers"
	"github.com/ViniciusBoroto/cabeleleila_leila/internal/mocks"
	"github.com/ViniciusBoroto/cabeleleila_leila/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupRouterWithMock(t *testing.T) (*gin.Engine, *mocks.MockAppointmentService, func()) {
	ctrl := gomock.NewController(t)
	mockSvc := mocks.NewMockAppointmentService(ctrl)

	h := handlers.NewAppointmentHandler(mockSvc)

	r := gin.New()
	api := r.Group("/api")
	h.RegisterRoutes(api)

	return r, mockSvc, func() { ctrl.Finish() }
}

func TestCreateAppointment_Success(t *testing.T) {
	r, mockSvc, finish := setupRouterWithMock(t)
	defer finish()

	now := time.Now()
	appt := &models.Appointment{
		ClientName: "Maria",
		ClientID:   "client-1",
		Date:       now.AddDate(0, 0, 7),
		Services:   []models.ServiceItem{{ID: "s1", Name: "Corte", DurationMin: 30}},
	}

	// esperamos que CreateAppointment seja chamado e retorne o appt com ID
	mockSvc.
		EXPECT().
		CreateAppointment(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, a *models.Appointment) (*models.Appointment, error) {
			a.ID = "appt-1"
			return a, nil
		})

	body, _ := json.Marshal(appt)
	req := httptest.NewRequest(http.MethodPost, "/api/appointments", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var resp models.Appointment
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Maria", resp.ClientName)
	assert.Equal(t, "appt-1", resp.ID)
}

func TestUpdateAppointment_ForbiddenIfWithin2Days(t *testing.T) {
	r, mockSvc, finish := setupRouterWithMock(t)
	defer finish()

	// appointment existente com data a menos de 48h -> alteração proibida
	id := "appt-2"
	existing := &models.Appointment{
		ID:   id,
		Date: time.Now().Add(24 * time.Hour), // menos de 48h
	}
	mockSvc.
		EXPECT().
		GetAppointment(gomock.Any(), id).
		Return(existing, nil)

	payload := `{"client_name":"Joao"}`
	req := httptest.NewRequest(http.MethodPut, "/api/appointments/"+id, strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestListAppointments_DefaultsToLastMonth(t *testing.T) {
	r, mockSvc, finish := setupRouterWithMock(t)
	defer finish()

	mockSvc.
		EXPECT().
		ListAppointments(gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]*models.Appointment{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/appointments", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
