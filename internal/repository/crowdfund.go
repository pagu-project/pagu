package repository

import "github.com/pagu-project/pagu/internal/entity"

func (db *Database) AddCrowdfundCampaign(campaign *entity.CrowdfundCampaign) error {
	tx := db.gormDB.Create(campaign)
	if tx.Error != nil {
		return WriteError{
			Message: tx.Error.Error(),
		}
	}

	return nil
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
