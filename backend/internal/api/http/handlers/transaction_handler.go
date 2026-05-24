package handlers

import (
	"backend/internal/domain"
	"backend/internal/services"
	"backend/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TransactionHandler struct {
	transactionService *services.TransactionService
}

func NewTransactionHandler(transactionService *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// Create godoc
// @Summary Create transaction
// @Tags transactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body domain.CreateTransactionRequest true "Create transaction"
// @Success 201 {object} domain.Transaction
// @Router /protected/transactions [post]
func (h *TransactionHandler) Create(c *gin.Context) {
	var req domain.CreateTransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetInt64("user_id")

	t := &domain.Transaction{
		UserID:      userID,
		Type:        req.Type,
		Amount:      req.Amount,
		CategoryID:  req.CategoryID,
		Date:        req.Date,
		Description: req.Description,
	}

	if err := h.transactionService.Create(c.Request.Context(), t); err != nil {
		utils.Log.Errorf("transaction create error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, t)
}

// GetAll godoc
// @Summary Get all transactions
// @Tags transactions
// @Security BearerAuth
// @Produce json
// @Success 200 {array} domain.Transaction
// @Router /protected/transactions [get]
func (h *TransactionHandler) GetAll(c *gin.Context) {
	userID := c.GetInt64("user_id")

	list, err := h.transactionService.GetAllByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

// GetByID godoc
// @Summary Get transaction by ID
// @Tags transactions
// @Security BearerAuth
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} domain.Transaction
// @Router /protected/transactions/{id} [get]
func (h *TransactionHandler) GetByID(c *gin.Context) {
	userID := c.GetInt64("user_id")

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	t, err := h.transactionService.GetByID(c.Request.Context(), id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, t)
}

// Update godoc
// @Summary Update transaction
// @Tags transactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Param request body domain.UpdateTransactionRequest true "Update transaction"
// @Success 200 {object} map[string]string
// @Router /protected/transactions/{id} [put]
func (h *TransactionHandler) Update(c *gin.Context) {
	var req domain.UpdateTransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetInt64("user_id")

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	t := &domain.Transaction{
		ID:          id,
		UserID:      userID,
		Type:        req.Type,
		Amount:      req.Amount,
		CategoryID:  req.CategoryID,
		Date:        req.Date,
		Description: req.Description,
	}

	if err := h.transactionService.Update(c.Request.Context(), t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// Delete godoc
// @Summary Delete transaction
// @Tags transactions
// @Security BearerAuth
// @Param id path int true "Transaction ID"
// @Success 200 {object} map[string]string
// @Router /protected/transactions/{id} [delete]
func (h *TransactionHandler) Delete(c *gin.Context) {
	userID := c.GetInt64("user_id")

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.transactionService.Delete(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
