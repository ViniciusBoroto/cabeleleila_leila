package repository

import (
    "errors"

    "github.com/ViniciusBoroto/cabeleleila_leila/internal/models"
    "gorm.io/gorm"
)

type sqlServiceRepository struct {
    db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) ServiceRepository {
    return &sqlServiceRepository{db: db}
}

func (r *sqlServiceRepository) Create(service models.Service) (models.Service, error) {
    if err := r.db.Create(&service).Error; err != nil {
        return models.Service{}, err
    }
    return service, nil
}

func (r *sqlServiceRepository) FindByID(id uint) (models.Service, error) {
    var service models.Service
    if err := r.db.First(&service, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return models.Service{}, errors.New("service not found")
        }
        return models.Service{}, err
    }
    return service, nil
}

func (r *sqlServiceRepository) FindAll() ([]models.Service, error) {
    var services []models.Service
    if err := r.db.Find(&services).Error; err != nil {
        return nil, err
    }
    return services, nil
}

func (r *sqlServiceRepository) Update(service models.Service) error {
    if service.ID == 0 {
        return errors.New("service ID is required for update")
    }
    return r.db.Model(&service).Updates(service).Error
}

func (r *sqlServiceRepository) Delete(id uint) error {
    if id == 0 {
        return errors.New("invalid service ID")
    }
    return r.db.Delete(&models.Service{}, id).Error
}

func (r *sqlServiceRepository) FindByName(name string) (models.Service, error) {
    var service models.Service
    if err := r.db.Where("name = ?", name).First(&service).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return models.Service{}, errors.New("service not found")
        }
        return models.Service{}, err
    }
    return service, nil
}