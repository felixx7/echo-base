package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// LoggerMiddleware returns logger middleware configuration
func LoggerMiddleware() echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${status} ${method} ${path} latency=${latency_human}\n",
	})
}

// RecoverMiddleware returns recover middleware configuration
func RecoverMiddleware() echo.MiddlewareFunc {
	return middleware.Recover()
}
