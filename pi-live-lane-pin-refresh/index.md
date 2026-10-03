---
title: Refresh the pi-live lane pins and substrate assertions for the Pi 1.0 family
status: ideation
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
                state: pending
sprint: pi-ux
group: tooling
sprint-readiness: ready
started: 2026-10-03T04:04:56Z
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
