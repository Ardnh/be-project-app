package proeject

import (
	"fmt"
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

// CreateProject godoc
// @Summary Create a new project
// @Description Create a new item within an existing project
// @Tags Projects
// @Accept  json
// @Produce  json
// @Param projectItem body model.CreateProjectRequest true "Create Project Request"
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
		ID:          uuid.New(),
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
		"message": "Successfully create project",
	})
}

// UpdateProject godoc
// @Summary Update an existing project
// @Description Update the details of an existing project by ID
// @Tags Projects
// @Accept  json
// @Produce  json
// @Param id path string true "Project ID"
// @Param project body model.UpdateProjectRequest true "Update Project Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /projects/{id} [put]
// @Security Bearer
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
		"message": "Successfully update project",
	})
}

// DeleteProject godoc
// @Summary Delete an existing project
// @Description Delete the project by its ID
// @Tags Projects
// @Accept  json
// @Produce  json
// @Param id path string true "Project ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /projects/{id} [delete]
// @Security Bearer
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

// GetAllProject godoc
// @Summary Get all projects with pagination, sorting, and filtering
// @Description Retrieve a paginated list of projects with optional sorting and filtering by project name and category name
// @Tags Projects
// @Accept  json
// @Produce  json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Number of items per page" default(10)
// @Param projectName query string false "Filter by project name"
// @Param sortDirection query string false "Sort order (asc or desc)" default(asc)
// @Param categoryName query string false "Filter by category name"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /projects [get]
// @Security Bearer
func (handler *ProjectHandlerImpl) GetAllProject(c *fiber.Ctx) error {

	page := c.Query("page", "1")
	pageInt, _ := strconv.Atoi(page)
	pageSize := c.Query("pageSize", "10")
	pageSizeInt, _ := strconv.Atoi(pageSize)
	projectName := c.Query("projectName", "")
	sortDirection := c.Query("sortDirection", "asc")
	categoryName := c.Query("categoryName", "")
	sortBy := c.Query("sortBy", "")

	fmt.Println(projectName)
	fmt.Println(categoryName)

	projects, totalItem, errResult := handler.ProjectRepository.GetAllProject(c, pageInt, pageSizeInt, sortDirection, projectName, categoryName)

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
		"code":    fiber.StatusOK,
		"message": "Successfully get projects",
		"data": fiber.Map{
			"items": projects,
			"pagination": fiber.Map{
				"currentPage":  page,
				"itemsPerPage": pageSize,
				"totalItems":   totalItem,
				"totalPages":   totalPages,
			},
			"filters": fiber.Map{
				"categoryName":  categoryName,
				"projectName":   projectName,
				"sortBy":        sortBy,
				"sortDirection": sortDirection,
			},
		},
	})
}

// CreateProjectItem godoc
// @Summary Create a new project item
// @Description Create a new item within an existing project
// @Tags Projects
// @Accept  json
// @Produce  json
// @Param projectItem body model.CreateProjectItem true "Create Project Item Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /project-item [post]
// @Security Bearer
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
		ID:         uuid.New(),
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

// UpdateProjectItem godoc
// @Summary Update an existing project item
// @Description Update the details of an existing project item by ID
// @Tags Projects
// @Accept  json
// @Produce  json
// @Param id path string true "Project Item ID"
// @Param projectItem body model.UpdateProjectItem true "Update Project Item Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /project-item/{id} [put]
// @Security Bearer
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
		"message": "Successfully updated project item",
	})
}

// @Summary Delete a project item
// @Description Delete a specific project item by its ID
// @Tags Projects
// @Param id path string true "Project Item ID"
// @Success 200 {object} map[string]interface{} "Successfully deleted the project item"
// @Failure 400 {object} map[string]interface{} "Invalid ID or error parsing ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /project-item/{id} [delete]
// @Security Bearer
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

// @Summary Get all project items by project ID
// @Description Retrieve all items associated with a specific project, with pagination, sorting, and filtering options.
// @Tags Projects
// @Param project_id path string true "Project ID"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Number of items per page" default(10)
// @Param projectItemName query string false "Filter by project item name"
// @Param sortDirection query string false "Sort order, either 'asc' or 'desc'" default(asc)
// @Success 200 {object} map[string]interface{} "Successfully retrieved project items"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /project-item/{project_id} [get]
// @Security Bearer
func (handler *ProjectHandlerImpl) GetAllProjectItemByProjectId(c *fiber.Ctx) error {

	page := c.Query("page", "1")
	pageInt, _ := strconv.Atoi(page)
	pageSize := c.Query("pageSize", "10")
	pageSizeInt, _ := strconv.Atoi(pageSize)
	projectItemName := c.Query("projectItemName", "")
	sortDirection := c.Query("sortDirection", "asc")
	sortBy := c.Query("sortBy", "")
	idString := c.Params("project_id", "")
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

	projectItems, totalItem, errResult := handler.ProjectRepository.GetAllProjectItemByProjectId(c, pageInt, pageSizeInt, sortDirection, projectItemName, uuidID)

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
		"code":    fiber.StatusOK,
		"message": "Successfully get category",
		"data": fiber.Map{
			"items": projectItems,
			"pagination": fiber.Map{
				"currentPage":  page,
				"itemsPerPage": pageSize,
				"totalItems":   totalItem,
				"totalPages":   totalPages,
			},
			"filters": fiber.Map{
				"projectItemName": projectItemName,
				"sortBy":          sortBy,
				"sortDirection":   sortDirection,
			},
		},
	})

}
