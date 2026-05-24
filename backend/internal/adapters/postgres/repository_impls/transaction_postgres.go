package repository_impls

import (
	"backend/internal/adapters/postgres"
	"backend/internal/domain"
	"context"
	"github.com/jackc/pgx/v5"
	"time"
)

type TransactionPostgresRepository struct{}

func NewTransactionPostgresRepository() *TransactionPostgresRepository {
	return &TransactionPostgresRepository{}
}

func (r *TransactionPostgresRepository) Create(ctx context.Context, t *domain.Transaction) error {
	query := `
		INSERT INTO transactions (
			user_id, type, amount, category_id, date, description
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	return postgres.DB.QueryRow(
		ctx,
		query,
		t.UserID,
		t.Type,
		t.Amount,
		t.CategoryID,
		t.Date,
		t.Description,
	).Scan(&t.ID, &t.CreatedAt)
}

func (r *TransactionPostgresRepository) GetByID(ctx context.Context, id int64, userID int64) (*domain.Transaction, error) {
	query := `
		SELECT id, user_id, type, amount, category_id, date, description, created_at
		FROM transactions
		WHERE id = $1 AND user_id = $2
	`

	t := &domain.Transaction{}

	err := postgres.DB.QueryRow(ctx, query, id, userID).Scan(
		&t.ID,
		&t.UserID,
		&t.Type,
		&t.Amount,
		&t.CategoryID,
		&t.Date,
		&t.Description,
		&t.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, pgx.ErrNoRows
	}

	return t, err

}

func (r *TransactionPostgresRepository) GetByPeriod(ctx context.Context, userID int64, from time.Time, to time.Time) ([]domain.Transaction, error) {

	query := `
		SELECT
			id,
			user_id,
			type,
			amount,
			category_id,
			date,
			description,
			created_at
		FROM transactions
		WHERE user_id = $1
		  AND date BETWEEN $2 AND $3
		ORDER BY date DESC
	`

	rows, err := postgres.DB.Query(
		ctx,
		query,
		userID,
		from,
		to,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var transactions []domain.Transaction

	for rows.Next() {
		var t domain.Transaction

		if err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.Type,
			&t.Amount,
			&t.CategoryID,
			&t.Date,
			&t.Description,
			&t.CreatedAt,
		); err != nil {
			return nil, err
		}

		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (r *TransactionPostgresRepository) GetAllByUser(ctx context.Context, userID int64) ([]domain.Transaction, error) {
	query := `
		SELECT id, user_id, type, amount, category_id, date, description, created_at
		FROM transactions
		WHERE user_id = $1
		ORDER BY date DESC
	`

	rows, err := postgres.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Transaction

	for rows.Next() {
		var t domain.Transaction

		if err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.Type,
			&t.Amount,
			&t.CategoryID,
			&t.Date,
			&t.Description,
			&t.CreatedAt,
		); err != nil {
			return nil, err
		}

		result = append(result, t)
	}

	return result, rows.Err()
}

func (r *TransactionPostgresRepository) Update(ctx context.Context, t *domain.Transaction) error {
	query := `
		UPDATE transactions
		SET type = $1,
			amount = $2,
			category_id = $3,
			date = $4,
			description = $5
		WHERE id = $6 AND user_id = $7
	`

	_, err := postgres.DB.Exec(
		ctx,
		query,
		t.Type,
		t.Amount,
		t.CategoryID,
		t.Date,
		t.Description,
		t.ID,
		t.UserID,
	)

	return err
}

func (r *TransactionPostgresRepository) Delete(ctx context.Context, id int64, userID int64) error {
	query := `
		DELETE FROM transactions
		WHERE id = $1 AND user_id = $2
	`

	_, err := postgres.DB.Exec(ctx, query, id, userID)
	return err
}
