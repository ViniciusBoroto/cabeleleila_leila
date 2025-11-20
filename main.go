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
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	_ "github.com/joho/godotenv/autoload"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	db := setupSqliteDB()
	migrateDatabase(db)

	// Setup repositories
	userRepo := repository.NewUserRepository(db)
	apRepo := repository.NewAppointmentRepository(db)

	// Setup services
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET environment variable not set")
	}
	authSvc := service.NewAuthService(jwtSecret)
	apSvc := service.NewAppointmentService(apRepo)

	// Setup handlers
	authHandler := handlers.NewAuthHandler(authSvc, userRepo)
	appointmentsHandler := handlers.NewAppointmentHandler(apSvc)

	// Public routes
	public := r.Group("/api")
	public.POST("/auth/login", authHandler.Login)
	public.POST("/auth/register", authHandler.Register)

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
		&models.Customer{},
		&models.Service{},
		&models.Appointment{},
	); err != nil {
		panic(err)
	}
}
