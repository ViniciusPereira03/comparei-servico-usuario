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

	// 🔓 ROTAS PÚBLICAS (sem API key)
	r.HandleFunc("/uploads/profile/{filename}", GetUserProfilePhoto).Methods("GET")

	// 🔒 ROTAS PROTEGIDAS
	api := r.PathPrefix("/").Subrouter()
	api.Use(middleware.APIKeyMiddleware)

	api.HandleFunc("/user", CreateUser).Methods("POST")
	api.HandleFunc("/users", GetUsers).Methods("GET")
	api.HandleFunc("/user/{id}", GetUser).Methods("GET")
	api.HandleFunc("/user/{id}", UpdateUser).Methods("PUT")
	api.HandleFunc("/user/{id}", DeleteUser).Methods("DELETE")
	api.HandleFunc("/user/{id}/photo", UpdateUserPhoto).Methods("POST")

	api.HandleFunc("/login", LoginHandler(userService)).Methods("POST")
	api.HandleFunc("/validate-token", ValidateTokenHandler()).Methods("GET")
	api.HandleFunc("/clear-redis", ClearRedisHandler(redisClient)).Methods("POST")

	return r
}
