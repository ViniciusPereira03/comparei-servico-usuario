package http

import (
	"comparei-servico-usuario/internal/app"
	"comparei-servico-usuario/internal/domain/user"
	"comparei-servico-usuario/internal/infrastructure/http/dto"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

var service *app.UserService

func InitHandlers(userService *app.UserService) {
	service = userService
}

func sendErrorResponse(w http.ResponseWriter, statusCode int, err error, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error":    err.Error(),
		"mensagem": message,
	})
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var userDTO dto.CreateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err, "JSON inválido")
		return
	}

	validate := validator.New()
	if err := validate.Struct(userDTO); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err, "Dados inválidos")
		return
	}

	user := userDTO.ParseToUser()
	err := service.CreateUser(user)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err, "Erro ao cadastrar usuário")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	user, err := service.GetUser(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := service.GetUsers()
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err, "Invalid user ID")
		return
	}

	var userDTO dto.UpdateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err, "JSON inválido")
		return
	}

	validate := validator.New()
	if err := validate.Struct(userDTO); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err, "Dados inválidos")
		return
	}

	user := userDTO.ParseToUser()
	user.ID = id

	err = service.UpdateUser(user)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err, "Erro ao atualizar usuário")
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	err = service.DeleteUser(id)
	if err != nil {
		http.Error(w, "Erro ao deletar usuário", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("usuário deletado com sucesso"))
}

// LoginHandler lida com a autenticação do usuário e geração de token JWT
func LoginHandler(userService *app.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds user.User
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Verificar credenciais do usuário
		user, err := userService.Authenticate(creds.Username, creds.Password)
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// Criar token JWT
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":              user.ID,
			"username":        user.Username,
			"email":           user.Email,
			"validation_hash": os.Getenv("VALIDATION_HASH"),
		})

		tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			log.Println("Error signing token: ", err)
			http.Error(w, "Could not generate token", http.StatusInternalServerError)
			return
		}

		// Retornar token JWT
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
	}
}

// ValidateTokenHandler verifica a validade de um token JWT
func ValidateTokenHandler() http.HandlerFunc {
	secret := os.Getenv("JWT_SECRET")
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// Remover o prefixo "Bearer " se existir
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// Verificar o token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Token é válido
		w.WriteHeader(http.StatusOK)
	}
}

func ClearRedisHandler(redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := redisClient.FlushAll(r.Context()).Err()
		if err != nil {
			http.Error(w, "Erro ao limpar o Redis", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Redis limpo com sucesso"))
	}
}
