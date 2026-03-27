package ratelimit

import (
	"sync"
	"time"
)

// RateLimiter 限流器
type RateLimiter struct {
	mu       sync.RWMutex
	limiters map[string]*TokenBucket
}

// TokenBucket 令牌桶
type TokenBucket struct {
	capacity      float64
	tokens        float64
	refillRate    float64
	lastRefillTime time.Time
	mu            sync.Mutex
}

// NewRateLimiter 创建限流器
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*TokenBucket),
	}
}

// NewTokenBucket 创建令牌桶
func NewTokenBucket(capacity, refillRate float64) *TokenBucket {
	return &TokenBucket{
		capacity:       capacity,
		tokens:         capacity,
		refillRate:     refillRate,
		lastRefillTime: time.Now(),
	}
}

// Allow 检查是否允许
func (tb *TokenBucket) Allow(tokens float64) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefillTime).Seconds()
	tb.tokens = min(tb.capacity, tb.tokens+elapsed*tb.refillRate)
	tb.lastRefillTime = now

	if tb.tokens >= tokens {
		tb.tokens -= tokens
		return true
	}
	return false
}

// AddLimiter 添加限流器
func (rl *RateLimiter) AddLimiter(key string, capacity, refillRate float64) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.limiters[key] = NewTokenBucket(capacity, refillRate)
}

// Allow 检查是否允许
func (rl *RateLimiter) Allow(key string, tokens float64) bool {
	rl.mu.RLock()
	limiter, ok := rl.limiters[key]
	rl.mu.RUnlock()

	if !ok {
		return true
	}
	return limiter.Allow(tokens)
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
