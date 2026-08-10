---
commissioned-by: spacedock@template
entity-type: artifact
entity-label: artifact
entity-label-plural: artifacts
id-style: sequential
stages:
  defaults:
    worktree: false
    concurrency: 2
  states:
    - name: draft
      initial: true
    - name: review
      gate: true
      feedback-to: draft
    - name: polish
    - name: done
      terminal: true
---

# Refinement Workflow Template

Track an artifact through rounds of improvement to a locked, shipped result. This is the universal base shape: an entity is drafted, reviewed for whether it is ready, polished if accepted, and marked done when finished. No layers active.

Use this template when the captain's mission is to track an artifact through rounds of improvement with a human-in-the-loop quality bar — design docs, PRDs, content pieces, outreach replies, integration records, anything where the work is "make this thing good enough." Variants in the Adoption section adapt the stage list and entity body to common end-use shapes (outreach, integration, content production, PRD authoring) without changing the underlying structure.

## File Naming

Each artifact lives as either:

- a flat markdown file `{slug}.md` (default — use this unless the artifact produces many side files), or
- a folder `{slug}/` containing `index.md` as the canonical entity file, when the artifact produces per-stage attachments (drafts, reviewer notes, transcripts, output files) that belong alongside the tracker.

Slugs are lowercase, hyphens, no spaces. Example: `q3-launch-narrative.md` or `q3-launch-narrative/index.md`.

## Schema

Every artifact file has YAML frontmatter. Fields are documented below; see **Artifact Template** for a copy-paste starter.

### Field Reference

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique identifier, format determined by id-style in README frontmatter |
| `title` | string | Human-readable artifact name |
| `status` | enum | One of: draft, review, polish, done |
| `source` | string | Where this artifact came from |
| `started` | ISO 8601 | When active work began |
| `completed` | ISO 8601 | When the artifact reached terminal status |
| `verdict` | enum | PASSED or REJECTED — set at final stage |
| `score` | number | Priority score, 0.0–1.0 (optional) |
| `worktree` | string | Worktree path while a dispatched agent is active, empty otherwise. Once set on first dispatch into a `worktree: true` stage, it stays set across all non-terminal advancements (stickiness) and clears at terminal merge. |
| `issue` | string | GitHub issue reference (optional cross-reference) |
| `pr` | string | GitHub PR reference (set when a PR is created) |

## Stages

### `draft`

The artifact is produced or revised — every time it enters the loop, including after a review bounce. The draft integrates prior reviewer notes and is complete enough to be evaluated end-to-end.

### `review`

A reviewer reads the whole draft and makes a clear accept/reject decision behind an approval gate: gate-approval to `polish`, or rejection back to `draft` with specific, actionable notes.

- **Gate content:** Show the review decision, non-empty actionable findings, and whether approval advances to polish or rejection returns to draft.

### `polish`

Final cleanup before the artifact is locked — formatting, copy edits, last-pass consistency. Cosmetic only; preserves the substance the reviewer accepted and reopens nothing structural.

### `done`

Terminal state: the artifact is locked and shipped (or filed, per the variant). `completed` set, `verdict: PASSED`, archived. Start a new entity rather than reopening a done one.

## Workflow-specific rules

The FO/ensign operating contract already governs generic stage semantics and proof discipline. Refinement is the universal base shape, so most of its discipline lives in the contract; the rules below add only the refinement-shape specifics.

- **Human-in-the-loop quality bar.** `review` is an approval gate a human reviewer owns — the reviewer reads the whole draft and makes an explicit accept/reject call with concrete notes, never a vague "this needs work." `polish` is cosmetic-only; new content that should have gone through review does not belong there.
- **No layers by default.** The base shape touches no repo and waits on no external event, so no structural layers fire. Variants that need them (e.g. outreach's `watching` stage) activate the layer through the variant, not the base.
- **Variant menu.** Common end-use shapes are refinement with adjusted stages and a different entity body — `outreach`, `integration`, `content-production`, `prd-authoring` (see `## Adoption` → Surface variants). A variant changes the stage list and snippet, never the underlying draft → review → ship structure.

## Workflow State

View the workflow overview:

```bash
spacedock status --workflow-dir {dir}
```

Output columns: ID, SLUG, STATUS, TITLE, SCORE, SOURCE.

Find dispatchable artifacts ready for their next stage:

```bash
spacedock status --workflow-dir {dir} --next
```

## Artifact Template

```yaml
---
id:
title: Artifact name here
status: draft
source:
started:
completed:
verdict:
score:
worktree:
issue:
pr:
---

Brief description of this artifact and what it aims to achieve.

## Draft

The current draft of the artifact lives here.

## Review notes

Reviewer notes accumulated across review rounds.

## Final

The locked, polished version (filled in at polish or done).
```

## Commit Discipline

- Commit status changes at dispatch and merge boundaries
- Commit artifact body updates when substantive

## Adoption

### Pre-fill stages

```yaml
- name: draft
  initial: true
- name: review
  gate: true
  feedback-to: draft
- name: polish
- name: done
  terminal: true
```

The default stage list above. Variants below adjust this list.

### Apply layers

None by default. The base refinement shape uses no structural layers; entities never touch the repo and never sit parked waiting on external events.

### Offer mods

None by default. Variants that activate the parked-stages layer (e.g., outreach with a `watching` stage) trigger the silence-watcher offer through the layer mechanism in the commission skill, not through this template.

### Inject entity-template snippet

Use the refinement snippet (draft / review notes / final) shown in the Artifact Template section above. Variants override this with their own snippet (see below).

### Surface variants

Refinement is the universal base — many common workflow shapes are refinement with adjusted stages and a different entity body. When trait detection lands on `refinement` and the cues match a variant below, surface it as a confirmation:

> Looks like an outreach pipeline. I'd suggest these stages instead of the default `draft / review / polish / done`. Want me to use this variant?

**outreach** — captain mentions contacts, leads, sending, followups, replies, drip, pipeline

- Stages: `research (initial)` → `draft (gate)` → `sent` → `watching (parked, gate)` → `followup (feedback-to: watching)` → `closed (terminal)`
- Layers: parked-stages (on `watching`)
- Mods: silence-watcher offered for `watching` (timeout/nudge semantics)
- Entity-template snippet: contact / message draft / sent-at / response / outcome

**integration** — captain mentions sync, ingest, enrichment, external system, record-by-record processing

- Stages: `intake (initial)` → `enrichment` → `sync (gate)` → `archived (terminal)`
- Layers: none (or repo-mutation if the sync target is a repo file)
- Entity-template snippet: incoming record / enrichment notes / sync target / sync result

**content-production** — captain mentions blog posts, articles, videos, publishing, editing, copy

- Stages: `drafting (initial)` → `editing (gate, feedback-to: drafting)` → `polish` → `shipping (gate)` → `published (terminal)`
- Layers: none (or parked-stages if `shipping` waits on an external publish window)
- Entity-template snippet: artifact draft / editor notes / final / publish target / published-at

**prd-authoring** — captain mentions design doc, PRD, RFC, spec, locked, approved

- Stages: `draft (initial)` → `review (gate, feedback-to: draft)` → `locked (terminal)`
- Layers: none
- Entity-template snippet: problem / proposal / open questions / decision (the final locked record)

If trait detection lands on `refinement` but the cues do not match any variant above, use the default stage list and the refinement snippet from the Artifact Template section.

**Variant per-stage detail is materialized at commission time, not pre-baked here.** Each variant entry above lists its stage shape, layer activation, mod offers, and entity-template snippet — but not per-stage `Inputs` / `Outputs` / `Good` / `Bad` bullets. Those are auto-generated by commission Phase 2b for the chosen variant's stages and then tightened by {captain} via the `review stages` flow before the first dispatch (see commission Phase 3 Step 1a). When a variant grows enough captain-specific detail to be worth pre-baking, specialize it into its own template file with a full `## Adoption` section and per-stage bullets — this template stays the universal-base shape that variants adapt from.

### Confirmation prose

Surface this in Phase 1 once the template is selected (substituting the chosen variant if one fired):

> I'll set this up as a **refinement** workflow{ — variant: {variant}}: each {entity_label} moves through `{stage_list}` until it is locked. No worktree stages and no PR/merge ritual — this workflow does not touch the repo.{ Parked-stages layer fires on `{parked_stages}` because the entity sits waiting on external response.{ I'll offer the silence-watcher mod when we get to mod offers.}}
