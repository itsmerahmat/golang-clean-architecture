package route

import (
	"golang-clean-architecture/internal/delivery/http"
	"golang-clean-architecture/internal/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type RouteConfig struct {
	App                  *fiber.App
	UserController       *http.UserController
	ContactController    *http.ContactController
	AddressController    *http.AddressController
	RoleController       *http.RoleController
	PermissionController *http.PermissionController
	UserRoleController   *http.UserRoleController
	AuthMiddleware       fiber.Handler
	Log                  *logrus.Logger
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
	c.SetupAdminRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.Post("/api/users", c.UserController.Register)
	c.App.Post("/api/users/_login", c.UserController.Login)
}

func (c *RouteConfig) SetupAuthRoute() {
	authMiddleware := c.AuthMiddleware
	
	c.App.Delete("/api/users", authMiddleware, c.UserController.Logout)
	c.App.Patch("/api/users/_current", authMiddleware, c.UserController.Update)
	c.App.Get("/api/users/_current", authMiddleware, c.UserController.Current)

	c.App.Get("/api/contacts", authMiddleware, c.ContactController.List)
	c.App.Post("/api/contacts", authMiddleware, c.ContactController.Create)
	c.App.Put("/api/contacts/:contactId", authMiddleware, c.ContactController.Update)
	c.App.Get("/api/contacts/:contactId", authMiddleware, c.ContactController.Get)
	c.App.Delete("/api/contacts/:contactId", authMiddleware, c.ContactController.Delete)

	c.App.Get("/api/contacts/:contactId/addresses", authMiddleware, c.AddressController.List)
	c.App.Post("/api/contacts/:contactId/addresses", authMiddleware, c.AddressController.Create)
	c.App.Put("/api/contacts/:contactId/addresses/:addressId", authMiddleware, c.AddressController.Update)
	c.App.Get("/api/contacts/:contactId/addresses/:addressId", authMiddleware, c.AddressController.Get)
	c.App.Delete("/api/contacts/:contactId/addresses/:addressId", authMiddleware, c.AddressController.Delete)
}

func (c *RouteConfig) SetupAdminRoute() {
	authMiddleware := c.AuthMiddleware
	adminRole := middleware.RequireRole(c.Log, "admin")
	
	// Role management endpoints - GET requires auth, POST/PUT/DELETE require admin
	c.App.Get("/api/roles", authMiddleware, c.RoleController.List)
	c.App.Get("/api/roles/:roleId", authMiddleware, c.RoleController.Get)
	c.App.Post("/api/roles", authMiddleware, adminRole, c.RoleController.Create)
	c.App.Put("/api/roles/:roleId", authMiddleware, adminRole, c.RoleController.Update)
	c.App.Delete("/api/roles/:roleId", authMiddleware, adminRole, c.RoleController.Delete)

	// Permission management endpoints - GET requires auth, POST/PUT/DELETE require admin
	c.App.Get("/api/permissions", authMiddleware, c.PermissionController.List)
	c.App.Get("/api/permissions/:permissionId", authMiddleware, c.PermissionController.Get)
	c.App.Post("/api/permissions", authMiddleware, adminRole, c.PermissionController.Create)
	c.App.Put("/api/permissions/:permissionId", authMiddleware, adminRole, c.PermissionController.Update)
	c.App.Delete("/api/permissions/:permissionId", authMiddleware, adminRole, c.PermissionController.Delete)

	// User role assignment endpoints - all require admin
	c.App.Get("/api/users/:userId/roles", authMiddleware, adminRole, c.UserRoleController.GetUserRoles)
	c.App.Post("/api/users/:userId/roles/:roleId", authMiddleware, adminRole, c.UserRoleController.AssignRole)
	c.App.Delete("/api/users/:userId/roles/:roleId", authMiddleware, adminRole, c.UserRoleController.RemoveRole)
}

