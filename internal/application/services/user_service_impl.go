package services

import (
	"github.com/Ardnh/be-project-app/internal/application/dto"
	repositories "github.com/Ardnh/be-project-app/internal/domain/repositories"
	"github.com/Ardnh/be-project-app/internal/domain/services"
)

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) services.UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUserByID(id string) (*dto.UserDto, error) {

	return nil, nil
}

func (s *userService) GetAllUsers() ([]*dto.UserDto, error) {

	return nil, nil
}

func (s *userService) CreateUser(user *dto.UserDto) error {

	return nil
}
