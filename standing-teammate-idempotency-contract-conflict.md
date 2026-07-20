---
id: g6c8zafrdm5s4vm6fnnveejr
title: Host-neutral core and Claude binding disagree on whether standing-teammate injection dedups
status: backlog
source: "Surfaced by member z7 during 0260 implementation and raised rather than silently reconciled; independently hit by the Commander FO at boot the same day, which had to resolve the ambiguity by judgement before its first dispatch."
started:
completed:
verdict:
score:
worktree:
issue:
---

The host-neutral dispatch core and the Claude runtime binding make opposite claims about whether standing-teammate injection is idempotent.

## Problem

`skills/first-officer/references/fo-dispatch-core.md` describes standing-teammate injection as "Idempotent (already-alive members omitted)".

`skills/first-officer/references/claude-fo-dispatch.md` says the call "does NOT dedup against a team config (there is none keyed by name), so it emits every declared standing teammate; idempotency is your own-roster concern — do not re-inject a standing teammate you already spawned this session."

These cannot both hold. The host-neutral claim says the mechanism guarantees idempotency; the host binding says the mechanism guarantees nothing and the FO must track it. An FO reading the core and trusting it will re-inject a duplicate standing teammate on a second dispatch; an FO reading only the binding carries bookkeeping the core says is unnecessary.

Observed, not theoretical: the 0260 Commander FO ran `spacedock dispatch spawn-standing-all` at boot, received a spawn spec for an already-declared teammate, and had to decide by judgement whether re-spawning was safe. It deferred the injection to point-of-need instead — a reasonable call, but one the contract should have made for it rather than leaving to improvisation.

## Proposed approach

Ideation fills this in. The decision is which claim is true, and that is a question about the binary's behavior, not about the prose: `spacedock dispatch spawn-standing-all` either filters already-alive members or it does not. Establish that by exercising the command against a session with a live standing teammate, then make the losing sentence agree with the observed behavior.

Note the shape: this is a prose-to-prose contradiction where one side is checkable against a real independent value (what the command actually emits). That makes it a candidate for a binding rather than a wording fix — the same treatment 0260's contractlint retirement applied to runtime-semantics claims. Whether that binding is worth its cost is an ideation question, not a foregone conclusion.

## Out of scope

Changing standing-teammate lifetime, naming, or the mod-declaration format. Adding a dedup mechanism to the binary if the honest answer is that the FO owns the bookkeeping — in that case the fix is to correct the host-neutral sentence, not to build the guarantee it wrongly promises.

## Acceptance criteria

Ideation fills these in. The value AC must rest on OBSERVED command behavior — what `spawn-standing-all` emits when a declared standing teammate is already alive — not on the two instruction files agreeing with each other, since two prose files agreeing is exactly the tautology 0260 retired.
