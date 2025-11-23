package databaseadapter

import (
	"errors"

	"gorm.io/gorm"
	"melodiapp/database"
	"melodiapp/models"
)

type GormUserRepository struct{}

func NewGormUserRepository() *GormUserRepository {
	return &GormUserRepository{}
}

func (r *GormUserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	result := database.DBConn.Find(&users)
	return users, result.Error
}

func (r *GormUserRepository) GetUserByID(id string) (*models.User, error) {
	var user models.User
	result := database.DBConn.Where("id = ?", id).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, result.Error
}

func (r *GormUserRepository) GetUserByUintID(id uint) (*models.User, error) {
	var user models.User
	result := database.DBConn.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, result.Error
}

func (r *GormUserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := database.DBConn.Where("email = ?", email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, result.Error
}

func (r *GormUserRepository) CreateUser(user *models.User) error {
	result := database.DBConn.Create(user)
	return result.Error
}

func (r *GormUserRepository) UpdateUser(user *models.User) error {
	result := database.DBConn.Save(user)
	return result.Error
}

func (r *GormUserRepository) DeleteUserByID(id string) error {
	result := database.DBConn.Delete(&models.User{}, id)
	return result.Error
}
