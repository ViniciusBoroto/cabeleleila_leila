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
	RefreshToken(tokenString string) (string, error)
}

type authService struct {
	secret string
}

type CustomClaims struct {
	UserID uint            `json:"user_id"`
	Email  string          `json:"email"`
	Role   models.UserRole `json:"role"`
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

func (s *authService) RefreshToken(tokenString string) (string, error) {
	// Parse token without validating expiration
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secret), nil
	})

	// Check for errors other than expiration
	if err != nil {
		if !errors.Is(err, jwt.ErrTokenExpired) {
			return "", err
		}
	}

	// Verify token signature is valid (even if expired)
	if token == nil || !token.Valid && !errors.Is(err, jwt.ErrTokenExpired) {
		return "", errors.New("invalid token")
	}

	// Check if token is too old to refresh (e.g., expired more than 7 days ago)
	if claims.ExpiresAt != nil {
		expirationTime := claims.ExpiresAt.Time
		if time.Since(expirationTime) > 7*24*time.Hour {
			return "", errors.New("token too old to refresh")
		}
	}

	// Generate new token with same user data
	newClaims := CustomClaims{
		UserID: claims.UserID,
		Email:  claims.Email,
		Role:   claims.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	return newToken.SignedString([]byte(s.secret))
}
