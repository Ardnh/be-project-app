package repositories

import (
	"context"

	"github.com/Ardnh/be-project-app/internal/domain/entities"
)

type ProjectExpensesItemRepository interface {
	Create(ctx context.Context, expenses *entities.ProjectExpenseItem) error
	Update(ctx context.Context, expenses *entities.ProjectExpenseItem) error
	Delete(ctx context.Context, id string) error
}
