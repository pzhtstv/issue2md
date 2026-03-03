# issue2md 任务列表

## 阶段说明
- **[P]**: 可并行执行的任务
- **依赖**: 前置任务标记

---

## Phase 1: Foundation (数据结构定义)

### 1.1 Parser 模块 - 数据结构与接口

| # | 任务 | 文件 | 依赖 |
|---|------|------|------|
| 1.1.1 | 创建 parser/types.go - 定义 URLType 常量和 ParsedURL 结构体 | `internal/parser/types.go` | - |
| 1.1.2 [P] | 创建 parser/errors.go - 定义解析错误 | `internal/parser/errors.go` | - |
| 1.1.3 [P] | 创建 parser/parser.go - 定义 Parser 接口 | `internal/parser/parser.go` | 1.1.1, 1.1.2 |

### 1.2 GitHub 模块 - 数据结构

| # | 任务 | 文件 | 依赖 |
|---|------|------|------|
| 1.2.1 | 创建 github/types.go - 定义 Author, Reaction, Comment, IssueData, PullRequestData, DiscussionData 结构体 | `internal/github/types.go` | - |
| 1.2.2 [P] | 创建 github/errors.go - 定义 GitHub API 错误 | `internal/github/errors.go` | - |

### 1.3 Convert 模块 - 数据结构

| # | 任务 | 文件 | 依赖 |
|---|------|------|------|
| 1.3.1 | 创建 convert/converter.go - 定义 Converter 接口和选项函数 | `internal/convert/converter.go` | 1.2.1 |

---

## Phase 2: GitHub Fetcher (API 交互逻辑，TDD)

### 2.1 Parser 实现

| # | 任务 | 文件 | 依赖 |
|---|------|------|------|
| 2.1.1 | 编写 parser/url_test.go - 表格驱动测试 (TDD: 先写失败测试) | `internal/parser/url_test.go` | 1.1.1, 1.1.2, 1.1.3 |
| 2.1.2 | 实现 parser/url.go - URL 解析逻辑 | `internal/parser/url.go` | 2.1.1 |

### 2.2 GitHub Client 实现

| # | 任务 | 文件 | 依赖 |
|---|------|------|------|
| 2.2.1 | 创建 github/githuber.go - 定义 GitHuber 接口 | `internal/github/githuber.go` | 1.2.1 |
| 2.2.2 | 编写 github/client_test.go - Mock HTTP 响应测试 | `internal/github/client_test.go` | 2.2.1 |
| 2.2.3 | 实现 github/client.go - GitHub GraphQL 客户端 | `internal/github/client.go` | 2.2.2 |

---

## Phase 3: Markdown Converter (转换逻辑，TDD)

### 3.1 Issue 转换

| # | 任务 | 文件 | 依赖 |
|---|------|------|------|
| 3.1.1 | 编写 convert/issue_test.go - Issue 转换表格驱动测试 | `internal/convert/issue_test.go` | 1.3.1, 1.2.1 |
| 3.1.2 | 实现 convert/issue.go - Issue 转 Markdown | `internal/convert/issue.go` | 3.1.1 |

### 3.2 Pull Request 转换

| # | 任务 | 文件 | 依赖 |
|---|------|------|------|
| 3.2.1 | 编写 convert/pr_test.go - PR 转换表格驱动测试 | `internal/convert/pr_test.go` | 1.3.1, 1.2.1 |
| 3.2.2 | 实现 convert/pr.go - PR 转 Markdown | `internal/convert/pr.go` | 3.2.1 |

### 3.3 Discussion 转换

| # | 任务 | 文件 | 依赖 |
|---|------|------|------|
| 3.3.1 | 编写 convert/discussion_test.go - Discussion 转换表格驱动测试 | `internal/convert/discussion_test.go` | 1.3.1, 1.2.1 |
| 3.3.2 | 实现 convert/discussion.go - Discussion 转 Markdown | `internal/convert/discussion.go` | 3.3.1 |

### 3.4 Converter 入口

| # | 任务 | 文件 | 依赖 |
|---|------|------|------|
| 3.4.1 | 实现 convert/markdown.go - 统一转换入口，组合所有转换器 | `internal/convert/markdown.go` | 3.1.2, 3.2.2, 3.3.2 |

---

## Phase 4: CLI Assembly (命令行入口集成)

### 4.1 命令行入口

| # | 任务 | 文件 | 依赖 |
|---|------|------|------|
| 4.1.1 | 创建 Makefile - build, test, run 命令 | `Makefile` | - |
| 4.1.2 | 实现 cmd/issue2md/main.go - CLI 参数解析和流程编排 | `cmd/issue2md/main.go` | 2.1.2, 2.2.3, 3.4.1 |

### 4.2 端到端验证

| # | 任务 | 文件 | 依赖 |
|---|------|------|------|
| 4.2.1 | 构建验证: `make build` | - | 4.1.1, 4.1.2 |
| 4.2.2 | 测试验证: `make test` | - | 4.2.1 |
| 4.2.3 | 功能验证: 测试真实 GitHub Issue | - | 4.2.1 |

---

## 任务依赖图

```
Phase 1: Foundation
├── 1.1.1 parser/types.go
├── 1.1.2 parser/errors.go
├── 1.1.3 parser/parser.go
├── 1.2.1 github/types.go
├── 1.2.2 github/errors.go
└── 1.3.1 convert/converter.go

Phase 2: GitHub Fetcher
├── 2.1.1 parser/url_test.go ─────┐
├── 2.1.2 parser/url.go ◄─────────┘
├── 2.2.1 github/githuber.go ─────┐
├── 2.2.2 github/client_test.go ──┤
└── 2.2.3 github/client.go ◄──────┘

Phase 3: Markdown Converter
├── 3.1.1 convert/issue_test.go ──┐
├── 3.1.2 convert/issue.go ◄──────┤
├── 3.2.1 convert/pr_test.go ─────┤
├── 3.2.2 convert/pr.go ◄─────────┤
├── 3.3.1 convert/discussion_test.go
├── 3.3.2 convert/discussion.go
└── 3.4.1 convert/markdown.go

Phase 4: CLI Assembly
├── 4.1.1 Makefile ───────────────┐
└── 4.1.2 cmd/issue2md/main.go ◄─┘
    │
    ├── 4.2.1 make build
    ├── 4.2.2 make test
    └── 4.2.3 功能验证
```

---

## 验收标准

每个任务完成后需满足：
1. **测试通过**: `go test ./...` 通过
2. **代码编译**: `go build ./...` 成功
3. **无 lint 错误**: `go vet ./...` 无警告
