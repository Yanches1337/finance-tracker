package main

import (
	_ "backend/docs"
	"backend/internal/api/http/middleware"
	"context"
	"fmt"
	_ "github.com/swaggo/files"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"backend/internal/adapters/postgres"
	"backend/internal/adapters/postgres/repository_impls"
	"backend/internal/adapters/redis"
	router "backend/internal/api/http"
	"backend/internal/services"
	"backend/internal/utils"
)

// @title Finance Tracker API
// @version 1.0
// @description API for personal finance tracking
// @host localhost:8000
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg, err := utils.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	utils.InitLogger(cfg.Logging.Level)
	utils.Log.Infof("Starting Finance Tracker...")

	// --- DB init ---
	if err := postgres.InitDB(&cfg.Database); err != nil {
		utils.Log.Fatalf("Failed to connect to database: %v", err)
	}
	utils.Log.Infof("Successfully connected to PostgreSQL")

	// --- Redis init ---
	if err := redis.InitRedis(&cfg.Redis); err != nil {
		utils.Log.Fatalf("Failed to connect to Redis: %v", err)
	}
	utils.Log.Infof("Successfully connected to Redis")

	// --- DI ---
	userRepo := repository_impls.NewUserPostgresRepository()
	categoryRepo := repository_impls.NewCategoryPostgresRepository()
	transactionRepo := repository_impls.NewTransactionPostgresRepository()
	goalRepo := repository_impls.NewGoalPostgresRepository()
	dashboardRepo := repository_impls.NewDashboardPostgresRepository()
	reportRepo := repository_impls.NewReportPostgresRepository()

	tokenService := services.NewTokenService(cfg)
	authService := services.NewAuthService(userRepo, tokenService, cfg)
	categoryService := services.NewCategoryService(categoryRepo)
	transactionService := services.NewTransactionService(transactionRepo)
	goalService := services.NewGoalService(goalRepo)
	dashboardService := services.NewDashboardService(dashboardRepo)
	reportService := services.NewReportService(reportRepo, transactionRepo)

	// --- HTTP server ---
	r := gin.New()

	// Подключаем CORS в самый верх цепочки, чтобы он обрабатывал абсолютно все роуты
	r.Use(middleware.CORSMiddleware())
	r.Use(gin.Recovery())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Use(gin.Recovery())

	router.SetupRouter(r, router.Dependencies{
		AuthService:        authService,
		CategoryService:    categoryService,
		TransactionService: transactionService,
		GoalService:        goalService,
		DashboardService:   dashboardService,
		ReportService:      reportService,
	})

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	// start server
	go func() {
		utils.Log.Infof("Server started on http://%s:%d", cfg.Server.Host, cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Log.Fatalf("Server error: %v", err)
		}
	}()

	// --- graceful shutdown ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	utils.Log.Warn("Shutdown signal received")
	if err := srv.Shutdown(context.Background()); err != nil {
		utils.Log.Fatalf("Server shutdown error: %v", err)
	}

	utils.Log.Infof("Server stopped gracefully")
}
