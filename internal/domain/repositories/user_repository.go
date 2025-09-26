package repositories

import (
	"github.com/Ardnh/be-project-app/internal/domain/entities"
)

type UserRepository interface {
	FindByID(id string) (*entities.User, error)
	FindAll() ([]*entities.User, error)
	Create(user *entities.User) error
}
