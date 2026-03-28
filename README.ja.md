# apa - AI Project Architect

## 言語

| 言語 | README |
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
[![CLI](https://img.shields.io/badge/Type-CLI-111111)](#コマンド一覧)
[![Skills](https://img.shields.io/badge/Repo%20Skills-apa--*-2f855a)](#commands--skills)

## `apa簡要說明`

- 可用觸發詞：`apa簡要說明`、`apa 簡要說明`、`apa說明`、`apa 說明`
- 実装ループを始める: `apa prompt` -> 出力を agent に貼る -> `apa-loop` + `apa-implement` を使うよう指示
- 先にドキュメントを詰める: `apa prompt --docs-only` -> 出力を agent に貼る -> `apa-doc-review` のみ使うよう指示
- Terminal loop wrapper: `bash scripts/apa-loop-setup.sh --max-iterations 30 --reviewer agent-self`
- wrapper を止める: `bash scripts/apa-loop-cancel.sh`

`apa` は、プロダクトの idea をすばやく開発開始可能なプロジェクトの土台に変える Go 製 CLI です。

プロジェクトの context、設計ドキュメント、実行可能な scaffold、そして継続的に回せる AI 反復フローを生成し、アイデアから実装までの距離を短くします。

1 つのツールで、プロジェクト bootstrap、ドキュメント起票、初期コード生成、AI 支援の継続開発ループへの移行までまとめて扱いたい場合に向いています。

## Why `apa`

- 新規プロジェクトを空ディレクトリではなく、ドキュメントとコード付きで始めたい
- AI 作業を、その場しのぎのプロンプトではなく構造化された repo 状態から始めたい
- repo-local の `apa-*` skills と `make test` のような native command で進めたい

## Installation

```bash
go build -o apa .
```

任意:

```bash
cp ./apa /usr/local/bin/apa
```

## Flow

```text
idea
  -> apa init
  -> docs + scaffold + context
  -> apa list-skills
  -> apa prompt
  -> agent implements
  -> make test
  -> repeat
```

## `apa` がやること

- 自然言語のプロダクト idea から新しいプロジェクトを bootstrap
- 実用的な技術スタックを推論し、必要なら flags で上書き可能
- `Phase 0` から始まるフェーズ型ドキュメントとして、PRD、SPEC、ARCHITECTURE、API、DB Schema、実装計画を生成
- 実行可能なスターターコード、テスト構成、Makefile、agent 設定を作成
- `apa prompt` で、AI agent が継続実装しやすい構造化プロンプトを出力
- `apa-loop` と自然に組み合わせられ、状態を読み、1-3 個の作業を選び、検証し、状態更新して繰り返すラウンド型デリバリーを進められる

## 推奨ワークフロー

```bash
# 1. apa をビルド
go build -o apa .

# 2. 対象 repo の外で新規プロジェクトを作成
./apa init --idea "SaaS レポーティングプラットフォーム" --name report-platform --path ~/projects

# 3. 生成された repo に入り反復開始
cd ~/projects/report-platform
./apa list-skills
./apa prompt
make test
```

基本ループ:

1. `apa init` を 1 回実行して初期プロジェクトを作る。
2. `Phase 0` から始まるフェーズ型ドキュメントを維持し、PRD/API/SPEC の範囲、テスト、ゲート、レポートを揃える。
3. 標準のデリバリーループとして `apa-loop` と `apa-implement` を使う。
4. `apa prompt` を実行して agent に作業させ、`make test` で確認する。
5. 出荷可能になるまで繰り返す。

## デリバリーループの状態と `apa-loop` の使い方

生成された repo では `docs/IMPLEMENTATION_STATUS.md` または `TASKS.md` を継続的に更新するべきです。
`apa-loop` と `apa-implement` を組み合わせて、agent が実装、テスト、修正、ドキュメント更新を繰り返し、完了ゲートを満たすまで進み続けるようにします。
`apa-loop` は、状態ファイルを読み、検証可能な 1-3 個の作業項目を選び、テストやチェックを実行し、状態を更新して、完了ゲートを満たすまで繰り返す repo-local skill です。
使い方:
- 実装前に文書を往復で詰めたい場合は、先に `apa prompt --docs-only` を実行し、agent へ `apa-doc-review` のみを使うよう明示する。
- Agent の主運用（Codex と Claude Code 共通）: `apa prompt` を実行したあと、agent に `apa-loop` と `apa-implement` を使うよう明示する
- 任意の terminal ラッパー（生成済み hook や slash command がある環境、たとえば Claude Code 向け）: `bash scripts/apa-loop-setup.sh --max-iterations 30 --reviewer agent-self`
- 任意の slash command: `/apa-loop --max-iterations 30 --reviewer agent-self`
- 任意のキャンセルコマンド: `/cancel-apa-loop`
- Review policy: interactive per round. Ask which reviewer to use (`agent-self`, `apa-codex-review`, or `apa-claude-review`) before review.

## クイック例

```bash
./apa init \
  --idea "チケット検索と AI 要約を備えた社内サポートダッシュボード" \
  --name support-hub \
  --path ~/projects \
  --backend go \
  --frontend next \
  --agent codex

cd ~/projects/support-hub
./apa list-skills
./apa prompt > prompt.md
make test
```

## コマンド一覧

| コマンド | 用途 |
|---|---|
| `apa init` | idea から新しいプロジェクトを作成 |
| `apa prompt` | 継続実装向けの構造化 AI プロンプトを生成 |
| `apa list-skills` | 利用可能な repo-local skills を表示 |
| `apa doctor` | ローカル環境と skills path を確認 |
| `apa version` | ビルド版情報を表示 |

詳細オプションは `apa <command> --help`。

## Commands & Skills

現在の CLI コマンド:

- `init`
- `prompt`
- `list-skills`
- `doctor`
- `version`

現在の repo-local skills:

- `apa-catalog`
- `apa-debug`
- `apa-devops`
- `apa-docs`
- `apa-doc-review`
- `apa-feature`
- `apa-implement`
- `apa-integration`
- `apa-loop`
- `apa-review`
- `apa-codex-review`
- `apa-claude-review`
- `apa-tdd`

## `apa init`

`apa init` は初回 bootstrap の中心コマンドです。

典型フロー:

```text
[1/4] infer tech stack from idea
[2/4] generate code scaffold
[3/4] generate design docs
[4/4] done
```

よくある使い方:

```bash
# 対話モード
./apa init
# Wizard に直接入り、最初の質問は Project idea です。

# 非対話モード
./apa init --idea "オンライン料理注文プラットフォーム" --name food-platform --path ~/projects

# 書き込みなしで確認
./apa init --dry-run
```

よく使う flags:

| Flag | 説明 |
|---|---|
| `--idea` | 技術スタック推論に使うプロダクト idea |
| `--name` | プロジェクト名 |
| `--path` | プロジェクトを作成する親ディレクトリ |
| `--type` | `cli`, `server`, `web-app-server`, `mobile-app-server`, `web-app`, `mobile-app` |
| `--agent` | `codex`、`claude-code`、または `universal` |
| `--backend` | `go`、`python`、`node`、`none` |
| `--frontend` | `react`、`next`、`nuxt`、`vue`、`pure-typescript`、`none` |
| `--unit-test` `--api-test` `--integration-test` `--e2e-test` | `yes` or `no` |
| `--skills` | コピーする repo-local skills。カンマ区切り |
| `--skills-path` | skills のソースディレクトリ。既定は `./skills` |
| `--force` | 既存ディレクトリを退避して再構築 |
| `--dry-run` | 書き込まずに計画だけ表示 |

代表的な生成物:

```text
.architect/context.json
docs/
backend/ frontend/ tests/
Makefile
CLAUDE.md or PROMPT.md
agents/ skills/
```

## `apa prompt`

`apa prompt` は現在の repo 状態を読み取り、AI agent 向けの構造化プロンプトを出力します。

実装前、開発中、回帰修正後のいずれでも使え、既存ドキュメント、タスク、制約に沿って agent が継続作業しやすくなります。

さらに `apa prompt` は既存ドキュメントが優先度順に揃った `Phase 0`, `Phase 1`, ... の段階構成になっているか確認します。そうでない場合は警告を出し、実装前に `apa-docs` で書き直すよう agent に指示します。

```bash
./apa prompt
./apa prompt --docs-only
./apa prompt --root ~/projects/report-platform
./apa prompt > prompt.md
```

## Repo-local Skills

この repository では、repo-local skill 名に `apa-*` 系列を使っています。

現在の例:

- `apa-catalog`
- `apa-debug`
- `apa-devops`
- `apa-docs`
- `apa-doc-review`
- `apa-feature`
- `apa-implement`
- `apa-integration`
- `apa-loop`
- `apa-review`
- `apa-tdd`

`apa-docs` は文書を優先度順の段階（`Phase 0`, `Phase 1`, ...）で作成します。`Phase 0` は常に最優先フェーズです。各フェーズには、範囲、対応する PRD/API/SPEC 内容、必要なテスト、確認項目、完了条件、明確な次フェーズ移行ゲート、フェーズ完了レポートを必ず含めます。
`apa-doc-review` は文書をユーザーとラウンドごとに詰めるための skill です。各修正後に止まってフィードバックを待ち、`docs approved` が出るまでは実装を始めません。

一覧表示:

```bash
./apa list-skills
```

## 開発とテスト

```bash
go test ./...
make test
make build TARGET=mac
make build TARGET=linux
make build TARGET=windows
make build-all
```

ビルド metadata 例:

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
