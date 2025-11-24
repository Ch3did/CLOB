package models

import "gorm.io/gorm"

type Side string

const (
	SideBuy  Side = "BUY"
	SideSell Side = "SELL"
)

type Order struct {
	gorm.Model

	AccountID  uint    `gorm:"not null;index"`
	Instrument string  `gorm:"size:20;not null;index"`
	Side       Side    `gorm:"type:varchar(4);not null"`
	Price      float64 `gorm:"not null"`
	Quantity   float64 `gorm:"not null"`
	Status     string  `gorm:"size:20;default:'open'"`
}
