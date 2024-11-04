package project

import (
	"project-app/helper"
	"project-app/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
	GetAllProjectItemByProjectId(ctx *fiber.Ctx, page int, pageSize int, sortOrder string, projectItemName string, projectID uuid.UUID) (*model.ProjectItemResponse, int64, error)
}

type ProjectRepositoryImpl struct {
	Db *gorm.DB
}

var projectTableName = "projects"
var projectItemTableName = "project_items"

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

	// 1. Ambil data project lama
	var existingProject model.Projects
	if err := tx.
		WithContext(ctx.Context()).
		Table(projectTableName).
		Where("id = ?", req.ID).
		First(&existingProject).
		Error; err != nil {
		return err
	}

	// 2. Cek apakah budget project lama sama dengan budget project baru
	if existingProject.Budget == req.Budget {
		// Budget sama, update project seperti biasa
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

	// 3. Budget berbeda, lakukan rekalkulasi
	// 3.1 Ambil semua project item yang memiliki project_id yang sama dengan req.ID
	var projectItems []model.ProjectItem
	if err := tx.
		WithContext(ctx.Context()).
		Table("project_items").
		Where("project_id = ?", req.ID).
		Find(&projectItems).
		Error; err != nil {
		return err
	}

	// 3.2 Jumlahkan semua budgetItem
	var totalBudgetItem int
	for _, item := range projectItems {
		totalBudgetItem += item.BudgetItem
	}

	// 3.3 Kurangi budget baru dengan total budget item
	newBudget := req.Budget - totalBudgetItem

	// 3.4 Update budget project dari hasil pengurangan budget baru - total budget item
	req.Budget = newBudget

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
	// Mulai transaksi
	tx := repository.Db.Begin()
	defer helper.CommitOrRollback(tx)

	// Log awal proses penghapusan project
	logrus.WithField("ProjectID", req).Info("Deleting project and related items")

	// Hapus project item berdasarkan project ID
	errDeleteProjectItem := tx.
		WithContext(ctx.Context()).
		Table(projectItemTableName).
		Where("project_id = ?", req). // Pastikan filter by project_id
		Delete(&model.ProjectItem{}).
		Error

	if errDeleteProjectItem != nil {
		logrus.WithError(errDeleteProjectItem).Error("Failed to delete project items")
		return errDeleteProjectItem
	}

	// Log sukses menghapus project items
	logrus.WithField("ProjectID", req).Info("Project items deleted successfully")

	// Hapus project berdasarkan project ID
	errDeleteProject := tx.
		WithContext(ctx.Context()).
		Table(projectTableName).
		Where("id = ?", req). // Pastikan filter by project ID
		Delete(&model.Projects{}).
		Error

	if errDeleteProject != nil {
		logrus.WithError(errDeleteProject).Error("Failed to delete project")
		return errDeleteProject
	}

	// Log sukses menghapus project
	logrus.WithField("ProjectID", req).Info("Project deleted successfully")

	return nil
}

func (repository *ProjectRepositoryImpl) GetAllProject(ctx *fiber.Ctx, page int, pageSize int, sortOrder string, projectName string, categoryName string) (*[]model.Projects, int64, error) {
	tx := repository.Db.WithContext(ctx.Context())
	defer helper.CommitOrRollback(tx)

	var projects []model.Projects
	var totalCount int64

	// Offset calculation
	offset := (page - 1) * pageSize

	// Base query with Model and Preload
	query := tx.Model(&model.Projects{}).
		Joins("JOIN categories ON categories.id = projects.category_id").
		Preload("Category")

	// Apply filters if projectName or categoryName is provided
	if projectName != "" || categoryName != "" {
		query = query.Where("projects.name ILIKE ? AND categories.name ILIKE ?", "%"+projectName+"%", "%"+categoryName+"%")
	}

	// Count total records
	err := query.Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply sorting based on budget
	switch sortOrder {
	case "asc":
		query = query.Order("budget ASC")
	case "desc":
		query = query.Order("budget DESC")
	default:
		query = query.Order("budget ASC") // Default sorting
	}

	// Pagination and fetch results with related Category
	err = query.Offset(offset).Limit(pageSize).Find(&projects).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &projects, 0, nil
		}
		return nil, 0, err
	}

	return &projects, totalCount, nil
}

func (repository *ProjectRepositoryImpl) CreateProjectItem(ctx *fiber.Ctx, req *model.ProjectItem) error {

	tx := repository.Db.WithContext(ctx.Context()).Begin()
	defer helper.CommitOrRollback(tx)

	// 1. Insert new project item to project_item table
	if err := tx.Table(projectItemTableName).Create(req).Error; err != nil {
		return err
	}

	// 2. Update current budget langsung di table projects tanpa mengambil data project terlebih dahulu
	if err := tx.Table(projectTableName).
		Where("id = ? AND budget >= ?", req.ProjectID, req.BudgetItem).
		UpdateColumn("budget", gorm.Expr("budget - ?", req.BudgetItem)).
		Error; err != nil {
		return err
	}

	return nil
}

func (repository *ProjectRepositoryImpl) UpdateProjectItem(ctx *fiber.Ctx, req *model.ProjectItem) error {

	tx := repository.Db.WithContext(ctx.Context()).Begin()
	defer helper.CommitOrRollback(tx)

	var oldProjectItem model.ProjectItem

	// 1. Get the old project item to calculate the budget difference
	err := tx.Table(projectItemTableName).
		Where("id = ?", req.ID).
		First(&oldProjectItem).
		Error

	if err != nil {
		return err
	}

	// 2. Update project item with the new data
	err = tx.Table(projectItemTableName).
		Where("id = ?", req.ID).
		Updates(req).
		Error

	if err != nil {
		return err
	}

	// 3. Calculate the budget difference
	budgetDifference := req.BudgetItem - oldProjectItem.BudgetItem

	// 4. Update the project's budget in the projects table
	err = tx.Table(projectTableName).
		Where("id = ?", req.ProjectID).
		UpdateColumn("budget", gorm.Expr("budget + ?", budgetDifference)).
		Error

	if err != nil {
		return err
	}

	return nil
}

func (repository *ProjectRepositoryImpl) DeleteProjectItem(ctx *fiber.Ctx, req uuid.UUID) error {

	tx := repository.Db.WithContext(ctx.Context()).Begin()
	defer helper.CommitOrRollback(tx)

	var projectItem model.ProjectItem

	// 1. Ambil project item berdasarkan req
	if err := tx.Table(projectItemTableName).
		Where("id = ?", req).
		First(&projectItem).
		Error; err != nil {
		return err
	}

	// 2. Jika budget pada project item tidak 0, tambahkan budget pada project item ke parent project
	if projectItem.BudgetItem != 0 {
		if err := tx.Table(projectTableName).
			Where("id = ?", projectItem.ProjectID).
			UpdateColumn("budget", gorm.Expr("budget + ?", projectItem.BudgetItem)).
			Error; err != nil {
			return err
		}
	}

	// 3. Hapus project item
	if err := tx.Delete(&projectItem).Error; err != nil {
		return err
	}

	return nil
}

func (repository *ProjectRepositoryImpl) GetAllProjectItemByProjectId(
	ctx *fiber.Ctx,
	page int,
	pageSize int,
	sortOrder string,
	projectItemName string,
	projectID uuid.UUID,
) (*model.ProjectItemResponse, int64, error) {

	// var projectItems []model.ProjectItem
	var totalCount int64
	var projectItemDetails model.ProjectItemResponse

	// Validasi Project ID dan ambiguitas id
	if err := repository.Db.WithContext(ctx.Context()).
		Table(projectTableName).
		Model(&model.Projects{}).
		Joins("JOIN categories ON categories.id = projects.category_id").
		Preload("Category").
		Where("projects.id = ?", projectID). // Tambahkan projects.id
		First(&projectItemDetails.Project).Error; err != nil {
		logrus.WithError(err).Error("Failed to fetch project details")
		return nil, 0, err
	}

	// Hitung offset untuk pagination
	offset := (page - 1) * pageSize

	// Query untuk mendapatkan project item berdasarkan project_id
	query := repository.Db.WithContext(ctx.Context()).
		Table(projectItemTableName).
		Where("project_id = ?", projectID).
		Offset(offset).
		Limit(pageSize)

	// Jika projectItemName tidak kosong, tambahkan filter pencarian
	if projectItemName != "" {
		query = query.Where("name LIKE ?", "%"+projectItemName+"%")
	}

	// Sorting berdasarkan budget item
	switch sortOrder {
	case "asc":
		query = query.Order("budget_item ASC")
	case "desc":
		query = query.Order("budget_item DESC")
	default:
		query = query.Order("budget_item ASC") // Default sorting
	}

	// Hitung total item
	if err := query.Count(&totalCount).Error; err != nil {
		logrus.WithError(err).Error("Failed to count project items")
		return nil, 0, err
	}

	// Ambil data project items
	if err := query.Find(&projectItemDetails.ProjectItems).Error; err != nil {
		logrus.WithError(err).Error("Failed to fetch project items")
		return nil, 0, err
	}

	// Log jumlah project items yang ditemukan
	logrus.WithFields(logrus.Fields{
		"totalItems": totalCount,
		"projectID":  projectID,
	}).Info("Successfully fetched project items")

	// Kembalikan response
	return &projectItemDetails, totalCount, nil
}
