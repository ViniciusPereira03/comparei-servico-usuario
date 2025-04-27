package user_interface

import "comparei-servico-usuario/internal/domain/user"

type UserRepositoryCache interface {
	GetUsers() ([]*user.User, error)
	SetUsers(users []*user.User) error

	GetUser(id string) (*user.User, error)
	SetUser(user *user.User) error

	DeleteUser(id string) error
}
