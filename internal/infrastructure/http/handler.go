package http

import (
	"comparei-servico-usuario/internal/app"
	"comparei-servico-usuario/internal/domain/user"
	"comparei-servico-usuario/internal/infrastructure/http/dto"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	new_user, err := service.CreateUser(user)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err, "Erro ao cadastrar usuário")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(new_user)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := service.GetUser(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	orderParam := queryValues.Get("order")

	users, err := service.GetUsers(orderParam)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

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

	err := service.UpdateUser(user)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err, "Erro ao atualizar usuário")
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func UpdateUserPhoto(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]

	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Limite de upload (5MB)
	r.ParseMultipartForm(5 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Arquivo não enviado", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validação simples de tipo
	ext := strings.ToLower(filepath.Ext(handler.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		http.Error(w, "Formato de imagem inválido", http.StatusBadRequest)
		return
	}

	// Gera nome único
	filename := primitive.NewObjectID().Hex() + ext

	// Cria diretório se não existir
	uploadDir := "uploads/profile"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		http.Error(w, "Erro ao criar diretório", http.StatusInternalServerError)
		return
	}

	filePath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Erro ao salvar imagem", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Erro ao escrever arquivo", http.StatusInternalServerError)
		return
	}

	// Caminho salvo no banco
	photoPath := fmt.Sprintf("/uploads/profile/%s", filename)

	// Atualiza usuário no Mongo
	user := &user.User{}
	user.ID = oid.Hex()
	user.Photo = photoPath

	err = service.UpdateUser(user)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err, "Erro ao atualizar foto de perfil do usuário")
		return
	}

	response := map[string]string{
		"photo": photoPath,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetUserProfilePhoto(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	filename := params["filename"]

	if filename == "" {
		http.Error(w, "Arquivo inválido", http.StatusBadRequest)
		return
	}

	// Segurança básica: impedir path traversal
	if strings.Contains(filename, "..") {
		http.Error(w, "Acesso inválido", http.StatusForbidden)
		return
	}

	filePath := filepath.Join("uploads", "profile", filename)

	// Verifica se o arquivo existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Imagem não encontrada", http.StatusNotFound)
		return
	}

	// Define Content-Type correto
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	default:
		http.Error(w, "Formato de imagem inválido", http.StatusBadRequest)
		return
	}

	http.ServeFile(w, r, filePath)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := service.DeleteUser(id)
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
		json.NewEncoder(w).Encode(map[string]string{"id": user.ID, "token": tokenString})
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

		// Acessar os dados (claims)
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			id := claims["id"]
			username := claims["username"]
			email := claims["email"]
			validation_hash := claims["validation_hash"]

			response := map[string]interface{}{
				"id":              id,
				"username":        username,
				"email":           email,
				"validation_hash": validation_hash,
			}

			fmt.Println(response)
			// w.Header().Set("Content-Type", "application/json")
			// json.NewEncoder(w).Encode(response)
			// return
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
