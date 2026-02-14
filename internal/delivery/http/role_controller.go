package http

import (
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type RoleController struct {
	Log     *logrus.Logger
	UseCase *usecase.RoleUseCase
}

func NewRoleController(useCase *usecase.RoleUseCase, logger *logrus.Logger) *RoleController {
	return &RoleController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *RoleController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateRoleRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create role : %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.RoleResponse]{Data: response})
}

func (c *RoleController) Get(ctx *fiber.Ctx) error {
	roleID := ctx.Params("roleId")

	request := &model.GetRoleRequest{
		ID: roleID,
	}

	response, err := c.UseCase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to get role")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.RoleResponse]{Data: response})
}

func (c *RoleController) List(ctx *fiber.Ctx) error {
	response, err := c.UseCase.List(ctx.UserContext())
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to list roles")
		return err
	}

	return ctx.JSON(model.WebResponse[[]model.RoleResponse]{Data: response})
}

func (c *RoleController) Update(ctx *fiber.Ctx) error {
	roleID := ctx.Params("roleId")

	request := new(model.UpdateRoleRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	request.ID = roleID
	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to update role")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.RoleResponse]{Data: response})
}

func (c *RoleController) Delete(ctx *fiber.Ctx) error {
	roleID := ctx.Params("roleId")

	request := &model.GetRoleRequest{
		ID: roleID,
	}

	response, err := c.UseCase.Delete(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to delete role")
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: response})
}
