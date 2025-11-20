package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ViniciusBoroto/cabeleleila_leila/internal/models"
	"github.com/ViniciusBoroto/cabeleleila_leila/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRouter(t *testing.T) *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestJWTAuthMiddleware_ValidToken(t *testing.T) {
	authSvc := service.NewAuthService("test-secret")
	router := setupTestRouter(t)
	user := models.User{
		ID:    1,
		Email: "test@example.com",
		Role:  models.RoleCustomer,
		Name:  "Test User",
	}

	token, err := authSvc.GenerateToken(user)
	require.NoError(t, err)

	router.GET("/protected", JWTAuthMiddleware(authSvc), func(c *gin.Context) {
		userID, exists := c.Get("userID")
		assert.True(t, exists)
		assert.Equal(t, uint(1), userID)

		email, exists := c.Get("email")
		assert.True(t, exists)
		assert.Equal(t, "test@example.com", email)

		role, exists := c.Get("role")
		assert.True(t, exists)
		assert.Equal(t, models.RoleCustomer, role)

		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJWTAuthMiddleware_MissingHeader(t *testing.T) {
	authSvc := service.NewAuthService("test-secret")
	router := setupTestRouter(t)

	router.GET("/protected", JWTAuthMiddleware(authSvc), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	// No Authorization header
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTAuthMiddleware_InvalidFormat(t *testing.T) {
	authSvc := service.NewAuthService("test-secret")
	router := setupTestRouter(t)

	router.GET("/protected", JWTAuthMiddleware(authSvc), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "InvalidToken")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTAuthMiddleware_InvalidToken(t *testing.T) {
	authSvc := service.NewAuthService("test-secret")
	router := setupTestRouter(t)

	router.GET("/protected", JWTAuthMiddleware(authSvc), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTAuthMiddleware_WrongSecret(t *testing.T) {
	authSvc1 := service.NewAuthService("secret-1")
	authSvc2 := service.NewAuthService("secret-2")
	router := setupTestRouter(t)

	user := models.User{
		ID:    1,
		Email: "test@example.com",
		Role:  models.RoleCustomer,
	}

	token, err := authSvc1.GenerateToken(user)
	require.NoError(t, err)

	router.GET("/protected", JWTAuthMiddleware(authSvc2), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRequireRole_AllowedRole(t *testing.T) {
	authSvc := service.NewAuthService("test-secret")
	router := setupTestRouter(t)

	user := models.User{
		ID:    1,
		Email: "admin@example.com",
		Role:  models.RoleAdmin,
	}

	token, err := authSvc.GenerateToken(user)
	require.NoError(t, err)

	router.GET("/admin", RequireRole(authSvc, models.RoleAdmin), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	req := httptest.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequireRole_UnauthorizedRole(t *testing.T) {
	authSvc := service.NewAuthService("test-secret")
	router := setupTestRouter(t)

	user := models.User{
		ID:    1,
		Email: "customer@example.com",
		Role:  models.RoleCustomer,
	}

	token, err := authSvc.GenerateToken(user)
	require.NoError(t, err)

	router.GET("/admin", RequireRole(authSvc, models.RoleAdmin), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	req := httptest.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRequireRole_MultipleRoles(t *testing.T) {
	authSvc := service.NewAuthService("test-secret")
	router := setupTestRouter(t)

	user := models.User{
		ID:    1,
		Email: "customer@example.com",
		Role:  models.RoleCustomer,
	}

	token, err := authSvc.GenerateToken(user)
	require.NoError(t, err)

	router.GET("/protected", RequireRole(authSvc, models.RoleAdmin, models.RoleCustomer), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequireRole_MissingHeader(t *testing.T) {
	authSvc := service.NewAuthService("test-secret")
	router := setupTestRouter(t)

	router.GET("/admin", RequireRole(authSvc, models.RoleAdmin), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	req := httptest.NewRequest("GET", "/admin", nil)
	// No Authorization header
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
