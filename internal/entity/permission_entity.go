package entity

// Permission is a struct that represents a permission entity
type Permission struct {
	ID          string `gorm:"column:id;primaryKey"`
	Name        string `gorm:"column:name;unique"`
	Resource    string `gorm:"column:resource"`
	Action      string `gorm:"column:action"`
	Description string `gorm:"column:description"`
	CreatedAt   int64  `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt   int64  `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
}

func (p *Permission) TableName() string {
	return "permissions"
}
