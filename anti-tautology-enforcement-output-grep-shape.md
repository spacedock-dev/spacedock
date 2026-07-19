---
id: ekmsd9c7b10br4xhfj09aez6
title: Extend anti-tautology enforcement to the output/prose-grep third shape
status: backlog
source: "This session (2026-07-19): a brittle help-output grep shipped on PR #526 and passed FOUR lenient Claude validation reviews; roborev's CODEX lens caught the analogous c6 AC-4 false-green (job 278) that the Claude reviewer missed. The two filed enforcement mechanisms do NOT cover this shape: testlint-assertion-free-gate catches assertion-free; anti-tautology-enforcement-and-template-gap routes the MIRROR shape to a reactive AC-provenance audit. Neither catches the THIRD shape (real sink + hand-written literal asserting rendered OUTPUT WORDING no consumer parses). Sibling verified sweep: fix-tautological-output-grep-tests (8 confirmed)."
started:
completed:
verdict:
score:
worktree:
---

Add a standing mechanism against the output/prose-grep third shape, and make the review lens that already catches it (roborev codex) a reliable backstop — because prose rules against prose-grep are themselves prose and were applied leniently 4x this session.

## Problem

The third tautology shape — a test with a real failure sink whose hand-written literals assert rendered command/help/doc OUTPUT WORDING that no machine consumer parses — is caught by neither filed gate. It recurs (help_test.go this session; the 8 in the sibling remediation task) because the rule against it lives only in prose (`.roborev.toml` Behavioral-proof / Test-proof-gate, AGENTS.md) applied by whoever reviews. The distinguishing rule: does a machine consumer PARSE the string, and would a real BEHAVIOR change (not a rewording) flip the assertion? Exit codes, JSON field presence/identity, machine-parsed contracts = behavior (legit); human wording with no parser = prose (tautology-risk).

## Proposed approach

{Ideation designs. Candidate directions: (a) a lint/heuristic gate (flag string-contains over --help/usage/stdout output — especially example command lines redundant with an execute-the-example test); (b) make roborev's codex-lens Test-proof-gate a REQUIRED check on every workflow PR — it caught c6 AC-4 (job 278) when the Claude reviewer did not — and fix the env so it runs (codex ~/.codex read; run against the commit from the main checkout, not the worktree); (c) tighten the validation-stage FO scope to hard-reject output/prose-grep and cross-check every test against its AC's stated verification method. Weigh code-gate (durable) vs reactive-review (cheaper) per shape, as anti-tautology-enforcement-and-template-gap did for the mirror shape.}

## Out of scope

The assertion-free gate (testlint-assertion-free-gate) and the mirror-shape reactive audit + commission-template gap (anti-tautology-enforcement-and-template-gap) — this EXTENDS that family to the third shape.

## Coordination note

Overlaps `anti-tautology-enforcement-and-template-gap` (same enforcement family, mirror shape). Captain to decide: fold into that entity, or keep separate with a cross-reference. Also records the durable session conclusion: the shipped-.md-prose-grep class is confined to contractlint (tracked) + pi_frontdoor (in the remediation-task 8); the non-contractlint .md sweep was clean.

## Acceptance criteria

**AC-1** — {ideation defines: a standing check or required-review-lens that fails/flags a newly-introduced output/prose-grep test which a real behavior change cannot break, verified against a fixture that the mechanism flags and a legitimate output-behavior assertion it does not.}
