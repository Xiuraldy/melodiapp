package song

import "melodiapp/models"

type SongRepository interface {
	GetAll() ([]models.Song, error)
	GetByID(id string) (*models.Song, error)
	Create(song *models.Song) error
	Update(song *models.Song) error
	DeleteByID(id string) error
}
