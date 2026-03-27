package reliability

import (
	"context"
	"errors"
	"time"
)

// RetryConfig 重试配置
type RetryConfig struct {
	MaxRetries      int
	InitialBackoff  time.Duration
	MaxBackoff      time.Duration
	BackoffMultiplier float64
}

// DefaultRetryConfig 默认重试配置
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:        2,
		InitialBackoff:    100 * time.Millisecond,
		MaxBackoff:        1 * time.Second,
		BackoffMultiplier: 2.0,
	}
}

// Retryable 判断是否可重试
func Retryable(err error) bool {
	if err == nil {
		return false
	}
	// 仅对连接失败和超时重试
	return errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, context.Canceled) == false
}

// RetryWithBackoff 带退避的重试
func RetryWithBackoff(ctx context.Context, cfg RetryConfig, fn func() error) error {
	var lastErr error
	backoff := cfg.InitialBackoff

	for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return ctx.Err()
			}
			if backoff < cfg.MaxBackoff {
				backoff = time.Duration(float64(backoff) * cfg.BackoffMultiplier)
				if backoff > cfg.MaxBackoff {
					backoff = cfg.MaxBackoff
				}
			}
		}

		lastErr = fn()
		if lastErr == nil || !Retryable(lastErr) {
			return lastErr
		}
	}

	return lastErr
}
