---
commissioned-by: spacedock@0.13.0-dev
entity-type: task
entity-label: task
entity-label-plural: tasks
id-style: sd-b32
state: .spacedock-state
trunk: main
stages:
  defaults:
    worktree: false
    concurrency: 3
  states:
    - name: backlog
      initial: true
      gate: true
    - name: ideation
      gate: true
      model: opus
    - name: implementation
      worktree: true
    - name: validation
      worktree: true
      fresh: true
      feedback-to: implementation
      gate: true
    - name: done
      terminal: true
---

# Build Spacedock v1 - Go Launcher Workflow

Spacedock v1 is the Go launcher and compatibility bridge for the next Spacedock command surface. This workflow tracks design and implementation tasks from initial concepts through validated, shippable behavior.

Runtime entities live in `.spacedock-state`, a per-workflow state checkout. During bootstrap, `.spacedock-state/README.md` may symlink to this README so current status tooling can operate against the state checkout directly.

This workflow registers lifecycle mods under `_mods/`: a `pr-merge` hook (opens a code-branch PR at the merge boundary and tracks it to merge) and a standing `comm-officer` prose-polish teammate.

## Sprints

This workflow tracks individual tasks. A **sprint** groups several tasks into one value-increment — a convention *stacked on top* of the per-task flow, not a builtin Spacedock construct. See `_proposals/sprint-roadmap-construct.md` for the decision record (ship as a skill + commission template; defer builtin).

**Before shaping or driving a sprint, load the discipline:** [`docs/roadmap/README.md`](../roadmap/README.md) carries the operational sprint construct: the Shaping-FO and Commander roles, and the canonical sprint-folder shape (`index.md` for durable strategy, `staff-review.md` for the readiness gap analysis, `dispatch-sprint-execution.md` for the cold-boot Commander package). The rules below summarize the per-task flow; that doc is the full prep-and-drive procedure.

- **Membership is a query, never a hard-coded list.** Entities carry `sprint: <slug>` (plus an optional `sprint-readiness:` filter). List members: `spacedock status --workflow-dir docs/dev --where sprint=<slug>`.
- **The index is durable strategy, not a tracker.** `docs/roadmap/<topic>/index.md` holds goal, scope, DoD, deliverable, and the movable target-train line only. It does NOT enumerate members or track their state — that is the query above, which never goes stale. A Commander's execution bookkeeping (what shipped, what's in flight, per-member status) belongs in the Commander's own chain of handoffs, never in the shared index.
- **Ownership is cross-session.** A sprint may be *driven* by one **Commander** session: a single FO boots `spacedock:first-officer`, takes the packaged sprint, and drives its members to the DoD. While a sprint is actively driven, **other FO sessions sharing this state checkout stay out of its members** — they report status and work unrelated entities rather than dispatching into a member and colliding on its worktree/state. Until an entity carries a machine-readable owner, ownership is coordinated out-of-band (the captain / a handoff); the durable per-entity owner signal is the graduation trigger tracked by `79` — `xp` covers only the FO/Commander message channel, not mutation authority.

## File Naming

Each task is a folder or markdown file named `{slug}` or `{slug}.md` - lowercase, hyphens, no spaces. Use folder-form entities when reports or artifacts may accumulate beside the task. Example: `native-go-status/index.md`.

## Schema

Every task file has YAML frontmatter. Fields are documented below; see **Task Template** for a copy-paste starter.

### Field Reference

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique 24-character Spacedock Base32 ID because this workflow uses `sd-b32` |
| `title` | string | Human-readable task name |
| `status` | enum | One of: backlog, ideation, implementation, validation, done |
| `source` | string | Where this task came from |
| `started` | ISO 8601 | When active work began |
| `completed` | ISO 8601 | When the task reached terminal status |
| `verdict` | enum | PASSED or REJECTED - set at final stage |
| `score` | number | Priority score, 0.0-1.0 (optional). Workflows can upgrade to a multi-dimension rubric in their README. |
| `worktree` | string | Worktree path while a dispatched agent is active, empty otherwise |
| `issue` | string | Optional external ticket reference, such as `ENG-123`, `kata:task-abc123`, or `owner/repo#42` |

## Proof policy

The FO contract's Operating Principles and Working Principles already govern generic proof discipline: prefer the cheapest check that can fail — a shipped guard's run, an existing mechanical check, a one-off falsifiable exercise recorded in the report, then the captain's judgment — with new standing enforcement as the last resort rather than the default; spike the riskiest mechanism first; and let the gate AC cross-check reject any criterion whose only proof is review of its own prose. Tasks in this workflow inherit that. The rules below add the dev-workflow specifics.

- **No prose-grep over instruction files.** A string, substring, or regex match over an instruction file the model reads (the FO/ensign contract, this README, a skill) never proves a behavioral claim. The matched text was written by the same implementer the check polices, so it asserts only that the file contains what we put in it. A valid paraphrase fails it and an inverted clause passes it. To settle a case, ask whether the expected value comes from outside the file under test; if it does not, the check is a tautology and is banned. A check that binds two independent values that can diverge, such as the plugin manifest's version sharing a major.minor with the binary's version, is legitimate and is not prose-grep. Captain ruling (2026-07-20, verbatim): prose-greps are one-off validation evidence, never committed tests. A grep whose output is pasted into the validation report is legitimate external evidence for that run; the same grep committed as a test is banned — it re-asserts that the file contains what we wrote and cannot fail. What such a grep may be evidence FOR is bounded by honesty, not by category: presence or absence is an existence fact, and a grep establishes it soundly when that fact is itself the claim. When it is not — when the claim is about what a program or an agent actually does, and the words being present says nothing about that — the grep is misleading, and the answer is to express the claim in a form that can be exercised, not to let a grep that technically passes stand in for proof it cannot give.
- **Evidence must be able to fail.** Each AC's cited evidence names the concrete change that would flip it — the falsifying edit. An author who cannot name what would make the evidence fail has not shown it can fail, and the criterion does not count. This is the author-side obligation behind the validation bar below (a static check counts only when it tests a real value against an independent source that can diverge); the gate reads that falsifying change, not a pass count.
- **Detached adversarial audit (high-stakes surfaces).** Before merging a change to one of four high-stakes surfaces (the front-door launcher, the `status` mutation and guard paths, the shipped contract and scaffolding, and the CI and release machinery), run a read-only audit on a throwaway checkout. It tries to refute the validation by constructing an adversarial edit the deliverable's own tests should catch, then confirms they do. A test that stays green under a claim-breaking edit is a hole. Material findings route back through validation to implementation, and "refuted nothing material" is a valid recorded outcome. The audit also fires on AC provenance: when an AC's expected value is derived from the same package's production functions or constants, run the adversarial-edit check on it — that provenance is the tautology tell. Scope it to that provenance form; the broader equality/byte-identity form over-fires on ordinary unit tests. The two triggers are independent: the four-surface trigger runs the full audit on a throwaway checkout, while the provenance trigger fires wherever such an AC appears — including on a change routine enough to skip the full audit — and covers only that AC's adversarial-edit check, not the whole change.
- **Required CI lanes are a function of the diff, not the FO's read of "relatedness."** Merging on the deterministic lanes (build/install/offline) alone is allowed ONLY when the diff provably touches nothing a live lane loads or drives. When the diff touches a file a live lane exercises — the shipped FO/ensign contract or a host adapter (`skills/**/references/**`), the dispatch/launch path, or the lane's own live tests — that lane is REQUIRED green before merge, and a flake there is grounds to re-run to green (serial, isolated), NEVER to skip, leave its deployment unapproved, or wave off as "the known flake." The path→lane mapping is the gate: a change to the Claude adapter requires `claude-live`; to a host adapter, that host's lane; to the host-neutral dispatch core, every host lane. The generic proof bar this realizes — what verification a change requires, when a result counts as a genuine green, and how a failure is read — lives once in the FO contract's Working Principles ("Hold your own gate, merge, and triage calls to the bar you impose on workers"); this bullet is its dev-lane realization, not a re-derivation.
- **Instruction-file read quarantine.** Tests do not read prompt or instruction files except in `internal/contractlint`, and there only for structural checks: reference closure, frontmatter validity, structural absence, and dedup. Prose-grep and prose-to-code consistency checks never substitute for running the behavior.
- **Trace every mechanism to value.** At implementation dispatch, and whenever implementation introduces another mechanism, name the user-visible/value AC it serves, the simplest available alternative, and the concrete reason that alternative cannot prove or deliver the value. A test harness is a mechanism: it should orchestrate and observe the supported runtime, not acquire its own protocol, daemon, lease, lifecycle state, recovery loop, or process controller and become a second implementation of the system under test. For terminal multiplexers, `setsid`, process-group control, raw PTY writes, or a second lifecycle supervisor require an explicit architecture review; try an existing real-terminal substrate first. Before approving an enabling task, require one end-to-end run through the simplest substrate. If it has not produced visible product proof within 90 minutes, stop and review the architecture instead of adding another coordination layer.

## Stages

### `backlog`

A task enters backlog when it is first proposed. It has a seed description but no design work has been done yet.

- **Inputs:** None - this is the initial state
- **Outputs:** A seed task file with title, source, brief description, acceptance criteria, and stage-specific test gates
- **Good:** Clear enough to understand what the task is about and what proof future stages must provide
- **Bad:** Mixing launcher, status, skill integration, and tracker work without a testable boundary

### `ideation`

A task moves to ideation when a pilot starts fleshing out the idea: clarify the problem, explore approaches, and produce a concrete description of what "done" looks like.

- **Inputs:** The seed description and any relevant context, including existing code, user feedback, related tasks, and current Spacedock behavior
- **Outputs:** A fleshed-out task body with problem statement, proposed approach, acceptance criteria, and a test plan
  - Acceptance criteria must include how each criterion will be tested.
  - The task body declares an **expected surface** — the files and LOC it expects to touch — and its tolerance, as part of the gated design. The ideation gate approves it as the baseline every later correction round calibrates against.
  - Acceptance criteria are **entity-level** - they describe properties of the finished task, not stage actions. Items that describe stage work belong in the stage report's checklist.
  - If an AC item reads as an imperative verb phrase, rewrite it as the end-state property it produces.
  - At least one AC must MEASURE the end-value the entity exists for, against an independent baseline that can move the wrong way (a number/delta/count/timing, a behavior, or resulting on-disk state). An AC that only asserts its mechanism shipped — "the prose updates to X", "the verb owns Y", "the section is rewritten to Z" — is end-state phrasing of a *means*; it counts only paired with the value-measuring AC it serves. This extends "prove behavior by exercising it" to the entity's reason-for-existing: "the mechanism shipped" is satisfied only by "the outcome it exists for is measured" (cf. `trim-dispatch-adapter-prose` AC-1: cumulative line delta vs origin/main is NEGATIVE).
  - Every task must produce a real, checkable change (code, a fixture, on-disk state, or instruction text whose effect a separate check can confirm). If the task's only output is a decision with nothing shipped, it does not belong in this queue; record the decision in the roadmap instead. Cleanup and overhaul qualify: the change is the new code plus passing tests.
  - When the design rests on an unverified mechanism (a parser round-trip, a runtime handoff, an on-disk format, a tool actually supporting a flag), spike the riskiest path first (see Proof policy above) and record the result in the task body. The throwaway exercise seeds the implementation's first test. If nothing is unverified, record "no spike needed: {the proven mechanisms it relies on}" so the determination is on the record.
  - Test plans should state what verifies the implementation, estimated cost/complexity, and whether fixture, CLI, or live workflow tests are needed.
  - Plans should describe intended behavior at the level a future worker or validator needs to reason about it. Prefer observable behavior over implementation internals unless the task is specifically about that internal representation.
  - For every new mechanism in the proposed approach or test plan, name the value AC it serves, the simplest alternative considered, and why that alternative is insufficient. An enabling mechanism is not justified by proving its own internals.
  - Prove behavior by exercising it and observing the outcome (output bytes, exit code, resulting on-disk state, or a test feeding many inputs and asserting uniform handling): Go unit tests for parser and command behavior, golden fixtures for status output, behavior fixtures that drive the binary for command-level claims, and live workflow smoke tests only when runtime behavior is the claim. See the Proof policy above for what counts as proof and what does not.
  - When captain feedback changes the target behavior, update the task body, acceptance criteria, and test plan together before re-validating.
  - For template or skill text changes: specific before/after wording, not just "change X".
  - When a task changes user-visible behavior — CLI output, command surfaces, startup banners, host integration, anything the docs site describes — ideation proposes the documentation changes: a concrete doc diff (before/after wording or a unified diff against the affected doc files) recorded in the task body. The ideation gate review includes this doc diff. Implementation applies it. Ideation runs without a worktree; the diff lives in the task body, not on a branch.
- **Good:** Clearly scoped, behavior-first, actionable, addresses a real need, considers edge cases, avoids unnecessary runtime-internal modeling, and uses tests that prove the intended behavior directly
- **Bad:** Vague hand-waving, scope creep, solving problems that do not exist yet, no clear definition of done, acceptance criteria without a test plan, static prose tests for behavioral requirements, or tests that pass while missing the intended behavior
- **Staff review:** When the FO assesses ideation as complex, such as native status parity, split-root behavior, or skill integration, it should request an independent review before presenting the ideation gate. The review checks design soundness, test plan sufficiency, gaps, and that the riskiest unverified mechanism was exercised first (or that the task records an auditable "no spike needed" with the proven mechanisms it relies on). A design whose soundness rests on an unexercised, unverified mechanism is not ready for the gate.

### `implementation`

A task moves to implementation once its design is approved. The work here is to produce the deliverable: write code, generate fixtures, update skill instructions, or make whatever changes the task describes. Implementation is complete when the deliverable exists and is ready for independent verification.

- **Inputs:** The fleshed-out task body from ideation with approach and acceptance criteria
- **Outputs:** The deliverable committed to the relevant repo or state checkout, with a summary of what was produced and where
- Implementation completion is not a stopping point: once the deliverable is committed and the stage report filed, the entity routes immediately to independent `validation` dispatch — a fresh validator, since `validation` is `fresh: true` — unless a gate, blocker, terminal ceremony, or captain decision intervenes. The FO does not park a completed implementation and wait.
- During implementation, re-run the mechanism/value check before adding a protocol, daemon, lease, lifecycle state, recovery loop, process controller, or another coordination layer. Stop and request a design reset when the simpler supported-runtime path still reaches the value; do not repair the enabling mechanism merely because it already exists.
- When consuming an in-stage review round's findings (a roborev panel, a detached-audit pass, or routed gate feedback), triage before fixing and record the triage as an advisory resolution on the round's briefing—advisory only, because a round never advances status. Use the `validation` release-scope classification (Material / Deferred risk / Polish / Needs decision) and `.roborev.toml`'s four-field evidence record (released user and workflow; observable harm; affected value AC or non-negotiable boundary; trigger evidence). Fix material findings; for a correct-but-disproportionate finding, record a decline naming its class, why it is not material, and its promote-to-material condition; escalate a needs-decision finding to the FO. A value AC narrowed to make a finding pass requires a captain-owned binding resolution through a gate attempt, not a fix (see the entry's AC-drift note).
- After the reviewer and worker entries are complete, invoke `${SPACEDOCK_BIN:-spacedock} gate record <entity> --round <stage>/<cycle> --briefing <path>/briefing.json --log <path>/briefing.review.jsonl --feedback-cycle <file>`; the round records no gate application and does not change status. Omit `--feedback-cycle` for a complete no-findings round.
- When an in-stage review round (a roborev re-panel, a detached-audit pass, or any similar in-stage reviewer loop) returns findings that trigger another pass, first append a conforming `### Feedback Cycles` entry — `- Cycle {N}: {verdict} — {reviewer}; surface {files}/{LOC} vs estimate {declared} ({P}%); AC {unchanged | narrowed: <note>}` — with the actuals from `git diff --numstat "$(git merge-base main HEAD)"..HEAD`. Read the deviation against the ideation estimate from the entry; past the declared tolerance, record a design-reset decision (see the bullet above) before another pass. In-stage and cross-stage rounds share one section and one convention, the one `feedback-rejection-flow` carries.
- **Good:** Minimal changes that satisfy acceptance criteria, clean Go packages, stable CLI output, tests that observe supported behavior through the simplest real substrate, and a self-contained deliverable
- **Bad:** Over-engineering, a test harness that becomes a second implementation of the runtime, another coordination layer without visible product proof, unrelated refactoring, skipping tests, ignoring edge cases identified in ideation, or leaving the deliverable incomplete for validation to finish

### `validation`

A task moves to validation after implementation is complete. The work here is to verify the deliverable meets the acceptance criteria defined in ideation. The validator checks what was produced - it does not produce the deliverable itself.

- **Inputs:** The implementation summary and the acceptance criteria from the task body
- **Outputs:**
  - Run applicable tests from the Testing Resources section and report results.
  - Verify each acceptance criterion with evidence.
  - Pull every `**AC-N**` item (including a value annotation inside the bold, e.g. `**AC-1 (VALUE)**`) from the entity body's `## Acceptance criteria` section; reproduce the evidence cited in each "Verified by" clause; flag any AC without evidence.
  - Reproduce each AC's cited evidence; reject any AC whose evidence is self-referential, or whose only deliverable is a decision with nothing shipped (that belongs in the roadmap, not a terminal dev task). Dev-workflow policy: an AC's proof is code, command, or state. A non-development workflow's AC proof may legitimately be a published artifact, a metric, or a human review.
  - Check that the task body, acceptance criteria, implementation, and tests reflect the latest captain feedback.
  - Reject when tests pass but prove an obsolete, over-specified, or wrong target behavior.
  - Before requesting review, perform one **semantic adversarial pass** over the changed behavior:
    - Trace each changed value or event through every representation and lifecycle phase. Check exact identity, cardinality, order, bytes, attribution, authority, and terminal state—not only field presence or counts.
    - Build a compact matrix of adjacent variants: empty and terminal states, repeated or out-of-order events, every input path, and relevant Unicode, EOF, size, visibility, and layout boundaries. Verify one invariant across the matrix.
    - When validating an existing format or protocol, use its canonical validator when possible. Otherwise validate the complete record atomically; do not add validation one field or failure case at a time.
    - Inspect changed hot paths and readers for multiplicative work, blocking I/O, unbounded allocation, and implicit size limits. Add one scaling or over-limit test when that risk exists.
    - Ask, “How could this test pass while the observable behavior is wrong?” Assert the exact result and the failure or cleanup behavior that distinguishes it.
  - Classify a rejection as either a defect in the intended product or a failure of the chosen mechanism. When the mechanism or harness architecture is wrong but the end value remains reachable by a simpler route, recommend a scope/design reset. Do not send it through another automatic implementation feedback cycle to repair the failed mechanism.
  - Classify every finding on two independent axes:
    - **Defect kind:** an **outcome defect** means delivered behavior or state fails a value AC; an **evidence defect** means the test, harness, or observation boundary cannot validly establish an AC. A “narrow fix” must cite the affected AC and exact failing boundary. If the correction spans both kinds, or an evidence fix requires another controller or lifecycle layer, recommend a design reset.
    - **Release scope:** a **material** finding has a supported, promised, common, or observed trigger and violates a value AC or non-negotiable safety, security, data-integrity, or compatibility boundary. A **deferred risk** is real, but its trigger is outside the current supported or promised workflow and no value AC fails under normal use.
    - Severity and release scope are independent. A High or Medium finding is not automatically material. “Edge case” alone is not a classification.
    - A deferred risk must record its exact trigger, why that trigger is outside the current promise, evidence that the supported path still satisfies its AC, and the concrete condition that promotes the risk to material. If the intended workload is unknown, measure or probe it before escalating.
    - Only material findings block the gate or consume a feedback cycle. When all promised ACs have valid evidence and no material finding remains, recommend PASSED and list any deferred risks with their revisit conditions.
  - A PASSED/REJECTED recommendation, with deferred risks listed separately from material and polish findings.
- **Good:** Thorough testing against acceptance criteria, clear evidence of pass/fail, honest assessment, and validation that tests prove the current intended behavior
- **Bad:** Rubber-stamping without testing, ignoring failing edge cases, validating against wrong criteria, accepting passing tests that encode stale prose or obsolete assumptions, or accepting a string/substring/regex match over an instruction file (the contract, this README, a skill) as proof of a behavioral claim. A check whose expected value is just the text the implementer wrote into the file under test proves nothing — it cannot fail. Proof of behavior must run the behavior and observe output, exit code, or on-disk state; a static check counts only when it tests a real value against an independent source that can diverge from it, not as a spelling check over a file the model reads.
- **Spot-check principle:** Before committing to an expensive live workflow or compatibility run, do a cheap fixture or single-command spot-check to verify the infrastructure works end-to-end.
- **Detached adversarial audit:** for the high-stakes surfaces named in the Proof policy above, run (or dispatch) the audit on a throwaway checkout, never the implementation worktree, before merging. Routine, low-blast-radius changes do not need it, except that the AC-provenance trigger in the Proof policy applies to them too, on that AC alone. Record material findings as a `### Feedback Cycles` entry naming the audit and its adversarial edit; do not present the gate as clean until they close. Note a clean audit in the gate's reviewer-findings block. Real catches on the record: #262 (two test-strength holes in `contract_gate_test.go`), `1x` and `external-tracker-checkpoint` AC-6 (self-referential ACs that can never fail), and `7h` AC-3 (a tag-cut that folded the release notes into the tag subject).

### `done`

A task reaches done when validation is complete and the captain approves the result. The task is closed with a verdict of PASSED or REJECTED.

- **Inputs:** The validation report with PASSED/REJECTED recommendation
- **Outputs:** Final verdict set in frontmatter, completed timestamp recorded
- **Good:** Clear resolution and lessons learned captured if relevant
- **Bad:** Closing without reading the validation report, overriding a REJECTED recommendation without reason, or marking PASSED a task that shipped nothing checkable (see Proof policy above). A design that concludes "do not build X" ships as a roadmap decision, not a PASSED dev-queue task. A contract or skill change is PASSED only when a live drive observed the behavior it claims.

## Workflow State

Entities live directly under `.spacedock-state/` (no `entities/` directory); the workflow README stays in the main repo. Workflow state is read from `.spacedock-state`. Read it with the launcher:

```bash
spacedock status --workflow-dir docs/dev
```

To list the tasks ready for dispatch (the query the first officer runs each loop):

```bash
spacedock status --workflow-dir docs/dev --next
```

### Reading sections

**Read one section, not the whole file.** `spacedock status --workflow-dir <wf> --read <entity-ref-or-path> --json` returns the file's frontmatter plus an ordered heading map (`text`, `level`, `offset`, `lines`). Pass a heading's `offset`/`lines` to `Read(path, offset, limit)` to load just that section — e.g. an entity's latest `## Stage Report`, or this README's `## Sprints`, without the rest.

## Runtime Live CI

The live runtime lanes prove host behavior, not text shape: a static grep over workflow YAML or skill prose never substitutes for launching the real host front door, observing its output, and checking the resulting workflow state. The full reference (shared scenarios, fixtures, assertions, per-host runners, local live execution, and the GitHub setup) lives in [`docs/runtime-live-ci.md`](../runtime-live-ci.md). Add or change a runtime scenario there.

## Task Template

```yaml
---
id:
title: Task name here
status: backlog
source:
started:
completed:
verdict:
score:
worktree:
issue:
---

Brief description of this task and what it aims to achieve.

## Problem

{What is broken or missing, and why it matters. Ideation fills this in.}

## Proposed approach

{How the task intends to solve the problem. Ideation fills this in.}

## Out of scope

{What this task deliberately does not cover, so the boundary is explicit.}

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - {End-state property.}**
Verified by: {test name / command output or exit code / file the change produces / resulting on-disk state — something outside this task body that a future reader can reproduce and that can fail; name the concrete change that would make it fail.}

## Test plan

{What verifies the implementation, estimated cost/complexity, and whether fixture, CLI, or live workflow tests are needed.}
```

## Testing Resources

Validation pilots should use these when verifying implementation work:

| Resource | Command or Path | Covers |
|----------|-----------------|--------|
| Go unit suite | `go test ./...` | CLI routing, parser behavior, status implementation, fixtures |
| Race-enabled Go suite | `go test ./... -race` | Concurrency hazards in Go code when relevant |
| Clean log + `-json` archive from one run | `gotestsum --jsonfile detail.jsonl --format pkgname -- ./pkg` | The live-CI test-output shape: a clean step log (per-package progress + an `=== Failed` recap with `file:line`) plus a full `-json` archive, from a single run, with the `go test` exit preserved. Locally: `go install gotest.tools/gotestsum@v1.13.0` (the version CI pins), or run `.github/scripts/install-gotestsum.sh` for the same sha256-verified prebuilt. Inspect the archive with `grep '"Action":"fail"' detail.jsonl` or a `go tool test2json`-aware reader. |
| Launcher help smoke test | `go run ./cmd/spacedock --help` | Basic command entrypoint behavior |
| Launcher version smoke test | `go run ./cmd/spacedock --version` | Basic version output behavior |
| Status validator | `spacedock status --workflow-dir docs/dev --validate` from the repo root, or pass an absolute workflow definition dir | Spacedock entity-contract validation; fails closed if `--workflow-dir` does not resolve to a commissioned workflow |
| Status table | `spacedock status --workflow-dir docs/dev` | Status enumeration output |
| State behavior extension | `docs/specs/state-behavior-extension.md` | Split-root state semantics and external tracker bridge principles |
| Bootstrap roadmap | `docs/roadmap/bootstrap-roadmap.md` | Stage-specific required tests |

Validators should pick the smallest test surface that proves the claim. Use Go unit tests for package behavior, golden fixtures for stable command output, and live workflow smoke tests only when the runtime integration itself is the claim.

## Commit Discipline

- Commit state changes at dispatch and archive boundaries.
- Commit task body updates when substantive.
- Keep main repo changes and `.spacedock-state` changes in their respective git repositories.
