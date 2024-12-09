package repository

import (
	"comparei-servico-usuario/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(id uint) (*models.User, error)
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
