package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Ardnh/be-project-app/internal/domain/entities"
	"github.com/Ardnh/be-project-app/internal/domain/repositories"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type projectTodolistItemRepositoryImpl struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewProjectTodolistItemRepository(db *gorm.DB, redis *redis.Client) repositories.ProjectTodolistItemRepository {
	return &projectTodolistItemRepositoryImpl{
		db:    db,
		redis: redis,
	}
}

func (r *projectTodolistItemRepositoryImpl) CreateProjectTodolistItem(ctx context.Context, todoItem *entities.ProjectTodolistItems) (*entities.ProjectTodolistItems, error) {
	// Create the item
	if err := r.db.WithContext(ctx).Model(&entities.ProjectTodolistItems{}).Create(&todoItem).Error; err != nil {
		return nil, err
	}

	return todoItem, nil
}

func (r *projectTodolistItemRepositoryImpl) UpdateProjectTodolistItem(ctx context.Context, todoItem *entities.ProjectTodolistItems) (*entities.ProjectTodolistItems, error) {
	// Validasi
	if todoItem.ID == "" {
		return nil, errors.New("project ID is required")
	}

	// Cek apakah record ada
	var exists bool
	if err := r.db.WithContext(ctx).
		Model(&entities.ProjectTodolistItems{}).
		Select("1").
		Where("id = ?", todoItem.ID).
		Limit(1).
		Find(&exists).Error; err != nil {
		return nil, fmt.Errorf("failed to check project existence: %w", err)
	}

	if !exists {
		return nil, gorm.ErrRecordNotFound
	}

	// Buat map untuk update (map tidak skip zero values)
	updates := map[string]interface{}{
		"name":          todoItem.Name,
		"category_name": todoItem.CategoryName,
		"is_completed":  todoItem.IsCompleted,
		"updated_at":    time.Now(),
	}

	result := r.db.WithContext(ctx).
		Model(&entities.ProjectTodolistItems{}).
		Where("id = ?", todoItem.ID).
		Updates(updates)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to update project: %w", result.Error)
	}

	// Fetch the updated item to return the latest data
	var updatedItem entities.ProjectTodolistItems
	if err := r.db.WithContext(ctx).
		Where("id = ?", todoItem.ID).
		First(&updatedItem).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch updated item: %w", err)
	}

	return &updatedItem, nil
}

func (r *projectTodolistItemRepositoryImpl) DeleteProjectTodolistItem(ctx context.Context, id string) error {

	// Validasi
	if id == "" {
		return errors.New("project ID is required")
	}

	// Cek apakah record ada
	var exists bool
	if err := r.db.WithContext(ctx).
		Model(&entities.ProjectTodolistItems{}).
		Select("1").
		Where("id = ?", id).
		Limit(1).
		Find(&exists).Error; err != nil {
		return fmt.Errorf("failed to check project existence: %w", err)
	}

	if !exists {
		return gorm.ErrRecordNotFound
	}

	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.ProjectTodolistItems{}).Error
}
