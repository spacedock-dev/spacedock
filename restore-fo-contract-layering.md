---
id: zyhwwm5a9a08vv6htpef1c34
title: Restore core/adapter and core/mod layering in the FO contract
status: ideation
source: '0221 captain review (2026-06-19, CL): auditing the layered-FO contract surfaced two pre-existing layering leaks in the host-neutral fo-dispatch-core.md. Defer 72 (tier-delegation); land this cleanup + the prose-function notation on a clean foundation first.'
started: 2026-06-19T18:31:45Z
completed:
verdict:
score: 0.5
worktree:
issue:
sprint: 0221-layered-fo
sprint-readiness: ready
group: cleanup
---

The host-neutral FO dispatch core has leaked host-specific and mod-specific detail; restore the core/adapter and core/mod boundaries before more contract is built on it.

## Problem

`skills/first-officer/references/fo-dispatch-core.md` is the host-neutral dispatch core, but two leaks violate that boundary:

1. **Claude model identifiers in the generic core.** Reuse-condition-4 names `opus`/`sonnet`/`haiku` and the `"opus[1m]"` fallback example directly. These are Claude-specific; the host-neutral core should refer generically to "the host's canonical model enum" and delegate the actual enum + the `[1m]` example to the Claude adapter (`claude-fo-dispatch.md`), where the model-for-member lookup and context-budget mapping already live. The behavioral contract (a stamped fallback value forces a one-time fresh re-stamp) must not change.

2. **A GitHub PR-pending scan in the generic event loop.** Event-loop step 1 ("Check PR-pending entities") runs `gh pr view` and advances merged PRs directly in the host-neutral core. PR lifecycle is the `pr-merge` mod's domain — its `idle` hook already does this scan (and its own prose calls itself "defense in depth" for the core's built-in scan). A workflow with no `pr-merge` mod (a `merge: local` or non-code workflow) should never reach for `gh` in its loop. The generic loop should only fire idle hooks; the `pr-merge` idle hook should own PR-pending entirely.

Both are pre-existing debt, but the layered-FO sprint is adding more contract on this foundation, so they are fixed first. A possible third item to assess in ideation: the `status --set` merge-hook guard falsely trips on a non-terminalizing `worktree=` clear (observed 2026-06-19) — confirm scope and fold in or split out.

## Out of scope

- The 72 (fo-tier-delegation) tier map's own model-name references — 72 is deferred; it rebuilds on this clean foundation in a later sprint.
- The prose-function `«»` notation restructure — that is the prose-function-restructure (czw) member, which lands after this cleanup.

## Proposed approach

Two text moves in the FO contract, each restoring a leaked boundary. Both are
behavior-neutral: no Go code reads the contract prose being moved — the model
enum + `[1m]` context window are CODE-resolved (`dispatch build`'s
`effective_model`, `dispatch context-budget`'s family rule), and the event loop
is pure FO instruction the FO runs but never code-reads. The deliverable is the
two doc edits plus a contractlint structural-absence guard (with a paired
discriminator control) for each, modeled on the package's existing
`TestDispatchCoreHasNoClaudeTeamImperative` + `...ScannerDiscriminates` pair.

### Move A — genericize the model enum (TWO leak sites in the core)

`fo-dispatch-core.md` names Claude model tokens directly in **two** places — the
checklist named only reuse-condition-4, but a whole-file scan finds a second site,
the Break-Glass Manual Dispatch template. Both leak the same `sonnet | opus |
haiku` enum (Codex/Pi load this core too and have no such enum). Move the
Claude-specific enum + the `[1m]` fallback example into `claude-fo-dispatch.md`
(the Claude adapter, where the model-for-member lookup and context-budget mapping
already live — `## Context Budget`, line ~142; and the break-glass realization
already lives — line ~48/52), leaving host-neutral statements that defer to "the
host's canonical model enum."

**A1 — reuse-condition-4.**

**Before** (`fo-dispatch-core.md`, reuse-condition-4, the second sentence):

> Members stamped with captain-session fallback values (e.g., `"opus[1m]"`) never
> match enum values (`sonnet`, `opus`, `haiku`) and force a one-time fresh
> dispatch that re-stamps the canonical enum.

**After** (`fo-dispatch-core.md`):

> A member stamped with a captain-session fallback value — one outside the host's
> canonical model enum — never matches an enum value and forces a one-time fresh
> dispatch that re-stamps a canonical enum value. The host's canonical enum and
> its fallback shapes are the runtime adapter's (see the runtime adapter).

**After** (`claude-fo-dispatch.md`, appended to the `## Context Budget and Dead
Ensign Handling` section's model-for-member paragraph, where the same
team-`config.json` member-model read already lives):

> The canonical model enum reuse-condition-4 compares against is `sonnet`, `opus`,
> `haiku` (the `dispatch build` effective_model values). A member stamped with a
> captain-session fallback value (e.g. `"opus[1m]"`) is outside this enum, so it
> never matches and forces the one-time fresh re-stamp.

**A2 — Break-Glass Manual Dispatch template (the second leak site).**

**Before** (`fo-dispatch-core.md`, `## Dispatch Adapter`, Break-Glass paragraph,
last sentence):

> The `model` slot is conditional — include it only when the stage (or
> `stages.defaults`) declares a model from `sonnet | opus | haiku`; otherwise omit
> the entire model argument.

**After** (`fo-dispatch-core.md`):

> The `model` slot is conditional — include it only when the stage (or
> `stages.defaults`) declares a model in the host's canonical model enum;
> otherwise omit the entire model argument. The concrete enum is the runtime
> adapter's.

**After** (`claude-fo-dispatch.md`, its Break-Glass note at line ~52, which
already points back to the core for the conditional `model=` slot — extend it to
own the enum it realizes):

> This is the concrete Claude form of fo-dispatch-core.md's Break-Glass template;
> the contract (what it omits, the conditional `model=` slot, "use only when the
> helper is unavailable") is stated there. The canonical enum the conditional slot
> draws from is `sonnet | opus | haiku`.

Behavior preserved verbatim: a stamped fallback value forces a one-time fresh
re-stamp, and the break-glass `model=` slot is included only for an enum-declared
stage. The host-neutral core keeps both rules; only the Claude-specific tokens
move to the Claude adapter.

### Move B — delete the PR-pending scan from the generic event loop

`fo-dispatch-core.md` `## Event Loop` step 1 runs a `gh pr view` PR-pending scan
in the host-neutral loop. PR lifecycle is the `pr-merge` mod's domain; its `idle`
hook already does this scan. Delete step 1 from the generic loop and renumber
2-4 → 1-3; the `pr-merge` idle hook becomes the sole PR scanner (fired by the
generic loop's existing "Fire `idle` hooks", step 4 today / step 3 after
renumber). Update the `pr-merge` idle hook's "defense in depth" framing (it no
longer backstops a core scan — it IS the scan) and the mod's line-96 detection-
path list (drop "the event loop PR check").

**Before** (`fo-dispatch-core.md`, `## Event Loop`, step 1):

> 1. **Check PR-pending entities** — Run `status --where "pr !=" --json --fields
>    id,slug,pr`. For each entity in `entities`, check PR state via `gh pr view`
>    and advance merged PRs. When advancing a merged PR, clear its `mod-block` if
>    set: `status --set {slug} mod-block=`.
> 2. **Check mod-blocked entities** — …
> 3. **Run `status --next …`** — …
> 4. **If nothing is dispatchable** — Fire `idle` hooks, re-run the host's step-0
>    reconcile sweep …

**After** (`fo-dispatch-core.md` — step 1 deleted, 2-4 renumbered to 1-3):

> 1. **Check mod-blocked entities** — …
> 2. **Run `status --next …`** — …
> 3. **If nothing is dispatchable** — Fire `idle` hooks, re-run the host's step-0
>    reconcile sweep …

The "PR-pending advancement on merge" clearing-`mod-block` detail relocates into
the `pr-merge` idle hook (the only place it now runs). The Claude adapter's
step-0 reconcile-sweep prose (`claude-fo-dispatch.md`) renumbers its
cross-references from "step 4" if any name the old step numbers — verified during
implementation; the only current reference is "step 4's re-run the host's step-0
reconcile sweep" which is positional and survives the renumber (idle remains the
last step).

**Before** (`mods/pr-merge.md` `## Hook: idle`, line 25):

> This provides a periodic re-check in case the event loop's built-in PR scan
> missed a state change (defense in depth).

**After** (`mods/pr-merge.md` `## Hook: idle`):

> This is the workflow's PR-pending scan: the generic event loop fires this idle
> hook and owns no PR scan of its own, so a workflow with no `pr-merge` mod never
> reaches for `gh` in its loop.

**Before** (`mods/pr-merge.md`, line 96):

> The FO handles advancement to the terminal stage and archival when it detects
> the merge (via the event loop PR check, idle hook, or startup hook).

**After** (`mods/pr-merge.md`):

> The FO handles advancement to the terminal stage and archival when it detects
> the merge (via this idle hook, the startup hook, or the reconcile sweep's
> un-advanced-pr class).

Behavior preserved: merged-PR advancement still happens — via the idle hook
(every idle cycle), the startup hook (boot), and the reconcile sweep's
`un-advanced-pr` class (`reconcile.go`, independent roster-derived safety net).
What changes: a `merge: local` / non-code workflow with no `pr-merge` mod no
longer issues `gh` from its event loop.

## Spike determination

No spike needed. The riskiest assumption — "no code reads the moved prose" — is
proven by the existing code path, not by an unexercised mechanism: the model enum
+ `[1m]` mapping are resolved in `dispatch build` / `dispatch context-budget`
(locked by `TestBuildModelPrecedence` and the context-budget family-rule tests),
and the event loop is FO instruction the FO runs but no Go code parses. The
contractlint structural-absence pattern this task reuses is already shipped and
green (`TestDispatchCoreHasNoClaudeTeamImperative` and its discriminator). The one
unverified-until-implementation detail (whether any adapter cross-reference names
the old event-loop step *numbers*) is a single grep at implementation time, not a
mechanism risk — the only current reference is positional ("step 4's re-run").

## Acceptance criteria

Each AC is an end-state property proven by behavior or a legitimate structural
(contractlint) check against an independent source — never a prose-grep
tautology. The independent source for each structural check is the *rule* (a
literal host-coupled token that must not appear in the host-neutral core), not
the file's own prose, so a host-neutral paraphrase passes and a re-introduced
host token fails — the same family as the shipped `claudeTeamDispatchTokens`
check.

**AC-1 — The host-neutral core carries zero Claude model tokens.**
`fo-dispatch-core.md` — scanned whole, so BOTH leak sites (reuse-condition-4 and
the Break-Glass template) are covered — contains the literal token `opus[1m]`
zero times, and contains `sonnet` / `opus` / `haiku` zero times. (After Move A the
enum words appear nowhere in the core, so a bare-word scan is exact, not
over-broad — verified by the pre-move grep finding the words only on the two lines
the move rewrites.) *Proof:* a contractlint structural-absence test
(`TestDispatchCoreHasNoClaudeModelToken`) scanning `fo-dispatch-core.md` for the
literal tokens, PLUS a paired discriminator control (`...ScannerDiscriminates`)
proving the scanner flags a planted host-coupled line (`stamped "opus[1m]" never
matches sonnet/opus/haiku`) and passes the host-neutral paraphrases ("outside the
host's canonical model enum", "a model in the host's canonical model enum"). The
token-set guards against the banned prose-grep: `opus[1m]` and the bare enum words
ARE the host-coupled defect (same as `spawn-standing-all`), not a paraphrasable
meaning. Discriminator keeps it non-vacuous.

**AC-2 — The move RELOCATED the enum, it did not DELETE it.** The Claude-specific
model contract still lives in the system, now in its declared home. *Proof — two
independent signals, NOT a standalone presence-grep:*

1. *Behavioral (the load-bearing proof):* the model enum the contract documents
   is CODE-resolved and stays correct — `TestBuildModelPrecedence` locks
   `sonnet`/`opus`/`haiku`/null as the `dispatch build` effective_model values, and
   the `dispatch context-budget` family-rule tests lock the `[1m]` → 1M mapping.
   The contract text is documentation OF this code; the code is unchanged, so the
   behavior the enum describes is provably intact regardless of which file's prose
   names it. This is shared with AC-5 and is what actually proves "not deleted."
2. *Structural (supplementary, honestly weak):* a contractlint presence assertion
   that `claude-fo-dispatch.md` names the relocated tokens (`opus[1m]` and the enum
   words). This is acknowledged as the weaker half — a paraphrase could drop a
   literal and the move would still be correct — so it is NOT relied on as the
   proof of relocation; it is a low-cost tripwire that the move landed text in the
   adapter, paired with AC-1's absence check on the core (the present-here /
   absent-there pattern the legacy-layering test uses, where the absence half
   carries the real weight). If staff review judges even the supplementary presence
   assertion too prose-grep-adjacent, DROP it and rely on AC-1 (absence from core) +
   AC-2.1 (behavioral) alone — the absence-from-core check plus unchanged code
   behavior already proves the move both happened and broke nothing.

**AC-3 — The host-neutral event loop carries zero `gh` PR-scan.**
`fo-dispatch-core.md` contains the token `gh pr view` zero times. *Proof:* a
contractlint structural-absence test (`TestEventLoopCoreHasNoPRScan`) scanning
`fo-dispatch-core.md` for `gh pr view`, PLUS a discriminator control proving it
flags a planted PR-scan line (`check PR state via gh pr view and advance merged
PRs`) and passes the idle-hook-firing line that remains (`Fire idle hooks, re-run
the host's step-0 reconcile sweep`). `gh pr view` IS the host-coupled defect (a
`merge: local` workflow can't satisfy it), not a paraphrasable meaning.
*Scoping note:* the token is `gh pr view`, NOT `pr !=` — `--where "pr !="` is a
legitimate status-query primitive documented in `first-officer-shared-core.md:48`
and used by the pr-merge mod; banning it would be over-broad. The `gh` reach is
the leak; the status filter is not.

**AC-4 — A no-`pr-merge`-mod workflow has a PR-free generic loop.** After the
change, `gh pr view` on the shipped instruction surface (skills/ + mods/, the same
walk `shippedInstructionMarkdown` already drives) appears ONLY in
`mods/pr-merge.md` (its startup + idle hooks). *Proof:* a contractlint
structural-absence test with an allowed-file map `{mods/pr-merge.md: true}` —
directly modeled on `TestNoUnexpectedModHookOrPRMergeIntroduced`'s
`allowedPRMergeFiles` map (the existing test already restricts PR-merge
*invocations* to the canonical mod; this adds `gh pr view` to that same allow-list
discipline). Any other shipped file naming `gh pr view` fails. A workflow whose
README registers no `pr-merge` mod therefore loads no instruction that reaches
`gh` in its loop — the end-to-end confirmation the checklist's item 3 asks for,
expressed as a checkable surface invariant rather than a live drive. *(The pre-
move grep confirms the only two current `gh pr view` sites are the core line being
deleted and the pr-merge mod; after the move only the mod remains.)*

**AC-5 — Behavior is unchanged.** The reuse re-stamp rule and merged-PR
advancement behave identically. *Proof:* the existing behavioral suites stay
green with no edits — `internal/dispatch` (`TestBuildModelPrecedence` and the
build/reuse hazards suite, the code that actually resolves `effective_model`),
`internal/status` (the merge-policy/terminal-guard suite), and the
`internal/contractlint` closure/ceremony suite (`TestBootResidentDeferredLoadPointsResolve`,
`TestHostNeutralCoresResolveAndCarryCeremony` — the event loop's `## Event Loop`
anchor and the cores' reference closure survive the step deletion). Green-after =
proof the move touched no behavioral seam. This AC owns no NEW test; it is the
gate's cross-check that the existing seams stayed green.

## Sibling finding: `status --set` merge-hook guard false-trip on `worktree=` clear

**Root cause (confirmed).** `internal/status/handlers.go`, `isTerminalUpdate()`
(lines 116-118): a bare `worktree=` clear is classified as a terminal update
unconditionally —

    if u.field == "worktree" && u.hasValue && u.value == "" {
        return true
    }

— so the merge-hook guard at line 203 (`isTerminalUpdate() && modBlock=="" &&
postUpdatePR==""`) fires and refuses a legitimate mid-flight `worktree=` clear
on an entity that is NOT transitioning to a terminal stage and NOT finalizing.
A `worktree=` clear only co-occurs with terminalization in the shared-core
terminalize shape (`completed verdict={v} worktree=`), where `completed`/`verdict`
ALREADY trip `isTerminalUpdate()` independently — so the `worktree`-clears-are-
terminal branch is redundant for the legitimate case and over-broad for the
standalone case.

**Decision: SPLIT OUT.** Rationale:

1. *Different layer, different proof.* This task is a contract-text layering move
   proven by structural-absence checks over markdown. The sibling is a Go control-
   flow bug in `handlers.go` proven by a behavioral golden test (a standalone
   `--set worktree=` on a non-terminal entity under a merge-hook workflow must exit
   0 and clear worktree; today it exits 1). Folding a behavioral status-guard fix
   into a doc-layering PR muddies both the diff and the gate review.
2. *Independent blast radius.* The guard touches the terminalization ceremony —
   the same code path the merge-policy/terminal-guard suite (15+ tests) pins. A fix
   risks regressing that suite and deserves its own validation pass, not a
   ride-along.
3. *Not on this task's critical path.* The layered-FO sprint is adding contract on
   the host-neutral core; this cleanup unblocks that. The status-guard bug does not
   touch the host-neutral core and does not block the contract work.

Recommended split-out: a new sprint task `status-set-worktree-clear-guard` —
fix `isTerminalUpdate()` to not classify a standalone `worktree=` clear as
terminal (gate on `completed`/`verdict`/terminal-`status`, which already cover the
legitimate terminalize shape), with a behavioral golden proving a non-terminal
`--set worktree=` under a registered merge hook now passes. Filed as a sibling,
not folded here.

## Staff review

**Flagged for staff review before the ideation gate.** This is a shipped-contract
+ dispatch-core change (`fo-dispatch-core.md` is the host-neutral dispatch core
every host loads; `claude-fo-dispatch.md` and `mods/pr-merge.md` ship). The review
should confirm: (a) the structural-absence ACs are genuine checks against an
independent rule, not prose-greps (the package's own `doc_test.go` policy is the
rubric); (b) the behavior-neutral claim holds — no code reads the moved prose;
(c) the SPLIT-OUT decision for the sibling finding is right; (d) the before/after
wording preserves the re-stamp rule and the merged-PR advancement paths verbatim.

## Stage Report: ideation

- DONE: Design the two layering moves precisely with before/after wording: (a) genericize fo-dispatch-core.md reuse-condition-4 to 'the host's canonical model enum'...; (b) remove the gh-pr-view PR-pending scan from generic event-loop step 1, delegating PR-pending entirely to the pr-merge mod's idle hook.
  `## Proposed approach` Move A (A1 reuse-condition-4 + A2 Break-Glass — a SECOND model-enum leak site found by whole-file grep) and Move B (event-loop step-1 deletion + pr-merge idle/line-96 rewrites), each with verbatim before/after blocks.
- DONE: Resolve the PROOF approach — behavior-unchanged rides existing suites; the layering invariant needs a LEGITIMATE structural check tied to an independent rule, NOT a prose-grep. State exactly what counts as proof for each AC.
  `## Acceptance criteria` AC-1/AC-3/AC-4 are contractlint structural-absence checks on literal host-coupled tokens (`opus[1m]`, enum words, `gh pr view`) with discriminator controls, modeled on the shipped `TestDispatchCoreHasNoClaudeTeamImperative` pair; AC-5 is the existing-suite green cross-check; AC-2 honestly demotes the presence-grep to a supplementary tripwire behind the code-behavioral proof.
- DONE: Confirm the end-state (PR-free generic loop with no pr-merge mod); assess the sibling status --set merge-hook guard finding and decide fold-in vs split-out; flag for staff review before the gate.
  AC-4 expresses the PR-free end-state as a `gh pr view`-only-in-pr-merge.md surface invariant; sibling root cause confirmed at `handlers.go:116-118` (`isTerminalUpdate()` over-classifies a bare `worktree=` clear) and decided SPLIT-OUT with rationale; `## Staff review` flags the shipped-contract + dispatch-core change.

### Summary

Designed both layering moves with verbatim before/after wording and found a third
leak site the checklist did not name (the Break-Glass template's `sonnet | opus |
haiku` slot), folded into Move A. Grounded the PROOF in the contractlint package's
own shipped structural-absence pattern (literal host-coupled tokens + discriminator
controls, keyed on `gh pr view` not the legitimate `pr !=` primitive), with
behavior-unchanged riding the green `dispatch`/`status`/`contractlint` suites; AC-2
honestly demotes the only presence-grep to a supplementary tripwire. Confirmed the
sibling `worktree=`-clear guard false-trip (`handlers.go:116-118`) and decided
SPLIT-OUT (different layer, independent blast radius, off critical path). No spike
needed — the no-code-reads-the-prose assumption is proven by the existing code
path. Flagged for staff review (shipped-contract + dispatch-core change).
