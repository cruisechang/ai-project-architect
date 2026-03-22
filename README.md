# apa - AI Project Architect

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![CLI](https://img.shields.io/badge/Type-CLI-111111)](#commands)
[![Skills](https://img.shields.io/badge/Repo%20Skills-apa--*-2f855a)](#commands--skills)

`apa` is a Go CLI that turns a product idea into a working project starter.

It generates project context, design docs, runnable scaffolding, and a repeatable AI iteration workflow so you can go from idea to implementation faster.

Use it when you want one tool to handle project bootstrap, documentation setup, starter code generation, and the handoff into an AI-assisted delivery loop.

## Why `apa`

- You want a new project to start with docs and code, not just empty folders
- You want AI work to begin from a structured repo state instead of ad hoc prompts
- You want repo-local `apa-*` skills and native commands like `make test` to drive delivery

## Installation

```bash
go build -o apa .
```

Optional:

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

## Languages

| Language | README |
|---|---|
| English | [README.md](README.md) |
| 简体中文 | [README.zh-CN.md](README.zh-CN.md) |
| 繁體中文 | [README.zh-TW.md](README.zh-TW.md) |
| Deutsch | [README.de.md](README.de.md) |
| Français | [README.fr.md](README.fr.md) |
| Español | [README.es.md](README.es.md) |
| 日本語 | [README.ja.md](README.ja.md) |
| 한국어 | [README.ko.md](README.ko.md) |

## What `apa` does

- Bootstraps a new project from a plain-language product idea
- Infers a practical stack, then lets you override it with flags
- Generates docs such as PRD, SPEC, ARCHITECTURE, API, DB schema, and implementation plan
- Creates runnable starter code, test setup, Makefile targets, and agent config
- Outputs an `apa iterate` prompt so an AI agent can keep working in a controlled loop

## Recommended Workflow

```bash
# 1. Build apa
go build -o apa .

# 2. Create a new project outside the target repo
./apa init --idea "SaaS reporting platform" --name report-platform --path ~/projects

# 3. Enter the generated repo and iterate
cd ~/projects/report-platform
./apa list-skills
./apa iterate
make test
```

Core loop:

1. Run `apa init` once to create the initial project.
2. Use repo-local `apa-*` skills to guide implementation work.
3. Run `apa iterate`, let the agent implement, then validate with `make test`.
4. Repeat until the repo is in a shippable state.

## Quick Example

```bash
./apa init \
  --idea "Internal support dashboard with ticket search and AI summaries" \
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

## Commands

| Command | Purpose |
|---|---|
| `apa init` | Create a new project from an idea |
| `apa iterate` | Generate the structured AI prompt for continued delivery |
| `apa list-skills` | Show available repo-local skills |
| `apa doctor` | Check local environment and skills path |
| `apa version` | Print build version info |

Run `apa <command> --help` for full options.

## Commands & Skills

Current CLI commands:

- `init`
- `iterate`
- `list-skills`
- `doctor`
- `version`

Current repo-local skills:

- `apa-catalog`
- `apa-debug`
- `apa-devops`
- `apa-docs`
- `apa-feature`
- `apa-integration`
- `apa-review`
- `apa-tdd`

## `apa init`

`apa init` is the main command. It handles first-time project bootstrap in one run.

Typical flow:

```text
[1/4] infer tech stack from idea
[2/4] generate code scaffold
[3/4] generate design docs
[4/4] done
```

Example usages:

```bash
# Interactive mode
./apa init

# Non-interactive mode
./apa init --idea "Online food ordering platform" --name food-platform --path ~/projects

# Preview only
./apa init --dry-run
```

Common flags:

| Flag | Description |
|---|---|
| `--idea` | Product idea used for stack inference |
| `--name` | Project name |
| `--path` | Parent directory where the project will be created |
| `--agent` | `codex` or `claude-code` |
| `--backend` | `go`, `python`, `node`, or `none` |
| `--frontend` | `react`, `next`, `nuxt`, `vue`, `pure-typescript`, or `none` |
| `--skills` | Comma-separated repo-local skills to copy into the generated project |
| `--skills-path` | Source directory for skills, defaults to this repo's `./skills` |
| `--force` | Back up an existing directory and rebuild |
| `--dry-run` | Show planned output without writing files |

Generated output usually includes:

```text
.architect/context.json
docs/
backend/ frontend/ tests/
Makefile
CLAUDE.md or PROMPT.md
agents/ skills/
```

## `apa iterate`

`apa iterate` reads the current repository and prints a structured prompt for the AI agent.

Use it before implementation, during development, or after regressions. It helps the agent stay aligned with the repo state, existing docs, queued tasks, and delivery constraints.

```bash
./apa iterate
./apa iterate --root ~/projects/report-platform
./apa iterate > prompt.md
```

## Repo-Local Skills

This repository uses the `apa-*` naming series for repo-local skills.

Current examples:

- `apa-catalog`
- `apa-debug`
- `apa-devops`
- `apa-docs`
- `apa-feature`
- `apa-integration`
- `apa-review`
- `apa-tdd`

List them with:

```bash
./apa list-skills
```

## Development

```bash
go test ./...
make test
make build TARGET=mac
make build TARGET=linux
make build TARGET=windows
make build-all
```

Build metadata example:

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
