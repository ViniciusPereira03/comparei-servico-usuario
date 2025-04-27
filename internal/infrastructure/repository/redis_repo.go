package repository

import (
	"comparei-servico-usuario/internal/domain/user"
	user_interface "comparei-servico-usuario/internal/domain/user/interface"
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type UserRepositoryCache struct {
	client *redis.Client
}

func NewUserRepositoryCache(client *redis.Client) user_interface.UserRepositoryCache {
	return &UserRepositoryCache{client: client}
}

func (r *UserRepositoryCache) GetUser(id string) (*user.User, error) {
	ctx := context.Background()
	key := "/user/" + id
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var user user.User
	err = json.Unmarshal([]byte(data), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryCache) SetUser(user *user.User) error {
	ctx := context.Background()
	key := "/user/" + user.ID
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, 5*time.Minute).Err()
}

func (r *UserRepositoryCache) GetUsers() ([]*user.User, error) {
	ctx := context.Background()
	data, err := r.client.Get(ctx, "users").Result()
	if err != nil {
		return nil, err
	}
	var users []*user.User
	err = json.Unmarshal([]byte(data), &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepositoryCache) SetUsers(users []*user.User) error {

	return nil
}

func (r *UserRepositoryCache) DeleteUser(id string) error {

	return nil
}
