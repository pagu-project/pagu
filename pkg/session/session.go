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
	Mtx           sync.RWMutex
	Sessions      map[string]Session
	SessionTTL    time.Duration
	CheckInterval time.Duration
	Ctx           context.Context
}

func NewSessionManager(ctx context.Context) *SessionManager {
	return &SessionManager{
		Mtx:      sync.RWMutex{},
		Sessions: make(map[string]Session),
		Ctx:      ctx,
	}
}

func (mgr *SessionManager) ExistSession(userID string) bool {
	mgr.Mtx.RLock()
	defer mgr.Mtx.RUnlock()

	_, exist := mgr.Sessions[userID]

	return exist
}

func (mgr *SessionManager) OpenSession(userID string, session Session) {
	mgr.Mtx.Lock()
	defer mgr.Mtx.Unlock()

	session.OpenTime = time.Now()
	mgr.Sessions[userID] = session
}

func (mgr *SessionManager) CloseSession(userID string) {
	mgr.Mtx.Lock()
	defer mgr.Mtx.Unlock()

	_, exist := mgr.Sessions[userID]
	if exist {
		delete(mgr.Sessions, userID)
	}
}

func (mgr *SessionManager) GetSession(userID string) *Session {
	mgr.Mtx.RLock()
	defer mgr.Mtx.RUnlock()

	session := mgr.Sessions[userID]

	return &session
}

func (mgr *SessionManager) RemoveExpiredSessions() {
	mgr.Mtx.Lock()
	defer mgr.Mtx.Unlock()

	for {
		now := time.Now()
		expiredSessions := []string{}
		select {
		case <-mgr.Ctx.Done():
			return
		default:
			for id, session := range mgr.Sessions {
				if now.Sub(session.OpenTime) > mgr.SessionTTL {
					expiredSessions = append(expiredSessions, id)
				}
			}

			// Now delete sessions with a write lock
			for _, id := range expiredSessions {
				delete(mgr.Sessions, id)
			}

			time.Sleep(mgr.CheckInterval)
		}
	}
}
