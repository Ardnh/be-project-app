package repository

import (
	"context"

	"github.com/Ardnh/be-project-app/internal/domain/entities"
	"github.com/Ardnh/be-project-app/internal/domain/repositories"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// Struct implementasi (private, lowercase)
type projectsRepositoryImpl struct {
	db    *gorm.DB
	redis *redis.Client
}

// Constructor yang return INTERFACE dari domain
// ⬇️ Return type adalah INTERFACE dari domain
func NewProjectsRepository(db *gorm.DB, redis *redis.Client) repositories.ProjectsRepository {
	return &projectsRepositoryImpl{
		db:    db,
		redis: redis,
	}
}

func (r *projectsRepositoryImpl) FindByUserID(ctx context.Context, userId string) ([]*entities.Project, error) {

	var projects []*entities.Project
	r.db.Preload("").Limit(10).Find(&projects).Where("user_id=?", userId)

	return projects, nil
}

func (r *projectsRepositoryImpl) FindProjectCategoryByUserID(ctx context.Context, userId string) ([]string, error) {

	return nil, nil
}

func (r *projectsRepositoryImpl) FindAll(ctx context.Context, limit int, offset int) ([]*entities.Project, error) {

	return nil, nil
}

func (r *projectsRepositoryImpl) Create(ctx context.Context, project *entities.Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}

func (r *projectsRepositoryImpl) Update(ctx context.Context, user *entities.Project) error {

	return nil
}

func (r *projectsRepositoryImpl) Delete(ctx context.Context, id string) error {

	return nil
}
