package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/ViniciusBoroto/cabeleleila_leila/internal/models"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("file:memdb%d?mode=memory&cache=shared", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err, "failed to create in-memory database")

	// Migrate the schema
	err = db.AutoMigrate(
		&models.User{},
		&models.Customer{},
		&models.Service{},
		&models.Appointment{},
	)
	require.NoError(t, err, "failed to migrate schema")

	return db
}

func createTestUser(t *testing.T, db *gorm.DB, email string) models.User {
	user := models.User{
		Email:    email,
		Password: "hashed_password",
		Role:     models.RoleCustomer,
		Name:     "Test User",
		IsActive: true,
	}
	result := db.Create(&user)
	require.NoError(t, result.Error, "failed to create test user")
	return user
}

func createTestCustomer(t *testing.T, db *gorm.DB, userID uint) models.Customer {
	customer := models.Customer{
		UserID:   userID,
		IsActive: true,
	}
	result := db.Create(&customer)
	require.NoError(t, result.Error, "failed to create test customer")
	return customer
}

func createTestService(t *testing.T, db *gorm.DB, name string, price float64, duration int) models.Service {
	service := models.Service{
		Name:            name,
		Price:           price,
		DurationMinutes: duration,
	}
	result := db.Create(&service)
	require.NoError(t, result.Error, "failed to create test service")
	return service
}

func createTestAppointment(t *testing.T, db *gorm.DB, customerID uint, services []models.Service, date time.Time) models.Appointment {
	appointment := models.Appointment{
		CustomerID: customerID,
		Services:   services,
		Date:       date,
		Status:     models.StatusPending,
	}
	result := db.Create(&appointment)
	require.NoError(t, result.Error, "failed to create test appointment")
	return appointment
}

// TestAppointmentRepository_Interface verifies the interface is implemented correctly
func TestAppointmentRepository_Interface(t *testing.T) {
	var _ AppointmentRepository = (*sqlAppointmentRepo)(nil)
}

// TestAppointmentRepository_Create tests creating an appointment
func TestAppointmentRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAppointmentRepository(db)

	user := createTestUser(t, db, "customer@example.com")
	customer := createTestCustomer(t, db, user.ID)
	service := createTestService(t, db, "Haircut", 50.00, 30)

	appointment := models.Appointment{
		CustomerID: customer.ID,
		Services:   []models.Service{service},
		Date:       time.Now().Add(24 * time.Hour),
		Status:     models.StatusPending,
	}

	created, err := repo.Create(appointment)
	assert.NoError(t, err)
	assert.NotZero(t, created.ID)
	assert.Equal(t, customer.ID, created.CustomerID)
	assert.Equal(t, models.StatusPending, created.Status)
}

// TestAppointmentRepository_FindByID tests retrieving an appointment by ID
func TestAppointmentRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAppointmentRepository(db)

	user := createTestUser(t, db, "customer@example.com")
	customer := createTestCustomer(t, db, user.ID)
	service := createTestService(t, db, "Haircut", 50.00, 30)

	tomorrow := time.Now().Add(24 * time.Hour)
	created, err := repo.Create(models.Appointment{
		CustomerID: customer.ID,
		Services:   []models.Service{service},
		Date:       tomorrow,
		Status:     models.StatusPending,
	})
	require.NoError(t, err)

	found, err := repo.FindByID(created.ID)
	assert.NoError(t, err)
	assert.Equal(t, created.ID, found.ID)
	assert.Equal(t, customer.ID, found.CustomerID)
	assert.Equal(t, models.StatusPending, found.Status)
	// Verify customer is preloaded
	assert.Equal(t, customer.ID, found.Customer.ID)
	// Verify services are preloaded
	assert.Len(t, found.Services, 1)
	assert.Equal(t, service.ID, found.Services[0].ID)
	assert.Equal(t, "Haircut", found.Services[0].Name)
}

// TestAppointmentRepository_FindByID_NotFound tests FindByID with non-existent ID
func TestAppointmentRepository_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAppointmentRepository(db)

	_, err := repo.FindByID(9999)
	assert.Error(t, err)
}

// TestAppointmentRepository_Update tests updating an appointment
func TestAppointmentRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAppointmentRepository(db)

	user := createTestUser(t, db, "customer@example.com")
	customer := createTestCustomer(t, db, user.ID)
	service := createTestService(t, db, "Haircut", 50.00, 30)

	tomorrow := time.Now().Add(24 * time.Hour)
	created, err := repo.Create(models.Appointment{
		CustomerID: customer.ID,
		Services:   []models.Service{service},
		Date:       tomorrow,
		Status:     models.StatusPending,
	})
	require.NoError(t, err)

	// Update status
	created.Status = models.StatusConfirmed
	err = repo.Update(created)
	assert.NoError(t, err)

	// Verify update
	found, err := repo.FindByID(created.ID)
	assert.NoError(t, err)
	assert.Equal(t, models.StatusConfirmed, found.Status)
}

// TestAppointmentRepository_Update_NonExistent tests updating non-existent appointment
func TestAppointmentRepository_Update_NonExistent(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAppointmentRepository(db)

	nonExistent := models.Appointment{
		ID:         9999,
		CustomerID: 1,
		Status:     models.StatusConfirmed,
	}

	// Update should not fail but also won't update anything
	err := repo.Update(nonExistent)
	assert.NoError(t, err)
}

// TestAppointmentRepository_FindCustomerAppointmentsInWeek tests filtering by customer and date range
func TestAppointmentRepository_FindCustomerAppointmentsInWeek(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAppointmentRepository(db)

	user1 := createTestUser(t, db, "customer1@example.com")
	customer1 := createTestCustomer(t, db, user1.ID)

	user2 := createTestUser(t, db, "customer2@example.com")
	customer2 := createTestCustomer(t, db, user2.ID)

	service := createTestService(t, db, "Haircut", 50.00, 30)

	// Create appointments for customer1 in a specific week
	baseTime := time.Date(2025, 11, 20, 12, 0, 0, 0, time.UTC) // Thursday
	weekStart := baseTime
	dayAfter := baseTime.Add(24 * time.Hour)      // Friday
	dayAfter2 := baseTime.Add(48 * time.Hour)     // Saturday
	dayAfter8 := baseTime.Add(8 * 24 * time.Hour) // next Thursday

	repo.Create(models.Appointment{
		CustomerID: customer1.ID,
		Services:   []models.Service{service},
		Date:       dayAfter,
		Status:     models.StatusPending,
	})
	repo.Create(models.Appointment{
		CustomerID: customer1.ID,
		Services:   []models.Service{service},
		Date:       dayAfter2,
		Status:     models.StatusPending,
	})
	repo.Create(models.Appointment{
		CustomerID: customer1.ID,
		Services:   []models.Service{service},
		Date:       dayAfter8,
		Status:     models.StatusPending,
	})

	// Create appointment for customer2 in the period
	repo.Create(models.Appointment{
		CustomerID: customer2.ID,
		Services:   []models.Service{service},
		Date:       dayAfter,
		Status:     models.StatusPending,
	})

	// Test: Get customer1's appointments for this week (Thursday to Sunday = 4 days)
	weekEnd := weekStart.Add(4 * 24 * time.Hour)
	appointments, err := repo.FindCustomerAppointmentsInWeek(customer1.ID, weekStart, weekEnd)
	assert.NoError(t, err)
	assert.Len(t, appointments, 2) // Friday and Saturday

	// Verify all appointments have services preloaded
	for _, apt := range appointments {
		assert.Len(t, apt.Services, 1)
		assert.Equal(t, service.ID, apt.Services[0].ID)
	}
}

// TestAppointmentRepository_FindCustomerAppointmentsInWeek_NoResults tests when no appointments exist
func TestAppointmentRepository_FindCustomerAppointmentsInWeek_NoResults(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAppointmentRepository(db)

	user := createTestUser(t, db, "customer@example.com")
	customer := createTestCustomer(t, db, user.ID)

	now := time.Now()
	weekStart := now.AddDate(0, 0, -int(now.Weekday())+1)
	weekEnd := weekStart.AddDate(0, 0, 6).Add(24 * time.Hour)

	appointments, err := repo.FindCustomerAppointmentsInWeek(customer.ID, weekStart, weekEnd)
	assert.NoError(t, err)
	assert.Len(t, appointments, 0)
}

// TestAppointmentRepository_ListByPeriod tests listing appointments by date range
func TestAppointmentRepository_ListByPeriod(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAppointmentRepository(db)

	user1 := createTestUser(t, db, "customer1@example.com")
	customer1 := createTestCustomer(t, db, user1.ID)

	user2 := createTestUser(t, db, "customer2@example.com")
	customer2 := createTestCustomer(t, db, user2.ID)

	service := createTestService(t, db, "Haircut", 50.00, 30)

	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)
	twoDaysLater := now.Add(48 * time.Hour)
	fiveDaysLater := now.Add(120 * time.Hour)

	repo.Create(models.Appointment{
		CustomerID: customer1.ID,
		Services:   []models.Service{service},
		Date:       tomorrow,
		Status:     models.StatusPending,
	})
	repo.Create(models.Appointment{
		CustomerID: customer2.ID,
		Services:   []models.Service{service},
		Date:       twoDaysLater,
		Status:     models.StatusConfirmed,
	})
	repo.Create(models.Appointment{
		CustomerID: customer1.ID,
		Services:   []models.Service{service},
		Date:       fiveDaysLater,
		Status:     models.StatusPending,
	})

	// List appointments in a 3-day period
	appointments, err := repo.ListByPeriod(now, twoDaysLater.Add(1*time.Hour))
	assert.NoError(t, err)
	assert.Len(t, appointments, 2)

	// Verify services are preloaded
	for _, apt := range appointments {
		assert.Len(t, apt.Services, 1)
		assert.Equal(t, service.ID, apt.Services[0].ID)
		// Verify customer is preloaded
		assert.NotZero(t, apt.Customer.ID)
	}
}

// TestAppointmentRepository_ListByPeriod_NoResults tests when no appointments exist in period
func TestAppointmentRepository_ListByPeriod_NoResults(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAppointmentRepository(db)

	user := createTestUser(t, db, "customer@example.com")
	customer := createTestCustomer(t, db, user.ID)
	service := createTestService(t, db, "Haircut", 50.00, 30)

	now := time.Now()
	repo.Create(models.Appointment{
		CustomerID: customer.ID,
		Services:   []models.Service{service},
		Date:       now.Add(10 * 24 * time.Hour),
		Status:     models.StatusPending,
	})

	appointments, err := repo.ListByPeriod(now, now.Add(3*24*time.Hour))
	assert.NoError(t, err)
	assert.Len(t, appointments, 0)
}

// TestAppointmentRepository_ListAll tests listing all appointments
func TestAppointmentRepository_ListAll(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAppointmentRepository(db)

	user1 := createTestUser(t, db, "customer1@example.com")
	customer1 := createTestCustomer(t, db, user1.ID)

	user2 := createTestUser(t, db, "customer2@example.com")
	customer2 := createTestCustomer(t, db, user2.ID)

	service1 := createTestService(t, db, "Haircut", 50.00, 30)
	service2 := createTestService(t, db, "Coloring", 80.00, 60)

	now := time.Now()

	// Create multiple appointments
	repo.Create(models.Appointment{
		CustomerID: customer1.ID,
		Services:   []models.Service{service1},
		Date:       now.Add(24 * time.Hour),
		Status:     models.StatusPending,
	})
	repo.Create(models.Appointment{
		CustomerID: customer2.ID,
		Services:   []models.Service{service1, service2},
		Date:       now.Add(48 * time.Hour),
		Status:     models.StatusConfirmed,
	})
	repo.Create(models.Appointment{
		CustomerID: customer1.ID,
		Services:   []models.Service{service2},
		Date:       now.Add(72 * time.Hour),
		Status:     models.StatusDone,
	})

	appointments, err := repo.ListAll()
	assert.NoError(t, err)
	assert.Len(t, appointments, 3)

	// Verify all appointments have customers and services preloaded
	for _, apt := range appointments {
		assert.NotZero(t, apt.Customer.ID)
		assert.NotEmpty(t, apt.Services)
	}

	// Verify specific appointments
	var apt2Services models.Appointment
	for _, apt := range appointments {
		if apt.Status == models.StatusConfirmed {
			apt2Services = apt
			break
		}
	}
	assert.Len(t, apt2Services.Services, 2)
}

// TestAppointmentRepository_ListAll_Empty tests listing when no appointments exist
func TestAppointmentRepository_ListAll_Empty(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAppointmentRepository(db)

	appointments, err := repo.ListAll()
	assert.NoError(t, err)
	assert.Len(t, appointments, 0)
}

// TestAppointmentRepository_StatusTransitions tests various appointment status transitions
func TestAppointmentRepository_StatusTransitions(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAppointmentRepository(db)

	user := createTestUser(t, db, "customer@example.com")
	customer := createTestCustomer(t, db, user.ID)
	service := createTestService(t, db, "Haircut", 50.00, 30)

	created, err := repo.Create(models.Appointment{
		CustomerID: customer.ID,
		Services:   []models.Service{service},
		Date:       time.Now().Add(24 * time.Hour),
		Status:     models.StatusPending,
	})
	require.NoError(t, err)

	// Transition: Pending -> Confirmed
	created.Status = models.StatusConfirmed
	err = repo.Update(created)
	assert.NoError(t, err)

	found, _ := repo.FindByID(created.ID)
	assert.Equal(t, models.StatusConfirmed, found.Status)

	// Transition: Confirmed -> Done
	created.Status = models.StatusDone
	err = repo.Update(created)
	assert.NoError(t, err)

	found, _ = repo.FindByID(created.ID)
	assert.Equal(t, models.StatusDone, found.Status)

	// Transition: Done -> Canceled (unusual but allowed)
	created.Status = models.StatusCanceled
	err = repo.Update(created)
	assert.NoError(t, err)

	found, _ = repo.FindByID(created.ID)
	assert.Equal(t, models.StatusCanceled, found.Status)
}

// TestAppointmentRepository_MultipleServices tests appointment with multiple services
func TestAppointmentRepository_MultipleServices(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAppointmentRepository(db)

	user := createTestUser(t, db, "customer@example.com")
	customer := createTestCustomer(t, db, user.ID)

	service1 := createTestService(t, db, "Haircut", 50.00, 30)
	service2 := createTestService(t, db, "Coloring", 80.00, 60)
	service3 := createTestService(t, db, "Treatment", 40.00, 45)

	created, err := repo.Create(models.Appointment{
		CustomerID: customer.ID,
		Services:   []models.Service{service1, service2, service3},
		Date:       time.Now().Add(24 * time.Hour),
		Status:     models.StatusPending,
	})
	require.NoError(t, err)

	found, err := repo.FindByID(created.ID)
	assert.NoError(t, err)
	assert.Len(t, found.Services, 3)

	// Verify all services are correctly associated
	serviceIDs := make(map[uint]bool)
	for _, svc := range found.Services {
		serviceIDs[svc.ID] = true
	}

	assert.True(t, serviceIDs[service1.ID])
	assert.True(t, serviceIDs[service2.ID])
	assert.True(t, serviceIDs[service3.ID])
}

// TestAppointmentRepository_CustomerRelation tests customer relationship preloading
func TestAppointmentRepository_CustomerRelation(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAppointmentRepository(db)

	user := createTestUser(t, db, "customer@example.com")
	customer := createTestCustomer(t, db, user.ID)
	service := createTestService(t, db, "Haircut", 50.00, 30)

	created, err := repo.Create(models.Appointment{
		CustomerID: customer.ID,
		Services:   []models.Service{service},
		Date:       time.Now().Add(24 * time.Hour),
		Status:     models.StatusPending,
	})
	require.NoError(t, err)

	found, err := repo.FindByID(created.ID)
	assert.NoError(t, err)

	// Verify customer is preloaded
	assert.NotZero(t, found.Customer.ID)
	assert.Equal(t, customer.ID, found.Customer.ID)
	assert.Equal(t, user.ID, found.Customer.UserID)
	assert.Equal(t, true, found.Customer.IsActive)
}

// TestAppointmentRepository_ConcurrentCustomers tests appointments for multiple customers
func TestAppointmentRepository_ConcurrentCustomers(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAppointmentRepository(db)

	// Create 5 customers
	customers := make([]models.Customer, 5)
	for i := 0; i < 5; i++ {
		user := createTestUser(t, db, "customer"+string(rune(i))+"@example.com")
		customers[i] = createTestCustomer(t, db, user.ID)
	}

	service := createTestService(t, db, "Haircut", 50.00, 30)
	now := time.Now()

	// Create 3 appointments for each customer
	for i, customer := range customers {
		for j := 0; j < 3; j++ {
			repo.Create(models.Appointment{
				CustomerID: customer.ID,
				Services:   []models.Service{service},
				Date:       now.Add(time.Duration(i*24+j) * time.Hour),
				Status:     models.StatusPending,
			})
		}
	}

	// Verify ListAll returns all 15 appointments
	all, _ := repo.ListAll()
	assert.Len(t, all, 15)

	// Verify each customer has exactly 3 appointments
	for _, customer := range customers {
		weekStart := now
		weekEnd := now.Add(30 * 24 * time.Hour)
		customerApts, _ := repo.FindCustomerAppointmentsInWeek(customer.ID, weekStart, weekEnd)
		assert.Len(t, customerApts, 3)
	}
}

// TestAppointmentRepository_EmptyAppointmentNoServices tests edge case with no services
func TestAppointmentRepository_EmptyAppointmentNoServices(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAppointmentRepository(db)

	user := createTestUser(t, db, "customer@example.com")
	customer := createTestCustomer(t, db, user.ID)

	// Create appointment without services
	created, err := repo.Create(models.Appointment{
		CustomerID: customer.ID,
		Services:   []models.Service{},
		Date:       time.Now().Add(24 * time.Hour),
		Status:     models.StatusPending,
	})
	require.NoError(t, err)

	found, err := repo.FindByID(created.ID)
	assert.NoError(t, err)
	assert.Len(t, found.Services, 0)
}

// TestAppointmentRepository_TimestampFields tests that CreatedAt and UpdatedAt are set
func TestAppointmentRepository_TimestampFields(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAppointmentRepository(db)

	user := createTestUser(t, db, "customer@example.com")
	customer := createTestCustomer(t, db, user.ID)
	service := createTestService(t, db, "Haircut", 50.00, 30)

	beforeCreate := time.Now()
	created, err := repo.Create(models.Appointment{
		CustomerID: customer.ID,
		Services:   []models.Service{service},
		Date:       time.Now().Add(24 * time.Hour),
		Status:     models.StatusPending,
	})
	afterCreate := time.Now()

	require.NoError(t, err)
	assert.True(t, created.CreatedAt.After(beforeCreate.Add(-1*time.Second)) && created.CreatedAt.Before(afterCreate.Add(1*time.Second)))
	assert.True(t, created.UpdatedAt.After(beforeCreate.Add(-1*time.Second)) && created.UpdatedAt.Before(afterCreate.Add(1*time.Second)))

	// Update appointment
	time.Sleep(10 * time.Millisecond)
	beforeUpdate := time.Now()
	created.Status = models.StatusConfirmed
	repo.Update(created)
	afterUpdate := time.Now()

	found, _ := repo.FindByID(created.ID)
	assert.True(t, found.UpdatedAt.After(beforeUpdate.Add(-1*time.Second)) && found.UpdatedAt.Before(afterUpdate.Add(1*time.Second)))
	assert.True(t, found.UpdatedAt.After(created.UpdatedAt) || found.UpdatedAt.Equal(created.UpdatedAt))
}
