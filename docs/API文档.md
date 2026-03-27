# API 接口文档

## 基础信息

- **Base URL**: `http://localhost:8080`
- **Content-Type**: `application/json`

## 接口列表

### 1. 文本生成

**接口**: `POST /api/generate`

**描述**: 调用 LLM 进行文本生成，支持 KV Cache 复用

**请求参数**:
```json
{
  "prompt": "Hello, how are you?",
  "session_id": "user-123"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| prompt | string | 是 | 输入文本 |
| session_id | string | 否 | 会话ID，默认为 "default" |

**响应示例**:
```json
{
  "success": true,
  "data": {
    "generated_text": "I'm doing well, thank you!",
    "cache_size": 1,
    "session_id": "user-123",
    "cache_reused": false
  }
}
```

**错误响应**:
```json
{
  "success": false,
  "error": "invalid request: prompt is required"
}
```

### 2. 清空缓存

**接口**: `POST /api/clear`

**描述**: 清空指定会话的 KV Cache

**请求参数**:
```json
{
  "session_id": "user-123"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| session_id | string | 否 | 会话ID，默认为 "default" |

**响应示例**:
```json
{
  "success": true,
  "data": {
    "message": "cache cleared",
    "session_id": "user-123"
  }
}
```

## 使用示例

### cURL
```bash
# 文本生成
curl -X POST http://localhost:8080/api/generate \
  -H "Content-Type: application/json" \
  -d '{"prompt": "Hello", "session_id": "test"}'

# 清空缓存
curl -X POST http://localhost:8080/api/clear \
  -H "Content-Type: application/json" \
  -d '{"session_id": "test"}'
```

### Python
```python
import requests

# 文本生成
response = requests.post(
    "http://localhost:8080/api/generate",
    json={"prompt": "Hello", "session_id": "test"}
)
print(response.json())

# 清空缓存
response = requests.post(
    "http://localhost:8080/api/clear",
    json={"session_id": "test"}
)
print(response.json())
```

## 状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 500 | 服务器内部错误 |
