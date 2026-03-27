package backend

import (
	"testing"
)

func TestNewVLLMBackend(t *testing.T) {
	b := NewVLLMBackend("http://localhost:8000")
	if b == nil {
		t.Fatal("NewVLLMBackend returned nil")
	}
	if b.baseURL != "http://localhost:8000" {
		t.Errorf("expected baseURL http://localhost:8000, got %s", b.baseURL)
	}
}

