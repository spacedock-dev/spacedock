---
id: tn1f2qr6kfq8mhp5wyat1v63
title: Contract files grow unbounded — add a per-file token-budget guard (the forcing function)
status: backlog
source: "captain (2026-06-04) — token-efficiency analysis: first-officer-shared-core.md (~9,730 tok) is ~6x the median superpowers skill (~1,600 tok) and ~1/3 of the entire 14-skill superpowers library, because nothing caps contract-file growth. The roadmap noted 'started thin, grew dense'; the comparison quantifies the drift. Without a forcing function, anything we slim re-accretes."
score: "0.32"
worktree:
started:
completed:
verdict:
issue:
---

The prose contract grows monotonically: every session adds nuance, nothing caps it. `first-officer-shared-core.md` reached ~9,730 tokens — far past the ~650–5,775 tok range superpowers holds each skill to. The two decompositions this session (zd, ep) slimmed specific files, but without a size budget they (and the rest) will re-accrete. This entity adds the missing forcing function.

## Problem

There is no automated cap on contract-file size. Slimming is a one-time manual effort that erodes. A code gate — the workflow's own "prefer a code gate over a prose-only rule" principle — would make the budget self-enforcing.

## Proposed approach (ideation firms)

A guard test (mirror the existing `TestNoTimeoutLiteralExceeds60s` budget-guard pattern) asserting no single operating-contract file exceeds a token budget `N`. Candidate `N` ~3k tokens (the superpowers ceiling, excluding the outlier `writing-skills`); the exact value + which files are in scope (FO + ensign shared cores + runtime adapters + the spacedock-owned skills) is an ideation decision. The guard must (a) name the offending file + its measured count on failure, (b) use a stable token estimate (e.g. a documented chars/N proxy or a real tokenizer), (c) fail loudly so a future PR that re-bloats a contract file is caught at CI.

## Out of scope

- The actual decomposition of any file (that is `fo-shared-core-decomposition` and the binary-command roadmap moves) — this entity only installs the cap.
- Imposing the budget on non-contract docs (READMEs, proposals).

## Acceptance criteria (seed)

- **AC-1 (seed):** A guard test fails when any in-scope contract file exceeds `N` tokens and passes otherwise — verified by the test going RED on a synthetic over-budget fixture (or a temporary appended block) and GREEN on the real tree once any current over-budget file is either decomposed or explicitly grandfathered with a documented exception list.
- **AC-2 (seed):** The failure message names the file and its measured token count (not just a boolean), so the forcing function is actionable.

## Test plan (seed)

- Go test in the package the existing budget guard lives in; fixture/synthetic over-budget input proves it bites. Decide grandfathering policy for currently-over-budget files at ideation (a shrinking allowlist is the honest interim if first-officer-shared-core is not yet decomposed).

## Notes

- Sibling to `fo-shared-core-decomposition` (the first file this guard would flag) and the binary-simplification roadmap (`docs/dev/_proposals/binary-simplification-roadmap.md`, refreshed 2026-06-04). The roadmap's refresh names this as the top structural move — install the cap first, then slim.
