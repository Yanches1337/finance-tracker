package handlers

import (
	"backend/internal/domain"
	"backend/internal/services"
	"backend/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type GoalHandler struct {
	goalService *services.GoalService
}

func NewGoalHandler(goalService *services.GoalService) *GoalHandler {
	return &GoalHandler{
		goalService: goalService,
	}
}

// Create godoc
// @Summary Create goal
// @Tags goals
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body domain.CreateGoalRequest true "Create goal"
// @Success 201 {object} domain.Goal
// @Failure 400 {object} map[string]string
// @Router /protected/goals [post]
func (h *GoalHandler) Create(c *gin.Context) {
	var req domain.CreateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetInt64("user_id")

	g := &domain.Goal{
		UserID:       userID,
		Name:         req.Name,
		TargetAmount: req.TargetAmount,
		TargetDate:   req.TargetDate,
		Description:  req.Description,
	}

	if err := h.goalService.Create(c.Request.Context(), g); err != nil {
		utils.Log.Errorf("goal create error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, g)
}

// GetAll godoc
// @Summary Get all goals
// @Tags goals
// @Security BearerAuth
// @Produce json
// @Success 200 {array} domain.Goal
// @Router /protected/goals [get]
func (h *GoalHandler) GetAll(c *gin.Context) {
	userID := c.GetInt64("user_id")

	list, err := h.goalService.GetAllByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

// GetByID godoc
// @Summary Get goal by ID
// @Tags goals
// @Security BearerAuth
// @Produce json
// @Param id path int true "Goal ID"
// @Success 200 {object} domain.Goal
// @Failure 404 {object} map[string]string
// @Router /protected/goals/{id} [get]
func (h *GoalHandler) GetByID(c *gin.Context) {
	userID := c.GetInt64("user_id")

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	g, err := h.goalService.GetByID(c.Request.Context(), id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, g)
}

// Update godoc
// @Summary Update goal
// @Tags goals
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Goal ID"
// @Param request body domain.UpdateGoalRequest true "Update goal"
// @Success 200 {object} map[string]string
// @Router /protected/goals/{id} [put]
func (h *GoalHandler) Update(c *gin.Context) {
	var req domain.UpdateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID := c.GetInt64("user_id")

	g := &domain.Goal{
		ID:            id,
		UserID:        userID,
		Name:          req.Name,
		TargetAmount:  req.TargetAmount,
		CurrentAmount: req.CurrentAmount,
		TargetDate:    req.TargetDate,
		Description:   req.Description,
		IsCompleted:   req.IsCompleted,
	}

	if err := h.goalService.Update(c.Request.Context(), g); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// Delete godoc
// @Summary Delete goal
// @Tags goals
// @Security BearerAuth
// @Param id path int true "Goal ID"
// @Success 200 {object} map[string]string
// @Router /protected/goals/{id} [delete]
func (h *GoalHandler) Delete(c *gin.Context) {
	userID := c.GetInt64("user_id")

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.goalService.Delete(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
