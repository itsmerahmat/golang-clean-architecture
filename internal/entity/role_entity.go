package entity

// Role is a struct that represents a role entity
type Role struct {
	ID          string       `gorm:"column:id;primaryKey"`
	Name        string       `gorm:"column:name;unique"`
	Description string       `gorm:"column:description"`
	CreatedAt   int64        `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt   int64        `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}

func (r *Role) TableName() string {
	return "roles"
}
