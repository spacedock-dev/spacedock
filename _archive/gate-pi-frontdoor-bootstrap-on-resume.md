---
id: 4avk4msa3ktyk1fdt6j5ktw1
title: Gate the Pi front door bootstrap on resume
status: done
source: "Captain CL, 2026-08-25: 'spacedock pi --resume didn't avoid loading the spacedock initial contract.' The Pi front door (internal/cli/pi.go) appends piBootstrapPrompt unconditionally — no containsResume gate, unlike the Claude/Codex front door (internal/cli/frontdoor.go:428,447) which suppresses its bootstrap prompt on --resume/-r/--continue/-c. The spacedock .pi extension session_start handler also re-injects FO_BOOTSTRAP_TEXT with no resume detection. CL hypothesis 'compaction hook leaked into general startup' checked and disproven: session_compact is correctly scoped to the compaction event (PR #738 / force-boot-at-compaction-boundary); the leak is the front door + session_start, neither resume-aware."
gates:
    version: 1
    records:
        - id: gate:4avk4msa3ktyk1fdt6j5ktw1:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:4avk4msa3ktyk1fdt6j5ktw1-backlog-1
              briefing:
                id: briefing:4avk4msa3ktyk1fdt6j5ktw1:backlog:attempt-1:revision-1
                digest: sha256:b3f27b5850f0b44aea082a43333494efec5ca22b1f64f3e0561a0372fe956e40
                request-digest: sha256:c7df81974554a473c4a131ead3f0471be0554004d090c66c69582482a99e48c0
                room-ref: ./gate-pi-frontdoor-bootstrap-on-resume/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:4avk4msa3ktyk1fdt6j5ktw1:backlog:1
                briefing: briefing:4avk4msa3ktyk1fdt6j5ktw1:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-27T00:11:22.702651Z"
                decision: approve
                reason: 'Captain approve: enter ideation to flesh out the approach and test plan'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:4avk4msa3ktyk1fdt6j5ktw1:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:4avk4msa3ktyk1fdt6j5ktw1-ideation-1
              briefing:
                id: briefing:4avk4msa3ktyk1fdt6j5ktw1:ideation:attempt-1:revision-1
                digest: sha256:fcdcd6469e168a567db59ca473a6e4e6eca09efd4721140dd5389984badb5337
                request-digest: sha256:7e4f6dcfec058013a65787a90c11eb9fe118d48ee5a90a05540cb511211d2c4a
                room-ref: ./gate-pi-frontdoor-bootstrap-on-resume/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:4avk4msa3ktyk1fdt6j5ktw1:ideation:1
                briefing: briefing:4avk4msa3ktyk1fdt6j5ktw1:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-27T01:59:14.369415Z"
                decision: approve
                reason: 'Captain approve: enter implementation'
              application:
                target-stage: implementation
                state: consumed
        - id: gate:4avk4msa3ktyk1fdt6j5ktw1:validation
          stage: validation
          attempts:
            - id: gate-attempt:4avk4msa3ktyk1fdt6j5ktw1-validation-1
              briefing:
                id: briefing:4avk4msa3ktyk1fdt6j5ktw1:validation:attempt-1:revision-1
                digest: sha256:e455516eef36fc670f79984a1425c5c500acae5f42803c08fa07dc443a6985f1
                request-digest: sha256:b11f2127b0897d6ec9335c9942e38718e8f03f269f555779c2b0d3af6854ace9
                room-ref: ./gate-pi-frontdoor-bootstrap-on-resume/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:4avk4msa3ktyk1fdt6j5ktw1:validation:1
                briefing: briefing:4avk4msa3ktyk1fdt6j5ktw1:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-27T06:20:20.333079Z"
                decision: approve
                reason: 'Captain approve: PASSED'
              application:
                target-stage: done
                state: consumed
started: 2026-08-27T00:11:57Z
worktree: .worktrees/spacedock-ensign-gate-pi-frontdoor-bootstrap-on-resume
mod-block:
pr: pr-merge:777
verdict: PASSED
completed: 2026-08-28T21:59:55Z
archived: 2026-08-28T21:59:56Z
---

`spacedock pi --resume` must not load the Spacedock first-officer contract as if starting fresh; a resume carries its own session intent and the contract survives in the system prompt via resources_discover. Today it loads the contract via two unconditional paths.

## Problem

`internal/cli/pi.go:271` appends `piBootstrapPrompt` ("Use $spacedock:first-officer for this whole Pi session.") to the Pi argv with no resume gate — `pi.go` has zero `containsResume` calls; `frontdoor.go` (Claude/Codex) has 3 and gates its bootstrap on `if !resume` (frontdoor.go:428,447). So `spacedock pi -- --resume` still injects the launch prompt as a durable user message that tells the resumed session to load the FO contract as if starting fresh.

Second path (`.pi/extensions/spacedock.ts`): the `session_start` handler (line 80) sets `injectBootstrap=true` unconditionally. The `context` hook then re-injects `FO_BOOTSTRAP_TEXT` on each turn. Pi's extension docs (`tmp/pi/packages/coding-agent/docs/extensions.md:653`) describe the `context` event's `event.messages` as a **deep copy, safe to modify** — the returned message array is what the LLM sees for that turn; modifications are **non-destructive** and do NOT persist to the session transcript. So `hasStructuralBootstrap` (which checks `event.messages` for a prior injection) never matches: the per-turn injection is transient, not a durable logged message, so the injection re-fires every turn (fresh and resume alike) but does not corrupt the resumed transcript with a new user message. Additionally, Pi's lifecycle docs (`extensions.md:281,319,398`) show `session_start` fires with `reason: "startup"` on CLI launch — `reason: "resume"` is only for mid-session switches via `/resume`, NOT for CLI `pi --resume`. So the extension's `session_start` cannot detect a CLI resume via `event.reason`. The extension has compaction-awareness (PR #738: `session_compact` → `injectBootRecord`) but this is a separate event untouched by the front-door fix.

## Proposed approach

**Decision: pi.go-only fix. No extension change.**

Gate the Pi front door bootstrap on `!resume`, mirroring the Claude/Codex front door: wrap the `launchPrompt(piBootstrapPrompt, fd)` append (pi.go:271) in `if !containsResume(fd.passthrough)`, reusing the existing shared `containsResume` function (frontdoor.go:553-562). Pi's resume flags (`-r`, `--resume`, `-c`, `--continue` — confirmed in `tmp/pi/packages/coding-agent/README.md:566-571`) match `containsResume`'s token set exactly. No new function or token set needed.

### Per-mechanism breakdown

1. **pi.go resume gate** (`containsResume(fd.passthrough)` guard on `piBootstrapPrompt`).
   - Value AC served: AC-1 (argv omits the prompt on resume), AC-2 (token-set parity).
   - Simplest alternative: suppress only on an exact `--resume` match.
   - Why insufficient: misses `-r`, `-c`, `--continue`, and `--resume=<id>` — all valid Pi resume forms.

2. **Extension `session_start` resume-awareness** — CONSIDERED AND REJECTED.
   - Value AC it would serve: AC-3 (no fresh `FO_BOOTSTRAP_TEXT` on resume).
   - Simplest alternative: `if (event.reason === "resume") return;` in the `session_start` handler.
   - Why unnecessary/insufficient: (a) CLI `pi --resume` fires `session_start` with `reason: "startup"`, not `"resume"` (Pi docs lifecycle: `reason: "resume"` is mid-session `/resume` only) — the extension cannot detect CLI resume via `event.reason`; (b) the `context` injection is transient (deep copy, non-destructive per Pi docs) — it does not add a durable message to the resumed transcript, so it does not re-load the contract as a fresh start; the FO sees a redundant per-turn directive in its context window, not a new bootstrap user message; (c) touching the extension risks the PR #738 compaction-boundary behavior (`session_compact` → `injectBootRecord`) for no gain.

## Risk evidence

No spike needed: proven mechanisms.
1. `containsResume` is a proven shared function with 3 live call sites in `frontdoor.go` (lines 428, 447, 553-562) — reusing it for pi.go adds no new mechanism.
2. Pi's `context` hook is transient: Pi extension docs (`extensions.md:653`) state `event.messages` is a "deep copy, safe to modify" and changes are "non-destructive" — the injection does not persist to the transcript, so it cannot re-load the contract as a fresh start on resume.
3. Pi's `session_start` fires with `reason: "startup"` on CLI launch (Pi docs `extensions.md:281,398`) — `reason: "resume"` is mid-session only, so the extension cannot detect CLI resume.
4. `piBootstrapPrompt` is the durable launch argv token (pi.go:271) — the only path that adds a fresh bootstrap user message to a resumed session. Gating it is sufficient.

## Out of scope

Changing the compaction-boundary behavior (PR #738 / `force-boot-at-compaction-boundary` owns that). The ensign child-session bootstrap exemption (`PI_SUBAGENT_CHILD=1`, owned by `pin-ensign-contract-entry-point`). The Claude/Codex front doors (already gated on resume). The extension's `session_start`/`context` handlers (transient injection, not harmful on resume — see Risk evidence).

## Expected surface and tolerance

**Fix is pi.go-only.** Insertions: ~+3 (pi.go: the `if !containsResume` guard wrapping the existing append). Deletions: 0. Net: +3 for the code change. Test: ~+30 (a table-driven resume-suppression test in `internal/cli/pi_frontdoor_test.go`). Total: ~+33 insertions, 0 deletions, net +33, across 2 files (`internal/cli/pi.go`, `internal/cli/pi_frontdoor_test.go`). Tolerance: ±50%.

Observable semantics: the Pi launch argv on resume (`--resume`, `-r`, `-c`, `--continue`) omits `piBootstrapPrompt`; a non-resume launch still appends it. No CLI grammar change. No stored-format change. The extension's `session_start`/`context` behavior is unchanged. No user-facing documentation describes the Pi bootstrap prompt or its resume suppression (checked `docs/site/` — no reference to `piBootstrapPrompt` or Pi resume behavior), so no doc diff is required.

Residual: Pi's `--session <path|id>` flag (loads a specific session file — also a resume form) is NOT covered by the shared `containsResume` token set (which mirrors the Claude/Codex front door). `--session` is Pi-specific and would require a Pi-specific gate addition. Flagged as a follow-up, not in this fix's scope (AC-2 pins parity with the existing front-door token set).

## Acceptance criteria

**AC-1 (value) — `spacedock pi --resume` does not inject the fresh-start FO contract as a durable user message.**
Verified by: a fixture in `internal/cli/pi_frontdoor_test.go` that launches the Pi front door with `--resume` in the passthrough and asserts the assembled argv does NOT contain `piBootstrapPrompt`, while a launch without `--resume` still contains it. Fails today (pi.go:271 appends unconditionally). Baseline: the non-resume argv contains the prompt — this test fails if the gate regresses or the non-resume path loses the prompt.

**AC-2 (serves AC-1) — the resume gate matches the Claude/Codex front door's resume token set.**
Verified by: a table-driven test asserting `--resume`, `--resume=<id>`, `-r`, `--continue`, `-c` all suppress the prompt (mirroring frontdoor.go:553-562), and non-resume passthrough (e.g. `--model`, a task string) does not. Fails today (no gate exists).

**AC-3 (serves AC-1) — the extension's transient `FO_BOOTSTRAP_TEXT` injection does not re-load the contract as a fresh start on resume.**
Verified by: code+docs reasoning recorded in Risk evidence — the `context` hook is transient (deep copy, non-destructive per Pi docs), so it does not add a durable message to the resumed transcript. No extension change needed. The residual is `--session <path|id>` (Pi resume form not covered by `containsResume`), flagged as a follow-up.

## Test plan

Fixture in `internal/cli/pi_frontdoor_test.go`: a table-driven resume-suppression test mirroring `frontdoor_test.go`'s codex resume tests (frontdoor_test.go:876-915). For each resume form (`--resume`, `--resume=<id>`, `-r`, `--continue`, `-c`), launch `runPi` with the form in the passthrough and assert `piBootstrapPrompt` is absent from `ops.launched`. For a non-resume passthrough (e.g. `--model google/gemini`), assert the prompt is present as the last argv token. Reuses the existing `fakePiRuntimeOps` harness and `healthyPiPackageStatus` fixtures. No live probe needed (mechanisms proven by code + Pi docs reading — see Risk evidence). Cost: low — the gate is ~3 lines + a ~30-line test.

### Feedback Cycles

## Stage Report: ideation

- DONE: Flesh out the four ideation sections (Problem, Proposed approach, Acceptance criteria, Test plan)
  Entity body updated with refined Problem (context hook transient per Pi docs; session_start reason "startup" on CLI launch), Proposed approach (firm pi.go-only decision with per-mechanism breakdown), AC-1/AC-2/AC-3 (value-measuring AC-1 against non-resume baseline), and Test plan (table-driven fixture in pi_frontdoor_test.go).
- DONE: For every new mechanism, name the value AC it serves, the simplest alternative, and why insufficient
  Mechanism 1 (pi.go resume gate): serves AC-1/AC-2; simplest alternative is exact `--resume` match; insufficient because it misses `-r`/`-c`/`--continue`/`--resume=<id>`. Mechanism 2 (extension session_start resume-awareness): serves AC-3; simplest alternative is `event.reason === "resume"` guard; rejected because CLI launch fires reason "startup" not "resume", and the context injection is transient (non-destructive).
- DONE: Record Risk evidence (exercise riskiest mechanism OR record "no spike needed")
  Recorded "no spike needed: proven mechanisms" — containsResume proven (3 frontdoor.go call sites); context hook transient (Pi docs: deep copy, non-destructive); session_start reason "startup" on CLI launch (Pi docs lifecycle); piBootstrapPrompt is the only durable bootstrap path (pi.go:271).
- DONE: Declare expected surface (net LOC + files + tolerance, insertions/deletions separately) and observable semantics
  ~+33 insertions, 0 deletions, net +33, across 2 files (pi.go, pi_frontdoor_test.go); tolerance ±50%. Observable: resume argv omits piBootstrapPrompt; no CLI grammar/stored-format change; extension unchanged. Fix is pi.go-only.
- DONE: Name whether the fix is pi.go-only or also touches the extension
  pi.go-only. Extension session_start change considered and rejected (see Per-mechanism breakdown #2).

### Summary

Fleshed out the ideation for gating the Pi front door bootstrap on resume. Key decision: pi.go-only fix (add `containsResume` guard wrapping the `piBootstrapPrompt` append), no extension change. The extension's `FO_BOOTSTRAP_TEXT` injection is transient (Pi `context` hook uses a deep copy per Pi docs) so it does not re-load the contract as a durable message on resume, and `session_start` fires with `reason: "startup"` on CLI launch so it cannot detect CLI resume anyway. No spike needed — all mechanisms proven by code + Pi docs reading. Residual: Pi's `--session <path|id>` flag is a resume form not covered by the shared `containsResume` token set, flagged as a follow-up.

## Stage Report: implementation

- DONE: the pi.go `containsResume` guard wrapping the `piBootstrapPrompt` append (pi.go:271) is in place
  Mirrors frontdoor.go's token set (`--resume`, `--resume=<id>`, `-r`, `--continue`, `-c`) via the shared `containsResume(fd.passthrough)` call — no new function or token set added.
- DONE: the table-driven `TestPiResumeSuppressesBootstrapPrompt` in `internal/cli/pi_frontdoor_test.go` asserts each resume token suppresses `piBootstrapPrompt` from `ops.launched`
  Non-resume passthrough (`--model google/gemini`, a task string) keeps the prompt as the last argv token (independent non-resume baseline).
- DONE: test and gofmt results
  `go test ./internal/cli/ -run TestPiResumeSuppressesBootstrapPrompt` is green; `gofmt -l` on the two changed files is clean; the full `go test ./internal/cli/` suite's only failure is `TestVersionAmbiguousMarkersExitZero`, a pre-existing environmental failure (`PI_CODING_AGENT=true` in the worker env) reproduced on the base branch via `git stash`, unrelated to this change.

### Summary

Implemented the Pi front door resume gate: wrapped the `launchPrompt(piBootstrapPrompt, fd)` append in `internal/cli/pi.go` with `if !containsResume(fd.passthrough)` so a resume launch omits the fresh-start FO bootstrap prompt, matching the Claude/Codex front door. Added a table-driven `TestPiResumeSuppressesBootstrapPrompt` in `internal/cli/pi_frontdoor_test.go` pinning AC-1/AC-2. pi.go-only change, no extension touch. Committed to `spacedock-ensign/gate-pi-frontdoor-bootstrap-on-resume`.

## Stage Report: validation

- DONE: AC-1 — `TestPiResumeSuppressesBootstrapPrompt` launches Pi with `--resume` in passthrough and asserts argv does NOT contain `piBootstrapPrompt`; non-resume keeps it as the last argv token
  `go test ./internal/cli/ -run TestPiResumeSuppressesBootstrapPrompt -v -count=1` → PASS (7/7 subtests: resume/--resume, resume/--resume=abc123, resume/-r, resume/--continue, resume/-c, nonresume/model_flag, nonresume/task_string). Per-token: resume tokens suppress the prompt (argv has no piBootstrapPrompt token); non-resume `--model google/gemini` and `review this code` keep piBootstrapPrompt as the last argv token.
- DONE: AC-2 — resume gate matches Claude/Codex front door token set; `containsResume` is the shared helper at frontdoor.go:553 (no new token set added)
  `grep -n "func containsResume" internal/cli/frontdoor.go` → 553:func containsResume; callsites: frontdoor.go:428 (Claude), frontdoor.go:447 (Codex, via `resume :=`), pi.go:277 (this change). Table-driven test covers all five tokens (`--resume`, `--resume=<id>`, `-r`, `--continue`, `-c`) plus two non-resume cases.
- DONE: AC-3 — extension `.pi/extensions/spacedock.ts` UNCHANGED (pi.go-only fix)
  `git diff --numstat main` → only internal/cli/pi.go (9/1) and internal/cli/pi_frontdoor_test.go (70/0); no `.pi/extensions/spacedock.ts` in the diff.
- DONE: pre-existing failure isolated
  `TestVersionAmbiguousMarkersExitZero` fails identically on main (`go test ./internal/cli/ -run TestVersionAmbiguousMarkersExitZero` on `43fc79a23 [main]` → same `want "Runtime: ambiguous (CODEX_THREAD_ID, CLAUDECODE)"` vs got `... (CODEX_THREAD_ID, CLAUDECODE, PI_CODING_AGENT)`) because the worker env sets `PI_CODING_AGENT=true`. Unrelated to this change.
- DONE: semantic adversarial pass
  Non-resume launch still gets the prompt (nonresume/model_flag, nonresume/task_string PASS — baseline that can move wrong is intact). The guard is exactly the `if !containsResume(fd.passthrough) { argv = append(...) }` wrap; no other argv mutation (diff confirms only the one append is wrapped). gofmt clean, go vet clean, -race clean on the focused test.

### Summary

Validation PASSED. The pi.go-only resume gate (`if !containsResume(fd.passthrough)` wrapping `launchPrompt(piBootstrapPrompt, fd)`) suppresses the fresh-start FO bootstrap prompt on all five resume tokens and preserves it on non-resume launches, matching the Claude/Codex front door. Extension unchanged. The sole suite failure (`TestVersionAmbiguousMarkersExitZero`) is a pre-existing environmental failure reproduced identically on main. Recommend PASSED.
