# issue2md 工具需求规范

## 1. 产品概述

### 1.1 项目简介

**issue2md** 是一个命令行工具，用于将 GitHub Issue、Pull Request 和 Discussion 转换为 GitHub Flavored Markdown (GFM) 格式。

### 1.2 核心功能

- 输入单个 GitHub URL（Issue / PR / Discussion）
- 自动识别类型（通过 URL 路径解析：`/issues/`、`/pull/`、`/discussions/`）
- 输出标题、作者、创建时间、状态、主楼内容、所有评论

### 1.3 目标用户

- 开发者需要将 GitHub 内容导出为本地 Markdown 文件
- 文档编写者需要备份 Issue/PR 内容
- 团队需要离线查看讨论内容

---

## 2. 命令行接口

### 2.1 使用方式

```
issue2md <github-url> [flags]
```

### 2.2 参数说明

| 参数 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `<github-url>` | string | 是 | GitHub Issue/PR/Discussion 的完整 URL |

### 2.3 Flags

| Flag | 短格式 | 默认值 | 说明 |
|------|--------|--------|------|
| `--output` | `-o` | (stdout) | 指定输出文件路径 |
| `--include-reactions` | - | false | 包含 reactions 统计信息 |
| `--user-links` | - | false | 将用户名渲染为 GitHub 主页链接 `[username](https://github.com/username)` |
| `--token` | - | - | GitHub Personal Access Token（优先级高于环境变量） |
| `--version` | `-v` | false | 显示版本信息 |
| `--help` | `-h` | false | 显示帮助信息 |

### 2.4 环境变量

| 变量名 | 必需 | 说明 |
|--------|------|------|
| `GITHUB_TOKEN` | 否 | GitHub Personal Access Token，用于提升 API rate limit（公有仓库需要） |

> 注意：`--token` 参数优先级高于 `GITHUB_TOKEN` 环境变量

### 2.5 退出码

| 退出码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1 | 通用错误（参数错误、API 错误等） |
| 2 | URL 解析错误（无效的 GitHub URL） |

---

## 3. 输出格式规范

### 3.1 Issue 输出格式

```markdown
# <title>

**Author:** @username | **Created:** <YYYY-MM-DD HH:MM:SS> | **Status:** Open/Closed

---

<body content>

---

## Comments (<count>)

### Comment #1 | <YYYY-MM-DD>

<comment content>

### Comment #2 | <YYYY-MM-DD>

<comment content>
```

### 3.2 Pull Request 输出格式

```markdown
# <title>

**Author:** @username | **Created:** <YYYY-MM-DD HH:MM:SS> | **Status:** Open/Closed/Merged

---

<body content>

---

## Timeline (<count>)

### Event | <YYYY-MM-DD>

<event content>
```

### 3.3 Discussion 输出格式

```markdown
# <title>

**Author:** @username | **Created:** <YYYY-MM-DD HH:MM:SS> | **Category:** <category>

---

<body content>

---

## Answers (<count>)

### Answer by @username | <YYYY-MM-DD>

<answer content>
```

### 3.4 内容处理规则

| 内容类型 | 处理方式 |
|----------|----------|
| `@username` | 保留原样（除非指定 `--link-user`） |
| emoji | 保留原始 emoji |
| 图片/附件 | 保留原始链接，不下载 |
| 代码块 | 保留原始格式，使用 ``` 包裹 |
| HTML | 转换为纯文本或保留安全的 HTML |

### 3.5 `--user-links` 行为

当指定 `--user-links` 时：
- `@username` 转换为 `[@username](https://github.com/username)`

---

## 4. 错误处理规范

### 4.1 错误输出

所有错误信息输出到 **stderr**，格式如下：

```
error: <error message>
```

### 4.2 错误类型

| 错误场景 | 错误信息 | 退出码 |
|----------|----------|--------|
| 无效 URL | `invalid GitHub URL: <url>` | 2 |
| 不支持的类型 | `unsupported URL type: <type>` | 2 |
| API 错误 | `GitHub API error: <details>` | 1 |
| 网络错误 | `network error: <details>` | 1 |
| 文件写入错误 | `failed to write file: <details>` | 1 |

---

## 5. 技术规范

### 5.1 项目结构

```
issue2md/
├── cmd/
│   └── issue2md/
│       └── main.go           # 入口文件
├── internal/
│   ├── github/
│   │   └── client.go         # GitHub API 客户端
│   ├── convert/
│   │   └── markdown.go        # Markdown 转换逻辑
│   └── parser/
│       └── url.go            # URL 解析
├── spec.md                   # 本规范文档
├── go.mod
├── go.sum
└── Makefile
```

### 5.2 依赖

- Go >= 1.24
- 标准库：`net/http`、`fmt`、`flag`、`os`、`time`、`strings`
- 可选：如有需要可使用 `github.com/google/go-github/vXX` 或自实现

### 5.3 版本信息

```
issue2md version <major>.<minor>.<patch>
```

版本号通过以下方式注入：
- 编译时使用 `-ldflags` 注入 Git tag 信息

---

## 6. 验收标准

### 6.1 功能验收

- [ ] 命令行解析正确识别所有 flags
- [ ] URL 解析正确识别 Issue/PR/Discussion
- [ ] 输出包含：标题、作者、创建时间、状态、主楼内容、评论
- [ ] `--output` 参数正确写入文件
- [ ] `--include-reactions` 正确包含 reactions 统计
- [ ] `--user-links` 正确转换用户名
- [ ] `--token` 正确覆盖环境变量
- [ ] 错误信息输出到 stderr
- [ ] 错误时返回非零退出码

### 6.2 测试验收

- [ ] 表格驱动测试覆盖核心转换逻辑
- [ ] Mock GitHub API 响应进行测试

### 6.3 构建验收

- [ ] `make build` 成功构建二进制文件
- [ ] `make test` 成功运行所有测试
