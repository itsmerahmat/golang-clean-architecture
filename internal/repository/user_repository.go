package repository

import (
	"golang-clean-architecture/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
	Log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) *UserRepository {
	return &UserRepository{
		Log: log,
	}
}

func (r *UserRepository) FindByToken(db *gorm.DB, user *entity.User, token string) error {
	return db.Where("token = ?", token).First(user).Error
}

func (r *UserRepository) FindByIdWithRoles(db *gorm.DB, user *entity.User, id string) error {
	return db.Preload("Roles").Preload("Roles.Permissions").Where("id = ?", id).First(user).Error
}

func (r *UserRepository) FindByTokenWithRoles(db *gorm.DB, user *entity.User, token string) error {
	return db.Preload("Roles").Preload("Roles.Permissions").Where("token = ?", token).First(user).Error
}

func (r *UserRepository) AssignRole(db *gorm.DB, user *entity.User, role *entity.Role) error {
	return db.Model(user).Association("Roles").Append(role)
}

func (r *UserRepository) RemoveRole(db *gorm.DB, user *entity.User, role *entity.Role) error {
	return db.Model(user).Association("Roles").Delete(role)
}

