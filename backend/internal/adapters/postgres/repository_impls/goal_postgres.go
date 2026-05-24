package repository_impls

import (
	"backend/internal/adapters/postgres"
	"backend/internal/domain"
	"context"
	"github.com/jackc/pgx/v5"
)

type GoalPostgresRepository struct{}

func NewGoalPostgresRepository() *GoalPostgresRepository {
	return &GoalPostgresRepository{}
}

func (r *GoalPostgresRepository) Create(ctx context.Context, g *domain.Goal) error {
	query := `
		INSERT INTO goals (
			user_id, name, target_amount, current_amount,
			target_date, description, is_completed
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id, created_at
	`

	return postgres.DB.QueryRow(
		ctx,
		query,
		g.UserID,
		g.Name,
		g.TargetAmount,
		g.CurrentAmount,
		g.TargetDate,
		g.Description,
		g.IsCompleted,
	).Scan(&g.ID, &g.CreatedAt)
}

func (r *GoalPostgresRepository) GetByID(ctx context.Context, id int64, userID int64) (*domain.Goal, error) {
	query := `
		SELECT id, user_id, name, target_amount, current_amount,
		       target_date, description, is_completed, created_at
		FROM goals
		WHERE id = $1 AND user_id = $2
	`

	g := &domain.Goal{}

	err := postgres.DB.QueryRow(ctx, query, id, userID).Scan(
		&g.ID,
		&g.UserID,
		&g.Name,
		&g.TargetAmount,
		&g.CurrentAmount,
		&g.TargetDate,
		&g.Description,
		&g.IsCompleted,
		&g.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, pgx.ErrNoRows
	}

	return g, err
}

func (r *GoalPostgresRepository) GetAllByUser(ctx context.Context, userID int64) ([]domain.Goal, error) {
	query := `
		SELECT id, user_id, name, target_amount, current_amount,
		       target_date, description, is_completed, created_at
		FROM goals
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := postgres.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goals []domain.Goal

	for rows.Next() {
		var g domain.Goal

		if err := rows.Scan(
			&g.ID,
			&g.UserID,
			&g.Name,
			&g.TargetAmount,
			&g.CurrentAmount,
			&g.TargetDate,
			&g.Description,
			&g.IsCompleted,
			&g.CreatedAt,
		); err != nil {
			return nil, err
		}

		goals = append(goals, g)
	}

	return goals, rows.Err()
}

func (r *GoalPostgresRepository) Update(ctx context.Context, g *domain.Goal) error {
	query := `
		UPDATE goals
		SET name = $1,
		    target_amount = $2,
		    current_amount = $3,
		    target_date = $4,
		    description = $5,
		    is_completed = $6
		WHERE id = $7 AND user_id = $8
	`

	cmdTag, err := postgres.DB.Exec(
		ctx,
		query,
		g.Name,
		g.TargetAmount,
		g.CurrentAmount,
		g.TargetDate,
		g.Description,
		g.IsCompleted,
		g.ID,
		g.UserID,
	)

	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *GoalPostgresRepository) Delete(ctx context.Context, id int64, userID int64) error {
	query := `DELETE FROM goals WHERE id = $1 AND user_id = $2`

	cmdTag, err := postgres.DB.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
