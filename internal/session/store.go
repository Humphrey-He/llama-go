package session

import (
	"sync"
	"time"
)

// Message 消息
type Message struct {
	Role    string
	Content string
}

// Session 会话
type Session struct {
	ID        string
	Messages  []Message
	CreatedAt time.Time
	UpdatedAt time.Time
	TTL       time.Duration
}

// SessionStore 会话存储
type SessionStore struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

// NewSessionStore 创建会话存储
func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]*Session),
	}
}

// Get 获取会话
func (s *SessionStore) Get(id string) *Session {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.sessions[id]
	if !ok {
		return nil
	}

	// 检查 TTL
	if session.TTL > 0 && time.Since(session.UpdatedAt) > session.TTL {
		return nil
	}

	return session
}

// Set 保存会话
func (s *SessionStore) Set(id string, session *Session) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session.UpdatedAt = time.Now()
	s.sessions[id] = session
}

// Clear 清空会话
func (s *SessionStore) Clear(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, id)
}

// AddMessage 添加消息
func (s *SessionStore) AddMessage(id string, role, content string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[id]
	if !ok {
		session = &Session{
			ID:        id,
			Messages:  []Message{},
			CreatedAt: time.Now(),
			TTL:       24 * time.Hour,
		}
		s.sessions[id] = session
	}

	session.Messages = append(session.Messages, Message{
		Role:    role,
		Content: content,
	})
	session.UpdatedAt = time.Now()
}

// GetMessages 获取会话消息
func (s *SessionStore) GetMessages(id string) []Message {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.sessions[id]
	if !ok {
		return nil
	}

	return session.Messages
}
