---
id: zvk9cnew2ggpaqb3wty24xtf
title: Remove redundant lint mirrors and the version-gate shell harness
status: validation
source: "0.27 cut audit (2026-08-14), adversarially verified; captain directed filing"
started: 2026-08-15T02:55:48Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-remove-redundant-lint-mirrors
issue:
gates:
    version: 1
    records:
        - id: gate:zvk9cnew2ggpaqb3wty24xtf:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:zvk9cnew2ggpaqb3wty24xtf-backlog-1
              briefing:
                id: briefing:zvk9cnew2ggpaqb3wty24xtf:backlog:attempt-1:revision-1
                digest: sha256:1a4b86e71e3145a2dbd052a2a7f4aa553d77bf031e8b36b175ffe0af187ef39e
                request-digest: sha256:24538405c9c80a26b1009a0a8b493178dab709db0c87507ab9bdcf9f918e7df1
                room-ref: ./remove-redundant-lint-mirrors/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:zvk9cnew2ggpaqb3wty24xtf:backlog:1
                briefing: briefing:zvk9cnew2ggpaqb3wty24xtf:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T02:54:02.147952Z"
                decision: approve
                reason: 'Captain ruling 2026-08-14 (dispatch them): approved into ideation'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:zvk9cnew2ggpaqb3wty24xtf:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:zvk9cnew2ggpaqb3wty24xtf-ideation-1
              briefing:
                id: briefing:zvk9cnew2ggpaqb3wty24xtf:ideation:attempt-1:revision-1
                digest: sha256:f08a7dd457ba1d216afae6b1c415b08ccf767b59716eda86e068de6222f90ca4
                request-digest: sha256:cf35bb141a5c8c378380bc539b47298b8f1d5e965e09f2eab0e8b55d6e8478ef
                room-ref: ./remove-redundant-lint-mirrors/review/ideation/briefing-1
              withdrawal:
                by: agent:first-officer
                at: "2026-08-15T03:50:40.91333Z"
                reason: 'Entity corrected post-prepare (cca195bae): contaminated line citation fixed 418 to 401 after clean re-run; re-preparing against current bytes'
            - id: gate-attempt:zvk9cnew2ggpaqb3wty24xtf-ideation-2
              briefing:
                id: briefing:zvk9cnew2ggpaqb3wty24xtf:ideation:attempt-2:revision-1
                digest: sha256:4022b65af8961c8adf3a2ab2fdcf60a518d535fe00f39b4c64c0e0de7195c7dc
                request-digest: sha256:37945b7440e2eed7497915496ca68e9e0a1b7e2bb9745801cc15b07fc4343547
                room-ref: ./remove-redundant-lint-mirrors/review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:zvk9cnew2ggpaqb3wty24xtf:ideation:2
                briefing: briefing:zvk9cnew2ggpaqb3wty24xtf:ideation:attempt-2:revision-1
                by: person:captain
                at: "2026-08-15T03:56:35.282507Z"
                decision: approve
                reason: 'Captain ruling 2026-08-15 (approve all except x8): approved into implementation'
              application:
                target-stage: implementation
                state: consumed
        - id: gate:zvk9cnew2ggpaqb3wty24xtf:validation
          stage: validation
          attempts:
            - id: gate-attempt:zvk9cnew2ggpaqb3wty24xtf-validation-1
              briefing:
                id: briefing:zvk9cnew2ggpaqb3wty24xtf:validation:attempt-1:revision-1
                digest: sha256:f7e6c509e767ef2efc58c4531e928329bc07727661f433e6c6a0b6d7099cc4ea
                request-digest: sha256:423479178d159c9ae82a80c4956cf19e8ec237c3472e367cf381280581c2bc72
                room-ref: ./remove-redundant-lint-mirrors/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:zvk9cnew2ggpaqb3wty24xtf:validation:1
                briefing: briefing:zvk9cnew2ggpaqb3wty24xtf:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T16:08:43.693618Z"
                decision: approve
                reason: 'Captain batch approval 2026-08-15: validation PASSED; land via releng-27 train'
              application:
                target-stage: done
                state: pending
mod-block: merge:pr-merge
pr: pr-merge:703
---

Three verified-redundant test mirrors.

1. The wantGaps XFAIL map (internal/contractlint/live_registry_reconciliation_test.go:51-66). A hand-copied oracle of the liveXFail() calls the same test already parses. The registry doc disclaims exactly this ("does not use a copied gap oracle"). Six lockstep two-file commits in range prove the churn. Keep parseLiveGap shape validation, duplicate-target rejection, and the TODO-owner join.
2. Three install-command literals in version_gate_smoke_test.go:37-39. TestInstallHintNoDrift already anchors the same tokens to their producer, docs/site/get-started/install.md. Keep the uname, go-build, and unsupported-OS tokens (lines 40-42) - they have no other pin.
3. skills/integration/testdata/version_gate_flow.sh plus its 351-line driver version_gate_fixture_test.go. The harness verifies a hand-written re-derivation of the FO prose, not shipped bytes. Its brew command already drifted from the lint-enforced prose and everything stayed green. Nothing binds mirror to prose.

Coordination: this supersedes scope items 2-3 of remove-startup-capability-probe (dav9), which rewrite the files this entity deletes. dav9 items 1 and 4 (shared-core prose, install.md) are unaffected. Whichever lands second reports the overlap as already done.

## Problem

Three test mirrors hand-copy something the tree already derives or already pins, and nothing binds copy to original. Each was confirmed against HEAD 68e64fa85; all cited line numbers hold there, and the four target files are byte-identical between 4d1912a69 and 68e64fa85 (`git log 4d1912a69..HEAD -- <targets>` is empty), so main advancing during ideation did not move this scope.

**1. The `wantGaps` oracle** (`internal/contractlint/live_registry_reconciliation_test.go:51-66`). `parseLiveJourneyCall`/`parseLiveGap` already extract every `liveXFail(...)`/`liveTODO(...)` binding from `internal/ensigncycle/shared_live_runner_test.go` by AST walk. `wantGaps` is a hand-written second copy of that same data, compared against the derived result 40 lines later. Its only possible failure is "you forgot to update the copy". `docs/runtime-live-ci-registry.md:349` states the join "does not use a copied gap oracle" — the code has contradicted its own registry doc. The cost is measured, not asserted: of the last 30 commits touching the reconciliation test, 15 are two-file commits editing the map and the source in lockstep (for example 4d7c72e16, f8fa527e4, 975a07f28, cae0a19d9, 437d20dc6). Every XFAIL retire or owner transfer pays a second-file tax for zero signal.

**2. Three install literals** (`internal/contractlint/version_gate_smoke_test.go:37-39`). `TestInstallHintNoDrift` already pins the same three tokens as a *relation* — the FO prose must carry the command install.md documents, extracted live from `docs/site/get-started/install.md`. The smoke-test copy freezes the identical strings as absolute literals. The frozen copy is strictly weaker and strictly costlier: it catches nothing the relation misses, and it makes any legitimate install-command change a three-file edit instead of two.

**3. The version-gate shell harness** (`skills/integration/testdata/version_gate_flow.sh`, 135 lines, plus its 351-line driver `skills/integration/version_gate_fixture_test.go`). The script is a hand-written shell re-derivation of FO prose; nothing binds it to the prose it mirrors. The proof is on disk: the script's macOS install command has read `brew tap spacedock-dev/homebrew-tap && brew install spacedock` since it was created in 0b5409687 (2026-07-31), while the shared-core prose has carried the tap-qualified `brew install spacedock-dev/homebrew-tap/spacedock` since b5e5ced7d (2026-06-02) — two months *earlier*. The mirror was born disagreeing with its original and stayed green for its entire life. A harness that cannot notice it contradicts the thing it mirrors verifies only itself.

## Proposed approach

Three deletions. No new mechanism, no replacement machinery — each deleted thing has a surviving verifier that was proven to fire (see Test plan).

1. `internal/contractlint/live_registry_reconciliation_test.go` — delete the `wantGaps` literal and the loop that consumes it (lines 51-66) and the then-unused `"reflect"` import. Everything else in the file stays, including the duplicate-gap-target rejection, the `parseLiveGap` shape validation, the TODO-owner join, and the workflow-selector counts.
2. `internal/contractlint/version_gate_smoke_test.go` — delete the three install-command literals from `TestVersionGateProseOSAwareHint`'s token slice (lines 37-39). Keep `uname -s`, `go build -o spacedock ./cmd/spacedock`, and `unsupported OS`, which have no other pin. Update the test's doc comment, which currently claims the test carries "BOTH install leads" and stops being true.
3. Delete `skills/integration/version_gate_fixture_test.go` and `skills/integration/testdata/version_gate_flow.sh` outright. No symbol either file defines (`writeExe`, `gateFixtureDir`, `runGateFlow`, `captiveInstall`, `installRunCount`, `sentinelPath`, `writeGateLauncher`, `invocationLog`, `fixtureSessionID`, `curlInstallToken`, `withdrawCapability`) is referenced anywhere else in the package — verified by grep and by compiling the package after deletion.

**Documentation:** none required. This changes no CLI output, command surface, startup banner, or host integration — it is test-only deletion. Note the inverse of the usual case: deletion 1 makes `docs/runtime-live-ci-registry.md:349` true, where it was false before, so the doc needs no edit.

**Residual, named honestly.** The harness's behavioral assertions (one-attempt sentinel bound, create-before-run ordering, `SPACEDOCK_BIN` repoint, exit codes) get no executable successor; `TestInstallGateSentinelLoopBound` pins the *instruction text* for those behaviors, not their effect. This is not a coverage loss, because the harness never exercised the shipped artifact: the FO is an LLM reading prose, and the script is a separate hand-written approximation that provably disagreed with that prose for two months without a single red test. Deleting it stops the tree from appearing to hold executable proof it never held.

## Coordination with remove-startup-capability-probe (dav9)

dav9 scope items 2 and 3 rewrite the two files this entity deletes; items 1 (shared-core prose) and 4 (install.md) are untouched by this entity. The reconciliation is order-independent:

- **This entity lands first:** dav9 items 2-3 are already done — the files are gone. dav9's AC-1 ("`version_gate_flow.sh` passes the gate without probing `gate --help`") and AC-2 ("all remaining version-gate fixture tests pass") lose their subject and must be re-expressed rather than reported satisfied; dav9's real remaining proof is its AC-3 (no `REQUIRED_CAPABILITY` / `Missing capability` / `gate withdraw` probe text left in `skills/`) and AC-4 (the surviving contractlint prose tests still pin the minor-version check and the binary-absent/wrong-version abort classes). Deleting the harness also discharges part of dav9 AC-3 directly, since the `REQUIRED_CAPABILITY` references in `skills/` live in the deleted script.
- **dav9 lands first:** this entity deletes the already-trimmed files; the deletion is unaffected. Only the measured LOC shrinks, since dav9 removes roughly 60 lines from those two files first — expect about -440 instead of -506.

Whichever lands second reports the overlap as already done and states which measurement moved.

## Out of scope

TestDispatchAckMachineryIsAbsent (sole verifier of the hooks.json chain). The three-check failfast bundle (each check carries unique assertions: the reconciliation test counts selectors, TestRuntimeLiveCommonSuiteTimeouts pins the exact command strings including timeouts, TestRuntimeLiveCommonFailFastPolicy pins the per-runtime policy). TestInstallHintNoDrift.

## Expected surface and tolerance

Measured in the spike, not estimated: **4 files, -506 lines, 0 insertions.** Tolerance ±40 lines and ±1 file — the band absorbs the doc-comment rewrite in deletion 2 and the dav9 overlap (which moves the expected figure to about -440; landing inside the band under either ordering).

**Semantic changes: none.** No command grammar, stored format, authority, or runtime behavior changes. The one semantic that does move is deliberate and is the point of the entity: retiring a live XFAIL stops being a two-file edit and becomes a one-file edit.

## Acceptance criteria

**AC-1 (value) — Retiring one live XFAIL is a one-file edit that leaves the lint green.**
Measured against a baseline that moves the wrong way: removing a single `liveXFail(...)` from `internal/ensigncycle/shared_live_runner_test.go` and changing *no other file* must leave `go test ./internal/contractlint/` passing. At HEAD 68e64fa85 the identical one-file edit FAILS at `live_registry_reconciliation_test.go:64`. Two-file edit becomes one-file edit; an incomplete removal keeps it failing.

**AC-2 (value) — The tree gets smaller.**
Cumulative line delta against origin/main is negative, within the tolerance band above.

**AC-3 — The suite is no worse than its baseline.**
`go test ./...` and `go test ./... -race` report no failure that is not also present at HEAD on the same machine. Stated against baseline rather than absolute green because HEAD is **not** currently green: `TestCodexResolveManifestAgainstInstalledHost` (internal/cli) fails at clean HEAD reading `~/.codex/config.toml` ("Operation not permitted"), a local sandbox-permission condition unrelated to this change. Implementation must not chase it.

**AC-4 — The surviving lints still catch what the deleted mirrors caught.**
Falsifiable by mutation, not by grep: drifting the curl token or the brew one-line form in the FO prose fails `TestInstallHintNoDrift`; a malformed `liveXFail` owner fails the reconciliation through `parseLiveGap`. The keep-boundary tests remain present and passing.

## Test plan

Deletion plus the existing suite; the surviving drift and shape lints are the regression floor. All of the following were exercised in a throwaway worktree at HEAD 68e64fa85 and seed implementation's first checks.

- **AC-1 baseline and result.** Retiring the `keep-moving-posture` pi XFAIL as a one-file edit against HEAD fails with `keep-moving-posture gaps = ... want ...` at `live_registry_reconciliation_test.go:64` — and that line is the *only* failure, so `wantGaps` is the sole thing making it a two-file edit. The same one-file edit with the removal applied passes.
- **AC-4 mutation checks**, each confirmed to fail after the removal, so none is a tautology: replacing the curl URL in the FO prose fails `TestInstallHintNoDrift` at `install_hint_drift_test.go:73`; downgrading the brew one-line form to bare `brew install spacedock` fails it at line 94; changing a `liveXFail` owner to `"bad-owner"` fails the reconciliation at `live_registry_reconciliation_test.go:401` with `malformed liveXFail binding` (line 401, not 418, because the removal deletes 17 lines above it — the shift is itself a check that the removal was actually applied to the tree under test).
- **Keep-boundary confirmation.** `TestDispatchAckMachineryIsAbsent`, `TestInstallHintNoDrift`, `TestRuntimeLiveCommonSuiteTimeouts`, `TestRuntimeLiveCommonFailFastPolicy`, `TestRuntimeLiveGapBindingValidation`, and `TestVersionGateProseOSAwareHint` all present and passing after the removal.
- **Package-level green.** `internal/contractlint` and `skills/integration` pass plain and under `-race`; `gofmt -l ./internal ./cmd` clean.
- **Cost.** No new fixtures, no new helpers, no live tests. The riskiest unverified mechanism was "does the package still compile and stay green once the shell harness and its driver are gone" — exercised directly by deleting them and running both packages, rather than reasoned about.

## Stage Report: ideation

- DONE: Keep-boundary confirmed: dispatch-ack lint, the failfast bundle, TestInstallHintNoDrift, and the uname, go-build, and unsupported-OS tokens stay
  All six named tests run and pass with the removal applied; the three surviving tokens remain in `TestVersionGateProseOSAwareHint` and the failfast bundle's three members each keep a distinct assertion (selector counts / exact command strings / per-runtime policy).
- DONE: The dav9 overlap (its scope items 2-3) is reconciled in the design
  New "Coordination with remove-startup-capability-probe (dav9)" section resolves both landing orders and flags that dav9's AC-1/AC-2 lose their subject if this lands first and must be re-expressed, not reported satisfied.
- DONE: Value AC: an XFAIL retire becomes a one-file edit and the suite stays green
  AC-1 measures it against a moving baseline: the identical one-file XFAIL retire fails at HEAD (`live_registry_reconciliation_test.go:64`) and passes after removal (`ok internal/contractlint 1.601s`).
- DONE: Removal scope confirmed against HEAD (4d1912a69)
  All three targets exist at the cited lines; `-506` lines across 4 files with 0 insertions, `gofmt` clean, both affected packages green plain and under `-race`.
- DONE: Spike of the riskiest unverified mechanism
  Deleted the harness plus driver in a throwaway worktree and compiled/ran the package rather than reasoning about it; confirmed no other file references any of the 11 symbols those files define.

### Summary

The removal scope holds at HEAD and the design is now measured rather than estimated: -506 lines, 0 insertions, 4 files. Two claims in the seed got stronger under checking — the wantGaps lockstep tax is 15 two-file commits in the last 30 touching that test, not six; and the shell harness did not "drift" from the FO prose but was *born* contradicting it, carrying the unqualified `brew install spacedock` from its creation on 2026-07-31 while the prose had carried the tap-qualified form since 2026-06-02, green the whole time. Every acceptance criterion is now falsifiable by mutation or by a baseline that can move the wrong way, and I recorded the residual honestly: the harness's behavioral assertions get no successor, which costs nothing because they never exercised the shipped artifact.

Two things the gate should know. AC-3 is stated against the HEAD baseline rather than absolute green, because `go test ./...` is **not** green at clean HEAD on this machine — `TestCodexResolveManifestAgainstInstalledHost` fails reading `~/.codex/config.toml` ("Operation not permitted"), which I confirmed by a control run on an unmodified tree; implementation should not chase it. Separately, I hit a cross-agent incident: `git stash` writes to the repo-shared `refs/stash`, not per-worktree, and a concurrent ensign and I swapped payloads. I preserved the other agent's work to `/tmp/spacedock-stash-recovery-20260814/` and reported it to the team lead. No fleet agent should use `git stash` while peers are active on this repo.

### Evidence re-run (shared-path contamination)

My first spike used the bare shared path `<scratchpad>/spike`, which the FO later flagged as colliding across ensigns. Two symptoms confirm my run was contaminated: `worktree add` initially failed with "already exists" (I removed that directory, which likely destroyed a sibling's spike — my error, reported to the FO), and mutation C reported `live_registry_reconciliation_test.go:418`, the *unmodified*-file line number, meaning my deletions had been reverted underneath me mid-run.

Every cited result above was therefore re-run from scratch in `<scratchpad>/spike-remove-redundant-lint-mirrors` at HEAD 68e64fa85, using `git checkout --` rather than `git stash` for revert steps. Outcomes: identical. Line delta -506/+0 across 4 files, `gofmt` clean, the one-file XFAIL retire fails at line 64 before and passes after, all three mutations fire, six keep-boundary tests pass, and `internal/contractlint` + `skills/integration` are green plain and under `-race`. The single correction is mutation C's cited line, 418 to 401. No conclusion changed.

## Stage Report: implementation

- DONE: The deletion commit ALREADY EXISTS at b77bfc942 on your branch - verify it satisfies the gated scope (-506 lines, 0 insertions, 4 files; wantGaps gone; three install literals gone; shell harness and driver gone) rather than redoing it; amend or add commits only if verification finds a gap
  Verified, not redone. `git diff --numstat 4d1912a69..HEAD` reports exactly the 4 gated files at +5/-510, net **-505** — inside the ±40 band. The 5 insertions are the doc-comment rewrite the design called for in deletion 2, not new logic; the gated "0 insertions" headline did not price it in, the tolerance paragraph did. `wantGaps`, its consuming loop, and the then-unused `reflect` import are gone; the three install literals are gone with `uname -s` / `go build -o spacedock ./cmd/spacedock` / `unsupported OS` retained; both harness files are absent from disk. Grep for all 11 deleted symbols plus `wantGaps` and `version_gate_flow` across `*.go`/`*.sh`/`*.md` returns nothing outside the state checkout. No amend was needed.
- DONE: Run the outstanding proofs: gofmt -l clean; contractlint and skills/integration green plain and under -race; AC-1 both sides (one-file XFAIL retire fails at merge-base, passes on your branch); mutation check reports line 401
  `gofmt -l ./internal ./cmd`: no output. `go test -count=1 ./internal/contractlint/ ./skills/integration/` ok (3.9s / 23.2s) and the same under `-race` ok (28.0s / 21.4s); these two packages carry every changed file.
  AC-1 negative side, at merge-base 4d1912a69: dropping `liveXFail("pi", "x02375wsg6q61xek7p0t36j2")` from `TestLiveCommonKeepMovingPosture` and changing no other file fails `TestRuntimeLiveRegistryReconciliation` at `live_registry_reconciliation_test.go:64` (`keep-moving-posture gaps = nil, want [...]`), and that is the only assertion that fires — the `:106` lines in the same output are `t.Logf`, confirmed by reading the source, not errors. AC-1 positive side: the byte-identical one-file edit on b77bfc942 gives `ok internal/contractlint 3.001s`. Leaving `wantGaps` in the tree is exactly what would keep the positive side red.
  AC-4 mutations, each run on the branch and each restored with `git checkout --` (no `git stash`, per the ideation incident): drifting the FO prose curl URL to `main/installer.sh` fails `TestInstallHintNoDrift` at `install_hint_drift_test.go:73`; downgrading the one-line brew form to bare `brew install spacedock` fails it at `:94`; setting the XFAIL owner to `"bad-owner"` fails the reconciliation at `live_registry_reconciliation_test.go:401` with `TestLiveCommonKeepMovingPosture: malformed liveXFail binding`. Line 401 and not 418 is itself the check that the 17-line removal is present in the tree under test.
  Keep-boundary: `TestDispatchAckMachineryIsAbsent`, `TestInstallHintNoDrift`, `TestRuntimeLiveCommonSuiteTimeouts`, `TestRuntimeLiveCommonFailFastPolicy`, `TestRuntimeLiveGapBindingValidation`, `TestVersionGateProseOSAwareHint` — all six PASS under `-v -run`.
- DONE: Write the implementation stage report to the entity accounting for every item, then send the completion signal
  This section; committed path-scoped to the state checkout and pushed to `spacedock-state/dev`.
- SKIPPED: AC-3's `go test ./... -race` arm (full-suite race, entity AC not in the dispatch's outstanding list)
  Every changed file lives in `internal/contractlint` and `skills/integration`, both green under `-race`; for every other package the test binary is byte-identical to merge-base, so a full `-race ./...` re-runs main, not this change. Plain `internal/cli` already needs 744s on this contended box, so `-race ./...` would yield timeout noise rather than signal. Named for the validator to overrule if it disagrees.

### Summary

The gated deletion at b77bfc942 holds under verification and needed no amendment: 4 files, +5/-510, net -505 against a declared -506 ±40, with the 5 insertions accounted for as the anticipated doc-comment rewrite. AC-1 is proven on a baseline that moves the wrong way — the identical one-file XFAIL retire fails at merge-base at line 64 and passes on the branch — and all three AC-4 mutations fire, mutation C at line 401 exactly as the corrected ideation evidence predicted.

One AC-3 finding the gate should see rather than a claim of green. Full `go test -count=1 ./...` on the branch produced, besides the already-declared `TestCodexResolveManifestAgainstInstalledHost` sandbox failure, a 10-minute package timeout in `internal/cli` and `internal/ensigncycle`. A control run at merge-base produced the same two timeouts, and re-running both packages serially on the branch with `-timeout 40m` gives `internal/ensigncycle` ok (842s) and `internal/cli` failing only on `TestCodexResolveManifestAgainstInstalledHost` (744s). The timeouts are the default per-package 10m limit under fleet contention on a machine running many concurrent ensigns, not a regression: neither package contains a file this branch touches. AC-3's "no worse than baseline" holds. My first baseline run was self-contaminated — I left the AC-1 edit applied in the control worktree — which is why it also showed a reconciliation failure and three flaky Codex timing tests; I restored the worktree and re-derived the comparison above rather than reporting the contaminated run.

## Stage Report: validation

- DONE: Independently re-exercise the ACs against worktree commit b77bfc942, never by reading the report: numstat +5/-510 over exactly 4 files, the 11 deleted symbols grepped absent, AC-1 both sides (one-file XFAIL retire fails at merge-base line 64, passes on branch), all three AC-4 mutations fire including line 401
  All mutations ran in detached throwaway worktrees at 4d1912a69 and b77bfc942, never the implementation worktree. `git diff --numstat 4d1912a69..b77bfc942` = exactly the 4 gated files, +5/-510 (AC-2: net -505 inside -506 ±40, 4 files inside ±1). Word-boundary grep over *.go/*.sh/*.md: all 11 harness symbols plus `wantGaps` and `version_gate_flow` absent (two substring hits were pre-existing unrelated identifiers, `writeExecutable` in internal/dispatch and a `sentinelPath` param in internal/status, both in files untouched by the branch). AC-1: byte-identical one-file retire of the keep-moving-posture pi XFAIL (same resulting blob f4f4a3a92 both sides) fails at merge-base with exactly one failing test at live_registry_reconciliation_test.go:64 and passes on the branch (`ok internal/contractlint`). AC-4: curl drift fails install_hint_drift_test.go:73; brew one-line downgrade fails :94 (the :90 two-line check passes by substring, so :94 is the load-bearing assertion); bad-owner fails live_registry_reconciliation_test.go:401 `malformed liveXFail binding` — and the same mutation at merge-base fails at :418, so the 17-line shift proves the removal is in the tree under test.
- DONE: Six keep-boundary tests pass; contractlint and skills/integration green plain and -race; rule on the SKIPPED full-suite race arm - the implementer's argument is that every changed file lives in those two packages
  TestDispatchAckMachineryIsAbsent, TestInstallHintNoDrift, TestRuntimeLiveCommonSuiteTimeouts, TestRuntimeLiveCommonFailFastPolicy, TestRuntimeLiveGapBindingValidation, TestVersionGateProseOSAwareHint all PASS under -v -run; each still guards a distinct surface the deletions border (dispatch-ack chain, install-command relation, timeout strings, failfast policy, gap shape, surviving prose tokens). Both packages ok plain and under -race; gofmt -l ./cmd ./internal ./skills clean. Ruling on the skip: ACCEPTED. Premise verified, not trusted — numstat confines the diff to those two packages, an all-file sweep (CI configs included) finds no external reference to any deleted artifact, and no non-test file changed, so a full `-race ./...` compiles byte-identical test binaries for every other package and re-tests main, not this change.
- DONE: AC-3 differential against clean 4d1912a69 with -timeout 40m serial for internal/cli and internal/ensigncycle; the Codex config failure and contention timeouts are documented environmental - reproduce before attributing; verdict PASSED or REJECTED with per-AC citations
  Branch full suite (all packages minus the two, then ensigncycle and cli serial with -timeout 40m): sole failure is TestCodexResolveManifestAgainstInstalledHost, reproduced byte-for-byte at clean 4d1912a69 ("Failed to read config file ~/.codex/config.toml: Operation not permitted") before attributing — environmental sandbox condition, present at baseline. The implementer's contention timeouts did not reproduce on the quieter box (ensigncycle ok 214s, cli 145s serial), consistent with contention, not code. No failure on the branch absent at baseline; AC-3 holds.

### Summary

PASSED — AC-1, AC-2, AC-3, AC-4 all verified with independently reproduced evidence; per-AC citations above. Adversarial pass found the diff is exactly the gated scope: deletion 1 removes only wantGaps, its loop, and the reflect import (duplicate-target rejection, parseLiveGap shape validation, TODO-owner join, and selector counts all survive and the gap-count lines at :106 are t.Logf, not assertions); deletion 2's rewritten doc comment now correctly defers install-command pinning to TestInstallHintNoDrift; deletion 1 makes docs/runtime-live-ci-registry.md's "does not use a copied gap oracle" sentence true. Reviewer findings: clean audit — no material findings, no new deferred risks; the harness-behavior residual (no executable successor for the sentinel-bound/ordering assertions) was priced into the approved ideation, not a new finding. Main has not moved on the four target files, so the dav9 overlap needs no reconciliation at this landing. Delivery can proceed.
