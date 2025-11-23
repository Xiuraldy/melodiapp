package servicesong

import "melodiapp/models"

type ServiceSongService interface {
	AssignSongs(serviceID uint, songIDs []uint) error
	ListByService(serviceID uint) ([]models.ServiceSong, error)
	Remove(serviceID uint, songID uint) error
}
