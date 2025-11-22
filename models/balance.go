package models

import "gorm.io/gorm"

type Balance struct {
	gorm.Model

	AccountID uint    `gorm:"not null;index"`
	Asset     string  `gorm:"size:20;not null;index"`
	Amount    float64 `gorm:"not null"`
}
