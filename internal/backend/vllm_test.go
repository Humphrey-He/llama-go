package backend

import (
	"context"
	"testing"
)

func TestNewVLLMBackend(t *testing.T) {
	backend := NewVLLMBackend("http://localhost:8000")
	if backend == nil {
		t.Fatal("NewVLLMBackend returned nil")
	}
	if backend.baseURL != "http://localhost:8000" {
		t.Errorf("expected baseURL http://localhost:8000, got %s", backend.baseURL)
	}
}

func TestVLLMBackendHealth(t *testing.T) {
	backend := NewVLLMBackend("http://localhost:8000")
	ctx := context.Background()

	// 这个测试需要实际的 vLLM 服务运行
	// err := backend.Health(ctx)
	// if err != nil {
	//     t.Errorf("Health check failed: %v", err)
	// }
}
