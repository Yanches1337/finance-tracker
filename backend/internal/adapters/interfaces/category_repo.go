package interfaces

import (
	"backend/internal/domain"
	"context"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *domain.Category) error
	GetByID(ctx context.Context, id int64) (*domain.Category, error)
	GetAllByUser(ctx context.Context, userID int64) ([]domain.Category, error)
	GetByUserAndType(ctx context.Context, userID int64, type_ domain.TransactionType) ([]domain.Category, error)
	Update(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, id int64, userID int64) error
}
