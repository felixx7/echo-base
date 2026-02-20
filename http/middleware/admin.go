package middleware

import (
	"fmt"

	"github.com/labstack/echo/v4"

	"echo-base/utils"
)

// AdminRoleMiddleware validates if user has admin role
func AdminRoleMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(401, "missing authorization header")
		}

		// Extract token from Bearer <token>
		parts := len(c.Request().Header.Values("Authorization"))
		if parts == 0 {
			return echo.NewHTTPError(401, "missing authorization header")
		}

		// Get token from context (set by BearerAuthMiddleware)
		userID := c.Get("user_id")
		roleID := c.Get("role_id")

		if userID == nil || roleID == nil {
			return echo.NewHTTPError(401, "unauthorized")
		}

		// Check if role_id is 2 (admin role)
		if roleID != int64(2) {
			return echo.NewHTTPError(403, "you don't have permission to access this resource")
		}

		return next(c)
	}
}

// BearerAuthMiddlewareWithRole validates bearer token and extracts role
func BearerAuthMiddlewareWithRole(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(401, "missing authorization header")
		}

		// Parse Bearer token
		var token string
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		} else {
			return echo.NewHTTPError(401, "invalid authorization header format")
		}

		// Validate token
		claims, err := utils.ValidateToken(token)
		if err != nil {
			return echo.NewHTTPError(401, fmt.Sprintf("invalid token: %v", err))
		}

		// Store claims in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("role_id", claims.RoleID)

		return next(c)
	}
}
