package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/ViniciusBoroto/cabeleleila_leila/internal/models"
	"github.com/ViniciusBoroto/cabeleleila_leila/internal/service"
	"github.com/gin-gonic/gin"
)

// AppointmentHandler contém dependências (service).
type AppointmentHandler struct {
	svc service.AppointmentService
}

// NewAppointmentHandler cria um handler com a service injetada.
func NewAppointmentHandler(svc service.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{svc: svc}
}

// RegisterRoutes registra rotas no router (group).
func (h *AppointmentHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/appointments", h.CreateAppointment)
	rg.PUT("/appointments/:id", h.UpdateAppointment)
	// rg.GET("/appointments/:id", h.GetAppointment)
	rg.GET("/appointments", h.ListAppointments)

	// Operacional
	rg.GET("/admin/incoming", h.ListIncoming)
	rg.POST("/admin/appointments/:id/confirm", h.ConfirmAppointment)
	// rg.PUT("/admin/appointments/:id/services/:serviceID/status", h.UpdateServiceStatus)

	// Gerencial
	// rg.GET("/admin/weekly-performance", h.WeeklyPerformance)
}

// CreateAppointment permite agendar um ou mais serviços.
func (h *AppointmentHandler) CreateAppointment(c *gin.Context) {
	var req struct {
		CustomerID uint             `json:"customer_id"`
		Services   []models.Service `json:"service"`
		Date       string           `json:"date"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, _ := time.Parse("2006-01-02 15:04", req.Date)

	ap, suggestion, err := h.svc.CreateAppointment(req.CustomerID, req.Services, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := gin.H{"appointment": ap}
	if suggestion != nil {
		resp["suggestion"] = suggestion.Format("2006-01-02 15:04")
	}

	c.JSON(http.StatusCreated, resp)
}

// UpdateAppointment aplica a regra: alteração permitida até 2 dias antes.
// Se a data agendada for menor que 2 dias, a alteração só por telefone (aqui retornamos 403).
func (h *AppointmentHandler) UpdateAppointment(c *gin.Context) {
	idParam := c.Param("id")
	var req models.Appointment
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := strconv.ParseUint(idParam, 10, 32) // Parse as base 10, target uint64
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid appointment ID"})
		return
	}
	req.UpdatedAt = time.Now()
	updated, err := h.svc.UpdateAppointment(uint(id), req)
	if err != nil {
		if errors.Is(err, models.ErrCannotUpdateWithingTwoDays) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// // GetAppointment retorna detalhes de um agendamento.
// func (h *AppointmentHandler) GetAppointment(c *gin.Context) {
// 	id := c.Param("id")
// 	appt, err := h.svc.GetAppointment(c.Request.Context(), id)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "appointment not found"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, appt)
// }

// ListAppointments retorna histórico / agendamentos em um período.
func (h *AppointmentHandler) ListAppointments(c *gin.Context) {
	var filter models.AppointmentFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if filter.StartDate == nil {
		dftStart := time.Now().AddDate(0, -1, 0)
		filter.StartDate = &dftStart
	}
	if filter.EndDate == nil {
		dftEnd := time.Now().AddDate(0, 1, 0)
		filter.EndDate = &dftEnd
	}

	list, err := h.svc.ListHistory(*filter.StartDate, *filter.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// ListIncoming - listagem operacional de agendamentos recebidos.
func (h *AppointmentHandler) ListIncoming(c *gin.Context) {
	list, err := h.svc.ListHistory(time.Now(), time.Now().AddDate(0, 0, 7))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// ConfirmAppointment - confirma um agendamento (operacional).
func (h *AppointmentHandler) ConfirmAppointment(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid appointment ID"})
		return
	}

	if err := h.svc.ConfirmAppointment(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// // WeeklyPerformance - endpoint gerencial com desempenho semanal.
// func (h *AppointmentHandler) WeeklyPerformance(c *gin.Context) {
// 	weekStartStr := c.Query("weekStart")
// 	var weekStart time.Time
// 	var err error
// 	if weekStartStr == "" {
// 		// calcula início da semana atual (segunda-feira)
// 		now := time.Now()
// 		weekday := int(now.Weekday())
// 		// Go: Sunday==0, Monday==1
// 		daysToMonday := (weekday + 6) % 7
// 		weekStart = time.Date(now.Year(), now.Month(), now.Day()-daysToMonday, 0, 0, 0, 0, now.Location())
// 	} else {
// 		weekStart, err = time.Parse("2006-01-02", weekStartStr)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid weekStart format, use YYYY-MM-DD"})
// 			return
// 		}
// 	}

// 	perf, err := h.svc.WeeklyPerformance(c.Request.Context(), weekStart)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, perf)
// }
