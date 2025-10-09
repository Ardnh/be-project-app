package services

import (
	"context"

	"github.com/Ardnh/be-project-app/internal/application/dto"
)

type UserService interface {
	GetUserByID(ctx context.Context, id string) (*dto.UserDto, error)
	GetAllUsers(ctx context.Context) ([]*dto.UserDto, error)
	CreateUser(ctx context.Context, user *dto.CreateUserDto) error
}
