package repository

import (
	"context"

	"github.com/Ardnh/be-project-app/internal/domain/entities"
	"github.com/Ardnh/be-project-app/internal/domain/repositories"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type ProjectExpensesItemRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewProjectExpensesItemRepository(db *gorm.DB, redis *redis.Client) repositories.ProjectExpensesItemRepository {
	return &ProjectExpensesItemRepository{
		db:    db,
		redis: redis,
	}
}

func (r *ProjectExpensesItemRepository) Create(ctx context.Context, expenses *entities.ProjectExpenseItem) error {
	return r.db.WithContext(ctx).Create(expenses).Error
}

func (r *ProjectExpensesItemRepository) Update(ctx context.Context, expenses *entities.ProjectExpenseItem) error {

	result := r.db.
		WithContext(ctx).
		Model(&entities.ProjectExpenseItem{}).
		Where("id = ?", expenses.ID).
		Updates(expenses)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *ProjectExpensesItemRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.ProjectExpenseItem{}).Error
}
