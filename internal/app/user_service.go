package app

import (
	"comparei-servico-usuario/internal/domain/user"
	user_interface "comparei-servico-usuario/internal/domain/user/interface"
	"comparei-servico-usuario/internal/infrastructure/messaging/publisher"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	mongoRepo user_interface.UserRepository
	redisRepo user_interface.UserRepositoryCache
}

func NewUserService(mongoRepo user_interface.UserRepository, redisRepo user_interface.UserRepositoryCache) *UserService {
	return &UserService{mongoRepo: mongoRepo, redisRepo: redisRepo}
}

func (s *UserService) CreateUser(user *user.User) (*user.User, error) {
	// Criptografar a senha do usuário
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("erro ao criptografar senha")
	}
	user.Password = string(hashedPassword)

	new_user, err := s.mongoRepo.CreateUser(user)
	if err == nil {
		s.redisRepo.SetUser(new_user)
		err_pub := publisher.PubCreateUser(new_user)
		if err_pub != nil {
			log.Println("[ERRO PUB] ", err_pub)
		}
	}

	return new_user, err
}

func (s *UserService) GetUser(id string) (*user.User, error) {
	user, err := s.redisRepo.GetUser(id)
	if err == nil {
		return user, nil
	}

	user, err = s.mongoRepo.GetUser(id)
	if err != nil {
		return nil, err
	}

	s.redisRepo.SetUser(user)
	return user, nil
}

func (s *UserService) GetUsers(order string) ([]*user.User, error) {
	users, err := s.redisRepo.GetUsers()
	if err == nil {
		userPtrs := make([]*user.User, len(users))
		for i := range users {
			userPtrs[i] = users[i]
		}
		return userPtrs, nil
	}

	users, err = s.mongoRepo.GetUsers(order)
	if err != nil {
		return nil, err
	}

	userPtrs := make([]*user.User, len(users))
	for i := range users {
		userPtrs[i] = users[i]
	}

	s.redisRepo.SetUsers(userPtrs)
	return userPtrs, nil
}

func (s *UserService) UpdateUser(user *user.User) error {
	err := s.mongoRepo.UpdateUser(user)
	if err == nil {
		user, err = s.mongoRepo.GetUser(user.ID)
		s.redisRepo.SetUser(user)

		err_pub := publisher.PubModifyUser(user)
		if err_pub != nil {
			log.Println("[ERRO PUB] ", err_pub)
		}
	}
	return err
}

func (s *UserService) UpdateLevelUser(user_id string, level int) error {
	err := s.mongoRepo.UpdateLevelUser(user_id, level)
	if err == nil {
		user, err := s.mongoRepo.GetUser(user_id)
		if err != nil {
			return err
		}
		s.redisRepo.SetUser(user)
		err_pub := publisher.PubModifyUser(user)
		if err_pub != nil {
			log.Println("[ERRO PUB] ", err_pub)
		}
	}
	return err
}

func (s *UserService) DeleteUser(id string) error {
	err := s.mongoRepo.DeleteUser(id)
	if err == nil {
		er := s.redisRepo.DeleteUser(id) // Remove do cache
		if er != nil {
			return er
		}
	}
	return err
}

// Authenticate verifica as credenciais do usuário e retorna o usuário se forem válidas
func (s *UserService) Authenticate(username, password string) (*user.User, error) {
	user, err := s.mongoRepo.GetUserByUsername(username)
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
