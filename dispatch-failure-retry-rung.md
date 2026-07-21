---
title: A failed dispatch is retried once before anything halts — replace Degraded Mode with a per-entity retry
source: "captain (CL), 2026-07-21 — ruling: an API transport failure must not trigger degradation; nudge the worker for retry. Split out of q4 on the evidence of a two-seat adversarial scope review, which found the defect is the missing rung BELOW the trigger, not the trigger itself."
status: backlog
id: 9q4x5hvyxthc41mt73txr23k
---

Today the FO contract has no first-failure response at all. It goes straight from one observed dispatch failure to "any SECOND dispatch failure within the session" tripping Degraded Mode — a session-wide, irreversible fallback to sequential bare dispatch. There is no retry, no nudge, and no way back within the session.

## Problem

**Observed twice, both times against a healthy substrate.** On 2026-07-20 two 2ae ensigns died with `API Error: Connection closed mid-response` / `"error":"server_error"` (read from the subagent jsonl, not inferred). Degraded Mode tripped on the second and the session lost concurrent dispatch for its remaining life. The same error pair, ~2 seconds apart, is recorded at `_debriefs/2026-07-04-01.md` issue 2.

**The zombie premise the irreversibility rests on does not describe this runtime.** Measured on 2026-07-21, not argued: the killed 2ae ensign was nudged and RESUMED FROM ITS TRANSCRIPT with 138 turns / 620 KB of context intact and continued its work. Three agents the FO had recorded as "confirmed terminated" were still alive a day later and approved shutdown in ~1 second each; the roster then went 5 → 2 with `reconcile` drift empty. Dead workers are recoverable, the roster is accurate, and shutdown reconciles cleanly.

**BARE MODE IS NOT THE DEFECT AND MUST SURVIVE.** Bare mode — sequential blocking `Agent()`, no `team_name` — is the legitimate teams-unavailable path, live-proven to work by `e3z`. Only the mid-session irreversible TRIGGER is wrong. Any change here must keep bare mode reachable.

**The counter has no mechanism advantage.** `claude-fo-dispatch.md:85` states the trigger has "no time window, no counter — the FO tracks it by its own observation". The contract's stated rationale for bluntness is that an FO under context pressure cannot reliably classify failures — but the counter is itself an unverified FO judgement from the same store. It is a smaller judgement, not a mechanism.

## The open question that decides the design

**Was the failure caused by our own fan-out?** The 2026-07-04 debrief correlates the identical error pair with 6 concurrent live sessions and says root cause was never determined (resource/concurrency ceiling vs. provider-side blip). If the cause is a fan-out ceiling, then retrying at unchanged concurrency reproduces the failure and Degraded Mode's concurrency reduction was the accidentally-correct remedy for the wrong reason. **Resolve this before designing the rung.** Sibling `r5y6` (live-ci-api-error-log-capture) states in its own Problem that a stall "cannot be classified: transient weather vs a CLI retry/hang defect" even with full logs on disk — so any classification-dependent design is currently unevidenced. Land r5y6 first, or accept a design that never classifies.

## Proposed approach

{Ideation fills this in. Both review seats independently converged on a failure-kind-AGNOSTIC bounded re-dispatch: the retry IS the classifier, so the FO never parses error strings. Expected shape: a dispatch failure is Agent() returning an error OR a worker session ending with no completion signal and no stage report; re-dispatch that (entity, stage) ONCE; a second consecutive failure of the SAME entity/stage halts THAT entity and reports to the captain while other entities continue and the session dispatch mode never changes. The seats disagreed on whether the session-wide trigger survives as a backstop — that disagreement is the ideation question.}

## Out of scope

- Retiring bare mode, or any change that makes it unreachable.
- Introducing transport-vs-structural error vocabulary into the contract (the repo has zero occurrences of "transient" under `skills/`).

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - A dispatch failure is retried once before anything degrades, and the retry is bounded.**
Verified by: an offline stream fixture asserting that after one `Agent` error the next `Agent` call carries the incremented suffix and NO third attempt appears. Trigger (1) has never had any test in either direction — the replacement must not ship equally unverified.

**AC-2 (VALUE) - Bare mode remains reachable after the change.**
Verified by: a live no-team drive reaching bare dispatch, the same way `e3z` established its baseline — an independent check that can move the wrong way.

**AC-3 - No session-wide irreversible degradation fires on a single worker's death.**
Verified by: the fixture in AC-1 plus a review-time read of the surviving triggers; no committed prose-grep.

## Known co-edits and collisions (found by scope review, verified)

- `internal/contractlint/boot_resident_closure_test.go:53` pins the literal `## Degraded Mode` anchor on disk — mandatory same-commit co-edit.
- `internal/contractlint/dispatch_recovery_value_binding_test.go:35-58` binds the captain-report blockquote VERBATIM to Go. Deleting or editing that prose breaks it.
- **`-cycleN` suffix is already double-booked** — zombie replacement (`fo-dispatch-recovery/SKILL.md:72`) and feedback rework (`fo-dispatch-core.md:55`). A retry burning cycle 2 can trip `feedback-rejection-flow/SKILL.md:17`'s "on cycle 3, escalate to the human" a cycle early. A distinct `-retry` axis avoids it but contends with the 64-char name cap.
- **`claude-fo-dispatch.md:68` bans re-dispatching a replacement ensign** ("you have no evidence the first ensign failed"). A retry rung needs an explicit carve-out for the observed-failure case or the two rules contradict and the FO improvises — the exact behavior the `## Awaiting Completion` anti-pattern list exists to suppress.
- **Shutdown collision:** `fo-dispatch-core.md:55` demands `«worker.shutdown»` of the prior cohort before a `-cycleN` dispatch; `fo-dispatch-recovery/SKILL.md:70` forbids cooperative shutdown to a dead ensign. A transport-killed worker fires both. Getting it wrong yields two live ensigns on one worktree — a correctness failure, worse than the throughput failure degradation costs.
- **`internal/claudeteam/claudeteam.go:68-79` `BareModeAdvisory`** tells the FO to run `ToolSearch select:TeamCreate` on every bare dispatch lacking recent team evidence. After q4 retires TeamCreate that names a tool which no longer exists, on the legitimate bare path.
- **The retry bookkeeping needs a durable answer, not prose.** "Retry once" is bounded only if the FO reliably remembers it retried; otherwise it is an unbounded respawn loop against a dead API with no captain-visible event. This is the load-bearing design question.

## Sequencing against q4

**q4 deletes the repo's only existing retry rung.** `using-legacy-claude-team/SKILL.md:50` («legacy-team.recover» rung 1, "Attempt one new TeamCreate with a fresh name") is the sole try-once-before-degrading step in the contract, and q4 retires that file. If q4 lands first without this entity, the contract is left with strictly ZERO retry surface. Land together, or land this first.

**Byte contention.** q4's AC-4 re-tightens `foFunctionReferenceBaselineBytes` to measured+1 at zero slack. A scope seat measured this change at roughly +164 to +204 B net (fits today's 507 B headroom with ~300 spare), but ~+900 B if the retry bookkeeping gets its own skill section — which does NOT fit and would require q4's removal of `using-legacy-claude-team/SKILL.md` (14065 B) from the measured surface first. Budget the two jointly.

## Test plan

{Ideation fills this in. Note the coverage vacuum: `TestLiveDegradedBareRecovery` is `//go:build live` and appears in ZERO CI `-run` filters, and `.github/workflows/runtime-live-e2e.yml:111` sets `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1`, so the teams-unavailable branch never runs in CI. Any change here currently lands with no oracle in either direction.}

## Captain rulings (2026-07-21, recorded at shaping)

1. **Degraded Mode is retired entirely — no session-wide backstop survives.** This answers the open ideation question (the two review seats' disagreement). The three triggers dissolve, each to its honest shape: (1) a second consecutive failure of the SAME entity/stage halts THAT entity and reports to the captain — never the session; (2) `/spacedock bare` becomes a plain captain instruction to dispatch bare from that point, not an irreversible mode transition; (3) `Agent`/`SendMessage` unavailable is not a mode at all — it is the ordinary teams-unavailable condition selecting bare dispatch, evaluated where dispatch happens. Bare mode itself survives untouched (already a hard constraint above). Co-edits stand as mapped: `fo-dispatch-recovery`'s Degraded Mode section, the two contractlint anchors, and the recovery-report verbatim binding.
2. **Title de-minted** (the metaphor family is banned vocabulary): the slug stays stable per the no-slug-churn precedent; the title now reads plain.
3. **Recommended to the captain, pending a nod:** accept the failure-kind-agnostic design both seats converged on (the retry IS the classifier; no error-string vocabulary enters the contract), which removes the r5y6 sequencing dependency — the "was it our own fan-out" question stops mattering because a retry that reproduces the failure halts that entity either way.
