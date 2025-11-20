package handlers

import (
	"net/http"
	"time"

	"github.com/ViniciusBoroto/cabeleleila_leila/internal/models"
	"github.com/ViniciusBoroto/cabeleleila_leila/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	rg.GET("/appointments/:id", h.GetAppointment)
	rg.GET("/appointments", h.ListAppointments)
	rg.GET("/appointments/suggest/:clientID", h.SuggestSameWeek)

	// Operacional
	rg.GET("/admin/incoming", h.ListIncoming)
	rg.POST("/admin/appointments/:id/confirm", h.ConfirmAppointment)
	rg.PUT("/admin/appointments/:id/services/:serviceID/status", h.UpdateServiceStatus)

	// Gerencial
	rg.GET("/admin/weekly-performance", h.WeeklyPerformance)
}

// CreateAppointment permite agendar um ou mais serviços.
func (h *AppointmentHandler) CreateAppointment(c *gin.Context) {
	var req models.Appointment
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// preenchimentos básicos
	req.ID = uuid.NewString()
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	created, err := h.svc.CreateAppointment(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

// UpdateAppointment aplica a regra: alteração permitida até 2 dias antes.
// Se a data agendada for menor que 2 dias, a alteração só por telefone (aqui retornamos 403).
func (h *AppointmentHandler) UpdateAppointment(c *gin.Context) {
	id := c.Param("id")
	var req models.Appointment
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existing, err := h.svc.GetAppointment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "appointment not found"})
		return
	}

	// Regra dos 2 dias
	now := time.Now()
	if existing.Date.Sub(now) < 48*time.Hour {
		c.JSON(http.StatusForbidden, gin.H{"error": "alteration allowed only by phone within 2 days of appointment"})
		return
	}

	req.UpdatedAt = time.Now()
	updated, err := h.svc.UpdateAppointment(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// GetAppointment retorna detalhes de um agendamento.
func (h *AppointmentHandler) GetAppointment(c *gin.Context) {
	id := c.Param("id")
	appt, err := h.svc.GetAppointment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "appointment not found"})
		return
	}
	c.JSON(http.StatusOK, appt)
}

// ListAppointments retorna histórico / agendamentos em um período.
func (h *AppointmentHandler) ListAppointments(c *gin.Context) {
	var filter models.AppointmentFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Se from/to não fornecido, default: último mês
	from := filter.From
	to := filter.To
	if from.IsZero() || to.IsZero() {
		to = time.Now()
		from = to.AddDate(0, -1, 0)
	}

	list, err := h.svc.ListAppointments(c.Request.Context(), from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// SuggestSameWeek sugere reagendamento na mesma semana (data do primeiro agendamento).
func (h *AppointmentHandler) SuggestSameWeek(c *gin.Context) {
	clientID := c.Param("clientID")
	weekOfStr := c.Query("weekOf") // opcional, formato yyyy-mm-dd
	var weekOf time.Time
	var err error
	if weekOfStr == "" {
		weekOf = time.Now()
	} else {
		weekOf, err = time.Parse("2006-01-02", weekOfStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid weekOf format, use YYYY-MM-DD"})
			return
		}
	}

	suggested, err := h.svc.SuggestSameWeekDate(c.Request.Context(), clientID, weekOf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if suggested == nil {
		c.JSON(http.StatusNoContent, nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{"suggested_date": suggested.Format(time.RFC3339)})
}

// ListIncoming - listagem operacional de agendamentos recebidos.
func (h *AppointmentHandler) ListIncoming(c *gin.Context) {
	list, err := h.svc.ListIncoming(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// ConfirmAppointment - confirma um agendamento (operacional).
func (h *AppointmentHandler) ConfirmAppointment(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.ConfirmAppointment(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// UpdateServiceStatus - atualiza status de um serviço específico dentro do agendamento.
func (h *AppointmentHandler) UpdateServiceStatus(c *gin.Context) {
	appointmentID := c.Param("id")
	serviceID := c.Param("serviceID")

	var body struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.UpdateServiceStatus(c.Request.Context(), appointmentID, serviceID, body.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// WeeklyPerformance - endpoint gerencial com desempenho semanal.
func (h *AppointmentHandler) WeeklyPerformance(c *gin.Context) {
	weekStartStr := c.Query("weekStart")
	var weekStart time.Time
	var err error
	if weekStartStr == "" {
		// calcula início da semana atual (segunda-feira)
		now := time.Now()
		weekday := int(now.Weekday())
		// Go: Sunday==0, Monday==1
		daysToMonday := (weekday + 6) % 7
		weekStart = time.Date(now.Year(), now.Month(), now.Day()-daysToMonday, 0, 0, 0, 0, now.Location())
	} else {
		weekStart, err = time.Parse("2006-01-02", weekStartStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid weekStart format, use YYYY-MM-DD"})
			return
		}
	}

	perf, err := h.svc.WeeklyPerformance(c.Request.Context(), weekStart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, perf)
}
