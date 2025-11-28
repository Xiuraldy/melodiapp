package databaseadapter

import (
	"log"

	"melodiapp/database"
	"melodiapp/models"
)

type GormServiceSongRepository struct{}

func NewGormServiceSongRepository() *GormServiceSongRepository {
	return &GormServiceSongRepository{}
}

func (r *GormServiceSongRepository) AddSongs(serviceID uint, songIDs []uint) error {
	log.Printf("[ServiceSongRepository] Replacing songs for service %d with %+v", serviceID, songIDs)

	// Reemplazar completamente el repertorio del servicio:
	// 1) Borrar asociaciones anteriores
	// 2) Crear solo las nuevas seleccionadas
	tx := database.DBConn.Begin()

	res := tx.Where("service_id = ?", serviceID).Delete(&models.ServiceSong{})
	if res.Error != nil {
		log.Printf("[ServiceSongRepository] Error deleting existing songs for service %d: %v", serviceID, res.Error)
		tx.Rollback()
		return res.Error
	}
	log.Printf("[ServiceSongRepository] Deleted %d existing service_songs rows for service %d", res.RowsAffected, serviceID)

	for _, sid := range songIDs {
		ss := models.ServiceSong{
			ServiceID: serviceID,
			SongID:    sid,
		}
		if err := tx.Create(&ss).Error; err != nil {
			log.Printf("[ServiceSongRepository] Error creating service_song (service=%d, song=%d): %v", serviceID, sid, err)
			tx.Rollback()
			return err
		}
		log.Printf("[ServiceSongRepository] Inserted service_song (service=%d, song=%d)", serviceID, sid)
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("[ServiceSongRepository] Error committing transaction for service %d: %v", serviceID, err)
		return err
	}

	log.Printf("[ServiceSongRepository] Successfully replaced songs for service %d", serviceID)
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
