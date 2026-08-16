#!/bin/bash
# usage: run_cell.sh <arm> <scenario>
set -uo pipefail
ARM="$1"; SCEN="$2"
BASE=/Users/clkao/.claude/jobs/4e49247e/tmp/fitgate
REPO=/Users/clkao/git/spacedock-research/spacedock-v1
OUT="$BASE/transcripts/${SCEN}-${ARM}.txt"
P="$BASE/prompts/${SCEN}-${ARM}.txt"
mkdir -p "$BASE/prompts" "$BASE/transcripts"
{
  echo "You are a Spacedock first officer. The following is your operating contract — the instruction files you have loaded for this session. Follow it."
  echo
  for f in "$REPO/skills/first-officer/references/first-officer-shared-core.md" "$BASE/arms/${ARM}.md" "$REPO/docs/dev/README.md"; do
    case "$f" in
      *first-officer-shared-core.md) n=first-officer-shared-core.md ;;
      *arms/*) n=fo-write-core.md ;;
      *) n=docs/dev/README.md ;;
    esac
    echo "===== BEGIN $n ====="
    cat "$f"
    echo "===== END $n ====="
    echo
  done
  echo "===== SITUATION ====="
  cat "$BASE/scenarios/${SCEN}.md"
} > "$P"
claude -p --model opus --disallowed-tools "Bash,Edit,Write,Agent,Task,NotebookEdit,WebFetch,WebSearch,Read,Grep,Glob" < "$P" > "$OUT" 2>&1
echo "CELL_DONE ${SCEN}-${ARM} exit=$? bytes=$(wc -c < "$OUT")"
