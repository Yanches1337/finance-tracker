package repository_impls

import (
	"backend/internal/adapters/postgres"
	"backend/internal/domain"
	"context"
	"github.com/jackc/pgx/v5"
)

type UserPostgresRepository struct{}

func NewUserPostgresRepository() *UserPostgresRepository {
	return &UserPostgresRepository{}
}

func (r *UserPostgresRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (email, password_hash, name)
		VALUES ($1, $2, $3) 
		RETURNING id, created_at, updated_at`

	err := postgres.DB.QueryRow(ctx, query, user.Email, user.PasswordHash, user.Name).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	return err
}

func (r *UserPostgresRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, name, created_at, updated_at 
		FROM users 
		WHERE email = $1`

	user := &domain.User{}
	err := postgres.DB.QueryRow(ctx, query, email).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.CreatedAt, &user.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, pgx.ErrNoRows
	}
	return user, err
}

func (r *UserPostgresRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := postgres.DB.QueryRow(ctx, query, email).Scan(&exists)
	return exists, err
}

func (r *UserPostgresRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, name, created_at, updated_at 
		FROM users 
		WHERE id = $1`

	user := &domain.User{}
	err := postgres.DB.QueryRow(ctx, query, id).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.CreatedAt, &user.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, pgx.ErrNoRows
	}
	return user, err
}
