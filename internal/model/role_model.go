package model

type RoleResponse struct {
	ID          string               `json:"id,omitempty"`
	Name        string               `json:"name,omitempty"`
	Description string               `json:"description,omitempty"`
	Permissions []PermissionResponse `json:"permissions,omitempty"`
	CreatedAt   int64                `json:"created_at,omitempty"`
	UpdatedAt   int64                `json:"updated_at,omitempty"`
}

type CreateRoleRequest struct {
	ID          string   `json:"id" validate:"required,max=100"`
	Name        string   `json:"name" validate:"required,max=100"`
	Description string   `json:"description,omitempty" validate:"max=255"`
	Permissions []string `json:"permissions,omitempty"`
}

type UpdateRoleRequest struct {
	ID          string   `json:"-" validate:"required,max=100"`
	Name        string   `json:"name,omitempty" validate:"max=100"`
	Description string   `json:"description,omitempty" validate:"max=255"`
	Permissions []string `json:"permissions,omitempty"`
}

type GetRoleRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}

type AssignRoleRequest struct {
	UserID string `json:"user_id" validate:"required,max=100"`
	RoleID string `json:"role_id" validate:"required,max=100"`
}
