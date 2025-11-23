package auth

import "melodiapp/models"

type AuthService interface {
	Register(input models.UserInput) (string, error)
	Login(input models.UserInput) (string, error)
	Logout(tokenStr string) error
}
