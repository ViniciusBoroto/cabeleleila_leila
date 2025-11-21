package repository

import (
    "github.com/ViniciusBoroto/cabeleleila_leila/internal/models"
)

type ServiceRepository interface {
    Create(service models.Service) (models.Service, error)
    FindByID(id uint) (models.Service, error)
    FindAll() ([]models.Service, error)
    Update(service models.Service) error
    Delete(id uint) error
    FindByName(name string) (models.Service, error)
}