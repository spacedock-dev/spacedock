---
id: zyhwwm5a9a08vv6htpef1c34
title: Restore core/adapter and core/mod layering in the FO contract
status: validation
source: '0221 captain review (2026-06-19, CL): auditing the layered-FO contract surfaced two pre-existing layering leaks in the host-neutral fo-dispatch-core.md. Defer 72 (tier-delegation); land this cleanup + the prose-function notation on a clean foundation first.'
started: 2026-06-19T18:31:45Z
completed:
verdict:
score: 0.5
worktree: .worktrees/spacedock-ensign-restore-fo-contract-layering
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
in the host-neutral loop AND does a second thing: it clears the merged entity's
`mod-block` (`status --set {slug} mod-block=`) so the subsequent terminalization
passes the merge-hook guard. PR lifecycle is the `pr-merge` mod's domain; its
`idle` hook already does the scan. Delete step 1 from the generic loop and
renumber 2-4 → 1-3; the `pr-merge` idle hook becomes the sole PR scanner (fired by
the generic loop's existing "Fire `idle` hooks", step 4 today / step 3 after
renumber).

**The mod-block-clear must be RELOCATED, not just dropped (gap closed).** Core
step-1 clears `mod-block` before terminalizing. The SHIPPED TEMPLATE
`./mods/pr-merge.md` does NOT clear `mod-block` in its merged-PR advancement
(grep: `mod-block` count ZERO) — it relies on the core's step-1 clear. So deleting
step-1 without relocating the clear would leave a shipped-template workflow's
merged-PR-while-mod-blocked finalize blocked by the terminalize guard
(`handlers.go:156-169` combined-clear / mod-block-pending refusal). The dev local
copy `docs/dev/_mods/pr-merge.md` (count 11) is the correct reference: its
startup/idle advancement already does the two-step clear-then-terminalize. Port
that two-step shape into the SHIPPED template's startup AND idle merged-PR
advancement (a STANDALONE `mod-block=` `--set`, separate from the terminal fields,
per `fo-merge-core.md:19` — the mechanism refuses combining `mod-block=` with
terminal fields).

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

The Claude adapter's step-0 reconcile-sweep prose (`claude-fo-dispatch.md`)
renumbers its cross-references from "step 4" if any name the old step numbers —
verified during implementation; the only current reference is "step 4's re-run the
host's step-0 reconcile sweep" which is positional and survives the renumber (idle
remains the last step).

**Before** (`mods/pr-merge.md` `## Hook: startup`, the `MERGED` advancement —
shipped template, which never clears mod-block):

> If `MERGED`, advance the entity to its terminal stage: set `status` to the
> terminal stage, `completed` to ISO 8601 now, `verdict: PASSED`, clear
> `worktree`, archive the file, and clean up any worktree/branch. Report each
> auto-advanced entity to the captain.

**After** (`mods/pr-merge.md` `## Hook: startup`):

> If `MERGED`, advance the entity to its terminal stage. Because a `mod-block` may
> be set while the PR is pending, the clear and the terminalization are two
> separate `--set` calls (the mechanism refuses combining `mod-block=` with
> terminal fields):
> 1. `spacedock status --workflow-dir {dir} --set {slug} mod-block=` when a
>    `mod-block` is set (skip when empty);
> 2. `spacedock status --workflow-dir {dir} --set {slug} status={terminal}
>    completed verdict=PASSED worktree=`, then `spacedock status --workflow-dir
>    {dir} --archive {slug}`.
>
> Clean up any worktree/branch. Report each auto-advanced entity to the captain.

**Before** (`mods/pr-merge.md` `## Hook: idle`):

> Check PR-pending entities using the same logic as the startup hook: scan entity
> files for non-empty `pr` and non-terminal status, run `gh pr view` for each, and
> advance merged PRs. This provides a periodic re-check in case the event loop's
> built-in PR scan missed a state change (defense in depth). Report any advanced
> entities to the captain.

**After** (`mods/pr-merge.md` `## Hook: idle`):

> Check PR-pending entities using the same logic as the startup hook: scan entity
> files for non-empty `pr` and non-terminal status, run `gh pr view` for each, and
> advance merged PRs (two-step `mod-block=` clear then terminalize). This is the
> workflow's PR-pending scan: the generic event loop fires this idle hook and owns
> no PR scan of its own, so a workflow with no `pr-merge` mod never reaches for
> `gh` in its loop. Report any advanced entities to the captain.

**Before** (`mods/pr-merge.md`, last paragraph — the detection-path list):

> The FO handles advancement to the terminal stage and archival when it detects
> the merge (via the event loop PR check, idle hook, or startup hook).

**After** (`mods/pr-merge.md`):

> The FO handles advancement to the terminal stage and archival when it detects
> the merge (via this idle hook, the startup hook, or the reconcile sweep's
> un-advanced-pr class).

Behavior preserved: merged-PR advancement still happens — via the idle hook
(every idle cycle), the startup hook (boot), and the reconcile sweep's
`un-advanced-pr` class (`reconcile.go`, independent roster-derived safety net) —
AND the mod-block clear that the terminalize guard depends on now rides the mod's
own advancement (relocated from core step-1), so a merged-PR-while-mod-blocked
finalize still passes the guard. What changes: a `merge: local` / non-code
workflow with no `pr-merge` mod no longer issues `gh` from its event loop.

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
model contract still works after the move. *Proof — behavioral, no presence-grep:*
the model enum the contract documents is CODE-resolved and stays correct —
`TestBuildModelPrecedence` locks `sonnet`/`opus`/`haiku`/null as the `dispatch
build` effective_model values, and the `dispatch context-budget` family-rule tests
lock the `[1m]` → 1M mapping. The contract text is documentation OF this code; the
code is unchanged, so the behavior the enum describes is provably intact
regardless of which file's prose names it. Paired with AC-1 (absence from the
core), this proves the move both happened (the tokens left the core) and broke
nothing (the behavior they documented still runs). The structural presence-grep
"the adapter names the relocated tokens" is DROPPED per staff review — it is the
banned prose-grep (a paraphrase could drop a literal and the move would still be
correct), and AC-1 + the behavioral proof already carry the full weight.

**AC-3 — The host-neutral event loop carries zero `gh` PR-scan.**
`fo-dispatch-core.md` contains the token `gh pr view` zero times. *Proof:* a
contractlint structural-absence test (`TestEventLoopCoreHasNoPRScan`) scanning
`fo-dispatch-core.md` for `gh pr view`. **The paired discriminator control
(`TestEventLoopPRScanScannerDiscriminates`) ships WITH it — the absence half
without the discriminator is not acceptable** (a typo'd token would pass
vacuously). The discriminator MUST: (a) flag a planted PR-scan line (`check PR
state via gh pr view and advance merged PRs`), and (b) pass the idle-hook-firing
line that legitimately remains (`Fire idle hooks, re-run the host's step-0
reconcile sweep`). `gh pr view` IS the host-coupled defect (a `merge: local`
workflow can't satisfy it), not a paraphrasable meaning. *Scoping note:* the token
is `gh pr view`, NOT `pr !=` — `--where "pr !="` is a legitimate status-query
primitive documented in `first-officer-shared-core.md:48` and used by the pr-merge
mod; banning it would be over-broad. The `gh` reach is the leak; the status filter
is not.

**AC-4 — A no-`pr-merge`-mod workflow has a PR-free generic loop.** After the
change, `gh pr view` on the shipped instruction surface (skills/ + mods/, the same
walk `shippedInstructionMarkdown` already drives) appears ONLY in
`mods/pr-merge.md` (its startup + idle hooks). *Proof:* a contractlint
structural-absence test with an allowed-file map `{mods/pr-merge.md: true}` —
directly modeled on `TestNoUnexpectedModHookOrPRMergeIntroduced`'s
`allowedPRMergeFiles` map (the existing test already restricts PR-merge
*invocations* to the canonical mod; this adds `gh pr view` to that same allow-list
discipline). Any other shipped file naming `gh pr view` fails. **Two discriminator
controls ship WITH the absence half** (modeled on
`TestPortabilityCheckDiscriminatesHostSpecific`, which proves the legitimately
host-specific form is present yet not flagged):

1. *Positive control:* assert `mods/pr-merge.md` DOES contain `gh pr view` (its
   startup + idle hooks legitimately scan), so the allow-list entry is load-bearing
   — if the mod ever stopped carrying the token the control reds, proving the
   allow-list exempts a real occurrence rather than a vacuous one.
2. *Negative control:* a planted non-allowed file carrying `gh pr view` is flagged
   by the same scan, proving the allow-list actually constrains.

A workflow whose README registers no `pr-merge` mod therefore loads no instruction
that reaches `gh` in its loop — the end-to-end confirmation the checklist's item 3
asks for, expressed as a checkable surface invariant rather than a live drive.
*(The pre-move grep confirms the only two current `gh pr view` sites are the core
line being deleted and the pr-merge mod; after the move only the mod remains.)*

**AC-5 — Behavior is unchanged.** The reuse re-stamp rule and merged-PR
advancement behave identically, INCLUDING the merged-PR-while-mod-blocked finalize
(the relocated mod-block clear). *Proof:* the existing behavioral suites stay green
with no edits — `internal/dispatch` (`TestBuildModelPrecedence` and the build/reuse
hazards suite, the code that actually resolves `effective_model`), `internal/status`
(the merge-policy/terminal-guard suite, which pins that a standalone `mod-block=`
clear followed by a terminalize passes and a combined clear+terminalize still
refuses — exactly the two-step shape the relocated mod prose now follows), and the
`internal/contractlint` closure/ceremony suite (`TestBootResidentDeferredLoadPointsResolve`,
`TestHostNeutralCoresResolveAndCarryCeremony` — the event loop's `## Event Loop`
anchor and the cores' reference closure survive the step deletion). Green-after =
proof the move touched no behavioral seam: the mod-block-clear mechanism the
terminalize guard depends on is unchanged Go code; only WHICH instruction issues
the clear moved (core step-1 → the mod's own advancement). This AC owns no NEW
test; it is the gate's cross-check that the existing seams stayed green.

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
legitimate terminalize shape `completed verdict={v} worktree=`), with a behavioral
golden proving a non-terminal `--set worktree=` under a registered merge hook now
passes (exit 0, worktree cleared).

**CONSTRAINT carried into the split-out spec (must not regress the combined-clear
guard).** The fix MUST NOT simply delete `handlers.go:116-118`. `fo-merge-core.md`
(lines 19 and 59) makes `worktree=` a terminal field WHEN COMBINED with
`mod-block=` / other terminal fields — the combined-clear refusal lists `worktree=`
explicitly. In code, the combined-clear guard (`handlers.go:156-169`) reads
`isTerminalUpdate()`: when `clearingModBlock` is true AND `isTerminalUpdate()` is
true, it refuses. A `mod-block= worktree=` combined call relies on the
`worktree`-is-terminal branch (116-118) to make `isTerminalUpdate()` true; deleting
that branch outright would let a combined `mod-block= worktree=` clear slip past the
combined-clear refusal — a regression the merge-policy/terminal-guard suite must
catch. So the fix must keep `worktree=` terminal WHEN the same `--set` also carries
`mod-block=` (or `completed`/`verdict`/terminal-`status`), and only stop
false-tripping on a STANDALONE `worktree=` clear. The split-out task's test plan
must include a regression golden: `--set {slug} mod-block= worktree=` (combined)
STILL refuses (exit 1, combined-clear reason). Filed as a sibling, not folded here.

## Staff review

**Staff review completed (cycle 1) — approach confirmed sound; four gaps closed in
cycle 2 (see the cycle-2 stage report).** This is a shipped-contract + dispatch-core
change (`fo-dispatch-core.md` is the host-neutral dispatch core every host loads;
`claude-fo-dispatch.md` and `mods/pr-merge.md` ship). The review confirmed the proof
pattern is legitimate, the behavior-neutral claim holds, and the split-out is right;
it raised four gaps, all now resolved: (1) Move B relocates the mod-block clear into
the shipped pr-merge template's startup + idle advancement (the shipped template
never cleared it; core step-1 did); (2) AC-3/AC-4 now mandate their discriminator
controls explicitly; (3) the AC-2 presence-grep is dropped; (4) the split-out spec
carries the combined-clear non-regression constraint. The original review rubric —
(a) ACs are genuine structural checks not prose-greps (the package's own
`doc_test.go` policy is the rubric); (b) no code reads the moved prose; (c) the
split-out is right; (d) the before/after wording preserves the re-stamp rule and the
merged-PR advancement paths verbatim — is satisfied by the reworked body.

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

## Stage Report: ideation (cycle 2)

Closed the four staff-review gaps before the gate.

- DONE: Gap 1 (MATERIAL — Move B mod-block-clear relocation).
  Confirmed via grep that the SHIPPED template `./mods/pr-merge.md` has `mod-block` count ZERO (relies on core step-1's clear) while the dev local `docs/dev/_mods/pr-merge.md` (count 11) does the two-step clear-then-terminalize. Reworked Move B: added the standalone `mod-block=` clear (separate `--set`, per `fo-merge-core.md:19`) to the shipped template's startup AND idle merged-PR advancement, with verbatim before/after. AC-5 now covers the merged-PR-while-mod-blocked finalize.
- DONE: Gap 2 (ship the AC-3/AC-4 discriminator controls).
  AC-3 now mandates `TestEventLoopPRScanScannerDiscriminates` (flags planted PR-scan line, passes idle-hook line). AC-4 now mandates two controls modeled on `TestPortabilityCheckDiscriminatesHostSpecific`: a positive control (pr-merge.md DOES carry `gh pr view`, allow-list load-bearing) and a negative control (planted non-allowed file is flagged).
- DONE: Gap 3 (drop AC-2.2 presence-grep).
  AC-2 rewritten to rely on AC-1 (absence from core) + the behavioral proof (`TestBuildModelPrecedence` + context-budget family rule) only. The presence-grep is explicitly dropped as the banned prose-grep.
- DONE: Gap 4 (combined-clear constraint into the split-out spec).
  Verified `fo-merge-core.md:19/59` list `worktree=` as a combined-clear terminal field, and `handlers.go:156-169` reads `isTerminalUpdate()` for the combined-clear refusal. Added a CONSTRAINT to the `status-set-worktree-clear-guard` spec: do NOT delete `handlers.go:116-118` outright — keep `worktree=` terminal when combined with `mod-block=`/terminal fields, stop false-tripping only on a STANDALONE clear, with a regression golden (`--set mod-block= worktree=` still refuses).

### Summary (cycle 2)

All four gaps closed. The material one (gap 1) was a real behavior-correctness miss
on my part: core step-1 clears `mod-block` AND the shipped pr-merge template never
did, so deleting step-1 without relocating the clear would have blocked
shipped-template merged-PR finalizes. Move B now relocates the two-step clear into
the shipped template's startup + idle advancement, grounded in the dev local copy
as reference. The discriminator controls are now mandated explicitly (not just
promised), the only presence-grep is dropped, and the split-out spec carries the
combined-clear non-regression constraint with concrete code/contract citations.

## Stage Report: implementation

- DONE: Move A: relocate the Claude model tokens (opus[1m] + the sonnet/opus/haiku enum) OUT of fo-dispatch-core.md (BOTH sites: reuse-condition-4 AND the Break-Glass template) INTO claude-fo-dispatch.md, genericizing the core to 'the host's canonical model enum'; the re-stamp behavior is unchanged. AC-1 (zero model tokens in core, with the ScannerDiscriminates control) + AC-2 (behavioral: TestBuildModelPrecedence + context-budget unchanged; the AC-2.2 presence-grep stays DROPPED) prove it.
  Core grep `opus|sonnet|haiku` now empty; both sites moved (reuse-cond-4 → Context-Budget paragraph, Break-Glass → Break-Glass note). AC-1 + ScannerDiscriminates green; TestBuildModelPrecedence + context-budget family-rule tests green and cached-unchanged (commit 16697bf3).
- DONE: Move B: delete the gh-pr-view PR scan from generic event-loop step 1, AND relocate the standalone mod-block-clear into the SHIPPED template ./mods/pr-merge.md's startup AND idle merged-PR advancement (separate --set, per fo-merge-core.md:19, with explicit status={terminal} for the generic case). Edit ONLY the shipped ./mods/pr-merge.md — leave docs/dev/_mods/pr-merge.md (the dev-local count-11 copy) UNTOUCHED; keep the literal gh pr view in the idle hook (AC-4's positive control depends on it).
  Step 1 deleted, 2-4 renumbered to 1-3; two-step `mod-block=` clear-then-terminalize added to shipped pr-merge.md startup + idle. `git diff` confirms docs/dev/_mods/pr-merge.md untouched; idle hook retains `gh pr view`. `gh pr view` now appears only in mods/pr-merge.md across skills/+mods/.
- DONE: Ship the discriminator controls, not just the absence halves: AC-3's TestEventLoopPRScanScannerDiscriminates (flags a planted PR-scan line, passes the idle-hook line); AC-4's positive control (pr-merge.md DOES carry gh pr view, allow-list load-bearing) + negative control (planted non-allowed file flagged); AC-5: the existing dispatch/status/contractlint suites stay green (no behavioral edits). Do NOT fold in the worktree= guard sibling — it splits out to status-set-worktree-clear-guard.
  All 7 new tests in internal/contractlint/layering_restore_test.go green (3 absence + 4 discriminator/control). AC-5: full `go test ./...` = 15/15 packages ok, dispatch+status cached-unchanged. Sibling worktree= guard NOT touched (split-out preserved).

### Summary

Two behavior-neutral text moves restoring leaked boundaries, plus contractlint
structural-absence guards with discriminator controls (TDD: tests written + confirmed
RED before the doc edits, then GREEN after). Move A genericized both model-enum leak
sites in fo-dispatch-core.md to "the host's canonical model enum" and relocated the
concrete `sonnet|opus|haiku` + `opus[1m]` tokens to claude-fo-dispatch.md. Move B
deleted the `gh pr view` scan from generic event-loop step 1 (renumber 2-4→1-3) and
relocated the two-step `mod-block=` clear-then-terminalize into the shipped
pr-merge.md startup + idle hooks. NOTE: the spec's pre-impl grep claimed the only
stale step-number cross-reference was the positional "step 4's re-run" phrase; I found
three additional bare "step 4" tokens in claude-fo-dispatch.md (lines 159/161/173,
now factually stale after the renumber) and updated all to step 3 — a faithful
execution of the renumber intent. No behavioral seam touched: dispatch+status suites
cached-unchanged; full `go test ./...` 15/15 green. Code committed to the worktree
(16697bf3); sibling worktree= guard left to its split-out task.

### Feedback Cycles

**Cycle 1 — detached adversarial audit (validation, 2026-06-19).** Audit run on a
THROWAWAY detached worktree at the impl SHA (16697bf3), never the impl worktree;
removed after. Eight claim-breaking edits; six caught, two findings recorded. Both
findings are non-blocking (the deliverable matches its ratified spec) but on the
record per the detached-audit instruction.

- *Caught (load-bearing confirmed):* re-planted model enum in core →
  `TestDispatchCoreHasNoClaudeModelToken` FAIL; re-planted `gh pr view` in core →
  `TestEventLoopCoreHasNoPRScan` + `TestNoUnexpectedPRViewScanIntroduced` FAIL;
  vacuous model-token scanner → `...ModelTokenScannerDiscriminates` FAIL; vacuous
  PR-scan scanner → `...PRScanScannerDiscriminates` FAIL; `gh pr view` planted in a
  REAL non-allowed shipped file (fo-merge-core.md) → `TestNoUnexpectedPRViewScanIntroduced`
  FAIL; `gh pr view` stripped from the mod → `TestPRViewAllowListIsLoadBearing` FAIL.
- *Finding A (no guard on the relocated mod-block clear — NON-BLOCKING).* Removing
  the two-step `mod-block=` clear from the shipped `mods/pr-merge.md` startup hook
  (reverting Move B's relocation) leaves the FULL `go test ./...` suite GREEN — nothing
  catches it. The `internal/status` terminal-guard suite pins the Go MECHANISM
  (`TestMergeLocalModBlockPendingStillBlocks`, `TestMergeLocalCombinedClearAndTerminalizeRefused`)
  but never reads the mod prose, so it cannot prove the shipped mod actually ISSUES the
  clear. Judged non-blocking: AC-5 scopes its claim narrowly and correctly to "the Go
  mechanism is unchanged," and staff review RATIFIED dropping the presence-grep over
  instruction prose (AC-2) as "the banned prose-grep." A presence-grep asserting the mod
  contains `mod-block=` is exactly that tautology. The relocation is verified by reading
  the diff (present in both startup + idle hooks). This is a principled proof-policy
  limit, not a defect — but a future deletion of the relocated clear would regress the
  merged-PR-while-mod-blocked finalize silently. Surfaced for the captain's awareness.
- *Finding B (AC-4 negative control is a tautology — NON-BLOCKING).*
  `TestPRViewAllowListConstrains` hardcodes `plantedRel` + `plantedContent` and
  re-implements the allow-list check inline against in-memory constants; it never reads
  disk. When the REAL scan's allow-list flag was defeated (`if false && …`), the negative
  control still PASSED — it cannot detect the real scan being weakened. Its expected value
  is the text the author wrote (the validation "Bad" clause: a static check counts only
  against an independent source that can diverge). Judged non-blocking because AC-4's
  REAL load-bearing proof is the other two trio members (`TestNoUnexpectedPRViewScanIntroduced`
  catches a real planted file; `TestPRViewAllowListIsLoadBearing` reds when the mod stops
  carrying the token) — both proven genuine above. The negative control adds nothing but
  weakens nothing; recommend either deleting it or rewriting it to plant a real temp file
  and drive the actual scan. Recorded, not gated.

## Stage Report: validation

- DONE: Detached adversarial audit on a THROWAWAY checkout — construct claim-breaking edits (re-introduce a model token into fo-dispatch-core.md, re-introduce gh pr view into the core, remove the shipped-template mod-block-block clear, make a discriminator control vacuous) and confirm a test catches EACH; refute that the absence checks AND discriminator controls are load-bearing, not prose-grep tautologies.
  Throwaway detached worktree at 16697bf3 (removed after). Mutations 1/2/4a/4b/4c/4d each caught by the named test (see `### Feedback Cycles` cycle 1). AC-1 absence check + discriminator share `lineLeaksClaudeModelToken` (genuine coupling). Two findings recorded: (A) the relocated mod-block clear has no guard — mutation 3 stays green (NON-BLOCKING, ratified proof-policy limit); (B) the AC-4 negative control is a tautology (NON-BLOCKING, AC-4's real weight is the other two trio members).
- DONE: Behavior-neutrality — TestBuildModelPrecedence + the context-budget family-rule + the dispatch/status/contractlint closure-ceremony suites stay green WITHOUT edits (AC-5); confirm the merged-PR-while-mod-blocked finalize works (shipped mods/pr-merge.md relocated mod-block-clear is a STANDALONE --set separate from the terminalize).
  Full `go test ./...` = 15/15 packages ok on clean throwaway HEAD. `TestBuildModelPrecedence` (sonnet/opus/haiku/null) green; context-budget family rule (`TestContextBudget1MFalseNegativeGone`) green; `TestHostNeutralCoresResolveAndCarryCeremony` (pins the `## Event Loop` anchor) green. Shipped pr-merge.md startup + idle each issue a standalone `mod-block=` clear then a separate terminalize `--set` (per fo-merge-core.md:19); the status terminal-guard suite confirms a combined clear+terminalize still refuses.
- DONE: AC-4 end-state — a no-pr-merge-mod workflow has a PR-free generic loop; gh pr view appears ONLY in mods/pr-merge.md across skills/+mods/ (allow-list load-bearing via the positive control); confirm the step 2-4→1-3 renumber left NO stale cross-references; the worktree=-guard sibling stays SPLIT-OUT.
  `grep -rn "gh pr view" skills/ mods/` → only mods/pr-merge.md:13,29. No stale event-loop step cross-refs remain (the four "step N" matches are the `## Dispatch`/merge-core/startup procedures, not the event loop). Sibling `handlers.go:116-118` unchanged and absent from the commit diff — split out, not folded.

### Summary

PASSED. Ran the detached adversarial audit on a throwaway checkout per the
high-stakes-surface proof policy. Six of eight claim-breaking edits were caught by
the named test, confirming the absence checks and the material discriminator/positive
controls are load-bearing (not prose-grep tautologies); both layering leaks are
genuinely closed (no model tokens and no `gh pr view` in the host-neutral core; `gh pr
view` only in mods/pr-merge.md). Behavior is provably unchanged: full `go test ./...`
15/15 green without edits, `TestBuildModelPrecedence` + the context-budget family rule
lock the code-resolved enum/mapping, and the relocated two-step mod-block clear matches
the terminal-guard mechanism the status suite pins. Two NON-BLOCKING findings recorded
in `### Feedback Cycles`: (A) the relocated mod-block clear has no test guard (a
ratified proof-policy limit — a presence-grep over the mod prose would be the banned
tautology; the relocation is diff-verified), and (B) the AC-4 negative control is a
tautology that passed even when the real scan was defeated (harmless — AC-4's real
weight is the other two trio members, both proven genuine). Neither rejects the gate;
both are on the record for the captain. Sibling worktree=-guard correctly split out.

## Stage Report: implementation (cycle 2)

Closed validation finding B; finding A accepted as a ratified proof-policy limit (no change).

- DONE: Fix finding B — make the AC-4 negative control `TestPRViewAllowListConstrains` load-bearing instead of a tautology.
  Factored the real allow-list scan into a shared helper `prViewLeaksOutsideAllowList` that BOTH the production check (`TestNoUnexpectedPRViewScanIntroduced`) and the negative control drive — no more inlined re-implementation. The control now plants a real non-allowed shipped file (`skills/first-officer/references/zz-pr-view-negative-control.md`, written + cleaned up via defer), runs the REAL scan over the actual shipped-surface walk, and asserts the scan flags the plant (commit 95b24e4d).
- DONE: Prove non-vacuity by mutation and record the outcome.
  MUTATION: stubbed the shared `prViewLeaksOutsideAllowList` to never flag (`if false && …`) → `TestPRViewAllowListConstrains` went RED via its scan-flagging assertion (`leaks=[]`, layering_restore_test.go:231), not via a setup guard. Restored the helper → GREEN. This proves the control reds when the real scan logic is defeated, the discipline the absence-check discriminators already follow.
- ACCEPTED (no change): finding A — the relocated mod-block clear in shipped `mods/pr-merge.md` has no test guard.
  Per the team-lead's direction: ratified proof-policy limit. A presence-grep over the mod prose would be the banned tautology; the relocation is diff-verified against the dev-local reference and gets incidental live-e2e coverage on the merged-PR path. Left as recorded in `### Feedback Cycles`; no presence-grep added.

### Summary (cycle 2)

The cycle-1 negative control was the exact defect the PR is about — a discriminator that
can't fail. Fixed by extracting the allow-list scan into one shared helper both callers
drive, then planting a real non-allowed file and asserting the REAL scan flags it; proven
load-bearing by a mutation that stubs the helper and turns the control RED via its scan
assertion. Finding A left untouched per the ratified proof policy. `go test ./...` 15/15
green; `go vet` clean. Committed to the worktree branch (95b24e4d).
