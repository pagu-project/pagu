package whatsapp

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
	sessionTTL    time.Duration
	checkInterval time.Duration
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		mtx:      sync.RWMutex{},
		sessions: make(map[string]Session),
	}
}

func (sessionManage *SessionManager) ExistSession(userID string) bool {
	sessionManage.mtx.RLock()
	_, exist := sessionManage.sessions[userID]
	sessionManage.mtx.RUnlock()

	return exist
}

func (sessionManage *SessionManager) OpenSession(userID string, session Session) {
	sessionManage.mtx.Lock()
	session.openTime = time.Now()
	sessionManage.sessions[userID] = session
	sessionManage.mtx.Unlock()
}

func (sessionManage *SessionManager) CloseSession(userID string) {
	_, exist := sessionManage.sessions[userID]
	if exist {
		sessionManage.mtx.Lock()
		delete(sessionManage.sessions, userID)
		sessionManage.mtx.Unlock()
	}
}

func (sessionManage *SessionManager) GetSession(userID string) *Session {
	sessionManage.mtx.Lock()
	session := sessionManage.sessions[userID]
	sessionManage.mtx.Unlock()

	return &session
}

func (sessionManage *SessionManager) removeExpiredSession() {
	for {
		sessionManage.mtx.RLock()
		now := time.Now()
		expiredSessions := []string{}

		for id, session := range sessionManage.sessions {
			if now.Sub(session.openTime) > sessionManage.sessionTTL {
				expiredSessions = append(expiredSessions, id)
			}
		}
		sessionManage.mtx.RUnlock() // Release read lock

		// Now delete sessions with a write lock
		sessionManage.mtx.Lock()
		for _, id := range expiredSessions {
			delete(sessionManage.sessions, id)
		}
		sessionManage.mtx.Unlock()

		time.Sleep(sessionManage.checkInterval)
	}
}
