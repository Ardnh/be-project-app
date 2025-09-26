package services

import (
	domainErrors "github.com/Ardnh/be-project-app/internal/domain"
	"github.com/Ardnh/be-project-app/internal/domain/entities"
	"github.com/Ardnh/be-project-app/internal/domain/repositories"
)

type UserService interface {
	GetUserByID(id string) (*entities.User, error)
	GetAllUsers() ([]*entities.User, error)
	CreateUser(user *entities.User) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUserByID(id string) (*entities.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, domainErrors.ErrUserNotFound
	}
	return user, nil
}

func (s *userService) GetAllUsers() ([]*entities.User, error) {
	return s.repo.FindAll()
}

func (s *userService) CreateUser(user *entities.User) error {
	// contoh validasi email
	users, _ := s.repo.FindAll()
	for _, u := range users {
		if u.Email == user.Email {
			return domainErrors.ErrEmailAlreadyExists
		}
	}
	return s.repo.Create(user)
}
