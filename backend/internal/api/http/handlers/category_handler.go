package handlers

import (
	"backend/internal/domain"
	"backend/internal/services"
	"backend/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CategoryHandler struct {
	categoryService *services.CategoryService
}

func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// Create godoc
// @Summary Create category
// @Tags categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body domain.CreateCategoryRequest true "Create category"
// @Success 201 {object} domain.Category
// @Failure 400 {object} map[string]string
// @Router /protected/categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	var req domain.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetInt64("user_id")

	category := &domain.Category{
		UserID: userID,
		Name:   req.Name,
		Type:   req.Type,
		Color:  req.Color,
		Icon:   req.Icon,
	}

	if err := h.categoryService.Create(c.Request.Context(), category); err != nil {
		utils.Log.Errorf("Create category error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      category.ID,
		"user_id": category.UserID,
		"name":    category.Name,
		"type":    category.Type,
		"color":   category.Color,
		"icon":    category.Icon,
	})
}

// GetByID godoc
// @Summary Get category by ID
// @Tags categories
// @Security BearerAuth
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} domain.Category
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /protected/categories/{id} [get]
func (h *CategoryHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_id"})
		return
	}

	category, err := h.categoryService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) GetAll(c *gin.Context) {
	userID := c.GetInt64("user_id")

	categories, err := h.categoryService.GetAllByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// Update godoc
// @Summary Update category
// @Tags categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body domain.UpdateCategoryRequest true "Update category"
// @Success 200 {object} domain.Category
// @Failure 400 {object} map[string]string
// @Router /protected/categories [put]
func (h *CategoryHandler) Update(c *gin.Context) {
	var req domain.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetInt64("user_id")

	category := &domain.Category{
		ID:     req.ID,
		UserID: userID,
		Name:   req.Name,
		Type:   req.Type,
		Color:  req.Color,
		Icon:   req.Icon,
	}

	if err := h.categoryService.Update(c.Request.Context(), category); err != nil {
		utils.Log.Errorf("Update category error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      category.ID,
		"user_id": category.UserID,
		"name":    category.Name,
		"type":    category.Type,
		"color":   category.Color,
		"icon":    category.Icon,
	})
}

// Delete godoc
// @Summary Delete category
// @Tags categories
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /protected/categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID := c.GetInt64("user_id")

	if err := h.categoryService.Delete(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "category deleted",
	})
}
