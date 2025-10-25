package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"ipl-be-svc/pkg/logger"
)

// LoggerMiddleware creates a gin middleware for logging
func LoggerMiddleware(logger *logger.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Log using structured logging
		logger.WithFields(logrus.Fields{
			"status_code":  param.StatusCode,
			"latency":      param.Latency,
			"client_ip":    param.ClientIP,
			"method":       param.Method,
			"path":         param.Path,
			"user_agent":   param.Request.UserAgent(),
			"error":        param.ErrorMessage,
			"body_size":    param.BodySize,
			"timestamp":    param.TimeStamp.Format(time.RFC3339),
		}).Info("HTTP Request")

		return ""
	})
}