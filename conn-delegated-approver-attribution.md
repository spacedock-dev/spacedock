---
id: j7jhntfa2ve8g6jwhatktrrv
title: Under a delegated conn, record the FO as approver and cite the grant
status: validation
source: "Captain CL, 2026-08-18, reframing the live-lane inventory: the auto-continue journey tests the conn and the FO behaved correctly under it; the defect is the approver label, not the approval. Corroborated by the in-tree audit note at internal/ensigncycle/shared_live_runner_test.go:139 — 'finding 9 — approval-actor alternation under a delegated conn: recording person:captain for a decision no captain made in-session grades green.'"
started: 2026-08-18T18:41:26Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-conn-delegated-approver-attribution
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
              resolution:
                type: Resolution
                id: resolution:spacedock:j7jhntfa2ve8g6jwhatktrrv:ideation:1
                briefing: briefing:j7jhntfa2ve8g6jwhatktrrv:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-18T19:32:07.591936Z"
                decision: approve
                reason: 'Captain approved in chat: ''approve those 4, and have them be on a pr stack.'' Accepts the retarget onto recorded-gate-lifecycle and the disjoint-shape citation design, including the grader inversion the captain called for: record the approver as FO under the conn and cite where the conn was given.'
              application:
                target-stage: implementation
                state: consumed
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

### Feedback Cycles

- Cycle 1: FO SCOPE CORRECTION (not a review rejection) — first-officer; surface 14 files/net +261 vs estimate 11 files±2/net +150±40 (174%); AC unchanged. Cause: the FO's dispatch scope-notes assigned `docs/site/reference/command-reference.md` to the sibling `merge-guard-requires-preceding-report` entity and fenced this entity out of three doc diffs its own approved ideation specified, plus the stale finding-9 audit note. The worker correctly stopped and flagged rather than editing outside its stated ownership. The FO granted the missing scope and routed it to the same live worker; cycle 2 applied exactly the four named items (4 files, net +1) with no collateral edits. The same scope-notes error also gave `internal/ensigncycle/recorded_gate_lifecycle_test.go` to both this entity and `decide-dispatch-build-count-bar`, which the FO resolves at restack.

## Stage Report: implementation

- DONE: AC-1: durable state distinguishes a captain decision from an FO decision under a grant — the conn citation is present and its quote appears verbatim in the granting runbook.
  `gates.Conn{Quote,Source}` on `Resolution` (model.go), required on FO chat closes by `recordChatLocked` and the CLI grammar; round-trip + forged-record proof in `TestConnCitationRoundTripsAndReadRefusesAForgedShape` (internal/gates); live-shape proof requires the quote verbatim inside `recordedGateDirective` in `assertRecordedGateLifecycle`. Commit d7fe02a03.
- DONE: AC-2: recorded-gate-lifecycle greens only on an FO-attributed resolution with a valid citation and reds by:person:* under conn-approval-misattributed; auto-continue still reds an FO-attributed cited resolution under human-gate-bypassed.
  Deleted `recordedGateExpectedActor` (finding 9's circular derivation); actor is pinned to `agent:first-officer` and a `by: person:*` reds under the distinct code, asserted in `TestRecordedGateLifecycleProvenanceMutants/actor-swap`. New `TestAutoContinueBypassRedsWithFOAttributedConnCitation` proves the citation confers no authority on the no-conn fixture.
- DONE: AC-3: the two record shapes are disjoint by grammar — every refusal case exits 2 with byte-identical entity and no mutation.
  `gate record` (cli.go) refuses FO-missing-citation, captain-with-citation, and round/briefing-with-citation, each exit 2 pre-mutation; matrix in `TestGateRecordConnCitationGrammarMatrix` (internal/cli) plus two live-binary cases in `TestRecordedGateLifecycleAC5RefusalMatrix`. `gates.Read` refuses a forged captain+conn record (same round-trip test as AC-1).

### Summary

`go test ./...` is green except the pre-existing, machine-local `TestCodexResolveManifestAgainstInstalledHost` (fails identically on `main` before this change — unrelated codex-install-state check); `go test ./internal/gates/... ./internal/cli/... ./internal/ensigncycle/... -race` is clean. Measured surface: 10 files, +373/−113 (net +260) — files are within the declared ±2 tolerance but net LOC is materially over the declared +150±40; the overage is thorough CLI/gates-package test coverage (AC-3's own "exit code, stderr, byte-identical entity" and "read-refusal" requirements), not scope creep. `docs/specs/gate-resolution-frontmatter-contract.md`, `docs/site/reference/command-reference.md`, and `docs/site/concepts/gates-and-decisions.md` still need the prose diffs the design specifies — they sit outside this entity's file ownership (command-reference.md is owned by the sibling `merge-guard-requires-preceding-report` entity in the stack) and need FO routing. The finding-9 audit note at `internal/ensigncycle/shared_live_runner_test.go:139` is now stale but that file is likewise outside my ownership (owned elsewhere in the stack) — flagging for retirement rather than editing it directly.

## Stage Report: implementation (cycle 2)

- DONE: AC-1: durable state distinguishes a captain decision from an FO decision under a grant — the conn citation is present and its quote appears verbatim in the granting runbook.
  No code change this cycle; documented the mechanism per the FO's granted scope (below).
- DONE: AC-2: recorded-gate-lifecycle greens only on an FO-attributed resolution with a valid citation and reds by:person:* under conn-approval-misattributed; auto-continue still reds an FO-attributed cited resolution under human-gate-bypassed.
  No code change this cycle; retired the now-stale finding-9 `AUDIT` note the FO named at `shared_live_runner_test.go:139` — the grader it described (actor-following expected-actor derivation) was deleted in cycle 1. No other line in that file touched.
- DONE: AC-3: the two record shapes are disjoint by grammar — every refusal case exits 2 with byte-identical entity and no mutation.
  No code change this cycle; documented the grammar per the FO's granted scope (below).

### Summary

The FO's cycle-1 file-ownership note was corrected and the missing doc scope granted; this cycle applies exactly the four items named. `docs/specs/gate-resolution-frontmatter-contract.md`: command-surface grammar line gains `[--conn-quote TEXT --conn-source TEXT]`; the delegated-resolution paragraph gains the citation requirement and its captain-side refusal. `docs/site/reference/command-reference.md` row 98 (gate record only — the sibling `vka` entity's `merge guard` row untouched): appended the citation-flag sentence verbatim from the ideation's "Documentation changes" section. `docs/site/concepts/gates-and-decisions.md`: replaced the "does not quote or authenticate" paragraph with the attributes-not-authenticates citation paragraph, verbatim from ideation. `internal/ensigncycle/shared_live_runner_test.go`: deleted the 3-line finding-9 `AUDIT` comment (lines 139-141 pre-edit) immediately above the `recorded-gate-lifecycle` live-journey registration; no other line in the file changed. Commit 15bb2d8ba. Updated surface total across both cycles: 14 files, +384/−123 (net +261) — cycle 2 alone is 4 files, +11/−10 (net +1), exactly the granted scope with no collateral edits. `go test ./...` green except the same pre-existing, machine-local `TestCodexResolveManifestAgainstInstalledHost`; `internal/contractlint` (byte caps, doc-reference, and live-registry structural checks) green.

## Review-finding disposition

- **Reviewer (validation, 2026-08-18) — Deferred risk: the grader's conn-SOURCE check is not scoped to the `conn:` block, so any `source:` line in the frontmatter satisfies it.** Released user and normal workflow: the `recorded-gate-lifecycle` live journey. Observable harm: `recordedGateConnCitation` (`recorded_gate_lifecycle_test.go:237`) regex-scans the whole authority block for the first `quote:`/`source:` line; an audit probe that stripped the conn block's own `source:` and added an ordinary top-level `source:` entity-provenance field graded the journey GREEN on a sourceless citation. Affected authority: `value-ac[AC-2]` — not violated. AC-2's substantive claims all hold and were reproduced (GREEN only on FO+valid citation; `by: person:*` RED under `conn-approval-misattributed`; auto-continue still RED), and the QUOTE half — the only half carrying authentication weight, since `source` is unauthenticated free text by design — is tightly bound by content, not position. Trigger evidence: none reachable on a shipped path; no writer can emit a sourceless conn block (CLI exit 2, `recordChatLocked` refuses, `gates.Read` refuses blank), and the harness-authored fixture `recordedGateEntity()` carries no top-level `source:`. Promotes to material if the recorded-gate fixture entity gains a top-level `source:` or `quote:` frontmatter field — which real Spacedock entities do carry, so a fixture refresh toward realism triggers it; the narrow fix is to slice the `conn:` block out of the authority before extracting (~3 lines, test-only).
- **Reviewer (validation, 2026-08-18) — Polish: two `internal/gates` defense-in-depth guards have zero test coverage.** Deleting `recordChatLocked`'s conn matrix leaves `./internal/gates` and `./internal/cli` fully green; the same holds for `RecordSemanticSummary`'s `--round`+conn incompatibility. No current user-visible loss: `internal/cli/cli.go` is the only non-test caller of `RecordSemantic*` and its guard refuses first at exit 2. But removing the CLI guard alone drops the refusal to exit 1 via the library guard, so the library guard is the real net — and no test would notice its deletion. Narrow fix: two `RecordSemantic` table cases in `internal/gates/gates_test.go`.
- **Reviewer (validation, 2026-08-18) — Polish: two command-surface texts still omit the now-required citation flags.** `docs/site/reference/command-reference.md`'s left-hand grammar cell still reads `--actor ID [--reason TEXT] [--consume]` while the spec grammar line and `gate --help` both gained `[--conn-quote TEXT --conn-source TEXT]`; `skills/present-gate/SKILL.md:13` shows the same unflagged form. Both match the approved design (the ideation specified exactly "row 98: append <sentence>" and did not enumerate `present-gate`), the command-reference row's own description states the requirement, and the authoritative FO template at `skills/fo-gate-lifecycle/SKILL.md:51` carries the flags. An FO following either text alone gets a clear exit-2 refusal, not silent corruption.

## Stage Report: validation

- DONE: AC-1: verify durable state distinguishes a captain decision from an FO decision under a grant — the conn citation round-trips, its quote must appear verbatim in the granting runbook, and a forged captain-plus-citation record is refused on read. Reproduce the evidence rather than trusting the report.
  Reproduced end-to-end, not synthetically: in a throwaway checkout I drove `TestRecordedGateLifecycleRealCLIReplay`'s REAL binary with an invented quote that still carries the granted phrase, and the grader red with "does not appear verbatim in the granting runbook". Deleting that one check greens the offline `quote-not-in-runbook` mutant; deleting the read-side `By != agent:first-officer` check fails both `TestPortableResolutionValidation/captain_with_conn_citation` and `TestConnCitationRoundTripsAndReadRefusesAForgedShape`. Historical no-conn resolutions still read.
- DONE: AC-2: verify the grader inversion actually inverted — recorded-gate-lifecycle greens ONLY on an FO-attributed resolution carrying a valid citation, a by:person:* resolution reds under the distinct code conn-approval-misattributed, and the auto-continue journey still reds an FO-attributed cited resolution under human-gate-bypassed. Confirm recordedGateExpectedActor is genuinely deleted, not merely bypassed.
  `recordedGateExpectedActor` is gone tree-wide (grep over all `*.go`), along with its own test; only a local `const expectedActor = "agent:first-officer"` pin remains at `:127`. Deleting the graded misattribution branch makes `actor-swap` grade under `""` instead of `conn-approval-misattributed`. The auto-continue boundary is non-vacuous: an audit probe put the injected fixture through `gates.Read` and it parses as a genuinely valid FO-attributed cited resolution (`Conn` round-trips) — and still reds under `human-gate-bypassed`.
- DONE: AC-3: verify the two record shapes are disjoint by grammar — FO-actor missing either citation flag, captain-actor with citation flags, and citation flags with --round/--briefing each exit 2 with a byte-identical entity and no mutation. Assert the exit code and the bytes, not just the error text.
  `TestGateRecordConnCitationGrammarMatrix` asserts `exit == 2` exactly, empty stdout, stderr text, `bytes.Equal(before, after)`, and no lock residue across 7 chat cases plus the `--round` case; the ensigncycle AC5 matrix adds live-binary cases under whole-tree `treeDigest` byte-identity. Deleting the CLI guard drops the exit code to 1 (the library guard still refuses) and the matrix catches it, so the exit-code assertion is load-bearing. The grammar is total: the actor whitelist is exactly `{person:captain, agent:first-officer}` (`operation.go:219`), so no third actor can silently drop the flags.

### Summary

The replacement is not circular in a new way. The grader authenticates the FO's quote against `recordedGateDirective`, which `recordedGatePrompt()` (`:858`) interpolates verbatim into the runbook the live FO receives — so the expectation comes from harness-authored bytes, never from the FO's own command log. An adjacent-variant matrix over "the quote must be authentic" holds across seven cases: the full grant and authentic fragments green, while an invented quote echoing the granted phrase, an invented quote without it, an authentic substring lacking the phrase, and — importantly — the real grant with appended self-granted authority all red. The citation confers no authority: a valid, parseable, FO-attributed cited record still reds on the no-conn fixture. `go test ./... -race` is green except `TestCodexResolveManifestAgainstInstalledHost`, which I reproduced failing byte-identically at the merge-base `a108559c9` — pre-existing, machine-local, unrelated; `contractlint`, `status`, and `gofmt -l ./cmd ./internal` are clean. Surface is 14 files, +384/−123 (net +261) against 11±2 files / +150±40 — past both tolerances (1 file over, +71 net over, 174%). I attributed it independently: product core is +78/−6 (net +72, the seeded "+70 across 3 files" shape), docs and skills are net +4 (UNDER declared), and the entire +88 overage is grader/test code — concentrated in `gate_test.go` (+91/−6, the CLI grammar matrix) and `gates_test.go` (+52, round-trip/forged/historical), which is exactly the proof AC-3 itself demands. No file landed outside the approved design: 12 of 14 are named in "Expected surface and tolerance", `shared_live_runner_test.go` is the audit-note retirement named at mechanism 5, and `gate_consume_sync_test.go` is a one-line forced call-site update. Cycle 2 independently checks out as 4 files / net +1 / exactly the granted scope. No AC was narrowed and no material finding remains. Recommendation: PASSED, with three deferred/polish findings recorded above and the surface overage as a captain-visible fact.
