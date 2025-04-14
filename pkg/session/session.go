package session

import (
	"context"
	"sync"
	"time"
)

type Session struct {
	OpenTime time.Time
	Commands []string
	Args     []string
}

type SessionManager struct {
	mtx           sync.RWMutex
	sessions      map[string]Session
	sessionTTL    time.Duration
	checkInterval time.Duration
	ctx           context.Context
}

func NewSessionManager(ctx context.Context) *SessionManager {
	return &SessionManager{
		mtx:      sync.RWMutex{},
		sessions: make(map[string]Session),
		ctx:      ctx,
	}
}

func NewSession(command []string) *Session {
	return &Session{
		Commands: command,
	}
}

func (mgr *SessionManager) SetConfig(sessionTTL, checkInterval time.Duration) {
	mgr.checkInterval = checkInterval
	mgr.sessionTTL = sessionTTL
}

func (mgr *SessionManager) ExistSession(userID string) bool {
	mgr.mtx.RLock()
	defer mgr.mtx.RUnlock()

	_, exist := mgr.sessions[userID]

	return exist
}

func (mgr *SessionManager) OpenSession(userID string, session Session) {
	mgr.mtx.Lock()
	defer mgr.mtx.Unlock()

	session.OpenTime = time.Now()
	mgr.sessions[userID] = session
}

func (mgr *SessionManager) CloseSession(userID string) {
	mgr.mtx.Lock()
	defer mgr.mtx.Unlock()

	_, exist := mgr.sessions[userID]
	if exist {
		delete(mgr.sessions, userID)
	}
}

func (mgr *SessionManager) GetSession(userID string) *Session {
	mgr.mtx.RLock()
	defer mgr.mtx.RUnlock()

	session := mgr.sessions[userID]

	return &session
}

func (mgr *SessionManager) RemoveExpiredSessions() {
	mgr.mtx.Lock()
	defer mgr.mtx.Unlock()

	for {
		now := time.Now()
		expiredSessions := []string{}
		select {
		case <-mgr.ctx.Done():
			return
		default:
			for id, session := range mgr.sessions {
				if now.Sub(session.OpenTime) > mgr.sessionTTL {
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
}
