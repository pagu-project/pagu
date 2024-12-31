package entity

type Role int

const (
	Admin     Role = 0
	Moderator Role = 1
	BasicUser Role = 2
)

type User struct {
	DBModel

	PlatformID PlatformID `gorm:"type:tinyint"`
	UserID     string
	Role       Role `gorm:"type:tinyint"`
}
