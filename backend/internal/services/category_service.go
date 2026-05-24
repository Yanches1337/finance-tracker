package services

import (
	"backend/internal/adapters/interfaces"
	"backend/internal/domain"
	"context"
	"errors"
)

type CategoryService struct {
	categoryRepo interfaces.CategoryRepository
}

func NewCategoryService(categoryRepo interfaces.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *CategoryService) Create(ctx context.Context, c *domain.Category) error {
	if c.Name == "" {
		return errors.New("category name is required")
	}

	if c.Type != "income" && c.Type != "expense" {
		return errors.New("invalid category type")
	}

	err := s.categoryRepo.Create(ctx, c)
	if err != nil {
		return err
	}

	return nil
}

func (s *CategoryService) GetByID(ctx context.Context, id int64) (*domain.Category, error) {
	c, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (s *CategoryService) GetAllByUser(ctx context.Context, userID int64) ([]domain.Category, error) {
	return s.categoryRepo.GetAllByUser(ctx, userID)
}

func (s *CategoryService) Update(ctx context.Context, c *domain.Category) error {
	if c.ID == 0 {
		return errors.New("category id is required")
	}

	return s.categoryRepo.Update(ctx, c)
}

func (s *CategoryService) Delete(ctx context.Context, id int64, userID int64) error {
	return s.categoryRepo.Delete(ctx, id, userID)
}
