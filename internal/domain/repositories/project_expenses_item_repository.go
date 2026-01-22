package repositories

import (
	"context"

	"github.com/Ardnh/be-project-app/internal/domain/entities"
)

type ProjectExpensesItemRepository interface {
	Create(ctx context.Context, expenses *entities.ProjectExpenseItems) (*entities.ProjectExpenseItems, error)
	Update(ctx context.Context, expenses *entities.ProjectExpenseItems) (*entities.ProjectExpenseItems, error)
	Delete(ctx context.Context, id string) error
}
