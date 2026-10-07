---
title: The commission skill lets an initial seed gate claim to review work no worker produced
status: implementation
source: "Captain, 2026-10-06: incident while commissioning an exploratory workflow on 0.28.0-pre4. An `exploration` stage was declared both initial and gated, and its prose implied the gate reviewed research no worker had run. Approving it would have advanced straight past the research worker."
id: 6dhxxd047s3c32n21h8ppf9y
gates:
    version: 1
    records:
        - id: gate:6dhxxd047s3c32n21h8ppf9y:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:6dhxxd047s3c32n21h8ppf9y-backlog-1
              briefing:
                id: briefing:6dhxxd047s3c32n21h8ppf9y:backlog:attempt-1:revision-1
                digest: sha256:a3c37779e5cbc8beb54e51b618936ef0718fe5fe474b989d565b69b2bbaffb1c
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:6dhxxd047s3c32n21h8ppf9y:backlog:1
                briefing: briefing:6dhxxd047s3c32n21h8ppf9y:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-10-06T23:39:04.936403Z"
                decision: approve
                reason: 'Captain, 2026-10-06, in session: "file the commission issue first and dispatch to astra for ideation"'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:6dhxxd047s3c32n21h8ppf9y:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:6dhxxd047s3c32n21h8ppf9y-ideation-1
              briefing:
                id: briefing:6dhxxd047s3c32n21h8ppf9y:ideation:attempt-1:revision-1
                digest: sha256:398b0bcf219806d53062113804426a2b3cb1811d92730b30543bbafb9dd9af3b
                room-ref: '@review/ideation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:6dhxxd047s3c32n21h8ppf9y:ideation:1
                briefing: briefing:6dhxxd047s3c32n21h8ppf9y:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-10-07T00:33:15.295714Z"
                decision: approve
                reason: 'Captain approved the presented ideation gate in session: one-file commission guidance, net +6, incident-shaped blind trial as limited evidence, live proof deferred to implementation'
              application:
                target-stage: implementation
                state: consumed
started: 2026-10-06T23:39:26Z
worktree: .worktrees/spacedock-ensign-commission-seed-approval-vs-result-review
---

## Problem

An initial gate can promise review of research that no worker has produced. Its approval advances to the successor, not into initial-stage work.

The captain's incident report describes `exploration (initial, gate) -> prototype` on `0.28.0-pre4`. The captain stopped that proposal before dispatch. It represents one proposed research bypass, not an observed skipped worker. The requested route was `backlog (initial, gate) -> exploration (gate) -> prototype`.

This task improves commissioning guidance. It does not mechanically prevent false claims or guarantee compliance. Acceptance requires zero observed research bypasses and zero false completed-research claims in the specified generation runs.

`internal/status/entered_stage.go` confirms that an initial gate needs a tracked, clean seed. A later gate needs a complete, committed current-stage report. These checks establish durability and report structure, not substantive research quality. `docs/dev/README.md` legitimately combines `initial: true` and `gate: true` for backlog permission to start design.

Presentation precedes decision mutation and can prevent bad advancement. It occurs after generation and cannot prevent a bad workflow from being generated. This task addresses generation only.

## Proposed approach

Change only `skills/commission/SKILL.md` during implementation. Teach one gate-timing rule, apply it during interactive and batch confirmation, and reconcile generated declarations against confirmed decisions.

The following exact before/after wording is the proposed deliverable. No skill file changes in this ideation round.

### Question 2 — Stages

Before:

> Store the confirmed stages as `{stages}`. The first stage is `{first_stage}` and the last is `{last_stage}`.

After:

> Store the confirmed stages as `{stages}`. The first stage is `{first_stage}` and the last is `{last_stage}`.
>
> **Gate timing:** An initial approval gate reviews the committed seed, not completed work assigned to that stage. No worker runs for that initial stage before seed approval. Approval advances to its successor. A later approval gate reviews its stage's completed, committed report and required artifacts.
>
> Separate permission to start research from permission to build from research results. When both approvals are requested, propose a seed stage, a research stage, then the next work stage. Put seed approval on the initial stage and result approval on the research stage. Name the reviewed artifact and successor work for each approval. Do not describe seed approval as review of research that has not run. If the supplied stages conflict with these decisions, explain the conflict and propose the separated route during Confirm Design.

### Confirm Design — approval defaults

Before:

> - `{approval_gates}` — default: gate before the terminal stage (e.g., the last stage before terminal).

After:

> - `{approval_gates}` — default: approval at the last stage before terminal. Apply the Gate timing rule from Question 2, including in batch mode. For each approval, name its reviewed artifact and successor work. Preserve requested research between seed approval and research-result approval.

### Confirm Design — captain-facing summary

Before:

> {for each stage: "{letter}. {stage_name} — {stage_description}"}

After:

> {for each stage: "{letter}. {stage_name} — {stage_description}"}
> {for each approval: "At {stage_name}, you review {artifact}. Approval starts {successor_work}."}

### Generated README guidance

Before:

> For every gated stage, add `- **Gate content:**` to its stage subsection and state the evidence needed for that decision. This rule also applies to custom stages and template variants.

After:

> For every gated stage, add `- **Gate content:**` to its stage subsection. Name the reviewed artifact, required evidence, and actual successor work. Apply the Gate timing rule to custom stages and template variants. For an initial gate, name the committed seed and permission to start successor work, not that initial stage's worker results. For a later gate, name its completed, committed stage report and required artifacts.

### Generation Checklist

Before:

> - [ ] Each seed entity file exists at `{dir}/{slug}.md` with valid YAML frontmatter

After:

> - [ ] Each seed entity file exists at `{dir}/{slug}.md` with valid YAML frontmatter
> - [ ] Reconcile confirmed approvals, reviewed artifacts, and successor work with generated stage declarations and Gate content, including requested initial approvals.

This is one author reconciliation check, not an independent validator. Do not add a recurring worked-example validation obligation. Implementation-added runnable examples still need validation in the implementation test plan.

## Expected surface and tolerance

Estimate net LOC change: +6, across 1 file.

| Deliverable file | Insertions | Deletions | Net |
| --- | ---: | ---: | ---: |
| `skills/commission/SKILL.md` | 8 | 2 | +6 |
| Total | 8 | 2 | +6 |

The estimate uses the literal paragraph layout above. Tolerance: net +4 to +18 LOC, exactly one deliverable file. Measure with `git diff --numstat` against the approved implementation base. State-body reports are excluded. This round changes only this entity body.

### Declared semantic changes

- **Instruction behavior:** commissioning separates seed permission from result approval in interactive and batch paths.
- **Generated output:** confirmation and README Gate content name reviewed artifacts and actual successors; declarations preserve requested approvals and work.
- **Command grammar and stored formats:** unchanged. No new schema, command, report format, or gate-record behavior.
- **Authority:** unchanged. The captain confirms the design and makes approval decisions under existing delegation rules.
- **Runtime behavior:** unchanged. Binary readiness, dispatch transport, host adapters, and presenter behavior remain intact.

## Out of scope

No presenter edits or criteria for repairing an already-existing malformed workflow during presentation. No binary edits, new lint, CI lane, template migration, startup redesign, host adapter changes, or ban on initial gates.

This task neither fixes presenter reload churn nor depends on that fix. Existing lifecycle and evidence rules remain unchanged. Future residency changes must preserve fresh gate evidence.

No repository product edits, worktree, PR, main push, or live workflow mutation in this ideation round. Throwaway exercises are outside the deliverable. Only this state entity is committed and pushed to `spacedock-state/dev`.

## Acceptance criteria

**AC-1 (VALUE) — Newly commissioned workflows preserve the requested research work.**
Both interactive and batch generation produce zero research bypasses against the fixed requested order. The independent oracle is seed approval, research with a reference scout, research-result approval, then prototype. The incident supplies one proposed bypass, not an observed baseline execution.
Verified by: two targeted live generation-to-FO drives using the unchanged Pi runtime. After seed approval, an observed worker spawn receives research and the reference scout assignment. No prototype spawn occurs before research completion and approval. After research-result approval, an observed prototype spawn occurs. Retain generated declarations, assignments, approvals, reports, commits, and runtime events. Making research initial and routing seed approval directly to prototype fails. A built assignment alone does not prove a spawn.

**AC-2 (VALUE) — Commissioned approval descriptions make zero false completed-research claims at the seed boundary.**
Confirmation summaries and generated Gate content identify the seed as the first reviewed artifact and research as successor work. They identify the completed research report and required artifacts as the later reviewed evidence, with prototype as successor. Both modes preserve legitimate initial approval gates.
Verified by: inspect the actual generated outputs from AC-1 against the independently fixed request and absent seed-stage worker evidence. Include a conflicting proposed route as generation input. The generated design must explain and correct that conflict before generation is accepted. A seed description asserting that researched findings already exist fails, even when declarations are valid. Validation must also run `status --validate` on both generated workflows. No criterion requires changing a presenter or repairing a previously commissioned workflow.

## Test plan

The primary proof owner is a targeted live commissioning exercise. No existing generation-behavior test owner was identified in the permitted reading surface. The status-owned report guard in `internal/status/entered_stage.go` remains unchanged; new readiness probes do not prove this authoring change and are not added.

Fix this request before generation:

> Create a workflow to assess an idea. I must approve its seed before research starts. Research must include a reference scout. I must approve the research results before any prototype work starts. A draft suggests exploration as the initial approval stage, followed by prototype, and says the first approval reviews researched findings. Reconcile that draft with my requested decisions.

Use the same request in both modes. For interactive mode, answer ordinary questions consistently with it. For batch mode, supply all inputs and a disposable inline location; retain explicit approval at each boundary. Use a single idea seed, no mods, and no repository mutation by research or prototype workers. Have workers write their assigned reports in the disposable workflow.

Run a pre-change generation smoke before editing the skill. Repeat against changed instructions and preserve both outputs; the old skill need not fail on demand. The incident is historical context; the fixed request is the oracle. Independent validation replays the changed behavior, rather than trusting the author's pass count.

| Check | Distinct falsifying change | Cost and mode |
| --- | --- | --- |
| AC-1 requested route and actual work | Make research initial and advance directly to prototype | Two live runs, interactive and batch, approximately 10–20 minutes each |
| AC-1 batch application | Apply gate timing only through interactive questions | Same batch run; no second harness |
| AC-2 truthful approval descriptions | Correct declarations but retain a seed claim of completed research | Inspect both runs' actual confirmation and README output; minutes, manual |
| AC-2 valid initial approvals | Omit the requested initial gate or emit invalid declarations | Compare generated declarations with request; run `status --validate`; seconds, CLI |

Retain commands, binary and host versions, skill revision, output, exit status, runtime events, and durable state commits. Schema validation proves schema validity, not instruction efficacy or actual dispatch. No permanent prose-grep test or new test harness is proposed.

Run `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal` during implementation verification. Confirm formatting causes no unrelated changes. Apply existing review and live-lane policy to the eventual one-file diff. These broad repository checks do not prove ideation behavior and are not run in this state-body-only round.

### Necessity

- Gate-timing guidance serves AC-1 and AC-2. Presenter-only correction can stop advancement but cannot ensure correct generated work order.
- Confirm Design covers batch mode, which skips Question 2. One canonical rule avoids duplicating timing policy.
- The reconciliation check compares confirmed intent with generated output. Guidance alone does not request that comparison; multiple checklist rows add no independent validation.
- Live generation and dispatch observations serve AC-1. Static schema checks cannot establish worker order.
- Manual inspection of generated descriptions serves AC-2. Correct declarations alone cannot establish truthful claims about absent research.

## Risk evidence and spike disposition

The principal risk is instruction efficacy, not parsing. The correction brief is `docs/dev/.spacedock-state/.briefings/commission-seed-approval-ideation-astra-review.md`. Its presenter recommendations are superseded by the captain's one-file scope. Its demand for a focused instruction trial remains applicable.

### Pi entry path and preflight

The commission source names the Claude runtime and Claude-shaped boot tools. The captain reports a successful Pi commission run on `0.28.0-pre4`. This is observed success reported by the captain, not a failure inferred from prose. The host-specific wording is a text gap, not an established functional block.

The supervising FO also reports this local Pi dispatch-builder output: `declared model 'opus' ignored on host pi`. A declared model outside Pi's settable space is ignored. This is separate from the captain's commissioning observation; no exit status or full handoff transcript was supplied for either observation.

Implementation enters commission from Pi, generates the disposable workflow, then loads the existing first-officer skill in that root session. It uses the existing Pi dispatch binding, without changing generated stages or pretending Claude tools are available. Retain the actual host events during implementation. Do not call separate generation and later Pi consumption a continuous handoff.

The FO explicitly accepted these observations as the ideation preflight; no further host preflight is required this round. Full live handoff and spawn proof belong to implementation. Any newly observed runtime incompatibility requires a scope decision, not an unapproved startup patch.

### Self-applied generation trial

The FO authorized a self-applied trial because this worker cannot launch subagents. I applied the candidate instructions to the fixed conflicting-route request. I generated `README.md`, `idea.md`, and a confirmation response in an isolated directory, outside the repository. This was generation by the current worker, not a blind model invocation.

The shell allocated the directory with `SPIKE=$(mktemp -d -t spike-commission-seed-approval-vs-result-review)`. A `python3` heredoc wrote the worker-generated files and printed the confirmation. No script inferred the expected route from the candidate wording. The generated declarations were:

```yaml
stages:
  defaults:
    worktree: false
  states:
    - name: backlog
      initial: true
      gate: true
    - name: exploration
      gate: true
    - name: prototype
    - name: done
      terminal: true
```

The seed had `title: Assess an idea`, `status: backlog`, and `source: commission seed`. The README used `id-style: slug` and `state: $inline`. Exploration outputs required research findings, a reference scout, and cited sources in the stage report. Prototype outputs were a plan and report, not repository changes.

Observed confirmation output:

```text
The draft puts research in the initial approval stage. Initial approval reviews a seed and advances to its successor; that draft skips research. Use backlog, exploration, prototype, done instead.
At backlog, review the committed seed. Approval starts exploration with a reference scout.
At exploration, review the completed, committed research report and cited sources. Approval starts prototype work.
```

Generated seed Gate content was: "Review the committed seed. No exploration worker has run. Approval starts exploration with a reference scout."
Generated exploration Gate content was: "Review the completed, committed exploration report and cited sources. Approval starts prototype work."

Executed validation command, with `SPIKE` bound to that disposable directory:

```sh
BIN=/opt/homebrew/Caskroom/spacedock@next/0.28.0-pre3/spacedock
"$BIN" --version
"$BIN" status --workflow-dir "$SPIKE" --validate
result=$?
printf 'VALIDATE_EXIT=%s\n' "$result"
exit "$result"
```

Observed version: `spacedock 0.28.0-pre3`, `darwin/arm64`, runtime `pi (PI_CODING_AGENT)`, contract 3. Observed validator output: `VALID` and `VALIDATE_EXIT=0`; shell exit status 0. The pinned executable is supplied by this assignment. For replay elsewhere, select an installed executable and record its version rather than assuming this local installation exists.

The generated route contains zero proposed research bypasses. The confirmation and Gate content contain zero completed-research claims at the seed boundary. These are observations of this worker's generated output, not mechanical truth guarantees.

Limit: you know a test exists, so this is a weak signal. It is not a blind trial. The validator establishes schema validity only. This trial did not execute approval, spawn workers, or measure interactive/batch reliability. The available binary is pre3; the captain's successful Pi commissioning observation used pre4.

A separate in-memory application of the exact proposed wording measured 8 insertions, 2 deletions, net +6. The command used `python3` with `difflib.ndiff`; exit 0. It wrote no product file. This establishes the surface estimate, not behavior. The inspected code revision was `c450d5d526fca30bbcd2c06cf2b40358e5a6ea8f`.

The pinned executable also ran `status --workflow-dir docs/dev --validate`: `VALID`, exit 0. It reported an existing, unrelated `withdrawn` verdict warning for `live-lane-gap-inventory`. That entity was not edited. Frontmatter equality, the exact acceptance heading, and `git diff --check` passed. These are artifact checks, not behavioral proof.

### Supervisor-run fresh-worker trials

The supervising FO supplied the following evidence. I did not run these trials. The FO invoked one fresh child for each run: agent `worker`, model `openai-codex/gpt-6-astra`, context `fresh`, in an empty throwaway directory. Each child received only the candidate instruction text and its request. Neither received test framing, incident context, or repository access.

Invocation record: fresh-child dispatch with the parameters above, not a shell command. No Spacedock CLI command ran, so no CLI version applies to these trials. The model identifier is recorded above; no additional host version was supplied. Run `d850a1e4` and run `2ba124cf` both returned `exit=completed`. This is the supervisor's completion status, not an invented numeric shell exit code.

**Trial 1 — weak evidence, run `d850a1e4`.**

Request:

> Create a workflow to assess an idea. I must approve its seed before research starts. Research must include a reference scout. I must approve the research results before any prototype work starts.

Output supplied verbatim by the FO, abridged to decision-relevant parts:

```text
stages: seed (initial, gate), research (gate), prototype (worktree), done (terminal)
seed gate content: "Review the committed idea seed ... No worker runs before approval. Approval starts research, including a reference scout—not prototype work."
research gate content: "Review the completed, committed research report and required reference-scout artifact ... Approval starts prototype work."
approval summary: "At seed, you review the committed idea seed. Approval starts research, including a reference scout." / "At research, you review the completed, committed research report and reference-scout artifact. Approval starts prototype work."
```

Weakness: the request itself listed both approvals in order, so the text got the answer for free.

**Trial 2 — primary spike evidence, run `2ba124cf`.**

Incident-shaped request:

> Create a workflow for open-ended exploration and rapid prototyping. I want to pick a seed direction first, then research it, including scouting bb.app as a reference, and then review the findings before we prototype.

This request names neither gates nor two approvals. It retains the captain's requested work sequence.

Output supplied verbatim by the FO:

```text
stages: defaults {worktree: false, concurrency: 2}; states: seed (initial: true, gate: true), research (gate: true), prototype (worktree: true), complete (terminal: true); transitions seed->research->prototype->complete
seed gate content: "Review the committed seed: proposed direction, exploration goals, and open questions. Approval selects the direction and permits research, including scouting bb.app as a reference. No seed-stage worker runs before approval."
research gate content: "Review the completed, committed research report and supporting artifacts: source links, bb.app scouting observations, findings, uncertainties, and recommended prototype experiments. Approval permits rapid prototyping based on those findings."
approval summary: "Seed approval: Approve the committed seed direction → start research and bb.app scouting." / "Research approval: Approve the committed research report and supporting evidence → start rapid prototyping."
```

The FO observed separate seed and research-result approvals, research between them, and no seed claim of completed research. The produced route has zero proposed research bypasses and zero false completed-research claims at the seed boundary. Applied to the reported incident, this route preserves the research worker. No actual worker spawn or controlled before/after effect was measured.

FO-supplied limits, retained verbatim:

1. Lab trial. I supplied the candidate rules directly. It measures whether an agent follows the text, not whether a real commission session delivers the text intact, and not whether the FO obeys at gate time.
2. Request order was still named. Trial 2 removed the approvals, not the sequence. The sequence is what a captain normally states.
3. Schema validity is unproven. The output matches the grammar I supplied, but no real workflow exists and `status --validate` was not run on it. Do not record it as VALID.
4. One model. Both runs used openai-codex/gpt-6-astra. Other models may differ.

The self-applied trial exercises reconciliation of a malformed proposed route. Trial 2 is stronger evidence that the candidate rules elicit the intended distinction without explicitly requested approvals. Neither replaces implementation's full interactive and batch drives. Reuse both the conflicting draft and incident-shaped request for implementation smoke inputs; retain the primary request separately from the expected generated output.

## Historical report

The following prior report is retained for chronology only. Its two-file surface, presenter scope, AC-3, and blanket no-spike disposition are superseded by this design.

## Stage Report: ideation

- DONE: Read the commission-seed entity body and the backlog briefing it points at.
  Read the seed and committed briefing `fe360ace7:.briefings/commission-seed-backlog.md`; distinguished the stopped proposal from an actual skipped dispatch.
- DONE: Verify the initial-gate semantics in internal/status/entered_stage.go against the seed's claim, and confirm this repo's own backlog stage uses initial+gate legitimately.
  `gatePreparable` uses seed durability only for initial stages; `docs/dev/README.md` backlog requests permission to start design.
- DONE: Rewrite the task body: problem statement, proposed approach, expected surface (files + net LOC with tolerance), declared semantic changes, and risk evidence.
  Declared two deliverable files, net +40 LOC with +25 to +55 tolerance, unchanged binary semantics, and retained risk evidence.
- DONE: Provide `## Acceptance criteria` as the exact literal heading, entity-level, with at least one criterion that measures the end value against an independent baseline.
  AC-1 measures zero research bypasses against the incident's one proposed bypass; the requested order is independent of generated instructions.
- DONE: Give specific before/after wording for skills/commission/SKILL.md, not "change X".
  Proposed exact replacements for stage design, batch confirmation, generated Gate content, and the generation checklist.
- DONE: Give specific before/after wording for the initial-gate case in skills/present-gate/SKILL.md.
  Proposed explicit seed/no-worker/successor wording and contradiction handling, without weakening required review facts.
- DONE: Reconcile the initial-gate presentation with the Pi adapter's per-gate load rule if they conflict.
  Kept initial handling inside the presenter; prepare, commit, reread, load, and one root presentation remain unchanged.
- SKIPPED: Validate any worked example YAML with `status --validate` and record the result.
  No worked YAML is proposed; the existing `docs/dev` workflow validated with exit 0 and one unrelated verdict warning.
- DONE: Write the test plan: proof owner, a distinct falsifying edit per additional check, estimated cost, and whether validation is deterministic or live.
  Named live commission/FO proof ownership, fixed independent input, readiness probes, distinct falsifiers, and estimated costs; no prose-grep proof.
- DONE: Record either the spike result for the riskiest unverified mechanism or an auditable "no spike needed" with the proven mechanisms.
  Recorded source baseline, this entity's committed backlog-to-ideation transition, runtime contract, and validator output; instruction efficacy awaits implementation live proof.
- DONE: Report in the canonical item form, one DONE/SKIPPED/FAILED line per checklist item with an evidence or rationale line, ending with a non-empty Summary.
  This report covers all eleven dispatch items; only the entity body was edited and frontmatter was preserved.

### Summary

Specified a two-skill documentation correction that separates seed permission from research-result approval. Preserved legitimate initial gates and Pi's presenter load rule, with independent behavioral criteria and a bounded test plan. Live instruction efficacy remains implementation work; no product files, worktree, or PR were created in this round.


## Stage Report: ideation (cycle 2)

- DONE: Cut the deliverable to skills/commission/SKILL.md only, remove the present-gate change, and remove every criterion that needs the presenter to catch an already-existing malformed workflow.
  Proposed one-file authoring correction; AC-1 and AC-2 cover newly generated workflows only; the prior report is marked superseded.
- DONE: Rework or drop AC-2 so that every remaining acceptance criterion is measurable inside the one-file scope.
  AC-2 measures actual confirmation and generated Gate content against absent seed-stage research, not presenter behavior; AC-3 was removed.
- DONE: Correct the prevention claim to a guidance claim, keep the observable zero-bypass and zero-false-claim thresholds, and fix the wrong statement that presentation is too late to prevent advancement.
  Problem and ACs distinguish guidance from mechanical prevention, and generation from later advancement.
- DONE: Recompute the expected surface and the net LOC estimate for one file, keeping the separate insertions and deletions.
  In-memory application measured 8 insertions, 2 deletions, net +6 in one file; tolerance is net +4 to +18.
- DONE: Exercise the riskiest mechanism in this round with a throwaway trial of malformed seed generation, and record the command, version, observed output, and exit status.
  Self-applied conflicting-draft generation validated on pre3 with exit 0; FO-supplied fresh-worker trial 2 is primary evidence, with invocation and completion status retained.
- DONE: Replace the blanket "no spike needed" with that trial result and its remaining limits.
  Risk evidence distinguishes the weak self-applied and explicit-approval trials from the primary incident-shaped trial; all four FO-supplied limits are retained.
- DONE: Record the test entry path: the commission source names the Claude runtime, but a Pi commission run succeeded in practice on 0.28.0-pre4. Preflight that path and do not treat Pi commissioning as broken.
  Recorded captain-observed success, FO-observed Pi model fallback, and FO acceptance of this ideation preflight; full live handoff proof remains implementation work.
- DONE: Trim recurring instructions: remove the worked-example validation obligation from the shipped checklist, consolidate the redundant generation checks, and drop the extra README read where declarations suffice.
  One reconciliation checklist item replaces five additions; no presenter read or recurring example-validation requirement remains.
- DONE: State the residency boundary explicitly: this task does not fix presenter reload churn and does not depend on that fix.
  Out of scope states both boundaries and preserves existing lifecycle and fresh-evidence rules.
- DONE: Report in the canonical item form, one DONE/SKIPPED/FAILED line per checklist item with an evidence or rationale line, ending with a non-empty Summary.
  This cycle covers all ten dispatch items; frontmatter is unchanged and only this entity body is edited.

### Summary

Narrowed the design to commissioning guidance in one skill, with observable authoring and work-order criteria. Recorded a self-applied malformed-route trial and two supervisor-run fresh-worker trials, with the incident-shaped trial as primary spike evidence. Pi commissioning is reported successful on pre4; full commissioning-to-dispatch proof remains implementation work.


## Review-finding disposition — implementation

Source findings remain verbatim in `.briefings/commission-seed-approval-ideation-astra-review.md`; the references below identify each finding or recommendation without rewriting its historical evidence. The worker proposed these dispositions before editing. The supervising FO authorized them through `contact_supervisor`: guidance/recurring-check/spike corrections addressed by the approved cycle-2 design and exact edits; presenter expansion declined under the captain's one-file scope; Pi observations retained with live evidence outstanding; residency/readiness changes declined outside scope. No new product decision or round-record mutation was made by this worker.

| Finding / reference | Released user and normal workflow | Observable harm | Affected value AC or non-negotiable boundary | Trigger evidence | Materiality; ownership; authorized disposition |
| --- | --- | --- | --- | --- | --- |
| §1 authoring guidance / §6.1: presentation is not too late to prevent advancement | Captain commissions research then prototype | Conflating generation with advancement overstates the remedy | value-ac[AC-1] Preserve requested research in generated workflows | Review cites the original entity's incorrect presentation timing claim; cycle 2 distinguishes generation from advancement | Material; design-owned; addressed by approved cycle-2 problem statement, no presenter edit |
| §1 Confirm Design coverage: batch skips Question 2 | Captain supplies all commissioning inputs in batch | A Question-2-only correction would miss batch generation | value-ac[AC-1] Both modes must preserve research between approvals | Existing Batch Mode jumps to Confirm Design | Material; task-owned; fixed by exact Confirm Design reference to canonical Gate timing rule in c4a605df3 |
| §1 recurring checklist / worked YAML P2 and §6.4 | Captain commissions an ordinary workflow | Recurring worked-example and redundant checks impose unnecessary work without independent proof | none: recurring prose burden alone does not establish a value loss | Original design added five checks and a worked-example obligation without proposing worked YAML | Polish; design/task-owned; addressed by approved single reconciliation check in c4a605df3; smoke workflows validated separately |
| §2 prevention P2 and §6.1 | Captain assesses seed before research exists | False certainty about guidance can misrepresent whether bypasses are impossible | value-ac[AC-2] Seed descriptions must not claim completed research | Review compares prose promises to unchanged initial durability guard | Material; design-owned; addressed by approved guidance framing; no mechanical prevention claim or binary change |
| §1 presenter clarification and §6.4 two-file recommendation | Captain presents an already-existing malformed workflow | Contradictory Gate content may mislead review, but repair is outside newly commissioned output scope | captain-ruling[2026-10-07] Approved deliverable is commission-only, not presenter repair | Original broader AC-2 is superseded by approved generation-only AC-2 | Deferred risk for existing malformed workflows; outside ownership; declined per FO/captain scope; promote only under separately approved existing-workflow repair scope |
| §1 extra README presentation read | Captain presents a gate with declarations already available | No demonstrated loss from omitting a redundant read | none: existing declaration availability does not establish harm | Review says a fresh mandatory read is not justified | Polish; outside ownership; declined; no presenter or residency change |
| §1 readiness probes | Later-stage gate checks existing reports | No observed readiness regression caused by this documentation change | none: unchanged guard behavior is not instruction efficacy proof | Review identifies probes as not testing authoring and their equivalence as unverified | Deferred risk; existing status owner/outside task; declined; promote on a demonstrated guard regression |
| §3 Pi path P1 / §6.3 | Captain commissions on Pi and hands off to FO | Missing continuous runtime evidence leaves actual spawn order unproved | value-ac[AC-1] Actual research and prototype spawn order must be observed | Captain reports successful pre4 commission; FO reports out-of-space model ignored on Pi; worker has no root spawn capability | Material evidence gap, not an established Pi defect; validation-owned; path observations accepted by FO, live proof remains outstanding |
| §5 no-spike P1 / §6.2 | Captain approves instruction-based correction | Untested instructions may still generate the bypass or false seed claim | value-ac[AC-1] Requested research must survive commissioning | Approved entity retains supervisor fresh-worker runs d850a1e4 and 2ba124cf plus limits | Material; design/task-owned; spike correction addressed in cycle 2; supervisor blind trial 2 remains primary spike, self-applied implementation smokes are weak supplements, not AC-1 |
| §4 residency / §6.5 sequencing | Captain reviews gates under current lifecycle | No present conflict or dependency established; future caching must retain fresh evidence | captain-ruling[2026-10-07] This task neither fixes presenter residency nor depends on it | Review explicitly excludes presenter loads from prior residency coverage | Deferred risk; outside ownership; declined; promote only if later residency scope changes freshness/lifecycle promises |

## Implementation verification evidence

Deliverable commit: `c4a605df3` on `spacedock-ensign/commission-seed-approval-vs-result-review`, pushed to that branch only. Base: `077fb3bc1a353b7f8384a3b37c1309c199d82ba1`. A Python comparison extracted the five approved before/after blocks from this entity and applied them to `git show 077fb3bc1:skills/commission/SKILL.md`; resulting bytes exactly equaled the deliverable. This is an artifact-equivalence check, not behavioral proof. `git diff --numstat 077fb3bc1 HEAD` reports `8  2  skills/commission/SKILL.md`: net +6, exactly the estimate and within +4 to +18. `git diff --check` passed. No permanent tests were added under the one-file constraint.

FO-authorized temporary artifacts: `/tmp/commission-seed-implementation.btjAj8`. Shell/Python serialized this worker's self-applied outputs before and after editing, not another model's output. `before-skill.md`, `after-skill.md`, `request.txt`, mode-input files, confirmation files, four workflow directories, validation logs, version metadata, and limitations are retained there, not committed to the deliverable. Baseline skill SHA256 `fd9d10719cbafd493e869659393a9c9bcb751b20acd9ca9547fe81411a921aa8`; changed skill SHA256 `19e6612ad727785a757525c19a46bd30d366b5cb1fbfaccca701b49b14274d1f`.

Primary request was the Test plan's fixed conflicting-draft request, unchanged, for both simulated modes. Simulated interactive answers specify one idea, corrected route, one seed; batch supplies those inputs and its disposable location directly. Both use slug identity, inline state, no mods, report-only work, and explicit approval boundaries. No live captain conversation or approvals occurred. An additional self-applied response to the incident-shaped request (including bb.app) is retained in `incident-self-applied-response.txt`; it is not a generated workflow or a schema-validity claim.

All four generated workflow declarations are `backlog (initial, gate) -> exploration (gate) -> prototype -> done (terminal)` with `worktree: false`. Exploration Outputs require a completed, committed report containing a reference scout, cited sources, findings, limitations, and an experiment recommendation. No research report was actually produced. Baseline and changed outputs both explicitly correct the conflicting draft. Baseline success is not a measured before/after improvement. Changed outputs reused the baseline schema scaffold; this is especially weak instruction-efficacy evidence from an author already aware of the expected design.

Changed confirmation: "At backlog, you review the committed seed, research scope, and open questions. Approval starts exploration including a reference scout. No backlog worker has run or produced findings." Later confirmation reviews the completed, committed exploration report, reference-scout observations, sources, and limitations, then starts prototype. Changed seed Gate content: "Review the committed seed, research scope, and open questions. No initial-stage worker has run and no research findings exist yet. Approval starts exploration, including a reference scout, not prototype work." Exploration Gate content reviews the completed, committed report and required research artifacts before prototype. Manual inspection of these outputs found zero proposed research bypasses and zero seed claims of completed research. Those observations do not satisfy AC-1 or establish live AC-2 efficacy.

Validation executable: `/opt/homebrew/Caskroom/spacedock@next/0.28.0-pre3/spacedock`. `--version` reported `spacedock 0.28.0-pre3`, `darwin/arm64`, runtime `pi (PI_CODING_AGENT)`, contract 3; exit 0. `pi --version` reported `1.0.2`; `go version` reported `go1.26.7 darwin/arm64`. The worktree plugin manifest is `0.28.0-pre4`; validation used pre3, not pre4.

| Exact validation command (same executable above) | Output | Exit | Local throwaway seed commit |
| --- | --- | --- | --- |
| `status --workflow-dir /tmp/commission-seed-implementation.btjAj8/before-interactive --validate` | `VALID` | 0 | `d47494f` |
| `status --workflow-dir /tmp/commission-seed-implementation.btjAj8/before-batch --validate` | `VALID` | 0 | `d47494f` |
| `status --workflow-dir /tmp/commission-seed-implementation.btjAj8/after-interactive --validate` | `VALID` | 0 | `936cedf` |
| `status --workflow-dir /tmp/commission-seed-implementation.btjAj8/after-batch --validate` | `VALID` | 0 | `5593590` |

The local smoke commits retain seeds and README files only; they are not runtime transition commits. No worker spawn, dispatch assignment, research-result approval, prototype execution, or continuous commission-to-FO handoff was observed. AC-1 remains outstanding and requires a live interactive root session for each mode. FO explicitly assigned that proof to validation and will surface the gap to the captain. The supervisor's blind trials recorded above remain the primary spike evidence; they were not rerun or claimed as this worker's work.

Repository checks: before editing, `go test ./internal/cli -run TestCommissionOrphanBranchScaffolding -count=1` passed (exit 0, 1.579s). It checks orphan-state creation, no inherited code files, clean code checkout, and status rendering; leaking inherited README files would fail it. It does not test the new instruction behavior. `gofmt -w ./cmd ./internal` exposed two pre-existing field-alignment changes in `internal/release/runtime_live_evidence_workflow_test.go`; those incidental edits were immediately restored, leaving no unrelated formatting changes. The initial foreground `go test ./...` tool call timed out after 120 seconds; the normal and race suites were then relaunched with polled logs rather than treating that interruption as a pass.


### FAILED — full repository test suites

Both `go test ./...` and `go test ./... -race` completed with exit 1. A focused pass does not cancel either full-suite failure. Logs are under `/tmp/commission-seed-implementation.btjAj8/`.

| Command / log | Failed package and duration | Observation at package timeout |
| --- | --- | --- |
| `go test ./...` / `go-test.log` | `internal/cli`, 600.482s | Go's 10m package deadline; `TestStateCommitMakesFlatRoundDurableInFreshHost` active for 2s |
| Same | `internal/ensigncycle`, 602.030s | Go's 10m package deadline; `TestDurableKeepMovingDelayedPersistence/corroborated_frontier` active for 2s |
| `go test ./... -race` / `go-test-race.log` | `internal/cli`, 600.821s | Go's 10m package deadline; `TestMergeGuardSiblingDirtNonFastForwardRemainsRecoverable` active for 1s |
| Same | `internal/ensigncycle`, 601.924s | Go's 10m package deadline; `TestDurableKeepMovingRequiresOverlappingJourneys` active for 51s |

**FAILED race assertion:** `TestCodexProcessRequiresFinalMessageForTerminalTurn` (1.37s) expected missing-terminal-evidence failure but observed a quiet timeout: `duration:619291125 exitCode:-1 terminal:false timedOut:true`, last event `item.started` / `collab_tool_call` / `wait_agent`. This is a distinct assertion failure, not just the later package timeout. Both suites passed `skills/integration`; no full-suite success is claimed.

Focused reruns, all exit 0 (these supplement rather than replace the failed broad checks):

- `go test ./internal/cli ./internal/ensigncycle -run '^(TestStateCommitMakesFlatRoundDurableInFreshHost|TestDurableKeepMovingDelayedPersistence)$' -count=1`: package durations 23.521s and 72.931s; `timeout-focused.log`.
- `go test ./internal/cli ./internal/ensigncycle -race -run '^(TestMergeGuardSiblingDirtNonFastForwardRemainsRecoverable|TestDurableKeepMovingRequiresOverlappingJourneys|TestCodexProcessRequiresFinalMessageForTerminalTurn)$' -count=1`: package durations 8.700s and 18.776s; `race-focused.log`. In particular the Codex assertion passed under race in this focused run.
- `go test -v ./internal/cli ./internal/ensigncycle -run '^(TestStateCommitMakesFlatRoundDurableInFreshHost|TestDurableKeepMovingDelayedPersistence|TestMergeGuardSiblingDirtNonFastForwardRemainsRecoverable|TestDurableKeepMovingRequiresOverlappingJourneys)$' -count=1`: all passed; `timeout-focused-verbose.log`. Per-test normal rerun durations: state durability 3.28s; delayed persistence 27.27s (corroborated frontier 4.98s); sibling-dirt recovery 6.90s; overlapping journeys 16.97s.

**Additional review-finding disposition — full-suite failures.** Released user and normal workflow: contributors running the required Go checks. Observable harm: no clean full-suite signal. Affected boundary: `contract[AGENTS.md#expected-commands]` required checks must be run and reported honestly. Trigger evidence: the two retained failure logs and three focused-pass logs above. Materiality authorized by FO: "unproven cause affecting every task in this workflow"; ownership: outside this one-skill deliverable; disposition: decline product fixes, validation reruns the broad checks independently and owns the verdict. FO authorized this after a distinct `need_decision` consultation and is raising it with the captain. No Go code changes were retained.

FO attribution, recorded as its conclusion rather than a baseline reproduction: "the candidate changes one markdown file, so it cannot cause a Go package timeout or a Codex assertion failure." Whether these failures reproduce on the untouched base is **UNVERIFIED**. This worker does not establish a root cause or label the failures pre-existing or load-sensitive.

## Stage Report: implementation

- DONE: Apply the five before/after blocks to skills/commission/SKILL.md exactly as the entity's Proposed approach specifies, and change no other file.
  Deliverable commit `c4a605df3` pushed to the assigned branch; exact-block byte comparison passed; only the separately authorized state report and temporary smoke files are outside that deliverable.
- DONE: Run each required check (gofmt -w ./cmd ./internal, go test ./..., and go test ./... -race), report each result, and retain no unrelated formatting changes.
  Formatting ran; incidental baseline field alignment was restored. Both full suites exited 1 with 10m cli/ensigncycle package timeouts; race also failed `TestCodexProcessRequiresFinalMessageForTerminalTurn`; exact durations, focused passes, and logs are recorded above.
- DONE: Measure the actual diff with git diff --numstat against the approved base and report insertions, deletions, and net against the +6 estimate and the +4 to +18 tolerance.
  Against `077fb3bc1`: one deliverable file, 8 insertions, 2 deletions, net +6 (0 deviation; within tolerance); `git diff --check` passed.
- DONE: Validate any generated workflow with status --validate and record the command, version, output, and exit status.
  All four temporary before/after simulated-mode workflows returned `VALID`, exit 0 on pre3; exact commands, runtime/version, and local seed commits appear above. This proves schema validity only.
- SKIPPED: Run the acceptance checks the design names with the fixed request, in interactive and batch modes. If a live interactive commission-to-FO drive cannot run from your context, say so plainly and record exactly what you ran instead. Do not claim a live worker spawn you did not observe.
  No live root conversation/FO spawn capability here. FO authorized weak self-applied pre/post simulated interactive/batch smokes instead; fixed request and outputs retained. AC-1 actual spawns and live AC-2 outputs remain validation-owned, not passed.
- DONE: Disposition every review finding on this task under the workflow's Review-finding disposition section, with the four evidence fields.
  Ten historical finding/recommendation rows plus the new full-suite verification finding have evidence, separate ownership/materiality/disposition, and distinct FO authorization; no presenter or Go fix was made.
- DONE: Report in the canonical item form, one DONE/SKIPPED/FAILED line per checklist item with an evidence or rationale line, ending with a non-empty Summary.
  This report covers all seven checklist items, preserves entity frontmatter, and distinguishes failed broad checks from focused passes and weak generation evidence.

### Summary

Implemented and pushed the exact one-file seed/result approval guidance, net +6 lines; four self-applied workflows schema-validated, but live commission-to-FO acceptance remains unproved. Both full Go suites failed and their focused reruns passed; the cause is unproven, and the FO classifies the failures as not attributable to this candidate (untouched-base reproduction remains UNVERIFIED). At the FO's direction, the check item is DONE because its obligation was to run and report each check, not obtain a pass; all failure evidence remains above, and validation owns independent broad-check reruns, actual interactive/batch worker observations, and the acceptance verdict.
