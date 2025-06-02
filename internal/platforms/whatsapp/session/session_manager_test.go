package session

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOpenAndExistSession(t *testing.T) {
	manager := NewSessionManager(context.Background(), time.Minute, time.Minute)
	userID := "user1"

	// Check that session does not exist
	assert.False(t, manager.ExistSession(userID), "Session should not exist initially")

	// Open session and check existence
	_ = manager.OpenSession(userID)
	assert.True(t, manager.ExistSession(userID), "Session should exist after opening")

	// Check another user that was never added
	otherUser := "nonexistent"
	assert.False(t, manager.ExistSession(otherUser), "Session should not exist for a different user")
}

func TestCloseSession(t *testing.T) {
	manager := NewSessionManager(context.Background(), time.Minute, time.Minute)
	userID := "user2"

	assert.False(t, manager.ExistSession(userID), "Session should not exist initially")

	// Open and then close session
	_ = manager.OpenSession(userID)
	manager.CloseSession(userID)

	assert.False(t, manager.ExistSession(userID), "Session should not exist after closing")
}

func TestCloseNonExistingSession(t *testing.T) {
	manager := NewSessionManager(context.Background(), time.Minute, time.Minute)
	userID := "user3"
	manager.CloseSession(userID)

	assert.False(t, manager.ExistSession(userID), "Session should not exist after closing")
}

func TestGetSession(t *testing.T) {
	manager := NewSessionManager(context.Background(), time.Minute, time.Minute)
	userID := "user3"
	session := manager.OpenSession(userID)
	session.AddCommand("cmd")
	session.AddArgName("arg")
	gotSession := manager.GetSession(userID)

	assert.NotNil(t, gotSession, "Expected to retrieve a session")
	assert.Equal(t, "cmd", gotSession.Commands[0])
	assert.Equal(t, "arg", gotSession.Args[0])

	nonExistent := "ghost"
	ghostSession := manager.GetSession(nonExistent)

	assert.Nil(t, ghostSession, "Expected a non-nil session object for non-existent user")
}

func TestRemoveExpiredSessions(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	mgr := &SessionManager{
		sessions:      make(map[string]*Session),
		sessionTTL:    3 * time.Second,
		checkInterval: 100 * time.Millisecond,
		ctx:           ctx,
	}

	mgr.sessions["expired"] = &Session{OpenTime: time.Now().Add(-5 * time.Second)}
	mgr.sessions["active"] = &Session{OpenTime: time.Now()}

	go mgr.RemoveExpiredSessions()
	time.Sleep(1 * time.Second)

	mgr.mtx.Lock()
	defer mgr.mtx.Unlock()

	_, exists1 := mgr.sessions["expired"]
	_, exists2 := mgr.sessions["active"]

	assert.False(t, exists1, "expired should be removed as expired")
	assert.True(t, exists2, "active should still exist as not expired")
}
