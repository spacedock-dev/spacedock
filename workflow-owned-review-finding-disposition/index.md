---
title: Make review-finding disposition workflow-owned and separate in-stage review from rejection routing
status: ideation
source: "Captain boundary audit, 2026-07-30: an end-of-implementation Roborev round never enters feedback-rejection-flow, and finding classification must happen when findings arrive rather than only after rejection."
started: 2026-07-30T01:34:32Z
completed:
verdict:
score: "0.90"
worktree:
issue:
pr:
sprint: durable-decisions
id: rhx820qrkn6vxpday10nch36
gates:
    version: 1
    current:
        gate: gate:rhx820qrkn6vxpday10nch36:backlog
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
2. **Worker — proposal authority.** Preserve the original finding, investigate it without changing the candidate, fill the workflow's evidence fields, and propose three separate facts: materiality, whether the approved task owns the remedy, and a disposition. The worker does not self-authorize its proposal.
3. **First Officer — current-work authority.** Compare the proposal with the active workflow policy and authorize `fix`, `decline`, `hold`, or `route for decision` for this work. This is operational authorization, not a binding gate Resolution and not permission to rewrite the approved design.
4. **Validator — recommendation authority.** Verify the delivered candidate against the accepted criteria and the authorized dispositions, then recommend `PASSED` or `REJECTED`. A validator may report a new finding, which re-enters step 1, but does not silently replace the authorized classification.
5. **Captain — design authority.** Decide any change to approved scope, accepted value, threshold, tolerance, acceptance criterion, or a genuine product/compatibility fork. A First Officer routes these decisions and holds the affected work.
6. **Rejection router — transport authority only.** Once a feedback gate or captain has selected revise, carry the rejected snapshot, finding evidence, existing workflow classifications, authorized dispositions, and concrete correction assignment unchanged to `feedback-to`; then re-run review and re-enter the gate. It never re-triages.

Materiality and task ownership are independent. A development `Material` finding caused by and inside the approved task normally yields an FO-authorized fix. A `Material` finding outside the task's approved surface or semantic boundary becomes `Needs decision` for this task: preserve the blocked claim and stop; do not expand scope. A deferred risk or polish item may be declined only after the FO authorizes that disposition. If evidence changes, the finding re-enters the checkpoint rather than inheriting stale authorization.

### Before-action boundary

Before FO authorization, the worker may inspect files and history, run non-mutating reproductions or existing tests, use a throwaway checkout for adversarial investigation, and report evidence/proposals to the FO. Control-plane writes needed to retain the original reviewer output or the FO's own disposition are allowed; they do not change the candidate. The candidate worktree must otherwise remain byte-identical and at the same Git HEAD, and no reviewer rerun may start.

After authorization, only the authorized disposition is allowed. Product edits and candidate commits must follow an authorized `fix`; a decline may leave the candidate unchanged; `hold` and `route for decision` forbid candidate mutation. Any new edit, commit, or rerun outside that disposition requires a new checkpoint. This keeps investigation useful without allowing review pressure to become implicit scope authority.

### Workflow ownership and loading

Add a named `## Review-finding disposition` section to the active development README and the reusable development template. The shared FO contract says, generically, to locate and read the active workflow's declared review policy when findings arrive. Existing fence-safe `status --read ... --json` heading discovery is sufficient; no new selector, status field, command, or schema is needed. The implementation and validation stage definitions point workers at this workflow section and state the pre-action checkpoint, so the rule is present both for an in-stage worker and an FO acting directly.

The development section owns:

- the Roborev `Material / Deferred risk / Polish / Needs decision` taxonomy;
- the four evidence fields and the separate task-ownership decision;
- expected files/LOC and semantic-boundary calibration, including the workflow's tolerance and AC-drift response;
- the FO consultation requirement and permitted read-only investigation;
- the advisory-round record after reviewer, worker, and FO consultation; and
- the rule that only the captain changes scope, accepted value, thresholds, tolerance, or ACs.

The reusable template carries the same development policy, while other commission templates remain free to declare different categories and evidence. The shared `present-gate` renderer uses the active workflow's categories verbatim (or a neutral unclassified findings list when none are declared). The shared `feedback-rejection-flow` accepts an already-authorized correction package and keeps category labels opaque.

### Four-field evidence remains four fields

Keep `.roborev.toml`, `docs/specs/check-finding-triage-materiality.sh`, and `docs/specs/testdata/finding-triage-materiality.tsv` as the independent development materiality oracle. Strengthen the third field in place: `Affected value AC or non-negotiable boundary` cites its authority, for example `AC-2`, a dated captain ruling, or a named contract. Do not add a fifth field. A missing authority citation cannot support a Material classification, and the existing material/decline red control remains.

### Exact instruction changes

- `first-officer-shared-core.md` / `fo-dispatch-core.md` — before: `«feedback.route»` “route findings back to the target stage.” After: “read the active workflow's review policy; allow investigation while holding candidate edits, commits, and reruns; authorize the current-work disposition or route a workflow-reserved decision. `«feedback.route»` forwards the existing disposition and revise assignment.” `--feedback-context-file` preserves existing category labels, authorization, and concrete work; it does not request reclassification.
- `feedback-rejection-flow/SKILL.md` — before: “Past the declared tolerance — 2× unless the entity declares its own” and “Material findings are fixed; each correct-but-disproportionate finding is declined.” After: “Read the active workflow's already-authorized disposition and correction authority; append the workflow's correction-round projection; route the package unchanged; re-run the reviewer; re-enter the gate. If a workflow-required disposition is absent, hold and consult the workflow owner.” Keep cycle transport, reuse/fresh dispatch, context-budget, rerun, and gate re-entry.
- `present-gate/SKILL.md` — before: “render ... in two tiers — `Material:` ... and `Polish:`.” After: “render findings in the active workflow's declared category order and exact labels; preserve the recorded category, omit empty categories, and use a neutral `Findings:` list when the workflow declares none. Presentation never classifies.”
- `docs/dev/README.md` / development template — before: “triage before fixing ... Material findings are fixed ... escalate a needs-decision finding.” After: “investigate read-only, preserve the four-field evidence, and propose materiality, task ownership, and disposition; obtain FO authorization before a candidate edit, commit, or reviewer rerun. The captain alone changes scope, accepted value, thresholds, tolerance, or ACs.” Implementation, validation, detached-audit, and correction-loop clauses invoke the named section; expected-surface and AC-drift policy remain here.
- `.roborev.toml` — before: “Affected value AC or non-negotiable boundary.” After: “Affected value AC or non-negotiable boundary, citing its authority (AC, dated captain ruling, or named contract).” It remains field three of four.

### Mechanism/value check

- **Named workflow policy read → AC-1/AC-2.** Simplest alternative: duplicate development classes into shared skills. Insufficient because it makes non-development workflows inherit Roborev policy and still leaves quick work outside stage-local text.
- **FO pre-action checkpoint → AC-1.** Simplest alternative: let workers triage and fix, then review at the gate. Insufficient because the unnecessary edit/commit/rerun has already happened before a gate exists.
- **Opaque rejection payload and workflow-category presentation → AC-3.** Simplest alternative: reclassify at rejection/presentation. Insufficient because it is too late, can contradict the in-stage decision, and gives generic layers development authority.
- **Ordered live event fixture → AC-1/AC-2/AC-3.** Simplest alternative: grep instruction prose. Insufficient because correct words do not establish agent action order; the fixture observes candidate bytes, HEAD, messages, and reviewer invocations.
- **Four-field fixture strengthened in place → AC-4.** Simplest alternative: add a fifth authority field. Insufficient because it breaks the established evidence record and lets field count drift instead of making the boundary claim accountable.

## Acceptance criteria

**AC-1 (VALUE, offline + live) — A finding cannot cause a candidate edit, candidate commit, or reviewer rerun before the First Officer authorizes its workflow-owned disposition, while read-only investigation remains available.**

Verified by: a shared live scenario emits ordered events for finding arrival, non-mutating reproduction, worker proposal, FO authorization, candidate mutation/commit, and reviewer rerun. The withheld-authorization case must leave candidate bytes, Git HEAD, and reviewer invocation count unchanged while still showing the evidence-producing read/test; the authorized-material case must mutate and rerun only after authorization. Offline negative controls reorder each side effect before authorization and must make the same state/event assertion fail. This goes red if a worker edits, commits, or reruns early, or if the checkpoint blocks the permitted investigation.

**AC-2 (VALUE, live) — Materiality, task ownership, and disposition remain separate: an owned Material finding produces one authorized fix, an out-of-scope Material finding holds unchanged for captain decision, and a supported deferred/polish decline leaves the candidate unchanged but permits an authorized rerun.**

Verified by: replay the retained duplicate-member incident shape against a fixture candidate: the unsupported/adversarial trigger and cited value AC produce an FO-authorized decline with zero product-line delta; changing only the trigger to a supported producer that breaks that AC produces one authorized fix; changing only ownership to outside the approved semantic boundary produces `Needs decision`, no diff, no commit, and no rerun. An always-fix, always-decline, or “Material means this task owns it” implementation fails at least one branch.

**AC-3 (offline + live) — Generic rejection routing and gate presentation preserve a workflow's existing category labels and authorized disposition without introducing development taxonomy or reclassification.**

Verified by: extend the real `dispatch.Run` reflow fixture with a non-development correction package (`Blocking / Advisory`) and assert byte-preserved payload output; a live gate/rejection scenario asserts the same labels and disposition reach `feedback-to` and the captain-facing review unchanged. A contractlint structural-absence guard with discriminator controls uses file-specific boundaries: the rejection skill excludes all development taxonomy/calibration tokens, the presenter excludes hard-coded category labels, and the new generic checkpoint excludes development classes without flagging the shared core's unrelated AC-verification vocabulary. Reintroducing a forbidden token at one of those policy-owner sites, or relabeling the fixture, makes these checks fail.

**AC-4 (offline) — Development findings still use exactly four evidence fields, and field three identifies both the affected value AC/boundary and its authority before it can support Material classification.**

Verified by: `docs/specs/check-finding-triage-materiality.sh` runs the checked-in TSV cases and red controls. Add controls for an uncited third field and an extra fifth evidence field; both must reject, while cited `AC`, captain-ruling, and named-contract examples pass. Removing the citation check, accepting a fifth field, or recording a material finding as declined turns the test red.

**AC-5 (offline + live) — Advisory in-stage review, validator recommendation, FO operational authorization, and captain binding design authority remain distinct, with no new recorder schema or status transition.**

Verified by: existing `internal/gates` advisory-round tests continue to show no gate application/status mutation; the finding-disposition live scenario fails if FO authorization edits AC/scope/tolerance or if a validator self-binds a verdict, and passes only when such a fork is held for the captain. `git diff --name-only` must show no `internal/gates`, schema, or stored-format production change. A round Resolution becoming binding, or an FO/validator changing the approved design, fails.

## Test plan

- **Offline policy boundary (low):** add `internal/contractlint/workflow_review_policy_boundary_test.go` for literal development-vocabulary absence from generic owners, with must-flag/must-pass discriminator cases. This proves ownership/structural absence only, not behavior.
- **Offline four-field oracle (low):** extend the existing shell fixture/check for authority citation and exact field count. Run `docs/specs/check-finding-triage-materiality.sh`.
- **Offline routing (low):** extend `internal/ensigncycle/feedback_test.go` so `dispatch.Run` transports opaque non-development classifications and disposition bytes. Mutation: normalize them to development labels; assertion fails.
- **Shared live behavior (medium/high):** add one host-neutral `finding-disposition-authority` fixture and ordered-event grader to the existing shared runtime-scenario framework, with Codex and Claude runners and an honest Pi coverage entry. It exercises read-only investigation, authorized fix/decline, outside-scope hold, opaque rejection routing, and category-preserving presentation. Negative fixtures spend no model.
- **Regression (medium):** `go test ./...`, `go test ./... -race`, live Codex/Claude shared-scenario lanes required by the host-neutral FO/skill diff, and existing advisory-round recorder tests. No new CLI E2E or recorder live test is needed because command grammar and storage are unchanged.

No spike needed: current `status --read docs/dev/README.md --json` returned fence-safe heading offsets; `TestFeedbackReflowRoutesFixRequest` exercised opaque dispatch payload transport; `TestAssertRejectionFlow` exercised the two-cycle seam; and the existing four-field shell fixture passed its positive and red controls on 2026-07-30. The risky claim is agent action order, so implementation's first test is the live ordered-event fixture rather than another prose check.

## Expected surface

Expected production/instruction files:

- `.roborev.toml`
- `docs/dev/README.md`
- `skills/commission/references/templates/development.md`
- `skills/first-officer/references/first-officer-shared-core.md`
- `skills/first-officer/references/fo-dispatch-core.md`
- `skills/feedback-rejection-flow/SKILL.md`
- `skills/present-gate/SKILL.md`

Expected proof/support files:

- `docs/specs/check-finding-triage-materiality.sh`
- `docs/specs/testdata/finding-triage-materiality.tsv`
- `internal/contractlint/workflow_review_policy_boundary_test.go` (new)
- `internal/ensigncycle/feedback_test.go`
- `internal/ensigncycle/shared_scenarios_test.go`
- `internal/ensigncycle/shared_fixtures_test.go`
- `internal/ensigncycle/shared_assertions_impl_test.go`
- `internal/ensigncycle/shared_assertions_test.go`
- `internal/ensigncycle/shared_scenarios_negative_test.go`
- `internal/ensigncycle/codex_live_runner_test.go`
- `internal/ensigncycle/claude_live_runner_test.go`
- `internal/ensigncycle/pi_shared_coverage_test.go`
- `docs/specs/scenario-testing-principles.md`

Estimate: 20 files, about 560 insertions and 150 deletions. Tolerance: at most 23 files and 730 insertions (+30%), only for shared-live-fixture plumbing or existing structural test registration. Any production change under `internal/gates`, `internal/status`, `internal/dispatch`, `docs/schema`, or any new state field/command is outside tolerance regardless of LOC and requires a design reset. `internal/ensigncycle/feedback_test.go` may change, but `internal/dispatch` production behavior may not.

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
