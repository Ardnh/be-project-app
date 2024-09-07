package category

import (
	"fmt"
	"project-app/helper"
	"project-app/model"
	categoryModel "project-app/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(ctx *fiber.Ctx, req *model.Category) error
	Update(ctx *fiber.Ctx, id uuid.UUID, req *model.Category) error
	Delete(ctx *fiber.Ctx, id uuid.UUID) error
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

func (repository *CategoryRepositoryImpl) Update(ctx *fiber.Ctx, id uuid.UUID, req *model.Category) error {

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

func (repository *CategoryRepositoryImpl) Delete(ctx *fiber.Ctx, id uuid.UUID) error {

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

	fmt.Println("search query")
	fmt.Println(searchQuery)

	tx := repository.Db.Begin()
	defer helper.CommitOrRollback(tx)

	// Offset
	offset := (page - 1) * pageSize

	// Query
	query := tx.WithContext(ctx.Context()).Table(tableName)

	if searchQuery != "" {
		query = query.Where("name ILIKE ? ", "%"+searchQuery+"%")
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

		if errResult == gorm.ErrRecordNotFound {
			return []categoryModel.Category{}, 0, nil
		}

		return nil, 0, errResult
	}

	return category, totalCount, nil
}
