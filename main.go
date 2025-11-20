package main

import (
	"log"
	"os"

	"github.com/ViniciusBoroto/cabeleleila_leila/internal/handlers"
	"github.com/ViniciusBoroto/cabeleleila_leila/internal/repository"
	"github.com/ViniciusBoroto/cabeleleila_leila/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()

	db := setupSqliteDB()
	apRepo := repository.NewAppointmentRepository(db)
	apSvc := service.NewAppointmentService(apRepo)

	h := handlers.NewAppointmentHandler(apSvc)
	api := r.Group("/api")
	h.RegisterRoutes(api)

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
