---
id: s0cqcf9hg4k0tgartr6ymf1b
title: spacedock install --host claude silently reuses stale marketplace ref (defeats fresh-install-from-next)
status: ideation
source: captain (2026-06-02) — observed `./spacedock install --host claude` reports success but the marketplace add no-ops on "already on disk", so the @next ref never replaces the existing pin; the uninstall+install cycle then re-pulls from the stale ref
score: "0.38"
worktree:
started: 2026-06-02T21:14:43Z
completed:
verdict:
issue:
---

`./spacedock install --host claude` runs a 3-command sequence (`internal/cli/host_exec.go:235-241` `installArgvSequence`):

1. `claude plugin marketplace add spacedock-dev/spacedock@next`
2. `claude plugin uninstall spacedock@spacedock`
3. `claude plugin install spacedock@spacedock`

Step 1 is a no-op when the marketplace `spacedock` is already declared in user settings — claude emits `Marketplace 'spacedock' already on disk — declared in user settings` and skips. The `@next` ref pin never lands. Steps 2-3 then uninstall and re-install the plugin sourced from whatever ref the marketplace was originally added with (the default branch, or a stale `@next` from a prior session). The user sees a successful-looking install and a stale plugin.

This silently defeats the recalibrated sprint goal "install path off `next` — fresh install works from `spacedock-dev/spacedock@next`."

## Problem

The 3-step sequence assumes `marketplace add` is idempotent and re-pins the ref. It is not: when the marketplace is already on disk, claude skips the add entirely with no exit-code signal, so the downstream uninstall/install pair operates on the stale ref. There is no observable failure for the user or for any test that just checks the install exit code.

## Proposed approach

Insert a `marketplace remove spacedock` before the add, making it a 4-step sequence:

1. `claude plugin marketplace remove spacedock` (no-op on a fresh box)
2. `claude plugin marketplace add spacedock-dev/spacedock@next`
3. `claude plugin uninstall spacedock@spacedock`
4. `claude plugin install spacedock@spacedock`

The remove step is benign on first install (claude either succeeds silently or emits "not declared," neither fatal). On an upgrade, it clears the stale on-disk pin so the add lands the @next ref cleanly.

Open design question for ideation: should the remove step's stderr be tolerated unconditionally (current `Install` aggregates `CombinedOutput` and fails on any non-zero exit), or should we wrap it to ignore the "not declared" exit? A small `installArgvSequence` refactor with a per-step "tolerate failure" flag covers it cleanly.

## Out of scope

- Codex install (the codex path is documented prose; the @ref is already correct on every printed invocation).
- The `--ref` flag surface itself (already correct — `marketplaceAddArg` composes `source@branch` correctly when `devBranch != ""`).

## Acceptance criteria

**AC-1 — `spacedock install --host claude` lands the `@next`-pinned plugin even when a stale marketplace declaration already exists.**
Verified by: a host-ops test that primes a fakeHost with an existing `spacedock` marketplace pinned at the default branch, runs `runInit`, and asserts (a) the issued argv sequence carries the remove step before the add, and (b) the resulting marketplace pin reflects `spacedock-dev/spacedock@next`. Plus a manual confirmation by the captain that the live install on a clean session pulls the `next` ref.

**AC-2 — A first install on a fresh box still succeeds end-to-end (the remove step's "not declared" outcome is tolerated).**
Verified by: a host-ops test with no pre-declared marketplace where `runInit` issues all four steps and returns success despite the remove step's non-zero exit (or a "not declared"-flavored stderr). The contract gate then reports the @next-installed contract.

## Test plan

- Host-ops unit tests via the existing `hostOps`/`fakeHost` seam in `internal/cli/`: extend the fake to record argv sequences and simulate the "already-on-disk" no-op, plus the "not declared" fresh-box case. Assertions are over the issued argv shape and the resulting per-host pin.
- Cost: low. Reuses the existing test seam; no live host needed.
- No live-workflow test required — the claim is the argv sequence emitted by the binary, not a runtime integration.

## Notes

- Captain-observed 2026-06-02 — `./spacedock install --host claude` after a prior session's install: `Marketplace 'spacedock' already on disk — declared in user settings` followed by a "successful" uninstall+install pair from the stale ref.
- Coordinates with packaging work (post-bootstrap repo migration to `spacedock-dev/spacedock`): a clean remove+add sequence is the right shape for `requires-contract` cross-repo verification too.
- Manual workaround until this lands: `claude plugin marketplace remove spacedock && ./spacedock install --host claude`.
