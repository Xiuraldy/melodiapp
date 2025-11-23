package databaseadapter

import (
	"errors"

	"gorm.io/gorm"
	"melodiapp/database"
	"melodiapp/models"
)

type GormSongRepository struct{}

func NewGormSongRepository() *GormSongRepository {
	return &GormSongRepository{}
}

func (r *GormSongRepository) GetAll() ([]models.Song, error) {
	var songs []models.Song
	result := database.DBConn.Find(&songs)
	return songs, result.Error
}

func (r *GormSongRepository) GetByID(id string) (*models.Song, error) {
	var song models.Song
	result := database.DBConn.Where("id = ?", id).First(&song)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &song, result.Error
}

func (r *GormSongRepository) Create(song *models.Song) error {
	return database.DBConn.Create(song).Error
}

func (r *GormSongRepository) Update(song *models.Song) error {
	return database.DBConn.Save(song).Error
}

func (r *GormSongRepository) DeleteByID(id string) error {
	return database.DBConn.Delete(&models.Song{}, id).Error
}
