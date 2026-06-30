package plugin

import (
	"context"
	"fmt"
	"sync"
)

// Manager 插件管理器
type Manager struct {
	plugins map[string]Plugin
	mu      sync.RWMutex
}

// NewManager 创建插件管理器
func NewManager() *Manager {
	return &Manager{
		plugins: make(map[string]Plugin),
	}
}

// Register 注册插件
func (m *Manager) Register(name string, p Plugin) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.plugins[name]; exists {
		return fmt.Errorf("plugin %s already registered", name)
	}

	m.plugins[name] = p
	return nil
}

// Get 获取插件
func (m *Manager) Get(name string) (Plugin, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	p, ok := m.plugins[name]
	if !ok {
		return nil, fmt.Errorf("plugin %s not found", name)
	}
	return p, nil
}

// List 列出所有插件
func (m *Manager) List() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.plugins))
	for name := range m.plugins {
		names = append(names, name)
	}
	return names
}

// InitAll 初始化所有插件
func (m *Manager) InitAll(ctx context.Context, configs map[string]map[string]interface{}) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for name, p := range m.plugins {
		config := configs[name]
		if err := p.Init(ctx, config); err != nil {
			return fmt.Errorf("init plugin %s: %w", name, err)
		}
	}
	return nil
}

// CloseAll 关闭所有插件
func (m *Manager) CloseAll() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, p := range m.plugins {
		if err := p.Close(); err != nil {
			return err
		}
	}
	return nil
}
