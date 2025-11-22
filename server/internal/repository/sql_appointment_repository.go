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
	if err != nil {
		return ap, err
	}

	// Associate services with appointment (many-to-many)
	if len(ap.Services) > 0 {
		err = r.db.Model(&ap).Association("Services").Append(ap.Services)
		if err != nil {
			return ap, err
		}
	}

	// Reload the appointment with User and Services preloaded
	err = r.db.Preload("User").Preload("Services").First(&ap, ap.ID).Error
	return ap, err
}

func (r *sqlAppointmentRepo) Update(ap models.Appointment) error {
	return r.db.Save(&ap).Error
}

func (r *sqlAppointmentRepo) FindByID(id uint) (models.Appointment, error) {
	var ap models.Appointment
	err := r.db.Preload("User").Preload("Services").First(&ap, id).Error
	return ap, err
}

func (r *sqlAppointmentRepo) FindUserAppointmentsInWeek(userID uint, weekStart, weekEnd time.Time) ([]models.Appointment, error) {
	var list []models.Appointment
	err := r.db.Preload("User").Preload("Services").Where("user_id = ? AND date BETWEEN ? AND ?", userID, weekStart, weekEnd).Find(&list).Error
	return list, err
}

func (r *sqlAppointmentRepo) ListByPeriod(start, end time.Time) ([]models.Appointment, error) {
	var list []models.Appointment
	err := r.db.Preload("User").Preload("Services").Where("date BETWEEN ? AND ?", start, end).Find(&list).Error
	return list, err
}

func (r *sqlAppointmentRepo) ListAll() ([]models.Appointment, error) {
	var list []models.Appointment
	err := r.db.Preload("User").Preload("Services").Find(&list).Error
	return list, err
}
