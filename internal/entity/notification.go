package entity

import (
	"github.com/pagu-project/pagu/pkg/notification"
	"gorm.io/datatypes"
)

type NotificationStatus int

const (
	NotificationStatusPending = iota
	NotificationStatusDone
	NotificationStatusFail
)

type Notification struct {
	DBModel

	Type      notification.NotificationType `gorm:"type:tinyint"`
	Recipient string                        `gorm:"size:255"`
	Data      datatypes.JSON
	Status    NotificationStatus `gorm:"type:tinyint"`
}

type VoucherNotificationData struct {
	Code      string  `json:"code"`
	Amount    float64 `json:"amount"`
	Recipient string  `json:"recipient"`
}
