# 0260 Commander package — cold-boot execution dispatch

You are the Commander: a fresh FO session driving 0260-proportionality from approved ideations to a cut 0.26.0. Shaping is complete; every driving member's ideation gate is closed-approved with a pending advance. Your job is implementation → validation → done per member, merge to main, and the pre-cut audit.

## Boot order

1. Load the FO contract (`spacedock:first-officer` skill) and engage the `docs/dev` workflow.
2. Read `docs/roadmap/0260-proportionality/index.md` IN FULL — its `## Constraints` (operating directives + captain rulings) bind this session's own conduct, not just the workers'.
3. Read the shaping debrief `docs/dev/.spacedock-state/_debriefs/2026-07-20-01-0260-shaping.md` — REQUIRED, not optional: it carries the gate-presentation float ritual (findings 12-13) that the 2ae validation condition depends on, plus the session's conduct findings.
4. Read the staff review pair: `staff-review.md` (codex seat) and `staff-review-fable-delta.md` (fable seat) in this folder. The folds they mandated are applied; their recorded declines stand with promote-to-material conditions.
5. Drivable set: `spacedock status --workflow-dir docs/dev --where sprint=0260-proportionality --where 'sprint-readiness != defer'` (returns eight members).

## Consuming the recorded approvals (3k notation)

Each member's approval lives in its entity `gates:` frontmatter — see the index `## Gate approvals` section for the rules. In short: a closed attempt with `resolution.decision: approve` and `application: {action: advance, target-stage: implementation, state: pending}` IS the captain's approval; apply it once via the normal transition/dispatch path and mark `state: consumed`; do not re-ask the captain for any closed gate.

**Drift check and its waiver.** Before applying a pending advance, compare the entity content against the attempt's `briefing.digest`. The digest covers the exact bytes of the briefing artifact named in the record. For float gates that is the frozen gate-summary in the review room — byte-verifiable. For chat-presented gates the digest is ADVISORY: it hashes the working file at recording time, which no single committed tree reproduces (an entity cannot self-bind its own gates record), so drift checking for a chat attempt diffs the entity BODY against the state commit that introduced the attempt, and never re-hashes the current file. Two kinds of post-approval change are NOT drift and do not void the approval: (1) edits the approval itself directed (recorded in the resolution reason or attempt note); (2) captain-approved staff-review folds recorded as superseding attempts in the same gates record. Anything else IS drift: mark the application `superseded`, open a new attempt, and present it — the no-re-ask rule applies to closed gates, not to drifted content, so presenting genuine drift is required, not forbidden.

## Members and landing notes

| ref | slug | deliverable surface | notes |
|---|---|---|---|
| 85 | merge-guard-arm-not-a-stopping-point | FO contract prose (`fo-merge-core.md`) | independent, small |
| ht | fix-tautological-output-grep-tests | 8 test fixes across 6 Go test files | independent; estimate discipline caught 6 files, not 5 |
| bw | feedback-cycle-record-command | `feedback-rejection-flow` skill prose + 2 shared-core lines + the skill's frontmatter description line | prose-only convention; the command, refusal machinery, and escalation backstop are DEFERRED — do not build them. Deviation beyond tolerance requires a RECORDED reconfirm/re-scope/park/escalate decision before further repair dispatch (captain-approved strengthen, 2026-07-20) |
| 02av | ensign-finding-triage-disposition | `feedback-rejection-flow` standing block + `docs/dev/README.md` implementation bullet | co-edits bw's file — see the shared-edit table; `ensign-shared-core` gets NOTHING |
| z7 | falsifiability-ladder | contract check-ordering prose + `first-officer-shared-core.md` + the full `docs/dev/README.md:74` sentence rewrite + the Claude-runtime pre-workflow declaration line | budget re-baselined +1055 boot-resident; lure catalog (six scenarios) is validation/pre-cut material, never a committed suite |
| az | anti-tautology-enforcement-and-template-gap | testlint-revert report protocol + evidence rule | Edit D resolution is recorded in az's gate record — read it there; no new committed enforcement regardless |
| 841 | contractlint-codex-runtime-semantics-retirement | 121 prose assertions → 0, 3 Go-source bindings | remaining four retirements are next-train backlog (no sprint label) |
| 2ae | template-rigor-propagation | commission templates + refit skill (4 files, ~18 net lines) | verbatim carriage from LANDED sibling text — see wave 4 |

## Landing waves and shared-edit seams

Execute in waves; within a wave, members run in parallel.

1. **Wave 1 — 85, ht, 841 in parallel.** Disjoint surfaces: `fo-merge-core.md`, six Go test files, contractlint tests.
2. **Wave 2 — z7 then az, serially** (or rebase the second immediately before its validation). Both touch `docs/dev/README.md` in the proof-policy region; z7 also edits `first-officer-shared-core.md`.
3. **Wave 3 — bw lands ALONE (captain decision 2026-07-21; supersedes the composed landing).** bw ships its own surface only: the entry-format convention with surface/estimate/AC-drift fields, steps 2-3, the description line, the shared-core lines, the README ideation + in-stage bullets. **02av is HELD — do not land its standing block, README implementation bullet, or findings field.** At 02av's next boundary, route it back to ideation for a captain-directed reframe: recording round dispositions (including the decline) as ADVISORY resolutions under the gate-recorder model (see the 3k scope cut of 2026-07-21 and 02av's `### Feedback Cycles` entry recording this reset). The finding-triage rule text remains approved-shaped; only its record/delivery mechanism is being reframed.
4. **Wave 4 — 2ae last.** Its "verbatim" sources are the LANDED post-wave-2/3 text of `docs/dev/README.md` and the contract files, not the ideation-time copies quoted in its body: re-anchor its Pieces 1-3 before/after diffs against the landed wording, then validate the assembled README before copying from it.

| Shared surface | Editors/readers | Rule |
|---|---|---|
| `first-officer-shared-core.md` | z7, then bw | different sections, same file — rebase/re-anchor bw after z7; never merge stale whole-file prose |
| `docs/dev/README.md` | z7, az, bw, 02av; read by 2ae | validate the assembled README between waves; 2ae copies only landed text |
| `feedback-rejection-flow/SKILL.md` | bw + 02av | one composed landing, one shared gate (wave 3) |
| commission/refit templates | 2ae | semantic dependency on landed sibling wording, not ideation-time quotes |

**Live-lane heads-up (wave 3).** Diffs touching `skills/**/references/**` require the host live lanes green (`docs/dev/README.md:78`). The shipped `feedback-3-cycle-escalation` scenario asserts the RETIRED record shape (bare `- Cycle N:` entries, escalation at third) under both runtimes. If that lane reds against bw's new entry format, it is a CAPTAIN reconciliation decision — the lane's fixture asserts a format the captain retired — never a silent fixture edit, and Go-test changes are outside every member's declared surface.

## Binding captain conditions (recorded in gates frontmatter; verbatim)

- **2ae validation gate:** the captain's approval reason reads "on validation gate, present the refitted delta on the workflow readme for human review" — the validation gate presentation MUST include the refit diff against the workflow README as a human-reviewable artifact, presented via the float ritual in the shaping debrief.
- **bw strengthen (captain-approved 2026-07-20):** the recorded-decision-before-re-dispatch rule above; prose-only.
- **Edit D:** resolved in az's gate record — apply as recorded there.
- **bw's `### Feedback Cycles` entry FORMAT in the template:** deferred by the re-lock; 2ae references the record by name only.
- **High-stakes surfaces** (`feedback-rejection-flow`, `first-officer-shared-core`, commission templates, refit skill, testlint/contractlint): detached adversarial audit at validation before merge, per the Proof policy.

## Definition of Done

The index `## Definition of Done` is the authority — replay fixtures live in `docs/dev/.spacedock-state/_evidence/0260-agent-derail-forensics/`. Every DoD line is a check that can fail; "the prose updated" is never evidence.

## Close-out

1. **Pre-cut audit** — independent staff-eng reviewer over the assembled sprint, AND the full lure-scenario drive: all six approved lure scenarios from z7's catalog, under BOTH runtimes (Claude and codex/gpt-5.6-sol), outcomes recorded in the pre-cut report. The catalog's home is that report artifact — no committed suite, no CI lane.
2. **Repository completion gates** — `go test ./...` AND `go test ./... -race` green; `gofmt -w ./cmd ./internal` clean; final clean `git status`.
3. **Cut 0.26.0** per `docs/releasing.md` *(captain authorizes the tag)*.
4. Seed-next-sprint is the Shaping FO's, not yours.
