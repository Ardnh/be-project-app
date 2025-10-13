package repository

import (
	"context"
	"errors"
	"fmt"

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

func (r *projectsRepositoryImpl) FindByUserID(ctx context.Context, userId string, params entities.GetProjectsParams) ([]*entities.Projects, error) {
	var projects []*entities.Projects
	// var projectWithTodolistAndExpenses *entities.ProjectWithTodolistAndExpenses

	query := r.
		db.
		WithContext(ctx).
		Table("projects").
		Where("user_id = ?", userId)

	// Filter by category name jika categoryName tidak kosong
	if params.CategoryName != "" {
		query = query.Where("category_name LIKE ?", params.CategoryName)
	}

	if params.Search != "" {
		query = query.Where("name LIKE ?", params.Search)
	}

	// Apply sorting
	orderClause := fmt.Sprintf("%s %s", params.SortBy, params.SortOrder)

	err := query.
		Order(orderClause).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&projects).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find projects by user: %w", err)
	}

	return projects, nil
}

func (r *projectsRepositoryImpl) Create(ctx context.Context, project *entities.Projects) error {
	return r.db.WithContext(ctx).Create(project).Error
}

func (r *projectsRepositoryImpl) FindByProjectId(ctx context.Context, projectId string) (*entities.Projects, error) {
	// Step 1: Get project
	var project entities.Projects
	err := r.db.
		WithContext(ctx).
		Where("id = ?", projectId).
		Where("deleted_at IS NULL").
		First(&project).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("project not found")
		}
		return nil, fmt.Errorf("failed to find project by id: %w", err)
	}

	// Step 2: Get project expenses
	var expenses []entities.ProjectExpenses
	err = r.db.
		WithContext(ctx).
		Where("project_id = ?", projectId).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&expenses).Error

	if err != nil {
		return nil, fmt.Errorf("failed to load expenses: %w", err)
	}

	// Step 3: Load items untuk setiap expense
	for i := range expenses {
		var items []entities.ProjectExpenseItem
		err = r.db.
			WithContext(ctx).
			Where("project_expense_id = ?", expenses[i].ID).
			Where("deleted_at IS NULL").
			Order("created_at DESC").
			Find(&items).Error

		if err != nil {
			return nil, fmt.Errorf("failed to load expense items: %w", err)
		}

		expenses[i].ProjectExpenseItem = items
	}

	project.ProjectExpenses = expenses

	return &project, nil
}

func (r *projectsRepositoryImpl) FindProjectCategoryByUserID(ctx context.Context, userId string) ([]string, error) {

	return nil, nil
}

func (r *projectsRepositoryImpl) FindAll(ctx context.Context, params entities.GetProjectsParams) ([]*entities.Projects, error) {

	var projects []*entities.Projects

	query := r.db.WithContext(ctx)

	// Filter by category name jika categoryName tidak kosong
	if params.CategoryName != "" {
		query = query.Where("category_name=?", params.CategoryName)
	}

	if params.Search != "" {
		query = query.Where("name=?", params.Search)
	}

	// Apply sorting
	orderClause := fmt.Sprintf("%s %s", params.SortBy, params.SortOrder)

	err := query.
		Order(orderClause).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&projects).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find projects: %w", err)
	}

	return projects, nil
}

func (r *projectsRepositoryImpl) Update(ctx context.Context, project *entities.Projects) error {

	result := r.db.WithContext(ctx).
		Model(project).
		Updates(project)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *projectsRepositoryImpl) Delete(ctx context.Context, id string) error {

	return r.db.WithContext(ctx).Where("id=?", id).Delete(entities.Projects{}).Error
}
