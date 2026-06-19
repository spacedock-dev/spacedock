---
title: Prose-function restructure — express the FO cores as «fn» invocations (shipped-verb backticks or hand-followed guillemets)
status: ideation
source: 0205 carve (2026-06-17, captain "stamp them") — index DoD candidate "prose-function-restructure"; 2y MERGED unblocks it.
score: 0.5
sprint: 0205-layered-fo
group: restructure
sprint-readiness: ready
id: czw7whmaqjkasx1sjbq0at4h
started: 2026-06-19T07:45:59Z
---

Express the boot-resident core + dispatch core + gate flow + merge flow as prose-function invocations («state.boot», «dispatch.next-action», «gate.assemble-verdict», «merge.guard», «feedback.route») with bodies that name a shipped verb (backticks) or carry the hand-followed recipe (guillemets). The migration substrate; every verb member flips one body guillemet→backtick. Verbs not shipping this sprint stay guillemets with a named target. Depends on 2y (the merge/dispatch cores it restructures, MERGED).

CRITICAL co-design: this re-touches `## Startup` / `## Completion and Gates` / `## FO Write Scope` — the SAME regions 6re's step-4 rewrite + 72's step-1.5 insert edit. The landing ORDER must be co-designed (captain-ratify): either 6re step-4 + 72 step-1.5 land first and the restructure re-expresses on top, or the restructure lands first and 6re/72 rebase onto it.

Ideation must resolve: a BEHAVIORAL (non-grep) AC — the proof is the live drive / the verb-body flip, NEVER a prose-grep over the restructured contract; the `«gate.assemble-verdict»` binarization-timing decision (ship the verdict verb this sprint, or keep prose-function with L3 escalation).
