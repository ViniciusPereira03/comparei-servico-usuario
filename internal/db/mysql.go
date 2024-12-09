package db

import (
	"comparei-servico-usuario/config"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	// Formatar a string de conexão
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Database,
	)

	// Conectar ao banco de dados
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao MySQL: %v", err)
	}

	// Retornar a conexão
	return db, nil
}

func CloseConnection(DB *gorm.DB) {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("failed to close the database connection: %v", err)
	}
	sqlDB.Close()
}
