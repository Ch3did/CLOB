package models

import "gorm.io/gorm"

type Account struct {
	gorm.Model        // ID, CreatedAt, UpdatedAt, DeletedAt
	Name       string `gorm:"size:255;not null"`
}
