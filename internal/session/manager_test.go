package session

import (
	"testing"
	"time"
)

func TestNewSessionManager(t *testing.T) {
	sm := NewSessionManager(100, 1*time.Hour)
	if sm == nil {
		t.Fatal("NewSessionManager returned nil")
	}
}

func TestCreateSession(t *testing.T) {
	sm := NewSessionManager(100, 1*time.Hour)
	meta := sm.CreateSession("test-session")

	if meta.SessionID != "test-session" {
		t.Errorf("expected session_id test-session, got %s", meta.SessionID)
	}
}

func TestGetSession(t *testing.T) {
	sm := NewSessionManager(100, 1*time.Hour)
	sm.CreateSession("test-session")

	meta := sm.GetSession("test-session")
	if meta == nil {
		t.Fatal("GetSession returned nil")
	}

	if meta.SessionID != "test-session" {
		t.Errorf("expected session_id test-session, got %s", meta.SessionID)
	}
}

func TestDeleteSession(t *testing.T) {
	sm := NewSessionManager(100, 1*time.Hour)
	sm.CreateSession("test-session")
	sm.DeleteSession("test-session")

	meta := sm.GetSession("test-session")
	if meta != nil {
		t.Errorf("expected nil after delete, got %v", meta)
	}
}

func TestSessionTTL(t *testing.T) {
	sm := NewSessionManager(100, 100*time.Millisecond)
	sm.CreateSession("test-session")

	time.Sleep(200 * time.Millisecond)

	meta := sm.GetSession("test-session")
	if meta != nil {
		t.Errorf("expected nil after TTL expiry, got %v", meta)
	}
}
