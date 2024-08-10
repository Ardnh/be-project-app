package proeject

import (
	"project-app/model"
	projectRepository "project-app/repository/project"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectHandler interface {
	CreateProject(c *fiber.Ctx) error
	UpdateProject(c *fiber.Ctx) error
	DeleteProject(c *fiber.Ctx) error
	GetAllProject(c *fiber.Ctx) error

	CreateProjectItem(c *fiber.Ctx) error
	UpdateProjectItem(c *fiber.Ctx) error
	DeleteProjectItem(c *fiber.Ctx) error
	GetAllProjectItemByProjectId(c *fiber.Ctx) error
}

type ProjectHandlerImpl struct {
	ProjectRepository projectRepository.ProjectRepository
	Validator         *validator.Validate
}

func NewProjectHandler(db *gorm.DB, validate *validator.Validate) ProjectHandler {
	projectRepository := projectRepository.NewProjectRepository(db)

	return &ProjectHandlerImpl{
		ProjectRepository: projectRepository,
		Validator:         validate,
	}
}

// CreateProject
// @Summary Create a new project
// @Description Create a new project with the given details
// @Tags Projects
// @Accept  json
// @Produce  json
// @Param project body model.CreateProjectRequest true "Create Project Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /projects [post]
// @Security Bearer
func (handler *ProjectHandlerImpl) CreateProject(c *fiber.Ctx) error {

	// Read body request
	var request model.CreateProjectRequest

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

	// Create Project
	createRequest := model.Projects{
		Name:        request.Name,
		CategoryID:  request.CategoryID,
		Description: request.Description,
		Budget:      request.Budget,
		UserID:      request.UserID,
	}

	err := handler.ProjectRepository.CreateProject(c, &createRequest)
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

func (handler *ProjectHandlerImpl) UpdateProject(c *fiber.Ctx) error {

	// Read body request
	var request model.UpdateProjectRequest
	idString := c.Params("id", "")
	uuidID, err := uuid.Parse(idString)

	if idString == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid ID",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Error parsing id",
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
	updateRequest := &model.Projects{
		ID:          uuidID,
		Name:        request.Name,
		CategoryID:  request.CategoryID,
		UserID:      request.UserID,
		Description: request.Description,
		Budget:      request.Budget,
	}

	errResult := handler.ProjectRepository.UpdateProject(c, updateRequest)
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

func (handler *ProjectHandlerImpl) DeleteProject(c *fiber.Ctx) error {

	idString := c.Params("id", "")
	uuidID, err := uuid.Parse(idString)

	if idString == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid ID",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Error parsing id",
		})
	}

	errResult := handler.ProjectRepository.DeleteProject(c, uuidID)

	if errResult != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": errResult.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Successfully delete project",
	})
}

func (handler *ProjectHandlerImpl) GetAllProject(c *fiber.Ctx) error {

	page := c.Query("page", "1")
	pageInt, _ := strconv.Atoi(page)
	pageSize := c.Query("pageSize", "10")
	pageSizeInt, _ := strconv.Atoi(pageSize)
	projectName := c.Query("projectName", "")
	sortOrder := c.Query("sortDirection", "asc")
	categoryName := c.Query("categoryName", "")

	projects, totalItem, errResult := handler.ProjectRepository.GetAllProject(c, pageInt, pageSizeInt, sortOrder, projectName, categoryName)

	if errResult != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": errResult.Error(),
		})
	}

	totalPages := int(totalItem) / pageSizeInt
	if int(totalItem)%pageSizeInt > 0 { // Tambahkan satu halaman jika ada sisa pembagian
		totalPages++
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":       fiber.StatusOK,
		"message":    "Successfully get category",
		"page":       page,
		"pageSize":   pageSize,
		"total":      totalItem,
		"totalPages": totalPages,
		"data":       projects,
	})
}

func (handler *ProjectHandlerImpl) CreateProjectItem(c *fiber.Ctx) error {

	var request model.CreateProjectItem
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

	createRequest := model.ProjectItem{
		ProjectID:  request.ProjectID,
		Name:       request.Name,
		BudgetItem: request.BudgetItem,
		Status:     request.Status,
	}

	errResult := handler.ProjectRepository.CreateProjectItem(c, &createRequest)
	if errResult != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": errResult.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Successfully created project item",
	})

}

func (handler *ProjectHandlerImpl) UpdateProjectItem(c *fiber.Ctx) error {

	// Read body request
	var request model.UpdateProjectItem
	idString := c.Params("id", "")
	uuidID, err := uuid.Parse(idString)

	if idString == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid ID",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Error parsing id",
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

	updateRequest := model.ProjectItem{
		ID:         uuidID,
		ProjectID:  request.ProjectID,
		Name:       request.Name,
		BudgetItem: request.BudgetItem,
		Status:     request.Status,
	}

	errResult := handler.ProjectRepository.UpdateProjectItem(c, &updateRequest)
	if errResult != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": errResult.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Successfully update project item",
	})
}

func (handler *ProjectHandlerImpl) DeleteProjectItem(c *fiber.Ctx) error {

	idString := c.Params("id", "")
	uuidID, err := uuid.Parse(idString)

	if idString == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid ID",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Error parsing id",
		})
	}

	errResult := handler.ProjectRepository.DeleteProjectItem(c, uuidID)
	if errResult != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": errResult.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Successfully Delete project item",
	})
}

func (handler *ProjectHandlerImpl) GetAllProjectItemByProjectId(c *fiber.Ctx) error {

	page := c.Query("page", "1")
	pageInt, _ := strconv.Atoi(page)
	pageSize := c.Query("pageSize", "10")
	pageSizeInt, _ := strconv.Atoi(pageSize)
	projectName := c.Query("projectName", "")
	sortOrder := c.Query("sortDirection", "asc")

	projectItems, totalItem, errResult := handler.ProjectRepository.GetAllProjectItemByProjectId(c, pageInt, pageSizeInt, sortOrder, projectName)

	if errResult != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": errResult.Error(),
		})
	}

	totalPages := int(totalItem) / pageSizeInt
	if int(totalItem)%pageSizeInt > 0 { // Tambahkan satu halaman jika ada sisa pembagian
		totalPages++
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":       fiber.StatusOK,
		"message":    "Successfully get category",
		"page":       page,
		"pageSize":   pageSize,
		"total":      totalItem,
		"totalPages": totalPages,
		"data":       projectItems,
	})

}
