package repository

import (
	"errors"

	"github.com/ViniciusBoroto/cabeleleila_leila/internal/models"
	"gorm.io/gorm"
)

type sqlUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &sqlUserRepository{db: db}
}

func (r *sqlUserRepository) Create(user models.User) (models.User, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *sqlUserRepository) FindByID(id uint) (models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}
	return user, nil
}

func (r *sqlUserRepository) FindByEmail(email string) (models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}
	return user, nil
}

func (r *sqlUserRepository) Update(user models.User) error {
	return r.db.Save(&user).Error
}

func (r *sqlUserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *sqlUserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *sqlUserRepository) FindByRole(role models.UserRole) ([]models.User, error) {
	var users []models.User
	if err := r.db.Where("role = ?", role).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
