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
        - id: gate:krbaeb3resfpbh1qvnb65krf:validation
          stage: validation
          attempts:
            - id: gate-attempt:krbaeb3resfpbh1qvnb65krf-validation-1
              briefing:
                id: briefing:krbaeb3resfpbh1qvnb65krf:validation:attempt-1:revision-1
                digest: sha256:05c661f7caa79f20a97c323cdbfba0f6de346d3d8bcef17924ac6d95241fba0e
                request-digest: sha256:0cccb28015e6e7fc60090c7ec8035c4195ef8d5716b1a28a40b0dfb92f41a113
                room-ref: ./review/validation/briefing-1
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

## Stage Report: validation (cycle 2)

- DONE: Independently verify the three-file candidate produces correct backlog, ideation, and validation evidence forms (3/3) with zero fabricated or empty result-class rows while retaining numeric validation totals.
  The retained direct grader remains 3/3 with zero fabricated rows and preserved `Assessment: 2 done, 0 skipped, 0 failed`; no candidate bytes or evidence were rerun or changed (AC-1/AC-2).
- DONE: Verify the bound Briefing/digest, exactly one recommendation and decision effect, unchanged gate-record authority, and the serialized live binding-to-presentation-to-recorder ordering; confirm the authorized rerun left candidate bytes and HEAD unchanged.
  The retained controlled and live evidence still proves the common spine and serialized authority sequence; candidate commit `a236c0074` and tree `6f62b9c1` remain unchanged (AC-3).
- FAILED: Cross-check AC-1 through AC-5, diff ownership and line tolerance, contract lint, strict docs build plus visual review where available, formatting, full tests, and race tests; reject any F6C/W5/XX, command, schema, recorder, or runtime semantic change.
  AC-4 passes by final diff classification: no F6C semantic-oracle cleanup, no W5 digest mechanism, no XX dispatch validation, and no command/schema/recorder/runtime change; AC-5 remains unmet because neither worker nor root session had an available browser backend for the required visual review.
- SKIPPED: Visual review of the two rendered documentation pages.
  Both browser-runtime checks listed zero available browsers; strict build and source inspection are not claimed as visual evidence.

### Summary

This cycle supersedes the prior report's AC-5-pass and PASSED statements. Validation recommends REJECTED for a material evidence defect at the AC-5 visual-review observation boundary; AC-1 through AC-4 pass, and no candidate product defect or candidate mutation is claimed.

## Review-finding disposition

### AC-5 unavailable visual observation

- Exact finding: AC-5 requires visual review of the two rendered documentation pages, but both the validator and root First Officer session observed zero available browser backends.
- Released user and normal workflow: the supported documentation validation path requires the declared visual review before this task can ship.
- Observable harm: validation cannot prove that the changed stage-specific and empty-row presentation renders correctly, so release evidence is incomplete even though no product defect is established.
- Authority: value-ac[AC-5]: strict docs build and visual review must establish correct rendered documentation presentation.
- Trigger evidence: both browser-runtime probes returned an empty browser list; strict build and source inspection are not visual evidence.
- Classification: Material evidence defect; ownership is outside KRB's approved three-file candidate surface.
- First Officer authorization: Needs decision; disposition route for decision. Candidate mutation and all test, browser, and live reruns remain forbidden pending Captain review.

## Stage Report: validation (cycle 3)

- DONE: Preserve the browser observation finding with all four workflow evidence fields and distinguish the unavailable observation boundary from a candidate product defect.
  The disposition above records released workflow, observable harm, `value-ac[AC-5]`, and the two empty-browser-list probes; it classifies an unavailable evidence boundary, not failed rendered behavior.
- DONE: Record the First Officer authorization as Needs decision with disposition route for decision; keep candidate commit a236c0074 and its bytes unchanged.
  The authorized hold is recorded verbatim, and candidate commit `a236c0074` / tree `6f62b9c1` remain unchanged with no candidate write or rerun.
- DONE: Append a replacement validation report that carries AC-1 through AC-4 evidence, names AC-5 as unmet, and performs no test, browser, or live rerun.
  Retained evidence remains: AC-1/AC-2 direct grader 3/3 with zero fabricated rows and numeric totals; AC-3 one snapshot/recommendation/decision plus serialized authority; AC-4 no F6C cleanup, W5 mechanism, XX validation, or command/schema/recorder/runtime change; AC-5 is unmet.

### Summary

Validation recommends REJECTED and routes the Material AC-5 evidence defect to the Captain as Needs decision. AC-1 through AC-4 retain their prior evidence, while no product defect, candidate mutation, or new execution is asserted.

## Review-finding disposition: Captain wording

- Exact finding: `A chat gate review has one concise evidence spine` is vague AI-style prose. Also inspect the two KRB documentation changes for the same defect, including `acceptance evidence is ready for delivery`. The Captain expands the finding to the 11-item presentation-discipline list in `skills/present-gate/SKILL.md`: "and i really don't like 11 step crap. if we are here, might as well clean it up."
- Released user and normal workflow: readers encounter both changed sentences in the published gate concept and first-workflow tutorial while learning what a gate review shows.
- Observable harm: the concept sentence hides the concrete evidence shown at each stage, the tutorial sentence does not name the validation checks or their observed result, and the presenter buries the rendering contract in 11 accumulated rules that repeat its template and mix evidence selection with style and lifecycle reminders; readers and First Officers cannot recover the promised stage-specific, omission-first behavior directly.
- Authority: value-ac[AC-5]: published gate documentation must describe the same stage-specific, omission-first review seen in the live presentation.
- Trigger evidence: commit `a236c0074` adds both quoted phrases and edits an 11-item `Captain-facing assembly rules` list; direct inspection shows the docs use abstract claims while the skill repeats lede, direction, report, findings, recommendation, bounce-back, format, worktree, length, label, and verification obligations after already presenting the template.
- Captain ruling: replace abstract phrases with direct stage-by-stage language that states what each review shows and that empty result groups are omitted; clean up the 11-step list while here.
- Proposed materiality: Material.
- Proposed task ownership: KRB's current three-file candidate; all affected prose is in the approved presenter and two documentation files and requires no F6C, W5, XX, command, schema, recorder, or runtime change.
- Proposed disposition: replace the 11-item list with the existing short presentation template plus a few direct rules: select only stage-available evidence, omit empty result/finding groups while retaining validation totals, preserve workflow-owned finding labels without classifying, emit exactly one recommendation/snapshot/decision effect, and keep `gate record --decision` as sole recorder. Rewrite the two documentation sentences to state the backlog, ideation, and validation evidence directly and name the passed validation checks; remove abstract prose across only these three files after distinct First Officer authorization. Candidate bytes, HEAD, and test state remain unchanged pending that authorization.

## Design reset: workflow-owned gate evidence

### Captain correction

The Captain supersedes the hardcoded stage table and the prior three-file fix proposal: "those per-stage thing should be defined in workflow. i don't understand how present-gate can generalize". `present-gate` may carry generic hints, but the current workflow stage definition is authoritative: read its declared inputs, outputs, Good/Bad criteria, and transition; show only evidence needed for that decision; never infer content from a stage name, invent missing evidence, or override an explicit workflow gate-content definition.

### Smallest declarative mechanism

- Keep stage meaning in the existing `### <stage>` subsection of the workflow README. A gated stage may add one optional `- **Gate content:** ...` instruction that names the evidence the Captain needs for that gate; this is workflow prose returned verbatim by the existing `dispatch show-stage-def` command, not a new YAML field, parser, stored form, or renderer schema.
- Add the current stage definition to `«gate.assemble-verdict»` inputs using the existing stage-def surface. If `Gate content` is present, it is the selection authority. If absent, select the minimum decision evidence from that stage's Inputs, Outputs, Good/Bad criteria, and declared advance/feedback transition. The checklist and AC scan remain evidence sources, not universal presentation sections.
- Keep `present-gate` generic. Its hints may distinguish a gate before direction selection, after direction selection, or after execution, but may not name or require any development stage. Render only non-empty workflow-requested evidence; omit an unavailable section rather than inserting `None`, `N/A`, a zero-result row, or an invented fact.
- Preserve the common gate authority unchanged: task label and current stage, bound Briefing identity/digest, workflow-owned reviewer finding labels without presenter classification, exactly one recommendation, exactly one decision effect derived from the declared transition, and `${SPACEDOCK_BIN:-spacedock} gate record --decision` as sole recorder.

### Migration surface

- Replace the hardcoded stage names and 11-item accumulated list in `skills/present-gate/SKILL.md` with one short generic template and the direct rules above.
- Extend `«gate.assemble-verdict»` in `skills/first-officer/references/first-officer-shared-core.md` to read the current workflow stage definition before rendering; reuse `dispatch show-stage-def`, checklist, AC scan, boot stage order, and `feedback-to` rather than adding a command or parser.
- Add explicit `Gate content` instructions to the gated stages in `docs/dev/README.md`, and to gated stages in the built-in development, refinement, and experiment commission templates. Custom commissioning guidance should ask what evidence makes each gate decidable and write that answer into the stage subsection.
- Rewrite the two KRB user-doc changes to explain that each workflow defines its gate evidence and that empty evidence groups are omitted; the example names the checks that passed. Do not teach a frontmatter field or dev-specific universal stage map.
- Candidate commit `a236c0074` is superseded by this design reset. No candidate byte or HEAD changes are authorized yet; expected file/line scope must be re-estimated before implementation, with any generated template fixtures included only when their existing propagation checks require them.

### Proof plan

1. Exercise the existing presenter with controlled workflows whose gated stages use non-development names. Grade that each visible review follows its own explicit `Gate content`, contains no development-stage vocabulary, and omits every unavailable or empty group.
2. Exercise an explicit-override case where Inputs/Outputs suggest extra material but `Gate content` requests a narrower metric; fail if the presenter adds the extra material. Exercise a fallback case with no `Gate content`; fail unless the review uses only the declared Inputs, Outputs, Good/Bad criteria, and transition.
3. Across the controlled reviews, require the supplied Briefing/digest, workflow-owned finding labels, exactly one recommendation, exactly one decision effect, and the unchanged sole recorder. A missing/duplicate common fact or presenter-authored finding class fails.
4. Use one live recorded gate with a non-development stage name to observe stage definition read before presentation and presentation before `gate record --decision`; grade output semantics and durable command ordering, not instruction-file substrings.
5. After authorized implementation, run contract lint, strict docs build plus available visual review, formatting, both repository suites, and final diff ownership. Add no semantic prose-presence oracle.

### F6C boundary

KRB owns the workflow-defined gate-content mechanism, generic evidence selection, built-in stage declarations needed to adopt it, and the two KRB user-doc sentences. F6C retains removal/replacement of semantic prose-grep assertions and the wider final gate/help documentation reconciliation; KRB adds no committed wording oracle and does not absorb F6C's inventory. W5 digest reliability, XX zero-item dispatch validation, command grammar, stored gate schema, recorder semantics, and runtime routing remain excluded.

## Re-estimate: workflow-owned Gate content

This re-estimate supersedes the earlier three-file/35-75-line surface and narrows the reset draft to the smallest end-to-end adoption. `present-gate` itself reads the current stage subsection through the existing `dispatch show-stage-def` surface, so `skills/first-officer/references/first-officer-shared-core.md`, Go code, YAML frontmatter, parser logic, and renderer schemas do not change.

### Revised mechanism

Each gated `### <stage>` subsection may declare one direct `- **Gate content:** ...` instruction. The presenter treats that instruction as authoritative; without it, the presenter uses only the subsection's Inputs, Outputs, Good/Bad criteria, and the workflow-declared advance/feedback transition. Generic hints distinguish gates before direction selection, after direction selection, and after execution without assigning meaning to a stage name. Missing evidence and empty result/finding groups are omitted, never replaced with placeholders or invented facts.

### Revised acceptance criteria

**AC-1 (VALUE) - Gate reviews follow their workflow's stage definition, improving from a hardcoded development-stage table to correct declared evidence for every controlled non-development gate with zero stage-name inference or fabricated rows.**
Verified by: exercise explicit Gate content across controlled refinement and experiment gates with non-development names; grade every visible row against the declaration. Renaming a stage without changing its definition must not change the evidence, while changing the Gate content must change the review.

**AC-2 - Explicit Gate content overrides generic hints, and an undeclared Gate content falls back only to declared stage criteria and transition.**
Verified by: one conflict case requests a narrow metric while Inputs/Outputs mention additional material, and one fallback case has no Gate content. The first fails if extra material appears; the second fails if output cannot be traced to Inputs, Outputs, Good/Bad, or transition.

**AC-3 - Empty evidence is absent while the common gate authority remains singular and unchanged.**
Verified by: controlled reviews with missing/zero evidence contain no empty heading, `None`, `N/A`, or zero-result row and contain exactly one task/stage, supplied Briefing identity/digest, recommendation, and decision effect; workflow finding labels remain exact and `${SPACEDOCK_BIN:-spacedock} gate record --decision` remains the sole recorder.

**AC-4 - Commission and refit propagate workflow-owned Gate content without rewriting old fixture truth.**
Verified by: commission one workflow from each built-in template and one custom/variant shape, then run refit Phase 3b against `skills/integration/testdata/refit-content-propagation/site-workflow/`. Generated gated stage subsections contain their template/mission-specific Gate content, and the refit diff exposes the additive declarations while the legacy fixture README remains byte-unchanged.

**AC-5 - Published documentation describes workflow-owned evidence directly and KRB does not absorb F6C or adjacent semantics.**
Verified by: strict docs rendering and visual review compare the user pages with controlled output; final diff classification shows no semantic prose-grep cleanup, F6C inventory, W5 digest mechanism, XX dispatch validation, command grammar, YAML/stored schema, parser, recorder, or runtime change.

### Revised expected surface

Expected implementation surface: **9 files / 55-105 changed lines** (insertions plus deletions):

- `skills/present-gate/SKILL.md` - remove the hardcoded stage table and 11-item list; add one generic template, stage-definition fetch, phase hints, omission, finding-label, singular-authority, and recorder rules (about 25-50 changed lines).
- `skills/commission/SKILL.md` - require every custom or variant gated stage definition to state the evidence that makes its decision possible (about 3-7 lines).
- `skills/commission/references/templates/development.md` - add Gate content to its three gated stage subsections (about 3-8 lines).
- `skills/commission/references/templates/refinement.md` - add Gate content to its gated review subsection; variants inherit the commission rule (about 1-4 lines).
- `skills/commission/references/templates/experiment.md` - add Gate content to its three gated stage subsections (about 3-8 lines).
- `docs/dev/README.md` - add explicit Gate content to this workflow's three gated stage definitions (about 3-8 lines).
- `skills/integration/testdata/refit-content-propagation/README.md` - name Gate content as an expected additive refit hunk (about 2-5 lines). Its legacy `site-workflow/README.md` deliberately remains unchanged; entity-label and scaffold fixtures are fallback workflows, not propagation fixtures.
- `docs/site/concepts/gates-and-decisions.md` - replace the abstract/hardcoded sentence with direct workflow-owned behavior (about 2-7 lines).
- `docs/site/get-started/first-workflow.md` - name the concrete passed checks and omit empty groups in the example (about 2-5 lines).

Tolerance is **plus or minus 1 file and +35/-20 changed lines** for existing template-propagation machinery discovered during implementation. Any Go production edit, new YAML/frontmatter field, parser, renderer schema, hardcoded stage semantic, second recorder, or F6C semantic-oracle cleanup is a design reset regardless of line count.

### Revised test plan

1. Run the AC-1/AC-2 controlled presenter matrix across development, refinement, experiment, and a custom shape. Grade semantic output rather than searching instruction prose; include rename invariance, explicit-override, fallback, and absent-evidence controls.
2. Run the AC-3 common-spine grader with empty classes/findings and a workflow-owned finding-label order. Fail on a placeholder, invented section, missing/duplicate authority fact, reclassification, or alternate recorder.
3. Drive commission for all three built-in templates plus one custom/variant gate, then drive refit Phase 3b against the retained old site-workflow fixture. Compare generated stage subsections and refit diff; do not update the old fixture to make the proof pass.
4. Run contract lint, strict MkDocs plus available visual review, `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`. Run the live gate journeys required by the shared-skill path mapping, including one visible non-development gate that proves stage-definition read before presentation and presentation before the unchanged recorder.
5. Audit final paths and line count against the estimate and F6C boundary. Add no committed semantic prose-presence test.

## Stage Report: implementation (cycle 2)

- DONE: Resolve the smallest workflow-owned Gate content mechanism without hardcoded stage names, new YAML, a parser, or a renderer schema.
  The revised mechanism uses one optional instruction in the existing stage subsection and the shipped stage-definition read; explicit content wins and fallback is bounded to declared criteria/transition.
- DONE: Re-estimate the exact files and changed lines, including only required built-in templates and propagation fixtures; state the F6C boundary and any tolerance decision.
  The new estimate is 9 files/55-105 changed lines with plus-or-minus-1-file and +35/-20-line tolerance; all three canonical templates and only the propagation descriptor change, while its old workflow fixture remains byte-stable.
- DONE: Update the durable design, expected surface, acceptance criteria, and test plan; keep candidate a236c0074 and all product bytes unchanged and run no tests.
  The sections above replace the obsolete design and name falsifiable matrix, commission/refit, docs, repository, live, and ownership proof; code HEAD remains `a236c0074` with a clean worktree and no test or live invocation in this cycle.

### Summary

The design now generalizes through workflow-owned Gate content carried by existing stage subsections and the existing read surface. The re-estimate identifies nine product/documentation/template files, retains the old refit workflow as the propagation control, and excludes F6C plus all command/schema/parser/recorder/runtime changes. Candidate `a236c0074` remains untouched pending a separately authorized implementation.

## Stage Report: implementation (cycle 3)

- DONE: Add a focused behavioral test before product edits for non-development stage names, explicit Gate content override, bounded fallback, rename invariance, and omitted empty evidence.
  `skills/integration/gate_content_live_test.sh` failed on `a236c0074` by leaking catalog evidence and zero classes; it now passes explicit `orbit-audit`, renamed `signal-check`, and no-instruction fallback, and fails if the stage fetch, override, bound, or omission behavior regresses.
- DONE: Implement the workflow-owned Gate content design across the re-estimated surface; keep present-gate stage-name agnostic, remove the 11-item list and abstract prose, and preserve all common gate authority.
  Commits `0d88168a8` and `6b36ecfe5` deliver 10 files/131 changed lines: generic presentation rules, three templates, dev workflow declarations, commission/refit propagation, direct user docs, and no command/schema/parser/recorder/runtime change.
- DONE: Run focused presenter, commission/refit, contract, strict documentation, formatting, full, race, and required live checks; compare generated documentation output with controlled gate output without requiring a browser-only observation.
  The three-case Codex presenter matrix, custom commission, refit diff with old hash `5565f9b3…` unchanged, built-in stage reads, contract/integration tests, `mkdocs build --strict`, `gofmt`, `go test ./...`, and `go test ./... -race` passed; rendered docs and live output both omit empty groups. Claude gate-guardrail was attempted but skipped for missing local auth, and Pi reported its existing registered TODO.

### Summary

Gate reviews now read workflow-owned `Gate content` and use a bounded declared-stage fallback without stage-name semantics. The presenter keeps one Briefing, recommendation, decision effect, finding-label authority, and recorder while omitting empty evidence; commission, refit, built-in templates, and user docs carry the same behavior.

## Stage Report: validation (cycle 4)

- DONE: Independently reproduce the explicit-content, renamed-stage, bounded-fallback, and empty-evidence matrix without trusting the implementation report.
  The current-checkout Codex presenter passed explicit `orbit-audit`, the definition-preserving rename `signal-check`, and the no-`Gate content` fallback. The explicit cases showed only the requested 4.2 ms drift and transition, never the conflicting raw-frame/catalog material; all cases omitted `SKIPPED`, `FAILED`, `None`, `N/A`, and zero-result rows. The stream recorded `dispatch show-stage-def` before presentation (AC-1/AC-2).
- DONE: Verify the common gate spine and workflow finding authority with a separate controlled review.
  A `release-check` workflow rendered exactly one task/stage, `briefing:spine:1`, `sha256:def456`, recommendation, and concrete decision effect; it preserved the non-empty `Urgent` label, omitted empty `Advisory`, fetched the stage definition, and invoked no alternate recorder. `${SPACEDOCK_BIN:-spacedock} gate record --decision` remains the sole declared recorder (AC-3).
- DONE: Drive all three built-in commission shapes, one custom variant, and refit Phase 3b while retaining the old fixture as the control.
  Fresh generated READMEs contained mission-specific Gate content at development's three gates, experiment's three gates, refinement's review gate, and the custom outreach variant's `review` and `watching` gates. Refit of a throwaway copy exposed additive Gate content at backlog, ideation, and validation. The repository control stayed byte-identical at Git blob `5565f9b3b6c183aa04c2b60969139daddcc24a1e` at both base `a929fcb60` and candidate HEAD (AC-4).
- DONE: Verify documentation output, exact candidate scope, formatting, and all required offline suites.
  Candidate `6b36ecfe5eef975741009fd4b8a6f7884dc678d1` / tree `9d9de65251` is clean and measures 10 files/131 changed lines from `main`. Contract and integration tests, pinned strict MkDocs, generated-visible-HTML comparison, exact `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` passed. The retired hardcoded evidence table, 11-item list, and abstract documentation phrases are absent; behavior was proved by rendered output, not by those absence checks. No browser-only observation is claimed or required by this reset's assigned documentation proof.
- DONE: Classify the final diff against AC-5 and the approved boundary.
  The final diff contains direct workflow-owned documentation and no F6C semantic-oracle cleanup, no F6C inventory, no W5 digest mechanism, no XX dispatch validation, and no command grammar, YAML/stored schema, parser, recorder, or runtime change. AC-5 therefore passes through strict rendering, generated-output comparison, and this final path classification.
- FAILED: Obtain passing evidence for every required live lane selected by the changed shared `present-gate` path.
  `TestLiveCommonGateGuardrail` passed on Claude Sonnet in 129.32s and Claude Opus in 166.20s. The same canonical journey explicitly skipped on Codex as `TODO(xp6c9qfe7y4wwp46enc3f85n): codex/gate-guardrail lacks passing live evidence`; Pi reported the same owned TODO. A zero-exit Go wrapper around a skipped required journey is not green evidence, so the required Codex lane remains unmet.

### Summary

Validation recommends REJECTED for one material verification defect: the required Codex `gate-guardrail` live journey has only an owned TODO skip. AC-1 through AC-5 otherwise pass, including controlled workflow-owned rendering, common authority, commission/refit propagation, strict generated documentation, and the explicit adjacent-surface exclusion. Candidate bytes and HEAD remain unchanged.

## Review-finding disposition: required Codex gate-guardrail evidence

- Exact finding: the shared `present-gate` path requires the canonical Codex `TestLiveCommonGateGuardrail` journey, but the executable registry skips it as `TODO(xp6c9qfe7y4wwp46enc3f85n)` and states that Codex lacks passing live evidence.
- Released user and normal workflow: pull requests exercise the current-checkout shared common journeys through both Claude Sonnet and Codex; a skipped required run is explicitly not green.
- Observable harm: validation cannot establish that Codex preserves the gate at the bound Briefing while presenting the revised shared presenter, even though the controlled Codex renderer and both Claude live lanes pass.
- Authority: value-ac[AC-1]: the workflow-owned presenter behavior must hold for the shared runtime surface; contract[docs/runtime-live-ci.md#runtime-live-ci]: Codex common journeys are current-checkout live evidence.
- Trigger evidence: `SPACEDOCK_LIVE_RUNTIME=codex go test -tags live -count=1 -run '^TestLiveCommonGateGuardrail$' ./internal/ensigncycle -v` emitted the owned TODO and `--- SKIP`; the Pi invocation emitted the same target-scoped TODO.
- Classification: Material evidence defect; the candidate does not own the target-scoped runtime TODO or its separate remediation.
- First Officer direction: preserve and route this new finding through the workflow checkpoint without candidate mutation. Disposition is route for decision; no KRB candidate edit or live-evidence substitution is authorized.

## Stage Report: validation (cycle 5)

- DONE: Record the First Officer's disposition for the required Codex gate-guardrail evidence defect without changing or re-exercising the candidate.
  Authorization is HOLD. KRB remains at candidate `6b36ecfe5`; no candidate bytes, HEAD, test result, or live evidence were changed, rerun, or substituted in this cycle.
- DONE: Assign the runtime evidence repair to its existing owner and keep it outside KRB's approved boundary.
  The material Codex `gate-guardrail` defect is owned by `xp6` (`restore-live-evidence-after-completed-repairs`) in the test-behavior-completeness sprint. KRB does not absorb the runtime TODO repair.
- DONE: State the exact re-entry condition.
  KRB is held for revalidation after xp6 supplies a real passing Codex `TestLiveCommonGateGuardrail` journey. An owned TODO skip or other substituted evidence does not satisfy that condition.

### Summary

Validation remains REJECTED and is now held by First Officer authorization. Candidate `6b36ecfe5` is unchanged; revalidation waits for xp6's real passing Codex gate-guardrail evidence.

## Stage Report: implementation (cycle 4)

- DONE: Remove only the Codex `gate-guardrail` TODO authorized by the Captain while preserving the Pi TODO and every other registry entry.
  Commit `4e3be373f` changes only `internal/ensigncycle/shared_live_runner_test.go`: it removes `liveTODO("codex", "xp6c9qfe7y4wwp46enc3f85n")` from `TestLiveCommonGateGuardrail` and leaves `liveTODO("pi", "xp6c9qfe7y4wwp46enc3f85n")` in place. The complete candidate is now 11 files/133 changed lines from `main`, within the reset tolerance.
- DONE: Run the focused registry and offline gate-authority controls before the live invocation.
  `TestRuntimeLiveRegistryReconciliation`, `TestAssertGateHeld`, `TestAssertGateHeldAcceptsPreparedFixtureBinding`, and `TestGateGuardrailNegativeBrokenStateTransition` passed after the one-line change; the same focused controls also passed before the edit.
- FAILED: Execute and pass the one authorized authenticated Codex `TestLiveCommonGateGuardrail` run through the real durable `assertGateHeld` check.
  The sole invocation was `SPACEDOCK_CODEX_LIVE_REQUIRED=1 SPACEDOCK_LIVE_RUNTIME=codex go test -tags live -count=1 -timeout 40m -run '^TestLiveCommonGateGuardrail$' ./internal/ensigncycle -v`. It did not skip, but failed before the scenario and assertion ran with exact evidence: `shared_live_runner_test.go:66: OPENAI_API_KEY is required for the approval-gated codex-live lane`; Go reported `--- FAIL: TestLiveCommonGateGuardrail (0.00s)`. Per the Captain's authorization, no rerun, runtime repair, harness change, fallback, alternate evidence, or further candidate mutation followed.

### Summary

The authorized Codex TODO removal is committed at `4e3be373f`, and all focused deterministic controls pass. The single required-auth live invocation could not authenticate, so it produced no real `assertGateHeld` evidence. Implementation stops with that exact failure preserved and the candidate otherwise unchanged.

## Stage Report: validation (cycle 6)

- DONE: Verify Runtime Live E2E run `31348183583` against the exact candidate rather than trusting its aggregate conclusion.
  GitHub records workflow `Runtime Live E2E`, event `workflow_dispatch`, head SHA `4e3be373f56be96f340d6322d07033865e99f31b`, and overall `success`. Jobs `offline` (`93333965496`), `claude-live (sonnet, claude-sonnet-5, max, CI-E2E)` (`93334128215`), and `codex-live` (`93334128222`) all completed successfully. The Codex checkout log fetches that exact SHA, requires `OPENAI_API_KEY`, sets `SPACEDOCK_CODEX_LIVE_REQUIRED=1`, and installs the current checkout as the Codex plugin.
- DONE: Inspect the target-level artifacts and prove that the formerly missing journey genuinely ran.
  Artifact `runtime-live-e2e-codex-live` (`9048334088`) records `TestLiveCommonGateGuardrail` as `run`, then `--- PASS` and `Action:"pass"` in 67.96s; its journey metric records `scenario_id: gate-guardrail`, runtime `codex`, `outcome.status: passed`, and run ID `31348183583`. The Sonnet detail artifact independently records the same target passing in 170.81s. Neither target was a TODO skip, and the durable held-gate assertion therefore clears the prior missing-auth evidence defect.
- FAILED: Regrade the retained visible presentations against the Captain-corrected bounded-fallback and empty-finding requirements.
  The exact Codex target fetches a validation stage definition containing only `Validate and present the retained package.`: it has no `Gate content`, Inputs, Outputs, Good/Bad criteria, or evidence declaration. Its visible `Evidence for validation` nevertheless adds checklist, AC-scan, and committed-review rows, then emits `- No material findings.` The retained Sonnet presentation likewise emits `- Findings: no material findings.` These are unsupported fallback rows and explicit empty finding rows; they violate AC-2's declared-criteria/transition bound and AC-3's requirement to omit empty result/finding groups. The green live assertion checks the durable held boundary, not these presentation semantics, so its pass cannot override the observed AC failure.
- DONE: Reconcile the remaining ACs, Captain wording correction, and final candidate classification without changing product bytes.
  Retained independent evidence for AC-1 and AC-4 remains valid: the non-development explicit/rename matrix follows workflow Gate content, all commission/refit shapes propagate it, and the legacy fixture stays at blob `5565f9b3b6c183aa04c2b60969139daddcc24a1e`. AC-5's direct docs, strict generated rendering, and Captain-requested removal of the abstract prose/11-item list remain valid. Exact candidate `4e3be373f` / tree `cd208b728b` is clean at 11 files/133 changed lines; its only delta from `6b36ecfe5` is the Captain-authorized test-only Codex TODO removal. The final product diff still contains no F6C semantic-oracle cleanup or inventory, W5 digest mechanism, XX dispatch validation, command/schema/parser/recorder change, or product runtime-routing change.

### Summary

Validation recommends REJECTED. Run `31348183583` supplies the real authenticated Codex evidence that was previously missing, but the evidence itself refutes AC-2 and AC-3: the presenter leaks undeclared fallback material and renders an empty findings row on both Codex and Sonnet. The live harness's durable-state grade passes while leaving that visible semantic defect ungraded. Candidate `4e3be373f` remains unchanged.

## Review-finding disposition: live fallback leakage and empty findings

- Exact finding: in exact-head run `31348183583`, Codex presents checklist/AC/review rows although the fetched fallback stage declares none of the permitted evidence fields, and both Codex and Sonnet present a `no material findings` row instead of omitting the empty finding group.
- Released user and normal workflow: workflow-owned `Gate content` is authoritative; absent it, fallback is limited to declared Inputs, Outputs, Good/Bad criteria, and transition, while missing and empty evidence/finding groups are omitted.
- Observable harm: a Captain receives evidence the workflow did not request and an empty findings row the corrected design promised to suppress; an aggregate green live job can conceal the regression because `assertGateHeld` grades only durable gate state.
- Authority: value-ac[AC-2]: undeclared fallback is bounded to declared criteria and transition; value-ac[AC-3]: empty evidence and finding groups are absent.
- Trigger evidence: Codex artifact `9048334088`, `codex-shared-scenarios-detail.jsonl` lines 6-9, and `live-artifacts/codex/codex-shared-scenarios/gate-guardrail/codex-final-message.txt`; Sonnet artifact `9048510715`, `live-e2e-detail.jsonl` target pass, and its gate-guardrail `claude-final-message.txt` line 14.
- Classification: Material product-behavior defect in KRB's presenter surface, plus a test-coverage gap because the canonical live journey does not grade the visible KRB semantics.
- Proposed disposition: route to implementation through the workflow checkpoint. Preserve exact candidate `4e3be373f` until distinct authorization; do not substitute the durable-state pass for presentation evidence or absorb unrelated runtime behavior into KRB.

## Worker proposal: live fallback leakage and empty findings

- Exact finding: Runtime Live E2E run 31348183583 shows Codex fallback presentation adding checklist/AC/review rows although the fetched stage declares no Gate content, Inputs, Outputs, Good/Bad criteria, or transition; both Codex and Sonnet also render `no material findings` instead of omitting the empty findings group.
- Released user and normal workflow: the normal pull-request Runtime Live E2E executes the shipped `present-gate` skill on Codex and Sonnet at a human gate whose workflow definition supplies no fallback evidence fields.
- Observable harm: the Captain is shown undeclared checklist, acceptance-criterion, and retained-review claims as decision evidence, plus a non-finding represented as a finding; the green job conceals both visible contract failures because it grades only the durable held state.
- Affected value AC or non-negotiable boundary: `value-ac[AC-2]` requires fallback to stay within declared Inputs, Outputs, Good/Bad criteria, and transition; `value-ac[AC-3]` requires empty result and finding groups to be absent.
- Trigger evidence: run `31348183583` is exact head `4e3be373f`. Codex artifact `9048334088` shows `dispatch show-stage-def` returning only `### validation` and `Validate and present the retained package.`, then the final message adds checklist, AC-1, committed-review, and `No material findings` rows. Sonnet artifact `9048510715` adds the same classes of rows and `Findings: no material findings`. Static inspection confirms `recordedGateReadme()` declares no evidence field in the validation subsection, while `recordedGateEntity()` and `recordedGateSourceReview()` contain the leaked distractors; `runGateStopScenario()` passes only before/after state to `assertGateHeld()` and never grades `result.finalMessage`.
- Root cause: the presenter says fallback uses only declared sources but does not state the zero-source terminal behavior or forbid entity checklist, AC scan, and retained package material from becoming fallback rows. It also says to omit empty findings without defining a negative statement such as `no material findings` as empty. Both hosts filled those gaps the same way, and the live oracle did not observe the presentation.
- Proposed materiality: Material. The trigger occurred in the released PR workflow on both supported hosts and directly violates the two value criteria at the Captain-facing decision boundary.
- Proposed task ownership: KRB owns the presenter correction and the semantic grade for its canonical gate journey. The remedy changes no runtime adapter, command, schema, recorder, gate-state assertion, or xp6-owned runtime mechanism.
- Proposed disposition: `fix`.

### Smallest concrete correction

1. In `skills/present-gate/SKILL.md`, state that when `Gate content` and every permitted fallback declaration are absent, the presenter emits no evidence block or row. Checklist, AC-scan, entity report, and retained review/package data may inform recommendation judgment but cannot become fallback evidence. State that negative summaries such as `no material findings` mean the finding group is empty and must not render.
2. Add one fixture-backed presentation assertion to `runGateStopScenario()` that rejects the recorded-gate fixture's checklist, AC, and retained-review distractors plus any `no material findings` row, without changing `assertGateHeld()`. Give that assertion default-build positive and negative table cases so the semantic oracle proves red offline without a model. This is test evidence for KRB output, not a production renderer or fallback implementation.

### Proof boundary after authorization

- Before the presenter edit, add and run only the focused default-build presentation-oracle cases; require them to fail on the retained Codex and Sonnet messages and pass on a review containing only the common task/stage, Briefing, recommendation, and decision facts.
- After the single skill correction, run the focused offline oracle, registry reconciliation, existing gate-state controls, contract/integration tests, formatting, full tests, race tests, and strict docs build. Do not alter fixture stage declarations to make the check pass.
- Trigger one normal PR Runtime Live E2E at the exact correction commit. Require both Codex and Sonnet `TestLiveCommonGateGuardrail` targets to execute and pass the new presentation assertion; inspect both retained final messages to confirm zero unsupported evidence rows and zero empty findings rows. Aggregate job success without those target-level facts is insufficient. No local model run, Pi expansion, runtime repair, fallback mechanism, or repeat-for-luck run belongs to this correction.

Candidate bytes and HEAD remain exactly `4e3be373f`; no edit, test, model rerun, or CI trigger occurred during this proposal.

## Review-finding disposition: race replay budget

- Exact finding: `go test ./... -race` fails `TestSonnetTeamDeleteHangReplay` at `streamwatch_regression_test.go:129` because the transcript has not reached the required `"terminal_reason":"completed"` beat after the prior beats.
- Released user and normal workflow: repository completion requires the full race suite after every implementation, including KRB's test-only assertion change.
- Observable harm: KRB cannot establish the required all-green deterministic boundary and therefore cannot spend the separately authorized targeted Codex live run or claim the correction ready for revalidation.
- Affected value AC or non-negotiable boundary: `contract[AGENTS.md#Expected Commands]` requires `go test ./... -race` before work is complete.
- Trigger evidence: the full race suite failed after 380.631s with `TestSonnetTeamDeleteHangReplay` reporting `replay: the recording must contain the diagnosed beat "\"terminal_reason\":\"completed\"" in order after the prior beat`. One focused diagnostic, `go test -race ./internal/ensigncycle -run '^TestSonnetTeamDeleteHangReplay$' -count=1 -v`, reproduced the same failure in 0.95s. Static inspection shows the untouched test drips a 341-line recording one line per 1 ms poll but gives `expectExit` only 150 ms under race instrumentation; KRB changes none of `streamwatch_regression_test.go`, its fixture, or the stream watcher.
- Proposed materiality: Material evidence defect at the repository completion boundary; no KRB product regression is established.
- Proposed task ownership: outside KRB. The remedy belongs to the existing streamwatch replay/test owner and is not part of the presenter, canonical gate semantic assertion, or authorized three-file correction.
- Proposed disposition: `route for decision`.

The authorized KRB correction is present but uncommitted; code HEAD remains `4e3be373f`. Focused semantic controls, registry reconciliation, contract/integration tests, formatting, and `go test ./...` pass. No targeted model run, GitHub Actions run, API-key use, broad live suite, unrelated repair, candidate commit, or test rerun beyond the single focused diagnostic followed this changed evidence.

### First Officer disposition

DECLINE any KRB candidate fix for `TestSonnetTeamDeleteHangReplay`. The timing test is outside KRB and byte-identical to current `main`, where the exact focused race command passed. The First Officer authorized one isolated serial focused rerun in the KRB worktree and, on success, one serial full race suite before the already-authorized targeted Codex OAuth journey. No timing-test mutation was authorized or made.

## Stage Report: implementation (cycle 5)

- DONE: Add the fixture-backed final-message semantic assertion before changing the presenter and prove the focused deterministic test red.
  `TestAssertRecordedGatePresentation` first failed to compile because `assertRecordedGatePresentation` was undefined. Its clean positive control contains only the common gate spine; five negative controls cover checklist, AC, retained-review/check rows, and `no material findings`.
- DONE: Apply only the authorized KRB correction for the zero-source fallback and empty negative findings, then bind the semantic assertion to the canonical journey.
  Commit `7e7651351` changes four files/54 changed lines: `present-gate` now omits the evidence block when no permitted source declares evidence, forbids checklist/AC/entity/package data from becoming fallback rows, and treats `no material findings` as empty. `runGateStopScenario()` now grades `result.finalMessage` after the unchanged durable `assertGateHeld()` check. No runtime adapter, command, schema, recorder, fixture stage declaration, or gate-state assertion changed.
- DONE: Pass all focused and repository-required deterministic checks without repairing the declined timing test.
  The semantic positive/negative controls, existing gate-state controls, registry reconciliation, contract/integration tests, `gofmt -w ./cmd ./internal`, and `go test ./...` passed. The initial parallel `go test ./... -race` and one diagnostic focused run exposed the recorded race-budget finding above. After the First Officer's DECLINE and serialized-rerun authorization, the exact focused test passed in 0.57s and `go test -p 1 -race ./...` passed with `internal/ensigncycle` completing in 288.486s; `streamwatch_regression_test.go` stayed byte-identical.
- DONE: Run only the targeted local Codex gate-guardrail journey through isolated OAuth after the deterministic boundary was green.
  `env -u OPENAI_API_KEY -u SPACEDOCK_CODEX_LIVE_REQUIRED SPACEDOCK_LIVE_RUNTIME=codex go test -tags live -count=1 -timeout 40m -run '^TestLiveCommonGateGuardrail$' ./internal/ensigncycle -v` executed and passed in 110.99s. The pass includes the new final-message semantic assertion and the existing durable held-gate assertion. No API key, GitHub Actions trigger, broad live suite, Pi lane, or second model run was used.

### Summary

Candidate `7e7651351` makes an empty declared fallback visibly empty and suppresses negative empty-finding summaries. The canonical live journey now fails on the exact checklist, AC, retained-review, and empty-finding leakage observed in run `31348183583`. All authorized deterministic and targeted OAuth evidence is green; the complete candidate is clean at 14 files/179 changed lines from `main` and is ready for independent revalidation.

## Stage Report: validation (cycle 7)

- DONE: Inspect exact candidate `7e765135172422b37c1d7fadae1d6347d452ff9c` and verify the authorized correction is limited to the rejected behavior.
  The cycle-5 delta is exactly four files/54 changed lines: one `present-gate` clarification plus the canonical journey hook, fixture-backed semantic oracle, and its positive/negative unit table. It changes no runtime adapter, command, schema, recorder, fixture stage declaration, durable gate-state assertion, documentation, commission/refit propagation, or declined streamwatch timing test. The complete clean candidate is 14 files/179 changed lines at tree `f72c3def20`; the surface beyond the reset estimate is precisely the separately authorized cycle-5 correction.
- DONE: Reproduce the fixture-backed semantic boundary without running another model journey.
  `TestAssertRecordedGatePresentation` passes its clean common-spine control and rejects all five controlled leak classes: checklist/detail, AC, retained review, retained review checks, and `no material findings`. The helper is wired after the unchanged `assertGateHeld` call in `runGateStopScenario`, so a canonical live target cannot pass on durable state alone while leaking the cycle-6 rows. Existing `TestAssertGateHeld`, prepared-binding mutants, and `TestGateGuardrailNegativeBrokenStateTransition` also pass (AC-2/AC-3).
- DONE: Verify all five acceptance criteria and Captain-authorized dispositions using retained and fresh evidence.
  AC-1 retains the independent explicit/non-development rename matrix; cycle 5 changes only empty fallback handling. AC-2 now has an explicit zero-source terminal rule and a canonical negative oracle. AC-3 treats negative findings summaries as empty, retains singular common authority, and has the positive/negative oracle. AC-4 is unchanged: built-in/custom commission and refit evidence remains valid and the legacy fixture remains blob `5565f9b3b6c183aa04c2b60969139daddcc24a1e`. AC-5's direct docs, strict generated rendering, Captain wording cleanup, and boundary classification remain unchanged. The final diff still contains no F6C semantic-oracle cleanup/inventory, W5 digest mechanism, XX dispatch validation, command/schema/parser/recorder change, or product runtime-routing change.
- DONE: Reproduce the required deterministic repository boundary at the exact candidate.
  `gofmt -d ./cmd ./internal` was empty; the focused semantic and gate-state controls, `TestRuntimeLiveRegistryReconciliation`, contract/integration tests, `go test ./...`, and `go test -p 1 -race ./...` all pass. The serial race run preserves the First Officer's DECLINE of an unrelated timing-test fix; `streamwatch_regression_test.go` is byte-identical to the merge base. No product write, GitHub Actions trigger, command/help prose comparison, or model journey occurred during validation.
- DONE: Confirm the authorized live evidence now exercises both the semantic and durable assertions.
  The sole retained cycle-5 invocation ran Codex with `OPENAI_API_KEY` and `SPACEDOCK_CODEX_LIVE_REQUIRED` unset and passed `TestLiveCommonGateGuardrail` in 110.99s. Because cycle 5 wires `assertRecordedGatePresentation(result.finalMessage)` into that exact journey before the hold-log grade, this pass clears the cycle-6 visible leakage/empty-findings defect as well as the existing durable held-gate assertion. No repeat-for-luck or substitute lane is claimed.

### Summary

Validation recommends PASSED. Candidate `7e7651351` closes the exact live fallback and empty-findings defect with a deterministic semantic oracle bound to the canonical journey; every AC and authorized disposition has green, proportionate evidence. Candidate bytes and HEAD remain unchanged.
