# API 设计草稿

## 1. 模块划分

```
internal/
├── parser/       # URL 解析
├── github/       # GitHub API 客户端
└── convert/      # Markdown 转换
```

---

## 2. Parser 模块

### 2.1 URL 解析器接口

```go
// parser/url.go

type URLType string

const (
    URLTypeIssue      URLType = "issue"
    URLTypePullRequest URLType = "pull"
    URLTypeDiscussion URLType = "discussion"
    URLTypeUnknown    URLType = "unknown"
)

type ParsedURL struct {
    Type       URLType
    Owner      string
    Repo       string
    Number     int
    RawURL     string
}

// Parser 接口
type Parser interface {
    Parse(rawURL string) (*ParsedURL, error)
}

// New 工厂函数
func New() Parser
```

### 2.2 URL 解析规则

| URL 路径示例 | Type | Owner | Repo | Number |
|-------------|------|-------|------|--------|
| `https://github.com/owner/repo/issues/123` | issue | owner | repo | 123 |
| `https://github.com/owner/repo/pull/456` | pull | owner | repo | 456 |
| `https://github.com/owner/repo/discussions/789` | discussion | owner | repo | 789 |

---

## 3. GitHub 模块

### 3.1 GitHub 客户端接口

```go
// internal/github/client.go

type IssueData struct {
    Title     string
    Body      string
    Author    string
    CreatedAt time.Time
    State     string
    Comments  []Comment
    Reactions []Reaction
}

type Comment struct {
    ID        int
    Body      string
    Author    string
    CreatedAt time.Time
    Reactions []Reaction
}

type Reaction struct {
    Name  string // "+1", "-1", "heart", "hooray", "laugh", "confused"
    Count int
}

type PullRequestData struct {
    Title     string
    Body      string
    Author    string
    CreatedAt time.Time
    State     string
    Merged    bool
    Timeline  []TimelineEvent
    Comments  []Comment
    Reactions []Reaction
}

type TimelineEvent struct {
    Type    string // "Event", "IssueComment", "PullRequestReview"
    Actor   string
    Body    string
    CreatedAt time.Time
}

type DiscussionData struct {
    Title     string
    Body      string
    Author    string
    CreatedAt time.Time
    Category  string
    Answers   []Answer
    Comments  []Comment
    Reactions []Reaction
}

type Answer struct {
    ID        int
    Body      string
    Author    string
    CreatedAt time.Time
    IsAnswer  bool
    Reactions []Reaction
}

// GitHuber 接口
type GitHuber interface {
    FetchIssue(owner, repo string, number int) (*IssueData, error)
    FetchPullRequest(owner, repo string, number int) (*PullRequestData, error)
    FetchDiscussion(owner, repo string, number int) (*DiscussionData, error)
}

// Option 配置选项
type Option func(*client)

func WithToken(token string) Option
func WithIncludeReactions(include bool) Option
func WithRateLimitHandler(handler RateLimitHandler) Option

// New 工厂函数
func New(opts ...Option) GitHuber
```

### 3.2 API 端点

| 资源 | 端点 | 方法 |
|------|------|------|
| Issue | `/repos/{owner}/{repo}/issues/{issue_number}` | GET |
| Issue Comments | `/repos/{owner}/{repo}/issues/{issue_number}/comments` | GET |
| Issue Reactions | `/repos/{owner}/{repo}/issues/{issue_number}/reactions` | GET |
| Pull Request | `/repos/{owner}/{repo}/pulls/{pull_number}` | GET |
| PR Timeline | `/repos/{owner}/{repo}/pulls/{pull_number}/timeline` | GET |
| PR Comments | `/repos/{owner}/{repo}/issues/{pull_number}/comments` | GET |
| Discussion | `/repos/{owner}/{repo}/discussions/{discussion_number}` | GET |
| Discussion Comments | `/repos/{owner}/{repo}/discussions/{discussion_number}/comments` | GET |

---

## 4. Convert 模块

### 4.1 Markdown 转换器接口

```go
// internal/convert/markdown.go

type DataItem interface {
    // Marker interface for IssueData, PullRequestData, DiscussionData
}

type MarkdownOption func(*converter)

func WithLinkUser(enabled bool) MarkdownOption
func WithPretty(enabled bool) MarkdownOption

// Converter 接口
type Converter interface {
    Convert(data DataItem) (string, error)
}

// New 工厂函数
func New(opts ...MarkdownOption) Converter
```

### 4.2 转换方法

| 输入类型 | 输出方法 |
|----------|----------|
| `IssueData` | `convertIssue(*IssueData) string` |
| `PullRequestData` | `convertPullRequest(*PullRequestData) string` |
| `DiscussionData` | `convertDiscussion(*DiscussionData) string` |

---

## 5. 命令行入口

### 5.1 Main 函数流程

```go
// cmd/issue2md/main.go

func main() {
    // 1. 解析命令行参数
    // 2. 解析 GitHub URL
    // 3. 创建 GitHub 客户端
    // 4. 获取数据
    // 5. 转换为 Markdown
    // 6. 输出
}
```

### 5.2 Flag 定义

```go
var (
    output    string
    includeReactions bool
    userLinks bool
    token     string
    version   bool
)
```

---

## 6. 错误定义

```go
// internal/github/errors.go

var (
    ErrInvalidURL     = errors.New("invalid GitHub URL")
    ErrUnsupportedURL = errors.New("unsupported URL type")
    ErrNotFound       = errors.New("resource not found")
    ErrRateLimited    = errors.New("API rate limit exceeded")
)
```
