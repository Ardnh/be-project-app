package repository

import (
	"context"
	"fmt"

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

func (r *projectsExpensesRepositoryImpl) Create(ctx context.Context, expenses *entities.ProjectExpenses) (*entities.ProjectExpenses, error) {

	if err := r.db.WithContext(ctx).Create(&expenses).Error; err != nil {
		return nil, err
	}

	return expenses, nil
}
func (r *projectsExpensesRepositoryImpl) Update(ctx context.Context, expenses *entities.ProjectExpenses) (*entities.ProjectExpenses, error) {

	var exists bool
	if err := r.db.WithContext(ctx).
		Model(&entities.ProjectExpenses{}).
		Select("1").
		Where("id = ?", expenses.ID).
		Limit(1).
		Find(&exists).Error; err != nil {
		return nil, err
	}

	if !exists {
		return nil, gorm.ErrRecordNotFound
	}

	result := r.db.
		WithContext(ctx).
		Model(&entities.ProjectExpenses{}).
		Where("id = ?", expenses.ID).
		Updates(expenses)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	var updatedExpenses entities.ProjectExpenses
	if err := r.db.WithContext(ctx).
		Where("id = ?", expenses.ID).
		First(&updatedExpenses).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch updated item: %w", err)
	}

	return &updatedExpenses, nil
}

func (r *projectsExpensesRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.ProjectExpenses{}).Error
}
