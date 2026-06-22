---
id: fea266y405b95053aq86q5d8
title: Codex edge install — per-channel marketplace source (edge branch + channel-specific binary source)
status: validation
source: handoff-codex-edge-0230 (z2 spacedock-marketplace-source-env follow-up)
sprint: 0230-stable-finalization
started: 2026-06-22T00:36:49Z
completed:
verdict:
score: 0.7
worktree: .worktrees/spacedock-ensign-codex-edge-channel-marketplace-source
issue:
---

The codex edge channel cannot install. The edge binary builds plugin id `spacedock@spacedock-edge`, but `codex plugin marketplace add spacedock-dev/marketplace` registers a marketplace NAMED `spacedock` (the repo's root `marketplace.json` `name`, which carries `spacedock-edge` only as an *entry*), so `codex plugin add spacedock@spacedock-edge` fails: `plugin spacedock was not found in marketplace spacedock-edge`. z2 shipped the binary half but validated its marketplace half only against a synthesized live-lane fixture — the real repo + real source were never wired, which masked this. Blocks pre.3 and v0.23.0.

## Problem

The binary's `channelMarketplace(edge)` resolves marketplace name `spacedock-edge`, but `internal/cli/init.go:22 marketplaceSource = "spacedock-dev/marketplace"` is a single var used for BOTH channels. The bare source's root `marketplace.json` registers under the name `spacedock`, so an edge install adds a marketplace named `spacedock`, then looks up `spacedock@spacedock-edge` and fails. Stable is unaffected (entry `spacedock` in marketplace `spacedock` → `spacedock@spacedock` matches).

## Proposed approach (captain pre-approved — handoff)

1. **edge branch in `spacedock-dev/marketplace`** (NOT a new repo). Its root `.claude-plugin/marketplace.json` is NAMED `spacedock-edge`, with a single entry NAMED `spacedock` whose `source` points at the product repo `spacedock.git` ref `next` (model the entry source on the current main `marketplace.json`'s `spacedock-edge` entry: url `https://github.com/spacedock-dev/spacedock.git`, ref `next`).
2. **VERIFY the mechanism FIRST (riskiest path — spike before any binary change):** `codex plugin marketplace add spacedock-dev/marketplace@edge` must register a DISTINCT marketplace named `spacedock-edge` (alongside the existing `spacedock`), and `codex plugin add spacedock@spacedock-edge` must SUCCEED (entry `spacedock` == the plugin's own `.codex-plugin/plugin.json` name → name-match passes). Record the transcript. If `owner/repo@ref` does NOT yield a distinct-named marketplace, STOP and report to the captain — the fallback is a separate `spacedock-dev/marketplace-edge` bare-source repo, which is the captain's call; do NOT create a new repo unilaterally.
3. **Binary: channel-specific marketplace source.** Make the edge install source `spacedock-dev/marketplace@edge` while stable stays `spacedock-dev/marketplace`; keep the `SPACEDOCK_MARKETPLACE_SOURCE` override working for both. TDD.
4. **Validate against the REAL marketplace** (an actual `codex plugin add` / `spacedock install --host codex` from the published edge), NOT a synthesized fixture — the fixture is exactly what masked this.

## Out of scope

- Cutting pre.3 / v0.23.0 (separate release step, gated on this fix).
- The Claude channel (already installs correctly).
- Re-pinning the edge entry version policy (tracked elsewhere); the edge entry tracks `next` HEAD.

## Acceptance criteria

**AC-1 — A real codex edge install succeeds end-to-end (end-value, measured against the broken baseline).**
Verified by: an actual codex command transcript captured in the validation stage report — `codex plugin marketplace add spacedock-dev/marketplace@edge` registers a marketplace named `spacedock-edge` distinct from `spacedock`, and `codex plugin add spacedock@spacedock-edge` exits 0 with `codex plugin list` showing `spacedock@spacedock-edge` installed. Baseline that moves the wrong way: today this errors `plugin spacedock was not found in marketplace spacedock-edge`.

**AC-2 — The binary resolves a channel-specific marketplace source.**
Verified by: a Go test asserting the edge devBranch resolves the install source `spacedock-dev/marketplace@edge`, stable resolves `spacedock-dev/marketplace`, and `SPACEDOCK_MARKETPLACE_SOURCE` overrides both. Code gate, not prose.

**AC-3 — `spacedock install --host codex` drives a real, succeeding edge install through the install seam.**
Verified by: a real `spacedock install --host codex` run against an edge binary, validated against the published edge marketplace (not a synthesized fixture), exit 0 with the edge plugin installed/refreshed.

## Test plan

- Spike first (step 2): run the manual codex commands against the new edge branch; record the transcript in the task body. Seeds the implementation's first test.
- Go unit test: channel-specific source resolution for edge vs stable + `SPACEDOCK_MARKETPLACE_SOURCE` override (AC-2).
- Real validation: an actual `codex plugin add spacedock@spacedock-edge` and `spacedock install --host codex` against the published edge (AC-1, AC-3).
- High-stakes surface (install/release machinery): detached adversarial audit at validation per the README proof policy.

## Stage Report: implementation

- DONE: edge branch in `spacedock-dev/marketplace` with root marketplace.json named `spacedock-edge` (single entry `spacedock`, source `spacedock.git` ref `next`); codex mechanism VERIFIED end-to-end FIRST.
  Created `edge` branch (commit 07f23aae) off main 2733eedd; main untouched. Spike: `codex plugin marketplace add spacedock-dev/marketplace@edge` registered a DISTINCT marketplace `spacedock-edge`, `codex plugin add spacedock@spacedock-edge` exited 0 (transcript in summary below).
- DONE: binary edge marketplace source is channel-specific (edge -> `…/marketplace@edge`, stable -> `…/marketplace`), `SPACEDOCK_MARKETPLACE_SOURCE` override preserved for both, Go unit test written test-first (TDD).
  `channelMarketplaceSource(devBranch)` in host_exec.go; `internal/cli/channel_marketplace_source_test.go` (3 tests: source-by-channel, override-verbatim, edge install seam) confirmed failing before impl, green after. Wired into runInit + front-door auto-install for both hosts. Worktree commit 6690a4af.
- DONE: a real `spacedock install --host codex` against the published edge marketplace succeeds end-to-end, validated against the REAL marketplace (not a fixture).
  Edge binary (`go build`, devBranch=next) `install --host codex` exit 0 against the live `edge` branch — both as refresh and in an isolated CODEX_HOME fresh box; plugin `spacedock@spacedock-edge` 0.23.0-pre installed. Baseline reproduced raw: bare source errors `plugin spacedock was not found in marketplace spacedock-edge` (exit 1).

### Summary

Root cause was one `marketplaceSource` var feeding the bare `spacedock-dev/marketplace` into BOTH host arms; the repo-root marketplace.json is named `spacedock`, so an edge `marketplace add` registered `spacedock` and the `spacedock@spacedock-edge` lookup failed. Fix: `channelMarketplaceSource(devBranch)` resolves the edge source to `spacedock-dev/marketplace@edge` (a new `edge` branch whose root marketplace.json is named `spacedock-edge`), stable stays bare; the `SPACEDOCK_MARKETPLACE_SOURCE` override wins verbatim on every channel (a chosen local marketplace gets no `@edge` appended). Host-neutral, so claude-edge is fixed too. The argv-sequence builders stay source-agnostic; the channel suffix is applied upstream. Real-marketplace transcript: `marketplace add …@edge` -> "Added marketplace `spacedock-edge`"; `plugin add spacedock@spacedock-edge` -> exit 0, installed 0.23.0-pre, ref `next`. Also corrected install.md (edge block now adds the `@edge` source — it previously documented the bug), releasing.md (branch-per-channel model), and the installArgvSequence comments. Notes for validation/cleanup: (1) created the permanent `edge` branch in `spacedock-dev/marketplace` (captain pre-approved; required for the channel) — do NOT delete it. (2) This box's `~/.codex` was RESTORED to its pre-spike baseline (local `spacedock` marketplace + `spacedock@spacedock` 0.22.0 only; no `spacedock-edge` artifacts; temp binary and all isolated CODEX_HOMEs removed). (3) `go test ./...` all green; the lone `go vet` finding (`pi_frontdoor_test.go:701 append with no values`) is pre-existing and out of scope.

## Stage Report: validation

- DONE: Reproduce AC-1 against the REAL published edge marketplace (success + baseline failure).
  Isolated CODEX_HOME, real GitHub source: `codex plugin marketplace add spacedock-dev/marketplace@edge` -> "Added marketplace `spacedock-edge` from …marketplace.git#edge" (distinct name), `codex plugin add spacedock@spacedock-edge` exit 0, `plugin list` shows installed/enabled 0.23.0-pre ref `next`. Baseline reproduced raw: bare `spacedock-dev/marketplace` registers marketplace named `spacedock`, then `plugin add spacedock@spacedock-edge` -> `Error: plugin spacedock was not found in marketplace spacedock-edge` (exit 1).
- DONE: AC-2 Go test green (edge->@edge, stable->bare, SPACEDOCK_MARKETPLACE_SOURCE override verbatim).
  `go test ./internal/cli/ -run 'ChannelMarketplaceSource|EdgeInstallSeam'` 3/3 PASS; the bare `-run ChannelMarketplaceSource` regex misses the install-seam test, so I ran the full file. `go test ./...` all green.
- DONE: AC-3 real `spacedock install --host codex` (edge binary, devBranch=next) exit 0 against the live edge branch.
  Plain `go build` yields an edge binary (default devBranch=next, no ldflags needed). Fresh isolated CODEX_HOME with no SPACEDOCK_* env: install exit 0, adds `…/marketplace@edge`, registers `spacedock-edge`, installs `spacedock@spacedock-edge` 0.23.0-pre, doctor "OK … compatible". Refresh run (second install on installed box) also exit 0.
- DONE: Detached adversarial audit on a THROWAWAY checkout (never the impl worktree).
  Fresh `git clone` of the branch to /tmp; 4 claim-breaking edits each go RED: edge->bare (edge + install-seam tests fail), override appends @edge (override/next fails), stable appends @edge (stable fails), gate inverted `==` (edge + override + seam fail). No test holes; the env-override-no-@edge guard the FO flagged is genuinely exercised. CLEAN audit.
- DONE: Confirm shipped docs/comments match real behavior; go vet finding pre-existing.
  install.md edge block now adds `marketplace add …@edge` (was the documented bug); releasing.md is branch-per-channel; install*ArgvSequence comments accurate (source is channel-resolved). Host-neutrality spot-check: documented `claude plugin marketplace add spacedock-dev/marketplace@edge` + install succeeds in isolated CLAUDE_CONFIG_DIR (registers `spacedock-edge`). `go vet` finding `pi_frontdoor_test.go:701 append with no values` confirmed pre-existing — that file is not in commit 6690a4af's changeset (last touched by #421).

### Summary

PASSED. All three ACs reproduced end-to-end against the REAL published `spacedock-dev/marketplace@edge` (branch at 07f23aae, root marketplace.json named `spacedock-edge`, single `spacedock` entry -> `spacedock.git` ref `next`, plugin manifest name `spacedock` 0.23.0-pre), not a fixture. AC-1 success and the broken baseline both reproduce; AC-2's 3 Go tests are green; AC-3's edge-binary `install --host codex` exits 0 fresh and on refresh. The detached adversarial audit on a throwaway clone found NO test holes — all four claim-breaking edits (including the env-override-appends-@edge case the FO singled out) go RED. Docs and comments match observed behavior and the fix is host-neutral (claude edge install verified too). All verification ran in isolated CODEX_HOME / CLAUDE_CONFIG_DIR dirs under /tmp, all removed afterward; the real `~/.codex` is untouched at baseline (no `spacedock-edge`, `spacedock@spacedock` 0.22.0 only), the permanent `edge` branch (07f23aae) and marketplace `main` (2733eedd) are intact, and the impl worktree is clean at 6690a4af.
