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
	c.App.Use(c.AuthMiddleware)
	c.App.Delete("/api/users", c.UserController.Logout)
	c.App.Patch("/api/users/_current", c.UserController.Update)
	c.App.Get("/api/users/_current", c.UserController.Current)

	c.App.Get("/api/contacts", c.ContactController.List)
	c.App.Post("/api/contacts", c.ContactController.Create)
	c.App.Put("/api/contacts/:contactId", c.ContactController.Update)
	c.App.Get("/api/contacts/:contactId", c.ContactController.Get)
	c.App.Delete("/api/contacts/:contactId", c.ContactController.Delete)

	c.App.Get("/api/contacts/:contactId/addresses", c.AddressController.List)
	c.App.Post("/api/contacts/:contactId/addresses", c.AddressController.Create)
	c.App.Put("/api/contacts/:contactId/addresses/:addressId", c.AddressController.Update)
	c.App.Get("/api/contacts/:contactId/addresses/:addressId", c.AddressController.Get)
	c.App.Delete("/api/contacts/:contactId/addresses/:addressId", c.AddressController.Delete)
}

func (c *RouteConfig) SetupAdminRoute() {
	// Role management endpoints
	adminRole := middleware.RequireRole(c.Log, "admin")
	
	c.App.Get("/api/roles", c.RoleController.List)
	c.App.Get("/api/roles/:roleId", c.RoleController.Get)
	c.App.Post("/api/roles", adminRole, c.RoleController.Create)
	c.App.Put("/api/roles/:roleId", adminRole, c.RoleController.Update)
	c.App.Delete("/api/roles/:roleId", adminRole, c.RoleController.Delete)

	// Permission management endpoints
	c.App.Get("/api/permissions", c.PermissionController.List)
	c.App.Get("/api/permissions/:permissionId", c.PermissionController.Get)
	c.App.Post("/api/permissions", adminRole, c.PermissionController.Create)
	c.App.Put("/api/permissions/:permissionId", adminRole, c.PermissionController.Update)
	c.App.Delete("/api/permissions/:permissionId", adminRole, c.PermissionController.Delete)

	// User role assignment endpoints - require admin role
	c.App.Get("/api/users/:userId/roles", adminRole, c.UserRoleController.GetUserRoles)
	c.App.Post("/api/users/:userId/roles/:roleId", adminRole, c.UserRoleController.AssignRole)
	c.App.Delete("/api/users/:userId/roles/:roleId", adminRole, c.UserRoleController.RemoveRole)
}

