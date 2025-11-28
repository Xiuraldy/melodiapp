package serviceoutfit

import "melodiapp/models"

type ServiceOutfitService interface {
	AssignOutfits(serviceID uint, outfitIDs []uint) error
	ListByService(serviceID uint) ([]models.ServiceOutfit, error)
	Remove(serviceID uint, outfitID uint) error
}
