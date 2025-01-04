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
