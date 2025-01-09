package repository

import "github.com/pagu-project/pagu/internal/entity"

func (db *Database) AddCrowdfundCampaign(campaign *entity.CrowdfundCampaign) error {
	tx := db.gormDB.Create(campaign)
	if tx.Error != nil {
		return WriteError{
			Message: tx.Error.Error(),
		}
	}

	campaign.ID = uint(tx.RowsAffected)

	return nil
}

func (db *Database) GetActiveCrowdfundCampaign() *entity.CrowdfundCampaign {
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

	purchase.ID = uint(tx.RowsAffected)

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
