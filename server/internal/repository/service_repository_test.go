package repository

import (
    "testing"

    "github.com/ViniciusBoroto/cabeleleila_leila/internal/models"
    "github.com/stretchr/testify/assert"
)

type MockServiceRepositoryBehavior struct {
    services map[uint]models.Service
    nextID   uint
}

func NewMockServiceRepositoryBehavior() *MockServiceRepositoryBehavior {
    return &MockServiceRepositoryBehavior{
        services: make(map[uint]models.Service),
        nextID:   1,
    }
}

func (m *MockServiceRepositoryBehavior) Create(service models.Service) (models.Service, error) {
    service.ID = m.nextID
    m.services[service.ID] = service
    m.nextID++
    return service, nil
}

func (m *MockServiceRepositoryBehavior) FindByID(id uint) (models.Service, error) {
    if service, ok := m.services[id]; ok {
        return service, nil
    }
    return models.Service{}, assert.AnError
}

func (m *MockServiceRepositoryBehavior) FindAll() ([]models.Service, error) {
    services := make([]models.Service, 0, len(m.services))
    for _, service := range m.services {
        services = append(services, service)
    }
    return services, nil
}

func (m *MockServiceRepositoryBehavior) Update(service models.Service) error {
    if _, ok := m.services[service.ID]; !ok {
        return assert.AnError
    }
    m.services[service.ID] = service
    return nil
}

func (m *MockServiceRepositoryBehavior) Delete(id uint) error {
    if _, ok := m.services[id]; !ok {
        return assert.AnError
    }
    delete(m.services, id)
    return nil
}

func (m *MockServiceRepositoryBehavior) FindByName(name string) (models.Service, error) {
    for _, service := range m.services {
        if service.Name == name {
            return service, nil
        }
    }
    return models.Service{}, assert.AnError
}

// Tests
func TestServiceRepository_Create(t *testing.T) {
    repo := NewMockServiceRepositoryBehavior()

    service := models.Service{
        Name:            "Corte de Cabelo",
        Price:           50.0,
        DurationMinutes: 30,
    }

    created, err := repo.Create(service)
    assert.NoError(t, err)
    assert.NotZero(t, created.ID)
    assert.Equal(t, "Corte de Cabelo", created.Name)
}

func TestServiceRepository_FindByID(t *testing.T) {
    repo := NewMockServiceRepositoryBehavior()

    service := models.Service{
        Name:            "Escova",
        Price:           40.0,
        DurationMinutes: 45,
    }

    created, _ := repo.Create(service)
    found, err := repo.FindByID(created.ID)
    assert.NoError(t, err)
    assert.Equal(t, created.ID, found.ID)
    assert.Equal(t, "Escova", found.Name)
}

func TestServiceRepository_FindByID_NotFound(t *testing.T) {
    repo := NewMockServiceRepositoryBehavior()

    _, err := repo.FindByID(999)
    assert.Error(t, err)
}

func TestServiceRepository_FindAll(t *testing.T) {
    repo := NewMockServiceRepositoryBehavior()

    service1 := models.Service{Name: "Corte", Price: 50.0, DurationMinutes: 30}
    service2 := models.Service{Name: "Coloração", Price: 100.0, DurationMinutes: 60}

    repo.Create(service1)
    repo.Create(service2)

    services, err := repo.FindAll()
    assert.NoError(t, err)
    assert.Len(t, services, 2)
}

func TestServiceRepository_Update(t *testing.T) {
    repo := NewMockServiceRepositoryBehavior()

    service := models.Service{
        Name:            "Corte",
        Price:           50.0,
        DurationMinutes: 30,
    }

    created, _ := repo.Create(service)
    created.Price = 60.0
    created.DurationMinutes = 35

    err := repo.Update(created)
    assert.NoError(t, err)

    updated, _ := repo.FindByID(created.ID)
    assert.Equal(t, 60.0, updated.Price)
    assert.Equal(t, 35, updated.DurationMinutes)
}

func TestServiceRepository_Delete(t *testing.T) {
    repo := NewMockServiceRepositoryBehavior()

    service := models.Service{
        Name:            "Corte",
        Price:           50.0,
        DurationMinutes: 30,
    }

    created, _ := repo.Create(service)
    err := repo.Delete(created.ID)
    assert.NoError(t, err)

    _, err = repo.FindByID(created.ID)
    assert.Error(t, err)
}

func TestServiceRepository_FindByName(t *testing.T) {
    repo := NewMockServiceRepositoryBehavior()

    service := models.Service{
        Name:            "Manicure",
        Price:           30.0,
        DurationMinutes: 20,
    }

    repo.Create(service)
    found, err := repo.FindByName("Manicure")
    assert.NoError(t, err)
    assert.Equal(t, "Manicure", found.Name)
}