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
	session := Session{}
	session.openTime = time.Now()
	sessionManage.mtx.Lock()
	sessionManage.sessions[userId] = session
	sessionManage.mtx.Unlock()
	return &session
}

func (sessionManage *SessionManager) EditSession(userId string, command []string, args []string) *Session {
	sessionManage.mtx.Lock()
	session := sessionManage.sessions[userId]
	if command != nil {
		session.commands = command
	}

	if args != nil {
		session.args = args
	}

	sessionManage.mtx.Unlock()
	return &session
}

func (sessionManage *SessionManager) CloseSession(userId string) {
	_, exist := sessionManage.sessions[userId]
	if exist {
		sessionManage.mtx.Lock()
		delete(sessionManage.sessions, userId)
		sessionManage.mtx.Unlock()
	}
	return
}

func (sessionManage *SessionManager) GetSession(userId string) *Session {
	sessionManage.mtx.Lock()
	session := sessionManage.sessions[userId]
	sessionManage.mtx.Unlock()
	return &session
}

func (sessionManage *SessionManager) removeExpiredSession() {
	for {
		sessionManage.mtx.RLock()
		now := time.Now()
		expiredSessions := []string{}

		for id, session := range sessionManage.sessions {
			if now.Sub(session.openTime) > sessionManage.sessionTtl {
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
