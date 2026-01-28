package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
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

func (r *projectsRepositoryImpl) FindByUserID(ctx context.Context, userId string, params entities.GetProjectsParams) ([]*entities.ProjectsByUserId, error) {
	var projects []*entities.ProjectsByUserId

	// Base CTE query
	query := r.db.WithContext(ctx).Raw(`
		WITH todolist_stats AS (
		    SELECT
		        p.id AS project_id,
		        COUNT(DISTINCT pt.id) AS total_todolist,
		        COUNT(pti.id) AS total_todolist_item,
		        COUNT(*) FILTER (WHERE pti.is_completed = true) AS total_todolist_item_done
		    FROM projects p
		    LEFT JOIN project_todolists pt ON p.id = pt.project_id
		    LEFT JOIN project_todolist_items pti ON pt.id = pti.project_todolist_id
		    WHERE p.user_id = $1
		    GROUP BY p.id
		),
		expense_stats AS (
		    SELECT
		        p.id AS project_id,
		        COUNT(pe.id) AS total_expenses
		    FROM projects p
		    LEFT JOIN project_expenses pe ON p.id = pe.project_id
		    WHERE p.user_id = $2
		    GROUP BY p.id
		)
		SELECT
		    p.id AS project_id,
		    p.user_id,
		    p.name,
		    p.budget,
		    p.is_completed,
		    p.category_name,
		    p.start_date,
		    p.end_date,
		    p.created_at,
		    COALESCE(ts.total_todolist, 0) AS total_todolist,
		    COALESCE(ts.total_todolist_item_done, 0) AS total_todolist_item_done,
		    COALESCE(ts.total_todolist_item, 0) AS total_todolist_item,
		    COALESCE(es.total_expenses, 0) AS total_expenses,
		    CASE
		        WHEN COALESCE(ts.total_todolist_item, 0) = 0 THEN 0
		        ELSE ROUND(
		            (ts.total_todolist_item_done::NUMERIC /
		             NULLIF(ts.total_todolist_item, 0)) * 100,
		            2
		        )
		    END AS completion_percentage
		FROM projects p
		LEFT JOIN todolist_stats ts ON p.id = ts.project_id
		LEFT JOIN expense_stats es ON p.id = es.project_id
		WHERE p.user_id = $3
	`, userId, userId, userId) // userId dipakai 3x

	// Apply search filter
	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("p.name ILIKE ? OR p.category_name ILIKE ?", searchPattern, searchPattern)
	}

	// Apply sorting
	if params.SortBy != "" && params.SortOrder != "" {
		sortField := params.SortBy
		// Map field names
		validSortFields := map[string]string{
			"name":                  "p.name",
			"category_name":         "p.category_name",
			"created_at":            "p.created_at",
			"budget":                "p.budget",
			"completion_percentage": "completion_percentage",
		}

		if mappedField, ok := validSortFields[sortField]; ok {
			sortField = mappedField
		}

		// Validate sort order
		sortOrder := strings.ToUpper(params.SortOrder)
		if sortOrder != "ASC" && sortOrder != "DESC" {
			sortOrder = "DESC"
		}

		query = query.Order(fmt.Sprintf("%s %s", sortField, sortOrder))
	} else {
		// Default sorting
		query = query.Order("p.created_at DESC")
	}

	// Apply pagination
	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}
	if params.Offset > 0 {
		query = query.Offset(params.Offset)
	}

	// Execute query
	err := query.Scan(&projects).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find projects by user: %w", err)
	}

	return projects, nil
}

func (r *projectsRepositoryImpl) FindProjectSummaryByUserID(ctx context.Context, userId string) (*entities.ProjectUserSummary, error) {
	var projects entities.ProjectUserSummary

	// err := query.Find(&projects).Error
	err := r.db.Model(&entities.Projects{}).
		Select(
			"COALESCE(SUM(budget), 0) as total_budget",
			"COUNT(id) as total_project",
			"COUNT(CASE WHEN is_completed = true THEN 1 END) as total_completed_project",
		).
		Where("user_id = ?", userId).
		Scan(&projects).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find projects by user: %w", err)
	}

	return &projects, nil
}

func (r *projectsRepositoryImpl) Create(ctx context.Context, project *entities.Projects) error {

	cacheKey := fmt.Sprintf("user_categories:%s", project.UserID)
	if err := r.db.WithContext(ctx).Create(project).Error; err != nil {
		return err
	}

	if err := r.redis.Del(ctx, cacheKey).Err(); err != nil {
		log.Printf("failed to delete cache %s: %v", cacheKey, err)
	}

	return nil
}

func (r *projectsRepositoryImpl) FindByProjectId(ctx context.Context, projectId string) (*entities.Projects, error) {
	// Step 1: Get project
	var project entities.Projects
	err := r.db.
		WithContext(ctx).
		Model(&entities.Projects{}).
		Preload("ProjectExpenses").
		Preload("ProjectTodolists").
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

	// Step 4: Get project todolist
	var todolists []entities.ProjectTodolists
	err = r.db.
		WithContext(ctx).
		Model(&entities.ProjectTodolists{}).
		Where("project_id = ?", projectId).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&todolists).Error

	if err != nil {
		return nil, fmt.Errorf("failed to load expenses: %w", err)
	}

	project.ProjectTodolists = todolists
	// Step 5: Load items untuk setiap todolist
	for i := range todolists {
		var items []entities.ProjectTodolistItems
		err = r.db.
			WithContext(ctx).
			Model(&entities.ProjectTodolistItems{}).
			Where("project_todolist_id = ?", todolists[i].ID).
			Where("deleted_at IS NULL").
			Order("created_at DESC").
			Find(&items).Error

		if err != nil {
			return nil, fmt.Errorf("failed to load todolist items: %w", err)
		}

		todolists[i].ProjectTodolistItems = items
	}

	return &project, nil
}

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
	// Check is user exist
	var user *entities.User
	errUser := r.db.WithContext(ctx).Where("id = ?", userId).First(&user).Error
	if errUser != nil {
		return nil, fmt.Errorf("failed to find user: %w", errUser)
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	var categories []*entities.ProjectCategorySummary
	err := r.db.WithContext(ctx).
		Model(&entities.Projects{}).
		Select("category_name, COUNT(*) as total").
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
			Where("name LIKE %?%", params.Search).
			Where("category_name LIKE %?%", params.Search)
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

	cacheKey := fmt.Sprintf("user_categories:%s", project.UserID)
	if err := r.db.WithContext(ctx).Create(project).Error; err != nil {
		return err
	}

	if err := r.redis.Del(ctx, cacheKey).Err(); err != nil {
		log.Printf("failed to delete cache %s: %v", cacheKey, err)
	}

	return nil
}

func (r *projectsRepositoryImpl) Delete(ctx context.Context, id string) error {
	// 1. Ambil project dulu untuk dapatkan UserID
	var project entities.Projects
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&project).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("project not found")
		}
		return fmt.Errorf("failed to find project: %w", err)
	}

	// 2. Delete project dari database
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.Projects{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete project: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("no project was deleted")
	}

	// 3. Hapus cache Redis
	cacheKey := fmt.Sprintf("user_categories:%s", project.UserID)

	// Log untuk debug
	log.Printf("Deleting cache for UserID: %s, Key: %s", project.UserID, cacheKey)

	deletedCount, err := r.redis.Del(ctx, cacheKey).Result()
	if err != nil {
		log.Printf("Failed to delete cache %s: %v", cacheKey, err)
	} else {
		log.Printf("Successfully deleted %d cache key(s): %s", deletedCount, cacheKey)
		if deletedCount == 0 {
			log.Printf("Warning: Cache key %s not found in Redis", cacheKey)
		}
	}

	return nil
}
