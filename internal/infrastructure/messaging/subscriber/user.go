package subscriber

import (
	"comparei-servico-usuario/config"
	"comparei-servico-usuario/internal/app"
	"comparei-servico-usuario/internal/domain/user"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client
var user_service *app.UserService

func init() {
	config.LoadConfig()
	host := fmt.Sprintf("%v:%v", os.Getenv("REDIS_MESSAGING_HOST"), os.Getenv("REDIS_MESSAGING_PORT"))
	rdb = redis.NewClient(&redis.Options{
		Addr: host,
	})
}

// Função para injetar o user_service
func SetUserService(service *app.UserService) {
	user_service = service
}

func SubCreateUser() error {
	ctx := context.Background()

	sub := rdb.Subscribe(ctx, "update_level_user")
	ch := sub.Channel()

	for msg := range ch {
		var user user.User
		err := json.Unmarshal([]byte(msg.Payload), &user)
		if err != nil {
			fmt.Println("[ERRO] Erro ao decodificar payload de mensageria:", err)
			continue
		}

		err = user_service.UpdateLevelUser(user.ID, user.Level)
		if err != nil {
			fmt.Println("[ERRO] Erro ao criar user nos logs:", err)
		}
	}

	return nil
}
