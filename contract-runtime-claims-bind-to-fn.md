---
title: Core-level contract claims must bind to «fn» present/absent, not assert runtime behavior flat — audit + fix, consolidating g6
source: "post-sprint 0260 (2026-07-21). Two instances found where host-neutral core prose asserts a runtime-varying behavior as universal, contradicting the per-adapter «fn» binding; each was caught by an FO having to improvise in the field. Captain directed a consolidated audit."
status: backlog
sprint:
id: j8s43ffcvdv6367td5v7d96e
---

The contract binds runtime-varying capabilities as `«fn»`s: `fo-dispatch-core.md:83` — "Runtime adapters bind the capability `«fn»`s below in their `## Runtime implementation` blocks." Present/absent and async/blocking are resolved per adapter (`«addressable-worker»` :91, `«async-dispatch»` :96). But some CORE-LEVEL prose asserts a runtime behavior FLAT — as if universal — which contradicts what a given adapter actually binds. An FO reading the core trusts a guarantee its runtime does not provide, and either acts wrongly or has to improvise the gap by judgment.

## Two known instances (both caught in the field, not by a check)

1. **Standing-teammate idempotency (entity g6).** `fo-dispatch-core.md` describes standing-teammate injection as "Idempotent (already-alive members omitted)"; `claude-fo-dispatch.md` says the call "does NOT dedup … idempotency is your own-roster concern — do not re-inject a standing teammate you already spawned." Both cannot hold. Observed: the 0260 Commander FO received a spawn spec for an already-alive teammate and "had to decide by judgement," deferring to point-of-need — a call the contract should have made for it.

2. **Async/streaming fan-out.** The v4dm fan-out clause draft asserted a "streaming per-member verify" as the universal failure mode, but streaming only exists when `«addressable-worker»` is PRESENT / `«async-dispatch»` is ASYNC; in bare mode dispatch BLOCKS and reviews batch automatically. The clause WORDING is fixed in v4dm (streaming now bound to `«async-dispatch»`). The RESIDUAL, deeper gap this entity owns: a Claude FO in bare `claude -p` (`«addressable-worker»` ABSENT) authored an async streaming fan-out — "member 1's findings verify while member 27 is still being read" — a shape its own runtime cannot execute. The FO planned against a capability it does not have; nothing in the contract stopped it.

## Problem shape

Both are the same class: a core-level claim about runtime-varying behavior that a binding contradicts (instance 1) or that an FO applies without checking its own binding (instance 2). The `«fn»` mechanism exists precisely to localize runtime variance; these leak it back into flat core prose or into FO planning that ignores the binding.

## Proposed approach (ideation)

1. **Sweep** the FO/ensign core (`first-officer-shared-core.md`, `fo-dispatch-core.md`, `fo-write-core.md`, `fo-merge-core.md`, `ensign-shared-core.md`) for claims asserting a runtime-varying behavior (async vs blocking, present vs absent, dedup/idempotency, back-channel, streaming, mid-run steering) FLAT rather than through a `«fn»`. Enumerate each with a verdict: already-bound / should-bind-to-«fn» / should-move-to-adapter / genuinely-host-neutral.
2. **The standing guard is a contractlint CONTAINMENT check (captain directive 2026-07-21): the shared / host-neutral contract files must contain NO specific-runtime-host mention.** No `Claude` / `Codex` / `Pi` host name and no host-specific tool token may appear in `first-officer-shared-core.md`, `fo-dispatch-core.md`, `fo-write-core.md`, `fo-merge-core.md`, `ensign-shared-core.md`; host-specifics live ONLY in the adapter `## Runtime implementation` blocks. This is NOT a banned prose-grep: it asserts an existence/containment fact (WHERE a host name may appear), which the captain's own grep ruling permits ("presence or absence is an existence fact; a grep establishes it soundly when that fact is itself the claim") and which existing contractlint already does for token containment (`internal/contractlint/runtime_binding_block_test.go` bounds where each runtime token may appear). The check binds an external structural property, not the meaning of any sentence. To make it pass, the legacy `→` host-coverage lines the core still carries (`fo-dispatch-core.md:83` names them "legacy host coverage this core still owns", e.g. the `→ Claude: … Codex: … Pi: …` lines at :94/:98) MOVE into the adapters' `## Runtime implementation` blocks — which is exactly the `«fn»` separation the mechanism intends. The semantic residue a grep cannot catch (a claim whose MEANING contradicts a binding, like g6's idempotency) stays a review-time item, but the containment check structurally prevents the largest class (host behavior stated in a shared file at all).
3. **Resolve the two seed instances**: g6's idempotency contradiction, and instance 2's FO-plans-against-missing-capability behavior.

## Relationship to g6 — CONSOLIDATE (captain to confirm at ideation)

This entity is the general audit; `g6` (standing-teammate-idempotency-contract-conflict) is its first known member. Recommend g6 be SUBSUMED here (its fix becomes AC-2 instance 1) and closed as folded-in, so the estate carries one audit rather than a scatter of one-offs — but that closure is a captain/ideation call, not done unilaterally. If kept separate, g6 owns the standing-teammate prose fix and this owns the general sweep + instance 2.

## Acceptance criteria

**AC-1 (VALUE) — A committed contractlint check reds when any shared / host-neutral contract file names a specific runtime host or host-specific tool token; it is green only when the shared files carry no host mention.**
Verified by: a planted divergence (add `Claude`/`Codex`/`Pi` or a host tool token to any of the five shared files) reds the check; removing it greens it. The check is a containment invariant over an enumerated shared-file set and an enumerated host-name/token set — an existence-fact check, not a prose-grep over meaning. Green requires the legacy `→` host-coverage lines to have moved from `fo-dispatch-core.md` into the adapters, so the check's passing is a real relocation, not a suppression.

**AC-2 — An enumerated list of every core-level claim that asserts runtime-varying behavior flat, each with a verdict (bound / should-bind / move-to-adapter / host-neutral).**
Verified by: the list, each row citing file:line and the adapter binding it agrees or conflicts with; a review can reproduce each classification by reading the two locations. This is the semantic residue the containment check (AC-1) cannot catch — a claim whose MEANING contradicts a binding without naming a host.

**AC-3 — The two seed instances are resolved.** g6's core/binding idempotency contradiction is reconciled (one owner, not both), and instance 2's async/bare-mode planning gap is addressed (bound in prose, or recorded as an accepted host-default gap with grounds).
Verified by: for g6, the core and binding no longer contradict — a review-read confirms a single source of truth; for instance 2, either the contract names the binding-consistency rule for fan-out authoring, or a recorded decline with grounds.

**AC-4 (optional/if in scope) — A bare-mode FO asked to author a fan-out does not assume a streaming/async shape its `«async-dispatch»` binding forbids.**
Verified by: a live check — a bare `claude -p` FO (no team tools) authoring a fan-out plan does not describe async streaming; an offline proxy proves the words, never the behavior. Firm this AC at ideation; it may be deferred as a known host-default gap if a live check is disproportionate.

## Out of scope

Re-adding BEHAVIORAL phrase/prose checks — a grep asserting what a program or agent DOES (the banned shape). The AC-1 containment check is the permitted kind: it asserts an existence/structural fact (a host name's presence in a shared file), which the captain's grep ruling explicitly allows and which existing contractlint already does for token containment. Also out of scope: changing any host's actual runtime policy — this reconciles the DESCRIPTION with the binding, it does not change what a runtime does. Not a 0.26.0 blocker: both instances are pre-existing, and instance 2's wording is already corrected in v4dm.
