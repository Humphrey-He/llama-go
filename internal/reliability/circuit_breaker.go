package reliability

import (
	"sync"
	"time"
)

// CircuitState 熔断器状态
type CircuitState string

const (
	StateClosed   CircuitState = "closed"
	StateOpen     CircuitState = "open"
	StateHalfOpen CircuitState = "half_open"
)

// CircuitBreakerConfig 熔断器配置
type CircuitBreakerConfig struct {
	FailureThreshold int
	SuccessThreshold int
	Timeout          time.Duration
}

// DefaultCircuitBreakerConfig 默认熔断器配置
func DefaultCircuitBreakerConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		FailureThreshold: 5,
		SuccessThreshold: 2,
		Timeout:          30 * time.Second,
	}
}

// CircuitBreaker 熔断器
type CircuitBreaker struct {
	mu              sync.RWMutex
	state           CircuitState
	failureCount    int
	successCount    int
	lastFailureTime time.Time
	config          CircuitBreakerConfig
}

// NewCircuitBreaker 创建熔断器
func NewCircuitBreaker(cfg CircuitBreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{
		state:  StateClosed,
		config: cfg,
	}
}

// Call 执行调用
func (cb *CircuitBreaker) Call(fn func() error) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == StateOpen {
		if time.Since(cb.lastFailureTime) > cb.config.Timeout {
			cb.state = StateHalfOpen
			cb.successCount = 0
		} else {
			return ErrCircuitOpen
		}
	}

	err := fn()

	if err != nil {
		cb.failureCount++
		cb.lastFailureTime = time.Now()
		if cb.failureCount >= cb.config.FailureThreshold {
			cb.state = StateOpen
		}
		return err
	}

	cb.failureCount = 0
	if cb.state == StateHalfOpen {
		cb.successCount++
		if cb.successCount >= cb.config.SuccessThreshold {
			cb.state = StateClosed
		}
	}

	return nil
}

// State 获取当前状态
func (cb *CircuitBreaker) State() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}
