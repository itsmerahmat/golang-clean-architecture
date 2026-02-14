package converter

import (
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"
)

func PermissionToResponse(permission *entity.Permission) *model.PermissionResponse {
	return &model.PermissionResponse{
		ID:          permission.ID,
		Name:        permission.Name,
		Resource:    permission.Resource,
		Action:      permission.Action,
		Description: permission.Description,
		CreatedAt:   permission.CreatedAt,
		UpdatedAt:   permission.UpdatedAt,
	}
}

func PermissionsToResponse(permissions *[]entity.Permission) []model.PermissionResponse {
	responses := make([]model.PermissionResponse, 0)
	for _, permission := range *permissions {
		responses = append(responses, *PermissionToResponse(&permission))
	}
	return responses
}
