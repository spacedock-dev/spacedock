---
id: 1d9k9a6eg54eg9k7kdvn2n1b
title: "`status --resolve` reports an archived entity as unknown, contradicting the archived-inclusive ID space"
status: backlog
source: "Caused two wrong inferences reported to the captain on 2026-07-26/27: a failed lookup was read as 'this entity does not exist' when it meant 'not in active scope'."
started:
completed:
verdict:
score: 0.5
worktree:
issue:
---

Make a reference lookup distinguish "no such entity" from "archived", so a failed resolve cannot be read as nonexistence.

## Problem

`spacedock status --workflow-dir docs/dev --resolve 02av` returns:

```
Error: unknown reference: 02av
```

The entity exists. It is `_archive/ensign-finding-triage-disposition`, carrying `id: 02avdajaz0q3hnjwycm5fq45`. The prefix is valid and unambiguous; it is simply archived.

This contradicts the ID model the write core defines. Identity for `sd-b32` is "the 24-char SD-B32 … Status output displays the shortest unique prefix **across active plus archived** for the ID column; collisions lengthen only affected entities." Prefixes are computed against an archived-inclusive space, so a prefix that is unique because of an archived sibling is only meaningful if archived entities are resolvable. Today the ID space includes the archive and resolution does not.

`--short-id` shares the surface and should be checked for the same behaviour.

## Why it matters more than it looks

`--resolve` is the sanctioned way to check a reference before acting on it. The write core says so explicitly: "use `status --resolve` before mutating any reference that came from a human or older transcript." A stale reference from a transcript, a roadmap document, or a debrief is *exactly* the case where the entity is most likely to have been archived in the meantime — so the tool is least reliable in the situation it exists for.

The failure mode is silent and reads as authoritative. Twice in one session a First Officer took `unknown reference` as evidence that a task had been dissolved:

- The `durable-decisions` roadmap package names members `3k`, `vn`, `h1`, `xb`, `02av`; none resolved, and the FO reported to the captain that the package was stale. The members had been archived. The captain corrected it.
- A landed product spec references `owner: 02av`. The FO reported it as "a pointer to a dissolved task nobody can look up". It resolves to an archived entity; the finding was real but its characterisation was wrong, and it reached a filed entity before being corrected.

An error that reliably produces a confident wrong conclusion is worse than one that produces a visible failure.

## What a fix needs to decide

The output format already models scope — a successful resolve prints `scope=active` — so the concept exists and only the search does not cover it. Options, and this is the design question:

- Resolve archived entities and report `scope=archived`, letting the caller decide. Consistent with the ID space and with the existing output field, but it means a mutating caller can now resolve something it must not mutate — so mutation guards must not assume a successful resolve implies an active target.
- Keep refusing, but with a message that distinguishes archived from absent, so the caller learns the truth without gaining a mutable handle.

Whichever is chosen, the archived-inclusive prefix rule and the resolution scope should agree, and the disagreement itself should be impossible to reintroduce.

## Out of scope

- Making archived entities mutable. `status` and `merge guard` correctly refuse them; this is about lookup, not authority.
- Changing how prefixes are computed.

## Acceptance criteria

Ideation fills these in. At minimum, a reference to an archived entity must produce an outcome a caller can distinguish from a reference to no entity at all, with a falsifier that collapsing the two back into one message turns the leg red.

## Test plan

Ideation fills this in. Fixtures with an archived-only slug, an active-only slug, a genuinely absent reference, and a prefix made unique only by an archived sibling.
