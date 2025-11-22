// @title Hair Salon API
// @version 1.0
// @description Hair Salon Management API with JWT Authentication
// @host localhost:8080
// @basePath /api
// @schemes http
package main

import (
	"log"
	"os"

	_ "github.com/ViniciusBoroto/cabeleleila_leila/docs"
	"github.com/ViniciusBoroto/cabeleleila_leila/internal/handlers"
	"github.com/ViniciusBoroto/cabeleleila_leila/internal/models"
	"github.com/ViniciusBoroto/cabeleleila_leila/internal/repository"
	"github.com/ViniciusBoroto/cabeleleila_leila/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	_ "github.com/joho/godotenv/autoload"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()

	// Setup CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 3600,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	db := setupSqliteDB()
	migrateDatabase(db)
	seed(db)

	// Setup repositories
	userRepo := repository.NewUserRepository(db)
	apRepo := repository.NewAppointmentRepository(db)
	serviceRepo := repository.NewServiceRepository(db)

	// Setup services
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET environment variable not set")
	}
	authSvc := service.NewAuthService(jwtSecret)
	apSvc := service.NewAppointmentService(apRepo)
	serviceSvc := service.NewServiceService(serviceRepo)

	// Setup handlers
	authHandler := handlers.NewAuthHandler(authSvc, userRepo)
	appointmentsHandler := handlers.NewAppointmentHandler(apSvc)

	// Public routes
	public := r.Group("/api")
	{
		public.POST("/auth/login", authHandler.Login)
		public.POST("/auth/register", authHandler.Register)
		public.GET("/services", handlers.ListServices(serviceSvc))
		public.GET("/services/:id", handlers.GetService(serviceSvc))
	}

	// Protected routes - all authenticated users
	protected := r.Group("/api")
	protected.Use(handlers.JWTAuthMiddleware(authSvc))
	{
		// Appointment routes (for all authenticated users)
		appointmentsHandler.RegisterRoutes(protected)

		// User management routes (admin only)
		admin := protected.Group("/admin")
		admin.GET("/users", handlers.GetAllUsers(userRepo))
		admin.POST("/users", handlers.CreateUser(userRepo))
		admin.GET("/users/:id", handlers.GetUser(userRepo))
		admin.PUT("/users/:id", handlers.UpdateUser(userRepo))
		admin.DELETE("/users/:id", handlers.DeleteUser(userRepo))

		// Service management routes (admin only)
		admin.POST("/services", handlers.CreateService(serviceSvc))
		admin.PUT("/services/:id", handlers.UpdateService(serviceSvc))
		admin.DELETE("/services/:id", handlers.DeleteService(serviceSvc))
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func setupSqliteDB() *gorm.DB {
	path := os.Getenv("DB_PATH")
	if path == "" {
		panic("DB_PATH environment variable not set")
	}
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func migrateDatabase(db *gorm.DB) {
	if err := db.AutoMigrate(
		&models.User{},
		&models.Service{},
		&models.Appointment{},
	); err != nil {
		panic(err)
	}
}

func seed(db *gorm.DB) {
	// Check if services already exist
	var count int64
	db.Model(&models.Service{}).Count(&count)
	if count > 0 {
		return
	}

	// Default services
	defaultServices := []models.Service{
		{
			Name:            "Corte de Cabelo",
			Price:           50.0,
			DurationMinutes: 30,
		},
		{
			Name:            "Escova",
			Price:           40.0,
			DurationMinutes: 45,
		},
		{
			Name:            "Coloração",
			Price:           100.0,
			DurationMinutes: 120,
		},
		{
			Name:            "Hidratação",
			Price:           60.0,
			DurationMinutes: 60,
		},
		{
			Name:            "Manicure",
			Price:           30.0,
			DurationMinutes: 30,
		},
		{
			Name:            "Pedicure",
			Price:           35.0,
			DurationMinutes: 40,
		},
	}

	for _, service := range defaultServices {
		if err := db.Create(&service).Error; err != nil {
			log.Printf("Error seeding service %s: %v", service.Name, err)
		}
	}

	log.Println("Default services seeded successfully")
}
