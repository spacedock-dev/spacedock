---
title: Make the shipped state/merge verbs the operative contract path + bind «fn» bodies to their verbs (oracle)
source: '0221-layered-fo rework (2026-06-19): validated findings 1+2 — the contract does not operatively USE the shipped verbs. `### Split-Root State Sync` names an abstract "status tool" not `spacedock state commit`, and `«state.commit»`s effect restates the hand git sequence; `«merge.guard»`s prose claims it "invoke[s] the registered merge hook" / runs "as one call" / "default-merge[s]", all FALSE of the shipped re-entrant partial envelope (merge.go doc-comment: "It does NOT invoke the merge hook"). Bundles the contract rewire (A) with a routing oracle (B) written test-first.'
status: ideation
score: 0.75
sprint: 0221-layered-fo
group: foundation
id: 4asxw7kxvdzdtf87w9rjkxwx
started: 2026-06-19T22:53:22Z
---

Complete the vertical slice `rgq`/`mz`/`czw` shipped empty: the verbs exist and route, but the contract instructs the hand sequence and mis-describes the merge verb. Make the verb the operative path AND guard it mechanically. (Per the captain's call, A and B are ONE bundled slice, not two layer-tasks.)

## Premise reconciliation (resolved in ideation — the seed was partly stale)

Three findings from reading main reshape the seed; the captain should see them before the gate.

**F1 — czw (#402) ALREADY shipped the routing oracle.** `internal/cli/prose_function_routing_test.go::TestProseFunctionNotationBindsToRouting` walks the three FO cores, extracts every `→ **shipped**`/`→ **prose**` migration target, invokes the named verb's argv against the compiled root, and asserts shipped verbs route while prose targets are rejected. It ships with two paired controls: `TestProseFunctionRoutingGuardFailsOnViolation` (RED control — plants `spacedock merge teleport`, an un-routed backtick, and a guillemet whose target *does* route) and `TestProseFunctionRoutingOracleDiscriminates` (non-vacuity). All GREEN on main today (`go test ./internal/cli -run TestProseFunctionNotationBindsToRouting` → ok). **Task-B as a NEW from-scratch routing oracle is already built.** What is NOT yet covered is the «state.commit» binding being a backtick-shipped target — see A2 below; once A1/A2 land, the existing oracle exercises it for free.

**F2 — the routing oracle is structurally BLIND to the «merge.guard» prose defect.** The seed says the oracle should be "RED on today's overstated «merge.guard»." It cannot be. The oracle reads only the `→ **shipped**: spacedock merge guard` line and checks the *verb routes* — it never reads the effect-bullet prose. The overstated claims ("invoke the registered merge hook", "default-merge", "as one call") live in lines 7 and 11, which the oracle does not parse. So "RED on the overstated prose" is mechanically incoherent for a routing oracle. The seed's own stated ALTERNATIVE — "a mutation pointing an «fn» at a non-existent verb" — IS the achievable non-vacuous RED, and it is ALREADY the shipped `TestProseFunctionRoutingGuardFailsOnViolation` (plants `merge teleport`). **The prose-accuracy defect (A3) and the routing oracle (B) are orthogonal mechanisms; the oracle cannot guard A3.** This slice does NOT manufacture a contrived prose-grep to fake a RED; it ships the prose fix (A3) and relies on the *already-present, already-non-vacuous* routing controls.

**F3 — `deferred_tier_absence_test.go` is NOT on main.** It exists only in the unmerged sibling branch `spacedock-ensign/strip-deferred-tier-vocabulary` (which ADDs it, +113 lines, branched off c25fee26). "B REPLACES the deferred-tier absence-grep" is a false coupling: there is nothing on main to replace, and the routing oracle that would have "replaced" it already shipped independently in #402. **Resolution (checklist item 3): B does NOT delete `deferred_tier_absence_test.go`.** The two tests are unrelated — the deferred-tier grep guards member-72 tier vocabulary out of the cores (a different concern entirely); the routing oracle guards «fn»↔verb routing. Neither subsumes the other. The "replaces" language in the seed is retired here as a mis-mapping; if the captain still wants the deferred-tier grep retired, that is the strip-deferred-tier branch's own scope, not this slice's.

## Proposed approach (firm)

**A — rewire the contract prose (dispatched worker on the implementation worktree — scaffolding guardrail under `skills/first-officer/references/`, NOT an FO live edit):**

The three rewrites, with exact before/after wording, are in the **Concrete rewrites** section below. In summary:
- **A1** — `### Split-Root State Sync` "Preferred" bullet (`first-officer-shared-core.md:189`): name `spacedock state commit <slug>` as the preferred path; demote the hand `git -C add/commit` to the named degraded/no-origin fallback.
- **A2** — `«state.commit»` effect bullet (`first-officer-shared-core.md:178`): delegate to the verb; stop restating the full hand git sequence as the effect (keep the hand sequence only as the documented fallback the verb automates).
- **A3** — `«merge.guard»` intro (line 7) + effect bullet (line 11) in `fo-merge-core.md`: replace the three false claims with the verb's actual re-entrant-partial-envelope behavior, sourced from `internal/status/merge.go` (doc-comment line 13: "It does NOT invoke the merge hook"; `signalArmed` line 213: "invoke the … merge hook, then re-run `merge guard`").

**B — the «fn»→shipped-verb routing oracle:** already shipped (#402). This slice's B-obligation is reduced to: (1) confirm the existing oracle stays GREEN after the A rewrites land (the «state.commit» backtick must keep routing; no «fn» migration-target line is changed by A — A touches effect/intro prose, not the `→` lines), and (2) confirm the existing RED control (`TestProseFunctionRoutingGuardFailsOnViolation`) and discriminator remain the non-vacuity proof. No new oracle test is authored — authoring a duplicate would violate YAGNI and the no-duplication rule.

## Concrete rewrites (before/after — load-bearing exact wording)

### A1 — Split-Root State Sync "Preferred" bullet (`first-officer-shared-core.md`, currently line 189)

BEFORE:
```
- **Preferred — tool-managed atomic state commits.** When the status tool owns `add`+`commit` under a lock, route through it.
- **Fallback — path-scoped commits per writer.** `git -C {state_checkout} add {entity_path} && git -C {state_checkout} commit -m "…" -- {entity_path}`. Never a bare `git add -A` or bare `git commit`. Retry on `index.lock` contention after ~2s.
```
AFTER:
```
- **Preferred — `spacedock state commit <slug>`.** The verb resolves the entity's path under the state checkout and runs the path-scoped commit → push → on-reject `pull --rebase` → re-push sequence (and the rebase-conflict halt) atomically. Invoke it directly; do not hand-roll the git sequence.
- **Fallback (no verb available / degraded) — path-scoped commits per writer.** `git -C {state_checkout} add {entity_path} && git -C {state_checkout} commit -m "…" -- {entity_path}`. Never a bare `git add -A` or bare `git commit`. Retry on `index.lock` contention after ~2s.
```
Rationale: the verb (`internal/cli/state_sync.go::runStateCommit`) already implements exactly the commit+sync+halt the prose described abstractly. The hand sequence stays as the named fallback (it is genuinely what the verb automates and what a no-`spacedock`-binary degraded FO falls back to).

### A2 — «state.commit» effect bullet (`first-officer-shared-core.md`, currently line 178)

BEFORE:
```
- **effect:** path-scoped commit + sync per **Split-Root State Sync** below — `git -C {state_checkout} add {entity_path} && git -C {state_checkout} commit -m "…" -- {entity_path}`, never a bare `git add -A` / `git commit`; then push; on push rejection `pull --rebase` then re-push; retry on `index.lock` contention after ~2s.
```
AFTER:
```
- **effect:** delegate to `spacedock state commit <slug>` — it path-scoped-commits the entity and runs the full sync (push; on push rejection `pull --rebase` then re-push; rebase-conflict halt) per **Split-Root State Sync** below. The hand `git -C add/commit … push/pull --rebase` sequence is the documented fallback the verb automates, not the FO's first move.
```
Note: the `→ **shipped**: spacedock state commit <slug>` migration line (currently line 181) is UNCHANGED — it already names the verb correctly, which is why the routing oracle is already GREEN on this binding. A2 changes only the effect bullet so the *delegation* is the stated behavior, not the hand sequence.

### A3 — «merge.guard» intro + effect bullet (`fo-merge-core.md`)

BEFORE (intro, line 7):
```
When an entity reaches its terminal stage, `«merge.guard»(slug)` runs the atomic merge-finalize ceremony — the mod-block set→invoke→clear→terminalize sequence below, as one call.
```
AFTER (intro):
```
When an entity reaches its terminal stage, `«merge.guard»(slug)` drives the terminal merge-finalize ceremony as a re-entrant partial envelope: it arms the mod-block and signals the FO to invoke the merge hook, then on re-run detects completion by state delta, clears the mod-block standalone, terminalizes, and archives. It does NOT itself invoke the hook, and the clear+terminalize is two separate `--set` calls, not one.
```
BEFORE (effect bullet, line 11, first clause through "and archive"):
```
- **effect:** run the terminal merge-finalize ceremony for the entity — set `mod-block=merge:{mod_name}`, invoke the registered merge hook, detect completion, clear the mod-block in its own call, default-merge if no hook handled it, terminalize (`completed verdict={verdict} worktree=`), and archive.
```
AFTER (effect bullet, first clause through "and archive"):
```
- **effect:** drive the terminal merge-finalize ceremony for the entity as a re-entrant envelope — **arm:** set `mod-block=merge:{mod_name}` and signal the FO to invoke the registered merge hook (the verb does NOT invoke it); the FO invokes the hook, then re-runs the verb. **finalize (on re-run):** detect hook completion by the `pr`/`mod-block`/`verdict` state delta, clear the mod-block in its own standalone `--set` call, terminalize (`completed verdict={verdict} worktree=`), and archive. (The default local merge — when no hook handled it — and worktree-removal/worker-teardown, steps 9–10, are the FO's, not the verb's.)
```
The remaining clause of the effect bullet ("The mechanism-level enforcement in **Mod-Block Enforcement** below guards every step …") is UNCHANGED — it is accurate. Steps 1–10 and the `→ **shipped**` line (line 14) are UNCHANGED.

Three false claims removed, with their ground truth:
1. "invoke the registered merge hook" (the verb does the inviting) → `merge.go:13` "It does NOT invoke the merge hook"; `merge.go:213` arm signal "invoke the … merge hook, then re-run `merge guard`" — the FO invokes, the verb signals.
2. "as one call" (intro) → the verb is re-entrant (arm, then a separate finalize re-run) AND the clear is standalone from terminalize (`merge.go::finalize` emits two `emitSet` calls; step 5 line 28 "The clear MUST be standalone").
3. "default-merge if no hook handled it" attributed to the verb → `merge.go::MergeGuard` never local-merges; `finalize` only clears+terminalizes+archives. The default local merge is the FO's (step 6).

## Acceptance criteria

- **AC-1 (state.commit primary path).** An FO reading `### Split-Root State Sync` "Preferred" bullet and the «state.commit» effect bullet is directed to `spacedock state commit <slug>` as the primary path, with hand-git named explicitly as the fallback. **Verified by:** the A1+A2 rewrite landing the AFTER text verbatim (diff in `fo`-references on the impl worktree) AND `TestProseFunctionNotationBindsToRouting` GREEN on the unchanged `→ **shipped**: spacedock state commit` binding (`go test ./internal/cli -run TestProseFunctionNotationBindsToRouting`).
- **AC-2 («merge.guard» prose matches shipped behavior).** The «merge.guard» intro + effect bullet describe the verb as a re-entrant partial envelope (arm → FO invokes hook → re-run → finalize; the verb does NOT invoke the hook, does NOT local-merge, clear-then-terminalize is two `--set`s). **Verified by:** the A3 rewrite landing the AFTER text verbatim AND a grep confirming the three false strings ("invoke the registered merge hook", "as one call", "default-merge if no hook handled it") are absent from `fo-merge-core.md` while the three true strings ("does NOT … invoke the hook", re-entrant/re-run framing, standalone clear) are present — a substring presence/absence check is acceptable proof here because the AC is *literally about specific prose tokens*, not behavior (the behavior is merge.go's, already tested in `internal/status/merge_guard_test.go`).
- **AC-3 (routing oracle stays the non-vacuous guard, no new test).** The shipped `TestProseFunctionNotationBindsToRouting` + its RED control `TestProseFunctionRoutingGuardFailsOnViolation` + discriminator `TestProseFunctionRoutingOracleDiscriminates` remain GREEN after the A rewrites (the rewrites touch no `→` migration line). **Verified by:** `go test ./internal/cli -run 'TestProseFunction'` GREEN post-rewrite; the RED control's planted `merge teleport`/routed-guillemet violations are its non-vacuity proof (already authored). **No `deferred_tier_absence_test.go` is deleted** (it is not on main; see F3).

## Test plan

- **No new Go test is authored.** The routing oracle and its non-vacuity controls already exist and pass (#402). Authoring a duplicate violates the no-duplication rule. The slice's mechanical proof is: (a) the existing `internal/cli` oracle suite GREEN after the rewrite, and (b) a prose presence/absence grep over `fo-merge-core.md` and `first-officer-shared-core.md` for AC-2's specific tokens.
- **Cost/complexity:** LOW. Three prose rewrites on the impl worktree + two `go test -run` invocations + one grep. No new mechanism, no spike.
- **Fixture/CLI/live:** CLI-level only (`go test ./internal/cli -run TestProseFunction`). No fixture, no live workflow drive. The behavior of `spacedock merge guard` itself is already covered by `internal/status/merge_guard_test.go` and `internal/cli/merge_test.go` (shipped #400) — this slice does not re-test the verb; it makes the *prose* describe it truthfully.
- **Pre-impl ordering:** run `go test ./internal/cli -run TestProseFunctionNotationBindsToRouting` BEFORE the rewrite to confirm GREEN baseline, then after to confirm no regression. The riskiest thing the rewrite could break is a `→` migration line; the rewrites are scoped to effect/intro prose precisely to avoid that.

## Spike determination

No spike needed. The proven mechanisms this slice relies on: (1) `internal/cli/prose_function_routing_test.go` routing oracle — exercised GREEN on main today; (2) `internal/status/merge.go` shipped behavior (re-entrant arm/finalize, no-hook-invoke, standalone clear) — read directly from source and confirmed against the doc-comment and signal strings; (3) `internal/cli/state_sync.go::runStateCommit` — confirmed to implement the commit+sync+halt the A1/A2 prose will delegate to. All three are shipped code on main, not unverified mechanisms.

## Live-usage proof

The end-value proof that a live FO actually CALLS the verbs as its path rides `kt` (haiku-drive-validation, id `kt96jb8yagkean75j0k1ep6n`, sprint-readiness=DEFER, sequenced LAST as the sprint gate proof). This slice ships the prose rewire (A) and confirms the mechanical oracle (B, already shipped); the live drive is `kt`'s. `kt` is blocked on all four verb-core members + 72 + 6re merged AND m4's tmux backend green — this slice is one of the upstream rewires it waits on.

## Worktree / split-root note

A is a `skills/first-officer/references/` edit (scaffolding guardrail), so implementation runs WITH a worktree and the rewrites land on the impl branch — this is NOT an FO live-contract edit. This ideation entity + report live in the split-root state checkout (`docs/dev/.spacedock-state/`), committed path-scoped.

## Stage Report: ideation

- DONE: Concrete before/after wording for all three rewrites (Split-Root State Sync Preferred bullet; «state.commit» effect bullet; «merge.guard» body) recorded in the entity body.
  Recorded as A1/A2/A3 in the "Concrete rewrites" section with verbatim BEFORE/AFTER blocks; exact source lines (first-officer-shared-core.md:189, :178; fo-merge-core.md:7, :11) and ground-truth citations from merge.go (:13, :213).
- DONE: The routing oracle's assertion is designed plus its non-vacuous RED control (the pre-rewrite overstated «merge.guard» prose, or a non-existent-verb mutation) named.
  Designed as a reconciliation, not a new build: the routing oracle already shipped in #402 (`TestProseFunctionNotationBindsToRouting`); its non-vacuous RED control is the already-shipped `TestProseFunctionRoutingGuardFailsOnViolation` (plants `spacedock merge teleport`, a non-existent-verb mutation). F2 documents WHY "RED on overstated prose" is mechanically impossible for a routing oracle (it never reads effect prose) and routes that proof to the AC-2 token grep instead.
- DONE: Confirmed: B replaces deferred_tier_absence_test.go, and the rewrites are reconciled with czw's restructured contract regions (the byte-preservation tension resolved per-region).
  Confirmed the OPPOSITE of the seed and recorded it (F1/F3): `deferred_tier_absence_test.go` is NOT on main (only in unmerged `strip-deferred-tier-vocabulary`), and the routing oracle that would "replace" it already shipped independently — so B does NOT delete it; the "replaces" coupling is retired as a mis-mapping. Per-region reconciliation: A2 and A3 are scoped to effect/intro prose and leave every `→` migration line and the steps 1–10 UNCHANGED, so czw's routing-bearing regions are byte-preserved and the existing oracle stays GREEN.

### Summary

Read main directly and found the seed's premise was partly stale: czw's #402 already shipped the routing oracle (B) with its non-vacuity controls, `deferred_tier_absence_test.go` is not on main (so nothing to "replace"), and a routing oracle is structurally blind to the «merge.guard» prose defect it was supposed to catch. Resolved by reducing B to "confirm the shipped oracle stays GREEN", routing the prose-accuracy proof to an AC-2 token grep (the three false strings — "invoke the registered merge hook", "as one call", "default-merge if no hook handled it" — confirmed present today), and recording all three before/after rewrites verbatim. No spike needed; all relied-on mechanisms are shipped code exercised against source. This is a low-cost prose-rewire slice with no new Go test.
