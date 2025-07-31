package entity

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type UserRole int

const (
	UserRole_Admin     UserRole = 0 //nolint // underscores used for UserRole
	UserRole_Moderator UserRole = 1 //nolint // underscores used for UserRole
	UserRole_BasicUser UserRole = 2 //nolint // underscores used for UserRole
)

type User struct {
	DBModel

	PlatformID     PlatformID `gorm:"type:tinyint;uniqueIndex:idx_platform_user"`
	PlatformUserID string     `gorm:"type:char(64);uniqueIndex:idx_platform_user"`
	Role           UserRole   `gorm:"type:tinyint"`
}

var UserRoleNameToID = map[string]UserRole{
	"Admin":     UserRole_Admin,
	"Moderator": UserRole_Moderator,
	"BasicUser": UserRole_BasicUser,
}

func AllUserRoles() []UserRole {
	return []UserRole{
		UserRole_Admin,
		UserRole_Moderator,
		UserRole_BasicUser,
	}
}

func (r *UserRole) UnmarshalYAML(value *yaml.Node) error {
	var roleName string
	if err := value.Decode(&roleName); err != nil {
		return err
	}

	role, exists := UserRoleNameToID[roleName]
	if !exists {
		return fmt.Errorf("invalid role name: %s", roleName)
	}

	*r = role

	return nil
}

func (r UserRole) String() string {
	for name, role := range UserRoleNameToID {
		if role == r {
			return name
		}
	}

	return fmt.Sprintf("%d", r)
}
