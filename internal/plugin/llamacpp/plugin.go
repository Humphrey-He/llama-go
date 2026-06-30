package llamacpp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"llama-go/internal/backend"
	"llama-go/internal/plugin"
)

// Plugin llama-cpp-python 插件
type Plugin struct {
	baseURL string
	client  *http.Client
}

func New() plugin.Plugin {
	return &Plugin{
		client: &http.Client{},
	}
}

func (p *Plugin) Init(ctx context.Context, config map[string]interface{}) error {
	url, ok := config["base_url"].(string)
	if !ok {
		return fmt.Errorf("missing base_url in config")
	}
	p.baseURL = url
	return p.Health(ctx)
}

func (p *Plugin) Close() error {
	return nil
}

func (p *Plugin) Health(ctx context.Context) error {
	req, _ := http.NewRequestWithContext(ctx, "GET", p.baseURL+"/health", nil)
	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("health check failed: %d", resp.StatusCode)
	}
	return nil
}

func (p *Plugin) Name() string    { return "llama-cpp-python" }
func (p *Plugin) Version() string { return "1.0.0" }
func (p *Plugin) Type() string    { return "llama-cpp" }

func (p *Plugin) Generate(ctx context.Context, req *backend.GenerateRequest) (*backend.GenerateResponse, error) {
	payload := map[string]interface{}{
		"model":       req.Model,
		"messages":    req.Messages,
		"temperature": req.Temperature,
		"max_tokens":  req.MaxTokens,
		"stream":      false,
	}

	body, _ := json.Marshal(payload)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/v1/chat/completions", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	choices := result["choices"].([]interface{})
	message := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	content := message["content"].(string)

	return &backend.GenerateResponse{
		ID:    result["id"].(string),
		Model: req.Model,
		Text:  content,
	}, nil
}

func (p *Plugin) GenerateStream(ctx context.Context, req *backend.GenerateRequest) (<-chan backend.StreamChunk, error) {
	ch := make(chan backend.StreamChunk)
	close(ch)
	return ch, fmt.Errorf("stream not implemented")
}

func (p *Plugin) ClearSession(ctx context.Context, sessionID string) error {
	return nil
}

func (p *Plugin) Info() backend.BackendInfo {
	return backend.BackendInfo{
		Name:           "llama-cpp-python",
		SupportsStream: false,
		MaxContextLen:  4096,
	}
}
