package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// RequireRole creates a middleware that checks if the user has the required role
func RequireRole(logger *logrus.Logger, roles ...string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		auth := GetUser(ctx)
		if auth == nil {
			logger.Warn("User not authenticated")
			return fiber.ErrUnauthorized
		}

		if !auth.HasAnyRole(roles...) {
			logger.Warnf("User %s does not have required role. Required: %v, Has: %v", auth.ID, roles, auth.Roles)
			return fiber.ErrForbidden
		}

		return ctx.Next()
	}
}

// RequirePermission creates a middleware that checks if the user has the required permission
func RequirePermission(logger *logrus.Logger, permissions ...string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		auth := GetUser(ctx)
		if auth == nil {
			logger.Warn("User not authenticated")
			return fiber.ErrUnauthorized
		}

		if !auth.HasAnyPermission(permissions...) {
			logger.Warnf("User %s does not have required permission. Required: %v, Has: %v", auth.ID, permissions, auth.Permissions)
			return fiber.ErrForbidden
		}

		return ctx.Next()
	}
}

// RequireAllRoles creates a middleware that checks if the user has all required roles
func RequireAllRoles(logger *logrus.Logger, roles ...string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		auth := GetUser(ctx)
		if auth == nil {
			logger.Warn("User not authenticated")
			return fiber.ErrUnauthorized
		}

		for _, role := range roles {
			if !auth.HasRole(role) {
				logger.Warnf("User %s does not have required role %s", auth.ID, role)
				return fiber.ErrForbidden
			}
		}

		return ctx.Next()
	}
}

// RequireAllPermissions creates a middleware that checks if the user has all required permissions
func RequireAllPermissions(logger *logrus.Logger, permissions ...string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		auth := GetUser(ctx)
		if auth == nil {
			logger.Warn("User not authenticated")
			return fiber.ErrUnauthorized
		}

		for _, permission := range permissions {
			if !auth.HasPermission(permission) {
				logger.Warnf("User %s does not have required permission %s", auth.ID, permission)
				return fiber.ErrForbidden
			}
		}

		return ctx.Next()
	}
}
