package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

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

	fmt.Println(params.Search)

	query := r.
		db.
		WithContext(ctx).
		Table("projects").
		Where("user_id = ?", userId)

	if params.Search != "" {
		searchPattern := "%" + params.Search + "%" // Tambahkan wildcard
		query = query.Where("name LIKE ? OR category_name LIKE ?", searchPattern, searchPattern)
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
		Preload("ProjectExpenses").
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
		Model(&entities.ProjectExpenses{}).
		Where("project_id = ?", projectId).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&expenses).Error

	if err != nil {
		return nil, fmt.Errorf("failed to load expenses: %w", err)
	}

	project.ProjectExpenses = expenses
	// Step 3: Load items untuk setiap expense
	for i := range expenses {
		var items []entities.ProjectExpenseItems
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

	return &project, nil
}

// Ubah return type ke []string (lebih idiomatic)
func (r *projectsRepositoryImpl) FindProjectCategoryByUserID(ctx context.Context, userId string) ([]*entities.ProjectCategorySummary, error) {

	if userId == "" {
		return nil, errors.New("user_id cannot be empty")
	}

	// Cache key
	cacheKey := fmt.Sprintf("user_categories:%s", userId)

	// Try get from cache
	if r.redis != nil {
		cached, err := r.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var categories []*entities.ProjectCategorySummary
			if err := json.Unmarshal([]byte(cached), &categories); err == nil {
				return categories, nil
			}
		}
	}

	// Query database
	var categories []*entities.ProjectCategorySummary
	err := r.db.WithContext(ctx).
		Model(&entities.Projects{}).
		Select("DISTINCT category_name, COUNT(*) as total").
		Where("user_id = ?", userId).
		Where("category_name IS NOT NULL AND category_name != ''").
		Order("category_name ASC").
		Group("category_name").
		Scan(&categories).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find categories: %w", err)
	}

	// Save to cache
	if r.redis != nil && len(categories) > 0 {
		if data, err := json.Marshal(categories); err == nil {
			r.redis.Set(ctx, cacheKey, data, 5*time.Minute) // TTL 5 menit
		}
	}

	return categories, nil
}

func (r *projectsRepositoryImpl) FindAll(ctx context.Context, params entities.GetProjectsParams) ([]*entities.Projects, error) {

	var projects []*entities.Projects
	query := r.db.WithContext(ctx)

	// Filter by category name jika categoryName tidak kosong
	if params.Search != "" {
		query = query.
			Where("name=?", params.Search).
			Where("category_name=?", params.Search)
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

	return r.db.WithContext(ctx).Where("id=?", id).Delete(&entities.Projects{}).Error
}
