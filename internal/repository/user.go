package repository

import "github.com/pagu-project/pagu/internal/entity"

type IUser interface {
	AddUser(u *entity.User) error
	HasUser(id string) bool
	GetUserByApp(appID entity.AppID, callerID string) (*entity.User, error)
}

func (db *Database) AddUser(u *entity.User) error {
	tx := db.Create(u)
	if tx.Error != nil {
		return WriteError{
			Message: tx.Error.Error(),
		}
	}

	return nil
}

func (db *Database) HasUser(id string) bool {
	var exists bool

	_ = db.Model(&entity.User{}).
		Select("count(*) > 0").
		Where("id = ?", id).
		Find(&exists).
		Error

	return exists
}

func (db *Database) GetUserByApp(appID entity.AppID, callerID string) (*entity.User, error) {
	var user *entity.User
	tx := db.Model(&entity.User{}).
		Where("application_id = ?", appID).
		Where("caller_id = ?", callerID).
		First(&user)

	if tx.Error != nil {
		return nil, ReadError{
			Message: tx.Error.Error(),
		}
	}

	return user, nil
}
