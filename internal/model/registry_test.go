package model

import (
	"context"
	"testing"

	"llama-go/internal/backend"
)

func TestModelRegistryRegister(t *testing.T) {
	registry := NewModelRegistry()
	mockBackend := backend.NewMockBackend()

	registry.Register("test-model", mockBackend)

	model := registry.Get("test-model")
	if model == nil {
		t.Fatal("expected model to be registered")
	}
	if model.ID != "test-model" {
		t.Errorf("expected id test-model, got %s", model.ID)
	}
}

func TestModelRegistryGet(t *testing.T) {
	registry := NewModelRegistry()
	mockBackend := backend.NewMockBackend()
	registry.Register("test-model", mockBackend)

	model := registry.Get("test-model")
	if model == nil {
		t.Fatal("expected model to exist")
	}

	notFound := registry.Get("nonexistent")
	if notFound != nil {
		t.Error("expected nil for nonexistent model")
	}
}

func TestModelRegistryList(t *testing.T) {
	registry := NewModelRegistry()
	mockBackend := backend.NewMockBackend()

	registry.Register("model1", mockBackend)
	registry.Register("model2", mockBackend)

	models := registry.List()
	if len(models) != 2 {
		t.Errorf("expected 2 models, got %d", len(models))
	}
}

func TestModelRegistryBackendAccess(t *testing.T) {
	registry := NewModelRegistry()
	mockBackend := backend.NewMockBackend()
	registry.Register("test-model", mockBackend)

	model := registry.Get("test-model")
	req := &backend.GenerateRequest{
		RequestID: "req-123",
		Model:     "test-model",
	}

	resp, err := model.Backend.Generate(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("expected response")
	}
}
