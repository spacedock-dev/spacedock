---
title: Refresh the pi-live lane pins and substrate assertions for the Pi 1.0 family
status: implementation
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
sprint: pi-ux
group: tooling
sprint-readiness: ready
started: 2026-10-03T04:04:56Z
worktree: .worktrees/spacedock-ensign-pi-live-lane-pin-refresh
pr: "#816"
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

**AC-1 (value) — pi-live installs and validates the published Pi 1.0 family.**
The three installed package names/versions equal the independently observed published
baseline `@earendil-works/pi-coding-agent@1.0.0`, `pi-subagents@0.75.0`, and
`pi-intercom@0.16.0`; the existing pi-live Runtime Live E2E lane completes green on
that candidate, with its existing common journeys and front-door smoke still selected.
Proof: the release Go guard compares every pair to the registry snapshot, and a
recorded CI run URL/SHA supplies installed-version logs, doctor output, and durable
journey/smoke results under the existing grading policy. No new skips/XFAILs count.
Named falsifying edit **revert-agent-pin**: set `PI_CODING_AGENT_VERSION` back to
`0.85.1`; the structural guard fails even if the older runtime still passes journeys.

**AC-2 — all four substrate assertions follow installed runtime declarations.**
Both setup checkpoints resolve every `pi.extensions` entry and the intercom bridge
runtime export from the installed manifest, fail for missing declarations/files,
and no longer require either old source path. Proof: Go structural guard binds both
checkpoint blocks; one-off real-package exercise checks success and failure with
an extension/bridge target removed in turn. Named falsifying edit
**restore-setup-source-assertion**: restore the line-777 extension assertion while
leaving the install check correct; the guard and published-package setup fail.
Additional negative **remove-bridge-target**: hide the manifest-resolved bridge file;
the Node check must exit nonzero, proving existence checks are load-bearing.

**AC-3 — corrupted integrity prevents installation.**
All three exact integrity pins are verified before their tarballs are installed.
Proof: Go structural checks bind each package spec and expected-integrity argument
to `verified_pack`; the existing function is exercised once with the real published
subagents package and a bad expected hash. Named falsifying edit
**corrupt-subagents-integrity**: replace `PI_SUBAGENTS_INTEGRITY` with an invalid
hash; both the independent pin guard and pack step fail before install. Removing
the mismatch exit is separately falsified by that negative pack exercise.

**AC-4 — current-family instructions agree with executable pins.**
The workflow compatibility comment identifies the Pi 1.0 family, preserves the
separate `0.83.0` floor explanation, and the two local-install commands in
`docs/runtime-live-ci.md` name the same three versions as the verified pins.
Proof: one-off comparison of those commands/comment against AC-1's registry snapshot
and workflow values, recorded during validation; no standing instruction prose test.
Named falsifying edit **restore-stale-doc-version**: restore `pi-subagents@0.35.1`
in the doc command; the independent agreement check must report a mismatch.
Historical floor comments are not stale pins and must not be indiscriminately removed.

**AC-5 (no regression) — the bounded change preserves existing non-Pi behavior.**
`go test ./internal/release/...`, `go test ./...`, `go test ./... -race`, and
`go build ./...` pass; the reviewed diff leaves all other lanes' CLI pins/action
majors, the Pi floor, and live selectors unchanged. Proof: existing release tests
plus a candidate-versus-base diff audit of these boundaries. Named falsifying edit
**downgrade-other-lane-checkout**: change claude-live's `actions/checkout@v5` to `@v4`;
`TestNode24ActionsPinnedAtMinimum` and the boundary audit fail. Any unauthorized
other-lane version change must also fail the diff audit, even if still above a floor.

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
