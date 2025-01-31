package repository

import (
	"github.com/pagu-project/pagu/internal/entity"
)

func (db *Database) AddFaucet(f *entity.PhoenixFaucet) error {
	tx := db.gormDB.Create(f)
	if tx.Error != nil {
		return WriteError{
			Message: tx.Error.Error(),
		}
	}

	return nil
}

func (db *Database) GetLastFaucet(user *entity.User) (*entity.PhoenixFaucet, error) {
	var lastFaucet *entity.PhoenixFaucet
	tx := db.gormDB.Model(&entity.PhoenixFaucet{}).Where("user_id = ?", user.ID).Order("id DESC").First(&lastFaucet)

	if tx.Error != nil {
		return nil, ReadError{
			Message: tx.Error.Error(),
		}
	}

	return lastFaucet, nil
}
