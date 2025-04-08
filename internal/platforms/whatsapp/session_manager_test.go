package whatsapp

import (
	"context"
	"testing"
	"time"
)

func TestOpenAndExistSession(t *testing.T) {
	manager := NewSessionManager()
	userID := "user1"
	session := Session{
		commands: []string{"start"},
		args:     []string{"arg1"},
	}

	// Check that session does not exist
	if manager.ExistSession(userID) {
		t.Errorf("Expected session to not exist initially")
	}

	// Open session and check existence
	manager.OpenSession(userID, session)
	if !manager.ExistSession(userID) {
		t.Errorf("Expected session to exist after opening")
	}

	// Check another user that was never added
	otherUser := "nonexistent"
	if manager.ExistSession(otherUser) {
		t.Errorf("Expected session to not exist for a different user")
	}
}

func TestCloseSession(t *testing.T) {
	manager := NewSessionManager()
	userID := "user2"
	session := Session{}

	if manager.ExistSession(userID) {
		t.Errorf("Expected session to not exist initially")
	}

	// Open and then close session
	manager.OpenSession(userID, session)
	manager.CloseSession(userID)

	if manager.ExistSession(userID) {
		t.Errorf("Expected session to be removed after closing")
	}

	// Try closing again (should be a no-op, but not crash)
	manager.CloseSession(userID)
}

func TestGetSession(t *testing.T) {
	manager := NewSessionManager()
	userID := "user3"
	session := Session{
		commands: []string{"cmd"},
		args:     []string{"arg"},
	}

	// Open session and retrieve it
	manager.OpenSession(userID, session)
	gotSession := manager.GetSession(userID)

	if gotSession == nil {
		t.Fatal("Expected to retrieve a session, got nil")
	}

	if gotSession.commands[0] != "cmd" || gotSession.args[0] != "arg" {
		t.Errorf("Session data mismatch: got %+v", gotSession)
	}

	// Get session for non-existent user
	nonExistent := "ghost"
	ghostSession := manager.GetSession(nonExistent)
	if ghostSession == nil {
		t.Fatal("Expected GetSession to return a pointer (even if empty)")
	}

	if len(ghostSession.commands) != 0 || len(ghostSession.args) != 0 {
		t.Errorf("Expected empty session data, got %+v", ghostSession)
	}
}

func TestRemoveExpiredSessions(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	mgr := &SessionManager{
		sessions:      make(map[string]Session),
		sessionTTL:    1 * time.Second,        // 1 second TTL for testing
		checkInterval: 100 * time.Millisecond, // short check interval
		ctx:           ctx,
		cancle:        cancel,
	}

	mgr.sessions["session1"] = Session{openTime: time.Now().Add(-2 * time.Second)}        // expired
	mgr.sessions["session2"] = Session{openTime: time.Now().Add(-500 * time.Millisecond)} // not expired

	go mgr.removeExpiredSessions()

	time.Sleep(2 * time.Second)

	mgr.mtx.Lock()
	defer mgr.mtx.Unlock()

	if _, exists := mgr.sessions["session1"]; exists {
		t.Errorf("Expected session 'session1' to be removed, but it still exists")
	}
	if _, exists := mgr.sessions["session2"]; exists {
		t.Errorf("Expected session 'session2' to still exist, but it was removed")
	}
}
