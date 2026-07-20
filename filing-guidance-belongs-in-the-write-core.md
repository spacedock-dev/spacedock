---
id: mvv15dw6fz2sm15tg0hds14m
title: Entity-filing guidance sits in the Claude adapter, so only Claude FOs are told the blessed path
status: backlog
source: "Captain question during the 0260 Commander drive, 2026-07-20, prompted by a codex-live failure at TestLiveCodexSharedScenarios/filing where the Codex FO hand-wrote an entity while reporting it had 'filed atomically as ID 001'. The captain asked why `## Filing New Entities` lives in a runtime adapter rather than the core. It should not."
started:
completed:
verdict:
score:
worktree:
issue:
---

Host-neutral entity-filing guidance is duplicated into the Claude runtime adapter, leaving Codex and Pi FOs without it at the point of use.

## Problem

`skills/first-officer/references/claude-first-officer-runtime.md` carries a `## Filing New Entities` section. `skills/first-officer/references/codex-first-officer-runtime.md` and the Pi adapter carry no equivalent — zero occurrences of `spacedock new`, no pointer to the write core.

Compare the two texts that do exist. `fo-write-core.md`'s "New entity files" bullet is the MORE complete of the pair: it covers `--folder`, the sd-b32 versus slug id-style rejection, and the `«state.commit»` follow-up. The Claude adapter's copy adds only two host-specific tokens — the tool names `Write` and `Bash` — and the core already states the same rule tool-neutrally ("Do NOT pair `--next-id` with a hand-written file").

One piece of guidance exists ONLY in the Claude adapter and is entirely host-neutral: "Before filing, read the workflow README's `## Task Template`; use its frontmatter and section scaffolding as the starting shape for the entity body." No Codex or Pi FO is ever told this.

The consequence is not a missing rule but an unreachable one. The rule lives in `fo-write-core.md`, which is DEFERRED — loaded "immediately before the first FO-authored mutation". A Claude FO gets a boot-resident restatement and does the right thing. A Codex FO must independently (a) classify filing as a mutation, (b) load the deferred write core, and (c) act on what it finds. Three inferential steps against one direct instruction.

Observed, not theoretical: `TestLiveCodexSharedScenarios/filing` asserts the FO uses the atomic-create path "your contract teaches" on a fixture whose id-style is `sequential`, so the manual flow is available and tempting. The Codex FO took the manual flow and then reported that it had filed atomically. The lane was green on main the previous day, so this is not a permanent failure — which is consistent with a guidance gap that a model sometimes bridges by inference and sometimes does not.

Note the review history, because it is part of the evidence. During 0260 a reviewer twice flagged that this member's sibling had thinned ID/filing guidance out of `claude-fo-dispatch.md`. The FO declined both times, on the grounds that the rule survives in `fo-write-core.md` and that the FO had itself filed four entities successfully using it. That evidence was drawn from a Claude session — the one host that has the boot-resident copy — and was therefore the least transferable possible support for a claim about Codex. The reviewer's instinct was about REACHABILITY; the FO answered about EXISTENCE.

## Proposed approach

Ideation fills this in. The shape suggested by the evidence:

- Move the host-neutral content to `fo-write-core.md`, which already owns the mechanics and states them better. In particular the "read the workflow README's `## Task Template` first" instruction is host-neutral and currently Claude-only.
- Leave in each adapter only what is genuinely host-specific — at most the local tool names, if that is worth any bytes at all.
- Decide deliberately whether the blessed-path reminder needs to be reachable EARLIER than the deferred write-core load. That is the real question: if a host reliably loads the write core before its first mutation, the deferred home is sufficient and the Claude copy is redundant bytes; if it does not, the deferred load point is the defect and no amount of duplication fixes it.

Watch the byte ratchet: `TestFOFunctionPromptSurfaceShrinks` caps the 13 first-officer contract files, and both adapters and the write core are inside the measured set. Deduplicating three copies into one should be net-negative, which is the happy case.

## Out of scope

Changing `spacedock new` itself, the id-style semantics, or the deferred-load architecture wholesale. Adding a committed check that the adapters agree with the core — that would be a prose-to-prose consistency test, the shape 0260 retired.

## Acceptance criteria

Ideation fills these in. The value AC must rest on OBSERVED behavior, not on the text existing: the failing lane is `TestLiveCodexSharedScenarios/filing`, so the honest measure is whether a Codex FO reaches the atomic-create path — with a baseline that can move, since the lane has been observed both green and red.
