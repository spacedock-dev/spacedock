---
title: Make review-finding disposition workflow-owned and separate in-stage review from rejection routing
status: implementation
source: "Captain boundary audit, 2026-07-30: an end-of-implementation Roborev round never enters feedback-rejection-flow, and finding classification must happen when findings arrive rather than only after rejection."
started: 2026-07-30T01:34:32Z
completed:
verdict:
score: "0.90"
worktree: .worktrees/spacedock-ensign-workflow-owned-review-finding-disposition
issue:
pr:
sprint: durable-decisions
id: rhx820qrkn6vxpday10nch36
gates:
    version: 1
    current:
        gate: gate:rhx820qrkn6vxpday10nch36:ideation
    records:
        - id: gate:rhx820qrkn6vxpday10nch36:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:rhx820qrkn6vxpday10nch36-backlog-1
              briefing:
                id: briefing:rhx820qrkn6vxpday10nch36:backlog:attempt-1:revision-1
                digest: sha256:f5c4251b4617310097bca5e1d14a03f5d177df80b0e452751b510e9b8eae30de
                digest-domain: canonical-bytes
                request-digest: sha256:50ec6a9d980e9ca88e45a5c22e5f7e394b6c0d6a5e3044a15a4dbf617515fcc7
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:rhx820qrkn6vxpday10nch36:backlog:1
                briefing: briefing:rhx820qrkn6vxpday10nch36:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-07-30T01:33:24.438725Z"
                decision: approve
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
        - id: gate:rhx820qrkn6vxpday10nch36:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:rhx820qrkn6vxpday10nch36-ideation-1
              briefing:
                id: briefing:rhx820qrkn6vxpday10nch36:ideation:attempt-1:revision-1
                digest: sha256:721699dca2d94770e8334bb36781e0365e4bc7c38191c1fd8fddc43a35a2f797
                digest-domain: canonical-bytes
                request-digest: sha256:de17345385e778d647dd1295ee0b2310571be8be85cd3091bb9784f79eaa7c6f
                room-ref: ./review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:rhx820qrkn6vxpday10nch36:ideation:1
                briefing: briefing:rhx820qrkn6vxpday10nch36:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-30T02:07:08.154809Z"
                decision: approve
                reason: 'Captain conn: staff review findings are closed; the revised 10-file design proves the riskiest operational seam, preserves authority boundaries, and makes stronger observability a design-reset trigger.'
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
---

Make the active workflow the authority for review-finding policy. Reviewers report findings; workers retain evidence and propose materiality, ownership, and disposition; the First Officer authorizes the disposition before finding-driven product edits, commits, or reviewer reruns; validators recommend gate verdicts; and the captain decides changes to scope, accepted value, thresholds, tolerance, or acceptance criteria. Rejection routing carries an existing disposition to `feedback-to`; it neither defines nor repeats the policy.

The shared rejection flow must not define Material, Deferred risk, Polish, Needs decision, correct-but-disproportionate, value-AC, ideation-estimate, or 2× policy. The development workflow and reusable development template own Roborev classification, task-ownership routing, expected-surface calibration, and advisory-round use. Shared gate preparation and presentation preserve workflow-declared finding categories rather than forcing development tiers.

## Problem

The shared `feedback-rejection-flow` currently carries development-specific finding classification even though it triggers only after a feedback gate rejects. Findings also arrive during implementation, validation, consequential quick work, detached audits, and routed feedback. Classification at rejection is therefore globally over-scoped and too late.

The current wording also blurs two separate questions: whether a finding matters and whether the active task owns its remedy. A worker may gather read-only evidence and prepare a rebuttal before consultation, but a Material finding outside the task's approved semantic boundary is `Needs decision`, not permission to expand the task. In-stage round Resolutions remain advisory; the First Officer's disposition authorizes or stops current work, while only an actual gate Resolution binds scope or status.

## Proposed approach

### One authority model

The checkpoint is event-based, not stage-based. It applies when a finding arrives during implementation, validation, a detached audit, consequential First-Officer quick work, or a correction already routed from a rejected gate.

1. **Reviewer — observation authority.** Report the exact finding and evidence. A reviewer may state severity and a tentative release-scope class, but cannot decide that the active task owns a remedy or authorize work.
2. **Worker — proposal authority.** Preserve the original finding, investigate it without changing the candidate, fill the workflow's evidence fields, and propose three separate facts: materiality, whether the approved task owns the remedy, and a disposition. A retained worker triage Resolution remains `actor:ensign`, advisory, and non-binding: it records the proposal/response, never FO authorization.
3. **First Officer — current-work authority.** Compare the proposal with the active workflow policy and send an explicit operational authorization (`fix`, `decline`, `hold`, or `route for decision`) through the current runtime's addressable-worker boundary. This message is not a binding gate Resolution, does not rewrite the approved design, and is not represented by the worker's advisory Resolution.
4. **Validator — recommendation authority.** Verify the delivered candidate against the accepted criteria and the authorized dispositions, then recommend `PASSED` or `REJECTED`. A validator may report a new finding, which re-enters step 1, but does not silently replace the authorized classification.
5. **Captain — design authority.** Decide any change to approved scope, accepted value, threshold, tolerance, acceptance criterion, or a genuine product/compatibility fork. A First Officer routes these decisions and holds the affected work.
6. **Rejection router — transport authority only.** Once a feedback gate or captain has selected revise, carry the rejected snapshot, finding evidence, existing workflow classifications, authorized dispositions, and concrete correction assignment unchanged to `feedback-to`; then re-run review and re-enter the gate. It never re-triages.

Materiality and task ownership are independent. A development `Material` finding caused by and inside the approved task normally yields an FO-authorized fix. A `Material` finding outside the task's approved surface or semantic boundary becomes `Needs decision` for this task: preserve the blocked claim and stop; do not expand scope. A correctly classified `Deferred risk` (real defect, unsupported/unpromised trigger) or `Polish` item (no user-visible loss or protected boundary) may be declined only after the FO authorizes that disposition. If evidence changes, the finding re-enters the checkpoint rather than inheriting stale authorization.

### Before-action boundary

Before FO authorization, the worker may inspect files and history, run non-mutating reproductions or existing tests, use a throwaway checkout for adversarial investigation, and report evidence/proposals to the FO. The candidate worktree remains byte-identical and at the same Git HEAD, and no reviewer rerun starts. Existing reviewer/worker evidence may be retained through the current advisory-round mechanism, but that retention does not manufacture FO authority.

After the worker receives the distinct FO message, only its authorized disposition is allowed. Product edits and candidate commits follow an authorized `fix`; a decline leaves the candidate unchanged but may authorize a reviewer rerun; `hold` and `route for decision` forbid candidate mutation and rerun. Any new finding, changed evidence, or side effect outside that disposition re-enters the checkpoint. If the current host/session cannot expose the authorization to the worker, hold and re-consult rather than infer it.

This deliberately adds one FO round-trip to every finding before candidate mutation. Immediate fixes become slower; that latency buys authority integrity after the e6j failure class, where a finding became an unowned rewrite before disposition. Keep-moving may use the wait to advance unrelated entities or other independent work, but it cannot parallelize through, speculate past, or bypass this entity's checkpoint.

### Observable and spike result

The exact existing operational observable is the ordered runtime boundary `worker completion-signal carrying proposal → FO addressable-worker authorization message → worker completion-signal after authorized work`, paired with external measurements of candidate bytes, Git HEAD, investigation evidence, and reviewer invocation count. The retained Codex spike at `artifacts/ideation-codex-fo-authorization-spike.md` observed:

- after read-only investigation/proposal: candidate SHA and HEAD unchanged, clean worktree, reviewer count `0`;
- after the distinct FO authorization: one candidate commit, changed candidate SHA, reviewer count `1`, and a passing result.

This is a host conversation observable, not an existing retained, host-neutral workflow record. Canonical retained round evidence includes reviewer Annotations, an `actor:ensign` advisory triage Resolution, candidate Git state, and reports; none is FO authorization. The one-host AC is therefore graded only by retained one-off evidence. Canonical durable host-neutral authorization depends on `workflow-neutral-advisory-round-recorder`; until that workflow supplies an observable, this task does not add a canonical event or claim standing/cross-host enforcement. Any implementation attempt to make that stronger claim stops for a design reset instead of inventing state.

Pi remains in semantic scope because its runtime exposes an addressable-worker route. Its checkpoint is the same workflow-policy act: the FO receives the worker proposal and sends the distinct authorization/dispatch message before mutation or rerun. This task adds no Pi-specific standing evidence, scenario registration, or state representation. If Pi cannot execute or expose that boundary in the implementation drive, stop for a design reset; do not silently exempt Pi or treat the advisory worker Resolution as the missing signal.

### Workflow ownership and loading

Add a named `## Review-finding disposition` section to the active development README and the reusable development template. The shared FO contract says, generically, to locate and read the active workflow's declared review policy when findings arrive. Existing fence-safe `status --read ... --json` heading discovery is sufficient; no new selector, status field, command, or schema is needed. The implementation and validation stage definitions point workers at this workflow section and state the pre-action checkpoint, so the rule is present both for an in-stage worker and an FO acting directly.

The development section owns:

- the Roborev `Material / Deferred risk / Polish / Needs decision` taxonomy;
- the four evidence fields and the separate task-ownership decision;
- expected files/LOC and semantic-boundary calibration, including the workflow's tolerance and AC-drift response;
- the FO consultation requirement and permitted read-only investigation;
- the advisory-round record after reviewer, worker, and FO consultation; and
- the rule that only the captain changes scope, accepted value, thresholds, tolerance, or ACs.

The reusable template carries the same development policy, while other commission templates remain free to declare different categories and evidence. The shared `present-gate` renderer uses the active workflow's categories verbatim (or a neutral unclassified findings list when none are declared). The shared `feedback-rejection-flow` accepts an already-authorized correction package, keeps category labels opaque, and transports any workflow-defined correction-round projection without defining its grammar.

### Four-field evidence remains four fields

Keep `.roborev.toml`, `docs/specs/check-finding-triage-materiality.sh`, and `docs/specs/testdata/finding-triage-materiality.tsv` as the independent development materiality oracle. When field 3 asserts an affected AC/boundary, it uses exactly one inline citation form followed by a nonblank claim:

- `value-ac[AC-N]: <affected outcome>`, where `N` is a positive decimal integer;
- `captain-ruling[YYYY-MM-DD]: <protected boundary>`, with exactly four/two/two digits; or
- `contract[repo/relative/path#anchor]: <protected boundary>`, with a nonblank repo-relative path and fragment, neither containing brackets or whitespace.

`none: <nonblank rationale>` remains the valid form when no AC/boundary is affected and can never establish Material. The POSIX ERE for a cited field is `^(value-ac\[AC-[1-9][0-9]*\]|captain-ruling\[[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]\]|contract\[[^][[:space:]#]+#[^][[:space:]]+\]):[[:space:]]+[^[:space:]].*$`; `none:` has the same nonblank-claim suffix. Rows remain exactly seven TSV columns: id + four evidence fields + recorded class + expected result. Red cases cover the old uncited form, blank locator, malformed locator/date/AC, blank claim, and an extra fifth evidence field (`NF != 7`), plus the existing Material-recorded-as-declined control.

### Exact instruction changes

- `first-officer-shared-core.md` rejection branch — before: “If gated and a reviewer recommends `REJECTED` at a configured feedback gate, invoke `«feedback.route»` before Captain presentation; otherwise complete ... revise invokes `«feedback.route»`.” After: “If a rejection carries a new or unresolved finding, first re-enter the active workflow's review-finding checkpoint. `«feedback.route»` is ineligible until the worker proposal has a distinct FO-authorized disposition and a concrete revise assignment. If either is missing, hold at the checkpoint. Once both exist, route before Captain presentation; routing preserves the finding/evidence/classification/disposition/assignment package and never classifies it. A later `revise` invokes the same eligibility check.” This is the always-on branch, not recorder state.
- `fo-dispatch-core.md` payload statement (separate from branch eligibility) — before: `--feedback-context-file` is added “when the stage has `feedback-to`.” After: “the feedback-context file carries the already-authorized package and concrete revise assignment with workflow labels unchanged; it never asks the target worker to classify again.”
- `feedback-rejection-flow/SKILL.md` — before: “Past the declared tolerance — 2× unless the entity declares its own,” the Material/correct-but-disproportionate rules, and the fixed files/LOC/estimate/AC Feedback Cycles format. After: “Read the already-authorized workflow package; if missing, hold at the checkpoint; transport it unchanged; append a workflow-defined correction-round projection verbatim when one is declared; apply the existing cycle-3 escalation/reuse/budget mechanics; re-run the reviewer; re-enter the gate.” The generic skill defines no taxonomy, tolerance, estimate, AC-drift rule, or projection grammar.
- `present-gate/SKILL.md` — before: “render ... in two tiers — `Material:` ... and `Polish:`.” After: “render findings in the active workflow's declared category order and exact labels; preserve the recorded category, omit empty categories, and use a neutral `Findings:` list when the workflow declares none. Presentation never classifies.”
- `docs/dev/README.md` / development template — before: “triage before fixing ... Material findings are fixed ... escalate a needs-decision finding.” After: “investigate read-only, preserve the four-field evidence, and propose materiality, task ownership, and disposition; obtain FO authorization before a candidate edit, commit, or reviewer rerun. The captain alone changes scope, accepted value, thresholds, tolerance, or ACs.” Implementation, validation, detached-audit, and correction-loop clauses invoke the named section; expected-surface and AC-drift policy remain here.
- `.roborev.toml` — before: “Affected value AC or non-negotiable boundary.” After: “Affected value AC or non-negotiable boundary using `value-ac[AC-N]`, `captain-ruling[YYYY-MM-DD]`, or `contract[path#anchor]`, followed by a nonblank claim.” It remains field 3 of 4.

### Mechanism/value check

- **Named workflow policy read → AC-1/AC-2.** Simplest alternative: duplicate development classes into shared skills. Insufficient because it makes non-development workflows inherit Roborev policy and still leaves quick work outside stage-local text.
- **Runtime addressable-worker authorization → AC-1.** Simplest alternative: treat the worker's advisory Resolution as authorization. Insufficient and false: that Resolution is `actor:ensign`, advisory, and non-binding. The Codex spike proves an operational message boundary exists but not durable host-neutral retention; a durable requirement triggers design reset.
- **Opaque rejection payload and workflow-category presentation → AC-3.** Simplest alternative: reclassify at rejection/presentation. Insufficient because it is too late, can contradict the in-stage decision, and gives generic layers development authority.
- **Retained one-off live drive → AC-1.** Simplest alternative: permanently register a shared scenario now. Insufficient because there is no canonical retained host-neutral authorization event to grade repeatably. Retain a one-host drive for the gate; do not add standing scenario/docs/Pi surface.
- **One-off file-scoped inspection → AC-3/AC-5.** Simplest alternative: add a contractlint lexical-absence test. Insufficient because it polices wording rather than behavior and adds permanent surface for a one-time ownership move.
- **Four-field fixture strengthened in place → AC-4.** Simplest alternative: add a fifth authority field. Insufficient because it breaks the established evidence record and lets field count drift instead of making the boundary claim accountable.

## Acceptance criteria

**AC-1 (VALUE, one-host live; cross-host semantic scope) — On a host exposing an addressable worker, a finding can produce read-only investigation and a worker proposal while candidate bytes, Git HEAD, and reviewer count remain unchanged; edit/commit/rerun occurs only after a distinct FO authorization message.**

Verified by: the pre-gate Codex spike retained at `artifacts/ideation-codex-fo-authorization-spike.md` records unchanged candidate SHA/HEAD, clean state, investigation evidence, and reviewer count `0` after proposal; then the distinct FO message, one candidate commit, changed SHA/HEAD, reviewer count `1`, and PASS. Implementation repeats this as a retained one-off drive against the changed contract. An early edit/commit/rerun or missing investigation evidence fails the measurement. Codex is the current evidence host; Pi is semantically bound through its addressable-worker route but gains no standing evidence here. Inability on either runtime to execute/observe the boundary requires design reset, not invented state or an exemption.

**AC-2 (offline + human review) — Materiality, task ownership, and disposition remain separate: an owned Material finding is eligible for an FO-authorized fix; an out-of-scope Material finding holds unchanged for captain decision; a correctly non-material Deferred risk or Polish finding is eligible for an FO-authorized decline.**

Verified by: the four-field TSV recomputes the unsupported/adversarial duplicate-member shape as non-material and the supported producer/value-AC control as Material; a one-off file-scoped review of the development policy checks the independent ownership fork and exact `Needs decision` hold. The review rejects always-fix, always-decline, “unsupported deferred risk is supported,” or “Material means this task owns it” wording.

**AC-3 (offline behavior + one-off inspection) — Generic rejection routing transports an already-authorized package and workflow-defined projection unchanged, while gate presentation preserves workflow labels; neither layer defines development taxonomy, tolerance, or Feedback Cycles grammar.**

Verified by: extend the real `dispatch.Run` reflow test with a non-development package (`Blocking / Advisory`) and assert byte-identical payload output. Validation performs a one-off, file-scoped inspection of the shared-core rejection branch, rejection skill, and presenter against the approved before/after text and retains its output in the stage report. Relabeling the fixture or finding development classification/calibration/projection grammar in those owner sites fails.

**AC-4 (offline) — Development findings still use exactly four evidence fields, and field three identifies both the affected value AC/boundary and its authority before it can support Material classification.**

Verified by: `docs/specs/check-finding-triage-materiality.sh` runs the checked-in TSV cases and red controls. Uncited, blank, malformed, blank-claim, and extra-fifth-field cases reject; cited AC/captain-ruling/contract examples and a non-material `none:` case pass. Removing citation or exact-column checks, treating `none:` as Material, or recording a material finding as declined turns the test red.

**AC-5 (offline + one-host evidence) — Advisory worker triage, validator recommendation, FO operational authorization, and captain binding design authority remain distinct, with no new recorder schema or status transition.**

Verified by: existing advisory-round tests continue to show `actor:ensign`, advisory, no gate application/status mutation; the retained spike identifies the distinct FO message rather than the worker Resolution. One-off diff inspection must show no `internal/gates`, schema, status, dispatch production, or stored-format change. Treating the worker Resolution as FO authorization or changing AC/scope/tolerance without captain authority fails.

## Test plan

- **Retained one-off live drive (medium, no committed harness):** repeat the proven Codex journey against the changed installed contract and retain transcript excerpts plus candidate SHA, HEAD, investigation evidence, and reviewer count. If the distinct FO message cannot be observed, stop for design reset. Pi shares the semantic requirement through its addressable-worker act, but this task adds no Pi-specific standing evidence; an observed Pi inability also triggers reset. Do not register a shared scenario, documentation ID, or Pi gap.
- **One-off policy inspection (low, no committed test):** inspect only the seven changed policy/instruction files for the approved ownership move and exact branch wording; retain command/output in validation evidence. No contractlint lexical test.
- **Offline four-field oracle (low):** extend the existing shell fixture/check for the exact inline grammar, exact seven-column rows, and missing/blank/malformed/fifth-field red cases.
- **Offline routing (low):** extend `internal/ensigncycle/feedback_test.go` so `dispatch.Run` transports opaque non-development classifications and disposition bytes. Mutation: normalize them to development labels; assertion fails.
- **Regression (medium):** `go test ./...`, `go test ./... -race`, format gate, and existing advisory-round tests. No new CLI E2E, recorder test, permanent live fixture, scenario registration, scenario documentation, or Pi coverage entry.

Spike result: the smallest Codex drive proved operational ordering and external measurements, but found no existing retained, host-neutral FO-authorization event. The retained artifact is sufficient one-off gate evidence, not a canonical mechanism. This finding removes the permanent shared-scenario proposal and makes durable/cross-host authorization a design-reset trigger.

## Expected surface

- `.roborev.toml`: 2 insertions.
- `docs/dev/README.md`: 38 insertions.
- `skills/commission/references/templates/development.md`: 38 insertions.
- `skills/first-officer/references/first-officer-shared-core.md`: 20 insertions.
- `skills/first-officer/references/fo-dispatch-core.md`: 5 insertions.
- `skills/feedback-rejection-flow/SKILL.md`: 16 insertions.
- `skills/present-gate/SKILL.md`: 10 insertions.
- `docs/specs/check-finding-triage-materiality.sh`: 18 insertions.
- `docs/specs/testdata/finding-triage-materiality.tsv`: 8 insertions.
- `internal/ensigncycle/feedback_test.go`: 32 insertions.

Baseline: exactly 10 files and 187 insertions. Tolerance is +1 file (10%) and +38 insertions (20% of 187 = 37.4, rounded up): hard ceilings are 11 files and 225 insertions. The one tolerated extra file may only extract an `internal/ensigncycle` routing fixture helper; it cannot add contractlint, permanent live-scenario, scenario-doc, or Pi-registration surface. Deletions are expected where development policy leaves generic skills and do not count toward the insertion ceiling; any semantic change outside the permitted list still requires reset.

Permitted semantic changes: finding-disposition authority and action ordering; workflow-category rendering at gates; workflow-opaque rejection payloads; development evidence field three gains an authority citation. Command grammar, exit codes, stored formats, advisory-round recorder semantics, gate application semantics, entity frontmatter, stage identity, and task acceptance values do not change.

## Out of scope

Do not redesign the advisory-round recorder's storage validation in this task; that is owned by `workflow-neutral-advisory-round-recorder`. Do not redesign the generic gate criteria-source interface, dispatch stage identity, Roborev provider setup, or the review-budget/recorder storage format. Do not make one held finding block unrelated sprint work.

## Stage Report: ideation

- DONE: Define one workflow-owned authority model across in-stage review, validation, detached audits, consequential quick work, and routed rejection.
  The body assigns observation, proposal, operational authorization, recommendation, design, and transport authority separately across all five arrival paths.
- DONE: Specify acceptance criteria and falsifiable tests that preserve read-only investigation while requiring FO disposition before finding-driven edits, commits, or reruns.
  AC-1/AC-2 grade candidate bytes, Git HEAD, reviewer invocation count, and ordered events; pre-authorization side effects and over-blocked investigation are explicit red cases.
- DONE: Declare expected files, insertions, tolerance, and permitted semantic changes; retain the four-field evidence test and exclude recorder redesign.
  The expected surface is 20 files/~560 insertions with +30% plumbing tolerance; field three gains an authority citation without a fifth field, and recorder/schema/command changes are forbidden.

### Summary

Ideation now makes the active workflow the policy owner and the First Officer the current-work disposition authority, while keeping captain design authority and treating rejection as transport of an existing decision. It specifies exact instruction moves, a category-neutral shared layer, a falsifiable live incident replay with material/ownership controls, and a bounded change surface that leaves the advisory recorder to its sibling task.

## Stage Report: ideation (cycle 2)

- DONE: Replace the false no-spike claim and define the available FO-authorization observable.
  The retained Codex spike measured unchanged pre-authorization SHA/HEAD/count with read-only evidence, then a distinct FO message, one commit, and reviewer count 1; it also proved no canonical retained host-neutral authorization event exists.
- DONE: Keep worker advisory triage distinct from FO operational authorization.
  The design states that `actor:ensign` advisory Resolution records a proposal only; FO authority is the distinct runtime message, and durable host-neutral authority requires design reset.
- DONE: Cut permanent proof surface to the approved baseline.
  Expected implementation is 10 files/187 insertions: seven policy files, shell+TSV oracle, and opaque routing test; contractlint/shared-scenario/docs/Pi additions are excluded.
- DONE: Define mechanical field-3 citation grammar and red cases.
  Three complete inline forms cover AC, dated captain ruling, and named contract; missing, blank, malformed, blank-claim, extra-field, and false-decline rows reject.
- DONE: Specify the shared-core rejection branch before/after and workflow-defined projection transport.
  New findings re-enter the checkpoint; missing FO disposition/assignment holds; eligible routing preserves the package without classification or generic Feedback Cycles grammar.

### Summary

Cycle 2 incorporates all five Material staff findings and the useful polish without changing product files or expanding recorder/schema/command semantics. The design now makes its proof limit explicit: operational ordering is demonstrated on Codex and retained as one-off evidence, while durable cross-host FO authorization is a future design-reset question.

## Stage Report: ideation (cycle 3)

- DONE: State the cadence trade and keep-moving boundary.
  Every finding pays one FO round-trip before mutation; unrelated work may proceed, but authority integrity is never bypassed.
- DONE: Bind standing observability to the workflow-neutral recorder dependency.
  One-host evidence grades this task; standing/cross-host enforcement waits on `workflow-neutral-advisory-round-recorder` or triggers design reset.
- DONE: State Pi operability without adding Pi-specific proof surface.
  Pi shares the addressable-worker checkpoint semantics; inability to execute or observe that act triggers reset, not exemption or invented state.

### Summary

Cycle 3 adds only the second-review discipline notes: explicit latency-for-integrity, the recorder dependency for stronger claims, and Pi's semantic scope. The approved implementation surface and all recorder/schema/command exclusions remain unchanged.

## Stage Report: implementation

- DONE: Implement the approved six-authority workflow policy across the seven instruction surfaces, with FO authorization before finding-driven edits/commits/reruns and category-neutral rejection/presentation plumbing.
  Code commit `00f8c203` moves development policy into the workflow/template, makes shared routing opaque, and preserves workflow labels; the file-scoped inspection caught and, after distinct FO authorization, removed the final generic `material` label leak.
- DONE: Strengthen the existing four-field oracle with mechanical field-3 citations and exact-column red controls, and prove opaque non-development feedback transport without new recorder/schema/command behavior.
  `check-finding-triage-materiality.sh` accepts cited AC/ruling/contract and `none:` controls while rejecting uncited/malformed/blank/eight-column/false-decline rows; `TestFeedbackReflowRoutesFixRequest` byte-compares `Blocking`/`Advisory` payload and would fail on normalization or projection rewrite.
- DONE: Stay within 10 files/187 insertions (hard ceilings 11/225), run the required offline checks plus retained one-host authorization drive, and stop for design reset if the FO boundary cannot be observed on Codex or Pi.
  The code diff is exactly 10 files/117 insertions; `go test ./...`, `go test ./... -race`, gofmt, focused contract/oracle checks, and diff inspection passed, while state commit `c63c0418` retains unchanged pre-authorization SHA/HEAD/count then one authorized commit, one reviewer run, and PASS on Codex.

### Summary

The implementation separates observation, proposal, current-work authorization, validation recommendation, captain design authority, and rejection transport without adding stored authority state or command behavior. Development keeps its calibrated taxonomy and strengthened four-field evidence, while generic layers carry workflow classifications and projections unchanged.

## Stage Report: validation

- DONE: Independently verify all five ACs against the exact 10-file/117-insertion candidate, including six-role authority separation, workflow-label preservation, and absence of recorder/schema/command drift.
  AC-1 FAIL: retained commit `c63c0418` proves ordered Codex authorization, but candidate `show-stage-def` omits the referenced policy on normal implementation/validation dispatches.
  AC-2 FAIL: the workflow text separates all six authorities and ownership outcomes, but dispatched workers receive only a pointer to the omitted section.
  AC-3 PASS: `TestFeedbackReflowRoutesFixRequest` byte-compares the exact `Blocking`/`Advisory` package; file-scoped inspection found no development taxonomy/calibration in generic owner sites.
  AC-4 FAIL: the current 18-row oracle passes, but in-memory deletion of either `exact_columns` or `valid_boundary` also exits 0, contradicting the promised red mutations.
  AC-5 PASS: focused advisory-round tests preserve `actor:ensign`, advisory/no-status semantics; the exact diff changes no gates/status/dispatch production, schema, command, or stored format.
- DONE: Re-run the full Go/race suites, four-field shell oracle and opaque-routing controls, and falsify the ordering/ownership claims with the approved retained-drive and negative cases.
  `gofmt -w ./cmd ./internal` left the worktree clean; `go test ./...` and `go test ./... -race` passed, as did focused routing, empty-payload, four-field, and advisory-round controls.
  The retained drive still has the recorded candidate SHA, candidate-only commit, and reviewer count `1`; guard-deletion probes independently exposed the AC-4 evidence defect.
- DONE: Run the required classified Roborev review; triage each finding by released workflow, observable harm, authority-citing value boundary, trigger, and task ownership before recommending PASSED or REJECTED.
  Roborev branch-final panel job `318` returned Changes requested.
  Material outcome — missing dispatch context: normal implementation/validation ensigns lose the policy body; `value-ac[AC-1]`/`value-ac[AC-2]`; every dispatch lacking `context-sections`; task-owned frontmatter/template fix.
  Material evidence — dead oracle guards: validators can delete exact-column or standalone-validity gates without red; `value-ac[AC-4]`; reproduced now; task-owned fixture fix.
  Deferred risk — legacy oracle class names: exact trigger is adding policy-category rows to the materiality-only fixture; current cited/non-material controls pass; promote if the script becomes the category serializer.
  Deferred risk — custom projection plus canonical gate recorder: no current workflow couples those surfaces and dev's declared line passes; promote if generic routing mandates `gate record --feedback-cycle`.
  Polish — payload terminator and stale NEG-A commentary: exact current emission is asserted by the positive byte comparison; revisit if dispatch section order changes.
  Deferred risk — file-input transport seam: production uses untrimmed `os.ReadFile` and downstream byte comparison passes; promote if file loading gains normalization.
  Deferred risk — removed default tolerance: current dev tasks declare expected surface/tolerance; promote when a supported workflow omits one and still enters another correction round.
  Polish — older audit-routing shorthand: the active generic route now enforces authorization; revise the shorthand when those docs next change.
  Deferred risk — gate-lifecycle three-label briefing shorthand: presenter preserves active labels and no omission was observed; promote on a reproduced `Needs decision` loss.

### Summary

Recommendation: **REJECTED** for two task-owned Material findings: the workflow policy is not delivered through the normal dispatched stage context, and AC-4's fixtures do not falsify removal of two promised guards. The runtime suites, opaque routing, six-authority wording, retained authorization drive, and no-drift checks otherwise pass; candidate commit `00f8c203` remains untouched.

### Feedback Cycles

- Cycle 1: REJECTED — validation/Roborev job 318; surface 10 files/117 insertions vs estimate 10 files/187 insertions (63%); AC unchanged
