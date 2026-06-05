---
id: 4qnn7dbzkyh9qv65t618vtxy
title: AC-3 sweep reader-axis — invert the in-package read predicate + bound the cross-package/package-var gaps
status: backlog
source: "first-officer (2026-06-05) — hwk (tautological-test-remediation) Option-C scope-fork (captain-decided). hwk's cycle-3 positive/taint redesign CLOSED the match axis (ingest⇒declare, M1/M2/M3 confirmed closed). A detached adversarial audit then found the READER axis reintroduced the same allow-list failure mode on three new shapes. Captain decision: ship hwk's match-axis closure + the 56 remediated tests now with the guard's guarantee honestly scoped; fork the reader-axis hardening here."
score: "0.22"
started:
completed:
verdict:
worktree:
issue:
---

The two AC-3 tautology sweeps (`TestNoUndeclaredTautologicalProof` in `skills/integration/`, `TestNoUndeclaredHostneutralityTautology` in `internal/hostneutrality/`) detect an undeclared instruction-file presence-check only when its READ is discovered. hwk closed the match axis (any inspection idiom over ingested bytes now triggers "must declare") but the reader-discovery side still rests on recognizers a genuine read can sit outside of. A detached adversarial audit (hwk cycle-3, HEAD `ea0441d6`) proved three evasions stay GREEN.

## The audit findings (concrete, reproducible)

- **M-A — `isInstructionPathLiteral` segment allow-list misses real instruction surfaces.** A path whose literal carries none of the recognized segments (`skills`/`references`/`agents`/`first-officer`/`ensign`/`commission`/`present-gate`/`SKILL.md`) evades. Confirmed GREEN: an undeclared `os.ReadFile`+match over `AGENTS.md` (repo root; `agents` ≠ case-sensitive `AGENTS`) and over `mods/pr-merge.md` (FO-facing `## Hook: startup`). Positive control over `skills/refit/SKILL.md` REDs, proving the sweep is alive and the gap is purely the segment list.
- **M-B — cross-package reads are invisible.** The sweep AST-scans only its own package's `*_test.go`. A read via another package's helper (e.g. `internal/dispatch.ParseModMetadata` over `skills/present-gate/SKILL.md`) stays GREEN even with a fully recognized instruction path — the read sink lives in another package.
- **M-C — package-level `var` path escapes per-function taint.** A recognized `.md` literal in a package-level `var` initializer (not the test function body), read via `os.ReadFile(thatVar)`, evades — `instructionTaintedNames` seeds only from the test fn body; HN's `instructionPathIdents` is itself a hand-listed allow-list.
- **M-D — `[]string`/`...string`-param and `range`-loop-variable flow escapes the taint** (hwk cycle-3 validation, detached audit at `ea0441d6`). The taint seeds only on bare-`string` params (`isStringyType`) and propagates only through `*ast.AssignStmt`; a `.md` literal sitting in a `[]string` composite read element-by-element in a `range` loop never reaches the read sink. Confirmed GREEN on BOTH sweeps: `paths := []string{"…/first-officer-shared-core.md"}; for _, p := range paths { b,_ := os.ReadFile(p); strings.Contains(string(b), "present-gate") }`, no marker. A `[]string`-param reader helper evades identically. (Closure-capture IS covered; the gap is the slice/range edge.)

These four are instances of ONE class: the reader-flow taint is an *enumerated set of AST flow shapes* (cycles 1–3 each closed the instances found, not the class). The recurrence is the signal that enumeration is the wrong tool here.

## Direction (for ideation)

The auditor's thesis (and the proof-policy's own lesson, already learned on the match axis): the segment/ident lists and same-package scan cannot be made exhaustive by enumeration — extending the lists just moves the hole. The honest fixes are scope-defining:

- **Invert the in-package read predicate** — treat ANY `os.ReadFile`/`Open`/`io.ReadAll`/`bufio` over ANY `.md`/instruction path as a read-to-declare, requiring a positive *non-instruction* exemption (a finite, reviewable exemption list — `.json` manifests, `docs/dev` recipes), instead of a positive instruction-recognizer. This structurally closes M-A and the whole unrecognized-path class.
- **M-B (cross-package) and M-C (package-var-from-another-file)**: decide per ideation — either extend discovery (cross-package reader resolution; package-var taint), OR formally BOUND the guard's guarantee to "in-package reads via recognized flows" and DOCUMENT M-B/M-C as out-of-scope structural limits of a per-package AST scan, with the detached adversarial audit as the named backstop for those classes. A per-package static meta-test structurally cannot see another package's read; claiming it does would be the false-universal the proof-policy bans.
- **Consider a go/types + SSA-backed taint instead of AST-shape enumeration** (raised by the hwk cycle-3 validator). Real type information + SSA value-flow would close the reader-flow class *definitionally* (M-A unrecognized-segment, M-D slice/range, and the param/field/method/var shapes all become "a value derived from an instruction-file path reaches a read sink," regardless of syntactic shape) rather than chasing one more AST shape per cycle. Weigh the cost (a `golang.org/x/tools/go/ssa` pass over the test packages) against the bounded-but-recurring AST approach on this high-stakes oracle. This is the candidate that actually stops the whack-a-mole; ideation should cost it.

## Out of scope

The match-axis closure, the 56 remediated tests, the four cycle-1 reader shapes, and the cycle-2/3 mutation controls — all verified sound and shipped with hwk. This task is solely the reader-axis robustness + an honest guarantee statement.

## Acceptance criteria

Each AC names a finished-state property + how an outside-the-body check verifies it (proof-policy: a planted-control test that REDs on the evasion shape then GREENs once caught — never a prose assertion).

**AC-1 — M-A closed: an undeclared read+match over any `.md` instruction surface the prior segment list missed (AGENTS.md, mods/*.md, and a path with no recognized segment literal) REDs both sweeps.**
Verified by: a planted-control case driving the AGENTS.md and mods/*.md shapes RED-then-GREEN in BOTH packages.

**AC-2 — M-B, M-C, and M-D are either closed or honestly bounded.**
Verified by: if closed — planted controls for a cross-package read, a package-var path, and a `[]string`/range flow REDing both sweeps; if bounded — a documented scope statement on each sweep naming the out-of-scope class + the audit backstop, and a test/assertion that the guard's doc no longer claims the falsified universal.

## Test plan

Reuse hwk's planted-control harness (real undeclared presence-check over a genuine instruction file, no marker, run the production sweep, observe RED/GREEN, restore). Offline Go, minutes. High-stakes shipped-test oracle → a detached adversarial audit on the result before merge.
