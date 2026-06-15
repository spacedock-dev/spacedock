---
title: Extract host-neutral merge + dispatch contract — codex/pi name no merge reference and break at terminalization
status: backlog
source: 'captain (2026-06-15, this session) — the shared first-officer merge/dispatch ceremony is siloed in claude-named refs, so codex/pi FOs are MISSING required terminal contract and will likely break at terminalization. Verified — first-officer-shared-core.md:128-130 ONLY defers (the ceremony lives in the runtime merge reference named by the runtime adapter) with no host-neutral fallback, while codex-first-officer-runtime.md and pi-first-officer-runtime.md name NO merge reference and carry no mod-block / ship-local / archive / worktree-removal contract (grep clean). Correction to an earlier FO analysis — codex (send_input + mailbox) and pi (message_dm via pi-agent-teams) DO message, so the reuse/feedback/await machinery is shared too, not Claude-only. This is a correctness gap, not a token cleanup. Captain routed it to 0.20.3 (0203-fo-efficiency; v0.20.2 is latest), otherwise codex breaks without the contract.'
started:
completed:
verdict:
score: 0.55
worktree:
issue:
id: 2yfsf01jf15fmts7xt7w71m2
sprint: 0203-fo-efficiency
---

The shared first-officer terminal-merge and dispatch ceremony is generic (pure `spacedock status` + `git`) but lives only in `claude-fo-merge.md` / `claude-fo-dispatch.md`. The boot-resident core defers to "the runtime's merge reference," but only Claude provides one. A codex or pi FO that reaches a terminal stage has no mod-block set->invoke->clear, no Ship-Local ceremony, no worktree-removal safety, and no archive sequence to follow. It improvises or breaks. Extract the host-neutral ceremony so every runtime has it; leave only the genuine host seam per file.

## Problem

`first-officer-shared-core.md:128-130` (`## Merge and Cleanup (deferred module)`) defers the entire terminal ceremony to a per-host merge reference and provides no host-neutral fallback. Only `claude-fo-merge.md` exists; `codex-first-officer-runtime.md` and `pi-first-officer-runtime.md` name no merge reference and contain no merge/mod-block/ship-local/archive/worktree-removal prose (verified by grep this session). So:
- A codex/pi FO reaching terminalization follows the shared-core pointer to a runtime merge reference that does not exist, and has no contract for the mod-block guard, the local-merge sequence, archival, or worktree-removal safety. This is the break.
- The same silo wastes tokens (the generic ceremony cannot be shared across hosts; the token-cleanup audit independently kept `claude-first-officer-runtime.md` as load-bearing precisely because it is the SOLE namer of `claude-fo-merge.md`), but the token cost is secondary to the correctness gap.

Dispatch has the same shape: the dispatch procedure, the event loop (PR/mod-block/next), worktree ownership, and the reuse principle are generic and duplicated in prose across the three runtime files. Codex and pi DO message (`send_input`/mailbox; `message_dm`), so the reuse/feedback-routing/await-completion contract is shared, not Claude-specific. The genuinely Claude-only residue is narrow: the Claude Code registry-desync recovery (#36806), the async teardown race + `TERMINAL_TEARDOWN_BOUNDED` marker, the reconcile sweep, and the context-budget probe.

## Proposed approach

{Ideation designs; the direction:}
- Extract a host-neutral `fo-merge.md` carrying the generic ceremony (mod-block enforcement, the set->invoke->clear sequence, Ship-Local, worktree-removal safety, archive). ALL three runtime adapters name it (or the boot-resident core names it directly, since it is host-agnostic). Each host file keeps ONLY its terminal teardown seam: Claude defers to `using-claude-team` `## Terminal Team Teardown` (already its owner); codex/pi carry their teardown analogues.
- Extract a host-neutral dispatch core (procedure, event loop, worktree ownership, reuse principle); each host file keeps only its dispatch substrate + genuine seams (Claude: team lifecycle + registry-desync recovery + budget probe + reconcile; codex: spawn_agent/wait_agent/no-reconcile; pi: subagent/fresh-default).
- Net: codex/pi GAIN the terminal-merge contract they currently lack; the generic ceremony is stated once; the Claude-only residue becomes visible and auditable instead of buried in generic prose.

## Out of scope

{Ideation pins. Likely: the Claude-only residue stays in claude files (no attempt to genericize the registry-desync/teardown-race/reconcile machinery); no behavioral change to the ceremony itself (same `spacedock`/`git` steps); the per-line token cuts (separate proposal `docs/dev/_proposals/fo-contract-token-cleanup-2026-06-15.md`).}

## Acceptance criteria

**AC-1 — A codex (and pi) FO has a terminal-merge contract: reaching a terminal stage, it follows a named, existing merge reference covering mod-block, local merge, archive, and worktree-removal safety.**
Verified by: {the riskiest mechanism — a live or faithfully-traced codex terminal drive that currently has NO merge contract is the smallest exercise proving the gap is real; after the extraction, the same drive follows the host-neutral ref. Ideation pins the exercise. Spike this FIRST: confirm codex terminalization today actually lacks the contract, not just structurally inferred.}

**AC-2 — The generic ceremony is stated once (host-neutral), and each host file carries only its seam, with no restated generic ceremony.**
Verified by: {a single-source structural check binding the host-neutral ceremony to the absence of a restated copy in any host runtime file — an independent-source check, not a prose-grep.}

## Test plan

{Ideation fills; the spike (AC-1, confirm codex currently breaks at terminalization) seeds it. Behavior-first over a real codex terminal drive plus the structural single-source check.}
