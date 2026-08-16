---
title: Gated initial stage deadlocks — no report means not gate-ready, gated means not dispatchable
status: ideation
source: "Field report from a pre5/0.27 consumer FO, relayed by the captain 2026-08-16 with focused-test confirmation; captain order: file and dispatch"
id: ra9qzfz94hzgsq938jz998mj
---

## Problem

A workflow whose INITIAL stage is gated deadlocks. The layers disagree (consumer citations from the 0.27.0 plugin cache; repo paths here, ideation verifies exact lines):

1. Contract: skills/fo-gate-lifecycle/SKILL.md:28 requires an exact-stage report before preparation and mandates report-incomplete otherwise; SKILL.md:24 requires --checklist and --ac-scan, both of which fail without a Stage Report.
2. Scheduler: internal/gates/model.go:155 labels a gated stage with no gate record as `validating`; internal/status/discover.go:261 promotes `validating` to needs-preparation only when hasCompleteCommittedStageReport succeeds; internal/status/entered_stage.go:28 defines that as a committed exact-stage report with checklist evidence and a summary; internal/status/format.go:150 excludes `validating` from ready_gates and format.go:296 excludes every current gated stage from ordinary dispatch.
3. Result: no report → not gate-ready; at a gate → not dispatchable. The entity sits forever, labeled `validating` — a lying label for a stage that has never validated anything.

The decisive intent proof: internal/gates/prepare.go:56 does NOT require a stage report — only a valid actionable gated stage and committed artifacts — and its own fixture (internal/gates/prepare_test.go:1002) explicitly prepares an initial gated entity with no report. The low-level gate machinery was designed for this; the scheduler and contract layers above it contradict the machinery's own test. Focused tests confirmed all cited behaviors.

The only workaround is a fabricated "backlog completion report" restating the committed seed — fake work of exactly the class the workflow-fit-gate entity names, except here the contract demands the ceremony rather than the FO inventing it.

## Proposed approach

The consumer FO's fix direction, with two scope disciplines bound:

1. Treat the committed initial-stage seed as the gate artifact: the scheduler promotes a gated INITIAL stage to needs-preparation when the entity is committed and clean in HEAD, without a completion report. The exception keys on the stage being initial (`initial: true`) ONLY — never on "no report exists" — so no other stage gains a path around its completion proof.
2. The clean-in-HEAD guard and the briefing's digest pinning stay: the captain reviews exactly the committed seed bytes.
3. Replace report-based checklist/AC extraction for the initial gate by scanning the seed's own acceptance-criteria section (field case: the consumer's 8e seed carries five ACs ready to scan).
4. Give the gated-no-record initial state an honest label or an honest needs-preparation path — `validating` is untrue there.
5. Amend fo-gate-lifecycle SKILL text with the matching initial-stage exception: prepare directly from the seed; present outcome, scope, exclusions, and required ideation proof.

Falsifiability both directions: a deadlock fixture (gated initial stage, committed seed) is red today — no scheduler row — and shows needs-preparation → prepare → present after the fix; a NON-initial gated stage without a report must still refuse, proving the guard survives.

## Out of scope

- The workflow-fit-gate write-core amendment (b8 — complementary entity, cross-referenced).
- The consumer workflow's own state (8e re-prepare is their field validation, referenced only).
