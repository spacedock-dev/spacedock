---
id: nzb7wbwgj93m25ayf9b226xn
title: Gate the release/flip on the live e2e suite (today it's PR-only behind manual env approval)
status: ideation
source: "FO investigation (2026-06-08, during the 0.19.7 cut) — release.yml triggers on tag push and runs goreleaser only; runtime-live-e2e.yml triggers on pull_request to next + manual workflow_dispatch, and its live lanes require per-environment approval (CI-E2E*). So the actual release/tag has NO e2e gate, and PR-time e2e sits 'waiting' for manual approval (it did not gate the 0.19.7 merges)."
score: "0.3"
started: 2026-06-08T15:29:12Z
completed:
verdict:
worktree:
issue:
sprint: 0200-flip
group: release-gating
sprint-readiness:
---

The release/tag path runs no end-to-end check. For the 0.20.0 marketplace flip especially, the release must not publish without the live e2e suite having passed for the exact commit being released.

## Problem

- `release.yml` (`on: push: tags: v*`) runs `goreleaser` (which creates the Release) plus the best-effort `journey-ledger` job (`needs: goreleaser`, non-fatal). Nothing in the cut path consults the live e2e — a tag publishes regardless.
- `runtime-live-e2e.yml` (`on: pull_request: [next]`, `workflow_dispatch`) runs the shared runtime scenarios. The secret-free `offline` job runs on every PR; each live lane (`claude-live` sonnet on `CI-E2E` + opus on `CI-E2E-OPUS`, `codex-live` on `CI-E2E-CODEX`, `pi-live` on `CI-E2E-PI`) carries an `environment:` approval gate (reviewer = clkao) and only runs after that environment is manually approved.
- Net: a release (and the flip) can ship without any live lane having run green. Confirmed in the spike below — many recent PR runs sit parked at `status: waiting` because nobody approved the spending environments.

## Spike — riskiest unknown exercised (2026-06-08)

The chosen gate (below) hinges on one mechanism: can a release-time check reliably tell "the full live matrix ran green for this commit" apart from "only the secret-free offline job passed while the env-gated live lanes sat waiting"? Exercised directly against `spacedock-dev/spacedock` with `gh`:

- `gh run view <id> --json jobs` on a green run (`27050060639`): all five legs `conclusion: success` (`offline`, `pi-live`, `codex-live`, `claude-live (sonnet)`, `claude-live (opus)`). A workflow-run `conclusion: success` therefore means every live lane was approved and passed.
- `gh run view 27118281803` on a parked run: overall `conclusion: ""`, `status: waiting`; `offline` + `pi-live` green but `codex-live` + both `claude-live` legs still `waiting`. Approval is **per-environment and independent per lane**, and the overall run conclusion stays non-`success` until all required jobs finish.
- `gh run list --workflow "Runtime Live E2E" --status success` returns only fully-green runs and excludes every `waiting` / partial run. This is the SAME query the existing `journey-ledger` job already runs in `release.yml`, so the mechanism is proven inside the release runner.
- `gh run ... --json headSha` exposes the run's commit, so the gate can bind the green run to the exact release commit rather than accepting "some green run somewhere."

Result: gate shape (c) composes already-proven mechanisms (the `gh run list --status success` query the release runner already uses, `headSha` binding, the workflow-guard test pattern already covering `release.yml`). The one risk — `--status success` false-passing on an offline-only/parked run — is disproven: a parked run is never `success`. No further spike needed for the chosen shape.

## Decision — gate shape (c): release-time precondition that goreleaser depends on

Add an `e2e-gate` job to `release.yml` that goreleaser `needs:`. It resolves the tagged commit SHA (`git rev-list -1 "$GITHUB_REF_NAME"`) and requires a **`conclusion: success` Runtime Live E2E run whose `headSha` equals that commit**. If none is found, the job fails and goreleaser never runs — the cut is blocked. A green run for the release commit means every live lane was approved and passed.

Rejected alternatives (recorded so the choice is on the record):

- **(a) Tag-triggered live lane blocking goreleaser.** The live lanes spend real API keys and are deliberately per-environment approval-gated; `pull_request` (not `pull_request_target`) withholds secrets from forks by design. Re-running the full matrix on the tag would (i) duplicate API spend already paid on the PR to `next`, (ii) still sit `waiting` for manual environment approval at release time (so it would NOT cleanly "block the cut" — it would hang it), and (iii) move secret-spending into the release path. Rejected.
- **(b) Required PR-time live check before merge to `next`.** Making the live lanes a required status check forces manual env approval + API spend on every PR to `next`, or forces auto-approving the spending environment (a security regression away from the deliberate selective-approval discipline the spike showed in use). Rejected for the routine PR path; the flip-time cut is where the gate must bite.

Gate (c) spends nothing extra, reuses the proven query, and binds the cut to a green live run that already happened on the line being released.

Scope: at minimum the **0.20.0 flip** cannot publish without a green live e2e for its commit. The gate applies to every `v*` tag (the flip is a `v*` tag), with a narrow, explicit captain-waiver escape hatch (`SPACEDOCK_E2E_GATE_WAIVER`) so an emergency cut is possible but auditable — mirroring the main-flip entity's "Final Runtime Live E2E on the prepared tip, unless the captain explicitly waives it."

### Implementation co-edit — required when the `needs: e2e-gate` edge lands (blast radius)

Adding `needs: e2e-gate` onto the goreleaser job changes the goreleaser job header from `  goreleaser:\n    runs-on: macos-latest` to a header carrying the new `needs:` line. `internal/release/journey_workflow_test.go` anchors on that byte-literal header in four places and asserts goreleaser has NO `needs:` in a fifth, so the bare edge reds 14 cases across 5 functions (empirically: `go test ./internal/release/` goes from green to **0 passed / 14 failed** with the edge alone). This co-edit belongs to **this entity** (it is the one adding the edge) and MUST land in the same change as the `release.yml` edit, or the implementer ships a red package believing AC-4 passed.

When the `needs: e2e-gate` edge lands, also update these sites in `internal/release/journey_workflow_test.go` so the anchors track the new goreleaser header (carrying `needs: e2e-gate`) and the no-`needs` assertion becomes a needs-only-`e2e-gate` assertion:

- The four byte-literal anchors on `  goreleaser:\n    runs-on: macos-latest`:
  - `TestReleaseWorkflowGuardRejectsGoreleaserNeedsJourneyLedger` — the `strings.Replace` anchor (≈:187).
  - `TestReleaseWorkflowGuardRejectsGoreleaserNeedsJourneyLedgerViaJobIdentityShapes` — the anchor/alias-needs mutation (≈:231) and the quoted-job-key mutation (≈:238).
  - `insertGoreleaserCarrier`'s `const anchor` (≈:274), used by the two multi-carrier tests.
- The no-`needs` assertion in `TestReleaseWorkflowJobGraphMatchesGitHubActions` (≈:424): `goreleaser` must now need exactly `["e2e-gate"]` (it currently asserts `len(needs["goreleaser"]) != 0` is an error). Keep the journey-ledger separation it pins (goreleaser still must NOT need `journey-ledger`).

The separation property AC-4 still protects is genuine — goreleaser needs `e2e-gate`, never the `journey-ledger` job — so the co-edit re-aligns the anchors to the new header rather than deleting the guards. Line numbers are approximate; locate by the byte-literal and the function names.

## Acceptance criteria

**AC-1 — A `v*` tag cannot reach goreleaser without a green live e2e for the release commit.**
The `goreleaser` job declares `needs: e2e-gate`, and the `e2e-gate` job fails when no `conclusion: success` Runtime Live E2E run exists whose `headSha` equals the tagged commit.
Verified by: a Go workflow-guard test (`internal/release`, the existing pattern) that parses `release.yml` and asserts (i) the goreleaser-carrying job's `needs:` includes `e2e-gate`, and (ii) the `e2e-gate` job's run block resolves the tagged commit SHA and gates on `gh run list/view --workflow "Runtime Live E2E" --status success` matched against that SHA — and a paired adversarial variant (string-substituted to drop the `needs:` edge / weaken the SHA match) that the guard REJECTS. The expected relationship is parsed from the real workflow YAML, independent of any instruction-file prose.

**AC-2 — The gate's SHA-matching predicate accepts a green run for the commit and rejects a parked/wrong-commit run.**
The decision logic (extracted as a pure function in `internal/release`, called by the `e2e-gate` step via a `spacedock-release` subcommand, mirroring `journey-costs`) returns "pass" only for a run with `conclusion == success` AND `headSha == release commit`, and returns "block" for a `waiting`/empty-conclusion run, a `success` run on a different `headSha`, and an empty run list.
Verified by: Go unit tests over the pure function feeding fixture run-list JSON (the four cases above), asserting pass/block per case. Expected values come from the constructed fixtures, not from the workflow file.

**AC-3 — The captain-waiver escape hatch bypasses the gate only when explicitly set, and is auditable.**
With `SPACEDOCK_E2E_GATE_WAIVER` set to a non-empty reason, the predicate returns "pass (waived)" and emits the reason to the step log / `$GITHUB_STEP_SUMMARY`; unset, the gate enforces normally.
Verified by: Go unit tests over the predicate for set/unset waiver, asserting the waived branch passes and records the reason, and the unset branch enforces. (Test-only proof of the predicate; the workflow wiring of the env var is covered by AC-1's guard.)

**AC-4 — `go test ./internal/release/` is green over the WHOLE package after the change, and the `e2e-gate` edge does not re-block the cut on the never-fired journey-ledger producer run.**
The `e2e-gate` job consults the live-e2e RUN history (the same `gh run list` the ledger uses) but does NOT carry the `journey-costs` builder and is NOT the ledger job, so goreleaser may `need` `e2e-gate` while still not needing the ledger (the journey-ledger separation stays a genuine constraint, not a side effect of the byte-literal anchors). The acceptance bar is the WHOLE `internal/release` package green — not merely "the separation guard does not false-trip."
Verified by: `go test ./internal/release/` passing over the entire package after the `needs: e2e-gate` edge AND its required co-edit (see the "Implementation co-edit" section above) both land. This is NOT satisfied by the byte-literal anchors surviving — adding the edge by itself makes the anchored header vanish and reds 14 cases; the co-edit re-aligns the anchors to the new header so the package is green for the right reason (goreleaser needs `e2e-gate` only, never the ledger).

## Test plan

- **What verifies it:** Go tests in `internal/release` only — no live run needed at implementation time. The workflow-guard tests parse the real `.github/workflows/release.yml` (via the existing `readWorkflow` helper) and assert the `needs:` edge + the SHA-gated query; paired adversarial string-substitution variants confirm the guard rejects a weakened workflow. The predicate unit tests drive fixture `gh run list` JSON through the pure decision function for the pass/block/waiver cases.
- **Cost/complexity:** Low. Reuses three established patterns (the `spacedock-release` subcommand shape, the `internal/release` workflow-guard tests, the `gh run list` query already in `journey-ledger`). New surface is one job in `release.yml`, one `spacedock-release` subcommand, one pure predicate + its tests, plus the required co-edit to `journey_workflow_test.go`'s four byte-literal goreleaser-header anchors + the no-`needs` assertion (see "Implementation co-edit" — the `needs: e2e-gate` edge reds 14 cases until those anchors track the new header).
- **Fixture vs CLI vs live:** Fixture + Go unit tests for the predicate; workflow-guard (offline YAML parse) for the wiring. **No live workflow test is required for this task** — the runtime behavior (live lanes actually passing) is `runtime-live-e2e.yml`'s job, already proven; this task only proves the CUT consults that result. The end-to-end observed gate is exercised for real at the flip: the 0.20.0 cut will only proceed after a green live run for the prepared tip exists, which is the main-flip entity's own acceptance step.
- **Observed/dry-run gate (satisfies the "not just prose" bar):** before/at the flip, confirm on a throwaway test tag (or a `--dry-run`-style invocation of the predicate against real `gh run list` output) that the gate BLOCKS when no green run matches the commit and PASSES once one does. This is the smallest end-to-end exercise of the cut-blocking path; the spike above already proved the underlying query distinguishes green from parked.

## Notes

Provenance: surfaced during the 0.19.7 cut; spiked 2026-06-08. Part of sprint `0198-pre-flip-hardening`, group `release-gating`. The `release-gate-job-separation-fix` (bqqr) is already landed in `release.yml` (the `journey-ledger` job no longer blocks goreleaser); this task adds the e2e precondition without reintroducing that coupling (AC-4). Directly serves the main-flip milestone's "Final Runtime Live E2E on the prepared tip, unless the captain explicitly waives it" by making that step an enforced gate rather than a remembered manual step.

## Stage Report: ideation

- DONE: Decide the gate shape: pick among (a) a pre-release/tag e2e gate that blocks goreleaser, (b) a required PR-time e2e check before merge to next (env auto-approval for trusted branches), (c) a manual workflow_dispatch e2e the captain runs before the flip — grounded against the actual release.yml + runtime-live-e2e.yml triggers.
  Chose (c) implemented as a release-time `e2e-gate` job that goreleaser `needs:`; (a) and (b) rejected with recorded reasons (API-spend/secret-path, env-approval would hang the cut, security regression). Grounded against the real workflow files (read release.yml jobs/triggers + runtime-live-e2e.yml's per-environment approval model).
- DONE: Spike the riskiest unknown: confirm the chosen trigger is actually achievable — record the result.
  Ran `gh run list/view` against the live repo: a green run = all 5 legs `success`; a parked run = `conclusion:""`/`waiting`; `--status success` excludes parked runs; `headSha` binds a run to a commit. The exact query is already used by the existing `journey-ledger` job. Result recorded in the entity's "Spike" section.
- DONE: Produce build-ready ACs + test plan: at minimum the 0.20.0 flip cannot publish without a green live e2e — verified by the workflow definition + an observed/dry-run gate, not prose.
  Four ACs (needs-edge guard, SHA-match predicate, captain waiver, no journey-ledger re-coupling) each verified by Go tests over real workflow YAML / fixture run-list JSON — expected values independent of the file under test. Test plan names the verifier, cost (low, reuses 3 existing patterns), fixture-vs-live split (no live test at impl time), and the observed/dry-run blocking exercise. Baseline `go test ./internal/release/` green (54 passed).

### Summary

Firmed gate shape (c): an `e2e-gate` job in `release.yml` that goreleaser `needs:`, requiring a `conclusion: success` Runtime Live E2E run whose `headSha` equals the tagged commit, with an auditable captain-waiver escape hatch. The riskiest unknown — whether a release-time check can tell a fully-approved green run from a parked offline-only run — was exercised directly against the live repo and disproven as a risk (parked runs are never `success`; the query is already proven inside the release runner). Proof is Go workflow-guard tests + a pure predicate's unit tests, reusing the established `internal/release` + `spacedock-release` subcommand patterns; the implementer starts from a green 54-test baseline.

### Staff-review fold

Folded preflight finding M2 (`docs/roadmap/0200-flip/staff-review.md`) into this banked ideation — a NARROW correction, not a re-ideation. The other ACs, the gate shape, and the spike were left untouched; the e2e-gate itself is NOT implemented here (that is the Commander's implementation, later).

- DONE: Reframe AC-4 from "the separation suite stays green" (empirically FALSE) to "`go test ./internal/release/` is green over the WHOLE package after the change."
  AC-4 title + Verified-by rewritten: the bar is the whole `internal/release` package green after the edge AND its co-edit land, NOT the byte-literal anchors surviving. The journey-ledger separation it protects is preserved (goreleaser needs `e2e-gate`, never the ledger).
- DONE: Add the explicit implementation co-edit to nzb's blast radius (the entity adding the edge owns it, at implementation time).
  New "Implementation co-edit" subsection after Scope names the four byte-literal `goreleaser:`/`runs-on: macos-latest` anchors in `internal/release/journey_workflow_test.go` (≈:187/:231/:238/:274) + the no-`needs` assertion in `TestReleaseWorkflowJobGraphMatchesGitHubActions` (≈:424) that must track the new `needs: e2e-gate` header. Also folded into the Test-plan cost line so blast radius is honest.
- DONE: Verified the anchor claim against the real files before writing (not re-read — exercised).
  Read `internal/release/journey_workflow_test.go` (the four anchors + the :423-425 no-`needs` assertion) and `.github/workflows/release.yml` (the real `  goreleaser:\n    runs-on: macos-latest` header at :90-91). Then reproduced M2 empirically: patched a scratch release.yml with `needs: e2e-gate`, `go test ./internal/release/` went 59→**45 passed / 14 failed** across exactly the 5 cited functions, then restored release.yml and confirmed 59 passed (clean, no diff). The 14-failure break is structural and matches M2 exactly.
