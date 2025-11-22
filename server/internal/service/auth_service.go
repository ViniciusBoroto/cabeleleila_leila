package service

import (
	"errors"
	"time"

	"github.com/ViniciusBoroto/cabeleleila_leila/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	GenerateToken(user models.User) (string, error)
	ValidateToken(tokenString string) (*CustomClaims, error)
	ValidateTokenWithRole(tokenString string, allowedRoles ...models.UserRole) (*CustomClaims, error)
}

type authService struct {
	secret string
}

type CustomClaims struct {
	UserID     uint            `json:"user_id"`
	Email      string          `json:"email"`
	Role       models.UserRole `json:"role"`
	CustomerID *uint           `json:"customer_id,omitempty"`
	jwt.RegisteredClaims
}

func NewAuthService(secret string) AuthService {
	return &authService{secret: secret}
}

func (s *authService) GenerateToken(user models.User) (string, error) {
	if user.ID == 0 {
		return "", errors.New("invalid user: missing ID")
	}

	claims := CustomClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *authService) ValidateToken(tokenString string) (*CustomClaims, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (s *authService) ValidateTokenWithRole(tokenString string, allowedRoles ...models.UserRole) (*CustomClaims, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	// If no specific roles required, allow any authenticated user
	if len(allowedRoles) == 0 {
		return claims, nil
	}

	for _, role := range allowedRoles {
		if claims.Role == role {
			return claims, nil
		}
	}

	return nil, errors.New("user role not authorized for this action")
}
