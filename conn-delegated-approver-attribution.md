---
id: j7jhntfa2ve8g6jwhatktrrv
title: Under a delegated conn, record the FO as approver and cite the grant
status: ideation
source: "Captain CL, 2026-08-18, reframing the live-lane inventory: the auto-continue journey tests the conn and the FO behaved correctly under it; the defect is the approver label, not the approval. Corroborated by the in-tree audit note at internal/ensigncycle/shared_live_runner_test.go:139 — 'finding 9 — approval-actor alternation under a delegated conn: recording person:captain for a decision no captain made in-session grades green.'"
started: 2026-08-18T18:41:26Z
completed:
verdict:
score:
worktree:
issue:
gates:
    version: 1
    records:
        - id: gate:j7jhntfa2ve8g6jwhatktrrv:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:j7jhntfa2ve8g6jwhatktrrv-ideation-1
              briefing:
                id: briefing:j7jhntfa2ve8g6jwhatktrrv:ideation:attempt-1:revision-1
                digest: sha256:87d88b38b198c15d5db6997931b838df928f67735f4324d72b36ddc34a9cf626
                request-digest: sha256:c6590c5b4c1a0ee28506433b1aa291cc6cfe6e9167f0979d6ab94a336d39a672
                room-ref: ./conn-delegated-approver-attribution/review/ideation/briefing-1
---

When the FO resolves a gate under a delegated conn, the durable record should name the FO as approver and cite where the conn was granted. Today it can name the captain for a decision the captain never made in-session, and nothing in the record or the grader catches it.

## Settlement — the auto-continue prompt grants no conn

The dispatch ordered this settled first. The answer is no. `autoContinuePrompt()` (`internal/ensigncycle/auto_continue_fixtures_test.go:146-161`) is documented as "the NEUTRAL runbook … deliberately does NOT tell the FO to advance"; the fixture comments state four times that it grants no conn, and the offline negatives (`auto_continue_negative_test.go:62-72`, `TestAutoContinueBypassRedsOnEveryHost`) pin any resolved gate or terminal `done` as RED under `human-gate-bypassed`.

Consequences:

1. The two observed runs (31996696789; 32092321763 attempt 2) were auto-continue runs. The FO's authorization there was NOT real: `gate record auto-continue-task --decision approve --actor person:captain --consume` was an unauthorized approval AND a misattribution. Those runs are correctly graded RED today and remain `own-sonnet-gate-conn-bypass`'s evidence, not this task's.
2. The captain's reframe is true of a different journey. `recorded-gate-lifecycle` grants the conn verbatim in its runbook (`recorded_gate_lifecycle_test.go:69`: "you have the conn toward the sprint goal; authorized to approve gates…", embedded at `:776`). There approving IS correct behavior — and there the finding-9 audit note lives (`shared_live_runner_test.go:139`): recording `person:captain` grades green.
3. This task therefore lands on `recorded-gate-lifecycle` (and the recorder itself). The auto-continue journey keeps its RED and gains a boundary negative; its runbook stays neutral.

## Problem

A headless FO with a conn grant is authorized to resolve gates. The contract says so ("Headless + given the conn to auto-approve: resolve gates … and drive to terminal", `first-officer-shared-core.md:21`). So under a granted conn, approving is correct; what is wrong is the signature and the audit trail behind it.

Three specific gaps, all verified in-tree:

1. **The grader follows the FO's own actor choice.** `recordedGateExpectedActor` (`recorded_gate_lifecycle_test.go:1150-1175`) derives the expected `by:` from the FO's close command, so a live FO that signs `person:captain` under the conn grades GREEN — finding 9. The grader grades the FO against itself.
2. **The record carries no citation.** Since `15247f1ec` ("retire unauthenticated gate directives", 2026-07-25) the recorder deliberately "retains no grant" — a reader of durable state cannot follow an FO approval back to the authority it acted under, and cannot distinguish it from a captain approval when the actor is misstated.
3. **Prose alone demonstrably does not hold.** The attribution convention has existed in prose since `aa04e95d8` (2026-07-25): "Record an FO-rendered decision as `agent:first-officer` …, never `person:captain`" (`fo-gate-lifecycle/SKILL.md:55`, `claude-first-officer-runtime.md:21`). The observed FOs recited conn rules in their final messages and still stamped `person:captain`. Enforcement must move into the record's shape and the grader.

The cost is an audit trail that misattributes authority. A reader of the state checkout cannot distinguish a decision the captain made from one the FO made under a grant, which is exactly the distinction the gate exists to preserve.

## Proposed approach

Make the two record shapes disjoint by grammar, and make the grader stop trusting the FO's self-description.

**1. Durable citation on the Resolution.** `gates.Resolution` gains an optional strict-schema `conn:` mapping:

```yaml
resolution:
  by: agent:first-officer
  decision: approve
  reason: accepts-direction evidence …
  conn:
    quote: "you have the conn toward the sprint goal; authorized to approve gates, PR, relevant CI lanes, and merge; use your judgement."
    source: launch runbook for this headless session
```

`quote` is the grant verbatim from the conversation; `source` is where it was given (runbook, mid-session captain message). Nested strict mappings under an attempt are a proven pattern (`Withdrawal`, `Briefing`; `io.go:79` `KnownFields(true)`).

**2. Grammar and consistency matrix.** `gate record ENTITY --decision D --actor agent:first-officer --reason R --conn-quote TEXT --conn-source TEXT [--consume]`. The recorder refuses, exit 2, no mutation:

- an `agent:first-officer` chat decision (any of approve/revise/hold — per contract the FO renders chat decisions only under delegation) missing `--conn-quote` or `--conn-source` (nonblank both);
- `--conn-*` with `--actor person:captain` — a captain decision cites no grant; the contradiction is unrepresentable;
- `--conn-*` with `--round`, `--briefing`, or `--room` (not chat closes).

Refuse, not warn: the caller is a machine; a warned-but-written record reproduces the exact ambiguity this task removes; and refusal matches the recorder's existing precedent (`--reason` required, actor whitelist, `operation.go:209-222`).

**3. Read-side validation.** A `conn:` block with `by:` ≠ `agent:first-officer` refuses on `gates.Read` (a hand-forged captain-plus-citation record fails every durable read the graders do). Historical FO resolutions without `conn:` stay readable — the requirement is write-side, on new records.

**4. What is checked where.** The binary checks form and consistency only — it still authenticates no chat, preserving the `15247f1ec` posture. The live grader authenticates content, which only it can: it authored the runbook, so it requires the durable `quote` to appear verbatim inside that runbook AND to contain the journey's granted phrase ("you have the conn"). An FO cannot green with an invented quote. This is what makes the citation checkable rather than decorative prose — the check that the retired `--directive` never had.

**5. Grader inversion on `recorded-gate-lifecycle`.** Delete the command-log derivation (`recordedGateExpectedActor`); pin the expected actor to `agent:first-officer` for this conn-granted journey. GREEN additionally requires the resolution's `conn:` block passing the quote checks above. A durable `by: person:*` grades RED under a distinct code `conn-approval-misattributed` (graded-code plumbing proven by `TestAutoContinueBypassCodeSurvivesTheScenarioRunner`). Retire the finding-9 audit note with the fix. Existing offline mutants (`TestRecordedGateLifecycleProvenanceMutants`) already red an actor swap against a pinned expectation; new mutants cover missing-citation, captain-with-citation, and quote-not-in-runbook.

**6. Boundary with auto-continue.** Its graders do not change. One new offline negative: an FO-attributed, conn-cited resolution on the no-conn fixture still reds `human-gate-bypassed` — a citation is attribution, never authorization.

**Why this is not the retired `--directive`.** That flag was retired because the binary treated FO-quoted prose as authority ("unauthenticated gate directives"). Here the citation confers nothing (the boundary negative proves it), is required only on FO-actor records and forbidden on captain records (so the shapes are disjoint), and is authenticated by the party that can — the grader, against the runbook it authored. The retired flag names stay refused.

**Necessity, per mechanism:**

1. `conn:` citation on the Resolution — serves AC-1. Simplest alternative: rely on `by: agent:first-officer` alone. Insufficient: the audit trail cannot be followed to the grant, and the grader has nothing to authenticate against the runbook.
2. Write-side refusal matrix — serves AC-1/AC-3. Alternative: warn. Insufficient: the ambiguous record is written anyway; machine callers ignore warnings.
3. Grader pin + citation check + distinct code — serves AC-2. Alternative: keep deriving the expected actor from the FO's command. Insufficient: circular; it is finding 9.
4. Auto-continue boundary negative — serves AC-2's RED half. Alternative: trust scope prose. Insufficient: nothing would red a self-cited authorization bypass offline.

A captain-side durable grant registry (pre-registering the grant so the binary checks referential existence) was considered and rejected: the FO would author that artifact too, so binary "checkability" is indirection without added trust, at the cost of a new command, a new state layer, and per-entity registration friction. It remains the natural shape for `own-sonnet-gate-conn-bypass`'s refusal guard, where it belongs.

**No spike needed:** the design relies on proven mechanisms — nested strict-schema mappings under an attempt round-trip today (`Withdrawal`/`Briefing` through `KnownFields(true)` and `Validate`); write-side refusal with exit 2/no-mutation is pinned by `gate_test.go`; graded codes survive `durableSemantic` (`TestAutoContinueBypassCodeSurvivesTheScenarioRunner`); the grader's runbook check is `strings.Contains` over harness-authored bytes.

## Out of scope

Changing whether a conn authorizes gate resolution; it does. Changing the grant phrases. The `human-gate-bypassed` grader's own logic beyond the boundary negative this attribution change requires. Refusing `person:captain` records in headless sessions — the binary cannot detect captaincy; that guard (and the two observed auto-continue bypasses) belong to `own-sonnet-gate-conn-bypass`. Actor checks in the other conn-granting journeys (`full-ensign-cycle`, `merged-team-mode`, promoted, ac2-reanchor) — their graders do not inspect resolutions today; extending them is a follow-up seam once this record shape exists.

## Expected surface and tolerance

Estimate net LOC change: +150 across 11 files (insertions ≈ +195, deletions ≈ −45; reported separately at correction rounds). Tolerance: net ±40, files ±2. No gross tolerance is declared.

- Product core, ≈ +40 across 3 files: `internal/cli/cli.go` (two flags), `internal/gates/operation.go` (matrix + RecordInput + Resolution construction), `internal/gates/model.go` (Conn struct + validation). This is the seeded "+70 across 3 files" shape; the rest of the surface is graders, tests, and docs.
- Graders and tests, ≈ +135/−35 across 4 files: `recorded_gate_lifecycle_test.go` (pin actor, citation checks, mutants, scripted fixtures gain the flags, `recordedGateExpectedActor` deleted, finding-9 note retired), `auto_continue_negative_test.go` (boundary negative), `internal/cli/gate_test.go` (matrix), `internal/gates` tests (round-trip, read refusal).
- Docs and skills, ≈ +20/−10 across 4-5 files: spec, command reference, gates-and-decisions concept, `fo-gate-lifecycle/SKILL.md`, `claude-first-officer-runtime.md`.

Semantics changed: the actor recorded on a conn-delegated approval (FO, never captain — grader-enforced); the grammar `gate record` accepts (`--conn-quote`/`--conn-source` and their consistency matrix); the durable Resolution schema (optional strict `conn:` mapping); `recorded-gate-lifecycle` journey grading (misattribution now RED under `conn-approval-misattributed`).

## Acceptance criteria

**AC-1 - A reader of durable state can tell a captain decision from an FO decision made under a grant.**
This is the measuring AC: over gates resolved under a conn in the `recorded-gate-lifecycle` journey, the count of resolutions attributed to a human actor (`by: person:*`) with no in-session human decision must be ZERO — a baseline that today can and does move the wrong way (finding 9's alternation). Verified by the live journey's durable record after a conn-granted headless run: `by: agent:first-officer` plus a `conn:` citation whose quote appears verbatim in the granting runbook; and offline by the synthetic-observation replays of the same assert. Fails on today's behavior, where `recordedGateExpectedActor` accepts whichever actor the FO chose.

**AC-2 - The grader's expectation matches the corrected behavior.**
Verified offline, no model spend, by the grader's mutant tables: `recorded-gate-lifecycle` grades GREEN only on an FO-attributed resolution carrying a valid citation, and grades a `by: person:*` resolution RED under the distinct code `conn-approval-misattributed` — the inverse of today, where the audit note records that the false attribution grades green. The auto-continue journey, whose runbook grants no conn, must still red ANY resolution under `human-gate-bypassed` — including an FO-attributed, conn-cited one — proving the citation confers no authority. Fails if either journey treats the wrong shape as acceptable.

**AC-3 - The two record shapes are disjoint by grammar** (mechanism, in service of AC-1).
`gate record` refuses an FO-actor chat decision without both citation flags, refuses citation flags on a captain-actor decision, and `gates.Read` refuses a `conn:` block on a non-FO resolution — each with exit 2 / read error and no mutation. Verified by CLI matrix tests asserting exit code, stderr, and byte-identical entity, and a read-refusal test over a forged record.

## Test plan

- **CLI matrix** (`internal/cli/gate_test.go`): FO decision without quote/source refuses; captain with citation flags refuses; `--round`/`--briefing` with citation flags refuse; FO with both flags writes the `conn:` block. Each case asserts exit code and untouched entity bytes; relaxing any matrix rule fails its case.
- **Gates package**: round-trip write/read of a cited resolution; read refusal of `conn:` under `by: person:captain`; historical FO resolution without `conn:` still reads. Fails if the schema loosens or the read validation is dropped.
- **Grader offline** (`recorded_gate_lifecycle_test.go` mutants + synthetic observations): actor-swap now reds with `conn-approval-misattributed`; missing `conn:` reds; quote absent from runbook reds; the valid cited baseline greens. Each mutant names the single durable-state edit that flips it.
- **Boundary offline** (`auto_continue_negative_test.go`): FO-attributed cited resolution on the no-conn fixture reds `human-gate-bypassed`.
- **Live**: the existing `recorded-gate-lifecycle` claude lane proves a real FO under the updated `fo-gate-lifecycle` text closes with the citation (no new lane; next scheduled runs). Auto-continue lane unchanged.
- **Skill/doc pins**: update `fo_function_reference_invariant_test.go` and skill smoke fixtures where they pin the `gate record` command lines (required by repo policy before changing FO command text).

Cost: offline tests are pure functions, cheap; live proof rides existing lanes.

## Documentation changes

`docs/specs/gate-resolution-frontmatter-contract.md` grammar line, before:
`spacedock gate record ENTITY --decision approve|revise|hold --actor ID [--reason TEXT] [--consume] [--workflow-dir DIR]`
after:
`spacedock gate record ENTITY --decision approve|revise|hold --actor ID [--reason TEXT] [--conn-quote TEXT --conn-source TEXT] [--consume] [--workflow-dir DIR]`

Same file, delegated paragraph, before: "New delegated chat resolutions use `by: agent:first-officer`, require a nonblank evidence reason, and reject `adoption-note` as an unknown prototype field. The recorder constructs the portable Resolution under the asserted identity that rendered the decision; it does not authenticate chat or apply the result." — after: "New delegated chat resolutions use `by: agent:first-officer`, require a nonblank evidence reason, and require a conn citation (`conn: {quote, source}`) naming the grant verbatim and where it was given; they reject `adoption-note` as an unknown prototype field. A citation is refused on `person:captain` resolutions. The recorder constructs the portable Resolution under the asserted identity that rendered the decision; it records the cited grant without authenticating chat, and does not apply the result."

`docs/site/reference/command-reference.md` row 98: append "Delegated First Officer decisions also require `--conn-quote` (the grant verbatim) and `--conn-source` (where it was given); citation flags are refused with `--actor person:captain`."

`docs/site/concepts/gates-and-decisions.md`, before: "Durable gate state records the first officer as the decision renderer and its evidence reason; it does not quote or authenticate the grant's wording or scope. Keep any required chat provenance in the host's own audit system." — after: "Durable gate state records the first officer as the decision renderer, its evidence reason, and a citation of the grant it acted under — the grant's wording and where it was given. The record attributes; it does not authenticate the grant. Captain decisions carry no citation, so the two shapes stay distinguishable."

`skills/fo-gate-lifecycle/SKILL.md`: the FO template line gains `--conn-quote GRANT_VERBATIM --conn-source WHERE_GRANTED`; prose "Recorder retains no grant." becomes "Cite the grant: quote it verbatim from the conversation and name where it was given; the recorder refuses an FO decision without the citation."

`skills/first-officer/references/claude-first-officer-runtime.md:21`: "record an FO-rendered delegated decision as `agent:first-officer`" gains ", citing the grant (`--conn-quote`/`--conn-source`)".

## Stage Report: ideation

- DONE: Settle first whether the auto-continue journey's prompt actually grants the conn — the whole task rests on it, and the FO's authorization is either real or it is not.
  Settled: no. `autoContinuePrompt()` is the neutral no-conn runbook (auto_continue_fixtures_test.go:146-161; negatives pin resolved gates RED). The observed runs were unauthorized (own-sonnet-gate-conn-bypass's evidence); the real under-conn attribution gap is finding 9 on recorded-gate-lifecycle, whose runbook grants the conn verbatim (recorded_gate_lifecycle_test.go:69,776) and whose grader follows the FO's own actor choice (:1150-1175). ACs retargeted accordingly.
- DONE: Design the grant citation as something the binary can check, since the conn is prose in a prompt today and a record that cites nothing checkable is not an improvement.
  Designed: strict `conn: {quote, source}` on the Resolution plus a refusal matrix — FO decisions require it, captain decisions refuse it, reads refuse forged shapes — so the binary checks form and disjointness; the live grader authenticates the quote against the runbook it authored (the check the retired `--directive` never had). Missing citation refuses (exit 2, no mutation), not warns. Grant-registry alternative considered and rejected on necessity.
- DONE: Invert the grader with the product change: a conn-delegated approval attributed to the FO with a citation must grade GREEN, and person:captain with no in-session captain must grade RED.
  Designed into AC-2/test plan: pin expectedActor to agent:first-officer (delete the command-log derivation), require the citation for GREEN, red person:* under distinct code `conn-approval-misattributed`; auto-continue keeps its RED and gains the cited-but-unauthorized boundary negative. Implementation ships the inversion; offline mutants prove both directions with no model spend.

### Summary

Settled the load-bearing factual question against the seeded premise: the auto-continue prompt grants no conn, so the observed runs stay with own-sonnet-gate-conn-bypass, and this task lands on recorded-gate-lifecycle where the conn is real and finding 9 shows misattribution grading green. Designed the smallest disjoint-shape mechanism: a required, strict `conn:` citation on FO-actor resolutions, refused on captain resolutions, with the grader pinned to the FO actor and authenticating the quote against its own runbook. Declared surface net +150 across 11 files (product core ≈ +40/3), with concrete doc diffs and no spike needed on named proven mechanisms.
