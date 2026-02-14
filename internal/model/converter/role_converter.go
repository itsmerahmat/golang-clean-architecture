package converter

import (
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"
)

func RoleToResponse(role *entity.Role) *model.RoleResponse {
	permissions := make([]model.PermissionResponse, 0)
	for _, permission := range role.Permissions {
		permissions = append(permissions, *PermissionToResponse(&permission))
	}

	return &model.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Permissions: permissions,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}

func RolesToResponse(roles *[]entity.Role) []model.RoleResponse {
	responses := make([]model.RoleResponse, 0)
	for _, role := range *roles {
		responses = append(responses, *RoleToResponse(&role))
	}
	return responses
}
