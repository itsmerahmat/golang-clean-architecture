package repository

import (
	"golang-clean-architecture/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PermissionRepository struct {
	Repository[entity.Permission]
	Log *logrus.Logger
}

func NewPermissionRepository(log *logrus.Logger) *PermissionRepository {
	return &PermissionRepository{
		Log: log,
	}
}

func (r *PermissionRepository) FindByName(db *gorm.DB, permission *entity.Permission, name string) error {
	return db.Where("name = ?", name).First(permission).Error
}

func (r *PermissionRepository) FindByResourceAndAction(db *gorm.DB, permission *entity.Permission, resource string, action string) error {
	return db.Where("resource = ? AND action = ?", resource, action).First(permission).Error
}
