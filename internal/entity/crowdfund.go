package entity

import "github.com/pagu-project/pagu/pkg/amount"

type Package struct {
	Name      string        `json:"name"`
	USDAmount int           `json:"usd_amount"`
	PACAmount amount.Amount `gorm:"type:double" json:"pac_amount"`
}

type CrowdfundCampaign struct {
	DBModel

	Title     string    `gorm:"type:char(128);not null"`
	Desc      string    `gorm:"type:text;not null"`
	Packages  []Package `gorm:"serializer:json"`
	CreatorID uint      // TODO: define foreign key here
	Active    bool
}

type CrowdfundPurchase struct {
	DBModel

	UserID    uint
	InvoiceID string
	USDAmount int
	PACAmount amount.Amount `gorm:"type:double"`
	TxHash    string        `gorm:"type:char(64);default:null"`
	Recipient string        `gorm:"type:char(42);default:null"`
}

func (p *CrowdfundPurchase) IsClaimed() bool {
	return p.TxHash != ""
}
