package repository

import (
	"comparei-servico-usuario/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) domain.UserRepository {
	return &RedisRepository{client: client}
}

func (r *RedisRepository) CreateUser(user *domain.User) error {
	ctx := context.Background()
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, strconv.Itoa(user.ID), data, 0).Err()
}

func (r *RedisRepository) GetUser(id int) (*domain.User, error) {
	ctx := context.Background()
	data, err := r.client.Get(ctx, strconv.Itoa(id)).Result()
	if err != nil {
		return nil, err
	}
	var user domain.User
	err = json.Unmarshal([]byte(data), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *RedisRepository) UpdateUser(user *domain.User) error {
	return r.CreateUser(user)
}

func (r *RedisRepository) DeleteUser(id int) error {
	ctx := context.Background()
	return r.client.Del(ctx, strconv.Itoa(id)).Err()
}

// GetUserByUsername busca um usuário pelo nome de usuário no Redis
func (r *RedisRepository) GetUserByUsername(username string) (*domain.User, error) {
	// A chave no Redis é baseada no ID, então precisamos de uma lógica para buscar por username
	// Isso pode ser uma limitação, pois Redis não é ideal para buscas complexas
	// Aqui, apenas retornamos um erro indicando que a operação não é suportada
	return nil, errors.New("operation not supported")
}
