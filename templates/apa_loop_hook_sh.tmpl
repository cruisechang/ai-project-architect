#!/bin/bash

# apa-loop Stop Hook
# Intercepts Claude's exit attempts when an apa-loop is active.
# Feeds a fresh prompt (built from IMPLEMENTATION_STATUS.md) back each iteration.

set -euo pipefail

HOOK_INPUT=$(cat)

STATE_FILE=".claude/apa-loop.local.md"

if [[ ! -f "$STATE_FILE" ]]; then
  exit 0
fi

FRONTMATTER=$(sed -n '/^---$/,/^---$/{ /^---$/d; p; }' "$STATE_FILE")
ITERATION=$(echo "$FRONTMATTER" | grep '^iteration:' | sed 's/iteration: *//' || true)
MAX_ITERATIONS=$(echo "$FRONTMATTER" | grep '^max_iterations:' | sed 's/max_iterations: *//' || true)
COMPLETION_PROMISE=$(echo "$FRONTMATTER" | grep '^completion_promise:' | sed 's/completion_promise: *//' | sed 's/^"\(.*\)"$/\1/' || true)

if [[ ! "$ITERATION" =~ ^[0-9]+$ ]]; then
  echo "⚠️  apa-loop: state file corrupted (iteration invalid)" >&2
  rm "$STATE_FILE"
  exit 0
fi

if [[ ! "$MAX_ITERATIONS" =~ ^[0-9]+$ ]]; then
  echo "⚠️  apa-loop: state file corrupted (max_iterations invalid)" >&2
  rm "$STATE_FILE"
  exit 0
fi

if [[ $MAX_ITERATIONS -gt 0 ]] && [[ $ITERATION -ge $MAX_ITERATIONS ]]; then
  echo "🛑 apa-loop: max iterations ($MAX_ITERATIONS) reached."
  rm "$STATE_FILE"
  exit 0
fi

TRANSCRIPT_PATH=$(echo "$HOOK_INPUT" | jq -r '.transcript_path')

if [[ -z "$TRANSCRIPT_PATH" ]] || [[ "$TRANSCRIPT_PATH" == "null" ]] || [[ ! -f "$TRANSCRIPT_PATH" ]]; then
  echo "⚠️  apa-loop: transcript not found" >&2
  rm "$STATE_FILE"
  exit 0
fi

if ! grep -q '"role":"assistant"' "$TRANSCRIPT_PATH"; then
  echo "⚠️  apa-loop: no assistant messages in transcript" >&2
  rm "$STATE_FILE"
  exit 0
fi

LAST_LINE=$(grep '"role":"assistant"' "$TRANSCRIPT_PATH" | tail -1)
LAST_OUTPUT=$(echo "$LAST_LINE" | jq -r '
  .message.content |
  map(select(.type == "text")) |
  map(.text) |
  join("\n")
' 2>&1) || true

if [[ "$COMPLETION_PROMISE" != "null" ]] && [[ -n "$COMPLETION_PROMISE" ]]; then
  PROMISE_TEXT=$(echo "$LAST_OUTPUT" | perl -0777 -pe 's/.*?<promise>(.*?)<\/promise>.*/$1/s; s/^\s+|\s+$//g; s/\s+/ /g' 2>/dev/null || echo "")
  if [[ -n "$PROMISE_TEXT" ]] && [[ "$PROMISE_TEXT" = "$COMPLETION_PROMISE" ]]; then
    echo "✅ apa-loop: completed — <promise>$COMPLETION_PROMISE</promise> detected."
    rm "$STATE_FILE"
    exit 0
  fi
fi

NEXT_ITERATION=$((ITERATION + 1))

TEMP_FILE="${STATE_FILE}.tmp.$$"
sed "s/^iteration: .*/iteration: $NEXT_ITERATION/" "$STATE_FILE" > "$TEMP_FILE"
mv "$TEMP_FILE" "$STATE_FILE"

SKILL_CONTENT=""
if [[ -f "skills/apa-loop/SKILL.md" ]]; then
  SKILL_CONTENT=$(cat "skills/apa-loop/SKILL.md")
fi

STATUS_CONTENT="(no status file yet — create docs/IMPLEMENTATION_STATUS.md in this round)"
if [[ -f "docs/IMPLEMENTATION_STATUS.md" ]]; then
  STATUS_CONTENT=$(cat "docs/IMPLEMENTATION_STATUS.md")
fi

NEXT_PROMPT="Continue working according to the apa-loop methodology below.

${SKILL_CONTENT}

---

Current implementation status (docs/IMPLEMENTATION_STATUS.md):

${STATUS_CONTENT}

---

Keep iterating until the Completion Gate conditions are all met.
When genuinely complete, output: <promise>${COMPLETION_PROMISE}</promise>
Do NOT output the promise unless it is completely and unequivocally true."

ITER_DISPLAY="$NEXT_ITERATION"
if [[ $MAX_ITERATIONS -gt 0 ]]; then
  ITER_DISPLAY="$NEXT_ITERATION / $MAX_ITERATIONS"
fi

jq -n \
  --arg prompt "$NEXT_PROMPT" \
  --arg msg "🔄 apa-loop iteration ${ITER_DISPLAY} | complete with <promise>${COMPLETION_PROMISE}</promise>" \
  '{
    "decision": "block",
    "reason": $prompt,
    "systemMessage": $msg
  }'

exit 0
