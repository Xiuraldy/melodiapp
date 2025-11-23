package serviceuser

import "melodiapp/models"

type ServiceUserRepository interface {
	AddUsers(serviceID uint, userIDs []uint) error
	ListByService(serviceID uint) ([]models.ServiceUser, error)
	UpdateStatus(serviceID uint, userID uint, status string) error
}
