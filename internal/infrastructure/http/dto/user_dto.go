package dto

import "comparei-servico-usuario/internal/domain/user"

type CreateUserDTO struct {
	Name     string `json:"name" validate:"required,min=2,max=255"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// Método para converter CreateUserDTO para user.User
func (dto *CreateUserDTO) ParseToUser() *user.User {
	return &user.User{
		Name:     dto.Name,
		Username: dto.Username,
		Email:    dto.Email,
		Password: dto.Password,
	}
}

type UpdateUserDTO struct {
	Name        *string `json:"name"`
	Username    *string `json:"username"`
	Email       *string `json:"email"`
	Status      *int    `json:"status"`
	Photo       *string `json:"photo"`
	RayDistance *int    `json:"ray_distance"`
	Level       *int    `json:"level"`
}

// Método para converter CreateUserDTO para user.User
func (dto *UpdateUserDTO) ParseToUser() *user.User {
	return &user.User{
		Name:        derefString(dto.Name),
		Username:    derefString(dto.Username),
		Email:       derefString(dto.Email),
		Photo:       derefString(dto.Photo),
		Status:      derefInt(dto.Status),
		RayDistance: derefInt(dto.RayDistance),
		Level:       derefInt(dto.Level),
	}
}

func derefString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func derefInt(i *int) int {
	if i != nil {
		return *i
	}
	return 0
}
