package interfaces

import (
	"backend/internal/domain"
	"context"
	"time"
)

type DashboardRepository interface {
	GetDashboard(
		ctx context.Context,
		userID int64,
		from time.Time,
		to time.Time,
	) (*domain.Dashboard, error)
}
