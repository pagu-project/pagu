package entity

type Role int

const (
	Admin     Role = 0
	Moderator Role = 1
	BasicUser Role = 2
)

type User struct {
	DBModel

	PlatformID     PlatformID `gorm:"type:tinyint;uniqueIndex:idx_platform_user"`
	PlatformUserID string     `gorm:"type:char(64);uniqueIndex:idx_platform_user"`
	Role           Role       `gorm:"type:tinyint"`
}
