---
id: 4qnn7dbzkyh9qv65t618vtxy
title: AC-3 sweep reader-axis — NET-REMOVE the detection machinery; rely on the detached audit backstop
status: ideation
source: "captain (2026-06-05) — re-scoped from 'invert/harden the reader-axis sweep' to NET REMOVAL. Captain directive: 'i want to see net removal, not more crap to mark crap or detect crap.' The hwk merge added ~1525 lines of go/ast sweep machinery (the bulk = the reader-axis taint/discovery), which is BOTH the heaviest part AND the incomplete part (M-A/B/C/D evade it). The detached adversarial audit caught every reader-axis hole the static sweep missed — so the audit, not a static guard, is the right backstop for that axis."
score: "0.40"
started: 2026-06-05T19:10:17Z
completed:
verdict:
worktree:
issue:
---

hwk shipped a standing AC-3 sweep whose MATCH axis ("a test that reads instruction-file bytes and inspects them must declare markNonAC/markCodeBoundInvariant, regardless of idiom") is small, universal, and sound — keep it. But its READER axis (the `readsInstructionContent` taint analysis: param/struct-field/method/closure flow + path-construction discovery + the transitive reader fixpoint) is the bulk of the ~1525 added lines AND is known-incomplete: a detached adversarial audit found four evasion classes that sail through (M-A unrecognized surfaces like AGENTS.md/mods, M-B cross-package reads, M-C package-var paths, M-D []string/range flow). Each hardening cycle bolted on more detection code and the audit still found the next hole — the classic enumeration trap the proof-policy itself warns against.

**The decision (captain): stop detecting. NET-REMOVE the reader-axis machinery.** A per-package go/ast static scan structurally CANNOT see a cross-package read or a path built in another file — so chasing completeness is futile. The detached adversarial audit already caught every reader-axis hole the sweep missed; the audit IS the right, complete-enough backstop for that axis. The deliverable here is a **net-negative diff**: less code, not more.

## Direction (for ideation)

- **Remove** the reader-axis taint/discovery machinery: `readsInstructionContent` and its helpers (param-flow, struct-field/method/closure tracking, path-construction reconstruction, the transitive reader fixpoint) plus the reader-axis planted-control tests that only exist to exercise that machinery.
- **Keep** the minimal MATCH-axis core: a test that reads an instruction file's bytes and inspects them must self-classify (markNonAC / markCodeBoundInvariant). This is the universal, sound, small part — do not regress it.
- **Replace** the reader-axis static guarantee with a documented stance: the reader axis (does an undeclared instruction-file presence-check hide via an undiscovered read shape?) is covered by the **detached adversarial audit** required at every high-stakes-surface gate (validation-stage policy), NOT by a static sweep. Record this explicitly so it is a deliberate, audited scope boundary, not a silent gap.
- Measure the removal: report the before/after line count of the two `nonac_marker_test.go` files; the net diff for this task MUST be negative.
- Supersedes the prior "invert / go-types+SSA taint" direction — do NOT build an SSA pass; that is still "more machinery to detect crap." If ideation concludes some minimal reader signal is genuinely worth keeping, that is a finding to bring back, but the default + the captain's stated intent is removal.

## Out of scope

The match-axis core, the 56 demotions/re-binds, and the standing match-axis sweep + its self-test (all shipped in hwk #306) stay. This task removes the reader-axis detection machinery only.

## Acceptance criteria

**AC-1 — the reader-axis detection machinery is removed; the diff is net-negative.**
Verified by: the `readsInstructionContent` taint machinery + its reader-axis planted controls are gone (grep/AST shows the helpers removed); `git diff --stat` for this task reports more deletions than insertions in the two `nonac_marker_test.go` files; offline `go test ./...` stays green.

**AC-2 — the match-axis core still enforces "ingest ⇒ declare" and the reader-axis is documented as audit-backstopped.**
Verified by: a planted undeclared match-axis tautology (read instruction bytes + inspect, no marker) still REDs the sweep (the match-axis guard is intact, mutation-controlled); and the sweep's doc + the validation-stage policy state the reader axis is covered by the detached adversarial audit, not the static sweep (a reader-shape evasion is a documented, audited boundary — no longer claimed as guarded).

## Test plan

Offline Go refactor: delete the reader-axis machinery + its controls, keep the match-axis sweep + its mutation control, confirm `go test ./...` green and the diff net-negative. High-stakes shipped-test surface → a detached adversarial audit on the result (confirming the match-axis core still catches its class, and the removal didn't break the demotions/re-binds) before merge.
