package proeject

import (
	projectRepository "project-app/repository/project"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ProjectHandler interface {
	CreateProject(c *fiber.Ctx) error
	UpdateProject(c *fiber.Ctx) error
	DeleteProject(c *fiber.Ctx) error
	GetByIdProject(c *fiber.Ctx) error
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

func (handler *ProjectHandlerImpl) CreateProject(c *fiber.Ctx) error
func (handler *ProjectHandlerImpl) UpdateProject(c *fiber.Ctx) error
func (handler *ProjectHandlerImpl) DeleteProject(c *fiber.Ctx) error
func (handler *ProjectHandlerImpl) GetByIdProject(c *fiber.Ctx) error
func (handler *ProjectHandlerImpl) GetAllProject(c *fiber.Ctx) error
func (handler *ProjectHandlerImpl) CreateProjectItem(c *fiber.Ctx) error
func (handler *ProjectHandlerImpl) UpdateProjectItem(c *fiber.Ctx) error
func (handler *ProjectHandlerImpl) DeleteProjectItem(c *fiber.Ctx) error
func (handler *ProjectHandlerImpl) GetAllProjectItemByProjectId(c *fiber.Ctx) error
