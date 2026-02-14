package model

type PermissionResponse struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Resource    string `json:"resource,omitempty"`
	Action      string `json:"action,omitempty"`
	Description string `json:"description,omitempty"`
	CreatedAt   int64  `json:"created_at,omitempty"`
	UpdatedAt   int64  `json:"updated_at,omitempty"`
}

type CreatePermissionRequest struct {
	ID          string `json:"id" validate:"required,max=100"`
	Name        string `json:"name" validate:"required,max=100"`
	Resource    string `json:"resource" validate:"required,max=100"`
	Action      string `json:"action" validate:"required,max=50"`
	Description string `json:"description,omitempty" validate:"max=255"`
}

type UpdatePermissionRequest struct {
	ID          string `json:"-" validate:"required,max=100"`
	Name        string `json:"name,omitempty" validate:"max=100"`
	Resource    string `json:"resource,omitempty" validate:"max=100"`
	Action      string `json:"action,omitempty" validate:"max=50"`
	Description string `json:"description,omitempty" validate:"max=255"`
}

type GetPermissionRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
