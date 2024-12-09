package service

import (
	"comparei-servico-usuario/internal/models"
	"comparei-servico-usuario/internal/repository"
)

type UserService interface {
	GetUserByID(id uint) (*models.User, error)
}

type userService struct {
	UserRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{UserRepo: userRepo}
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {
	return s.UserRepo.GetUserByID(id)
}
