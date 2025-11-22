package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"central_limit_order_book/models"
)

func main() {
	db, err := gorm.Open(sqlite.Open("clob.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(
		&models.Account{},
		&models.Balance{},
		&models.Order{},
		&models.Trade{},
	)
	if err != nil {
		log.Fatal(err)
	}
}
