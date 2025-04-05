package http

import (
	"comparei-servico-usuario/internal/app"
	"comparei-servico-usuario/internal/infrastructure/http/middleware"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func NewRouter(userService *app.UserService, redisClient *redis.Client) *mux.Router {
	_ = godotenv.Load()

	r := mux.NewRouter()

	r.Use(middleware.APIKeyMiddleware)

	r.HandleFunc("/user", CreateUser).Methods("POST")
	r.HandleFunc("/user/{id}", GetUser).Methods("GET")
	r.HandleFunc("/user/{id}", UpdateUser).Methods("PUT")
	r.HandleFunc("/user/{id}", DeleteUser).Methods("DELETE")

	// Rota de login
	r.HandleFunc("/login", LoginHandler(userService)).Methods("POST")

	// Rota de validação de token
	r.HandleFunc("/validate-token", ValidateTokenHandler()).Methods("GET")

	// Rota para limpar o Redis
	r.HandleFunc("/clear-redis", ClearRedisHandler(redisClient)).Methods("POST")

	return r
}
