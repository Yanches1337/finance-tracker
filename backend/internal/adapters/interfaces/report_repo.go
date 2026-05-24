package interfaces

import (
	"backend/internal/domain"
	"context"
)

type ReportRepository interface {
	Create(ctx context.Context, report *domain.Report) error
	GetByID(ctx context.Context, id int64, userID int64) (*domain.Report, error)
	GetAllByUser(ctx context.Context, userID int64) ([]domain.Report, error)
	Delete(ctx context.Context, id int64, userID int64) error
}
