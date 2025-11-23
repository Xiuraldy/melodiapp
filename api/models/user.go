package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	Password          string    `json:"password"`
	Celphone          string    `json:"celphone"`
	Role              string    `json:"role"`
	Lastname          string    `json:"lastname"`
	ProfilePictureURL string    `json:"profile_picture_url" gorm:"column:profile_picture_url"`
	SecondaryRole     string    `json:"secondary_role" gorm:"column:secondary_role"`
}

type UserInput struct {
	Username      string `json:"username"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Role          string `json:"role"`
	Celphone      string `json:"celphone"`
	Lastname      string `json:"lastname"`
	SecondaryRole string `json:"secondary_role"`
}

func (user *User) BeforeSave(*gorm.DB) error {
	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(passwordHashed)
	return nil
}
