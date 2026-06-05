---
id: 4qnn7dbzkyh9qv65t618vtxy
title: AC-3 sweep reader-axis — NET-REMOVE the detection machinery; rely on the detached audit backstop
status: validation
source: "captain (2026-06-05) — re-scoped from 'invert/harden the reader-axis sweep' to NET REMOVAL. Captain directive: 'i want to see net removal, not more crap to mark crap or detect crap.' The hwk merge added ~1525 lines of go/ast sweep machinery (the bulk = the reader-axis taint/discovery), which is BOTH the heaviest part AND the incomplete part (M-A/B/C/D evade it). The detached adversarial audit caught every reader-axis hole the static sweep missed — so the audit, not a static guard, is the right backstop for that axis."
score: "0.40"
started: 2026-06-05T19:10:17Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-sweep-guard-reader-axis-invert
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

---

# Ideation product

## Concrete removal plan

Two files, identical shape (`internal/hostneutrality/nonac_marker_test.go` = 731 lines, `skills/integration/nonac_marker_test.go` = 794 lines; 1525 total = the "~1525 lines" the seed cites). The reader axis is the **taint-flow** layer; the deletable bulk is the flow-tracking, NOT the read detection itself.

### Reader-axis machinery to DELETE (both files)

The flow/taint apparatus — what grows the reader set transitively and reconstructs paths across statements/params/fields:

- `instructionTaintedNames` (HN 330-367 / integ 341-387) — param-flow seeding + the `:=`/`=` local-taint fixpoint.
- `instructionTaintedFields` (HN 300-329 / integ 296-321) — package-wide struct-field-name taint scan.
- `readsTaintedField` (HN 285-299 / integ 279-295) — the `s.path` / method-receiver field-read flow.
- `lvalueName` (HN 368-387 / integ 388-399) — taint-key rendering for `recv.field`.
- `isStringyType` (integ 471-474, only in integration) — param-type taint gate.
- The **transitive reader fixpoint** loop in the sweep (HN ~169-183 / integ ~179-193) and its first-pass `helperCalls` discovery scaffolding that exists only to grow the reader set; the second pass keeps only the named-allowlist + direct-read check.
- The **name-taint branches** of `exprInstructionTainted` (the `tainted[x.Name]` / `tainted[recv.field]` ident+selector cases) — leaving only the literal/segment-and-(HN)path-ident branch, which is the direct-read predicate the minimal core keeps.
- The `taintedFields map[string]bool` parameter threaded through `readsInstructionContent` and the sweep — drops out once the field-taint flow is gone.

### Reader-axis planted controls to DELETE

These fixtures exist ONLY to exercise the flow machinery; without it they cannot RED, so they go:

- `TestSweepDetectsEvasionShapes` + its fixture helpers `readArg`/`walkSkills`/`readArg2`/`wrapHop`/`fixt.read` (integ 578-735, ~158 lines).
- `TestHostneutralitySweepDetectsEvasionShapes` + `wrapHop`/`readArg`/`walkSkills`/`fixt.read`/`assertRedThenGreenHN` (HN 552-725, ~174 lines) — but `assertRedThenGreenHN`/`writeFixture` are reused by the kept match-axis control, so the helper stays; only the flow-shape cases (path-arg, WalkDir, split-md, Join, struct-method) drop.
- The multi-hop case inside `TestHostneutralitySweepDetectsAnUndeclaredTautology` (HN 532-549) and Shape-1 multi-hop in the integration evasion test — these prove the fixpoint, which is gone.

### Match-axis core to KEEP (intact, mutation-controlled)

- `markNonAC`, `markCodeBoundInvariant` — the two declaration seams.
- `ingestedFileReaders` / `instructionFileReaders` — the named-reader allowlist.
- `instructionPathSegments` + `isInstructionPathLiteral`, `readSinks`, `fnFiltersInstructionMarkdown` — the **direct-read predicate** (a read sink whose path arg subtree carries a recognized instruction literal/segment, or a WalkDir `.md` collector). HN additionally keeps `instructionPathIdents` (the recognized package-var path allowlist — a declared list, not flow).
- The trimmed `readsInstructionContent` reduced to that direct predicate (no taint param).
- `TestNoUndeclaredTautologicalProof` / `TestNoUndeclaredHostneutralityTautology` — the offline sweep, undeclared-count metric.
- `TestSweepDetectsAnUndeclaredTautology` / `TestHostneutralitySweepDetectsAnUndeclaredTautology` (minus the multi-hop case) — the match-axis mutation control: a planted undeclared read via a named reader still REDs, declared GREENs, code-scan stays unflagged.
- `collectCalls`, `sortedUnique`/`sortedUniqueHN`, `containsStr`/`containsStrHN`, `writeFile`/`writeFixture`.

### Net-negative accounting

Deleted: ~158 lines (integ controls) + ~174 lines (HN controls, less the reused helper) + the flow functions above (integ: 17+30+47+22+25 = 141; HN: 15+30+38+20 = 103) + the fixpoint/taint-branch trims. Estimated **~400-450 net lines removed across the two files**; the minimal core re-uses functions that already exist. The implementation reports actual `git diff --stat` per file; the per-file diff for each `nonac_marker_test.go` MUST be net-negative.

## Riskiest-mechanism spike (done in ideation — the result seeds the implementation's first test)

The unverified mechanism: *does the match-axis guarantee survive deleting all the flow machinery, or does some real declared test silently fall out of the sweep's view?* Exercised both packages with a throwaway minimal sweep (named allowlist + direct literal/segment/path-ident read + WalkDir, NO flow):

- **integration: 0 undeclared offenders** on the real package; all 28 named-reader-detected + 17 previously-taint-only tests stay covered, because every integration taint-only test reads via a literal path. A planted undeclared tautology REDs both via a named reader (`foCore`) AND via an inline direct-literal `os.ReadFile`.
- **hostneutrality: 0 undeclared offenders**, BUT a reader-detection probe exposed a real nuance — two declared tests stop being *detected as readers* without flow:
  - `TestDevDisciplinesSurviveInDevHomes` reads `os.ReadFile(h.path)` ranging over the package var `devHomePresence` (M-D range/slice flow). The read arg is `h.path`, not the `devHomePresence` ident → invisible to a pure direct-read core.
  - `TestLiveScenarioRecommendedPracticePresent` reads a local `path` assigned from `filepath.Join(literals)` in a separate statement (cross-statement local taint) → invisible to a pure direct-read core.
  - (`TestClaudeAdapterOwnsRelocatedCommands` reads `os.ReadFile(claudeRuntimePath)`, a recognized `instructionPathIdents` var → STILL detected.)

This is exactly the silent-coverage-loss the proof policy warns about, surfaced before committing to the plan. It is the one finding to bring to the gate (below), not a unilateral decision.

## Finding for the gate — the minimal reader signal worth keeping, and two HN tests

Per checklist item 3 ("default is removal; if a minimal reader signal is genuinely worth keeping, bring it as an explicit finding to the gate"):

1. **The minimal reader signal worth keeping = the DIRECT-read predicate**, not flow tracking. It is small (already implemented by `isInstructionPathLiteral`+`readSinks`+`fnFiltersInstructionMarkdown`+`instructionPathIdents`), sound, and catches every real test that reads via a literal/segment/recognized-var path. It is NOT an SSA/go-types taint pass — it adds zero new detection machinery; it is what remains after the flow layer is deleted. Recommendation: keep it.
2. **Two HN declared tests (`TestDevDisciplinesSurviveInDevHomes`, `TestLiveScenarioRecommendedPracticePresent`) lose sweep coverage** under removal because they read via range-var / cross-statement-local flow. Both are `markNonAC` text-consistency lints; `TestLiveScenarioRecommendedPracticePresent` reads `skills/commission/references/templates/development.md` and `TestDevDisciplinesSurviveInDevHomes` ranges over `docs/dev/*` dev-home paths (not shipped LLM-instruction surfaces). Captain's call — three honest options, presented at the gate:
   - **(a) Accept the boundary** (default per directive): these two reader shapes become audit-backstopped, same as M-A/B/C/D. They keep their `markNonAC` markers; the sweep simply no longer re-derives them as readers each run. The detached adversarial audit is the declared backstop.
   - **(b) Re-home the read to a literal** so the kept direct-read core still sees them with zero flow machinery: `TestLiveScenarioRecommendedPracticePresent` is a one-line inline-literal change; `TestDevDisciplinesSurviveInDevHomes` would inline its dev-home paths as literals.
   - **(c) Name the helper** — n/a here, these inline-read in the Test body (no reusable helper to add to the allowlist).
   Ideation recommends **(a)**: it matches the captain's stated intent (net removal, audit as backstop) and avoids re-introducing per-test detection accommodation. (b) is a cheap fallback if the captain wants zero coverage loss.

## Refined acceptance criteria

**AC-1 — reader-axis flow machinery removed; per-file diff net-negative; offline tests green.**
Verified by: `grep` shows `instructionTaintedNames`/`instructionTaintedFields`/`readsTaintedField`/`lvalueName`/`isStringyType` gone from both `nonac_marker_test.go` files and the reader-axis evasion-shape controls removed; `git diff --stat` reports more deletions than insertions for EACH `nonac_marker_test.go`; `go test ./internal/hostneutrality/ ./skills/integration/` (and `go test ./...`) green.

**AC-2 — match-axis core intact and mutation-controlled.**
Verified by: the kept mutation control RUNS and REDs a planted undeclared read via a named reader, GREENs once it declares `markNonAC`/`markCodeBoundInvariant`, and leaves a code-scanning invariant unflagged (the existing `TestSweepDetectsAnUndeclaredTautology` / `TestHostneutralitySweepDetectsAnUndeclaredTautology` minus the dropped fixpoint case, run via `go test`); the real-package undeclared-offender count stays zero in both packages.

**AC-3 — the reader axis is documented as detached-audit-backstopped, not statically guarded.**
Verified by: the sweep docstrings + the validation-stage policy state the reader axis (an undeclared instruction-file read hiding via an undiscovered read shape) is covered by the detached adversarial audit at the high-stakes gate, NOT the static sweep — the M-A/B/C/D shapes AND (if the captain accepts gate option (a)) the two HN flow-read tests are a documented, audited boundary. This is a text-consistency property of the shipped sweep doc + workflow README (a `markNonAC`-class lint, not a behavioral AC); its real proof is the detached adversarial audit required before merge (test plan).

## Spike determination

NOT "no spike needed" — the design's soundness rested on the unverified claim that the match-axis guarantee survives flow-machinery deletion. That spike is DONE above (both packages exercised with a minimal sweep + planted tautologies); it surfaced the two-HN-test coverage nuance that became the gate finding. The implementation's first test is the AC-2 mutation control run against the trimmed core.

## Stage Report: ideation

- DONE: Concrete removal plan naming the reader-axis machinery to delete (readsInstructionContent + its param/struct-field/method/closure flow, path-construction reconstruction, transitive reader fixpoint) and the reader-axis planted controls; the net-negative diff is the deliverable (report before/after line counts of the two nonac_marker_test.go files)
  See "Concrete removal plan": named the flow functions (instructionTaintedNames/Fields, readsTaintedField, lvalueName, isStringyType, the fixpoint loop, the name-taint branches) + the evasion-shape controls; before line counts recorded (HN 731, integ 794 = 1525); est. ~400-450 net lines removed, per-file diff MUST be net-negative.
- DONE: Keep the match-axis core intact and mutation-controlled (a planted undeclared match-axis tautology still REDs the sweep) and document the reader axis as detached-audit-backstopped, not statically guarded
  AC-2 keeps the existing mutation control (planted named-reader tautology REDs / declared GREENs / code-scan unflagged); AC-3 documents the reader axis as audit-backstopped. Spike confirmed a planted match-axis tautology still REDs against the trimmed core in both packages.
- DONE: Default is removal: if ideation concludes a minimal reader signal is genuinely worth keeping, bring it as an explicit finding to the gate; do NOT design an SSA/go-types taint pass
  Finding for the gate: the minimal signal worth keeping is the existing DIRECT-read predicate (literal/segment/path-ident + read sink + WalkDir) — zero new machinery, NOT an SSA pass. Plus a surfaced nuance: two HN markNonAC tests (range-var / cross-statement-local reads) lose sweep coverage; three honest options presented, ideation recommends accepting the audit-backstopped boundary (option a).

### Summary

Re-scoped the captain's net-removal directive into a concrete deletion plan: the deletable reader axis is the taint-FLOW layer (param/field/method/closure flow + cross-statement path reconstruction + the transitive reader fixpoint) plus the planted controls that only exist to exercise it; the kept match-axis core is the named-reader allowlist + the small DIRECT-read predicate (a read sink whose path arg carries a recognized instruction literal/segment/var, or a WalkDir `.md` collector) + the mutation control. Ran the riskiest-mechanism-first spike (throwaway minimal sweep against both real packages + planted tautologies): integration loses zero coverage; hostneutrality surfaced two declared markNonAC tests that read via range-var / cross-statement-local flow and would silently fall out of the sweep — brought to the gate as a finding with three options (recommend: accept the audit-backstopped boundary). Refined the ACs so each is independently checkable by grep/git-diff-stat/go-test plus the required detached adversarial audit before merge.

## Stage Report: implementation

- DONE: Remove the reader-axis taint-FLOW machinery from BOTH nonac_marker_test.go files (instructionTaintedNames/Fields, readsTaintedField, lvalueName, isStringyType, the transitive reader fixpoint, the name-taint branches of exprInstructionTainted) plus the reader-axis planted-control tests that only exercise it; per-file diff net-negative (report before/after line counts)
  Commit 37ba4cf7 (worktree branch spacedock-ensign/sweep-guard-reader-axis-invert). All named helpers gone (grep -E '...' over both files exits 1 / no matches). Net-negative per file: HN 731->407 (-324), integration 794->430 (-364); `git diff --stat` = 149 insertions / 836 deletions. Removed the fixpoint loop + helperCalls scaffolding + the taintedFields param thread; readsInstructionContent/exprInstructionTainted reduced to the direct predicate. Also dropped TestSweepDetectsEvasionShapes / TestHostneutralitySweepDetectsEvasionShapes + assertRedThenGreen[HN] and the multi-hop fixpoint case in the kept HN control.
- DONE: Keep the match-axis core intact and mutation-controlled: a planted undeclared match-axis tautology still REDs the sweep, declared GREENs, a code-scan invariant stays unflagged (the kept direct-read predicate + named-reader allowlist); offline go test ./... green
  TestSweepDetectsAnUndeclaredTautology + TestHostneutralitySweepDetectsAnUndeclaredTautology run via `go test` and PASS: a planted named-reader read REDs, GREENs once it declares markNonAC/markCodeBoundInvariant, code-scan (scanFile over ../dispatch) stays unflagged. Real-package undeclared-offender count stays zero in both (verified pre-edit with a throwaway flow-free probe — integration 34 detected-readers/0 offenders, HN 9/0 — and post-edit by the offline sweeps passing). `go test ./...` = 1164 passed in 15 packages; build success; gofmt clean.
- DONE: Document the reader axis as detached-audit-backstopped (sweep docstring + the validation-stage policy): option (a) accepted — the 2 HN flow-read tests keep their markNonAC markers but become an audited boundary, not a silent gap
  Both sweep docstrings (TestNoUndeclaredTautologicalProof / TestNoUndeclaredHostneutralityTautology) rewritten: MATCH axis = STATICALLY GUARDED, READER axis = DETACHED-AUDIT-BACKSTOPPED (names M-A/B/C/D + the two HN range-var/cross-statement reads as the audited boundary). docs/dev/README.md validation-stage policy gains an "AC-3 sweep reader axis — this audit IS the backstop" sub-bullet under the detached adversarial audit block. Workflow contract still VALID (`spacedock status --workflow-dir docs/dev --validate`). The two HN tests (TestDevDisciplinesSurviveInDevHomes, TestLiveScenarioRecommendedPracticePresent) retain their markNonAC markers; confirmed they fall out of static reader detection per option (a).

### Summary

NET-REMOVED the reader-axis taint-flow machinery from both AC-3 sweep files per captain option (a): -688 net lines across the two files (each net-negative), zero new machinery. Validated the riskiest mechanism FIRST — a throwaway flow-free probe confirmed the trimmed direct-read core finds zero undeclared offenders in both real packages before committing to deletion. Kept the match-axis core + its mutation control (both PASS), reduced readsInstructionContent/exprInstructionTainted to the direct predicate, and documented the reader axis as detached-adversarial-audit-backstopped in both sweep docstrings and the validation-stage README policy. `go test ./...` green (1164 passed); the required detached adversarial audit on the merge result is the declared pre-merge backstop.
