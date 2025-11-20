package handler

import (
	"net/http"

	"ipl-be-svc/internal/service"
	"ipl-be-svc/pkg/logger"

	"github.com/gin-gonic/gin"
)

// BulkBillingRequest represents the request for bulk billing creation
type BulkBillingRequest struct {
	UserIDs []uint `json:"user_ids,omitempty"`                        // Empty means all penghuni users
	Month   int    `json:"month" binding:"required,min=1,max=12"`     // Month 1-12
	Year    int    `json:"year" binding:"required,min=2020,max=2100"` // Reasonable year range
}

// BulkBillingHandler handles bulk billing-related HTTP requests
type BulkBillingHandler struct {
	billingService service.BillingService
	logger         *logger.Logger
}

// NewBulkBillingHandler creates a new BulkBillingHandler instance
func NewBulkBillingHandler(billingService service.BillingService, logger *logger.Logger) *BulkBillingHandler {
	return &BulkBillingHandler{
		billingService: billingService,
		logger:         logger,
	}
}

// CreateBulkMonthlyBillings creates monthly billings for specified users or all penghuni users
// @Summary Create bulk monthly billings
// @Description Create monthly billings for specified user IDs or all penghuni users if user_ids is empty. Requires auth-token cookie.
// @Tags billings
// @Accept json
// @Produce json
// @Param request body BulkBillingRequest true "Bulk billing request with month and year"
// @Success 200 {object} service.BulkBillingResponse "Bulk billing creation result"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/billings/bulk-monthly [post]
func (h *BulkBillingHandler) CreateBulkMonthlyBillings(c *gin.Context) {
	var req BulkBillingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Request body must be valid JSON",
		})
		return
	}

	var response *service.BulkBillingResponse
	var serviceErr error

	if len(req.UserIDs) > 0 {
		// Create for specific users
		response, serviceErr = h.billingService.CreateBulkMonthlyBillings(req.UserIDs, req.Month, req.Year)
	} else {
		// Create for all penghuni users
		response, serviceErr = h.billingService.CreateBulkMonthlyBillingsForAllUsers(req.Month, req.Year)
	}

	if serviceErr != nil {
		h.logger.WithError(serviceErr).Error("Failed to create bulk billings")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create billings",
			"message": serviceErr.Error(),
		})
		return
	}

	h.logger.WithFields(map[string]interface{}{
		"total_users":    response.TotalUsers,
		"total_billings": response.TotalBillings,
		"success_count":  response.SuccessCount,
		"failed_count":   response.FailedCount,
	}).Info("Bulk billings created successfully")

	c.JSON(http.StatusOK, response)
}
