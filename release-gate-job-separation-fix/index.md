---
id: bqqr8vzz152n8bk2dqf1rw4q
title: Release-gate fix — job-separation + tag-body OPTION B so the 0.19.6/0.20.0 cut doesn't fail like 0.19.5
status: validation
source: "captain + handoff (2026-06-05) — THE cut blocker, hard prerequisite for the main-flip/0.20.0 milestone. `release.yml` hard-requires a Runtime-Live-E2E on `next` that never auto-runs there, so the cut fails (like 0.19.5). Grounded analysis: Workflow task w569jug9c (session #08)."
score: "0.42"
started: 2026-06-05T19:05:55Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-release-gate-job-separation-fix
issue:
---

## Problem

`.github/workflows/release.yml` runs a single `goreleaser` job. Before the goreleaser action, that job has a hard step — **"Download latest journey metrics artifacts"** — that queries `gh run list --workflow "Runtime Live E2E" --branch next --status success --limit 1` and `exit 1`s when no successful run is found. `runtime-live-e2e.yml` fires ONLY on `workflow_dispatch` and `pull_request` (branches: `[next]`); there is **no `push: next` trigger** (verified: `.github/workflows/runtime-live-e2e.yml` `on:` block, and `next-publish.yml` is `workflow_dispatch`-only). So a merge into `next` never records a success-on-`next` run, the download step `exit 1`s inside the goreleaser job, goreleaser never runs, and the whole cut fails. This is what failed the 0.19.5 cut and would fail 0.20.0. It is the hard blocker for the `main-flip-0200-marketplace` milestone.

## Why "move steps after goreleaser" is wrong (and what the guard actually checks)

`internal/release/workflow_exec_guard_test.go`'s `assertReleaseWorkflowPublishesJourneyCosts` parses `release.yml` into a **flat, job-unaware step list** (it splits on `- name:` regardless of which job the step is in) and enforces, by **document order**:
- a journey-cost builder step exists (`go run ./cmd/spacedock-release journey-costs … --metrics-dir … --out …`) and is non-empty-checked (`test -s …`),
- a goreleaser publish step exists,
- `builderStep <= goreleaserStep` (else: "builds journey costs after goreleaser"),
- a `gh release upload "$GITHUB_REF_NAME" "$RUNNER_TEMP/journey-costs-v${RELEASE_VERSION}.json" --clobber` step exists with `publishStep > builderStep`.

Two consequences, both **spike-confirmed** (see the spike below):
1. Moving the journey steps to a sibling job placed **after** `goreleaser` in document order makes the guard FAIL (`builderStep > goreleaserStep`). The sibling job MUST be authored **before** the `goreleaser` job in `release.yml`'s text.
2. The guard is **`needs:`-blind**: it never inspects job DAG edges, so a `goreleaser: needs: journey-ledger` edge — which would re-block the cut on the never-fired run, the exact bug — passes the guard unchanged. The current guard therefore cannot distinguish a correctly-separated graph from a re-coupled one; the AC-1 separation invariant has no code gate today.

## Proposed approach (spike-validated)

**Job-separation:** split `release.yml`'s `goreleaser` job into two sibling jobs:
- **`journey-ledger`** — carries the three journey steps (Download metrics / Build cost ledger / Publish ledger), authored **before** `goreleaser` in the file so the flat-parser document-order guard stays green. Its download step degrades to a **non-fatal skip** (`echo "::warning::…"; exit 0`) when no successful Runtime-Live-E2E run exists, instead of `exit 1` — so a missing producer run no longer fails anything. It carries `needs: goreleaser` (one-way, the SAFE direction): the `gh release upload` requires the Release to exist, which only goreleaser creates, so journey-ledger waits for goreleaser; goreleaser does NOT wait for journey-ledger. The dispatch's "goreleaser does NOT `needs:` the sibling" is satisfied; the reverse one-way edge is required for upload correctness and the guard tolerates it.
- **`goreleaser`** — keeps Checkout / Set up Go / Extract release notes / Run goreleaser / Stamp plugin manifests. Cuts the release regardless of journey-ledger.

**Guard-hardening (closes the `needs:`-blind gap):** extend `assertReleaseWorkflowPublishesJourneyCosts` (or add a sibling assertion) to parse the job graph and FAIL if the goreleaser job declares `needs:` on the journey-ledger job. This turns AC-1's separation invariant into a code gate the binary enforces, not an instruction. Add an adversarial test in `journey_workflow_test.go` that injects `goreleaser: needs: journey-ledger` into the real `release.yml` and asserts the guard REDS — mirroring the existing comment-only adversarial tests.

**Tag-body OPTION B — RECONCILIATION (see "Finding" below):** OPTION B is ALREADY implemented and tested in the current tree. `release.AnnotatedTagArgs` already cuts a double-`-m` tag (`Release <version>` subject + body paragraph); `notes_extract_test.go` already round-trips `%(contents:body)` and locks the empty-body guard against the real `release.yml` awk/cat-file extraction; `notes.go`'s doc comments already describe the double-`-m` shape accurately. The 7h release-notes catch already landed (commits `42684bb3` "land notes in the tag BODY", `6fb6a36f`, `f676b955`, `bce54524`). The AC-2 "fix the now-false Go doc/tests" describes work that is already done — there are no false doc/tests to correct. AC-2 is recorded below as a verify-and-confirm AC (the suite proves the shape holds against the post-separation `release.yml`), NOT new corrective authoring. If the captain knows of a residual tag-body defect not visible in the tree, that needs to be named at the gate.

## Spike — riskiest unknown EXERCISED (result recorded)

**Question:** does the separated-job DAG actually let a cut succeed while `assertReleaseWorkflowPublishesJourneyCosts` stays green?

**Spike (throwaway, this session):** authored two candidate `release.yml`s differing only in document order of the `journey-ledger` sibling job (preserved at `/tmp/release-gate-spike/candidate-{before,after}.yml`), plus a one-way-`needs` variant, and ran the REAL unexported `assertReleaseWorkflowPublishesJourneyCosts` over each via a throwaway in-package test (`internal/release/spike_separation_test.go`, since removed). Result:

```
candidate "after"  (sibling job after goreleaser):    "release.yml builds journey costs after goreleaser"  (FAIL)
candidate "before" (sibling job before goreleaser):   <nil>                                                (PASS)
candidate "before+needs" (goreleaser needs: journey): <nil>                                                (PASS — guard is needs-blind)
```

**Findings:** (1) the separated-job fix is sound **only with the sibling job authored before `goreleaser`** in document order — non-obvious because the guard is job-unaware; (2) the guard does NOT enforce separation (the `needs:`-blind pass), so the fix's core invariant needs a new code gate, which becomes part of this task's scope. The full `internal/release/` suite stayed green (24/24) before the spike and after the throwaway spike test was removed.

## Peer coordination — DO NOT author across this boundary

The Codex peer's entity **4n owns the journey-ledger CONSUMER** (the producer/consumer of journey metrics on `next`). This release-gate fix touches only the `release.yml` job-graph that READS the metrics; it must ALIGN with 4n's consumer, not redefine it. The dependency to confirm with the peer:
- **What makes the Runtime-Live-E2E producer actually emit a success-on-`next` run** that `gh run list --branch next --status success` can find — or whether the producer side is being changed such that the journey-ledger job should query differently. This fix makes the absence non-fatal (skip), so the cut is unblocked either way; but whether the ledger is ever POPULATED depends on 4n's consumer/producer alignment. The implementation must NOT edit the consumer side; it must document the exact handoff and get the peer's confirmation that the non-fatal-skip download contract matches what 4n expects.

## Out of scope

The actual 0.20.0 cut + marketplace flip — that's the `main-flip-0200-marketplace` (pj) milestone (this fix is its prerequisite). A live tag push / real cut is out of scope; the DAG is proven by the guard + a dry-run, not by pushing a tag.

## Acceptance criteria

Each AC's proof comes from a source OTHER than the workflow/skill text under change — a Go test parsing the real `release.yml`, or git's own `%(contents:body)` round-trip — never a substring match over instruction prose.

**AC-1 — `release.yml` separates the journey-ledger gate from goreleaser so a cut is not blocked by a never-fired Runtime-Live-E2E run, and the separation is code-gated.**
Verified by: `internal/release/workflow_exec_guard_test.go`'s `TestWorkflowsPreserveAndPublishJourneyCosts` stays green against the new two-job `release.yml` (the journey-ledger sibling is authored before goreleaser, so the document-order guard passes); AND a new adversarial test asserts the hardened guard REDS when the real `release.yml` is mutated to add `goreleaser: needs: journey-ledger` (proves the separation is enforced by the binary, not merely instructed). The download step's missing-run path is a non-fatal skip (`exit 0`), exercised by a shell-level check that the skip branch exits 0 on an empty run-list and a guard check that accepts the skip form.

**AC-2 — the tag body uses OPTION B and the Go doc/tests match it (verify-and-confirm; the shape already landed).**
Verified by: `internal/release/notes_extract_test.go`'s `TestAnnotatedTagBodyRoundTrips` and `TestReleaseYAMLGuardRejectsEmptyBody` pass against the post-separation `release.yml` (the awk/cat-file extraction + double-`-m` body round-trip — notes land in the body, not the subject), and `go test ./internal/release/` is fully green. No new corrective authoring is expected unless the audit/captain names a residual defect not present in the current tree.

## Test plan

- **AC-1:** extend `assertReleaseWorkflowPublishesJourneyCosts` with a job-DAG parse that rejects a goreleaser→journey-ledger `needs:` edge; keep the existing flat-step ordering checks. Add the `needs:`-injection adversarial test alongside the existing comment-only ones in `journey_workflow_test.go`. Cost: ~1-2h (Go-only; the parser already extracts steps, job/`needs:` parsing is a small addition). Fixture-level (parses the real `release.yml`), no live workflow run.
- **AC-2:** the existing `notes_extract_test.go` real-git round-trip + the empty-body guard test read from the real `release.yml`, so editing the job structure re-exercises them; run `go test ./internal/release/`. Cost: ~0 new (verify-only).
- **Spike already paid** (the riskiest DAG/guard interaction). The implementation's first test seeds from the spike: the `needs:`-blind pass becomes the new adversarial red.
- **Peer handoff:** confirm with the Codex peer (4n) that the non-fatal-skip download contract aligns with the consumer/producer side before merge.
- **High-stakes CI/release machinery → detached adversarial audit before merge** (the dispatch's standing requirement). No live tag push in this task.

## Stage Report: ideation

- DONE: Job-separation design: a SIBLING journey-ledger job that goreleaser does NOT need (NOT moving steps after goreleaser, which breaks workflow_exec_guard_test.go); the exec-guard test stays green against the new release.yml
  Designed two-sibling-job split with the journey-ledger job authored BEFORE goreleaser in document order (the spike-validated requirement); guard-hardening added to scope because the spike found the guard is needs-blind. See "Proposed approach" + AC-1.
- DONE: Spike-first the riskiest unknown — the goreleaser needs-DAG + workflow_exec_guard_test interaction — via the smallest end-to-end (the guard test against the new YAML) and record the result in the task body
  Ran the REAL assertReleaseWorkflowPublishesJourneyCosts over before/after/before+needs candidates via a throwaway in-package test; result recorded in "Spike" section. Key finding: only the before-ordering passes, and the guard is needs-blind. Candidates preserved at /tmp/release-gate-spike/.
- FAILED: Tag-body OPTION B codified (single -m / awk) with the now-false Go doc/tests corrected so notes land in the tag body not the subject (cf. the 7h release-notes catch)
  Not fabricated: OPTION B is ALREADY implemented (double -m in release.AnnotatedTagArgs) and tested (notes_extract_test.go round-trips contents:body; the awk/cat-file extraction + empty-body guard already landed in 42684bb3/6fb6a36f/f676b955/bce54524). There are NO now-false doc/tests to correct — go test ./internal/release/ is 24/24 green. Recast as a verify-and-confirm AC-2; this discrepancy is surfaced for the gate.

### Summary

The real cut blocker is release.yml's pre-goreleaser "Download latest journey metrics" step that exit-1s when no success-on-next Runtime-Live-E2E run exists (and none auto-fires on next — no push trigger). Designed a two-sibling-job split (journey-ledger authored before goreleaser, one-way needs: goreleaser, non-fatal skip on missing run) and proved via spike that the fix is sound ONLY in before-ordering and that the existing guard is needs-blind, so guard-hardening is in scope. The dispatch's tag-body OPTION B work is already shipped/tested in the tree, so AC-2 is verify-and-confirm — flagged for captain/audit in case a residual defect exists that the tree doesn't show.

## Stage Report: implementation

- DONE: Implement the two-sibling release.yml split (journey-ledger authored BEFORE goreleaser, needs: goreleaser) AND policy-1: SKIP the Build/Publish-ledger steps when the Download step finds no producer run, so the journey-ledger job is GREEN/SKIPPED (never RED) on a producer-less cut
  `.github/workflows/release.yml` now has `journey-ledger` (needs: goreleaser) authored first, then `goreleaser` (no needs). Download step `id: download_metrics` exits 0 emitting `found=false` on empty/no-artifact run-list; Build + Publish are gated `if: steps.download_metrics.outputs.found == 'true'`. yq confirms: jobs order [journey-ledger, goreleaser], journey-ledger.needs=goreleaser, goreleaser.needs=null. Commit 45c2ea43 on spacedock-ensign/release-gate-job-separation-fix.
- DONE: Code-gate the separation: a new adversarial test REDs when release.yml is mutated to goreleaser:needs:journey-ledger, with the guard binding needs: to the OWNING job (not match-anywhere); workflow_exec_guard_test.go stays green; AC-1 tests the JOB-LEVEL skip consequence (journey-ledger job clean on a producer-less cut), not just download exit 0
  Added `parseWorkflowJobs` + `assertGoreleaserDoesNotNeedJourneyLedger` (binds needs: to the goreleaser-owning job vs the builder-owning job; safe reverse edge does NOT trip — verified by a throwaway DAG-print spike, since removed). `TestReleaseWorkflowGuardRejectsGoreleaserNeedsJourneyLedger` REDs on the injected edge. `assertReleaseLedgerStepsSkipWhenNoProducerRun` (AC-1 job-level) asserts Build+Publish gate on the download output; `TestReleaseWorkflowGuardRejectsUngatedLedgerBuild` REDs on a stripped gate; `TestReleaseDownloadSkipBranchExitsZeroOnEmptyRunList` EXERCISES the real download script with a stubbed empty `gh` and asserts exit 0 + found=false. `go test ./internal/release/` 28/28 green; full repo 1145/15 green.
- DONE: Verify AC-2 (tag-body OPTION B already shipped) holds against the post-separation release.yml; do NOT edit the 4n journey-ledger consumer/producer — flag the non-fatal-skip contract for the Codex peer's confirmation
  `TestAnnotatedTagBodyRoundTrips` + `TestReleaseYAMLGuardRejectsEmptyBody` pass against the post-separation file. The job reordering put a second `if [ ... ]; then` (the download run-id check) ahead of the empty-body guard, so I anchored `guardConditionRe` on `release-notes.txt` to keep selecting the right condition — selector repair my change necessitated, NOT corrective tag-body authoring. Did NOT touch runtime-live-e2e.yml or the `journey-costs` consumer. NON-FATAL-SKIP CONTRACT FOR 4n: release.yml's download step queries `gh run list --workflow "Runtime Live E2E" --branch next --status success --limit 1`; absence → exit 0 / found=false / ledger steps skip (cut still succeeds, ledger simply unpopulated). Whether the ledger is ever POPULATED depends on 4n making the producer emit a success-on-next run findable by that query — peer confirmation requested before merge.

### Summary

Split release.yml into journey-ledger + goreleaser siblings: journey-ledger authored first (document-order guard stays green), one-way needs: goreleaser, and POLICY 1 — the Download step is a non-fatal skip emitting found=false, with Build/Publish gated on it so a producer-less/empty-dir cut SKIPS cleanly (job green, never red) while goreleaser cuts the release regardless. Hardened the guard to bind needs: to the owning job (REDs goreleaser→journey-ledger, tolerates the safe reverse edge) and added a job-level skip-consequence guard, two adversarial tests, and a shell-level skip-branch exercise; AC-2 verified unchanged. Open handoff: the Codex peer (4n) must confirm the non-fatal-skip download contract matches the producer/consumer side before merge — the fix unblocks the cut either way but ledger population depends on 4n.

## Stage Report: validation

- DONE: The cut-unblock works — producer-less cut no longer fails; journey-ledger green/skipped (Build+Publish gated on download_metrics.found), goreleaser cuts regardless (no needs on journey-ledger). Reproduce via guard tests + the real download-skip-branch exit-0 exercise
  yq (real YAML parser) over release.yml: job order [journey-ledger, goreleaser], journey-ledger.needs=goreleaser, goreleaser.needs=null, Build/Publish both gated `if: steps.download_metrics.outputs.found == 'true'`. `TestReleaseDownloadSkipBranchExitsZeroOnEmptyRunList` runs the REAL download script vs a stubbed empty `gh` → exit 0 + found=false. Mutation C: inverting the no-run branch exit 0→exit 1 in the real file REDs that test ("want exit 0 on empty run list"); restored clean.
- DONE: Separation is code-gated — assertGoreleaserDoesNotNeedJourneyLedger REDs goreleaser:needs:journey-ledger AND tolerates the safe reverse edge (owning-job binding, not match-anywhere); workflow_exec_guard_test.go green; AC-1 job-level skip-consequence test present; mutation-confirm by stripping the gate / inverting the guard
  Mutation A: injected `goreleaser: needs: journey-ledger` into the REAL release.yml (scalar form) → TestWorkflowsPreserveAndPublishJourneyCosts REDs with the re-coupling message; restored. Adversarial refutation: flow-list forms `needs: [journey-ledger]` and `needs: [some-other-job, journey-ledger]` also caught (parseNeeds handles scalar+flow). Real file (safe edge on journey-ledger job) stays GREEN — owning-job binding discriminates direction. Mutation B: stripping the Build-step gate in the real file REDs TestReleaseWorkflowSkipsLedgerWhenNoProducerRun. Full internal/release suite 28/28 green.
- DONE: AC-2 (tag-body OPTION B) holds against the post-separation release.yml (round-trip + empty-body guard green); the 4n non-fatal-skip contract is flagged — validate the CODE/cut-unblock side, do NOT block on ledger population
  `TestAnnotatedTagBodyRoundTrips` + `TestReleaseYAMLGuardRejectsEmptyBody` green; AnnotatedTagArgs confirmed double-`-m` (notes.go:130). Probed old vs new `guardConditionRe`: the OLD greedy regex would have mis-selected the download step's `[ -z "$run_id" ] ...` condition (proven by direct regex match over the real file); the new `release-notes.txt` anchor correctly selects the empty-body guard — the selector repair was load-bearing, not corrective tag-body authoring. 4n producer/consumer dependency surfaced for the gate, NOT blocked on.

### Summary

PASSED. The cut-unblock and its code gate hold against the post-separation release.yml, validated by exercise not prose: yq confirms the DAG contract (goreleaser.needs=null, journey-ledger one-way needs: goreleaser, document-order [journey-ledger, goreleaser], Build/Publish gated on download_metrics.found); full repo suite green (28/28 in internal/release, all packages ok, vet+build clean). Adversarial mutation-confirm refuted nothing material: re-coupling edge (scalar+flow), stripped gate, and inverted skip-exit each RED the owning suite, and the real file's safe reverse edge stays green. AC-2 OPTION B round-trips through real git; the selector-regex repair is justified (the old greedy regex mis-selected the download condition). Standing handoff unchanged: ledger POPULATION depends on the Codex peer's 4n producer — out of this fix's scope, not a blocker, flagged for the gate.

### Feedback Cycles

**Cycle 1 — validation PASSED, detached adversarial audit found one MATERIAL guard hole (2026-06-05).** The validator recommended PASSED (cut-unblock + skip-gating sound; 28/28 + offline green). The detached audit (separate checkout of f1d2d871) found **M1**: the owning-job guard's `parseNeeds` (`internal/release/workflow_exec_guard_test.go:371`) only parses scalar/flow `needs:` — a future `goreleaser: needs:` written as a YAML BLOCK-LIST (`needs:` then `- journey-ledger`) EVADES `assertGoreleaserDoesNotNeedJourneyLedger` (parseWorkflowJobs reports goreleaser needs=[]), so the guard greenlights a real re-coupling edge, defeating its purpose. Scalar/flow/`[journey-ledger]`/multi-element forms are all correctly rejected; only the block-list shape evades, and the committed test injects only the scalar form. Routed to implementation (fresh dispatch — prior impl worker lost in the session crash; same worktree). Polish P1 (non-blocking): the download skip-branch reds on gh-error/missing-gh under `set -euo pipefail` (only empty-list is non-fatal) — matches the "no run exists" contract, does not block the cut.

## Stage Report: implementation (cycle 2)

- DONE: M1 fix — teach parseNeeds the YAML block-list needs: form so a block-list `goreleaser: needs: [- journey-ledger]` re-coupling edge is REJECTED; add the block-list shape as an adversarial case and mutation-confirm it REDs; the safe block-list reverse edge stays tolerated
  `parseWorkflowJobs` now consumes the deeper-indented `- name` entries following a bare `needs:` line (workflow_exec_guard_test.go); `parseNeeds` doc updated to match. `TestReleaseWorkflowGuardRejectsGoreleaserNeedsJourneyLedger` is now table-driven over scalar/flow/block-list and all three RED. Mutation-confirm: breaking the block-list consumption (`- ` → `ZZZNOMATCH`) reds ONLY the block_list subtest (scalar/flow stay green) — not tautological. Throwaway probe confirmed a block-list reverse edge on journey-ledger is attributed to journey-ledger (needs=[goreleaser]), goreleaser needs=[], guard does NOT trip — safe direction tolerated. Commit bd1fb72b.
- DONE: All existing bq tests stay green (skip-gating, AC-1 job-level skip-consequence, AC-2 tag-body, the scalar/flow guard cases); offline go test ./... green
  `go test ./internal/release/` 34/34 (was 28; +2 guard subtests reorganized table-driven, +4 P1 subtests); did NOT churn skip-gating, the AC-2 release-notes.txt anchor, or the existing mutation controls. `go vet ./internal/release/` clean; `go build ./...` clean; full offline `go test ./...` 1176 passed across 15 packages.
- DONE: P1 (optional, non-blocking) — tolerate gh-error/missing-gh in the release.yml download skip-branch as found=false
  Added `|| true` to the `run_id` command substitution so a gh error / missing gh degrades to an empty run_id → the existing no-run skip branch fires (found=false, exit 0) instead of aborting under `set -euo pipefail`. `TestReleaseDownloadSkipBranchToleratesGhError` exercises the REAL download script vs a non-zero gh stub AND a PATH with no gh, asserting exit 0 + found=false in both. Mutation-confirm: removing `|| true` reds both subtests (exit 1 / exit 127), the empty-list test stays green. Commit bd1fb72b.

### Summary (cycle 2)

Closed the M1 block-list evasion the cycle-1 detached audit found: `parseWorkflowJobs` now reads a `needs:` block-list (`needs:` then deeper-indented `- name`), and the re-coupling adversarial test injects the edge in all three YAML shapes (scalar/flow/block-list), each mutation-confirmed to RED specifically on the fix. The safe reverse edge (journey-ledger needs goreleaser) stays tolerated, including in block-list form (probe-verified). Also took P1: the download skip-branch now tolerates a gh error / missing gh as found=false rather than aborting under set -e, exercised by a new behavior-level test. Everything the cycle-1 audit verified sound was left untouched. internal/release 34/34, offline ./... 1176 passed, vet+build clean. Ready for FO re-validation + re-audit (confirm the block-list edge is now caught) before merge.
