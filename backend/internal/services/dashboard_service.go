package services

import (
	"backend/internal/adapters/interfaces"
	"backend/internal/domain"
	"context"
	"time"
)

type DashboardService struct {
	dashboardRepo interfaces.DashboardRepository
}

func NewDashboardService(dashboardRepo interfaces.DashboardRepository) *DashboardService {
	return &DashboardService{
		dashboardRepo: dashboardRepo,
	}
}

func (s *DashboardService) GetDashboard(
	ctx context.Context,
	userID int64,
	from time.Time,
	to time.Time,
) (*domain.Dashboard, error) {
	return s.dashboardRepo.GetDashboard(ctx, userID, from, to)
}
