package session

import (
	"testing"
	"time"
)

func TestNewSessionStore(t *testing.T) {
	store := NewSessionStore()
	if store == nil {
		t.Fatal("NewSessionStore returned nil")
	}
}

func TestAddAndGetMessages(t *testing.T) {
	store := NewSessionStore()
	sessionID := "test-session"

	store.AddMessage(sessionID, "user", "Hello")
	store.AddMessage(sessionID, "assistant", "Hi there")

	messages := store.GetMessages(sessionID)
	if len(messages) != 2 {
		t.Errorf("expected 2 messages, got %d", len(messages))
	}

	if messages[0].Role != "user" || messages[0].Content != "Hello" {
		t.Errorf("first message incorrect")
	}
}

func TestClearSession(t *testing.T) {
	store := NewSessionStore()
	sessionID := "test-session"

	store.AddMessage(sessionID, "user", "Hello")
	store.Clear(sessionID)

	messages := store.GetMessages(sessionID)
	if messages != nil {
		t.Errorf("expected nil after clear, got %v", messages)
	}
}

func TestStoreSessionTTL(t *testing.T) {
	store := NewSessionStore()
	sessionID := "test-session"

	session := &Session{
		ID:        sessionID,
		Messages:  []Message{},
		CreatedAt: time.Now(),
		TTL:       1 * time.Millisecond,
	}

	store.Set(sessionID, session)
	time.Sleep(10 * time.Millisecond)

	retrieved := store.Get(sessionID)
	if retrieved != nil {
		t.Errorf("expected nil after TTL expiry, got %v", retrieved)
	}
}
