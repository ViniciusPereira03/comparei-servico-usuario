package app

import (
	"comparei-servico-usuario/internal/domain"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	mysqlRepo domain.UserRepository
	redisRepo domain.UserRepository
}

func NewUserService(mysqlRepo domain.UserRepository, redisRepo domain.UserRepository) *UserService {
	return &UserService{mysqlRepo: mysqlRepo, redisRepo: redisRepo}
}

func (s *UserService) CreateUser(user *domain.User) error {
	// Criptografar a senha do usuário
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("erro ao criptografar senha")
	}
	user.Password = string(hashedPassword)

	fmt.Println("SENHA: ", user.Password)

	err = s.mysqlRepo.CreateUser(user)
	if err == nil {
		s.redisRepo.CreateUser(user)
	}
	return err
}

func (s *UserService) GetUser(id int) (*domain.User, error) {
	user, err := s.redisRepo.GetUser(id)
	if err == nil {
		return user, nil
	}

	user, err = s.mysqlRepo.GetUser(id)
	if err != nil {
		return nil, err
	}

	s.redisRepo.CreateUser(user)
	return user, nil
}

func (s *UserService) UpdateUser(user *domain.User) error {
	err := s.mysqlRepo.UpdateUser(user)
	if err == nil {
		s.redisRepo.CreateUser(user)
	}
	return err
}

func (s *UserService) DeleteUser(id int) error {
	err := s.mysqlRepo.DeleteUser(id)
	if err == nil {
		er := s.redisRepo.DeleteUser(id) // Remove do cache
		if er != nil {
			return er
		}
	}
	return err
}

// Authenticate verifica as credenciais do usuário e retorna o usuário se forem válidas
func (s *UserService) Authenticate(username, password string) (*domain.User, error) {
	user, err := s.mysqlRepo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	// Comparação segura de senhas
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
