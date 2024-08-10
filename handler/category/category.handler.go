package category

import (
	"project-app/model"
	categoryRepository "project-app/repository/category"

	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
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

	// Read body request
	var request model.CategoryCreateRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	// Validate incoming request
	errValidate := handler.Validator.Struct(request)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": errValidate.Error(),
		})
	}

	// Create Category
	createRequest := model.Category{
		Name: request.Name,
	}

	err := handler.CategoryRepository.Create(c, &createRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Successfully create category",
	})
}

// Update category
// @Summary Update category
// @Description Update category
// @Tags Category
// @Accept json
// @Produce json
// @Param body body model.CategoryUpdateRequest true "Update category"
// @Success 200 {object} map[string]interface{} "Success update category"
// @Failure 400 {object} map[string]interface{} "Invalid request body or missing required fields"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /category/{} [put]
// @Security Bearer
func (handler *CategoryHandlerImpl) Update(c *fiber.Ctx) error {

	// Read body request
	var request model.CategoryUpdateRequest
	idString := c.Params("id", "")
	idInt, errConv := strconv.Atoi(idString)

	if errConv != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": errConv.Error(),
		})
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	// Validate incoming request
	errValidate := handler.Validator.Struct(&request)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": errValidate.Error(),
		})
	}

	// Update request
	updateRequest := &model.Category{
		Name: request.Name,
	}

	errResult := handler.CategoryRepository.Update(c, idInt, updateRequest)
	if errResult != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": errResult.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Successfully update category",
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

	// Read body request
	idString := c.Params("id")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": err.Error(),
		})
	}

	errResult := handler.CategoryRepository.Delete(c, idInt)
	if errResult != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": errResult.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Successfully delete category",
	})
}

// Get all category
// @Summary Get all category
// @Description Get all category
// @Tags Category
// @Accept json
// @Produce json
// @Param page path string false "page"
// @Param pageSize path string false "pageSize"
// @Param categoryName path string false "categoryName"
// @Success 200 {object} map[string]interface{} "Success update category"
// @Failure 400 {object} map[string]interface{} "Invalid request body or missing required fields"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /category/ [get]
// @Security Bearer
func (handler *CategoryHandlerImpl) FindAll(c *fiber.Ctx) error {

	page := c.Query("page", "1")
	pageInt, _ := strconv.Atoi(page)
	pageSize := c.Query("pageSize", "10")
	pageSizeInt, _ := strconv.Atoi(pageSize)
	categoryName := c.Query("categoryName", "")

	category, totalEntries, errResult := handler.CategoryRepository.FindAll(c, pageInt, pageSizeInt, categoryName)
	if errResult != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": errResult.Error(),
		})
	}

	totalPages := int(totalEntries) / pageSizeInt
	if int(totalEntries)%pageSizeInt > 0 {
		totalPages++
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":       fiber.StatusOK,
		"message":    "Successfully get category",
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": totalEntries,
		"data":       category,
	})
}
