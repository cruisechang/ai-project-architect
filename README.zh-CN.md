# apa - AI Project Architect

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![CLI](https://img.shields.io/badge/Type-CLI-111111)](#命令总览)
[![Skills](https://img.shields.io/badge/Repo%20Skills-apa--*-2f855a)](#commands--skills)

`apa` 是一个 Go CLI，用来把产品 idea 快速变成可开发的项目起点。

它会生成项目 context、设计文档、可运行骨架，以及后续可持续迭代的 AI 工作流，让你从想法到实现更快落地。

适合用于你想用一个工具同时完成项目 bootstrap、文档起稿、初始代码生成，以及衔接 AI 持续交付流程的场景。

## Why `apa`

- 你希望新项目一开始就同时具备文档和代码，而不是只有空目录
- 你希望 AI 从结构化的 repo 状态开始工作，而不是临时拼凑提示词
- 你希望用 repo-local `apa-*` skills 和原生命令如 `make test` 推进交付

## Installation

```bash
go build -o apa .
```

可选：

```bash
cp ./apa /usr/local/bin/apa
```

## Flow

```text
idea
  -> apa init
  -> docs + scaffold + context
  -> apa list-skills
  -> apa iterate
  -> agent implements
  -> make test
  -> repeat
```

## 语言

| 语言 | README |
|---|---|
| English | [README.md](README.md) |
| 简体中文 | [README.zh-CN.md](README.zh-CN.md) |
| 繁體中文 | [README.zh-TW.md](README.zh-TW.md) |
| Deutsch | [README.de.md](README.de.md) |
| Français | [README.fr.md](README.fr.md) |
| Español | [README.es.md](README.es.md) |
| 日本語 | [README.ja.md](README.ja.md) |
| 한국어 | [README.ko.md](README.ko.md) |

## `apa` 可以做什么

- 从自然语言产品想法 bootstrap 新项目
- 先自动推断技术栈，再允许你通过 flags 覆盖
- 生成 PRD、SPEC、ARCHITECTURE、API、DB Schema、实现计划等文档
- 创建可运行的起始代码、测试结构、Makefile 和 agent 配置
- 通过 `apa iterate` 输出结构化提示词，让 AI agent 在可控流程里持续推进

## 推荐工作流

```bash
# 1. 先构建 apa
go build -o apa .

# 2. 在目标 repo 外创建新项目
./apa init --idea "SaaS 报表平台" --name report-platform --path ~/projects

# 3. 进入生成出的 repo 开始迭代
cd ~/projects/report-platform
./apa list-skills
./apa iterate
make test
```

核心循环：

1. 用 `apa init` 创建第一版项目。
2. 用 repo-local 的 `apa-*` skills 引导实现工作。
3. 执行 `apa iterate`，让 agent 实现，再用 `make test` 验证。
4. 重复直到项目可交付。

## 快速示例

```bash
./apa init \
  --idea "内部客服仪表盘，支持工单搜索和 AI 摘要" \
  --name support-hub \
  --path ~/projects \
  --backend go \
  --frontend next \
  --agent codex

cd ~/projects/support-hub
./apa list-skills
./apa iterate > prompt.md
make test
```

## 命令总览

| 命令 | 用途 |
|---|---|
| `apa init` | 从 idea 创建新项目 |
| `apa iterate` | 生成后续持续交付用的 AI 提示词 |
| `apa list-skills` | 列出可用的 repo-local skills |
| `apa doctor` | 检查本地环境和 skills 路径 |
| `apa version` | 显示构建版本信息 |

完整选项请使用 `apa <command> --help`。

## Commands & Skills

当前 CLI 命令：

- `init`
- `iterate`
- `list-skills`
- `doctor`
- `version`

当前 repo-local skills：

- `apa-catalog`
- `apa-debug`
- `apa-devops`
- `apa-docs`
- `apa-feature`
- `apa-integration`
- `apa-review`
- `apa-tdd`

## `apa init`

`apa init` 是核心命令，负责一次完成首次 bootstrap。

典型流程：

```text
[1/4] infer tech stack from idea
[2/4] generate code scaffold
[3/4] generate design docs
[4/4] done
```

常见用法：

```bash
# 交互模式
./apa init

# 非交互模式
./apa init --idea "在线点餐平台" --name food-platform --path ~/projects

# 仅预览，不写文件
./apa init --dry-run
```

常用 flags：

| Flag | 说明 |
|---|---|
| `--idea` | 用于推断技术栈的产品 idea |
| `--name` | 项目名称 |
| `--path` | 项目创建所在的父目录 |
| `--agent` | `codex` 或 `claude-code` |
| `--backend` | `go`、`python`、`node`、`none` |
| `--frontend` | `react`、`next`、`nuxt`、`vue`、`pure-typescript`、`none` |
| `--skills` | 要复制进新项目的 repo-local skills，逗号分隔 |
| `--skills-path` | skills 来源目录，默认为本 repo 的 `./skills` |
| `--force` | 备份现有目录后重建 |
| `--dry-run` | 显示计划输出但不写入 |

常见产物：

```text
.architect/context.json
docs/
backend/ frontend/ tests/
Makefile
CLAUDE.md or PROMPT.md
agents/ skills/
```

## `apa iterate`

`apa iterate` 会读取当前 repo 状态，并输出一段给 AI agent 的结构化提示词。

它可以在实现前、开发中，或回归修复后使用，帮助 agent 按照现有文档、任务和约束持续推进到交付。

```bash
./apa iterate
./apa iterate --root ~/projects/report-platform
./apa iterate > prompt.md
```

## Repo-Local Skills

本仓库的 skill 命名遵循 `apa-*` 系列。

当前示例：

- `apa-catalog`
- `apa-debug`
- `apa-devops`
- `apa-docs`
- `apa-feature`
- `apa-integration`
- `apa-review`
- `apa-tdd`

可以用下面命令查看：

```bash
./apa list-skills
```

## 开发与测试

```bash
go test ./...
make test
make build TARGET=mac
make build TARGET=linux
make build TARGET=windows
make build-all
```

构建 metadata 示例：

```bash
VERSION=v1.2.3 COMMIT=abc1234 BUILD_DATE=2026-03-19T10:00:00Z ./build.sh mac
```

## Repository Layout

```text
ai-project-architect/
├── apa/           # CLI commands
├── internal/      # generation, planner, runtime, config, output
├── skills/        # repo-local apa-* skills
├── templates/     # embedded templates
├── scripts/
├── build.sh
└── Makefile
```
