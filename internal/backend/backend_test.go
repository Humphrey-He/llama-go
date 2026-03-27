package backend

import (
	"context"
	"testing"
)

func TestPythonBackendInfo(t *testing.T) {
	pb := NewPythonBackend("http://localhost:8000")
	info := pb.Info()

	if info.Name != "python" {
		t.Errorf("expected name python, got %s", info.Name)
	}
	if !info.SupportsStream {
		t.Error("expected SupportsStream true")
	}
	if info.MaxContextLen != 2048 {
		t.Errorf("expected MaxContextLen 2048, got %d", info.MaxContextLen)
	}
}

func TestMockBackendGenerate(t *testing.T) {
	mb := NewMockBackend()
	req := &GenerateRequest{
		RequestID: "req-123",
		Model:     "mock-model",
		Prompt:    "test",
	}

	resp, err := mb.Generate(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ID != "req-123" {
		t.Errorf("expected id req-123, got %s", resp.ID)
	}
	if resp.FinishReason != "stop" {
		t.Errorf("expected finish_reason stop, got %s", resp.FinishReason)
	}
}

func TestMockBackendStream(t *testing.T) {
	mb := NewMockBackend()
	req := &GenerateRequest{
		RequestID: "req-123",
		Model:     "mock-model",
		Stream:    true,
	}

	chunks, err := mb.GenerateStream(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	count := 0
	for range chunks {
		count++
	}

	if count == 0 {
		t.Error("expected at least one chunk")
	}
}

func TestMockBackendSetResponse(t *testing.T) {
	mb := NewMockBackend()
	mb.SetResponse("custom-model", "Custom response")

	req := &GenerateRequest{
		RequestID: "req-123",
		Model:     "custom-model",
	}

	resp, err := mb.Generate(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Text != "Custom response" {
		t.Errorf("expected text 'Custom response', got %s", resp.Text)
	}
}
