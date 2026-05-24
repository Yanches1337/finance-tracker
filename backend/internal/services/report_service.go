package services

import (
	"backend/internal/adapters/interfaces"
	"backend/internal/domain"
	"backend/internal/exporters"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type ReportService struct {
	reportRepo      interfaces.ReportRepository
	transactionRepo interfaces.TransactionRepository
}

func NewReportService(reportRepo interfaces.ReportRepository, transactionRepo interfaces.TransactionRepository) *ReportService {
	return &ReportService{
		reportRepo:      reportRepo,
		transactionRepo: transactionRepo,
	}
}

func (s *ReportService) Generate(ctx context.Context, userID int64, req domain.GenerateReportRequest) (*domain.Report, error) {
	if req.Format != "csv" {
		return nil, errors.New("currently only csv format is supported")
	}

	transactions, err := s.transactionRepo.GetByPeriod(
		ctx,
		userID,
		req.FromDate,
		req.ToDate,
	)

	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll("storage/reports", os.ModePerm); err != nil {
		return nil, err
	}

	fileName := fmt.Sprintf(
		"report_%d_%d.csv",
		userID,
		time.Now().Unix(),
	)

	filePath := filepath.Join("storage/reports", fileName)

	if err := exporters.ExportTransactionsToCSV(filePath, transactions); err != nil {
		return nil, err
	}

	report := &domain.Report{
		UserID:    userID,
		Format:    req.Format,
		FromDate:  req.FromDate,
		ToDate:    req.ToDate,
		FilePath:  filePath,
		FileName:  fileName,
		CreatedAt: time.Now(),
	}

	if err := s.reportRepo.Create(ctx, report); err != nil {
		return nil, err
	}

	return report, nil
}

func (s *ReportService) GetByID(ctx context.Context, id int64, userID int64) (*domain.Report, error) {
	return s.reportRepo.GetByID(ctx, id, userID)
}

func (s *ReportService) GetAllByUser(ctx context.Context, userID int64) ([]domain.Report, error) {
	return s.reportRepo.GetAllByUser(ctx, userID)
}

func (s *ReportService) Delete(ctx context.Context, id int64, userID int64) error {
	return s.reportRepo.Delete(ctx, id, userID)
}
