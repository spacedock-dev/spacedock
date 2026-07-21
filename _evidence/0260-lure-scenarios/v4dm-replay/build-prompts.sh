#!/bin/bash
set -euo pipefail
SCRATCH=/private/tmp/claude-501/-Users-clkao-git-spacedock-research-spacedock-v1/4fc7973c-f2ca-4d34-a2b1-cf24471173af/scratchpad
WT=/Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-fanout-dedupe-ordering/skills/first-officer/references
SCEN=/Users/clkao/git/spacedock-research/spacedock-v1/docs/dev/.spacedock-state/_evidence/0260-lure-scenarios/scenarios/s6.md
HEADER='You are a Spacedock first officer. The following is your operating contract — the instruction files you have loaded for this session. Follow it.'

emit_file() {
  local path="$1"
  local name
  name=$(basename "$path")
  printf '===== BEGIN %s =====\n' "$name"
  cat "$path"
  printf '\n===== END %s =====\n' "$name"
}

# s6 prompt: shared-core + fo-dispatch-core (patched)
{
  printf '%s\n\n' "$HEADER"
  emit_file "$WT/first-officer-shared-core.md"
  printf '\n'
  emit_file "$WT/fo-dispatch-core.md"
  printf '\n===== SITUATION =====\n\n'
  cat "$SCEN"
} > "$SCRATCH/s6-patched-prompt.txt"

# s6c prompt: shared-core + fo-dispatch-core (patched) + claude-fo-dispatch
{
  printf '%s\n\n' "$HEADER"
  emit_file "$WT/first-officer-shared-core.md"
  printf '\n'
  emit_file "$WT/fo-dispatch-core.md"
  printf '\n'
  emit_file "$WT/claude-fo-dispatch.md"
  printf '\n===== SITUATION =====\n\n'
  cat "$SCEN"
} > "$SCRATCH/s6c-patched-prompt.txt"

echo "s6  prompt bytes: $(wc -c < "$SCRATCH/s6-patched-prompt.txt")"
echo "s6c prompt bytes: $(wc -c < "$SCRATCH/s6c-patched-prompt.txt")"
echo "=== patched imperative present (want 1 each) ==="
grep -c "Collapse demonstrably-identical findings in a barrier stage BEFORE the verifier fan-out" "$SCRATCH/s6-patched-prompt.txt" "$SCRATCH/s6c-patched-prompt.txt"
echo "=== claude adapter marker (want s6=0 s6c=1) ==="
grep -c "Worker Back-Channel" "$SCRATCH/s6-patched-prompt.txt" "$SCRATCH/s6c-patched-prompt.txt" || true
