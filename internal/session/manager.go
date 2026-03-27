package session

import (
	"sync"
	"time"
)

// SessionMetadata 会话元数据
type SessionMetadata struct {
	SessionID   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Model       string
	TokenCount  int
}

// SessionManager 会话管理器
type SessionManager struct {
	sessions map[string]*SessionMetadata
	mu       sync.RWMutex
	maxSessions int
	ttl      time.Duration
}

// NewSessionManager 创建会话管理器
func NewSessionManager(maxSessions int, ttl time.Duration) *SessionManager {
	sm := &SessionManager{
		sessions:    make(map[string]*SessionMetadata),
		maxSessions: maxSessions,
		ttl:         ttl,
	}

	// 启动 TTL 清理任务
	go sm.cleanupExpired()

	return sm
}

// CreateSession 创建会话
func (sm *SessionManager) CreateSession(sessionID string) *SessionMetadata {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// LRU 淘汰
	if len(sm.sessions) >= sm.maxSessions {
		oldest := ""
		oldestTime := time.Now()
		for id, meta := range sm.sessions {
			if meta.UpdatedAt.Before(oldestTime) {
				oldest = id
				oldestTime = meta.UpdatedAt
			}
		}
		if oldest != "" {
			delete(sm.sessions, oldest)
		}
	}

	meta := &SessionMetadata{
		SessionID: sessionID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	sm.sessions[sessionID] = meta
	return meta
}

// GetSession 获取会话
func (sm *SessionManager) GetSession(sessionID string) *SessionMetadata {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	meta, ok := sm.sessions[sessionID]
	if !ok {
		return nil
	}

	// 检查 TTL
	if time.Since(meta.UpdatedAt) > sm.ttl {
		return nil
	}

	return meta
}

// UpdateSession 更新会话
func (sm *SessionManager) UpdateSession(sessionID string, tokenCount int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if meta, ok := sm.sessions[sessionID]; ok {
		meta.UpdatedAt = time.Now()
		meta.TokenCount = tokenCount
	}
}

// DeleteSession 删除会话
func (sm *SessionManager) DeleteSession(sessionID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.sessions, sessionID)
}

// cleanupExpired 清理过期会话
func (sm *SessionManager) cleanupExpired() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		sm.mu.Lock()
		now := time.Now()
		for id, meta := range sm.sessions {
			if now.Sub(meta.UpdatedAt) > sm.ttl {
				delete(sm.sessions, id)
			}
		}
		sm.mu.Unlock()
	}
}

// GetStats 获取统计信息
func (sm *SessionManager) GetStats() map[string]interface{} {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return map[string]interface{}{
		"total_sessions": len(sm.sessions),
	}
}
