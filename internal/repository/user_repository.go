package repository

import (
	"github.com/ViniciusBoroto/cabeleleila_leila/internal/models"
)

type UserRepository interface {
	Create(user models.User) (models.User, error)
	FindByID(id uint) (models.User, error)
	FindByEmail(email string) (models.User, error)
	Update(user models.User) error
	Delete(id uint) error
	FindAll() ([]models.User, error)
	FindByRole(role models.UserRole) ([]models.User, error)
}
