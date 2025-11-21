package handlers

import (
    "net/http"
    "strconv"

    "github.com/ViniciusBoroto/cabeleleila_leila/internal/models"
    "github.com/ViniciusBoroto/cabeleleila_leila/internal/service"
    "github.com/gin-gonic/gin"
)

type CreateServiceRequest struct {
    Name            string  `json:"name" binding:"required"`
    Price           float64 `json:"price" binding:"required,min=0"`
    DurationMinutes int     `json:"duration_minutes" binding:"required,min=1"`
}

type UpdateServiceRequest struct {
    Name            string  `json:"name" binding:"required"`
    Price           float64 `json:"price" binding:"required,min=0"`
    DurationMinutes int     `json:"duration_minutes" binding:"required,min=1"`
}

type ServiceResponse struct {
    ID              uint    `json:"id"`
    Name            string  `json:"name"`
    Price           float64 `json:"price"`
    DurationMinutes int     `json:"duration_minutes"`
}

// ListServices godoc
// @Summary      List all services
// @Description  Retrieve all available services
// @Tags         services
// @Produce      json
// @Success      200  {array}   ServiceResponse
// @Failure      500  {object}  map[string]string
// @Router       /services [get]
func ListServices(svc service.ServiceService) gin.HandlerFunc {
    return func(c *gin.Context) {
        services, err := svc.ListServices()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        response := make([]ServiceResponse, len(services))
        for i, s := range services {
            response[i] = ServiceResponse{
                ID:              s.ID,
                Name:            s.Name,
                Price:           s.Price,
                DurationMinutes: s.DurationMinutes,
            }
        }
        c.JSON(http.StatusOK, response)
    }
}

// GetService godoc
// @Summary      Get service by ID
// @Description  Retrieve a specific service by ID
// @Tags         services
// @Produce      json
// @Param        id   path      int  true  "Service ID"
// @Success      200  {object}  ServiceResponse
// @Failure      404  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /services/{id} [get]
func GetService(svc service.ServiceService) gin.HandlerFunc {
    return func(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 32)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service ID"})
            return
        }

        srv, err := svc.GetService(uint(id))
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
            return
        }

        response := ServiceResponse{
            ID:              srv.ID,
            Name:            srv.Name,
            Price:           srv.Price,
            DurationMinutes: srv.DurationMinutes,
        }
        c.JSON(http.StatusOK, response)
    }
}

// CreateService godoc
// @Summary      Create a new service (admin only)
// @Description  Create a new service
// @Tags         services
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        service  body      CreateServiceRequest  true  "Service data"
// @Success      201      {object}  ServiceResponse
// @Failure      400      {object}  map[string]string
// @Failure      403      {object}  map[string]string
// @Router       /admin/services [post]
func CreateService(svc service.ServiceService) gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("role")
        if !exists || role != models.RoleAdmin {
            c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
            return
        }

        var req CreateServiceRequest
        if err := c.BindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        srv := models.Service{
            Name:            req.Name,
            Price:           req.Price,
            DurationMinutes: req.DurationMinutes,
        }

        created, err := svc.CreateService(srv)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        response := ServiceResponse{
            ID:              created.ID,
            Name:            created.Name,
            Price:           created.Price,
            DurationMinutes: created.DurationMinutes,
        }
        c.JSON(http.StatusCreated, response)
    }
}

// UpdateService godoc
// @Summary      Update service (admin only)
// @Description  Update a specific service
// @Tags         services
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        id       path      int                    true  "Service ID"
// @Param        service  body      UpdateServiceRequest   true  "Updated service data"
// @Success      200      {object}  ServiceResponse
// @Failure      400      {object}  map[string]string
// @Failure      403      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Router       /admin/services/{id} [put]
func UpdateService(svc service.ServiceService) gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("role")
        if !exists || role != models.RoleAdmin {
            c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
            return
        }

        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 32)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service ID"})
            return
        }

        var req UpdateServiceRequest
        if err := c.BindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        srv := models.Service{
            ID:              uint(id),
            Name:            req.Name,
            Price:           req.Price,
            DurationMinutes: req.DurationMinutes,
        }

        updated, err := svc.UpdateService(srv)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        response := ServiceResponse{
            ID:              updated.ID,
            Name:            updated.Name,
            Price:           updated.Price,
            DurationMinutes: updated.DurationMinutes,
        }
        c.JSON(http.StatusOK, response)
    }
}

// DeleteService godoc
// @Summary      Delete service (admin only)
// @Description  Delete a specific service
// @Tags         services
// @Security     Bearer
// @Param        id  path  int  true  "Service ID"
// @Success      204
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /admin/services/{id} [delete]
func DeleteService(svc service.ServiceService) gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("role")
        if !exists || role != models.RoleAdmin {
            c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
            return
        }

        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 32)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service ID"})
            return
        }

        if err := svc.DeleteService(uint(id)); err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
            return
        }

        c.Status(http.StatusNoContent)
    }
}