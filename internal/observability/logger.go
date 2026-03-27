package observability

import (
	"encoding/json"
	"os"
	"time"
)

// LogEntry 日志条目
type LogEntry struct {
	Timestamp    string `json:"ts"`
	Level        string `json:"level"`
	RequestID    string `json:"request_id,omitempty"`
	SessionID    string `json:"session_id,omitempty"`
	APIKeyID     string `json:"api_key_id,omitempty"`
	Route        string `json:"route,omitempty"`
	Method       string `json:"method,omitempty"`
	Model        string `json:"model,omitempty"`
	Backend      string `json:"backend,omitempty"`
	Stream       bool   `json:"stream,omitempty"`
	LatencyMS    int64  `json:"latency_ms,omitempty"`
	TTFTMS       int64  `json:"ttft_ms,omitempty"`
	PromptTokens int    `json:"prompt_tokens,omitempty"`
	CompTokens   int    `json:"completion_tokens,omitempty"`
	FinishReason string `json:"finish_reason,omitempty"`
	Status       int    `json:"status,omitempty"`
	ErrorCode    string `json:"error_code,omitempty"`
	ErrorType    string `json:"error_type,omitempty"`
	Message      string `json:"message,omitempty"`
}

// Logger 结构化日志记录器
type Logger struct {
	encoder *json.Encoder
}

// NewLogger 创建日志记录器
func NewLogger() *Logger {
	return &Logger{
		encoder: json.NewEncoder(os.Stdout),
	}
}

// Log 记录日志
func (l *Logger) Log(level string, entry LogEntry) {
	entry.Timestamp = time.Now().UTC().Format(time.RFC3339)
	entry.Level = level
	l.encoder.Encode(entry)
}

// Info 记录信息日志
func (l *Logger) Info(entry LogEntry) {
	l.Log("info", entry)
}

// Error 记录错误日志
func (l *Logger) Error(entry LogEntry) {
	l.Log("error", entry)
}

// Warn 记录警告日志
func (l *Logger) Warn(entry LogEntry) {
	l.Log("warn", entry)
}
