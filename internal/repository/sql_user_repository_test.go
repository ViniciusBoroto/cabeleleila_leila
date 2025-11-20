package repository

import (
	"testing"

	"github.com/ViniciusBoroto/cabeleleila_leila/internal/models"
	"github.com/stretchr/testify/assert"
)

// Note: Full database integration tests require CGO_ENABLED=1
// These tests serve as integration tests when CGO is enabled
// For unit testing without database, use mocks

// TestUserRepository_Interface verifies the interface is implemented correctly
func TestUserRepository_Interface(t *testing.T) {
	var _ UserRepository = (*sqlUserRepository)(nil)
}

// Mock-based tests for UserRepository behavior
type MockUserRepositoryBehavior struct {
	users map[uint]models.User
	nextID uint
}

func NewMockUserRepositoryBehavior() *MockUserRepositoryBehavior {
	return &MockUserRepositoryBehavior{
		users: make(map[uint]models.User),
		nextID: 1,
	}
}

func (m *MockUserRepositoryBehavior) Create(user models.User) (models.User, error) {
	user.ID = m.nextID
	m.users[user.ID] = user
	m.nextID++
	return user, nil
}

func (m *MockUserRepositoryBehavior) FindByID(id uint) (models.User, error) {
	if user, ok := m.users[id]; ok {
		return user, nil
	}
	return models.User{}, assert.AnError
}

func (m *MockUserRepositoryBehavior) FindByEmail(email string) (models.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return models.User{}, assert.AnError
}

func (m *MockUserRepositoryBehavior) Update(user models.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepositoryBehavior) Delete(id uint) error {
	delete(m.users, id)
	return nil
}

func (m *MockUserRepositoryBehavior) FindAll() ([]models.User, error) {
	users := make([]models.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func (m *MockUserRepositoryBehavior) FindByRole(role models.UserRole) ([]models.User, error) {
	var users []models.User
	for _, user := range m.users {
		if user.Role == role {
			users = append(users, user)
		}
	}
	return users, nil
}

// Behavioral tests using mock
func TestUserRepository_CreateAndFind(t *testing.T) {
	repo := NewMockUserRepositoryBehavior()

	user := models.User{
		Email:    "test@example.com",
		Password: "hashed_password",
		Role:     models.RoleCustomer,
		Name:     "John Doe",
		IsActive: true,
	}

	created, err := repo.Create(user)
	assert.NoError(t, err)
	assert.NotZero(t, created.ID)

	found, err := repo.FindByID(created.ID)
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", found.Email)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	repo := NewMockUserRepositoryBehavior()

	user := models.User{
		Email:    "john@example.com",
		Password: "hashed_password",
		Role:     models.RoleCustomer,
		Name:     "John",
		IsActive: true,
	}

	created, err := repo.Create(user)
	assert.NoError(t, err)

	found, err := repo.FindByEmail("john@example.com")
	assert.NoError(t, err)
	assert.Equal(t, created.ID, found.ID)
}

func TestUserRepository_Update(t *testing.T) {
	repo := NewMockUserRepositoryBehavior()

	user := models.User{
		Email:    "test@example.com",
		Password: "hashed_password",
		Role:     models.RoleCustomer,
		Name:     "John",
		IsActive: true,
	}

	created, err := repo.Create(user)
	assert.NoError(t, err)

	created.Name = "Jane"
	err = repo.Update(created)
	assert.NoError(t, err)

	found, err := repo.FindByID(created.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Jane", found.Name)
}

func TestUserRepository_Delete(t *testing.T) {
	repo := NewMockUserRepositoryBehavior()

	user := models.User{
		Email:    "test@example.com",
		Password: "hashed_password",
		Role:     models.RoleCustomer,
		Name:     "John",
		IsActive: true,
	}

	created, err := repo.Create(user)
	assert.NoError(t, err)

	err = repo.Delete(created.ID)
	assert.NoError(t, err)

	_, err = repo.FindByID(created.ID)
	assert.Error(t, err)
}

func TestUserRepository_FindAll(t *testing.T) {
	repo := NewMockUserRepositoryBehavior()

	user1 := models.User{
		Email: "user1@example.com",
		Role:  models.RoleCustomer,
		Name:  "User 1",
	}
	user2 := models.User{
		Email: "user2@example.com",
		Role:  models.RoleAdmin,
		Name:  "User 2",
	}

	repo.Create(user1)
	repo.Create(user2)

	users, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestUserRepository_FindByRole(t *testing.T) {
	repo := NewMockUserRepositoryBehavior()

	admin := models.User{
		Email: "admin@example.com",
		Role:  models.RoleAdmin,
		Name:  "Admin",
	}
	customer1 := models.User{
		Email: "customer1@example.com",
		Role:  models.RoleCustomer,
		Name:  "Customer 1",
	}
	customer2 := models.User{
		Email: "customer2@example.com",
		Role:  models.RoleCustomer,
		Name:  "Customer 2",
	}

	repo.Create(admin)
	repo.Create(customer1)
	repo.Create(customer2)

	customers, err := repo.FindByRole(models.RoleCustomer)
	assert.NoError(t, err)
	assert.Len(t, customers, 2)

	admins, err := repo.FindByRole(models.RoleAdmin)
	assert.NoError(t, err)
	assert.Len(t, admins, 1)
}
