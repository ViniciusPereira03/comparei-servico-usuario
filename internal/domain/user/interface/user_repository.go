package user_interface

import "comparei-servico-usuario/internal/domain/user"

type UserRepository interface {
	CreateUser(user *user.User) (*user.User, error)
	GetUser(id string) (*user.User, error)
	UpdateUser(user *user.User) error
	UpdateLevelUser(id string, level int) error
	DeleteUser(id string) error
	GetUserByUsername(username string) (*user.User, error)
	GetUsers(order string) ([]*user.User, error)
}
