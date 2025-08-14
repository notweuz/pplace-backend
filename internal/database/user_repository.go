package database

import (
	"errors"
	"gorm.io/gorm"
	"pplace_backend/internal/model"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(database *gorm.DB) UserRepository {
	return UserRepository{db: database}
}

func (ur *UserRepository) Create(user *model.User) (*model.User, error) {
	result := ur.db.Create(user)
	return user, result.Error
}

func (ur *UserRepository) GetById(id uint) (*model.User, error) {
	var user model.User
	result := ur.db.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, result.Error
}

func (ur *UserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	result := ur.db.Where("username = ?", username).Limit(1).Find(&user)

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &user, result.Error
}
