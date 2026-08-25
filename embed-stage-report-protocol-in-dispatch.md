---
id: t4rqqmmrqh69adgb1bbdx838
title: Embed the stage-report protocol in the dispatch artifact
status: ideation
source: "/tmp/spacedock-pi-dispatch-diagnosis.md (2026-08-25): 3 of 4 Pi-dispatched ensigns completed implementation without writing a ## Stage Report; the dispatch build artifact's ## First action claims the file contains the stage-report format, but the body has 0 such tokens. Same class as archived pin-ensign-contract-entry-point (2026-08-01), which fixed the spawn binding but not the artifact-body gap."
gates:
    version: 1
    records:
        - id: gate:t4rqqmmrqh69adgb1bbdx838:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:t4rqqmmrqh69adgb1bbdx838-backlog-1
              briefing:
                id: briefing:t4rqqmmrqh69adgb1bbdx838:backlog:attempt-1:revision-1
                digest: sha256:e0f39a8f65f844ee0ee6111067ffed3fd4c9fce97862567591001bd1e52c8b1f
                request-digest: sha256:7c11d27c46f2d428b7f1f5f7b67b75645e3b5b84f106d283e0ea7b8294b5a86f
                room-ref: ./embed-stage-report-protocol-in-dispatch/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:t4rqqmmrqh69adgb1bbdx838:backlog:1
                briefing: briefing:t4rqqmmrqh69adgb1bbdx838:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-25T20:19:55.599104Z"
                decision: approve
                reason: 'Captain approve: live-reproduced artifact-body gap on v0.28.0-pre0; scope matches the archived pin-ensign-contract-entry-point neighbor and the systemic self-describing-checklist coverage hole; direction accepted for ideation.'
              application:
                target-stage: ideation
                state: consumed
started: 2026-08-25T20:22:24Z
---

The dispatch artifact must carry the stage-report protocol so a worker writes a complete `## Stage Report:` section on a real dispatch, not only when a smoke-checklist tells it to.

## Problem

`internal/dispatch/build.go` step 8 (around line 673) emits `### Completion checklist` + `### Summary\n{brief description...}` — no `## Stage Report:` template, no DONE/SKIPPED/FAILED structure. The Pi `firstActionBlock` (around lines 891-901) claims "This file contains the shared ensign discipline entry points (stage-report format, polling, worktree ownership, and completion protocol)" — a false claim; the body contains none of it (reproduced on v0.28.0-pre0: 0 `Stage Report`, 0 `DONE:`/`SKIPPED:`/`FAILED:`). Claude's `Skill()` auto-loads the skill; Codex has a load-bearing `$spacedock:ensign; then` bootstrap edge (proven by `internal/dispatch/codex_bootstrap_test.go:60-76`, a mutation test that fails the child before any read when the edge is removed). Pi has neither — `skill="ensign"` is discoverable, not loaded — so 3 of 4 Pi workers missed the stage report on a real dispatch. The format lives in `skills/ensign/references/ensign-shared-core.md:66-95`, never surfaced to a Pi worker whose checklist is the entity's real acceptance criteria.

The coverage gap is systemic: every live lane (`internal/ensigncycle/pi_live_runner_test.go:131-137`, `internal/ensigncycle/shared_fixtures_test.go:41,135-146`) supplies the format via its own checklist or stage-def, then asserts the worker used it — a tautology that cannot fail on the real mode.

## Proposed approach

Embed the stage-report protocol (the `## Stage Report: {stage}` template with the DONE/SKIPPED/FAILED/Summary structure, sourced from `ensign-shared-core.md:66-95`) into the `dispatch build` artifact body for every host, so the `## First action` claim is literally true for the format. Narrow the Pi `firstActionBlock` wording so it no longer overclaims polling/worktree/completion (those stay in the skill) — unless ideation chooses to embed those too. Add a non-self-describing live-lane variant per host (checklist = the entity's real acceptance criteria, with no format or skill-path hints) that asserts the worker still writes a complete stage report. Ideation decides Pi-only vs all-host embed and whether to keep the skill-load instruction alongside the body.

## Risk evidence

Backlog: the gap is live-reproduced — the artifact body has 0 stage-report tokens on v0.28.0-pre0, and the prior fix `pin-ensign-contract-entry-point` proved the worker loads the skill only when the checklist tells it to, not on a real dispatch. Riskiest unverified mechanism: whether embedding the format in the body is sufficient for a Pi worker to write a complete report when its checklist does not mention the format — the non-self-describing live lane exercises exactly this.

## Out of scope

Auto-loading the ensign skill on Pi (the deeper root cause — `skill="ensign"` is discoverable, not loaded) is a pi-subagents / `.pi`-extension concern outside this repo; not this entity. Changing the ensign stage-report contract itself. Claude/Codex skill-load mechanisms, which already work.

## Expected surface and tolerance

Estimate net LOC change: ~+120, across ~4 files (`internal/dispatch/build.go`, a new or extended `internal/dispatch` fixture test, a non-self-describing live-lane variant in `internal/ensigncycle`, and possibly a skill-text reference). Tolerance: ±50%.

## Acceptance criteria

**AC-1 (value) — A Pi-dispatched ensign writes a complete stage report on a real, non-self-describing dispatch.**
Verified by: a live Pi lane that dispatches a worker through `dispatch build --host pi` with a checklist equal to a real entity's acceptance criteria (no "First read ensign/SKILL.md", no "append a stage report with heading..."), then asserts the entity file carries a complete `## Stage Report: implementation` (heading + one DONE/SKIPPED/FAILED per checklist item + `### Summary`) and a clean state-checkout commit. Fails today: with the current artifact body, the worker has no format source on Pi.

**AC-2 (serves AC-1) — The `dispatch build` artifact body carries the stage-report protocol.**
Verified by: a fixture test in `internal/dispatch` that builds an artifact with a non-self-describing checklist for host=pi (and claude, codex) and asserts the body contains `## Stage Report:`, `- DONE:`, `- SKIPPED:`, `- FAILED:`, and `### Summary`. Fails today (0 hits, reproduced). The template tokens' presence in the generated body is the claim — a presence check over generated output, not a behavioral prose-grep over an instruction file.

**AC-3 (serves AC-1) — A non-self-describing live lane exists and closes the tautology.**
Verified by: at least one live lane per host runs with a checklist that does NOT name the ensign skill path, the stage-report heading, or the DONE/Summary structure, and still passes the stage-report assertion; reverting the AC-2 body embed makes that lane RED — the hole that was never tested.

## Test plan

Fixture in `internal/dispatch`: build an artifact with a non-self-describing checklist, assert the body carries the protocol (AC-2). Live: extend `runPiSmokeDispatchBuild` / `assertPiLiveSmokeResult` with a non-self-describing-checklist variant (AC-1, AC-3); add Claude/Codex equivalents if ideation picks all-host scope. Adversarial: revert the body embed, confirm the non-self-describing lane goes RED while the self-describing lane stays green — proving the new lane tests the real mode, not the fixture's hint. `internal/dispatch/build.go` is a front-door launcher surface, so the README's detached-adversarial-audit trigger may fire; ideation scopes it.

### Feedback Cycles
