package repository

import "github.com/pagu-project/Pagu/internal/entity"

type IZealy interface {
	GetZealyUser(id string) (*entity.ZealyUser, error)
	AddZealyUser(u *entity.ZealyUser) error
	UpdateZealyUser(id string, txHash string) error
	GetAllZealyUser() ([]*entity.ZealyUser, error)
}

func (db *Database) GetZealyUser(id string) (*entity.ZealyUser, error) {
	var user *entity.ZealyUser
	tx := db.Model(&entity.ZealyUser{}).First(&user, "discord_id = ?", id)
	if tx.Error != nil {
		return &entity.ZealyUser{}, ReadError{
			Message: tx.Error.Error(),
		}
	}

	return user, nil
}

func (db *Database) AddZealyUser(u *entity.ZealyUser) error {
	tx := db.Create(u)
	if tx.Error != nil {
		return WriteError{
			Message: tx.Error.Error(),
		}
	}

	return nil
}

func (db *Database) UpdateZealyUser(id, txHash string) error {
	tx := db.Model(&entity.ZealyUser{
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
	tx := db.Find(&users)
	if tx.Error != nil {
		return nil, ReadError{
			Message: tx.Error.Error(),
		}
	}

	return users, nil
}
