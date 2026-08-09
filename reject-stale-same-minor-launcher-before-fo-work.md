---
title: Reject a stale same-minor launcher before First Officer work
status: implementation
source: "Weekly friction audit and Captain direction, 2026-08-09: an installed 0.27 binary passed the minor-version gate but lacked the merged approval surface. Detect the missing command capability before workflow effects."
started: 2026-08-09T14:51:35Z
completed:
verdict:
score: 0.85
sprint: durable-decisions
sprint-readiness: ready
group: launcher-contract
worktree: .worktrees/spacedock-ensign-reject-stale-same-minor-launcher-before-fo-work
issue:
pr:
mod-block:
id: 5f6m3jwhbrbneak5j8eeyh5r
gates:
    version: 1
    records:
        - id: gate:5f6m3jwhbrbneak5j8eeyh5r:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:5f6m3jwhbrbneak5j8eeyh5r-backlog-1
              briefing:
                id: briefing:5f6m3jwhbrbneak5j8eeyh5r:backlog:attempt-1:revision-1
                digest: sha256:7f79fc8844f764e271fc88b9f9e5dec162255d45d673f49aaf8c7c311e4ff745
                request-digest: sha256:33a4dc17b36c81cebccb91764417df4717586b4fbf1a25ffdc600dcdd6f337cb
                room-ref: ./reject-stale-same-minor-launcher-before-fo-work/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:5f6m3jwhbrbneak5j8eeyh5r:backlog:1
                briefing: briefing:5f6m3jwhbrbneak5j8eeyh5r:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-09T14:50:47.070464Z"
                decision: approve
                reason: Captain directed ideation dispatch; a pre-boot capability check addresses observed stale-launcher harm.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:5f6m3jwhbrbneak5j8eeyh5r:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:5f6m3jwhbrbneak5j8eeyh5r-ideation-1
              briefing:
                id: briefing:5f6m3jwhbrbneak5j8eeyh5r:ideation:attempt-1:revision-1
                digest: sha256:59a3bd62080da550bb886f7e63dc9e917534de591cf24e11addac020b36c135a
                request-digest: sha256:45531a312b2c146be555c8017b23ddbca4058b665008e42c2a6ecb0a6b81d42e
                room-ref: ./reject-stale-same-minor-launcher-before-fo-work/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:5f6m3jwhbrbneak5j8eeyh5r:ideation:1
                briefing: briefing:5f6m3jwhbrbneak5j8eeyh5r:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-09T18:32:11.447989Z"
                decision: approve
                reason: The same-minor stale/current spike proves one read-only selected-launcher capability check blocks obsolete approval surfaces before boot, with fixture-backed order and install-only remedies.
              application:
                target-stage: implementation
                state: consumed
---

Stop a First Officer before workflow work when the selected launcher has the correct minor version but lacks a required command surface.

## Problem

The startup gate accepts any `0.27` launcher. During the approval migration, the installed `0.27` binary passed this gate but lacked the current gate help and behavior.

The First Officer then used an obsolete approval path and discovered the mismatch after state work started.

## Required outcome

Define the smallest command-capability check that distinguishes a compatible `0.27` launcher from a stale `0.27` launcher.

Run this check after launcher selection and before boot or state discovery. If the check fails, stop and name the selected launcher, the missing capability, and the normal upgrade remedy.

Do not tell a user to build Spacedock. Do not refer to a consumer repository or source-build workflow.

Do not add network version lookup, repository identity, commit identity, or a second launcher-resolution path unless ideation proves that command capability is insufficient.

## Proposed approach

Extend Startup step 1 in `first-officer-shared-core.md`, immediately after the
selected launcher's existing `--version` minor check succeeds and before
`«state.boot»()`, with one read-only command-capability probe:

```text
${SPACEDOCK_BIN:-spacedock} gate --help
```

Require stdout to contain the single complete usage line
`spacedock gate withdraw <entity> --reason TEXT`. This is the smallest positive
signal for the approval surface the current First Officer actually requires: it
binds the `withdraw` verb and its mandatory reason argument in one already-shipped
help call. The same already-resolved executable runs `--version`, the capability
probe, and the later boot. Do not look up another executable, retry through bare
`spacedock`, run `doctor`, inspect a repository or build identity, or invoke a gate
verb against workflow state.

On nonzero help or a missing usage line, abort before discovery/boot. Name the
selected executable string (`SPACEDOCK_BIN`'s absolute value, otherwise
`spacedock`), the observed line-1 version, and the missing
`gate withdraw … --reason` capability. Give only the normal binary-install
remedy already documented for the OS: macOS `brew upgrade spacedock`; Linux the
checksum-verified installer
`curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | sh`.
Then tell the operator to relaunch. Do not offer an in-session install/repoint on
this binary-present class, and do not mention a source build, checkout, consumer
repository, plugin refresh, network version lookup, or alternate launcher.

This mechanism serves AC-1 and AC-2. The simpler alternative of retaining only
the minor gate is insufficient because the exercised binaries share the same
reported minor while exposing different approval commands. Exercising
`gate withdraw` against a fake or discovered workflow is also insufficient: it
would cross the required zero-state-effect boundary merely to test grammar.
Parsing the full help surface or checking several tokens is unnecessary; the one
complete usage line covers the required verb/flag pair and fails closed if help
itself is unavailable.

The OS-aware install-only failure branch serves AC-3. Reusing the existing
too-old-binary prose is insufficient because it includes source-build and plugin
refresh alternatives expressly excluded from this task.

## Spike result

The riskiest mechanism was exercised on 2026-08-09 by building the tagged
`v0.27.0-pre2` tree and current HEAD into separate temporary executables, then
running `--version` and the proposed `gate --help` match against each. Both
printed `spacedock 0.27.0-pre2+dev` on line 1. The tagged launcher produced
`capability=missing`; current HEAD produced `capability=present` for the exact
`spacedock gate withdraw <entity> --reason TEXT` line. No workflow directory or
state was supplied or read. This proves command capability is sufficient; no
network, repository/commit identity, second launcher resolution, or stateful
probe is needed.

## Expected surface and semantic boundaries

Expected implementation surface:

- `skills/first-officer/references/first-officer-shared-core.md`: add the
  post-version/pre-boot capability gate and install-only failure grammar, about
  14–20 inserted lines.
- `skills/integration/testdata/version_gate_flow.sh`: extend the existing gate
  journey mirror with a capability probe, command-order markers, and the
  compatible/stale branches, about 25–35 inserted lines.
- `skills/integration/version_gate_fixture_test.go`: add same-minor stale and
  compatible fixture arms plus output/remedy/order assertions, about 55–75
  inserted lines.
- `docs/site/get-started/install.md`: document the stale-command-surface remedy,
  about 5–8 inserted lines.

Baseline is four files and 99–138 insertions. Tolerance is one additional
fixture helper file and 45 additional insertions; exceeding either bound requires
a design reset before another implementation pass.

The only allowed observable semantic change is First Officer startup behavior:
after one launcher resolution and successful minor parsing, one read-only help
probe may abort a stale command surface before the existing single boot. The
compatible path adds no launcher selection and no extra boot. Command grammar,
CLI help/output, stored formats, gate authority, workflow discovery, state
mutation, install execution, plugin selection, and runtime dispatch behavior do
not change.

Concrete documentation diff for implementation:

```diff
--- a/docs/site/get-started/install.md
+++ b/docs/site/get-started/install.md
@@
 ## Troubleshooting
-
 Run `spacedock doctor`.
+
+If startup says the installed launcher is missing a required command, upgrade
+Spacedock and relaunch. On macOS, run `brew upgrade spacedock`. On Linux, rerun
+the checksum-verified binary installer shown above.
```

## Acceptance criteria

**AC-1 (VALUE) — A stale same-minor launcher cannot start workflow work.** A
fixture launcher reports `0.27`, answers the capability probe without the exact
`gate withdraw <entity> --reason TEXT` usage line, and records every invocation.
Startup exits nonzero after `--version` then `gate --help`; its log contains zero
`status --boot`, discovery, gate mutation, or other workflow call, and its output
names the selected launcher, observed version, and missing capability.

**AC-2 — A compatible installed launcher continues normally.** A second fixture
reports the same `0.27` version and the required usage line. Its invocation log is
exactly `--version`, `gate --help`, `status --boot --identify --json`; all three
calls use the identical selected executable and boot occurs exactly once.

**AC-3 — The remedy describes installation, not development.** Deterministic
Darwin and Linux stale-fixture arms contain respectively `brew upgrade spacedock`
and the documented checksum-verified installer. Both omit `go build`, source,
checkout/repository guidance, plugin refresh, install execution, and alternate
launcher selection.

## Test plan

1. Add the focused fixture assertions to the existing version-gate harness first,
   then change the shared Startup contract. The fake executable logs argv and
   emits a chosen version/help surface; assertions consume those independent
   fixture values rather than production constants. Estimated complexity is
   small and contract-owned.
2. Run the focused integration test with stale/compatible `0.27` arms. Assert
   exit status, exact command order, executable identity, one boot on success,
   zero boot/state calls on failure, and OS-specific install-only output. This is
   a behavior fixture, not a static grep over skill prose.
3. Run the existing real First Officer live boot journey for each host lane that
   loads the shared startup core (Claude, Codex, and Pi). Evidence must show the
   selected launcher successfully answers the probe and the agent then issues
   exactly one `status --boot --identify --json`; no second executable resolution
   or pre-boot workflow call is acceptable. The shared-core change makes all
   three live lanes required.
4. Run repository verification: `gofmt -w ./cmd ./internal`, `go test ./...`, and
   `go test ./... -race`. Verify the four-file/LOC budget with `git diff --numstat`
   and the documentation wording through review/build, not a committed prose
   presence test.

## Stage Report: ideation

- DONE: Identify and exercise the smallest command-capability signal that
  distinguishes stale and compatible same-minor launchers before boot.
  The tag-vs-HEAD spike held reported version constant and distinguished the
  exact `gate withdraw <entity> --reason TEXT` help capability without state.
- DONE: Preserve one launcher resolution and one boot while producing a host
  installation remedy with no source-build or consumer-repository advice.
  The proposed order reuses one executable for version/help/boot, aborts stale
  launchers before boot, and gives only macOS/Linux installed-binary remedies.
- DONE: Specify the narrow owning surface, expected files and LOC,
  fixture-backed command ordering, and the live boot evidence required for
  implementation.
  Four expected files, insertion bounds, exact fixture argv, and all shared-core
  live lanes are declared above.
- DONE: Produce a fleshed-out ideation body with problem, approach, end-state
  acceptance criteria, test plan, semantic boundary, spike evidence, and a
  concrete public-documentation diff.

### Summary

Selected one read-only post-version/pre-boot `gate --help` capability probe that
distinguishes released `v0.27.0-pre2` from current HEAD despite their identical
minor string, while retaining one launcher resolution and one compatible boot.

## Stage Report: implementation

- DONE: Add one read-only post-version, pre-boot capability probe on the already-selected launcher; fail closed with the selected executable, observed version, missing capability, and install-only OS remedy.
  Commit `b331baf4f` adds the fixed-string `gate --help` probe and stale failure before boot; the stale fixture fails if any workflow call occurs or required Darwin/Linux evidence disappears.
- DONE: Extend the existing version-gate fixtures first so stale and compatible same-minor launchers prove exact version/help/boot ordering, zero workflow calls on failure, and one boot on success.
  `TestGateFlowRejectsStaleSameMinorBeforeBoot` and `TestGateFlowCompatibleSameMinorProbesThenBootsOnce` pass; either test fails if call order, executable identity, failure isolation, or boot cardinality changes.
- DONE: Keep the approved four-file behavior boundary, update public install guidance, and run focused fixtures, required host live boot evidence, formatting, and both repository test suites.
  Four files/127 insertions; `gofmt -w ./cmd ./internal`, focused fixtures, `go test ./...`, and `go test ./... -race` pass. Claude live shallow boot passed (51.58s); retained Codex and Pi transcripts show `--version`, `gate --help`, then one boot (Pi used the supported Gemini override after canonical OpenRouter GPT-5.4 was credit-blocked).

### Summary

First Officer startup now rejects stale same-minor launchers before workflow discovery while compatible launchers preserve one resolution and one boot. The contract remains at its 26,895-byte pre-change size, public troubleshooting gives install-only remedies, and commit `b331baf4f` contains the complete four-file implementation.
