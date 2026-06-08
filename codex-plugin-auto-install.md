---
id: z9pqhvtj21smtxka8p88j23r
title: Codex plugin auto-install on `spacedock codex` (mirror #311; codex 0.137.0 supports marketplace + plugin add)
status: implementation
source: "captain (2026-06-08) — codex-cli 0.137.0 adds `codex plugin marketplace add` + `codex plugin add` (0.132.0 lacked them). The front-door auto-install is currently Claude-only (frontdoor.go:315 'codex has nothing to auto-install'); now the Codex binary-first path can ensure the plugin exists too, like #311 did for Claude."
score: "0.32"
started: 2026-06-08T15:48:51Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-codex-plugin-auto-install
issue:
sprint: 0198-pre-flip-hardening
group: binary-ux
sprint-readiness: ready
---

Bring the front-door plugin auto-install to the Codex binary-first path, now that codex-cli 0.137.0 supports adding a marketplace + plugin non-interactively. Mirror the #311 Claude shape: on the front door, a missing plugin (`NoPluginFound`) auto-installs and proceeds to launch so the single command the user typed yields a working first-officer session; `--no-install` opts out; a contract mismatch still fails fast.

## Problem

- `spacedock claude` already auto-installs a missing plugin (`NoPluginFound → ops.Install("claude", …) → launch`, `frontdoor.go:177-188` — the #311 pattern).
- `spacedock codex` does NOT: `frontdoor.go:317-321` keeps the all-or-nothing gate and `host_exec.go:271-274` `execHost.Install` errors for any non-claude host ("programmatic install is only supported for claude"), because earlier codex (0.132.0) could not add a marketplace/plugin from the CLI.
- codex-cli **0.137.0** now has `codex plugin marketplace add` / `plugin add` / `list` — so the Codex path can ensure the plugin exists on first launch, the same way Claude does.

## Proposed approach

Three changes, all behind the existing `hostOps` seam (no new interface):

1. **`execHost.Install` (host_exec.go:271): add a codex branch.** Today it `return …fmt.Errorf("programmatic install is only supported for claude…")` for non-claude. Add a codex arm that runs a codex-specific argv sequence. Keep the claude arm byte-identical.
2. **A codex install sequence helper** alongside `installArgvSequence` (host_exec.go:254). The spike (below) settled the exact shape:
   - `plugin remove spacedock@spacedock` — tolerated (cleanup; exit 0 even fresh-box, see spike)
   - `plugin marketplace remove spacedock` — **tolerated** (cleanup; exits **1** "is not configured or installed" on a fresh box, see spike)
   - `plugin marketplace add spacedock-dev/spacedock --ref <branch>` — fail-fast (the real-failure backstop)
   - `plugin add spacedock@spacedock` — fail-fast
   The `--ref <branch>` form is codex's own pinning flag (its verb is a separate `--ref`, NOT the claude `source@branch` shorthand that `marketplaceAddArg` builds). The branch is `devBranch` (default `next`). When `branch == ""`, omit `--ref` and its value.
3. **`runCodex` (frontdoor.go:317): wire `NoPluginFound → auto-install → launch`.** Replace the all-or-nothing `if gateHost(…) != contract.Compatible { return 1 }` with the same `switch` `runClaude` uses (`Compatible` → proceed; `NoPluginFound` → `--no-install`? fail : `ops.Install("codex", marketplaceSource, devBranch)` then proceed; default → fail fast). Reuse the `--no-install` flag (already parsed in `frontDoorArgs`) and the existing `codexEntryInstalled`/`codexCacheManifest` resolver for the post-install gate verdict — no resolver change needed (spike confirmed the cache manifest lands where the resolver looks).

The `--no-install`, `--skip-contract-check`, and "any non-Compatible-non-NoPluginFound fails fast" semantics are identical to the Claude branch; this is the codex analog of #311.

### Out of scope (recorded, not done)

- **`spacedock install --host codex` prose path (`runInit`, init.go:41-58).** That surface still prints `printCodexInstallProse` rather than calling `ops.Install("codex", …)`. Flipping it to a programmatic install is a natural follow-up now that the codex Install arm exists, but the task and its checklist target the **front door**. Flagged so a follow-up task can pick it up; not in this task's ACs.
- **A codex `--json`-based resolver.** Spike found `codex plugin list --json` now works in 0.137.0 (the host_exec.go:32-34 comment "Codex 0.136.0 has no --json (rejects it, exit 2)" is stale for 0.137.0), with a different schema (`{"installed":[{"pluginId":…,"installed":true}]}`). The existing text-listing + cache-layout resolver still resolves correctly, so YAGNI — no resolver rewrite here.

## Riskiest unknown (spike — DONE)

**Q:** Do the 0.137.0 commands run NON-INTERACTIVELY (no prompt) for an unattended auto-install, and what is the exact marketplace+add sequence + source/ref flags?

**Run (codex-cli 0.137.0 on this box, isolated `CODEX_HOME`, stdin closed `</dev/null`):**

- `codex plugin marketplace add <local-path>` → exit 0, no prompt ("Added marketplace `spacedock`").
- `codex plugin add spacedock@spacedock` → exit 0, no prompt ("Added plugin `spacedock`"). Installs to `<CODEX_HOME>/plugins/cache/spacedock/spacedock/<version>/.codex-plugin/plugin.json` — exactly where `codexCacheManifest()` resolves.
- Re-running both when already installed → exit 0 (idempotent).
- `codex plugin remove spacedock@spacedock` on a fresh box → exit 0 (idempotent cleanup).
- `codex plugin marketplace remove spacedock` on a fresh box → **exit 1** "marketplace `spacedock` is not configured or installed" → this cleanup step MUST be tolerated.
- **Full end-to-end against the real `spacedock-dev/spacedock@next`:** `codex plugin marketplace add spacedock-dev/spacedock --ref next` (exit 0) → `codex plugin add spacedock@spacedock` (exit 0, installs 0.19.7) → `codex plugin list --json` reports `spacedock@spacedock installed:true`. Both `--ref next` and the `owner/repo@next` shorthand resolve to `https://github.com/spacedock-dev/spacedock.git#next`.

**Verdict:** No gate finding. The mechanism is sound — every step runs unattended (exit 0, no prompt). The design only needs the tolerance asymmetry (both cleanup steps tolerated; both pinning steps fail-fast), mirroring the claude sequence.

## Acceptance criteria

Each AC names a check whose expected value comes from outside the changed file (the fake host's recorded calls, the launch-seam argv, the exit code, or a real codex install's on-disk state) — never a substring of the source under test.

- **AC-1 — `execHost.Install` issues the codex sequence in the right order with the right tolerance.** Driven by a per-PATH `codex` stub (mirroring `writeClaudeStub`/`install_tolerance_test.go`). Verified by: (a) a stub that exits 0 on every step → `Install("codex", "spacedock-dev/spacedock", "next")` returns nil error and combined output contains all four step markers in order, including `plugin marketplace add spacedock-dev/spacedock --ref next` and `plugin add spacedock@spacedock`; (b) a stub that exits 1 on `plugin marketplace remove` → still nil error (tolerated); (c) a stub that exits 1 on `plugin marketplace add` → non-nil error wrapping that subcommand (fail-fast). The stub is the independent source of truth (its echoed argv), not the function body.
- **AC-2 — `spacedock codex` with no plugin auto-installs then launches.** Verified by a `fakeHost{manifest: ""}` front-door test: `runCodex` records `installCmds` beginning `codex spacedock-dev/spacedock next` (the seam's {host, source, branch}) AND `launchedArg` is the assembled codex launch argv. This INVERTS the current `TestCodexFrontDoorNoPluginFailsFastWithoutInstalling` (frontdoor_test.go:414), which asserts no install + no launch — that test migrates to the `--no-install` arm. The phantom-installPath case (resolved-but-missing manifest, the `missing` sub-case) auto-installs too, matching #311.
- **AC-3 — `--no-install` opts out (refuse-and-instruct).** Verified by a `fakeHost{manifest: ""}` + `--no-install` front-door test: `runCodex` exits non-zero, `installCmds` is empty, `launchedArg` is nil, and stderr carries the `codex plugin` no-plugin remedy (the migrated fail-fast assertion).
- **AC-4 — a contract mismatch still fails fast without installing.** Verified by the existing `TestCodexFrontDoorFailFastOnMismatch` (frontdoor_test.go:394) staying green: `tooOldBinaryManifest` → exit non-zero, `installCmds` empty, `launchedArg` nil. (No code change to this path; the AC guards that the new switch's `default` arm preserves it.)
- **AC-5 — the codex install is host-native, non-interactive, end-to-end.** A real isolated-`CODEX_HOME` behavioral test (mirroring `TestClaudePluginInstallIsHostNative`, skipped when `codex` not on PATH): build a local-path marketplace, run `execHost{}.Install("codex", <local-marketplace-path>, "")`, then assert `codex plugin list --json` reports `spacedock@spacedock installed:true` and the cache manifest exists at the resolver's path. Source is the local path (no `--ref`, branch `""`) so the test stays hermetic/offline. This proves the real codex CLI accepts the sequence unattended, not just the stub.

## Test plan

- **Unit / stub (AC-1, AC-2, AC-3, AC-4):** Go tests in `internal/cli`. Fast (<1s), no network. AC-1 adds a `codex` analog of `writeClaudeStub` (or generalizes it to take a binary name) in a new `install_tolerance` codex case or alongside it. AC-2/AC-3 extend `frontdoor_test.go`'s codex cases against `fakeHost`. AC-4 is an existing test kept green. Cost: low; this is the bulk of coverage and the build-time gate.
- **Behavioral / real CLI (AC-5):** one Go test, `t.Skip` when `codex` not on PATH, isolated `CODEX_HOME` + local-path marketplace, real `execHost{}.Install`. Mirrors `install_behavior_test.go`'s claude analog. Cost: low when codex present (~1-2s, offline local-path); skipped otherwise.
- **No live-workflow test needed:** the claim is install-sequence + front-door wiring (command/exec behavior), not runtime FO behavior. The launch itself is exercised by the `fakeHost.launchedArg` assertion; `syscall.Exec` is unchanged.
- **Fixtures reused:** `compatibleManifest`, `tooOldBinaryManifest`, `safehouseFixtureDir`, `wantCodexBootstrapPrompt`, `fakeHost`.

## Notes

High-stakes: this is the **front-door launcher** — validation gets a detached adversarial audit. Codex analog of `rbjk` (#311). 0198 binary-ux; couples with `qa` (the install/launch journey). The spike was run live against codex-cli 0.137.0 and the real `spacedock-dev/spacedock@next`; the implementation's first test should be AC-1 (the install sequence), seeded by the spike's exact argv + tolerance findings.

## Stage Report: ideation

- DONE: SPIKE first (riskiest): confirm codex-cli 0.137.0's `plugin marketplace add` + `plugin add` run NON-INTERACTIVELY for an unattended auto-install, and nail the exact sequence/flags (source/ref).
  Ran live against codex-cli 0.137.0 (isolated CODEX_HOME, stdin closed): marketplace add + plugin add both exit 0 with no prompt, idempotent on re-run; full end-to-end against real `spacedock-dev/spacedock@next` installed 0.19.7 and `plugin list --json` reported `installed:true`. Sequence/flags recorded in "Riskiest unknown (spike — DONE)". No gate finding.
- DONE: Design the auto-install mirroring #311: extend execHost.Install + installArgvSequence to a codex sequence, and wire the codex NoPluginFound → auto-install → launch path, respecting --no-install and reusing the codexEntryInstalled idempotency check.
  "Proposed approach" specifies the three changes (Install codex arm at host_exec.go:271, a codex sequence helper, the runCodex switch at frontdoor.go:317) with exact argv, the `--ref <branch>` pinning (codex's own flag, not the claude `source@branch` shorthand), and the tolerance asymmetry the spike settled (both cleanup steps tolerated — marketplace remove exits 1 fresh-box; both pinning steps fail-fast).
- DONE: Produce build-ready ACs + test plan: a host-fixture/behavior test that codex NoPluginFound triggers the install sequence then launches (mismatch still fails fast; --no-install opts out). Note front-door = high-stakes → validation gets a detached adversarial audit.
  Five entity-level ACs (AC-1..AC-5), each bound to an independent source of truth (stub-echoed argv, fakeHost recorded calls, launch-seam argv, exit code, real-CLI on-disk state) — none a substring of the changed file. Test plan covers stub/unit + one offline behavioral test; high-stakes note carried in "Notes".

### Summary

Fleshed out the codex front-door plugin auto-install as the codex analog of #311. The riskiest unknown — whether codex-cli 0.137.0's plugin commands run unattended — was settled by a live spike (all steps exit 0, no prompt, idempotent, end-to-end against the real next branch), so the design composes proven mechanisms with no open mechanism risk. Two items were explicitly scoped out and recorded: flipping the `spacedock install --host codex` prose path to a programmatic install (a natural follow-up), and a `--json`-based codex resolver (the existing text+cache resolver still works; the host_exec.go:32-34 comment about 0.136.0 lacking `--json` is now stale for 0.137.0 but the resolver itself is correct).
