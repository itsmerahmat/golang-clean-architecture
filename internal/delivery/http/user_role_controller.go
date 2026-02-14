package http

import (
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRoleController struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	UserRepository *repository.UserRepository
	RoleRepository *repository.RoleRepository
}

func NewUserRoleController(db *gorm.DB, logger *logrus.Logger,
	userRepository *repository.UserRepository, roleRepository *repository.RoleRepository) *UserRoleController {
	return &UserRoleController{
		DB:             db,
		Log:            logger,
		UserRepository: userRepository,
		RoleRepository: roleRepository,
	}
}

func (c *UserRoleController) AssignRole(ctx *fiber.Ctx) error {
	userID := ctx.Params("userId")
	roleID := ctx.Params("roleId")

	tx := c.DB.WithContext(ctx.UserContext()).Begin()
	defer tx.Rollback()

	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, userID); err != nil {
		c.Log.Warnf("User not found: %+v", err)
		return fiber.ErrNotFound
	}

	role := new(entity.Role)
	if err := c.RoleRepository.FindById(tx, role, roleID); err != nil {
		c.Log.Warnf("Role not found: %+v", err)
		return fiber.ErrNotFound
	}

	if err := c.UserRepository.AssignRole(tx, user, role); err != nil {
		c.Log.Warnf("Failed to assign role: %+v", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return fiber.ErrInternalServerError
	}

	return ctx.JSON(model.WebResponse[bool]{Data: true})
}

func (c *UserRoleController) RemoveRole(ctx *fiber.Ctx) error {
	userID := ctx.Params("userId")
	roleID := ctx.Params("roleId")

	tx := c.DB.WithContext(ctx.UserContext()).Begin()
	defer tx.Rollback()

	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, userID); err != nil {
		c.Log.Warnf("User not found: %+v", err)
		return fiber.ErrNotFound
	}

	role := new(entity.Role)
	if err := c.RoleRepository.FindById(tx, role, roleID); err != nil {
		c.Log.Warnf("Role not found: %+v", err)
		return fiber.ErrNotFound
	}

	if err := c.UserRepository.RemoveRole(tx, user, role); err != nil {
		c.Log.Warnf("Failed to remove role: %+v", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return fiber.ErrInternalServerError
	}

	return ctx.JSON(model.WebResponse[bool]{Data: true})
}

func (c *UserRoleController) GetUserRoles(ctx *fiber.Ctx) error {
	userID := ctx.Params("userId")

	tx := c.DB.WithContext(ctx.UserContext()).Begin()
	defer tx.Rollback()

	user := new(entity.User)
	if err := c.UserRepository.FindByIdWithRoles(tx, user, userID); err != nil {
		c.Log.Warnf("User not found: %+v", err)
		return fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return fiber.ErrInternalServerError
	}

	roles := make([]string, 0)
	for _, role := range user.Roles {
		roles = append(roles, role.Name)
	}

	return ctx.JSON(model.WebResponse[[]string]{Data: roles})
}
