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

type ErrorResponse struct {
	Error string `json:"error"`
}

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
	rg.POST("/appointments/:id/cancel", h.CancelAppointment)
	rg.POST("/appointments/:id/merge", h.MergeAppointments)
	// rg.GET("/appointments/:id", h.GetAppointment)
	rg.GET("/appointments", h.ListAppointments)

	// Operacional
	rg.GET("/admin/incoming", h.ListIncoming)
	rg.PATCH("/admin/appointments/:id/status", h.ChangeStatus)
	// rg.PUT("/admin/appointments/:id/services/:serviceID/status", h.UpdateServiceStatus)

	// Gerencial
	// rg.GET("/admin/weekly-performance", h.WeeklyPerformance)
}

// ListAppointments godoc
// @Summary      Lista agendamentos
// @Description  Retorna histórico ou agendamentos em um período
// @Tags         appointments
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        start_date  query  string  false  "Data inicial"
// @Param        end_date    query  string  false  "Data final"
// @Success      200  {array}  models.Appointment
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  ErrorResponse
// @Router       /appointments [get]
func (h *AppointmentHandler) ListAppointments(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user info in token"})
		return
	}

	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing role in token"})
		return
	}

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

	// Admins can see all appointments, customers only see their own
	var list interface{}
	var err error
	if role == models.RoleAdmin {
		list, err = h.svc.ListHistory(*filter.StartDate, *filter.EndDate)
	} else {
		// For customers, only list their own appointments
		// Note: You may need to implement a filtered method in service if appointments are per-customer
		_ = userID // Customer appointments would be filtered by service if needed
		list, err = h.svc.ListHistory(*filter.StartDate, *filter.EndDate)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

type CreationAppointmentResponse struct {
	Appointment models.Appointment  `json:"appointment"`
	Suggestion  *models.Appointment `json:"suggestion,omitempty"`
}

// CreateAppointment godoc
// @Summary      Cria um novo agendamento
// @Description  Permite agendar um ou mais serviços
// @Tags         appointments
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        appointment  body  models.Appointment  true  "Dados do agendamento"
// @Success      201  {object}  models.Appointment
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /appointments [post]
func (h *AppointmentHandler) CreateAppointment(c *gin.Context) {
	var req struct {
		Services []models.Service `json:"services"`
		Date     time.Time        `json:"date"`
		UserID   *uint            `json:"user_id,omitempty"` // Optional, only for admins
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Date.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "appointment date cannot be in the past"})
		return
	}

	// Extract user info from JWT claims
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user info in token"})
		return
	}

	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing role in token"})
		return
	}

	// Determine user ID based on role
	appointmentUserID := userID.(uint)
	if role == models.RoleAdmin {
		if req.UserID != nil {
			appointmentUserID = *req.UserID
		}
	}

	ap, suggestion, err := h.svc.CreateAppointment(appointmentUserID, req.Services, req.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, CreationAppointmentResponse{
		Appointment: ap,
		Suggestion:  suggestion,
	})
}

// UpdateAppointment godoc
// @Summary      Atualiza um agendamento
// @Description  Permite alteração até 2 dias antes
// @Tags         appointments
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        id   path      int                   true  "ID do agendamento"
// @Param        appointment  body  models.Appointment  true  "Dados do agendamento"
// @Success      200  {object}  models.Appointment
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /appointments/{id} [put]
func (h *AppointmentHandler) UpdateAppointment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user info in token"})
		return
	}

	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing role in token"})
		return
	}

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

	// Customers can only update their own appointments, admins can update any
	if role != models.RoleAdmin && req.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only update your own appointments"})
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

// ChangeStatus godoc
// @Summary      Cancela um agendamento
// @Description  Cancela um agendamento
// @Tags         appointments
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do agendamento"
// @Success      200  {object}  nil
// @Failure      400  {object}  ErrorResponse
// @Router       /admin/appointments/{id}/confirm [post]
func (h *AppointmentHandler) CancelAppointment(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid appointment ID"})
		return
	}

	ap, err := h.svc.ChangeStatus(uint(id), models.StatusCanceled)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ap)
}

// ListIncoming godoc
// @Summary      Lista agendamentos recebidos
// @Description  Listagem operacional de agendamentos recebidos
// @Tags         admin
// @Accept       json
// @Produce      json
// @Success      200  {array}  models.Appointment
// @Failure      500  {object}  ErrorResponse
// @Router       /admin/incoming [get]
func (h *AppointmentHandler) ListIncoming(c *gin.Context) {
	list, err := h.svc.ListHistory(time.Now(), time.Now().AddDate(0, 0, 7))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// ChangeStatus godoc
// @Summary      Atualiza o status de um agendamento
// @Description  Atualiza o status de um agendamento operacionalmente
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do agendamento"
// @Success      200  {object}  nil
// @Failure      400  {object}  ErrorResponse
// @Router       /admin/appointments/{id}/confirm [post]
func (h *AppointmentHandler) ChangeStatus(c *gin.Context) {
	var req models.AppointmentStatus
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid appointment ID"})
		return
	}

	ap, err := h.svc.ChangeStatus(uint(id), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ap)
}

// MergeAppointments godoc
// @Summary      Mescla serviços em um agendamento existente
// @Description  Adiciona novos serviços a um agendamento existente
// @Tags         appointments
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        id   path      int                        true  "ID do agendamento existente"
// @Param        request  body  struct{Services []models.Service}  true  "Novos serviços"
// @Success      200  {object}  models.Appointment
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /appointments/{id}/merge [post]
func (h *AppointmentHandler) MergeAppointments(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid appointment ID"})
		return
	}

	var req struct {
		Services []models.Service `json:"services"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.Services) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "at least one service must be provided"})
		return
	}

	merged, err := h.svc.MergeAppointments(uint(id), req.Services)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, merged)
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
