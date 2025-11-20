package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"ipl-be-svc/internal/service"
	"ipl-be-svc/pkg/logger"
)

// Routes sets up all API routes
func SetupRoutes(
	router *gin.Engine,
	menuService service.MenuService,
	paymentService service.PaymentService,
	userService service.UserService,
	billingService service.BillingService,
	logger *logger.Logger,
) {
	// Initialize handlers
	menuHandler := NewMenuHandler(menuService)
	paymentHandler := NewPaymentHandler(paymentService, logger)
	userHandler := NewUserHandler(userService, logger)
	bulkBillingHandler := NewBulkBillingHandler(billingService, logger)

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

		// User routes
		users := v1.Group("/users")
		{
			users.GET("/profile/:user_id", userHandler.GetUserDetailByProfileID)
			users.GET("/penghuni", userHandler.GetPenghuniUsers)
		}

		// Billing routes
		billings := v1.Group("/billings")
		{
			billings.POST("/bulk-monthly", bulkBillingHandler.CreateBulkMonthlyBillings)
		}
	}
}

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"message": "Server is running",
		"service": "IPL Backend Service",
	})
}
