package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ViniciusBoroto/cabeleleila_leila/internal/mocks"
	"github.com/ViniciusBoroto/cabeleleila_leila/internal/models"
	"github.com/ViniciusBoroto/cabeleleila_leila/internal/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

func TestLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	authSvc := service.NewAuthService("test-secret")
	handler := NewAuthHandler(authSvc, mockUserRepo)

	password := "password123"
	hashedPassword := hashPassword(password)

	user := models.User{
		ID:       1,
		Email:    "test@example.com",
		Password: hashedPassword,
		Role:     models.RoleCustomer,
		Name:     "Test User",
		IsActive: true,
	}

	mockUserRepo.EXPECT().
		FindByEmail("test@example.com").
		Return(user, nil)

	router := setupTestRouter(t)
	router.POST("/login", handler.Login)

	req := LoginRequest{
		Email:    "test@example.com",
		Password: password,
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var response LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Token)
	assert.Equal(t, "test@example.com", response.User.Email)
	assert.Equal(t, models.RoleCustomer, response.User.Role)
}

func TestLogin_InvalidEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	authSvc := service.NewAuthService("test-secret")
	handler := NewAuthHandler(authSvc, mockUserRepo)

	router := setupTestRouter(t)
	router.POST("/login", handler.Login)

	req := LoginRequest{
		Email:    "invalid",
		Password: "password123",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_MissingEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	authSvc := service.NewAuthService("test-secret")
	handler := NewAuthHandler(authSvc, mockUserRepo)

	router := setupTestRouter(t)
	router.POST("/login", handler.Login)

	req := LoginRequest{
		Email:    "",
		Password: "password123",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_ShortPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	authSvc := service.NewAuthService("test-secret")
	handler := NewAuthHandler(authSvc, mockUserRepo)

	router := setupTestRouter(t)
	router.POST("/login", handler.Login)

	req := LoginRequest{
		Email:    "test@example.com",
		Password: "pass", // Too short
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	authSvc := service.NewAuthService("test-secret")
	handler := NewAuthHandler(authSvc, mockUserRepo)

	mockUserRepo.EXPECT().
		FindByEmail("nonexistent@example.com").
		Return(models.User{}, error_NotFound())

	router := setupTestRouter(t)
	router.POST("/login", handler.Login)

	req := LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func error_NotFound() error {
	return assert.AnError
}

func TestLogin_InactiveUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	authSvc := service.NewAuthService("test-secret")
	handler := NewAuthHandler(authSvc, mockUserRepo)

	user := models.User{
		ID:       1,
		Email:    "inactive@example.com",
		Password: hashPassword("password123"),
		Role:     models.RoleCustomer,
		Name:     "Inactive User",
		IsActive: false,
	}

	mockUserRepo.EXPECT().
		FindByEmail("inactive@example.com").
		Return(user, nil)

	router := setupTestRouter(t)
	router.POST("/login", handler.Login)

	req := LoginRequest{
		Email:    "inactive@example.com",
		Password: "password123",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogin_WrongPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	authSvc := service.NewAuthService("test-secret")
	handler := NewAuthHandler(authSvc, mockUserRepo)

	user := models.User{
		ID:       1,
		Email:    "test@example.com",
		Password: hashPassword("correctPassword123"),
		Role:     models.RoleCustomer,
		Name:     "Test User",
		IsActive: true,
	}

	mockUserRepo.EXPECT().
		FindByEmail("test@example.com").
		Return(user, nil)

	router := setupTestRouter(t)
	router.POST("/login", handler.Login)

	req := LoginRequest{
		Email:    "test@example.com",
		Password: "wrongPassword123",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogin_AdminUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	authSvc := service.NewAuthService("test-secret")
	handler := NewAuthHandler(authSvc, mockUserRepo)

	password := "adminPass123"
	hashedPassword := hashPassword(password)

	user := models.User{
		ID:       2,
		Email:    "admin@example.com",
		Password: hashedPassword,
		Role:     models.RoleAdmin,
		Name:     "Admin User",
		IsActive: true,
	}

	mockUserRepo.EXPECT().
		FindByEmail("admin@example.com").
		Return(user, nil)

	router := setupTestRouter(t)
	router.POST("/login", handler.Login)

	req := LoginRequest{
		Email:    "admin@example.com",
		Password: password,
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var response LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, models.RoleAdmin, response.User.Role)
}

func TestRegister_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	authSvc := service.NewAuthService("test-secret")
	handler := NewAuthHandler(authSvc, mockUserRepo)

	// First call: FindByEmail returns error (user doesn't exist)
	mockUserRepo.EXPECT().
		FindByEmail("newcustomer@example.com").
		Return(models.User{}, assert.AnError)

	// Second call: Create user
	createdUser := models.User{
		ID:       1,
		Email:    "newcustomer@example.com",
		Password: hashPassword("password123"),
		Role:     models.RoleCustomer,
		Name:     "New Customer",
		Phone:    "123456789",
		IsActive: true,
	}

	mockUserRepo.EXPECT().
		Create(gomock.Any()).
		Return(createdUser, nil)

	router := setupTestRouter(t)
	router.POST("/register", handler.Register)

	req := RegisterRequest{
		Email:    "newcustomer@example.com",
		Password: "password123",
		Name:     "New Customer",
		Phone:    "123456789",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response RegisterResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Token)
	assert.Equal(t, "newcustomer@example.com", response.User.Email)
	assert.Equal(t, models.RoleCustomer, response.User.Role)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	authSvc := service.NewAuthService("test-secret")
	handler := NewAuthHandler(authSvc, mockUserRepo)

	existingUser := models.User{
		ID:    1,
		Email: "existing@example.com",
		Role:  models.RoleCustomer,
	}

	mockUserRepo.EXPECT().
		FindByEmail("existing@example.com").
		Return(existingUser, nil)

	router := setupTestRouter(t)
	router.POST("/register", handler.Register)

	req := RegisterRequest{
		Email:    "existing@example.com",
		Password: "password123",
		Name:     "Someone Else",
		Phone:    "987654321",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestRegister_InvalidEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	authSvc := service.NewAuthService("test-secret")
	handler := NewAuthHandler(authSvc, mockUserRepo)

	router := setupTestRouter(t)
	router.POST("/register", handler.Register)

	req := RegisterRequest{
		Email:    "invalid-email",
		Password: "password123",
		Name:     "Test",
		Phone:    "123456789",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegister_MissingFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	authSvc := service.NewAuthService("test-secret")
	handler := NewAuthHandler(authSvc, mockUserRepo)

	router := setupTestRouter(t)
	router.POST("/register", handler.Register)

	req := RegisterRequest{
		Email: "test@example.com",
		// Missing password, name, phone
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegister_ShortPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	authSvc := service.NewAuthService("test-secret")
	handler := NewAuthHandler(authSvc, mockUserRepo)

	router := setupTestRouter(t)
	router.POST("/register", handler.Register)

	req := RegisterRequest{
		Email:    "test@example.com",
		Password: "pass", // Too short
		Name:     "Test",
		Phone:    "123456789",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegister_AlwaysCreatesCustomer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	authSvc := service.NewAuthService("test-secret")
	handler := NewAuthHandler(authSvc, mockUserRepo)

	mockUserRepo.EXPECT().
		FindByEmail("customer@example.com").
		Return(models.User{}, assert.AnError)

	// Create user returns customer role
	mockUserRepo.EXPECT().
		Create(gomock.Any()).
		Return(models.User{
			ID:       1,
			Email:    "customer@example.com",
			Role:     models.RoleCustomer,
			Name:     "New Customer",
			IsActive: true,
		}, nil)

	router := setupTestRouter(t)
	router.POST("/register", handler.Register)

	req := RegisterRequest{
		Email:    "customer@example.com",
		Password: "password123",
		Name:     "New Customer",
		Phone:    "123456789",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response RegisterResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, models.RoleCustomer, response.User.Role)
}
