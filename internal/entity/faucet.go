package entity

import "gorm.io/gorm"

type Faucet struct {
	Address         string
	Amount          int
	TransactionHash string
	UserID          string `gorm:"size:255"`

	gorm.Model
}