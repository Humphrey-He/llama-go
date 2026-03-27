package service

import (
	"context"
	"testing"

	"llama-go/internal/backend"
	"llama-go/internal/session"
)

// MockBackend 模拟后端
type MockBackend struct{}

func (m *MockBackend) Chat(ctx context.Context, req backend.ChatRequest) (*backend.ChatResponse, error) {
	return &backend.ChatResponse{
		ID:      "test-id",
		Content: "Test response",
		Usage: backend.Usage{
			PromptTokens:     10,
			CompletionTokens: 5,
			TotalTokens:      15,
		},
	}, nil
}

func (m *MockBackend) StreamChat(ctx context.Context, req backend.ChatRequest) (<-chan backend.StreamChunk, error) {
	ch := make(chan backend.StreamChunk)
	go func() {
		ch <- backend.StreamChunk{Delta: "Test"}
		ch <- backend.StreamChunk{Done: true}
		close(ch)
	}()
	return ch, nil
}

func (m *MockBackend) Health(ctx context.Context) error {
	return nil
}

func TestNewChatService(t *testing.T) {
	backend := &MockBackend{}
	store := session.NewSessionStore()
	service := NewChatService(backend, store)

	if service == nil {
		t.Fatal("NewChatService returned nil")
	}
}

func TestChat(t *testing.T) {
	backend := &MockBackend{}
	store := session.NewSessionStore()
	service := NewChatService(backend, store)

	req := backend.ChatRequest{
		Model:     "test",
		Messages:  []backend.Message{{Role: "user", Content: "Hello"}},
		SessionID: "test-session",
	}

	resp, err := service.Chat(context.Background(), req)
	if err != nil {
		t.Errorf("Chat failed: %v", err)
	}

	if resp.Content != "Test response" {
		t.Errorf("expected 'Test response', got '%s'", resp.Content)
	}
}

func TestClearSession(t *testing.T) {
	backend := &MockBackend
	store := session.NewSessionStore()
	service := NewChatService(backend, store)

	sessionID := "test-session"
	store.AddMessage(sessionID, "user", "Hello")

	err := service.ClearSession(sessionID)
	if err != nil {
		t.Errorf("ClearSession failed: %v", err)
	}

	messages := store.GetMessages(sessionID)
	if messages != nil {
		t.Errorf("expected nil after clear, got %v", messages)
	}
}
