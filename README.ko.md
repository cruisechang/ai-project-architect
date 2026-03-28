# apa - AI Project Architect

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

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![CLI](https://img.shields.io/badge/Type-CLI-111111)](#명령어)
[![Skills](https://img.shields.io/badge/Repo%20Skills-apa--*-2f855a)](#commands--skills)

## `apa簡要說明`

- 可用觸發詞：`apa簡要說明`、`apa 簡要說明`、`apa說明`、`apa 說明`
- 구현 루프 시작: `apa prompt` -> 출력을 agent에 붙여넣기 -> `apa-loop` + `apa-implement` 사용 지시
- 먼저 문서 반복 수정: `apa prompt --docs-only` -> 출력을 agent에 붙여넣기 -> `apa-doc-review`만 사용 지시
- Terminal loop wrapper: `bash scripts/apa-loop-setup.sh --max-iterations 30 --reviewer agent-self`
- wrapper 중지: `bash scripts/apa-loop-cancel.sh`

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
  -> apa prompt
  -> agent implements
  -> make test
  -> repeat
```

## `apa`가 하는 일

- 자연어 제품 idea로 새 프로젝트를 bootstrap
- 현실적인 기술 스택을 추론하고 flags로 덮어쓸 수 있음
- `Phase 0`부터 시작하는 단계형 문서로 PRD, SPEC, ARCHITECTURE, API, DB Schema, 구현 계획 생성
- 실행 가능한 시작 코드, 테스트 구조, Makefile, agent 설정 생성
- `apa prompt`로 AI agent가 계속 작업할 수 있는 구조화된 프롬프트 출력
- `apa-loop`와 자연스럽게 결합되어 상태 읽기, 1-3개 작업 선택, 검증, 상태 업데이트, 반복의 라운드 기반 전달을 진행할 수 있음

## 권장 워크플로

```bash
# 1. apa 빌드
go build -o apa .

# 2. 대상 repo 밖에서 새 프로젝트 생성
./apa init --idea "SaaS 리포팅 플랫폼" --name report-platform --path ~/projects

# 3. 생성된 repo로 이동해서 반복 작업 시작
cd ~/projects/report-platform
./apa list-skills
./apa prompt
make test
```

기본 반복 루프:

1. `apa init`으로 첫 프로젝트를 만든다.
2. `Phase 0`부터 시작하는 단계형 문서를 유지해 PRD/API/SPEC의 범위, 테스트, gate, 보고서를 정렬한다.
3. 기본 전달 루프로 `apa-loop`와 `apa-implement`를 사용한다.
4. `apa prompt`를 실행하고 agent가 작업한 뒤 `make test`로 검증한다.
5. 배포 가능한 상태가 될 때까지 반복한다.

## 전달 루프 상태와 `apa-loop` 사용법

생성된 repo는 `docs/IMPLEMENTATION_STATUS.md` 또는 `TASKS.md`를 계속 업데이트해야 합니다.
`apa-loop`와 `apa-implement`를 함께 사용해서 agent가 구현, 테스트, 수정, 문서 업데이트를 계속 순환하고 완료 gate가 충족될 때까지 진행하도록 합니다.
`apa-loop`는 상태 파일을 읽고, 검증 가능한 작업 1-3개를 고르고, 테스트/체크를 실행하고, 상태를 갱신한 뒤 완료 gate가 충족될 때까지 반복하도록 강제하는 repo-local skill입니다.
사용 방법:
- 구현 전에 문서를 먼저 반복해서 다듬고 싶다면 `apa prompt --docs-only`를 실행하고 agent에게 `apa-doc-review`만 사용하라고 명시한다.
- 기본 agent 사용 흐름(Codex 와 Claude Code 공통): `apa prompt`를 실행한 뒤 agent에게 `apa-loop`와 `apa-implement`를 명시적으로 사용하라고 지시
- 선택 가능한 terminal 래퍼(생성된 hook 또는 slash command 가 있는 환경, 예: Claude Code): `bash scripts/apa-loop-setup.sh --max-iterations 30 --reviewer agent-self`
- 선택 가능한 slash command: `/apa-loop --max-iterations 30 --reviewer agent-self`
- 선택 가능한 취소 명령: `/cancel-apa-loop`
- Review policy: interactive per round. Ask which reviewer to use (`agent-self`, `apa-codex-review`, or `apa-claude-review`) before review.

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
./apa prompt > prompt.md
make test
```

## 명령어

| 명령어 | 용도 |
|---|---|
| `apa init` | idea로 새 프로젝트 생성 |
| `apa prompt` | 지속 구현용 구조화 AI 프롬프트 생성 |
| `apa list-skills` | 사용 가능한 repo-local skills 표시 |
| `apa doctor` | 로컬 환경과 skills 경로 점검 |
| `apa version` | 빌드 버전 정보 출력 |

전체 옵션은 `apa <command> --help`에서 확인할 수 있습니다.

## Commands & Skills

현재 CLI 명령어:

- `init`
- `prompt`
- `list-skills`
- `doctor`
- `version`

현재 repo-local skills:

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
# Wizard로 바로 들어가며 첫 질문은 Project idea입니다.

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
| `--type` | `cli`, `server`, `web-app-server`, `mobile-app-server`, `web-app`, `mobile-app` |
| `--agent` | `codex`, `claude-code`, 또는 `universal` |
| `--backend` | `go`, `python`, `node`, `none` |
| `--frontend` | `react`, `next`, `nuxt`, `vue`, `pure-typescript`, `none` |
| `--unit-test` `--api-test` `--integration-test` `--e2e-test` | `yes` or `no` |
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

## `apa prompt`

`apa prompt`는 현재 repo 상태를 읽고 AI agent용 구조화 프롬프트를 출력합니다.

구현 전, 개발 중, 회귀 수정 이후에도 사용할 수 있으며, 기존 문서, 작업, 제약사항에 맞춰 agent가 계속 일하도록 돕습니다.

또한 기존 문서가 우선순위 기준으로 정렬된 `Phase 0`, `Phase 1`, ... 구성을 따르는지 검사합니다. 그렇지 않으면 `apa prompt`가 경고를 출력하고, 구현 전에 `apa-docs`로 다시 쓰라고 agent에 지시합니다.

```bash
./apa prompt
./apa prompt --docs-only
./apa prompt --root ~/projects/report-platform
./apa prompt > prompt.md
```

## Repo-local Skills

이 repository의 repo-local skill 이름은 `apa-*` 시리즈를 따릅니다.

현재 예시:

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

`apa-docs`는 문서를 우선순위 기반 단계(`Phase 0`, `Phase 1`, ...)로 작성합니다. `Phase 0`는 항상 가장 높은 우선순위 단계입니다. 각 단계에는 범위, 정렬된 PRD/API/SPEC 내용, 필요한 테스트, 점검 항목, 완료 기준, 명확한 다음 단계 gate, 단계 완료 보고서가 반드시 포함되어야 합니다.
`apa-doc-review`는 사용자와 문서를 라운드별로 다듬는 skill입니다. 매 수정 후 피드백을 기다리며, `docs approved`가 나오기 전에는 구현을 시작하지 않습니다.

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
