package repository

import "github.com/pagu-project/pagu/internal/entity"

func (db *Database) GetZealyUser(id string) (*entity.ZealyUser, error) {
	var user *entity.ZealyUser
	tx := db.gormDB.Model(&entity.ZealyUser{}).First(&user, "discord_id = ?", id)
	if tx.Error != nil {
		return &entity.ZealyUser{}, ReadError{
			Message: tx.Error.Error(),
		}
	}

	return user, nil
}

func (db *Database) AddZealyUser(u *entity.ZealyUser) error {
	tx := db.gormDB.Create(u)
	if tx.Error != nil {
		return WriteError{
			Message: tx.Error.Error(),
		}
	}

	return nil
}

func (db *Database) UpdateZealyUser(id, txHash string) error {
	tx := db.gormDB.Model(&entity.ZealyUser{
		DiscordID: id,
	}).Where("discord_id = ?", id).Update("tx_hash", txHash)
	if tx.Error != nil {
		return WriteError{
			Message: tx.Error.Error(),
		}
	}

	return nil
}

func (db *Database) GetAllZealyUser() ([]*entity.ZealyUser, error) {
	var users []*entity.ZealyUser
	tx := db.gormDB.Find(&users)
	if tx.Error != nil {
		return nil, ReadError{
			Message: tx.Error.Error(),
		}
	}

	return users, nil
}
