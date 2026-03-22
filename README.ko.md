# apa - AI Project Architect

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![CLI](https://img.shields.io/badge/Type-CLI-111111)](#명령어)
[![Skills](https://img.shields.io/badge/Repo%20Skills-apa--*-2f855a)](#commands--skills)

`apa`는 제품 idea를 빠르게 개발 가능한 프로젝트 시작점으로 바꿔주는 Go CLI입니다.

프로젝트 context, 설계 문서, 실행 가능한 scaffold, 그리고 반복 가능한 AI 작업 흐름을 생성해서 아이디어에서 구현까지 가는 시간을 줄여줍니다.

프로젝트 bootstrap, 문서 초안, 초기 코드 생성, 그리고 AI 보조 전달 루프까지 하나의 도구로 묶고 싶을 때 적합합니다.

## Why `apa`

- 새 프로젝트를 빈 폴더가 아니라 문서와 코드가 있는 상태로 시작하고 싶을 때
- AI 작업을 임시 프롬프트가 아니라 구조화된 repo 상태에서 시작하고 싶을 때
- repo-local `apa-*` skills와 `make test` 같은 native command로 전달 루프를 돌리고 싶을 때

## Installation

```bash
go build -o apa .
```

선택 사항:

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

## 언어

| 언어 | README |
|---|---|
| English | [README.md](README.md) |
| 简体中文 | [README.zh-CN.md](README.zh-CN.md) |
| 繁體中文 | [README.zh-TW.md](README.zh-TW.md) |
| Deutsch | [README.de.md](README.de.md) |
| Français | [README.fr.md](README.fr.md) |
| Español | [README.es.md](README.es.md) |
| 日本語 | [README.ja.md](README.ja.md) |
| 한국어 | [README.ko.md](README.ko.md) |

## `apa`가 하는 일

- 자연어 제품 idea로 새 프로젝트를 bootstrap
- 현실적인 기술 스택을 추론하고 flags로 덮어쓸 수 있음
- PRD, SPEC, ARCHITECTURE, API, DB Schema, 구현 계획 생성
- 실행 가능한 시작 코드, 테스트 구조, Makefile, agent 설정 생성
- `apa iterate`로 AI agent가 계속 작업할 수 있는 구조화된 프롬프트 출력

## 권장 워크플로

```bash
# 1. apa 빌드
go build -o apa .

# 2. 대상 repo 밖에서 새 프로젝트 생성
./apa init --idea "SaaS 리포팅 플랫폼" --name report-platform --path ~/projects

# 3. 생성된 repo로 이동해서 반복 작업 시작
cd ~/projects/report-platform
./apa list-skills
./apa iterate
make test
```

기본 반복 루프:

1. `apa init`으로 첫 프로젝트를 만든다.
2. repo-local `apa-*` skills로 구현 작업을 진행한다.
3. `apa iterate`를 실행하고 agent가 작업한 뒤 `make test`로 검증한다.
4. 배포 가능한 상태가 될 때까지 반복한다.

## 빠른 예시

```bash
./apa init \
  --idea "티켓 검색과 AI 요약 기능이 있는 내부 지원 대시보드" \
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

## 명령어

| 명령어 | 용도 |
|---|---|
| `apa init` | idea로 새 프로젝트 생성 |
| `apa iterate` | 지속 구현용 구조화 AI 프롬프트 생성 |
| `apa list-skills` | 사용 가능한 repo-local skills 표시 |
| `apa doctor` | 로컬 환경과 skills 경로 점검 |
| `apa version` | 빌드 버전 정보 출력 |

전체 옵션은 `apa <command> --help`에서 확인할 수 있습니다.

## Commands & Skills

현재 CLI 명령어:

- `init`
- `iterate`
- `list-skills`
- `doctor`
- `version`

현재 repo-local skills:

- `apa-catalog`
- `apa-debug`
- `apa-devops`
- `apa-docs`
- `apa-feature`
- `apa-integration`
- `apa-review`
- `apa-tdd`

## `apa init`

`apa init`은 최초 bootstrap을 담당하는 핵심 명령입니다.

일반적인 흐름:

```text
[1/4] infer tech stack from idea
[2/4] generate code scaffold
[3/4] generate design docs
[4/4] done
```

자주 쓰는 예시:

```bash
# 인터랙티브 모드
./apa init

# 비인터랙티브 모드
./apa init --idea "온라인 음식 주문 플랫폼" --name food-platform --path ~/projects

# 미리보기만
./apa init --dry-run
```

주요 flags:

| Flag | 설명 |
|---|---|
| `--idea` | 기술 스택 추론에 사용할 제품 idea |
| `--name` | 프로젝트 이름 |
| `--path` | 프로젝트를 생성할 부모 디렉터리 |
| `--agent` | `codex` 또는 `claude-code` |
| `--backend` | `go`, `python`, `node`, `none` |
| `--frontend` | `react`, `next`, `nuxt`, `vue`, `pure-typescript`, `none` |
| `--skills` | 복사할 repo-local skills, 쉼표 구분 |
| `--skills-path` | skills 소스 디렉터리, 기본값은 `./skills` |
| `--force` | 기존 디렉터리를 백업한 뒤 재생성 |
| `--dry-run` | 파일을 쓰지 않고 계획만 표시 |

대표 산출물:

```text
.architect/context.json
docs/
backend/ frontend/ tests/
Makefile
CLAUDE.md or PROMPT.md
agents/ skills/
```

## `apa iterate`

`apa iterate`는 현재 repo 상태를 읽고 AI agent용 구조화 프롬프트를 출력합니다.

구현 전, 개발 중, 회귀 수정 이후에도 사용할 수 있으며, 기존 문서, 작업, 제약사항에 맞춰 agent가 계속 일하도록 돕습니다.

```bash
./apa iterate
./apa iterate --root ~/projects/report-platform
./apa iterate > prompt.md
```

## Repo-local Skills

이 repository의 repo-local skill 이름은 `apa-*` 시리즈를 따릅니다.

현재 예시:

- `apa-catalog`
- `apa-debug`
- `apa-devops`
- `apa-docs`
- `apa-feature`
- `apa-integration`
- `apa-review`
- `apa-tdd`

목록 확인:

```bash
./apa list-skills
```

## 개발 및 테스트

```bash
go test ./...
make test
make build TARGET=mac
make build TARGET=linux
make build TARGET=windows
make build-all
```

빌드 metadata 예시:

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
