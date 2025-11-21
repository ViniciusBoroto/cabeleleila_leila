package service

import (
    "errors"

    "github.com/ViniciusBoroto/cabeleleila_leila/internal/models"
    "github.com/ViniciusBoroto/cabeleleila_leila/internal/repository"
)

type ServiceService interface {
    CreateService(service models.Service) (models.Service, error)
    GetService(id uint) (models.Service, error)
    ListServices() ([]models.Service, error)
    UpdateService(service models.Service) (models.Service, error)
    DeleteService(id uint) error
}

type serviceService struct {
    repo repository.ServiceRepository
}

func NewServiceService(repo repository.ServiceRepository) ServiceService {
    return &serviceService{repo: repo}
}

func (s *serviceService) CreateService(service models.Service) (models.Service, error) {
    if service.Name == "" {
        return models.Service{}, errors.New("service name is required")
    }
    if service.Price < 0 {
        return models.Service{}, errors.New("service price cannot be negative")
    }
    if service.DurationMinutes <= 0 {
        return models.Service{}, errors.New("service duration must be greater than 0")
    }
    return s.repo.Create(service)
}

func (s *serviceService) GetService(id uint) (models.Service, error) {
    if id == 0 {
        return models.Service{}, errors.New("invalid service ID")
    }
    return s.repo.FindByID(id)
}

func (s *serviceService) ListServices() ([]models.Service, error) {
    return s.repo.FindAll()
}

func (s *serviceService) UpdateService(service models.Service) (models.Service, error) {
    if service.ID == 0 {
        return models.Service{}, errors.New("service ID is required")
    }
    if service.Name == "" {
        return models.Service{}, errors.New("service name is required")
    }
    if service.Price < 0 {
        return models.Service{}, errors.New("service price cannot be negative")
    }
    if service.DurationMinutes <= 0 {
        return models.Service{}, errors.New("service duration must be greater than 0")
    }

    if err := s.repo.Update(service); err != nil {
        return models.Service{}, err
    }
    return s.repo.FindByID(service.ID)
}

func (s *serviceService) DeleteService(id uint) error {
    if id == 0 {
        return errors.New("invalid service ID")
    }
    _, err := s.repo.FindByID(id)
    if err != nil {
        return err
    }
    return s.repo.Delete(id)
}