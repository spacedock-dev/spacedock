---
id: ajg9afxxb1gdksfnh6jd032n
title: status --ac-scan extracts explicit per-AC evidence rows, not only checklist citations
status: backlog
source: "Captain-filed spacedock-dev/spacedock#766, 2026-08-26: a report with a three-item checklist and an ### Acceptance results subsection carrying AC-1..AC-5 PASS rows reported AC-2/3/4 unevidenced — the extractor scans only checklist item ranges and a range sentence (AC-1 through AC-5) cites only its endpoint tokens"
started:
completed:
verdict:
score:
worktree:
issue: "spacedock-dev/spacedock#766"
pr:
mod-block:
---

`--ac-scan` treats completion-checklist item ranges as the only citation source. The dispatch contract caps checklists at three linchpins while real entities carry more ACs, so validators naturally report one row per AC beneath a structured subsection — which the extractor ignores. The result is loss of structured evidence at the extraction boundary, not a missing-evidence finding. This session hit the artifact class at two gates and substituted disclosed manual verification.

## Problem

{Ideation fills this in. Seeded from the issue: define ONE unambiguous per-AC evidence form inside the latest matching Stage Report, extracted in addition to existing checklist citations. Safety boundary preserved: negative mentions ("AC-3 has no evidence") in Summary, Feedback Cycles, or reviewer prose never count as evidence; a range sentence never manufactures citations for interior criteria. Sibling extractor gap observed this session, for ideation to rule in or out with the captain: the heading matcher also misses ACs whose bold spans wrap lines (`**AC-1 (value) — ...` multi-line), which returned an EMPTY scan on a five-AC entity — same boundary, different symptom.}

## Proposed approach

{Ideation fills this in. The issue's acceptance boundary is the spec: explicit rows cite every named AC; ranges do not; negative prose stays excluded; latest-report scoping and --checklist output unchanged; removing one row marks only that criterion unevidenced.}

## Risk evidence

{Backlog: the issue's reproduction plus this session's two disclosed gate artifacts decide design should start.}

## Out of scope

The producer-side body-schema transport (#765, its own task). The checklist cap.

## Expected surface and tolerance

{Backlog seed; ideation estimates with the production/proof split.}

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) - A latest Stage Report with a three-item checklist and explicit AC-N PASS/FAIL rows reports citations for every named AC, and removing one row marks only that criterion unevidenced.**
Verified by: {ideation refines — seed: fixture reports per the issue's boundary cases (range-only, explicit rows, negative prose, per-row removal); failing-today baseline: the issue's reproduction returns AC-2/3/4 unevidenced.}

## Test plan

{Ideation fills this in.}

### Feedback Cycles

{First officer appends one `- Cycle {N}: ...` line per correction round; the validation gate reads reviewer findings from here.}
