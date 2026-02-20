package middleware

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"

	"echo-base/utils"
)

// BearerAuthMiddleware validates bearer token in Authorization header
func BearerAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(401, "missing authorization header")
		}

		// Extract token from Bearer <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return echo.NewHTTPError(401, "invalid authorization header format")
		}

		token := parts[1]

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

// OptionalBearerAuthMiddleware validates bearer token if provided
func OptionalBearerAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return next(c)
		}

		// Extract token from Bearer <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return next(c)
		}

		token := parts[1]

		// Validate token
		claims, err := utils.ValidateToken(token)
		if err == nil {
			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)
			c.Set("role_id", claims.RoleID)
		}

		return next(c)
	}
}
