package user

import (
	userports "melodiapp/internal/ports/user"
	"melodiapp/models"
)

type Service struct {
	repo userports.UserRepository
}

func NewService(repo userports.UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAllUsers()
}

func (s *Service) GetUserByID(id string) (*models.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *Service) GetUserByUintID(id uint) (*models.User, error) {
	return s.repo.GetUserByUintID(id)
}

func (s *Service) CreateUser(user *models.User) (*models.User, error) {
	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) UpdateUser(id string, updated *models.User) (*models.User, error) {
	existing, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	existing.Username = updated.Username
	existing.Email = updated.Email
	if updated.Password != "" {
		existing.Password = updated.Password
	}

	if err := s.repo.UpdateUser(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *Service) DeleteUser(id string) error {
	return s.repo.DeleteUserByID(id)
}
