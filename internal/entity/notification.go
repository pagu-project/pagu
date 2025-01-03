package entity

import (
	"github.com/pagu-project/pagu/pkg/notification"
)

type NotificationStatus int

const (
	NotificationStatusPending = iota
	NotificationStatusDone
	NotificationStatusFail
)

type VoucherNotificationData struct {
	Code      string  `json:"code"`
	Amount    float64 `json:"amount"`
	Recipient string  `json:"recipient"`
}

type Notification struct {
	DBModel

	Type      notification.NotificationType `gorm:"type:tinyint"`
	Recipient string                        `gorm:"size:255"`
	Data      VoucherNotificationData       `gorm:"serializer:json"`
	Status    NotificationStatus            `gorm:"type:tinyint"`
}
