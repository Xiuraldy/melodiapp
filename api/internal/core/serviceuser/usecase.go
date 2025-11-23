package serviceuser

import (
	serviceuserports "melodiapp/internal/ports/serviceuser"
	"melodiapp/models"
)

type Service struct {
	repo serviceuserports.ServiceUserRepository
}

func NewService(repo serviceuserports.ServiceUserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AssignUsers(serviceID uint, userIDs []uint) error {
	return s.repo.AddUsers(serviceID, userIDs)
}

func (s *Service) ListByService(serviceID uint) ([]models.ServiceUser, error) {
	return s.repo.ListByService(serviceID)
}

func (s *Service) ChangeStatus(serviceID uint, userID uint, status string) error {
	return s.repo.UpdateStatus(serviceID, userID, status)
}
