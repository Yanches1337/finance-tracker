package postgres

import (
	"backend/internal/utils"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

var DB *pgxpool.Pool

func InitDB(cfg *utils.DatabaseConfig) error {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)

	var err error

	for i := 0; i < 10; i++ {
		DB, err = pgxpool.New(context.Background(), dsn)
		if err != nil {
			utils.Log.Warnf("DB connect attempt %d failed: %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}

		err = DB.Ping(context.Background())
		if err == nil {
			utils.Log.Info("DB connected successfully")
			return nil
		}

		utils.Log.Warnf("DB ping attempt %d failed: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("failed to connect to database after retries: %w", err)
}
