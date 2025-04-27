package user_interface

import "comparei-servico-usuario/internal/domain/user"

type UserRepository interface {
	CreateUser(user *user.User) (*user.User, error)
	GetUser(id int) (*user.User, error)
	UpdateUser(user *user.User) error
	DeleteUser(id int) error
	GetUserByUsername(username string) (*user.User, error)
	GetUsers() ([]*user.User, error)
}
