package service

import (
	"context"
	"testing"

	"llama-go/internal/backend"
	"llama-go/internal/session"
)

// MockBackend 模拟后端
type MockBackend struct{}

func (m *MockBackend) Generate(ctx context.Context, req *backend.GenerateRequest) (*backend.GenerateResponse, error) {
	return &backend.GenerateResponse{
		ID:           "test-id",
		Model:        req.Model,
		Text:         "Test response",
		FinishReason: "stop",
		Usage: backend.Usage{
			PromptTokens:     10,
			CompletionTokens: 5,
			TotalTokens:      15,
		},
	}, nil
}

func (m *MockBackend) GenerateStream(ctx context.Context, req *backend.GenerateRequest) (<-chan backend.StreamChunk, error) {
	ch := make(chan backend.StreamChunk)
	go func() {
		ch <- backend.StreamChunk{Content: "Test", Done: false}
		ch <- backend.StreamChunk{Content: "", Done: true}
		close(ch)
	}()
	return ch, nil
}

func (m *MockBackend) ClearSession(ctx context.Context, sessionID string) error {
	return nil
}

func (m *MockBackend) Info() backend.BackendInfo {
	return backend.BackendInfo{
		Name:           "mock",
		SupportsStream: true,
		MaxContextLen:  2048,
	}
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
	b := &MockBackend{}
	store := session.NewSessionStore()
	service := NewChatService(b, store)

	req := &backend.GenerateRequest{
		Model:     "test",
		Messages:  []backend.Message{{Role: "user", Content: "Hello"}},
		SessionID: "test-session",
	}

	resp, err := service.Chat(context.Background(), req)
	if err != nil {
		t.Errorf("Chat failed: %v", err)
	}

	if resp.Text != "Test response" {
		t.Errorf("expected 'Test response', got '%s'", resp.Text)
	}
}

func TestClearSession(t *testing.T) {
	backend := &MockBackend{}
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
