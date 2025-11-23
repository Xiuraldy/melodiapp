package song

import (
	songports "melodiapp/internal/ports/song"
	"melodiapp/models"
)

type Service struct {
	repo songports.SongRepository
}

func NewService(repo songports.SongRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll() ([]models.Song, error) {
	return s.repo.GetAll()
}

func (s *Service) GetByID(id string) (*models.Song, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Create(song *models.Song) (*models.Song, error) {
	if err := s.repo.Create(song); err != nil {
		return nil, err
	}
	return song, nil
}

func (s *Service) Update(id string, input *models.Song) (*models.Song, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil
	}

	existing.Name = input.Name
	existing.Author = input.Author
	existing.SongKey = input.SongKey
	existing.BPM = input.BPM
	existing.TimeSignature = input.TimeSignature
	existing.Duration = input.Duration
	existing.Structure = input.Structure
	existing.HasSequence = input.HasSequence
	existing.HasChart = input.HasChart
	existing.HasScore = input.HasScore
	existing.YoutubeURL = input.YoutubeURL
	existing.VoiceURL = input.VoiceURL
	existing.GuitarURL = input.GuitarURL
	existing.PianoURL = input.PianoURL
	existing.DrumsURL = input.DrumsURL
	existing.BassURL = input.BassURL
	existing.WindURL = input.WindURL

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) Delete(id string) error {
	return s.repo.DeleteByID(id)
}
