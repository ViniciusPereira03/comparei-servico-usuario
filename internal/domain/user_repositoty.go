package domain

type UserRepository interface {
	CreateUser(user *User) error
	GetUsers() (*[]User, error)
	GetUser(id int) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id int) error
	GetUserByUsername(username string) (*User, error)
}
