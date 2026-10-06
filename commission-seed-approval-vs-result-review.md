---
title: The commission skill lets an initial seed gate claim to review work no worker produced
status: ideation
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
started: 2026-10-06T23:39:26Z
---

## Problem

A commissioned initial gate can promise review of research that no worker has produced. Approval then advances past the intended research work.

The independent incident baseline is `.briefings/commission-seed-backlog.md` in this state checkout, committed at `fe360ace7`. On `0.28.0-pre4`, the proposed route was `exploration (initial, gate) -> prototype`. Its gate promised approval of a researched direction. The captain stopped it before dispatch. The corrected route was `backlog (initial, gate) -> exploration (gate) -> prototype`.

The baseline contains one intended research stage bypassed by the proposed seed approval, not an observed completed bad dispatch. The finished change must produce zero bypasses in the live checks below. It must also prevent a seed presentation from claiming completed research.

`internal/status/entered_stage.go` confirms the distinction. `gatePreparable` accepts a tracked, clean seed when `stage.initial` is true. Other gated stages require a complete, committed current-stage report. Missing reports never turn later stages into seed gates.

The combination `initial: true` plus `gate: true` is legitimate. `docs/dev/README.md` uses it for backlog. Its Gate content asks whether design can start, not whether design has finished.

## Proposed approach

Change instruction text only, in two skills. Keep seed approval and result approval as separate decisions.

1. Teach gate timing during stage design and Confirm Design, including batch mode. Name the reviewed artifact and actual successor work.
2. Check those names during generation. Preserve legitimate initial gates; separate requested research from the initial seed when both approvals are needed.
3. Render initial gates through the existing presenter. State that the committed seed is under review, with no worker result from that initial stage.
4. Preserve the existing result-gate completion requirement and decision authority. Do not invent a new validator or runtime branch.

### Concrete instruction changes

These are the proposed documentation changes for implementation, not edits made during ideation. Keep existing unrelated text. The added wording also defines the captain-facing behavior; no separate site or CLI change is proposed.

**`skills/commission/SKILL.md`, Question 2 — Stages**

Before:

> Store the confirmed stages as `{stages}`. The first stage is `{first_stage}` and the last is `{last_stage}`.

After:

> Store the confirmed stages as `{stages}`. The first stage is `{first_stage}` and the last is `{last_stage}`.
>
> An initial approval gate reviews the committed seed, not work assigned to that stage. No worker runs for that initial stage before seed approval. Approval advances to its successor. A later approval gate reviews its stage's completed, committed worker report and artifacts.
>
> Separate permission to start research from permission to build from research results. When both decisions are required, use a seed stage, a research stage, then the next work stage. Put the seed approval on the initial stage and the result approval on the research stage. Name the reviewed artifact and successor work for each approval. Do not describe seed approval as review of research that has not run.

**`skills/commission/SKILL.md`, Confirm Design**

Before:

> - `{approval_gates}` — default: gate before the terminal stage (e.g., the last stage before terminal).

After:

> - `{approval_gates}` — default: approval at the last stage before terminal. For each approval, name its reviewed artifact and successor work. Apply the gate-timing rule from Question 2, including in batch mode. If seed approval precedes research, keep research as a separate successor stage. Approval of the seed must not skip that work.

In the captain-facing summary, replace this line:

> `{for each stage: "{letter}. {stage_name} — {stage_description}"}`

With:

> `{for each stage: "{letter}. {stage_name} — {stage_description}"}`
> `{for each approval: "At {stage_name}, you review {artifact}. Approval starts {successor_work}."}`

**`skills/commission/SKILL.md`, generated README guidance**

Before:

> For every gated stage, add `- **Gate content:**` to its stage subsection and state the evidence needed for that decision. This rule also applies to custom stages and template variants.

After:

> For every gated stage, add `- **Gate content:**` to its stage subsection. Name the reviewed artifact, required evidence, and actual successor work. Apply this rule to custom stages and template variants. For an initial gate, name the committed seed and permission to start the successor work. Do not promise review of that initial stage's worker results. For a later gate, name its completed, committed stage report and required artifacts.

**`skills/commission/SKILL.md`, Generation Checklist**

Before:

> - [ ] Each seed entity file exists at `{dir}/{slug}.md` with valid YAML frontmatter

After:

> - [ ] Each seed entity file exists at `{dir}/{slug}.md` with valid YAML frontmatter
> - [ ] Each approved gate is declared on its stage, including the initial stage when requested.
> - [ ] Each gate names its reviewed artifact and actual successor work in the generated README.
> - [ ] Initial gates review committed seeds. Seed approval advances to the requested first work stage without skipping it.
> - [ ] Later gates require their completed, committed stage reports and artifacts before result review.
> - [ ] Each worked YAML example passes `spacedock status --workflow-dir {example_dir} --validate` in a disposable workflow.

**`skills/present-gate/SKILL.md`, Captain-facing assembly rules**

Before:

> - Fetch the current stage with `${SPACEDOCK_BIN:-spacedock} dispatch show-stage-def --workflow-dir {workflow_dir} --stage {stage}` before selecting evidence. Never infer content from the stage name.

After:

> - Fetch the current stage with `${SPACEDOCK_BIN:-spacedock} dispatch show-stage-def --workflow-dir {workflow_dir} --stage {stage}` before selecting evidence. Never infer content from the stage name.
> - Read the workflow README's stage declarations to identify an initial gate. Do not infer it from a missing report. For an initial gate, identify the committed seed as the reviewed artifact. For a newly commissioned seed, state: "No worker has run for this initial stage. Approval starts {successor_work}; it does not approve completed research." Name the actual successor from the workflow. Do not present absent worker results as evidence.

Before:

> - Treat a `Gate content` instruction as the workflow's authoritative presentation preference and override. When it is absent, show a concise, meaningful, decision-relevant subset of the current stage definition, selected Artifact and References, current stage report, checklist and AC evidence, and findings. Do not fabricate facts or dump every source.

After:

> - Use `Gate content` to select decision-relevant evidence. It cannot override gate timing or turn missing results into evidence. If an initial gate promises unproduced worker results, explain the mismatch and recommend revision of the workflow design. Otherwise, follow its presentation preference. Without `Gate content`, select relevant stage definition, Artifact, References, stage report, checklist, AC evidence, and findings. Do not fabricate facts or dump every source.

### Pi loading and authority

Keep the Pi adapter unchanged. Its `«gate.lifecycle»` requires prepare, commit, reread, presenter load, then one root-session review block. The initial-gate case runs **inside** `spacedock:present-gate`; it is not a lightweight bypass or a load exemption. The shared core also delegates gate rendering to that skill. Both rules remain satisfied.

Keep the entity, stage, bound Briefing identity/digest, recommendation, and concrete decision ask in the initial review. Omit nonexistent result rows, not these required facts. A contradiction in Gate content is explained before any decision; the FO does not silently rewrite workflow state. The captain retains design authority. The presenter adds no recording authority and does not change `gate record`.

## Expected surface and tolerance

Estimate net LOC change: +40, across 2 files.

| Deliverable file | Insertions | Deletions | Net |
| --- | ---: | ---: | ---: |
| `skills/commission/SKILL.md` | 35 | 5 | +30 |
| `skills/present-gate/SKILL.md` | 13 | 3 | +10 |
| Total | 48 | 8 | +40 |

Tolerance: net +25 to +55 LOC; exactly these two deliverable files. Measure against the implementation branch's approved base, using `git diff --numstat`. State-body reports are separate from this deliverable estimate. This ideation round changes only this entity body.

### Declared semantic changes

- **Instruction behavior:** commission distinguishes seed permission from result approval, in interactive and batch paths.
- **Presentation behavior:** initial reviews identify the seed and successor work; contradictory result claims cause an explicit revision recommendation.
- **Command grammar and stored formats:** unchanged. No schema, CLI, report format, or gate-record changes.
- **Authority:** unchanged. The FO presents; the captain or existing explicit delegation decides. Workflow presentation preferences cannot authorize false evidence.
- **Runtime behavior:** binary transitions, report guards, Pi load timing, and dispatch transport stay unchanged. Only agent authoring and presentation behavior changes.

## Out of scope

No ban on initial gates. No binary edits, new lint, CI lane, validator, template migration, host adapter edit, or gate-loading optimization. No generic redesign of commission's runtime startup. No worktree, PR, main push, or live workflow mutation in this ideation round.

No worked YAML example is included. The route above is explanatory prose, not a new schema example. If implementation adds runnable YAML, validate each complete disposable workflow and record the command, version, output, and exit status.

## Acceptance criteria

**AC-1 (VALUE) — Seed approval preserves the requested research work.**
The commissioned workflow has zero skipped research stages in both interactive and batch checks, against the incident's one proposed bypass. The independent oracle is the captain's requested order: seed approval, research and reference scout, research-result approval, then prototype. Expected order is fixed before generation, not derived from the generated README.
Verified by: a targeted live commission-to-FO drive in each mode. The first ready gate belongs to the seed. After approval, an observed worker spawn receives the research assignment, including the requested reference scout. No prototype worker starts before completed research and its approval. After research-result approval, an observed prototype spawn occurs. Retain generated README, assignments, root transcript, state, reports, and commit IDs. A generated route that marks research initial and advances directly to prototype fails, even if YAML validation passes.

**AC-2 (VALUE) — Gate presentations distinguish seed permission from result approval.**
The first review names the committed seed, says no worker ran for the initial stage, and names research as the next work. It makes zero claims of completed research. A later research-result review is unavailable without the complete, committed research report. With that report, it reviews research evidence and names prototype as the successor.
Verified by: the same live drive's root-session review blocks and durable gate state. Before research completion, attempt preparation on an isolated copy; it must refuse and create no open result gate. Repeat with an incomplete report and a complete but uncommitted report. Then observe successful preparation after a complete committed report. Each result review follows that success. A Gate content override claiming researched results at the seed gate must produce an explicit mismatch and revision recommendation, not invented evidence.

**AC-3 — Legitimate initial gates retain the standard Pi lifecycle and decision authority.**
The corrected seed gate remains valid and uses the existing presenter after commit and reread. Each engaged gate has exactly one qualifying root-session presentation before its decision mutation. Required identity, snapshot, recommendation, and ask remain present. There is no initial-gate load exemption or duplicate fallback presentation.
Verified by: ordered Pi root events and state commits from the AC-1/AC-2 drive, plus `status --validate` on each generated workflow. The existing `docs/dev` backlog declaration also remains valid. Omitting the presenter load only for the initial gate, or presenting again after recording, fails. This criterion supports both value criteria; schema validity alone does not establish either.

## Test plan

**Primary proof owners.** The status-owned completion guard in `internal/status/entered_stage.go` owns mechanical readiness. Existing Go tests remain its regression protection; no new binary mechanism or test suite is proposed. No existing live test is identified in the permitted search surface as owning commission's authored stage order. A one-off targeted live drive is therefore the primary proof of the changed skills, not prose matching.

The implementation owner first records the fixed input and observation sheet, then runs the smoke before changing skill text. Repeat against changed skills; preserve both outcomes. Do not require the old skill to fail on demand. The historical incident is the independent bad-order baseline; the fixed requested order is the acceptance oracle. A lucky pre-change live pass does not replace post-change evidence.

Use this fixed request in both modes: "Create a workflow to assess an idea. I must approve its seed before research starts. Research must include a reference scout. I must approve the research results before any prototype work starts." For the interactive run, answer the ordinary design questions consistently with that request. For batch, provide all inputs and a disposable inline workflow location; do not auto-approve result gates. The test operator supplies explicit approvals at the two observed boundaries. Use Pi for the root-session load and presentation evidence. Do not claim actual spawn from a built assignment alone.

| Check and owner | Distinct falsifying edit | Cost and mode |
| --- | --- | --- |
| AC-1: implementation smoke, independently replayed by validation | Make the commissioned research stage initial and route approval straight to prototype | Two bounded live runs, one interactive and one batch; approximately 10–20 minutes each |
| AC-1 batch coverage within those runs | Put timing guidance only under Question 2 and omit its application in Confirm Design | Same batch run; no additional harness |
| AC-2: validation inspects actual root review blocks | Let initial Gate content assert research findings despite absent worker evidence | One contradictory-seed presentation, approximately 5 minutes, live |
| AC-2: status readiness negative cases | Allow non-initial gates to use clean seed durability instead of a complete report | Three isolated readiness probes, approximately 2 minutes, deterministic CLI; no production mutation |
| AC-3: validation inspects the same root event trace | Skip the present-gate load at initial gates while retaining later loads | Same Pi runs; no separate live loop |
| AC-3: generated workflows pass schema validation | Emit an invalid stage name or omit the terminal declaration | `status --validate` per generated workflow; seconds, deterministic CLI |

Retain each probe's command, exit status, resulting state, and commit identity in the validation evidence. A readiness probe proves the guard, not whether an agent invents results; root-session evidence covers that separate claim. A schema pass cannot prove stage order or truthful presentation.

Run `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal` during implementation verification as required by the repository. Verify that formatting produces no unrelated changes. Apply the workflow's existing live-lane and independent-review policy to the final two-file diff; do not introduce a new lane. This ideation does not modify code or run those broad checks.

### Mechanism necessity

- Design-time timing guidance serves AC-1. Presenter-only wording is cheaper but detects a skipped research stage after generation, too late to prevent it.
- The generation checklist serves AC-1. Design prose alone does not check whether generated declarations preserve the confirmed approvals, especially in batch mode.
- Initial-gate assembly serves AC-2. Commission-only text cannot correct a misleading Gate content instruction in an existing workflow.
- One-off live drives serve AC-1 through AC-3. Reading the new wording or validating YAML cannot observe the authored route, actual spawn, or root presentation.
- Existing CLI readiness probes serve AC-2. A positive completed-report run alone cannot expose the missing-report bypass. A new standing guard duplicates the shipped one.

## Risk evidence and spike disposition

No spike needed: the design adds no parser, stored format, runtime handoff, or tool flag. It relies on existing initial-stage readiness, committed report checks, ordinary successor dispatch, and the existing presenter load.

Auditable evidence:

- Code baseline `c450d5d526fca30bbcd2c06cf2b40358e5a6ea8f`: `gatePreparable` branches only on `stage.initial`. `stageReportFailure` rejects missing or incomplete reports; `entityGitFailure` requires tracked, clean bytes.
- The committed incident briefing `fe360ace7:.briefings/commission-seed-backlog.md` records the bad proposal and corrected route. It does not claim a worker was actually skipped.
- This entity already followed the legitimate initial-gate route: state commit `b31f7a5e4` records backlog approval, `7cc73f250` consumes it into ideation, and `5dde37bd2` records dispatch. Its current read reports `status: ideation`; this worker is executing that assignment.
- The supplied launcher reports `spacedock 0.28.0-pre3`, contract 3. `status --workflow-dir docs/dev --validate` returned `VALID`, exit 0. It warned about unrelated `live-lane-gap-inventory` having verdict `withdrawn`; that entity was not changed.
- The Pi adapter and shared core already require the presenter. The proposed initial case preserves those boundaries rather than adding an exception.

Instruction efficacy remains unproven until the post-change live drive. This is the principal implementation risk, not a claimed ideation success. Version drift also remains visible: the incident names pre4; this round's available validation binary is pre3. Record the implementation binary and skill revision together. Stop for a scope decision if current runtime behavior differs from these retained semantics; do not repair the binary under this documentation task.

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
