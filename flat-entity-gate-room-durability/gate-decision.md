# Backlog gate: make flat-ticket review rooms durable and round-recordable

## Recommendation

Approve ideation. Flat tickets are a supported entity shape, but their sibling review rooms are not handled consistently by state commits or advisory round recording. The resulting dead end currently blocks `wy` after a valid correction.

## Outcome

- A flat ticket can retain and publish its sibling review room through supported commands.
- Advisory correction rounds can be recorded for a flat ticket without moving it to folder form.
- Existing closed-gate `room-ref` values keep their current meaning and remain readable.
- A fresh clone can retrieve and validate the exact room named by the ticket.

## Scope boundaries

- Do not migrate flat tickets to folder form.
- Do not rewrite frozen historical gate bindings.
- Do not change Artifact or Reference `git-root` identity semantics.
- Do not redesign room retention or request-digest ordering.

## Proof owed at ideation

Ideation must select one shared entity-review-home rule for flat and folder shapes, identify every recorder/commit/validation consumer that needs it, estimate net LOC and files with tolerance, and define two-clone tests that fail if the sibling room is absent, a historical `room-ref` changes meaning, or a flat correction round cannot be recorded and replayed.
