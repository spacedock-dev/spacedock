---
title: "codex --plugin-dir version-masquerade advisory leaks the internal branch name into shipped CLI output"
status: done
source: "0240 pre-cut antipattern audit (2026-07-01, lens 1). internal/cli/codex_marketplace.go installCodexLocalPluginDir prints a version-masquerade advisory on every --plugin-dir codex/pi install ending '...see next-post-release-preversion-bump' — an internal roadmap/branch identifier an end user cannot act on. The advisory itself is legitimate (version stamping is deferred); the dangling internal-branch pointer is the issue. Non-blocking; captain fast-tracked."
group: tooling
id: acy7gdv88md7jgzsfea85zx6
started: 2026-07-02T01:37:54Z
worktree: .worktrees/spacedock-ensign-codex-plugin-dir-advisory-internal-ref
mod-block:
pr: pr-merge:458
verdict: passed
completed: 2026-07-02T01:54:41Z
archived: 2026-07-02T01:54:41Z
---

## Problem
`installCodexLocalPluginDir` (internal/cli/codex_marketplace.go, ~L132-142) prints a version-masquerade advisory to stderr on every `--plugin-dir` codex/pi install, ending with "...see next-post-release-preversion-bump" — an internal roadmap/branch name an end user cannot act on.

## Proposed approach
Drop the trailing `— see next-post-release-preversion-bump` clause from the advisory and end the sentence after `not necessarily its current HEAD`. End with **nothing**, not a docs pointer: the advisory is purely informational (there is no user action — the version shown is the checked-in manifest's, and the stamping fix that would change it is deferred), and no public docs page describes it, so a pointer would promise a destination that does not exist. The `version-masquerade advisory:` prefix is retained so the existing presence/absence test keys off unchanged text.

The internal branch token also appears in the same function's docstring (`codex_marketplace.go:124`, "the full stamping fix lives in next-post-release-preversion-bump"). This is not shipped output, but it is the identical dangling internal-branch reference 18 lines above the string being cleaned, and a roadmap/branch identifier in a comment is exactly the temporal-context naming we avoid. It is folded into this change as a companion edit — reworded to "the full stamping fix is deferred", preserving the deferred-fix meaning without the branch pointer. It is deliberately NOT covered by a measured AC (a comment cannot be exercised); the gate reviewer may pare it back if they prefer to scope strictly to shipped bytes.

### Reworded advisory string — before/after

Shipped advisory (`internal/cli/codex_marketplace.go`, `installCodexLocalPluginDir`):

```diff
 	fmt.Fprintf(stderr,
 		"Installed codex plugin from %s.\n"+
 			"version-masquerade advisory: the reported version reflects the checkout's "+
-			"checked-in .codex-plugin/plugin.json, not necessarily its current HEAD — "+
-			"see next-post-release-preversion-bump.\n",
+			"checked-in .codex-plugin/plugin.json, not necessarily its current HEAD.\n",
 		checkout)
```

Rendered stderr, before → after:

```
Installed codex plugin from <checkout>.
version-masquerade advisory: the reported version reflects the checkout's checked-in .codex-plugin/plugin.json, not necessarily its current HEAD — see next-post-release-preversion-bump.
```
```
Installed codex plugin from <checkout>.
version-masquerade advisory: the reported version reflects the checkout's checked-in .codex-plugin/plugin.json, not necessarily its current HEAD.
```

Companion docstring cleanup (same function, `codex_marketplace.go:121-124`):

```diff
 // invocation. It prints the version-masquerade advisory on every call: a
 // `--plugin-dir` install reports the checkout's checked-in .codex-plugin/plugin.json
-// version, not necessarily its current HEAD (the full stamping fix lives in
-// next-post-release-preversion-bump).
+// version, not necessarily its current HEAD (the full stamping fix is deferred).
```

## Acceptance criteria
Each AC names the test that proves it.

- **AC-1 (value, measured against a baseline that can move the wrong way):** Running a `--plugin-dir` codex install and capturing its stderr, the emitted advisory contains **no** occurrence of the internal branch token `next-post-release-preversion-bump`. This measures the entity's reason-for-existing (shipped output does not leak the internal branch name) against a baseline that moves the wrong way: on current `main` the assertion FAILS (the string emits the token today — confirmed at `codex_marketplace.go:142`, and the presence test proves stderr is produced on `--plugin-dir`); after the reword it passes; a regression re-introducing the token flips it back to failing.
  - **Test:** a new `strings.Contains` guard added to `TestCodexPluginDirAdvisoryPresenceAndAbsence`'s "present on --plugin-dir" subtest — `if strings.Contains(stderr.String(), "next-post-release-preversion-bump") { t.Fatalf(...) }` — asserted over the real stderr bytes `runCodex` writes, not static prose. `next-post-release-preversion-bump` is the concrete, auditable form of "no internal branch identifier"; it is the specific token the source note flags.
- **AC-2 (meaning preserved; presence/absence contract intact):** On a `--plugin-dir` install, stderr still contains both `version-masquerade advisory` and the meaning-bearing clause `not necessarily its current HEAD`; on a plain (non---plugin-dir) launch it contains neither. The reword shortens the advisory without gutting what it communicates or breaking the present-on-install / absent-on-plain-launch pair.
  - **Test:** the existing `TestCodexPluginDirAdvisoryPresenceAndAbsence` present/absent pair, still green on the reworded string (present subtest keys off `version-masquerade advisory`, which the reword retains), plus a strengthening assertion added to the present subtest that stderr contains `not necessarily its current HEAD` (so the test guards the surviving meaning, not just the label).
- **AC-3 (build + suite green):** `go build ./...` and `go test ./internal/cli/` both succeed on the reworded string and updated test.
  - **Test:** the two commands themselves.

## Test plan
- **Scope/kind:** Go unit tests only, in `internal/cli/codex_plugin_dir_test.go`, driving the command through `runCodex` with a `fakeHost` and capturing `stderr`. No new fixtures; no live codex host required for the wording change. The unrelated live edge-channel resolve test (`TestInstallCodexLocalPluginDirResolvesOnEdgeChannel`, skips when `codex` is absent) is untouched.
- **What verifies it:** AC-1's absence guard and AC-2's presence/meaning assertions run in-process over real emitted bytes; AC-3 is the build + package test run.
- **Cost/complexity:** trivial — sub-second unit tests, a two-line string edit plus a one-line docstring edit and ~2 added test assertions. No CLI/live/workflow tests needed.

## Spike
**No spike needed.** Proven mechanisms this rests on: the advisory is emitted by a plain `fmt.Fprintf(stderr, …)` in `installCodexLocalPluginDir`, and its presence/absence is already exercised by `TestCodexPluginDirAdvisoryPresenceAndAbsence` driving `runCodex` over captured stderr. Both were run green during ideation (`present on --plugin-dir`, `absent on plain launch`, and `TestInstallCodexPluginDirInstallsViaSharedHelper` all PASS), so the print-and-assert path is verified — no unverified parser round-trip, runtime handoff, or on-disk format is involved.

## Documentation
**No doc diff needed.** A repo-wide grep for `version-masquerade` hits only `.go` files (source string + test); no docs-site page describes this advisory string, so there is nothing to update. The change removes an internal reference from user-facing output and adds no new user-facing surface — output becomes strictly less internal-leaky.

## Stage Report: ideation

- DONE: Flesh the exact reworded advisory string (drop "next-post-release-preversion-bump" / any internal branch identifier; preserve the version-masquerade meaning; end user-actionable — a public docs pointer or nothing), with concrete before/after
  "## Proposed approach" + "### Reworded advisory string — before/after": ends after "…not necessarily its current HEAD." (chose nothing over a docs pointer, rationale recorded); before/after diff + rendered stderr shown; companion docstring cleanup at codex_marketplace.go:124 folded in and flagged as non-gated.
- DONE: Tighten into measured ACs each naming its test: (a) shipped advisory contains no internal branch/roadmap identifier (asserted), (b) the AC-3 presence/absence pair in codex_plugin_dir_test.go still passes on the new wording
  AC-1 measures the shipped stderr bytes contain no `next-post-release-preversion-bump` (baseline fails on current main, flips on regression) via a new `strings.Contains` guard; AC-2 keeps the present/absence pair green (prefix retained) + adds a "not necessarily its current HEAD" meaning assertion; AC-3 is build + `go test ./internal/cli/`.
- DONE: Spike: none needed for a wording change — record "no spike needed" with the proven mechanism (the existing advisory print + its presence/absence test)
  "## Spike" records no spike needed; ran the existing `fmt.Fprintf` + presence/absence mechanism green during ideation (present/absent + shared-helper all PASS), so the print-and-assert path is verified.

### Summary
Fleshed the ideation for the version-masquerade advisory internal-branch leak. Chose to end the reworded advisory with nothing (not a docs pointer) since the advisory is purely informational and no docs page describes it, and recorded concrete before/after diffs for both the shipped string and the twin docstring reference. ACs are pinned to the real shipped stderr bytes (AC-1 measures absence of the branch token against a baseline that fails on current main), and "no spike needed" is backed by a green run of the existing advisory presence/absence test. Open decision for the gate: whether the non-gated docstring companion cleanup stays in scope.

## Stage Report: implementation

- DONE: Apply the reworded advisory string + the companion docstring cleanup in internal/cli/codex_marketplace.go — drop "— see next-post-release-preversion-bump" (advisory) and "(the full stamping fix lives in next-post-release-preversion-bump)" -> "(the full stamping fix is deferred)" (docstring), per the ideation before/after; no other behavior change
  Both edits applied exactly as the ideation diff specifies; commit caf1f54e (worktree branch spacedock-ensign/codex-plugin-dir-advisory-internal-ref). `version-masquerade advisory:` prefix and `not necessarily its current HEAD` clause retained.
- DONE: Update TestCodexPluginDirAdvisoryPresenceAndAbsence: add AC-1's no-internal-branch-token guard (strings.Contains stderr, "next-post-release-preversion-bump") and AC-2's surviving-meaning assertion ("not necessarily its current HEAD"); keep the present-on-install / absent-on-plain-launch pair green
  Both assertions added to the "present on --plugin-dir" subtest; present/absent pair still green. AC-1 baseline exercised: reintroducing the token made the guard FAIL at codex_plugin_dir_test.go:100 ("advisory leaks the internal branch identifier"), then restored to green.
- DONE: go build ./... and go test ./internal/cli/ green from the worktree root; commit in the worktree
  `go build ./...` exit 0; `go test ./internal/cli/` -> ok. Committed on the worktree branch (caf1f54e). Only residual `next-post-release-preversion-bump` in source is the test guard string itself.

### Summary
Removed the internal branch identifier `next-post-release-preversion-bump` from the `--plugin-dir` version-masquerade advisory's shipped stderr and from the twin reference in the function docstring, per the ideation before/after; no other behavior changed. Guarded the shipped bytes in TestCodexPluginDirAdvisoryPresenceAndAbsence (absence of the token + survival of the "not necessarily its current HEAD" clause) with the present/absent pair kept green. The AC-1 baseline was exercised by temporarily reintroducing the token and confirming the guard fails, then reverting; build and package tests are green in the worktree.

## Stage Report: validation

- DONE: MEASURE AC-1: the shipped advisory (real stderr on a --plugin-dir install) contains NO "next-post-release-preversion-bump" — reproduce the new guard test; confirm it is baseline-fails-on-main (the token IS present on origin/main), so the assertion can move the wrong way
  Guard at codex_plugin_dir_test.go:100 passes on branch over real runCodex stderr. Composing origin/main's codex_marketplace.go with the branch guard FAILS at :101 ("advisory leaks the internal branch identifier"; emitted bytes carry "…current HEAD — see next-post-release-preversion-bump"). The assertion moves the wrong way on main and is not self-referential — it tests live emitted output, an independent source that diverged.
- DONE: AC-2: the presence/absence pair still passes (present on --plugin-dir incl. "not necessarily its current HEAD"; absent on plain launch); AC-3: go build ./... + go test ./internal/cli/ green from the worktree root
  TestCodexPluginDirAdvisoryPresenceAndAbsence PASS, both subtests (-count=1): present asserts "version-masquerade advisory" + "not necessarily its current HEAD"; absent-on-plain-launch has neither. AC-3: `go build ./...` exit 0; `go test ./internal/cli/ -count=1` → ok 21.951s exit 0.
- DONE: Confirm the ONLY change is the advisory string + its docstring (no install/launch behavior change)
  `git diff <merge-base 59251f18>..HEAD` touches only codex_marketplace.go (advisory string + docstring, both dropping the branch token) and codex_plugin_dir_test.go (added assertions). No MkdirAll / WriteCodexLocalMarketplace / Install / launch-argv lines touched. Only `next-post-release-preversion-bump` residual in any .go source is the guard string itself (test:100).

### Summary
PASSED. All three ACs verified against real behavior on branch caf1f54e. AC-1's absence guard is proven able to fail: main's source composed with the branch guard fails over emitted stderr, so the assertion tracks shipped output and flips on regression rather than matching static source. AC-2's presence/absence pair stays green with the surviving-meaning clause asserted; AC-3 build + package suite are green uncached; the diff confirms the change is confined to the advisory string + its docstring with no install/launch behavior change.
