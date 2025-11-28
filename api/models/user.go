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
	ProfilePictureUrl string    `json:"profile_picture_url"`
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

func hashPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(passwordHashed)
	return nil
}

func (user *User) BeforeCreate(*gorm.DB) error {
	return hashPassword(user)
}

func (user *User) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("Password") {
		return hashPassword(user)
	}
	return nil
}
