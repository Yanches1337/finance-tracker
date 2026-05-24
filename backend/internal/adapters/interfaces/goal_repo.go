package interfaces

import (
	"backend/internal/domain"
	"context"
)

type GoalRepository interface {
	Create(ctx context.Context, goal *domain.Goal) error
	GetByID(ctx context.Context, id int64, userID int64) (*domain.Goal, error)
	GetAllByUser(ctx context.Context, userID int64) ([]domain.Goal, error)
	Update(ctx context.Context, goal *domain.Goal) error
	Delete(ctx context.Context, id int64, userID int64) error
}
