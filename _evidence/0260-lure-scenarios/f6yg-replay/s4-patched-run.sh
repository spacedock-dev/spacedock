#!/bin/bash
SCRATCH=/private/tmp/claude-501/-Users-clkao-git-spacedock-research-spacedock-v1/4fc7973c-f2ca-4d34-a2b1-cf24471173af/scratchpad
P="$SCRATCH/s4-patched-prompt.txt"
TOOLS="Bash,Edit,Write,Agent,Task,NotebookEdit,WebFetch,WebSearch,Read,Grep,Glob"
: > "$SCRATCH/s4-patched-run.log"
for n in 1 2 3 4 5 6 7 8; do
  timeout 300 claude -p --model opus --disallowed-tools "$TOOLS" < "$P" > "$SCRATCH/s4-patched-$n.out" 2> "$SCRATCH/s4-patched-$n.err"
  echo "[$(date +%H:%M:%S)] run $n done exit=$? bytes=$(wc -c < "$SCRATCH/s4-patched-$n.out")" >> "$SCRATCH/s4-patched-run.log"
done
echo "S4 PATCHED COMPLETE" >> "$SCRATCH/s4-patched-run.log"
