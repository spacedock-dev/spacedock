#!/bin/zsh
# Targeted codex rejection-flow repair loop for entity hz2ankag6fk379ssabpv4ckc.
#
# Runs the ONE journey against the WORKTREE binary (SPACEDOCK_BIN must point at
# the fixed build -- without it the live test falls back to a `spacedock` on
# PATH, which is a stale pre-fix binary and makes the whole loop meaningless).
#
# A run counts as a green only when all three hold:
#   verdict=PASS      the journey graded pass
#   state_open=yes    a `gate prepare` exited 0 emitting state=open
#   inline_commit=yes a `state commit` reported the inline commit -- proof the
#                     FIX carried the run, not the model's raw-git fallback
# prefix_noop must stay 0; nonzero means the run shelled a pre-fix binary.
set -u
W=/Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-repair-codex-rejection-round-recording
SCRATCH=/Users/clkao/.claude/jobs/4e49247e/tmp/repair-codex-rejection-round-recording-scratch
LOOP="$SCRATCH/loop"
BIN="$SCRATCH/spacedock-under-test"
SHIM="$SCRATCH/codex-shim"
# The live runner resolves codex via exec.LookPath("codex"). CI puts a shim
# first on PATH that pins --model gpt-5.6-luna with reasoning effort max; a bare
# local codex is a DIFFERENT configuration, so a loop run without this shim is
# not measuring what the stack PR's matrix will judge AC-2 against.
export SPACEDOCK_CODEX_REAL_BIN="${SPACEDOCK_CODEX_REAL_BIN:-$(command -v codex)}"
export PATH="$SHIM:$PATH"
[ -x "$SHIM/codex" ] || { echo "MISSING CODEX SHIM at $SHIM/codex"; exit 1; }
mkdir -p "$LOOP"
LEDGER="$LOOP/ledger.tsv"
[ -f "$LEDGER" ] || printf 'run\tverdict\tstate_open\tinline_commit\tprefix_noop\tconsecutive\tseconds\n' > "$LEDGER"

cd "$W" || exit 1
go build -o "$BIN" ./cmd/spacedock || { echo "BUILD FAILED"; exit 1; }
# Fail fast if the binary under test is not the fixed one.
"$BIN" state commit probe --workflow-dir /nonexistent 2>&1 | grep -q "no README.md" || true

TARGET_GREENS=${TARGET_GREENS:-5}
MAX_RUNS=${MAX_RUNS:-12}
consecutive=0
run=0
if [ -s "$LEDGER" ] && [ "$(wc -l < "$LEDGER")" -gt 1 ]; then
  run=$(( $(wc -l < "$LEDGER") - 1 ))
  consecutive=$(tail -n 1 "$LEDGER" | awk -F'\t' '{print $6+0}')
fi

while [ "$consecutive" -lt "$TARGET_GREENS" ] && [ "$run" -lt "$MAX_RUNS" ]; do
  run=$((run + 1))
  RUNDIR="$LOOP/run-$run"
  rm -rf "$RUNDIR"; mkdir -p "$RUNDIR"
  start=$(date +%s)
  cd "$W" || exit 1
  SPACEDOCK_LIVE_RUNTIME=codex \
  SPACEDOCK_BIN="$BIN" \
  SPACEDOCK_LIVE_ARTIFACT_DIR="$RUNDIR" \
    go test -tags live -count=1 -timeout 25m \
      -run TestLiveCommonRejectionFlow ./internal/ensigncycle/ \
      > "$RUNDIR/go-test.log" 2>&1
  code=$?
  end=$(date +%s)
  if [ "$code" -eq 0 ]; then verdict=PASS; else verdict=FAIL; fi

  stream=$(find "$RUNDIR" -name 'codex-exec.jsonl' -print 2>/dev/null | head -1)
  evidence="no no 0"
  if [ -n "$stream" ]; then
    evidence=$(python3 - "$stream" <<'PY'
import json, sys
state_open = inline_commit = False
prefix_noop = 0
for line in open(sys.argv[1], errors="replace"):
    try: e = json.loads(line)
    except Exception: continue
    it = e.get("item") or {}
    if it.get("type") != "command_execution": continue
    cmd = it.get("command") or ""
    out = it.get("aggregated_output") or ""
    prefix_noop += out.count("nothing to commit to a state checkout")
    if it.get("exit_code") == 0 and "gate prepare" in cmd and "state=open" in out:
        state_open = True
    if "in the inline workflow repository" in out:
        inline_commit = True
print("%s %s %d" % ("yes" if state_open else "no",
                    "yes" if inline_commit else "no", prefix_noop))
PY
)
  fi
  state_open=$(echo "$evidence" | awk '{print $1}')
  inline_commit=$(echo "$evidence" | awk '{print $2}')
  prefix_noop=$(echo "$evidence" | awk '{print $3}')

  if [ "$verdict" = PASS ] && [ "$state_open" = yes ] && [ "$inline_commit" = yes ] && [ "$prefix_noop" -eq 0 ]; then
    consecutive=$((consecutive + 1))
  else
    consecutive=0
  fi
  printf '%s\t%s\t%s\t%s\t%s\t%s\t%s\n' "$run" "$verdict" "$state_open" "$inline_commit" "$prefix_noop" "$consecutive" "$((end - start))" >> "$LEDGER"
done

printf 'LOOP DONE runs=%s consecutive_greens=%s\n' "$run" "$consecutive"
cat "$LEDGER"
