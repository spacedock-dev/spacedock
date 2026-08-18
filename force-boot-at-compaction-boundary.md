---
id: 4ctz4sfybfk0mfsbnrjcc7bv
title: Force one boot at the compaction boundary — a compacted FO resumes on stale bindings
status: ideation
source: "Captain CL, 2026-08-18, in chat: 'file the compaction improvement, detail the problem diagnosis and proposed solution.' Raised after the FO opened three PRs by hand without reading the pr-merge mod, and diagnosed the root cause as never having run Startup in a compaction-resumed session."
started: 2026-08-18T23:15:43Z
completed:
verdict:
score:
worktree:
issue:
gates:
    version: 1
    records:
        - id: gate:4ctz4sfybfk0mfsbnrjcc7bv:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:4ctz4sfybfk0mfsbnrjcc7bv-backlog-1
              briefing:
                id: briefing:4ctz4sfybfk0mfsbnrjcc7bv:backlog:attempt-1:revision-1
                digest: sha256:8c097b36dc8addb17335f90ea018c35bb0f1b083f31bce9d390c80df9e543608
                request-digest: sha256:79ea2c142d7fb3cbbd6ec4afb6b981efee79195296b7750c8b3b4166550a49f9
                room-ref: ./force-boot-at-compaction-boundary/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:4ctz4sfybfk0mfsbnrjcc7bv:backlog:1
                briefing: briefing:4ctz4sfybfk0mfsbnrjcc7bv:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-18T23:14:57.630325Z"
                decision: approve
                reason: 'Captain approved in chat: ''ok dispatch 4c and nr.'' Accepts the seed into ideation; the captain independently identified this as the highest-value item on the board because it makes every future session more reliable rather than fixing one test.'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:4ctz4sfybfk0mfsbnrjcc7bv:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:4ctz4sfybfk0mfsbnrjcc7bv-ideation-1
              briefing:
                id: briefing:4ctz4sfybfk0mfsbnrjcc7bv:ideation:attempt-1:revision-1
                digest: sha256:dde54526427e9d671d87cc588f41ea34487238b0f323dd1282f2a9a820f93afd
                request-digest: sha256:32bf3d22c3138621ab5a5deaf78881ec0a961edaa88d8f5ad9cc3f03bad99a68
                room-ref: ./force-boot-at-compaction-boundary/review/ideation/briefing-1
---

A compacted First Officer session inherits the previous session's **narrative** but not its **bindings**. The transcript summary preserves what the FO was discussing and drops what it was standing on: which binary, which contract version, which mods are registered, which workers are alive, and what durable state actually says.

## Problem — diagnosis

The FO contract's Startup runs a binary gate, one `«state.boot»()`, and `«engage»`, which runs `state ready`, `state sweep`, and `«hooks.run»("startup")` exactly once. A compaction-resumed session does not know it is a new session. It believes it is a continuation, so none of that runs.

Six failures in one session (2026-08-18), all traceable to that single omission:

1. **Binary gate never ran.** The first command failed against `/opt/homebrew/Caskroom/spacedock@next/0.27.0-pre7/spacedock`, deleted by the pre8 cask upgrade. Startup step 1 resolves the launcher before anything else.
2. **The `mods` map was never read**, because no boot record was obtained. The FO then opened three GitHub PRs by hand, bypassing the `pr-merge` mod completely: no gate recorded, no `pr:` field, no `mod-block`, no merge guard armed. Three PRs existed that the workflow could not see. Repaired only when the captain asked the FO to check the mod.
3. **`state ready` / `state sweep` never ran**, so the session never converged with the remote state checkout.
4. **`«hooks.run»("startup")` never fired**, so `pr-merge`'s PR-pending scan and `comm-officer` injection never happened.
5. **The session ran the `0.27.0-pre1` contract while `pre8` was installed**, so `fo-gate-lifecycle` was absent from its loaded set. Four gates were "presented" in chat and never `gate prepare`d. Durable readiness stayed `needs-preparation` for hours while the FO reported them to the captain as awaiting decision.
6. **An orphaned worker kept running** — the ideation ensign for a withdrawn entity. A roster reconcile at startup would have reclaimed it.

Item 5 is the sharpest: the FO's prose and durable state disagreed, and only the binary knew. That is the same defect class as two entities already in flight — `conn-delegated-approver-attribution` (a grader that trusted the agent's self-description) and `merge-guard-requires-preceding-report` (an FO that trusted a completion message over durable state) — one level up. Here the FO trusts **its own summary** over the binary.

## Why a prose rule will not fix this

The gate-attribution convention has lived in `skills/fo-gate-lifecycle/SKILL.md` since `aa04e95d8` (2026-07-25). Live First Officers recited it in their own final messages and still stamped `person:captain` — recorded in this repo's finding-9 audit note. Prose in the contract demonstrably does not hold against a behavior the session is already inclined toward.

Worse, the compaction summary carries an explicit counter-instruction: *"Resume directly — do not acknowledge the summary, do not recap, do not preface."* That is correct for conversational continuity and wrong for a stateful dispatcher. Any mechanism here must survive an instruction telling the session to keep going.

## Proposed approach

Bind exactly **one `«state.boot»()`** to the compaction boundary, before any other action. Nothing else.

The boot record already carries `mods`, `ready_gates`, `dispatchable`, `pr_state`, `team_state`, `state_backend`, `definition_dir`/`entity_dir`, and the binary version gate. Every one of the six failures above is answered by that single call. The contract does not need re-injecting; the state needs re-reading.

### Relationship to issue #595

#595 is this same problem filed for Codex, and it reaches the same diagnosis in its own words: a teammate must retain contract and authority "without trusting a lossy transcript summary." It records that the Codex surface exposes no compaction callback, and proposes a digest-pinned bundle of contract, adapter, README, entity state, worker identity, gates, PR, and CI.

This entity is the Claude-side instance and deliberately proposes far less. If Claude exposes a usable boundary, one boot call ships without any bundle. #595's bundle remains the right shape only where no callback exists at all.

## Ideation must settle

1. **Whether Claude Code exposes a usable compaction boundary, proven by a captured event rather than documentation.** This repo's own `hook-events.jsonl` contains exactly one `PreToolUse` record and no compaction evidence, so nothing local proves it today. Prove the boundary first; design second.
2. **Whether a hook can force an action before the next tool call, or only inject text.** Injected text is prose and fails for the reason above. If text is all the host offers, say so and record the boundary as unsupported — do not ship a prose rule dressed as a mechanism.
3. **Where the mechanism belongs** — a new mod hook point, the runtime adapter, or the binary. Prefer the smallest that holds.
4. **Whether a host-independent fail-closed tell can substitute.** If the FO must be able to prove it booted — for example by holding a boot record whose identity matches the current session — then a session that cannot prove it booted refuses gate and merge actions until it does. This needs no host callback and would cover Codex and Pi as well. Weigh it against the hook approach on necessity, not preference.

## Out of scope

The Codex boundary — #595 owns it. Any context-reinjection bundle. Post-compaction worker-identity recovery. Changing what Startup does; this entity changes only whether it is forced to run.

## Value

The First Officer mutates durable state, resolves gates, and drives merges. A compacted session doing that work on stale bindings is the highest-authority actor in the system operating on a summary it wrote about itself. The value is that the FO cannot act on workflow state it has not re-read.

## Ideation findings (2026-08-18, Claude Code v2.1.226)

The four questions the seed left open are settled, each on captured evidence — verbatim records in `force-boot-at-compaction-boundary/ideation-spike-evidence.md` beside this file.

**1. The boundary exists, captured two ways.** A controlled live session (tmux-driven scratch session, capture hooks registered via `--settings`, `/compact` issued after five turns) produced the full event sequence: `PreCompact{trigger:"manual"}` and `SessionStart{source:"compact"}` hook events, plus a durable `{"type":"system","subtype":"compact_boundary",…,"compactMetadata":{trigger, preTokens, postTokens,…}}` record in the session transcript. The incident session's own transcript carries the same boundary record (801,920 → 19,650 tokens, 2026-08-18T18:59:30Z, `trigger:"manual"`). The boundary is real, host-observable, and — decisively — durable on disk.

**2. Hooks cannot force an action, and the shipped prose injection demonstrably did not hold.** The hook surfaces at this boundary are: inject text into the next context (SessionStart `additionalContext`) or run host-side code (any command hook). Nothing forces the model's next tool call. The plugin has shipped a SessionStart `^compact$` prose-injection hook (`codex_session_start_compact.sh`) since before the incident — registered in the 0.27.0-pre1 hooks.json the incident session ran. At the incident boundary it produced no observable effect (zero injection text in the transcript, while 216 `stop_hook_summary` records prove hooks were active in that session), and the six failures followed. Whether it failed to fire or fired and was ignored, prose injection at this boundary is empirically insufficient — as the seed predicted.

**3. Session identity does NOT change at compaction — the naive fail-closed tell is dead.** `sessionId` is identical across the boundary in both the incident record and the controlled capture, and the transcript file continues. A guard of the shape "boot record must match the current session id" cannot detect compaction. (It still catches `--resume` and never-booted fresh sessions, which get new ids.)

**4. The tell survives in a stronger form, because the boundary is durable state.** Compaction leaves a timestamped `compact_boundary` record in the session's own transcript, and the binary can read it at verb time. `CLAUDE_CODE_SESSION_ID` is set per-session by the harness, visible to model-side Bash, and survives the boundary (evidence file §4). So the binary can answer "has this session compacted since it last booted?" from durable state alone — no hook callback, no model cooperation.

One more captured nuance: PreCompact fired at a `/compact` the host then refused ("Not enough messages to compact"). A PreCompact-based detector would mark sessions stale for compactions that never happened. The transcript record has no such false positive — it exists only when compaction completed.

## Design (settled): boot receipt + fail-closed guard, no hook

The seed said "bind exactly one `«state.boot»()` to the compaction boundary." Ideation refines where the binding lives: refusal-side, not callback-side. The first workflow effect after the boundary cannot proceed until a fresh boot runs — enforced by the binary, the one component that knew the truth in the incident.

**Boot receipt.** `spacedock status --boot --identify --json` (the shipped `«state.boot»()`) additionally writes a session boot receipt to `.spacedock/boot/{session_id}.json`: `{session_id, transcript_path, booted_at, binary_version}`. Written only when session identity is resolvable (env `CLAUDE_CODE_SESSION_ID` present); silent no-op otherwise. The receipt is host scratch, not workflow state — the state checkout is untouched, and boot's "mutates nothing" contract still holds for workflow/entity state (contract wording gains a clarifying line; doc diff below).

**Guard.** FO-held mutating verbs — `gate prepare|record|consume|withdraw`, `merge guard`, `dispatch build`, `new`, `status --set` — run a shared preflight:

- No receipt for the current session id → refuse.
- Latest `compact_boundary` timestamp in the receipt's recorded transcript > `booted_at` → refuse.
- Refusal: distinct exit code (4 proposed; `state commit` owns 3 — implementation confirms 4 is unclaimed) and one stderr paragraph naming the condition and the exact remedy: re-run `${SPACEDOCK_BIN:-spacedock} status --boot --identify --json` and consume the fresh record. Re-running boot rewrites the receipt and clears the guard.
- `state commit` is NOT guarded: ensigns run it, and a subagent's environment carries the root session's id, so an FO compaction mid-flight would wrongly block every ensign commit.

**Fail direction and degradation.** Identity unresolvable (no env — captains in plain terminals, Codex, Pi today) → guard silently no-ops: exactly today's behavior. Transcript missing/unreadable or format drift → fail open with one stderr warning line; malformed transcript lines are skipped. The guard fails CLOSED only on the one condition captured evidence supports. Declared limitation: a host transcript-format change fails open until the pinned fixture is refreshed.

**No escape hatch.** No env kill-switch, no bypass flag. A false refusal's remedy is one cheap idempotent boot — the intended action anyway. An override would hand the failure mode (a session inclined to keep going) its bypass, the same way prose failed.

**Why not the hook (necessity, not preference).** The guard alone delivers the value sentence — the FO *cannot* act on workflow state it has not re-read — for every workflow effect that flows through the binary. A hook adds a host-specific moving part that delivers at most "reminded", requires registration to survive plugin and host changes, produced nothing observable at the incident boundary, and (PreCompact) false-positives on refused compactions. The detector the guard needs is already durable state. Codex and Pi get the same guard architecture with detectors absent today (guard no-ops there); #595 owns the Codex boundary. The existing `codex_session_start_compact.sh` registration is left as-is — removing or fixing it is a separate decision.

**Residual risk (declared).** Actions that never touch the binary — the incident's hand-opened `gh` PRs — are not intercepted. The window is bounded: the FO's first `dispatch build`, gate verb, status mutation, or merge hits the wall, forces the re-boot, and the re-read mods map is what tells the FO that PRs route through `pr-merge`. Closing the `gh` path itself is out of scope (and impossible binary-side).

**No further spike needed:** the design rests on three mechanisms, all exercised — the boundary record's existence and parse (captured incident + spike records, parsed during ideation), env session identity (echoed live post-compaction inside the spike session), and receipt-file round-trip (ordinary file I/O). The captured records seed the implementation's transcript fixture.

## Acceptance criteria

1. **(Value measure)** In the incident-replay fixture — a booted receipt, then a transcript whose tail is the captured-format `compact_boundary` record newer than `booted_at` — the count of guarded verbs that succeed before re-boot is **0 of N** (N = the guarded set above; each exits non-zero with the boot-stale stderr). Baseline on the current binary: N of N succeed — the number that moves the wrong way today. Test: CLI behavior fixture driving the real binary against a temp `$HOME`/projects layout and env.
2. After `status --boot` re-runs in the same fixture, every previously refused verb succeeds unchanged; and in sessions with no boundary (or no resolvable identity), behavior is bit-identical to today — existing golden fixtures pass unmodified. Test: fixture continuation plus full `go test ./...`.
3. The verdict derives from durable state only: deleting the receipt or appending the boundary record flips refusal/pass with no conversational input of any kind. Test: unit verdict matrix (no receipt / fresh / stale / missing transcript / no env / malformed lines) plus fixture file-toggling. (Mechanism AC serving AC 1.)
4. The doc diff below is applied at implementation. Test: the contract lines appear in the built plugin skill and `docs/runtime-support.md` renders the new section.

## Test plan

- Go unit, `internal/bootguard` (new package): receipt round-trip; transcript scan over a fixture seeded with the captured incident and spike records plus malformed and irrelevant lines; the verdict matrix. Cheap; no live host.
- CLI behavior fixture: the AC-1/AC-2 journey (refuse-all → boot → pass-all) with golden stderr for the refusal text.
- No live workflow smoke test: the only runtime claim — the boundary is observable durable state — is what ideation's live captures already prove; every implementation claim is binary behavior, fixture-provable.

## Expected surface

Estimate net LOC change: **+540**, across **12 files** (insertions ≈ +560, deletions ≈ −20). Tolerance: 2x net (workflow default). Breakout, per the 2026-08-18 calibration (6ht/j7j/vka all landed product near estimate and tests at ~2x instinct — the test budget below is twice first instinct, stated explicitly):

- Product ≈ +145: `internal/bootguard/` (receipt, scan, verdict, identity resolution) ~90; receipt write in the `status --boot` path ~15; preflight wiring at guarded-verb entry points ~30; refusal text ~10.
- Tests ≈ +340 (first instinct ~170, doubled): unit verdict matrix, transcript fixture, CLI behavior fixture, goldens.
- Docs/contract ≈ +55: `docs/runtime-support.md` section, `skills/first-officer/references/first-officer-shared-core.md` boot-section lines, exit-code note.

Files (~12): bootguard.go, bootguard_test.go, transcript fixture(s), status-boot wiring, ~4 CLI entry-point wirings, runtime-support.md, first-officer-shared-core.md, behavior-fixture files.

## Semantics that may change (declared)

- Runtime behavior / authority: guarded verbs gain refusal authority (exit 4 + stderr) when the running session's boot is stale or absent — including for a human driving spacedock inside any Claude session who has not run boot (the refusal names the one command to run).
- Stored format: new project-local receipt `.spacedock/boot/{session_id}.json` (host scratch; gitignored if the repo tracks an ignore file).
- `status --boot` gains a write side effect (the receipt); workflow and entity state remain read-only.
- Command grammar: unchanged — no new flags, verbs, or hook registrations.

## Doc diff (proposed, applied at implementation)

`skills/first-officer/references/first-officer-shared-core.md`, `«state.boot»()` section — after the `- **effect:** …` line, add:

```
- **receipt:** the shipped command also writes a session boot receipt (`.spacedock/boot/{session_id}.json`, host scratch — workflow state stays read-only). FO-held mutating verbs refuse with exit 4 (BOOT_STALE) when the running session has no receipt or compacted after it booted; the remedy the stderr names is exactly this call — re-run it and consume the fresh record.
```

`docs/runtime-support.md`, new subsection after "Runtime layers":

```
## Boot guard at the compaction boundary

A compaction-resumed session keeps its session id and transcript; the host
records the boundary durably (a `compact_boundary` record in the session
transcript). `status --boot` writes a per-session receipt; gate, merge-guard,
dispatch-build, new, and status-mutation verbs refuse (exit 4) when the
receipt is missing or older than the latest boundary, until boot re-runs.
Detection resolves per host: Claude Code via `CLAUDE_CODE_SESSION_ID` plus the
project transcript path; hosts without a resolvable identity degrade to a
silent no-op (Codex: #595). The guard fails open on unreadable transcripts and
never needs a hook.
```

## Stage Report: ideation

- DONE: Prove the Claude compaction boundary exists with a CAPTURED EVENT, not documentation — this repo's own hook-events.jsonl holds exactly one PreToolUse record and no compaction evidence, so nothing local proves it today. If no usable boundary exists, record it as unsupported rather than shipping a prose rule dressed as a mechanism.
  Captured live on Claude Code v2.1.226: PreCompact + SessionStart(source:compact) hook events and the durable `compact_boundary` transcript record, plus the incident session's own boundary record — verbatim in `force-boot-at-compaction-boundary/ideation-spike-evidence.md`.
- DONE: Decide between the host-hook mechanism and the host-independent fail-closed tell (the FO must hold a boot record whose identity matches the current session, else it refuses gate and merge actions until it re-boots) on necessity, not preference — the second needs no host callback and would cover Codex and Pi too.
  Fail-closed guard chosen; the naive session-identity tell is disproven by capture (sessionId survives compaction), so the guard compares a boot receipt against the durable `compact_boundary` record instead. The hook is rejected on necessity: advisory-only, produced nothing at the incident boundary, and PreCompact false-positives on refused compactions (captured).
- DONE: Declare net LOC change and file count with a REALISTIC test budget. Today three consecutive entities blew their estimate and every single overrun was under-budgeted test-side work; product code landed near its number each time. Budget test code at roughly twice your first instinct and say so explicitly.
  Net +540 across ~12 files (≈ +560 / −20), tolerance 2x net; breakout product +145 / tests +340 (first instinct ~170, explicitly doubled) / docs +55.

### Summary

Ran a live capture spike (tmux-driven Claude 2.1.226 session with PreCompact/SessionStart capture hooks) and proved the compaction boundary three ways: hook events fire, the transcript records a durable timestamped `compact_boundary`, and session identity survives the boundary — which kills the naive session-id tell and the seed's hook assumption in one stroke. Settled the design as a binary-side fail-closed guard: `status --boot` writes a per-session receipt, and FO-held mutating verbs refuse until the receipt postdates the latest boundary, with no hook, no escape hatch, and honest fail-open degradation on non-Claude hosts. Declared surface +540 net across ~12 files with the test budget explicitly doubled per today's calibration.
