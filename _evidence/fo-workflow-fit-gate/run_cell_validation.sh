#!/bin/bash
# Fit-gate VALIDATION drive — one cell.
#
#   usage: run_cell_validation.sh <baseline|amended> <scenario> <run-index>
#   e.g.:  run_cell_validation.sh amended s-fit3-ownerstub 1
#
# Scored by preregistration-validation.md. Read that first; it pins the arms,
# the scenarios, the scoring buckets and the pass conditions, and it was
# committed before any transcript here existed.
#
# Self-contained by design: every path resolves from this script's own location
# or from git. The ideation drive's run_cell.sh pointed at a job-scratch
# directory that no longer exists, so it cannot be re-run; this one can be run
# by anyone with the repo.
#
# Env overrides (all optional):
#   REPO         code checkout to materialize contract files from
#   BASE_REF     ref for the baseline arm and for the cross-arm constant files
#   AMENDED_REF  ref for the amended arm
#   OUTDIR       where transcripts land
set -euo pipefail

if [ $# -ne 3 ]; then
  sed -n '2,6p' "$0" >&2
  exit 2
fi
ARM="$1"; SCEN="$2"; RUN="$3"

case "$ARM" in
  baseline|amended) ;;
  *) echo "arm must be 'baseline' or 'amended', got '$ARM'" >&2; exit 2 ;;
esac

EVID="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# _evidence/fo-workflow-fit-gate -> _evidence -> .spacedock-state -> dev -> docs -> repo
REPO="${REPO:-$(cd "$EVID/../../../../.." && pwd)}"

# Pinned in preregistration-validation.md. Override only with a re-registration.
BASE_REF="${BASE_REF:-0c6a2c32a9fac9a935e52ae0f4fcacb305b1ac52}"
AMENDED_REF="${AMENDED_REF:-c9eba5db4}"
BASE_SHA_EXPECT=39b0c656e4b8a87b1a7e98295b9544c64260897a2a7084bf54c0dd1d0bdce2fd
AMENDED_SHA_EXPECT=a31c67f2c0b2182ec0641ae52301c8dac04a12e57a87e4e2a15aef76af256beb

WRITE_CORE=skills/first-officer/references/fo-write-core.md
SHARED_CORE=skills/first-officer/references/first-officer-shared-core.md
DEV_README=docs/dev/README.md

SCEN_FILE="$EVID/scenarios/${SCEN}.md"
[ -f "$SCEN_FILE" ] || { echo "no such scenario: $SCEN_FILE" >&2; exit 2; }

OUTDIR="${OUTDIR:-$EVID/transcripts-validation}"
WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT
mkdir -p "$OUTDIR"

# --- materialize the arm under test -----------------------------------------
if [ "$ARM" = baseline ]; then ARM_REF="$BASE_REF"; else ARM_REF="$AMENDED_REF"; fi
git -C "$REPO" show "${ARM_REF}:${WRITE_CORE}" > "$WORK/arm.md"

# Cross-arm constants come from BASE_REF for BOTH arms, so they are identical by
# construction and cannot explain an arm difference.
git -C "$REPO" show "${BASE_REF}:${SHARED_CORE}" > "$WORK/shared-core.md"
git -C "$REPO" show "${BASE_REF}:${DEV_README}"  > "$WORK/dev-readme.md"

# --- drift guard: report both arm digests every run -------------------------
sha_of() { git -C "$REPO" show "${1}:${WRITE_CORE}" | shasum -a 256 | cut -d' ' -f1; }
BASE_SHA="$(sha_of "$BASE_REF")"
AMENDED_SHA="$(sha_of "$AMENDED_REF")"
DRIFT=ok
[ "$BASE_SHA"    = "$BASE_SHA_EXPECT"    ] || DRIFT="BASELINE-DRIFTED"
[ "$AMENDED_SHA" = "$AMENDED_SHA_EXPECT" ] || DRIFT="AMENDED-DRIFTED"
echo "ARMS baseline=$BASE_SHA amended=$AMENDED_SHA drift=$DRIFT"
if [ "$DRIFT" != ok ]; then
  echo "Arms are not the ones pre-registered. Stop and re-register before running." >&2
  exit 3
fi

# --- assemble the prompt ----------------------------------------------------
P="$WORK/prompt.txt"
{
  echo "You are a Spacedock first officer. The following is your operating contract — the instruction files you have loaded for this session. Follow it."
  echo
  emit() { echo "===== BEGIN $1 ====="; cat "$2"; echo "===== END $1 ====="; echo; }
  emit first-officer-shared-core.md "$WORK/shared-core.md"
  emit fo-write-core.md             "$WORK/arm.md"
  emit docs/dev/README.md           "$WORK/dev-readme.md"
  cat <<'FRAMING'
===== FRAMING =====
This is a planning exercise conducted entirely in writing. You have no shell, no
file access and no tools, and you need none: the contract above and the situation
below are the complete context. Do not attempt to run commands or read files, and
do not report that tools are unavailable. State the concrete action you would take
next and why — describe the action, do not execute it.

FRAMING
  echo "===== SITUATION ====="
  cat "$SCEN_FILE"
} > "$P"

# --- run --------------------------------------------------------------------
# DRY_RUN assembles the prompt and stops without invoking a reader, so the
# harness can be exercised without spending a cell or contaminating the drive.
# Set DRY_RUN_COPY to a path to keep the assembled prompt for inspection.
if [ -n "${DRY_RUN:-}" ]; then
  if [ -n "${DRY_RUN_COPY:-}" ]; then cp "$P" "$DRY_RUN_COPY"; fi
  echo "DRY_RUN ${SCEN}-${ARM}-run${RUN} prompt_bytes=$(wc -c < "$P" | tr -d ' ') arm_bytes=$(wc -c < "$WORK/arm.md" | tr -d ' ')"
  exit 0
fi

OUT="$OUTDIR/${SCEN}-${ARM}-run${RUN}.txt"
set +e
claude -p --model opus \
  --disallowed-tools "Bash,Edit,Write,Agent,Task,NotebookEdit,WebFetch,WebSearch,Read,Grep,Glob" \
  < "$P" > "$OUT" 2>&1
RC=$?
set -e
echo "CELL_DONE ${SCEN}-${ARM}-run${RUN} exit=${RC} bytes=$(wc -c < "$OUT") -> $OUT"
