package serviceoutfit

import "melodiapp/models"

type ServiceOutfitRepository interface {
	AddOutfits(serviceID uint, outfitIDs []uint) error

	ListByService(serviceID uint) ([]models.ServiceOutfit, error)

	Remove(serviceID uint, outfitID uint) error
}
