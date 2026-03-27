package kvcache

import (
	"llama-go/internal/config"
	"sync"
)

// KV 缓存单条数据
type CacheEntry struct {
	Keys  [][]float32 `json:"keys"`  // Key 张量
	Vals  [][]float32 `json:"vals"`  // Value 张量
	Token int         `json:"token"` // 对应 token 下标
}

// KVCache 全局缓存结构体（并发安全）
type KVCache struct {
	// key: session id / sequence id
	// value: 该序列的所有 KV 缓存
	cache map[string][]*CacheEntry
	mu    sync.RWMutex // 读写锁：高读低写，Go 核心优势
}

// NewKVCache 初始化缓存
func NewKVCache() *KVCache {
	return &KVCache{
		cache: make(map[string][]*CacheEntry),
	}
}

// Set 写入缓存（写锁 + 滑动窗口）
func (k *KVCache) Set(sessionID string, entry *CacheEntry) {
	k.mu.Lock()
	defer k.mu.Unlock()

	entries := k.cache[sessionID]

	// 滑动窗口：超过最大长度时删除最旧的
	if len(entries) >= config.MaxCacheSize {
		entries = entries[1:]
	}

	k.cache[sessionID] = append(entries, entry)
}

// Get 读取缓存（读锁，并发安全）
func (k *KVCache) Get(sessionID string) []*CacheEntry {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.cache[sessionID]
}

// Clear 清空缓存
func (k *KVCache) Clear(sessionID string) {
	k.mu.Lock()
	defer k.mu.Unlock()
	delete(k.cache, sessionID)
}

// GetCacheSize 获取缓存长度
func (k *KVCache) GetCacheSize(sessionID string) int {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return len(k.cache[sessionID])
}
