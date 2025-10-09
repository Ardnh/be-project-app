package services

import (
	"context"

	"github.com/Ardnh/be-project-app/internal/application/dto"
	"github.com/Ardnh/be-project-app/internal/application/mapper"
	repositories "github.com/Ardnh/be-project-app/internal/domain/repositories"
	"github.com/Ardnh/be-project-app/internal/domain/services"
)

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) services.UserService {
	return &userService{
		userRepo: repo,
	}
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*dto.UserDto, error) {

	// Pass context ke repository
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Map ke DTO
	userDto := mapper.ToUserDTO(user)
	return userDto, nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]*dto.UserDto, error) {

	return nil, nil
}

func (s *userService) CreateUser(ctx context.Context, user *dto.CreateUserDto) error {

	err := s.userRepo.Create(ctx)

	return nil
}
