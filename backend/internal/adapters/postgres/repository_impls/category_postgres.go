package repository_impls

import (
	"backend/internal/adapters/postgres"
	"backend/internal/domain"
	"context"
	"github.com/jackc/pgx/v5"
)

type CategoryPostgresRepository struct{}

func NewCategoryPostgresRepository() *CategoryPostgresRepository {
	return &CategoryPostgresRepository{}
}

func (r *CategoryPostgresRepository) Create(ctx context.Context, category *domain.Category) error {
	query := `
		INSERT INTO categories (user_id, name, type, color, icon)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	return postgres.DB.QueryRow(
		ctx,
		query,
		category.UserID,
		category.Name,
		category.Type,
		category.Color,
		category.Icon,
	).Scan(&category.ID)
}

func (r *CategoryPostgresRepository) GetByID(ctx context.Context, id int64) (*domain.Category, error) {
	query := `
		SELECT id, user_id, name, type, color, icon
		FROM categories
		WHERE id = $1
	`

	c := &domain.Category{}

	err := postgres.DB.QueryRow(ctx, query, id).Scan(
		&c.ID,
		&c.UserID,
		&c.Name,
		&c.Type,
		&c.Color,
		&c.Icon,
	)

	if err == pgx.ErrNoRows {
		return nil, pgx.ErrNoRows
	}

	return c, err
}

func (r *CategoryPostgresRepository) GetAllByUser(ctx context.Context, userID int64) ([]domain.Category, error) {
	query := `
		SELECT id, user_id, name, type, color, icon
		FROM categories
		WHERE user_id = $1
		ORDER BY id DESC
	`

	rows, err := postgres.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Category

	for rows.Next() {
		var c domain.Category

		if err := rows.Scan(
			&c.ID,
			&c.UserID,
			&c.Name,
			&c.Type,
			&c.Color,
			&c.Icon,
		); err != nil {
			return nil, err
		}

		result = append(result, c)
	}

	return result, rows.Err()
}

func (r *CategoryPostgresRepository) GetByUserAndType(
	ctx context.Context,
	userID int64,
	type_ domain.TransactionType,
) ([]domain.Category, error) {

	query := `
		SELECT id, user_id, name, type, color, icon
		FROM categories
		WHERE user_id = $1 AND type = $2
		ORDER BY name
	`

	rows, err := postgres.DB.Query(ctx, query, userID, type_)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Category

	for rows.Next() {
		var c domain.Category

		if err := rows.Scan(
			&c.ID,
			&c.UserID,
			&c.Name,
			&c.Type,
			&c.Color,
			&c.Icon,
		); err != nil {
			return nil, err
		}

		result = append(result, c)
	}

	return result, rows.Err()
}

func (r *CategoryPostgresRepository) Update(ctx context.Context, category *domain.Category) error {
	query := `
		UPDATE categories
		SET name = $1,
		    type = $2,
		    color = $3,
		    icon = $4
		WHERE id = $5 AND user_id = $6
	`

	cmd, err := postgres.DB.Exec(
		ctx,
		query,
		category.Name,
		category.Type,
		category.Color,
		category.Icon,
		category.ID,
		category.UserID,
	)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *CategoryPostgresRepository) Delete(ctx context.Context, id int64, userID int64) error {
	query := `
		DELETE FROM categories
		WHERE id = $1 AND user_id = $2
	`

	cmd, err := postgres.DB.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
