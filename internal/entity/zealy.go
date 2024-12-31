package entity

import (
	"github.com/pagu-project/pagu/pkg/amount"
)

type ZealyUser struct {
	DBModel

	Amount    amount.Amount `gorm:"column:amount"`
	DiscordID string        `gorm:"column:discord_id"`
	TxHash    string        `gorm:"type:char(64);unique;default:null"`
}

func (z *ZealyUser) IsClaimed() bool {
	return z.TxHash != ""
}
