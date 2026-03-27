package reliability

import (
	"context"
	"sync"
)

// HealthChecker 健康检查器
type HealthChecker struct {
	mu       sync.RWMutex
	checks   map[string]func(context.Context) error
	statuses map[string]string
}

// NewHealthChecker 创建健康检查器
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		checks:   make(map[string]func(context.Context) error),
		statuses: make(map[string]string),
	}
}

// Register 注册检查项
func (hc *HealthChecker) Register(name string, check func(context.Context) error) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	hc.checks[name] = check
	hc.statuses[name] = "unknown"
}

// Check 执行单个检查
func (hc *HealthChecker) Check(ctx context.Context, name string) error {
	hc.mu.RLock()
	check, ok := hc.checks[name]
	hc.mu.RUnlock()

	if !ok {
		return nil
	}

	err := check(ctx)
	hc.mu.Lock()
	if err != nil {
		hc.statuses[name] = "unhealthy"
	} else {
		hc.statuses[name] = "ok"
	}
	hc.mu.Unlock()

	return err
}

// CheckAll 执行所有检查
func (hc *HealthChecker) CheckAll(ctx context.Context) map[string]string {
	hc.mu.RLock()
	names := make([]string, 0, len(hc.checks))
	for name := range hc.checks {
		names = append(names, name)
	}
	hc.mu.RUnlock()

	for _, name := range names {
		hc.Check(ctx, name)
	}

	hc.mu.RLock()
	defer hc.mu.RUnlock()
	result := make(map[string]string)
	for k, v := range hc.statuses {
		result[k] = v
	}
	return result
}

// IsHealthy 检查是否健康
func (hc *HealthChecker) IsHealthy() bool {
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	for _, status := range hc.statuses {
		if status != "ok" {
			return false
		}
	}
	return true
}
