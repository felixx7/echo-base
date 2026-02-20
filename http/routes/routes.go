package routes

import (
	"github.com/labstack/echo/v4"

	"echo-base/http/handler"
	"echo-base/http/middleware"
)

// RegisterRoutes registers all HTTP routes for the application
func RegisterRoutes(e *echo.Echo, h *handler.UserHandler) {
	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	const apiVersion = "/api/v1"

	api := e.Group(apiVersion)

	// Auth routes
	authRoutes := api.Group("/auth")
	authRoutes.POST("/register", h.Register)
	authRoutes.POST("/login", h.Login)

	// Admin routes (admin only)
	adminRoutes := api.Group("/admin")
	adminRoutes.Use(middleware.BearerAuthMiddleware)
	adminRoutes.Use(middleware.AdminRoleMiddleware)
	// Add admin-only endpoints here

	// User routes (protected)
	userRoutes := api.Group("/users")
	userRoutes.Use(middleware.BearerAuthMiddleware)
	userRoutes.GET("", h.GetAll)
	userRoutes.GET("/pagination", h.GetAllPagination)
	userRoutes.GET("/:id", h.GetByID)
	userRoutes.PUT("/:id", h.Update)
	userRoutes.DELETE("/:id", h.Delete)

	// Profile route (protected)
	apiRoutes := api.Group("/profile")
	apiRoutes.Use(middleware.BearerAuthMiddleware)
	apiRoutes.GET("", h.GetProfile)
}
