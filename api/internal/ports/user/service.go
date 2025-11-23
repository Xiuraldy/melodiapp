package user

import "melodiapp/models"

type UserService interface {
	GetAllUsers() ([]models.User, error)
	GetUserByID(id string) (*models.User, error)
	GetUserByUintID(id uint) (*models.User, error)
	CreateUser(user *models.User) (*models.User, error)
	UpdateUser(id string, updated *models.User) (*models.User, error)
	DeleteUser(id string) error
}
