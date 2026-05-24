package handlers

import (
	"backend/internal/domain"
	"backend/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ReportHandler struct {
	reportService *services.ReportService
}

func NewReportHandler(reportService *services.ReportService) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
	}
}

// Generate godoc
// @Summary Generate report
// @Tags reports
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body domain.GenerateReportRequest true "Generate report"
// @Success 201 {object} domain.Report
// @Router /protected/reports [post]
func (h *ReportHandler) Generate(c *gin.Context) {
	var req domain.GenerateReportRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetInt64("user_id")

	report, err := h.reportService.Generate(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, report)
}

// GetAll godoc
// @Summary Get all reports
// @Tags reports
// @Security BearerAuth
// @Produce json
// @Success 200 {array} domain.Report
// @Router /protected/reports [get]
func (h *ReportHandler) GetAll(c *gin.Context) {
	userID := c.GetInt64("user_id")

	reports, err := h.reportService.GetAllByUser(
		c.Request.Context(),
		userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reports)
}

// Download godoc
// @Summary Download report file
// @Tags reports
// @Security BearerAuth
// @Param id path int true "Report ID"
// @Produce application/octet-stream
// @Success 200 {file} file
// @Router /protected/reports/{id}/download [get]
func (h *ReportHandler) Download(c *gin.Context) {
	userID := c.GetInt64("user_id")

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid report id"})
		return
	}

	report, err := h.reportService.GetByID(
		c.Request.Context(),
		id,
		userID,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.FileAttachment(report.FilePath, report.FileName)
}
