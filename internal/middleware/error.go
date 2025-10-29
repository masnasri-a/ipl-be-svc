package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ipl-be-svc/pkg/utils"
)

// ErrorHandler is a middleware that handles panics and errors
func ErrorHandler() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if _, ok := recovered.(string); ok {
			utils.InternalServerErrorResponse(c, "Internal server error", nil)
		} else {
			utils.InternalServerErrorResponse(c, "Internal server error", nil)
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}

// NoRouteHandler handles 404 errors
func NoRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.NotFoundResponse(c, "Route not found")
	}
}

// NoMethodHandler handles 405 errors
func NoMethodHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.ErrorResponse(c, http.StatusMethodNotAllowed, "Method not allowed", nil)
	}
}