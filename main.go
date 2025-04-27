package main

import (
	"comparei-servico-usuario/config"
	"comparei-servico-usuario/internal/app"
	customHTTP "comparei-servico-usuario/internal/infrastructure/http"
	"comparei-servico-usuario/internal/infrastructure/repository"
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Erro ao carregar configurações")
	}

	// Testar conexão com Redis de mensageria
	redisMessaging := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_MESSAGING_HOST") + ":" + os.Getenv("REDIS_MESSAGING_PORT"),
	})
	ctx := context.Background()
	_, err := redisMessaging.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Não foi possível conectar ao Redis de mensageria:", err)
	}

	// Configuração da conexão com o MySQL usando variáveis de ambiente
	dsn := os.Getenv("MYSQL_USER") + ":" + os.Getenv("MYSQL_PASSWORD") + "@tcp(" + os.Getenv("MYSQL_HOST") + ")/" + os.Getenv("MYSQL_DB")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Verificar a conexão com o MySQL
	if err := db.Ping(); err != nil {
		log.Fatal("Não foi possível conectar ao MySQL: ", err)
	}

	// Configuração do cliente Redis usando variáveis de ambiente
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
	})

	// Verificar a conexão com o Redis
	if _, err := redisClient.Ping(redisClient.Context()).Result(); err != nil {
		log.Fatal("Não foi possível conectar ao Redis: ", err)
	}

	mysqlRepo := repository.NewMySQLRepository(db)
	redisRepo := repository.NewUserRepositoryCache(redisClient)

	if mysqlRepo == nil {
		log.Fatal("mysqlRepo está nil")
	}
	if redisRepo == nil {
		log.Fatal("redisRepo está nil")
	}

	userService := app.NewUserService(mysqlRepo, redisRepo)
	customHTTP.InitHandlers(userService)

	router := customHTTP.NewRouter(userService, redisClient)

	log.Println("Servidor iniciado na porta " + os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), router)
}
