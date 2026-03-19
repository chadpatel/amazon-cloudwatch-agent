#!/bin/bash
# ralph.sh - Autonomous task loop for feature implementation
# Usage: ./ralph.sh <feature-path> [--max-iterations N] [--agent AGENT]
set -e

if [[ -z "$1" ]] || [[ "$1" == --* ]]; then
  echo "Usage: ./ralph.sh <feature-path> [--max-iterations N] [--agent AGENT]"
  echo "Example: ./ralph.sh my-feature"
  echo "         ./ralph.sh my-feature --max-iterations 20"
  exit 1
fi

FEATURE="$1"
shift

if [[ -d "docs/$FEATURE" ]]; then
  DOCS_DIR="docs/$FEATURE"
elif [[ -d "docs/design/$FEATURE" ]]; then
  DOCS_DIR="docs/design/$FEATURE"
else
  echo "Error: Could not find docs directory for feature '$FEATURE'"
  echo "Searched: docs/$FEATURE, docs/design/$FEATURE"
  exit 1
fi

MAX_ITERATIONS=50
STATE_DIR="$DOCS_DIR/.ralph"
AGENT="ralph-execution"

while [[ $# -gt 0 ]]; do
  case $1 in
    --max-iterations) MAX_ITERATIONS="$2"; shift 2 ;;
    --agent) AGENT="$2"; shift 2 ;;
    *) echo "Unknown option: $1"; exit 1 ;;
  esac
done

mkdir -p "$STATE_DIR"

FEATURE_NAME=$(basename "$DOCS_DIR")
FEATURE_UPPER=$(echo "$FEATURE_NAME" | tr '[:lower:]' '[:upper:]' | tr '-' '_')

log() { echo "[$(date '+%H:%M:%S')] $1"; }

notify() {
  command -v osascript &>/dev/null && osascript -e "display notification \"$1\" with title \"ralph: $FEATURE_NAME\"" || true
  command -v notify-send &>/dev/null && notify-send "ralph: $FEATURE_NAME" "$1" || true
}

cleanup() {
  log "Stopped at iteration $(cat "$STATE_DIR/.iteration" 2>/dev/null || echo 0)"
  notify "stopped"
  exit 130
}

trap cleanup INT TERM

log "Starting ralph for '$FEATURE_NAME' ($MAX_ITERATIONS iterations max)"
log "Docs: $DOCS_DIR"
notify "started"

PROMPT="You are implementing the $FEATURE_NAME feature. Work through the tasks in $DOCS_DIR/plan.md.

Instructions:
1. Read $DOCS_DIR/plan.md to find the next unchecked task (- [ ])
2. Read $DOCS_DIR/design.md for technical specifications (if exists)
3. Read $DOCS_DIR/worklog.md for context from previous iterations (if exists)
4. Read AGENTS.md for project conventions
5. Implement that ONE task only
6. Mark the task as done in $DOCS_DIR/plan.md (- [ ] → - [x])
7. Append an entry to $DOCS_DIR/worklog.md with:
   - Timestamp, task name, status
   - What you completed
   - Blockers (if any)
   - Next steps

CRITICAL RULES:
- Do NOT commit unless explicitly instructed in the plan
- NEVER use --no-verify when committing
- Only work on ONE task per iteration

FINAL ITERATION (when all tasks are done):
- If all tasks in plan.md are checked [x], output <promise>${FEATURE_UPPER}_COMPLETE</promise>

Begin by reading the plan and worklog to determine current state."

check_catalog_entry() {
  local feature_name="$1"
  local catalog_dir="docs/design/_catalog"
  [[ ! -d "$catalog_dir" ]] && return 0
  grep -rq "$feature_name" "$catalog_dir"/*.md 2>/dev/null && return 0
  # Heading format uses spaces around first dash: "260224 - feature-name"
  grep -rq "${feature_name/-/ - }" "$catalog_dir"/*.md 2>/dev/null && return 0
  return 1
}

for i in $(seq 1 "$MAX_ITERATIONS"); do
  echo "$i" > "$STATE_DIR/.iteration"
  log "=== Iteration $i/$MAX_ITERATIONS ==="

  output=$(kiro-cli chat --agent "$AGENT" --trust-all-tools --no-interactive -- "$PROMPT" 2>&1 | tee /dev/tty) || true

  if [[ "$output" == *"<promise>${FEATURE_UPPER}_COMPLETE</promise>"* ]]; then
    if [[ "$DOCS_DIR" == docs/design/* ]] && ! check_catalog_entry "$FEATURE_NAME"; then
      log "⚠ Warning: No catalog entry found for '$FEATURE_NAME'"
      log "  Add an entry to docs/design/_catalog/<category>.md"
      notify "missing catalog entry"
    fi
    log "✓ $FEATURE_NAME complete at iteration $i"
    if [[ "$FEATURE_NAME" == *docs-onboarding* ]]; then
      log ""
      log "📋 Run @chief-review to evaluate the onboarding quality:"
      log "   kiro-cli chat --trust-all-tools --no-interactive -- '@chief-review'"
    fi
    notify "complete!"
    exit 0
  fi

  sleep 2
done

log "Reached max iterations ($MAX_ITERATIONS)"
notify "max iterations reached"
