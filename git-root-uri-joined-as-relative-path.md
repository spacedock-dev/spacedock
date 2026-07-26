---
id: x9aed1cdyhpwxpgesa55jypm
title: "Subspace package mode joins a `git-root://` artifact URI onto the briefing directory instead of resolving it"
status: backlog
source: "Reproduced 2026-07-26 floating the 79 ideation gate briefing in Review v1 package mode; exact stderr and retained provider diagnostics below. Evidence for `rq` (git-root-review-v1-materialization) — fold in or close as duplicate at that owner's discretion."
started:
completed:
verdict:
score: 0.5
worktree:
issue:
---

Record the exact mechanism by which a recorder-valid `git-root://` artifact fails to present in Subspace Review v1 package mode, so the fix is designed against a path-join defect rather than an unspecified "not renderable".

## Problem

`rq` records that "recorder-valid `git-root://` sources are not renderable by current Subspace package mode". That is a symptom. The mechanism, reproduced deliberately, is narrower and more actionable: the provider has no `git-root://` scheme resolver, so it treats the URI as a **relative filesystem path and joins it onto the briefing file's own directory**.

Observed on `subspace-tui 0.10.0-beta.6` via the `review-zellij` fixed entry, presenting a canonical Briefing whose sole artifact was pinned as `git-root://state/<sha>/entity-session-claim-lease.md`:

```
artifact "artifact:entity-session-claim-lease-skeleton":
  resolve "git-root://state/7db811a7d5a45c03237ebf766166e60f943777d6/entity-session-claim-lease.md":
  open <room>/git-root:/state/7db811a7d5a45c03237ebf766166e60f943777d6/entity-session-claim-lease.md:
  no such file or directory
```

Note the collapsed `git-root:/` in the attempted path — the scheme survives into the filename, confirming a join rather than a failed lookup.

## What this rules in and out

- The capability handshake is NOT the problem: `capability.exit` was `0`, so `review-v1-provider-package-v1` is present and the child launched with the correct argv.
- The failure is per-artifact and occurs during resolution, before any presentation. The reviewer sees nothing; the entry exits 2 and retains a recoverable provider package.
- Nothing was recorded and no entity state was written, so the failure is safe — but it is safe by fail-closed resolution, not by design intent.

## Why the obvious workaround is a real decision, not a fix

Swapping the artifact URI to a filesystem path makes the float work immediately, and trades away exactly what the recorder depends on: `git-root://state/<sha>/<path>` pins an immutable commit, so a Briefing stays byte-verifiable after the entity moves or archives. A working-tree path is neither immutable nor valid from another checkout. Choosing between materialising git-root sources for the provider, teaching the provider the scheme, or changing what the recorder pins is `rq`'s decision — this entity only supplies the mechanism it should be designed against.

## Out of scope

- Choosing the resolution strategy — that is `rq`'s.
- The gate-room `request.json` scaffold and the `/subspace:r gate <room>` entry, which are separate unshipped surfaces.
- The binding-approval affordance hazard (`Tab`+`q` reaching a binding approve), recorded in the subspace repo's own shaping record.

## Acceptance criteria

Ideation fills these in, if this is not folded into `rq` first. The end state is that a Briefing pinning a `git-root://` artifact either presents its artifact bytes in package mode, or is refused with an error naming the unsupported scheme rather than a misleading missing-file path.

## Test plan

Ideation fills this in. A fixture Briefing pinning one `git-root://` artifact, presented through the package-mode entry, is the whole substrate.
