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

## A working fix, demonstrated end to end (2026-07-27)

Materialisation resolves this without giving up provenance, and it has now been run twice through a real review rather than proposed. It is what `rq`'s title already says — *materialize* git-root sources for provider presentation — so this section supplies the recipe, not a competing direction.

The recipe, and it is three steps:

1. Freeze a copy of the git-root source into the gate room beside `briefing.json` (`review/<stage>/briefing-<n>/design.md`).
2. Point the artifact `uri` at that copy as a **plain relative filename** — `design.md`, not a `git-root://` URI. The provider joins relative paths onto the briefing's own directory, which is exactly what broke the original attempt and exactly what makes this work.
3. Keep `rev` as the `sha256` of the frozen copy. The Briefing stays digest-pinned.

Verified twice on `subspace-tui 0.10.0-beta.6` via the `review-zellij` entry, same terminal, same entry, only the artifact changed: the git-root form exits 2 at resolution before presentation, the materialised form renders and returns a validated Result. Both round trips produced captain-rendered resolutions against `gate:docs-dev:cn:ideation` — a `revise` carrying two annotations, then an `approve` — recorded through `gate record --decision --actor person:captain`.

**Provenance is preserved, not traded away.** The earlier framing in this entity treated a filesystem path as the only alternative and called it a real loss: a working-tree path is neither immutable nor valid from another checkout. A frozen room copy is neither of those things. It is committed with the entity's own state commit, digest-pinned in the Briefing, and travels with the room through archive — the same durability the `git-root://` pin was reaching for, reached by copy rather than by reference. The cost is duplicated bytes in the room, which is the retention `9t` (minimum recoverable gate-room retention) exists to bound.

**What this does NOT fix, and must not be mistaken for fixed.** The provider still writes its Result into its own allocated scratch package, not into the gate room, so `gate record --room` never sees it. Recording still goes through the chat decision form. That is a separate gap (`krd`), untouched by materialisation.

## The remaining decision, which is still `rq`'s

Materialisation is proven to work; whether it is the *right* answer is not this entity's call. The alternatives remain teaching the provider the `git-root://` scheme, or changing what the recorder pins. Materialisation trades room bytes and a copy step for zero provider change; a scheme resolver trades provider work for no duplication. This entity supplies the mechanism, the failure diagnosis, and a demonstrated working recipe — `rq` chooses.

## Out of scope

- Choosing the resolution strategy — that is `rq`'s.
- The gate-room `request.json` scaffold and the `/subspace:r gate <room>` entry, which are separate unshipped surfaces.
- The binding-approval affordance hazard (`Tab`+`q` reaching a binding approve), recorded in the subspace repo's own shaping record.

## Acceptance criteria

Ideation fills these in, if this is not folded into `rq` first. The end state is that a Briefing pinning a `git-root://` artifact either presents its artifact bytes in package mode, or is refused with an error naming the unsupported scheme rather than a misleading missing-file path.

## Test plan

Ideation fills this in. A fixture Briefing pinning one `git-root://` artifact, presented through the package-mode entry, is the whole substrate.
