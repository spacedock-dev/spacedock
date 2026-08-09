#!/bin/sh
set -eu

: "${SPACEDOCK_BIN:?set SPACEDOCK_BIN to the current spacedock binary}"
repo=${SPACEDOCK_REPO_ROOT:-$(git rev-parse --show-toplevel)}
stage=${SPACEDOCK_GATE_TEST_STAGE:-orbit-audit}
mode=${SPACEDOCK_GATE_TEST_MODE:-explicit}
gate_content='- **Gate content:** Show only the median orbit drift and that approval releases the observation.'
[ "$mode" = fallback ] && gate_content=
fixture=$(mktemp -d "${TMPDIR:-/tmp}/spacedock-gate-content.XXXXXX")
trap 'rm -rf "$fixture"' EXIT

cat >"$fixture/README.md" <<EOF
---
commissioned-by: spacedock@0.27.0-pre3
entity-type: observation
entity-label: observation
entity-label-plural: observations
id-style: slug
stages:
  defaults: {worktree: false, concurrency: 1}
  states:
    - name: $stage
      initial: true
      gate: true
    - name: released
      terminal: true
---
# Orbital workflow

### $stage

- **Inputs:** Raw frames from the sensor array.
- **Outputs:** A full anomaly catalog.
$gate_content
- **Good:** The median orbit drift is below 5 ms.
- **Bad:** The median orbit drift is 5 ms or more.
EOF
cat >"$fixture/reading.md" <<EOF
---
id: reading
title: Orbital reading
status: $stage
---
# Orbital reading

## Stage Report: $stage

- DONE: Measured median orbit drift
  The median orbit drift is 4.2 ms.

### Summary

The full anomaly catalog contains 19 entries from raw frames.
The operator's favorite color is blue.
EOF

git -C "$fixture" init -q
git -C "$fixture" -c user.name=Spacedock -c user.email=test@example.invalid add README.md reading.md
git -C "$fixture" -c user.name=Spacedock -c user.email=test@example.invalid commit -qm fixture

prompt="Use \$spacedock:present-gate. The explicit workflow directory is $fixture; pass it as --workflow-dir to every Spacedock helper. Present the gate for reading at $stage. The reviewed Briefing is briefing:orbit:1 with digest sha256:abc123; recommend approve. Do not record a decision or mutate files."
(
	cd "$fixture"
	"$SPACEDOCK_BIN" codex --plugin-dir "$repo" --skip-compat-check "$prompt" -- exec --json --dangerously-bypass-approvals-and-sandbox --cd "$fixture" --output-last-message "$fixture/final.txt"
) >"$fixture/stream.jsonl"

for required in "$stage" "4.2 ms" "briefing:orbit:1" "sha256:abc123" approve release; do
	grep -Fq "$required" "$fixture/final.txt" || { echo "missing gate evidence: $required" >&2; exit 1; }
done
forbidden='favorite color|blue|SKIPPED|FAILED|None|N/A|0 skipped|0 failed'
[ "$mode" = explicit ] && forbidden="raw frames|anomaly catalog|$forbidden"
if grep -Eqi "$forbidden" "$fixture/final.txt"; then
	echo "gate review contains undeclared or empty evidence" >&2
	cat "$fixture/final.txt" >&2
	exit 1
fi
grep -Fq 'dispatch show-stage-def' "$fixture/stream.jsonl" || { echo "presenter did not fetch the stage definition" >&2; exit 1; }
