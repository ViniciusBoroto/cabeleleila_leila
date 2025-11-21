package handlers

import (
	"net/http"
	"strconv"

	"github.com/ViniciusBoroto/cabeleleila_leila/internal/models"
	"github.com/ViniciusBoroto/cabeleleila_leila/internal/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserRequest struct {
	Email    string          `json:"email" binding:"required,email"`
	Password string          `json:"password" binding:"required,min=6"`
	Name     string          `json:"name" binding:"required"`
	Phone    string          `json:"phone"`
	Role     models.UserRole `json:"role" binding:"required,oneof=admin customer"`
}

type UserResponse struct {
	ID       uint            `json:"id"`
	Email    string          `json:"email"`
	Name     string          `json:"name"`
	Phone    string          `json:"phone"`
	Role     models.UserRole `json:"role"`
	IsActive bool            `json:"is_active"`
}

// GetAllUsers godoc
// @Summary      List all users (admin only)
// @Description  Retrieve all users from the system
// @Tags         admin
// @Security     Bearer
// @Produce      json
// @Success      200  {array}   UserResponse
// @Failure      403  {object}  map[string]string
// @Router       /admin/users [get]
func GetAllUsers(userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != models.RoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			return
		}

		users, err := userRepo.FindAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := make([]UserResponse, len(users))
		for i, user := range users {
			response[i] = UserResponse{
				ID:       user.ID,
				Email:    user.Email,
				Name:     user.Name,
				Phone:    user.Phone,
				Role:     user.Role,
				IsActive: user.IsActive,
			}
		}

		c.JSON(http.StatusOK, response)
	}
}

// CreateUser godoc
// @Summary      Create a new user (admin only)
// @Description  Create a new user with the provided details
// @Tags         admin
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        user  body      CreateUserRequest  true  "User data"
// @Success      201   {object}  UserResponse
// @Failure      400   {object}  map[string]string
// @Failure      403   {object}  map[string]string
// @Router       /admin/users [post]
func CreateUser(userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != models.RoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			return
		}

		var req CreateUserRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}

		user := models.User{
			Email:    req.Email,
			Password: string(hashedPassword),
			Name:     req.Name,
			Phone:    req.Phone,
			Role:     req.Role,
			IsActive: true,
		}

		created, err := userRepo.Create(user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response := UserResponse{
			ID:       created.ID,
			Email:    created.Email,
			Name:     created.Name,
			Phone:    created.Phone,
			Role:     created.Role,
			IsActive: created.IsActive,
		}

		c.JSON(http.StatusCreated, response)
	}
}

// GetUser godoc
// @Summary      Get user by ID (admin only)
// @Description  Retrieve a specific user by ID
// @Tags         admin
// @Security     Bearer
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  UserResponse
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /admin/users/{id} [get]
func GetUser(userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != models.RoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			return
		}

		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			return
		}

		user, err := userRepo.FindByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		response := UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			Name:     user.Name,
			Phone:    user.Phone,
			Role:     user.Role,
			IsActive: user.IsActive,
		}

		c.JSON(http.StatusOK, response)
	}
}

type UpdateUserRequest struct {
	Name     *string          `json:"name"`
	Phone    *string          `json:"phone"`
	Role     *models.UserRole `json:"role"`
	IsActive *bool            `json:"is_active"`
}

// UpdateUser godoc
// @Summary      Update user (admin only)
// @Description  Update a specific user's details
// @Tags         admin
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        id    path      int                  true  "User ID"
// @Param        user  body      UpdateUserRequest    true  "Updated user data"
// @Success      200   {object}  UserResponse
// @Failure      403   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Router       /admin/users/{id} [put]
func UpdateUser(userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != models.RoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			return
		}

		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			return
		}

		user, err := userRepo.FindByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		var req UpdateUserRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Name != nil {
			user.Name = *req.Name
		}
		if req.Phone != nil {
			user.Phone = *req.Phone
		}
		if req.Role != nil {
			user.Role = *req.Role
		}
		if req.IsActive != nil {
			user.IsActive = *req.IsActive
		}

		if err := userRepo.Update(user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			Name:     user.Name,
			Phone:    user.Phone,
			Role:     user.Role,
			IsActive: user.IsActive,
		}

		c.JSON(http.StatusOK, response)
	}
}

// DeleteUser godoc
// @Summary      Delete user (admin only)
// @Description  Delete a specific user
// @Tags         admin
// @Security     Bearer
// @Param        id  path  int  true  "User ID"
// @Success      204
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /admin/users/{id} [delete]
func DeleteUser(userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != models.RoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			return
		}

		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			return
		}

		// Verify user exists
		_, err = userRepo.FindByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		if err := userRepo.Delete(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}
