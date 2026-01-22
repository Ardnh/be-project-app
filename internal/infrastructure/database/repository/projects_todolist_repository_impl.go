package repository

import (
	"context"
	"fmt"

	"github.com/Ardnh/be-project-app/internal/domain/entities"
	"github.com/Ardnh/be-project-app/internal/domain/repositories"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type projectTodolistRepositoryImpl struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewProjectTodolistRepository(db *gorm.DB, redis *redis.Client) repositories.ProjectTodolistRepository {
	return &projectTodolistRepositoryImpl{
		db:    db,
		redis: redis,
	}
}

func (r *projectTodolistRepositoryImpl) CreateProjectTodolist(ctx context.Context, todo *entities.ProjectTodolists) (*entities.ProjectTodolists, error) {

	if err := r.db.WithContext(ctx).Model(&entities.ProjectTodolists{}).Create(&todo).Error; err != nil {
		return nil, err
	}

	return todo, nil
}

func (r *projectTodolistRepositoryImpl) UpdateProjectTodolist(ctx context.Context, todo *entities.ProjectTodolists) (*entities.ProjectTodolists, error) {

	var exists bool
	if err := r.db.WithContext(ctx).
		Model(&entities.ProjectTodolists{}).
		Select("1").
		Where("id = ?", todo.ID).
		Limit(1).
		Find(&exists).Error; err != nil {
		return nil, err
	}

	if !exists {
		return nil, gorm.ErrRecordNotFound
	}

	result := r.db.WithContext(ctx).Updates(&todo)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	var updatedTodo entities.ProjectTodolists
	if err := r.db.WithContext(ctx).
		Where("id = ?", todo.ID).
		First(&updatedTodo).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch updated item: %w", err)
	}

	return &updatedTodo, nil
}

func (r *projectTodolistRepositoryImpl) DeleteProjectTodolist(ctx context.Context, id string) error {

	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.ProjectTodolists{}).Error
}
