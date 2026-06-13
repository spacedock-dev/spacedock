---
id: gja5htstcgjxydcz5h2051wc
title: Show sandbox state on startup
status: ideation
source: captain (2026-06-12, UX improvement)
started: 2026-06-13T04:08:48Z
completed:
verdict:
score:
worktree:
issue:
sprint: 0201-post-flip-release-model
group: ux-cleanup
sprint-readiness: ready
---

UX improvement: on startup, Spacedock should also show whether sandboxing is enabled — distinguishing enabled / available-but-not-enabled / unavailable — so the operator knows the execution-isolation posture before dispatching work.

Captain amendment (2026-06-12): the same sandbox info should also appear in `--version`. In addition, `--version` should show each detected runtime (claude / codex / pi) with its install and spacedock-enablement status, e.g.:

```
codex: installed
claude: installed, spacedock enabled
pi: not installed
```

## Problem

The startup surfaces report mods, ID style, orphans, PR state, and dispatchables, but say nothing about whether the session is sandboxed. The operator has no at-a-glance signal of the execution-isolation posture before dispatching work — a `.safehouse` profile may be present but the launch silently unsandboxed because the `safehouse` binary is missing from PATH (this exact split is live on the dev machine: a `.safehouse` profile sits in the repo root, yet `which safehouse` returns not-found, so any launch here runs unsandboxed). Separately, `--version` prints only `spacedock <ver> (contract <N>)` and gives no view of which runtimes (claude / codex / pi) are installed or whether the Spacedock plugin is enabled in each, so an operator debugging "why won't `spacedock codex` launch" has no single command that answers it.

## Proposed approach

Three surfaces carry the new information; each renders the three-way sandbox state **enabled / available-but-not-enabled / unavailable** from the two existing safehouse probes (`safehouse.Present(dir)` = a `.safehouse` profile exists in the launch dir; `safehouse.Available(lookPath)` = the `safehouse` binary resolves on PATH):

- **enabled** — `Present(dir)` is true (or `--safehouse` / a `--safehouse-*` knob was given) AND `Available` is true: the launch will be wrapped through safehouse.
- **available-but-not-enabled** — `Available` is true but no profile/flag selects it: the binary is installed, this launch is not sandboxed.
- **unavailable** — `Available` is false: the `safehouse` binary is not on PATH; a present profile cannot take effect.

**Surface 1 — launcher banner** (`launchBanner` in `internal/cli/frontdoor.go`, written to stderr before the host takes over). Add one line after the workflow line: `Sandbox: <state>`. The banner already takes `host` and `dir`; sandbox state is computed from the same `dir` and a `lookPath` (the launchers already hold one). Rendered to stderr, suppressed on resume (matching the existing banner suppression).

- enabled → `Sandbox: enabled (safehouse)`
- available-but-not-enabled → `Sandbox: available, not enabled (no .safehouse profile)`
- unavailable → `Sandbox: unavailable (safehouse not on PATH)`

**Surface 2 — `status --boot`** (`printBoot` in `internal/status/boot.go`). Add a `SANDBOX` section alongside the existing `TEAM_STATE` / `STATE_BACKEND` sections, sourced in `gatherBoot` so the text and JSON renderers share one source of truth. Detection runs from the `status` package, which already does PATH lookups (`lookupExecutable`) and file stats — it does not import `internal/cli`, so the probe is a small local helper (`os.Stat(.safehouse)` + a PATH scan for `safehouse`), mirroring the safehouse package's two checks rather than calling across the package boundary. Rendered as:
  ```
  SANDBOX: <state>
  ```
  with the same three state strings as the banner.

**Surface 3 — `--version`** (`printVersion` in `internal/cli/cli.go`). The first line is **unchanged** — `spacedock <ver> (contract <N>)` — because the FO/ensign skills parse the `contract <N>` token from it (`first-officer-shared-core.md` step 1). Below it, append the same `Sandbox: <state>` line (computed against the cwd), then a per-runtime block:
  ```
  spacedock 0.20.0 (contract 1)
  Sandbox: unavailable (safehouse not on PATH)
  claude: installed, spacedock enabled
  codex: installed
  pi: not installed
  ```
  Per-runtime detection differs by host because the hosts differ:
  - **claude / codex** — `installed` = the host binary resolves on PATH. `spacedock enabled` = the Spacedock plugin is installed AND enabled in that host. For claude this reads `claude plugin list --json`: the spacedock@spacedock entry carries both `installPath` (install) and an `enabled` boolean (today `host_exec.go`'s `pluginListEntry` reads `id`+`installPath` but NOT `enabled` — this task adds the field). For codex the existing `codexEntryInstalled` text-parse gives install; enablement reuses the same listing's status field.
  - **pi** — pi has **no plugin-list model** (verified: `pi plugin list` is not a command). `installed` = the `pi` binary resolves; `spacedock enabled` = the existing `checkPiRuntime` readiness check (skills + extension present). So pi's "enabled" reuses `pi.go`'s `piRuntimeLaunchReady`, NOT a plugin probe.
  - **Probe failure ≠ not-installed.** A host CLI can fail for reasons other than absence — in a sandboxed session `codex plugin list` returns "Operation not permitted" and `pi` emits EPERM warnings (both observed live on this machine). When the host binary resolves but the enablement probe errors, render `<host>: installed, enablement unknown` rather than silently downgrading to "not installed". Binary-absent renders `<host>: not installed`.

The detection mechanisms are mostly already in the tree (`safehouse.Present`/`Available`, `ResolveManifest`/`codexEntryInstalled`, `checkPiRuntime`); the one new read is the claude `enabled` field, exercised in the spike below.

## Out of scope

- No new sandbox/safehouse behavior — this task only **reports** posture, it does not change when a launch is wrapped.
- No JSON schema for `--version` (it stays plain text the skills line-parse). `status --boot --json` does gain a `sandbox` field, since boot already has a JSON renderer that must stay in sync.
- No caching/daemonizing of the per-runtime probes. `--version` accepts the ~0.24s `claude plugin list` cost (measured) and bounded codex/pi probes; if a probe exceeds a small timeout it renders `enablement unknown`.
- No change to the `spacedock <ver> (contract <N>)` first line of `--version`.
- pi is not added to the launcher-banner sandbox line beyond the shared three-state computation (pi launches do not currently route through safehouse; the banner reports the same profile/binary state regardless of host).

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 — The launcher banner carries a `Sandbox:` line whose text matches the three-way state computed from `(profile-present, binary-available)`.**
Verified by: a Go test in `internal/cli` that calls `launchBanner` (or its renderer) with each of the three `(present, available)` input combinations and asserts the exact rendered stderr line for each — independent expected strings, the same proof shape `launch_banner_wording_test.go` already uses. The expected value is the state string, supplied by the test, not read from the file under render.

**AC-2 — `status --boot` emits a `SANDBOX:` section whose state matches the profile/binary inputs.**
Verified by: a Go test in `internal/status` that drives `printBoot` over a fixture workflow with (a) a `.safehouse` file present and (b) absent, against a `lookPath` stub pinning safehouse found / not-found, and asserts the rendered `SANDBOX:` line for each combination. The `--boot --json` form gains a `sandbox` field asserted in the same test. Expected strings are test-supplied independent values.

**AC-3 — `--version` keeps `spacedock <ver> (contract <N>)` as its first line AND appends the `Sandbox:` line plus a per-runtime block.**
Verified by: a Go test in `internal/cli` that captures `printVersion` output and asserts (1) the first line still matches `^spacedock .* \(contract \d+\)$` (the token the FO skill parses — this guards against the regression of breaking skill contract-detection), and (2) the per-runtime lines render `installed` / `installed, spacedock enabled` / `not installed` / `installed, enablement unknown` from injected host-probe stubs. Host detection is behind an injected `hostOps`-style seam so the test pins each outcome without shelling real CLIs.

**AC-4 — claude `spacedock enabled` reflects the plugin's `enabled` field, not merely its presence.**
Verified by: a Go test feeding `resolveClaude…` (the enablement reader) a `plugin list --json` fixture with the spacedock entry's `enabled` set true vs false and asserting `enabled` vs `installed (not enabled)`. The fixture JSON is the independent source of truth; a fixture with `enabled:false` must NOT render "enabled". (This is the field added in the spike below — the existing `pluginListEntry` does not read it.)

**AC-5 — A host whose enablement probe errors (not absence) renders `enablement unknown`, never silently `not installed`.**
Verified by: a Go test injecting a host probe that returns (binary-resolves=true, probe-error) and asserting the line is `<host>: installed, enablement unknown`; a separate case with binary-absent asserts `<host>: not installed`. Distinguishes the sandboxed-session failure mode observed live (`codex plugin list` → "Operation not permitted") from genuine absence.

## Riskiest-mechanism spike (run first)

The design composes mostly-proven seams (`safehouse.Present`/`Available`, `ResolveManifest`, `codexEntryInstalled`, `checkPiRuntime` all exist). The one unverified read is the claude plugin `enabled` field. Exercised on the dev machine during ideation:

- `claude plugin list --json` spacedock entry (live): `{"id":"spacedock@spacedock", ... "enabled":true, "installPath":"/Users/clkao/.claude/plugins/cache/spacedock/spacedock/0.19.9", ...}` — the `enabled` boolean is present and distinct from `installPath`. Probe cost ~0.24s (measured via `time`).
- `pi plugin list` is **not a command** — pi emits EPERM/no-API-key warnings; there is no plugin-list model. Confirms pi enablement must come from `checkPiRuntime`, not a plugin probe.
- `codex plugin list` in this sandboxed session returns "Operation not permitted (os error 1)" reading its config — confirms probe-failure must be a distinct rendered state, not "not installed".
- `which safehouse` → not found, while `.safehouse` exists in the repo root — confirms the `available-but-not-enabled` vs `unavailable` distinction is a real, currently-live state, not hypothetical.

These four observations seed the first implementation tests (AC-4 fixture from the live JSON shape; AC-5 from the codex permission-error mode).

## Test plan

All AC proofs are Go unit/golden tests driven by the built binary's renderers with injected probe seams — no live host CLIs in the test path (the live probes were the ideation spike; tests use fixtures). Estimated complexity: low-to-moderate.

- `internal/cli`: banner test (AC-1), version test (AC-3, AC-5), claude-enabled reader test (AC-4). Reuse the existing `hostOps` injection seam in `cli`/`host_exec.go` and the `lookPath` injection already used by the launchers.
- `internal/status`: boot `SANDBOX` text + JSON test (AC-2), extending the existing golden-fixture harness (`golden_harness_test.go`) and `lookupExecutable`/stat probes already in `boot.go`.
- No live-workflow smoke test is needed — none of the new output depends on a running team or a real host session; runtime behavior is not the claim. A one-off manual `./spacedock --version` after the build is a sanity check, not an AC proof.
- Cost note: the per-runtime probes run only on `--version`, not on every command; the FO contract-gate path reads only the first line and is unaffected.

## Documentation diff

User-visible CLI output changes, so per the ideation doc-diff rule the concrete before/after is recorded here for implementation to apply.

**`docs/install-journey.md`** — two `--version` descriptions (Homebrew step 2 at line 34, source-build step at line 75) currently read:

> Prints the installed version, e.g. `spacedock 0.20.0`.

After (both occurrences):

> Prints the installed version plus the sandbox posture and per-runtime install/enablement, e.g.:
>
> ```
> spacedock 0.20.0 (contract 1)
> Sandbox: unavailable (safehouse not on PATH)
> claude: installed, spacedock enabled
> codex: installed
> pi: not installed
> ```

**`README.md`** — line 139 currently reads:

> `spacedock --version                                       # print the installed version`

After:

> `spacedock --version                                       # version + sandbox posture + per-runtime status`

No other user-facing doc describes the banner or `status --boot` output verbatim, so those surfaces need no doc edit beyond the above.

## Stage Report: ideation

- DONE: Pin exactly which surfaces carry sandbox state and the rendered format of each: startup banner, `status --boot`, and `--version` — plus the per-host (claude/codex/pi) "spacedock enabled in this runtime" detection.
  Proposed approach pins all three surfaces with verbatim render strings and the three-way state mapping; per-host detection differs (claude/codex plugin list `enabled`, pi via `checkPiRuntime`); probe-failure rendered as a distinct `enablement unknown` state.
- DONE: Run the riskiest-mechanism check FIRST: how sandbox state and per-runtime install/enablement are actually detected on this machine.
  Spike section records live probes: claude `enabled` field present (~0.24s); `pi plugin list` is not a command; `codex plugin list` denied in-sandbox ("Operation not permitted"); `which safehouse` not-found while `.safehouse` present (the available-vs-unavailable distinction is live).
- DONE: Propose the concrete doc diff for the user-visible CLI-output change and a test plan whose AC proof is build/command output, not prose.
  Documentation diff section gives before/after for `docs/install-journey.md` (two spots) and `README.md:139`; all five ACs proven by Go renderer tests with injected probe seams (banner stderr, boot golden+JSON, version first-line regex + per-runtime stubs, claude `enabled` fixture, probe-error case).

### Summary
Pinned three surfaces (launcher banner, `status --boot` SANDBOX section, `--version` Sandbox line + per-runtime block) reporting the three-way sandbox state from the existing `safehouse.Present`/`Available` probes. The riskiest unverified mechanism — the claude plugin `enabled` field — was exercised live (it exists and is distinct from `installPath`); the spike also surfaced two design-shaping facts: pi has no plugin-list model (enablement reuses `checkPiRuntime`) and host probes can be sandbox-denied (so probe-failure renders as `enablement unknown`, never silent `not installed`). Key constraint recorded: `--version`'s first line must stay `spacedock <ver> (contract <N>)` so the FO/ensign skills' token parse keeps working — AC-3 guards it with a first-line regex.
