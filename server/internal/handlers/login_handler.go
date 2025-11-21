package handlers

import (
	"net/http"

	"github.com/ViniciusBoroto/cabeleleila_leila/internal/models"
	"github.com/ViniciusBoroto/cabeleleila_leila/internal/repository"
	"github.com/ViniciusBoroto/cabeleleila_leila/internal/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	authSvc service.AuthService
	userRepo repository.UserRepository
}

func NewAuthHandler(authSvc service.AuthService, userRepo repository.UserRepository) *AuthHandler {
	return &AuthHandler{
		authSvc: authSvc,
		userRepo: userRepo,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  UserInfo    `json:"user"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
}

type RegisterResponse struct {
	Token string      `json:"token"`
	User  UserInfo    `json:"user"`
}

type UserInfo struct {
	ID    uint                `json:"id"`
	Email string              `json:"email"`
	Name  string              `json:"name"`
	Role  models.UserRole     `json:"role"`
}

// Login godoc
// @Summary      Realiza login e retorna JWT token
// @Description  Autentica o usuário e retorna um token JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      LoginRequest  true  "Credenciais de login"
// @Success      200          {object}  LoginResponse
// @Failure      400          {object}  map[string]string
// @Failure      401          {object}  map[string]string
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepo.FindByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// Check if user is active
	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user account is inactive"})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// Generate token
	token, err := h.authSvc.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User: UserInfo{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Role:  user.Role,
		},
	})
}

// Register godoc
// @Summary      Registra um novo cliente
// @Description  Cria uma nova conta de cliente e retorna um token JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      RegisterRequest  true  "Dados de registro"
// @Success      201          {object}  RegisterResponse
// @Failure      400          {object}  map[string]string
// @Failure      409          {object}  map[string]string
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	_, err := h.userRepo.FindByEmail(req.Email)
	if err == nil {
		// User found - email already exists
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process password"})
		return
	}

	// Create new customer user
	user := models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
		Phone:    req.Phone,
		Role:     models.RoleCustomer, // New registrations are always customers
		IsActive: true,
	}

	created, err := h.userRepo.Create(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate token
	token, err := h.authSvc.GenerateToken(created)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, RegisterResponse{
		Token: token,
		User: UserInfo{
			ID:    created.ID,
			Email: created.Email,
			Name:  created.Name,
			Role:  created.Role,
		},
	})
}
