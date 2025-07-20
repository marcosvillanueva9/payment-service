package model

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	Name    string  `gorm:"not null" json:"name"`
	Balance float64 `gorm:"not null" json:"balance"`
}
