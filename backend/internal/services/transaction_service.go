package services

import (
	"backend/internal/adapters/interfaces"
	"backend/internal/domain"
	"context"
	"errors"
)

type TransactionService struct {
	transactionRepo interfaces.TransactionRepository
}

func NewTransactionService(transactionRepo interfaces.TransactionRepository) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
	}
}

func (s *TransactionService) Create(ctx context.Context, t *domain.Transaction) error {
	if t.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	if t.Type != domain.Income && t.Type != domain.Expense {
		return errors.New("invalid transaction type")
	}

	// TODO: позже сюда добавим проверку category ownership (важно для безопасности)

	return s.transactionRepo.Create(ctx, t)
}

func (s *TransactionService) GetAllByUser(ctx context.Context, userID int64) ([]domain.Transaction, error) {
	return s.transactionRepo.GetAllByUser(ctx, userID)
}

func (s *TransactionService) GetByID(ctx context.Context, id int64, userID int64) (*domain.Transaction, error) {
	t, err := s.transactionRepo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	if t.UserID != userID {
		return nil, errors.New("access denied")
	}

	return t, nil
}

func (s *TransactionService) Update(ctx context.Context, t *domain.Transaction) error {
	if t.ID == 0 {
		return errors.New("transaction id is required")
	}

	if t.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	if t.Type != domain.Income && t.Type != domain.Expense {
		return errors.New("invalid transaction type")
	}

	return s.transactionRepo.Update(ctx, t)
}

func (s *TransactionService) Delete(ctx context.Context, id int64, userID int64) error {
	return s.transactionRepo.Delete(ctx, id, userID)
}
