package publisher

import (
	"comparei-servico-usuario/config"
	"comparei-servico-usuario/internal/domain/user"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func init() {
	config.LoadConfig()
	host := fmt.Sprintf("%v:%v", os.Getenv("REDIS_MESSAGING_HOST"), os.Getenv("REDIS_MESSAGING_PORT"))
	rdb = redis.NewClient(&redis.Options{
		Addr: host,
	})
}

func PubCreateUser(u *user.User) error {
	ctx := context.Background()

	payload, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao codificar payload: %v", err)
	}
	_, err = rdb.Publish(ctx, "user_created", string(payload)).Result()
	if err != nil {
		return fmt.Errorf("erro ao publicar mensagem no Redis: %v", err)
	}

	return nil
}

func PubModifyUser(u *user.User) error {
	ctx := context.Background()

	payload, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao codificar payload: %v", err)
	}
	_, err = rdb.Publish(ctx, "user_modified", string(payload)).Result()
	if err != nil {
		return fmt.Errorf("erro ao publicar mensagem no Redis: %v", err)
	}

	return nil
}
