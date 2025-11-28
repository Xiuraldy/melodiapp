package databaseadapter

import (
	"log"

	"melodiapp/database"
	"melodiapp/models"
)

type GormServiceUserRepository struct{}

func NewGormServiceUserRepository() *GormServiceUserRepository {
	return &GormServiceUserRepository{}
}

// AddUsers reemplaza completamente el equipo de un servicio:
// 1) Elimina todas las filas existentes de service_users para ese serviceID
// 2) Inserta únicamente los userIDs recibidos, con estado "pending" por defecto
func (r *GormServiceUserRepository) AddUsers(serviceID uint, userIDs []uint) error {
	log.Printf("[ServiceUserRepository] Replacing users for service %d with %+v", serviceID, userIDs)

	tx := database.DBConn.Begin()

	// Borrar asociaciones anteriores
	res := tx.Where("service_id = ?", serviceID).Delete(&models.ServiceUser{})
	if res.Error != nil {
		log.Printf("[ServiceUserRepository] Error deleting existing users for service %d: %v", serviceID, res.Error)
		tx.Rollback()
		return res.Error
	}
	log.Printf("[ServiceUserRepository] Deleted %d existing service_users rows for service %d", res.RowsAffected, serviceID)

	// Crear nuevas asociaciones
	for _, uid := range userIDs {
		su := models.ServiceUser{
			ServiceID: serviceID,
			UserID:    uid,
			Status:    "pending",
		}
		if err := tx.Create(&su).Error; err != nil {
			log.Printf("[ServiceUserRepository] Error creating service_user (service=%d, user=%d): %v", serviceID, uid, err)
			tx.Rollback()
			return err
		}
		log.Printf("[ServiceUserRepository] Inserted service_user (service=%d, user=%d)", serviceID, uid)
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("[ServiceUserRepository] Error committing transaction for service %d: %v", serviceID, err)
		return err
	}

	log.Printf("[ServiceUserRepository] Successfully replaced users for service %d", serviceID)
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
