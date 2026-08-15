---
id: zvk9cnew2ggpaqb3wty24xtf
title: Remove redundant lint mirrors and the version-gate shell harness
status: ideation
source: "0.27 cut audit (2026-08-14), adversarially verified; captain directed filing"
started: 2026-08-15T02:55:48Z
completed:
verdict:
score:
worktree:
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
---

Three verified-redundant test mirrors.

1. The wantGaps XFAIL map (internal/contractlint/live_registry_reconciliation_test.go:51-66). A hand-copied oracle of the liveXFail() calls the same test already parses. The registry doc disclaims exactly this ("does not use a copied gap oracle"). Six lockstep two-file commits in range prove the churn. Keep parseLiveGap shape validation, duplicate-target rejection, and the TODO-owner join.
2. Three install-command literals in version_gate_smoke_test.go:37-39. TestInstallHintNoDrift already anchors the same tokens to their producer, docs/site/get-started/install.md. Keep the uname, go-build, and unsupported-OS tokens (lines 40-42) - they have no other pin.
3. skills/integration/testdata/version_gate_flow.sh plus its 351-line driver version_gate_fixture_test.go. The harness verifies a hand-written re-derivation of the FO prose, not shipped bytes. Its brew command already drifted from the lint-enforced prose and everything stayed green. Nothing binds mirror to prose.

Coordination: this supersedes scope items 2-3 of remove-startup-capability-probe (dav9), which rewrite the files this entity deletes. dav9 items 1 and 4 (shared-core prose, install.md) are unaffected. Whichever lands second reports the overlap as already done.

## Problem

Three test mirrors hand-copy something the tree already derives or already pins, and nothing binds copy to original. Each was confirmed against HEAD 4d1912a69.

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
Measured against a baseline that moves the wrong way: removing a single `liveXFail(...)` from `internal/ensigncycle/shared_live_runner_test.go` and changing *no other file* must leave `go test ./internal/contractlint/` passing. At HEAD 4d1912a69 the identical one-file edit FAILS at `live_registry_reconciliation_test.go:64`. Two-file edit becomes one-file edit; an incomplete removal keeps it failing.

**AC-2 (value) — The tree gets smaller.**
Cumulative line delta against origin/main is negative, within the tolerance band above.

**AC-3 — The suite is no worse than its baseline.**
`go test ./...` and `go test ./... -race` report no failure that is not also present at HEAD on the same machine. Stated against baseline rather than absolute green because HEAD is **not** currently green: `TestCodexResolveManifestAgainstInstalledHost` (internal/cli) fails at clean HEAD reading `~/.codex/config.toml` ("Operation not permitted"), a local sandbox-permission condition unrelated to this change. Implementation must not chase it.

**AC-4 — The surviving lints still catch what the deleted mirrors caught.**
Falsifiable by mutation, not by grep: drifting the curl token or the brew one-line form in the FO prose fails `TestInstallHintNoDrift`; a malformed `liveXFail` owner fails the reconciliation through `parseLiveGap`. The keep-boundary tests remain present and passing.

## Test plan

Deletion plus the existing suite; the surviving drift and shape lints are the regression floor. All of the following were exercised in a throwaway worktree at HEAD 4d1912a69 and seed implementation's first checks.

- **AC-1 baseline and result.** Retiring the `keep-moving-posture` pi XFAIL as a one-file edit against HEAD fails with `keep-moving-posture gaps = ... want ...` at `live_registry_reconciliation_test.go:64` — and that line is the *only* failure, so `wantGaps` is the sole thing making it a two-file edit. The same one-file edit with the removal applied passes (`ok internal/contractlint 1.601s`).
- **AC-4 mutation checks**, each confirmed to fail after the removal, so none is a tautology: replacing the curl URL in the FO prose fails `TestInstallHintNoDrift` at `install_hint_drift_test.go:73`; downgrading the brew one-line form to bare `brew install spacedock` fails it at line 94; changing a `liveXFail` owner to `"bad-owner"` fails the reconciliation at line 418 with `malformed liveXFail binding`.
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
