package usecase

import (
	"context"
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/model/converter"
	"golang-clean-architecture/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RoleUseCase struct {
	DB                   *gorm.DB
	Log                  *logrus.Logger
	Validate             *validator.Validate
	RoleRepository       *repository.RoleRepository
	PermissionRepository *repository.PermissionRepository
}

func NewRoleUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	roleRepository *repository.RoleRepository, permissionRepository *repository.PermissionRepository) *RoleUseCase {
	return &RoleUseCase{
		DB:                   db,
		Log:                  logger,
		Validate:             validate,
		RoleRepository:       roleRepository,
		PermissionRepository: permissionRepository,
	}
}

func (c *RoleUseCase) Create(ctx context.Context, request *model.CreateRoleRequest) (*model.RoleResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Check if role already exists
	total, err := c.RoleRepository.CountById(tx, request.ID)
	if err != nil {
		c.Log.Warnf("Failed count role from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if total > 0 {
		c.Log.Warnf("Role already exists")
		return nil, fiber.ErrConflict
	}

	role := &entity.Role{
		ID:          request.ID,
		Name:        request.Name,
		Description: request.Description,
	}

	if err := c.RoleRepository.Create(tx, role); err != nil {
		c.Log.Warnf("Failed create role to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Assign permissions if provided
	if len(request.Permissions) > 0 {
		for _, permID := range request.Permissions {
			permission := new(entity.Permission)
			if err := c.PermissionRepository.FindById(tx, permission, permID); err != nil {
				c.Log.Warnf("Permission not found: %s", permID)
				return nil, fiber.NewError(fiber.StatusBadRequest, "Permission not found: "+permID)
			}
			if err := c.RoleRepository.AssignPermission(tx, role, permission); err != nil {
				c.Log.Warnf("Failed to assign permission: %+v", err)
				return nil, fiber.ErrInternalServerError
			}
		}
		// Reload role with permissions
		if err := c.RoleRepository.FindWithPermissions(tx, role, role.ID); err != nil {
			c.Log.Warnf("Failed to reload role with permissions: %+v", err)
			return nil, fiber.ErrInternalServerError
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.RoleToResponse(role), nil
}

func (c *RoleUseCase) Get(ctx context.Context, request *model.GetRoleRequest) (*model.RoleResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	role := new(entity.Role)
	if err := c.RoleRepository.FindWithPermissions(tx, role, request.ID); err != nil {
		c.Log.Warnf("Failed find role by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.RoleToResponse(role), nil
}

func (c *RoleUseCase) List(ctx context.Context) ([]model.RoleResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var roles []entity.Role
	if err := c.RoleRepository.FindAllWithPermissions(tx, &roles); err != nil {
		c.Log.Warnf("Failed to list roles : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.RolesToResponse(&roles), nil
}

func (c *RoleUseCase) Update(ctx context.Context, request *model.UpdateRoleRequest) (*model.RoleResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	role := new(entity.Role)
	if err := c.RoleRepository.FindById(tx, role, request.ID); err != nil {
		c.Log.Warnf("Failed find role by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	if request.Name != "" {
		role.Name = request.Name
	}

	if request.Description != "" {
		role.Description = request.Description
	}

	if err := c.RoleRepository.Update(tx, role); err != nil {
		c.Log.Warnf("Failed update role : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Update permissions if provided
	if request.Permissions != nil {
		// Clear existing permissions
		if err := tx.Model(role).Association("Permissions").Clear(); err != nil {
			c.Log.Warnf("Failed to clear permissions: %+v", err)
			return nil, fiber.ErrInternalServerError
		}

		// Assign new permissions
		for _, permID := range request.Permissions {
			permission := new(entity.Permission)
			if err := c.PermissionRepository.FindById(tx, permission, permID); err != nil {
				c.Log.Warnf("Permission not found: %s", permID)
				return nil, fiber.NewError(fiber.StatusBadRequest, "Permission not found: "+permID)
			}
			if err := c.RoleRepository.AssignPermission(tx, role, permission); err != nil {
				c.Log.Warnf("Failed to assign permission: %+v", err)
				return nil, fiber.ErrInternalServerError
			}
		}
	}

	// Reload role with permissions
	if err := c.RoleRepository.FindWithPermissions(tx, role, role.ID); err != nil {
		c.Log.Warnf("Failed to reload role with permissions: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.RoleToResponse(role), nil
}

func (c *RoleUseCase) Delete(ctx context.Context, request *model.GetRoleRequest) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return false, fiber.ErrBadRequest
	}

	role := new(entity.Role)
	if err := c.RoleRepository.FindById(tx, role, request.ID); err != nil {
		c.Log.Warnf("Failed find role by id : %+v", err)
		return false, fiber.ErrNotFound
	}

	if err := c.RoleRepository.Delete(tx, role); err != nil {
		c.Log.Warnf("Failed delete role : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	return true, nil
}
