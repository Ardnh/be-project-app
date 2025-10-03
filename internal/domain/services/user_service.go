package services

import (
	"github.com/Ardnh/be-project-app/internal/application/dto"
)

type UserService interface {
	GetUserByID(id string) (*dto.UserDto, error)
	GetAllUsers() ([]*dto.UserDto, error)
	CreateUser(user *dto.UserDto) error
}
