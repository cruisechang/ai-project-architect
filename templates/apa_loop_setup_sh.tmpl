#!/bin/bash

# apa-loop Setup Script
# Activates the apa-loop Stop hook for the current session.

set -euo pipefail

MAX_ITERATIONS=0
COMPLETION_PROMISE="COMPLETE"

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
    -h|--help)
      cat <<'HELP'
apa-loop-setup — activate the apa delivery loop for this session

USAGE:
  scripts/apa-loop-setup.sh [OPTIONS]

OPTIONS:
  --max-iterations <n>        Stop after N iterations (default: unlimited)
  --completion-promise <text> Exact phrase Claude must output to exit (default: COMPLETE)
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

mkdir -p .claude

cat > .claude/apa-loop.local.md <<EOF
---
active: true
iteration: 1
max_iterations: ${MAX_ITERATIONS}
completion_promise: "${COMPLETION_PROMISE}"
started_at: "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
---
EOF

echo "🔄 apa-loop activated"
echo ""
echo "   Max iterations : $(if [[ $MAX_ITERATIONS -gt 0 ]]; then echo "$MAX_ITERATIONS"; else echo "unlimited"; fi)"
echo "   Completion      : <promise>${COMPLETION_PROMISE}</promise>"
echo ""
echo "The stop hook is now active. Each time you exit, apa-loop feeds back a"
echo "fresh prompt built from skills/apa-loop/SKILL.md and docs/IMPLEMENTATION_STATUS.md."
echo ""
echo "Begin: read skills/apa-loop/SKILL.md and docs/IMPLEMENTATION_STATUS.md,"
echo "then start the first implementation round."
echo ""
echo "To stop early: delete .claude/apa-loop.local.md"
