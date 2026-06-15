---
title: "Team-mode FO non-deterministically omits `verdict:` during the multi-turn team finalize under headless `-p`"
status: backlog
source: "Surfaced by 7e (headless-dispatch-mode-intent) cycle-2 live runs, 2026-06-15: -count=3 TestLiveEnsignCycle on sonnet ran twice -> both PASS/FAIL/PASS, the lone FAIL always the TEAM run failing the verdict: end-state. Team FO sets status:done, archives, emits TERMINAL_TEARDOWN_BOUNDED, but ~1/3 omits the verdict: value; bare reliably writes it. Confirmed a real omission, not a snapshot race (barrier required verdict -> team run sat 3m silent, still none)."
sprint:
id: reeppr990pyzzaejmbnyrvt7
---

A headless `-p` FO driving in **team mode** non-deterministically omits the `verdict:` frontmatter value when it terminalizes an entity, even though it sets `status: done`, archives the entity, lands the path-scoped commit, and emits the `TERMINAL_TEARDOWN_BOUNDED` teardown marker. In **bare mode** the verdict is written reliably. So a terminalized-and-archived entity can end up with no recorded PASS/REJECT verdict — an incomplete terminal record.

## Problem

`verdict:` is the recorded PASS/REJECT outcome of a terminal entity. An archived entity with `status:done` but no `verdict:` is a data-integrity gap in the terminal record. It also means `verdict:` is NOT a team-independent end-state fact — which is why 7e's team-agnostic live smoke gates on the genuinely mode-invariant facts (`status:done` + path-scoped commit + stage-report shape) and this omission is tracked here instead of papered over.

## Evidence (2026-06-15, surfaced by 7e cycle-2)

- `-count=3 TestLiveEnsignCycle` on sonnet, run TWICE: both PASS/FAIL/PASS; both times the lone FAIL was the TEAM run, failing the `verdict:` end-state check. Both BARE runs each wrote a verdict and passed.
- NOT a snapshot/reaping race: 7e's implementation strengthened the test barrier to require `verdict:` before snapshot+kill (worktree commit 9a2b5290); the team run then sat at the 3-minute quiet budget with still no verdict — the FO never writes it, given full time.
- The forced-team `TestLiveEnsignCycleTeamTeardown` happened to write a verdict in its observed run — luck, not reliability; team-mode verdict-write is non-deterministic, not uniformly absent.

## Proposed approach (ideation to confirm)

Hypothesis: the FO's terminalize/merge ceremony spans several turns; in team mode under the upstream Claude Code teardown fragility (`-p`), the `verdict:` write (part of the merge/finalize) is sometimes not reached before the cycle ends, while `status:done` + the path-scoped commit land earlier. Likely fix territory: make the terminalize ceremony atomic w.r.t. `status:done` + `verdict:` (write the verdict no later than the status/commit), or order the verdict write before the teardown marker. This is FO terminalize/merge-ceremony behavior (contract/dispatch), NOT live-test-harness territory.

## Acceptance criteria (sketch — ideation fleshes out, externally proven)

- A team-mode headless `-p` FO that terminalizes an entity writes a `verdict:` value deterministically across repeated runs (proven by live or harness runs showing `verdict:` set on EVERY team-mode terminalize).
- The fix lands in the FO terminalize/merge ceremony, not the smoke's assertion shape.

## Out of scope

- The team-vs-bare DISPATCH-mode determination (7e — proven done).
- The live-smoke assertion shape (7e cycle-2 gates on mode-invariant facts).
- The contract prose changed by 7e (d3f0196d) — validated-clean and unrelated.

## Notes

Fast-follow from 7e (`headless-dispatch-mode-intent`). Same team-mode-under-`-p` fragility theme as 7e's determination (7e body line ~75: team is fragile under `-p`; bare is the robustness choice). Captain may pull into a sprint or keep as fast-follow.
