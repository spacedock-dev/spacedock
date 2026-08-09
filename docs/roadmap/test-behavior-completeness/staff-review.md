# Test-behavior-completeness staff review

Date: 2026-08-09

Verdict: **NOT READY**. The sprint has five Material findings. The final shape
needs three new product-repair tasks and two task folds before gate lock.

This review excludes `g3` and all durable-decisions work. It does not approve
product changes or create task entities.

## Evidence read

The review read the shaping handoff, the roadmap rules, the development workflow,
the live guide, and the complete desired registry. It also read every active
sprint task except `g3`.

The latest ideation records were read for `ts`, `98a`, `6x5`, `zh`, `9a`, and
`xp6`. The banked `0a` design was also read. The source bindings and durable
assertions were checked against those records.

The membership source is this query:

```bash
spacedock status --workflow-dir docs/dev \
  --where sprint=test-behavior-completeness \
  --where 'sprint-readiness != defer' \
  --archived
```

## Material findings

### M1 — Strict XFAIL has no rule for two semantic failures

`ts` keeps infrastructure failures outside the semantic grade. That boundary is
sound. Its result matrix also makes XPASS fail the lane.

The design does not define the result when two durable assertions fail. The gate
state assertion and command-log assertion run in sequence. A first-match rule can
hide a second, different product failure as XFAIL.

Fold this rule into `ts` before its gate:

- Run all durable semantic assertions after infrastructure succeeds.
- XFAIL only when the expected code is the sole semantic failure.
- XPASS when no semantic failure exists for an XFAIL binding.
- FAIL when any other code exists, including an additional code.
- Emit the existing metrics record after classification.
- Never send launch, auth, timeout, fixture, parse, state-read, or metric errors
  through the classifier.

This fold also defines exact candidate order. The first `ts` landing must include
the Sonnet and Codex `default-headless-gate-stop` XFAIL records. A helper-only
landing is not valid.

### M2 — The rejection task combines two mechanisms

`zh` proposes one recorder-order change. The Pi evidence supports this mechanism:
the recorder ran with two entries before the complete four-entry log existed.

The Codex evidence does not reach that boundary. It stops after the first rejected
validation and prepares an ordinary gate. The task itself calls this result
`rejection-flow-not-completed` and excludes it from recorder XFAIL evidence.

Recarve `zh` to the stable recorder-publication failures only. A Sonnet or Opus
cell can join only after one exact run supplies a stable recorder code. A passing
cell can remove its TODO after exact-candidate evidence.

Create a separate proposed task with this scope:

**Proposed slug:** `continue-codex-rejection-after-first-validation`

**Scope:** Route the first rejected Codex validation through correction,
independent validation/2, and the fresh final gate. Do not change recorder bytes,
round format, or the Pi recorder-order repair.

**Acceptance criteria:**

- The Codex `rejection-flow` cell starts from strict XFAIL code
  `rejection-flow-not-completed` on a committed baseline.
- The repaired exact candidate leaves two implementation reports, two validation
  reports, and one fresh open final gate.
- XPASS fails before source removes the binding.
- The final Codex run passes after binding removal.

**Dependencies:** `ts` first. Land after the `zh` recorder repair if both tasks
touch feedback-flow text.

**Visible value:** Codex completes the correction journey instead of stopping at
the first rejected candidate.

**Estimate:** 3 existing files, about 18 insertions and 8 deletions, net `+10`.
Tolerance: one file and 12 net lines. Ideation must revise this estimate after the
required live spike.

### M3 — The two Pi failures do not share one owner or mechanism

The Pi `gate-guardrail` cell fails its command-log boundary. It prepares the gate
without a later successful state commit. The existing gate-lifecycle skill already
requires prepare, commit, reread, and present in that order.

The Pi `default-headless-gate-stop` cell fails the final entity boundary. The
entity is not held at the open validation gate. This result does not prove the
missing-worker mechanism that `98a` owns. `98a` explicitly excludes Pi.

Do not put either repair in `xp6`. Do not combine the two Pi failures.

Create this first proposed task:

**Proposed slug:** `commit-pi-gate-prepare-before-presentation`

**Scope:** Make Pi commit and reread a successful prepared gate before the root
session presents it. Reuse `gate prepare`, `state commit`, and the existing room.
Do not change gate storage or command grammar.

**Acceptance criteria:**

- The Pi `gate-guardrail` cell starts from committed strict XFAIL code
  `gate-prepare-state-commit-missing`.
- The exact command log orders one successful prepare before one successful state
  commit and state-head read.
- The prepared Briefing stays open and unchanged after presentation.
- The repaired exact candidate passes after binding removal.

**Dependencies:** `ts` first. It can ideate with the next Pi task, but their live
runs must use isolated artifact roots.

**Visible value:** A Pi operator receives a review bound to committed gate state.

**Estimate:** 3 existing files, about 20 insertions and 8 deletions, net `+12`.
Tolerance: one file and 12 net lines.

Create this second proposed task:

**Proposed slug:** `hold-pi-default-headless-validation-gate`

**Scope:** Repair the Pi-only path that fails to leave the entity at the open
validation gate. First identify which final-state clause fails. Do not assume the
Sonnet and Codex worker-dispatch mechanism.

**Acceptance criteria:**

- A committed strict-XFAIL baseline names one refined semantic code for the exact
  failed clause in `assertGateHeld`.
- The Pi journey dispatches and completes required work, reaches validation,
  binds one open Briefing, and stops without decision or successor dispatch.
- XPASS fails before binding removal.
- The exact Pi candidate passes after binding removal.

**Dependencies:** `ts`, then `98a`. Running after `98a` proves whether the shared
worker guard changes the Pi symptom.

**Visible value:** A headless Pi run stops at the first open validation gate with
durable evidence.

**Estimate:** spike allowance of 3 existing files, about 24 insertions and 10
deletions, net `+14`. Tolerance: one file and 14 net lines. The ideation gate must
reject this estimate if the mechanism differs.

### M4 — Current ownership and file order disagree

The source still assigns all three `smallest-sufficient-mechanism` TODO rows to
`9a`. The recarved design assigns this journey to `6x5`.

The `6x5` XFAIL baseline must change those owners before its product edit. The
mutable owner join checks active state, but it cannot detect the wrong active
task. The sprint review must therefore check the semantic owner mapping.

Six tasks can edit `internal/ensigncycle/shared_live_runner_test.go`. Three tasks
also edit `skills/first-officer/references/fo-dispatch-core.md`. Parallel product
branches without an ordered rebase can overwrite source bindings or contract
rules.

Use the landing order in the sprint index. Each repair branch must rebase onto
the prior landing before it commits its XFAIL baseline. Run exact-candidate live
proof after that baseline commit and before product bytes.

### M5 — The capstone cannot carry repaired or deferred product scope

`xp6` is correctly evidence-only. Its product estimate is net zero. Keep this
boundary.

The Codex `withdrawn-gate-recovery` probe passed. This evidence permits `xp6` to
remove that TODO after an exact-candidate rerun. Task `47g` stays outside this
sprint. No deferred task must remain in scope to remove a stale TODO.

The two Opus rows remain TODO only until an authenticated run exists. Missing
local authentication is an execution-path reason, not product evidence.

`xp6` runs last. It removes passing bindings, preserves honest TODO rows, and
leaves product repairs with their own tasks.

## Claim-by-claim refutation result

| Claim | Result |
| --- | --- |
| Each task has independent visible value. | Refuted for `zh` until Codex is split. Confirmed for the other current tasks. |
| `ts` stays a small test mechanism. | Conditional. Its `+210` cap is acceptable only with the sole-failure rule and no second artifact. |
| Infrastructure cannot appear as XFAIL. | Conditional pass. The written boundary is sound, but multi-error behavior needs the M1 fold. |
| XPASS cannot remain green. | Confirmed by `ts` AC-2. The Commander must keep the binding through the XPASS run. |
| TODO and XFAIL owners join to active tasks. | Conditional. `ts` extends the join, but `6x5` must replace the stale `9a` owner mapping. |
| TODO cells change only after execution evidence. | Confirmed by the proposed order. Opus stays TODO until authenticated execution. |
| Product repairs start from committed XFAIL. | Conditional. Every branch needs a baseline commit after its latest rebase. |
| Shared surfaces have one landing order. | Refuted in the task bodies. The index now supplies one order. |
| The capstone contains no product repair. | Confirmed if M3 creates two separate Pi owners. |
| Durable-decisions work is excluded. | Confirmed. `g3` is excluded from review, dispatch, and merge authority. |

## Net-line and implementability audit

`0a` is large because it restores a complete retained CI job. Its limit remains
8 files and 380 insertions after its approved design reset.

`ts` has a hard `+210` net cap. This is the largest acceptable common mechanism.
A protocol, new artifact, copied scenario map, or host switch is a design reset.

`98a`, `6x5`, and `zh` each have small product deltas. Their live proof is larger
than their product change, but it reuses the common suite.

`9a` estimates net `+228`. Most growth is focused CLI evidence. Implementation
must prefer existing test files and stay below the declared 25% tolerance. A new
controller, lifecycle model, or test-only state machine is a design reset.

`xp6` remains net zero in product source. The three proposed tasks add about 36
net lines before their ideation refinements.

No proposed interface lacks command grammar or stored-field domains. The three
new tasks reuse current commands and formats. Their ideation spikes must name the
exact semantic branch before gate approval.

## Required closure before gate presentation

1. Fold M1 into `ts`.
2. Recarve `zh` and create the separate Codex repair task.
3. Create and ideate the two Pi repair tasks.
4. Fold the owner and merge order from M4 into affected tasks.
5. Update the membership query result and net-line budget.
6. Record fresh ideation gate attempts for all changed or new tasks.

After these actions, one independent delta review must make sure that each
Material finding is present in durable task state. The sprint can then become
Commander-ready.
