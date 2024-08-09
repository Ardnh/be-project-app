package project

import (
	"project-app/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	CreateProject(ctx *fiber.Ctx, req *model.Project) error
	UpdateProject(ctx *fiber.Ctx, req *model.Project) error
	DeleteProject(ctx *fiber.Ctx, req uuid.UUID) error
	GetByIdProject(ctx *fiber.Ctx, req uuid.UUID) (*model.Project, error)
	GetAllProject(ctx *fiber.Ctx) (*[]model.Project, int, error)

	CreateProjectItem(ctx *fiber.Ctx, req *model.ProjectItem) error
	UpdateProjectItem(ctx *fiber.Ctx, req uuid.UUID) error
	DeleteProjectItem(ctx *fiber.Ctx, req uuid.UUID) error
	GetAllProjectItemByProjectId(ctx *fiber.Ctx) (*[]model.ProjectItem, int, error)
}

type ProjectRepositoryImpl struct {
	Db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &ProjectRepositoryImpl{
		Db: db,
	}
}

func (repository *ProjectRepositoryImpl) CreateProject(ctx *fiber.Ctx, req *model.Project) error {

}

func (repository *ProjectRepositoryImpl) UpdateProject(ctx *fiber.Ctx, req *model.Project) error
func (repository *ProjectRepositoryImpl) DeleteProject(ctx *fiber.Ctx, req uuid.UUID) error
func (repository *ProjectRepositoryImpl) GetByIdProject(ctx *fiber.Ctx, req uuid.UUID) (*model.Project, error)
func (repository *ProjectRepositoryImpl) GetAllProject(ctx *fiber.Ctx) (*[]model.Project, int, error)
func (repository *ProjectRepositoryImpl) CreateProjectItem(ctx *fiber.Ctx, req *model.ProjectItem) error
func (repository *ProjectRepositoryImpl) UpdateProjectItem(ctx *fiber.Ctx, req uuid.UUID) error
func (repository *ProjectRepositoryImpl) DeleteProjectItem(ctx *fiber.Ctx, req uuid.UUID) error
func (repository *ProjectRepositoryImpl) GetAllProjectItemByProjectId(ctx *fiber.Ctx) (*[]model.ProjectItem, int, error)
