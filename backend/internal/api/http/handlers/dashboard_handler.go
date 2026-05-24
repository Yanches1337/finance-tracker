package handlers

import (
	"backend/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type DashboardHandler struct {
	dashboardService *services.DashboardService
}

func NewDashboardHandler(dashboardService *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

// GetDashboard godoc
// @Summary Get dashboard analytics
// @Tags dashboard
// @Security BearerAuth
// @Produce json
// @Param from query string true "From date (YYYY-MM-DD)"
// @Param to query string true "To date (YYYY-MM-DD)"
// @Success 200 {object} domain.Dashboard
// @Failure 400 {object} map[string]string
// @Router /protected/dashboards [get]
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	userID := c.GetInt64("user_id")

	fromStr := c.Query("from")
	toStr := c.Query("to")

	if fromStr == "" || toStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "from and to query params are required",
		})
		return
	}

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid from date format",
		})
		return
	}

	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid to date format",
		})
		return
	}

	dashboard, err := h.dashboardService.GetDashboard(
		c.Request.Context(),
		userID,
		from,
		to,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dashboard)
}
