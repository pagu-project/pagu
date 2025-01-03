package entity

import (
	"github.com/pagu-project/pagu/pkg/amount"
)

type Voucher struct {
	DBModel

	Creator     uint
	Code        string        `gorm:"type:char(8);unique"`
	Amount      amount.Amount `gorm:"column:amount"`
	Desc        string
	Email       string
	Recipient   string
	ValidMonths uint8
	TxHash      string `gorm:"type:char(64);default:null"`
	ClaimedBy   uint
}

// TODO: rename me to "voucher" (just remove this function is enough)
func (Voucher) TableName() string {
	return "voucher"
}

func (v *Voucher) IsClaimed() bool {
	return v.TxHash != ""
}
