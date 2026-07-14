---
title: Generalize the smallest-concrete-mechanism rule
status: ideation
source: "Captain request 2026-07-13."
started: 2026-07-12T23:19:48Z
completed:
verdict:
score:
worktree:
issue:
id: 5wmx24qwygh711wsmmqb1qwj
milestone: 0.26.0
---

Capture a concise, reusable anti-overengineering rule in the shared Spacedock
contract, with a developer-specific adaptation as an example rather than a
second universal policy.

## Problem

The contract already says to choose the smallest sufficient mechanism, but its
dispatch wording and write-authority wording leave a live edge ambiguous. On
2026-07-13 the captain named an existing task and supplied its exact amendment.
The FO dispatched a worker whose assignment repeated the target and mutation.
There was no remaining judgment, fan-out, required isolation, independent
adversarial verification, or safety boundary. The expected trace was one direct
edit and zero dispatches; the observed trace contained a needless dispatch.
That run is the RED baseline, not evidence that the proposed repair works.

The same missing discriminator permits implementation layers with no present
need. A recent Safehouse change represented one future tmux case with a registry
and custom parser/merger. Nothing required the implementer or validator to name
a second current case or a boundary that made the direct scalar mechanism fail.

Policy placement is another instance, not a separate universal rule. A shared
contract should not absorb one workflow's prose-polish routing policy merely
because the shared layer supplies teammate mechanics. The concrete routing
integration remains owned by `workflow-owned-prose-polish-routing`; this task
owns only the narrowest-policy-owner rule that reveals the wrong layer.

The 7h incident exposed the same gap after dispatch. A disposable Zellij/WASM
smoke grew into a test-owned PTY/session controller with foreground
process-group checks, raw-byte canaries, lease publication, signal coordination,
readiness races, and cleanup state. It consumed about 4.25 active hours across
13.5 elapsed hours; the final packet failed 3/3 on harness readiness while
focused profile runs passed, and nothing merged. The successful replacement
used the supported boundary: isolated tmux -> real Zellij -> literal keys ->
captured screen/native state -> cleanup. The missing general rule is not “never
use a PTY.” It is that a mechanism must keep naming the value it serves and the
simpler route it beat; when the mechanism fails but that value remains reachable
more simply, the feedback loop resets scope instead of repairing the mechanism.

## Proposed direction

Replace the existing `Smallest sufficient mechanism` paragraph in
`skills/first-officer/references/first-officer-shared-core.md` with one compact
rule that covers both abstraction and routing:

| Condition | Mechanism |
| --- | --- |
| Declared commissioned stage | Dispatch through the standing loop. |
| Other work | Classify FO write authority first. |
| Captain grants the exact target and exact mutation | Authorize only that mutation. |
| Remaining judgment, fan-out, required isolation, independent adversarial verification, or safety boundary remains | Dispatch; otherwise apply the authorized mutation directly. |

The two-consumer threshold below governs only a new reusable or general-purpose
abstraction. It is never a dispatch threshold.

> **Smallest concrete mechanism.** Prefer the direct mechanism that achieves the
> current outcome. Do not introduce a new reusable or general-purpose
> abstraction—a generic layer, workflow template, agent role, recurring stage or
> gate, reusable state/configuration surface, or extensibility scheme—without two
> present concrete consumers. One consumer can justify reusable structure only
> for a proven safety, compatibility, required isolation, or external-contract
> boundary; cite it. Separately, before choosing a heavier mechanism for one
> current task, name the direct alternative and the remaining judgment, fan-out,
> required isolation, independent adversarial verification, or safety boundary
> that makes it fail. An exact captain-authorized target and mutation with none
> of those constraints is a direct edit, not a worker assignment. This check does
> not fire for a commissioned stage: commissioning has already fixed dispatch,
> so the declared loop dispatches every ready task without per-task
> justification. Nor does it require narration for the direct choice or other
> deterministic filing work. Put policy at the narrowest layer with its present
> consumers; workflow-specific behavior stays in that workflow or mod until a
> second current consumer needs the shared policy. For example, a developer uses
> a scalar or direct branch for one case and introduces a reusable type,
> collection, parser, or extensibility layer only for two present consumers or
> one cited boundary.

This is one decision rule with two distinct questions, not a scoring system.
The two-consumer check applies only to proposed reusable or general-purpose
structure; it does not prohibit a one-off agent, stage, gate, plain function,
ordinary file, or other task-specific mechanism when one task has remaining
judgment, fan-out, required isolation, independent adversarial verification, or
a safety boundary. The direct-versus-heavy check then chooses the mechanism for
that task. Neither list is a vocabulary escape hatch such as “substantive” or
“future proof.”

The First Officer wording remains general. At dispatch, while assessing a
completed implementation, and before routing a rejection, it requires a new
mechanism to name the value AC it serves, the simplest available alternative,
and the concrete reason that alternative is insufficient. When rejection shows
the mechanism failed while the value remains reachable by the simpler route,
the FO records rejected experimental evidence and escalates a scope/design
reset; it does not increment the feedback cycle or redispatch implementation to
repair the mechanism. Runtime- or test-specific tripwires do not belong in the
shared FO contract.

### Decision order and write authority

Use the following order so the general rule cannot override a stronger existing
boundary:

1. If the work is a declared stage of a commissioned workflow, dispatch it by
   the standing loop. The stage worker remains correct for this ideation task
   and for implementation or validation that requires its worktree or role.
2. Otherwise classify the intended write under `fo-write-core`. Routine filing
   through `spacedock new`, state transitions, and other already-allowed FO
   mutations stay direct without new ceremony.
3. An instruction that supplies the exact target and exact mutation is the
   existing exact captain grant for that mutation; it needs no magic word such
   as “override.” Refine `fo-write-core` to say that this grant can authorize the
   exact otherwise-off-limits body or product edit, while a broad “fix it
   directly” instruction cannot.
4. For that exact authorized mutation, dispatch only if remaining judgment,
   fan-out, required isolation, independent adversarial verification, or a
   safety boundary blocks the direct edit. Name that blocker once at the climb.
   If none exists, edit directly. The grant changes authority for only the named
   mutation; it does not widen FO ownership afterward.

Thus an exact amendment to an existing backlog body can be direct when the
captain supplied both target and replacement. An underspecified body rewrite
still belongs to a stage worker. A blocked product path still routes to a worker
under a generic direct-work prompt, but an exact path plus exact mutation can use
the existing override. A required isolated worktree or independent adversarial
check still dispatches even when the desired text is exact.

### Concrete documentation diff

In `skills/first-officer/references/fo-write-core.md`, keep the current mutation
table and replace the prose below it with this authority clarification:

> A blocked target routes through a worker unless the captain grants the exact
> target and exact mutation. Supplying both is the grant; no special “override”
> phrase is required. The grant authorizes only that mutation and may also cover
> otherwise-off-limits entity-body content. Broad direction such as “fix it
> directly,” a target without the mutation, or an FO's self-declared override is
> not a grant. After authority is established, apply the smallest-concrete-
> mechanism rule: direct work stays direct when no remaining judgment, fan-out,
> required isolation, independent adversarial verification, or safety boundary
> remains.

Update the `override` row's rule to match: `Exact captain-authorized target and
mutation; scoped to that mutation only.` Do not add a new class, policy engine,
mandatory decision log, per-edit template, or runtime flag.

In `skills/commission/references/templates/development.md`, add the
developer-specific realization. Tests orchestrate and observe the supported
runtime rather than owning a parallel protocol, daemon, lease, lifecycle state,
recovery loop, or process controller. For terminal multiplexers, `setsid`,
process-group control, raw PTY writes, or a second lifecycle supervisor trigger
explicit architecture review and an existing real-terminal harness is tried
first. Require one end-to-end run through the simplest substrate before
approving an enabling layer; if no visible product proof appears within 90
minutes, stop for architecture review rather than adding coordination. A
mechanism-architecture rejection recommends a design reset, while a product
defect still routes through the ordinary implementation feedback cycle.

When either the FO or the development workflow recommends a “narrow fix,” it
must label an **outcome defect** (the delivered behavior or state fails the value
AC) or an **evidence defect** (the test, harness, or observation boundary cannot
validly establish the AC), cite the affected AC, and name the exact failing
boundary. A correction spanning both paths, or an evidence correction that adds
another controller or lifecycle layer, is a design reset rather than a narrow
fix.

In `docs/specs/scenario-testing-principles.md`, revise only the
`smallest-sufficient-mechanism` seed description so it names all three arms:
exact authorized mutation -> direct; remaining judgment, fan-out, required
isolation, independent adversarial verification, or safety -> dispatch;
commissioned stage -> standing dispatch. The generic policy-placement example
belongs in this task's fixture; prose-polish routing details remain in
`workflow-owned-prose-polish-routing`.

## Acceptance criteria

**AC-1 (observed busywork turns green): Given an exact captain-authorized target and mutation with no remaining judgment, fan-out, required isolation, independent adversarial verification, or safety boundary, the FO applies the mutation directly and dispatches zero workers for it.**

Verified by: extend the shared smallest-mechanism fixture with the 2026-07-13 RED shape
and grade the host trace plus resulting file bytes. The correct trace has one FO
edit, zero matching dispatches, and the exact replacement on disk. The recorded
RED trace must fail, and deleting the new discriminator from the operating prompt
must make the live/captured arm fail again.

**AC-2 (real constraints still dispatch): The same exact desired mutation routes to exactly one worker when the prompt adds required isolation or independent adversarial verification; the FO does not edit the target itself.**

Verified by: add the paired constrained arm to the fixture and grade one matching worker
dispatch, zero FO edits to that target, the exact expected bytes, and a target
edit plus commit attributable to the dispatched worker between its start and
completion. The worker-scoped transcript and worker-branch commit are the
attribution evidence; FO narration that a worker owns the result is not. Offline
negative cases must reject direct FO editing, dispatch-only narration with no
worker edit/commit, a commit attributable to the FO, wrong resulting bytes, and
suppression of the required dispatch.

**AC-3 (commissioned stages remain mandatory): Every ready task in a commissioned stage is dispatched by the standing loop without a per-task smallest-mechanism justification, including the genuine commissioned ideation and isolated-worktree cases.**

Verified by: retain the current two-ready-task scope guard in
`shared_smallest_mechanism_*`; its positive trace dispatches both, while existing
negative traces for a suppressed dispatch and a per-task justification stay red.

**AC-4 (authority is exact and low-ceremony): Exact target-plus-mutation captain authorization permits only that otherwise-off-limits edit without magic override wording; a broad direct-work instruction, target-only instruction, self-declared override, or adjacent mutation still routes through the existing worker guard. Routine `spacedock new` filing and allowed state transitions gain no new justification step.**

Verified by: extend `fo_product_edit_guard_test.go` for both hosts with positive exact
target-plus-mutation user messages and negative broad, target-only,
self-declared, and adjacent-file/mutation transcripts. The positive case must
compare the target's resulting bytes with the exact authorized replacement.
Add a same-file negative whose trace claims the exact grant but writes different
bytes; it must fail even though target and authorization narration match. Keep
the existing filing and allowed-state tests green; add no assertion that
searches contract prose.

**AC-5 (one generalized rule, one illustrative developer case): The shared contract contains one smallest-concrete-mechanism rule with the evidence-based single-case exceptions; scalar-versus-registry developer wording appears inside it as an example, not as an independently triggered policy.**

Verified by: the shared behavior fixture presents a one-case scalar choice and a
two-current-case or proven-boundary variant. Captured host traces and resulting
fixture files must show the direct representation for one case and permit the
layer only in the qualifying variant. Add an explicit adversarial negative in
which the one-case arm produces the same user-visible result through a registry,
custom parser, or equivalent reusable layer; the grader must reject that trace
and on-disk representation. A mutation that removes or bypasses the one-case
layer rejection must make this negative pass incorrectly and therefore fail the
mutation test. Review the implementation diff for one resident rule; do not add
a prose-grep or style-lint assertion.

**AC-6 (narrowest policy owner): A policy with one workflow-specific consumer remains in that workflow or mod while the shared layer retains only its generic mechanism; a second present consumer may promote the policy. No prose-polish routing integration moves from `workflow-owned-prose-polish-routing` into this task.**

Verified by: add a generic policy-placement arm to the shared fixture and grade the
paths actually changed: the one-consumer run changes only the workflow-local
policy file; the two-consumer variant may change the shared fixture policy. A
negative trace that writes the one-consumer policy into the shared file must
fail.

**AC-7 (failed mechanism resets scope): When a rejected enabling mechanism fails but its served value remains reachable through a simpler alternative, the FO records rejected experimental evidence, adds no feedback cycle, dispatches no implementation repair, and surfaces a scope/design reset; a paired product defect still follows the ordinary feedback route.**

Verified by: extend the existing smallest-mechanism and rejection-flow scenario
infrastructure with paired mechanism-failure and product-defect arms. Grade
dispatch calls, feedback-cycle state, and the resulting routing decision. A
mechanism failure bounced to implementation and a product defect suppressed as
a design reset are mandatory red controls; FO narration alone is not evidence.

**AC-8 (developer adaptation stays in the development template): A freshly commissioned development workflow carries the observer-not-second-implementation harness boundary, real-substrate-first rule, architecture tripwires, timebox, and design-reset exception, while the shared FO contract contains only the general mechanism/value rule.**

Verified by: commission a development workflow and inspect the generated stage
structure for the dev-specific obligations, then drive the paired AC-7 behavior
through the existing scenario harness. Contract structure may prove placement;
it cannot substitute for the routing behavior. Review the shared FO diff to
reject PTY, Zellij, tmux, or 90-minute development-specific wording there.

**AC-9 (narrow-fix diagnosis is explicit): Every narrow-fix recommendation labels exactly one failing path—outcome defect or evidence defect—cites its value AC, and names the failing boundary; a mixed-path correction or evidence change that adds another controller/lifecycle layer is surfaced as a design reset.**

Verified by: paired gate fixtures accept one precise outcome-path recommendation
and one precise evidence-path recommendation, then reject an unlabeled “narrow
fix,” a recommendation with no served AC, a mixed-path recommendation, and an
evidence fix that grows another supervisor. Grade the recommendation fields and
routing outcome, not wording in the contract file.

## Test plan

Extend the existing `smallest-sufficient-mechanism` shared scenario rather than
adding a parallel harness. Its host trace, on-disk fixture, offline negative
tests, and live runners already cover direct edits and commissioned dispatch.
Estimated cost is medium: distinguish direct, constrained, abstraction, and
policy-placement targets in both host extractors and add paired fixture
variants. No new dependency or runtime surface is needed.

Extend the existing rejection-flow fixture for AC-7/AC-9 and the existing
commission-template fixture for AC-8. Do not create a dedicated policy daemon,
controller, lifecycle, or live harness. The general FO behavior is proven by
routing/state outcomes; the template-specific placement is a structural check.

Run focused tests first:

```text
go test ./internal/ensigncycle -run 'Smallest|FOProductEditGuard'
```

Then run the repository gates required by `AGENTS.md`:

```text
gofmt -w ./cmd ./internal
go test ./...
go test ./... -race
```

The behavior fixture and offline mutation/negative cases are required. Run the
existing live shared scenario for Claude and Codex when credentials are
available because the claim concerns model routing; record an unavailable live
lane as unavailable, never as green. Do not accept contract substring checks,
self-review of the new prose, or the 2026-07-13 RED transcript as proof of the
fix.

Ideation-stage verification on 2026-07-13:

```text
$ go test ./internal/ensigncycle -run 'Smallest|FOProductEditGuard'
ok  github.com/spacedock-dev/spacedock/internal/ensigncycle  0.403s

$ go test ./...
exit 0; all listed packages passed

$ go test ./... -race
exit 0; all listed packages passed

$ git -C docs/dev/.spacedock-state diff --check -- smallest-concrete-mechanism-rule/index.md
exit 0; no output
```

The clean race result is the standalone rerun. An earlier race run concurrent
with the non-race suite failed `TestSonnetTeamDeleteHangReplay`; that exact test
passed immediately in isolation under `-race`, and the subsequent standalone
full race suite passed.

## Spike determination

No spike needed. The proposal reuses proven mechanisms already exercised in
this repo: exact-override transcript grading in
`fo_product_edit_guard_test.go`, host-neutral mechanism traces and on-disk
fixtures in `shared_smallest_mechanism_*`, and the commissioned-dispatch scope
guard in the live shared scenario. Implementation extends those inputs and
graders; it does not assume a new parser, handoff, file format, or CLI flag.

## Stage Report: ideation

- DONE: Define the generalized direct-versus-dispatch rule so the observed busywork case fails while genuine commissioned or constrained work still dispatches.
  The decision order makes exact authorized closed mutations direct, while commissioned stages and named isolation, verification, or safety constraints still dispatch.
- DONE: Resolve the filing, write-authority, and commissioned-stage boundaries without adding ritual or a policy engine.
  The proposal refines the existing exact grant, preserves `spacedock new` and state mutations, and adds no classifier, flag, log, or required direct-choice narration.
- DONE: Produce behavior-first acceptance criteria and tests.
  Six criteria specify trace and on-disk outcomes, paired positives/negatives, mutation coverage, focused tests, and the full repository gates.
- DONE: Run the repository verification required for this stage.
  Focused tests, `go test ./...`, and a clean standalone `go test ./... -race` passed; an earlier concurrent race run's replay failure also passed isolated.
- SKIPPED: Route the captain-facing draft through the workflow comm officer before commit.
  `/root/comm_officer` remained `pending_init` beyond the declared two-minute bound; proceeded with a behavior-preserving self-edit as the routing contract directs.

### Summary

The ideation generalizes the rule across abstraction, dispatch, write authority,
and policy placement while keeping one precise discriminator: exact target plus
exact mutation plus no remaining constraint means direct work. It reuses the
existing cross-runtime smallest-mechanism and FO write-guard harnesses, keeps
prose-routing integration in its separate task, and records the bounded comm-
officer fallback.

## Stage Report: ideation (cycle 2)

- DONE: Define the generalized direct-versus-dispatch rule so the observed busywork case fails while genuine commissioned or constrained work still dispatches.
  The resident rule now separates reusable-abstraction consumers from one-task mechanism blockers and uses one normalized direct-edit discriminator.
- DONE: Resolve the filing, write-authority, and commissioned-stage boundaries without adding ritual or a policy engine.
  A four-row decision table orders commissioned dispatch, write classification, exact scoped authority, and the remaining-constraint test.
- DONE: Produce behavior-first acceptance criteria and tests.
  AC-2 requires a worker-attributed edit and commit, AC-4 compares exact bytes and rejects a same-file wrong mutation, and AC-5 adds an adversarial one-case layer negative plus mutation.
- DONE: Repair the independent review's reusable-layer contradiction.
  Two present consumers govern only new reusable or general-purpose abstraction; one-off agents, stages, or gates remain valid for concrete judgment, fan-out, isolation, verification, or safety needs.
- DONE: Route the captain-facing draft through the workflow comm officer before commit.
  `/root/comm_officer` found no structural rewrite necessary; its `required isolation` normalization, decision table, and concise harness wording were incorporated without changing behavior.
- DONE: Re-run focused verification after the cycle-2 repair.
  `go test ./internal/ensigncycle -run 'Smallest|FOProductEditGuard'` returned `ok ... (cached)` and `git ... diff --check` exited 0 with no output.

### Summary

Cycle 2 removes the only contradictory universal reading and makes each
behavioral proof inspect the actual mechanism: worker attribution, commit,
resulting bytes, or forbidden layer. Commissioned stages still dispatch,
closed exact mutations remain direct, cs retains prose-routing integration, and
the comm officer's behavior-preserving polish is folded into the proposal.

## Stage Report: ideation (captain amendment, cycle 3)

- DONE: Fold the 7h mechanism failure into the general rule without putting test-specific wording in the FO contract.
  The FO owns mechanism/value tracing and the design-reset exception; the development template owns PTY/process/lifecycle tripwires and the 90-minute enabling-task timebox.
- DONE: Distinguish outcome defects from evidence defects in narrow-fix recommendations.
  AC-9 requires one labeled path, the affected value AC, and the exact failing boundary; mixed-path or controller-growing corrections reset design.
- DONE: Keep normal product-defect feedback routing intact.
  AC-7 requires paired behavior: mechanism failure resets scope, product failure still routes to implementation.
- DONE: Reuse existing proof infrastructure.
  AC-7/AC-9 extend smallest-mechanism plus rejection-flow; AC-8 uses the commission-template fixture. No new harness or controller is proposed.

### Summary

The captain amendment generalizes the Zaphod learning at the correct layers:
strong value-to-mechanism and reset semantics in the FO contract, concrete
harness discipline in the development template, and behavior proof through the
existing scenario infrastructure.
