package serviceoutfit

import (
	serviceoutfitports "melodiapp/internal/ports/serviceoutfit"
	"melodiapp/models"
)

type Service struct {
	repo serviceoutfitports.ServiceOutfitRepository
}

func NewService(repo serviceoutfitports.ServiceOutfitRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AssignOutfits(serviceID uint, outfitIDs []uint) error {
	return s.repo.AddOutfits(serviceID, outfitIDs)
}

func (s *Service) ListByService(serviceID uint) ([]models.ServiceOutfit, error) {
	return s.repo.ListByService(serviceID)
}

func (s *Service) Remove(serviceID uint, outfitID uint) error {
	return s.repo.Remove(serviceID, outfitID)
}
