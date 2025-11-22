package database

import (
	"fmt"

	"CLOB/internal/config"
	"CLOB/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// AutoMigrate simples
	if err := db.AutoMigrate(
		&models.Account{},
		&models.Balance{},
		&models.Order{},
		&models.Trade{},
	); err != nil {
		return nil, err
	}

	return db, nil
}
