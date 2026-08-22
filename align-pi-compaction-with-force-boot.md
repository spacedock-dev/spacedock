---
title: "Align the Pi extension's compaction hook with force-boot (re-read state, not re-inject contract)"
status: ideation
source: "Captain (2026-08-21): the Pi extension's session_compact hook re-injects FO_BOOTSTRAP_TEXT (a contract pointer), but PR #738 (force-boot-at-compaction-boundary, merged) established the opposite mechanism — re-read durable state via one «state.boot»(), do NOT re-inject the contract. The Pi extension did not follow #738."
started: 2026-08-22T00:41:38Z
completed:
verdict:
score: 0.8
worktree:
issue:
id: h9nn5brc1dp0m82x5en21d56
gates:
    version: 1
    records:
        - id: gate:h9nn5brc1dp0m82x5en21d56:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:h9nn5brc1dp0m82x5en21d56-backlog-1
              briefing:
                id: briefing:h9nn5brc1dp0m82x5en21d56:backlog:attempt-1:revision-1
                digest: sha256:7192343da7f4ba50ac81545ad7b13d6164748bb92cb20b68fa21cd1b74602ad2
                request-digest: sha256:e0f99158010258233aa1c4ccae34bb9b70ea5d506ba6a10e722397693488f7d6
                room-ref: ./align-pi-compaction-with-force-boot/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:h9nn5brc1dp0m82x5en21d56:backlog:1
                briefing: briefing:h9nn5brc1dp0m82x5en21d56:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-22T00:34:54.852202173Z"
                decision: approve
                reason: 'Conn-held. Seed clearly identifies the #738 misalignment: Pi re-injects contract text, #738 says re-read state. Advance to ideation for the alignment fix, stacking on 753.'
              application:
                target-stage: ideation
                state: consumed
---

The Pi extension's `session_compact` hook (`.pi/extensions/spacedock.ts`) re-injects `FO_BOOTSTRAP_TEXT` at the compaction boundary — a contract pointer telling the FO to re-satisfy load preconditions and re-read durable state. PR #738 (`force-boot-at-compaction-boundary`, merged) established the opposite for Claude/Codex: fire one `«state.boot»()` (re-read durable state); the contract does not need re-injecting. The Pi extension is misaligned — it does the thing #738 rejected (re-inject the contract) and only points at re-reading state rather than doing it.

Direction: change the Pi `session_compact` hook to fire a `«state.boot»()` (the `spacedock status --boot --identify --json` read) at the compaction boundary, aligning with #738's "re-read durable state, don't re-inject contract." The boot record carries mods/ready_gates/dispatchable/pr_state/team_state/state_backend — the same answers #738 provides Claude/Codex. Keep the bootstrap text re-injection only if a follow-on load actually requires the contract pointer (it likely does not — #738's conclusion is the contract survives compaction; only the state goes stale).

Acceptance sketch: value — a compacted Pi FO re-reads durable state on resume (the boot record), matching Claude/Codex; the bootstrap-text re-injection is removed or reduced to what #738's mechanism doesn't cover. mechanism — a behavior test asserting the compaction boundary triggers a boot read, not a contract re-injection. Expected surface: `.pi/extensions/spacedock.ts` + test; small. Stacks on PR 753 (same Pi-extension file).

## Stage Report

### Problem analysis

The Pi extension (`.pi/extensions/spacedock.ts`) uses a single `injectBootstrap` flag set by both `session_start` and `session_compact`. The `context` handler injects `FO_BOOTSTRAP_TEXT` — a contract pointer ("Load the $spacedock:first-officer skill … re-satisfy every load precondition") plus a state directive ("re-read durable state before the next workflow effect"). After compaction the system prompt is rebuilt from the skill (discovered via `resources_discover`), so the contract is already present. Re-injecting the contract pointer is the exact mechanism #738 rejected.

PR #738's Claude/Codex implementation (`hooks/session_start_compact_reminder.sh`, registered via `hooks.json` as a `SessionStart(compact|clear)` hook) does NOT re-inject the contract. It outputs a reminder: "your bindings are stale … take a fresh `spacedock status --boot --identify --json` … resume the loop where it stopped." The FO then runs `«state.boot»()` as part of its contract (the shared core says "Invoke `«state.boot»()` once and retain its boot record"). The Pi extension's `FO_BOOTSTRAP_TEXT` already says "re-read durable state" — the problem is it ALSO re-injects the contract, which #738 says is unnecessary, and it doesn't actually run the read.

### Selected approach

**Replace the compaction-path injection with an extension-run boot read.** Add a separate `injectBootRecord` flag (set by `session_compact` only). When the `context` handler sees `injectBootRecord`, it runs `spacedock status --boot --identify --json` via `pi.exec()` (the Pi ExtensionAPI's built-in async subprocess executor) and injects the boot record as a user message with a minimal framing directive. The `FO_BOOTSTRAP_TEXT` contract re-injection is fully removed from the compaction path — #738 says the contract survives compaction (it's in the system prompt via `resources_discover` + the skill); only state goes stale.

The injected message format:
```
[SPACEDOCK-FO-BOOT-v2] Durable state boot record — re-read before the next workflow effect. The compacted summary is not authoritative. Resume the loop where it stopped; do NOT greet or re-present a session summary. Pi tool mapping: read/write/edit/bash/grep/find/ls; load skills via read; subagent via pi-subagents when available.

{boot-record JSON}
```

This mirrors #738's reminder language ("bindings are stale, not just narrative"; "resume the loop") while adding the actual boot record — going one step further than Claude/Codex, which only remind the FO to run the read. The Pi extension can do this because `pi.exec()` is available on the ExtensionAPI.

**`session_start` path:** unchanged. The initial boot legitimately needs the contract pointer (the FO is starting fresh, the skill is in the system prompt but the FO needs the "you are the first officer" framing + Pi tool mapping). `FO_BOOTSTRAP_TEXT` is retained for `session_start` only. The two paths use separate flags (`injectBootstrap` for `session_start`, `injectBootRecord` for `session_compact`), both cleared on `agent_end`.

**Bootstrap text re-injection: fully removed on the compaction path.** No minimal pointer is kept — the boot record's framing directive carries the "resume the loop" + "compacted summary is not authoritative" language, and the contract is in the system prompt. A minimal contract pointer would be the thing #738 rejected.

### Simplest alternative considered and why insufficient

**Alternative:** keep `FO_BOOTSTRAP_TEXT` on the compaction path, just strip the contract sentences ("Load the skill", "re-satisfy load preconditions") and keep the state directive ("re-read durable state"). This is the minimal text edit.

**Why insufficient:** this keeps the reminder-only pattern — the FO is told to re-read state but the extension doesn't actually do it. The current `FO_BOOTSTRAP_TEXT` already says "re-read durable state before the next workflow effect" and that hasn't prevented stale-state resume failures (the thing #738 was filed to fix). The entity's direction is to *fire* `«state.boot»()`, not to improve the reminder. Just editing the text doesn't change the mechanism; it only removes the contract re-injection without adding the state read.

### Riskiest-mechanism spike

**No spike needed: proven mechanisms.**

The Pi ExtensionAPI provides `exec(command: string, args: string[], options?: ExecOptions): Promise<ExecResult>` (confirmed in `@earendil-works/pi-coding-agent/dist/core/extensions/types.d.ts`). `ExecResult` has `{ stdout: string; stderr: string; code: number; killed: boolean }`. The `context` handler type is `ExtensionHandler<ContextEvent, ContextEventResult>` = `(event, ctx) => Promise<R | void> | R | void` — it can be async. The Pi runtime awaits async handlers before the LLM call. `pi` is captured in the factory closure, so `pi.exec()` is available inside `context` handler bodies.

The boot read (`spacedock status --boot --identify --json`) is a fast local binary call (~100ms). No event-loop blocking concern for a one-shot post-compaction read. The `--identify` flag adds `discovery`, `stages`, and `ready_gates` to the boot record — the same fields #738's reminder tells the FO to read.

### Acceptance criteria

**AC-1 (value-measuring):** A compacted Pi FO resumes on re-read durable state (the boot record injected at the compaction boundary), not a contract re-injection. Measured against the current baseline: the compaction path re-injects `FO_BOOTSTRAP_TEXT` (contract pointer) without reading state — a baseline that can move the wrong way (the FO re-satisfies load preconditions from a stale compacted summary instead of re-reading actual state). Test: fire `session_compact`, then `context`; assert the injected message contains boot-record fields (`"command":"boot"`, `"mods"`, `"ready_gates"`) and does NOT contain `FO_BOOTSTRAP_MARKER` ("SPACEDOCK-FO-BOOTSTRAP-v1") or "Load the $spacedock:first-officer skill".

**AC-2 (mechanism):** The `session_compact` → `context` path fires a boot read via `pi.exec("spacedock", ["status", "--boot", "--identify", "--json"])` and injects the parsed JSON as a user message. Test: mock `pi.exec` to return a fixed boot JSON; fire `session_compact`; `await` the `context` handler; assert `pi.exec` was called with the expected args and the injected message contains the mock boot record's `"command":"boot"` field.

**AC-3 (scope boundary):** The `session_start` path is unchanged — it still injects `FO_BOOTSTRAP_TEXT` (contract pointer + state directive + Pi tool mapping). Test: fire `session_start`; fire `context`; assert the injected message contains `FO_BOOTSTRAP_MARKER` and "Load the $spacedock:first-officer skill" (unchanged behavior).

**AC-4 (child exemption):** A pi-subagents child session (`PI_SUBAGENT_CHILD=1`) injects zero bootstrap on both `session_start` and `session_compact`. Test: existing `TestSpacedockPiExtensionChildExemption` covers `session_compact`; verify it still passes with the new `injectBootRecord` flag (child sessions skip both flags).

**AC-5 (deduplication):** If the boot record is already present in the context messages, the `context` handler does not inject a second copy. Test: fire `session_compact`; fire `context` (injects boot record); fire `context` again with the same messages; assert the second call returns `undefined` (no duplicate injection).

### Test plan

Update `internal/piruntime/spacedock_extension_test.go` — the existing JS harness:

1. **Mock `pi.exec`:** The harness's `pi` object adds `exec(command, args, options)` returning a `Promise<ExecResult>` with a fixed boot JSON in `stdout` (e.g., `{"command":"boot","mods":{},"ready_gates":[],"state_backend":"single-root"}`).
2. **Async `context` handler:** The harness `await`s the `context` handler's return value (it's now async on the boot-record path). The `session_start` path can remain sync (no `pi.exec` call).
3. **Compaction-path assertions:** Fire `session_compact`; `await` `context`; assert the injected message contains `"command":"boot"` and does NOT contain `FO_BOOTSTRAP_MARKER` or "Load the $spacedock:first-officer skill".
4. **Dedup assertion:** Fire `context` again with the boot-record message already present; assert `undefined`.
5. **`session_start`-path regression:** Fire `session_start`; fire `context`; assert `FO_BOOTSTRAP_TEXT` is injected (unchanged).
6. **Child exemption regression:** Existing child-exemption test passes unchanged.

Cost: fixture test (no live run needed). The test runs the real `.pi/extensions/spacedock.ts` through a Node.js harness — same pattern as the existing test, extended with `pi.exec` mock + async `await`.

### Expected surface

**Files:** `.pi/extensions/spacedock.ts` (implementation), `internal/piruntime/spacedock_extension_test.go` (test harness update).

**Net LOC change:** +15–20, across 2 files.
- `.pi/extensions/spacedock.ts`: +~20 (add `injectBootRecord` flag, `hasBootRecord` dedup check, boot-record injection in `context` handler with `pi.exec` call + framing directive), -~0 (`FO_BOOTSTRAP_TEXT` constant retained for `session_start`). Net: +20.
- `internal/piruntime/spacedock_extension_test.go`: +~15 (mock `pi.exec`, async `await`, boot-record assertions), -~5 (remove compaction-path `FO_BOOTSTRAP_TEXT` assertions, replaced by boot-record assertions). Net: +10.

**Observable-semantics declaration:** No change to the FO bootstrap content itself (the first-officer skill `SKILL.md` is untouched). No change to gate/dispatch/state mechanics. No change to `resources_discover`, `session_start` injection, or `agent_end` flag-clearing. No change to command grammar, stored formats, or authority. The only observable semantic change: **the compaction boundary injects a boot record (state) instead of contract text** — a compacted Pi FO receives `spacedock status --boot --identify --json` output as a context message, not `FO_BOOTSTRAP_TEXT`. This changes post-compaction context content (what the FO sees on resume), not any command surface or persistent format.

### Summary

Diagnosed the `session_compact` misalignment: the Pi extension re-injects `FO_BOOTSTRAP_TEXT` (contract pointer) at the compaction boundary, but #738's mechanism (implemented for Claude/Codex via `hooks/session_start_compact_reminder.sh`) re-reads durable state and does NOT re-inject the contract. Proposed: add a separate `injectBootRecord` flag for the compaction path; the `context` handler runs `spacedock status --boot --identify --json` via `pi.exec()` and injects the boot record with a minimal framing directive; remove `FO_BOOTSTRAP_TEXT` from the compaction path entirely (retain for `session_start`). No spike needed — `pi.exec()` is a proven ExtensionAPI method and the `context` handler can be async. 5 ACs (value-measuring + mechanism + scope boundary + child exemption + dedup). Test plan updates the existing Node.js harness with `pi.exec` mock + async `await`.
