package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	App struct {
		Port int `yaml:"port"`
	} `yaml:"app"`
	Database struct {
		Host     string `yaml:"host"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Port     int    `yaml:"port"`
		Database string `yaml:"database"`
	} `yaml:"database"`
	Redis struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"redis"`
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return nil, err
	}

	configFile, err := os.ReadFile("config/config.yaml")
	if err != nil {
		return nil, err
	}

	configContent := os.Expand(string(configFile), func(key string) string {
		return os.Getenv(key)
	})

	var config Config
	err = yaml.Unmarshal([]byte(configContent), &config)
	if err != nil {
		return nil, err
	}

	if config.App.Port == 0 {
		appPortStr := os.Getenv("APP_PORT")
		if appPortStr != "" {
			port, err := strconv.Atoi(appPortStr)
			if err != nil {
				return nil, fmt.Errorf("invalid APP_PORT: %v", err)
			}
			config.App.Port = port
		}
	}

	if config.Database.Port == 0 {
		dbPortStr := os.Getenv("DB_PORT")
		if dbPortStr != "" {
			port, err := strconv.Atoi(dbPortStr)
			if err != nil {
				return nil, fmt.Errorf("invalid DB_PORT: %v", err)
			}
			config.Database.Port = port
		}
	}

	if config.Redis.Port == 0 {
		redisPortStr := os.Getenv("REDIS_PORT")
		if redisPortStr != "" {
			port, err := strconv.Atoi(redisPortStr)
			if err != nil {
				return nil, fmt.Errorf("invalid REDIS_PORT: %v", err)
			}
			config.Redis.Port = port
		}
	}

	return &config, nil
}
