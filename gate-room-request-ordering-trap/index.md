---
id: z74qkypyq1zwhsjcnbjk1nk4
title: "Binding a briefing before its `request.json` exists permanently forecloses provider-room recording, and nothing documents the ordering"
status: backlog
source: "Independently hit in two repositories: four attempts wedged in spacedock-v1 on 2026-07-27, and two entities wedged in spacedock-subspace, both by following the order the command list implies."
started:
completed:
verdict:
score: 0.8
worktree:
issue:
---

Stop the obvious command order from silently and permanently disabling provider-room recording on a gate attempt.

## Problem

`gate record <entity> --briefing PATH/briefing.json` freezes a request digest at bind time, read from `request.json` **in the briefing's own directory**. If that file is absent, the attempt is bound with an empty frozen request digest. `internal/gates/operation.go:263` then refuses every later `--room` record:

```go
if attempt.Briefing.RequestDigest == "" || requestDigest != attempt.Briefing.RequestDigest {
    return fmt.Errorf("gate room request does not match the frozen request digest")
}
```

An empty digest can never be satisfied, so the attempt cannot accept a provider Result for the rest of its life. Adding `request.json` afterwards does not repair it — the binding is frozen and rebinding is refused.

This is confirmed by shipped behaviour, not inference. `TestGateRoomRejectsBriefingOnlyAttemptWithoutFrozenRequest` removes `request.json`, binds with `--briefing`, and asserts the bind **must exit 0**; then adding the file and re-recording is refused with `binding is frozen`, and `--room` is refused with `frozen request digest`.

So briefing-first is a fully supported path that quietly forecloses a capability, and it is the order the command list implies:

```
spacedock gate record <entity> --briefing PATH/briefing.json
spacedock gate record <entity> --room PATH
```

Nothing in `gate --help`, the lifecycle skill, or the contract states that `request.json` must exist before the first line runs.

## Evidence: two repositories, independently

- **spacedock-v1, 2026-07-27.** Four attempts bound briefing-first and are permanently chat-only: `79` backlog and ideation, `cn` backlog and ideation. The First Officer discovered the constraint by reading `operation.go` after a float failed, not from any documentation.
- **spacedock-subspace.** `subspace-r-one-question-synthesis` is blocked this way, and `one-question-in-a-review`'s ideation attempt was blocked identically, "reproduced from scratch by following the documented order."

Two First Officers, two repositories, same wedge, neither warned.

## Adjacent frictions from the same command family

Recorded here because they were found in the same sessions and share the surface. They are not this task's deliverable unless ideation folds them in deliberately.

1. **`digest-domain: canonical-bytes` is undocumented.** It is compact JSON with sorted keys, not the file's `shasum`. The mismatch error names neither the expected value nor the domain, so the only route is guessing serialisations.
2. **`--room` requires a `request.json` with no published schema.** Both sessions reverse-engineered the shape from a sibling entity's on-disk artifacts.
3. **A label-only `Reference` is accepted by `gate record --briefing` and rejected by the canonical loader** (`Reference "…" requires uri`). Same bytes, two verdicts — and the digest froze the version the presenter cannot load, which is how one subspace attempt became unpresentable.
4. **`request.json` requires `actor` and `approver` and requires both to equal `person:captain`.** The subspace report notes this contradicts `internal/reviewv1/mode.go`, which states the mode is never inferred by comparing two identities. Worth checking before filing as one defect: the `subspace:r` skill *does* describe mode-by-identity-comparison ("equal actor/approver authority retains Review v1's minimal binding Result; distinct authority carries the advisory wrapper"), so the contradiction may be between the skill text and the implementation rather than in the recorder alone.

## Asks, from the subspace proposal

Document the ordering, or compute the request digest lazily so briefing-first works; name the digest domain and the expected value in the mismatch error; publish the request schema; make the two `Reference` validators agree; and reduce the request to one declared authority rather than a pair.

Lazy computation is the option worth evaluating first, because it removes the trap rather than warning about it — a documented ordering constraint still fails closed for anyone who does not read the document, which is the population that hit this.

## Why this entity exists separately

`gate-legibility-summary-versus-excerpt` (`t2`) retains the originating proposal but explicitly scopes this out, and `gate-agent-ergonomics` (`sk`) does not name it. It was therefore described only in an attachment to a task that disclaims it, while continuing to wedge attempts in both repositories. This entity owns it.

## Out of scope

- Gate presentation legibility, which is `t2`'s.
- The provider package's location and the fact that `gate record --room` never sees a Subspace Result, which is a separate cross-repo gap.
- Materialising `git-root://` artifacts for provider presentation, which has its own entity.

## Acceptance criteria

Ideation fills these in. At minimum, the briefing-first order must either work or refuse at bind time with a message naming the consequence — a silently foreclosed capability must not be reachable by following the documented commands, with a falsifier that restoring the silent path turns the leg red.

## Test plan

Ideation fills this in. The existing gate-room fixtures in `internal/cli/gate_test.go` are the substrate.
