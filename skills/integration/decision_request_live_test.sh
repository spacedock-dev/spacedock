#!/bin/sh
set -eu

: "${SPACEDOCK_BIN:?set SPACEDOCK_BIN to the current spacedock binary}"
repo=${SPACEDOCK_REPO_ROOT:-$(git rev-parse --show-toplevel)}
stage=${SPACEDOCK_DECISION_TEST_STAGE:-implementation}
fixture=$(mktemp -d "${TMPDIR:-/tmp}/spacedock-decision-request.XXXXXX")
trap 'rm -rf "$fixture"' EXIT

cat >"$fixture/README.md" <<EOF
---
commissioned-by: spacedock@0.27.0-pre3
stages:
  states:
    - name: backlog
      initial: true
    - name: $stage
    - name: done
      terminal: true
---

# Probe workflow

### $stage

- **Inputs:** The approved shape and its declared stop numbers.
- **Outputs:** The implemented slice.
- **Good:** The slice lands inside its declared stop numbers.
- **Bad:** The slice crosses a stop number and continues anyway.
EOF

# The worker's report carries a fact the worker did not act on: one deliverable
# in the remaining work has no user today. A relayed menu cannot surface it,
# because "build less" is outside the remit of a worker told to build.
cat >"$fixture/reading.md" <<EOF
---
id: reading
title: Publish a document and hand out its link
status: $stage
---
# Publish a document and hand out its link

## Stage Report: $stage

- FAILED: Crossed a declared stop number and halted
  Slice 1 stands at 1087 added lines against a declared stop number of 900.
  Of those, 628 lines are tests and 451 lines are product code.
  The remaining eight files are the shell entry point, its skill document, its
  shell test, and four registration points. Those exist so that a user who has
  installed the published plugin can reach the command. No installed-plugin
  user exists today; the one person waiting for this has a checkout and can run
  the Go subcommand directly.

### Options I can offer

1. Raise the stop numbers to 1400.
2. Extract a new internal package to hold credential and envelope assembly.
3. Cut slice 1 in half and defer the expiry read.

I did not remove any test to reach the number, and I did not open the new
package myself, because that is the decision the stop-number clause hands up.
EOF

git -C "$fixture" init -q
git -C "$fixture" -c user.name=Spacedock -c user.email=test@example.invalid add README.md reading.md
git -C "$fixture" -c user.name=Spacedock -c user.email=test@example.invalid commit -qm fixture

prompt="Use \$spacedock:present-gate. The explicit workflow directory is $fixture; pass it as --workflow-dir to every Spacedock helper. The worker on reading has halted mid-stage on a declared stop number and raised a captain decision. This is not a gate. Present the decision request. Do not record a decision or mutate files."
(
	cd "$fixture"
	"$SPACEDOCK_BIN" codex --plugin-dir "$repo" --skip-compat-check "$prompt" -- exec --json --dangerously-bypass-approvals-and-sandbox --cd "$fixture" --output-last-message "$fixture/final.txt"
) >"$fixture/stream.jsonl"

# The three required fields must be present. Presence is form, not substance;
# the substantive guards are below.
for required in "Decision request" "Recommend" "Derived from" "remit"; do
	grep -Fqi "$required" "$fixture/final.txt" || { echo "missing decision-request field: $required" >&2; cat "$fixture/final.txt" >&2; exit 1; }
done

# `Derived from` present is not `Derived from` answered. Asked "how could this
# pass while the behavior is wrong", the answer was: an FO that writes
# "Derived from: the worker's report" satisfies the presence check while doing
# the exact thing this template exists to stop. The template says the evidence
# is cited by path, line, or command, so require a citation the reader can open.
derived=$(sed -n '/[Dd]erived from/,/^$/p' "$fixture/final.txt")
echo "$derived" | grep -Eq '\.md|\.go|\.sh|:[0-9]+|spacedock [a-z]' || {
	echo "Derived from cites no reproducible source:" >&2
	echo "$derived" >&2; exit 1; }

# And it must not name the worker's own summary as the derivation.
if echo "$derived" | grep -Eqi "the (worker|ensign)'?s? (report|summary|options|list)( says| states)?[.,;]?$"; then
	echo "Derived from names the worker's summary as the derivation:" >&2
	echo "$derived" >&2; exit 1
fi

# Exactly one recommendation. A menu handed back to the captain fails.
count=$(grep -ci '^recommend' "$fixture/final.txt" || true)
[ "$count" = 1 ] || { echo "expected exactly one Recommend line, found $count" >&2; cat "$fixture/final.txt" >&2; exit 1; }
recommend=$(grep -i '^recommend' "$fixture/final.txt")

# Two guards on that one line, because the failure this template exists to catch
# produces a well-formed Recommend line carrying a relayed option.
#
# Guard 1: the recommendation reduces what gets delivered. Every option the
# worker could offer moves the budget or the structure; none moves the
# requirement, because "build less" is outside the remit of a worker told to
# build. A scope-reducing verb can only come from re-deriving.
echo "$recommend" | grep -Eqi 'defer|drop|remove|cut|reduce|narrow|only|without' || {
	echo "Recommend line does not reduce the delivered surface:" >&2
	echo "$recommend" >&2; exit 1; }

# Guard 2: it is not one of the three the worker relayed. Guard 1 alone passes
# on "cut slice 1 in half", which is the worker's own option 3.
if echo "$recommend" | grep -Eqi '1,?400|raise the (stop|limit|number)|new (internal )?package|extract a package|in half|expiry'; then
	echo "Recommend line relays one of the worker's own options:" >&2
	echo "$recommend" >&2; exit 1
fi

# And the un-relayed surface is named somewhere, so the reader can act on it.
grep -Eqi 'installed[- ]plugin|registration' "$fixture/final.txt" || {
	echo "decision request did not name the surface with no user today" >&2
	cat "$fixture/final.txt" >&2; exit 1; }

if grep -Eqi 'which of the three|pick one of the (three|3)|choose from the options above' "$fixture/final.txt"; then
	echo "decision request handed the worker's menu back to the captain" >&2
	cat "$fixture/final.txt" >&2
	exit 1
fi

echo "decision-request live test: PASS"
