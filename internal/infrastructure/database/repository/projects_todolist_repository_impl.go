package repository

import (
	"context"

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

func (r *projectTodolistRepositoryImpl) CreateProjectTodolist(ctx context.Context, todo *entities.ProjectTodolists) error {

	return r.db.WithContext(ctx).Model(&entities.ProjectTodolists{}).Create(&todo).Error
}

func (r *projectTodolistRepositoryImpl) UpdateProjectTodolist(ctx context.Context, todo *entities.ProjectTodolists) error {

	result := r.db.WithContext(ctx).Updates(&todo)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *projectTodolistRepositoryImpl) DeleteProjectTodolist(ctx context.Context, id string) error {

	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.ProjectTodolists{}).Error
}
