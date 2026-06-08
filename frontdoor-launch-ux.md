---
id: yqf0amtyecjcft0vsw6nbqtk
title: spacedock codex/claude front-door launch UX — auto-install messaging, pre-launch info banner, neutral bootstrap prompt
status: implementation
source: "captain live-test of 0.19.8 (2026-06-08). The first real `spacedock codex` run surfaced three front-door UX issues (A/B/D below). The pre-cut audit verified z9's auto-install SEAM via tests but nobody ran the front door end-to-end — so a real `spacedock codex` live drive is this task's load-bearing proof."
started: 2026-06-08T21:00:22Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-frontdoor-launch-ux
issue:
group: binary-ux
sprint: 0199-pre-flip-mechanics
sprint-readiness: ready
---

Make the `spacedock codex` (and `claude`) launch experience honest and useful: don't tell the user to install manually right before silently auto-installing, show a short pre-launch info banner, and ship a neutral bootstrap prompt.

## Problem

- **A — the auto-install message is self-contradicting.** On `NoPluginFound`, `gateHost` (`frontdoor.go:124-128`) prints the *manual-install remedy* — "no installed codex plugin found. Run `spacedock install --host codex` (or `spacedock claude --skip-contract-check` to bootstrap)" — and then `runCodex`/`runClaude` **auto-install anyway** (`ops.Install`, silently — `execHost.Install` writes no progress). So the tool says "install it yourself," then quietly does it. Two defects: the manual remedy fires on the auto-install path, and it hardcodes a `spacedock claude` hint that is wrong in a codex run.
- **B — no useful pre-launch info.** `spacedock codex`/`claude` launches straight into the host with no Spacedock context (version, which workflow it detected, anything actionable).
- **D — the shipped bootstrap prompts carry personal flavor text.** Both `bootstrapPrompt` (`frontdoor.go:24`) and `codexBootstrapPrompt` (`:289`) literally read "You totally got this. Take your time. I love you. And tell all subagents and team members you love them too. Engage." — so the FO tries to *relay* that at an empty-team session start ("no team to relay your note to"). The product default should be neutral.

## Proposed approach

Three independent edits in `internal/cli/frontdoor.go`, each with an exact before/after below. All three compose already-shipped front-door behavior (the contract gate, the auto-install seam, the launch seam, `stderr` writes, and `status.DiscoverWorkflowDir`); none introduces a new mechanism. See **No spike needed** below.

### A — honest auto-install messaging (move the NoPluginFound remedy to the caller)

The defect's root is *who owns the NoPluginFound message*. Today `gateHost` prints a remedy for **every** non-Compatible verdict, including `NoPluginFound` — but the caller's response to `NoPluginFound` is not uniform: it auto-installs by default and only fails fast under `--no-install`. So `gateHost` prints "install it yourself" before the caller silently auto-installs, and it hardcodes a `spacedock claude` hint that `gateHost` knows is wrong (it has `host`, but the literal is fixed).

Fix: `gateHost` stops emitting the `NoPluginFound` remedy. It keeps emitting for the verdicts where the caller ALWAYS fails fast (resolve-error → `MalformedRange`, mismatch → not-Compatible: line 139's `res.Message`). The caller now owns the `NoPluginFound` message, which lets it pick the right one:

- on auto-install (the default): print `Installing the {host} plugin…` to stderr immediately before `ops.Install` — `{host}` is `"claude"`/`"codex"`, so a codex run never says "claude";
- on `--no-install`: print the manual remedy with the **host-correct** bootstrap hint — `spacedock install --host {host}` and `spacedock {host} --skip-contract-check` (no hardcoded `claude`).

Exact changes in `gateHost` (frontdoor.go:116-142):
- **DELETE** the `manifestPath == ""` `Fprintf` (lines 125-127) — return `NoPluginFound` with no print.
- **DELETE** the `res.Verdict == NoPluginFound` `Fprintf` (lines 132-135) — return `NoPluginFound` with no print. (The phantom-manifest path collapses to the same caller-owned message; the operator does not need to know "manifest path missing" vs "no entry" — both mean "no plugin, installing".)
- **KEEP** the resolve-error `Fprintf` (lines 119-121) and the `res.Verdict != Compatible` `Fprintln(res.Message)` (line 139) — these are the always-fail-fast verdicts.
- Doc-comment update: `gateHost` now prints the remedy for non-Compatible verdicts **except `NoPluginFound`, whose message the caller owns** (auto-install vs refuse).

Exact changes in `runClaude` / `runCodex` (the `case contract.NoPluginFound:` arm, lines 177-188 / 325-335):
- BEFORE the `ops.Install(...)` call, on the auto-install branch (i.e. NOT `--no-install`): `fmt.Fprintf(stderr, "Installing the %s plugin…\n", host)`.
- In the `if fd.noInstall { … }` branch: print the host-correct manual remedy before `return 1` — `fmt.Fprintf(stderr, "Spacedock: no installed %s plugin found. Run `+"`spacedock install --host %s`"+` (or `+"`spacedock %s --skip-contract-check`"+` to bootstrap).\n", host, host, host)`. Factor this into a small `noPluginRemedy(host string) string` helper so both launchers share one wording and the message-oracle has one source.

Before/after of the user-visible behavior:
- **Before:** `spacedock codex` (no plugin) → stderr `Spacedock: no installed codex plugin found. Run `spacedock install --host codex` (or `spacedock claude --skip-contract-check` to bootstrap).` then silent install then launch.
- **After:** `spacedock codex` (no plugin) → stderr `Installing the codex plugin…` then install then launch. `spacedock codex --no-install` → stderr `Spacedock: no installed codex plugin found. Run `spacedock install --host codex` (or `spacedock codex --skip-contract-check` to bootstrap).` then exit 1 (note: `codex`, not `claude`).

### B — pre-launch info banner

A short banner to stderr before the host launches, emitted in both `runClaude` and `runCodex` immediately before assembling `inner`/`argv` (after the gate, after `warnStrayPromptAfterDash`, before `ops.Launch` — `Launch` execs and never returns, so anything after it is unreachable). Banner is suppressed on `--resume`/`codexResume` (the operator is continuing a session, not starting one) to avoid noise — though if simpler to always print, that is acceptable; the AC pins the no-resume case.

Content (3 lines, stderr; keep it to the version + detected workflow + one orientation line):

    Spacedock v{Version} · first officer launching {host}
    Workflow: {detected}        ← "docs/dev" (the rel path from cwd) when status.DiscoverWorkflowDir(dir) hits, else "none detected (launching anyway)"
    {host} is starting as your first officer; run `spacedock status` inside the session for the queue.

`{detected}` uses `status.DiscoverWorkflowDir(dir)` (already public, already the walk-up that recognizes a commissioned workflow by its README's `commissioned-by: spacedock@`). Render the workflow as its path relative to `dir` (falling back to the absolute on a `filepath.Rel` error). Factor the banner into a `launchBanner(host, dir string, w io.Writer)` helper so both launchers share it and the golden/string oracle has one source. `Version` is `cli.Version` (already in package `cli`).

### D — neutral bootstrap prompts

Replace the personal/relay flavor in both shipped constants; keep the functional FO-selection.

- `bootstrapPrompt` (frontdoor.go:24): **before** `"You totally got this. Take your time. I love you. And tell all subagents and team members you love them too. Engage."` → **after** `"Engage as the Spacedock first officer for this session."` (claude selects the FO via the `--agent spacedock:first-officer` flag already in the argv, so the prompt only needs a neutral launch-and-go.)
- `codexBootstrapPrompt` (frontdoor.go:289): **before** the same text + `" Assume $spacedock:first-officer for the entire session."` → **after** `"Engage as the Spacedock first officer for this session. Assume $spacedock:first-officer for the entire session."` (codex has no `--agent` flag, so the `Assume $spacedock:first-officer` clause is load-bearing and MUST stay — confirmed by the existing const comment at :283-288.)

The exact replacement wording is the implementer's to finalize within these constraints: NO personal/relay text ("I love you", "tell all subagents/team members"), the codex const KEEPS `Assume $spacedock:first-officer for the entire session.` verbatim. The two test-file oracle copies (`wantBootstrapPrompt` / `wantCodexBootstrapPrompt`, `safehouse_frontdoor_test.go:18,156`) are updated to the same new literals — they are an INDEPENDENT source of truth (a second hand-written copy in the test file), which is what makes the argv oracle a real check: production drift fails the argv-shape test, and a residual "I love you" fails the absence assertion.

## Out of scope

- The codex plugin install MECHANISM (z9 — works; this is messaging/UX around it).
- `th` (safehouse-preserves-spacedock-bin) territory: the wrap-path argv prefix and `launchEnv`. See **Overlap with `th`**.

## No spike needed

No spike needed: this task composes only already-proven front-door behavior — the contract gate + verdict switch (`gateHost`/`runClaude`/`runCodex`, proven by the AC-1a/1b auto-install tests), the auto-install seam (`ops.Install`, proven), stderr writes before `ops.Launch` (proven — the gate already writes there), `status.DiscoverWorkflowDir` (proven walk-up, used by native status), and a string-constant swap. There is no parser round-trip, on-disk format, runtime handoff, or tool-flag assumption whose failure would invalidate the plan. The riskiest thing is not a mechanism but a *coverage gap*: A shipped without anyone running the front door end-to-end. That gap is closed by this task's load-bearing proof — the **live drive** in the test plan — not by a throwaway spike. Host CLIs are present on the dev box (`which codex` → `/opt/homebrew/bin/codex`, `claude` resolves), so the live drive is runnable at implementation/validation.

## Acceptance criteria

All behavioral halves are proven by a real front-door **live drive** plus an independent-source argv/message oracle. A grep/substring over `frontdoor.go` does NOT satisfy any AC (the matched string was written by the implementer; the absence of a real front-door drive is exactly what let A ship).

- **AC-A — honest, host-correct no-plugin messaging.** With no installed plugin, `spacedock codex` prints `Installing the codex plugin…` (not a manual remedy) and contains NO `spacedock claude` hint; `spacedock codex --no-install` prints the manual remedy naming `spacedock install --host codex` and `spacedock codex --skip-contract-check` (host-correct — never `claude`) and exits non-zero without launching. The claude analog holds for `spacedock claude`.
  *Verified by:* (1) `runCodex`/`runClaude` message oracles in `frontdoor_test.go` driving `fakeHost{manifest: ""}` — auto-install arm asserts stderr contains `Installing the codex plugin…` AND `installCmds`/`launchedArg` are non-empty AND stderr does NOT contain `spacedock claude`; `--no-install` arm asserts stderr contains `spacedock install --host codex` and `spacedock codex --skip-contract-check`, does NOT contain `spacedock claude`, `installCmds` empty, `launchedArg` nil, exit non-zero. (2) `TestGateRemedyNamesLiveInstallCommand` updated: the `MalformedRange` (resolve-error) case still asserts `gateHost` prints `spacedock install`; the no-plugin / phantom-manifest cases move to assert the remedy on the launcher's `--no-install` output (gateHost no longer prints for NoPluginFound). (3) a real `spacedock codex` live drive (no plugin) observing `Installing the codex plugin…` on stderr and no `spacedock claude` text. The expected values live in the test file / live terminal, not in `frontdoor.go`.

- **AC-B — pre-launch banner renders version + detected workflow.** Before the host launches, stderr carries a banner whose first line names `spacedock {Version}` and whose workflow line reads the path of the discovered workflow (e.g. `docs/dev`) when launched inside a commissioned workflow, or `none detected` when launched outside one.
  *Verified by:* (1) a `launchBanner` unit/oracle test driving it with (a) a temp dir whose ancestor README declares `commissioned-by: spacedock@…` → asserts the rendered workflow line names that dir's relative path, and (b) a bare temp dir → asserts `none detected`; both assert the version line contains `cli.Version`. The discovered-workflow expectation comes from the fixture README the test writes, an independent source. (2) a `runClaude`/`runCodex` oracle asserting the banner reaches stderr before `Launch` (stderr non-empty with the version line; `launchedArg` set). (3) a real live drive run from inside `docs/dev` observing the `docs/dev` workflow line and from `/tmp` observing `none detected`.

- **AC-D — neutral launched prompt, FO selection intact.** The launched inner argv's bootstrap-prompt token contains NO personal/relay text (`I love you`, `tell all subagents`, `tell all … team members`); the codex prompt token still contains `Assume $spacedock:first-officer for the entire session.`
  *Verified by:* (1) the existing argv oracles (`frontdoor_test.go` / `safehouse_frontdoor_test.go` / `frontdoor_stray_prompt_test.go`) with `wantBootstrapPrompt` / `wantCodexBootstrapPrompt` updated to the new neutral literals — the launched argv must equal the new value (independent test-file copy), AND a new assertion that the launched prompt token does NOT contain `I love you` / `tell all subagents`, AND that the codex token DOES contain `Assume $spacedock:first-officer`. (2) a real live drive: confirm the launched session starts the FO and the prompt shown is neutral (no relay attempt — the original report was the FO trying to relay "tell all team members you love them"; the drive confirms that relay no longer happens).

## Test plan

- **Message/argv oracles (`internal/cli`, Go unit tests) — low cost, no network, no host CLI.** AC-A: extend `frontdoor_test.go`'s `fakeHost{manifest: ""}` cases with the stderr-message assertions above (auto-install line, host-correct remedy, no-`claude`-hint), update `TestGateRemedyNamesLiveInstallCommand` for the moved NoPluginFound responsibility. AC-B: new `launchBanner` test with two fixture dirs (commissioned README vs bare) + a launcher-level "banner reaches stderr before Launch" assertion. AC-D: update the four `want*BootstrapPrompt` oracle sites to the new literals and add the personal-text-absence / FO-clause-presence assertions. These are the same `fakeHost.launchedArg` / stderr-buffer oracle shape the package already uses; expected values come from the test file and fixtures, never from `frontdoor.go`.
- **Live drive (load-bearing — the front-door run nobody did before A shipped).** Build the dev binary; run, against a scratch dir with no installed plugin and one inside `docs/dev`:
  - `spacedock codex` (no plugin) → observe `Installing the codex plugin…`, the banner (version + workflow line), and a neutral launched prompt; observe NO `spacedock claude` text.
  - `spacedock codex --no-install` (no plugin) → observe the host-correct manual remedy and non-zero exit.
  - `spacedock claude` analog for the claude paths.
  - from `/tmp` (outside any workflow) → observe `none detected`.
  Capture terminal output as the evidence. This is the end-to-end front-door exercise; it is the proof, not a spike.
- **Detached adversarial audit at validation.** Front-door is high-stakes (the surface every user hits first, and the surface a passing test suite already failed to protect once). Flag a detached adversarial audit at the validation stage — independent of the implementer — re-running the live drive and probing the `--no-install`, resume, and outside-workflow edges.
- **Estimated complexity:** small-to-medium. Three localized edits in one file (`frontdoor.go`), three test files touched, two small helpers factored (`noPluginRemedy`, `launchBanner`). No fixtures beyond the temp-dir READMEs the banner test writes. The live drive is the only non-unit cost.

## Overlap with `th` (safehouse-preserves-spacedock-bin)

Confirmed clean — both edit `frontdoor.go` but disjoint funcs:

- **`th`** (id `thdf2…`, currently in **validation**) edits the **wrap path**: it factors `resolvedLauncherBin()` and adds `launcherBinArgvPrefix()` (a `/usr/bin/env SPACEDOCK_BIN=<bin>` token group) injected at the `safehouse.Wrap(...)` call sites in `runClaude` (≈line 216) and `runCodex` (≈line 362), plus `launchEnv` (lines 55-61). It touches argv ASSEMBLY on the wrapped branch and the env helper.
- **this task (yq)** edits: `gateHost` message emission (lines 116-142), the two NoPluginFound caller arms (lines 177-188 / 325-335), a banner call sited *before* argv assembly, and the two bootstrap CONSTS (lines 24, 289).

The only shared lines are the launcher function bodies, but the edits are to non-overlapping statements: `th` works on the `safehouse.Wrap`/`launchEnv` argv-and-env machinery; yq works on the gate-message branch, a pre-assembly banner print, and the prompt constants. No shared statement is edited by both. **No real collision.** Sequencing note: `th` is already in validation and will likely land first; if it merges before yq's implementation, the implementer rebases onto it and the banner print sits cleanly between the gate switch and `th`'s wrap-path prefix. If yq lands first, `th`'s rebase is equally clean. Either order, the conflict is a trivial textual adjacency at worst, not a semantic one.

## Stage Report: ideation

- DONE: Scope A/B/D precisely with exact before/after: A — gateHost (frontdoor.go:124-128) suppresses the manual-install remedy on the AUTO-INSTALL path and prints "Installing the {host} plugin…" before ops.Install; host-aware hint (no `spacedock claude` in a codex run); manual remedy stays ONLY on the --no-install / hard-fail branch. B — the pre-launch banner content (Spacedock vX + detected workflow `docs/dev`/none + 1-2 useful lines). D — the neutral bootstrap-prompt text replacing the personal/relay flavor (keep `Assume $spacedock:first-officer`). Record "no spike needed: {the proven front-door paths it composes}" or exercise the riskiest unknown first.
  Proposed approach now carries line-anchored before/after for all three: A moves NoPluginFound message ownership from gateHost to the caller (delete the two NoPluginFound Fprintf's at :125-127/:132-135, keep the always-fail-fast prints; caller prints `Installing the %s plugin…` pre-Install and a host-correct `noPluginRemedy(host)` on --no-install); B is a `launchBanner(host,dir,w)` over `status.DiscoverWorkflowDir` (import + export confirmed present); D swaps both consts to neutral text keeping the codex `Assume $spacedock:first-officer` clause. "No spike needed" section records the composed proven paths (gate switch, ops.Install seam, stderr-before-Launch, DiscoverWorkflowDir, const swap) — the risk was a coverage gap, closed by the live drive, not a mechanism.
- DONE: AC + test plan proven by a real `spacedock codex`/`claude` LIVE DRIVE (observe "Installing…", the banner, and NO personal/relay text in the launched prompt) PLUS argv/message oracles for the bootstrap-prompt const + the install-path message — NEVER a grep (the absence of a real front-door live drive is exactly what let A ship). Front-door = high-stakes → flag the detached adversarial audit at validation.
  AC-A/AC-B/AC-D each pair an independent-source oracle (test-file `want*` copies, fixture READMEs, stderr-buffer message assertions over `fakeHost`) with a real live-drive leg; test plan names the exact files/cases and explicitly bans grep. Detached adversarial audit flagged at validation (re-run the live drive + probe --no-install/resume/outside-workflow edges).
- DONE: Confirm the th overlap is clean: both edit frontdoor.go but disjoint funcs (th = launchEnv/argv-prefix on the wrap path; yq = gateHost message + the banner + the bootstrap const). Name any real collision.
  "Overlap with th" section: read th (safehouse-preserves-spacedock-bin, id thdf2…, in validation) — it edits the wrap-path argv assembly (`launcherBinArgvPrefix` at the `safehouse.Wrap` sites) + `launchEnv`; yq edits the gate-message branch, a pre-assembly banner print, and the two prompt consts. No shared STATEMENT is edited by both. No real collision; worst case is a trivial textual adjacency on either merge order.

### Summary

Firmed the ideation for the three front-door UX defects into three localized `internal/cli/frontdoor.go` edits, each with line-anchored before/after. The load-bearing insight for A is that the defect is a message-ownership problem: `gateHost` prints a remedy for every non-Compatible verdict, but `NoPluginFound` is the one verdict whose response is non-uniform (auto-install by default, refuse under --no-install), so message ownership moves to the caller, which simultaneously fixes the hardcoded-`claude` hint (the caller has `host`). Recorded "no spike needed" — every path composes already-proven front-door behavior; the real risk was a coverage gap (A shipped with no end-to-end front-door run), closed by the live-drive proof, not a throwaway spike. Confirmed the `th` overlap is clean (disjoint statements in shared functions) and flagged the detached adversarial audit at validation since the front door is the highest-stakes surface.
