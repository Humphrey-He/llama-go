# RFC-001: KV Cache 优化机制

## 元信息
- **RFC 编号**：001
- **标题**：KV Cache 优化机制设计
- **状态**：已实现
- **创建日期**：2026-03-27
- **作者**：Llama-Go Team

## 背景

### 问题描述
Transformer 模型在生成文本时，每生成一个新 token 都需要重新计算所有历史 token 的 Attention 权重，导致：
- 计算复杂度：O(n²)
- 推理速度慢
- 资源浪费严重

### 目标
实现 KV Cache 机制，缓存已计算的 Key-Value 状态，避免重复计算。

## 设计方案

### 核心思路
在 Transformer 的 Self-Attention 中：
```
Attention(Q, K, V) = softmax(QK^T / √d) V
```

对于已生成的 token，其 K 和 V 是固定的，可以缓存复用。

### 数据结构
```go
type CacheEntry struct {
    Keys  [][]float32  // [num_heads, seq_len, head_dim]
    Vals  [][]float32  // [num_heads, seq_len, head_dim]
    Token int          // 对应的 token 位置
}
```

### 存储策略
- **Key**：sessionID（会话标识）
- **Value**：CacheEntry 列表（按 token 顺序）
- **并发控制**：sync.RWMutex 读写锁

## 实现细节

### Python 端（推理层）
```python
# 启用 KV Cache
outputs = model(**inputs, use_cache=True)

# 提取缓存
past_kv = outputs.past_key_values
kv = {
    "keys": past_kv[0][0].tolist(),
    "vals": past_kv[0][1].tolist()
}
```

### Go 端（缓存层）
```go
// 写入缓存
func (k *KVCache) Set(sessionID string, entry *CacheEntry) {
    k.mu.Lock()
    defer k.mu.Unlock()
    k.cache[sessionID] = append(k.cache[sessionID], entry)
}

// 读取缓存（并发安全）
func (k *KVCache) Get(sessionID string) []*CacheEntry {
    k.mu.RLock()
    defer k.mu.RUnlock()
    return k.cache[sessionID]
}
```

## 性能分析

### 理论提升
- **首次推理**：O(n²) → 无优化
- **后续推理**：O(n²) → O(n)
- **加速比**：2-5倍（取决于序列长度）

### 内存开销
- 每个 token：~4KB（1.1B 模型）
- 100 token 对话：~400KB
- 可接受范围

## 风险与限制

### 内存限制
- 长对话会占用大量内存
- **缓解方案**：设置最大缓存长度，超出后清理旧缓存

### 会话管理
- 当前使用固定 sessionID="default"
- **改进方向**：支持多会话并发

### 缓存一致性
- 缓存与模型状态需保持同步
- **保障措施**：每次推理后立即更新缓存
