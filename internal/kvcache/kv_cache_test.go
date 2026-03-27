package kvcache

import (
	"sync"
	"testing"
)

func TestNewKVCache(t *testing.T) {
	cache := NewKVCache()
	if cache == nil {
		t.Fatal("NewKVCache returned nil")
	}
	if cache.cache == nil {
		t.Fatal("cache map not initialized")
	}
}

func TestSetAndGet(t *testing.T) {
	cache := NewKVCache()
	sessionID := "test-session"

	entry := &CacheEntry{
		Keys:  [][]float32{{1.0, 2.0}, {3.0, 4.0}},
		Vals:  [][]float32{{5.0, 6.0}, {7.0, 8.0}},
		Token: 10,
	}

	cache.Set(sessionID, entry)

	entries := cache.Get(sessionID)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	if entries[0].Token != 10 {
		t.Errorf("expected token 10, got %d", entries[0].Token)
	}
}

func TestClear(t *testing.T) {
	cache := NewKVCache()
	sessionID := "test-session"

	cache.Set(sessionID, &CacheEntry{Token: 1})
	cache.Clear(sessionID)

	entries := cache.Get(sessionID)
	if len(entries) != 0 {
		t.Errorf("expected 0 entries after clear, got %d", len(entries))
	}
}

func TestGetCacheSize(t *testing.T) {
	cache := NewKVCache()
	sessionID := "test-session"

	for i := 0; i < 5; i++ {
		cache.Set(sessionID, &CacheEntry{Token: i})
	}

	size := cache.GetCacheSize(sessionID)
	if size != 5 {
		t.Errorf("expected size 5, got %d", size)
	}
}

func TestConcurrentAccess(t *testing.T) {
	cache := NewKVCache()
	sessionID := "concurrent-test"
	var wg sync.WaitGroup

	// 100 个并发写
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			cache.Set(sessionID, &CacheEntry{Token: id})
		}(i)
	}

	// 100 个并发读
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cache.Get(sessionID)
		}()
	}

	wg.Wait()

	size := cache.GetCacheSize(sessionID)
	if size != 100 {
		t.Errorf("expected 100 entries, got %d", size)
	}
}

func TestSlidingWindow(t *testing.T) {
	cache := NewKVCache()
	sessionID := "sliding-test"

	// 写入超过 MaxCacheSize 的数据
	for i := 0; i < 1100; i++ {
		cache.Set(sessionID, &CacheEntry{Token: i})
	}

	size := cache.GetCacheSize(sessionID)
	if size > 1000 {
		t.Errorf("expected size <= 1000, got %d", size)
	}

	// 验证最旧的被删除
	entries := cache.Get(sessionID)
	if entries[0].Token < 100 {
		t.Errorf("expected oldest token >= 100, got %d", entries[0].Token)
	}
}
