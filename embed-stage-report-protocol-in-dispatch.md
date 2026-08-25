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
        - id: gate:t4rqqmmrqh69adgb1bbdx838:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:t4rqqmmrqh69adgb1bbdx838-ideation-1
              briefing:
                id: briefing:t4rqqmmrqh69adgb1bbdx838:ideation:attempt-1:revision-1
                digest: sha256:aad9a5cd7a059bef30d67afa7fab9179d845617404cfd574bf164e711eb8aeca
                request-digest: sha256:b0275dbec35d2beacd6813e25faf0ee9a0794e29e93ac45ff686e5222349fac3
                room-ref: ./embed-stage-report-protocol-in-dispatch/review/ideation/briefing-1
started: 2026-08-25T20:22:24Z
---

The dispatch artifact must carry the stage-report protocol so a worker writes a complete `## Stage Report:` section on a real dispatch, not only when a smoke-checklist tells it to.

## Problem

`internal/dispatch/build.go` step 8 (around line 673) emits `### Completion checklist` + `### Summary\n{brief description...}` — no `## Stage Report:` template, no DONE/SKIPPED/FAILED structure. The Pi `firstActionBlock` (around lines 891-901) claims "This file contains the shared ensign discipline entry points (stage-report format, polling, worktree ownership, and completion protocol)" — a false claim; the body contains none of it (reproduced on v0.28.0-pre0: 0 `Stage Report`, 0 `DONE:`/`SKIPPED:`/`FAILED:`). Claude's `Skill()` auto-loads the skill; Codex has a load-bearing `$spacedock:ensign; then` bootstrap edge (proven by `internal/dispatch/codex_bootstrap_test.go:60-76`, a mutation test that fails the child before any read when the edge is removed). Pi has neither — `skill="ensign"` is discoverable, not loaded — so 3 of 4 Pi workers missed the stage report on a real dispatch. The format lives in `skills/ensign/references/ensign-shared-core.md:66-95`, never surfaced to a Pi worker whose checklist is the entity's real acceptance criteria.

The coverage gap is systemic: every live lane (`internal/ensigncycle/pi_live_runner_test.go:131-137`, `internal/ensigncycle/shared_fixtures_test.go:41,135-146`) supplies the format via its own checklist or stage-def, then asserts the worker used it — a tautology that cannot fail on the real mode.

## Proposed approach

**Decision: Pi-only embed.** Embed the stage-report protocol — the `## Stage Report: {stage}` template with the DONE/SKIPPED/FAILED/Summary structure, sourced from `ensign-shared-core.md:66-95` — into the `dispatch build` artifact body only when `host=pi`. Rationale: the gap is Pi-specific — Claude auto-loads the skill via `Skill()`, and Codex has the proven `$spacedock:ensign; then` bootstrap edge (`codex_bootstrap_test.go:60-76`). The adversarial revert (AC-3) can only go RED for Pi; for Claude/Codex the skill supplies the format, so an all-host embed would be untestable redundancy. All-host embed considered and rejected: it adds the same template block to hosts that already work, with no adversarial test that can fail on revert for them.

**Decision: keep the skill-load instruction for Claude/Codex; the embedded format is the format source for Pi.** For Pi, there is no skill-load instruction to keep or remove — Pi does not auto-load the skill, and auto-loading is out of scope. The embedded format block replaces the skill as the format source for Pi only. For Claude/Codex, the existing `Skill()` / `$spacedock:ensign` bootstrap stays unchanged; no body embed is added.

**build.go step-8 edit.** Step 8 (`build.go:672-674`) currently emits `### Completion checklist` + `### Summary\n{brief description...}`. The edit adds a stage-report format template block after the Summary slot, conditional on `host == "pi"`:

```markdown
### Stage Report format

Append a `## Stage Report: {stage}` section at the end of the entity file using this structure:

- DONE: {item text}
  {one-line evidence or reference}
- SKIPPED: {item text}
  {one-line rationale}
- FAILED: {item text}
  {one-line details}

### Summary
{2-3 sentences: what was done, key decisions, anything notable}

Every checklist item must appear. Use `- DONE:` / `- SKIPPED:` / `- FAILED:` markers. Do not use checkbox markers. Append at the end of the entity file.
```

This block is sourced from `ensign-shared-core.md:66-95` (the `## Stage Report Protocol` section), trimmed to the template + rules a worker needs to produce the report.

**firstActionBlock wording correction (Pi only).** The Pi `firstActionBlock` (`build.go:897-904`) currently claims:

> This file contains the shared ensign discipline entry points (stage-report format, polling, worktree ownership, and completion protocol) plus the stage-specific assignment.

After the embed, the body carries the stage-report format but NOT polling/worktree/completion. Corrected wording:

> This file carries the stage-report format template plus the stage-specific assignment. The ensign skill supplies the remaining shared discipline (polling, worktree ownership, completion protocol); on Pi it is discoverable (`skill="ensign"`), not auto-loaded.

The Claude and Codex `firstActionBlock` strings are unchanged — their claims already correctly attribute the discipline to the skill.

### Mechanism justification

1. **Embedded protocol block (Pi body).** Serves AC-1 (worker writes complete report) and AC-2 (body carries protocol). Simplest alternative: add a checklist hint ("append a stage report with heading `## Stage Report:`...") to every Pi dispatch. Why insufficient: that is the self-describing tautology this entity exists to close — the real dispatch checklist is the entity's acceptance criteria, which does not mention the format.
2. **Non-self-describing live-lane variant (Pi).** Serves AC-1 and AC-3. Simplest alternative: extend the existing `runPiSmokeDispatchBuild` with a second checklist variant that mentions the format. Why insufficient: a format-mentioning checklist is the tautology — it cannot fail on revert because the hint, not the body embed, is what makes the worker produce the report.
3. **firstActionBlock rewording (Pi).** Serves AC-1 (the worker trusts the body as the format source because the First-action claim says so). Simplest alternative: delete the overclaim, say nothing about what the file contains. Why insufficient: the worker needs to know the body has the format template; without the claim, the worker may not look for it in the body.

## Risk evidence

Backlog: the gap is live-reproduced — the artifact body has 0 stage-report tokens on v0.28.0-pre0, and the prior fix `pin-ensign-contract-entry-point` proved the worker loads the skill only when the checklist tells it to, not on a real dispatch. Riskiest unverified mechanism: whether embedding the format in the body is sufficient for a Pi worker to write a complete report when its checklist does not mention the format — the non-self-describing live lane exercises exactly this. No spike needed: the build pipeline already emits text into the dispatch body (step 8 emits the checklist + summary slot), and Pi workers already read and follow the dispatch file's instructions (proven by every existing Pi live lane in `internal/ensigncycle`). The non-self-describing live lane (AC-1/AC-3) is the first test of whether the embedded format alone suffices — it is the implementation's first test, seeded by this risk evidence.

## Out of scope

Auto-loading the ensign skill on Pi (the deeper root cause — `skill="ensign"` is discoverable, not loaded) is a pi-subagents / `.pi`-extension concern outside this repo; not this entity. Changing the ensign stage-report contract itself. Claude/Codex skill-load mechanisms, which already work.

## Expected surface and tolerance

Estimate net LOC change: ~+120, across ~3 files. Insertions: ~+130, deletions: ~-10, net ~+120. Files: `internal/dispatch/build.go` (+~30 / -~10: stage-report format block function, step-8 conditional, Pi firstActionBlock rewording), a new fixture test in `internal/dispatch` (+~50: build artifact with non-self-describing checklist, assert body contains protocol tokens for host=pi, assert absence for host=claude and host=codex), a non-self-describing live-lane variant in `internal/ensigncycle` (+~50: Pi non-self-describing checklist + assertion). No skill-text reference change needed (Pi-only embed does not touch skill files). Tolerance: ±50%.

Observable semantics this task may change: dispatch artifact body contents (Pi host gains a stage-report format block); the Pi `## First action` claim (narrowed from full ensign discipline to stage-report format only). No CLI grammar change, no stored-format change, no authority change, no runtime behavior change. The Claude and Codex dispatch artifacts are unchanged.

## Acceptance criteria

**AC-1 (value) — A Pi-dispatched ensign writes a complete stage report on a real, non-self-describing dispatch.**
Verified by: a live Pi lane that dispatches a worker through `dispatch build --host pi` with a checklist equal to a real entity's acceptance criteria (no "First read ensign/SKILL.md", no "append a stage report with heading..."), then asserts the entity file carries a complete `## Stage Report: implementation` (heading + one DONE/SKIPPED/FAILED per checklist item + `### Summary`) and a clean state-checkout commit. Fails today: with the current artifact body, the worker has no format source on Pi.

**AC-2 (serves AC-1) — The `dispatch build` artifact body carries the stage-report protocol for Pi.**
Verified by: a fixture test in `internal/dispatch` that builds an artifact with a non-self-describing checklist for host=pi and asserts the body contains `## Stage Report:`, `- DONE:`, `- SKIPPED:`, `- FAILED:`, and `### Summary`. The same test asserts host=claude and host=codex artifacts do NOT contain the embedded block (Pi-only scope). Fails today (0 hits, reproduced). The template tokens' presence in the generated body is the claim — a presence check over generated output, not a behavioral prose-grep over an instruction file.

**AC-3 (serves AC-1) — A non-self-describing Pi live lane exists and closes the tautology.**
Verified by: a Pi live lane runs with a checklist that does NOT name the ensign skill path, the stage-report heading, or the DONE/Summary structure, and still passes the stage-report assertion (`## Stage Report: implementation` + `- DONE:` + `### Summary`); reverting the AC-2 body embed makes that lane RED — the hole that was never tested.

## Test plan

Fixture in `internal/dispatch`: build an artifact with a non-self-describing checklist for host=pi, assert the body carries the protocol tokens (AC-2); assert host=claude and host=codex artifacts do NOT carry the embedded block. Cost: low — pure output assertion over generated text, no runtime. Complexity: ~50 lines.

Live: add a non-self-describing-checklist variant to the Pi live smoke runner (AC-1, AC-3) — a checklist equal to a real entity's acceptance criteria with no format or skill-path hints. Cost: medium — requires a live Pi dispatch round-trip. Complexity: ~50 lines, reuses `runPiSmokeDispatchBuild` / `assertPiLiveSmokeResult` infrastructure.

Adversarial: revert the body embed in `build.go`, confirm the non-self-describing Pi lane goes RED while the existing self-describing Pi lane stays green — proving the new lane tests the real mode, not a fixture hint. Cost: low — a code revert + test run. `internal/dispatch/build.go` is a front-door launcher surface, so the README's detached-adversarial-audit trigger may fire; ideation scopes it to the fixture + live lane only (no detached audit needed for an ideation-stage body-text embed).

### Feedback Cycles

## Stage Report: ideation

- DONE: Four-section design completion (problem statement, proposed approach, acceptance criteria, test plan)
  Entity body carries Problem (live-reproduced artifact-body gap on v0.28.0-pre0), Proposed approach (Pi-only embed via build.go step-8 conditional + firstActionBlock wording correction), Acceptance criteria (AC-1/AC-2/AC-3 with test methods), and Test plan (fixture + live + adversarial). Commit ae0ffd8bb.
- DONE: Surface declaration with insertions/deletions/net
  "Expected surface and tolerance" section declares net ~+120, insertions ~+130, deletions ~-10, across ~3 files (build.go, new fixture test in internal/dispatch, non-self-describing live-lane in internal/ensigncycle); tolerance ±50%; observable semantics declared (Pi artifact body, Pi First-action claim).
- DONE: Risk-evidence/no-spike-needed determination
  "Risk evidence" section records no spike needed: build pipeline already emits text into the dispatch body (step 8), and Pi workers already read and follow the dispatch file's instructions (proven by existing Pi live lanes in internal/ensigncycle); the non-self-describing live lane (AC-1/AC-3) is seeded as the implementation's first test of whether the embedded format alone suffices.

### Summary

Ideation design is complete and committed (ae0ffd8bb). The design embeds the stage-report protocol into the dispatch build artifact body for Pi only (build.go step-8 conditional block + firstActionBlock wording correction), closing the gap where Pi-dispatched ensigns missed the stage report because `skill="ensign"` is discoverable, not loaded. Three acceptance criteria gate the implementation: a value AC (AC-1, a non-self-describing Pi live lane writes a complete report), a serves AC-1 fixture (AC-2, body carries protocol tokens for Pi and not for Claude/Codex), and a tautology-closing live lane (AC-3, reverts RED on body-embed removal). No spike needed — the build pipeline and Pi worker behavior are already proven; the non-self-describing live lane is the implementation's first test.
