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

type PermissionUseCase struct {
	DB                   *gorm.DB
	Log                  *logrus.Logger
	Validate             *validator.Validate
	PermissionRepository *repository.PermissionRepository
}

func NewPermissionUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	permissionRepository *repository.PermissionRepository) *PermissionUseCase {
	return &PermissionUseCase{
		DB:                   db,
		Log:                  logger,
		Validate:             validate,
		PermissionRepository: permissionRepository,
	}
}

func (c *PermissionUseCase) Create(ctx context.Context, request *model.CreatePermissionRequest) (*model.PermissionResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Check if permission already exists
	total, err := c.PermissionRepository.CountById(tx, request.ID)
	if err != nil {
		c.Log.Warnf("Failed count permission from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if total > 0 {
		c.Log.Warnf("Permission already exists")
		return nil, fiber.ErrConflict
	}

	permission := &entity.Permission{
		ID:          request.ID,
		Name:        request.Name,
		Resource:    request.Resource,
		Action:      request.Action,
		Description: request.Description,
	}

	if err := c.PermissionRepository.Create(tx, permission); err != nil {
		c.Log.Warnf("Failed create permission to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.PermissionToResponse(permission), nil
}

func (c *PermissionUseCase) Get(ctx context.Context, request *model.GetPermissionRequest) (*model.PermissionResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	permission := new(entity.Permission)
	if err := c.PermissionRepository.FindById(tx, permission, request.ID); err != nil {
		c.Log.Warnf("Failed find permission by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.PermissionToResponse(permission), nil
}

func (c *PermissionUseCase) List(ctx context.Context) ([]model.PermissionResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var permissions []entity.Permission
	if err := c.PermissionRepository.FindAll(tx, &permissions); err != nil {
		c.Log.Warnf("Failed to list permissions : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.PermissionsToResponse(&permissions), nil
}

func (c *PermissionUseCase) Update(ctx context.Context, request *model.UpdatePermissionRequest) (*model.PermissionResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	permission := new(entity.Permission)
	if err := c.PermissionRepository.FindById(tx, permission, request.ID); err != nil {
		c.Log.Warnf("Failed find permission by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	if request.Name != "" {
		permission.Name = request.Name
	}

	if request.Resource != "" {
		permission.Resource = request.Resource
	}

	if request.Action != "" {
		permission.Action = request.Action
	}

	if request.Description != "" {
		permission.Description = request.Description
	}

	if err := c.PermissionRepository.Update(tx, permission); err != nil {
		c.Log.Warnf("Failed update permission : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.PermissionToResponse(permission), nil
}

func (c *PermissionUseCase) Delete(ctx context.Context, request *model.GetPermissionRequest) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return false, fiber.ErrBadRequest
	}

	permission := new(entity.Permission)
	if err := c.PermissionRepository.FindById(tx, permission, request.ID); err != nil {
		c.Log.Warnf("Failed find permission by id : %+v", err)
		return false, fiber.ErrNotFound
	}

	if err := c.PermissionRepository.Delete(tx, permission); err != nil {
		c.Log.Warnf("Failed delete permission : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	return true, nil
}
