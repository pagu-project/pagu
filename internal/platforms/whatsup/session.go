package whatsup

import (
	"sync"
	"time"
)

type Session struct {
	openTime time.Time
	commands []string
	args     []string
}

type SessionManager struct {
	mtx           sync.RWMutex
	sessions      map[string]Session
	sessionTtl    time.Duration
	checkInterval time.Duration //or check ticker
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		mtx:      sync.RWMutex{},
		sessions: make(map[string]Session),
	}
}

func (sessionManage *SessionManager) OpenSession(userId string) *Session {

	return &Session{}
}

func (sessionManage *SessionManager) CloseSession(userId string) {

	return
}

func (sessionManage *SessionManager) GetSession(userId string) *Session {

	return &Session{}
}

// in a gorutine , base on checkInterval
func (sessionManage *SessionManager) removeExpiredSession() {

	return
}
