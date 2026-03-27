package service

import (
	"context"
	"llama-go/internal/backend"
	"llama-go/internal/session"
)

// ChatService 聊天服务
type ChatService struct {
	backend backend.InferenceBackend
	store   *session.SessionStore
}

// NewChatService 创建聊天服务
func NewChatService(b backend.InferenceBackend, s *session.SessionStore) *ChatService {
	return &ChatService{
		backend: b,
		store:   s,
	}
}

// Chat 聊天
func (cs *ChatService) Chat(ctx context.Context, req *backend.GenerateRequest) (*backend.GenerateResponse, error) {
	// 获取会话历史
	messages := cs.store.GetMessages(req.SessionID)

	// 构建完整消息列表
	allMessages := make([]backend.Message, len(messages))
	for i, msg := range messages {
		allMessages[i] = backend.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}
	allMessages = append(allMessages, req.Messages...)

	req.Messages = allMessages

	// 调用后端
	resp, err := cs.backend.Generate(ctx, req)
	if err != nil {
		return nil, err
	}

	// 保存消息到会话
	for _, msg := range req.Messages {
		cs.store.AddMessage(req.SessionID, msg.Role, msg.Content)
	}
	cs.store.AddMessage(req.SessionID, "assistant", resp.Text)

	return resp, nil
}

// StreamChat 流式聊天
func (cs *ChatService) StreamChat(ctx context.Context, req *backend.GenerateRequest) (<-chan backend.StreamChunk, error) {
	// 获取会话历史
	messages := cs.store.GetMessages(req.SessionID)

	// 构建完整消息列表
	allMessages := make([]backend.Message, len(messages))
	for i, msg := range messages {
		allMessages[i] = backend.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}
	allMessages = append(allMessages, req.Messages...)

	req.Messages = allMessages

	// 调用后端流式接口
	ch, err := cs.backend.GenerateStream(ctx, req)
	if err != nil {
		return nil, err
	}

	// 保存用户消息
	for _, msg := range req.Messages {
		cs.store.AddMessage(req.SessionID, msg.Role, msg.Content)
	}

	// 包装通道，收集完整响应
	wrappedCh := make(chan backend.StreamChunk)
	go func() {
		defer close(wrappedCh)

		var fullContent string
		for chunk := range ch {
			wrappedCh <- chunk
			if chunk.Content != "" {
				fullContent += chunk.Content
			}
		}

		// 保存助手响应
		if fullContent != "" {
			cs.store.AddMessage(req.SessionID, "assistant", fullContent)
		}
	}()

	return wrappedCh, nil
}

// ClearSession 清空会话
func (cs *ChatService) ClearSession(sessionID string) error {
	cs.store.Clear(sessionID)
	return nil
}
