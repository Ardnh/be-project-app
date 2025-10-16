package services

import (
	"context"

	"github.com/Ardnh/be-project-app/internal/application/dto"
)

type AuthService interface {
	Login(ctx context.Context, req *dto.LoginDto) (*string, error)
	Register(ctx context.Context, req *dto.RegisterDto) error
}
