package router

import (
	"backend/internal/api/http/handlers"
	"backend/internal/api/http/middleware"
	"backend/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Dependencies struct {
	AuthService        *services.AuthService
	CategoryService    *services.CategoryService
	TransactionService *services.TransactionService
	GoalService        *services.GoalService
	DashboardService   *services.DashboardService
	ReportService      *services.ReportService
}

func SetupRouter(r *gin.Engine, deps Dependencies) {
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "finance-tracker",
		})
	})

	// TODO: Здесь будем добавить группы роутов (auth, transactions, goals и т.д.)
	api := r.Group("/api/v1")
	{
		// -------------------------
		// AUTH
		// -------------------------
		authHandler := handlers.NewAuthHandler(deps.AuthService)

		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)
		}

		// -------------------------
		// PROTECTED
		// -------------------------
		protected := api.Group("/protected")
		protected.Use(middleware.AuthMiddleware(deps.AuthService.Cfg.JWT.SecretKey))
		{
			protected.GET("/me", func(c *gin.Context) {
				userID := c.GetInt64("user_id")
				c.JSON(http.StatusOK, gin.H{"user_id": userID})
			})
			categoryHandler := handlers.NewCategoryHandler(deps.CategoryService)
			categories := protected.Group("/categories")
			{
				categories.POST("", categoryHandler.Create)
				categories.GET("", categoryHandler.GetAll)
				categories.GET("/:id", categoryHandler.GetByID)
				categories.PUT("", categoryHandler.Update)
				categories.DELETE("/:id", categoryHandler.Delete)
			}
			transactionHandler := handlers.NewTransactionHandler(deps.TransactionService)
			transactions := protected.Group("/transactions")
			{
				transactions.POST("", transactionHandler.Create)
				transactions.GET("", transactionHandler.GetAll)
				transactions.GET("/:id", transactionHandler.GetByID)
				transactions.PUT("", transactionHandler.Update)
				transactions.DELETE("/:id", transactionHandler.Delete)
			}

			goalHandler := handlers.NewGoalHandler(deps.GoalService)
			goals := protected.Group("/goals")
			{
				goals.POST("", goalHandler.Create)
				goals.GET("", goalHandler.GetAll)
				goals.GET("/:id", goalHandler.GetByID)
				goals.PUT("/:id", goalHandler.Update)
				goals.DELETE("/:id", goalHandler.Delete)
			}

			dashboardHandler := handlers.NewDashboardHandler(deps.DashboardService)
			dashboard := protected.Group("/dashboards")
			{
				dashboard.GET("", dashboardHandler.GetDashboard)
			}

			reportHandler := handlers.NewReportHandler(deps.ReportService)
			reports := protected.Group("/reports")
			{
				reports.POST("", reportHandler.Generate)
				reports.GET("", reportHandler.GetAll)
				reports.GET("/:id/download", reportHandler.Download)
			}
		}

	}
}
