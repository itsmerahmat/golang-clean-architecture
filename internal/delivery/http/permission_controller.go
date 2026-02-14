package http

import (
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type PermissionController struct {
	Log     *logrus.Logger
	UseCase *usecase.PermissionUseCase
}

func NewPermissionController(useCase *usecase.PermissionUseCase, logger *logrus.Logger) *PermissionController {
	return &PermissionController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *PermissionController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreatePermissionRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create permission : %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.PermissionResponse]{Data: response})
}

func (c *PermissionController) Get(ctx *fiber.Ctx) error {
	permissionID := ctx.Params("permissionId")

	request := &model.GetPermissionRequest{
		ID: permissionID,
	}

	response, err := c.UseCase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to get permission")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.PermissionResponse]{Data: response})
}

func (c *PermissionController) List(ctx *fiber.Ctx) error {
	response, err := c.UseCase.List(ctx.UserContext())
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to list permissions")
		return err
	}

	return ctx.JSON(model.WebResponse[[]model.PermissionResponse]{Data: response})
}

func (c *PermissionController) Update(ctx *fiber.Ctx) error {
	permissionID := ctx.Params("permissionId")

	request := new(model.UpdatePermissionRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	request.ID = permissionID
	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to update permission")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.PermissionResponse]{Data: response})
}

func (c *PermissionController) Delete(ctx *fiber.Ctx) error {
	permissionID := ctx.Params("permissionId")

	request := &model.GetPermissionRequest{
		ID: permissionID,
	}

	response, err := c.UseCase.Delete(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to delete permission")
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: response})
}
