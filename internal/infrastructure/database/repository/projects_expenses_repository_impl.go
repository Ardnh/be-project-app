package repository

import (
	"context"

	"github.com/Ardnh/be-project-app/internal/domain/entities"
	"github.com/Ardnh/be-project-app/internal/domain/repositories"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// Struct implementasi (private, lowercase)
type projectsExpensesRepositoryImpl struct {
	db    *gorm.DB
	redis *redis.Client
}

// Constructor yang return INTERFACE dari domain
// ⬇️ Return type adalah INTERFACE dari domain
func NewProjectsExpensesRepository(db *gorm.DB, redis *redis.Client) repositories.ProjectsExpensesRepository {
	return &projectsExpensesRepositoryImpl{
		db:    db,
		redis: redis,
	}
}

func (r *projectsExpensesRepositoryImpl) Create(ctx context.Context, expenses *entities.ProjectExpenses) error {

	return r.db.WithContext(ctx).Create(expenses).Error
}
func (r *projectsExpensesRepositoryImpl) Update(ctx context.Context, expenses *entities.ProjectExpenses) error {

	return nil
}
func (r *projectsExpensesRepositoryImpl) Delete(ctx context.Context, id string) error {

	return nil
}
