---
id: s6j7mn45sf8d62e4xz181s93
title: "`gate consume` writes the successor status but not `started`, so an entity enters its first working stage with no start time"
status: backlog
source: "Patched by hand twice on 2026-07-27 driving 79 and cn through their backlog gates; both entered ideation with `started` empty until the First Officer noticed and stamped it."
started:
completed:
verdict:
score: 0.4
worktree:
issue:
---

Make the verb that performs a status transition own the lifecycle timestamp that transition implies, so an entity's start time does not depend on the First Officer remembering to add it.

## Problem

The dispatch core's transition step writes status and the start stamp together:

```
spacedock status --workflow-dir {workflow_dir} --set {slug} status={next_stage} worktree=… started
```

`gate consume` took over that transition for gated stages — it "atomically writes successor status plus consumed state" — but replicates only the status half. The entity lands in its first working stage with `started:` still empty.

Observed twice on 2026-07-27, both at a backlog gate:

- `79` (`entity-session-claim-lease`) — consumed to `ideation`, `started` empty, stamped by hand.
- `cn` (`version-output-runtime-and-sandbox-state`) — same shape, same manual repair.

`started` is entity-level and idempotent (`status --set … started` auto-fills and is skipped when already set), so the omission only bites on the *first* consumed transition into working. That is precisely the transition the field exists to record, which is why the miss is easy to overlook and lands wrong data rather than no data.

## The inconsistency that makes this a defect rather than a preference

The two verbs that own status transitions disagree about whether they own lifecycle timestamps.

- `merge guard` terminalizes with `status` + `verdict=passed` + `completed` **in one `--set`**, per the `pr-merge` mod. It owns its timestamp.
- `gate consume` writes successor status alone. It does not.

So `completed` is verb-owned and `started` is First-Officer-owned, with no stated reason for the split. `internal/status/mutate.go:19` treats them as one class — `timestampFields = map[string]bool{"started": true, "completed": true}` — so the storage layer already regards them as symmetric while the transition layer does not.

## Cost

Small but silent, and it corrupts rather than omits. Any consumer computing a duration from `started` gets a wrong answer rather than a missing one when the field is blank but the entity has plainly been working. The First Officer has no prompt to notice: `consume` reports success, `status` shows the new stage, and nothing surfaces the empty field.

## Relation to `gqs`

`dispatch-entered-stage-after-gate-consume` (`gqs`, at validation) is adjacent and distinct. It concerns whether the First Officer *dispatches* the stage it has entered after consuming; this concerns whether the transition stamps its timestamp. They share the post-consume boundary and could be fixed together, but neither subsumes the other, and `gqs` is a `durable-decisions` member while this is unlabelled.

## What a fix needs to decide

- Whether `consume` should stamp `started` itself, or whether the post-consume dispatch step should, with the contract naming the owner either way. Both are defensible; the current state is that neither does and it works by FO memory.
- Whether the symmetry with `merge guard`'s `completed` should be made explicit in the contract, so the next transition verb inherits a rule rather than a precedent.
- Whether an entity at a non-initial stage with an empty `started` should be a validation warning. That would catch the class rather than this instance, but it is a new standing check and needs the captain's approval on those grounds.

## Out of scope

- Dispatching the entered stage, which is `gqs`.
- `completed` and `verdict` at terminalization, which `merge guard` already owns correctly.

## Acceptance criteria

Ideation fills these in. At minimum: an entity consumed from its initial gate into a working stage carries a `started` timestamp without a separate First Officer command, with a falsifier that removing the stamp leaves the field empty and turns the leg red.

## Test plan

Ideation fills this in. Existing gate consume fixtures are the substrate; this needs no live lane.
