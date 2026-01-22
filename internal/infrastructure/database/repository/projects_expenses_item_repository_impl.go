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

func (r *ProjectExpensesItemRepository) Create(ctx context.Context, expensesItem *entities.ProjectExpenseItems) (*entities.ProjectExpenseItems, error) {

	if err := r.db.WithContext(ctx).Model(&entities.ProjectExpenseItems{}).Create(&expensesItem).Error; err != nil {
		return nil, err
	}

	return expensesItem, nil
}

func (r *ProjectExpensesItemRepository) Update(ctx context.Context, expenses *entities.ProjectExpenseItems) (*entities.ProjectExpenseItems, error) {

	result := r.db.
		WithContext(ctx).
		Model(&entities.ProjectExpenseItems{}).
		Where("id = ?", expenses.ID).
		Updates(&expenses)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	var updatedExpense entities.ProjectExpenseItems
	if err := r.db.WithContext(ctx).
		Where("id = ?", expenses.ID).
		First(&updatedExpense).Error; err != nil {
		return nil, err
	}

	return &updatedExpense, nil
}

func (r *ProjectExpensesItemRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.ProjectExpenseItems{}).Error
}
