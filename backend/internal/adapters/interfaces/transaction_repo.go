package interfaces

import (
	"backend/internal/domain"
	"context"
	"time"
)

type TransactionRepository interface {
	Create(ctx context.Context, t *domain.Transaction) error
	GetByID(ctx context.Context, id int64, userID int64) (*domain.Transaction, error)
	GetByPeriod(ctx context.Context, userID int64, from time.Time, to time.Time) ([]domain.Transaction, error)
	GetAllByUser(ctx context.Context, userID int64) ([]domain.Transaction, error)
	Update(ctx context.Context, t *domain.Transaction) error
	Delete(ctx context.Context, id int64, userID int64) error
}
