package databaseadapter

import (
	"errors"

	"gorm.io/gorm"
	"melodiapp/database"
	"melodiapp/models"
)

type GormServiceRepository struct{}

func NewGormServiceRepository() *GormServiceRepository {
	return &GormServiceRepository{}
}

func (r *GormServiceRepository) GetAll() ([]models.Service, error) {
	var services []models.Service
	result := database.DBConn.Find(&services)
	return services, result.Error
}

func (r *GormServiceRepository) GetByID(id string) (*models.Service, error) {
	var svc models.Service
	result := database.DBConn.Where("id = ?", id).First(&svc)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &svc, result.Error
}

func (r *GormServiceRepository) Create(svc *models.Service) error {
	return database.DBConn.Create(svc).Error
}

func (r *GormServiceRepository) Update(svc *models.Service) error {
	return database.DBConn.Save(svc).Error
}

func (r *GormServiceRepository) DeleteByID(id string) error {
	return database.DBConn.Delete(&models.Service{}, id).Error
}
