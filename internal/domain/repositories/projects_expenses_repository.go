package repositories

import (
	"context"

	"github.com/Ardnh/be-project-app/internal/domain/entities"
)

type ProjectsExpensesRepository interface {
	Create(ctx context.Context, expenses *entities.ProjectExpenses) error
	Update(ctx context.Context, expenses *entities.ProjectExpenses) error
	Delete(ctx context.Context, id string) error
}
