#!/bin/zsh
# AC-1 loop for zqb683j8jth0tyr2eme231e2 (run the rejection journey in team mode).
#
# A run is a CONFORMING GREEN only when BOTH hold:
#   1. the focused journey test exits 0, and
#   2. it persisted a topology digest whose rows are that run's branch chain, in
#      order -- exit 0 alone never counted a chain, which is how two self-reviewing
#      single-worker runs graded green before this layer.
set -u
RUNTIME="${1:?usage: loop.sh <codex|claude> [target] [max]}"
TARGET="${2:-3}"
MAX="${3:-8}"
W=/Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-run-rejection-journey-in-team-mode
SCRATCH=/Users/clkao/.claude/jobs/4e49247e/tmp/zqb-loop
LOOP="$SCRATCH/$RUNTIME"
BIN="$SCRATCH/spacedock-under-test"
mkdir -p "$LOOP"
LEDGER="$LOOP/ledger.tsv"
[ -f "$LEDGER" ] || printf 'run\tverdict\tdigest_ok\tbranch\tchain\tconsecutive\tseconds\n' > "$LEDGER"

# The live runner resolves codex via exec.LookPath("codex"); CI puts a shim first
# on PATH pinning the model+effort the matrix judges against, so a bare local
# codex would measure a different configuration.
export SPACEDOCK_CODEX_REAL_BIN="${SPACEDOCK_CODEX_REAL_BIN:-$(command -v codex)}"
export PATH="$SCRATCH/shim:$PATH"

cd "$W" || exit 1
go build -o "$BIN" ./cmd/spacedock || { echo "BUILD FAILED"; exit 1; }

consecutive=0; run=0
if [ "$(wc -l < "$LEDGER")" -gt 1 ]; then
  run=$(( $(wc -l < "$LEDGER") - 1 ))
  consecutive=$(tail -n 1 "$LEDGER" | awk -F'\t' '{print $6+0}')
fi

while [ "$consecutive" -lt "$TARGET" ] && [ "$run" -lt "$MAX" ]; do
  run=$((run + 1))
  RUNDIR="$LOOP/run-$run"; rm -rf "$RUNDIR"; mkdir -p "$RUNDIR"
  start=$(date +%s)
  cd "$W" || exit 1
  if [ "$RUNTIME" = "codex" ]; then
    SPACEDOCK_LIVE_RUNTIME=codex SPACEDOCK_BIN="$BIN" SPACEDOCK_LIVE_ARTIFACT_DIR="$RUNDIR" \
      go test -tags live -count=1 -timeout 25m -run TestLiveCommonRejectionFlow ./internal/ensigncycle/ > "$RUNDIR/go-test.log" 2>&1
  else
    SPACEDOCK_LIVE_RUNTIME=claude SPACEDOCK_LIVE_MODEL=sonnet SPACEDOCK_BIN="$BIN" SPACEDOCK_LIVE_ARTIFACT_DIR="$RUNDIR" \
      go test -tags live -count=1 -timeout 25m -run TestLiveCommonRejectionFlow ./internal/ensigncycle/ > "$RUNDIR/go-test.log" 2>&1
  fi
  code=$?; end=$(date +%s)
  [ "$code" -eq 0 ] && verdict=PASS || verdict=FAIL

  DIGEST=$(find "$RUNDIR" -name rejection-topology.tsv 2>/dev/null | head -1)
  branch="-"; chain="-"; digest_ok=no
  if [ -n "$DIGEST" ]; then
    cp "$DIGEST" "$RUNDIR/rejection-topology.tsv" 2>/dev/null
    branch=$(awk -F'\t' '$1=="branch"{print $2}' "$DIGEST")
    chain=$(awk -F'\t' '$1!="branch"{printf "%s/%s ", $3, $4}' "$DIGEST" | sed 's/ $//')
    if [ "$branch" = "reuse" ]; then
      want="spawn/implementation done/implementation spawn/validation done/validation reuse/implementation done/implementation reuse/validation done/validation"
    else
      want="spawn/implementation done/implementation spawn/validation done/validation spawn/implementation done/implementation spawn/validation done/validation"
    fi
    [ "$chain" = "$want" ] && digest_ok=yes
  fi

  if [ "$verdict" = "PASS" ] && [ "$digest_ok" = "yes" ]; then
    consecutive=$((consecutive + 1))
  else
    consecutive=0
  fi
  printf '%s\t%s\t%s\t%s\t%s\t%s\t%s\n' "$run" "$verdict" "$digest_ok" "$branch" "$chain" "$consecutive" "$((end-start))" >> "$LEDGER"
  echo "[$RUNTIME] run $run: $verdict digest=$digest_ok branch=$branch consecutive=$consecutive ($((end-start))s)"
done
echo "[$RUNTIME] DONE runs=$run consecutive=$consecutive"
