package reliability

import (
	"context"
	"time"
)

// TimeoutConfig 超时配置
type TimeoutConfig struct {
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	RequestTimeout   time.Duration
	BackendTimeout   time.Duration
	StreamIdleTimeout time.Duration
	ShutdownTimeout  time.Duration
}

// DefaultTimeoutConfig 默认超时配置
func DefaultTimeoutConfig() TimeoutConfig {
	return TimeoutConfig{
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		RequestTimeout:    30 * time.Second,
		BackendTimeout:    25 * time.Second,
		StreamIdleTimeout: 5 * time.Minute,
		ShutdownTimeout:   10 * time.Second,
	}
}

// WithTimeout 为 context 添加超时
func WithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, timeout)
}
