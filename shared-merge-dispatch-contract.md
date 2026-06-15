---
title: Extract host-neutral merge + dispatch contract — codex/pi name no merge reference and break at terminalization
status: ideation
source: 'captain (2026-06-15, this session) — the shared first-officer merge/dispatch ceremony is siloed in claude-named refs, so codex/pi FOs are MISSING required terminal contract and will likely break at terminalization. Verified — first-officer-shared-core.md:128-130 ONLY defers (the ceremony lives in the runtime merge reference named by the runtime adapter) with no host-neutral fallback, while codex-first-officer-runtime.md and pi-first-officer-runtime.md name NO merge reference and carry no mod-block / ship-local / archive / worktree-removal contract (grep clean). Correction to an earlier FO analysis — codex (send_input + mailbox) and pi (message_dm via pi-agent-teams) DO message, so the reuse/feedback/await machinery is shared too, not Claude-only. This is a correctness gap, not a token cleanup. Captain routed it to 0.20.3 (0203-fo-efficiency; v0.20.2 is latest), otherwise codex breaks without the contract.'
started: 2026-06-15T15:06:22Z
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

## Spike result (gap-before, on the record)

RISKIEST claim — "codex AND pi actually lack a terminal-merge contract today" — confirmed by a structural reference-closure trace this session (account under sustained 429; a live codex drive is a noted follow-up, not the gate proof — a faithful trace is sufficient per the Proof policy). Evidence:

1. **The boot-resident core only defers, with no host-neutral fallback.** `first-officer-shared-core.md:128-130` (`## Merge and Cleanup (deferred module)`) and `:97-99` (`## Dispatch (deferred module)`) both say the machinery "lives in the runtime's merge/dispatch reference … The runtime adapter names the [load point/reference]." Neither carries the ceremony nor a fallback — they are pure pointers.
2. **Only the Claude adapter names a merge/dispatch reference.** `grep -rln claude-fo-merge .` and `claude-fo-dispatch` both return exactly ONE file: `claude-first-officer-runtime.md` (its `## Merge reference (load at terminalization)` :11-13 and `## Dispatch reference (load at first dispatch)` :5-9). No other file names either reference.
3. **codex and pi name NO merge reference and carry no ceremony prose.** `grep -niE "merge|mod-block|ship-local|archive|worktree remov|teardown"` over `codex-first-officer-runtime.md` + `pi-first-officer-runtime.md` returns only two codex hits, both unrelated to terminal merge: line 68 ("a `wait_agent` timeout … is not a teardown trigger") and line 108 (the Backstop: "a missed terminal teardown … stays missed" — it *names* the boundary step but carries zero ceremony). ZERO hits for `mod-block` / `ship-local` / `archive` / `worktree remov` in either file. pi has zero hits entirely.

So a codex/pi FO reaching its terminal stage follows the shared-core pointer to a host merge reference **that does not exist for its host**, and has no contract for the mod-block set→invoke→clear guard, the local-merge sequence, the post-merge sentinel, archival, or worktree-removal safety. The gap is real, not inferred. (A live codex terminal drive proving the break end-to-end is recorded as a follow-up for when the 429 clears; it is not required to gate this design.)

**Seam-boundary spike (where the genuine host residue sits).** Confirmed the Claude-only residue is narrow and already lives outside the moved ceremony, so the cut is clean:
- The terminal-teardown marker `TERMINAL_TEARDOWN_BOUNDED: …` and the `TeamDelete`-race / settle / attempt-cap machinery already live in `using-claude-team` SKILL.md `## Terminal Team Teardown` (:98-109). `claude-fo-merge.md:24` (step 10) only *narrates a pointer* into it. So the marker text never needs to enter a host-neutral file.
- `reconcile` and `context-budget` live under `internal/claudeteam/` (they read `~/.claude/teams/*/config.json`) — host-gated at the binary, genuinely Claude-only.
- `status --set` / `status --archive` / `dispatch trunk` / `dispatch build` are NOT host-gated — host-neutral CLI surfaces the ceremony already calls identically on every host.

## Proposed approach

Extract TWO host-neutral deferred cores and re-home the genuine host residue. Two cores, not one: merge loads at terminalization, dispatch loads at first dispatch — different lazy load points, so keeping them separate preserves the existing lazy-load economics (a boot that greets and stops reads neither; a session that dispatches but never terminalizes reads only the dispatch core). Collapsing them into one file would force the merge ceremony to load at first dispatch.

**New files (product scaffolding, `skills/first-officer/references/`):**
- `fo-merge-core.md` — the host-neutral terminal merge-and-cleanup ceremony.
- `fo-dispatch-core.md` — the host-neutral dispatch procedure, reuse contract, worktree ownership, and event-loop skeleton.

**Composition — the boot-resident core names the host-neutral cores directly; each adapter names only its seam reference.** The shared-core deferred-module pointers stop saying "lives in the runtime's reference" and instead say "the host-neutral ceremony lives in `fo-merge-core.md` / `fo-dispatch-core.md`; the runtime adapter names its host seam." Because the cores are host-agnostic, the boot-resident core can name them directly (they are loaded lazily at the same load points as today — naming ≠ loading). Each runtime adapter keeps a thin `## Merge seam` / `## Dispatch seam` section naming only its host residue. Net load order at terminalization: read `fo-merge-core.md` (the ceremony) + the adapter's merge-seam pointer (Claude → `using-claude-team` `## Terminal Team Teardown`; codex/pi → their teardown analogue).

### Merge: what moves to `fo-merge-core.md` (host-neutral) vs stays per-host (seam)

HOST-NEUTRAL (moves to `fo-merge-core.md`, verbatim — same `spacedock`/`git` steps):
- `## Merge and Cleanup` steps **1-9** (mod-block set-before-invoke, run hooks, detect completion via state delta, blocked-leaves-mod-block, standalone mod-block clear, default local merge, terminalize frontmatter, archive, worktree+local-branch removal).
- `### Ship-Local Ceremony` (all 7 steps — `dispatch trunk`, the post-merge `pr=local-merge:{sha}` sentinel, the set→invoke→clear mandate).
- `### Worktree removal safety` (the no-`--force` default, the audit-then-force ladder).
- `## Mod-Block Enforcement` (the whole section — the mechanism-level invariant in `status --set`/`--archive`, the empty-pr/empty-mod-block recovery, session-resume scan).

PER-HOST SEAM (stays in / moves to each adapter):
- **Merge-and-Cleanup step 10 (terminal agent teardown).** The host-neutral core states the boundary obligation generically — "at the terminal boundary, derive the entity's worker cohort and run the host's terminal-teardown ceremony (cooperative shutdown of the cohort, then the host's bounded team/worker teardown); teardown is mandatory whether merge ran locally or via a PR host." It names NO marker, NO settle interval, NO cap. Each adapter's merge-seam names the concrete teardown:
  - **Claude:** `using-claude-team` `## Terminal Team Teardown` (the `TeamDelete`-race, the ~2s settle, the ≈5 cap, the `TERMINAL_TEARDOWN_BOUNDED` marker). UNCHANGED — already there.
  - **codex:** no team registry → cohort cooperative-shutdown via the mailbox surface, no `TeamDelete`, no marker (the existing codex Backstop already states "a missed terminal teardown stays missed"; the seam states the cohort-shutdown obligation positively).
  - **pi:** `pi-subagents` → a completed child needs no shutdown (mark closed in FO memory); `pi-agent-teams` → the adapter's `member_shutdown`/`team_done` mapping. (Both already in `pi-first-officer-runtime.md` `## Shutdown` — the seam just becomes the terminal-boundary entry point for it.)

### Dispatch: what moves to `fo-dispatch-core.md` (host-neutral) vs stays per-host (seam)

HOST-NEUTRAL (moves to `fo-dispatch-core.md`):
- The per-entity `## Dispatch` procedure (steps 1-9: read entity+stage def, build the ≤3 linchpin checklist, conflict check, resolve `dispatch_agent_id`, the dispatch frontmatter `--set`, the state-transition commit, worktree creation on first worktree-stage dispatch, dispatch via the runtime adapter, wait-for-result).
- `## Worker Resolution` (the `dispatch_agent_id` / `worker_key` split, `:`→`-`, worktree/branch naming) — pure string mechanics, identical on every host.
- `## Worktree Ownership` + `### Split-Root Worktree Contract` — about where entity state lives, not how the host messages.
- The `## Reuse and Fresh Dispatch` **contract** — the reuse *conditions* 1-4 (not-bare, no `fresh: true`, worktree-routing match, model-match) and the reuse/fresh decision, stated host-neutrally with condition 0 (budget probe) and the messaging mechanism (`SendMessage` advance vs `send_input` vs `subagent` re-dispatch) delegated to the adapter seam. The captain's correction is load-bearing here: codex (`send_input`) and pi (`message_dm`/`subagent`) DO message, so the reuse/feedback-routing/await contract is host-neutral; only the concrete handle call is the seam.
- `## Dispatch Adapter` — the `spacedock dispatch build`-mandatory rule, the write-inputs-to-files discipline, the parse-stdout-call-`Agent()`-verbatim flow, and the Break-Glass fallback are host-neutral (codex/pi runtimes ALREADY restate `dispatch build` — that restatement is the duplication this collapses). The host substitutes its spawn call (`Agent()` / `spawn_agent` / `subagent`) at one named point.
- The `## Event Loop` **skeleton** steps 1-4 (PR-pending check → mod-blocked check → `status --next` → idle), which are pure `spacedock status --where`/`--next` reads.

PER-HOST SEAM (stays in / moves to each adapter):
- **Claude:** `## Team Creation` (`using-claude-team`, TeamCreate-first sequencing), the no-pre-dispatch-probe registry-desync rule (#36806), `## Degraded Mode (spacedock seams)`, `## Context Budget and Dead Ensign Handling` (the probe that realizes reuse-condition-0), Event-Loop **step 0 reconcile sweep** (Class A-E) and `### Backstop (Claude)`, the spawn-call (`Agent()` + `spacedock dispatch build`), the `SendMessage` reuse-advance handle.
- **codex:** `spawn_agent`/`send_input`/`wait_agent`, the mailbox completion signal, "declares no budget probe" (reuse-condition-0 satisfied), no reconcile (`### Backstop (Codex)` already says so).
- **pi:** `subagent(...)` / `pi-agent-teams`, `context: "fresh"` default, the epoch-based stale-reuse guard, the live-harness isolation. (All already in `pi-first-officer-runtime.md`.)

**Net:** codex/pi GAIN the terminal-merge contract and the dispatch procedure they currently lack (by naming the new host-neutral cores); the generic ceremony is stated ONCE; the Claude-only residue (registry-desync, teardown-race+marker, reconcile, budget probe) becomes a visible, auditable per-host seam instead of being buried inside generic prose only Claude can reach.

## Out of scope

- **The Claude-only residue is NOT genericized.** The `TERMINAL_TEARDOWN_BOUNDED` marker, the `TeamDelete`-race/settle/cap, registry-desync #36806, the reconcile sweep, and the context-budget probe stay in Claude files (`using-claude-team`, the Claude adapter, `internal/claudeteam/`). This task re-homes the host-neutral ceremony; it does not touch the host residue's behavior.
- **No behavioral change to the ceremony itself.** Same `spacedock status --set`/`--archive`, `dispatch trunk`, `dispatch build`, `git worktree` steps, same order, same MUSTs. This is a re-home (move text + repoint references), not a redesign.
- **The per-line token cuts** belong to the parallel proposal `docs/dev/_proposals/fo-contract-token-cleanup-2026-06-15.md` (the trim). This task is the structural re-home; word-level reduction is that proposal's concern. (The re-home will *enable* dedup savings, but measuring/claiming them is out of scope here.)
- **No codex/pi teardown *mechanism* is built.** codex/pi gain the merge *contract* (they now name an existing ceremony). Their teardown seams describe the obligation against their existing shutdown surfaces (codex mailbox / pi `member_shutdown`); no new binary support or marker is added for them.

## Documentation impact

No user-visible behavior changes — no CLI output, banner, command surface, or docs-site content changes. The edits are entirely to internal `skills/first-officer/references/` scaffolding the FO loads at runtime; the captain-facing behavior (the same merge ceremony runs, the same gate flow) is byte-for-byte preserved on Claude and newly-available (not changed) on codex/pi. No doc diff against the docs site is required. (Per the ideation stage def's doc-diff rule: it applies when a task changes user-visible behavior; this one does not.)

## Out of scope

{Ideation pins. Likely: the Claude-only residue stays in claude files (no attempt to genericize the registry-desync/teardown-race/reconcile machinery); no behavioral change to the ceremony itself (same `spacedock`/`git` steps); the per-line token cuts (separate proposal `docs/dev/_proposals/fo-contract-token-cleanup-2026-06-15.md`).}

## Acceptance criteria

**AC-1 — Reference-closure: from each host's terminal-merge entry, the named merge reference resolves to an existing file carrying the mod-block / local-merge / archive / worktree-removal-safety ceremony — for codex and pi, not only Claude.**
Verified by: a **reference-closure walk** (not a prose-grep): start at each adapter's merge-seam entry (the shared-core deferred-module pointer → the adapter's `## Merge seam` → the named file), follow each named reference, and assert the walk terminates at a file that EXISTS and contains the four ceremony anchors (mod-block enforcement, local-merge/Ship-Local, archive, worktree-removal safety). The independent source is the reference graph itself — a Go test (or a committed `scripts/` checker invoked by a Go test) that parses the three adapter files + the two new cores, builds the name→file edge set, and FAILS if any host's walk dead-ends at a missing file or a file lacking a ceremony anchor. Binds to file existence + the graph, not to a phrase appearing somewhere. The **gap-before** half is already on the record (Spike result above: codex/pi walks dead-end today); AC-1 is satisfied when the same walk for all three hosts terminates at the ceremony file post-extraction. Live codex terminal drive: noted follow-up (account 429), not required for this AC.

**AC-2 — Single-source: the host-neutral ceremony text exists in exactly one file; no host adapter restates it.**
Verified by: a **single-source structural check** binding the host-neutral text to the absence of a restated copy. The independent source is a content-fingerprint of the moved ceremony blocks (the mod-block set→invoke→clear sequence, the Ship-Local steps, the worktree-removal ladder): a Go test asserts each fingerprinted block appears in `fo-merge-core.md` / `fo-dispatch-core.md` AND that none of `claude-/codex-/pi-first-officer-runtime.md` contains it. Not a tautological "grep finds the phrase in the core" — it is "the block lives in the core XOR in any adapter," failing if a host file carries a restated copy. (Distinguishes a true single-source extraction from a copy-paste-into-each-host regression.)

**AC-3 — Behavior-preserving on Claude: the existing live-e2e stays green and no MUST / qualifier is dropped from the moved text.**
Verified by: (a) the existing `internal/ensigncycle` live-e2e cycle (`TestLiveEnsignCycle` and the teardown-grade watcher tests) stays GREEN — the `TERMINAL_TEARDOWN_BOUNDED` marker is asserted byte-exact in `teardown_grade_watcher_test.go:13` / `streamwatch_test.go:57`, so any drift in the moved step-10 narration fails an existing test (this is the load-bearing behavior gate, independent of the new structural checks); and (b) a **word-level no-MUST-dropped audit** of the moved blocks: diff the relocated text against the pre-move `claude-fo-merge.md` / `claude-fo-dispatch.md` and assert zero MUST / refusal / qualifier deletions (same audit method 7e used at `headless-dispatch-mode-intent.md:260` — "contract diff vs main drops NO MUST/qualifier"). User-visible behavior: NONE changes → no doc-diff required (see Documentation impact).

## Test plan

**Cost/complexity:** moderate-high — this is a boot-resident contract touching all three runtime adapters plus two new deferred cores. The risk is behavior drift in the moved Claude text (caught by AC-3a's existing live-e2e) and an incomplete extraction leaving a restated copy (caught by AC-2). The gap-proof spike (AC-1 before-half) is already done structurally this session.

1. **AC-1 reference-closure test (new Go test, fixture-level — no live drive).** Parse the three adapter files + the two cores, build the name→file reference graph, walk from each host's merge/dispatch entry, assert termination at an existing ceremony file with all four anchors. RED before extraction for codex/pi (walks dead-end — the recorded gap), GREEN after. This is the TDD seed: write it first, watch codex/pi fail, do the extraction, watch them pass.
2. **AC-2 single-source fingerprint test (new Go test, fixture-level).** Fingerprint the moved ceremony blocks; assert each in exactly one core and absent from every adapter. RED if a copy-paste regression leaves a restated block in a host file.
3. **AC-3a live-e2e (existing, re-run).** Run the `internal/ensigncycle` cycle on a real Claude drive after the extraction — green confirms the moved step-10 teardown narration still drives the byte-exact marker the watcher grades. This is `live` (real Claude), gated behind the existing harness; it is the one real-drive cost and it already exists.
4. **AC-3b word-level audit (manual, recorded in the implementation stage report).** `git worktree add --detach` the pre-extraction tree, diff moved blocks, assert no MUST/qualifier dropped; record the net word delta (expected: a reduction from dedup, but the trim is the parallel proposal's claim — here we only assert *no loss*).

**Fixture vs CLI vs live:** AC-1/AC-2 are fixture-level Go tests over the instruction files (fast, deterministic, the structural gates). AC-3a is the single existing live test (real Claude, no new live infra). No new CLI surface, so no golden-fixture CLI test. A live codex terminal drive proving the end-to-end merge under the new ref is a **noted follow-up** (account 429) — not in this task's gate.

**Rebase ordering (coordination).** 7e (`headless-dispatch-mode-intent`, id `7ea4…`) merges to main FIRST — it is at `validation` and validated-clean. 7e touched `first-officer-shared-core.md`'s Single-Entity Scope + step 9 and the Claude-adapter exception; it did NOT touch `claude-fo-merge.md`/`claude-fo-dispatch.md`, the deferred-module pointers (:97-99, :128-130), or the adapter load-point sections — confirmed disjoint from 2y's edit set (grep over `headless-dispatch-mode-intent.md` returns no `claude-fo-merge`/`claude-fo-dispatch`/`deferred module`/`Merge and Cleanup` hit). So 2y branches off post-7e main; no conflict expected. Note: 7e's validation also confirmed `claude-fo-merge.md:24`'s `TERMINAL_TEARDOWN_BOUNDED` marker is byte-identical to `teardown_grade_watcher_test.go:13` — the extraction MUST keep that line byte-identical (it is the step-10 teardown seam that STAYS in the Claude file, not moved).

**Complexity assessment for staff review.** This is complex by the ideation stage definition's own examples (skill integration + a boot-resident contract spanning all three runtimes). Recommend an independent staff review before the gate: the review should check (a) the one-core-vs-two call and the boot-resident-names-the-core composition, (b) that the seam boundary correctly leaves the genuine Claude residue per-host (especially that step-10's marker stays in the Claude file), (c) AC-1/AC-2 bind to independent sources and are not tautological prose-greps, and (d) the AC-3a live-e2e is the right behavior gate. The riskiest unverified mechanism (does codex/pi actually lack the contract) was exercised structurally first and recorded in the Spike result above.

## Stage Report: ideation

- DONE: RISKIEST SPIKE FIRST (AC-1 seed): confirm codex AND pi ACTUALLY lack a terminal-merge contract today — structurally trace that shared-core defers with no host-neutral fallback, codex/pi name no merge reference and carry no ceremony prose, and the ceremony exists only in claude-fo-merge.md/claude-fo-dispatch.md.
  Reference-closure trace recorded in `## Spike result`: `grep -rln claude-fo-merge/dispatch` → only `claude-first-officer-runtime.md`; `grep -niE "merge|mod-block|ship-local|archive|worktree remov|teardown"` over codex/pi runtimes → 2 codex hits both unrelated to terminal merge, 0 in pi; shared-core :97-99/:128-130 are pure pointers. Gap confirmed on the record. Live codex drive deferred (account 429) as a noted follow-up.
- DONE: Pin the extraction design with before/after structure: which sections are HOST-NEUTRAL vs the genuine CLAUDE-ONLY seams; whether dispatch+merge are ONE core or TWO; and exactly how shared-core / each runtime adapter names + composes host-neutral-core + per-host seam.
  `## Proposed approach`: TWO cores (`fo-merge-core.md` + `fo-dispatch-core.md` — distinct lazy load points). Boot-resident core names the cores directly; each adapter keeps a thin merge/dispatch SEAM. Per-section move/stay tables for both merge (steps 1-9 + Ship-Local + worktree-safety + Mod-Block Enforcement move; step-10 teardown marker stays per-host) and dispatch (procedure/reuse-contract/worker-resolution/event-loop-skeleton move; team-lifecycle/reconcile/budget-probe/registry-desync stay Claude). Seam-boundary spike confirmed the marker + reconcile + context-budget already sit outside the moved ceremony.
- DONE: Behavioral ACs (NO prose-grep): AC-1 codex/pi gain a NAMED existing terminal-merge contract (reference-closure); AC-2 generic ceremony stated ONCE with no restated per-host copy (single-source structural check); AC-3 behavior-preserving for claude (live-e2e green + word-level no-MUST-dropped audit). State whether user-visible behavior changes.
  `## Acceptance criteria` AC-1 = reference-graph closure walk (Go test, binds to file existence + graph); AC-2 = content-fingerprint single-source (block in core XOR adapter, not a tautological grep); AC-3 = existing `internal/ensigncycle` live-e2e + byte-exact `TERMINAL_TEARDOWN_BOUNDED` watcher assertion + word-level MUST-audit. `## Documentation impact`: NO user-visible change → no doc-diff required.

### Summary

Designed the extraction of two host-neutral deferred cores (`fo-merge-core.md`, `fo-dispatch-core.md`) named directly by the boot-resident core, with each runtime adapter reduced to a thin host seam. Confirmed the gap structurally first (codex/pi merge-reference walks dead-end today; on the record) and pinned the seam boundary precisely — the genuine Claude-only residue (TERMINAL_TEARDOWN_BOUNDED marker, reconcile, context-budget, registry-desync #36806) already lives outside the moved ceremony, so the cut is clean and the step-10 marker line stays byte-identical to its Go-test assertion. ACs bind to independent sources (reference graph, content-fingerprint, existing live-e2e) — no prose-grep. Assessed COMPLEX (boot-resident contract across all three runtimes); recommended an independent staff review before the gate, with the four review focal points named in the test plan.
