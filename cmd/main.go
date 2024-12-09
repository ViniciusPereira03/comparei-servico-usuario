package main

import (
	"comparei-servico-usuario/config"
	"comparei-servico-usuario/internal/api"
	"comparei-servico-usuario/internal/controller"
	"comparei-servico-usuario/internal/db"
	"comparei-servico-usuario/internal/models"
	"comparei-servico-usuario/internal/repository"
	"comparei-servico-usuario/internal/service"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	time.Sleep(time.Second * 30)
	log.Println("Iniciando comparei-serviço-usuario...")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	fmt.Println("DB Host:", cfg.Database.Host)
	fmt.Println("Redis Host:", cfg.Redis.Host)

	dbConnection, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.CloseConnection(dbConnection)

	redisClient, err := db.ConnectRedis(cfg)
	if err != nil {
		log.Fatalf("erro ao conectar ao Redis: %v", err)
	}
	defer redisClient.Close()

	log.Println("Iniciando migrações...")
	errMigrationUser := dbConnection.AutoMigrate(&models.User{})
	if errMigrationUser != nil {
		log.Fatalf("failed to migrate: %v", errMigrationUser)
	}
	log.Println("Migrações realizadas com sucesso!")

	userRepo := repository.NewUserRepository(dbConnection)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	router := gin.Default()

	api.RegisterRoutes(router, userController)

	router.Run(fmt.Sprintf(":%d", cfg.App.Port))
}
