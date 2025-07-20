package model

import (
	"gorm.io/gorm"
)

type Transfer struct {
	gorm.Model
	OriginAccountID      uint      `gorm:"not null" json:"origin_account_id"`
	DestinationAccountID uint      `gorm:"not null" json:"destination_account_id"`
	Amount               float64   `gorm:"not null" json:"amount"`
	Status               string    `gorm:"not null;default:'PENDING'" json:"status"`
}

type TransferRequest struct {
	OriginAccountID      uint    `json:"origin_account_id" binding:"required"`
	DestinationAccountID uint    `json:"destination_account_id" binding:"required"`
	Amount               float64 `json:"amount" binding:"required"`
}

type TransferResponse struct {
	ID                   uint    `json:"id"`
	OriginAccountID      uint    `json:"origin_account_id"`
	DestinationAccountID uint    `json:"destination_account_id"`
	Amount               float64 `json:"amount"`
	Status               string  `json:"status"`
}

type TransferUpdateRequest struct {
	Status string `json:"status" binding:"required"`
}

type TransferUpdateResponse struct {
	ID     uint   `json:"id"`
	Status string `json:"status"`
}