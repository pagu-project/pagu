package entity

import (
	"time"

	"github.com/pagu-project/pagu/pkg/amount"
)

type PhoenixFaucet struct {
	DBModel

	UserID  uint          `gorm:"type:bigint"`
	Address string        `gorm:"type:char(43)"`
	Amount  amount.Amount `gorm:"column:amount"`
	TxHash  string        `gorm:"type:char(64);unique;not null"`
}

func (*PhoenixFaucet) TableName() string {
	// TODO: rename me to "faucets" (just remove this function is enough).
	return "phoenix_faucet"
}

func (f *PhoenixFaucet) ElapsedTime() time.Duration {
	return time.Since(f.CreatedAt)
}
