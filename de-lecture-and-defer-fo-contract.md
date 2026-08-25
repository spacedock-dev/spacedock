---
id: wn3dg7txnrte0jrcxf56b859
title: De-lecture the FO contract and defer non-boot shared-core sections (folds s6q)
status: implementation
source: "Captain CL, 2026-08-24: contract audit ruling after the #757 fo-install overbuilt review - 'file and fold s6q. dispatch. keep it light in ideation as this is clear about what to cut'; folds defer-shared-core-non-boot-sections (s6qamkh7efky9zh5jh6ba6xq, superseded into this task)"
started: 2026-08-24T23:37:58Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-de-lecture-and-defer-fo-contract
issue:
gates:
    version: 1
    records:
        - id: gate:wn3dg7txnrte0jrcxf56b859:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:wn3dg7txnrte0jrcxf56b859-backlog-1
              briefing:
                id: briefing:wn3dg7txnrte0jrcxf56b859:backlog:attempt-1:revision-1
                digest: sha256:b879d0be71511f84291a906609ac6f12c9d70175e7665971caf86203204700c5
                request-digest: sha256:e0160113b31146057ccfa33cf4ef0a3ece370f9b5cc84ad1c782dd9516f8467a
                room-ref: ./de-lecture-and-defer-fo-contract/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:wn3dg7txnrte0jrcxf56b859:backlog:1
                briefing: briefing:wn3dg7txnrte0jrcxf56b859:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-24T23:37:42.284397Z"
                decision: approve
                reason: 'Captain CL in chat 2026-08-24: ''file and fold s6q. dispatch. keep it light in ideation as this is clear about what to cut and we should have stacked PR of the two to run live CI'' - accepts the seed with light-ideation and stacked-delivery directives'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:wn3dg7txnrte0jrcxf56b859:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:wn3dg7txnrte0jrcxf56b859-ideation-1
              briefing:
                id: briefing:wn3dg7txnrte0jrcxf56b859:ideation:attempt-1:revision-1
                digest: sha256:f8d70f5565f1d662a98de3b20c6e7a311a00e338e66a5a115a75ad528903c3d4
                request-digest: sha256:80b53391d0b55fa47c67da34fd6501e30dafccca843fe6e85832dc8dc4db6590
                room-ref: ./de-lecture-and-defer-fo-contract/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:wn3dg7txnrte0jrcxf56b859:ideation:1
                briefing: briefing:wn3dg7txnrte0jrcxf56b859:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-25T00:42:16.785549Z"
                decision: approve
                reason: 'Captain CL in chat 2026-08-24: ''dispatch. keep the validation lean. i''ll review on github'' at ideation attempt-1 (digest f8d70f55) - accepts the design including the cap-raise-as-follow-up recommendation; sets a lean validation posture with captain GitHub review as the human check'
              application:
                target-stage: implementation
                state: consumed
---

Two movements over the FO contract with one lens: the contract carries rules and commands; rationale, narration, design history, and repeated justification live in entity bodies and tests. Movement 1 (delete): the 2026-08-24 FO audit enumerated the lecture, file by file. Movement 2 (defer, folded from s6q): move boot-unreachable shared-core sections behind their existing triggers, measured by boot occupancy in tokens (s6q's archived body carries the full measured profile: shared core ~10.1k tok of a 40.4k greet boot; candidate span 9,105 of 26,092 bytes = 35%).

The audit cut list (chat, 2026-08-24, all files read in full):
- first-officer-shared-core.md: the "future multi-workflow form" design-history sentence in engage's scope; the stopping-point paragraph's re-argument (~40% compressible, exception list is the content); "Dense durable boundaries keep unhinted auto-compaction survivable"; Working Principles justification tails ("...is the asymmetry that lets a means-accurate, end-missed stage pass", "a scheme minted in a dispatch prompt becomes every downstream artifact's vocabulary" second clause, third restatements in the smallest-mechanism bullet). KEEP the banned-excuse enumerations - prohibitions are behavior-shaping, their justifications are the lecture.
- fo-dispatch-core.md: dispatch.checklist's "This is not a work breakdown..." justification; consent-stop tail ("The license hangs off... bites hardest in non-dev workflows"); second-verifier's third restatement ("N agents... agreement raises cost, not confidence"); fan-out checkpoint why-clauses (~25%).
- claude-fo-dispatch.md: Awaiting Completion states one rule four ways with narrated failure psychology - compress the narration, KEEP the anti-pattern list (bought with a real incident; captain-reviewed compression, not deletion); reconcile binding's interleaved rationale ("Cost of a miss..." etc.) - spec stands alone.
- fo-write-core.md: Workflow Fit Gate's two arguing paragraphs (rule = check fit, name the output's existing home, ask when ambiguous); the --next-id-preview/new-closes-window lecture stated three times - once plus pointers.
- fo-merge-core.md: "the toil merge.guard eliminates" tail.

The deferral groups and open Working Principles residency question are s6q's, verbatim in its archived body - ideation rules on Working Principles explicitly rather than leaving it unexamined.

## Problem

The FO contract carries two kinds of weight the reader pays for and does not use. Lecture: rationale, design history, and third restatements of a rule already stated — every session that loads the file pays it. Boot residency: sections in `first-officer-shared-core.md` that cannot fire at a greet-and-stop, the most common interactive shape.

Measured at the stack base (`spacedock-ensign/install-gate-channel-aware-hint`), the five files in the cut list:

| file | bytes | lines | load point |
|---|---|---|---|
| first-officer-shared-core.md | 26,103 | 202 | boot |
| fo-dispatch-core.md | 29,444 | 224 | first dispatch |
| claude-fo-dispatch.md | 19,934 | 130 | first dispatch (Claude) |
| fo-write-core.md | 7,643 | 53 | first FO-authored mutation |
| fo-merge-core.md | 4,296 | 28 | first terminal boundary |

s6q's archived profile is the boot number: session afd74765, 2026-07-25, cold baseline 23,721 tok, greet-and-stop occupancy 40,450 tok, shared core 26,092 B ≈ 10,120 tok of it, candidate deferral span 9,105 B = 35% of the boot-resident core.

Two honesty notes carried forward from that profile. This is deferral, not reduction, for Movement 2: a session that dispatches and gates pays the same tokens later. The shapes optimized are greet-only, question-only, and gate-only. And Movement 1's deletions ARE a reduction, on every load path, which is why they are the larger and more certain half of this task.

### Ideation finding: three of s6q's four "first-dispatch group" members fail AC-3

s6q grouped `## Completion and Gates`, `«gate.ac-cross-check»`, `«gate.assemble-verdict»`, and `«feedback.route»` as rideable on the existing `fo-dispatch-core.md` read, on the reasoning that "none can fire before a worker completes." That reasoning holds for `## Completion and Gates` and fails for the other three.

`«engage»` selects gate-first: "Immediately load `spacedock:fo-gate-lifecycle` before any gate evidence… Only when no gate wins, load the dispatch owner." A gate left ready by a prior session is reached at engage with `fo-dispatch-core.md` never read. `skills/fo-gate-lifecycle/SKILL.md:57` and `:69` invoke `«feedback.route»` directly on that path, and `«gate.ac-cross-check»` is declared "at every gate". Moving those three behind the dispatch-core trigger puts them behind a trigger they can be reached before — the exact condition AC-3 forbids.

Their only always-loaded-at-gate home is `skills/fo-gate-lifecycle/SKILL.md`, which is 7,645 B against a captain-set 7,700 B cap (`internal/contractlint/fo_function_reference_invariant_test.go`) — 55 bytes of headroom for 2,653 B of content. That is a captain cap decision, not a mechanical cleanup, so it is out of scope here and filed as a follow-up. **Gate decision point:** approving a `fo-gate-lifecycle` cap raise to ~10,400 B would add ~2,653 B ≈ 1,020 tok to the boot win; the recommendation is to keep this task a clean de-lecture plus one provably-safe move and take the cap raise on its own evidence.

### Working Principles residency: RESIDENT (s6q's open question, ruled)

`## Working Principles` is 6,301 B, 24% of the boot-resident core, and stays boot-resident. Four of its five posture bullets bind the greet itself, not a later boundary: "Lead with a recommendation the captain can say yes to" and "Speak the workflow's declared label, not the generic entity" govern the greet's own captain-facing prose; "Prefer the cheapest check that can fail" and "Do obvious reversible work without ceremony" govern the first decision after it. Deferring the section would strip posture from the one shape this task exists to make cheaper.

The counter-argument s6q raised — that the gate/merge/triage bar fires only at those boundaries — is true of exactly one sub-block, the `Hold your own gate, merge, and triage calls…` blockquote (~1,050 B). Relocating it needs an anchor in both `fo-gate-lifecycle` and `fo-merge-core`, or a duplicated copy: two anchors for ~400 tok. Below the bar. The section keeps its residency and gives up only the ~370 B of justification tails the audit named.

Also ruled out, for the record: the engage group (2,572 B) and the narrow group (1,202 B) both require minting at least one new deferred reference file plus a new trigger, because no existing deferred file is read at engage time — `state ready`, `«hooks.run»("startup")`, and the exit-3 halt all precede the gate/dispatch branch that owns the only existing engage-time loads. New file plus new trigger for ~980 tok, and three separate triggers for the narrow group's ~460 tok, is the worst mechanism-per-token ratio in the candidate set. Not in this task.

One member of the narrow group is a delete rather than a move. `## Compaction continuity`'s "After" bullet duplicates `hooks/session_start_compact_reminder.sh`, which already fires at `SessionStart(compact|clear)` and injects the same instruction ("A summary claiming you read or did something is not the reading or the doing… re-satisfy each load precondition and state read at its existing trigger"). The hook is registered in `hooks.json` under `${CLAUDE_PLUGIN_ROOT:-${PLUGIN_ROOT}}`, which covers Claude and Codex but not Pi, so the bullet compresses to one clause rather than disappearing; the full delete becomes available when `align-pi-compaction-with-force-boot` lands.

## Proposed approach

### Sequencing: deletes first (one pass per file), then the one move, then tests

Three commits, in this order.

1. **Deletes**, one pass per file across all five files.
2. **The move**: `## Completion and Gates` out of `first-officer-shared-core.md` into `fo-dispatch-core.md`, plus its pointer.
3. **Test updates**.

Deletes lead because one of them — the stopping-point paragraph — lives inside the section that Movement 2 relocates. Deleting first means the move relocates already-thinned text and its diff reads as a pure relocation. Move-first would land that deletion in a different file from the one the per-file table prices it against. Deletes-first also gives the separable measurement AC-2 needs: commit 1's byte delta is the reduction, commit 2's is ~zero by construction.

### Movement 1 — the deletions, with the exact text struck

`first-officer-shared-core.md`:

- **L32** (`«engage»` scope bullet), strike the design-history clause: `; the `workflow` argument is present so a future multi-workflow form extends this signature rather than replacing it`. After: `- **scope:** ONE workflow per invocation.` (-116 B)
- **L77** (stopping-point paragraph, 797 B), compress the re-argument, keep the exception list — the list is the content. Before: `**A completed non-gated, non-terminal stage is not a stopping point.** After verifying the report, the FO MUST advance the entity to the next stage and dispatch it (reuse-or-fresh per the dispatch module's reuse conditions) BEFORE ending its turn. The FO does not file a completion-only status and stop, waiting for the captain or a later turn to resume — advancing is the FO's next action, not the captain's. The only conditions that legitimately halt the turn here are: {list}. Absent one of those, stopping after a completion-only report is a contract violation.` After: `**A completed non-gated, non-terminal stage is not a stopping point.** After verifying the report, the FO MUST advance the entity to the next stage and dispatch it (reuse-or-fresh per the dispatch module's reuse conditions) BEFORE ending its turn. Only these conditions legitimately halt the turn here: {list unchanged}. Absent one, stopping after a completion-only report is a contract violation.` (~-320 B)
- **L163**, strike ` Dense durable boundaries keep unhinted auto-compaction survivable.` (-67 B)
- **L164**, compress the "After" bullet to its instruction, dropping the restated justification the hook already injects. (~-120 B)
- **L183** (smallest-mechanism bullet), strike the third restatement `Raising an answer's thoroughness never raises the mechanism's weight.` and the tail `, because commissioning already fixed the mechanism`. KEEP the banned-excuse enumeration (`Never by "it's substantive," "Ultracode is on," …`) — prohibitions are behavior-shaping. (~-180 B)
- **L192**, strike `Naming the end without gating it is the asymmetry that lets a means-accurate, end-missed stage pass.` (-103 B)
- **L195**, strike the second clause `: a scheme minted in a dispatch prompt becomes every downstream artifact's vocabulary`. (-85 B)

`fo-dispatch-core.md`:

- **L35**, delete the whole `This is not a work breakdown…` paragraph. The rule it justifies already stands in the `done-when` immediately below it (`no more than three outcome signals and no structural task boilerplate`), and the AC-vs-checklist distinction is owned by `«gate.ac-cross-check»`. (-330 B, -2 lines)
- **L218** (consent stop), strike the tail `The license hangs off the captain wanting it, never an inference that it would help; it bites hardest in non-dev workflows, where every check is new process.` (~-155 B)
- **L220** (fan-out checkpoint), strike the why-clauses — `— dispatched workers are the counted unit, whether they file entities or open PRs`, `since a running script reaches no later moment to catch`, `Keep-moving speeds independent, already-scoped work; it does not authorize an unbounded spawn chain off one thread.` KEEP the declaration rule, the stop condition, the barrier-stage rule, and the async-binding rule. (~-350 B)
- **L222** (second verifier), strike the third restatement `And N agents reaching the same answer is one confirmation observed N times, not N independent confidences — agreement raises cost, not confidence.` (-146 B)

`claude-fo-dispatch.md`:

- **L61-70** (`## Awaiting Completion` narration, 2,156 B), compress ~40%. The section states one rule four ways with narrated failure psychology. Keep: the three-item completion-signal enumeration (L55-59, untouched), the first-turn decision procedure's rule, the four `Do not:` bullets stripped of their psychological narration, and `Just emit end_turn with empty content.` Strike the narration: `this converts idle-polling into a multi-turn generation loop that drifts into hallucination on subsequent wake-ups`, `the runtime handles the wait for you; sleeping in Bash wastes time and does not accelerate delivery`, and the `**A new system init entry…**` paragraph's second and third sentences. KEEP the `## Anti-patterns` list at L72-78 intact — captain-reviewed compression, not deletion; it was bought with a real incident. (~-860 B, -2 lines)
- **L118** (reconcile step-0), strike the interleaved parentheticals `(a stale prior-session or parallel-session config must never be mistaken for the live team)` and `(the member is registered at the subagents/agent-*.meta.json path instead, and reading that roster is not implemented)`. The behavior statements around them stand alone. (~-200 B)
- **L130**, strike `Cost of a miss: one extra event-loop cycle.` (-44 B)

`fo-write-core.md`:

- **L20** (Workflow Fit Gate, first arguing paragraph), keep sentence 1, strike sentences 2-3: `Write authorization is not workflow-fit authorization. The write classifier is not evidence of fit either: a path's class says who may write it, never whether this workflow should be tracking this work.` (~-215 B)
- **L26**, delete the whole second arguing paragraph (`A fit failure is not repaired by adding a shippable mechanism…`). (-295 B, -2 lines)

  What survives is the rule in three clauses, which is what the gate is: check fit (L22), name the output's existing home (L24, with its enumeration kept), ask when ambiguous (L28).
- **L35** (new-entity bullet), replace the duplicated preview-vs-mint explanation with a pointer. Before: `…and `new` mints the id, stamps it into the frontmatter, and atomically writes the stamped entity in one call, so no `--next-id` candidate can drift between preview and write. … Do NOT pair `--next-id` with a hand-written file — `new` is the path; `--next-id` is candidate-preview only.` After: `…and `new` mints the id, stamps it into the frontmatter, and atomically writes the stamped entity in one call (see `## ID Styles` for the preview-vs-mint window).` The full statement stays once at L53, its natural home. (~-300 B)

`fo-merge-core.md`:

- **L11**, strike the tail ` — the toil `«merge.guard»` eliminates`. The preceding clause (a stale binary makes the subcommand unknown and you silently fall back to the hand ceremony) is a real failure mode and stays. (-37 B)

### Movement 2 — the one move

`## Completion and Gates` (shared-core L62-81, 2,678 B before deletes, ~2,358 B after) moves verbatim into `fo-dispatch-core.md`, placed immediately before `## Reuse and Fresh Dispatch` — which already opens by deferring to it (`fo-dispatch-core.md:42`, "The gate-presentation spine is in the boot-resident core's `## Completion and Gates`"). That cross-reference becomes an intra-file one.

The trigger is provably not-before: the section's entry condition is "When a worker completes," a completion signal requires a same-session `«worker.spawn»`, and every spawn path is inside `fo-dispatch-core.md`, which the contract requires read before dispatch. No new file, no new trigger, no extra round trip.

Two edits fall out of the move:

- The shared core's `## Deferred load points` entry for `references/fo-dispatch-core.md` gains the completion clause, so the section is named at its trigger: `read before `«dispatch.next-action»()`, worker dispatch, or dispatch-state mutation, and it owns `## Completion and Gates` — the worker-completion procedure and gate routing.` (+~230 B, +1 line)
- `«interaction.boundary»`'s headless-with-conn branch (L21) says "resolve gates per `## Completion and Gates`". That reference now crosses files and becomes `per `references/fo-dispatch-core.md`'s `## Completion and Gates``. (~0 B)

`«gate.ac-cross-check»`, `«gate.assemble-verdict»`, and `«feedback.route»` stay boot-resident, per the AC-3 finding above.

Expected boot-resident result: `first-officer-shared-core.md` 26,103 → ~22,984 B, a -3,119 B / ~-1,200 tok reduction (-12% of the file, ~3.0% of the 40,450-tok greet occupancy). Deletions across all five files total ~-3,923 B ≈ ~-1,520 tok, paid back on every load path, not only the greet. Neither number is a transformation, and the gate should read them as such: the honest claim is a ~1,200-token greet win plus a ~1,500-token cut that every session pays less for.

### Mechanisms introduced, and why the cheaper rung will not do

One, and it is a test. `internal/contractlint/boot_resident_closure_test.go` gains a moved-section table asserting, for each `(heading, from, to)` row: the heading is ABSENT in the boot-resident file, present exactly once in the deferred target, and the target is named in the shared core's `## Deferred load points` block (reusing `deferredLoadPointsBlock`, already in `first_officer_eager_references_test.go`). It serves AC-3.

The cheaper rungs, in the order Working Principles requires: no shipped guard covers section residency; the nearest existing mechanical check, `TestBootResidentDeferredLoadPointsResolve`, proves a named target *exists* but says nothing about where a section *lives*, so it cannot fail on a section left duplicated in both files or moved behind an unnamed target; and a falsifiable one-off exercise (the live-lane boot-through-gate run) proves this revision but leaves the next edit unguarded. This is not the "new standing check" the last-resort rule gates — it is a test function in an existing test file exercising the behavior in hand, which that same rule names as ordinary work the proof policy already requires.

One number changes in a second existing test: `TestFOInstructionComponentCaps`'s `first-officer-shared-core.md` cap ratchets 26,900 → 23,500 B (measured post-change ~22,984 plus ~520 headroom). This is an edit to an existing check, not a new one, and it is what keeps AC-2's win durable — without it the ~3,900 B of freed cap headroom refills silently on the next contract edit.

### Semantic changes declared

None to command grammar, stored formats, authority, or runtime behavior. Every rule keeps its wording-equivalent semantics and its current trigger.

One load-timing change: `## Completion and Gates` moves from boot to the first-dispatch read. The declared risk is an FO that reaches a worker completion without having read `fo-dispatch-core.md`; the contract already forbids that state (dispatch requires the read, completion requires a dispatch), and AC-3 exists to prove it. No user-visible surface changes — no CLI output, banner, host integration, or docs-site content — so no doc diff is owed.

### No spike needed

Every mechanism this task relies on is already proven in-repo: the deferred-load-point discipline itself (the dispatch, write, merge, status, and dispatch-recovery cores all measurably stay unread at a greet, per s6q's profile); `boot_resident_closure_test.go`'s filesystem-oracle reference closure, which already ships with its own negative control (`TestBootResidentDeferredLoadPointGuardFailsOnDanglingTarget`); `deferredLoadPointsBlock`'s section extraction; `TestFOInstructionComponentCaps`'s byte measurement; and the `boot-forensics` per-turn occupancy extraction that produced the 40,450-tok baseline. Nothing here is an unverified parser round-trip, runtime handoff, or on-disk format.

## Out of scope

fo-install.md (#757 already rebuilt it under this lens). Any change to what the FO does - load timing and prose weight only; every rule keeps its semantics and trigger. The cold baseline (system prompt, tool schemas). Orphan cleanup (tracked separately).

Added by ideation: the `fo-gate-lifecycle` cap raise and the relocation of `«gate.ac-cross-check»` / `«gate.assemble-verdict»` / `«feedback.route»` (follow-up, needs captain cap approval on its own evidence); s6q's engage group and narrow group (new file plus new trigger, ruled below the bar); the full delete of the `## Compaction continuity` "After" bullet (blocked on Pi hook coverage via `align-pi-compaction-with-force-boot`).

## Expected surface and tolerance

**Estimate net LOC change: +32, across 7 files** (insertions ~+82, deletions ~-50). The net is positive because the one new contractlint test (~+38 lines) is larger than the contract-prose reduction, which is dominated by intra-paragraph clause strikes that count +1/-1.

The contract surface AC-2 measures, `skills/first-officer/references/**` alone: **-3,700 bytes / -5 net lines across 5 files.**

| file | Δ bytes | Δ lines | what |
|---|---|---|---|
| first-officer-shared-core.md | -3,119 | -19 | 7 strikes (-991, of which -320 sit inside the moving section) + `## Completion and Gates` out (-2,358 post-strike) + pointer (+230) |
| fo-dispatch-core.md | +1,377 | +18 | 4 strikes (-981) + `## Completion and Gates` in (+2,358) |
| claude-fo-dispatch.md | -1,104 | -2 | Awaiting Completion compression, reconcile parentheticals, Cost-of-a-miss |
| fo-write-core.md | -810 | -2 | Workflow Fit Gate arguing paragraphs, next-id pointer |
| fo-merge-core.md | -37 | 0 | the toil tail |
| **contract subtotal** | **-3,693** | **-5** | |
| internal/contractlint/boot_resident_closure_test.go | +~1,300 | +38 | moved-section table + fo-dispatch-core in `deferredModuleBodies` |
| internal/contractlint/fo_function_reference_invariant_test.go | ~0 | 0 | cap 26,900 → 23,500 |

Tolerance: **±25 net LOC, ±2 files** overall; **±800 bytes** on the contract subtotal. The byte figure is the one that binds — see AC-2.

**Baseline caveat.** All figures are measured against `spacedock-ensign/install-gate-channel-aware-hint`, which is in flight. Implementation re-measures at the branch point and re-confirms each cut-list span still exists before striking it; a span that moved or was already removed by #757 is reported, not hunted.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (value) — A greet-and-stop FO boot occupies measurably fewer context tokens after the change than before it, on the same host and session shape.**
Verified by: the `boot-forensics` per-turn occupancy extraction run twice on Claude Code — once at the stack base, once at the stack tip — each on a fresh greet-and-stop session, reported in tokens as an A/B delta and additionally against the historical 40,450-tok figure. Expected: ~-1,200 tok. The number lives outside every file this task edits.
Falsifying change: a pointer that restates the moved section instead of naming it, or any edit that lands one extra boot-time Read, leaves the A/B delta flat or positive.
Why A/B and not the bare 40,450 baseline: that figure was measured 2026-07-25 on `0.26.0+dev`, and the contract has changed since. Comparing a tip measurement to a month-old absolute would credit this task with unrelated drift in either direction. Both numbers are reported; the A/B is the one that decides.

**AC-2 (value) — The contract's total prose weight goes DOWN, not sideways: `skills/first-officer/references/**` is net smaller in BYTES against the merge base, with behavior parity.**
Verified by: summed `wc -c` over `skills/first-officer/references/*.md` at the merge base versus the tip, expected ~-3,700 B (tolerance ±800). Parity evidence: the offline suite green (`go test ./...` and `-race`), and all three host live lanes green at the stack tip.
Falsifying change: relocating a paragraph into a sibling file inside the same directory instead of deleting it leaves the summed byte total flat — that is exactly the failure this AC exists to catch, and exactly the failure a line count would miss.
Why bytes, not the seeded line delta: nearly every cut here is an intra-paragraph clause strike in a one-paragraph-per-line markdown file, so it scores +1 insertion / -1 deletion and nets zero. A file that sheds 3,000 bytes of lecture can show a net line delta of zero. Lines cannot falsify this claim; bytes can. The Task Template's net-LOC declaration is kept above for the expected-surface contract; the value measurement is byte-denominated.

**AC-3 — No section is moved behind a trigger it can be reached before.**
Verified by two forms, neither a prose-grep over the contract.

*Static (structural absence + reference closure).* A new table-driven case in `internal/contractlint/boot_resident_closure_test.go` asserts for the one moved row (`## Completion and Gates`, from `first-officer-shared-core.md`, to `fo-dispatch-core.md`): the heading is absent in the source file, present exactly once in the target, and the target is named in the shared core's `## Deferred load points` block. `fo-dispatch-core.md` joins `deferredModuleBodies` so the existing `TestBootResidentDeferredLoadPointsResolve` keeps resolving the `spacedock:fo-gate-lifecycle` reference the moved section carries.
Falsifying changes, each independently: leave a copy of the heading in the shared core (absence assertion fails); move the section to a file the deferred-load-points block does not name (reachability assertion fails); point the block at a file that does not exist (closure fails, with its existing dangling-target control proving the check can fail).

*Behavioral (boot through gate).* A live-lane run drives an FO through boot → engage → dispatch → worker completion → gate with no missing-section failure, on all three hosts. This is what catches the class of error the design review already found once statically: the gate-first engage path that reaches `«feedback.route»` without loading `fo-dispatch-core.md`.
Falsifying change: had the three gate/feedback `«fn»` blocks been moved as s6q proposed, the gate-first lane would reach `spacedock:fo-gate-lifecycle`'s `«feedback.route»` invocation with the definition unloaded.

**AC-3 serves AC-1 and AC-2.** It is the mechanism AC: deferral is only a legitimate way to buy AC-1's token reduction if the moved section is still reachable at its trigger. Per the gate's end re-anchor, AC-3 passing while AC-1's A/B delta is flat or positive is a REJECT, not a pass.

## Test plan

**Offline (every commit).** `go test ./...` then `go test ./... -race`, plus `gofmt -w ./cmd ./internal`. The load-bearing packages are `internal/contractlint` (the AC-3 static form, the cap ratchet, the eager-import topology, the reconcile drift-class binding against `claude-fo-dispatch.md`) and `internal/cli` (`prose_function_routing_test.go` walks shared-core, dispatch-core, and merge-core for `→` migration-target lines; the moved `«fn»` blocks stay inside that walked set, so a strike that damages a `→` line fails there). Cost: minutes, no new harness.

Confirmed clear before implementation: no Go test asserts on any cut-list phrase. The only in-repo occurrences of `not a work breakdown`, `bites hardest`, `agreement raises cost`, `Cost of a miss`, `Dense durable boundaries`, and `future multi-workflow` are recorded live-stream fixtures under `internal/ensigncycle/testdata/`, which are historical transcripts used as inputs, not assertions against the current contract text.

**Live lanes (required green at the stack tip, not optional).** The diff is `skills/**/references/**`, the shipped-contract high-stakes surface, and four of the five files are host-neutral core, so every host lane is required: `claude-live` (both matrix models), `codex-live`, and `pi-live` in `.github/workflows/runtime-live-e2e.yml`. `claude-fo-dispatch.md` is Claude-specific, which makes `claude-live` the tightest lane on the compression, but it does not narrow the requirement — the other four files bind all three. These lanes supply AC-3's behavioral form and AC-2's parity evidence in one run.

**Detached adversarial audit.** Required by the proof policy for this surface, and the right rung here for a reason no mechanical check covers: the question "does this deletion remove a justification or a rule?" is judgment no lint owns and no diff read settles. Its brief is the necessity question first — for each strike, is what remains still sufficient to produce the behavior the struck sentence was defending — and specifically to attack the four KEEP calls (banned-excuse enumerations, the anti-pattern list, the fan-out declaration rule, the stale-binary failure mode) and the AC-3 reachability argument for the single move. Detached from the implementing worker.

**AC-1 measurement.** Two `boot-forensics` runs, stack base and stack tip, Claude Code, fresh greet-and-stop each, occupancy in tokens. Recorded in the validation stage report with both the A/B delta and the reading against 40,450.

**AC-2 measurement.** `git show {merge-base}:skills/first-officer/references/*.md | wc -c` summed against the tip's `wc -c`, plus the per-file table above re-measured as actuals versus this ideation's declared figures.

## Delivery

Stacked layer on the existing stack, per the captain's order: branch off spacedock-ensign/install-gate-channel-aware-hint (#757, which already carries the shared-core edits this task touches next). Stack becomes #756 -> #757 -> this; the live lanes run once at this tip and prove all three layers. gh pr create + gh stack link per the pr-merge mod's Stacked mode.

This is stacked layer 3. Implementation branches off `spacedock-ensign/install-gate-channel-aware-hint` at its then-current tip, not off `main`, and re-measures the per-file baseline there per the caveat under Expected surface. #757 renames `fo-install-gate.md` to `fo-install.md` and edits shared-core lines 10 and 47; neither line is in this cut list, so the layers do not conflict, but layer 3 must not resurrect the old filename in the `## Deferred load points` block it edits.

## Stage Report: ideation

- DONE: The task body prices the audit cut list per file (net LOC table with tolerance), sequences deletes-vs-moves, and rules explicitly on Working Principles residency - LIGHT per captain order: the cut list and s6q's archived measured profile are the spec; no re-derivation, no new mechanisms, "no spike needed" is expected with the proven deferral machinery named.
  Per-file table under `## Expected surface and tolerance` carries measured byte/line deltas for all 5 contract files plus 2 test files (contract subtotal -3,693 B / -5 lines; overall +32 net LOC across 7 files, tolerance ±25 LOC / ±2 files / ±800 B). Every cut-list span was measured at the stack base, not estimated. Sequencing ruled deletes-first-then-move under `### Sequencing`, because the stopping-point strike lives inside the section Movement 2 relocates. Working Principles ruled RESIDENT with the reason (4 of its 5 posture bullets bind the greet itself) and the counter-argument priced (the one boundary-scoped blockquote is ~1,050 B needing two anchors). "No spike needed" recorded with five named proven mechanisms.
- DONE: ACs re-anchored: the greet-occupancy token AC and the net-negative-LOC-with-parity AC each name their falsifying change; AC-3's not-a-prose-grep verification is designed (contractlint structural absence + reference closure static form, live-lane boot-through-gate behavioral form).
  AC-1 falsifying change: a restating pointer or one extra boot-time Read leaves the A/B delta flat or positive; re-anchored to a stack-base-vs-tip A/B because the 40,450 figure is a month old on a different binary, so a bare comparison would credit unrelated drift. AC-2 re-denominated from lines to BYTES with the reason stated: nearly every cut is an intra-paragraph clause strike scoring +1/-1, so a 3,000-byte reduction can show a zero line delta - lines cannot falsify this claim, bytes can. AC-3 designed on existing machinery: a moved-section table in `boot_resident_closure_test.go` asserting absence-in-source, exactly-once-in-target, and target-named-in-`## Deferred load points` (reusing `deferredLoadPointsBlock`), with three independent falsifying changes named; behavioral form is the boot→engage→dispatch→completion→gate live run. AC-3 explicitly declared as serving AC-1/AC-2 for the gate's end re-anchor.
- DONE: The test plan names the required lanes (all three hosts at the stack tip, detached audit) and the Delivery section pins stacked layer 3: branch off spacedock-ensign/install-gate-channel-aware-hint per the pr-merge Stacked mode.
  Test plan names `claude-live` (both matrix models), `codex-live`, and `pi-live` from `.github/workflows/runtime-live-e2e.yml` as required-green (4 of 5 files are host-neutral core), the offline suite with its two load-bearing packages, and the detached adversarial audit with a necessity-first brief targeting the four KEEP calls and the AC-3 reachability argument. Verified no Go test asserts on any cut-list phrase - the only in-repo hits are recorded stream fixtures under `internal/ensigncycle/testdata/`. Delivery pins layer 3 off the #757 tip with the non-conflict argument (#757 touches shared-core lines 10 and 47; neither is in the cut list) and the one hazard (do not resurrect `fo-install-gate.md` in the deferred-load-points block).

### Summary

The design changed on one finding: s6q grouped four sections as rideable on the existing `fo-dispatch-core.md` read, but `«engage»` selects gate-first and loads `spacedock:fo-gate-lifecycle` *instead of* the dispatch owner, so a gate left ready by a prior session reaches `fo-gate-lifecycle/SKILL.md:57`'s `«feedback.route»` invocation with dispatch-core never read. Three of the four would have been moved behind a trigger they can be reached before - the exact AC-3 violation. Only `## Completion and Gates` moves; the three gate/feedback `«fn»` blocks stay resident, and their proper home (`fo-gate-lifecycle`, 55 bytes under a captain-set cap) is filed as a follow-up needing its own cap approval.

That shrinks the deferral win from s6q's 5,331 B to 2,678 B, so the honest headline is a ~1,200-token greet reduction plus ~1,500 tokens of lecture deleted off every load path - stated plainly in the body rather than rounded up. Two secondary findings: AC-2's seeded line-delta measurement cannot falsify its own claim on this diff shape and was re-denominated to bytes, and the `## Compaction continuity` "After" bullet duplicates `hooks/session_start_compact_reminder.sh` verbatim in substance, but the hook is registered only under `${CLAUDE_PLUGIN_ROOT:-${PLUGIN_ROOT}}` and so does not cover Pi - it compresses rather than disappearing, with the full delete blocked on `align-pi-compaction-with-force-boot`.

## Stage Report: implementation

- DONE: Three commits in the ruled order on a branch REBASED onto the #757 tip (deletes one-pass-per-file striking exactly the priced Movement-1 spans; the single Completion-and-Gates move into fo-dispatch-core.md plus its one-line pointer; test updates) - contract subtotal within -3,693 B +/-800 B, no cut beyond the priced list, fo-install-gate.md NOT resurrected in the deferred-load-points block.
  Branch reset onto parent tip `ef43c7ca8` before any edit; verified `fo-install.md` present and `fo-install-gate.md` absent first. Commits `7d7b16b7d` (strikes), `faae94380` (move), `9d6b948e8` (tests). Contract subtotal **-3,210 B** against the parent (target -3,693 +/-800 = -4,493..-2,893). `grep -c fo-install-gate` on the shared core = 0. Per file, strikes then move: shared-core 25,871 -> 25,123 -> 22,545 (-3,326); fo-dispatch-core 29,444 -> 28,548 -> 31,271 (+1,827); claude-fo-dispatch -1,027; fo-write-core -642; fo-merge-core -42. Every cut-list line number matched the ideation at the new baseline, so no span was hunted.
- DONE: The AC-3 machinery lands and is falsified: the moved-section table in boot_resident_closure_test.go (absence-in-source, exactly-once-in-target, target named in Deferred load points) with all three named falsifying changes re-red then restored; full UNFILTERED go test ./... and -race green (known machine-local codex failure excepted, reproduced at parent).
  `TestMovedSectionsLiveBehindANamedTrigger` asserts the three properties per row and reuses `deferredLoadPointsBlock`; `fo-dispatch-core.md` joined `deferredModuleBodies`. Falsification transcript, each restored to green after: (1) appended a duplicate `## Completion and Gates` to the shared core -> `still appears 1 times in ...first-officer-shared-core.md`; (2) physically moved the section into `claude-fo-dispatch.md` (a file the block does not name) and repointed the table -> `deferred load-points block does not name "references/claude-fo-dispatch.md"`; (3) renamed the block entry to `fo-dispatch-core-nope.md` -> `TestBootResidentDeferredLoadPointsResolve` failed on the dangling stat. Unfiltered `go test ./...` and `go test ./... -race`: every package `ok` except `internal/cli`'s `TestCodexResolveManifestAgainstInstalledHost`, which reproduces byte-identically on a clean checkout of `ef43c7ca8` (machine-local codex install state, touches no contract file).
- DONE: AC-1/AC-2 evidence recorded in the stage report: the commit-1 byte delta as the measured reduction with commit-2 ~zero by construction, and the report names the required lanes at the stack tip plus the lean-validation posture the captain set.
  AC-2: commit 1 alone is **-3,355 B** across `skills/first-officer/references/**` - the reduction, paid on every load path. Commit 2 is **+145 B** (the deferred-load-points clause plus two cross-reference repairs), confirming the move is sideways by construction, exactly the failure mode AC-2 exists to catch. AC-1 is validation's measurement, not mine; the mechanism it prices is the boot-resident core at **22,545 B**, down 3,326 B from 25,871 (~-1,280 tok), with the cap ratcheted 26,900 -> 23,500 (955 B headroom) so the win cannot silently refill. Required green at this tip per the test plan: `claude-live` (both matrix models), `codex-live`, `pi-live`, plus the detached adversarial audit. Captain's posture (ideation gate, 2026-08-24): "keep the validation lean. i'll review on github" - branch pushed as stacked layer 3; PR creation stays the FO's merge-stage ceremony.

### Summary

Three commits on `spacedock-ensign/de-lecture-and-defer-fo-contract`, rebased onto `ef43c7ca8`: the priced strikes, the single `## Completion and Gates` move, and the AC-3 guard with the cap ratchet. Contract subtotal -3,210 B, boot-resident core -3,326 B.

Four things the reviewer should look at rather than take on trust. First, net LOC is **+59 across 7 files** against the declared +32 +/-25 - a 2-LOC breach, entirely in the test files (the moved-section table came in at 56 lines, not the estimated 38, and the cap ratchet carries an 8-line provenance comment the estimate priced at zero). The contract subtotal, which the ideation declared binding, is comfortably inside its +/-800 B band; the byte figures also ran ~20% under the per-strike estimates across the board, which is why the subtotal landed at -3,210 rather than -3,693. Second, the `fo-write-core.md` next-id cut follows the ideation's literal After text, which drops the prohibition `Do NOT pair --next-id with a hand-written file` and not just its justification. The full preview-vs-mint statement survives at `## ID Styles`, but this is the one strike that removes a rule rather than a rationale, and it is squarely the necessity question the detached audit's brief owns. Third, the ideation named two fallout edits for the move; there were three - shared-core's Working Principles bullet also carried a bare `(see Completion and Gates)`, now pointing across files. Fourth, `gofmt -w ./internal` reformatted `internal/release/runtime_live_evidence_workflow_test.go`, a pre-existing violation from unrelated commit `e68aa5b7a`; reverted to keep the diff at the declared surface, and it is still unformatted on the parent.
