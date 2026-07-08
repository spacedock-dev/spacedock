---
id: 85z12f0ywkzy47akg9gwh6hm
title: "merge-guard arm phase has no keep-moving clause: armed reads as a stopping point"
status: backlog
source: "FO self-diagnosis, 2026-07-08 live session. After the captain approved three validation gates and said \"push it,\" the FO ran `spacedock merge guard <slug> --verdict passed` for each entity, which only ARMS the merge (sets mod-block=merge:pr-merge) — then stopped to read the pr-merge.md hook file instead of immediately constructing and presenting the PR draft in the same turn. Before finishing even one entity's draft, the FO got pulled into an unrelated task and the arm sat untouched. When the captain later asked \"what did you do when I said push it,\" the honest answer was: armed three entities, pushed nothing."
started:
completed:
verdict:
score:
worktree:
issue:
---

**Problem:** `skills/first-officer/references/fo-merge-core.md`'s `«merge.guard»` capability describes the phase-invocation mechanics ("invoke it directly per phase; its own stdout/stderr name the FO's next action") but never states that an "armed" result is not a stopping point. Contrast `fo-dispatch-core.md`, which is explicit for stage completions: "Implementation completion is not a stopping point... The FO does not park a completed implementation and wait." The merge ceremony is exactly as sequential and stateful as the dispatch stage sequence (arm → invoke hook → finalize), but only the dispatch side carries the "keep moving" clause. The asymmetry reads as an invitation to treat "armed" as a natural pause.

**Cause:** the FO's own behavior is only as good as what the contract prescribes at each state-machine transition. `fo-dispatch-core.md` earns correct keep-moving behavior at stage boundaries because it says so explicitly, in words, at the exact transition point. `fo-merge-core.md` never says the analogous thing at the arm→hook transition, so nothing in the contract text pushed the FO to continue past "armed" into constructing and presenting the PR draft in the same turn.

**Recommended fix:** add a line to `«merge.guard»`'s effect/done-when description mirroring the dispatch-core language: an "armed" result is not a stopping point — the FO proceeds to construct and present the PR draft (or invoke whatever hook action is named) in the same turn, not a later one, exactly as a completed non-gated dispatch stage routes immediately to the next action.
