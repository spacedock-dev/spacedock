---
id: 42chs9dh7nq22f8at4szvbxp
title: Tune the dev task template for gated stages
status: ideation
source: Captain direction, 2026-08-13
sprint-readiness: ready
score: 0.8
gates:
    version: 1
    records:
        - id: gate:42chs9dh7nq22f8at4szvbxp:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:42chs9dh7nq22f8at4szvbxp-backlog-1
              briefing:
                id: briefing:42chs9dh7nq22f8at4szvbxp:backlog:attempt-1:revision-1
                digest: sha256:b8e080bbf96fbf61ee83cbefef2c0f8e0d1e497f4484c9b486e0094587a4e2e3
                request-digest: sha256:7f4b651b2b1eff1b7e4dbc257f9a1dee993d0e12c148f6d8ef4e9614fecb3705
                room-ref: ./tune-dev-template-for-gated-stages/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:42chs9dh7nq22f8at4szvbxp:backlog:1
                briefing: briefing:42chs9dh7nq22f8at4szvbxp:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T20:36:58.019759Z"
                decision: approve
                reason: 'Captain directive 2026-08-15: fold stack support for pr-mod into 42c and dispatch on top of the stack'
              application:
                target-stage: ideation
                state: consumed
---

Tune the reusable dev task template so task authors supply the decision evidence required by each gated stage without making task files verbose.

## Problem

The stage definitions declare authoritative `Gate content`, but the reusable task template does not prompt authors for six of the thirteen decision inputs those rows demand. Authors invent a section per gate instead: `3w1`, `r5`, and `pnc` each hand-rolled `## Selected approach`, `## Risk evidence`, and `## Semantic changes`; `repair-codex-smallest-sufficient-mechanism-regression` coined `## Risk evidence and dependency` and `## Declared semantic changes` for the same rows; and `### Feedback Cycles` lands in three different places across four entities (nested under a Stage Report in two, top-level after `## Test plan` in one). Every invention is an author guessing at a contract the template should have handed them.

The same shape appears in the `pr-merge` mod. The 0.27 stack (PRs #699-#710) ran the merge ceremony against stacked PRs, and the mod's front half assumed the base is always the trunk and the PR is always created by `gh pr create`. The FO worked around it with `gh stack submit` followed by an after-the-fact body repair, which left PR #699 carrying the branch-derived title `stack27/01 trim dispatch core stale prose` instead of its entity title `Trim stale prose in the dispatch core contract`. The workaround is undocumented, so the next stack repeats it.

## Proposed approach

Two edits, one file each. Neither restates a stage contract; both point at the contract that already exists.

**`docs/dev/README.md` — close the six uncovered gate-content inputs with two new prompts and two amended ones.** The gated stages are `backlog`, `ideation`, and `validation`. Their `Gate content` rows decompose into fifteen decision inputs, two of which the validation Stage Report owns rather than the task body (`non-empty Stage Report results`, `whether delivery can proceed`), leaving thirteen the template owes:

| # | Stage | Decision input | Today | After |
|---|---|---|---|---|
| B1 | backlog | seed outcome | lead paragraph | unchanged |
| B2 | backlog | included scope | *missing* | `## Problem` prompt gains "what a fix must cover" |
| B3 | backlog | excluded scope | `## Out of scope` | unchanged |
| B4 | backlog | proof needed to decide whether design should start | *missing* | `## Risk evidence` backlog line |
| I1 | ideation | selected approach | `## Proposed approach` | prompt marked ideation-owned, asks for the rejected alternative |
| I2 | ideation | risk evidence | *missing* | `## Risk evidence` ideation line |
| I3 | ideation | expected files and lines | `## Expected surface and tolerance` | unchanged |
| I4 | ideation | tolerance | *missing* | `tolerance {±NN%}` field |
| I5 | ideation | semantic changes | *missing* | `Semantics this may change:` field |
| I6 | ideation | proposed proof per AC | `Verified by:` | unchanged |
| V1 | validation | checks run | `## Test plan` | unchanged |
| V2 | validation | evidence per AC | `Verified by:` | unchanged |
| V3 | validation | reviewer findings under workflow labels | *missing* | `### Feedback Cycles` stub |

Coverage moves 7/13 to 13/13 for +10 template lines. The design adds exactly two headings and reuses the names authors already coined (`## Risk evidence`, three entities; `### Feedback Cycles`, four entities and two references in the shipped FO contract) rather than inventing a third vocabulary.

Three deliberate economies keep the template lean. B2 folds into the existing `## Problem` prompt instead of renaming `## Out of scope` to `## Scope`, which would strand the 65 entities using the current heading for no coverage gain. B4 and I2 share one `## Risk evidence` section with a labelled line each, because both ask the same question — what evidence settles the open question at this gate — at different stages. I4 and I5 become fields inside `## Expected surface and tolerance`, whose heading already promises a tolerance the body never asked for.

Chosen direction stays in the ideation-owned `## Proposed approach` prompt and appears nowhere else, per the captain's earlier ruling that it is not a generic task-template field.

**`docs/dev/README.md` — extend the freeze rule rather than duplicate it.** The `## Testing Resources` freeze paragraph already says a worker routes post-completion findings to the first officer. It does not say what the FO then does with them. Append that missing half to the same paragraph: repeated clarification on one gate routes by snapshot currency — changed entity evidence means withdraw and re-prepare, presentation-only thinness means leave the bound snapshot alone and sharpen that stage's `Gate content` for future gates. No new section; the freeze rule and the read-quarantine exemption are cited, never restated.

**`docs/dev/_mods/pr-merge.md` — add a `### Stacked mode` subsection under `## Hook: merge`.** Four amendments, all verified against the live repository:

1. **A stack-sibling base is a valid base.** Resolve the layer's parent from `gh stack view --json` (`branches[].name`, `.base`, `.needsRebase`, `.pr.number`) instead of `dispatch trunk`. Re-record `CANDIDATE_SHA` after any restack, because a lower layer merging replaces every recorded candidate above it — observed when #709 was restacked onto `main`.
2. **Create layers with `gh pr create`, join them with `gh stack link`.** `gh stack submit` exposes no `--title` or `--body-file` flag; a non-interactive run (every FO run) takes auto-generated titles and opens drafts unless `--open`. The reviewed bytes therefore cannot reach GitHub through it. `gh pr create --base {sibling}` already takes the reviewed body file, and `gh stack link` documents that it reuses existing open PRs and chains their bases. This keeps the mod's existing exact-reviewed-bytes discipline intact instead of carving an exception into it.
3. **Repair bodies with the REST API, never `gh pr edit`.** `gh pr edit` fails closed on this repository.
4. **Live-lane approval is what the stack economizes, not lane registration.** The lanes register on every layer and sit `WAITING`; approving the bottom layer (first increment against the trunk) and the top layer (the composed tree) exercises every tree the intermediate layers carry. Stated as the stacked realization of the README's diff-to-lane rule, with the carve-out that a layer whose diff reaches a lane the top does not compose still needs its own approval.

The back half needs no stacked case at all, and the mod will say so explicitly: `MERGED` detection, the `pr-merge:` sentinel, and `merge guard` read the PR's state and the entity's row, never its base. Encoding that as a statement is the point — today it is an accident nobody has written down.

## Out of scope

Do not change gate authority, lifecycle behavior, product code, or live-runtime grading. Do not migrate the 65 existing entities that use the current headings. Do not edit the three `Gate content` rows — the template moves to meet them, never the reverse. Do not add a committed test that reads `docs/dev/README.md` or the mod (see `## Risk evidence`).

## Risk evidence

Backlog: superseded — design was approved to start at the backlog gate on 2026-08-15.

Ideation: four unverified mechanisms carried the design; all four were spiked against the live repository before this body was written, on `gh` 2.68.1.

- **`gh pr edit` vs `gh api PATCH`.** Round-tripped PR #710's own body bytes through both. `gh pr edit --body-file` exits **1** with `GraphQL: Projects (classic) is being deprecated ... (repository.pullRequest.projectCards)` and writes nothing. `gh api --method PATCH repos/spacedock-dev/spacedock/pulls/710 -f body=...` exits **0** and returns the updated timestamp. The captain's "observed failing / observed working" is confirmed, with the exact error text the mod prose now cites.
- **`gh stack submit` cannot carry a title or body.** Its full flag set is `--auto`, `--open`, `--remote`, `--help`; the help states a non-interactive terminal skips the editor and uses auto-generated titles, creating drafts unless `--open`. This refutes the "build bodies BEFORE `gh stack submit`" option in the captain's scope note as written — there is no channel to hand it a body — and is why the design creates layers with `gh pr create` and joins with `gh stack link` instead.
- **`gh stack link` accepts already-open PRs.** Its help states arguments are given bottom-to-top, branches are pushed automatically, and "for branches that already have open PRs, those PRs are used." So creating first and linking second is a supported path, not a workaround.
- **Live-lane distribution across layers.** Queried #699 (bottom), #703 (middle), #710 (top): all three register `claude-live` and `codex-live` in `status=WAITING`, with `pi-live` `SKIPPED`. This **partially refutes** the captain's note that "required checks fire on the bottom and top layers" — they fire everywhere and wait everywhere. What the stack actually economizes is deployment *approval*, so the mod prose states the approval rule, not a firing rule.

The template half is instruction text consumed by a model, so its one machine-facing property was exercised directly rather than assumed: the proposed after-template was piped to `spacedock status --new` in a scratch `sd-b32` workflow, which exited **0**, minted an id, and preserved all nine sections including the two new ones; `status --validate` then reported **VALID**. Its 51-line body was measured the same way the 41-line baseline was. That property also stays guarded in the suite by `internal/ensigncycle/filing_readme_template_test.go`, which runs against its own synthetic fixture README and therefore needs no edit.

AC-5's mechanism needs no spike: the freeze rule already records three occurrences on 2026-08-15 of a post-`prepare` entity commit invalidating the briefing digest, so digest sensitivity to entity bytes is proven behavior this design only routes around, never introduces.

## Expected surface and tolerance

Estimate: +23 net LOC across 2 files, tolerance ±40%.
Semantics this may change: entity-body section vocabulary gains `## Risk evidence` and `### Feedback Cycles`, and the FO's stacked merge ceremony gains a documented create path. No command grammar, no stored format, no gate authority, and no lifecycle behavior changes.

## Acceptance criteria

**AC-1 (VALUE) — Every gate-content decision input the task body owes is prompted by the template, at 13/13 versus 7/13 today, without the template body exceeding 55 lines or the `Gate content` rows moving.**
Verified by: walk each of the thirteen rows in the `## Proposed approach` table against the rendered template and name the prompt covering it; `awk` the fenced template body's line count (41 today, 51 as designed, cap 55); `git diff` against the stack tip shows no change to the `- **Gate content:**` lines of `backlog`, `ideation`, and `validation`. This is the one-off comparison the Proof policy permits as validation evidence, not a committed test. Falsifying edits: deleting the `## Risk evidence` block drops coverage to 11/13; expanding each prompt to a paragraph pushes the body past 55 lines; weakening a `Gate content` row so a thin prompt covers it fails the byte-unchanged check.

**AC-2 — Chosen direction is requested only by the ideation-owned `## Proposed approach` prompt.**
Verified by: the rendered template's backlog-owned prompts (`## Problem`, `## Out of scope`, and the `## Risk evidence` backlog line) request scope and decision proof and never a direction; the validation-owned prompt (`### Feedback Cycles`) requests correction-round findings and never a direction. Falsifying edit: adding "and the approach chosen" to the `## Problem` prompt puts direction in a backlog-owned field.

**AC-3 — The estimate prompt requests one signed net figure, with tolerance and declared semantics as separate fields.**
Verified by: the rendered `## Expected surface and tolerance` block reads `Estimate: {+NNN} net LOC across {M} files, tolerance {±NN%}.` followed by a `Semantics this may change:` line. Falsifying edit: dropping the tolerance field returns ideation input I4 to uncovered and fails AC-1 at 12/13.

**AC-4 — The change is confined to `docs/dev/README.md` and `docs/dev/_mods/pr-merge.md`, with no existing task file, gate record, or test rewritten.**
Verified by: `git diff --name-only` against the stack tip lists exactly those two paths; `spacedock status --workflow-dir docs/dev --validate` reports VALID; `go test ./...` and `go test ./... -race` stay green with no test edited, including `internal/ensigncycle/filing_readme_template_test.go`, which reads its own synthetic fixture README and so must pass untouched. Falsifying edit: migrating the 65 entities using the current headings puts entity paths in the diff.

**AC-5 — Repeated clarification on one gate routes by snapshot currency, carried as a continuation of the existing freeze rule rather than a second rule.**
Verified by: the `## Testing Resources` freeze paragraph carries the routing sentences and no new section restates the freeze rule or the read-quarantine exemption; and a behavior exercise on a scratch split-root workflow shows the two routes diverge — a `gate prepare`, then an entity-body commit, then a re-read shows the briefing digest no longer matches the entity bytes (the withdraw-and-re-prepare route), while editing a `Gate content` row with no entity commit leaves the bound digest matching (the presentation-only route). Falsifying edits: adding a standalone "Gate clarification" section fails the no-duplication half; if the digest were insensitive to entity bytes, both routes would look identical and the exercise would show no divergence.

**AC-6 (VALUE) — A stack layer's PR carries the captain-reviewed title and body from creation, with no post-creation repair.**
Verified by: the PR that ships this task — dispatched as the next layer on `stack27/11-self-contained-contracts` — has `title` equal to this entity's frontmatter `title:` and `body` byte-equal to the reviewed `PR_BODY_FILE`, both established by `gh pr create --base {stack sibling}` and joined with `gh stack link`, with no `gh api PATCH` in the ceremony record. The independent baseline that moved the wrong way is PR #699: created through `gh stack submit`, it still carries the branch-derived title `stack27/01 trim dispatch core stale prose`. Falsifying edit: creating this layer with `gh stack submit --auto` reproduces a branch-derived title and a draft. If the captain elects not to stack this task, this AC is unverifiable as written and returns to the gate rather than being satisfied on a synthetic stack.

**AC-7 — The mod states that the back half needs no stacked case, and routes body repair away from `gh pr edit`.**
Verified by: the diff leaves the `## Hook: startup` and `## Hook: idle` prose byte-unchanged, which is the claim "the back half needs no stacked case" made checkable rather than asserted; and the `gh pr edit` exit-1 / `gh api PATCH` exit-0 contrast recorded in `## Risk evidence` is reproducible on any open PR in this repository. Falsifying edit: adding a stacked branch to the startup hook puts those lines in the diff and contradicts the statement the same commit makes.

## Test plan

No new standing test. The template half's machine-facing property is already covered by `internal/ensigncycle/filing_readme_template_test.go` against its own fixture, and the read-quarantine exemption in the Proof policy does not stretch to a test whose oracle would be this README's wording — that is the banned prose-grep, and the captain's 2026-07-20 ruling puts such a grep in the validation report, never in the suite.

Validation therefore runs: the thirteen-row coverage walk and template line count (AC-1); the rendered-template reads for AC-2 and AC-3; `git diff --name-only`, `spacedock status --workflow-dir docs/dev --validate`, `go test ./...`, and `go test ./... -race` (AC-4, AC-7); the scratch-workflow `gate prepare` digest exercise (AC-5); and `gh pr view --json title,body` on the delivered stack layer against this entity's `title:` and the reviewed body file (AC-6). Cost is low: one scratch workflow, one Go suite run, and reads against PRs that already exist.

## Captain-directed scope addition (2026-08-15)

Fold stacked-PR support into the pr-merge mod (docs/dev/_mods/pr-merge.md), from the 0.27 stack experience:

- The mod's back half (MERGED detection, sentinel, merge guard) worked unchanged on stacked PRs and needs no edit - encode that as an explicit statement, not an accident.
- The front half must gain a stacked mode: the FO builds template-conformant PR bodies BEFORE `gh stack submit` (or patches them immediately after via the REST API - the GraphQL path chokes on deprecated projectCards), records the per-layer candidate SHA, and accepts a stack-sibling base as valid where the mod today assumes the trunk.
- Body edits on existing PRs must use `gh api --method PATCH repos/{repo}/pulls/{n}` (observed working) rather than `gh pr edit` (observed failing).
- Live-lane economics note for the mod prose: required checks fire on the bottom and top layers; the top's tree is the composition, so one live approval covers the stack.
- Evidence base: PRs #699-#710, stack #707, and this session's ceremony record.

Design constraint: the mod is blocked-product (dispatched-worker territory); the README template is FO-owned process. Ideation designs both edits; implementation is dispatched on top of the current stack (worktree based on the stack tip) and lands as the next stack layer.

## Proposed edits

Exact before/after wording, per the ideation stage's "specific before/after wording, not just change X" rule.

### Edit 1 — `docs/dev/README.md`, the fenced body of `## Task Template` (lines 242-269)

Before:

    Brief description of this task and what it aims to achieve.

    ## Problem

    {What is broken or missing, and why it matters. Ideation fills this in.}

    ## Proposed approach

    {How the task intends to solve the problem. Ideation fills this in.}

    ## Out of scope

    {What this task deliberately does not cover, so the boundary is explicit.}

    ## Expected surface and tolerance

    Estimate net LOC change: {+NNN or -NNN}, across {M} files.

    ## Acceptance criteria

    Each AC names a property of the finished entity, not a stage action, and how it is verified.

    **AC-1 - {End-state property.}**
    Verified by: {test name / command output or exit code / file the change produces / resulting on-disk state — something outside this task body that a future reader can reproduce and that can fail; name the concrete change that would make it fail.}

    ## Test plan

    {What verifies the implementation, estimated cost/complexity, and whether fixture, CLI, or live workflow tests are needed.}

After:

    Brief description of this task and what it aims to achieve.

    ## Problem

    {What is broken or missing, why it matters, and what a fix must cover. Backlog seeds it; ideation sharpens it.}

    ## Proposed approach

    {Ideation: the direction chosen, and the simplest alternative rejected with the reason it cannot deliver the value.}

    ## Out of scope

    {What this task deliberately does not cover, so the boundary is explicit.}

    ## Risk evidence

    {Backlog: the check, artifact, or observation that decides whether design should start.}
    {Ideation: the riskiest unverified mechanism and what exercising it showed, or `no spike needed: {the proven mechanisms this relies on}`.}

    ## Expected surface and tolerance

    Estimate: {+NNN} net LOC across {M} files, tolerance {±NN%}.
    Semantics this may change: {command grammar, stored formats, authority, runtime behavior, or `none`}.

    ## Acceptance criteria

    Each AC names a property of the finished entity, not a stage action, and how it is verified. At least one measures the end-value against a baseline that can move the wrong way.

    **AC-1 - {End-state property.}**
    Verified by: {test name / command output or exit code / file the change produces / resulting on-disk state — something outside this task body that a future reader can reproduce and that can fail; name the concrete change that would make it fail.}

    ## Test plan

    {What verifies the implementation, estimated cost/complexity, and whether fixture, CLI, or live workflow tests are needed.}

    ### Feedback Cycles

    {First officer appends one `- Cycle {N}: ...` line per correction round; the validation gate reads reviewer findings from here.}

Net +10 lines: 41 to 51, against the 55-line cap in AC-1. The frontmatter block above the body is unchanged.

### Edit 2 — `docs/dev/README.md`, the freeze paragraph in `## Testing Resources` (line 294)

Before, the paragraph ends:

    ... and a revise decision that re-dispatches the entity lifts the freeze for that rework until its next completion signal.

After, the same paragraph continues:

    ... and a revise decision that re-dispatches the entity lifts the freeze for that rework until its next completion signal. When one gate draws repeated clarification, the missing decision input decides the route: if the entity's evidence changed, withdraw the gate and prepare a new snapshot; if only the presentation was thin, leave the bound snapshot alone and sharpen that stage's `Gate content` for future gates. Never alter a bound snapshot in place.

Net 0 lines added; one source line modified.

### Edit 3 — `docs/dev/_mods/pr-merge.md`, new `### Stacked mode` between the `gh pr create` paragraph (line 69) and `### PR body template` (line 71)

Insert:

    ### Stacked mode

    When the candidate is one layer of a stack, the front half above still owns the title and body; only the base and the create path change. Resolve the layer's parent from `gh stack view --json` — `branches[]` carries `name`, `base`, `needsRebase`, and `pr.number` in stack order — and use the branch below as `$BASE` in place of `dispatch trunk`. A stack-sibling base is a valid base, not an error. Record `CANDIDATE_SHA` per layer and re-record it after any restack: when a lower layer merges, every candidate above it is replaced.

    Create each layer with the same `gh pr create --base "$BASE" --head "$BRANCH" --title "$PR_TITLE" --body-file "$PR_BODY_FILE"` call the front half already prescribes, then join the layers with `gh stack link {pr} {pr} ...` in stack order, bottom to top; it reuses branches that already have open PRs. Do NOT create layers with `gh stack submit`: it exposes no title or body flag, and a non-interactive run takes auto-generated titles and opens drafts, so the captain-reviewed bytes never reach GitHub. PR #699 was created that way and still carries the branch-derived title `stack27/01 trim dispatch core stale prose`.

    To repair an existing PR's title or body, use `gh api --method PATCH repos/{owner}/{repo}/pulls/{N} -f title=... -f body=...`. Do not use `gh pr edit`: its GraphQL read requests the deprecated `repository.pullRequest.projectCards` field and exits 1 without writing (observed on gh 2.68.1).

    Required live lanes register on every layer and sit at `status=WAITING` for deployment approval. Approve the bottom layer, whose tree is the first increment against the trunk, and the top layer, whose tree is the whole stack composed; no intermediate layer carries a tree those two do not already exercise. This is the stacked realization of the workflow's diff-to-lane rule, not an exemption from it: a layer whose diff reaches a lane the top layer does not compose still needs that lane approved on its own PR.

    The back half needs no stacked case. `MERGED` detection, the `pr-merge:` sentinel, and `merge guard` read the PR's state and the entity's row, never its base, so a stack-sibling base flows through the startup and idle hooks unchanged. Do not add a stacked branch to them.

Net +13 lines.

## Stage Report: ideation

- DONE: Design both halves: the gated-stage template prompts (original 42c scope, AC-1 to AC-5) and the pr-merge mod stacked-mode amendment (captain scope addition in the body) - reconcile with the freeze rule and read-quarantine exemption already in the README rather than duplicating them
  `## Proposed edits` carries verbatim before/after for all three edits; the freeze rule is extended in place (Edit 2, 0 added lines) and the read-quarantine exemption is cited in `## Test plan` as the reason no committed test ships.
- DONE: Design against the stack tip (branch stack27/11-self-contained-contracts, head d1d8f745f); implementation will land as the next stack layer
  `git diff main d1d8f745f` is empty for both target files, so the line numbers in `## Proposed edits` hold at the tip; AC-6 verifies on the delivered layer itself.
- DONE: AC evidence per the Proof policy: template prompts verified against each Gate content row; mod prose changes get the falsifying-edit statement, never a committed prose-grep
  AC-1 is the 13-row walk as a one-off validation exercise; all seven ACs name a falsifying edit; `## Test plan` ships no standing test.

### Summary

Two edits, one file each, closing the same gap in two places: authors and the FO both invent structure the contract should have handed them. The template covers 7 of the 13 gate-content decision inputs a task body owes; the design takes it to 13/13 for +10 lines by adding two headings authors already coined (`## Risk evidence` in three entities, `### Feedback Cycles` in four) rather than a third vocabulary. The `pr-merge` mod gains a `### Stacked mode` subsection and an explicit statement that its back half needs no stacked case.

Four unverified mechanisms were spiked live before the design was written, and one partly refuted the captain's scope note. `gh stack submit` exposes no title or body flag and takes auto-generated titles non-interactively, so "build bodies BEFORE `gh stack submit`" has no channel to work through — the design creates layers with `gh pr create --base {sibling}` and joins them with `gh stack link`, which preserves the mod's existing exact-reviewed-bytes discipline instead of carving an exception into it. `gh pr edit` exits 1 on the deprecated `projectCards` field where `gh api --method PATCH` exits 0, both observed on PR #710. Live lanes register and sit `WAITING` on every layer, not just bottom and top, so the mod prose states an approval rule rather than a firing rule.

Two points want the captain's eye at the gate. The live-lane approval rule touches the README's "required lanes are a function of the diff" policy: the design frames per-layer approval as that rule's stacked realization, with a carve-out for any layer whose diff reaches a lane the top does not compose, but it is the one place this task edges into CI policy. And AC-6 is verifiable only if this task actually ships as a stack layer; if the captain elects not to stack it, that AC returns to the gate rather than being satisfied on a synthetic stack.
