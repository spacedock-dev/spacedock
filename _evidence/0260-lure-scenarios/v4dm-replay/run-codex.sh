#!/bin/bash
SCRATCH=/private/tmp/claude-501/-Users-clkao-git-spacedock-research-spacedock-v1/4fc7973c-f2ca-4d34-a2b1-cf24471173af/scratchpad
CELL="$1"; N="$2"
P="$SCRATCH/${CELL}-patched-prompt.txt"
LOG="$SCRATCH/${CELL}-codex-run.log"
timeout 400 codex exec -m gpt-5.6-sol --sandbox read-only --skip-git-repo-check - < "$P" > "$SCRATCH/${CELL}-codex-$N.out" 2> "$SCRATCH/${CELL}-codex-$N.err"
echo "[$(date +%H:%M:%S)] ${CELL} codex run $N done exit=$? bytes=$(wc -c < "$SCRATCH/${CELL}-codex-$N.out")" >> "$LOG"
