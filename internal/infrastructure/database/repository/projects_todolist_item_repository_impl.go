package repository

import (
	"context"

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

func (r *projectTodolistItemRepositoryImpl) CreateProjectTodolistItem(ctx context.Context, todoItem *entities.ProjectTodolistItems) error {

	return r.db.WithContext(ctx).Model(&entities.ProjectTodolistItems{}).Create(&todoItem).Error
}

func (r *projectTodolistItemRepositoryImpl) UpdateProjectTodolistItem(ctx context.Context, todoItem *entities.ProjectTodolistItems) error {

	result := r.db.WithContext(ctx).Model(&entities.ProjectTodolistItems{}).Updates(&todoItem)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *projectTodolistItemRepositoryImpl) DeleteProjectTodolistItem(ctx context.Context, id string) error {

	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.ProjectTodolistItems{}).Error
}
