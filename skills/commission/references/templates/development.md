---
commissioned-by: spacedock@template
entity-type: task
entity-label: task
entity-label-plural: tasks
id-style: sd-b32
stages:
  defaults:
    worktree: false
    concurrency: 2
  states:
    - name: backlog
      initial: true
      gate: true
    - name: ideation
      gate: true
    - name: implementation
      worktree: true
      context-sections:
        - Review-finding disposition
    - name: validation
      worktree: true
      fresh: true
      feedback-to: implementation
      gate: true
      context-sections:
        - Review-finding disposition
    - name: done
      terminal: true
---

# Development Workflow Template

Tasks move from a captain-curated backlog through ideation, get built in a dedicated worktree, are independently validated against the acceptance criteria, and land via PR review. This is the refinement shape specialized for code that ships via PR/merge: the repo-mutation layer is active on `implementation` and `validation`, and the `pr-merge` mod handles the PR lifecycle.

Use this template when the captain's mission is to ship code in a repo where work is reviewed and merged via PR. The bucket-noun stage names (`implementation`, `validation`, `done`) describe where the entity is sitting; the `pr-merge` mod removes any temptation to add `pr_open` or `awaiting_merge` stages because PR state is tracked on the `pr` field, not as a stage.

## File Naming

Each task lives as either:

- a flat markdown file `{slug}.md` (default), or
- a folder `{slug}/` containing `index.md` when the task produces sibling artifacts (transcripts, design notes, comparison tables) that belong with the tracker.

Slugs are lowercase, hyphens, no spaces. Example: `pluggable-id-style-collaboration-friendly-ids.md`.

## Schema

Every task file has YAML frontmatter. Fields are documented below; see **Task Template** for a copy-paste starter.

### Field Reference

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique identifier; SD-B32 by default for collaborative cross-branch creation |
| `title` | string | Human-readable task name |
| `status` | enum | One of: backlog, ideation, implementation, validation, done |
| `source` | string | Where this task came from (issue, captain note, retrospective) |
| `started` | ISO 8601 | When active work began |
| `completed` | ISO 8601 | When the task reached terminal status |
| `verdict` | enum | PASSED or REJECTED — set at validation |
| `score` | number | Priority score, 0.0–1.0 (optional) |
| `worktree` | string | Worktree path while a dispatched agent is active, empty otherwise. Once set on first dispatch into a `worktree: true` stage, it stays set across all non-terminal advancements (stickiness) and clears at terminal merge. |
| `issue` | string | GitHub issue reference (e.g., `#42`) — optional cross-reference |
| `pr` | string | GitHub PR reference (e.g., `#57`) — set when a PR is opened for this task's branch |
| `mod-block` | string | Pending mod-declared blocking action, format `{lifecycle_point}:{mod_name}` |

## Stages

### `backlog`

A task enters backlog when it is first proposed: a seed description, no design work. Captain-curated holding stage — the gate decides which tasks advance to ideation.

### `ideation`

The captain greenlights a task for design: flesh out the problem, propose an approach, define acceptance criteria as entity-level end-state properties with `Verified by:` clauses, and write a test plan that matches the AC's level of abstraction.

- Split each acceptance criterion by how it is verified: **offline** (a test, command, or on-disk state a fresh agent reproduces) or **interactive** (requires a human or a live drive to judge). Declare the split at ideation. A plan that would build a harness to automate an interactive AC is visible here, at the gate, before the harness is built — interactive ACs are validated by a live drive or the captain, not by new automation.

### `implementation`

The design is approved and the deliverable is built in a dedicated worktree on a feature branch — minimal changes that satisfy the AC, self-contained for validation.

- When a finding arrives, follow `## Review-finding disposition`: investigate read-only, preserve its evidence, propose materiality/ownership/disposition, and obtain distinct FO authorization before any candidate edit, commit, or reviewer rerun.

### `validation`

A `fresh` agent independently verifies the deliverable against the ideation AC, reproducing each `Verified by:` clause rather than trusting the implementation's self-report. The validator checks what was produced; it does not produce it. Either gate-approval to `done` or rejection back to `implementation` with concrete fixes.

- **Small-change fast path.** Scale the validation checks to the diff's blast radius. A routine, low-blast-radius change (a doc line, a one-line fix, a rename) does not need the full checklist or the detached adversarial audit — the same "routine changes exempt" carve-out the audit already grants, applied to validation as a whole. Match the rigor to the change; a trivial diff over-validated is its own waste.

### `done`

Terminal state: the task's PR is merged (tracked via the `pr` field and the `pr-merge` mod), `completed` set, `verdict: PASSED`, entity archived. Reached via real merge, not a manual flag flip.

## Review-finding disposition

Every finding enters this checkpoint when it arrives during implementation, validation, a detached audit, consequential FO quick work, or a correction routed from a rejected gate.

1. The reviewer owns observation, not task ownership or authorization.
2. The worker preserves the finding, investigates without candidate mutation, records the four evidence fields, and proposes materiality, task ownership, and disposition separately. Its `actor:ensign` round Resolution is advisory.
3. The FO sends a distinct `fix`, `decline`, `hold`, or `route for decision` authorization through the runtime's addressable-worker boundary.
4. The validator recommends `PASSED` or `REJECTED`; a new finding re-enters step 1.
5. Only the captain changes approved scope, accepted value, thresholds, tolerance, or acceptance criteria.
6. After revise is selected, rejection routing transports the evidence, workflow classifications, authorized dispositions, and concrete assignment unchanged; it never re-triages.

Before FO authorization, candidate bytes and Git HEAD stay unchanged, no candidate commit is made, and no reviewer rerun starts. Read-only file/history inspection, non-mutating reproductions, existing tests, and adversarial work in a throwaway checkout are allowed. After authorization, perform only that disposition; `hold` and `route for decision` forbid mutation and rerun. Changed evidence re-enters the checkpoint, and an unobservable runtime authorization means hold and re-consult.

The four evidence fields are released user and normal workflow; observable harm; affected value AC or non-negotiable boundary; and trigger evidence. Field 3 uses `value-ac[AC-N]`, `captain-ruling[YYYY-MM-DD]`, or `contract[repo/relative/path#anchor]` plus a nonblank claim; `none:` plus a rationale cannot establish Material.

- **Material:** all four fields establish supported-workflow harm to a value AC or protected boundary.
- **Deferred risk:** the trigger is hypothetical, unsupported, unobserved, or outside current promises; record its promote-to-material condition.
- **Polish:** no current user-visible loss or protected boundary is at risk.
- **Needs decision:** the task cannot own the required scope, product, or compatibility decision.

Materiality and task ownership are independent. Owned Material is eligible for an FO-authorized fix; out-of-scope Material holds unchanged as Needs decision. Deferred risk or Polish may be declined only after FO authorization.

After reviewer and worker entries and FO consultation, retain the advisory round with `${SPACEDOCK_BIN:-spacedock} gate record`; it applies no gate or status change. A workflow correction round uses `- Cycle {N}: {verdict} — {reviewer/loop}; surface {files}/{LOC} vs estimate {declared} ({P}%); AC {unchanged | narrowed: <note>}`. Compare `git diff --numstat "$(git merge-base main HEAD)"..HEAD` with the ideation estimate; beyond declared tolerance or on narrowed AC, require a captain-visible design reset. Cycle 3 escalates.

## Workflow-specific rules

The FO/ensign operating contract already governs generic stage semantics and proof discipline: prefer the cheapest check that can fail — a shipped guard's run, an existing mechanical check, a one-off falsifiable exercise recorded in the report, then the captain's judgment — with new standing enforcement as the last resort rather than the default; prove by exercising rather than re-reading; and reject any AC whose only proof is a review of its own prose. Tasks in a commissioned development workflow inherit these rules from the contract their FO loads at boot; the rules below add only the dev-shape specifics.

- **Repo-mutation worktree layer.** `implementation` and `validation` run in a worktree against the codebase, and `validation` is `fresh` so an independent agent checks the AC. PR state lives on the `pr` field, managed by the `pr-merge` mod — there is no `pr_open` or `awaiting_merge` stage.
- **No prose-grep over instruction files.** A string, substring, or regex match over an instruction file the model reads (the FO/ensign contract, this README, a skill) never proves a behavioral claim. The matched text was written by the same implementer the check polices, so it asserts only that the file contains what we put in it. A valid paraphrase fails it and an inverted clause passes it. To settle a case, ask whether the expected value comes from outside the file under test; if it does not, the check is a tautology and is banned. A check that binds two independent values that can diverge, such as the plugin manifest's version sharing a major.minor with the binary's version, is legitimate and is not prose-grep. Captain ruling (2026-07-20, verbatim): prose-greps are one-off validation evidence, never committed tests. A grep whose output is pasted into the validation report is legitimate external evidence for that run; the same grep committed as a test is banned — it re-asserts that the file contains what we wrote and cannot fail. What such a grep may be evidence FOR is bounded by honesty, not by category: presence or absence is an existence fact, and a grep establishes it soundly when that fact is itself the claim. When it is not — when the claim is about what a program or an agent actually does, and the words being present says nothing about that — the grep is misleading, and the answer is to express the claim in a form that can be exercised, not to let a grep that technically passes stand in for proof it cannot give.
- **Evidence must be able to fail.** Each AC's cited evidence names the concrete change that would flip it — the falsifying edit. An author who cannot name what would make the evidence fail has not shown it can fail, and the criterion does not count.
- **Opt-in proof disciplines (copy into the `validation` stage when commissioning).** These are recommended dev-shape practices, not universal — a non-development workflow's acceptance proof may legitimately be a published artifact, a metric, or a human review. Adopt the ones the mission needs by folding them into the `validation` stage's Outputs and Bad lists:
  - **Test-first authoring** — for a code or fixture deliverable, write the failing test first, watch it fail for the right reason, then write the minimum code to pass. The test is what the gate judges.
  - **External-proof acceptance criteria** — each AC's evidence must come from a check outside the task body (a test, a command's output or exit code, a file the change produces, on-disk state). Reject self-referential ACs whose only proof is review of the task's own prose; if nothing ships, the decision belongs in the roadmap, not a terminal dev task.
  - **Detached adversarial audit** — for high-stakes surfaces (a front-door launcher, status/guard mutation paths, shipped contract/scaffolding, CI/release machinery), run a read-only audit on a throwaway checkout that tries to refute the validation with an edit the deliverable's own tests should catch. `Material:` findings route back through the validation→implementation feedback flow; "refuted nothing material" is a valid recorded outcome. The audit also fires on AC provenance: when an AC's expected value is derived from the same package's production functions or constants, run the adversarial-edit check on it — that provenance is the tautology tell. Scope it to that provenance form; the broader equality/byte-identity form over-fires on ordinary unit tests.
  - **Live scenario for runtime claims** — when an AC's truth is what an agent or model *does* at runtime, prove it with a scripted live scenario graded on durable before→after state plus observed output, with a negative case that reds the grade. Mark the AC `Verified by: live <ci-run:<id> | session:<path>>`; an offline proxy or a contract-text check proves the watcher or the words, never the runtime behavior.

  The generic rationale for each — why a prose-only proof never satisfies a behavioral claim — lives in the FO/ensign contract; the disciplines above are the dev-shape applications a captain opts into.

- **Declaring a posture (optional).** A workflow that wants a single findable answer to "how much engineering does this project want?" may declare it here: project maturity (prototype / product), default test depth, infra-addition policy (may a worker add a CI lane or lint unasked?), and review-finding priority. This is a place to write an existing posture down, not a new required concept — omit it unless the workflow benefits from a stated posture.

## Workflow State

View the workflow overview:

```bash
spacedock status --workflow-dir {dir}
```

Output columns: ID, SLUG, STATUS, TITLE, SCORE, SOURCE.

Find dispatchable tasks ready for their next stage:

```bash
spacedock status --workflow-dir {dir} --next
```

## Task Template

```yaml
---
id:
title: Task title here
status: backlog
source:
started:
completed:
verdict:
score:
worktree:
issue:
pr:
mod-block:
---

## Problem

{What is broken or missing, and why it matters now.}

## Proposed approach

{How the implementation will address the problem. Concrete enough that a worker can start.}

## Acceptance criteria

Each AC names a property of the finished task (not a stage action) and how it is verified.

**AC-1 — {End-state property.}**
Verified by: {test name / command output or exit code / file the change produces / resulting on-disk state — something outside this task body that a future reader can reproduce and that can fail; name the concrete change that would make it fail.}

## Test plan

{What tests verify the implementation, estimated cost, whether E2E is needed.}

## Out of scope

{What this task deliberately does not address.}
```

## Commit Discipline

- Commit status changes at dispatch and merge boundaries
- Commit task body updates when substantive
- Implementation commits land on the worktree branch; merge to main happens via the `pr-merge` mod after PR review

## Adoption

### Pre-fill stages

```yaml
- name: backlog
  initial: true
  gate: true
- name: ideation
  gate: true
- name: implementation
  worktree: true
  context-sections:
    - Review-finding disposition
- name: validation
  worktree: true
  fresh: true
  feedback-to: implementation
  gate: true
  context-sections:
    - Review-finding disposition
- name: done
  terminal: true
```

### Apply layers

- **repo-mutation**: fires on `implementation` and `validation` because both stages operate inside a worktree against the codebase. The `worktree: true` flag on those stages is the structural consequence; `validation` adds `fresh: true` to ensure independent perspective.

### Offer mods

- **pr-merge**: install by default. The shipping ritual for this template is open-a-PR-get-review-merge-to-main, which is exactly what `pr-merge` automates. Surface the offer in Phase 1 with this framing:

  > Because this workflow ships code via PR review, I'll install the **pr-merge** mod. This is the structural reason your stages can stay clean — you don't need a `pr_open` or `awaiting_merge` stage to model the PR step. The mod tracks PR state on the `pr` field, watches for merges in the background, and advances the entity to `done` when the PR lands. Stages describe where work is happening; the PR lifecycle is mod-managed.

  Skip the offer only if the captain explicitly says "no PR review" or "we commit straight to main."

### Inject entity-template snippet

Use the development snippet (problem / proposed approach / acceptance criteria / test plan / out-of-scope) shown in the Task Template section above.

### Surface variants

None. Development is a single coherent shape. Captains who want a different code-shipping cadence (e.g., direct-to-main, trunk-based, release-train) can adjust the stage list and skip the pr-merge mod, but those are edits to this template, not separate variants worth pre-baking.

### Confirmation prose

Surface this in Phase 1 once the template is selected:

> I'll set this up as a **development** workflow: tasks move through `backlog → ideation → implementation → validation → done`, with worktrees on `implementation` and `validation` and `validation` running with a fresh agent so it independently checks the AC.
>
> ID style is **sd-b32** by default, because development workflows usually involve multiple worktree branches creating tasks in parallel and you want the IDs to reconcile without coordination. If this workflow is single-writer (just you, no concurrent branches), I can switch to sequential — let me know.
>
> Default mod: **pr-merge** (PR lifecycle automation, framing above). I'll confirm the install at file generation time.
