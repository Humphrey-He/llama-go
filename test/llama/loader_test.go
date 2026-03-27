package llama

import (
	"testing"
	"llama-go/pkg/llama"
)

func TestNewModelLoader(t *testing.T) {
	loader := llama.NewModelLoader("test.gguf")
	if loader == nil {
		t.Fatal("NewModelLoader returned nil")
	}
}

func TestLoadModel(t *testing.T) {
	loader := llama.NewModelLoader("nonexistent.gguf")
	err := loader.Load()
	// TODO: 完善测试逻辑
	_ = err
}
