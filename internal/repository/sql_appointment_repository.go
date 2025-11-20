package repository

import (
	"time"

	"github.com/ViniciusBoroto/cabeleleila_leila/internal/models"
	"gorm.io/gorm"
)

type sqlAppointmentRepo struct {
	db *gorm.DB
}

func NewAppointmentRepository(db *gorm.DB) AppointmentRepository {
	return &sqlAppointmentRepo{db}
}

func (r *sqlAppointmentRepo) Create(ap models.Appointment) (models.Appointment, error) {
	err := r.db.Create(&ap).Error
	return ap, err
}

func (r *sqlAppointmentRepo) Update(ap models.Appointment) error {
	return r.db.Save(&ap).Error
}

func (r *sqlAppointmentRepo) FindByID(id uint) (models.Appointment, error) {
	var ap models.Appointment
	err := r.db.Preload("Customer").First(&ap, id).Error
	return ap, err
}

func (r *sqlAppointmentRepo) FindCustomerAppointmentsInWeek(customerID uint, weekStart, weekEnd time.Time) ([]models.Appointment, error) {
	var list []models.Appointment
	err := r.db.Where("customer_id = ? AND date BETWEEN ? AND ?", customerID, weekStart, weekEnd).Find(&list).Error
	return list, err
}

func (r *sqlAppointmentRepo) ListByPeriod(start, end time.Time) ([]models.Appointment, error) {
	var list []models.Appointment
	err := r.db.Where("date BETWEEN ? AND ?", start, end).Find(&list).Error
	return list, err
}

func (r *sqlAppointmentRepo) ListAll() ([]models.Appointment, error) {
	var list []models.Appointment
	err := r.db.Find(&list).Error
	return list, err
}
