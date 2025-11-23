package servicesong

import "melodiapp/models"

type ServiceSongRepository interface {
	AddSongs(serviceID uint, songIDs []uint) error
	ListByService(serviceID uint) ([]models.ServiceSong, error)
	Remove(serviceID uint, songID uint) error
}
