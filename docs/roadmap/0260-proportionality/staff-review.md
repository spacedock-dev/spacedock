# 0260 proportionality — sprint-wide preflight staff review

## Scope, method, and evidence base

This is an independent sprint review, not a review of the shaping FO's own work by its authors. I tried to refute whether the eight-member drivable set composes into the 0.26.0 deliverable and whether a fresh Commander can drive it without session context. I used the six lenses required by `docs/roadmap/README.md`: DoD ownership, dependency order, cross-member wiring and collisions, set-wide blast radius, missing or excess scope, and cold-boot readiness. I surfaced a member-level defect only where it breaks the sprint package.

The live membership queries both return the same eight members because the second re-lock removed deferred work from the sprint rather than retaining deferred members inside it:

```text
85   merge-guard-arm-not-a-stopping-point
ht   fix-tautological-output-grep-tests
bw   feedback-cycle-record-command
02av ensign-finding-triage-disposition
z7   falsifiability-ladder
az   anti-tautology-enforcement-and-template-gap
841  contractlint-codex-runtime-semantics-retirement
2ae  template-rigor-propagation
```

I read the sprint index and Commander dispatch in full; the current body of every member above, using `template-rigor-propagation/index.md` for the in-progress folder move; the 0240 staff-review precedent; and the 0260 forensic synthesis and remedy digest. I checked the claims against the live status query, the four now-unlabelled contractlint backlog headers, the committed `.roborev.toml` materiality taxonomy, and the current `docs/dev/README.md` proof and design-reset rules.

## Verdict: READY AFTER FOLDS

0260 is not ready to drive as packaged. Four Material defects remain: the reframe DoD promises a stop and a unit-tested refusal that the approved member explicitly defers; the recurrence DoD promises a `testlint` gate that the approved design explicitly rejects; the index still advertises deferred or removed scope with no owner; and the Commander dispatch omits load-bearing shared-edit and pre-cut proof wiring.

The underlying reduced design is coherent. The smallest repairs are wording and wiring folds into `bw`, the sprint index, and the Commander dispatch. No new member or binary mechanism is required. After those folds, the eight-member set is proportionate and drivable.

Finding count: **4 Material, 3 Polish**.

## DoD-to-owner coverage

| Sprint DoD bullet | Owner(s) | Coverage | Staff assessment |
|---|---|---|---|
| e6j replay reaches a design-reset decision by cycle 2; dispatch-refusal guard has a fixture-fed unit test (`index.md:35`) | `bw`, with AC-narrowing rule from `02av` and replay pressure from `z7` | **Partial / contradictory** | `bw` makes deviation visible by round 2, but defers the command and every refusal/backstop. Its first cut asks the FO to *weigh* a reset, not record one before another repair. No member owns the promised refusal unit test. Material finding 1. |
| Correct-but-disproportionate finding produces a recorded decline and zero-line diff (`index.md:36`) | `02av`, composed with `bw`'s record | **Owned** | AC-1 supplies the decline path and a material control; AC-2 can fail on a misclassified decline. The archived incident-13 shape and current four-field taxonomy support the no-stakes-field design. |
| PTY-harness brief trips consent before dispatch (`index.md:37`) | `z7` | **Owned** | AC-1 is a branch-vs-main live FO drive over the archived harness lure. The consent rule lands at the boot principle and dispatch decision. |
| `testlint` reds on the reverted 11-phrase presence test; eight output-grep tests fixed; gate reads assertion content (`index.md:38`) | `ht` owns the eight fixes; `az` owns the report/evidence rule | **Partial / contradictory** | The eight fixes and assertion-content report change are owned. `az` deliberately rejects a standing automated gate and makes “no new committed gate” AC-4. No member owns the `testlint` claim. Material finding 2. |
| Runtime-semantics phrase checks retire; remaining four retirements are deferred (`index.md:39`) | `841` owns Codex + Pi runtime semantics | **Owned for this train; package text stale for the rest** | `841` maps the ten in-scope functions to Go bindings, fixtures, live lanes, or structural residue. The four backlog entities have blank `sprint` and `sprint-readiness`, not `defer`; the DoD contradicts the second re-lock. Material finding 3. |
| Commissioned workflow carries materiality + fixed Verified-by; refit carries content (`index.md:40`) | `2ae`, consuming settled text from `z7`, `02av`, and `az` | **Owned** | AC-1 drives commission against `main`; AC-2 drives the refit content delta against the version-only control. The refit spike established the mechanism. `2ae` must land last. |

`85` supports the incident-9 merge-arm correction but owns no sprint DoD bullet. Its one-file change is independent and small; it is not needed to close a coverage hole in the six DoD bullets.

## Dependency waves and shared-edit seams

The correct landing sequence is:

1. `85`, `ht`, and `841` may run in parallel. They touch `fo-merge-core.md`, six Go test files, and contractlint tests, respectively.
2. Land `z7` and `az` serially, or rebase the second immediately before validation. Both touch `docs/dev/README.md` near the proof-policy region; `z7` also changes `first-officer-shared-core.md`.
3. Land `bw` and `02av` as one composed change and one shared validation gate. They co-edit `feedback-rejection-flow/SKILL.md`, the implementation stage definition in `docs/dev/README.md`, and the same `### Feedback Cycles` entry shape. The composed entry must contain `bw`'s surface/estimate/AC fields and `02av`'s findings disposition, including the all-declines case.
4. Land `2ae` last. It carries settled sibling wording into commission/refit scaffolding and must validate against the post-`z7`/`az`/`02av` text, not an earlier branch copy.

The shared-edit and semantic-seam map is:

| Shared surface | Editors/readers | Collision or seam |
|---|---|---|
| `skills/first-officer/references/first-officer-shared-core.md` | `z7`, then `bw` | Different sections, same file. Rebase/re-anchor `bw` after `z7`; do not merge stale whole-file prose. |
| `docs/dev/README.md` | `z7`, `az`, `bw`, `02av`; read as source by `2ae` | `z7` and `az` are adjacent in the proof-policy area; `bw` and `02av` compose in the implementation stage. Validate the assembled README before `2ae` copies from it. |
| `skills/feedback-rejection-flow/SKILL.md` | `bw` + `02av` | One semantic record and one delivery path. Separate landings risk a format without `findings` or a triage block writing a field the base format does not define. |
| Commission/refit templates | `2ae` reads sibling outputs | A semantic dependency, not only a text dependency. “Verbatim” must mean the landed sibling wording at implementation time. |

The planned blast radius is otherwise bounded: three contract/skill prose lanes, one dev-workflow README seam, six test-only cleanup files, contractlint test changes, and four template/refit instruction files. No member except the existing contractlint lane adds Go behavior; the reduced `bw`, `02av`, `z7`, `az`, and `2ae` designs expressly avoid new binary enforcement.

## Findings

### 1. Material — the central reframe outcome is stronger than the approved mechanism

**Confirmed defect.** The sprint says the feedback flow “halts at a design-reset decision” (`index.md:15`) and requires the e6j replay to halt by cycle 2 plus a fixture-fed dispatch-refusal unit test (`index.md:35`). The approved `bw` body says the first cut is prose only and defers the record command, escalation marker, refusal, `--force`, and enforcement backstop (`feedback-cycle-record-command/index.md:66-83`). Its generic rule says narration “prompts” a reconfirm-or-reframe judgment and explicitly defers the hard stop (`:97-101`); AC-1 proves only that arithmetic surfaces ≥200% at round 2 (`:171-172`), and AC-4 proves only that the FO is told to *weigh* a reset (`:179-180`). The dispatch correctly warns that the command is deferred (`dispatch-sprint-execution.md:22`), so it cannot also satisfy the index's guard-unit-test requirement.

This matters beyond wording. The forensic record says the four runaway reviewer loops remained contract-legal one round at a time (`synthesis.md:17-20,33`); awareness without a required decision can repeat that failure. The current README already says not to send a mechanism failure through another automatic repair (`docs/dev/README.md:145-147`), yet the incidents occurred. `bw` must carry that discipline into the actual reflow decision.

**Smallest fold:** in `bw`, change the step-3 result from “weigh a design-reset decision” to “record reconfirm / re-scope / park / escalate before another reflow; no automatic repair re-dispatch while that decision is absent.” Keep it prose-only. In `index.md`, replace the deferred refusal/unit-test sentence with the actual falsifiable proof: the e6j fixture reaches threshold by round 2 and a live replay records the required decision before any further dispatch. Mirror the barrier in the dispatch landing note. If the captain intends narration only, the alternative fold is to weaken “halts” explicitly; that would abandon the re-lock's first thesis and is not recommended.

### 2. Material — the evidence-recurrence DoD promises the automated gate the approved design rejects

**Confirmed defect.** The sprint goal says the tautological-test estate is “fixed and gated against recurrence” (`index.md:18`), and the DoD requires `testlint` to red on the reverted 11-phrase presence test (`index.md:38`). `az` reaches the opposite approved decision: no new standing automated prose gate (`anti-tautology-enforcement-and-template-gap/index.md:75-87`), no new committed test/gate/lint in scope (`:89-102`), and AC-4 requires zero new enforcement (`:126-127`). `ht` owns only six deletions and two narrowings across the eight named tests (`fix-tautological-output-grep-tests/index.md:43,58-72`). No drivable member owns the `testlint` red.

The reduced mechanism can still embody “evidence can fail”: `az` removes the pass-count shortcut, requires a falsifying change, and proves the distinction with a mutation exercise; `ht` proves retained behavior with a mutation matrix. But that is review-time discipline, not a `testlint` recurrence gate. Edit D, which would widen the existing detached audit's trigger, is expressly outside the captain's approval and awaits a yes/no (`dispatch-sprint-execution.md:25,34`; `anti-tautology-enforcement-and-template-gap/index.md:157-161`).

**Smallest fold:** resolve Edit D before `az` advances and record the answer in its gate application. Then rewrite `index.md:18,38` to the mechanism actually approved: eight fixes; the report protocol exposes assertion and falsifying change; the representative mutation stays green only for the tautology and reds for the control; the existing audit uses the chosen trigger if Edit D is approved. Delete the unowned `testlint` claim. Add the resolved Edit-D state—not an open yes/no—to the dispatch.

### 3. Material — the re-locked index still advertises removed or unowned scope

**Confirmed defect.** The index correctly says deferred members leave the sprint label entirely (`index.md:6`), but the DoD says four retirements “carry `sprint-readiness: defer`” (`:39`). The four relevant backlog headers—`contractlint-fo-write-classifier-retirement`, `contractlint-reuse-advance-pointer-retirement`, `contractlint-mixed-structural-boundary`, and `contractlint-prose-function-backstop-retirement`—currently have blank `sprint`, `group`, and `sprint-readiness`. The live query therefore cannot prove the DoD text as written.

The layer map also says 0260 lands a project `AGENTS.md` router scaffold and maintenance (`index.md:26`) plus roborev alignment and a stale `AGENTS.md ## Priorities` refresh (`:28`). The current template member expressly defers the router scaffold and maintenance (`template-rigor-propagation/index.md:54-64`), and `02av` leaves the reviewer/production side unchanged (`ensign-finding-triage-disposition.md:102-109`). No drivable member owns the stale-priorities or roborev-alignment edits.

**Smallest fold:** update the index only. State that the four contractlint items are next-train backlog outside the 0260 membership query; remove the false `sprint-readiness: defer` acceptance condition. Move the AGENTS router/maintenance, stale priorities refresh, and any roborev production-side change to Out of scope/next-train unless the captain adds an owner. Keep 0260's actual packet/reviewer reach: the feedback-context triage block, assertion-aware gate evidence, and template propagation.

### 4. Material — the Commander dispatch omits load-bearing composition and proof instructions

**Confirmed defect.** The dispatch notes only the `bw`/`02av` co-edit (`dispatch-sprint-execution.md:22-23`) and says `2ae` lands last (`:27-29`). It omits the `z7`/`bw` collision in `first-officer-shared-core.md`, the four-editor `docs/dev/README.md` seam, and the requirement that the composed feedback entry include both members' fields. A cold Commander using separate worktrees could merge approved prose that is individually correct but semantically incomplete.

The close-out also reduces the pre-cut work to a generic antipattern audit (`dispatch-sprint-execution.md:42-44`). `z7`'s approved test plan requires the lure catalog at validation and at the existing pre-cut slot, under both Claude and Codex/`gpt-5.6-sol` (`falsifiability-ladder.md:255-273`), and the captain added a fifth reviewer-side trap (`:360-362`). The catalog's durable home remains unresolved (`:268-273,353-354`). A cold Commander can finish the written index DoD without running the minting, mechanism-climb, or reviewer-means-AC trap at pre-cut.

**Smallest fold:** expand `dispatch-sprint-execution.md` with the four-row shared-edit table above, require rebase/re-anchor and assembled-README validation between waves, and state that `bw`+`02av` share one implementation/validation gate. Add a pre-cut checklist that runs all five approved lure scenarios under both runtimes and records outcomes in the pre-cut report. Resolve the catalog home before drive; a validation/pre-cut report artifact is sufficient—no new suite or CI lane is needed.

### 5. Polish — live member metadata still carries superseded concepts

`bw`'s title still promises a “binary-owned count, diff-growth refusal” even though its authoritative body defers both (`feedback-cycle-record-command/index.md:3,64-66`). `85` still uses `group: ladder` (`merge-guard-arm-not-a-stopping-point/index.md:12-14`) after the captain rejected that minted name. The dispatch corrects the first and neither changes behavior, but a status-first Commander sees stale concepts before reading the bodies. Retitle/regroup with the normal state mutator when convenient; do not rename stable slugs merely for polish.

### 6. Polish — lifecycle accounting is historical but reads as current

The lifecycle says 23 members were stamped, seven deferred, and nine driving ideations completed (`index.md:76-79`), while the live sprint query returns eight and the second re-lock removed deferred labels. Preserve the history, but mark the counts as first-lock history and add “current drivable set: eight” so cold readers do not infer a hidden roster.

### 7. Polish — the release close-out omits repository-wide completion gates

The dispatch names `go test ./...` and the release runbook (`dispatch-sprint-execution.md:44`). Repository instructions also require `go test ./... -race` and `gofmt -w ./cmd ./internal` before claiming completion (`AGENTS.md`, “Expected Commands”). Add those to the close-out, with a final clean-status check, while keeping `docs/releasing.md` authoritative for the tag and release sequence.

## What held under refutation

- **Triage does not require the parked stakes field.** `02av` tests a real non-material symlink finding and a material control against the committed four-field taxonomy and per-entity value ACs (`ensign-finding-triage-disposition.md:111-124`). The classifier distinguishes them, and `.roborev.toml` currently carries the cited fields.
- **The template/refit mechanism is proportionate.** `2ae` proved the current version-only refit control and a content-bearing regenerated diff (`template-rigor-propagation/index.md:172-182`). It uses existing full-body generation plus `diff`, not a new command.
- **Runtime semantics have an honest independent source.** `841` routes each in-scope claim to Go source, an existing build fixture, or an existing live lane, retains only structural checks with discriminators, and found the `member_spawn` divergence rather than hiding it (`contractlint-codex-runtime-semantics-retirement.md:44-69,98-109`).
- **The independent cleanup wave is low-collision.** `85`, `ht`, and `841` have disjoint surfaces and can lead. The eight test fixes retain named behavior checks and a mutation matrix; they do not widen production scope.
- **The re-lock removed most speculative machinery.** `bw`, `z7`, `02av`, `az`, and `2ae` all use existing delivery, review, fixture, or diff paths. The remaining defects are promises and wiring left behind by the earlier shape, not evidence that the reduced mechanisms require a new binary subsystem.

## Re-review and closure condition

This verdict becomes **READY** after one independent closure pass confirms all of the following in the working tree:

1. `bw` and `index.md` agree on a prose-only recorded-decision barrier and its live/fixture proof; neither promises a deferred refusal command or unit test.
2. The evidence DoD matches `az`'s resolved Edit-D decision and contains no unowned `testlint` claim.
3. The index's DoD, layer map, and Out of scope describe the eight-member query exactly; removed backlog work is not represented as deferred sprint membership or as an owned 0260 deliverable.
4. The dispatch contains the shared-edit/wave table, the composed `bw`+`02av` gate, the five-scenario/two-runtime pre-cut proof, a resolved catalog artifact home, and the repository completion gates.
5. The live query still returns the intended eight members, every pending gate digest matches the body it approves, and no new member or binary enforcement surface entered through the folds.

Re-review need not repeat the per-member ideation gates. It should diff only the sprint index, this dispatch, the folded `bw` body/gate record, and any metadata changed to close Polish items, then rerun the two membership queries and verify the four Material contradictions are gone.
