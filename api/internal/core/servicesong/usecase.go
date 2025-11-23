package servicesong

import (
	servicesongports "melodiapp/internal/ports/servicesong"
	"melodiapp/models"
)

type Service struct {
	repo servicesongports.ServiceSongRepository
}

func NewService(repo servicesongports.ServiceSongRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AssignSongs(serviceID uint, songIDs []uint) error {
	return s.repo.AddSongs(serviceID, songIDs)
}

func (s *Service) ListByService(serviceID uint) ([]models.ServiceSong, error) {
	return s.repo.ListByService(serviceID)
}

func (s *Service) Remove(serviceID uint, songID uint) error {
	return s.repo.Remove(serviceID, songID)
}
