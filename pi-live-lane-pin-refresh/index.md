---
title: Refresh the pi-live lane pins and substrate assertions for the Pi 1.0 family
status: validation
score: 0.75
source: "Captain directive, 2026-10-03: Pi 1.0 shipped; update the CI pin and the relevant Pi extensions."
id: mh698y3ht6ydmr6ethaw9hg9
gates:
    version: 1
    records:
        - id: gate:mh698y3ht6ydmr6ethaw9hg9:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:mh698y3ht6ydmr6ethaw9hg9-backlog-1
              briefing:
                id: briefing:mh698y3ht6ydmr6ethaw9hg9:backlog:attempt-1:revision-1
                digest: sha256:6041994faf9dbfafbc683fc9a5b47729d9b34825fab9707546d87835726e2b53
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:mh698y3ht6ydmr6ethaw9hg9:backlog:1
                briefing: briefing:mh698y3ht6ydmr6ethaw9hg9:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-10-03T04:04:23.251551Z"
                decision: approve
                reason: Covers validating the published Pi family and deriving substrate assertions from the installed package.
                conn:
                    quote: i already said dispatch to ideation, but don't present the ideation gate until staff review finishes
                    source: captain instruction, this session, 2026-10-03
              application:
                target-stage: ideation
                state: consumed
        - id: gate:mh698y3ht6ydmr6ethaw9hg9:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:mh698y3ht6ydmr6ethaw9hg9-ideation-1
              briefing:
                id: briefing:mh698y3ht6ydmr6ethaw9hg9:ideation:attempt-1:revision-1
                digest: sha256:67395f0bd45e31a6c2842ca0be396b50f674b3dce43940f7af9728518f69f6f8
                room-ref: '@review/ideation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:mh698y3ht6ydmr6ethaw9hg9:ideation:1
                briefing: briefing:mh698y3ht6ydmr6ethaw9hg9:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-10-03T04:58:35.738946Z"
                decision: approve
                reason: Covers refreshing the pi-live lane to the published Pi family and deriving substrate assertions from the installed package, as the fastest way to observe the real matrix.
                conn:
                    quote: just update the pi and run the ci in pr
                    source: captain instruction, this session, 2026-10-03
              application:
                target-stage: implementation
                state: consumed
        - id: gate:mh698y3ht6ydmr6ethaw9hg9:validation
          stage: validation
          attempts:
            - id: gate-attempt:mh698y3ht6ydmr6ethaw9hg9-validation-1
              briefing:
                id: briefing:mh698y3ht6ydmr6ethaw9hg9:validation:attempt-1:revision-1
                digest: sha256:d8ff7c2b4faf19191eff5360afd23d248d7bcd5c2742a21a41df0c57e697fd34
                room-ref: '@review/validation/briefing-1'
              withdrawal:
                by: agent:first-officer
                at: "2026-10-05T01:43:58.731238Z"
                reason: The prepared question bundled the validation approval with three merge actions. Corrected to ask about the validation only.
            - id: gate-attempt:mh698y3ht6ydmr6ethaw9hg9-validation-2
              briefing:
                id: briefing:mh698y3ht6ydmr6ethaw9hg9:validation:attempt-2:revision-1
                digest: sha256:0867dd2a8a0d09af5fc9f688dd6f919d0be239b4641b8ddeeac9dfb2519ad333
                room-ref: '@review/validation/briefing-2'
              resolution:
                type: Resolution
                id: resolution:spacedock:mh698y3ht6ydmr6ethaw9hg9:validation:2
                briefing: briefing:mh698y3ht6ydmr6ethaw9hg9:validation:attempt-2:revision-1
                by: person:captain
                at: "2026-10-05T02:32:10.426728Z"
                decision: approve
              application:
                target-stage: done
                state: superseded
            - id: gate-attempt:mh698y3ht6ydmr6ethaw9hg9-validation-3
              briefing:
                id: briefing:mh698y3ht6ydmr6ethaw9hg9:validation:attempt-3:revision-1
                digest: sha256:9b37cf018a05b95675c24bb56e851c685f700232d4d983046923a7045a1371fc
                room-ref: '@review/validation/briefing-3'
              resolution:
                type: Resolution
                id: resolution:spacedock:mh698y3ht6ydmr6ethaw9hg9:validation:3
                briefing: briefing:mh698y3ht6ydmr6ethaw9hg9:validation:attempt-3:revision-1
                by: person:captain
                at: "2026-10-05T16:32:28.119506Z"
                decision: approve
              application:
                target-stage: done
                state: pending
sprint: pi-ux
group: tooling
sprint-readiness: ready
started: 2026-10-03T04:04:56Z
worktree: .worktrees/spacedock-ensign-pi-live-lane-pin-refresh
pr:
---

The pi-live Runtime Live E2E lane validates a superseded Pi family. Refresh its three
version/integrity pairs and replace **all four** hardcoded TypeScript layout assertions
with checks of the installed package's declared runtime artifacts.

## Problem

At code baseline `cdfa462d1`, `.github/workflows/runtime-live-e2e.yml` installs
`@earendil-works/pi-coding-agent@0.85.1`, `pi-subagents@0.67.0`, and
`pi-intercom@0.13.0`. npm's published latest family, checked on 2026-10-03, is
`1.0.0` / `0.75.0` / `0.16.0`. The current lane therefore does not measure
compatibility with the family users install today.

There are FOUR stale assertions, not two (line numbers at that baseline):

| Step | Line | Hardcoded assertion |
|---|---:|---|
| Install Pi CLI and substrates | 669 | `test -f "$pi_npm_root/node_modules/pi-subagents/src/extension/index.ts"` |
| Install Pi CLI and substrates | 671 | `test -f "$pi_npm_root/node_modules/pi-subagents/src/intercom/intercom-bridge.ts"` |
| Verify Pi current-checkout setup | 777 | `test -f "$PI_SUBAGENTS_PACKAGE_ROOT/src/extension/index.ts"` |
| Verify Pi current-checkout setup | 779 | `test -f "$PI_SUBAGENTS_PACKAGE_ROOT/src/intercom/intercom-bridge.ts"` |

Both source paths are absent from pi-subagents 0.75.0. Its manifest instead declares
`pi.extensions = ["./index.js"]` and
`exports["./intercom-bridge"].default = "./src/api/intercom-bridge.js"`.
Changing only the pins leaves both setup checkpoints broken.

## Proposed approach

1. Update the three exact version/integrity pairs together using the published npm
   metadata below. Keep the existing `verified_pack` check, tarball installs,
   installed name/version assertions, and compatibility guard. Rewrite the obsolete
   `0.85.x` compatibility comment to identify the Pi 1.0 family while distinguishing
   it from the unchanged launcher floor `0.83.0`.
2. At BOTH setup checkpoints, replace the extension and bridge `.ts` assertions with
   a small inline Node check reading the installed pi-subagents `package.json`.
   Require a nonempty `pi.extensions` array and the runtime string at
   `exports["./intercom-bridge"].default`; resolve the entries relative to the package
   root and require each to be a regular file. Missing declarations or files fail
   setup. Check every declared extension, not just element zero. Do not substitute
   `.d.ts` entries or hardcode the observed `.js` destinations. The published family
   needs no generalized export-condition resolver. Retain the existing skill checks,
   package-root exports, doctor invocation, and pi-ai compatibility import.
3. Add a focused Go structural guard in `internal/release/pi_live_pins_guard_test.go`,
   reusing the existing workflow-job parser. Its independent oracle is the reviewed
   published version/integrity snapshot, not values read back from the same YAML.
   Bind checks to `pi-live` and its named executable setup steps; comments and a
   matching snippet in another job cannot satisfy them. Cover all three pin pairs,
   verified-pack wiring, and manifest-check wiring at both checkpoints with negative
   mutations. This is structural proof, not a claim that Go executed Node or Pi.
4. Refresh only the Pi local-install instructions in `docs/runtime-live-ci.md`.
   The concrete proposed doc diff is:

   ```diff
   -npm install -g @earendil-works/pi-coding-agent@0.80.10
   -npm install --prefix "$HOME/.pi/agent/npm" pi-subagents@0.35.1 pi-intercom@0.6.0
   +npm install -g @earendil-works/pi-coding-agent@1.0.0
   +npm install --prefix "$HOME/.pi/agent/npm" pi-subagents@0.75.0 pi-intercom@0.16.0
   ```

All mechanisms serve AC-1's value: the actual published family is installed AND
validated. A pin-only bump is insufficient because all four source assertions fail.
Replacing `.ts` with fixed `.js` paths is simpler but encodes another private layout;
the already-shipped manifest gives the package's public runtime entrypoints. One
inline check at each existing checkpoint avoids a new helper/script protocol. A
manual-only wiring review cannot catch subsequent pin or checkpoint regressions;
the small release guard is the approved deterministic proof owner. Do not build a
workflow replay harness or add another live scenario; the existing lane owns runtime
compatibility. Doc-version agreement is a one-off validation, not a prose-grep test.

## Risk evidence

Published metadata queried with `npm view <package>@<version> name version dist.integrity
pi exports engines --json`; `npm view <package> dist-tags.latest` independently returned
`1.0.0`, `0.75.0`, and `0.16.0` on 2026-10-03:

| Package | Version | Published integrity |
|---|---|---|
| `@earendil-works/pi-coding-agent` | `1.0.0` | `sha512-/FtbxoSQU/mEv1QnichJjRjqteqaIaMWxmhB4G367+MwZfX7/DI5B9YAg5lqbN7nztFskBEtUSZ+FlmMBECtMw==` |
| `pi-subagents` | `0.75.0` | `sha512-RO4DiTJM6pnK8y9PnD7Y6TLeiX2c8Kh6QhmkecFo+ou5qtnz5qiM+vUKru48UmJG1nB9eJ6rALd1M04uBOR4XQ==` |
| `pi-intercom` | `0.16.0` | `sha512-ClGQuovPsz7r1iQwMRjEN+8wxywfrDrMILAkCSf/z19Nezzyfz77U8362Wb/VFNoZOwmF4PBsuGWCtB9AsEJMQ==` |

Ideation exercised the new manifest mechanism against the REAL `npm pack
pi-subagents@0.75.0 --ignore-scripts --json` tarball, extracted in a temporary directory:

- Independent SHA-512 of tarball bytes equals both the registry snapshot and pack
  metadata. Resolving the manifest entries finds regular files `./index.js` and
  `./src/api/intercom-bridge.js`; both old source paths are absent.
- Temporarily renaming each resolved runtime file in turn makes the same Node check
  throw; restoring it passes. This seeds implementation's missing-target negative
  exercise without inventing a fixture package layout.
- A one-off extraction of the workflow's existing `verified_pack` function accepts
  the published subagents integrity and returns the tarball path. Supplying
  `sha512-corrupted-for-negative-probe` exits 1 with `integrity mismatch`, before any
  install. No committed shell replay harness is needed.

The seed also records a local Pi 1.0.0 doctor probe: the unchanged version floor
accepts 1.0.0 and other checks pass except the two stale doctor file probes. That is
inherited evidence, NOT a green live-lane run produced by ideation. Node's published
engine remains `>=22.19.0`, satisfied by this lane's existing Node 24 setting.

Dependency: `pi-doctor-probes-stale-subagents-layout` (`mc`) owns the doctor's stale
probe behavior in `internal/cli/pi.go`. Its fix must be present for this lane's
current-checkout doctor check and live validation to pass; do not bypass doctor or
absorb its fix here. A green pin install alone does not discharge AC-1. Any additional
Pi 1.0 runtime incompatibility is a surfaced blocker, not permission to shrink the
scenario set or weaken the lane. No further spike is needed before ideation review:
the only newly proposed resolution mechanism and existing integrity failure path
were exercised; end-to-end compatibility remains the implementation/validation proof.

## Expected surface and tolerance

Estimate net LOC change: +160, across 3 files (`.github/workflows/runtime-live-e2e.yml`,
`internal/release/pi_live_pins_guard_test.go`, and `docs/runtime-live-ci.md`).
Insertions ~+185, deletions ~-25. Tolerance: ±60 net LOC, ±1 file. This replaces the
seed's +45 estimate to account for four assertion replacements and falsifiable guard
coverage. Reuse the existing test parser rather than add another parsing framework.

Declared semantic changes: which Pi family pi-live installs and which installed
substrate artifacts its two setup checkpoints accept. No CLI grammar, stored format,
authority, launcher floor, scenario selection, failure policy, other-lane pins,
GitHub action majors, or Node setting changes are authorized.

## Acceptance criteria

**AC-1 (value) — pi-live installs and validates the Pi family it pins.**
The three installed package names/versions equal the pins the layer declares once in
`internal/pilive/pilive.go` — `@earendil-works/pi-coding-agent@1.0.2`,
`pi-subagents@0.75.0`, and `pi-intercom@0.16.0` — and the pi-live Runtime Live E2E lane
completes green on that candidate, with its existing common journeys and front-door smoke
still selected. Proof: the lane run at the frozen tip, which installs those pins and
supplies installed-version logs, doctor output, and durable journey/smoke results under
the existing grading policy. No new skips count.
**Guarantee lost by the shipped design (recorded, not claimed):** the registry-comparison
Go guard is gone, so drift between a pin and the published registry is caught only by the
lane's install step, not by a test that can fail on its own. The named falsifying edit
**revert-agent-pin** is retired with it. The criterion's original literal `1.0.0` was the
observed baseline at ideation; the family moved to `1.0.2` before implementation.

**AC-2 — the substrate assertions follow installed runtime declarations.**
Both setup checkpoints resolve every `pi.extensions` entry and the intercom bridge
runtime export from the installed manifest, fail for missing declarations or files, and
require neither old source path. Proof: the port's `verify-manifest` command and its tests,
and the live lane's green run at the frozen tip.
**Guarantee lost by the shipped design (recorded, not claimed):** the Go structural guard is
gone, and with it the **restore-setup-source-assertion** and **remove-bridge-target**
falsifiers. The remaining demonstration is a test authored by the same change, so it does not
independently bind the workflow's two checkpoint blocks.

**AC-3 — corrupted integrity prevents installation.**
All three integrity pins are verified before their tarballs are installed. Proof: the port's
install path and its tests, and the live lane's install of the pinned tarballs at the frozen
tip.
**Guarantee lost by the shipped design (recorded, not claimed):** the `verified_pack` wiring
and the **PI_SUBAGENTS_INTEGRITY** falsifier are gone. The comparison reads npm-reported
metadata against the pin rather than the tarball bytes independently, and its demonstration is
a test authored by the same change.

**AC-4 — the pins live in one place and the instructions agree with them.**
The three versions are declared once, in `internal/pilive/pilive.go`, and the local-install
documentation delegates to `spacedock-release print-install` instead of naming versions.
Proof: one-off comparison of that documentation against the declared pins, recorded during
validation.
**Guarantee lost by the shipped design (recorded, not claimed):** the workflow compatibility
comment no longer names the Pi 1.0 family or explains the separate `0.83.0` floor, and the named
**restore-stale-doc-version** falsifier is gone. The floor survives as a code value only, and no
standing check catches documentation that drifts from the pins.

**AC-5 (no regression) — the bounded change preserves existing non-Pi behavior.**
`go test ./internal/release/...`, `go test ./...`, `go test ./... -race`, and
`go build ./...` pass; the reviewed diff leaves all other lanes' CLI pins/action
majors and every live selector unchanged, and raises the Pi readiness floor to the Pi 1.0 family by the captain's direction. Proof: existing release tests
plus a candidate-versus-base diff audit of these boundaries. Named falsifying edit
**downgrade-other-lane-checkout**: change claude-live's `actions/checkout@v5` to `@v4`;
`TestNode24ActionsPinnedAtMinimum` and the boundary audit fail. Any unauthorized
other-lane version change must also fail the diff audit, even if still above a floor.

### Criteria amendment — captain-approved, 2026-10-04

The validation returned AC-1 to AC-4 as **not met**, because each named an artifact the
captain's simplification rounds removed: the registry-comparison Go guard, the Go structural
guard with its falsifiers, the `verified_pack` wiring with `PI_SUBAGENTS_INTEGRITY`, and the
version-naming workflow comment and documentation lines. The work was not at fault; those
criteria were written against a design that no longer ships.

The captain directed the amendment. Each criterion above is restated against the shipped
mechanism, and each records the guarantee it loses. Per the science officer's ruling, a
deliberately removed guarantee is recorded rather than described as equivalent. This layer
therefore claims the pinned family, the manifest-resolved checks, the integrity comparison,
and the single pin source. It does not claim a standalone pin-revert falsifier, an independent
structural guard, an independent integrity falsifier, or a standing doc-versus-pin check.

AC-5 to AC-10 were met as written. Of those, AC-2, AC-3, AC-7, AC-8 and AC-10 rest on
demonstrations authored by the same change, and AC-6 is the only criterion holding an
independent live artifact.

Two lane defects found by the same validation are recorded and owned elsewhere: the
`default-headless-gate-stop` binding engaged on `implementation-worker-not-dispatched`, which its
own comment says it does not cover, and two bindings (`smallest-sufficient-mechanism`,
`keep-moving-posture`) are stale because their journeys now pass.

### Folded scope — isolated-home Pi discovery (M1)

This task also carries the M1 isolated-home discovery fold: one isolated-home
setup contract (`seedPiIsolatedHome`/`seedPiDefaultExtensions`), a non-live
helper seam, both substrate npm registrations in the isolated home's settings,
the Spacedock package as one absolute checkout path, independent explicit
overrides, and no substrate extension path supplied by the harness in default
discovery. Keep the negative control: both root symlinks present, both npm
registrations removed, both substrate tools disappear.

**AC-6 (VALUE) — an isolated-home Pi run with BOTH `PI_SUBAGENTS_PACKAGE_ROOT` and `PI_INTERCOM_PACKAGE_ROOT` unset resolves both required extensions through Pi's own package discovery, with no substrate extension path supplied by the harness, and with the Spacedock package and ensign skill still loaded.** A successful explicit fallback is NOT evidence for this criterion. Proof: the authorized `TestLivePiFrontDoorSmoke` run with both variables unset, whose `assertPiEnsignBootContract` grading still requires the Spacedock/ensign boot contract plus both substrate tools. Named falsifying edit **remove-registration-seeding**: delete the isolated settings' npm registrations (or the seeded symlinks); the run loses the substrate tools.

**AC-7 — discovery reads the real installed location, never a hard-coded path.** `settings.json` `packages` entries resolve through the agentDir npm path (`<agentDir>/npm/node_modules/<name>`), and local sources match by `package.json` name. The isolated home's settings must register BOTH `npm:pi-subagents` and `npm:pi-intercom`, and the Spacedock package must stay exactly one absolute checkout-path entry, never a `file:` entry. Proof: deterministic `internal/ensigncycle` tests (`TestPiDefaultExtensionRoots…`, `TestPiIsolatedHomeRegistersBothSubstratesAndAbsoluteSpacedock`). Named falsifying edits **point-settings-elsewhere**: point settings at another directory and require that directory; **write-file-prefix**: write `file:`+path and require zero extensions loaded.

**AC-8 — each explicit override still wins independently when set, and the other package still discovers normally.** In default mode both package-root variables are scrubbed from the child env and are not re-added, not even as empty assignments. Proof: `TestPiLiveEnvHonorsIndependentOverrides` and `TestPiLiveEnvDefaultScrubsPackageRoots`, plus the retained `TestPiLiveEnvDropsForeignRuntimeMarkers`. Named falsifying edits **set-one-override**: set one override and require its precedence while the other variable stays absent; **re-add-empty**: re-add either variable as an empty assignment in default mode.

**AC-9 — ordinary non-live helper tests plus gofmt, live-tagged vet and build, and the authorized front-door smoke pass with both variables unset, with durable entity report and commit evidence.** The helper seam must NOT carry a live build constraint, because a live-tagged definition is invisible to ordinary Go tests. Proof: `internal/ensigncycle/pi_default_extensions_test.go` has no build constraint; `gofmt -w ./cmd ./internal`, `go vet -tags live ./internal/ensigncycle`, and `go build -tags live ./internal/ensigncycle` pass; the smoke run carries durable entity report/commit evidence. Named falsifying edit **live-tag-the-helper**: define the helper under the live tag and watch the ordinary tests fail to see it.

**AC-10 (no regression) — the offline `go test ./...` gate is deterministic and tests only our decision logic, not the machine.**
(a) The reset-on-activity decision must not be tested by racing real sleeps.
`TestCodexProcessActivityResetsQuietBudget` (`internal/ensigncycle/codex_single_run_test.go`)
must supply the activity through a test-controlled signal or an injected deadline so it
asserts the decision logic, and must not decide on a real-time window: remove the
`const quietBudget = 250 * time.Millisecond` wall-clock race and the
`result.duration > 4*quietBudget` machine-speed assertion. The kill-on-silence sibling
`TestCodexProcessQuietTimeoutPreservesFaultEvidence` (mode `stall`, silent by
construction) stays as the deterministic offline test of the kill path; the
real-pacing property against the real watchdog remains covered by the live lane and is
not re-asserted offline.
(b) `TestCodexResolveManifestAgainstInstalledHost` and `TestVersionAmbiguousMarkersExitZero`
(`internal/cli`) and `TestSurveyCodexPresenceThroughSync` (`skills/integration`) must build
their own isolated host home — using the existing `internal/ensigncycle/codex_liveenv.go`
machinery (isolated `CODEX_HOME`, scrubbed parent `CODEX_HOME`, `HOME`,
`PI_CODING_AGENT_DIR`) or an equivalent fake host — instead of reading the operator's
machine.
Proof: the offline suite passes in the CI offline job on a loaded runner, and the same
tests pass locally with an operator `CODEX_HOME`/`PI_CODING_AGENT_DIR` set. Named
falsifying edits **race-real-sleeps**: reintroduce the real-sleep reset test;
**inherit-operator-home**: read the operator's `CODEX_HOME`/`PI_CODING_AGENT_DIR` again
and require the local plugin cache.

## Out of scope

Doctor resolution (`mc`), the harness Node-runtime task (`3g8`), claude-live/codex-live
CLI pins, the launcher `piVersionFloor`, substrate implementation, new live scenarios,
workflow replay infrastructure, and generic package export resolution. Do not fix
unrelated stale runtime-support examples in this slice.

## Test plan

Primary proof owner: Go structural guard under `internal/release/`, following
`node24_actions_guard_test.go` and `workflow_exec_guard_test.go`. Reuse
`parseWorkflowJobs` to inspect the owning job and named setup steps. This follows the
release-machinery policy: Go tests for decision logic and structural YAML wiring,
not workflow-shell replay against fixture repositories. Add the focused tests before
editing workflow pins.

- Deterministic guard (AC-1/2/3): new `TestPiLivePinsAndSubstrateAssertions` and a
  table-driven mutation companion. Reject each old version/hash, a removed or
  bypassed verified-pack call, a removed checkpoint check, and a restored hardcoded
  source assertion at EACH of lines 669/671/777/779. Mutating only the later setup
  checkpoint must fail. Expected cost: roughly 100–160 test LOC, seconds, no npm,
  network, Node, or credentials required by Go unit tests.
- One-off package exercise (AC-2/3): execute the candidate manifest checks against
  the packed published family, then hide each target and remove each required
  manifest declaration in turn; record nonzero exits. Exercise the actual
  `verified_pack` function with good/bad expected integrity. A few minutes plus npm
  network; temporary directories only. Distinct falsifiers are absent declarations,
  absent files, and a disabled integrity mismatch exit, not another text-match test.
- One-off doc/base audit (AC-4/5): record version agreement and the protected-boundary
  diff. `restore-stale-doc-version` and `downgrade-other-lane-checkout` name concrete
  failures. Minutes, no new persistent checker.
- Existing deterministic suites/build and registry reconciliation:
  `go test ./internal/release/...`, `go test ./...`, `go test ./... -race`,
  `go build ./...`, `gofmt -w ./cmd ./internal`, and
  `go test ./internal/contractlint -run '^TestRuntimeLiveRegistryReconciliation$'`.
- Live proof (AC-1): run the unchanged pi-live lane via `live:pi` or
  `workflow_dispatch live_cadence=pi`, after `mc` is available. Record candidate SHA,
  run URL, three installed versions, setup success, common-journey grades and
  `TestLivePiFrontDoorSmoke` durable results. Requires existing CI-E2E-PI approval
  and auth; budget tens of minutes at current lane caps. Existing XFAIL policy is
  unchanged; installation success alone or a skipped runtime suite is not green proof.
- Before merge, obtain the required detached adversarial audit for CI machinery.
  It tests the above claim-breaking edits on a throwaway checkout and checks that
  the advertised owner actually fails. Review-only; not another standing mechanism.

## Stage Report: ideation

- DONE: Correct the seed's surface count: the pi-live lane hardcodes FOUR TypeScript layout assertions (lines 669, 671, 777, and 779), not two, and the body must state all four.
  The Problem table identifies all four executable assertions at code baseline `cdfa462d1`; both named setup checkpoints are in scope.
- DONE: Task body states the problem, the chosen approach, and criteria under the exact `## Acceptance criteria` heading, with every AC paired to a named falsifying edit, and a value AC measuring that the lane installs and validates the published Pi family.
  AC-1 measures installed versions and existing live results against an independent npm snapshot; AC-1 through AC-5 each name a falsifying edit.
- DONE: The test plan names a Go structural guard under internal/release/ as the proof owner, consistent with the repository's release-machinery proof posture, and the expected surface declares net LOC change, files, and tolerance.
  `pi_live_pins_guard_test.go` owns structural proof; proposed baseline is +160 net LOC (+185/-25), 3 files, tolerance ±60 net LOC/±1 file; no replay harness.
- DONE: Exercise the risky manifest and integrity mechanisms before proposing them.
  Real pi-subagents 0.75.0 tarball resolves both runtime files; hiding either fails; the actual verified_pack function accepts the published hash and exits 1 on a corrupted hash.
- DONE: Propose the documentation change and preserve stage ownership.
  Body contains the exact two-line install-command diff, leaves frontmatter byte-identical, and keeps the doctor's `mc` fix and live green proof explicitly separate.
- DONE: Run focused baseline validation.
  `go test ./internal/release/...`, registry reconciliation, and `go build ./...` pass; downgraded checkout wiring is the existing guard's concrete falsifier, not a new test claimed shipped here.
- FAILED: Repository-wide baseline `go test ./...` and `go test ./... -race` validation.
  Baseline fails Codex plugin resolution, inherited PI_CODING_AGENT ambiguity, and survey Codex cwd expectations; ensigncycle times out at 10m. Race repeats Codex/survey failures and times out CLI/ensigncycle at 10m; no product-code fix is included.
- DONE: AC-1 proof plan: independent npm version/hash snapshot plus Go pin guard and unchanged green pi-live run; revert-agent-pin must fail (live proof remains pending implementation).
- DONE: AC-2 proof plan/evidence: guard both setup checkpoints and exercise actual manifest files; real-tarball missing-target probes fail; restore-setup-source-assertion and remove-bridge-target falsify it.
- DONE: AC-3 proof plan/evidence: guard all verified_pack call sites; actual function exits 1 for corrupt-subagents-integrity and accepts the published hash.
- DONE: AC-4 proof plan: one-off doc-command/comment comparison to the independent registry snapshot; restore-stale-doc-version must produce a mismatch.
- DONE: AC-5 proof plan/evidence: baseline suites/build and protected-boundary audit; focused release/registry/build pass, broad baseline failures noted above; downgrade-other-lane-checkout must fail the existing guard.
- SKIPPED: Implement the pin refresh and obtain a green live-lane acceptance run.
  This dispatch is ideation only; implementation follows the reviewed design and live acceptance requires the separate doctor fix plus CI environment approval.

### Summary

Fleshed out a bounded three-file Pi 1.0 pin refresh, covering all four stale assertions and proving manifest resolution against the published tarball. The proposal includes independent version/integrity evidence, falsifiable acceptance criteria, a concrete documentation diff, and a larger explicit LOC estimate for the guard; staff review and the ideation gate remain with the first officer.

### AC scan evidence

Command: `/opt/homebrew/Caskroom/spacedock@next/0.28.0-pre3/spacedock status --read /Users/clkao/git/spacedock-research/spacedock-v1/docs/dev/.spacedock-state/pi-live-lane-pin-refresh/index.md --ac-scan --json --workflow-dir docs/dev`. All five criteria have citations.

```json
{"command":"read","stage":"ideation","acs":[{"id":"AC-1","line":"162","unevidenced":"false","citations":[{"line":"262","text":"  AC-1 measures installed versions and existing live results against an independent npm snapshot; AC-1 through AC-5 each name a falsifying edit."},{"line":"262","text":"  AC-1 measures installed versions and existing live results against an independent npm snapshot; AC-1 through AC-5 each name a falsifying edit."},{"line":"273","text":"- DONE: AC-1 proof plan: independent npm version/hash snapshot plus Go pin guard and unchanged green pi-live run; revert-agent-pin must fail (live proof remains pending implementation)."}]},{"id":"AC-2","line":"173","unevidenced":"false","citations":[{"line":"274","text":"- DONE: AC-2 proof plan/evidence: guard both setup checkpoints and exercise actual manifest files; real-tarball missing-target probes fail; restore-setup-source-assertion and remove-bridge-target falsify it."}]},{"id":"AC-3","line":"184","unevidenced":"false","citations":[{"line":"275","text":"- DONE: AC-3 proof plan/evidence: guard all verified_pack call sites; actual function exits 1 for corrupt-subagents-integrity and accepts the published hash."}]},{"id":"AC-4","line":"193","unevidenced":"false","citations":[{"line":"276","text":"- DONE: AC-4 proof plan: one-off doc-command/comment comparison to the independent registry snapshot; restore-stale-doc-version must produce a mismatch."}]},{"id":"AC-5","line":"203","unevidenced":"false","citations":[{"line":"262","text":"  AC-1 measures installed versions and existing live results against an independent npm snapshot; AC-1 through AC-5 each name a falsifying edit."},{"line":"277","text":"- DONE: AC-5 proof plan/evidence: baseline suites/build and protected-boundary audit; focused release/registry/build pass, broad baseline failures noted above; downgrade-other-lane-checkout must fail the existing guard."}]}]}
```


## Stage Report: implementation

- DONE: The three pins name the published family with verified sha512 integrity values, and the lane installs that family.
  Code commit `ed19357ef` (branch `spacedock-ensign/pi-live-lane-pin-refresh`): the extracted install step succeeded in an isolated local HOME/npm prefix on macOS, Node 24.13.1/npm 11.8.0; installed-name/version checks logged 1.0.0 / 0.75.0 / 0.16.0. This is not a green CI claim (AC-1).
- DONE: All four hardcoded TypeScript assertions at lines 669, 671, 777, and 779 are replaced by assertions derived from the installed package's own declarations, and no TypeScript source path remains.
  Both executable checkpoint blocks read installed package.json; no old source assertions remain in the workflow. Their exact replacement checks are enumerated below (AC-2).
- DONE: The workflow comments and docs/runtime-live-ci.md state the pinned family consistently, and the repository's release-machinery checks pass.
  `go test ./internal/release/...` passed; one-off doc/comment comparison agrees with the independent npm snapshot, retains the distinct 0.83.0 floor, and rejects restore-stale-doc-version (AC-4).
- DONE: Add deterministic tests before editing pins, with an independent published snapshot and negative wiring mutations.
  New guard first failed on the old workflow, then passed with 34 mutation cases covering all old versions/hashes, corrupt hashes, commented pins, removed/bypassed pack calls, tarball-install bypasses, wrong owning job, both checks' declarations/file validation, and all four restored source assertions (AC-1/2/3).
- DONE: Exercise both actual Node bodies against the real packed and installed pi-subagents 0.75.0 package.
  Both passed; each exited 1 for either hidden runtime target, either target replaced by a directory, missing/empty extensions, missing bridge default, types-only bridge, and a missing second extension; restored packages passed. These fail if existence/declaration enforcement or all-extension iteration is removed (AC-2).
- DONE: Exercise the actual verified_pack integrity failure path.
  Extracted workflow function accepted the published subagents hash (exit 0), rejected `sha512-corrupted-for-negative-probe` (exit 1 before installation); disabling its mismatch exit in the temporary probe returned 0 and falsified the bad-hash expectation. Independent SHA-512 of all three packed tarballs matched registry values (AC-3).
- DONE: Preserve compatibility guard, doctor invocation, launcher, other lanes, journeys, XFAILs, grading, and Node/action settings.
  Extracted unchanged compatibility step passed its real pi-ai 1.0.0 /compat dynamic import; protected-boundary comparison against base `25a67d219` found all other steps/jobs byte-identical and separately confirmed doctor invocation unchanged; a checkout downgrade failed that audit (AC-5).
- DONE: Run focused validation and formatting/build checks.
  `go test ./internal/contractlint -run '^TestRuntimeLiveRegistryReconciliation$' -count=1`, `go build ./...`, `gofmt -w ./cmd ./internal`, and `git diff --check` passed; unrelated pre-existing formatting was restored to keep exactly three changed code files (+200/-15).
- FAILED: Repository-wide `go test ./...` and `go test ./... -race` (both exit 1).
  Both reproduced known baseline failures: TestCodexResolveManifestAgainstInstalledHost (local plugin cache), TestVersionAmbiguousMarkersExitZero (ambient PI_CODING_AGENT), TestSurveyCodexPresenceThroughSync (blank_cwd=0), and internal/ensigncycle package timeout at 10m. No unrelated fixes were made (AC-5 remains incomplete).
- FAILED: Current-checkout doctor against the installed published family.
  Strict-shell extracted setup probe exited 1: only the unchanged doctor's extension/bridge source probes were MISSING; CLI 1.0.0, auth, skills, intercom root, and floor checks passed. The separate `mc` dependency owns this fix; nothing was bypassed (AC-1).
- SKIPPED: Green pi-live CI run and detached adversarial reviewer audit.
  No CI run URL/SHA or journey/smoke grades are claimed; full live proof awaits the separate doctor fix and CI approval, and independent review remains required before merge (AC-1).

### Registry commands and returned integrity values

- `npm view @earendil-works/pi-coding-agent@1.0.0 name version dist.integrity --json` → `sha512-/FtbxoSQU/mEv1QnichJjRjqteqaIaMWxmhB4G367+MwZfX7/DI5B9YAg5lqbN7nztFskBEtUSZ+FlmMBECtMw==`
- `npm view pi-subagents@0.75.0 name version dist.integrity --json` → `sha512-RO4DiTJM6pnK8y9PnD7Y6TLeiX2c8Kh6QhmkecFo+ou5qtnz5qiM+vUKru48UmJG1nB9eJ6rALd1M04uBOR4XQ==`
- `npm view pi-intercom@0.16.0 name version dist.integrity --json` → `sha512-ClGQuovPsz7r1iQwMRjEN+8wxywfrDrMILAkCSf/z19Nezzyfz77U8362Wb/VFNoZOwmF4PBsuGWCtB9AsEJMQ==`

### Exact replacement assertions

- Install, former 669: read `$pi_npm_root/node_modules/pi-subagents/package.json`; require nonempty `pi.extensions`; every entry must be a nonempty string resolving relative to that package root to a regular file.
- Install, former 671: read that same manifest's `exports["./intercom-bridge"].default`; require a nonempty runtime string resolving relative to that root to a regular file (not the types declaration).
- Current-checkout setup, former 777: read `$PI_SUBAGENTS_PACKAGE_ROOT/package.json`; require nonempty `pi.extensions`; every entry must be a nonempty string resolving relative to that root to a regular file.
- Current-checkout setup, former 779: read that same manifest's `exports["./intercom-bridge"].default`; require a nonempty runtime string resolving relative to that root to a regular file (not the types declaration).

### Summary

Committed the bounded three-file Pi 1.0 refresh, with independently registry-verified pins, manifest-driven assertions at both checkpoints, and a 34-negative-case structural guard. Local install, manifest/integrity experiments, compatibility import, release checks, reconciliation, and build pass; full suites reproduce baseline failures and live acceptance remains blocked by the unchanged doctor probes. Local raw validation logs are under `/tmp/pi-live-pin-refresh.TRTQoZ/`; the copied temporary auth file was removed, and no launcher/journey/grading changes are included.

## Stage Report: implementation (cycle 2)

- DONE: Implement the folded isolated-home setup contract for the required pi-subagents and pi-intercom extensions: one setup contract, a non-live helper seam, both npm registrations in the isolated home's settings, the Spacedock package as one absolute checkout path, independent explicit overrides, and no substrate extension path in default discovery. Keep the negative control from the scope note.
  Code commit `631b482e6` (rebased onto the merge-augmented remote tip): `seedPiIsolatedHome`/`seedPiDefaultExtensions` in non-live `pi_default_extensions_test.go`; isolated settings = `["npm:pi-subagents","npm:pi-intercom",repo]` (never `file:`); both live fixtures call it; `piLiveEnv` scrubs both root vars by default and forwards only nonempty explicit overrides. The removed `.ts` source assertion was replaced by root discovery, not restored.
- DONE: AC-7 (discovery reads the real installed location) — deterministic proof.
  `TestPiDefaultExtensionRootsReadsRealInstalledLocation` points settings at a custom agentDir and requires that agentDir's npm path; `TestPiDefaultExtensionRootsMatchesLocalSourceByPackageName` requires the package.json-named local source and rejects a same-named decoy; falsifiers point-settings-elsewhere and write-file-prefix.
- DONE: AC-8 (independent explicit overrides) — deterministic proof.
  `TestPiLiveEnvHonorsIndependentOverrides` covers subagents-only/intercom-only/both; `TestPiLiveEnvDefaultScrubsPackageRoots` requires neither variable present (not even empty) in default mode; falsifiers set-one-override and re-add-empty.
- DONE: AC-6 negative control + deterministic preconditions.
  `TestPiIsolatedHomeNegativeControlDropsSubstrateRegistrations`: both root symlinks present, both npm registrations removed, both substrates disappear from discovery. A no-model real Pi 1.0.0 loader check on the exact isolated layout (Node 24.13.1, both root vars deleted, no extra extension paths, parent `PI_SUBAGENT_*` scrubbed) loaded 3 extensions (`spacedock.ts`, `pi-subagents/index.js`, `pi-intercom/index.ts`), tools `subagent`/`bg_wait`/`subagents_enable`/`intercom`, and the `ensign` skill with zero loader errors; removing both npm registrations left only the Spacedock extension and zero substrate tools (AC-6/AC-7).
- DONE: Record the folded scope and the four criteria in the task body, and cite each criterion in the stage report. Record the disposition of every live Pi binding (cleared, kept, or re-anchored) with the evidence that decided it.
  Body `### Folded scope — isolated-home Pi discovery (M1)` adds AC-6..AC-9 under the existing numbering, each with its own falsifying edit; all four are cited above and in the binding section below.
- DONE: AC-1 through AC-5 carried forward unchanged from the implementation cycle-1 report.
  AC-1 lane pins, AC-2 manifest-derived assertions at both setup checkpoints, AC-3 verified integrity pins, AC-4 doc/comment agreement, and AC-5 bounded non-Pi behavior remain as committed (`094e06ff7`, rebased as `818b9c4ad`); this cycle neither narrows nor renumbers them.
- SKIPPED: AC-6/AC-9 authorized front-door model smoke (`TestLivePiFrontDoorSmoke`) with both package-root variables unset.
  Not run this cycle: it is a live, model-auth, multi-minute run outside the FO's bounded-proof budget. The default-discovery mechanism is proven by the no-model loader check above; the model run's durable entity report/commit evidence and the resulting child-env absence proof remain owed by the live lane.
- DONE: AC-9 non-live half — ordinary helper tests plus formatting and live-tagged vet/build.
  `go test ./internal/ensigncycle -run 'TestPiDefaultExtensionRoots|TestPiIntercomPackageRootDiscoversIndependently|TestPiIsolatedHome|TestPiLiveEnv' -count=1` passed; `gofmt -w ./internal/ensigncycle`, `go vet -tags live ./internal/ensigncycle`, `go build -tags live ./internal/ensigncycle` passed; falsifier live-tag-the-helper is contradicted by the helper file carrying no build constraint.
- SKIPPED: Repository-wide `go test ./...` / `go test ./... -race` (the prior cycle's FAILED item), resolved with bounded evidence.
  `TestCodexResolveManifestAgainstInstalledHost` + `TestVersionAmbiguousMarkersExitZero` (`internal/cli`) and `TestSurveyCodexPresenceThroughSync` (`skills/integration`) fail identically on this tree at `631b482e6` and on base `1f41f289f` (`/tmp/spacedock-base`): local codex plugin cache, ambient `PI_CODING_AGENT`, and `blank_cwd=0`. Pre-existing, so the item is skipped; full/race runs were not repeated per the budget rule.

### Live Pi binding disposition

Four Pi XFAILs in `internal/ensigncycle/shared_live_runner_test.go`; no XFAIL added, no selector narrowed.

- re-anchored: `TestLiveCommonRejectionFlow` -> `6h3teccccn3qh71yqcmjbjx4` (was `p17swb3375rt525fn7f8xt7e`). Evidence: commit `d8f758bd1` (pre-rebase `b4ef3a90f`); the prior owner was the archived timeout repair, while `6h3teccccn3qh71yqcmjbjx4` owns the live Pi rejection-worker topology fault.
- kept: `TestLiveCommonOwnedConflictOwnerHandoff` -> `fe7bfjz9sb8wyckmnnm3ncjx`.
- kept: `TestLiveCommonKeepMovingPosture` -> `x02375wsg6q61xek7p0t36j2`.
- kept: `TestLiveCommonSmallestSufficientMechanism` -> `h30c9jrfcf21fdh2qs5z58sd`.
- cleared: none. Evidence: staff-review M1 outcome row ``d52`` records "Four retained bindings, one re-anchor, zero clearances", and this cycle's diff does not touch shared_live_runner_test.go's binding owners.

### Summary

Folded the M1 isolated-home discovery contract into the pi-live lane: one setup contract now links and registers both substrate packages in the isolated home (Spacedock kept as one absolute checkout path), discovery reads the real installed roots from the agentDir settings, and explicit overrides stay independent while default mode scrubs both package-root variables. Deterministic non-live tests, gofmt, live-tagged vet/build, the retained negative control, and a no-model real Pi 1.0.0 loader check pass; AC-6/AC-9's live front-door model run and the full/race suites remain unclaimed. The prior cycle's FAILED suite item is resolved as skipped on bounded cross-tree evidence.

Rebase note: the remote branch had advanced past the worktree (pin refresh rebased onto the `mc` merge #817); cycle-2 commits were rebased onto the remote tip, so the pre-rebase SHAs `094e06ff7`/`b4ef3a90f`/`115e37fad` are now `818b9c4ad`/`d8f758bd1`/`631b482e6`. The lane pins and all four manifest-derived assertions are unchanged.

### AC scan evidence

Command: `/opt/homebrew/Caskroom/spacedock@next/0.28.0-pre3/spacedock status --read /Users/clkao/git/spacedock-research/spacedock-v1/docs/dev/.spacedock-state/pi-live-lane-pin-refresh/index.md --ac-scan --json --workflow-dir docs/dev`.

| AC | Body line | Falsifier | Cited in report |
| --- | --- | --- | --- |
| AC-6 | 246 | remove-registration-seeding | implementation (cycle 2), AC-6 negative control + SKIPPED smoke |
| AC-7 | 248 | point-settings-elsewhere / write-file-prefix | implementation (cycle 2), AC-7 line |
| AC-8 | 250 | set-one-override / re-add-empty | implementation (cycle 2), AC-8 line |
| AC-9 | 252 | live-tag-the-helper | implementation (cycle 2), AC-9 line + SKIPPED smoke |


## Review-finding disposition

Validation at candidate `631b482e6`; advisory classifications only. No candidate edits or FO authorization are implied.

### V1 — Missing promised live evidence

- Observation: implementation cycle 2 explicitly SKIPPED the authorized front-door model smoke; no candidate green pi-live CI URL/SHA, journey grades, or durable child report/commit evidence was supplied. Independent no-model loading passes, but never launches or observes an ensign.
- Released user and normal workflow: Pi 1.0 users running the supported default-discovery first-officer/ensign path and the existing pi-live lane.
- Observable harm: acceptance cannot establish actual delegation, child boot, durable completion, or current-family journey compatibility; this is not evidence that the product fails at runtime.
- Authority: value-ac[AC-6] requires the authorized model-backed isolated-home run, not merely package loading. AC-1 and AC-9 also retain unfulfilled live requirements.
- Trigger evidence: the cycle-2 SKIPPED entry, absence of a candidate run reference/results, and this validation's intentionally no-model observation boundary.
- Proposal: **evidence defect / Material**; ownership = this task's live acceptance evidence; disposition = hold acceptance pending separately authorized proof, not an automatic implementation retry. A budget prohibition is not an out-of-promise trigger and does not turn a promised value AC into a deferred risk.
- Closure: after V2's observation gap is addressed or the captain explicitly revises the criterion, record an authorized `TestLivePiFrontDoorSmoke` at the candidate SHA with BOTH root variables actually absent, no supplied substrate extension paths, observed Spacedock/ensign and both substrate tools, passing boot grade, committed entity stage report, clean entity state, and root/child session artifacts. AC-1 additionally needs the candidate green existing pi-live CI URL/SHA, installed 1.0.0/0.75.0/0.16.0 logs, doctor success, unchanged common-journey grades and smoke results under the existing XFAIL policy. Do not substitute override-enabled success or this loader probe.

### V2 — AC-6 names observation that the shipped smoke grader does not make

- Observation: `internal/ensigncycle/pi_live_runner_test.go:320-432` and `pi_evidence_grade_impl_test.go:57-99` grade artifact forwarding, ordered ensign reads, absence of first-officer reads, sessions and durable output. Neither records/asserts an inventory containing both substrate tools or absence of the two package-root variables. AC-6's assertion that this grader “still requires … both substrate tools” does not describe the shipped implementation.
- Released user and normal workflow: default-discovery smoke acceptance on the supported Pi lane.
- Observable harm: a passing grade can be presented as default-discovery/both-tools proof without observing those conditions; subagent execution establishes subagent availability, not independent intercom availability. The workflow also still exports both override variables into GITHUB_ENV, and its unchanged smoke step does not unset them.
- Authority: value-ac[AC-6] promises both required extensions through native discovery with both variables unset; the named proof boundary must distinguish that from explicit fallback.
- Trigger evidence: `go test ./internal/ensigncycle -run '^TestPiFrontDoorEvidenceGrade' -count=1 -v` passes with the existing model/timestamp/cost-only session fixtures (no substrate inventory/environment evidence); source tracing confirms the outer boot grader does not add those observations.
- Proposal: **evidence defect / Material**; ownership = this task's folded AC-6 proof; disposition = FO-authorized narrow evidence correction at the existing smoke/grade boundary, or route to captain if the intended criterion differs. No new controller/lifecycle layer is proposed. Candidate remains unchanged.
- Closure: an independently falsifiable observation in the authorized default-mode smoke must reject a missing intercom/subagent tool and any present override (including empty assignments), preserve the ensign/durable-state checks, and run with the two workflow overrides explicitly absent. Current no-model positive/negative proof remains complementary, not a replacement for the promised live observation.

### Evidence quality and residual scope

- AC-2's Go guard deliberately duplicates the reviewed Node body: it proves executable wiring, not Node semantics. The separately executed real-package matrix supplies the independent behavioral oracle.
- `TestPiIsolatedHomeNegativeControlDropsSubstrateRegistrations` invokes the harness's own settings reader, not Pi. Standing alone its “Pi discovery” interpretation is self-referential and is rejected as runtime proof. The independent Pi 1.0 loader experiment below closes that specific load-bearing mechanism claim.
- No other AC depends solely on checking its own prose. AC-7/8 describe the shipped resolver, registrations and environment behavior; AC-9 correctly places the helper outside the live build tag, but its live half remains unproved. AC-5's explicit full/race-pass requirement remains unestablished, despite the prior report's “resolved as skipped” wording; this validation neither reruns nor claims those suites.
- No newly established outcome defect or deferred-risk finding. Existing four registered journey exceptions remain owned; retained XFAILs are not new proof of product correctness. Required independent reviewer gate remains outstanding.

## Stage Report: validation

- DONE: Reproduce the deterministic evidence for every criterion in the body, AC-1 through AC-9: run the ordinary helper tests in pi_default_extensions_test.go, the focused controls tests, and the live-tagged vet and build. For each criterion, confirm that the criterion describes the shipped implementation and not obsolete prose, and flag any criterion whose evidence is self-referential.
  Candidate `631b482e6`; bounded evidence reproduced below. AC-1/5/6/9 are incomplete, not silently redefined as green; V2 identifies the obsolete AC-6 grader claim and Evidence quality identifies the helper-only circular runtime interpretation.
- DONE: AC-1 deterministic pin proof; AC-2 declarations and runtime files; AC-3 integrity enforcement.
  `go test ./internal/release/... -count=1` PASS (22.6s), including all 34 mutations; installed local names/versions independently read as 1.0.0/0.75.0/0.16.0. Reverting a pin, removing either check, or bypassing verified-pack wiring fails the guards; no CI-green inference.
- DONE: AC-2 independently execute both candidate Node heredocs against a fresh real `npm pack pi-subagents@0.75.0` extraction.
  Each checkpoint: published/restored package exits 0; either missing runtime file, either directory in place of a file, missing/empty extensions, missing bridge, types-only bridge and missing second extension each exit 1 (18 negative results total). Removing existence/regular-file/all-entry checks falsifies these expectations.
- DONE: AC-3 independently execute the candidate's actual `verified_pack` shell function.
  Published hash exits 0; `sha512-corrupted-for-negative-probe` exits 1 before installation; disabling only the mismatch exit returns 0 and falsifies the negative. Independent SHA-512 of tarball bytes equals the body snapshot; no global installation was changed.
- DONE: AC-4 current-family instructions/comment and AC-5 protected boundaries.
  One-off comparison against the independent three-version snapshot passes; stale doc-version mutation fails. Against `12b695f26`, only the two authorized Pi setup steps differ: other lanes/actions/Node/selectors, compatibility step, doctor call and CLI floor are unchanged (binding owner re-anchor audited separately).
- DONE: Detached semantic adversarial pass over the changed CI assertions.
  In a disposable `git archive HEAD` snapshot under the assigned worktree, revert-agent-pin, restore-setup-source-assertion and corrupt-subagents-integrity each fail `TestPiLivePinsAndSubstrateAssertions`; downgrade-other-lane-checkout fails `TestNode24ActionsPinnedAtMinimum`. Restored snapshot passes both guards plus the mutation companion. This is validator evidence, not an independent reviewer sign-off.
- DONE: Independently reproduce the no-model loader check on the exact isolated layout: both package-root variables unset, no substrate extension path supplied, and require spacedock plus pi-subagents plus pi-intercom loaded with the subagent and intercom tools and the ensign skill. Then remove both npm registrations and require zero substrate tools. That negative control is the load-bearing fact.
  Node 24.13.1 / Pi 1.0.0 `DefaultResourceLoader({cwd: emptyWorkflow, agentDir: cleanHome+'/.pi/agent'})`, then `await reload()`; no additional paths/factories. Both actual substrate roots symlinked at agentDir/npm/node_modules; settings packages exactly `["npm:pi-subagents","npm:pi-intercom",absoluteCandidateCheckout]`; both root variables and parent PI_SUBAGENT_* absent before import. Inspect `getExtensions()` errors/paths/tool-map keys and `getSkills()` names.
  Positive: exactly spacedock.ts, pi-subagents/index.js, pi-intercom/index.ts; exact tool set `{subagent,bg_wait,subagents_enable,intercom}`; ensign present; zero errors. Remove both npm entries, retain both symlinks: exactly spacedock.ts, ZERO substrate tools, ensign retained, zero errors. Replace sole checkout entry with `file:`+path: ZERO extensions/skills/tools. Fresh Node process per variant; all assertions pass; no model/auth/session calls made.
- DONE: AC-7 and AC-8 focused ordinary helper/controls tests.
  `go test ./internal/ensigncycle -run 'TestPiDefaultExtensionRoots|TestPiIntercomPackageRootDiscoversIndependently|TestPiIsolatedHome|TestPiLive' -count=1 -v` PASS (0.25s): custom-agentDir npm roots, package-name-matched local roots, non-sibling intercom, exact settings order, registration negative, unset-not-empty env and each independent override; changing those expected paths/registrations/presence makes the respective test fail.
- DONE: AC-9 non-live visibility, formatting and live-tagged checks; bounded AC-5 build.
  Ordinary tests discover the helper without `-tags live`; `go vet -tags live ./internal/ensigncycle`, `go build -tags live ./internal/ensigncycle`, `go build ./...`, and `git diff --check` PASS. Read-only `gofmt -l` on all six changed Go files returns empty; no repository-wide formatting mutation in this validation-only stage.
- DONE: Verify the recorded disposition of every live Pi binding against the live registry and the owner bodies, and confirm the rejection-flow binding names an active owner rather than the archived one.
  `SPACEDOCK_LIVE_STATE_DIR=/Users/clkao/git/spacedock-research/spacedock-v1/docs/dev/.spacedock-state go test ./internal/contractlint -run '^TestRuntimeLive(RegistryReconciliation|TODOOwnersAreActive)$' -count=1 -v` PASS; derived Pi XFAIL count exactly 4; active-owner join runs (not skipped).
  Re-anchored rejection-flow -> `6h3teccccn3qh71yqcmjbjx4` (active backlog topology owner); old `p17swb3375rt525fn7f8xt7e` is archived done/PASSED for timeout only. Kept owner-handoff -> `fe7bfjz9sb8wyckmnnm3ncjx`, keep-moving -> `x02375wsg6q61xek7p0t36j2`, smallest-mechanism -> `h30c9jrfcf21fdh2qs5z58sd`; all active backlog bodies own the exact semantics and require exact PASS before removal. Keep-moving's historical XPASS alone does not authorize clearance. Cleared: zero.
- DONE: Report the missing live model proof explicitly as a finding, with the defect-kind and release-scope classification that this stage defines. Do not run that smoke. State what evidence would close it.
  V1 = evidence defect / Material, not demonstrated runtime failure or deferred risk; V2 = evidence defect / Material at AC-6's named observation boundary. Exact closure evidence and advisory ownership/dispositions recorded above; no candidate correction authorized or attempted.
- SKIPPED: Full repository suite, race suite, live model front-door smoke and green pi-live CI execution.
  Explicit FO budget prohibition; no prohibited command run. Prior cross-tree failures remain historical evidence only. AC-1/6/9 live acceptance and AC-5 full/race-pass clauses are not claimed satisfied.

### Summary

Recommend **REJECTED / hold acceptance** for material evidence findings V1/V2, not a product-failure diagnosis or a mechanism/design reset. Deterministic pin, package, environment and independent loader positive/negative checks pass, including the load-bearing zero-tools result with both symlinks retained; all four Pi bindings have active semantic owners. Code HEAD and candidate files remain unchanged; only this state report is committed, and live proof plus the independent reviewer gate remain owed.

## Stage Report: implementation (cycle 3)

- DONE: Require both extension tools, subagent and intercom, in the inventory the grade records, and fail the grade on a missing tool.
  `piBootContractEvidence` gains `Tools`; `buildPiFrontDoorEvidenceGrade` rejects an inventory lacking either name and records it as `isolated_discovery.tools`. The live smoke builds the inventory from the run's recorded toolCalls across the parent and child JSONL transcripts via `recordedPiToolInventory`. New negative subtests `missing subagent tool` / `missing intercom tool` fail. (AC-6 / review V2)
- DONE: Require neither PI_SUBAGENTS_PACKAGE_ROOT nor PI_INTERCOM_PACKAGE_ROOT present in the child environment, including as an empty assignment; fail the grade on a present override.
  `piBootContractEvidence` gains `ChildEnv`; the grade rejects any KEY=VALUE pair for either variable (empty value included) and records `isolated_discovery.package_root_env_absent`. The smoke passes the env it launches the child with. New negative subtests `present package root` / `empty package root assignment` fail. (AC-6 / review V2)
- DONE: Stop the workflow exporting the two overrides into the run.
  Removed the two `>> "$GITHUB_ENV"` echoes from `Install Pi CLI and substrates`. The `Verify Pi current-checkout setup` step now sets `PI_SUBAGENTS_PACKAGE_ROOT` / `PI_INTERCOM_PACKAGE_ROOT` shell-locally (never exported) so the explicit-override path stays deliberately exercised there while the smoke step sees neither variable.
- DONE: Smoke records the intercom tool so the inventory can carry it.
  `piLiveSmokePrompt` now instructs the FO to call `intercom({action:"list"})` once; the parent transcript records that toolCall, so a native-discovery run records both subagent and intercom.
- DONE: Keep every existing check; add no framework or lane; do not touch pins, guard, or manifest assertions.
  Ensign read order, durable state commits, session checks, and the four manifest checks are unchanged; `internal/release/pi_live_pins_guard_test.go` and the pins are untouched. Changed only the smoke, its grade, and the workflow env.
- DONE: Focused validation per the FO budget.
  `go test ./internal/ensigncycle -run 'TestPiFrontDoorEvidenceGrade|TestPiLiveEnv|TestPiDefaultExtensionRoots|TestPiIsolatedHome|TestPiIntercomPackageRoot' -count=1` PASS; `go vet -tags live ./internal/ensigncycle` PASS; `go build -tags live ./internal/ensigncycle` PASS; `go test ./internal/release/... -run TestPiLivePinsAndSubstrateAssertions -count=1` PASS; `go test -tags live ./internal/ensigncycle -run 'TestPiFrontDoorEvidenceGrade|TestPiLiveEnv' -count=1` PASS; `gofmt`/`git diff --check` clean. Full suite and race suite not run per the budget prohibition.
- SKIPPED: `TestLivePiFrontDoorSmoke` model run with both variables absent.
  Out of the bounded budget; the observation is implemented and unit-falsified but not exercised against a live model. AC-6's live acceptance still owes a run URL/SHA and the durable grade JSON.

### Changed evidence boundary

The grade now observes two facts it previously did not, and both are falsified by a named mutation:

- tool inventory: `isolated_discovery.tools` must contain `subagent` and `intercom`; removing either tool from the recorded calls fails the grade.
- environment: `isolated_discovery.package_root_env_absent` must hold; any present override, including an empty assignment, fails the grade.

### What the smoke can now tell apart

- Native discovery with both tools: both variables absent from the workflow, the isolated home's `settings.json` registers `npm:pi-subagents` + `npm:pi-intercom`, and the run records `subagent` and `intercom` -> the grade passes.
- An override we ran on purpose: with either `PI_SUBAGENTS_PACKAGE_ROOT` or `PI_INTERCOM_PACKAGE_ROOT` set (empty or nonempty), the child environment carries it and the grade fails on `package_root_env_absent`; a run that needs the override can no longer be presented as native discovery. Explicit overrides remain independently honored by `piLiveEnv`/`piIsolatedExtensionRoots` for that deliberate path.

### Summary

Implemented the review-V2 evidence-gap fix at the existing grade boundary: the smoke now requires both extension tools in the recorded inventory and requires both package-root overrides absent from the child environment, and the workflow no longer exports those overrides. Focused tests plus live-tagged vet/build pass; the live model smoke and the full/race suites were not run per the FO budget. Code commit `32049bdbd`.

## Stage Report: implementation (cycle 4)

- DONE: Restore native package-root discovery for the Pi live lane: in the "Install Pi CLI and substrates" step of .github/workflows/runtime-live-e2e.yml, after installing to $pi_npm_root, create-or-merge $HOME/.pi/agent/settings.json registering packages ["npm:pi-subagents","npm:pi-intercom"], so internal/ensigncycle/pi_default_extensions_test.go's piDefaultExtensionRoots resolves both roots from the real agent dir with neither PI_SUBAGENTS_PACKAGE_ROOT nor PI_INTERCOM_PACKAGE_ROOT set. Keep both variables unexported.
  Code commit `7e13e49b4`: the install step now writes `<agentDir>/settings.json` via a node create-or-merge block after the substrate install, appending `npm:pi-subagents` and `npm:pi-intercom` to `packages` while preserving any existing packages/keys. agentDir is `$HOME/.pi/agent`, so `piDefaultExtensionRoots` resolves `npm:pi-subagents` -> `$HOME/.pi/agent/npm/node_modules/pi-subagents` (the same tree `$pi_npm_root` installs into). Removing it is the falsifier `remove-registration-seeding`.
- DONE: Add a deterministic guard that the registration cannot be dropped: extend internal/release/pi_live_pins_guard_test.go's workflow-structure check to require the install step writes the settings.json registration and still does not export the two root variables.
  `assertPiLivePinsAndSubstrateAssertions` now requires the executable commands `settings_path="$HOME/.pi/agent/settings.json"` and `node - "$settings_path" <<'NODE'`, requires both `'npm:pi-subagents'`/`'npm:pi-intercom'` literals in the install step run, and rejects any executable install command carrying either root variable with `GITHUB_ENV` or `export`. Six new negative mutations fail the guard: `settings-registration-removed`, `settings-path-not-agent-dir`, `settings-registration-subagents-missing`, `settings-registration-intercom-missing`, `re-export-subagents-root` (GITHUB_ENV echo), `re-export-intercom-root` (`export`).
- DONE: Bounded proof, no live model: go test ./internal/release/... -run TestPiLivePinsAndSubstrateAssertions, go test ./internal/ensigncycle -run 'TestPiDefaultExtensionRoots|TestPiIsolatedHome|TestPiLiveEnv|TestPiIntercomPackageRoot' -count=1, gofmt on changed Go files, go build ./..., and one go test ./... in the worktree.
  All bounded commands pass. `go test ./...` exit 1 with ONLY the pre-existing failures this dispatch named: internal/cli `TestCodexResolveManifestAgainstInstalledHost` + `TestVersionAmbiguousMarkersExitZero`, and skills/integration `TestSurveyCodexPresenceThroughSync`. `internal/release`, `internal/ensigncycle` (243.7s, no timeout), `internal/contractlint`, and `go build ./...` are green. Result: the workflow now reproduces exactly the settings/root-symlink layout the deterministic tests encode, so discovery no longer resolves nothing.
- DONE: Commit on spacedock-ensign/pi-live-lane-pin-refresh and push; touch nothing outside the pi-live workflow install step and that guard.
  Pushed `32049bdbd..7e13e49b4`. Diff is two files, +58/-0: 26 lines in the install step, 32 lines in the guard. No other lane, pin, guard assertion, docs file, or ensigncycle test was touched.

### Independent checks

- The extracted node registration block was run directly in a temp HOME: fresh create wrote `{"packages":["npm:pi-subagents","npm:pi-intercom"]}`; merging into `{"theme":"dark","packages":["npm:other"]}` preserved both the unrelated key and the existing package while appending the two registrations. This is the create-or-merge behavior the guard pins.
- Falsification: reverting the registration, pointing settings outside the agent dir, dropping either package literal, or re-adding either root variable as an export/GITHUB_ENV echo each makes `TestPiLivePinsAndSubstrateAssertions` fail (6 new subtests, all pass on the candidate).

### Summary

Fixed the live-lane setup failure at its real cause: commit `32049bdbd` stopped exporting the two package-root variables but never registered the substrates in the real agent dir, so isolated-home discovery resolved nothing. The install step now create-or-merges `$HOME/.pi/agent/settings.json` with both `npm:` registrations and keeps both variables unexported; the release guard and six mutations make the registration non-droppable. Bounded suites, build, gofmt, and a direct create-or-merge exercise pass; the only `go test ./...` failures are the known pre-existing ambient ones. No live run was started.

#### Known pre-existing failures (not fixed, out of scope)

- `internal/cli`: `TestCodexResolveManifestAgainstInstalledHost` (ambient local codex plugin cache), `TestVersionAmbiguousMarkersExitZero` (ambient `PI_CODING_AGENT=true`).
- `skills/integration`: `TestSurveyCodexPresenceThroughSync` (`blank_cwd=0`).
- `internal/ensigncycle`: `TestCodexProcessActivityResetsQuietBudget` is a known CI race outside scope; it did not fire in this worktree run.

## Stage Report: implementation (cycle 5)

AC-10 (no regression): make the offline `go test ./...` gate deterministic and test our decision logic, not the machine.

- DONE: Make the reset-on-activity offline decision deterministic and machine-independent: rewrite TestCodexProcessActivityResetsQuietBudget so the helper's activity is supplied by a test-controlled signal or an injected deadline (assert the decision logic), remove the `const quietBudget = 250 * time.Millisecond` wall-clock race and the `result.duration > 4*quietBudget` machine-speed assertion, and keep TestCodexProcessQuietTimeoutPreservesFaultEvidence (mode `stall`) as the deterministic offline kill-path test.
  `streamwatch_test.go` adds `now`/`sleep` clock seams (production default `time.Now`/`time.Sleep`); `codex_single_run_test.go` drains `drainCodexToTerminal` against a fake line source and never-exiting fake proc, injecting activity while a fake clock advances one poll step per sleep. The test asserts the no-progress decision (`stepTimeout`, kill happens only after the activity window closes) with no real timer; the old `progress-then-exit` helper mode and the duration/count machine assertions are gone. `stall` kill-path test unchanged. `go test ./internal/ensigncycle -run 'TestCodexProcessActivityResetsQuietBudget|TestCodexProcessQuietTimeoutPreservesFaultEvidence' -count=5` PASS (4.2s).
- DONE: Isolate the three environment-dependent tests from the operator's machine (isolated CODEX_HOME / fake host, scrubbed markers, existing HOME/data-dir isolation).
  `TestCodexResolveManifestAgainstInstalledHost` now sets an isolated `CODEX_HOME`, installs `spacedock@spacedock` from a local-path marketplace via the production `execHost.Install`, and asserts the resolver returns that isolated install's `.codex-plugin/plugin.json` under `codexHomeDir` — never the operator `~/.codex`. `TestVersionAmbiguousMarkersExitZero` scrubs all four marker vars (`PI_CODING_AGENT`, `PI_CODING_AGENT_DIR`, plus its two controlled ones). `TestSurveyCodexPresenceThroughSync` keeps `HOME`/`AGENTSVIEW_DATA_DIR`/`CODEX_SESSIONS_DIR`/`CLAUDE_PROJECTS_DIR`. All three pass with an operator `CODEX_HOME`/`PI_CODING_AGENT`/`PI_CODING_AGENT_DIR` set.
- DONE: Survey-test decision (pre-approved, not escalated): keep the sync-ingest property (repo-project Codex row count = 2).
  Chose option (2): deleted the `blank_cwd > 0` expectation. Whether agentsview persists or blanks a Codex cwd is an ambient binary-version behavior, not this repo's code; the real sync here persists the cwd (`blank_cwd=0`). The query's blank-cwd counting stays deterministically covered by the fixture-only codex-presence test in `survey_queries_test.go`, which injects `cwd=''` rows. Kept the count-2 assertion and the HOME/data-dir isolation; test still runs, not skipped.
- DONE: Bounded proof.
  `go test ./internal/ensigncycle -run '...' -count=5` PASS; `go test ./internal/cli -run '...' -count=1` PASS, and again with operator `CODEX_HOME`/`PI_CODING_AGENT`/`PI_CODING_AGENT_DIR` set PASS; `go test ./skills/integration -run TestSurveyCodexPresenceThroughSync -count=1` PASS; `gofmt -l` on the five changed files empty; `go build ./...` PASS; one `go test ./...` in the worktree EXIT=0 (internal/cli 126s, internal/ensigncycle 204s, no package timeout).
- DONE: Write the cycle-5 report citing AC-10, commit on `spacedock-ensign/pi-live-lane-pin-refresh` and push; touch only named tests plus the needed clock helper.
  Code commit `4e95513e3` pushed (`7e13e49b4..4e95513e3`); five files changed, +114/-52, all named tests plus the `streamwatch_test.go` clock seam. The three previously-ambient failures reported in cycle 4 now pass offline; the cycle-4 note that they were pre-existing ambient failures was correct and item 2 is what makes them machine-independent.
- SKIPPED: Real Codex pacing against the real watchdog.
  Deliberately not re-asserted offline; it stays covered by the live lane. The offline direction that remains is the deterministic `stall` kill-path test.

### Summary

AC-10 is satisfied: the reset-on-activity test now asserts the no-progress decision through an injected clock and test-controlled activity instead of racing real sleeps, and the three host-dependent tests build their own isolated home so the offline gate no longer reads the operator's machine. The blank-cwd expectation was deleted under the pre-approved option (2) because the blank-cwd behavior belongs to the ambient agentsview version, not our code, and is already covered deterministically by the fixture-only query test. Code commit `4e95513e3` is pushed; `go test ./...` is green in the worktree.

### AC scan evidence

AC-10 cited above with its named falsifiers: **race-real-sleeps** is contradicted because the reset test no longer sleeps against a real timer; **inherit-operator-home** is contradicted because the three tests pass with an operator `CODEX_HOME`/`PI_CODING_AGENT` set.

## Stage Report: implementation (cycle 6)

- DONE: Single source of truth for the Pi family stamp: create one Go source holding, for each of @earendil-works/pi-coding-agent, pi-subagents and pi-intercom, the spec, version, sha512 integrity, plus the readiness floor(s) the gate enforces. Choose the shape ... make .github/workflows/runtime-live-e2e.yml consume it: no version literal, no integrity literal, and no JavaScript ... Move the build step if the install step needs the binary.
  Chose `internal/pilive` (dependency-free leaf package) + the `cmd/spacedock-pilive` helper the workflow builds and calls, NOT a user `spacedock` subcommand group: this is CI-lane tooling, and the helper leaves the user CLI surface, grouped help, and its tests untouched. `pilive.go` holds each spec/version/integrity and the floors (Node 22.19.0, pi-coding-agent 0.83.0, pi-subagents 0.53.0); `internal/cli`'s `piVersionFloor` is `pilive.PiCodingAgentFloor`. Workflow now invokes `spacedock-pilive install|guard|verify-manifest`; the build step moved before install and also builds the helper.
- DONE: Move all four inline Node heredocs and the JS one-liners into Go: (a) npm-pack integrity verification, (b) installed ... name+version verification, (c) the installed-substrate manifest assertions ..., (d) the settings.json create-or-merge registration ..., (e) the compatibility guard ... The workflow may only invoke commands.
  (a) `PackTarball`+`VerifyIntegrity`; (b) `VerifyInstalled`; (c) `VerifyRuntimeManifest` at both checkpoints; (d) `MergeSettings`; (e) `Guard` (engine+version floors; `CompatExportPath` resolves, `CompatExportLoads` loads). No `node`/heredoc/JS-one-liner remains in the workflow. The one unavoidable Node use is `CompatExportLoads` (importing an ESM module is inherently Node) — invoked from Go, not the workflow.
- DONE: Refresh the pins from the registry at this moment and record the exact commands ...
  `npm view <spec> dist-tags.latest version dist.integrity --json` for each of `@earendil-works/pi-coding-agent`, `pi-subagents`, `pi-intercom` on 2026-10-04 returned 1.0.2 / 0.75.0 / 0.16.0 with integrities 3ZdIghMS… / RO4DiTJM… / ClGQuovP… (full values in `internal/pilive/pilive.go`). Discrepancy flagged: the captain named 1.0.1, the registry had moved to 1.0.2 by fetch time, so 1.0.2 is pinned; pi-subagents and pi-intercom did not move, so the family moved by one package only. `latest` is recorded as the resolution dist-tag, never as the pin.
- DONE: Delete every duplicate copy of the stamp ...
  Deleted `internal/release/pi_live_pins_guard_test.go`; floors removed from the guard (now Go, reading `pilive`); version literals/comments removed from `internal/cli/pi_frontdoor_test.go` (0.85.1, 0.83.0 → `pilive.PiCodingAgentFloor` + synthetic 99.0.0), `internal/ensigncycle/pi_default_extensions_test.go` (fixture → `npm:pi-subagents@0.0.0-synthetic`), `internal/ensigncycle/pi_live_runner_test.go` and `claude_runtime_helpers_test.go` (0.53.0 comments reworded); family version list removed from `docs/runtime-live-ci.md` (now `go run ./cmd/spacedock-pilive print-install`). The only remaining 0.75.0 mentions are historical prose in `docs/roadmap/pi-ux/staff-review-2.md`, not consumed.
- DONE: Add behavioural Go tests for the moved logic, with negative cases ... Keep exactly one independent oracle: the registry snapshot compared against the single source ...
  `internal/pilive/pilive_test.go` (behavioural, temp-dir fixtures) and the single oracle `registry_oracle_live_test.go::TestPiLivePinsMatchRegistry` (live-tagged), which queries `npm view` at CHECK TIME and stores no copy of the numbers (per captain steering) — it fails only when the registry's `latest`/integrity actually move. Kept the freshness check; did not store a snapshot.
- DONE: Prove: offline `go test ./...` green in the worktree; the new Go tests pass and fail on the right negatives; and a LOCAL live Pi lane run ... Do NOT start a CI lane run.
  `go test ./...` EXIT=0 at b6018c893 (cli 218s, ensigncycle 349s, release 31s). `go test -tags live ./internal/pilive -run TestPiLivePinsMatchRegistry` PASS. Helper exercised locally: `install` installed 1.0.2/0.75.0/0.16.0 and registered both substrates (create-or-merge preserved unrelated keys); `guard` verified the compat load; `verify-manifest` resolved both real 0.75.0 runtime files. No CI lane was started.

### Local live Pi run (captain's credentials; CODEX_AUTH_JSON/OPENAI_API_KEY unset)

- `TestLivePiFrontDoorSmoke` PASS (81.5s): grade artifact `isolated_discovery.tools` carries `subagent` + `intercom`, `package_root_env_absent: true` — the AC-6 observation the cycle-3/validation V1/V2 findings asked for.
- `TestLiveCommon` 17 journeys: 15 pass; `owned-conflict-owner-handoff`, `smallest-sufficient-mechanism`, `keep-moving-posture` XPASS alerts (recorded as-is); `rejection-flow` XFAIL engaged (owner 6h3teccccn3qh71yqcmjbjx4).
- Known deferred reds (reproduced on retry, `owner=` empty because their bindings live above on 9w/080d37f23): `TestLiveCommonAutoContinueAfterImplementation` (validation-worker-not-dispatched — observer evidence defect, owned by repair-pi-worker-lifecycle-observation) and `TestLiveCommonDefaultHeadlessGateStop` (gate-not-held / implementation-worker-not-dispatched — reference-path fragility, owned by repair-pi-recorded-gate-lifecycle). Not regressions from this diff (live test code unchanged) and not repaired here.

### Added-test falsifiability (each can fail for a reason other than editing a copy)

- `VerifyRuntimeManifest` tests: fail if a missing/empty `pi.extensions`, a missing/bare-string/types-only bridge, a missing file, a missing second extension, or a directory target is accepted.
- `MergeSettings` tests: fail if the merge drops unrelated keys or existing packages, duplicates on re-run, or overwrites malformed JSON.
- `ParsePackMetadata` tests: fail if malformed/empty/multi/partial npm output is accepted. `PackTarball` test: fails if the integrity check is dropped (corrupt hash accepted) or the spec/path is wrong.
- `VersionAtLeast`/`NodeEngineAtLeast`: fail if the compare direction/logic is wrong (0.82.9 vs floor). `CompatExportPath`: fails on a missing `./compat` export or missing file or wrong root. `CompatExportLoads`: fails if node's non-zero exit is ignored.
- `Guard` floor tests: fail if a floor check is dropped/inverted. `Install` test: fails if any installed name/version, manifest, skill, or settings-registration check is dropped.
- Per captain steering the replacement for the deleted fingerprint test is behavioural, not textual; the workflow-text structural guard I first wrote was deleted because its only failure mode was the workflow text changing.

### Summary

Moved the whole Pi family stamp and every piece of pi-live lane logic into `internal/pilive`/`cmd/spacedock-pilive`; the workflow now carries no version, integrity, or JavaScript. Pins are 1.0.2/0.75.0/0.16.0 (captain's 1.0.1 had moved to 1.0.2 by fetch time; the other two did not move). Behavioural negative tests plus one check-time registry oracle replace the deleted fingerprint test; offline `go test ./...` is green, the helper ran end-to-end, and a local live Pi lane run executed with the front-door smoke green (both substrate tools, overrides absent) and two known deferred reds owned by the repair worktrees. Code HEAD `d4fd8ead1` pushed; no CI lane started.

## Stage Report: implementation (cycle 7)

- DONE: Fold the pi-live helper into the existing spacedock-release binary: delete cmd/spacedock-pilive and expose its commands (pins, print-install, install, guard, verify-manifest) through cmd/spacedock-release; update .github/workflows/runtime-live-e2e.yml to build and call the release binary instead of a second artifact; keep every behaviour identical. The release must stop shipping a second binary.
  Code commit `a24c489bc`: `cmd/spacedock-pilive` deleted; the five commands live in `cmd/spacedock-release/pilive.go` and dispatch from `main.go`; the workflow builds `./spacedock-release` and calls `spacedock-release install|guard|verify-manifest`; `docs/runtime-live-ci.md` calls `go run ./cmd/spacedock-release print-install`. No second binary source remains under `cmd/`.
- DONE: Collapse internal/pilive/runner.go into internal/pilive/pilive.go wherever the split is only organisational; keep a file split only where it carries a real reason.
  `runner.go` deleted; its `Runner`/`ExecRunner`/`PackTarball`/`CompatExportLoads`/`Install`/`Guard` moved into `pilive.go` with merged imports. The split (pure logic vs exec orchestration) was organisational within one package, so it now has no reason.
- DONE: Cut tests that restate rather than check: remove duplication and tautological/textual assertions, keep the behavioural coverage and the negative cases. Do not re-add any guard previously deleted (the workflow-text guard and the internal/release fingerprint test stay deleted).
  `pilive_test.go`: deleted `TestParsePackMetadataAcceptsSingleEntry` (its expectations were its own input; `TestPackTarballVerifiesIntegrityAgainstThePin` exercises the same mapping), folded the directory-entry manifest case into the negative table, merged the Node-floor cases into `TestVersionAtLeastComparesFloors` (dropping the `NodeEngineAtLeast` dupe), and made the compat test actually check nested-then-global. All negative tables and behavioural cases stay. Grep confirms neither deleted guard is back.
- DONE: Do not weaken the two independent oracles: the Pi catalog check and the Codex models_cache check. Preserve the Pi registry oracle unchanged in behaviour (it must still query the registry at check time and store no copy).
  `internal/pilive/registry_oracle_live_test.go` is byte-unchanged; `go test -tags live ./internal/pilive -run TestPiLivePinsMatchRegistry` PASS — it still queries `npm view` at check time and stores no copy. The sibling `pin-lane-models-in-one-place` catalog/models_cache oracles are not present in this worktree and were not touched.
- DONE: Target a materially smaller net increment versus main, not a specific number. Report the new increment versus main (files and lines) and exactly what the trim removed.
  Measured: `git diff --numstat 4436ec14c` = **22 files, +1767/-296, net +1471** at `a24c489bc`, versus **22 files, net +1511** at pre-trim `9febb0c50` — only **-40 net**, not material. Removed: cmd fold net -19 (delete 143-line helper, add 106-line file + 18 wiring lines), runner collapse net -10 (runner.go -192 absorbed as +182 in pilive.go), pilive test dedupe -12, workflow +1. The remaining ~1471 lines are AC-6..AC-10 deliverable (`pi_default_extensions_test.go` +353, `pilive.go` +495, `pilive_test.go` +357, the ensigncycle grade/live tests, AC-10 clock/host rewrites). The dispatch requires keeping the behavioural coverage and negative cases, so a materially smaller surface needs deleting AC-named deterministic tests — a captain scope decision, not done here.
- DONE: Prove: offline `go test ./...` green in the worktree, and the folded helper still works end to end (build the release binary and run its pi-live install/guard/verify-manifest paths against the local pinned family). Do NOT start a CI lane run. Write the implementation stage report and commit code + state.
  `go test ./...` EXIT=0 (ensigncycle 370s, release 20.6s, integration 9.6s, no timeout); `go build ./...`, `go vet -tags live ./internal/pilive ./internal/ensigncycle`, `gofmt -l` on changed files, `git diff --check` all clean. Built `./spacedock-release` and ran the lane paths against the local pinned family: `install` EXIT=0 (1.0.2/0.75.0/0.16.0 installed, both npm substrates registered, unrelated settings keys preserved), `guard` EXIT=0 (compat resolved+loaded), `verify-manifest` EXIT=0 (resolved `index.js` + `src/api/intercom-bridge.js`). No CI lane run started.

### Summary

Folded the pi-live helper into `spacedock-release` (one CI binary, no second artifact), collapsed the two-file `internal/pilive` split into one, and deduped the restating pilive tests. Offline `go test ./...`, build, vet, and the real end-to-end install/guard/verify-manifest runs are green; the registry oracle is unchanged. The measured increment is net +1471 (22 files), only 40 lines below pre-trim, because the remainder is AC-6..AC-10 behavioural coverage the dispatch says to keep; a "materially smaller" surface requires a captain decision to narrow those ACs, so it is reported rather than silently done.

### Trim disposition (captain decision)

The captain called the volume and sent mh back for a trim. The trim executed the three
specified moves — fold the helper into `spacedock-release`, collapse
`internal/pilive/runner.go` into `pilive.go`, cut tests that restate rather than check —
and measured **mh's own increment at net +1453 across 21 files** against its base (9w
`c5ca95e74`), down from +1493 before the trim. The three moves relocate and merge code
rather than delete behaviour, so the reduction is only about 40 lines.

Cause and disposition: the declaration said +160 across 3 files; the real content is a
faithful Go port of the four moved heredocs plus mh's earlier AC-6..AC-10 folds, so the
estimate was the fiction, not the code. The captain accepted the volume (decision c): do
not cut behavioural tests and do not simplify load-bearing logic. Neither deleted guard
(the workflow-text guard, the `internal/release` fingerprint test) is re-added, and the
Pi registry oracle is untouched. The final increment is **net +1453 across 21 files**, and
the captain owns the tolerance decision.

## Stage Report: implementation (cycle 8)

- DONE: Cut mh by giving the port the specified shape: merge internal/pilive/*.go and cmd/spacedock-release/pilive.go into ONE file, keep ONE test file, drop registry_oracle_live_test.go unless a behaviour depends on it, and do not change what the workflow calls or how install/guard/verify-manifest behave.
  Code commit `2c061b2e9` on `spacedock-ensign/pi-live-lane-pin-refresh`. `internal/pilive/pilive.go` is the single port file (pins + floors + `Command`/`install`/`guard`/`verifyManifest` + helpers); `internal/pilive/pilive_test.go` is the single test file; `cmd/spacedock-release/pilive.go` (106 lines) and `internal/pilive/registry_oracle_live_test.go` (43 lines) are deleted. `cmd/spacedock-release/main.go` routes the five pi-live commands to `pilive.Command`. The workflow still calls `spacedock-release install|guard|verify-manifest|print-install`; `pins`/`print-install`/`verify-manifest` stdout and exit codes diff clean against the pre-merge binary `a24c489bc` (built side by side); install/guard keep the same npm/node invocations and verification order as `a24c489bc` and are exercised by the fake-runner unit tests, which error on any unexpected command.
- SKIPPED: The "roughly 250 lines for the port including its tests" target (1001 -> 759).
  The port is now ONE logic file at 464 lines (398 code) plus ONE test file at 295 lines (271 code). Reaching 250 total would require deleting load-bearing install/guard/verify-manifest behaviour (integrity verify, manifest resolution, settings create-or-merge, floor + compat guard) or the named negative cases; the note's escape hatch ("if a cut would remove real coverage, say so and keep it") is invoked, so the shape bound is met and the line bound is reported, not forced.
- DONE: Trim internal/ensigncycle/pi_default_extensions_test.go from 353 lines while keeping the deterministic setup contract and the negative control that removing both registrations loses both extensions, and cut tests that restate rather than check.
  File is now 308 lines. Kept `seedPiIsolatedHome`/`seedPiDefaultExtensions` + their helpers, `TestPiIsolatedHomeRegistersBothSubstratesAndAbsoluteSpacedock` (setup contract), and `TestPiIsolatedHomeNegativeControlDropsSubstrateRegistrations` (negative control). Deleted two discovery tests that re-asserted the resolver under alternate settings shapes: `TestPiDefaultExtensionRootsMatchesLocalSourceByPackageName` and `TestPiIntercomPackageRootDiscoversIndependently`.
- DONE: Report the new increment per file and the total, and name what was deleted.
  Port, before -> after: `internal/pilive/pilive.go` 495 -> 464, `internal/pilive/pilive_test.go` 357 -> 295, `cmd/spacedock-release/pilive.go` 106 -> 0, `internal/pilive/registry_oracle_live_test.go` 43 -> 0; port total 1001 -> 759. `internal/ensigncycle/pi_default_extensions_test.go` 353 -> 308. Entity net increment vs its base `c5ca95e74` +1453 -> +1159; vs `origin/main` (4436ec14c) +1471 -> +1177. Deleted: `cmd/spacedock-release/pilive.go`, `internal/pilive/registry_oracle_live_test.go`, and the two ensigncycle discovery tests named above.
- DONE: Validate the merge without the repository-wide suite or a CI lane.
  `go build ./...`, `go test ./internal/pilive/...`, `go test ./cmd/spacedock-release/...`, `go test ./internal/cli -run 'TestPiVersionAtLeast|TestPiVersionFloor'`, `go test ./internal/ensigncycle -run 'TestPiDefaultExtensionRoots|TestPiIsolatedHome|TestPiLiveEnv|TestPiIntercomPackageRoot'`, `go vet -tags live ./internal/pilive ./internal/ensigncycle ./cmd/spacedock-release`, and `go build -tags live ./internal/ensigncycle ./internal/pilive` all pass; `gofmt -w cmd internal` clean for changed files. No repository-wide `go test ./...`, race suite, or CI run was started.

### Coverage note (what stayed)

`TestPiDefaultExtensionRootsReadsRealInstalledLocation` is kept because AC-7 names `TestPiDefaultExtensionRoots…` as its deterministic proof; only the two resolver variants the fold no longer needs were cut. The port keeps every pilive negative table (`verifyManifest` 9 cases, `mergeSettings` malformed/merge/idempotent, `pack` corrupt/invalid integrity, version floors, guard floors, compat resolve/load, install registration).

### Risk

Deleting `registry_oracle_live_test.go` removes the only check-time query against the npm registry, so a pin drift no longer fails any Go test (no production behaviour depends on it, per the dispatch). The single source of the pins remains `internal/pilive/pilive.go`, but nothing now compares it to the published family.


## Stage Report: implementation (cycle 9)

Hard-bounded cut of the port (`internal/pilive/pilive.go` + `pilive_test.go`) from 759
to 400 lines or fewer, implementation plus tests. Priced every behaviour first.

- DONE: Price all seven behaviours with line cost and what breaks without each.
  Costs measured at HEAD `2c061b2e9` (pilive.go 464 lines): install 60; guard 41;
  verify-manifest 40 (verifyManifest) + ~11 dispatcher; integrity check (pack) 20;
  manifest resolution 40 (same verifyManifest body); settings create-or-merge
  (mergeSettings) 52; floor guard (versionAtLeast 9 + parseVersion 13 + guard floor
  branches ~15) 37; compatibility guard (compatExportPath 31 + compatExportLoads 11 +
  guard compat branch ~8) 50. Every one is load-bearing — consequences below — so the
  "no consequence" cut rule removes none.
- DONE: Remove what is not behaviour: per-function comment blocks deleted, the four
  repeated read/unmarshal sites replaced by one `readJSON`, `mergeSettings` rewritten
  to a single `map[string]any` round-trip, `parseVersion` folded into `versionAtLeast`,
  and the single-use `agentDir`/`fail` kept only where reused.
  Commit `d6900386c`: -303/+141 across the two files. `gofmt`/`go vet`/`go build`
  clean.
- DONE: Cut the test file to the behaviours the workflow exercises plus one negative
  each, deleting the copy-restating tables.
  `pilive_test.go` 295 -> 185: verify-manifest, integrity (pack), settings merge,
  floor comparison, guard floor, compat load, and install each keep one positive and
  one load-bearing negative. No acceptance criterion cites a committed pilive
  negative by name (AC-2/AC-3 proofs are the one-off real-package exercise), so none
  was retained for a citation.
- DONE: Freeze the call surface.
  The workflow still calls `spacedock-release install|guard|verify-manifest`; the
  pre-cut binary `2c061b2e9` and the cut binary produce byte-identical `pins` and
  `print-install` stdout and identical `verify-manifest` usage exit 2 and resolved
  "verified runtime file <path>" lines/exit 0 (diffed side by side).
- SKIPPED: Cut to 400 lines or fewer (reached 597, not 400).
  BLOCKED. The implementation alone is 412 physical lines after removing every
  per-function comment block (381 code + 9 comments + 22 blank); 412 > 400, so even
  deleting the entire test file cannot reach the target. The seven behaviours above
  account for the implementation; none has no consequence, so 400 requires deleting a
  behaviour. Cheapest blockers and what breaks: integrity check 20 lines (AC-3:
  corrupted tarballs install), verify-manifest 51 lines (AC-2: setup checkpoint
  loses the installed-manifest assertion). Reached number is 597 = 412 impl + 185
  tests. Not deleting the integrity check or the manifest checkpoint to hit a number.
- DONE: Validate without the repository-wide suite or a CI lane run.
  `go build ./...`, `go vet ./internal/pilive ./cmd/spacedock-release`,
  `go test ./internal/pilive/... ./cmd/spacedock-release/... -count=1`, and
  `go test ./internal/cli -run 'TestPi|TestPiFrontDoor' -count=1` all pass; `gofmt -l`
  clean; call-surface diff frozen. No repo-wide suite, race suite, or CI lane started.
- DONE: Commit code on `spacedock-ensign/pi-live-lane-pin-refresh` and push.
  `d6900386c` pushed (`2c061b2e9..d6900386c`); only the two pilive files changed, no
  workflow, pin, launcher, or ensigncycle change.

### Pricing (line cost at HEAD -> what breaks without it)

- install — 60 lines. Removing it fails the "Install Pi CLI and substrates" step: no
  pinned family installed, no agent-dir settings registration -> AC-1/AC-6/AC-7 break.
- guard — 41 lines. Removing it fails the "Guard Pi substrate compatibility" step.
- verify-manifest — 51 lines (verifyManifest 40 + dispatcher). Removing it fails the
  current-checkout manifest checkpoint -> AC-2's installed-manifest assertion gone.
- integrity check (pack) — 20 lines. Removing it installs tarballs without matching
  the published sha512 -> AC-3 breaks.
- manifest resolution (verifyManifest body) — 40 lines. Removing it lets setup accept
  missing/stale runtime paths -> AC-2 breaks.
- settings create-or-merge (mergeSettings) — 52 lines. Removing it leaves the isolated
  home with no `npm:pi-subagents`/`npm:pi-intercom` registrations -> AC-6/AC-7 break.
- floor guard (versionAtLeast+parseVersion+guard floors) — 37 lines. Removing it
  accepts Node < 22.19.0 or pi-coding-agent < 0.83.0 or pi-subagents < 0.53.0.
- compatibility guard (compatExportPath+compatExportLoads+guard branch) — 50 lines.
  Removing it stops verifying that `@earendil-works/pi-ai/compat` resolves AND loads.

### Summary

Trimmed the port 759 -> 597 with the workflow-observable surface frozen and all
behavioural checks intact. 400 is not reachable without deleting a load-bearing
behaviour because the behaviour-only implementation is already 412 lines; reported
the exact blocking behaviours (integrity check 20 lines, verify-manifest 51 lines and
their AC-3/AC-2 consequences) rather than deleting a check to hit the number. Code
commit `d6900386c` pushed; focused build/vet/tests/gofmt and the call-surface diff are
green; no repo-wide suite or CI lane was run.

### Cut-falsifiability

- `TestVerifyManifestResolvesDeclaredRuntimeFiles` fails if a missing declared runtime
  file or a manifest without `pi.extensions` is accepted.
- `TestPackRejectsIntegrityMismatch` fails if a tarball whose integrity does not match
  the pin is accepted.
- `TestMergeSettingsCreateOrMergeRejectsMalformed` fails if the merge drops unrelated
  keys/existing packages, is non-idempotent, or accepts malformed JSON.
- `TestGuardRejectsBelowFloorNodeAndSubstrate` fails if a Node or pi-subagents version
  below floor is accepted; `TestCompatExportLoadsPropagatesNodeFailure` fails if a
  non-zero node import is ignored.
- `TestInstallVerifiesAndRegistersThePinnedFamily` fails if any installed
  name/version, manifest, skill, or settings registration check is dropped.

## Stage Report: implementation (cycle 10)

An independent review of the 597-line port found three correctness faults and said
explicitly not to rewrite for a smaller number. Fixed those three; changed nothing
else. Call surface frozen: same `install`/`guard`/`verify-manifest` commands, same
exit codes, same printed output. Code commit `52223fe29`.

- DONE: `mergeSettings` now preserves every existing `packages` entry, object entries included, and appends only missing registrations.
  `internal/pilive/pilive.go`: the string-only `[]string` accumulator became
  `pkgs []any` (kept verbatim) with a parallel `names []string` used only for the
  missing-registration check (`switch` reads a string entry or a map's `"source"`).
  Mutation A (restore the old string-only filter): `TestMergeSettingsPreservesEntriesAndIsIdempotent` fails at `pilive_test.go:84` with
  `merged packages=[npm:other npm:pi-subagents npm:pi-intercom]` — the `{"source":"npm:other-object"}` entry is gone. On the fixed code the same test returns `npm:other,npm:other-object,npm:pi-subagents,npm:pi-intercom`.
- DONE: `compatExportPath` selects the installed nested pi-ai copy first and fails when its `./compat` export is invalid, instead of falling through to a working global copy.
  It now `os.Stat`s each candidate's `package.json`, skips only a genuinely absent
  copy, and returns an error when an existing installed copy has no valid
  `./compat` import; the global copy is consulted only when nothing is installed.
  Mutation B (restore `continue` on a missing `./compat`): `TestCompatExportPathPrefersInstalledAndRejectsIncompatible` fails at
  `pilive_test.go:161`, proving the global copy was masking the incompatible
  installed one.
- DONE: Added the negative settings fixture: seed settings holding an object package entry, merge, and require that entry to survive.
  `TestMergeSettingsPreservesEntriesAndIsIdempotent` seeds
  `{"theme":"dark","packages":["npm:other",{"source":"npm:other-object"}]}`, then asserts the written file still contains
  `"source": "npm:other-object"`.
- DONE: Added the negative compatibility fixture: installed copy present, no usable `./compat`, working global copy present, resolution must fail.
  `TestCompatExportPathPrefersInstalledAndRejectsIncompatible` also positively locks
  that the installed copy wins when valid and that the global copy is used only when
  nothing is installed.
- DONE: Strengthened the idempotence assertion so it fails on duplicate registrations, not just on a returned error.
  The second merge is now followed by a byte comparison of the settings file before
  and after (`data` vs `again`). Mutation D (make the append unconditional):
  `TestMergeSettingsPreservesEntriesAndIsIdempotent` fails with the duplicated
  packages array; the old error-only assertion would have passed it.
- DONE: Strengthened the install test so deleting a check or mis-wiring an npm call fails it.
  (a) Exact-argument assertions: the fake now records `npm install` args and the
  test requires exactly two installs — `npm install -g <pinned pi-coding-agent
  tarball>` and `npm install --prefix <agentDir>/npm <pinned pi-subagents tarball>
  <pinned pi-intercom tarball>` — so a dropped install, a wrong tarball, or a
  wrong prefix fails. (b) `npm pack` now errors on any spec/version other than the
  pin. (c) Three negatives re-run `install` against a broken tree and require
  failure: installed version not the pin, a packaged skill file removed, and a
  declared runtime file removed. Mutation C (delete the packaged-skill existence
  check): the test fails at `pilive_test.go:233` (`install must fail when a
  packaged skill file is missing`); the old positive-only test passed it.
- DONE: Froze the call surface and held the port at its accepted behaviour set.
  `Command` is untouched; the pre-fix and fixed binaries print identical `pins` /
  `print-install` and keep `verify-manifest` exit codes. The port is now 669 lines
  (427 impl + 242 tests) versus the accepted 597 (412 + 185): +15 impl for the two
  fix bodies and +57 tests for the required negative fixtures and strengthened
  assertions. No behaviour was cut and no size-driven rewrite was done; the growth
  is exactly the negatives the review asked to add.
- DONE: Bounded validation, no repository-wide suite and no CI lane run.
  `go build ./...`, `go vet ./internal/pilive`, `gofmt -l internal/pilive` (empty),
  `git diff --check`, and `go test ./internal/pilive/ -count=1` all pass; the four
  mutation probes above were run and reverted. `go test ./...` and any race/live/
  CI lane were deliberately not run per the dispatch.

### Why each touched test can now fail

- `TestMergeSettingsPreservesEntriesAndIsIdempotent` (renamed from
  `TestMergeSettingsCreateOrMergeRejectsMalformed`): fails when an existing object
  `packages` entry is dropped (Mutation A) and when a second merge duplicates a
  registration (Mutation D). The old body only checked that the second call returned
  no error and seeded no object entry.
- `TestCompatExportPathPrefersInstalledAndRejectsIncompatible` (new): fails when an
  installed pi-ai without `./compat` silently resolves to the global copy
  (Mutation B), when the installed copy does not win over the global one, or when
  the global fallback is used with nothing installed.
- `TestInstallVerifiesAndRegistersThePinnedFamily`: fails when an install call is
  dropped, when a tarball/prefix is not the pinned one, when any `npm pack` spec is
  not the pin, and when the installed-version, manifest, or packaged-skill check is
  removed (Mutation C). The old body pre-created every valid file, accepted any npm
  arguments, and asserted no failure.

### Summary

Fixed the three reviewed faults at 669 lines: `mergeSettings` preserves object
`packages` entries and dedupes registrations, `compatExportPath` treats the installed
pi-ai copy as authoritative and fails on an invalid `./compat` instead of masking it
with a global copy, and the settings idempotence plus install tests are now
falsifiable (four mutation probes each fail the intended test). Call surface and all
other behaviours unchanged; focused build/vet/gofmt/tests green; no repository-wide
suite or CI lane was started.

## Implementation cycle 12 — the confounded version assertion

- DONE: De-confounded the installed-version rejection in
  `TestInstallVerifiesAndRegistersThePinnedFamily`. The fixture previously replaced the
  package manifest with only a name and a version, which also removed `pi.extensions`, so
  `verifyManifest` rejected it for its own reason and the assertion passed even with the
  version check deleted.
  Evidence: the fixture now changes only the identity inside an otherwise valid manifest
  (`pilive_test.go:229-232`) and requires the version-mismatch diagnostic.
- DONE: Proved the assertion can fail. Removing the version check fails the test at
  `pilive_test.go:231` with `install must fail with the version-mismatch diagnostic, got
  <nil>`; the source was then restored byte-identical to HEAD.
- DONE: Added the wrong-name case the same way, since the cycle-10 report claims name
  enforcement is proven. Removing the name check fails at `pilive_test.go:235`.
- DONE: Test-only change, `5a2ad16e9`, `internal/pilive/pilive_test.go` +16/-7, pushed.
- SKIPPED: No state report was written by the worker; the first officer records it here.

## First-officer gate record — frozen tip `5a2ad16e9`

- The lane run `37189752489` is in flight at this commit and carries `go test ./...` on a
  runner with disk. One duplicate run at the same tip and one at the superseded
  `52223fe29` were cancelled.
- Local offline suite: **one failure**, `TestCodexProcessRecognizesTerminalTurnBeforeOSExit`
  in `internal/ensigncycle` — the 250 ms no-progress quiet budget under machine pressure.
  The test passes 3 of 3 in about 0.5 s when run alone; in the failing run its package took
  447 s, with 289 MiB free on the disk. This is machine thrash, not a fault at the tip, and
  the local suite is not trustworthy evidence on this machine in this condition.
- Race on the packages this layer touches — `internal/pilive`, `internal/cli`,
  `internal/release`: **PASS**. The full race suite is unrun because a race build of every
  package does not fit in the free space.
- `gofmt`: this layer's files are clean. One file is unformatted **on `main`**
  (`internal/release/runtime_live_evidence_workflow_test.go`), untouched by this layer, and
  no workflow runs a format gate.

## Stage Report: validation (frozen tip 5a2ad16e9)

Validation of the frozen candidate `5a2ad16e9` against `/tmp/mh-validation-evidence.md`, the
exception register at `/tmp/pi-art/pi-coverage-detail.jsonl`, and the code in the worktree.
Criteria were read from this entity (AC-1..AC-10). No repository-wide suite, race suite, or CI
run was started.

- DONE: Lane assertions reachable from the register: the run passed, 17 journeys executed, 5 XFAIL engagements, 2 XPASS alerts.
  `/tmp/pi-art/pi-coverage-detail.jsonl` carries 18 `pass` actions (17 `TestLiveCommon*` journeys + package `PASS`), exactly 5 `XFAIL` lines and 2 `XPASS ALERT` lines; `/tmp/pi-art/pi-front-door-smoke-detail.jsonl` shows `TestLivePiFrontDoorSmoke` PASS. Run `37189752489` is `completed success` at `5a2ad16e9` per the pack.
- DONE: Bound-journey reason check (point 1) — 4 of the 5 XFAIL engagements failed for the fault their binding names.
  `owned-conflict-owner-handoff` observed=[conflict-owner-handoff-violation] matches owner `fe7bfjz9sb8wyckmnnm3ncjx` (repair-pi-owner-handoff); `auto-continue-after-implementation` single-root and split-root observed=[validation-worker-not-dispatched] match owner `mk72bnt1b5hsp9sfv83979xs` (repair-pi-worker-lifecycle-observation); `rejection-flow` observed=[rejection-worker-topology] matches the active owner `6h3teccccn3qh71yqcmjbjx4` (own-pi-rejection-worker-topology).
- DONE: Bound-journey reason check — the exception: `pi/default-headless-gate-stop` did NOT fail for the reason its binding names.
  Register line: owner=`gcmfwfjd9735b58sbzw7xsb8` observed=[implementation-worker-not-dispatched]. The binding's own comment (`internal/ensigncycle/shared_live_runner_test.go:138-147`) says it "covers the missing-prepare reds" (`gate-hold-violation`/`gate-not-held`) and that `implementation-worker-not-dispatched` "is noted, not bound here" (owned by `mk72bnt1b5hsp9sfv83979xs`); the owner body `repair-pi-recorded-gate-lifecycle.md:112` records the same observed code and that the missing-reference fault did not reproduce. `gradeLive` only checks that an XFAIL failed, never that the observed codes match the bound owner, so this binding engaged for an unbound fault.
- DONE: Stale bindings (point 2) — confirmed for both, journeys named exactly.
  `pi/smallest-sufficient-mechanism` (owner `h30c9jrfcf21fdh2qs5z58sd`) and `pi/keep-moving-posture` (owner `x02375wsg6q61xek7p0t36j2`) each emitted `XPASS ALERT … observed=[]` and each test ended `--- PASS`, so the retained XFAIL bindings are stale. The smallest-sufficient-mechanism comment's claim that it "can never XPASS" (`shared_live_runner_test.go:155-157`) is contradicted by the run.
- DONE: The eight internal/pilive behaviours are present and wired into the workflow.
  install (`install`), guard (`guard`), verify-manifest (`verifyManifest` + `Command`), integrity check (`pack` rejects `entries[0].Integrity != p.integrity`), manifest resolution (`verifyManifest` requires nonempty `pi.extensions` and `exports["./intercom-bridge"].default` resolving to regular files), settings create-or-merge (`mergeSettings` keeps existing entries incl. object entries and is idempotent), floor guard (`versionAtLeast` + Node/Pi/pi-subagents floors), compatibility guard (`compatExportPath` prefers the installed copy and fails on an invalid `./compat`; `compatExportLoads`). The workflow calls only `spacedock-release install|guard|verify-manifest`; `go test ./internal/pilive/... ./cmd/spacedock-release/... -count=1` PASS.
- DONE: AC-1 — NOT MET.
  `/tmp/pi-art/live-artifacts/pi/pi-doctor.txt` logs `OK pi version: 1.0.2 (floor 0.83.0)`, but AC-1 names the baseline `@earendil-works/pi-coding-agent@1.0.0` (integrity `sha512-/FtbxoSQU…`); the shipped pin is `1.0.2` (`sha512-3ZdIghMS…`). AC-1's proof "the release Go guard compares every pair to the registry snapshot" has no implementation — no registry oracle or pin guard remains (grep finds none; `internal/pilive/registry_oracle_live_test.go` and `internal/release/pi_live_pins_guard_test.go` were deleted), so the named falsifier `revert-agent-pin` fails no test. The lane is green and doctor passes, but the installed-version condition is false and the stated proof owner is gone.
- DONE: AC-2 — NOT MET as written.
  Behaviour is implemented (install calls `verifyManifest`; the setup step calls `spacedock-release verify-manifest`), but AC-2's proof "Go structural guard binds both checkpoint blocks" and falsifier "restore … the line-777 extension assertion" refer to a guard deleted in `2c061b2e9` and a line that no longer exists. At the tip the only demonstration is the same-change `internal/pilive/pilive_test.go`.
- DONE: AC-3 — NOT MET as written.
  Integrity enforcement exists and is exercised (`TestPackRejectsIntegrityMismatch`), but AC-3's proof names a Go structural check binding each spec/integrity to `verified_pack`, and its falsifier `corrupt-subagents-integrity` names `PI_SUBAGENTS_INTEGRITY` — all removed. Only demonstration at the tip is the same-change pilive test.
- DONE: AC-4 — NOT MET.
  The pi-live job comments no longer identify the "Pi 1.0 family" or preserve a separate `0.83.0` floor explanation, and `docs/runtime-live-ci.md` no longer names the three versions — it runs `eval "$(go run ./cmd/spacedock-release print-install)"`. AC-4 requires the compatibility comment plus two doc install commands naming the same three versions.
- DONE: AC-5 — MET (residual note).
  Pack: race `go test ./... -race -count=1` 21/21 packages ok; offline CI job success; `gofmt` clean for this layer. Boundary audit: `git diff c5ca95e74 HEAD -- .github/workflows/runtime-live-e2e.yml` touches only the pi-live job; `piVersionFloor = pilive.PiCodingAgentFloor` keeps the value `0.83.0`; no other lane's pins, action majors, or live selectors change. Residual: the pack records one local offline failure `TestCodexProcessRecognizesTerminalTurnBeforeOSExit` on a real 250 ms budget in a test this layer does not touch.
- DONE: AC-6 — MET.
  `/tmp/pi-art/live-artifacts/pi/pi-frontdoor-smoke/run/pi-ensign-boot-grade.json`: verdict `pass`, `boot_contract` true, skills `["ensign"]`, `isolated_discovery.tools` contains both `subagent` and `intercom`, `package_root_env_absent` true; the workflow no longer exports either override. This is a live model run, not a same-change unit test.
- DONE: AC-7 — MET (same-change tests).
  `TestPiDefaultExtensionRootsReadsRealInstalledLocation` and `TestPiIsolatedHomeRegistersBothSubstratesAndAbsoluteSpacedock` (with `TestPiIsolatedHomeNegativeControlDropsSubstrateRegistrations`) pass; the resolver matches candidates by `package.json` name and resolves `<agentDir>/npm/node_modules/<name>`. Caveat: the dedicated local-source-by-name test was deleted in cycle 8, so that clause rests on the code, and every remaining test was authored by this change.
- DONE: AC-8 — MET (same-change tests).
  `TestPiLiveEnvDefaultScrubsPackageRoots`, `TestPiLiveEnvHonorsIndependentOverrides`, and `TestPiLiveEnvDropsForeignRuntimeMarkers` exist and pass (`go test ./internal/ensigncycle -run 'TestPiDefaultExtensionRoots|TestPiIsolatedHome|TestPiLiveEnv' -count=1` PASS).
- DONE: AC-9 — MET (non-live half + the AC-6 smoke).
  `internal/ensigncycle/pi_default_extensions_test.go` carries no `//go:build` constraint; the helper/controls tests pass; the front-door smoke passed; the pack records gofmt clean and live-tagged vet/build pass.
- DONE: AC-10 — MET on clauses (a)/(b) with a residual determinism gap.
  `drainCodexToTerminal` now uses `w.now()`/`w.sleep()`, and `TestCodexProcessActivityResetsQuietBudget` drives a fake clock + fake line source with no real sleeps and no duration assertion; the `stall` kill-path test is unchanged and the three host tests build isolated homes. Residual: the pack's own recorded local offline failure (`TestCodexProcessRecognizesTerminalTurnBeforeOSExit`, untouched, real 250 ms budget) means the offline gate is not yet fully machine-independent.
- DONE: Same-change-only demonstrations flagged.
  AC-2, AC-3, AC-7, AC-8, and AC-10 are demonstrated solely by tests authored in this change; AC-6 is the only value criterion with an independent live artifact. AC-1's independent artifacts (green lane, doctor) plus its own version condition disagree with its text.
- DONE: Evidence-pack inaccuracies.
  The pack's "Port … `pilive_test.go` 242 = 669 lines" is stale: at the frozen tip `pilive.go` is 427 and `pilive_test.go` is 251 (sum 678), and the layer is +1367/-289 (net +1078) across 19 files vs base `c5ca95e74`, not the pack's "net +1069". The pack's criteria note also omits AC-10, which the entity declares.
- SKIPPED: Repository-wide `go test ./...`, `go test ./... -race`, and any CI lane run.
  Explicit FO budget prohibition; those results are taken from the pack only. No prohibited command was run.

### Summary

The frozen candidate is green and all eight port behaviours are present, but the entity's AC-1,
AC-2, AC-3, and AC-4 are not met as written: AC-1 installs 1.0.2 against a criterion naming 1.0.0
and whose registry-comparison guard was deleted; AC-2/AC-3 describe proof owners and falsifiers
that no longer exist; AC-4 requires workflow/doc prose that was deliberately removed. AC-5 through
AC-10 are met on their named requirements, with the caveats that several rest on same-change tests
and that the recorded local offline flake leaves AC-10's determinism not absolute. Recommendation:
return the four stale criteria to the owner for revision (or a captain decision), keep the
mechanism, and do not treat the green lane as evidence for criteria the evidence cannot reach.


## Stage Report: implementation (cycle 3)

- DONE: Change the pi-coding-agent readiness floor to 1.0.0. It is PiCodingAgentFloor in internal/pilive/pilive.go. Change that constant and nothing else about the value.
  `PiCodingAgentFloor = "1.0.0"` (was `"0.83.0"`); the launcher gate still reads it through `const piVersionFloor = pilive.PiCodingAgentFloor` (internal/cli/pi.go:44). Code commit `5ead9b85c` on `spacedock-ensign/pi-live-lane-pin-refresh` (rebased onto the remote tip `eef9bef49`; the pre-rebase SHA was `665835f87`).
- DONE: Write no version string in any comment. The constant stays the only place the number appears. Do not add a comment naming the old floor or the new one.
  The const-block comment is unchanged and still reads only "readiness floors (pi-coding-agent's is also the launcher ready-gate floor): the ONLY place these numbers may be written" (internal/pilive/pilive.go:19-20); `git diff` shows no comment line touched. The new number is written only at the constant.
- DONE: Keep every other constant as it is, including the pi-subagents floor and the node floor.
  `PiSubagentsFloor = "0.53.0"` and `NodeEngineFloor = "22.19.0"` are byte-identical, as are every Spec/Version/Integrity pin; the diff is one constant line plus two test literals.
- DONE: Update any test that assumes the old value, but prefer the existing constant over a new literal. A test whose only purpose is to name a version must use the constant or a value clearly above or below the floor.
  Two cases in `internal/cli/pi_frontdoor_test.go` `TestPiVersionAtLeast` asserted true for a version just above the OLD floor and now fall below 1.0.0; each moved to a value clearly above the new floor, none naming the floor as a new literal. Cases already at/above or below the new floor were left untouched: `PiCodingAgentFloor`, `1.0.2`, `piTestVersionAboveFloor` (99.0.0), `0.82.9` (pilive + cli), `0.73.1`, `0.8.9`, `0.9.0`.
- DONE: Report which tests you touched and why each still fails when its behaviour is removed.
  - `{"1.0.1", true}` (was `{"0.83.1", true}`): falsified if `piVersionAtLeast` stops treating a version above the floor as ready (e.g. compares only major/minor or uses strict `<`), because 1.0.1 is one patch above 1.0.0.
  - `{"garbage 1.2.3 trailing", true}` (was `{"garbage 0.90.0 trailing", true}`): falsified if the semver regexp stops extracting the first embedded triple or the parser fails closed on non-bare output.
  Both remain true falsifiers of the same behaviour they had before; only the numeric value was lifted above the new floor.

### Summary

Raised `PiCodingAgentFloor` to `1.0.0` in `internal/pilive/pilive.go` as a one-constant change; no comment names any version and no other pin/floor moved. Under the new floor, two `TestPiVersionAtLeast` cases in `internal/cli/pi_frontdoor_test.go` that asserted true for a version just above the OLD floor sat below 1.0.0: `{"0.83.1", true}` was lifted to `{"1.0.1", true}` by inspection and the focused run then caught `{"garbage 0.90.0 trailing", true}` still failing, which was lifted to `{"garbage 1.2.3 trailing", true}`. Focused `go test ./internal/cli/... ./internal/pilive/...`, `go build ./...`, and `gofmt -l` pass. Code commit `5ead9b85c` (pre-rebase `665835f87`).

## Stage Report: validation (successor — readiness floor `1.0.0` at `5ead9b85c`)

Successor validation of tip `5ead9b85c`, which adds exactly one change over the previously
validated tip: `PiCodingAgentFloor` `0.83.0` -> `1.0.0`. Lane evidence is run `37338782614
at chain top `238a2d59a` (a descendant of `5ead9b85c`). Run state checked once: the run is
`in_progress` — `offline` job `completed/success`, `pi-live` job `in_progress`,
`claude-live`/`codex-live` skipped. No repository-wide suite, race suite, or CI run was
started by this stage; results for lane-dependent criteria are therefore reported as
awaiting the run.

- DONE: The floor constant is `1.0.0`.
  `internal/pilive/pilive.go:32` reads `PiCodingAgentFloor = "1.0.0"`; it is the launcher
  ready-gate through `internal/cli/pi.go:44` `const piVersionFloor = pilive.PiCodingAgentFloor`.
- DONE: No comment carries a version string.
  `grep -nE '//.*[0-9]+\.[0-9]+\.[0-9]+' internal/pilive/ internal/cli/` returns nothing; the
  const-block comment (`pilive.go:19-20`) and the floor comment (`pi.go:39-43`) name no number;
  `git diff 5ead9b85c^ 5ead9b85c` changes no comment line.
- DONE: The pi-subagents floor and the node floor are unchanged.
  `PiSubagentsFloor = "0.53.0"` and `NodeEngineFloor = "22.19.0"` are byte-identical; the floor
  commit's only logic line is the `PiCodingAgentFloor` value.
- DONE: Tests that assumed the old floor were corrected, not deleted.
  In `internal/cli/pi_frontdoor_test.go` `TestPiVersionAtLeast`, `{"0.83.1", true}` ->
  `{"1.0.1", true}` and `{"garbage 0.90.0 trailing", true}` -> `{"garbage 1.2.3 trailing", true}`;
  the case count is unchanged (no cases removed), and the remaining cases use
  `pilive.PiCodingAgentFloor` or values clearly above/below. No literal `0.83.0`/`0.83.1`/`0.90.0`
  remains in code, tests, docs, or workflow.
- DONE: Focused checks pass.
  `go build ./...`, `go test ./internal/pilive/...`, `go test ./internal/cli -run
  'TestPiVersionAtLeast|TestPiFrontDoor|TestPiVersionFloor|TestPiDoctor'`,
  `go test ./cmd/spacedock-release/...`, `go vet ./internal/pilive ./cmd/spacedock-release`, and
  `gofmt -l` on the changed files (empty) all pass.
- DONE: Amended AC-4 holds at this tip and needs no lane result.
  The three versions are declared once in `internal/pilive/pilive.go`; `go run
  ./cmd/spacedock-release print-install` emits exactly `@earendil-works/pi-coding-agent@1.0.2`,
  `pi-subagents@0.75.0`, `pi-intercom@0.16.0`; `docs/runtime-live-ci.md:95` delegates with
  `eval "$(go run ./cmd/spacedock-release print-install)"`. The only non-declaration occurrences
  of the pin numbers are a synthetic comparison-test literal and non-consumed roadmap prose.
- DONE: Amended AC-1 code half holds; the installed-name/version condition is unaffected by the
  floor raise.
  Pins are `1.0.2` / `0.75.0` / `0.16.0`; the lane's `guard` requires the installed
  pi-coding-agent version >= `PiCodingAgentFloor` (`1.0.0`), and the pinned `1.0.2` satisfies it.
- DONE: Amended AC-2 code half holds; unaffected by the floor change.
  `verifyManifest` requires nonempty `pi.extensions` and `exports["./intercom-bridge"].default`
  resolving to regular files; `install` calls it for the installed subagents root and the
  "Verify Pi current-checkout setup" step calls `spacedock-release verify-manifest`.
- DONE: Amended AC-3 code half holds; unaffected by the floor change.
  `pack` rejects any `npm pack --json` entry whose integrity differs from the pin before either
  tarball install; covered by `internal/pilive` tests that pass.
- SKIPPED: Amended AC-1 lane half (green pi-live run on the pinned candidate).
  Run `37338782614` is `in_progress` with its `pi-live` job running; the installed-version logs,
  doctor output, and journey/smoke grades are not yet available. Awaiting that run.
- SKIPPED: Amended AC-2 lane half ("the live lane's green run at the frozen tip").
  Same run; the pi-live job has not completed. Awaiting that run.
- SKIPPED: Amended AC-3 lane half ("the live lane's install of the pinned tarballs").
  Same run; the pi-live job has not completed. Awaiting that run.
- DONE: The amended criteria's recorded lost guarantees are still absent, so no criterion is
  silently strengthened.
  `grep` for `registry_oracle`, `pi_live_pins_guard`, `verified_pack`, and
  `PI_SUBAGENTS_INTEGRITY` finds no match in Go or workflow files; AC-1..AC-3 rest only on the
  pins, the port's commands/tests, and the pending lane run, exactly as amended.
- SKIPPED: AC-5's suite clauses could not be re-run here.
  FO budget forbids the repository-wide and race suites. The run's `offline` job ran
  `go test ./...` and is green at `238a2d59a`; the race clause is not evidenced by this run's
  visible jobs.
- OPEN RISK: AC-5's text says the reviewed diff "leaves ... the Pi floor ... unchanged", but the
  captain-directed successor change raises that floor. The floor raise is authorized, so AC-5's
  Pi-floor clause is stale against it and needs an FO/captain criteria amendment; the other
  boundary clauses (other lanes' pins/action majors/live selectors) still hold, since the floor
  commit touches only `internal/pilive/pilive.go` and `internal/cli/pi_frontdoor_test.go`.

### Summary

Settled the floor change against the code: the constant is `1.0.0`, no comment names a version,
`PiSubagentsFloor`/`NodeEngineFloor` are unchanged, the two old-floor test cases were corrected
in place rather than deleted, and focused build/vet/gofmt/tests pass. Amended AC-4 holds
outright; amended AC-1/AC-2/AC-3 hold on their code halves but their lane halves await run
`37338782614`, which is still `in_progress`. Flagged that AC-5's "Pi floor ... unchanged" clause
is contradicted by the authorized floor raise and needs a criteria amendment.
