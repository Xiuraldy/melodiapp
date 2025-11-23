package serviceuser

import "melodiapp/models"

type ServiceUserService interface {
	AssignUsers(serviceID uint, userIDs []uint) error
	ListByService(serviceID uint) ([]models.ServiceUser, error)
	ChangeStatus(serviceID uint, userID uint, status string) error
}
