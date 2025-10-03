package repositories

import (
	"github.com/Ardnh/be-project-app/internal/domain/entities"
)

// Interface ini implementasi nya ada di infra -> database -> postgresql
type UserRepository interface {
	FindByID(id string) (*entities.User, error)
	FindAll(limit int, offset int) ([]*entities.User, error)
	Create(user *entities.User) error
	Update(user *entities.User) error
	Delete(userId string) error
	ExistsByEmail(email string) (bool, error)
}
