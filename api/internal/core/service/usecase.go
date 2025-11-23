package service

import (
	serviceports "melodiapp/internal/ports/service"
	"melodiapp/models"
)

type ServiceUsecase struct {
	repo serviceports.ServiceRepository
}

func NewServiceUsecase(repo serviceports.ServiceRepository) *ServiceUsecase {
	return &ServiceUsecase{repo: repo}
}

func (s *ServiceUsecase) GetAll() ([]models.Service, error) {
	return s.repo.GetAll()
}

func (s *ServiceUsecase) GetByID(id string) (*models.Service, error) {
	return s.repo.GetByID(id)
}

func (s *ServiceUsecase) Create(svc *models.Service) (*models.Service, error) {
	if err := s.repo.Create(svc); err != nil {
		return nil, err
	}
	return svc, nil
}

func (s *ServiceUsecase) Update(id string, input *models.Service) (*models.Service, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil
	}

	existing.Name = input.Name
	existing.StartTime = input.StartTime
	existing.EndTime = input.EndTime

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *ServiceUsecase) Delete(id string) error {
	return s.repo.DeleteByID(id)
}
