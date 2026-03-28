#!/bin/bash

# apa-loop Setup Script
# Activates the apa-loop Stop hook for the current session.

set -euo pipefail

MAX_ITERATIONS=0
COMPLETION_PROMISE="COMPLETE"
REVIEWER=""

is_valid_reviewer() {
  case "$1" in
    agent-self|apa-codex-review|apa-claude-review) return 0 ;;
    *) return 1 ;;
  esac
}

while [[ $# -gt 0 ]]; do
  case $1 in
    --max-iterations)
      if [[ -z "${2:-}" ]] || ! [[ "$2" =~ ^[0-9]+$ ]]; then
        echo "❌ --max-iterations requires a non-negative integer" >&2
        exit 1
      fi
      MAX_ITERATIONS="$2"
      shift 2
      ;;
    --completion-promise)
      if [[ -z "${2:-}" ]]; then
        echo "❌ --completion-promise requires a text argument" >&2
        exit 1
      fi
      COMPLETION_PROMISE="$2"
      shift 2
      ;;
    --reviewer)
      if [[ -z "${2:-}" ]]; then
        echo "❌ --reviewer requires one of: agent-self, apa-codex-review, apa-claude-review" >&2
        exit 1
      fi
      if ! is_valid_reviewer "$2"; then
        echo "❌ invalid reviewer: $2 (allowed: agent-self, apa-codex-review, apa-claude-review)" >&2
        exit 1
      fi
      REVIEWER="$2"
      shift 2
      ;;
    -h|--help)
      cat <<'HELP'
apa-loop-setup — activate the apa delivery loop for this session

USAGE:
  scripts/apa-loop-setup.sh [OPTIONS]

OPTIONS:
  --max-iterations <n>        Stop after N iterations (default: unlimited)
  --completion-promise <text> Exact phrase Claude must output to exit (default: COMPLETE)
  --reviewer <name>           Default reviewer (agent-self | apa-codex-review | apa-claude-review)
  -h, --help                  Show this help

STOPPING:
  Claude must output <promise>COMPLETE</promise> (or your custom phrase).
  Delete .claude/apa-loop.local.md to stop immediately.
HELP
      exit 0
      ;;
    *)
      echo "❌ Unknown argument: $1" >&2
      exit 1
      ;;
  esac
done

if [[ -z "$REVIEWER" ]]; then
  if [[ -t 0 && -t 1 ]]; then
    echo "Select default reviewer for apa-loop:"
    echo "  1) agent-self"
    echo "  2) apa-codex-review"
    echo "  3) apa-claude-review"
    read -r -p "Reviewer [agent-self]: " REVIEWER_INPUT
    case "${REVIEWER_INPUT:-agent-self}" in
      1) REVIEWER="agent-self" ;;
      2) REVIEWER="apa-codex-review" ;;
      3) REVIEWER="apa-claude-review" ;;
      agent-self|apa-codex-review|apa-claude-review) REVIEWER="${REVIEWER_INPUT}" ;;
      *)
        echo "❌ invalid reviewer: ${REVIEWER_INPUT}" >&2
        exit 1
        ;;
    esac
  else
    REVIEWER="agent-self"
    echo "ℹ️  non-interactive shell detected, default reviewer: ${REVIEWER}"
  fi
fi

mkdir -p .claude

cat > .claude/apa-loop.local.md <<EOF
---
active: true
iteration: 1
max_iterations: ${MAX_ITERATIONS}
completion_promise: "${COMPLETION_PROMISE}"
reviewer: "${REVIEWER}"
started_at: "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
---
EOF

echo "🔄 apa-loop activated"
echo ""
echo "   Max iterations : $(if [[ $MAX_ITERATIONS -gt 0 ]]; then echo "$MAX_ITERATIONS"; else echo "unlimited"; fi)"
echo "   Completion      : <promise>${COMPLETION_PROMISE}</promise>"
echo "   Reviewer        : ${REVIEWER}"
echo ""
echo "The stop hook is now active. Each time you exit, apa-loop feeds back a"
echo "fresh prompt built from skills/apa-loop/SKILL.md and docs/IMPLEMENTATION_STATUS.md."
echo "It also requires interactive reviewer selection each round."
echo ""
echo "Begin: read skills/apa-loop/SKILL.md and docs/IMPLEMENTATION_STATUS.md,"
echo "then start the first implementation round."
echo ""
echo "To stop early: delete .claude/apa-loop.local.md"
