package stream

import (
	"encoding/json"
	"fmt"
	"io"
)

// SSEWriter 服务端事件写入器
type SSEWriter struct {
	w io.Writer
}

// NewSSEWriter 创建 SSE 写入器
func NewSSEWriter(w io.Writer) *SSEWriter {
	return &SSEWriter{w: w}
}

// WriteChunk 写入数据块
func (sw *SSEWriter) WriteChunk(data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(sw.w, "data: %s\n\n", string(bytes))
	return err
}

// WriteDone 写入完成标记
func (sw *SSEWriter) WriteDone() error {
	_, err := fmt.Fprint(sw.w, "data: [DONE]\n\n")
	return err
}

// Flush 刷新输出
func (sw *SSEWriter) Flush(flusher interface{ Flush() }) {
	if f, ok := flusher.(interface{ Flush() }); ok {
		f.Flush()
	}
}
