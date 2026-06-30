# llama-go 插件化改造完成

## ✅ 改造成果

已成功将 llama-go 改造为"内核插件化，接口标准化"的云原生推理网关。

### 核心变化

**Before (硬编码):**
```go
ollamaBackend := backend.NewOllamaBackend("http://localhost:11434")
registry.Register("llama3.1:8b", ollamaBackend)
```

**After (插件化):**
```yaml
plugins:
  llama-cpp:
    enabled: true
    base_url: "http://localhost:8000"
    models: ["tinyllama-chat"]
```

### 新增文件

```
internal/plugin/
├── interface.go           # 标准插件接口
├── manager.go             # 插件管理器
├── llamacpp/plugin.go     # llama-cpp 插件
├── ollama/plugin.go       # Ollama 插件
└── vllm/plugin.go         # vLLM 插件

internal/config/
└── loader.go              # YAML 配置加载

config.yaml                # 配置文件
docs/
├── 插件化架构设计.md
└── 改造总结.md
```

## 🚀 快速开始

### 1. 启动后端

```bash
# Windows
.\venv\Scripts\Activate.ps1
python python/start_llama_cpp.py
```

### 2. 启动网关

```bash
# 方式 1: 使用脚本
.\start.ps1

# 方式 2: 直接运行
go run cmd/inference/main.go
```

### 3. 测试

```bash
curl http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{"model":"tinyllama-chat","messages":[{"role":"user","content":"你好"}]}'
```

## 📋 架构优势

1. **插件化**: 后端完全解耦，新增插件零侵入
2. **配置驱动**: 修改 YAML 即可切换后端
3. **标准接口**: 统一 Plugin 接口，降低集成成本
4. **生命周期**: Init/Health/Close 完整管理
5. **云原生**: 适合 Docker/K8s 部署

## 🔧 切换后端

编辑 `config.yaml`:

```yaml
plugins:
  llama-cpp:
    enabled: true    # 使用 llama-cpp
  ollama:
    enabled: false   # 禁用 Ollama
  vllm:
    enabled: false   # 禁用 vLLM
```

重启服务即可，无需修改代码。

## 📚 文档

- [插件化架构设计](docs/插件化架构设计.md)
- [改造总结](docs/改造总结.md)
