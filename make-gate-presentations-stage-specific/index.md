---
title: Make gate presentations stage-specific and omit empty result classes
status: validation
id: krbaeb3resfpbh1qvnb65krf
score: 0.8
source: "Captain feedback on 2026-08-08 after repeated gate reviews rendered FAILED: None."
issue:
sprint: durable-decisions
sprint-readiness: ready
started: 2026-08-09T14:51:39Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-make-gate-presentations-stage-specific
pr:
mod-block:
gates:
    version: 1
    records:
        - id: gate:krbaeb3resfpbh1qvnb65krf:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:krbaeb3resfpbh1qvnb65krf-backlog-1
              briefing:
                id: briefing:krbaeb3resfpbh1qvnb65krf:backlog:attempt-1:revision-1
                digest: sha256:87dcf5d1606cb23b261dc62b4e261e7b18d1f6859af6627461f2e184900619be
                request-digest: sha256:f53ac4ab11d714aafe7349c2071e67c1894e4b14f71ae946013cde46f3fc2ce8
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:krbaeb3resfpbh1qvnb65krf:backlog:1
                briefing: briefing:krbaeb3resfpbh1qvnb65krf:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-09T14:50:51.962348Z"
                decision: approve
                reason: Captain directed ideation dispatch; stage-specific evidence removes fabricated rows without changing gate authority.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:krbaeb3resfpbh1qvnb65krf:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:krbaeb3resfpbh1qvnb65krf-ideation-1
              briefing:
                id: briefing:krbaeb3resfpbh1qvnb65krf:ideation:attempt-1:revision-1
                digest: sha256:c028a4ccb394f861e45474aeeaa0b7380b04800516d0c246763e00dd0f79b967
                request-digest: sha256:8c4f31c653488181d11e32aee392450effb9e87a9cb780a7841e838e014e23f3
                room-ref: ./review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:krbaeb3resfpbh1qvnb65krf:ideation:1
                briefing: briefing:krbaeb3resfpbh1qvnb65krf:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-09T18:33:35.248081Z"
                decision: approve
                reason: The bounded three-file design keeps present-gate and gate record as sole owners, proves stage-specific omission behavior through controlled and live presentations, and excludes adjacent F6C, W5, and XX semantics.
              application:
                target-stage: implementation
                state: consumed
---

## Outcome

Each gate review presents only the evidence that can exist at its current stage. Across controlled backlog, ideation, and validation reviews, all three reviews use the existing `present-gate` path and zero reviews invent an empty result-class row.

## Problem

The generic `skills/present-gate/SKILL.md` template displays one placeholder each for `DONE`, `SKIPPED`, and `FAILED`. First Officers repeatedly convert an empty class into a row such as `FAILED: None`; the published first-workflow example similarly foregrounds `0 SKIPPED, 0 FAILED` instead of evidence.

The same template assumes every gate has a Stage Report. A backlog gate reviews a seed task, not a completed stage report. Ideation needs the selected design and its proof plan, while validation needs actual execution results. Treating all three as a checklist roll-up caused a fabricated backlog checklist during task `0y` dispatch.

## Proposed approach

Keep `skills/present-gate/SKILL.md` as the one presenter and replace only its generic evidence middle with a stage-selection table. The common spine remains: task and stage, one recommendation, the already-bound Briefing identity and digest, optional classified findings, and one decision-effect line. Chat and Subspace continue to render that same presentation and return semantic decision/reason input to the First Officer; `${SPACEDOCK_BIN:-spacedock} gate record --decision` remains the only recorder.

| Gate stage | Evidence rendered by the existing presenter | Evidence deliberately absent |
|---|---|---|
| backlog | Outcome; scope boundary (included work and explicit cuts); proof readiness (the proposed observable proof and any known unknown) from the seed body | Stage Report citation, result classes, checklist counts, executed checks, delivery claim |
| ideation | Chosen direction; riskiest-mechanism spike result or the recorded no-spike basis; expected files/LOC/tolerance and semantic-change declaration; each AC's proposed observable proof | Fabricated execution results, delivery readiness, empty result classes |
| validation | Non-empty `DONE`/`SKIPPED`/`FAILED` classes from the real Stage Report; reviewer findings when present; checks actually executed; AC evidence; delivery readiness | Any result-class heading or row whose item count is zero; any empty findings block |

For validation, keep the numeric `Assessment: {N} done, {N} skipped, {N} failed` summary because it reports the whole result set compactly. Render bullets only beneath classes with at least one item. Backlog and ideation do not render the validation assessment or result-class rows at all. The words `None`, `N/A`, and equivalent zero-result placeholders are never substitutes for absent evidence.

This is selection guidance inside the existing prose presenter, not a stage enum, renderer function, output schema, new command, or stored form. `«gate.assemble-verdict»(slug, stage)` already supplies the stage and calls `present-gate`; `status --read --checklist` and `--ac-scan` remain its existing deterministic inputs. No spike needed: current gates already pass the stage into the one presenter, the presenter already owns conditional omission of empty reviewer categories, and the live recorded-gate journey already proves the bound-package-to-presentation-to-recorder sequence.

## Invariants and ownership

- **KRB (`krbaeb3resfpbh1qvnb65krf`)** owns stage-specific evidence selection and omission of empty presentation classes only.
- **F6C (`f6cvn0s87ywbs158yy0b5q7k`)** owns removal/replacement of semantic prose-grep assertions and the wider final gate/help documentation reconciliation. KRB adds no committed prose-presence oracle and does not absorb F6C's inventory.
- **W5 (`w5bfnrvpcphw857nzz93340c`)** owns reliable presentation of the exact `sha256:` Briefing digest in repeated live runs. KRB preserves the existing reviewed-snapshot requirement and does not weaken, duplicate, or claim to fix its reliability.
- **XX (`xxqk1kq7v8h2am9cvm6y8gyw`)** owns whether dispatch build accepts a zero-item dispatch checklist. KRB changes only captain-facing gate evidence; it neither accepts nor rejects an empty worker-dispatch checklist.
- Recorder authority, Briefing/package format, digest computation, command grammar, stored frontmatter, decision authority, runtime routing, and Subspace behavior do not change.

## Concrete presentation contract

The implementation replaces the current single `Checklist` block with one common spine plus the table-driven stage evidence above. The common before/after is:

```text
Before: Checklist (from ## Stage Report ...): DONE/SKIPPED/FAILED placeholders at every stage.
After:  Stage evidence ({backlog | ideation | validation}): only that stage's available evidence; validation emits only non-empty result classes.
```

The `Reviewed snapshot:` line continues to identify the bound Briefing and its digest. `Recommend ...` continues to appear exactly once. `Decision:` continues to name the concrete effect of approve/revise/hold. The assembly rules continue to omit an absent reviewer-findings block and forbid the presentation layer from classifying findings. After the presentation, the First Officer passes the semantic decision and reason to the existing recorder exactly as before.

## Documentation diff

`docs/site/concepts/gates-and-decisions.md`, under **What you see at a gate**:

```diff
- A chat gate review has one concise evidence spine.
+ A chat gate review has one concise evidence spine. Its evidence matches the gate: backlog shows the seed outcome, boundary, and proof readiness; ideation shows the chosen direction, risk evidence, expected surface, and acceptance proof; validation shows actual results, checks, acceptance evidence, and delivery readiness. Empty result classes are omitted.
```

`docs/site/get-started/first-workflow.md`, in the example gate review:

```diff
- Test and evidence: the Stage Report records 2 DONE, 0 SKIPPED, 0 FAILED.
+ Test and evidence: both validation checks passed; the acceptance evidence is ready for delivery.
```

No command-reference or lifecycle-spec diff belongs to KRB because the command surface and authority sequence do not change.

## Acceptance criteria

**AC-1 (VALUE) - All controlled gate reviews contain only decision-relevant evidence available at their current stage, improving from the generic one-form baseline to 3/3 correct stage forms with zero fabricated rows.**
Verified by: exercise one backlog, ideation, and validation input through the existing `present-gate` path; grade the visible outputs against the table above and count wrong-stage or fabricated rows. The baseline generic template cannot correctly represent backlog and exposes all three result placeholders; the candidate must score 3/3 forms and 0 fabricated rows.

**AC-2 - Empty checklist result classes produce no presentation row.**
Verified by: a validation input with two DONE items, zero SKIPPED, and zero FAILED visibly contains the two DONE bullets and the numeric assessment, but no SKIPPED/FAILED heading, bullet, `None`, `N/A`, or equivalent placeholder.

**AC-3 - The stage-specific forms preserve gate authority.**
Verified by: grade each controlled presentation for task/stage, exactly one recommendation, the supplied bound Briefing identity and digest, and exactly one decision-effect line; in one live First Officer validation presentation, observe the review after the binding commit and before the unchanged `gate record --decision` mutation.

**AC-4 - The change does not duplicate adjacent gate tasks.**
Verified by: diff classification shows no F6C semantic-oracle cleanup, no W5 digest extraction/reliability mechanism, no XX dispatch validation, and no command/schema/recorder/runtime change; only KRB's expected surface changes.

**AC-5 - Published gate documentation describes the same stage-specific, omission-first review seen in the live presentation.**
Verified by: `mkdocs build --strict` renders the two concrete wording changes, then visual review compares their validation example with the live output; no documentation claims a new renderer or recorder.

## Expected surface

Expected implementation surface: **3 files / 35-75 changed lines**, mostly replacement text and net-neutral in mechanisms:

- `skills/present-gate/SKILL.md` — replace the generic evidence block and checklist-only assembly rule with the three stage forms and omission rule (about 25-50 changed lines).
- `docs/site/concepts/gates-and-decisions.md` — add the stage evidence map in reader-facing prose (about 8-18 changed lines).
- `docs/site/get-started/first-workflow.md` — replace the zero-class example (about 2-7 changed lines).

Tolerance is **±1 file and +35/−20 changed lines** for nearby wording or a focused live-fixture adjustment. Any second renderer, new Go production code, command/schema/stored-format/runtime change, or edit in F6C/W5/XX-owned semantics is a design reset regardless of line count.

## Test plan

1. Exercise the existing presentation path with controlled backlog, ideation, and validation bodies. Grade visible output against the stage table: required evidence present, wrong-stage evidence absent, 3/3 forms correct, zero fabricated rows. This is the direct AC-1 proof, not a prose-presence test.

2. In the validation case, provide two DONE, zero SKIPPED, and zero FAILED results. Require two DONE bullets and `Assessment: 2 done, 0 skipped, 0 failed`, with no empty-class row or placeholder (AC-2).

3. Run `go test ./internal/contractlint` to preserve skill structure/lazy-loading and `mkdocs build --strict` for the published examples. Then run the repository-required `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`; Go formatting should be a no-op because KRB expects no Go edits.

4. Run one live First Officer presentation using the supported recorded-gate fixture (`SPACEDOCK_LIVE_RUNTIME=claude go test -tags=live ./internal/ensigncycle -run '^TestLiveCommonRecordedGateLifecycle$' -count=1`, with a supported live model). Retain the visible root-session review and command log. Grade stage-specific validation evidence and empty-class omission, then prove ordering: bound Briefing commit before presentation, presentation before decision mutation, and the same `gate record --decision` recorder after it (AC-3).

5. Audit the final diff by ownership and line count (AC-4), then render and visually inspect the two site pages (AC-5). Add no semantic prose-search test, parser, second presentation renderer, or new live lane.

## Stage Report: ideation

- DONE: Map backlog, ideation, and validation evidence to one existing presentation path, with no empty result-class rows or second renderer.
  The proposed stage table keeps `skills/present-gate/SKILL.md` as the sole renderer and defines stage-available evidence plus omission semantics.
- DONE: Prove the proposed form preserves the bound Briefing, digest, recommendation, decision effect, and recorder authority.
  The common spine and AC-3 preserve all five facts and retain `${SPACEDOCK_BIN:-spacedock} gate record --decision` as the sole recorder.
- DONE: Separate KRB from F6C, W5, and XX ownership, and specify the expected files, LOC, focused checks, and one live presentation proof.
  The ownership section, 3-file/35-75-line estimate, focused checks, and recorded-gate live command set the gated boundary.

### Summary

Ideation narrows KRB to evidence selection inside the existing `present-gate` presenter: backlog, ideation, and validation each expose the evidence that exists, and validation omits empty result classes while retaining numeric totals. The design preserves the bound Briefing and sole recorder path, documents the visible change, and leaves prose-oracle cleanup, exact-digest reliability, and zero-item dispatch semantics with F6C, W5, and XX respectively.

## Stage Report: implementation

- DONE: Update only the existing present-gate renderer so backlog, ideation, and validation show stage-available evidence; validation keeps numeric totals but omits every empty result class and placeholder.
  Commit `a236c0074` changes the one presenter and two approved docs; one controlled live presenter invocation scored 3/3 forms and zero fabricated rows, and would fail if backlog/ideation emitted validation rows or validation emitted an empty class/placeholder.
- DONE: Preserve the bound Briefing/digest, exactly one recommendation and decision effect, and gate record as sole recorder while excluding all F6C, W5, XX, command, schema, and runtime semantics.
  All three controlled forms contained exactly one `Recommend`, `Reviewed snapshot`, and `Decision`; the serialized recorded-gate lane proved binding/presentation/unchanged-recorder ordering, and a duplicate/missing spine line or different recorder would fail the grade.
- DONE: Stay within the approved three-file surface and verify controlled 3/3 presentations with zero fabricated rows, contract lint, strict docs build, one live ordering proof, formatting, and both repository test suites.
  Scope was 3 files/22 changed lines; contractlint, pinned `mkdocs build --strict`, gofmt, `go test ./...`, `go test ./... -race`, controlled 3/3, and the serialized Claude live lane passed. Initial live finding was `successor dispatch build attempts/successes = 2/2, want 1/1`; output reported concurrent `/tmp/spacedock-dispatch` clobber, so the FO classified it Material to `value-ac[AC-3]`, outside KRB ownership, and authorized decline candidate fix plus one serialized rerun with candidate bytes/HEAD unchanged.

### Summary

The existing presenter now selects backlog seed evidence, ideation design/proof evidence, or validation execution evidence and omits unavailable or empty classes. The published pages match the validation form, and commit `a236c0074` contains no adjacent command, schema, recorder, or runtime mechanism. Strict rendering passed; in-app visual inspection was unavailable because this session exposed no browser backend, so that limitation remains explicit for independent validation.

## Stage Report: validation

- DONE: Independently verify the three-file candidate produces correct backlog, ideation, and validation evidence forms (3/3) with zero fabricated or empty result-class rows while retaining numeric validation totals.
  A fresh grader over the retained `present-gate` invocation scored 3/3, required two DONE bullets and `Assessment: 2 done, 0 skipped, 0 failed`, and would fail on wrong-stage evidence, an empty class, `None`, `N/A`, or an empty findings block (AC-1/AC-2).
- DONE: Verify the bound Briefing/digest, exactly one recommendation and decision effect, unchanged gate-record authority, and the serialized live binding-to-presentation-to-recorder ordering; confirm the authorized rerun left candidate bytes and HEAD unchanged.
  Each controlled form had one snapshot, `Recommend`, and `Decision`; the retained live transcript records the FO-authorized serialized rerun passing in 216.43s after the concurrent-pointer finding, while candidate commit `a236c0074` and tree `6f62b9c1` remain unchanged (AC-3).
- DONE: Cross-check AC-1 through AC-5, diff ownership and line tolerance, contract lint, strict docs build plus visual review where available, formatting, full tests, and race tests; reject any F6C/W5/XX, command, schema, recorder, or runtime semantic change.
  AC-1–AC-5 pass: 3 approved files/22 changed lines, no adjacent semantic surface, gofmt no-op, contractlint, strict MkDocs, `go test ./...`, and `go test ./... -race` all passed; an empty-class or authority/surface drift would fail the direct grader, lifecycle assertions, or diff audit.
- SKIPPED: Visual review of the two rendered documentation pages.
  `mkdocs build --strict` rendered both pages, but browser discovery returned no available backend; no source-text inspection is claimed as visual evidence.

### Summary

Validation recommends PASSED: all five acceptance criteria have executable or retained durable evidence, the candidate remains commit `a236c0074`, and no material, deferred-risk, or polish finding remains. The earlier `2/2` dispatch-attempt result was an external concurrent `/tmp/spacedock-dispatch` collision, not a KRB candidate defect; the First Officer authorized no candidate fix and one serialized rerun, which passed.
