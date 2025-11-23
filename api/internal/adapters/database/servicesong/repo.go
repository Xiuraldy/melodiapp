package databaseadapter

import (
	"melodiapp/database"
	"melodiapp/models"
)

type GormServiceSongRepository struct{}

func NewGormServiceSongRepository() *GormServiceSongRepository {
	return &GormServiceSongRepository{}
}

func (r *GormServiceSongRepository) AddSongs(serviceID uint, songIDs []uint) error {
	for _, sid := range songIDs {
		ss := models.ServiceSong{
			ServiceID: serviceID,
			SongID:    sid,
		}
		if err := database.DBConn.Create(&ss).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *GormServiceSongRepository) ListByService(serviceID uint) ([]models.ServiceSong, error) {
	var list []models.ServiceSong
	result := database.DBConn.Where("service_id = ?", serviceID).Find(&list)
	return list, result.Error
}

func (r *GormServiceSongRepository) Remove(serviceID uint, songID uint) error {
	return database.DBConn.Where("service_id = ? AND song_id = ?", serviceID, songID).
		Delete(&models.ServiceSong{}).Error
}
