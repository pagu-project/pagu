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
	defer sessionManage.mtx.RUnlock()

	_, exist := sessionManage.sessions[userID]

	return exist
}

func (sessionManage *SessionManager) OpenSession(userID string, session Session) {
	sessionManage.mtx.Lock()
	defer sessionManage.mtx.Unlock()

	session.openTime = time.Now()
	sessionManage.sessions[userID] = session
}

func (sessionManage *SessionManager) CloseSession(userID string) {
	sessionManage.mtx.Lock()
	defer sessionManage.mtx.Unlock()

	_, exist := sessionManage.sessions[userID]
	if exist {
		delete(sessionManage.sessions, userID)
	}
}

func (sessionManage *SessionManager) GetSession(userID string) *Session {
	sessionManage.mtx.RLock()
	defer sessionManage.mtx.RUnlock()

	session := sessionManage.sessions[userID]

	return &session
}

func (mgr *SessionManager) removeExpiredSessions() {
	mgr.mtx.Lock()
	defer mgr.mtx.Unlock()

	for {
		now := time.Now()
		expiredSessions := []string{}

		for id, session := range mgr.sessions {
			if now.Sub(session.openTime) > mgr.sessionTTL {
				expiredSessions = append(expiredSessions, id)
			}
		}

		// Now delete sessions with a write lock
		for _, id := range expiredSessions {
			delete(mgr.sessions, id)
		}

		time.Sleep(mgr.checkInterval)
	}
}
