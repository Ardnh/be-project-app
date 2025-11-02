package services

import (
	"context"
	"time"

	"github.com/Ardnh/be-project-app/internal/application/dto"
	"github.com/Ardnh/be-project-app/internal/application/mapper"
	"github.com/Ardnh/be-project-app/internal/domain"
	"github.com/Ardnh/be-project-app/internal/domain/entities"
	repositories "github.com/Ardnh/be-project-app/internal/domain/repositories"
	"github.com/Ardnh/be-project-app/internal/domain/services"
	"golang.org/x/crypto/bcrypt"
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

	existingUser, err := s.userRepo.ExistsByEmail(ctx, user.Email)
	if err == nil && existingUser {
		return domain.ErrEmailAlreadyExists
	}

	userEntities := &entities.User{
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	userEntities.Password = string(hashedPassword)

	errCreate := s.userRepo.Create(ctx, userEntities)
	if errCreate != nil {
		return errCreate
	}

	return nil
}

func (s *userService) UpdateUser(ctx context.Context, user *dto.UpdateUserDto) error {

	existingUser, err := s.userRepo.ExistsByEmail(ctx, user.Email)
	if err == nil && existingUser {
		return domain.ErrEmailAlreadyExists
	}

	userEntities := &entities.User{
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	errCreate := s.userRepo.Create(ctx, userEntities)
	if errCreate != nil {
		return errCreate
	}

	return nil
}
