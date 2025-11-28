package databaseadapter

import (
	"melodiapp/database"
	"melodiapp/models"
)

type GormServiceOutfitRepository struct{}

func NewGormServiceOutfitRepository() *GormServiceOutfitRepository {
	return &GormServiceOutfitRepository{}
}

func (r *GormServiceOutfitRepository) AddOutfits(serviceID uint, outfitIDs []uint) error {
	for _, oid := range outfitIDs {
		so := models.ServiceOutfit{
			ServiceID: serviceID,
			OutfitID:  oid,
		}
		if err := database.DBConn.Create(&so).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *GormServiceOutfitRepository) ListByService(serviceID uint) ([]models.ServiceOutfit, error) {
	var list []models.ServiceOutfit
	result := database.DBConn.Where("service_id = ?", serviceID).Find(&list)
	return list, result.Error
}

func (r *GormServiceOutfitRepository) Remove(serviceID uint, outfitID uint) error {
	return database.DBConn.Where("service_id = ? AND outfit_id = ?", serviceID, outfitID).
		Delete(&models.ServiceOutfit{}).Error
}
