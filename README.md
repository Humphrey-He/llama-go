# Llama-Go

基于 Go + Python 的轻量级 LLM 推理服务，实现 KV Cache 优化

## 特性

- ✅ Go 高并发 API 服务
- ✅ Python 模型推理
- ✅ KV Cache 缓存优化
- ✅ 多会话管理
- ✅ 滑动窗口缓存
- ✅ 并发安全设计
- ✅ Docker 容器化部署

## 快速开始

```bash
# 启动服务
docker-compose up --build

# 测试接口
curl -X POST http://localhost:8080/api/generate \
  -H "Content-Type: application/json" \
  -d '{"prompt": "Hello, how are you?"}'
```

## 项目结构

```
llama-go/
├── api/                    # API 路由层
├── internal/
│   ├── backend/           # Python 服务客户端
│   ├── kvcache/           # KV Cache 管理
│   ├── config/            # 配置管理
│   └── logger/            # 日志系统
├── python/
│   ├── kv_server.py       # Python 推理服务
│   └── requirements.txt
├── docs/                  # 项目文档
├── main.go               # 主入口
└── docker-compose.yml
```

## 文档

- [项目说明](docs/01-项目说明.md)
- [技术选型](docs/02-技术选型.md)
- [架构设计](docs/03-架构与系统设计.md)
- [API 文档](docs/API文档.md)
- [部署文档](docs/部署文档.md)

## 技术栈

- **Go**: Gin Web Framework
- **Python**: FastAPI + PyTorch + Transformers
- **模型**: TinyLlama-1.1B
- **部署**: Docker Compose

## 开发

```bash
# 运行测试
go test ./...

# 竞态检测
go test -race ./internal/kvcache
```

## License

MIT
