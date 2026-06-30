package plugin

import (
	"context"
	"llama-go/internal/backend"
)

// Plugin 后端插件接口（扩展 InferenceBackend）
type Plugin interface {
	backend.InferenceBackend

	// 插件生命周期
	Init(ctx context.Context, config map[string]interface{}) error
	Close() error
	Health(ctx context.Context) error

	// 插件元数据
	Name() string
	Version() string
	Type() string // "llama-cpp", "ollama", "vllm"
}

// Metadata 插件元数据
type Metadata struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	Capabilities []string `json:"capabilities"`
}
