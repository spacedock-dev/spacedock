---
title: Release process must advance the `next` edge line (reconcile + calendar bump) on every release, not just `stable`
status: validation
source: "0.24.0-pre1 prerelease cut, 2026-07-01 (Commander). Cutting v0.24.0-pre1 advanced the binary + stamped main to 0.24.0-pre1, but left the `next` edge line stale — it had diverged 40 commits during the 0240 sprint, and release.yml's hyphenated-tag carve-out advances NEITHER `stable` (correctly skipped) NOR `next`. So `spacedock install --host codex` (the spacedock-edge marketplace serves the plugin from `next`) kept serving the old 0.23.0-pre plugin, and the binary(0.24.0-pre1)/plugin(0.23.0-pre) version-compat check hard-blocked `spacedock codex`. Manually reconciled next -> main@0.24.0-pre1 (merge favoring main) + bumped the marketplace calendar key (origin/next now 1bb3da06), which unblocked the edge install. The release process should do this automatically."
group: tooling
id: s20pdb1pzexwkbp5b4cz30av
sprint: 0240-lean-contract
started: 2026-07-01T02:38:47Z
worktree: .worktrees/spacedock-ensign-release-advances-edge-next-line
---

## Problem
`release.yml` advances the `stable` ref only on a stable (non-hyphenated) `vX.Y.Z` tag; on a prerelease (hyphenated `-pre`) tag it advances NEITHER `stable` (correctly, per the carve-out) NOR the `next` edge line. So the edge channel (the `spacedock-edge` marketplace, which serves the plugin from `next`) does not advance with a prerelease cut — it drifts behind `main` until someone manually reconciles it. The 0240 sprint left `next` 40 commits behind main; the 0.24.0-pre1 cut then produced a binary(0.24.0-pre1)/edge-plugin(0.23.0-pre) skew that the strict version-compat check hard-blocked `spacedock codex`.

## Spike: proving the reconcile mechanism (riskiest path, exercised first)

The riskiest assumption is whether release.yml can advance `next` to the release
content in CI, deterministically, with no manual conflict resolution and no
`--force`. This was exercised directly rather than assumed.

**Precedent (this session's manual fix):** `96bf2243` reconciled `next`'s
40-commit-stale tip (`36bcd692`) to `main@0.24.0-pre1` (`6aa200e3`) via a merge
favoring main; `1bb3da06` then bumped the marketplace calendar key. `git diff
96bf2243 6aa200e3 --stat` is empty — the merge commit's tree is byte-for-byte
`main`'s tree, and `96bf2243`'s first parent is `36bcd692`, so pushing it was a
plain fast-forward (no force).

**Reproduction (this ideation pass):** in a disposable worktree, checked out
`36bcd692` (the pre-reconcile `next` tip, which carries 15 commits `main` never
saw, including real content changes to `plugin.json`/`codex-plugin/plugin.json`)
and ran the exact mechanism the design below automates:

```
git merge -X theirs origin/main -m "spike: reconcile edge line to main (favor main)"
```

Result: **clean, zero-conflict merge** (`Auto-merging .claude-plugin/plugin.json`,
`Auto-merging .codex-plugin/plugin.json`, `Merge made by the 'ort' strategy`,
exit 0 — no manual intervention). Verified after the merge:
- `git diff origin/main HEAD --stat` → empty (tree matches `main` exactly, same
  as the real precedent).
- `git merge-base --is-ancestor 36bcd692 HEAD` → true (the pre-merge `next` tip
  is a first-parent ancestor of the new commit), so `git push origin
  <sha>:next` is a **plain fast-forward — never `--force`**.

This confirms: `git merge -X theirs <release-commit>` run from `next`'s current
tip is a deterministic, non-interactive, force-free mechanism that reproduces
the manual precedent exactly, including on genuinely divergent history (not
just a fast-forward case). It is the mechanism the design below wires into CI.
Worktree and throwaway branch were removed after the spike; no repo state
changed.

## Design

### Prerelease (`-pre`) tag path

New sibling job in `release.yml`, `needs: goreleaser` (mirrors the
`journey-ledger` job's isolation rationale already in the file: a failure here
must not block or unwind a release that already published):

```yaml
  # Advances `next` — the branch `spacedock-edge` resolves — on every tag, not
  # just stable. See docs/releasing.md "Advancing the Edge Line".
  edge-advance:
    needs: goreleaser
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v5
        with:
          fetch-depth: 0
          fetch-tags: true
      - uses: actions/setup-go@v6
        with:
          go-version: "1.22"

      - name: Reconcile the edge line to the prerelease commit
        if: "contains(github.ref, '-')"
        run: |
          set -euo pipefail
          RELEASE_COMMIT="$(git rev-list -1 "$GITHUB_REF_NAME")"
          git fetch origin next
          git switch -c edge-advance origin/next
          git config user.name "github-actions[bot]"
          git config user.email "github-actions[bot]@users.noreply.github.com"
          git merge -X theirs --no-edit "$RELEASE_COMMIT" \
            -m "next: reconcile edge line to $GITHUB_REF_NAME"

      - name: Reconcile the edge line past the stable release
        if: "!contains(github.ref, '-')"
        run: |
          set -euo pipefail
          RELEASE_COMMIT="$(git rev-list -1 "$GITHUB_REF_NAME")"
          RELEASE_VERSION="${GITHUB_REF_NAME#v}"
          DEV_VERSION="$(go run ./cmd/spacedock-release dev-preversion "$RELEASE_VERSION")"
          git fetch origin next
          git switch -c edge-advance origin/next
          git config user.name "github-actions[bot]"
          git config user.email "github-actions[bot]@users.noreply.github.com"
          git merge -X theirs --no-edit "$RELEASE_COMMIT" \
            -m "next: reconcile edge line to $GITHUB_REF_NAME"
          go run ./cmd/spacedock-release stamp-version "$DEV_VERSION" \
            .claude-plugin/plugin.json .codex-plugin/plugin.json
          git commit -m "next: bump dev pre-version to $DEV_VERSION" \
            -- .claude-plugin/plugin.json .codex-plugin/plugin.json

      - name: Bump the marketplace calendar key and push the edge line
        run: |
          set -euo pipefail
          git config user.name "github-actions[bot]"
          git config user.email "github-actions[bot]@users.noreply.github.com"
          go run ./cmd/spacedock-release bump-calendar .claude-plugin/marketplace.json
          git commit -m "next: bump marketplace calendar version" \
            -- .claude-plugin/marketplace.json
          git push origin edge-advance:next
```

The two `if:` branches are the same mutually-exclusive condition the existing
"Stamp plugin manifests" step already uses (`contains`/`!contains(github.ref,
'-')`), so exactly one runs per tag. `bump-calendar` runs unconditionally last
(both paths need it; the stable path also needs the extra `stamp-version`
commit first). `bump-calendar` and `stamp-version` are unchanged, existing
`spacedock-release` subcommands. The only NEW piece of logic is `dev-preversion`.

### `spacedock-release dev-preversion` (new subcommand)

`internal/release/release.go` — new pure function alongside `StampVersion` /
`BumpCalendarVersion`:

```go
var stableVersionRe = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)$`)

// DevPreVersion computes the post-release dev pre-version for the `next` edge
// line from a just-released stable version: X.Y.Z -> X.(Y+1).0-pre1. Input must
// be a bare stable semver (no hyphen) — release.yml only calls this on the
// `!contains(github.ref, '-')` branch, which already guarantees that shape.
func DevPreVersion(stableVersion string) (string, error) {
	m := stableVersionRe.FindStringSubmatch(stableVersion)
	if m == nil {
		return "", fmt.Errorf("stable version %q is not X.Y.Z", stableVersion)
	}
	minor, _ := strconv.Atoi(m[2])
	return fmt.Sprintf("%s.%d.0-pre1", m[1], minor+1), nil
}
```

`cmd/spacedock-release/main.go` — new `case "dev-preversion":` dispatching to a
`devPreversion(args []string) int` that takes exactly one `<stable-version>` arg,
calls `release.DevPreVersion`, and prints the result to stdout (mirrors
`stampVersion`/`bumpCalendar`'s arg-count-then-call-then-print shape).

**Version-format note:** existing `-pre` tags are inconsistent — `v0.23.0-pre.1`
(dotted) vs the current `v0.24.0-pre1` (no dot, live in production today,
`.claude-plugin/plugin.json` currently reads `0.24.0-pre1`). `dev-preversion`
standardizes on the no-dot form (`-pre1`) to match current live practice; it
does not attempt to fix the historical dotted tags.

**Scope boundary with `next-post-release-preversion-bump`:** that backlog task's
concerns (Codex accepting a prerelease manifest version, binary-reported dev
version parity) are broader than this entity. This entity implements only the
mechanical half — computing and stamping the dev pre-version onto `next` as part
of the automated advance. `0.24.0-pre1` already installs live today, so "no
spike needed: prerelease manifest version strings are already proven in
production" for the Codex-acceptance half.

## Acceptance criteria

**AC-1 (VALUE)** — after a release cut, `next`'s commit divergence from the
release closes to zero (`git rev-list --count <release-commit>..next` counts
only the expected calendar/dev-preversion commits; `git rev-list --count
next..<release-commit>` is 0), replacing the 40-commit drift baseline that
hard-blocked `spacedock codex` this session.
Verified by: `TestEdgeLineReconcileClosesDivergence` (new,
`internal/release/`, `t.TempDir()` + `exec.Command("git", ...)`, following the
`TestAnnotatedTagBodyRoundTrips` pattern in `notes_extract_test.go`) — builds a
temp repo shaped like the real incident (`next` diverged N commits behind
`main` with some `next`-exclusive commits), runs the literal command sequence
above, and asserts the divergence count drops from N>0 to 0 and the tree
matches the release commit's tree.

**AC-2** — the reconcile never force-pushes or resets `next`, on either path.
Verified by: (a) the same fixture test asserting `git merge` exits 0 with no
`CONFLICT` markers and the pre-merge `next` tip remains a first-parent ancestor
of the pushed commit; (b) `TestReleaseWorkflowEdgeAdvanceNeverForces` (new,
`internal/release/`, following `manifest_tag_gate_workflow_test.go`'s
adversarial-mutation pattern) — greps the `edge-advance` job's `run:` blocks for
`--force`/`-f`/`reset --hard` and reds if found; the adversarial twin injects one
and asserts the guard catches it.

**AC-3** — on a prerelease tag, `next`'s plugin manifests read exactly the tag's
version and the marketplace calendar key advances past its pre-cut value.
Verified by: the AC-1 fixture test's post-merge manifest-content assertion, plus
existing `bump-calendar` unit coverage (monotonic-advance case already tested).

**AC-4** — on a stable tag, `next`'s plugin manifests read the computed
post-release dev pre-version, never the just-released stable version.
Verified by: `TestDevPreVersion` (new, table-driven, `internal/release/`):
`"0.24.0"→"0.25.0-pre1"`, `"0.9.6"→"0.10.0-pre1"`, `"1.2.3"→"1.3.0-pre1"`, error
on a hyphenated/malformed input; plus the AC-1 fixture test's stable-path
manifest-content assertion.

**AC-5** — `docs/releasing.md` documents the edge-line advance for both
prerelease and stable cuts, naming the merge-favoring-release mechanism and the
never-force invariant. Verified by: the doc diff below, reviewed at the ideation
gate (per stage-def: proposed here, applied in implementation).

## Test plan

Cost/complexity: low-to-medium. All new coverage is `go test`-speed (temp-repo
git fixtures + table-driven unit tests + text-parsing workflow guards) — no live
CI dry-run or GitHub Actions dispatch needed, matching this file's existing
`manifest-tag-gate`/`e2e-gate` test precedents. No install-smoke test is planned:
the real edge marketplace (`spacedock-dev/marketplace@edge`) is a standalone
repo already documented as "unreachable from here" (docs/releasing.md's existing
Notes section) — install-path behavior for a repointed marketplace source is
already covered by `internal/cli/marketplace_source_override_test.go` /
`channel_marketplace_source_test.go` and is not re-proven here.

1. `internal/release/edge_reconcile_test.go` (new) — `TestEdgeLineReconcileClosesDivergence`, `TestReleaseWorkflowEdgeAdvanceNeverForces` (AC-1, AC-2).
2. `internal/release/release_test.go` (extend) — `TestDevPreVersion` table test (AC-4).
3. `internal/release/edge_advance_workflow_test.go` (new) — wiring guard: `edge-advance` job exists, `needs: goreleaser`, both `if:` branches present and mutually exclusive, `bump-calendar` invoked unconditionally (AC-3, AC-4 wiring half).
4. `cmd/spacedock-release/main_test.go` or a new `dev_preversion_test.go` — CLI arg-count/usage coverage for `dev-preversion`, mirroring existing subcommand tests.

## Documentation diff (docs/releasing.md)

Append a bullet to "## What the Tag Push Does":

```diff
 - stamps the plugin manifests' `version` on `main`, then advances the stable
   channel ref (see below).
+- advances the `next` edge line to match the release — reconciled on a
+  prerelease tag, reconciled plus bumped to the post-release dev pre-version on
+  a stable tag (see "Advancing the Edge Line" below).
```

New section, inserted after step 9 of "## Cutting a Stable Release" and before
"## Dev-Only `next` Publishing":

```diff
+## Advancing the Edge Line (`next`)
+
+Every tag push also advances `next` — the branch the `spacedock-edge`
+marketplace resolves — in a job that `needs: goreleaser` (a sibling job, so a
+rare conflict here cannot unwind or block the release that already published):
+
+- **Prerelease (`-pre`) tag:** `next` is reconciled to the tagged commit's
+  content — `git merge -X theirs "$RELEASE_COMMIT"`, favoring the release over
+  whatever `next` had drifted to — then the marketplace calendar key is bumped
+  (`spacedock-release bump-calendar`) so `claude plugin update` / `codex`
+  re-pull. This is the automated form of the manual reconcile the 0.24.0-pre1
+  cut required (`next` had drifted 40 commits behind `main`, hard-blocking
+  `spacedock codex` on a binary/plugin version-compat check).
+- **Stable (`vX.Y.Z`) tag:** `next` is reconciled the same way, then stamped
+  PAST the release to the post-release dev pre-version
+  (`spacedock-release dev-preversion X.Y.Z` → `X.(Y+1).0-pre1`), so the edge
+  line never masquerades as the stable version it just shipped, then the
+  calendar key is bumped.
+
+The reconcile is a merge, never a reset or force-push: the previous `next` tip
+is always a first-parent ancestor of the new commit, so `git push origin
+<sha>:next` is a plain fast-forward. A real conflict (two sides changing the
+same file in incompatible, non-superseded ways) fails the step loudly instead
+of guessing — the same manual reconciliation this replaces remains the escape
+hatch.
+
 ## Dev-Only `next` Publishing
```

Append a sentence to "## Dev-Only `next` Publishing":

```diff
 Keep `next` for development. Source builds may use
 `go install github.com/spacedock-dev/spacedock/cmd/spacedock@next`, local
 checkouts may use `--plugin-dir`, and the deliberate `next-publish` workflow may
 bump the marketplace calendar key for dev testers.
+
+Every release tag now advances `next` and bumps its calendar key automatically
+(see "Advancing the Edge Line" above); `next-publish` stays for an out-of-band
+re-pull between releases (e.g. a `next`-only fix that isn't worth a full cut).
```

## Related
- `next-post-release-preversion-bump` — broader dev/edge line coherence (Codex
  acceptance, binary-reported dev version); this entity implements only the
  mechanical dev-preversion stamp + edge-ref advance (see Scope boundary above).
- `minor-version-compat-coupling` — why the skew hard-blocks; a laxer
  contract-compatible check softens the failure mode.
- The 0.24.0-pre1 cut (this session) that exposed it; the manual reconcile
  (`96bf2243` merge + `1bb3da06` calendar bump, `origin/next` @ 1bb3da06) is the
  proof precedent the Spike section reproduces.
- `release.yml` prerelease carve-out (`if: "!contains(github.ref, '-')"`),
  `next-publish.yml` (calendar bump, stays for out-of-band re-pulls).

## Stage Report: ideation

- DONE: Flesh the next-line-advance design for both prerelease and stable-tag paths, with concrete release.yml / spacedock-release before/after
  New `edge-advance` sibling job (both `if:` branches) plus a new `spacedock-release dev-preversion` subcommand + `release.DevPreVersion` function, spec'd in the "Design" section.
- DONE: Tighten the rough acceptance sketch into measured ACs each naming its test, plus the docs/releasing.md diff
  AC-1..AC-5 in "Acceptance criteria" (AC-1 is the VALUE ac: divergence count N>0 -> 0); concrete diff hunks in "Documentation diff".
- DONE: Spike the riskiest mechanism first — can release.yml advance next in CI safely without a force-push footgun
  Reproduced in a disposable worktree (`36bcd692` + `git merge -X theirs origin/main`): clean 0-conflict merge, resulting tree byte-identical to main's tree, old `next` tip stays a first-parent ancestor of the new commit (fast-forward push, no `--force`) — matches the real `96bf2243`/`1bb3da06` precedent exactly. Worktree and throwaway branch removed after.

### Summary

Reproduced the manual reconcile mechanism empirically (`git merge -X theirs`) in a throwaway worktree before designing around it, confirming it is deterministic, conflict-free even on genuinely divergent history, and never requires `--force`. Designed a single new `edge-advance` CI job covering both the prerelease (reconcile + calendar bump) and stable (reconcile + dev-preversion stamp + calendar bump) paths, adding one new pure function (`DevPreVersion`) while reusing existing `stamp-version`/`bump-calendar`. Tightened the AC sketch into 5 measured ACs each naming a specific new or existing test (fixture git-repo tests following `notes_extract_test.go`'s pattern, a workflow-guard adversarial test following `manifest_tag_gate_workflow_test.go`, and a table-driven unit test), and wrote the concrete docs/releasing.md diff.

## Stage Report: implementation

- DONE: Add the edge-advance job to release.yml (needs: goreleaser; prerelease -pre branch: reconcile via git merge -X theirs + calendar bump; stable branch: reconcile + dev-preversion stamp + calendar bump), force-free (fast-forward push only, never --force/reset)
  `edge-advance` sibling job added after goreleaser (commit 978aed89); both mutually-exclusive `if:` branches present, calendar bump + `git push origin edge-advance:next` (fast-forward) run unconditionally last.
- DONE: Add release.DevPreVersion (X.Y.Z -> X.(Y+1).0-pre1, error on hyphenated/malformed) + the cmd/spacedock-release dev-preversion subcommand, reusing stamp-version/bump-calendar unchanged
  `DevPreVersion` in internal/release/release.go (regexp-guarded); `dev-preversion` case + `devPreversion()` in cmd/spacedock-release/main.go print the value to stdout; stamp-version/bump-calendar untouched.
- DONE: Land the 4 tests (edge divergence-closes fixture, never-forces workflow guard, DevPreVersion table, dev-preversion CLI) + apply the docs/releasing.md diff; go build ./... and go test ./internal/release/ green
  `TestEdgeLineReconcileClosesDivergence` (both paths) + `TestReleaseWorkflowEdgeAdvanceNeverForces` (2 adversarial twins) in edge_reconcile_test.go; `TestDevPreVersion` in release_test.go; CLI tests in dev_preversion_test.go; docs/releasing.md diff applied. `go build ./...` exit 0; `go test ./internal/release/ ./cmd/spacedock-release/` green; gofmt/vet clean.

### Summary

Added an `edge-advance` sibling job (`needs: goreleaser`) that reconciles `next`
to the tagged commit via `git merge -X theirs` on every tag — prerelease paths
reconcile + bump calendar, stable paths also stamp `next` past the release to
the dev pre-version — with the push always a fast-forward, never a force. The
only new logic is the pure `DevPreVersion` (X.Y.Z -> X.(Y+1).0-pre1) plus its
`dev-preversion` CLI wrapper; `stamp-version`/`bump-calendar` are reused
unchanged. The temp-repo fixture exercises the real reconcile mechanism on both
paths (genuinely-diverged `next`, clean zero-conflict merge, tree byte-matching
the release, divergence closing to zero, pre-merge tip staying a first-parent
ancestor), and the workflow guard reds if any edge-advance step force-pushes or
resets. Notable decision: the never-forces guard also asserts `needs:
goreleaser` (the sibling-isolation invariant), and the fixture exercises both
`if:` branches end-to-end rather than text-grepping the conditions.

## Stage Report: validation

- DONE: MEASURE each AC's end-value against its "Verified by" test — reproduce, do not assert: AC-1 (edge divergence count drops N>0 to 0 + tree matches release commit), AC-2 (no --force/reset; pre-merge next tip stays first-parent ancestor), AC-4 (DevPreVersion table: 0.24.0->0.25.0-pre1 etc + error on hyphenated)
  AC-1/AC-2a reproduced two ways: `TestEdgeLineReconcileClosesDivergence` (both paths PASS fresh) AND an independent throwaway repo — behind_before=3/ahead_before=2, conflict=no, tree_matches_release=yes, behind_after=0, pre-tip-is-ancestor-of-pushed=yes. Real precedent independently confirmed: `git diff 96bf2243 6aa200e3 --stat` empty, `96bf2243^1`=36bcd692, 36bcd692 ancestor of 96bf2243 (fast-forward). AC-2b: `TestReleaseWorkflowEdgeAdvanceNeverForces` PASS — parses real release.yml via yaml.v3, adversarial twins inject `--force`/`reset --hard` and the guard reds (non-vacuous). AC-4: `TestDevPreVersion` PASS + real CLI `go run ... dev-preversion` emits 0.25.0-pre1 / 0.10.0-pre1 / 1.3.0-pre1 and exits 1 on `0.24.0-pre1`.
- DONE: go build ./... and go test ./internal/release/ ./cmd/spacedock-release/ green from the repo root
  From the worktree root: `go build ./...` exit 0; `go test -count=1 ./internal/release/ ./cmd/spacedock-release/` both ok (release-suite 1.18s, CLI-suite 0.14s, all subtests PASS). gofmt clean on all 5 changed Go files; `go vet` clean.
- DONE: Confirm the release.yml edge-advance job wiring (needs: goreleaser, both mutually-exclusive if: branches, bump-calendar unconditional-last) and the docs/releasing.md "Advancing the Edge Line" diff are applied as designed
  release.yml: `edge-advance` `needs: goreleaser` (L259); prerelease step `if: "contains(github.ref, '-')"` (L279) and stable step `if: "!contains(github.ref, '-')"` (L294) are logical negations (mutually exclusive); calendar-bump+push step has no `if:` and is last (L315-323). docs/releasing.md: "Advancing the Edge Line (`next`)" section (L154) documents both paths, names `git merge -X theirs`, and the never-force invariant ("never a reset or force-push... first-parent ancestor... plain fast-forward", L173-175); plus the tag-push bullet (L19-21) and Dev-Only sentence (L187-189). Matches the design.

### Summary

PASSED. All five ACs reproduced by exercising real behavior, not assertion: the reconcile mechanism (AC-1/AC-2a) verified both via the fixture test and an independent throwaway-repo run plus the real-repo precedent, closing divergence 3->0 with a byte-matching tree and a fast-forward-safe pushed tip on both prerelease and stable paths; the never-forces guard (AC-2b) parses the real workflow and its adversarial twins prove it can fail; DevPreVersion (AC-4) confirmed via table test and the real CLI. No AC is self-referential or vacuous. Non-blocking note: the test-plan's optional separate `edge_advance_workflow_test.go` (mutual-exclusivity / unconditional-bump wiring guard) was not landed, but no AC's "Verified by" clause depends on it — `needs: goreleaser` is guarded by the never-forces test, and the mutual-exclusive `if:` branches + unconditional-last bump-calendar are confirmed by inspection and exercised end-to-end by the fixture running both paths.

### Feedback Cycles

**Cycle 1 (2026-07-01) — validation REJECTED by the detached adversarial audit (2 material holes; routed to implementation).** The validator recommended PASSED but the required detached adversarial audit (`s2-detached-adversarial-audit`, 4 lenses on throwaway checkouts of the s2 tip) refuted its "non-blocking" wiring note by exercising it: two claim-breaking mutations to the deliverable's guarded surfaces stayed GREEN.

- **M-1 (merge-strategy coupling gap) — `TestEdgeLineReconcileClosesDivergence` / `internal/release/edge_reconcile_test.go:130`.** Adversarial edit: flip BOTH edge-advance reconcile steps in `.github/workflows/release.yml` from `git merge -X theirs` to `git merge -X ours` (the exact 0.24.0-pre1 divergence-incident reintroduction). The fixture stayed GREEN even at `-count=1`, because it hand-codes its OWN `exec.Command("git","merge","-X","theirs",...)` re-implementation rather than extracting the strategy from release.yml — no assertion ties the fixture's `-X theirs` to the workflow's. (Probe A confirmed the fixture's tree-match assertion is non-vacuous on its own path: flipping the fixture's own merge ref reds at `:307`.) The doc-comment falsely claims it "runs the exact reconcile sequence release.yml's edge-advance job runs." **Fix:** in the `readWorkflow`/`os.ReadFile(../../.github/workflows/release.yml)` harness `TestReleaseWorkflowEdgeAdvanceNeverForces` already establishes, assert both reconcile steps' merge command contains `-X theirs` (rejects `-X ours`); better, derive the strategy string from the parsed workflow step so any drift reds the fixture directly. Correct the overclaiming doc-comment.
- **M-2 (mutual-exclusivity unguarded) — release.yml wiring.** Adversarial edit: change the stable-reconcile step's `if: "!contains(github.ref, '-')"` to `if: "always()"` so BOTH branches fire on a prerelease (the stable dev-preversion stamp wrongly runs on prereleases). Full suite stayed GREEN — `assertEdgeAdvanceIsForceFreeSibling` inspects job existence + `needs: goreleaser` + per-step push/reset only, never the `if:` conditions or the unconditional bump-calendar. "Exactly one branch per tag" + "bump-calendar runs unconditionally last" are prose-only. **Fix:** extend the same `readWorkflow`-based harness to assert complementary `if:` guards (prerelease `contains(github.ref,'-')`, stable `!contains(...)`) and that bump-calendar carries no `if:`, asserting on the literal `if:` strings from the on-disk workflow so `always()`/copy-paste conditions red.

Two lenses refuted nothing (kept as-is): `TestReleaseWorkflowEdgeAdvanceNeverForces` (never-force; RED under `--force`/`reset --hard`, non-vacuous twins) and `TestDevPreVersion` (minor-bump; concrete want-values + `bad[]` error-guard slice). Consolidate all edge-advance step-level wiring assertions into the one `readWorkflow` harness so future workflow edits are guarded uniformly.

## Stage Report: implementation (cycle 1)

- DONE: M-1 — couple the divergence fixture to the real release.yml (kill the -X theirs coupling gap)
  `runEdgeReconcileFixture` now derives its merge `-X <strategy>` from the parsed on-disk edge-advance reconcile step (`edgeAdvanceReconcileStep`/`mergeStrategyOption`) instead of hardcoding `theirs`; corrected the overclaiming doc-comment. Confirmed on-disk `-X theirs`→`-X ours` (both steps) now reds the fixture tree-match (edge_reconcile_test.go:319, both paths) AND the wiring guard. Commit 0683c11e.
- DONE: M-2 — guard the mutual-exclusivity + unconditional bump-calendar wiring
  New `assertEdgeAdvanceWiring` asserts both reconcile steps merge `-X theirs`, complementary `if:` guards (prerelease `contains(github.ref,'-')`, stable `!contains(...)` → exactly one fires), and bump-calendar carries no `if:`. Confirmed on-disk stable-step `if:`→`always()` now reds the wiring guard ("has if: always(), want contains/!contains").
- DONE: Verify + consolidate
  `TestReleaseWorkflowEdgeAdvanceWiring` carries 4 adversarial twins (-X ours, always(), copy-paste guard, gated bump) each proven to red; all edge-advance wiring assertions share the one `readWorkflow` harness + `edgeAdvanceJob` lookup. `go build ./...` exit 0; `go test -count=1 ./internal/release/ ./cmd/spacedock-release/` green; gofmt + go vet clean. The two clean lenses (`TestReleaseWorkflowEdgeAdvanceNeverForces`, `TestDevPreVersion`) unchanged and still pass.

### Summary

Closed both material coupling holes from the cycle-1 audit. The fixture no longer re-implements a hardcoded `-X theirs`; it reads the strategy from the on-disk workflow step, so the exact 0.24.0-pre1 incident (`-X ours`) now reds it directly — verified by flipping release.yml on disk (fixture tree-match reds on both paths). The new wiring guard, in the same `readWorkflow` harness, locks the merge strategy, the complementary `if:` guards, and the unconditional bump-calendar; the `always()` widening now reds — also verified on disk. Mutations confirmed red: `-X ours` (fixture + wiring), `always()` on the stable step (wiring); the twins additionally exercise copy-paste and gated-bump. release.yml was reverted clean after each on-disk probe.

## Stage Report: validation (cycle 1)

- DONE: MEASURE each AC against its "Verified by" test — reproduce, do not assert: AC-1 (edge divergence N>0 to 0 + tree matches release), AC-2 (never --force/reset; first-parent-ancestor push), AC-4 (DevPreVersion table + error on hyphenated)
  AC-1/AC-2a/AC-3: `TestEdgeLineReconcileClosesDivergence` PASS both paths (behind closes to 0, tree byte-matches release, pre-merge tip stays first-parent ancestor, no CONFLICT, calendar advances). AC-2b: `TestReleaseWorkflowEdgeAdvanceNeverForces` PASS (parses real release.yml; `--force`/`reset --hard` twins red — non-vacuous). AC-4: `TestDevPreVersion` PASS + real CLI `go run ... dev-preversion` emits 0.25.0-pre1 / 0.10.0-pre1 / 1.3.0-pre1 and exits 1 on `0.24.0-pre1`.
- DONE: VERIFY the feedback-cycle-1 fixes bite on disk — reproduce BOTH mutations going red
  (a) M-1: flipped on-disk release.yml `git merge -X theirs`→`-X ours` (both steps) — `TestEdgeLineReconcileClosesDivergence` reds the tree-match at edge_reconcile_test.go:319 on BOTH paths AND `TestReleaseWorkflowEdgeAdvanceWiring` reds ("merges with -X \"ours\", want theirs"), confirming the fixture now derives its strategy from the workflow. (b) M-2: flipped the stable step `if: "!contains(...)"`→`"always()"` — `TestReleaseWorkflowEdgeAdvanceWiring` reds ("has if: \"always()\", want contains/!contains"). release.yml reverted clean after each probe (`git status --porcelain` empty).
- DONE: go build ./... and go test -count=1 ./internal/release/ ./cmd/spacedock-release/ green from the repo root; gofmt + go vet clean
  From the worktree root: `go build ./...` exit 0; `go test -count=1 ./internal/release/ ./cmd/spacedock-release/` both ok. `gofmt -l` clean on all changed Go files; `go vet` clean.

### Summary

PASSED (recommendation). Re-validating after the cycle-1 rejection: both material audit holes are proven closed by reproduction on the real on-disk release.yml, not by assertion. M-1 — the divergence fixture now derives its `-X` strategy from the parsed workflow, so the exact 0.24.0-pre1 incident (`-X ours`) reds it directly on both paths (plus the wiring guard). M-2 — the new `TestReleaseWorkflowEdgeAdvanceWiring` locks the merge strategy, complementary `if:` guards, and the unconditional bump-calendar; widening the stable `if:` to `always()` reds it. Both wiring/fixture checks parse the real workflow through yaml.v3 and compare against independent expected values (can genuinely diverge — not spelling checks). All five ACs measured against real behavior; none self-referential or vacuous. Build/test/gofmt/vet all green; working tree left clean.
