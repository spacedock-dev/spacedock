---
id: 42chs9dh7nq22f8at4szvbxp
title: Tune the dev task template for gated stages
status: validation
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
        - id: gate:42chs9dh7nq22f8at4szvbxp:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:42chs9dh7nq22f8at4szvbxp-ideation-1
              briefing:
                id: briefing:42chs9dh7nq22f8at4szvbxp:ideation:attempt-1:revision-1
                digest: sha256:ea7199bbce8f4e01dedf066c870f82a43b2045229f35d90e4e9856526fc06212
                request-digest: sha256:a0b4129f1df539c8b2a8a43d221268a58637c614bee3420704c755ff098c8c60
                room-ref: ./tune-dev-template-for-gated-stages/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:42chs9dh7nq22f8at4szvbxp:ideation:1
                briefing: briefing:42chs9dh7nq22f8at4szvbxp:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T21:24:21.17569Z"
                decision: approve
                reason: 'Captain approved 2026-08-15 with both rulings as recommended: per-layer live approval framed as the diff-rule''s stacked realization with the named carve-out; ships as a real stack layer'
              application:
                target-stage: implementation
                state: consumed
        - id: gate:42chs9dh7nq22f8at4szvbxp:validation
          stage: validation
          attempts:
            - id: gate-attempt:42chs9dh7nq22f8at4szvbxp-validation-1
              briefing:
                id: briefing:42chs9dh7nq22f8at4szvbxp:validation:attempt-1:revision-1
                digest: sha256:8007ccce0da037accf88b7d620de2edf8af84611a2da8ebe65453e17144dc4ef
                request-digest: sha256:59a8744a48133f67bb0769e5e9f0edd25647144c526286a635deca1c4587a9e0
                room-ref: ./tune-dev-template-for-gated-stages/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:42chs9dh7nq22f8at4szvbxp:validation:1
                briefing: briefing:42chs9dh7nq22f8at4szvbxp:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T22:59:38.494149Z"
                decision: approve
                reason: 'Captain 2026-08-15: behavior assumed - validation PASSED approved, the merged-ancestor Stacked-mode amendment ratified, and the exercised deviation shape (worker holds, FO authorizes with distinct scope reasoning, captain ratifies at the next gate) adopted for the undiscoverable-at-ideation class'
              application:
                target-stage: done
                state: superseded
            - id: gate-attempt:42chs9dh7nq22f8at4szvbxp-validation-2
              briefing:
                id: briefing:42chs9dh7nq22f8at4szvbxp:validation:attempt-2:revision-1
                digest: sha256:a3d6c50e877c40237e702f17b702ee65cca1228d48acde2fc5e2596f6cda016f
                request-digest: sha256:90784a609379b6c84274589a91dbf47e7cfd57dc89491e5585355912135aa163
                room-ref: ./tune-dev-template-for-gated-stages/review/validation/briefing-2
worktree: .worktrees/spacedock-ensign-tune-dev-template-for-gated-stages
started: 2026-08-15T21:25:11Z
pr:
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

## Stage Report: implementation

- DONE: Execute the gated design exactly: 13/13 template prompts reusing author-coined headings, and the pr-merge Stacked mode subsection with the create-then-link shape; both captain rulings approved as recommended
  Commit `d89916600`, the three `## Proposed edits` blocks applied verbatim. A scripted 13-row walk over the rendered template reads 7/13 at the stack tip `d1d8f745f` and 13/13 after, closing exactly B2, B4, I2, I4, I5, V3; fenced body 41 to 51 lines against the 55 cap. Zero `- **Gate content:**` lines appear in the diff, so the rows were met rather than moved. Both rulings landed as recommended: live-lane approval is stated as the diff-rule's stacked realization with the named carve-out, and the task ships as a real stack layer (item 3).
- DONE: Verify contractlint green including the generalized write-core fixture after the README edits; falsifying-edit statements recorded, no committed prose-grep
  `go test ./internal/contractlint/` ok. Its write-core fixture (`TestFOWriteCoreMutationGateClassifiesTargets`, `TestFOWriteCoreMutationGateRequiresExactOverride`) asserts the classifier sorts `{workflow}/README.md` to allowed-process and `{workflow}/_mods/pr-merge.md` to blocked-product across two workflow dirs; it would red if either edited path changed class or if the table only worked for the repo it shipped from. `go test ./...` and `go test ./... -race` both exit 0 across 22 packages. `internal/ensigncycle` carries `filing_readme_template_test.go`, which asserts a workflow-local Task Template's literal fenced body round-trips through `new` (mints an id, preserves frontmatter) while an indented fence is rejected; it reads its own synthetic fixture and passed untouched, which is the claim AC-4 makes. `git diff --name-only` lists two doc paths and no `*_test.go`, so no prose-grep was committed; the falsifying edits live in AC-1 through AC-7 above.
- DONE: Land as the next stack layer using your own designed mechanics (gh pr create --base onto the current top, then gh stack link); base worktree is the stack tip
  PR #711, candidate `d89916600`. Worktree branch `spacedock-ensign/tune-dev-template-for-gated-stages` is based on `d1d8f745f`, the head of `stack27/11-self-contained-contracts` and the last unmerged layer. Created with `gh pr create --base stack27/11-self-contained-contracts`, a stack-sibling base rather than the trunk, carrying the reviewed title and `--body-file`. `gh pr view 711` returns `title` equal to this entity's frontmatter `title:`, `baseRefName` equal to the sibling, `isDraft` false, and a body whose sha256 equals the reviewed `PR_BODY_FILE` after trailing-newline normalisation (the 1-byte delta is the shell redirect's newline, not content). No `gh stack submit`, no `gh pr edit`, no `gh api PATCH` anywhere in the record, so the create-then-link shape and the no-post-creation-repair half of AC-6 both hold. The first officer owns the `pr:` field.
- DONE: the `gh stack link` join clause of AC-6 — executed, its failure mode discovered and guarded, ungrouped state accepted under FO authorization
  The link ran and did not produce stack membership: `gh stack link 699 ... 711` printed `⚠ failed to update base branch for PR #709 to stack27/09-trim-version-output: HTTP 422 PullRequest.base is invalid` and then `✓ Updated stack to 12 PRs (stack #707)`, while `gh stack view 707 --json` still returns 11 branches with #711 absent; a narrower `gh stack link 707 711` separately claimed `ℹ PR #711 is already in stack #707 — skipping`, so gh's membership check contradicts its own enumeration. The cause is that #709's predecessors merged, making `stack27/09-trim-version-output` an ancestor of `main` (`git merge-base --is-ancestor`, true), so GitHub rejects resetting #709's base onto it and the record rewrite aborts. Per FO authorization (2026-08-15) this is the healthy end-state of a stack, not a defect: the PR base chain #709 on `main` to #710 to #711 is coherent, stack-object membership is bookkeeping, and #709 is explicitly not to be restacked. The failure mode is now guarded in `### Stacked mode` — verify membership with `gh stack view --json` after every link, treat a merged-ancestor layer as complete rather than re-linkable. Recorded as a declared deviation below for the captain's ruling at the validation gate.

### Summary

The three approved edits landed verbatim in `d89916600`: the task template closes the six uncovered gate-content inputs (7/13 to 13/13, 41 to 51 lines) using the two headings authors had already coined, the freeze paragraph gains snapshot-currency routing in place with no added line, and the pr-merge mod gains `### Stacked mode` plus an explicit statement that its back half needs no stacked case.

Two measurements are worth the validator's eye. The 13-row walk is only honest with wording-neutral probes: a first pass using the new template's phrasing scored the baseline 5/13 instead of the designed 7/13, because the probes for I1 and I3 matched only post-edit wording. Re-running with probes that detect a prompt for the input in either version reproduces 7/13 to 13/13 exactly. And actual surface is +22 net LOC across 2 files against the declared estimate of +23 with tolerance +/-40%, a -4% deviation; the mod insert measured 12 added lines where the design said 13, an off-by-one in the design's own count with no content difference.

AC-5's digest-divergence exercise and AC-6's delivered-PR check are the validation stage's to run, per this task's `## Test plan`; the mechanics AC-6 inspects are recorded in item 3.

One thing needs the captain's ruling at the validation gate: `### Stacked mode` carries a paragraph that was not in the verbatim-approved ideation text, guarding a `gh stack link` failure mode that only became observable once a layer of this stack had merged. The first officer authorized it and it is recorded as a declared deviation below, at +2 lines and no change to any acceptance criterion.

### Finding — resolved under first-officer authorization (2026-08-15)

`gh stack link` reports success it did not achieve, and the layer this task ships is not in the stack record.

- Released user and normal workflow: the first officer running the `pr-merge` ceremony on a stacked candidate, which is the path the `### Stacked mode` subsection this task just added prescribes.
- Observable harm: `gh stack link` exits reporting `✓ Updated stack to 12 PRs` while `gh stack view 707 --json` still returns 11; the new layer is created and correctly based but never grouped, and the FO gets no failure signal to act on.
- Affected value AC: `value-ac[AC-6]` — its verification clause reads "both established by `gh pr create --base {stack sibling}` and joined with `gh stack link`". The title and body halves hold; the join half does not.
- Trigger evidence: the 422 on #709 quoted in the stack-link item above, plus the 11-vs-12 disagreement between `gh stack link` and `gh stack view --json`, both reproducible now against stack #707.

Worker proposal, as three separate facts. Materiality: **Material** — a released workflow's documented step silently no-ops. Ownership: **not this task**. Disposition: **Needs decision**. Candidate bytes and `HEAD` were held unchanged and no reviewer rerun was started while this was outstanding.

First-officer authorization, received 2026-08-15 after consultation, and what it corrected. Materiality **Material** was agreed. Ownership was **corrected**: the prose remedy is this task's, because `docs/dev/_mods/pr-merge.md` is in the approved surface; only the #709-restack remedy is out of scope, and it was declined. The worker's framing of #709 as "base drift" was also **corrected** — a layer sits on the trunk because its predecessors merged, which is a stack's healthy end-state rather than a defect. Disposition authorized: accept the ungrouped state, do not restack #709, amend the `### Stacked mode` prose with the false-success guard, and restate the join clause of AC-6 from FAILED to DONE under that resolution.

**Declared deviation for the validation gate.** The `### Stacked mode` wording was approved verbatim at the ideation gate and now carries one paragraph that was not in that approved text: the merged-ancestor 422, the false `✓ Updated stack to {N} PRs`, the instruction to verify membership with `gh stack view --json` after every link, and the rule that a merged-ancestor layer is complete rather than re-linkable. It is an implementation-discovered truth, not a design change — the mechanism was unverifiable before a layer of this stack had actually merged under a `gh stack link`. Surface effect is +2 lines, taking the task to +24 net LOC against the declared estimate of +23 with tolerance +/-40%, a +4.3% deviation still inside the band. Acceptance criteria are unchanged and unnarrowed. The captain rules on this amendment at the validation gate.

## Stage Report: validation

- DONE: Independently re-exercise the ACs against code commit ab112b6d3, never by reading the report: 13/13 gate-content coverage at 51 lines against the 55 cap, exactly two files changed, the startup and idle hooks byte-unchanged, PR #711 body sha256 e3ffd75ec90b5ab14a31fa42 matching the reviewed bytes with no post-creation edit in the record
  All re-run from the entity worktree at ab112b6d3; per-AC citations below.
- DONE: Verify the Stacked mode prose including the authorized merged-ancestor amendment states observed behavior (the 422, the false success, gh stack view --json as the membership read); the amendment is a declared FO-authorized deviation for the captain to ratify at this gate
  `git merge-base --is-ancestor` confirms stack27/09-trim-version-output is an ancestor of main; the open chain #709 (base main) -> #710 -> #711 is coherent; `gh stack view --json` from the ceremony clone enumerates 11 branches with #711 absent, matching the recorded 11-vs-12 false success; #711's lanes sit claude-live/codex-live WAITING, pi-live SKIPPED, as the approval-rule prose assumes. The amendment is +2 lines (commit ab112b6d3), ACs unchanged, surface +24 net LOC vs the +23 +/-40% estimate.
- DONE: Suites green plain and -race; status validate VALID; verdict PASSED or REJECTED with per-AC citations; flag for the captain: the undiscoverable-at-ideation deviation class question the implementer raised
  `go test ./...` and `go test ./... -race` both exit 0 across all 22 packages in the worktree; `status --validate` prints VALID, exit 0 (four pre-existing archived superseded-verdict warnings, owned by the peer vocabulary entity). Verdict and captain flags below.

### Per-AC citations

- AC-1 PASS: wording-neutral 13-row walk reads 7/13 at d1d8f745f (B1 B3 I1 I3 I6 V1 V2) and 13/13 at ab112b6d3, closing exactly B2 B4 I2 I4 I5 V3; fenced body 41 -> 51 lines against the 55 cap; zero `- **Gate content:**` lines changed in the diff (the one grep hit is the freeze paragraph citing the term). Would fail if a row were weakened or the `## Risk evidence` block dropped.
- AC-2 PASS: the rendered template requests a direction only in the ideation-owned `## Proposed approach`; `## Problem`, `## Out of scope`, and the `## Risk evidence` backlog line ask scope and decision proof, `### Feedback Cycles` asks correction rounds.
- AC-3 PASS: the block reads `Estimate: {+NNN} net LOC across {M} files, tolerance {±NN%}.` with a separate `Semantics this may change:` line; the 51-line body piped to `new` in a scratch workflow exits 0, mints an id, preserves all nine sections including the ± byte, and validates.
- AC-4 PASS: `git diff --name-only d1d8f745f..ab112b6d3` lists exactly docs/dev/README.md and docs/dev/_mods/pr-merge.md; no entity or test file in the diff; filing_readme_template_test.go untouched and green in both suite runs.
- AC-5 PASS: the freeze paragraph carries the routing sentences in place with no new section; scratch split-root exercise: the briefing pin (rev 07c903cb...) still matches the committed entity bytes after sharpening a `Gate content` row with no entity commit, and diverges (8ad13181...) after a post-prepare entity commit - the two routes are observably distinct.
- AC-6 PASS: #711 title equals the frontmatter `title:`; the live body minus its single trailing newline hashes e3ffd75ec90b5ab14a31fa42b507eee6867d1c3c76d8c534cabfb0a7aab9cd9f; baseRefName is the stack sibling stack27/11-self-contained-contracts, headRefOid ab112b6d3; lastEditedAt null and userContentEdits 0, so no post-creation repair; baseline #699 still carries the branch-derived title. The join half holds as the FO-authorized ungrouped end-state.
- AC-7 PASS: the `## Hook: startup` and `## Hook: idle` spans are byte-identical by `cmp` across the diff range (single hunk, confined to `## Hook: merge`); `gh pr edit` round-tripping #710's own bytes exits 1 with the projectCards error on gh 2.68.1 and leaves no new edit record, while #710's three 2026-08-15 userContentEdits corroborate the spike's `gh api PATCH` writes.

### Reviewer findings

- Deferred risk (recommendation for Review-finding disposition; no candidate change made): the prescribed membership read is clone-local. `gh stack view --json` succeeds only in the clone holding `.git/gh-stack` (the FO ceremony clone, 11 branches enumerated); a fresh clone on member branch stack27/10 and the candidate branch itself both exit 2 `not part of a stack`. Released workflow: FO ceremony in its own clone, which works. Harm: a takeover FO on a fresh clone cannot enumerate membership and could misread the error as an ungrouped stack. Authority: none: no value AC fails on the supported path. Trigger evidence: both exits reproducible now. Promote to material if a ceremony must run from a clone lacking `.git/gh-stack` for an existing stack.
- No material findings; no polish findings.

### Summary

All seven ACs re-exercised independently against ab112b6d3 hold: recommend PASSED. Two rulings for the captain at this gate: ratify the merged-ancestor amendment to `### Stacked mode` (declared FO-authorized deviation, +2 lines, ACs unchanged and unnarrowed, +24 net LOC inside the +/-40% band); and the implementer's deviation-class question - whether an implementation-discovered truth that amends verbatim-approved gate text, undiscoverable at ideation because the mechanism only became observable once a stack layer merged, should be a recognized deviation class (FO authorization plus gate-time ratification, as exercised here) or remain ad hoc. One deferred risk recorded above with its promote condition; delivery can proceed.

## Stage Report: implementation (cycle 2)

Captain rejection at the validation feedback gate, 2026-08-15. Two material findings; both dispositions executed.

- DONE: Execute the gated design exactly: 13/13 template prompts reusing author-coined headings, and the pr-merge Stacked mode subsection with the create-then-link shape; both captain rulings approved as recommended
  Retargeted per Finding 1. Commit `cc728f119` reverts both `docs/dev` edits so the branch touches neither our installed workflow nor its installed mod; `git diff --name-only` against the stack tip lists exactly `skills/commission/references/templates/development.md` and `mods/pr-merge.md`, and `git diff -- docs/dev` is empty. Commit `e49154829` carries the corrected deliverable. Coverage was re-derived against the shipped template's own structure rather than assumed from our README, and the two differ: the shipped template starts at **5/13**, not 7/13, because it carries neither a lead description line (B1) nor any `## Expected surface and tolerance` section (I3, I4, I5). Closing all thirteen takes its fenced body 37 to 53 lines against the 55 cap. Its three `- **Gate content:**` rows are byte-unchanged — zero appear in the diff — so the template still moves to meet the contract. `## Risk evidence` and `### Feedback Cycles` are the author-coined headings, unchanged from the approved design. **One deviation to declare:** the live-lane approval ruling cannot ship. Finding 2 excludes lane economics from the shipped mod, so that captain ruling now has no home in this PR; it is guidance for our own README and belongs to the captain's later refit, not to a mod that ships to arbitrary repositories.
- DONE: Verify contractlint green including the generalized write-core fixture after the README edits; falsifying-edit statements recorded, no committed prose-grep
  Read now as "after the shipped-artifact edits", the README edits having been backed out. `go test ./...` and `go test ./... -race` both exit 0 across 20 packages with zero `FAIL` lines, `internal/contractlint` green in both. Its write-core fixture (`TestFOWriteCoreMutationGateClassifiesTargets`, `TestFOWriteCoreMutationGateRequiresExactOverride`) asserts the classifier sorts `{workflow}/README.md` to allowed-process and `{workflow}/_mods/pr-merge.md` to blocked-product across two workflow dirs, and would red if a path changed class or the table only worked for its origin repo. `git diff --name-only` lists no `*_test.go`, so nothing was committed that greps this prose; the falsifying edits remain stated in AC-1 through AC-7.
- DONE: Land as the next stack layer using your own designed mechanics (gh pr create --base onto the current top, then gh stack link); base worktree is the stack tip
  PR #711 remains OPEN with `baseRefName` `stack27/11-self-contained-contracts` and head `e49154829`; `git merge-base --is-ancestor d1d8f745f HEAD` confirms the branch still descends from the stack tip. The layer was created in cycle 1 by `gh pr create --base {sibling}`; for this corrected round the body was rewritten through `gh api --method PATCH`, which is the repair path this mod itself prescribes and explicitly not `gh pr edit`. Stack-object membership is unchanged from cycle 1 and remains the accepted ungrouped state.

### Summary

The captain's two findings moved the deliverable, not its design. Finding 1 relocated it from our installed copies to the artifacts `commission` and `refit` actually ship, and measuring the shipped template fresh proved that instruction load-bearing: its baseline is 5/13 where our README's was 7/13, so assuming the mirror would have shipped a template still missing three ideation inputs and its own lead line. Finding 2 stripped the mod back to generic stack mechanics.

Corrected surface is **+28 net LOC across 2 files** (template +16, mod +12), against the ideation estimate of +23 with tolerance +/-40% — a +21.7% deviation, inside the 13.8 to 32.2 band, with acceptance criteria unchanged and unnarrowed. The composition moved even though the total barely did: the template needed +16 rather than +10 because of the lower baseline, while the mod needed +12 rather than +14 because the live-lane paragraph was dropped.

The CI-agnostic audit was run as a check rather than an assertion: every added mod line was grepped for `lane`, `e2e`, `WAITING`, `environment`, `runtime-live`, `stack27`, this repository's PR numbers and its owner name, all clean, and no added sentence names this repository. What survived is the six mechanics the captain listed as belonging — create-then-link, `gh stack view --json` read-back over the success banner, merged-ancestor-is-complete, the base chain as the ceremony's dependency, reviewed-bytes discipline, and REST PATCH over `gh pr edit`. The `gh stack submit` warning kept its behavioural claim but lost its citation of our PR #699, which a reader of a shipped mod cannot verify.

### Feedback Cycles

- Cycle 1: REJECTED — Captain wrong-target and CI-agnostic scope correction (shipped template + canonical mod, not the installed workflow); surface +28 net/2 files vs estimate +23 (+21.7%, in band); AC unchanged

## Stage Report: validation (cycle 2)

- DONE: Cycle-2 validation, independently re-exercised against code commit e49154829 (revert cc728f119 + corrected deliverable), never from the report: diff touches ONLY skills/commission/references/templates/development.md and mods/pr-merge.md; the docs/dev diff is EMPTY; template coverage 13/13 at 53 lines against the 55 cap from the true 5/13 baseline; the three Gate content rows byte-unchanged
  `git diff --name-only d1d8f745f..e49154829` lists exactly the two shipped paths; `git diff -- docs/dev` empty and docs/dev at cc728f119 is byte-identical to the stack tip; wording-neutral 13-row walk scripted fresh reads 5/13 before (B3 I1 I6 V1 V2) and 13/13 after; fenced body 37 -> 53 vs cap 55; the three `- **Gate content:**` rows hash identically (sha256 d7acb0e2) at unchanged line numbers 74/81/94.
- DONE: CI-agnosticism as a check: added mod lines grep clean for lane, e2e, WAITING, environment, runtime-live, stack27, this repo's PR numbers and owner; no added sentence names this repository; the six generic mechanics survive; the live-lane ruling is correctly ABSENT (homeless by design, deferred to the local refit)
  All 12 added mod lines grep clean for every banned token (only bare number: 422, the HTTP status); no added sentence in either file names this repository or the product; create-then-link, view --json read-back over the success banner, merged-ancestor-is-complete, base-chain-over-membership, reviewed-bytes discipline, and REST PATCH over `gh pr edit` are all present; no lane/approval/WAITING sentence exists in the added text.
- DONE: PR #711 body byte-equal to reviewed file (sha256 d8b44634007308513ac6bff0) via the mod's own REST-PATCH path; suites green plain and -race; contractlint green; verdict PASSED or REJECTED with per-AC citations
  Live body minus the jq-appended newline hashes d8b44634007308513ac6bff0b986dc3991e9914c208f390581e320c1f4985456; edit trail totalCount 2 = creation revision (21:45:29Z = createdAt) + exactly one rewrite (23:45:03Z), so cycle-1's "0 edits" read stands and the one rewrite is the cycle-2 PATCH; `gh pr edit` re-exercised on #711's own bytes exits 1 (projectCards) leaving totalCount 2. `go test ./...` and `-race` both exit 0, 40 ok / 0 FAIL lines, contractlint ok in both; `status --validate` VALID exit 0 (125 pre-existing archived-scope warnings, peer-owned).

### Per-AC citations

- AC-1 PASS: 13/13 from the true 5/13 shipped baseline at 53/55 lines; Gate rows byte-unchanged. The AC's "versus 7/13 today" described the docs/dev copy; the shipped baseline is lower, so the value property holds a fortiori.
- AC-2 PASS: direction is requested only by the ideation-owned `## Proposed approach` prompt; the backlog-owned prompts and `### Feedback Cycles` ask scope, decision proof, and correction rounds, never a direction.
- AC-3 PASS: `Estimate: {+NNN} net LOC across {M} files, tolerance {±NN%}.` plus separate `Semantics this may change:` line; the 53-line body piped to `status --new` in a scratch sd-b32 workflow exits 0, mints an id, preserves all eight headings and the ± byte, and validates.
- AC-4 PASS (retargeted per captain Finding 1): exactly the two shipped paths in the diff; no entity, gate record, or `*_test.go`; filing_readme_template_test.go untouched and green.
- AC-5 FAILED as written — flagged for Review-finding disposition, see finding below. The behavior half (digest divergence) remains proven by cycle-1's exercise on unchanged machinery; the prose half was reverted with docs/dev and has no shipped home.
- AC-6 PASS under the correction-round reading this dispatch directs: title equals the entity title from creation; body byte-equal to the reviewed file via exactly one rewrite through `gh api PATCH` (the mod's own repair path, not `gh pr edit`); baseline #699 (now MERGED) still carries the branch-derived title.
- AC-7 PASS: the mod diff is a single hunk inside `## Hook: merge`; startup and idle spans byte-unchanged; `gh pr edit` fail-closed re-observed live rather than trusted from cycle 1.

### Reviewer findings

- New finding (recommendation; no candidate change made): AC-5's deliverable half is homeless after the captain's wrong-target correction — the freeze-routing sentences were reverted with docs/dev and no shipped artifact carries a freeze rule to extend (grep over skills/ and mods/ is empty). Same class as the live-lane ruling, but the cycle-2 deviation declaration named only the live-lane ruling. Proposed: not material (field 3 is `none:` — AC-5 is not a value AC and no protected boundary is at risk); ownership the captain's local refit; disposition route for decision — captain ratifies the deferral or narrows AC-5 at this gate.
- Deferred risk, carried from cycle 1, now attached to the shipped mod's prose: the prescribed `gh stack view --json` membership read is clone-local (works only where `.git/gh-stack` exists). Promote to material if a ceremony must run from a fresh clone for an existing stack.
- No material findings; no polish findings.

### Summary

Recommend PASSED, contingent on the captain's gate-time ruling on AC-5: six of seven ACs hold re-exercised against e49154829, the deliverable is confined to the two shipped artifacts with docs/dev byte-clean, coverage is 13/13 at 53/55 from the true 5/13 baseline, the added mod prose is fully CI-agnostic with the six generic mechanics intact and the live-lane ruling correctly absent, and PR #711 carries the reviewed bytes through the mod's own REST-PATCH path. AC-5 fails as written because its only home left the approved surface with the captain's own Finding 1 — a rejection cycle cannot fix it (any fix is out of surface or captain-only), so it routes for decision alongside the already-declared live-lane deferral. One deferred risk carried forward with its promote condition.
