---
id: xdhbscjmpv0prk3k5yz1ngpw
title: Paseo-wide FO runtime support and generic unsupported-harness fallback hints
status: implementation
source: captain (2026-11-23)
started: 2026-11-23T23:47:27Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-paseo-and-generic-harness-fo-support
issue:
pr:
mod-block:
gates:
    version: 1
    records:
        - id: gate:xdhbscjmpv0prk3k5yz1ngpw:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:xdhbscjmpv0prk3k5yz1ngpw-backlog-1
              briefing:
                id: briefing:xdhbscjmpv0prk3k5yz1ngpw:backlog:attempt-1:revision-1
                digest: sha256:69a37a507dadd7f2fba03e3f3b7bdaea1f1a17f05c5036e44db8d3a34a79d253
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:xdhbscjmpv0prk3k5yz1ngpw:backlog:1
                briefing: briefing:xdhbscjmpv0prk3k5yz1ngpw:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-11-23T23:47:18.861991Z"
                decision: approve
                reason: 'Captain directive: dispatch xdh and sweep unnecessary operating-contract loads on Paseo path'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:xdhbscjmpv0prk3k5yz1ngpw:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:xdhbscjmpv0prk3k5yz1ngpw-ideation-1
              briefing:
                id: briefing:xdhbscjmpv0prk3k5yz1ngpw:ideation:attempt-1:revision-1
                digest: sha256:c7bb97ee19aa2a5df68a07d06161f185a29ed2cb57fa2cd6bd15215c1e54d898
                room-ref: '@review/ideation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:xdhbscjmpv0prk3k5yz1ngpw:ideation:1
                briefing: briefing:xdhbscjmpv0prk3k5yz1ngpw:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-11-24T02:21:35.72354Z"
                decision: approve
                reason: 'Captain directive: proceed on all'
              application:
                target-stage: implementation
                state: consumed
---

Make any Paseo-hosted first officer boot and run cleanly regardless of inner harness (fixing `F4` and Paseo-vs-native dispatch selection), eliminate unnecessary or contradictory operating-contract loads on the Paseo FO & ensign path while preserving the 23,500-byte shared-core cap, and give unsupported harnesses clear fallback guidance to map Spacedock's worker capabilities onto their own Claude/Codex-like native features.

This design composes cleanly atop PR #893 (`936` / `paseo-dispatch-backend`) and PR #895 (`x9m` / `paseo-worker-adoption-and-recovery`).

## Problem

1. **Non-{Claude,Codex,Pi} Paseo harnesses abort on Boot Step 1 (`F4`):** `internal/cli/cli.go:913-930` (`printVersion`) returns early after 3 lines — omitting `Runtime:`, `Sandbox:`, and `contract 3` — whenever `runtimehost.Detect` sees none of `CLAUDECODE`, `CODEX_THREAD_ID`, or `PI_CODING_AGENT_DIR`. Any FO running inside Paseo on Antigravity (`agy`), Goose, OpenCode, Copilot, Cursor, etc. immediately aborts on `first-officer-shared-core.md:9`'s `contract 3` binary gate.
2. **Wasted & contradictory inner-host runtime loads under active Paseo:** In `skills/first-officer/SKILL.md` (atop PR #893), Claude (`CLAUDECODE`), Codex (`CODEX_THREAD_ID`), and Pi (`PI_CODING_AGENT_DIR`) FOs unconditionally load `claude-first-officer-runtime.md` (5,087 B), `codex-first-officer-runtime.md` (8,483 B), or `pi-first-officer-runtime.md` (8,259 B) at boot BEFORE checking Paseo. 95–98% of those files are native subagent transport rules (`Agent()`/`SendMessage`, `spawn_agent`/`wait_agent`/`followup_task`, `subagent`/`subagent_wait`/`intercom`) that are unused by Paseo and actively contradict it (e.g. `codex-first-officer-runtime.md:30-38` unconditional `wait_agent` and L52–54 mandatory reviewer reuse; `claude-first-officer-runtime.md:13,19` `SendMessage(shutdown_request)` and keep-worker-alive-at-gate).
3. **Unconditional `PASEO_AGENT_ID` selection breaks sandboxed FOs & late load misses resume reconciliation:**
   - Once upstream `paseo-safehouse` (`tv`) preserves `PASEO_AGENT_ID` for all sandboxed agents, a normal sandboxed Claude/Codex/Pi FO running under `paseo-access=workspace` (no daemon access) will have `PASEO_AGENT_ID` set too. Selecting Paseo without checking daemon reachability (`paseo status --json` → `connectedDaemon == "reachable"`) breaks native fallback with `EPERM`/`command not found` — and on the ensign side, naive `PASEO_AGENT_ID`-first matching in `skills/ensign/SKILL.md` would misroute native subagents spawned by that fallback FO onto `paseo-ensign-runtime.md` instead of `claude`/`codex`/`pi-ensign-runtime.md`.
   - `skills/first-officer/SKILL.md:38` and `paseo-fo-dispatch.md:3` currently say "at first dispatch", which is too late for resume reconciliation (`«roster-reconcile»()` at `fo-dispatch-core.md:231` and `paseo-fo-dispatch.md` Steps 0–3, flagged in `pre-merge-audit.md` §4.1).
4. **Remaining runtime-specific contradictions in shared cores:**
   - `fo-dispatch-core.md:188` asserts `` `host` derives from the runtime (`--host` is for tests/cross-host tooling only) `` — directly contradicting `paseo-fo-dispatch.md:16`, which requires `--host paseo` on every production Paseo dispatch.
   - `first-officer-shared-core.md` (L45, L46, L52, L82), `fo-dispatch-core.md` (L61), and `feedback-rejection-flow/SKILL.md` (L27) still embed Claude-only `Skill(skill="spacedock:...")` tool syntax instead of host-neutral `` `spacedock:...` ``.
   - `first-officer-shared-core.md:76` assumes `the worker kept alive` unconditionally at a gate, conflicting with Paseo v1's Row 4 archival of verified fresh workers.
   - `ensign-shared-core.md` `## DISPATCH_FILE Bootstrap` documents Claude (`Skill(skill=...)`), Codex (`$spacedock:ensign`), and Pi, but omits the Paseo / host-neutral pointer pattern.
5. **No fallback guidance for unsupported harnesses:** A harness without a dedicated adapter and without reachable Paseo (e.g. standalone Antigravity, Cursor, OpenCode, Goose) hits two blockers: (a) 3-line `spacedock --version` output omits `contract 3`, tripping Boot Step 1; and (b) `SKILL.md` gives no hint on mapping Spacedock's capabilities onto native Claude/Codex-like subagent/wait tools or using `--host paseo` for host-neutral dispatch artifacts.

## Proposed approach

### 1. Binary `--version` Paseo fallback (`F4`, `internal/cli/cli.go`)

Keep `internal/runtimehost.Detect` untouched (preserving zero ambiguity when `PASEO_AGENT_ID` coexists with `CLAUDECODE`, `CODEX_THREAD_ID`, or `PI_CODING_AGENT*`, and preserving `dispatch build`'s explicit `--host paseo` invariant).

In `internal/cli/cli.go` `printVersion`:

```go
	host, markers, ambiguous := runtimehost.Detect(getenv)
	if !ambiguous && host == "" {
		if getenv("PASEO_AGENT_ID") != "" {
			host = "paseo"
			markers = []string{"PASEO_AGENT_ID"}
		} else {
			// Outside every runtime — a human at a terminal. Three lines: the version
			// line, the OS line, and the Channel line, nothing else.
			return
		}
	}
```

Observable behavior matrix:
- **Outside every runtime** (no native marker, `PASEO_AGENT_ID` unset): unchanged **3 lines** (`spacedock <version>`, `OS:`, `Channel:`).
- **Non-{Claude,Codex,Pi} Paseo session** (`PASEO_AGENT_ID` set, no native marker): **6 lines**, emitting `Runtime: paseo (PASEO_AGENT_ID)`, `Sandbox: <SessionState>`, and `contract 3`. Never renders the UUID value of `PASEO_AGENT_ID`, never execs `paseo` or any host CLI.
- **Claude/Codex/Pi inside Paseo** (native marker + `PASEO_AGENT_ID`): `host != ""` takes precedence, emitting `Runtime: claude (CLAUDECODE)` / `codex` / `pi` unchanged with zero ambiguity.
- **Ambiguous native markers** (e.g. `CLAUDECODE` + `CODEX_THREAD_ID`, with or without `PASEO_AGENT_ID`): emits `Runtime: ambiguous (...)` at exit `0`, unchanged.

**Mechanism-to-value check:** Serves **AC-1 (VALUE)**. Simplest alternative considered: adding `PASEO_AGENT_ID` to `runtimehost.Detect` or keying on `PASEO_HOME`. Rejected because adding `PASEO_AGENT_ID` to `runtimehost.Detect` turns every Paseo-hosted Claude/Codex/Pi session into `ambiguous` (breaking `dispatch build` auto-detection and `--version`), while `PASEO_HOME` can be set in a human terminal profile outside any agent session (breaking 3-line human output).

### 2. Full Operating-Contract Audit on the Paseo FO & Ensign Path

Every file in the Paseo FO and ensign load graph was audited line-by-line (sizes measured against `main` + PR #893/#895):

| File | Current Size | Loaded When (Paseo) | Audit Finding | Action |
|---|---:|---|---|---|
| `skills/first-officer/SKILL.md` | 2,811 B (3,180 B on `936`) | FO boot | Loads inner-host adapters unconditionally before Paseo; no daemon-reachability check; "at first dispatch" misses resume reconciliation; no unsupported-harness fallback. | **Update** priority-ordered `## Runtime adapter` (exact wording below). |
| `claude-first-officer-runtime.md` | 5,087 B | **Skipped** under active Paseo | ~95% native Claude transport (`Agent`/`SendMessage`/`claude-fo-dispatch.md`) or duplicated from `fo-write-core.md` (`spacedock new`) & `fo-gate-lifecycle`. Only unique non-transport item is Claude `Read`-then-Bash cache-write quirk (optional micro-opt). | **Leave untouched; skip at boot** when Paseo daemon is reachable (**-5,087 B** at boot). |
| `codex-first-officer-runtime.md` | 8,483 B | **Skipped** under active Paseo | ~98% native Codex transport (`spawn_agent`/`wait_agent`/`followup_task`/`list_agents`, mandatory reviewer reuse) — actively contradicts Paseo (`pre-merge-audit.md` §4.1). | **Leave untouched; skip at boot** when Paseo daemon is reachable (**-8,483 B** at boot). |
| `pi-first-officer-runtime.md` | 8,259 B | **Skipped** under active Paseo | ~98% native Pi transport (`--host pi`, `subagent`, `subagent_wait`, `intercom`, `pi-agent-teams`, Pi test harness). `«gate.lifecycle»` is already owned by `fo-gate-lifecycle/SKILL.md`. | **Leave untouched; skip at boot** when Paseo daemon is reachable (**-8,259 B** at boot). |
| `first-officer-shared-core.md` | 23,497 B (cap 23,500 B — 3 B left!) | FO boot | 4 Claude-only `Skill(skill="spacedock:...")` wrappers (L45, L46, L52, L82); L76 `the worker kept alive` conflicts with Paseo Row 4 archival. `fo-dispatch-recovery` (L52) is deferred and never triggered by Paseo. | **Tighten in-place** (-43 B → **23,454 B**, 46 B headroom under 23,500 B cap; zero cap raise). |
| `fo-dispatch-core.md` | 33,609 B | Deferred (`«engage»` / dispatch / resume) | L188 `` `--host` is for tests/cross-host tooling only `` directly contradicts `paseo-fo-dispatch.md:16`; L61 has `Skill(skill="spacedock:fo-gate-lifecycle")`. ~4.5 KB of reuse/advance/`«reuse.model-match»` prose is unused on Paseo v1 (short-circuits cleanly on `«addressable-worker»` ABSENT). | **Minimal 2-line fix** (L61 & L188). **Do NOT split** reuse prose out (necessity-first: splitting adds a new deferred load trigger to the capped shared core and an extra read hop on every Claude/Codex stage advance). |
| `paseo-fo-dispatch.md` (`x9m`) | 20,103 B | Deferred (alongside `fo-dispatch-core.md`) | L3 says "at the first worker dispatch" (contradicting its own L101 resume reconciliation) and lists only 4 inner harnesses. Rest of file (5 verbs, local lease, Steps 0–3, 11-row Recovery Matrix) is clean and host-neutral. | **Minimal 1-paragraph alignment** on L3 only (preserves PR #895's Recovery Matrix & lease intact). |
| `skills/ensign/SKILL.md` | 578 B (700 B on `936`) | Ensign boot | Lacks explicit first-match-wins rule; naive `PASEO_AGENT_ID`-first rule misroutes native subagents spawned by a sandboxed `paseo-access=workspace` fallback FO onto `paseo-ensign-runtime.md`. | **Update** with exact first-match-wins precedence + generic fallback (exact wording below). |
| `ensign-shared-core.md` | 11,006 B | Ensign boot | `## DISPATCH_FILE Bootstrap` (L99–129) omits the Paseo/host-neutral pointer pattern. ~900 B `Advance bootstrap` is unused by Paseo v1 fresh workers (harmlessly gated). | **Add 2-line Paseo/host-neutral pointer pattern** to `## DISPATCH_FILE Bootstrap`; do NOT split `Advance bootstrap`. |
| `paseo-ensign-runtime.md` (`936`) | 1,595 B | Ensign boot (Paseo / generic) | 10 lines, zero host leaks, clean Git completion gate. | **Leave untouched** (or optional 1-clause note that generic unsupported harnesses reuse it; no structural change needed). |
| Other deferred FO skills (`fo-write-core`, `fo-merge-core`, `fo-install`, `fo-gate-lifecycle`, `fo-status-viewer`, `present-gate`, `feedback-rejection-flow`, `fo-dispatch-recovery`) | — | Triggered on demand | All host-neutral except one `Skill(skill="spacedock:fo-gate-lifecycle")` in `feedback-rejection-flow/SKILL.md:27`; `fo-dispatch-recovery` is Claude-only and never loaded on Paseo. | **1-line cleanup** in `feedback-rejection-flow/SKILL.md:27` only. |

**Mechanism-to-value check (Paseo priority, daemon reachability check, and skipping inner-host adapters):** Serves **AC-2 (VALUE)** and **AC-3**.
- Simplest alternative considered for daemon reachability: (a) unconditional `PASEO_AGENT_ID` selection without probing, or (b) checking only `paseo status --json` exit code. Both fail: (a) breaks sandboxed `paseo-access=workspace` Claude/Codex/Pi FOs once `paseo-safehouse` `tv` preserves `PASEO_AGENT_ID`; (b) fails on a stopped daemon because `paseo status --json` exits `0` with `"connectedDaemon": "not_probed"` (proven by live spike below). Requiring exit `0` AND `"connectedDaemon": "reachable"` — and running the probe ONLY when `PASEO_AGENT_ID` is set or Paseo is explicitly requested — gives zero overhead on non-Paseo sessions and exact fallback on sandboxed/stopped daemons.
- Simplest alternative considered for `fo-dispatch-core.md` / `ensign-shared-core.md` reuse prose: splitting reuse/advance into a new deferred file (`fo-reuse-core.md`). Rejected as disproportionate: `fo-dispatch-core.md` is already deferred off boot, the reuse rules already short-circuit cleanly on `«addressable-worker»` ABSENT, and splitting adds a new trigger to the 23,500 B capped shared core plus an extra read hop on Claude/Codex.

### 3. Unsupported-Harness Fallback Hint (without adding `--host generic`)

When neither reachable Paseo nor a dedicated Claude/Codex/Pi adapter applies (e.g. standalone Antigravity, Cursor, OpenCode, Goose, or an unsupported harness inside a sandboxed Paseo session without daemon access), `skills/first-officer/SKILL.md` provides a concise fallback branch that:
1. Directs captain interaction to direct conversation text.
2. Authorizes Boot Step 1 to accept 3-line `spacedock --version` output matching binary minor `0.29` without aborting on the absent `contract 3` line.
3. Directs worker dispatch to build host-neutral assignments with `spacedock dispatch build --host paseo` (paired with `references/paseo-ensign-runtime.md` via `skills/ensign/SKILL.md`'s generic fallback) and map Spacedock's capabilities onto the harness's own native subagent features resembling Claude (`Agent` / background subagent) or Codex (`spawn_agent` / `wait_agent` / `list_agents`), keeping durable Git commits + `## Stage Report` as the authoritative completion gate.

**Mechanism-to-value check:** Serves **AC-5**. Simplest alternative considered: adding a fifth `--host generic` target to `internal/dispatch/build.go` and a new `generic-ensign-runtime.md`. Rejected because `--host paseo` already emits a completely host-neutral schema-v2 envelope (`model: null`, no Claude/Codex/Pi syntax, `Done: ...` completion signal, Git-state completion gate) — adding `--host generic` would duplicate 100% of that code path and test suite for zero behavioral gain.

### 4. Exact Before/After Skill Text Changes

#### A. `skills/first-officer/SKILL.md` (`## Runtime adapter`, lines 31–40 atop PR #893)

**Before (on PR #893):**
```markdown
## Runtime adapter

Load the runtime adapter for your platform:
- Claude Code (`CLAUDECODE` env var is set): read exactly `{first_officer_base}/references/claude-first-officer-runtime.md`; do not list or search the references directory.
- Codex (`CODEX_THREAD_ID` env var is set): read `references/codex-first-officer-runtime.md`
- Pi (`PI_CODING_AGENT_DIR` is set, or this session is running under Pi without the Claude/Codex markers above): read `references/pi-first-officer-runtime.md`

**Paseo dispatch backend:** When `PASEO_AGENT_ID` is set in the environment or the captain explicitly selects Paseo dispatch, bind worker dispatch at first dispatch by reading `{first_officer_base}/references/paseo-fo-dispatch.md` alongside `{first_officer_base}/references/fo-dispatch-core.md` instead of the host's native subagent transport (on hosts without a dedicated adapter above, communicate with the captain via direct conversation text).

Then begin the Startup procedure from the shared core.
```

**After:**
```markdown
## Runtime adapter

Select the first matching runtime below (load only one transport; never mix native subagent tools with Paseo):

1. **Paseo (priority when `PASEO_AGENT_ID` is set or captain explicitly selects Paseo):** Skip this check (run no `paseo` command) when `PASEO_AGENT_ID` is unset and Paseo was not requested. Otherwise verify daemon reachability with `command -v paseo >/dev/null && paseo status --json` and require exit `0` with `"connectedDaemon": "reachable"`.
   - **Reachable:** Paseo is the active runtime. Do NOT read `claude-first-officer-runtime.md`, `codex-first-officer-runtime.md`, or `pi-first-officer-runtime.md` at boot. Communicate with the captain in direct conversation text. Defer `{first_officer_base}/references/paseo-fo-dispatch.md` alongside `{first_officer_base}/references/fo-dispatch-core.md` until the dispatch/resume boundary (before resume reconciliation, `«dispatch.next-action»()`, worker dispatch, or terminal teardown — not on a greet-and-stop boot). Paseo owns all worker lifecycle bindings (`spawn`, `wait`, `roster-reconcile`, `shutdown`).
   - **Unreachable:** If the captain explicitly requested Paseo dispatch, report the unreachable daemon and stop. If only ambient `PASEO_AGENT_ID` was set (for example a sandboxed `paseo-access=workspace` session), fall through to step 2.
2. **Dedicated host adapter (when Paseo is not active):**
   - Claude Code (`CLAUDECODE` env var is set): read exactly `{first_officer_base}/references/claude-first-officer-runtime.md`; do not list or search the references directory.
   - Codex (`CODEX_THREAD_ID` env var is set): read `references/codex-first-officer-runtime.md`
   - Pi (`PI_CODING_AGENT_DIR` is set, or this session is running under Pi without the Claude/Codex markers above): read `references/pi-first-officer-runtime.md`
3. **Unsupported-harness fallback (neither reachable Paseo nor Claude/Codex/Pi):** Communicate with the captain in direct conversation text. In Startup step 1, `spacedock --version` prints 3 lines (`spacedock <version>`, `OS:`, `Channel:`) without `Runtime:` or `contract 3` when no host marker or `PASEO_AGENT_ID` is set — accept a matching binary minor (`0.29`) on line 1 instead of aborting on the missing `contract 3` line. At the dispatch boundary, read `{first_officer_base}/references/fo-dispatch-core.md`, build host-neutral assignments with `spacedock dispatch build --host paseo`, and map worker capabilities onto this harness's native features resembling Claude (`Agent` / background subagent) or Codex (`spawn_agent` / `wait_agent` / `list_agents`): launch fresh isolated workers with the intact pointer (`«worker.spawn»`), wait/observe asynchronously when supported (`«async-dispatch»`), treat `«addressable-worker»` and `«context-budget»` as ABSENT unless a live follow-up handle exists, use native list/close tools when available (`«roster-reconcile»` / `«worker.shutdown»`), and treat committed Git state plus `## Stage Report` as the authoritative completion gate.

Then begin the Startup procedure from the shared core.
```

#### B. `skills/first-officer/references/first-officer-shared-core.md` (net **-43 B**: 23,497 B → **23,454 B**, cap 23,500 B)

- L45: ``- `Skill(skill="spacedock:fo-status-viewer")` `` → ``- `spacedock:fo-status-viewer` ``
- L46: ``- `Skill(skill="spacedock:fo-gate-lifecycle")` `` → ``- `spacedock:fo-gate-lifecycle` ``
- L52: ``- `Skill(skill="spacedock:fo-dispatch-recovery")` `` → ``- `spacedock:fo-dispatch-recovery` ``
- L76: `the worker kept alive.` → `an addressable worker kept alive.`
- L82: ``- **effect:** invoke `Skill(skill="spacedock:feedback-rejection-flow")`;`` → ``- **effect:** load `spacedock:feedback-rejection-flow`;``

#### C. `skills/first-officer/references/fo-dispatch-core.md` (2-line fix)

- L61: ``Other gates complete `Skill(skill="spacedock:fo-gate-lifecycle")`, then `«gate.lifecycle»(slug, stage)`.`` → ``Other gates load `spacedock:fo-gate-lifecycle`, then run `«gate.lifecycle»(slug, stage)`.``
- L188: `` `host` derives from the runtime (`--host` is for tests/cross-host tooling only). `` → `` `host` derives from the runtime unless the active dispatch adapter requires an explicit `--host` (such as `--host paseo`). ``

#### D. `skills/first-officer/references/paseo-fo-dispatch.md` (L3 atop PR #895)

**Before:**
```markdown
The Paseo dispatch parts (deferred from `fo-dispatch-core.md`), read alongside `fo-dispatch-core.md` at the first worker dispatch when `PASEO_AGENT_ID` is set or the captain selects Paseo dispatch. Any Paseo-hosted first officer (Claude, Codex, Pi, or Antigravity) drives fresh workers through five `paseo` CLI verbs and `spacedock dispatch build --host paseo`.
```

**After:**
```markdown
The Paseo dispatch parts (deferred from `fo-dispatch-core.md`), read alongside `fo-dispatch-core.md` at the dispatch/resume boundary (before resume reconciliation, `«dispatch.next-action»()`, worker dispatch, or terminal teardown) when Paseo is active (reachable `paseo` daemon with `PASEO_AGENT_ID` set, or explicit captain selection). Any Paseo-hosted first officer (regardless of inner harness) drives fresh workers through five `paseo` CLI verbs and `spacedock dispatch build --host paseo`.
```

#### E. `skills/ensign/SKILL.md` (`## Runtime adapter`, lines 10–18 atop PR #893)

**Before:**
```markdown
## Runtime adapter

Load the runtime adapter for your platform:
- Paseo (`PASEO_AGENT_ID` env var is set, or the dispatch prompt names `references/paseo-ensign-runtime.md`): read `references/paseo-ensign-runtime.md`
- Claude Code (`CLAUDECODE` env var is set): read `references/claude-ensign-runtime.md`
- Codex (`CODEX_THREAD_ID` env var is set): read `references/codex-ensign-runtime.md`
- Pi (`PI_CODING_AGENT_DIR` is set, or this session is running under Pi without the Claude/Codex markers above): read `references/pi-ensign-runtime.md`

Then read your assignment and begin work.
```

**After:**
```markdown
## Runtime adapter

Load ONLY the first matching runtime adapter (do not load a second inner-host adapter):
- Paseo or host-neutral dispatch (the dispatch prompt names `references/paseo-ensign-runtime.md`, OR `PASEO_AGENT_ID` is set with none of `CLAUDECODE`, `CODEX_THREAD_ID`, `PI_CODING_AGENT_DIR`): read `references/paseo-ensign-runtime.md`
- Claude Code (`CLAUDECODE` env var is set): read `references/claude-ensign-runtime.md`
- Codex (`CODEX_THREAD_ID` env var is set): read `references/codex-ensign-runtime.md`
- Pi (`PI_CODING_AGENT_DIR` is set, or this session is running under Pi without the Claude/Codex markers above): read `references/pi-ensign-runtime.md`
- Other harness (none of the above): read `references/paseo-ensign-runtime.md`

Then read your assignment and begin work.
```

*(Why this exact precedence matters: a native Claude/Codex/Pi subagent spawned by a sandboxed `paseo-access=workspace` fallback FO inherits ambient `PASEO_AGENT_ID` AND its native marker, but its prompt does NOT name `references/paseo-ensign-runtime.md` — so it skips bullet 1 and correctly loads its native `claude`/`codex`/`pi-ensign-runtime.md`!)*

#### F. `skills/ensign/references/ensign-shared-core.md` (`## DISPATCH_FILE Bootstrap`, L104–112)

Insert right after the `Codex:` example (before the prose paragraph at L112):
```markdown
Paseo / host-neutral:

    Load the Spacedock ensign skill (with references/ensign-shared-core.md and references/paseo-ensign-runtime.md); then Read /tmp/spacedock-dispatch/{name}.md and treat its content as your assignment.
```
And change L120 `When the initial prompt begins with either host-specific pattern` → `When the initial prompt begins with any host-specific pattern above`.

#### G. `skills/feedback-rejection-flow/SKILL.md` (L27)

- L27: ``Invoke `Skill(skill="spacedock:fo-gate-lifecycle")` and run `«gate.lifecycle»` `` → ``Load `spacedock:fo-gate-lifecycle` and run `«gate.lifecycle»` ``

### 5. Proposed User-Facing Documentation Diff (ASD-STE100)

```diff
--- a/docs/site/reference/command-reference.md
+++ b/docs/site/reference/command-reference.md
@@ -4,7 +4,7 @@ Run `spacedock <command> --help` for the full options. Run `spacedock` for the c
 
 ## Version
 
-`spacedock --version` reports the version, platform and release channel. Inside an agent session, it also reports the runtime and sandbox.
+`spacedock --version` reports the version, platform, and release channel. Inside a supported or Paseo-hosted agent session, it also reports the runtime, sandbox, and contract level.
```

## Risk evidence

**Live spike executed during ideation (2026-11-23) on `paseo status --json` reachability semantics:**
1. Against the running local Paseo 0.17.1 daemon (`127.0.0.1:6769`): `paseo status --json` exited **`0`** with `"localDaemon": "running"` and `"connectedDaemon": "reachable"`.
2. Against an unreachable/stopped daemon (`paseo --home /tmp/spike-paseo-and-generic-harness-fo-support-empty status --json`, cleaned up after check): `paseo status --json` ALSO exited **`0`**, with `"localDaemon": "stopped"` and `"connectedDaemon": "not_probed"`.
3. **Design impact proven by spike:** Checking exit code alone would false-green on a stopped daemon and skip native fallback. The `SKILL.md` reachability rule therefore MUST require BOTH exit `0` AND `"connectedDaemon": "reachable"`.

**Proven mechanisms relied on without further spike:**
- `internal/cli/cli.go:printVersion` already threads injectable `getenv` and `lookPath` seams tested in `internal/cli/version_session_test.go`.
- `internal/dispatch/build.go` `--host paseo` (PR #893) already emits a host-neutral schema-v2 envelope and prompt with zero Claude/Codex/Pi tokens when no native host marker is set (`internal/dispatch/build_paseo_host_test.go`).
- `internal/contractlint/boot_resident_closure_test.go:42` matches lazy skills via `spacedock:([a-z0-9-]+)` without requiring `Skill(skill="...")`, and `internal/contractlint/fo_function_reference_invariant_test.go` enforces the 23,500-byte cap on `first-officer-shared-core.md`.

## Out of scope

- Adding `PASEO_AGENT_ID` or `PASEO_HOME` to `internal/runtimehost.Detect`.
- Adding a fifth `--host generic` target to `internal/dispatch/build.go` or a separate `generic-ensign-runtime.md`.
- Splitting reuse/advance/`«reuse.model-match»` prose out of `fo-dispatch-core.md` or `Advance bootstrap` out of `ensign-shared-core.md`.
- Changing PR #895's Paseo Recovery Matrix semantics, local lease schema, or candidate-scoping follow-ups (`pre-merge-audit.md` §4.2–4.5).
- Adding Paseo v2 addressable-worker / worker-reuse / child-to-parent messaging support.
- Upstream `paseo-safehouse` changes (`tv` env preservation or `F7` worktree `.git`/`.spacedock-state` grants).

## Expected surface and tolerance

Estimate net LOC change: `+90`, across `10` files.
- **Insertions:** ~113 lines
- **Deletions:** ~23 lines
- **Tolerance:** net `+40` to `+160` LOC, across `8` to `12` files.
- **Files expected:**
  1. `internal/cli/cli.go` (~+9 / -2)
  2. `internal/cli/version_session_test.go` (~+55 / -0)
  3. `skills/first-officer/SKILL.md` (~+26 / -6)
  4. `skills/first-officer/references/first-officer-shared-core.md` (~+5 / -5; **net -43 bytes**, 23,454 B <= 23,500 B cap)
  5. `skills/first-officer/references/fo-dispatch-core.md` (~+2 / -2)
  6. `skills/first-officer/references/paseo-fo-dispatch.md` (~+2 / -2)
  7. `skills/ensign/SKILL.md` (~+6 / -3)
  8. `skills/ensign/references/ensign-shared-core.md` (~+5 / -1)
  9. `skills/feedback-rejection-flow/SKILL.md` (~+1 / -1)
  10. `docs/site/reference/command-reference.md` (~+2 / -1)
- **Declared observable semantic changes:**
  - `spacedock --version` output when `PASEO_AGENT_ID` is set and no native host marker (`CLAUDECODE`, `CODEX_THREAD_ID`, `PI_CODING_AGENT*`) is set changes from 3 lines to 6 lines (`Runtime: paseo (PASEO_AGENT_ID)`, `Sandbox: ...`, `contract 3`).
  - Paseo-hosted FOs with a reachable `paseo` daemon skip reading `claude`/`codex`/`pi-first-officer-runtime.md` at boot and load `paseo-fo-dispatch.md` at the dispatch/resume boundary; sandboxed/unreachable-daemon FOs fall back to native adapters.
  - Unsupported harnesses outside Paseo accept 3-line `spacedock --version` (minor `0.29`) and dispatch via `--host paseo` mapped to native subagent tools.
- **Semantics that MUST NOT change:**
  - `spacedock --version` 3-line human terminal output when neither native host markers nor `PASEO_AGENT_ID` are set.
  - `internal/runtimehost.Detect` return values and `spacedock dispatch build` host auto-detection / ambiguity refusal.
  - `first-officer-shared-core.md` byte cap (`<= 23,500` bytes in `internal/contractlint/fo_function_reference_invariant_test.go`).
  - Native Claude, Codex, and Pi FO/ensign behavior when `PASEO_AGENT_ID` is unset.

## Acceptance criteria

**AC-1 (VALUE) - Any Paseo-hosted non-{Claude,Codex,Pi} harness passes the `--version` `contract 3` boot gate while bare terminal output and `runtimehost.Detect` stay unchanged.**
Verified by: `go test -race -count=1 ./internal/cli ./internal/runtimehost ./internal/dispatch` — specifically `TestVersionSessionRender` (new cases: `paseo-only-inside-safehouse`, `paseo-only-unsandboxed`, `paseo-marker-value-never-rendered`, `claude-inside-paseo-no-ambiguity`, plus unchanged `outside-every-runtime-is-three-lines`), `TestVersionContractTokenPlacement`, `TestVersionExecutesNoHostCLI`, and `TestBuildHostResolution` (`internal/dispatch/build_host_test.go`).
Falsifying edit: removing the `PASEO_AGENT_ID` fallback in `printVersion` flips `paseo-only-*` red; moving `PASEO_AGENT_ID` into `runtimehost.Detect` flips `claude-inside-paseo-no-ambiguity` and `dispatch build` with `CLAUDECODE`+`PASEO_AGENT_ID` red; rendering `getenv("PASEO_AGENT_ID")` flips `paseo-marker-value-never-rendered` red.

**AC-2 (VALUE) - Active Paseo FOs eliminate unnecessary inner-host adapter loads at boot (5.1–8.5 KB saved) while `first-officer-shared-core.md` shrinks strictly below its 23,497-byte baseline and Paseo end-to-end dispatch/adoption remains green.**
Verified by: (1) `wc -c skills/first-officer/references/first-officer-shared-core.md` reports `<= 23496` bytes AND `go test -race -count=1 ./internal/contractlint` passes `TestFOInstructionComponentCaps` (23,500 B cap unchanged); (2) targeted live Paseo run `SPACEDOCK_PASEO_LIVE_REQUIRED=1 go test -tags live -count=1 -timeout 20m -run '^TestLivePaseo(Dispatch|Adoption)Smoke$' ./internal/ensigncycle/` passes with committed code, committed split-root stage report, zero duplicate spawns, and worker archival, and inspecting the live FO session trace confirms zero reads of `claude-first-officer-runtime.md`, `codex-first-officer-runtime.md`, or `pi-first-officer-runtime.md` (serves end-value; note: requires local Paseo daemon + provider credentials outside an ordinary ensign sandbox — report explicitly if unavailable).
Falsifying edit: adding >46 bytes to `first-officer-shared-core.md` flips `TestFOInstructionComponentCaps` red; reverting `SKILL.md` to load `codex-first-officer-runtime.md` at boot causes the live Paseo Codex FO trace to record a read of `codex-first-officer-runtime.md`.

**AC-3 - Sandboxed `paseo-access=workspace` fallback and non-Paseo native lanes work cleanly without daemon errors or ensign misrouting.**
Verified by: (1) deterministic spike/control confirming `paseo --home <empty> status --json` exits `0` with `connectedDaemon != "reachable"` (demonstrating why the JSON field check is load-bearing) and `skills/ensign/SKILL.md` first-match precedence routes a prompt lacking `references/paseo-ensign-runtime.md` to the inner host's native adapter even when `PASEO_AGENT_ID` is set; (2) targeted native live non-regression runs on affected shared-contract lanes (`SPACEDOCK_LIVE_RUNTIME=claude SPACEDOCK_LIVE_MODEL=sonnet go test -tags live -count=1 -timeout 30m -run '^TestLiveCommon(ShallowBoot|FullEnsignCycle)$' ./internal/ensigncycle/` and the equivalent `codex`/`pi` shallow-boot checks as required by `docs/dev/README.md:85` — report any sandbox/quota unavailability explicitly rather than skipping silently).
Falsifying edit: gating Paseo solely on `PASEO_AGENT_ID` or `paseo status --json` exit code `0` breaks a sandboxed `paseo-access=workspace` Claude/Codex FO; reversing rule 1 in `skills/ensign/SKILL.md` to trigger on bare `PASEO_AGENT_ID` before `CLAUDECODE` causes a native Claude ensign with ambient `PASEO_AGENT_ID` to omit `SendMessage(to="main")` and stall `TestLiveCommonFullEnsignCycle`.

**AC-4 - Shared FO/ensign contract files are free of Claude-only `Skill(skill=...)` syntax and `--host` contradictions while preserving all `contractlint` structural invariants.**
Verified by: `go test -race -count=1 ./internal/contractlint` (`TestBootResidentDeferredLoadPointsResolve`, `TestUserSkillReferenceClosureResolves`, `TestFirstOfficerEntryEagerlyImportsOnlySharedCore`, `TestPaseoRuntimeAdaptersResolveOnDisk`, `TestFOInstructionComponentCaps`) passes, AND one-off validation grep `grep -n 'Skill(skill=' skills/first-officer/references/first-officer-shared-core.md skills/first-officer/references/fo-dispatch-core.md skills/feedback-rejection-flow/SKILL.md` returns empty (one-off validation evidence only — NOT a committed prose-grep test).
Falsifying edit: breaking any `references/...` or `spacedock:...` target in the edited skill files flips `TestBootResidentDeferredLoadPointsResolve` or `TestUserSkillReferenceClosureResolves` red; leaving a `Skill(skill=` occurrence causes the validation grep to exit `0` with a match.

**AC-5 - Unsupported harnesses have an executable host-neutral dispatch path and accurate ASD-STE100 `--version` documentation.**
Verified by: (1) `go test -race -count=1 ./internal/dispatch -run '^TestBuildPaseoHost'` passes with scrubbed native host markers (`CLAUDECODE=""`, `CODEX_THREAD_ID=""`, `PI_CODING_AGENT=""`, `PI_CODING_AGENT_DIR=""`, `PASEO_AGENT_ID=""`), proving an unsupported harness invoking `spacedock dispatch build --host paseo` gets a valid host-neutral schema-v2 envelope and dispatch file referencing `references/paseo-ensign-runtime.md` with zero Claude/Codex/Pi tokens; (2) live verification of a non-{Claude,Codex,Pi} Paseo-hosted FO (e.g. `antigravity` running `spacedock --version` with `PASEO_AGENT_ID` set and booting the workflow without aborting on `contract 3`, plus one-off unsupported-harness dry-run confirming 3-line `--version` acceptance and `--host paseo` build); (3) read-back of `docs/site/reference/command-reference.md` against actual `spacedock --version` output across bare terminal, Claude, and `PASEO_AGENT_ID`-only environments, plus ASD-STE100 sentence-length (<=25 words) and modal check.
Falsifying edit: requiring a native host marker inside `--host paseo` flips `TestBuildPaseoHostPromptShapeAndSplitRoot` red; mismatching `command-reference.md` against real `spacedock --version` output fails the read-back check.

## Test plan

- **Primary proof owners (reuse existing, no new standing enforcement):**
  1. `internal/cli/version_session_test.go` — owns AC-1 (`printVersion` matrix: Paseo-only sandboxed/unsandboxed, UUID non-leakage, native+Paseo coexistence, 3-line human terminal, zero host CLI exec).
  2. `internal/runtimehost/runtimehost_test.go` & `internal/dispatch/build_host_test.go` — own the invariant that `runtimehost.Detect` and `dispatch build` auto-detection are unaffected by `PASEO_AGENT_ID` (AC-1).
  3. `internal/dispatch/build_paseo_host_test.go` — owns host-neutral `--host paseo` envelope assembly under scrubbed host markers (AC-5).
  4. `internal/contractlint/*` (`fo_function_reference_invariant_test.go`, `boot_resident_closure_test.go`, `structural_checks_test.go`, `first_officer_eager_references_test.go`) — own structural closure and the 23,500 B shared-core cap (AC-2, AC-4).
  5. `internal/ensigncycle/paseo_live_smoke_test.go` & `paseo_adoption_live_test.go` — own live Paseo FO boot, skipped inner-host adapter verification, dispatch, and adoption (AC-2).
  6. `internal/ensigncycle/shared_live_runner_test.go` — owns native Claude/Codex/Pi non-regression (`TestLiveCommonShallowBoot`, `TestLiveCommonFullEnsignCycle`) because `skills/first-officer/references/first-officer-shared-core.md` and `fo-dispatch-core.md` are touched (AC-3, per `docs/dev/README.md:85`).
- **Additional check & distinct falsifying edit:** Add 4 subtests in `internal/cli/version_session_test.go` (inside existing `TestVersionSessionRender` + `TestVersionContractTokenPlacement`) covering `PASEO_AGENT_ID`-only sandboxed, `PASEO_AGENT_ID`-only unsandboxed, `PASEO_AGENT_ID` value redaction, and `CLAUDECODE`+`PASEO_AGENT_ID` coexistence; plus 1 subtest in `internal/dispatch/build_host_test.go` confirming `CLAUDECODE=1` + `PASEO_AGENT_ID=<uuid>` still auto-detects `claude` without ambiguity. Each fails under its own distinct mutation (removing fallback, leaking env value, or mistakenly adding `PASEO_AGENT_ID` to `runtimehost.Detect`).
- **Estimated cost & complexity:**
  - Deterministic Go unit + race suite (`go test ./...` and `go test ./... -race`): ~4–5 minutes, low complexity, fully offline.
  - Targeted live Paseo smoke + adoption (`TestLivePaseoDispatchSmoke`, `TestLivePaseoAdoptionSmoke`): ~5–8 minutes, requires local Paseo daemon + provider credentials (Commander / CI).
  - Required native host live non-regression (`claude-live`, `codex-live`, `pi-live` per `docs/dev/README.md:85` because shared core & `fo-dispatch-core.md` change): targeted `ShallowBoot` + `FullEnsignCycle` (~5–15 minutes per host locally, or CI matrix before merge).
  - One-off manual check: real `spacedock --version` read-back across 3 env shapes (bare terminal, `PASEO_AGENT_ID`-only, `CLAUDECODE`-set) + ASD-STE100 check on `docs/site/reference/command-reference.md` (<1 minute).
- **Detached adversarial audit:** Required before merge per `docs/dev/README.md:84` because this change touches shipped FO/ensign contract scaffolding (`skills/**/references/**`).

### Feedback Cycles

- **Cycle 1 (2026-11-23, detached `safehouse-codex` adversarial audit `b7f52561` + Cycle 1 validator `7320a70c`):** REJECTED on `284d19193` for 5 surgical fixes:
  1. `skills/ensign/SKILL.md:13` (Material): exclude `PI_CODING_AGENT` alongside `PI_CODING_AGENT_DIR` in the Paseo-first predicate so native Pi ensigns with `PI_CODING_AGENT=1` + ambient `PASEO_AGENT_ID` do not misroute onto `paseo-ensign-runtime.md`.
  2. `skills/first-officer/references/first-officer-shared-core.md:76` (Material): make keeping an addressable worker alive conditional on addressability (staying strictly `<= 23,496` B / `23,500` B cap) so fresh-worker Paseo gates satisfy the literal done condition.
  3. `internal/cli/version_session_test.go:335` (Deferred Risk -> fix now): scrub `CLAUDECODE`, `CODEX_THREAD_ID`, `PI_CODING_AGENT`, and `PI_CODING_AGENT_DIR` in `TestVersionExecutesNoHostCLI` when testing `PASEO_AGENT_ID` fallback so the Paseo-only branch is genuinely exercised.
  4. `internal/dispatch/build_paseo_host_test.go` (Deferred Risk -> fix now): add a marker-scrubbed subtest (all native host markers and `PASEO_AGENT_ID` empty, not `runNative`'s `CLAUDECODE=1`) proving AC-5 unsupported-harness `--host paseo` assembly.
  5. `skills/first-officer/SKILL.md:42` (Polish / drift guard): split the dense fallback paragraph into short sentences and refer to Startup step 1's required binary minor instead of hardcoding `0.29`.

## Stage Report: ideation

- DONE: Design internal/cli/cli.go printVersion fix (F4) so any Paseo-hosted harness (PASEO_AGENT_ID set) gets contract 3 without adding PASEO_AGENT_ID to runtimehost.Detect or changing human 3-line terminal output
  Specified exact 5-line `PASEO_AGENT_ID` fallback in `printVersion` (when `!ambiguous && host == ""`), emitting `Runtime: paseo (PASEO_AGENT_ID)`, `Sandbox: ...`, and `contract 3` with zero `runtimehost.Detect` changes, zero UUID leakage, and unchanged 3-line terminal output (`## Proposed approach` §1, `AC-1`).
- DONE: Audit every file loaded on the Paseo FO & ensign path (SKILL.md, inner-host *-first-officer-runtime.md adapters, first-officer-shared-core.md, fo-dispatch-core.md, paseo-fo-dispatch.md, ensign-shared-core.md, paseo-ensign-runtime.md) to eliminate unnecessary or contradictory loads, define Paseo-vs-native priority & sandboxed fallback (reachable daemon check), and check whether any runtime-specific constructs remain
  Audited all 10+ files atop PR #893 (`936`) and PR #895 (`x9m`), spiked `paseo status --json` live (proving stopped daemons also exit `0` with `connectedDaemon: "not_probed"`, so `"connectedDaemon": "reachable"` is required), skipped 5.1–8.5 KB inner-host FO adapters under active Paseo, fixed `fo-dispatch-core.md:188` `--host` contradiction and `Skill(skill=...)` residue while cutting `first-officer-shared-core.md` by 43 B (23,497 → 23,454 B <= 23,500 B cap), and locked first-match ensign precedence so sandboxed `paseo-access=workspace` native subagents never misroute (`AC-2`, `AC-3`, `AC-4`).
- DONE: Design the unsupported-harness fallback hint (guiding harnesses without a dedicated adapter to map Spacedock worker capabilities onto Claude/Codex-like native subagent/wait features) and define Expected surface, scanner-exact Acceptance criteria with falsifying Verified by clauses, and Test plan
  Authored exact before/after wording for `skills/first-officer/SKILL.md` and `skills/ensign/SKILL.md` (reusing `--host paseo` + `paseo-ensign-runtime.md` without minting `--host generic`, plus 3-line `--version` minor acceptance), concrete ASD-STE100 doc diff for `docs/site/reference/command-reference.md`, net `+90` LOC / 10-file expected surface, scanner-exact `## Acceptance criteria` (`AC-1`, `AC-2`, `AC-3`, `AC-4`, `AC-5` with VALUE ACs and falsifying edits), and primary-owner `## Test plan`.

### Summary

Completed ideation for `paseo-and-generic-harness-fo-support` atop PR #893 (`936`) and PR #895 (`x9m`). A live spike on `paseo status --json` revealed that a stopped daemon still exits `0` with `"connectedDaemon": "not_probed"`, making the `"connectedDaemon": "reachable"` check essential for safe sandboxed fallback; the operating-contract sweep saves 5.1–8.5 KB of contradictory inner-host boot prose under active Paseo, removes `Skill(skill=...)` and `--host` contradictions while trimming `first-officer-shared-core.md` 43 bytes further below its 23,500-byte cap, and reuses `--host paseo` for unsupported harnesses without adding unnecessary CLI mechanisms.


## Stage Report: implementation

- DONE: Implement the 5-line PASEO_AGENT_ID fallback in internal/cli/cli.go printVersion (F4, leaving runtimehost.Detect and 3-line human terminal output unchanged) with focused unit tests in internal/cli/version_session_test.go and internal/dispatch/build_host_test.go (AC-1)
  Commit `284d19193` adds the `PASEO_AGENT_ID` fallback in `internal/cli/cli.go:printVersion` without modifying `internal/runtimehost.Detect`. `TestVersionSessionRender` (`paseo-only-inside-safehouse`, `paseo-only-unsandboxed`, `paseo-marker-value-never-rendered`, `claude-inside-paseo-no-ambiguity`, `outside-every-runtime-is-three-lines`), `TestVersionContractTokenPlacement`, and `TestVersionExecutesNoHostCLI` assert 6-line Paseo output with `contract 3`, UUID non-leakage, 3-line bare terminal output, and zero `paseo`/host CLI exec (removing the fallback, leaking the UUID, or execing `paseo` flips them red); `TestResolveBuildHostFromMarkers` (`claude-inside-paseo-no-ambiguity`, `paseo-agent-id-only-requires-explicit-host`) asserts `runtimehost.Detect` ignores `PASEO_AGENT_ID` (adding `PASEO_AGENT_ID` to `runtimehost.Detect` flips both red).
- DONE: Update skills/first-officer/SKILL.md (priority-ordered Paseo reachability probe requiring exit 0 and connectedDaemon == reachable, skipping inner-host FO adapters when reachable, clean sandboxed fallback, and unsupported-harness fallback reusing --host paseo), skills/ensign/SKILL.md (first-match precedence + generic fallback), and skills/ensign/references/ensign-shared-core.md (2-line Paseo/host-neutral DISPATCH_FILE bootstrap example) (AC-2, AC-3, AC-5)
  Commit `284d19193` applies exact approved wording across `skills/first-officer/SKILL.md`, `skills/ensign/SKILL.md`, and `skills/ensign/references/ensign-shared-core.md`. `internal/dispatch` (`TestBuildPaseoHost*`) passes under scrubbed host markers confirming host-neutral `--host paseo` assembly (requiring a native marker inside `--host paseo` flips it red); live Paseo/native-host suites (`TestLivePaseo*`, `TestLiveCommon*`) could not run inside this Safehouse ensign sandbox (`F8`: `paseo status --json` exited 1 with `DAEMON_STATE_READ_FAILED` / `EPERM` on `paseo.pid`), which also empirically confirms why sandboxed sessions fall through to step 2.
- DONE: Clean up Skill(skill=...) and --host contradictions in skills/first-officer/references/first-officer-shared-core.md (shrinking strictly below 23,497 B / 23,500 B cap), skills/first-officer/references/fo-dispatch-core.md (L61, L188), skills/first-officer/references/paseo-fo-dispatch.md (L3), skills/feedback-rejection-flow/SKILL.md (L27), and update docs/site/reference/command-reference.md per ASD-STE100 (AC-2, AC-4, AC-5)
  Commit `284d19193` removes all 6 `Skill(skill=...)` wrappers and aligns `fo-dispatch-core.md:188` and `paseo-fo-dispatch.md:3`; `wc -c skills/first-officer/references/first-officer-shared-core.md` measures `23446` B (-51 B vs 23,497 B baseline, 54 B below the 23,500 B cap), one-off `grep -n 'Skill(skill='` across the 3 edited contract files returned empty, `internal/contractlint` (`TestFOInstructionComponentCaps`, `TestBootResidentDeferredLoadPointsResolve`, `TestUserSkillReferenceClosureResolves`, `TestFirstOfficerEntryEagerlyImportsOnlySharedCore`, `TestPaseoRuntimeAdaptersResolveOnDisk`) passed (exceeding 23,500 B or breaking any `references/`/`spacedock:` target flips them red), and `docs/site/reference/command-reference.md` L7 matches ASD-STE100 (10 and 18 words, active present tense) and all 3 real `--version` shapes.
- DONE: Run gofmt, go vet, go test ./... and go test ./... -race, internal/contractlint, real spacedock --version read-back across 3 env shapes, and spacedock status --validate; commit cleanly in the worktree and append ## Stage Report: implementation to docs/dev/.spacedock-state/paseo-and-generic-harness-fo-support/index.md
  `gofmt -l` clean on all changed Go files, `go vet ./...` clean, full `go test ./...` and `go test ./... -race` green across all 22 test packages (including `internal/contractlint`, `internal/cli`, `internal/runtimehost`, `internal/dispatch`), real built binary `--version` read-back verified 3 lines (bare terminal) and 6 lines (`PASEO_AGENT_ID`-only -> `Runtime: paseo (PASEO_AGENT_ID)` + `contract 3`; `CLAUDECODE=1` and `CLAUDECODE=1 PASEO_AGENT_ID=...` -> `Runtime: claude (CLAUDECODE)` + `contract 3`), and `spacedock status --validate --workflow-dir docs/dev` returned `VALID`. Code committed at `284d19193` (+124/-30, net +94 LOC across 11 files, within +40..+160 LOC / 8..12 file tolerance).

### Summary

Implemented Paseo-wide FO boot support (`F4` `PASEO_AGENT_ID` fallback in `printVersion`), priority-ordered Paseo reachability probing with inner-host adapter skip and sandboxed/unsupported-harness fallback, host-neutral shared-core cleanup (`first-officer-shared-core.md` reduced from 23,497 B to 23,446 B), and ASD-STE100 `--version` doc update in commit `284d19193` (11 files, +94 net LOC). All deterministic unit, race, `contractlint`, real binary `--version` read-back, and `status --validate` checks passed; live Paseo/native-host suites require daemon/provider access outside this Safehouse ensign sandbox (`F8`, `paseo status --json` returned `EPERM` on `paseo.pid`) and remain for Commander/CI validation.


## Stage Report: validation

- DONE: Audit commit 284d19193 against xdh/index.md (AC-1 through AC-5, Expected surface +40..+160 LOC / 8..12 files, first-officer-shared-core.md <= 23,496 B / 23,500 B cap, zero Skill(skill=...) residue, zero runtimehost.Detect changes, zero edits to x9m files)
  Commit `284d19193` vs `3d134b255` (`main`) is `+124/-30` (net `+94` LOC, `11` files; within `+40..+160` LOC / `8..12` files tolerance); `first-officer-shared-core.md` shrunk from `23,497` B to `23,446` B (`-51` B, `<= 23,496` B and `54` B below the `23,500` B cap); `grep -n 'Skill(skill='` across `first-officer-shared-core.md`, `fo-dispatch-core.md`, and `feedback-rejection-flow/SKILL.md` returned empty; `git diff 3d134b255..284d19193 -- internal/runtimehost` is empty (zero `runtimehost.Detect` changes); zero edits to `x9m`-only files (`paseo_adoption_live_test.go`, `docs/runtime-live-ci.md`) or `paseo-fo-dispatch.md` verbs/matrix/bindings (only L3 intro aligned per §4.D).
- DONE: Run gofmt -l, go vet ./..., go test -race -count=1 ./internal/cli ./internal/runtimehost ./internal/dispatch ./internal/contractlint, full go test -count=1 ./... and go test -race -count=1 ./..., real built spacedock --version 3-shape read-back, ASD-STE100 check on docs/site/reference/command-reference.md, and spacedock status --validate
  `gofmt -l` clean on all 3 changed Go files (pre-existing untouched `internal/release/runtime_live_evidence_workflow_test.go` only in repo-wide scan); `go vet ./...` clean; targeted `-race -count=1` green (`internal/runtimehost` 1.25s, `internal/dispatch` 75.14s incl. scrubbed-env `TestBuildPaseoHost*` 2.37s, `internal/contractlint` 5.57s, `internal/cli` 221.52s with `PATH=/opt/homebrew/bin:/usr/bin:/bin` + `TestVersion*` 6.28s); full `go test -count=1` and `go test -race -count=1` green across all 22 test packages (`internal/cli` with clean `PATH`, all other packages default `PATH`); real built temp binary `--version` read-back confirmed exact 3-line bare/empty-`PASEO_AGENT_ID` output and 6-line `PASEO_AGENT_ID`-only (`Runtime: paseo (PASEO_AGENT_ID)`, `contract 3`, zero UUID leakage), native+`PASEO_AGENT_ID` (`claude`/`codex`/`pi` win without ambiguity), and ambiguous-native+`PASEO_AGENT_ID` (`Runtime: ambiguous (CODEX_THREAD_ID, CLAUDECODE)` exit 0) shapes; `docs/site/reference/command-reference.md:7` passes ASD-STE100 (descriptive, 8 and 16 words <= 25, simple present, active voice, zero banned modals/semicolons, matches real `--version` output); `spacedock status --validate --workflow-dir docs/dev` returned `VALID`; AC-3 spike reproduced (`paseo --home <empty> status --json` exited `0` with `localDaemon: "stopped"` and `connectedDaemon: "not_probed"`).
- DONE: Record ## Stage Report: validation (Recommendation: PASSED or REJECTED) with per-AC evidence, falsifying edits, Review-finding disposition (Material / Deferred / Polish), and named F8 sandbox boundary for live Paseo/native-host suites
  **Recommendation: PASSED** (deterministic proof complete, zero Material findings; merge remains gated on Shaping FO's concurrent live `TestLivePaseoDispatchSmoke` / non-{Claude,Codex,Pi} check, the parallel detached adversarial audit, and CI native live lanes per `docs/dev/README.md:85`). **Per-AC & falsifying edits:** **AC-1 (VALUE)** proven by `TestVersionSessionRender`, `TestVersionContractTokenPlacement`, `TestVersionExecutesNoHostCLI`, `TestResolveBuildHostFromMarkers`, and real 9-shape `--version` binary read-back (removing `PASEO_AGENT_ID` fallback flips `paseo-only-*` red; moving `PASEO_AGENT_ID` into `runtimehost.Detect` flips `claude-inside-paseo-no-ambiguity` & `TestResolveBuildHostFromMarkers` red; rendering the UUID flips `paseo-marker-value-never-rendered` red); **AC-2 (VALUE)** deterministic half proven (`first-officer-shared-core.md` = `23,446` B + `TestFOInstructionComponentCaps` PASS; adding >54 B flips it red), live half (`TestLivePaseo(Dispatch|Adoption)Smoke` + zero inner-host adapter reads in FO trace) blocked by `F8` ensign sandbox (`paseo status --json` exited `1` with `DAEMON_STATE_READ_FAILED` / `EPERM` on `paseo.pid`) and run concurrently outside sandbox by Shaping FO; **AC-3** deterministic half proven (`paseo --home <empty> status --json` exited `0` with `connectedDaemon: "not_probed"` + `skills/ensign/SKILL.md` first-match precedence), native live non-regression (`TestLiveCommon(ShallowBoot|FullEnsignCycle)` on `claude`/`codex`/`pi`) blocked by `F8` sandbox / provider credentials and required on CI before merge; **AC-4** proven by `internal/contractlint` (`TestBootResidentDeferredLoadPointsResolve`, `TestUserSkillReferenceClosureResolves`, `TestFirstOfficerEntryEagerlyImportsOnlySharedCore`, `TestPaseoRuntimeAdaptersResolveOnDisk`, `TestFOInstructionComponentCaps` PASS; breaking any `references/`/`spacedock:` target flips them red) + empty `grep -n 'Skill(skill='` across the 3 edited contract files; **AC-5** deterministic + doc halves proven by scrubbed-env `TestBuildPaseoHost*` (requiring a native marker flips `TestBuildPaseoHostPromptShapeAndSplitRoot` red) + real `--version` read-back + ASD-STE100 check on `command-reference.md:7`, live non-{Claude,Codex,Pi} FO boot check running concurrently under Shaping FO. **Review-finding disposition:** *Material:* none. *Deferred risk (1):* `skills/first-officer/SKILL.md:42` hardcodes `binary minor (0.29)` instead of referencing the required binary minor from `first-officer-shared-core.md:9` (trigger: future `v0.30.0` bump via `spacedock-release`, since `internal/release/prose_stamp.go:12-16` rewrites only `first-officer-shared-core.md`'s single pinned literal; current `0.29` AC-5 passes and `xdh/index.md:167` approved exact wording; promote-to-material condition: cutting `0.30`). *Polish (2):* (a) `internal/cli/version_session_test.go:334-335` (`TestVersionExecutesNoHostCLI`) sets `PASEO_AGENT_ID` without scrubbing ambient `CLAUDECODE`/`CODEX_THREAD_ID`/`PI_CODING_AGENT*` (runs the Paseo branch in CI/bare shells, but nested Claude/Codex/Pi dev shells match the outer host first; `printVersion` has zero exec calls); (b) `internal/cli/version_session_test.go:153` test comment uses a semicolon (minor ASD-STE100 Rule 8.1 nit in non-user-facing test prose).

### Summary

Read-only validation of commit `284d19193` recommends **PASSED** with zero Material findings (code worktree left untouched and clean). All deterministic unit, `-race`, `contractlint`, real binary `--version` read-back, ASD-STE100, `status --validate`, and stopped-daemon spike checks passed across all 22 test packages; one Deferred risk (`SKILL.md:42` hardcoded `0.29` on future `0.30` minor bump) and two Polish items are recorded above. Live Paseo and native-host suites (`TestLivePaseo*`, `TestLiveCommon*`) cannot execute inside this Safehouse ensign sandbox (`F8`: `paseo status --json` exited `1` with `EPERM` on `paseo.pid`) and remain gated on Shaping FO's concurrent live run, the parallel detached adversarial audit, and required native live CI lanes (`docs/dev/README.md:85`).
