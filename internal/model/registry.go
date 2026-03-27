package model

import (
	"llama-go/internal/backend"
	"sync"
)

// ModelRegistry 模型注册表
type ModelRegistry struct {
	models map[string]*ModelInfo
	mu     sync.RWMutex
}

// ModelInfo 模型信息
type ModelInfo struct {
	ID              string
	Backend         backend.InferenceBackend
	ContextLength   int
	SupportsStream  bool
	SupportedModels []string
}

// NewModelRegistry 创建模型注册表
func NewModelRegistry() *ModelRegistry {
	return &ModelRegistry{
		models: make(map[string]*ModelInfo),
	}
}

// Register 注册模型
func (mr *ModelRegistry) Register(id string, b backend.InferenceBackend) {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	info := b.Info()
	mr.models[id] = &ModelInfo{
		ID:              id,
		Backend:         b,
		ContextLength:   info.MaxContextLen,
		SupportsStream:  info.SupportsStream,
		SupportedModels: info.SupportedModels,
	}
}

// Get 获取模型
func (mr *ModelRegistry) Get(id string) *ModelInfo {
	mr.mu.RLock()
	defer mr.mu.RUnlock()
	return mr.models[id]
}

// List 列出所有模型
func (mr *ModelRegistry) List() []*ModelInfo {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	models := make([]*ModelInfo, 0, len(mr.models))
	for _, m := range mr.models {
		models = append(models, m)
	}
	return models
}
