# Docker 部署指南

## 快速开始

### 1. 准备模型文件

```bash
mkdir -p models
# 下载 GGUF 模型到 models/ 目录
# 或使用 huggingface-hub:
pip install huggingface-hub
huggingface-cli download TheBloke/Dolphin-3.0-Llama-3.1-8B-GGUF \
  dolphin-3.0-llama-3.1-8b.Q4_K_M.gguf \
  --local-dir ./models
```

### 2. 启动所有服务

```bash
docker-compose up -d
```

### 3. 验证服务

```bash
# 检查后端健康状态
curl http://localhost:8000/health

# 检查 API 健康状态
curl http://localhost:8080/healthz

# 查看 Prometheus 指标
curl http://localhost:9090/api/v1/targets

# 访问 Grafana
# 浏览器打开: http://localhost:3000
# 用户名: admin, 密码: admin
```

### 4. 测试推理

```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "tinyllama-chat",
    "messages": [{"role": "user", "content": "Hello"}],
    "stream": false
  }'
```

## 服务端口

| 服务 | 端口 | 用途 |
|------|------|------|
| llama-go API | 8080 | 推理 API |
| vLLM 后端 | 8000 | 模型推理 |
| Prometheus | 9090 | 指标收集 |
| Grafana | 3000 | 可视化监控 |

## 常见问题

**Q: 模型加载很慢?**
A: 首次加载需要时间，后续请求会快速响应

**Q: 如何修改模型?**
A: 编辑 docker-compose.yml 中的 MODEL_PATH，重启服务

**Q: 如何查看日志?**
A: `docker-compose logs -f llama-backend`
