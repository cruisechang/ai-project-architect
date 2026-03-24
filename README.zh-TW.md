# apa - AI Project Architect

## 語言

| 語言 | README |
|---|---|
| English | [README.md](README.md) |
| 简体中文 | [README.zh-CN.md](README.zh-CN.md) |
| 繁體中文 | [README.zh-TW.md](README.zh-TW.md) |
| Deutsch | [README.de.md](README.de.md) |
| Français | [README.fr.md](README.fr.md) |
| Español | [README.es.md](README.es.md) |
| 日本語 | [README.ja.md](README.ja.md) |
| 한국어 | [README.ko.md](README.ko.md) |

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![CLI](https://img.shields.io/badge/Type-CLI-111111)](#指令總覽)
[![Skills](https://img.shields.io/badge/Repo%20Skills-apa--*-2f855a)](#commands--skills)

`apa` 是一個 Go CLI，目標是把產品 idea 快速變成可開發的專案起點。

它會產出專案 context、設計文件、可執行骨架，以及後續可持續迭代的 AI 工作流程，讓你從想法到實作更快落地。

適合用在你想用單一工具完成專案 bootstrap、文件起稿、初始程式碼產生，以及銜接 AI 持續交付流程的情境。

## Why `apa`

- 你希望新專案一開始就同時有文件與程式碼，不是只有空目錄
- 你希望 AI 從有結構的 repo 狀態開始工作，而不是臨時拼湊提示詞
- 你希望用 repo-local `apa-*` skills 與原生指令如 `make test` 推進交付

## Installation

```bash
go build -o apa .
```

可選：

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

## `apa` 會做什麼

- 從自然語言的產品想法 bootstrap 新專案
- 先自動推論技術棧，再允許你用 flags 覆蓋
- 產生從 `Phase 0` 開始的分階段文件，包括 PRD、SPEC、ARCHITECTURE、API、DB Schema、實作計畫
- 建立可執行的起始程式碼、測試結構、Makefile 與 agent 設定
- 透過 `apa iterate` 產出結構化提示詞，讓 AI agent 持續在可控流程內工作
- 可自然搭配 `apa-loop`，用分輪交付方式推進：讀取狀態、選 1-3 項任務、驗證、更新狀態、再重複

## 建議工作流程

```bash
# 1. 先把 apa build 起來
go build -o apa .

# 2. 在目標 repo 外建立新專案
./apa init --idea "SaaS 報表平台" --name report-platform --path ~/projects

# 3. 進入產生出的 repo 開始迭代
cd ~/projects/report-platform
./apa list-skills
./apa iterate
make test
```

核心循環：

1. 用 `apa init` 建立第一版專案。
2. 維持從 `Phase 0` 開始的分階段文件，讓 PRD/API/SPEC 的範圍、測試、gate 與報告保持對齊。
3. 預設以 `apa-loop` 搭配 `apa-implement` 推進每一輪交付。
4. 跑 `apa iterate`，讓 agent 實作，再用 `make test` 驗證。
5. 反覆執行直到可交付。

## 交付循環狀態與 `apa-loop` 使用方式

產生出的 repo 應持續更新 `docs/IMPLEMENTATION_STATUS.md` 或 `TASKS.md`。
搭配 `apa-loop` 與 `apa-implement` 使用，讓 agent 持續在實作、測試、修正與文件更新之間循環，直到完成 gate 被滿足。
`apa-loop` 是用來強制執行每輪交付循環的 repo-local skill：先讀狀態檔、選 1-3 個可驗證工作項、跑測試或檢查、更新狀態，再持續重複直到完成 gate 被滿足。
使用方式：
- Codex 專案：先執行 `apa iterate`，再明確要求 agent 使用 `apa-loop` 與 `apa-implement`
- Claude Code 專案：`/apa-loop --max-iterations 30`
- Claude Code 專案：`/cancel-apa-loop`

## 快速範例

```bash
./apa init \
  --idea "內部客服儀表板，支援工單搜尋與 AI 摘要" \
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

## 指令總覽

| 指令 | 用途 |
|---|---|
| `apa init` | 從 idea 建立新專案 |
| `apa iterate` | 產生後續持續交付用的 AI 提示詞 |
| `apa list-skills` | 列出可用的 repo-local skills |
| `apa doctor` | 檢查本機環境與 skills 路徑 |
| `apa version` | 顯示建置版本資訊 |

完整選項請用 `apa <command> --help` 查看。

## Commands & Skills

目前 CLI 指令：

- `init`
- `iterate`
- `list-skills`
- `doctor`
- `version`

目前 repo-local skills：

- `apa-catalog`
- `apa-debug`
- `apa-devops`
- `apa-docs`
- `apa-feature`
- `apa-implement`
- `apa-integration`
- `apa-loop`
- `apa-review`
- `apa-tdd`

## `apa init`

`apa init` 是核心指令，負責一次完成首次 bootstrap。

典型流程：

```text
[1/4] infer tech stack from idea
[2/4] generate code scaffold
[3/4] generate design docs
[4/4] done
```

常見用法：

```bash
# 互動模式
./apa init

# 非互動模式
./apa init --idea "線上訂餐平台" --name food-platform --path ~/projects

# 只預覽，不寫檔
./apa init --dry-run
```

常用 flags：

| Flag | 說明 |
|---|---|
| `--idea` | 用來推論技術棧的產品 idea |
| `--name` | 專案名稱 |
| `--path` | 專案要建立在哪個父目錄 |
| `--type` | `web-app`、`ai-app`、`devops-tool`、`internal-tool`、`platform-service` |
| `--agent` | `codex` 或 `claude-code` |
| `--backend` | `go`、`python`、`node`、`none` |
| `--frontend` | `react`、`next`、`nuxt`、`vue`、`pure-typescript`、`none` |
| `--skills` | 要複製進新專案的 repo-local skills，逗號分隔 |
| `--skills-path` | skills 來源目錄，預設為本 repo 的 `./skills` |
| `--force` | 備份既有目錄後重建 |
| `--dry-run` | 顯示預計產物但不寫入 |

常見產物：

```text
.architect/context.json
docs/
backend/ frontend/ tests/
Makefile
CLAUDE.md or PROMPT.md
agents/ skills/
```

## `apa iterate`

`apa iterate` 會讀取目前 repo 狀態，輸出一段給 AI agent 的結構化提示詞。

它可以在實作前、開發中、或修回歸問題後使用，協助 agent 依照現有文件、任務與限制持續往交付推進。

它也會檢查既有文件是否採用依優先度對齊的 `Phase 0`、`Phase 1`... 段落；如果不是，`apa iterate` 會提醒使用者，並要求 agent 先用 `apa-docs` 重寫文件，再繼續實作。

```bash
./apa iterate
./apa iterate --root ~/projects/report-platform
./apa iterate > prompt.md
```

## Repo-Local Skills

此 repo 的 skill 命名遵循 `apa-*` 系列。

目前範例：

- `apa-catalog`
- `apa-debug`
- `apa-devops`
- `apa-docs`
- `apa-feature`
- `apa-implement`
- `apa-integration`
- `apa-loop`
- `apa-review`
- `apa-tdd`

`apa-docs` 會用依優先度對齊的階段方式寫文件（`Phase 0`、`Phase 1`...），其中 `Phase 0` 永遠是最高優先。每個階段都必須寫出範圍、對應的 PRD/API/SPEC 內容、所需測試、檢測項目、完成標準、明確的下一階段 gate，以及階段完成報告。

可用以下指令查看：

```bash
./apa list-skills
```

## 開發與測試

```bash
go test ./...
make test
make build TARGET=mac
make build TARGET=linux
make build TARGET=windows
make build-all
```

建置 metadata 範例：

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
