package session

import (
	"context"
	"sync"

	tdtelegram "github.com/gotd/td/telegram"
	"github.com/sirupsen/logrus"
)

type (
	MessageUpdate struct {
		MessageID int64
		From      string
		FromID    int64
		Text      string
		Timestamp int64
	}

	TelegramSession struct {
		ID        string
		Client    *tdtelegram.Client
		Ctx       context.Context
		Cancel    context.CancelFunc
		isActive  bool
		UpdatesCh chan *MessageUpdate
		QRChan    chan string
		mu        sync.RWMutex
	}

	ManagerSession struct {
		sessions map[string]*TelegramSession
		mu       sync.RWMutex
	}
)

func NewTelegramSession(id string, client *tdtelegram.Client) *TelegramSession {
	ctx, cancel := context.WithCancel(context.Background())
	return &TelegramSession{
		ID:        id,
		Client:    client,
		Ctx:       ctx,
		Cancel:    cancel,
		isActive:  true,
		UpdatesCh: make(chan *MessageUpdate, 100),
		QRChan:    make(chan string, 1),
	}
}

func (s *TelegramSession) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isActive {
		s.isActive = false
		if s.Cancel != nil {
			s.Cancel()
		}
		if s.UpdatesCh != nil {
			close(s.UpdatesCh)
		}
		if s.QRChan != nil {
			close(s.QRChan)
		}
	}
}

func (s *TelegramSession) IsActive() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isActive
}

func (s *TelegramSession) SendUpdate(update *MessageUpdate) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.isActive {
		return
	}
	select {
	case s.UpdatesCh <- update:
	default:
	}
}

func NewSessionManager() *ManagerSession {
	return &ManagerSession{
		sessions: make(map[string]*TelegramSession),
	}
}

func (sm *ManagerSession) CreateSession(id string, client *tdtelegram.Client) *TelegramSession {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sess := NewTelegramSession(id, client)
	sm.sessions[id] = sess
	return sess
}

func (sm *ManagerSession) GetSession(id string) (*TelegramSession, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	sess, ok := sm.sessions[id]
	return sess, ok
}

func (sm *ManagerSession) DeleteSession(id string) bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sess, ok := sm.sessions[id]
	if !ok {
		return false
	}

	if sess.Client != nil && sess.IsActive() {
		api := sess.Client.API()
		_, err := api.AuthLogOut(sess.Ctx)
		if err != nil {
			logrus.Fatal("Failed to logout session %s: %v", id, err)
		}
	}

	sess.Stop()

	delete(sm.sessions, id)
	return true
}

func (sm *ManagerSession) ListSessions() []string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	ids := make([]string, 0, len(sm.sessions))
	for id, sess := range sm.sessions {
		if sess.IsActive() {
			ids = append(ids, id)
		}
	}
	return ids
}
