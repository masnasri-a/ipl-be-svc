package handler

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"

	"ipl-be-svc/internal/service"
	"ipl-be-svc/pkg/logger"
)

// Routes sets up all API routes
func SetupRoutes(
	router *gin.Engine,
	menuService service.MenuService,
	paymentService service.PaymentService,
	logger *logger.Logger,
) {
	// Initialize handlers
	menuHandler := NewMenuHandler(menuService)
	paymentHandler := NewPaymentHandler(paymentService, logger)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", HealthCheck)

		// Menu routes
		menus := v1.Group("/menus")
		{
			menus.GET("/user/:id", menuHandler.GetMenusByUserID)
		}

		// Payment routes
		payments := v1.Group("/payments")
		{
			payments.POST("/billing/:id/link", paymentHandler.CreatePaymentLink)
		}
	}
}

// HealthCheck handles GET /api/v1/health
// @Summary Health check
// @Description Check if the service is running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Service is running"
// @Router /api/v1/health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"message": "Server is running",
		"service": "IPL Backend Service",
	})
}