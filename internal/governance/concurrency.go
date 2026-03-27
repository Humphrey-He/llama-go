package governance

import (
	"sync"
)

// ConcurrencyLimiter 并发限制器
type ConcurrencyLimiter struct {
	mu       sync.RWMutex
	limiters map[string]int32
	limits   map[string]int32
}

// NewConcurrencyLimiter 创建并发限制器
func NewConcurrencyLimiter() *ConcurrencyLimiter {
	return &ConcurrencyLimiter{
		limiters: make(map[string]int32),
		limits:   make(map[string]int32),
	}
}

// SetLimit 设置限制
func (cl *ConcurrencyLimiter) SetLimit(key string, limit int32) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.limits[key] = limit
	if _, ok := cl.limiters[key]; !ok {
		cl.limiters[key] = 0
	}
}

// Acquire 获取许可
func (cl *ConcurrencyLimiter) Acquire(key string) bool {
	cl.mu.RLock()
	limit, ok := cl.limits[key]
	cl.mu.RUnlock()

	if !ok {
		return true
	}

	cl.mu.Lock()
	current := cl.limiters[key]
	cl.mu.Unlock()

	if current >= limit {
		return false
	}

	cl.mu.Lock()
	cl.limiters[key]++
	cl.mu.Unlock()
	return true
}

// Release 释放许可
func (cl *ConcurrencyLimiter) Release(key string) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	if cl.limiters[key] > 0 {
		cl.limiters[key]--
	}
}

// Current 获取当前并发数
func (cl *ConcurrencyLimiter) Current(key string) int32 {
	cl.mu.RLock()
	defer cl.mu.RUnlock()
	return cl.limiters[key]
}
