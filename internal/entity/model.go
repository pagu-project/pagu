package entity

import (
	"database/sql"
	"time"
)

type DBModel struct {
	ID        uint         `gorm:"primarykey"`
	CreatedAt time.Time    `gorm:"not null"`
	UpdatedAt time.Time    `gorm:"not null"`
	DeletedAt sql.NullTime `gorm:"null"`
}
