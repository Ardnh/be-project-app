package category

import (
	"project-app/model"
	categoryRepository "project-app/repository/category"

	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CategoryHandler interface {
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	FindAll(c *fiber.Ctx) error
}

type CategoryHandlerImpl struct {
	CategoryRepository categoryRepository.CategoryRepository
	Validator          *validator.Validate
}

func NewCategoryHandler(db *gorm.DB, validate *validator.Validate) CategoryHandler {
	categoryRepository := categoryRepository.NewCategoryRepository(db)
	return &CategoryHandlerImpl{
		CategoryRepository: categoryRepository,
		Validator:          validate,
	}
}

// Create category
// @Summary Create category
// @Description Create a new category
// @Tags Category
// @Accept json
// @Produce json
// @Param body body model.CategoryCreateRequest true "Create category"
// @Success 200 {object} map[string]interface{} "Success create category"
// @Failure 400 {object} map[string]interface{} "Invalid request body or missing required fields"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /category [post]
// @Security Bearer
func (handler *CategoryHandlerImpl) Create(c *fiber.Ctx) error {
	// Logging request start
	logrus.Info("Create category request received")

	// Read body request
	var request model.CategoryCreateRequest
	if err := c.BodyParser(&request); err != nil {
		logrus.WithError(err).Error("Failed to parse request body")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Invalid request body",
		})
	}

	logrus.WithFields(logrus.Fields{
		"category_name": request.Name,
	}).Info("Request body parsed successfully")

	// Validate incoming request
	errValidate := handler.Validator.Struct(request)
	if errValidate != nil {
		logrus.WithError(errValidate).Warn("Validation error in create category request")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": errValidate.Error(),
		})
	}

	// Prepare to create category
	createRequest := model.Category{
		ID:   uuid.New(),
		Name: request.Name,
	}

	logrus.WithFields(logrus.Fields{
		"category_id":   createRequest.ID.String(),
		"category_name": createRequest.Name,
	}).Info("Creating category")

	// Create category in the repository
	err := handler.CategoryRepository.Create(c, &createRequest)
	if err != nil {
		logrus.WithError(err).Error("Failed to create category in the repository")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	logrus.WithFields(logrus.Fields{
		"category_id":   createRequest.ID.String(),
		"category_name": createRequest.Name,
	}).Info("Category created successfully")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Successfully created category",
	})
}

// Update category
// @Summary Update category
// @Description Update category by ID
// @Tags Category
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param body body model.CategoryUpdateRequest true "Update category"
// @Success 200 {object} map[string]interface{} "Success update category"
// @Failure 400 {object} map[string]interface{} "Invalid request body or missing required fields"
// @Failure 404 {object} map[string]interface{} "Category not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /category/{id} [put]
// @Security Bearer
func (handler *CategoryHandlerImpl) Update(c *fiber.Ctx) error {
	// Logging request start
	logrus.Info("Update category request received")

	// Extract ID from URL
	idString := c.Params("id", "")
	if idString == "" {
		logrus.Warn("Invalid ID: empty")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid ID",
		})
	}

	// Parse UUID
	uuidID, errParsing := uuid.Parse(idString)
	if errParsing != nil {
		logrus.WithError(errParsing).Warn("Error parsing ID")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Error parsing id",
		})
	}

	logrus.WithFields(logrus.Fields{
		"category_id": uuidID.String(),
	}).Info("ID parsed successfully")

	// Read body request
	var request model.CategoryUpdateRequest
	if err := c.BodyParser(&request); err != nil {
		logrus.WithError(err).Error("Failed to parse request body")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Invalid request body",
		})
	}

	logrus.WithFields(logrus.Fields{
		"category_name": request.Name,
	}).Info("Request body parsed successfully")

	// Validate incoming request
	errValidate := handler.Validator.Struct(&request)
	if errValidate != nil {
		logrus.WithError(errValidate).Warn("Validation error in update category request")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": errValidate.Error(),
		})
	}

	// Update request object
	updateRequest := &model.Category{
		Name: request.Name,
	}

	logrus.WithFields(logrus.Fields{
		"category_id":   uuidID.String(),
		"category_name": request.Name,
	}).Info("Updating category")

	// Update category in the repository
	errResult := handler.CategoryRepository.Update(c, uuidID, updateRequest)
	if errResult != nil {
		logrus.WithError(errResult).Error("Failed to update category in repository")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": errResult.Error(),
		})
	}

	logrus.WithFields(logrus.Fields{
		"category_id":   uuidID.String(),
		"category_name": request.Name,
	}).Info("Category updated successfully")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Successfully updated category",
	})
}

// Delete category
// @Summary Delete category
// @Description Delete category
// @Tags Category
// @Accept json
// @Produce json
// @Param id path string true "category id"
// @Success 200 {object} map[string]interface{} "Success update category"
// @Failure 400 {object} map[string]interface{} "Invalid request body or missing required fields"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /category/{id} [delete]
// @Security Bearer
func (handler *CategoryHandlerImpl) Delete(c *fiber.Ctx) error {
	// Logging request start
	logrus.Info("Delete category request received")

	// Extract ID from URL
	idString := c.Params("id", "")
	if idString == "" {
		logrus.Warn("Invalid ID: empty")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid ID",
		})
	}

	// Parse UUID
	uuidID, errParsing := uuid.Parse(idString)
	if errParsing != nil {
		logrus.WithFields(logrus.Fields{
			"id":    idString,
			"error": errParsing,
		}).Warn("Error parsing ID")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Error parsing id",
		})
	}

	logrus.WithFields(logrus.Fields{
		"category_id": uuidID.String(),
	}).Info("ID parsed successfully, proceeding with deletion")

	// Delete category in the repository
	errResult := handler.CategoryRepository.Delete(c, uuidID)
	if errResult != nil {
		logrus.WithFields(logrus.Fields{
			"category_id": uuidID.String(),
			"error":       errResult,
		}).Error("Failed to delete category from repository")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": errResult.Error(),
		})
	}

	logrus.WithFields(logrus.Fields{
		"category_id": uuidID.String(),
	}).Info("Category deleted successfully")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Successfully deleted category",
	})
}

// Get all category
// @Summary Get all category
// @Description Get all category
// @Tags Category
// @Accept json
// @Produce json
// @Param page query string false "Page number"
// @Param pageSize query string false "Number of items per page"
// @Param categoryName query string false "Filter by category name"
// @Success 200 {object} map[string]interface{} "Success get all categories"
// @Failure 400 {object} map[string]interface{} "Invalid request body or missing required fields"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /category/ [get]
// @Security Bearer
func (handler *CategoryHandlerImpl) FindAll(c *fiber.Ctx) error {
	// Extract query parameters
	page := c.Query("page", "1")
	pageSize := c.Query("pageSize", "10")
	categoryName := c.Query("categoryName", "")
	sortDirection := c.Query("sortDirection", "")
	sortBy := c.Query("sortBy", "")

	// Convert page and pageSize to integers
	pageInt, errPage := strconv.Atoi(page)
	pageSizeInt, errPageSize := strconv.Atoi(pageSize)

	logrus.WithFields(logrus.Fields{
		"page":          page,
		"pageSize":      pageSize,
		"categoryName":  categoryName,
		"sortDirection": sortDirection,
		"sortBy":        sortBy,
	}).Info("FindAll category request received")

	if errPage != nil || errPageSize != nil {
		logrus.WithFields(logrus.Fields{
			"page":     page,
			"pageSize": pageSize,
		}).Warn("Failed to parse page or pageSize into integers")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid page or pageSize",
		})
	}

	// Call the repository to find categories
	category, totalEntries, errResult := handler.CategoryRepository.FindAll(c, pageInt, pageSizeInt, categoryName)
	if errResult != nil {
		logrus.WithError(errResult).Error("Failed to retrieve categories from repository")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": errResult.Error(),
		})
	}

	// Calculate total pages
	totalPages := int(totalEntries) / pageSizeInt
	if int(totalEntries)%pageSizeInt > 0 {
		totalPages++
	}

	logrus.WithFields(logrus.Fields{
		"totalEntries": totalEntries,
		"totalPages":   totalPages,
		"page":         pageInt,
		"pageSize":     pageSizeInt,
	}).Info("Categories retrieved successfully")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Successfully retrieved categories",
		"data": fiber.Map{
			"items": category,
			"pagination": fiber.Map{
				"page":       pageInt,
				"pageSize":   pageSizeInt,
				"totalPages": totalPages,
				"totalItems": totalEntries,
			},
			"filters": fiber.Map{
				"categoryName":  categoryName,
				"sortDirection": sortDirection,
				"sortBy":        sortBy,
			},
		},
	})
}
