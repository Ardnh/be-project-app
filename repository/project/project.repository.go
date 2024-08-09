package project

import (
	"errors"
	"project-app/helper"
	"project-app/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	CreateProject(ctx *fiber.Ctx, req *model.Projects) error
	UpdateProject(ctx *fiber.Ctx, req *model.Projects) error
	DeleteProject(ctx *fiber.Ctx, req uuid.UUID) error
	GetAllProject(ctx *fiber.Ctx, page int, pageSize int, sortOrder string, projectName string, categoryName string) (*[]model.Projects, int64, error)

	CreateProjectItem(ctx *fiber.Ctx, req *model.ProjectItem) error
	UpdateProjectItem(ctx *fiber.Ctx, req *model.ProjectItem) error
	DeleteProjectItem(ctx *fiber.Ctx, req uuid.UUID) error
	GetAllProjectItemByProjectId(ctx *fiber.Ctx, page int, pageSize int, sortOrder string, projectItemName string) (*[]model.ProjectItem, int64, error)
}

type ProjectRepositoryImpl struct {
	Db *gorm.DB
}

var projectTableName = "project"
var projectItemTableName = "project_item"

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &ProjectRepositoryImpl{
		Db: db,
	}
}

func (repository *ProjectRepositoryImpl) CreateProject(ctx *fiber.Ctx, req *model.Projects) error {
	tx := repository.Db.Begin()
	defer helper.CommitOrRollback(tx)

	err := tx.
		WithContext(ctx.Context()).
		Table(projectTableName).
		Create(req).
		Error

	if err != nil {
		return err
	}

	return nil
}

func (repository *ProjectRepositoryImpl) UpdateProject(ctx *fiber.Ctx, req *model.Projects) error {

	tx := repository.Db.Begin()
	defer helper.CommitOrRollback(tx)

	err := tx.
		WithContext(ctx.Context()).
		Table(projectTableName).
		Where("id = ?", req.ID).
		Updates(req).
		Error

	if err != nil {
		return err
	}

	return nil
}

func (repository *ProjectRepositoryImpl) DeleteProject(ctx *fiber.Ctx, req uuid.UUID) error {

	tx := repository.Db.Begin()
	defer helper.CommitOrRollback(tx)

	// hapus project item berdasarkan project ID
	errDeleteProjectItem := tx.
		WithContext(ctx.Context()).
		Table(projectItemTableName).
		Delete(&model.ProjectItem{}, req).
		Error

	if errDeleteProjectItem != nil {
		return errDeleteProjectItem
	}

	// hapus project berdasarkan project ID
	errDeleteProject := tx.
		WithContext(ctx.Context()).
		Table(projectItemTableName).
		Delete(&model.Projects{}, req).
		Error

	if errDeleteProject != nil {
		return errDeleteProject
	}

	return nil
}

func (repository *ProjectRepositoryImpl) GetAllProject(ctx *fiber.Ctx, page int, pageSize int, sortOrder string, projectName string, categoryName string) (*[]model.Projects, int64, error) {

	tx := repository.Db.Begin()
	defer helper.CommitOrRollback(tx)

	var projects []model.Projects
	var totalCount int64

	// Offset
	offset := (page - 1) * pageSize

	// Query
	query := tx.WithContext(ctx.Context()).Table(projectTableName)

	if projectName != "" || categoryName != "" {
		query = query.
			Joins("Join categories ON categories.id = projects.category_id").
			Where("projects.name LIKE ? OR categories.name LIKE ? ", "%"+projectName+"%", "%"+categoryName+"%")
	}

	err := query.Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	// Sorting by budget
	if sortOrder == "asc" {
		query = query.Order("budget ASC")
	} else if sortOrder == "desc" {
		query = query.Order("budget DESC")
	} else {
		query = query.Order("budget ASC") // Default sorting
	}

	errResult := query.
		Offset(offset).
		Limit(pageSize).
		Find(&projects).
		Error

	if errResult != nil {
		return nil, 0, errResult
	}

	if len(projects) == 0 {
		return nil, 0, errors.New("Project Not Found")
	}

	return &projects, totalCount, nil
}

func (repository *ProjectRepositoryImpl) CreateProjectItem(ctx *fiber.Ctx, req *model.ProjectItem) error {

	tx := repository.Db.Begin()
	defer helper.CommitOrRollback(tx)

	err := tx.
		WithContext(ctx.Context()).
		Table(projectItemTableName).
		Create(req).
		Error

	if err != nil {
		return err
	}

	return nil

}

func (repository *ProjectRepositoryImpl) UpdateProjectItem(ctx *fiber.Ctx, req *model.ProjectItem) error {

	tx := repository.Db.Begin()
	defer helper.CommitOrRollback(tx)

	err := tx.
		WithContext(ctx.Context()).
		Table(projectItemTableName).
		Where("id = ?", req.ID).
		Updates(req).
		Error

	if err != nil {
		return err
	}

	return nil

}

func (repository *ProjectRepositoryImpl) DeleteProjectItem(ctx *fiber.Ctx, req uuid.UUID) error {

	tx := repository.Db.Begin()
	defer helper.CommitOrRollback(tx)

	errDeleteProjectItem := tx.
		WithContext(ctx.Context()).
		Table(projectItemTableName).
		Delete(&model.ProjectItem{}, req).
		Error

	if errDeleteProjectItem != nil {
		return errDeleteProjectItem
	}

	return nil
}

func (repository *ProjectRepositoryImpl) GetAllProjectItemByProjectId(ctx *fiber.Ctx, page int, pageSize int, sortOrder string, projectItemName string) (*[]model.ProjectItem, int64, error) {

	tx := repository.Db.Begin()
	defer helper.CommitOrRollback(tx)

	var projectItems []model.ProjectItem
	var totalCount int64

	// Offset
	offset := (page - 1) * pageSize

	// Query
	query := tx.WithContext(ctx.Context()).Table(projectItemTableName)

	if projectItemName != "" {
		query = query.Where("name LIKE ? ", "%"+projectItemName+"%")
	}

	err := query.Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	// Sorting by budget
	if sortOrder == "asc" {
		query = query.Order("budget ASC")
	} else if sortOrder == "desc" {
		query = query.Order("budget DESC")
	} else {
		query = query.Order("budget ASC") // Default sorting
	}

	errResult := query.
		Offset(offset).
		Limit(pageSize).
		Find(&projectItems).
		Error

	if errResult != nil {
		return nil, 0, errResult
	}

	if len(projectItems) == 0 {
		return nil, 0, errors.New("Project Item Not Found")
	}

	return &projectItems, totalCount, nil
}
