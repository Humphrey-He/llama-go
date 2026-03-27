# Contributing Guide

## 开发流程

### 分支策略
- `main` - 主分支，发布版本
- `dev` - 开发分支，集成点
- `feature/*` - 功能分支，从 dev 创建

### 提交规范
```
<type>: <subject>

<body>

<footer>
```

**Type 类型**：
- `feat` - 新功能
- `fix` - 修复
- `refactor` - 重构
- `test` - 测试
- `docs` - 文档
- `chore` - 构建/工具

### 代码规范

**格式检查**：
```bash
go fmt ./...
```

**Lint 检查**：
```bash
golangci-lint run ./...
```

**测试**：
```bash
go test -v -race ./...
```

## PR 流程

1. Fork 项目
2. 创建功能分支：`git checkout -b feature/xxx`
3. 提交代码：`git commit -m "feat: xxx"`
4. 推送分支：`git push origin feature/xxx`
5. 创建 Pull Request
6. 等待 review 和 CI 通过

## 代码质量要求

- 单元测试覆盖率 ≥ 70%
- 无 lint 错误
- 代码注释完整
- 文档同步更新
