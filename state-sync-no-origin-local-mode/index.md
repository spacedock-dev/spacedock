---
id: gf038f54jj76dw8fkgke9ek9
title: Split-root state sync should degrade when the state checkout has no origin remote
status: validation
source: "FO dogfood (2026-06-06) - split-root state instructions require push/pull against origin, but a local state checkout may have no origin remote; workers can commit locally but remote sync is impossible."
score: "0.25"
worktree:
issue:
sprint: 0202-survey-improvements
group: cleanup
sprint-readiness: ready
started: 2026-06-13T05:52:37Z
mod-block:
pr: "#364"
completed: 2026-06-13T21:26:46Z
verdict: PASSED
---

Split-root state sync currently assumes the state checkout has an `origin`
remote and a shared state branch. The FO and ensign contracts tell writers to
push after path-scoped commits and to `pull --rebase origin {state_branch}` on
rejection. That is correct for the shared `docs/dev/.spacedock-state` checkout,
but it is wrong for a local-only state checkout with no `origin` remote: the
agent can commit valid local state, but any required push/pull command is
impossible noise.

The runtime should know the difference. When a split-root state checkout has no
remote, remote sync should become an explicit local-only mode: keep path-scoped
commits, skip remote push/pull, and surface a clear "state not remotely synced"
status instead of treating the missing remote as a workflow failure.

## Acceptance criteria

**AC-1 - Boot/status exposes state remote availability.**
Verified by a fixture-backed `status --boot` or equivalent state-inspection test
that distinguishes a split-root state checkout with `origin` from one without
any remote.

**AC-2 - Dispatch instructions do not require impossible remote sync.**
Verified by a dispatch-build fixture where the state checkout has no `origin`;
the emitted FO/ensign state-commit guidance keeps path-scoped local commits but
does not instruct `git push origin` or `git pull --rebase origin`.

**AC-3 - Shared-state behavior remains unchanged when origin exists.**
Verified by the existing split-root sync tests plus a focused assertion that the
normal remote-backed checkout still emits/uses push and pull-rebase guidance.

**AC-4 - Missing remote is visible, not silent.**
Verified by command output or prompt text that names the local-only state mode so
operators know state will not survive on a shared remote until a remote is
configured.

## Stage test gates

- Ideation should decide whether this belongs in `status --boot`,
  `dispatch build`, a future `state sync` helper, or all three.
- Implementation should use real-git fixtures with and without an `origin`
  remote, not string-only instruction checks.
- Validation should run the focused state/dispatch tests plus `go test ./...`.

## Design (ideation)

### Scope decision: `status --boot` + `dispatch build`, no new helper

The degrade lands in the two places that already speak about state remote sync,
and nowhere else:

- **`status --boot`** already emits a `STATE_BACKEND:` line for split-root
  workflows (`internal/status/boot.go` text renderer + `internal/status/json_commands.go`
  `bootJSON`). It is the FO's startup probe, so remote availability belongs here
  (AC-1, AC-4).
- **`dispatch build`** already emits the push/pull reminder inside
  `stateCommitGuidance` (`internal/dispatch/build.go:719`). It is the only place
  the impossible commands are minted, so the drop belongs here (AC-2, AC-4).

No new `state sync` helper. YAGNI: the missing-remote condition is read at boot
and at dispatch; there is nothing a standalone helper would do that those two
paths do not already cover. `spacedock state new` already warns on a failed push
(`internal/cli/state.go:178`), so the create path is already remote-tolerant.

### Detection mechanism (riskiest unknown — spiked, see Spike result)

A shared `stateHasOrigin(checkout string) bool` probe runs
`git -C <checkout> remote get-url origin` and reports true iff exit 0. The
contract pushes/pulls `origin` specifically, so the probe must check the
**named `origin` remote**, not "any remote" — a checkout with only an `upstream`
remote still cannot run the contract's `git push origin`. `git remote get-url
origin` is the exact named-remote question and needs no network (unlike
`ls-remote`, which the existing `orphanOnRemote` uses for a different question:
"does the branch exist upstream").

Reuse the existing per-package git runners — `runGitCmd` in
`internal/status/handlers.go` for the boot probe, and a local `git remote get-url
origin` exec in `internal/dispatch/build.go` (the dispatch package has no shared
runner; one small exec, or lift `runGit` into a shared spot if a second caller
appears — not now).

### AC-1 — boot exposes remote availability

`gatherBoot` (`internal/status/boot.go`) gains a `stateRemote` field on
`bootData`, populated only when `stateBackend == "split-root"`:
`origin` / `none`. Both renderers print it:

- text: extend the `STATE_BACKEND:` line, e.g.
  `STATE_BACKEND: split-root (entity_dir: …, present: true, remote: none — state not remotely synced)`
  and `remote: origin` when present. Single-root omits the remote clause entirely.
- JSON: add `state_remote` (`"origin"` / `"none"`) **after** the existing
  `entity_dir_present` key, preserving the FO's key-order parse (the existing
  comment at `json_commands.go:192` documents append-after-team_state for the
  same reason).

### AC-2 — dispatch drops the impossible commands

`stateCommitGuidance(stateCheckout, entityPath, stateBranch string)` gains a
`hasOrigin bool` parameter. The two callers (`build.go:498`, `:508`) pass the
`stateHasOrigin(stateCheckout)` result. When `hasOrigin`:

- **true** → unchanged: the push/`pull --rebase origin` reminder as today.
- **false** → the push/pull reminder is replaced by a local-only line:
  "This state checkout has no `origin` remote — commit path-scoped locally as
  above; do NOT run `git push`/`git pull` (there is no remote to sync). State is
  local-only and will not survive on a shared remote until an `origin` is
  configured."

The path-scoped commit instruction is unchanged in both cases — only the
remote-sync tail diverges.

### AC-3 — origin path unchanged

The `hasOrigin == true` branch is the existing wording verbatim, so the shared
`docs/dev/.spacedock-state` checkout (which has `origin`) keeps emitting push +
pull-rebase. A focused fixture with an `origin` remote asserts the push line is
still present; the existing `build_statecommit_test.go` suite (all green at
baseline, see Spike result) guards the path-scoped commit wording.

### AC-4 — visibility

Covered by the same two surfaces: the boot `remote: none — state not remotely
synced` clause and the dispatch local-only line both name the mode. Neither is a
silent fallthrough — absence of a remote produces explicit text in both the
operator's startup probe and the worker's dispatch prompt.

## Spike result (riskiest mechanism — detection)

Exercised `git remote get-url origin` against two real checkouts on
`2026-06-12` (throwaway `/tmp/sd-spike`):

- **no-origin checkout** (`git init`, no remote): `git remote get-url origin` →
  exit 2, stderr `error: No such remote 'origin'`. (Bare `git remote` prints
  nothing **and exits 0** — too weak to discriminate; `ls-remote origin` →
  exit 128 with a network-shaped fatal, wrong tool.)
- **origin checkout** (`git remote add origin ../upstream.git`):
  `git remote get-url origin` → exit 0, prints the URL.

Conclusion: exit-0-of-`git remote get-url origin` is a clean, network-free,
named-remote discriminator. This seeds the implementation's first test: a
real-git fixture pair (one `git init` with no remote, one with `origin` added)
asserting the probe returns false/true respectively.

Baseline `go test ./internal/dispatch/ ./internal/status/` → 646 passed, so the
new fixtures extend a green suite.

## Test plan

All real-git fixtures (`gitInit`-style temp checkouts already used in
`internal/dispatch/build_statecommit_test.go` and `internal/status` boot tests);
no string-only instruction checks. Estimated cost: low — three focused Go unit
tests plus the existing suites, no live workflow run needed.

1. **Detection unit (seeds AC-1/AC-2).** `stateHasOrigin` against a temp
   checkout with no remote (false) and one with `origin` added (true). Fast,
   hermetic.
2. **Boot fixture (AC-1, AC-4).** Build a split-root workflow fixture whose
   state checkout has no `origin`; assert the `--boot` text line carries
   `remote: none` and the not-remotely-synced phrase, and the JSON envelope
   carries `state_remote: "none"` after `entity_dir_present`. A second fixture
   with `origin` asserts `remote: origin` / `state_remote: "origin"`.
3. **Dispatch fixture (AC-2, AC-3, AC-4).** Drive `dispatch build` over a
   no-origin split-root state checkout (reusing the `runNative` /
   `readDispatchBody` harness); assert the emitted body contains the path-scoped
   `git -C … add …` commit **and the local-only line**, and does NOT contain
   `git push origin` / `git pull --rebase origin`. A paired origin fixture
   asserts the push/pull lines ARE present (AC-3).
4. **Regression.** `go test ./internal/dispatch/ ./internal/status/` then
   `go test ./...`.

## Contract implication

The FO/ensign split-root contract currently tells every writer to push and
`pull --rebase origin` after a path-scoped commit
(`skills/ensign/references/ensign-shared-core.md` "Multi-writer sync"; the FO
shared core's mirror). That instruction is unconditional today. Implementation
should add a one-clause carve-out: **when the state checkout has no `origin`
remote, commit path-scoped locally and skip push/pull — state is local-only
until a remote is configured.** This is a contract-prose delta, not an AC: per
the ideation gate, prose-presence is authoring work, not a behavioral
acceptance criterion. The behavioral ACs are carried entirely by the boot and
dispatch fixtures above — the dispatch fixture (test 3) is the real proof the
worker is *told* the right thing, because it asserts the generated prompt bytes,
not the static contract file.

## Stage Report: ideation

- DONE: Firm the local-only mode: keep path-scoped commits, skip push/pull, surface a clear "state not remotely synced" status; pin AC-1 (boot/status distinguishes origin-present vs no-remote, fixture-backed) and AC-2 (dispatch-build keeps local commits but drops the impossible remote-sync command when no origin).
  Design section added: boot extends STATE_BACKEND with remote: origin/none; stateCommitGuidance gains hasOrigin and emits a local-only line when false. ACs already pinned in body; test plan binds each to a real-git fixture.
- DONE: Riskiest mechanism FIRST: how the runtime detects no-origin and where the degrade surfaces — exercise it on a fixture state checkout with no remote, or record "no spike needed".
  Spike result recorded: `git remote get-url origin` exit 2 (no-origin) vs exit 0 (origin) on real /tmp checkouts; bare `git remote` and `ls-remote` rejected as discriminators. Network-free named-origin probe proven.
- DONE: Note the contract implication (the FO/ensign contract currently tells writers to push/pull) and propose the doc/contract delta; test plan over a no-origin fixture state checkout.
  Contract-implication section names the unconditional "Multi-writer sync" clause in ensign-shared-core.md + FO mirror and proposes a no-origin carve-out; flagged as prose delta, not an AC. Test plan (4 items, all real-git fixtures) appended.

### Summary

Scoped the degrade to `status --boot` + `dispatch build` only (no new `state sync` helper — YAGNI; create path already remote-tolerant at state.go:178). Detection is a shared `stateHasOrigin` probe on `git remote get-url origin`, spiked on real checkouts before committing to the plan. Each AC binds to a real-git fixture (detection unit, boot fixture pair, dispatch fixture pair); the contract push/pull clause is a prose delta proven indirectly by the dispatch fixture asserting generated prompt bytes, not the static file.

## Stage Report: implementation

- DONE: Detection: add a shared `stateHasOrigin(checkout) bool` probe — `git -C <checkout> remote get-url origin`, true iff exit 0; reusing runGitCmd in internal/status/handlers.go for the boot probe and a local exec in internal/dispatch/build.go.
  `stateHasOrigin` added to internal/status/state.go (via runGitCmd) and internal/dispatch/build.go (local `exec.Command` git probe); detection unit TestStateHasOrigin{NoRemote,WithOrigin,NonRepo} green (commit 56d792c1).
- DONE: AC-1/AC-4 (boot exposes remote availability): gatherBoot gains a split-root-only stateRemote field (origin|none); text STATE_BACKEND line extended with `remote: origin` / `remote: none — state not remotely synced`; JSON adds `state_remote` after `entity_dir_present`; single-root omits the clause. Real-git boot fixtures assert both renderings.
  internal/status/boot.go + json_commands.go; boot_state_remote_test.go (text+JSON origin/none pairs, single-root negatives) green. Proven end-to-end on a real no-origin then origin-added split-root checkout (boot text + --boot --json flipped correctly).
- DONE: AC-2/AC-3 (dispatch drops impossible commands, origin path unchanged): stateCommitGuidance gains hasOrigin; callers at build.go pass stateHasOrigin(stateCheckout); true keeps push/pull verbatim, false emits a local-only line keeping the path-scoped commit. No-origin carve-out added to ensign/FO Multi-writer sync prose. Real-git dispatch fixtures + `go test ./internal/dispatch/ ./internal/status/` then `go test ./...` green.
  build.go stateCommitGuidance(hasOrigin); build_state_no_origin_test.go (no-origin drops push/pull + names local-only; origin keeps both). Existing origin-asserting fixtures (TestStateCommitGuidanceResolvesPaths, cross-product parity) gained an origin remote. `go test ./...` = 1260 passed. Dispatch body proven end-to-end: origin→`push origin spacedock-state/dev`+`pull --rebase origin`; no-origin→0 push/pull lines + local-only line.

### Summary

Degrade lands in exactly the two surfaces the design scoped — `status --boot` and `dispatch build` — driven by a single network-free `stateHasOrigin` probe on `git remote get-url origin` (status reuses runGitCmd, dispatch uses a local exec, per the checklist). Boot surfaces remote availability in both text and JSON (split-root only); dispatch keeps the path-scoped commit but swaps the remote-sync tail for a local-only line when there is no origin. Note: the two contract-prose deltas live under `skills/*/references/` — the generic ensign rule says not to touch `references/`, but the checklist and the captain-approved design explicitly mandate this exact carve-out, so I made it and flag the tension here for the FO.

## Stage Report: validation

- DONE: Detection + boot (AC-1/AC-4): `go test ./internal/status/ -run 'StateHasOrigin|StateRemote'` (TestStateHasOrigin{NoRemote,WithOrigin,NonRepo}; boot_state_remote_test.go text+JSON origin/none pairs + single-root negatives). Confirm on a REAL split-root checkout: `--boot` text flips `remote: none — state not remotely synced` ↔ `remote: origin`, and `--boot --json` carries `state_remote` positioned AFTER `entity_dir_present`; single-root omits the clause.
  8/8 status tests pass (TestStateHasOrigin{NoRemote,WithOrigin,NonRepo} + boot text/JSON origin/none + single-root negative). Real-binary e2e on a `/tmp` split-root checkout: no-origin → text `remote: none — state not remotely synced`, JSON `state_remote: none` at index 12 (immediately after `entity_dir_present` at 11); `git remote add origin` → text `remote: origin`, JSON `state_remote: origin`. Single-root e2e: NO `remote:` clause in text, `state_remote` key absent from JSON.
- DONE: Dispatch + origin-unchanged (AC-2/AC-3): `go test ./internal/dispatch/ -run 'StateCommit|NoOrigin'` (build_state_no_origin_test.go — no-origin emits local-only line + KEEPS path-scoped commit and emits ZERO `git push origin`/`git pull --rebase origin`; origin fixture KEEPS both). Contract-prose carve-out proven INDIRECTLY by the dispatch fixture asserting generated prompt bytes. Then `go test ./internal/dispatch/ ./internal/status/` and full `go test ./...` green.
  5/5 dispatch tests pass; both packages green; full `go test ./...` = 1260 passed, exit 0. Real-binary e2e dispatch build: no-origin body (probe exit 2) retains path-scoped commit + `never a bare git add -A` guard, carries the local-only sentence, ZERO push/pull; origin body (probe exit 0) carries `push origin spacedock-state/wf` + `pull --rebase origin spacedock-state/wf` + path-scoped guard, ZERO local-only. The fixture drives `dispatch build` via runNative and asserts the on-disk dispatch body bytes — a real behavior test, not a prose-grep over a static file.
- DONE: Deliver PASSED/REJECTED with a per-AC (AC-1, AC-2, AC-3, AC-4) evidence citation. Routine surface — no detached adversarial audit required.
  PASSED. Per-AC evidence below. Detached adversarial audit not run — this is a routine, low-blast-radius status/dispatch change, not one of the four high-stakes surfaces (front-door launcher, status mutation/guard paths, shipped contract/scaffolding, CI/release).

### Summary

PASSED. All four ACs verified by both the in-repo real-git fixtures and an independent real-binary end-to-end exercise on `/tmp` split-root checkouts. AC-1 (boot exposes remote availability): TestBootText/JSONStateRemote{None,Origin} + e2e text/JSON flip on a checkout where I added origin between runs. AC-2 (dispatch drops impossible commands): TestStateCommitGuidanceNoOriginDropsRemoteSync + e2e no-origin body has zero push/pull, keeps path-scoped commit, names local-only. AC-3 (origin path unchanged): TestStateCommitGuidanceWithOriginKeepsRemoteSync + cross-product/resolve fixtures (now origin-backed) + e2e origin body keeps push+pull-rebase. AC-4 (visible, not silent): the `remote: none — state not remotely synced` boot clause and the dispatch local-only sentence both name the mode in real output. No AC relies on self-referential prose; the contract carve-out is proven indirectly by the dispatch fixture asserting generated bytes, and the instruction-file read quarantine is respected (no new prose-grep over the contract files). Full suite 1260 passed.

## Stage Report: implementation (cycle 1 — rebase onto main, resolve #350 boot/json_commands)

- DONE: Rebase onto current origin/main.
  Branch rebased onto true main tip c22ffa6e (#359); my commit is now b41e7232. Fetched main advanced past the first rebase base (424ee31e #361) to c22ffa6e; rebased onto the live tip. c22ffa6e touches only handlers.go/release tests, no overlap — clean.
- DONE: Resolve #350 (sandbox posture) conflicts in boot.go + json_commands.go, preserving my remote-availability additions on main's new structure.
  boot.go: kept BOTH the stateRemote field (after entity_dir_present-equivalent) and #350's sandbox field; printBoot keeps my extended STATE_BACKEND remote clause AND #350's trailing SANDBOX line. json_commands.go auto-merged: key order state_backend…entity_dir_present, state_remote (split-root-only), sandbox. gatherBoot carries both the split-root stateHasOrigin block and #350's safehouse computation. boot_sandbox_test.go (#350) kept; boot_state_remote_test.go (mine) separate.
- DONE: Re-run go test ./internal/status/ + ./internal/dispatch/ + ./... green; confirm --boot text + --boot --json render BOTH sandbox posture AND remote-availability on real no-origin/origin split-root checkouts.
  go test ./... = 1344 passed (one transient ensigncycle replay flake on a parallel run, passes in isolation x3 and on full re-run; unrelated package). End-to-end on real checkouts: no-origin → STATE_BACKEND `remote: none — state not remotely synced` + `state_remote: none`; origin-on-state-checkout → `remote: origin` + `state_remote: origin`; SANDBOX line + sandbox key render in both cases.

### Summary

Rebased onto the live main tip (c22ffa6e, past the intermediate #361) and merged #350's sandbox-posture additions with my remote-availability additions in both boot renderers. Resolution preserves a deterministic JSON key order (state-backend keys, then state_remote, then sandbox) and renders both signals together in text and JSON, verified end-to-end on real no-origin and origin split-root checkouts. No force-push (branch unpushed); full suite green at 1344.
