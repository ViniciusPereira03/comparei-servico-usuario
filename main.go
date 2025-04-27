package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"comparei-servico-usuario/config"
	"comparei-servico-usuario/internal/app"
	customHTTP "comparei-servico-usuario/internal/infrastructure/http"
	"comparei-servico-usuario/internal/infrastructure/messaging/subscriber"
	"comparei-servico-usuario/internal/infrastructure/repository"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// carrega .env, etc.
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Erro ao carregar configurações:", err)
	}

	// --- Redis de mensageria ---
	redisMessaging := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_MESSAGING_HOST") + ":" + os.Getenv("REDIS_MESSAGING_PORT"),
	})
	ctx := context.Background()
	if _, err := redisMessaging.Ping(ctx).Result(); err != nil {
		log.Fatal("Não foi possível conectar ao Redis de mensageria:", err)
	}

	// --- MongoDB ---
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal("Erro ao criar cliente MongoDB:", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := mongoClient.Connect(ctx); err != nil {
		log.Fatal("Erro ao conectar no MongoDB:", err)
	}
	// opcional: ping para certificar
	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatal("Ping no MongoDB falhou:", err)
	}

	// instancia repositório Mongo (hexagonal)
	mongoRepo := repository.NewMongoRepository(
		mongoClient,
		os.Getenv("MONGO_DB_NAME"),
		os.Getenv("MONGO_COLLECTION"),
	)

	// --- Redis de cache (outra instância) ---
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
	})
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		log.Fatal("Não foi possível conectar ao Redis de cache:", err)
	}
	cacheRepo := repository.NewUserRepositoryCache(redisClient)

	// --- Service ---
	userService := app.NewUserService(mongoRepo, cacheRepo)

	// subscriber de eventos
	subscriber.SetUserService(userService)
	go func() {
		fmt.Println("Inicializando subscriber...")
		if err := subscriber.SubCreateUser(); err != nil {
			log.Println("Erro no subscriber:", err)
		}
	}()

	// HTTP
	customHTTP.InitHandlers(userService)
	router := customHTTP.NewRouter(userService, redisClient)

	log.Println("Servidor iniciado na porta", os.Getenv("PORT"))
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), router); err != nil {
		log.Fatal(err)
	}
}
