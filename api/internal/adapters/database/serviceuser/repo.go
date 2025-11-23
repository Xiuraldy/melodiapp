package databaseadapter

import (
	"melodiapp/database"
	"melodiapp/models"
)

type GormServiceUserRepository struct{}

func NewGormServiceUserRepository() *GormServiceUserRepository {
	return &GormServiceUserRepository{}
}

func (r *GormServiceUserRepository) AddUsers(serviceID uint, userIDs []uint) error {
	for _, uid := range userIDs {
		su := models.ServiceUser{
			ServiceID: serviceID,
			UserID:    uid,
			Status:    "pending",
		}
		if err := database.DBConn.Create(&su).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *GormServiceUserRepository) ListByService(serviceID uint) ([]models.ServiceUser, error) {
	var list []models.ServiceUser
	result := database.DBConn.Where("service_id = ?", serviceID).Find(&list)
	return list, result.Error
}

func (r *GormServiceUserRepository) UpdateStatus(serviceID uint, userID uint, status string) error {
	return database.DBConn.Model(&models.ServiceUser{}).
		Where("service_id = ? AND user_id = ?", serviceID, userID).
		Update("status", status).Error
}
