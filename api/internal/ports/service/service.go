package service

import "melodiapp/models"

type ServiceService interface {
	GetAll() ([]models.Service, error)
	GetByID(id string) (*models.Service, error)
	Create(svc *models.Service) (*models.Service, error)
	Update(id string, input *models.Service) (*models.Service, error)
	Delete(id string) error
}
