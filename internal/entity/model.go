package entity

import (
	"database/sql"
	"time"
)

type DBModel struct {
	ID        uint         `gorm:"primarykey"`
	CreatedAt time.Time    `gorm:"not_null"`
	UpdatedAt time.Time    `gorm:"not_null"`
	DeletedAt sql.NullTime `gorm:"index"`
}
