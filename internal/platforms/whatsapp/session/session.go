package session

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/pagu-project/pagu/pkg/log"
)

type Session struct {
	OpenTime time.Time
	Commands []string
	Args     []string
}

func (s *Session) AddCommand(command string) {
	s.Commands = append(s.Commands, command)
}

func (s *Session) AddArgName(argName string) {
	s.Args = append(s.Args, fmt.Sprintf("--%s=", argName))
}

func (s *Session) AddArgValue(argValue string) {
	if len(s.Args) == 0 {
		log.Warn("No arg name found")

		return
	}

	s.Args[len(s.Args)-1] = fmt.Sprintf("%s%s", s.Args[len(s.Args)-1], argValue)
}

func (s *Session) GetLastCommand() string {
	if len(s.Commands) == 0 {
		return ""
	}

	return s.Commands[len(s.Commands)-1]
}

func (s *Session) GetLastArg() string {
	return s.Args[len(s.Args)-1]
}

func (s *Session) GetNumberOfArgs() int {
	return len(s.Args)
}

func (s *Session) GetCommandLine() string {
	if len(s.Commands) < 1 {
		log.Warn("No commands found", "commands", s.Commands)

		return ""
	}

	// Exclude the root command (`pagu`)
	commandWithoutRoot := s.Commands[1:]
	commandLine := strings.Join(commandWithoutRoot, " ")
	commandLine += " " + strings.Join(s.Args, " ")

	return commandLine
}

type SessionManager struct {
	ctx           context.Context
	mtx           sync.RWMutex
	sessions      map[string]*Session
	sessionTTL    time.Duration
	checkInterval time.Duration
}

func NewSessionManager(ctx context.Context, sessionTTL, checkInterval time.Duration) *SessionManager {
	return &SessionManager{
		ctx:           ctx,
		mtx:           sync.RWMutex{},
		sessions:      make(map[string]*Session),
		checkInterval: checkInterval,
		sessionTTL:    sessionTTL,
	}
}

func (mgr *SessionManager) ExistSession(userID string) bool {
	mgr.mtx.RLock()
	defer mgr.mtx.RUnlock()

	_, exist := mgr.sessions[userID]

	return exist
}

func (mgr *SessionManager) OpenSession(userID string) *Session {
	mgr.mtx.Lock()
	defer mgr.mtx.Unlock()

	session := &Session{
		OpenTime: time.Now(),
	}

	mgr.sessions[userID] = session

	return session
}

func (mgr *SessionManager) CloseSession(userID string) {
	mgr.mtx.Lock()
	defer mgr.mtx.Unlock()

	delete(mgr.sessions, userID)
}

func (mgr *SessionManager) GetSession(userID string) *Session {
	mgr.mtx.RLock()
	defer mgr.mtx.RUnlock()

	session := mgr.sessions[userID]

	return session
}

func (mgr *SessionManager) RemoveExpiredSessions() {
	for {
		now := time.Now()
		expiredSessions := []string{}
		select {
		case <-mgr.ctx.Done():
			return
		default:
			mgr.mtx.RLock()
			for id, session := range mgr.sessions {
				if now.Sub(session.OpenTime) > mgr.sessionTTL {
					expiredSessions = append(expiredSessions, id)
				}
			}
			mgr.mtx.RUnlock()

			// Now delete sessions with a write lock
			mgr.mtx.Lock()
			for _, id := range expiredSessions {
				delete(mgr.sessions, id)
			}
			mgr.mtx.Unlock()

			time.Sleep(mgr.checkInterval)
		}
	}
}
