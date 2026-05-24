package services

import (
	"backend/internal/adapters/interfaces"
	"backend/internal/domain"
	"context"
	"errors"
)

type GoalService struct {
	goalRepo interfaces.GoalRepository
}

func NewGoalService(goalRepo interfaces.GoalRepository) *GoalService {
	return &GoalService{goalRepo: goalRepo}
}

func (s *GoalService) Create(ctx context.Context, g *domain.Goal) error {
	if g.Name == "" {
		return errors.New("goal name is required")
	}
	if g.TargetAmount <= 0 {
		return errors.New("target amount must be > 0")
	}

	g.CurrentAmount = 0
	g.IsCompleted = false

	return s.goalRepo.Create(ctx, g)
}

func (s *GoalService) GetByID(ctx context.Context, id int64, userID int64) (*domain.Goal, error) {
	return s.goalRepo.GetByID(ctx, id, userID)
}

func (s *GoalService) GetAllByUser(ctx context.Context, userID int64) ([]domain.Goal, error) {
	return s.goalRepo.GetAllByUser(ctx, userID)
}

func (s *GoalService) Update(ctx context.Context, g *domain.Goal) error {
	if g.ID == 0 {
		return errors.New("goal id is required")
	}

	return s.goalRepo.Update(ctx, g)
}

func (s *GoalService) Delete(ctx context.Context, id int64, userID int64) error {
	return s.goalRepo.Delete(ctx, id, userID)
}
