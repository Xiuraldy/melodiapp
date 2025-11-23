package song

import "melodiapp/models"

type SongService interface {
	GetAll() ([]models.Song, error)
	GetByID(id string) (*models.Song, error)
	Create(song *models.Song) (*models.Song, error)
	Update(id string, input *models.Song) (*models.Song, error)
	Delete(id string) error
}
