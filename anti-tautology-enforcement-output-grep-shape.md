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
sprint: 0260-proportionality
group: test-cleanups
---

Add a standing mechanism against the output/prose-grep third shape, and make the review lens that already catches it (roborev codex) a reliable backstop — because prose rules against prose-grep are themselves prose and were applied leniently 4x this session.

## Problem

The third tautology shape — a test with a real failure sink whose hand-written literals assert rendered command/help/doc OUTPUT WORDING that no machine consumer parses — is caught by neither filed gate. It recurs (help_test.go this session; the 8 in the sibling remediation task) because the rule against it lives only in prose (`.roborev.toml` Behavioral-proof / Test-proof-gate, AGENTS.md) applied by whoever reviews. The distinguishing rule: does a machine consumer PARSE the string, and would a real BEHAVIOR change (not a rewording) flip the assertion? Exit codes, JSON field presence/identity, machine-parsed contracts = behavior (legit); human wording with no parser = prose (tautology-risk).

**We already have the mechanism and did not use it.** `internal/contractlint` is an AST-based structural-check framework that rides the existing `go test ./...` sweep, and `testlint-assertion-free-gate` already proposes mirroring it for the assertion-free shape. So the capability to MECHANICALLY gate this shape exists in the repo today. It was not applied. Instead, every instance this session was kept alive by a rationalized carve-out — "it's runtime-loaded," "it's a contract-shape check," "it honestly asserts only the shape" — used to defend the 11-phrase compaction-presence test and the `dispatch build --help` grep. Those carve-outs ARE the anti-pattern in disguise: a test that asserts a string literal it read out of a repo file is checking file CONTENT, not behavior, regardless of the file's role or how the assertion is framed.

## Current tension

Established this session (four lenient reviews missed this class; roborev's codex lens + the captain caught it):

- **Prose-driven discipline does not hold.** The rules against this already exist in prose (`.roborev.toml` Behavioral-proof / Test-proof-gate, `AGENTS.md`, the Proof policy) and were applied leniently 4× in one session. Adding another prose rule — including one that says these tests are "shameful" — just repeats the failure.
- **Only mechanical enforcement is not-prose-driven, and it is not free:**
  - **Structural AST gate (cheap, partial) — reuse the existing contractlint/testlint infra.** Flag the tell of this shape: an assertion whose expected value is a string literal that appears verbatim in a repo file the test itself read (the read-a-file-then-`strings.Contains`-a-literal-from-it form — asserting a file contains its own bytes). Catches the 11-phrase presence test, the `runtime-support.md` grep, any read-and-grep. Blind spot: the execute-a-static-passthrough-then-grep-the-constant form (the help-grep runs `--help` but greps a frozen const), and semantically-tautological tests that look structurally fine.
  - **Scoped mutation testing (general, costs compute).** Mutate the PR's changed production code; a test that kills no mutant is tautological by machine verdict — no author claim, no reviewer judgment. Catches what structure cannot see. Cost: compute + equivalent-mutant noise, so scope to changed files.
- **The prose-driven layers decay.** Author-states-a-bite ritual, reviewer-reproduces-the-bite, and a "named anti-pattern / shame ledger" are backup, not enforcement — a lazy or adversarial author fakes/omits the bite, a lenient reviewer skips reproducing it (exactly what happened 4×).
- **No carve-outs — design mandate.** "Runtime-loaded contract," "contract-shape check," and "honestly scoped" do NOT exempt a prose-grep. A contract-shape AC is verified by a REVIEW-TIME grep cited as evidence — never a committed grep-test. The enforcement must reuse the existing contractlint AST infra rather than rationalizing a new exception, which is the exact failure this task exists to end.

## Proposed approach

{Ideation designs. Candidate directions: (a) a lint/heuristic gate (flag string-contains over --help/usage/stdout output — especially example command lines redundant with an execute-the-example test); (b) make roborev's codex-lens Test-proof-gate a REQUIRED check on every workflow PR — it caught c6 AC-4 (job 278) when the Claude reviewer did not — and fix the env so it runs (codex ~/.codex read; run against the commit from the main checkout, not the worktree); (c) tighten the validation-stage FO scope to hard-reject output/prose-grep and cross-check every test against its AC's stated verification method. Weigh code-gate (durable) vs reactive-review (cheaper) per shape, as anti-tautology-enforcement-and-template-gap did for the mirror shape.}

## Out of scope

The assertion-free gate (testlint-assertion-free-gate) and the mirror-shape reactive audit + commission-template gap (anti-tautology-enforcement-and-template-gap) — this EXTENDS that family to the third shape.

## Coordination note

Overlaps `anti-tautology-enforcement-and-template-gap` (same enforcement family, mirror shape). Captain to decide: fold into that entity, or keep separate with a cross-reference. Also records the durable session conclusion: the shipped-.md-prose-grep class is confined to contractlint (tracked) + pi_frontdoor (in the remediation-task 8); the non-contractlint .md sweep was clean.

## Acceptance criteria

**AC-1** — {ideation defines: a standing check or required-review-lens that fails/flags a newly-introduced output/prose-grep test which a real behavior change cannot break, verified against a fixture that the mechanism flags and a legitimate output-behavior assertion it does not.}
