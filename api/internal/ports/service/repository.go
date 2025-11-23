package service

import "melodiapp/models"

type ServiceRepository interface {
	GetAll() ([]models.Service, error)
	GetByID(id string) (*models.Service, error)
	Create(svc *models.Service) error
	Update(svc *models.Service) error
	DeleteByID(id string) error
}
