package databaseadapter

import (
	"melodiapp/database"
	"melodiapp/models"
)

type GormServiceUserRepository struct{}

func NewGormServiceUserRepository() *GormServiceUserRepository {
	return &GormServiceUserRepository{}
}

// AddUsers agrega usuarios al servicio, EVITANDO DUPLICADOS
func (r *GormServiceUserRepository) AddUsers(serviceID uint, userIDs []uint) error {
	for _, uid := range userIDs {
		// 1. Verificar si ya existe la relación
		var count int64
		err := database.DBConn.Model(&models.ServiceUser{}).
			Where("service_id = ? AND user_id = ?", serviceID, uid).
			Count(&count).Error

		if err != nil {
			return err // Retorna si hubo error de conexión a la DB
		}

		// 2. Si count es 0, significa que NO existe, entonces lo creamos
		if count == 0 {
			su := models.ServiceUser{
				ServiceID: serviceID,
				UserID:    uid,
				Status:    "pending",
			}
			if err := database.DBConn.Create(&su).Error; err != nil {
				return err
			}
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
