package proeject

import (
	"project-app/model"
	projectRepository "project-app/repository/project"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
	// Parsing body request
	var request model.CreateProjectRequest
	if err := c.BodyParser(&request); err != nil {
		logrus.WithError(err).Error("Failed to parse CreateProject request body")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Invalid request body",
		})
	}

	// Validasi request
	if errValidate := handler.Validator.Struct(request); errValidate != nil {
		logrus.WithError(errValidate).Warn("Validation error in CreateProject request")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Validation failed",
		})
	}

	// Membuat project request
	createRequest := model.Projects{
		ID:          uuid.New(),
		Name:        request.Name,
		CategoryID:  request.CategoryID,
		Description: request.Description,
		Budget:      request.Budget,
		UserID:      request.UserID,
	}

	// Simpan project ke repository
	if err := handler.ProjectRepository.CreateProject(c, &createRequest); err != nil {
		logrus.WithError(err).Error("Failed to create project")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Project creation failed",
		})
	}

	// Log keberhasilan pembuatan proyek
	logrus.WithFields(logrus.Fields{
		"ProjectID":  createRequest.ID,
		"Name":       createRequest.Name,
		"CategoryID": createRequest.CategoryID,
		"UserID":     createRequest.UserID,
	}).Info("Project created successfully")

	// Kembalikan response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Project created successfully",
		"data":    createRequest,
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
	// Baca ID dari parameter URL
	idString := c.Params("id", "")
	if idString == "" {
		logrus.Warn("Invalid project ID: ID is empty")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid ID",
		})
	}

	// Parsing UUID
	uuidID, err := uuid.Parse(idString)
	if err != nil {
		logrus.WithError(err).Warn("Error parsing project ID")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Error parsing ID",
		})
	}

	// Parsing body request
	var request model.UpdateProjectRequest
	if err := c.BodyParser(&request); err != nil {
		logrus.WithError(err).Error("Failed to parse UpdateProject request body")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	// Log request parsed successfully
	logrus.WithFields(logrus.Fields{
		"Name":        request.Name,
		"CategoryID":  request.CategoryID,
		"Description": request.Description,
		"Budget":      request.Budget,
		"UserID":      request.UserID,
	}).Info("Parsed UpdateProject request successfully")

	// Validasi request
	errValidate := handler.Validator.Struct(&request)
	if errValidate != nil {
		logrus.WithError(errValidate).Warn("Validation error in UpdateProject request")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": errValidate.Error(),
		})
	}

	// Log proses update project dimulai
	logrus.WithField("ProjectID", uuidID).Info("Updating project...")

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
		logrus.WithError(errResult).Error("Failed to update project")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": errResult.Error(),
		})
	}

	// Log success
	logrus.WithFields(logrus.Fields{
		"ProjectID":  updateRequest.ID,
		"Name":       updateRequest.Name,
		"CategoryID": updateRequest.CategoryID,
		"UserID":     updateRequest.UserID,
	}).Info("Project updated successfully")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Successfully updated project",
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
	// Ambil ID dari parameter URL
	idString := c.Params("id", "")
	if idString == "" {
		logrus.Warn("DeleteProject request failed: invalid ID parameter")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid ID",
		})
	}

	// Parsing UUID
	uuidID, err := uuid.Parse(idString)
	if err != nil {
		logrus.WithError(err).Warn("Error parsing project ID")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Error parsing ID",
		})
	}

	// Log ID proyek yang akan dihapus
	logrus.WithField("ProjectID", uuidID).Info("Deleting project")

	// Hapus project melalui repository
	errResult := handler.ProjectRepository.DeleteProject(c, uuidID)
	if errResult != nil {
		logrus.WithError(errResult).Error("Failed to delete project")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Failed to delete project",
		})
	}

	// Log jika project berhasil dihapus
	logrus.WithField("ProjectID", uuidID).Info("Project deleted successfully")

	// Mengembalikan response berhasil
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Project deleted successfully",
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
	// Ambil query parameters dan parsing jika diperlukan
	page := c.Query("page", "1")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		logrus.WithError(err).Warn("Invalid page parameter")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid page parameter",
		})
	}

	pageSize := c.Query("pageSize", "10")
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		logrus.WithError(err).Warn("Invalid pageSize parameter")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid pageSize parameter",
		})
	}

	projectName := c.Query("projectName", "")
	sortDirection := c.Query("sortDirection", "asc")
	categoryName := c.Query("categoryName", "")
	sortBy := c.Query("sortBy", "")

	// Log detail query parameters
	logrus.WithFields(logrus.Fields{
		"page":          pageInt,
		"pageSize":      pageSizeInt,
		"projectName":   projectName,
		"categoryName":  categoryName,
		"sortDirection": sortDirection,
		"sortBy":        sortBy,
	}).Info("Fetching projects with filters")

	// Panggil repository untuk mendapatkan data project
	projects, totalItem, errResult := handler.ProjectRepository.GetAllProject(c, pageInt, pageSizeInt, sortDirection, projectName, categoryName)
	if errResult != nil {
		logrus.WithError(errResult).Error("Failed to fetch projects")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Failed to fetch projects",
		})
	}

	// Hitung total halaman
	totalPages := totalItem / int64(pageSizeInt)
	if totalItem%int64(pageSizeInt) > 0 {
		totalPages++
	}

	// Log sukses mendapatkan proyek
	logrus.WithFields(logrus.Fields{
		"totalItems": totalItem,
		"totalPages": totalPages,
	}).Info("Successfully fetched projects")

	// Kirimkan respons
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Successfully fetched projects",
		"data": fiber.Map{
			"items": projects,
			"pagination": fiber.Map{
				"currentPage":  pageInt,
				"itemsPerPage": pageSizeInt,
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
	// Log permintaan CreateProjectItem diterima
	logrus.Info("CreateProjectItem request received")

	// Parsing body request
	var request model.CreateProjectItem
	if err := c.BodyParser(&request); err != nil {
		logrus.WithError(err).Error("Failed to parse CreateProjectItem request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid request body",
		})
	}

	// Validasi request
	errValidate := handler.Validator.Struct(request)
	if errValidate != nil {
		logrus.WithError(errValidate).Warn("Validation error in CreateProjectItem request")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Validation failed",
		})
	}

	// Log proses pembuatan project item dimulai
	logrus.WithFields(logrus.Fields{
		"ProjectID":  request.ProjectID,
		"Name":       request.Name,
		"BudgetItem": request.BudgetItem,
		"Status":     request.Status,
	}).Info("Creating project item")

	// Membuat ProjectItem
	createRequest := model.ProjectItem{
		ID:         uuid.New(),
		ProjectID:  request.ProjectID,
		Name:       request.Name,
		BudgetItem: request.BudgetItem,
		Status:     request.Status,
	}

	// Simpan ProjectItem ke repository
	errResult := handler.ProjectRepository.CreateProjectItem(c, &createRequest)
	if errResult != nil {
		logrus.WithError(errResult).Error("Failed to create project item")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Failed to create project item",
		})
	}

	// Log keberhasilan pembuatan project item
	logrus.WithField("ProjectItemID", createRequest.ID).Info("Project item created successfully")

	// Kembalikan response berhasil
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Successfully created project item",
		"data":    createRequest,
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
	// Log permintaan UpdateProjectItem diterima
	logrus.Info("UpdateProjectItem request received")

	// Baca ID dari parameter URL
	idString := c.Params("id", "")
	if idString == "" {
		logrus.Warn("Invalid project item ID: ID is empty")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid ID",
		})
	}

	// Parsing UUID
	uuidID, err := uuid.Parse(idString)
	if err != nil {
		logrus.WithError(err).Warn("Error parsing project item ID")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Error parsing ID",
		})
	}

	// Parsing body request
	var request model.UpdateProjectItem
	if err := c.BodyParser(&request); err != nil {
		logrus.WithError(err).Error("Failed to parse UpdateProjectItem request body")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Invalid request body",
		})
	}

	// Validasi request
	errValidate := handler.Validator.Struct(&request)
	if errValidate != nil {
		logrus.WithError(errValidate).Warn("Validation error in UpdateProjectItem request")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Validation failed",
		})
	}

	// Log proses update dimulai
	logrus.WithFields(logrus.Fields{
		"ProjectItemID": uuidID,
		"ProjectID":     request.ProjectID,
		"Name":          request.Name,
		"BudgetItem":    request.BudgetItem,
		"Status":        request.Status,
	}).Info("Updating project item")

	// Update project item
	updateRequest := model.ProjectItem{
		ID:         uuidID,
		ProjectID:  request.ProjectID,
		Name:       request.Name,
		BudgetItem: request.BudgetItem,
		Status:     request.Status,
	}

	// Simpan perubahan
	errResult := handler.ProjectRepository.UpdateProjectItem(c, &updateRequest)
	if errResult != nil {
		logrus.WithError(errResult).Error("Failed to update project item")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Failed to update project item",
		})
	}

	// Log keberhasilan update
	logrus.WithField("ProjectItemID", updateRequest.ID).Info("Project item updated successfully")

	// Kembalikan respons berhasil
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Successfully updated project item",
		"data":    updateRequest,
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
	// Ambil ID dari parameter URL
	idString := c.Params("id", "")
	if idString == "" {
		logrus.Warn("DeleteProjectItem request failed: invalid ID parameter")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid ID",
		})
	}

	// Parsing UUID
	uuidID, err := uuid.Parse(idString)
	if err != nil {
		logrus.WithError(err).Warn("Error parsing project item ID")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Error parsing ID",
		})
	}

	// Log penghapusan project item dimulai
	logrus.WithField("ProjectItemID", uuidID).Info("Deleting project item")

	// Hapus ProjectItem dari repository
	errResult := handler.ProjectRepository.DeleteProjectItem(c, uuidID)
	if errResult != nil {
		logrus.WithError(errResult).Error("Failed to delete project item")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Failed to delete project item",
		})
	}

	// Log sukses menghapus project item
	logrus.WithField("ProjectItemID", uuidID).Info("Project item deleted successfully")

	// Mengembalikan respons berhasil
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Successfully deleted project item",
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
	// Ambil query parameter dan parsing
	page := c.Query("page", "1")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		logrus.WithError(err).Warn("Invalid page parameter")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid page parameter",
		})
	}

	pageSize := c.Query("pageSize", "10")
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		logrus.WithError(err).Warn("Invalid pageSize parameter")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid pageSize parameter",
		})
	}

	projectItemName := c.Query("projectItemName", "")
	sortDirection := c.Query("sortDirection", "asc")
	sortBy := c.Query("sortBy", "")

	// Ambil dan validasi parameter project_id
	idString := c.Params("project_id", "")
	if idString == "" {
		logrus.Warn("Invalid project ID: ID is empty")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid project ID",
		})
	}

	uuidID, err := uuid.Parse(idString)
	if err != nil {
		logrus.WithError(err).Warn("Error parsing project ID")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Error parsing project ID",
		})
	}

	// Log detail pencarian
	logrus.WithFields(logrus.Fields{
		"page":            pageInt,
		"pageSize":        pageSizeInt,
		"projectItemName": projectItemName,
		"sortDirection":   sortDirection,
		"sortBy":          sortBy,
		"ProjectID":       uuidID,
	}).Info("Fetching project items by project ID")

	// Panggil repository untuk mengambil project items
	projectItems, totalItem, errResult := handler.ProjectRepository.GetAllProjectItemByProjectId(c, pageInt, pageSizeInt, sortDirection, projectItemName, uuidID)
	if errResult != nil {
		logrus.WithError(errResult).Error("Failed to fetch project items")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Failed to fetch project items",
		})
	}

	// Hitung total halaman
	totalPages := int(totalItem) / pageSizeInt
	if int(totalItem)%pageSizeInt > 0 {
		totalPages++
	}

	// Log keberhasilan pengambilan project items
	logrus.WithFields(logrus.Fields{
		"totalItems": totalItem,
		"totalPages": totalPages,
	}).Info("Successfully fetched project items")

	// Kembalikan respons berhasil
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Successfully fetched project items",
		"data": fiber.Map{
			"items": projectItems,
			"pagination": fiber.Map{
				"currentPage":  pageInt,
				"itemsPerPage": pageSizeInt,
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
