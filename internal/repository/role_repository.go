package repository

import (
	"golang-clean-architecture/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RoleRepository struct {
	Repository[entity.Role]
	Log *logrus.Logger
}

func NewRoleRepository(log *logrus.Logger) *RoleRepository {
	return &RoleRepository{
		Log: log,
	}
}

func (r *RoleRepository) FindByName(db *gorm.DB, role *entity.Role, name string) error {
	return db.Where("name = ?", name).First(role).Error
}

func (r *RoleRepository) FindWithPermissions(db *gorm.DB, role *entity.Role, id string) error {
	return db.Preload("Permissions").Where("id = ?", id).First(role).Error
}

func (r *RoleRepository) FindAllWithPermissions(db *gorm.DB, roles *[]entity.Role) error {
	return db.Preload("Permissions").Find(roles).Error
}

func (r *RoleRepository) AssignPermission(db *gorm.DB, role *entity.Role, permission *entity.Permission) error {
	return db.Model(role).Association("Permissions").Append(permission)
}

func (r *RoleRepository) RemovePermission(db *gorm.DB, role *entity.Role, permission *entity.Permission) error {
	return db.Model(role).Association("Permissions").Delete(permission)
}
