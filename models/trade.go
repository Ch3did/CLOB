package models

import "gorm.io/gorm"

type Trade struct {
	gorm.Model

	Instrument  string  `gorm:"size:20;not null;index"`
	BuyOrderID  uint    `gorm:"not null;index"`
	SellOrderID uint    `gorm:"not null;index"`
	Price       float64 `gorm:"not null"`
	Quantity    float64 `gorm:"not null"`
}
