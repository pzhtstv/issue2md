# issue2md 技术实现方案

## 1. 技术上下文总结

### 1.1 技术栈

| 组件 | 技术选型 | 说明 |
|------|----------|------|
| 语言 | Go >= 1.21.0 | 项目要求 |
| CLI 框架 | Go 标准库 `flag` | 遵循"标准库优先"原则 |
| GitHub API | `google/go-github` v4 (GraphQL) | 结合 GraphQL API |
| Markdown 处理 | 标准库 + 自实现 | 尽量不使用第三方库 |
| 测试 | Go testing + 表格驱动 | 遵循 TDD 铁律 |

### 1.2 架构概览

```
┌─────────────────────────────────────────────────────────┐
│                      cmd/issue2md                       │
│                      (main.go)                          │
└─────────────────────┬───────────────────────────────────┘
                      │
        ┌─────────────┼─────────────┐
        ▼             ▼             ▼
   ┌─────────┐  ┌─────────┐  ┌─────────┐
   │ parser/ │  │ github/ │  │convert/ │
   │         │  │         │  │         │
   └─────────┘  └─────────┘  └─────────┘
```

---

## 2. "合宪性"审查

### 2.1 第一条：简单性原则

| 条款 | 审查结果 | 实施措施 |
|------|----------|----------|
| 1.1 YAGNI | ✅ 符合 | 只实现 spec.md 中的 6 个 flags 功能 |
| 1.2 标准库优先 | ✅ 符合 | CLI 使用 `flag`，不使用 cobra/clap |
| 1.3 反过度工程 | ✅ 符合 | 扁平化包结构，避免过度接口抽象 |

### 2.2 第二条：测试先行铁律

| 条款 | 审查结果 | 实施措施 |
|------|----------|----------|
| 2.1 TDD 循环 | ✅ 符合 | 每个功能先写失败测试，再实现 |
| 2.2 表格驱动 | ✅ 符合 | URL 解析、Markdown 转换均采用表格驱动 |
| 2.3 拒绝 Mocks | ✅ 符合 | 使用 `httptest` 模拟 HTTP 响应 |

### 2.3 第三条：明确性原则

| 条款 | 审查结果 | 实施措施 |
|------|----------|----------|
| 3.1 错误处理 | ✅ 符合 | 所有错误使用 `fmt.Errorf("...: %w", err)` 包装 |
| 3.2 无全局变量 | ✅ 符合 | 依赖通过参数或结构体注入 |

---

## 3. 项目结构细化

### 3.1 目录结构

```
issue2md/
├── cmd/
│   └── issue2md/
│       └── main.go              # 入口文件，负责 CLI 解析和流程编排
├── internal/
│   ├── github/
│   │   ├── client.go            # GitHub GraphQL 客户端
│   │   ├── types.go             # 数据结构定义
│   │   └── errors.go            # 错误定义
│   ├── convert/
│   │   └── markdown.go          # Markdown 转换逻辑
│   └── parser/
│       ├── url.go               # URL 解析逻辑
│       └── url_test.go          # 表格驱动测试
├── spec/
│   ├── spec.md                  # 需求规范
│   ├── api-sketch.md            # API 设计草稿
│   └── plan.md                  # 本技术方案
├── go.mod
└── Makefile
```

### 3.2 包职责与依赖

| 包 | 职责 | 依赖 |
|---|------|------|
| `cmd/issue2md` | CLI 解析、参数校验、流程编排 | parser, github, convert |
| `internal/github` | GitHub API 调用（GraphQL） | 标准库 |
| `internal/convert` | 数据结构转换为 Markdown 字符串 | 标准库 |
| `internal/parser` | GitHub URL 解析 | 标准库 |

**依赖方向**：cmd → (parser | github | convert)，无循环依赖

---

## 4. 核心数据结构

### 4.1 Parser 模块

```go
// internal/parser/url.go

type URLType string

const (
    TypeIssue      URLType = "issue"
    TypePullRequest URLType = "pull"
    TypeDiscussion URLType = "discussion"
    TypeUnknown    URLType = "unknown"
)

type ParsedURL struct {
    Type   URLType
    Owner  string
    Repo   string
    Number int
}
```

### 4.2 GitHub 模块

```go
// internal/github/types.go

type Author struct {
    Login string
    URL   string // 当 --user-links 时使用
}

type Reaction struct {
    Content string // "+1", "-1", "heart", "hooray", "laugh", "confused"
    Count   int
}

type Comment struct {
    ID        int
    Body      string
    Author    Author
    CreatedAt time.Time
    Reactions []Reaction // 当 --include-reactions 时填充
}

type IssueData struct {
    Title     string
    Body      string
    Author    Author
    CreatedAt time.Time
    State     string // "open" | "closed"
    Number    int
    Comments  []Comment
    Reactions []Reaction // 当 --include-reactions 时填充
}

type PullRequestData struct {
    Title     string
    Body      string
    Author    Author
    CreatedAt time.Time
    State     string // "open" | "closed" | "merged"
    Number    int
    Comments  []Comment
    Reactions []Reaction
}

type DiscussionData struct {
    Title     string
    Body      string
    Author    Author
    CreatedAt time.Time
    Category  string
    Number    int
    Answers   []Comment
    Comments  []Comment
    Reactions []Reaction
}
```

### 4.3 Convert 模块

```go
// internal/convert/markdown.go

type ConverterOption func(*converter)

type Converter struct {
    userLinks bool
}

func WithUserLinks(enabled bool) ConverterOption
func New(opts ...ConverterOption) *Converter

func (c *Converter) ConvertIssue(data *github.IssueData) (string, error)
func (c *Converter) ConvertPullRequest(data *github.PullRequestData) (string, error)
func (c *Converter) ConvertDiscussion(data *github.DiscussionData) (string, error)
```

---

## 5. 接口设计

### 5.1 Parser 接口

```go
// internal/parser/parser.go

type Parser interface {
    Parse(rawURL string) (*ParsedURL, error)
}
```

### 5.2 GitHuber 接口

```go
// internal/github/githuber.go

type GitHuber interface {
    FetchIssue(owner, repo string, number int) (*IssueData, error)
    FetchPullRequest(owner, repo string, number int) (*PullRequestData, error)
    FetchDiscussion(owner, repo string, number int) (*DiscussionData, error)
}
```

### 5.3 Converter 接口

```go
// internal/convert/converter.go

type Converter interface {
    Convert(data any) (string, error)
}
```

---

## 6. 实现计划

### Phase 1: 项目初始化

- [ ] 创建 Makefile（build, test, run 命令）
- [ ] 添加 `google/go-github` 依赖

### Phase 2: Parser 模块

- [ ] 实现 URL 解析逻辑
- [ ] 编写表格驱动测试（URL 解析）

### Phase 3: GitHub 模块

- [ ] 实现 GraphQL 客户端
- [ ] 实现 Issue/PR/Discussion 获取方法
- [ ] 编写集成测试（使用 httptest）

### Phase 4: Convert 模块

- [ ] 实现 Markdown 转换逻辑
- [ ] 实现 --user-links 用户名链接转换
- [ ] 编写表格驱动测试

### Phase 5: CLI 集成

- [ ] 实现 main.go 命令行解析
- [ ] 实现文件输出逻辑

### Phase 6: 端到端测试

- [ ] 使用真实 GitHub Issue 测试

---

## 7. 验证方案

### 7.1 构建验证

```bash
make build  # 成功构建 issue2md 二进制
```

### 7.2 测试验证

```bash
make test   # 运行所有测试
```

### 7.3 功能验证

```bash
# Issue 测试
./issue2md https://github.com/golang/go/issues/123 -o issue.md

# PR 测试
./issue2md https://github.com/golang/go/pull/456 -o pr.md

# Discussion 测试
./issue2md https://github.com/golang/go/discussions/789 -o discussion.md

# 带 token 测试
./issue2md https://github.com/owner/repo/issues/1 --token ghp_xxx

# 带用户链接
./issue2md https://github.com/owner/repo/issues/1 --user-links
```

### 7.4 错误验证

```bash
# 无效 URL
./issue2md "invalid-url"
# 预期：error: invalid GitHub URL: invalid-url (exit 2)

# 不支持的类型
./issue2md https://github.com/owner/repo/blob/main/README.md
# 预期：error: unsupported URL type (exit 2)
```
