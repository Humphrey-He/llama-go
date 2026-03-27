package auth

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

// APIKey API密钥信息
type APIKey struct {
	ID       string
	KeyHash  string
	Enabled  bool
	TenantID string
}

// KeyStore API密钥存储
type KeyStore struct {
	mu   sync.RWMutex
	keys map[string]*APIKey
}

// NewKeyStore 创建密钥存储
func NewKeyStore() *KeyStore {
	return &KeyStore{
		keys: make(map[string]*APIKey),
	}
}

// AddKey 添加密钥
func (ks *KeyStore) AddKey(id, key, tenantID string) {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	hash := sha256.Sum256([]byte(key))
	ks.keys[id] = &APIKey{
		ID:       id,
		KeyHash:  fmt.Sprintf("%x", hash),
		Enabled:  true,
		TenantID: tenantID,
	}
}

// ValidateKey 验证密钥
func (ks *KeyStore) ValidateKey(id, key string) (*APIKey, bool) {
	ks.mu.RLock()
	defer ks.mu.RUnlock()

	apiKey, ok := ks.keys[id]
	if !ok || !apiKey.Enabled {
		return nil, false
	}

	hash := sha256.Sum256([]byte(key))
	if apiKey.KeyHash != fmt.Sprintf("%x", hash) {
		return nil, false
	}

	return apiKey, true
}

// GetKey 获取密钥信息
func (ks *KeyStore) GetKey(id string) (*APIKey, bool) {
	ks.mu.RLock()
	defer ks.mu.RUnlock()

	apiKey, ok := ks.keys[id]
	return apiKey, ok && apiKey.Enabled
}

// DisableKey 禁用密钥
func (ks *KeyStore) DisableKey(id string) {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	if apiKey, ok := ks.keys[id]; ok {
		apiKey.Enabled = false
	}
}
