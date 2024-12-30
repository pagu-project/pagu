package repository

import (
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/notification"
)

type INotification interface {
	AddNotification(v *entity.Notification) error
	GetPendingMailNotification() (*entity.Notification, error)
	UpdateNotificationStatus(id uint, status entity.NotificationStatus) error
}

func (db *Database) AddNotification(v *entity.Notification) error {
	tx := db.Create(v)
	if tx.Error != nil {
		return WriteError{
			Message: tx.Error.Error(),
		}
	}

	return nil
}

func (db *Database) GetPendingMailNotification() (*entity.Notification, error) {
	var notif *entity.Notification
	tx := db.Model(&entity.Notification{}).
		Where("status = ?", entity.NotificationStatusPending).
		Where("type = ?", notification.NotificationTypeMail).
		First(&notif)

	if tx.Error != nil {
		return nil, ReadError{
			Message: tx.Error.Error(),
		}
	}

	return notif, nil
}

func (db *Database) UpdateNotificationStatus(id uint, status entity.NotificationStatus) error {
	tx := db.Model(&entity.Notification{}).Where("id = ?", id).Update("status", status)
	if tx.Error != nil {
		return WriteError{
			Message: tx.Error.Error(),
		}
	}

	return nil
}
