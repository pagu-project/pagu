package repository

import (
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/log"
)

func (db *Database) AddCrowdfundCampaign(campaign *entity.CrowdfundCampaign) error {
	tx := db.gormDB.Create(campaign)
	if tx.Error != nil {
		return WriteError{
			Message: tx.Error.Error(),
		}
	}

	return nil
}

func (db *Database) GetCrowdfundActiveCampaign() *entity.CrowdfundCampaign {
	var campaign *entity.CrowdfundCampaign
	tx := db.gormDB.Model(&entity.CrowdfundCampaign{}).
		Where("active = true").
		First(&campaign)

	if tx.Error != nil {
		return nil
	}

	return campaign
}

func (db *Database) GetCrowdfundCampaign(campaignID uint) (*entity.CrowdfundCampaign, error) {
	var campaign *entity.CrowdfundCampaign
	tx := db.gormDB.Model(&entity.CrowdfundCampaign{}).
		Where("id = ?", campaignID).
		First(&campaign)

	if tx.Error != nil {
		return nil, ReadError{
			Message: tx.Error.Error(),
		}
	}

	return campaign, nil
}

func (db *Database) AddCrowdfundPurchase(purchase *entity.CrowdfundPurchase) error {
	tx := db.gormDB.Create(purchase)
	if tx.Error != nil {
		return WriteError{
			Message: tx.Error.Error(),
		}
	}

	return nil
}

func (db *Database) UpdateCrowdfundPurchase(purchase *entity.CrowdfundPurchase) error {
	tx := db.gormDB.Save(purchase)
	if tx.Error != nil {
		return WriteError{
			Message: tx.Error.Error(),
		}
	}

	return nil
}

func (db *Database) GetCrowdfundPurchases(userID uint) ([]*entity.CrowdfundPurchase, error) {
	var purchases []*entity.CrowdfundPurchase
	tx := db.gormDB.Model(&entity.CrowdfundPurchase{}).
		Where("user_id = ?", userID).
		Find(&purchases)

	if tx.Error != nil {
		return nil, ReadError{
			Message: tx.Error.Error(),
		}
	}

	return purchases, nil
}

func (db *Database) GetTotalPurchasedPackages() int64 {
	var count int64
	tx := db.gormDB.Model(&entity.CrowdfundPurchase{}).
		Where("tx_hash <> ''").
		Count(&count)

	if tx.Error != nil {
		log.Error("failed to get total purchased packages", "error", tx.Error.Error())

		return -1
	}

	return count
}

// GetTotalCrowdfundedAmount returns the total crowdfunded amount in USD with successfully claimed transactions.
// It does not include any paid but unclaimed purchases.
func (db *Database) GetTotalCrowdfundedAmount() int64 {
	var totalUSD int64
	tx := db.gormDB.Model(&entity.CrowdfundPurchase{}).
		Where("tx_hash <> ''").
		Select("SUM(usd_amount)").
		Scan(&totalUSD)
	if tx.Error != nil {
		log.Error("failed to get total crowdfunded amount", "error", tx.Error.Error())

		return -1
	}

	return totalUSD
}
