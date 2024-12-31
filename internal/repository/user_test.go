package repository

import (
	"testing"

	"github.com/pagu-project/pagu/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserByID(t *testing.T) {
	db, err := NewDB("sqlite:file::memory:")
	require.NoError(t, err)

	user := entity.User{
		PlatformID:     entity.PlatformIDDiscord,
		PlatformUserID: "1234",
		Role:           entity.BasicUser,
	}
	err = db.AddUser(&user)
	assert.NoError(t, err)

	foundUser, err := db.GetUserByPlatformID(user.PlatformID, user.PlatformUserID)
	require.NoError(t, err)

	assert.Equal(t, foundUser.PlatformID, user.PlatformID)
	assert.Equal(t, foundUser.PlatformUserID, user.PlatformUserID)
	assert.Equal(t, foundUser.Role, user.Role)

	hasUser := db.HasUser(foundUser.DBModel.ID)
	assert.True(t, hasUser)
}
