package repositories

import (
	"context"

	"github.com/Ardnh/be-project-app/internal/domain/entities"
)

// Interface ini implementasi nya ada di infra -> database -> postgresql
type UserRepository interface {
	FindByID(ctx context.Context, id string) (*entities.User, error)
	FindAll(ctx context.Context, limit int, offset int) ([]*entities.User, error)
	Create(ctx context.Context, user *entities.User) error
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, userId string) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
