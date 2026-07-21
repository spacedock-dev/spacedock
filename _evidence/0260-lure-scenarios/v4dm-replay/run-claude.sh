#!/bin/bash
SCRATCH=/private/tmp/claude-501/-Users-clkao-git-spacedock-research-spacedock-v1/4fc7973c-f2ca-4d34-a2b1-cf24471173af/scratchpad
CELL="$1"   # s6 or s6c
P="$SCRATCH/${CELL}-patched-prompt.txt"
TOOLS="Bash,Edit,Write,Agent,Task,NotebookEdit,WebFetch,WebSearch,Read,Grep,Glob"
LOG="$SCRATCH/${CELL}-patched-run.log"
: > "$LOG"
for n in 1 2 3 4; do
  timeout 300 claude -p --model opus --disallowed-tools "$TOOLS" < "$P" > "$SCRATCH/${CELL}-patched-$n.out" 2> "$SCRATCH/${CELL}-patched-$n.err"
  echo "[$(date +%H:%M:%S)] ${CELL} run $n done exit=$? bytes=$(wc -c < "$SCRATCH/${CELL}-patched-$n.out")" >> "$LOG"
done
echo "${CELL} PATCHED COMPLETE" >> "$LOG"
