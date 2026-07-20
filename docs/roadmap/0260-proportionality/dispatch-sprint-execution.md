# 0260 Commander package — cold-boot execution dispatch

You are the Commander: a fresh FO session driving 0260-proportionality from approved ideations to a cut 0.26.0. Shaping is complete; every driving member's ideation gate is closed-approved with a pending advance. Your job is implementation → validation → done per member, merge to main, and the pre-cut audit.

## Boot order

1. Load the FO contract (`spacedock:first-officer` skill) and engage the `docs/dev` workflow.
2. Read `docs/roadmap/0260-proportionality/index.md` IN FULL — its `## Constraints` (operating directives + captain rulings) bind this session's own conduct, not just the workers'.
3. Drivable set: `spacedock status --workflow-dir docs/dev --where sprint=0260-proportionality --where 'sprint-readiness != defer'`.
4. Session context if needed: `docs/dev/.spacedock-state/_debriefs/2026-07-20-01-0260-shaping.md` (gate history, float wiring findings).

## Consuming the recorded approvals (3k notation)

Each member's approval lives in its entity `gates:` frontmatter — see the index `## Gate approvals` section for the rules. In short: a closed attempt with `resolution.decision: approve` and `application: {action: advance, target-stage: implementation, state: pending}` IS the captain's approval; apply it once via the normal transition/dispatch path and mark `state: consumed`; on entity-content drift from `briefing.digest`, mark `superseded` and re-present instead of advancing. Do not re-ask the captain for any closed gate.

## Members and landing notes

| ref | slug | deliverable surface | notes |
|---|---|---|---|
| 85 | merge-guard-arm-not-a-stopping-point | FO contract prose | independent, small |
| ht | fix-tautological-output-grep-tests | 8 test fixes | independent; estimate discipline caught 6 files, not 5 |
| bw | feedback-cycle-record-command | `feedback-rejection-flow` skill prose | prose-only convention; the command is DEFERRED — do not build it |
| 02av | ensign-finding-triage-disposition | `feedback-rejection-flow` standing block + `docs/dev/README.md` implementation bullet | co-edits bw's file — land bw and 02av together at a shared gate; `ensign-shared-core` gets NOTHING |
| z7 | falsifiability-ladder | contract check-ordering prose | budget re-baselined +1055 boot-resident; lure catalog is validation-time material, never a committed suite |
| az | anti-tautology-enforcement-and-template-gap | testlint revert-fix + evidence rule | Edit D (audit-trigger widening) is NOT approved — requires its own explicit captain yes/no before any Edit D work |
| 841 | contractlint-codex-runtime-semantics-retirement | 121 prose assertions → 0, 3 Go-source bindings | remaining four retirements carry `sprint-readiness: defer` |
| 2ae | template-rigor-propagation | commission templates + refit skill (4 files, ~18 net lines) | verbatim carriage from siblings; see binding condition below |

Suggested order: (85, ht, 841) in parallel first — independent and mechanical; then z7 and az; then bw+02av as the co-edit pair; 2ae LAST so its verbatim sources are landed contract text.

## Binding captain conditions (recorded in gates frontmatter; verbatim)

- **2ae validation gate:** the captain's approval reason reads "on validation gate, present the refitted delta on the workflow readme for human review" — the validation gate presentation MUST include the refit diff against the workflow README as a human-reviewable artifact.
- **az Edit D:** pending its own captain yes/no; not covered by the az approval.
- **bw's `### Feedback Cycles` entry FORMAT in the template:** deferred by the re-lock; 2ae references the record by name only.
- **High-stakes surfaces** (`feedback-rejection-flow`, commission templates, refit skill, testlint/contractlint): detached adversarial audit at validation before merge, per the Proof policy.

## Definition of Done

The index `## Definition of Done` is the authority — replay fixtures live in `docs/dev/.spacedock-state/_evidence/0260-agent-derail-forensics/`. Every DoD line is a check that can fail; "the prose updated" is never evidence.

## Close-out

Pre-cut antipattern audit (independent staff-eng reviewer over the assembled sprint) → `go test ./...` green → `docs/releasing.md` for the 0.26.0 cut (captain authorizes the cut). Seed-next-sprint is the Shaping FO's, not yours.
