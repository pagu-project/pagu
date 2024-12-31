package repository

import (
	"time"

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

func (db *Database) CanGetFaucet(user *entity.User) bool {
	var lastFaucet entity.PhoenixFaucet
	err := db.gormDB.Model(&entity.PhoenixFaucet{}).Where("user_id = ?", user.ID).Order("id DESC").First(&lastFaucet).Error
	if err != nil {
		return true
	}

	if lastFaucet.ElapsedTime() > 24*time.Hour {
		return true
	}

	return false
}
