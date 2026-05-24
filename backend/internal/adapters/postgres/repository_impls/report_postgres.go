package repository_impls

import (
	"backend/internal/adapters/postgres"
	"backend/internal/domain"
	"context"
	"github.com/jackc/pgx/v5"
)

type ReportPostgresRepository struct{}

func NewReportPostgresRepository() *ReportPostgresRepository {
	return &ReportPostgresRepository{}
}

func (r *ReportPostgresRepository) Create(ctx context.Context, report *domain.Report) error {
	query := `
		INSERT INTO reports (
			user_id,
			format,
			file_path,
		    file_name,
			from_date,
			to_date
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	return postgres.DB.QueryRow(
		ctx,
		query,
		report.UserID,
		report.Format,
		report.FilePath,
		report.FileName,
		report.FromDate,
		report.ToDate,
	).Scan(
		&report.ID,
		&report.CreatedAt,
	)
}

func (r *ReportPostgresRepository) GetByID(
	ctx context.Context,
	id int64,
	userID int64,
) (*domain.Report, error) {

	query := `
		SELECT
			id,
			user_id,
			format,
			file_path,
			file_name,
			from_date,
			to_date,
			created_at
		FROM reports
		WHERE id = $1
		  AND user_id = $2
	`

	report := &domain.Report{}

	err := postgres.DB.QueryRow(
		ctx,
		query,
		id,
		userID,
	).Scan(
		&report.ID,
		&report.UserID,
		&report.Format,
		&report.FilePath,
		&report.FileName,
		&report.FromDate,
		&report.ToDate,
		&report.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, pgx.ErrNoRows
	}

	return report, err
}

func (r *ReportPostgresRepository) GetAllByUser(
	ctx context.Context,
	userID int64,
) ([]domain.Report, error) {

	query := `
		SELECT
			id,
			user_id,
			format,
			file_path,
			file_name,
			from_date,
			to_date,
			created_at
		FROM reports
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := postgres.DB.Query(
		ctx,
		query,
		userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var reports []domain.Report

	for rows.Next() {
		var report domain.Report

		if err := rows.Scan(
			&report.ID,
			&report.UserID,
			&report.Format,
			&report.FilePath,
			&report.FileName,
			&report.FromDate,
			&report.ToDate,
			&report.CreatedAt,
		); err != nil {
			return nil, err
		}

		reports = append(reports, report)
	}

	return reports, nil
}

func (r *ReportPostgresRepository) Delete(
	ctx context.Context,
	id int64,
	userID int64,
) error {

	query := `
		DELETE FROM reports
		WHERE id = $1
		  AND user_id = $2
	`

	_, err := postgres.DB.Exec(
		ctx,
		query,
		id,
		userID,
	)

	return err
}
