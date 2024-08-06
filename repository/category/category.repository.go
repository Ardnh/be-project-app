package category

import (
	"project-app/helper"
	"project-app/model"
	categoryModel "project-app/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(ctx *fiber.Ctx, req *model.Category) error
	Update(ctx *fiber.Ctx, id int, req *model.Category) error
	Delete(ctx *fiber.Ctx, id int) error
	FindAll(ctx *fiber.Ctx, page int, pageSize int, searchQuery string) ([]categoryModel.Category, int64, error)
}

type CategoryRepositoryImpl struct {
	Db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &CategoryRepositoryImpl{
		Db: db,
	}
}

var tableName = "categories"

func (repository *CategoryRepositoryImpl) Create(ctx *fiber.Ctx, req *model.Category) error {

	tx := repository.Db.Begin()
	defer helper.CommitOrRollback(tx)

	err := tx.
		WithContext(ctx.Context()).
		Table(tableName).
		Create(req).
		Error

	if err != nil {
		return err
	}

	return nil
}

func (repository *CategoryRepositoryImpl) Update(ctx *fiber.Ctx, id int, req *model.Category) error {

	tx := repository.Db.Begin()
	defer helper.CommitOrRollback(tx)

	err := tx.WithContext(ctx.Context()).
		Table(tableName).
		Where("id = ? ", id).
		Updates(req).
		Error

	if err != nil {
		return err
	}

	return nil
}

func (repository *CategoryRepositoryImpl) Delete(ctx *fiber.Ctx, id int) error {

	tx := repository.Db.Begin()
	defer helper.CommitOrRollback(tx)

	err := tx.WithContext(ctx.Context()).
		Table(tableName).
		Delete(&categoryModel.Category{}, id).Error

	if err != nil {
		return err
	}

	return nil
}

func (repository *CategoryRepositoryImpl) FindAll(ctx *fiber.Ctx, page int, pageSize int, searchQuery string) ([]categoryModel.Category, int64, error) {

	var category []categoryModel.Category
	var totalCount int64

	tx := repository.Db.Begin()
	defer helper.CommitOrRollback(tx)

	// Offset
	offset := (page - 1) * pageSize

	// Query
	query := tx.WithContext(ctx.Context()).Table(tableName)

	if searchQuery != "" {
		query = query.Where("category LIKE ? ", "%"+searchQuery+"%")
	}

	err := query.Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	errResult := query.
		Offset(offset).
		Limit(pageSize).
		Find(&category).
		Error

	if errResult != nil {
		return nil, 0, errResult
	}

	return category, totalCount, nil
}
