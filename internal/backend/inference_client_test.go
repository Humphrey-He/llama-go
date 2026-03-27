package backend

import (
	"context"
	"testing"
)

func TestNewInferenceClient(t *testing.T) {
	client := NewInferenceClient("http://localhost:8000")
	if client == nil {
		t.Fatal("NewInferenceClient returned nil")
	}
	if client.baseURL != "http://localhost:8000" {
		t.Errorf("expected baseURL http://localhost:8000, got %s", client.baseURL)
	}
}

func TestGenerateRequest(t *testing.T) {
	req := GenerateRequest{
		SessionID:    "test-session",
		Prompt:       "Hello",
		MaxNewTokens: 128,
		Temperature:  0.7,
		TopP:         0.9,
	}

	if req.SessionID != "test-session" {
		t.Errorf("expected session_id test-session, got %s", req.SessionID)
	}
}
