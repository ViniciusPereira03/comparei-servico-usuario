package dto

import "comparei-servico-usuario/internal/domain/user"

type CreateUserDTO struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3"`
	Status   int    `json:"status" validate:"required"`
}

// Método para converter CreateUserDTO para user.User
func (dto *CreateUserDTO) ParseToUser() *user.User {
	return &user.User{
		Name:     dto.Name,
		Username: dto.Username,
		Email:    dto.Email,
		Password: dto.Password,
		Status:   dto.Status,
	}
}

type UpdateUserDTO struct {
	Name        string `json:"name" validate:"required"`
	Username    string `json:"username" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Status      int    `json:"status" validate:"required"`
	Photo       string `json:"photo"`
	RayDistance int    `json:"ray_distance"`
	Level       int    `json:"level"`
}

// Método para converter CreateUserDTO para user.User
func (dto *UpdateUserDTO) ParseToUser() *user.User {
	return &user.User{
		Name:        dto.Name,
		Username:    dto.Username,
		Email:       dto.Email,
		Status:      dto.Status,
		Photo:       dto.Photo,
		RayDistance: dto.RayDistance,
		Level:       dto.Level,
	}
}
