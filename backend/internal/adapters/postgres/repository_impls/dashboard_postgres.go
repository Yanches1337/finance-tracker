package repository_impls

import (
	"backend/internal/adapters/postgres"
	"backend/internal/domain"
	"context"
	"time"
)

type DashboardPostgresRepository struct{}

func NewDashboardPostgresRepository() *DashboardPostgresRepository {
	return &DashboardPostgresRepository{}
}

func (r *DashboardPostgresRepository) GetDashboard(
	ctx context.Context,
	userID int64,
	from time.Time,
	to time.Time,
) (*domain.Dashboard, error) {

	dashboard := &domain.Dashboard{}

	incomeQuery := `
		SELECT COALESCE(SUM(amount), 0)
		FROM transactions
		WHERE user_id = $1
			AND type = 'income'
			AND date BETWEEN $2 AND $3
	`
	err := postgres.DB.QueryRow(ctx, incomeQuery, userID, from, to).
		Scan(&dashboard.TotalIncome)
	if err != nil {
		return nil, err
	}

	expenseQuery := `
		SELECT COALESCE(SUM(amount), 0)
		FROM transactions
		WHERE user_id = $1
		  AND type = 'expense'
		  AND date BETWEEN $2 AND $3
	`

	err = postgres.DB.QueryRow(
		ctx,
		expenseQuery,
		userID,
		from,
		to,
	).Scan(&dashboard.TotalExpense)

	if err != nil {
		return nil, err
	}

	dashboard.Balance =
		dashboard.TotalIncome - dashboard.TotalExpense

	expenseCategoryQuery := `
		SELECT
			c.id,
			c.name,
			COALESCE(SUM(t.amount), 0)
		FROM transactions t
		JOIN categories c ON c.id = t.category_id
		WHERE t.user_id = $1
		  AND t.type = 'expense'
		  AND t.date BETWEEN $2 AND $3
		GROUP BY c.id, c.name
		ORDER BY SUM(t.amount) DESC
	`

	rows, err := postgres.DB.Query(
		ctx,
		expenseCategoryQuery,
		userID,
		from,
		to,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var stat domain.CategoryStats

		if err := rows.Scan(
			&stat.CategoryID,
			&stat.CategoryName,
			&stat.Total,
		); err != nil {
			return nil, err
		}

		dashboard.ExpensesByCategory =
			append(dashboard.ExpensesByCategory, stat)
	}

	incomeCategoryQuery := `
		SELECT
			c.id,
			c.name,
			COALESCE(SUM(t.amount), 0)
		FROM transactions t
		JOIN categories c ON c.id = t.category_id
		WHERE t.user_id = $1
		  AND t.type = 'income'
		  AND t.date BETWEEN $2 AND $3
		GROUP BY c.id, c.name
		ORDER BY SUM(t.amount) DESC
	`

	rows2, err := postgres.DB.Query(
		ctx,
		incomeCategoryQuery,
		userID,
		from,
		to,
	)

	if err != nil {
		return nil, err
	}

	defer rows2.Close()

	for rows2.Next() {
		var stat domain.CategoryStats

		if err := rows2.Scan(
			&stat.CategoryID,
			&stat.CategoryName,
			&stat.Total,
		); err != nil {
			return nil, err
		}

		dashboard.IncomeByCategory =
			append(dashboard.IncomeByCategory, stat)
	}

	return dashboard, nil
}
