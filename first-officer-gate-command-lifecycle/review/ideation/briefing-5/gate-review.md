# IDEATION GATE — First Officer recorded gate lifecycle (`6y`)

Recommendation: **APPROVE cycle 13 and resume bounded implementation.**

## Capability and change

Keep the three authority mutations and their commit barriers. After consume, route by target kind: nonterminal targets dispatch; terminal targets enter the existing merge ceremony with no successor dispatch. Probe `gate --help` once immediately before each lifecycle and remove the session executable cache.

## Test and evidence

- The existing xb record is consumed at terminal `done`, has no dispatchable successor, and remains merge-blocked on open PR #564.
- Existing real lifecycle/resume and merge-sentinel tests pass at `d55ddca3`.
- The reset reuses the current lifecycle, Pi, merge-guard, and contract fixtures; intended change is nine existing files, +75/-99, with +88 additions as the hard stop.
- Pi evidence becomes root-session `assistant` text only; nested and nonassistant rows must fail.

## Reviewed snapshot

Cycle-13 design in `first-officer-gate-command-lifecycle/index.md`, including canonical AC-1–AC-8, obligation delta, test plan, and exact expected surface.

## Findings

Material findings are resolved in design: terminal approvals were ambiguously sent to dispatch; the PR hook excluded terminal in-flight rows; cache tests proved a helper rather than FO behavior; Pi parsing admitted the wrong authority. Exact-child/artifact transport attribution remains outside the approved AC-8 boundary.

## Recommendation

Approve the net-deleting reset. Reject if implementation adds another lifecycle, harness, runtime protocol, compatibility path, tenth file, or changes recorder/application semantics.

## Decision ask

Approve to resume implementation in the existing 6y worktree, or revise/hold with a concrete boundary.
