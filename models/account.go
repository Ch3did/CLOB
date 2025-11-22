package models

import "gorm.io/gorm"

type Account struct {
	gorm.Model

	FirstName    string `gorm:"size:255;not null"`
	LastName     string `gorm:"size:255;not null"`
	Email        string `gorm:"size:255;uniqueIndex;not null"`
	PasswordHash string `gorm:"size:255;not null"`
}
