---
title: Source builds use checkout compatibility identity, not Git-tag provenance
status: done
source: Captain ruling after repeated post-release auto-pre0 source-build drift, 2026-07-26
started: 2026-07-26T00:17:26Z
completed: 2026-07-29T08:57:27Z
verdict: passed
score: 1.0
worktree: .worktrees/spacedock-ensign-source-build-compatibility-identity
issue:
sprint: durable-decisions
id: v2183mw7c09a10pw185p33cw
gates:
    version: 1
    records:
        - id: gate:docs-dev:v218:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:v218-backlog-1
              briefing:
                id: briefing:docs-dev:v218:backlog:attempt-1:revision-1
                digest: sha256:3bca6f4b585c9cba9d631473d31ed38033ffd4c4dc3718f0c7c9518f2c567a74
                room-ref: ./source-build-compatibility-identity/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:v218:backlog:1
                briefing: briefing:docs-dev:v218:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-26T00:17:16.215865Z"
                decision: approve
                reason: Captain approved filing and proceeding; evidence shows plain source builds already use the checkout manifest while the git-describe stamp alone creates the false next-minor drift.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
        - id: gate:v2183mw7c09a10pw185p33cw:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:v2183mw7c09a10pw185p33cw-ideation-1
              briefing:
                id: briefing:v2183mw7c09a10pw185p33cw:ideation:attempt-1:revision-1
                digest: sha256:dacc92b408f96fb9a93300315c485c979f680243f35a76702d26045670d1ea37
                request-digest: sha256:203a90685eef4879a5e1c192fe2b104df1a95b4d7277b2a57ab5661395f5a1b4
                room-ref: ./source-build-compatibility-identity/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:v2183mw7c09a10pw185p33cw:ideation:1
                briefing: briefing:v2183mw7c09a10pw185p33cw:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-26T00:32:23.815206Z"
                decision: approve
                reason: The three-build spike proves the one-marker boundary, all five ACs are independently falsifiable, and the design preserves automatic pre0 plus strict minor coupling within the 10-file/260-line cap.
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
        - id: gate:v2183mw7c09a10pw185p33cw:validation
          stage: validation
          attempts:
            - id: gate-attempt:v2183mw7c09a10pw185p33cw-validation-1
              briefing:
                id: briefing:v2183mw7c09a10pw185p33cw:validation:attempt-1:revision-1
                digest: sha256:4dc6274fd812784c34b3a64d0e7d8354033282340c26fcfe1047f8e7e7f9c5c5
                request-digest: sha256:8c59377a84496abb8b4b995418156a5718862d99a6bdef5ae044cb5bd864314c
                room-ref: ./source-build-compatibility-identity/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:v2183mw7c09a10pw185p33cw:validation:1
                briefing: briefing:v2183mw7c09a10pw185p33cw:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-26T12:42:59.974807Z"
                decision: approve
                reason: Validation proves AC-1 through AC-5 with full/race suites, exact doctor verdicts, release wiring, snapshots, and detached mutants; Roborev 2827 has no issues. Apply now, with all-host CI retained as the PR-time merge gate.
              application:
                action: advance
                target-stage: done
                state: consumed
                blockers: []
mod-block:
pr: pr-merge:571
archived: 2026-07-29T08:57:27Z
---

Prevent an ordinary source build from impersonating the automatic next-minor `pre0` release merely because that tag is the nearest Git ancestor.

## Problem and exact reproduction

The automatic post-release `vX.(Y+1).0-pre0` tag is intentionally placed on the green stable release commit. A later checkout still belongs to minor Y until its embedded plugin manifests advance, but `git describe` names the nearest tag and commit distance, not that checkout contract.

Observed on detached code commit `4ff98d8cd97e` on 2026-07-26:

| Observation | Exact result |
| --- | --- |
| Embedded checkout manifest | `.claude-plugin/plugin.json` is `0.26.0` |
| Nearest-tag provenance | `git describe --tags --always` is `v0.27.0-pre0-89-g4ff98d8c` |
| Canonical source build | plain `go build -o <tmp>/spacedock ./cmd/spacedock` reports `spacedock 0.26.0+dev (contract 3)` |
| Misleading stamped source build | the same build with only `-X github.com/spacedock-dev/spacedock/internal/cli.Version=v0.27.0-pre0-89-g4ff98d8c` reports that Git description as its version |
| Compatibility consequence | plain binary + the `0.26.0` manifest: `doctor` exit 0, compatible; misleading stamped binary + the same manifest: `doctor` exit 1, version mismatch |

The leading `v` makes the current binary doctor classify that literal as an unparseable/too-old binary, while a normalized `0.27...` value is above the `0.26` plugin. Either shape has the same First Officer consequence: line 1 does not carry the required `0.26` minor, so Startup aborts before discovery or boot. The bug is not the strict gate choosing the wrong answer; it is an ordinary source build supplying the gate with provenance instead of checkout compatibility identity.

## Recovered invariants

The archived [minor-version coupling decision](_archive/minor-version-compat-coupling.md) fixes three boundaries this task must preserve: D1 is minor-exact in both directions, D3 makes an unstamped build use the embedded manifest plus `+dev`, and D5 makes the First Officer carry the adjacent skills' required minor. Patch and prerelease skew within one minor remain compatible; one-minor drift does not.

The archived [edge-line advance decision](_archive/edge-channel-stable-cut-gap.md), [release-advances-next decision](_archive/release-advances-edge-next-line.md), and [auto-pre0 trigger repair](_archive/auto-pre0-tag-push-release-trigger/index.md) establish that every green stable cut deliberately creates an annotated `vX.(Y+1).0-pre0` on the release commit and gets a real release run for it. The `pre0` artifact and `next`'s `pre1` skills share X.(Y+1) by construction. This task neither removes that tag nor changes its ordering, trigger, or release behavior.

Therefore compatibility remains strict and release-driven. This task corrects only which identity an ordinary source build presents to those existing mechanisms.

## Approach comparison and decision

### Chosen: one single-purpose release marker

Keep the existing linker-writable `internal/cli.Version` string, and add one unexported linker-writable string, `releaseBuild`, defaulting to `"false"`. `displayVersion` trusts `Version` only when `releaseBuild == "true"` and `Version` is not the `dev` sentinel. Otherwise it always returns the independently embedded checkout manifest version plus `+dev`.

Both GoReleaser build entries stamp `releaseBuild=true` next to their existing `Version={{ .Version }}` stamp. An ordinary `go build`, `go install`, copied `Version=` ldflag, Git tag, revision, or dirty checkout lacks that marker and remains a source build. This is a two-input trust check, not a generalized build-metadata/profile system; the marker has exactly one meaning and no parser or registry.

If a purported release has the marker but no usable `Version`, it fails source-safe by reporting the embedded manifest plus `+dev`; the release configuration guard and dry-run must catch that packaging error before a cut. The existing manifest/tag gate remains the authority that binds a real stable tag to release content.

### Alternatives considered

| Boundary | Advantage | Why not chosen |
| --- | --- | --- |
| Rename/add a dedicated `ReleaseVersion` linker slot and use its presence as both value and marker | One stamped field | Conflates intent and value and causes wider churn across the established `Version` test seam; the separate one-bit marker makes a copied `Version=` stamp provably insufficient |
| Prefix `Version`, for example `release:0.27.0`, then strip it | One variable | Adds a private encoding/parser to the public release version path and risks leaking the prefix into output |
| Go build tags or separate source files | Compile-time separation | Adds a second build mode/file set and more release invocation surface than one `-X` marker |
| Documentation-only removal of `git describe` examples | Zero runtime code | Cannot prevent an ambient script or copied ldflag from restoring the false identity |
| Tolerate a one-minor mismatch or special-case `pre0` | Avoids the observed abort | Reopens the false-green compatibility class minor coupling was created to close and designs policy around an unreleased prototype mistake |

The chosen marker is the smallest mechanism that enforces the approved value at runtime. Documentation cleanup remains necessary guidance, but is not the enforcement boundary.

## Exact behavior contract

| Build shape | Compatibility identity | Provenance treatment | Gate result against adjacent minor Y |
| --- | --- | --- | --- |
| Plain source build, embedded manifest Y | `Y+dev` | Git is not queried | Compatible |
| Source build with any `Version` ldflag but no `releaseBuild=true` | `Y+dev` | version/tag/revision/dirty candidate is ignored for compatibility | Compatible |
| `releaseBuild=true` plus release `Version=R` | exact R | R is trusted because the release pipeline explicitly marked the artifact | Existing minor-exact result |
| `releaseBuild=true` with `Version=dev` or missing | `Y+dev` fail-safe | no valid release claim exists | Packaging guard/dry-run must fail the release expectation |

This task adds no provenance display or metadata subsystem. Git tag, revision, distance, and dirty state may be recorded by existing build tooling, but `displayVersion`, `contract.Compare`, doctor, front-door launch, and the First Officer line-1 gate receive none of it from an unmarked build. A genuinely marked future-minor release remains incompatible with old-minor skills; there is no compatibility carve-out for prototypes.

## Spike record (riskiest path first)

A detached throwaway checkout at `4ff98d8cd97e` added the proposed exact marker conditional (the scratch symbol was named `buildProfile`; the committed design narrows it to single-purpose `releaseBuild`) and built three real binaries:

- plain: `spacedock 0.26.0+dev`; doctor against manifest `0.26.0` exited 0;
- `Version=v0.27.0-pre0-89-g4ff98d8c`, no marker: still `spacedock 0.26.0+dev`; doctor exited 0;
- `Version=0.27.0-pre0` plus the release marker: `spacedock 0.27.0-pre0`; doctor against `0.26.0` exited 1 with the existing update-plugin remedy.

This closes the unverified mechanism: Go's linker can stamp the unexported string, the exact conditional ignores the misleading stamp without it, and the strict gate still rejects a genuinely marked next-minor release. The throwaway checkout was removed; no spike code is part of this entity.

## Acceptance criteria

**AC-1 (VALUE) — A source build remains compatible with its adjacent checkout even when its Git-derived version candidate is a future-minor automatic `pre0`.**

Verified by a real three-build fixture that reads minor Y independently from `.claude-plugin/plugin.json`, constructs a misleading `vX.(Y+1).0-pre0-N-gSHA` ldflag, runs `--version` and doctor, and observes `Y+dev` plus exit 0 both with no ldflag and with the unmarked misleading ldflag. It fails if `displayVersion` trusts `Version` without the marker or the source build stops deriving Y from the manifest.

**AC-2 — Only an explicitly marked release pipeline artifact claims its exact release version.**

Verified by the same real-build fixture observing the exact tag version only with `Version=R` plus `releaseBuild=true`, and by a parsed `.goreleaser.yaml` guard requiring both stamps on every `spacedock-stable` and `spacedock-edge` build. Adversarial twins remove the marker from either build and must red. It fails if a version ldflag alone is sufficient or either release channel omits the marker.

**AC-3 — Source-build compatibility is invariant under tag, revision, distance, and dirty provenance.**

Verified by a table of unmarked source builds/units carrying a clean release tag, future-minor `git describe`, bare SHA, and dirty suffix: every case emits the same manifest-derived identity and identical doctor verdict. The independent baseline is the manifest read from disk. It fails if any provenance string changes output bytes or exit code.

**AC-4 — Published source-build guidance has one canonical unstamped command and explains the identity boundary accurately.**

Verified by executing the documented `go build -o spacedock ./cmd/spacedock` command and observing the AC-1 result, plus human review of the concrete doc/comment diff below. No committed prose-grep substitutes for behavior. It fails review if source guidance injects a version/provenance ldflag or implies a nearest tag defines compatibility.

**AC-5 — Automatic `pre0` release behavior and the strict First Officer minor gate remain unchanged.**

Verified by the marked-release arm rejecting Y+1 against a Y manifest, the existing minor-exact contract table, and the existing edge-advance/pre0 decision and wiring tests. It fails if one-minor drift becomes compatible, the auto-pre0 test algebra changes, or the release marker becomes an unreleased-source compatibility exception.

## Proposed mechanisms traced to value

| Mechanism | Value AC served | Simplest alternative replaced | Why the alternative is insufficient |
| --- | --- | --- | --- |
| Exact `releaseBuild=true` trust check in `displayVersion` | AC-1, AC-3, AC-5 | Remove misleading docs only | Prose cannot constrain ambient/copy-pasted ldflags |
| Marker in both GoReleaser builds | AC-2, AC-5 | Assume all `Version` stamps are releases | That assumption is the defect reproduced above |
| Real three-build + doctor fixture | AC-1, AC-2, AC-3 | Unit-test only | A unit call does not prove linker wiring or emitted CLI bytes |
| Parsed release-config guard with adversarial twins | AC-2 | One manual dry-run | A future channel-specific deletion could silently regress after the one run |
| Source/release docs and comments | AC-4, supporting AC-1/AC-2 | Leave stale `git describe` wording | It would teach operators to recreate the now-ignored and conceptually wrong stamp |

## Expected surface

Gate baseline: **9 files, about 150–210 lines added/changed, no production package outside `internal/cli`.**

| File | Intended change | Expected churn |
| --- | --- | --- |
| `internal/cli/cli.go` | replace the Git-provenance-as-version comment; add `releaseBuild` | 8–14 lines |
| `internal/cli/dev_version.go` | require exact release marker before trusting `Version` | 5–10 lines |
| `internal/cli/dev_version_test.go` | unit matrix and real three-build/doctor fixture | 65–95 lines |
| `internal/cli/frontdoor_test.go` | make the existing `withVersion` helper explicitly model release builds | 4–8 lines |
| `.goreleaser.yaml` | stamp the marker in stable and edge builds; correct comments | 8–14 lines |
| `internal/release/goreleaser_guard_test.go` | parse/guard both build entries and add marker-removal twins | 45–70 lines |
| `.github/workflows/release.yml` | correct checkout/version-resolution comment only | 4–8 lines |
| `docs/releasing.md` | define the two-stamp release contract and provenance separation | 10–18 lines |
| `docs/site/contributing/build-from-source.md` | make plain `go build` canonical and explain `manifest+dev` | 6–12 lines |

Tolerance: up to **10 files and 260 total changed lines** if the release-config parser is cleaner in a dedicated `internal/release/*_test.go`. Re-gate before touching `internal/contract`, First Officer prose, auto-pre0/edge-advance logic, introducing a metadata/profile abstraction, exceeding 10 files/260 changed lines, or adding compatibility for unmarked prototypes.

## Test and validation plan

1. **Focused RED first:** extend `dev_version_test.go` so a real binary carrying only the computed future-minor `git describe`-shape stamp is expected to emit manifest Y+dev and pass doctor. It must fail on current code. Add the marked release control and unit provenance table before production edits. Cost: low, three local builds, fixture/CLI.
2. **Implement the two-line trust boundary and test helper update:** make focused `internal/cli` tests green. No contract changes.
3. **Release wiring:** parse `.goreleaser.yaml` by build ID and ldflag target; require `Version={{ .Version }}` and `releaseBuild=true` for both stable and edge, with marker-removal/wrong-value adversarial twins. Run GoReleaser 2.16 `build --snapshot --clean --single-target --id spacedock-stable` and `--id spacedock-edge`, execute both artifacts, and observe the snapshot release identity. Cost: medium, release-config fixture plus real CLI dry-run.
4. **Policy non-regression:** run focused `./internal/contract`, `./internal/cli`, and `./internal/release` tests, including `TestPre0MinorMatchesDevRequiredMinor`, edge-advance decision/wiring, and the marked Y+1 versus Y doctor refusal. No new live workflow scenario is needed; this is existing launch behavior with corrected identity input.
5. **Repository checks:** `gofmt -w ./cmd ./internal`, then `go test ./...`, then `go test ./... -race`. Record named claims and falsifying edits, not pass counts.
6. **Detached adversarial audit:** in a throwaway checkout, (a) revert the marker check to `Version != "dev"` and require the source fixture to red, (b) remove the marker from only the edge GoReleaser build and require the config guard to red, and (c) relax the contract by one minor or bypass marked-release comparison and require AC-5 tests to red. Audit is mandatory because both front-door compatibility and release machinery are high-stakes.
7. **Diff-derived live CI:** `internal/cli/dev_version.go` feeds the host-neutral front-door gate exercised by all supported hosts, so `claude-live`, `codex-live`, and `pi-live` are all required green before merge. This follows the actual launch-path diff, even though `.goreleaser.yaml` alone would not require host lanes. The first real tag cut after merge additionally records stable/edge artifact line 1 as release evidence.

## Concrete documentation and comment diff

`docs/site/contributing/build-from-source.md`, after the plain build command, changes:

> Prints `spacedock <version>` for your local build.

to:

> Keep this build unstamped: do not add a Git tag, revision, or `git describe` ldflag. A source build reports the embedded checkout manifest version plus `+dev` (for example `spacedock 0.26.0+dev`), so it stays compatible with the adjacent skills even when a future-minor release tag is the nearest Git ancestor.

`docs/releasing.md` changes the release summary from:

> stamping `git describe --tags` into `internal/cli.Version`, for BOTH channels

to:

> deriving the pushed tag version through GoReleaser and stamping both `internal/cli.Version={{ .Version }}` and the explicit `internal/cli.releaseBuild=true` marker for BOTH channels; only artifacts carrying that marker may use the stamped release version as compatibility identity

The Dev-Only `next` section additionally states: plain `go build`/`go install` uses the embedded manifest plus `+dev`; tag/revision/dirty data is provenance and does not affect compatibility. Release stamps belong only to `.goreleaser.yaml`, not source-build examples.

The matching comments in `internal/cli/cli.go`, `.goreleaser.yaml`, and `.github/workflows/release.yml` stop calling `git describe` the binary's unconditional source of truth. They state the mechanical two-input contract: `Version` is trusted only beside `releaseBuild=true`; full tags are still fetched for GoReleaser tag resolution and annotated release-note bodies. These are comment/doc changes reviewed manually and exercised through the behavior tests, never policed by a committed prose grep.

## Stage Report: ideation

- DONE: Reproduce and cite the exact source-build version drift: embedded checkout manifest, ordinary unstamped go build, misleading git-describe stamp, and the First Officer minor-version consequence.
  Problem table records `0.26.0`, `v0.27.0-pre0-89-g4ff98d8c`, both binary outputs, doctor exits 0/1, and the pre-boot FO abort.
- DONE: Compare the smallest mechanical source-versus-release identity boundaries, choose one, and explain why documentation-only cleanup or a relaxed compatibility gate cannot enforce the value.
  Approach comparison chooses one unexported `releaseBuild` string marker and rejects dedicated-slot, prefix, build-tag, docs-only, and relaxed-gate alternatives.
- DONE: Specify exact source-build, release-build, and provenance behavior without changing the automatic pre0 tag policy or adding compatibility for unreleased prototypes.
  Exact behavior table makes unmarked builds manifest+dev, marked releases exact, provenance inert, and a genuine Y+1 release incompatible with Y.
- DONE: Turn the seed criteria into entity-level value ACs with independent, falsifiable evidence, including a misleading future-minor tag/stamp case and explicit release-pipeline wiring.
  AC-1 through AC-5 name independent manifest/config baselines and the exact edits that would make each proof fail.
- DONE: Declare the intended files and LOC with a tolerance, trace every proposed mechanism to a value AC, and name the simplest alternative it replaces.
  Expected surface is 9 files/150–210 changed lines, tolerance 10/260; the mechanism table maps each mechanism to ACs and rejected simpler choices.
- DONE: Define the riskiest-first spike, focused tests, full/race/format checks, detached adversarial audit, and diff-derived live-CI lane requirements.
  Spike exercised three real binaries; plan requires focused tests, gofmt/full/race, three adversarial edits, and Claude/Codex/Pi live lanes.
- DONE: Propose concrete documentation/comment changes for the canonical source build and release stamping contract.
  Concrete before/after wording covers the source guide, release guide, CLI/config/workflow comments, and bans release stamps from source examples.
- DONE: Append a complete ideation stage report with every checklist item marked DONE, SKIPPED, or FAILED; commit and push state, then stop at the ideation gate without implementation.
  This report is appended to the entity body; only split-root state will be committed and pushed.

### Summary

Ideation reproduces the false next-minor identity and chooses a single-purpose release marker that makes copied Git provenance inert while preserving exact release identity. The detached spike proves the boundary and strict-gate control; the design now has falsifiable ACs, bounded surface, release wiring, docs, adversarial audit, and diff-derived live-lane requirements, with no implementation performed.

## Implementation Intended Surface Declaration

Before implementation edits, the intended deliverable surface is exactly these nine files:

- `internal/cli/cli.go` — add the single linker-writable `releaseBuild` marker and correct the version comment (estimated 8–12 changed lines; mechanism serves AC-1, AC-2, AC-3, and AC-5).
- `internal/cli/dev_version.go` — trust `Version` only when `releaseBuild == "true"` and `Version != "dev"` (estimated 4–8 changed lines; mechanism serves AC-1, AC-2, AC-3, and AC-5).
- `internal/cli/dev_version_test.go` — add the unit provenance matrix and real source/misleading-stamp/marked-release binary-plus-doctor proof before production edits (estimated 75–95 changed lines; proof serves AC-1, AC-2, AC-3, and AC-5).
- `internal/cli/frontdoor_test.go` — make the existing version helper explicitly stamp marked releases (estimated 4–8 changed lines; preserves the release-shaped AC-5 fixtures).
- `.goreleaser.yaml` — stamp `releaseBuild=true` in both stable and edge builds and correct the surrounding comment (estimated 6–10 changed lines; mechanism serves AC-2 and AC-5).
- `internal/release/goreleaser_guard_test.go` — parse the release configuration by build ID/ldflag target and prove stable-only and edge-only marker-removal twins fail (estimated 45–60 changed lines; proof serves AC-2).
- `.github/workflows/release.yml` — correct the checkout/tag-resolution comment without changing workflow behavior (estimated 2–4 changed lines; documentation supports AC-4 and AC-5).
- `docs/releasing.md` — document the two-stamp release contract and separate provenance from compatibility identity (estimated 10–15 changed lines; documentation serves AC-4 and supports AC-2).
- `docs/site/contributing/build-from-source.md` — make plain unstamped `go build` canonical and explain manifest-plus-dev identity (estimated 8–12 changed lines; documentation serves AC-4 and supports AC-1).

Estimated total is nine files and 190–220 changed lines. This is one file and 40–70 lines below the approved 10-file/260-line cap; there is no difference in file identity from the approved nine-file baseline and no requested tolerance expansion. Implementation will not touch `internal/contract`, First Officer prose, automatic pre0 or edge-advance logic, add a metadata/profile abstraction, or add compatibility for unmarked prototypes; any such need triggers a design reset before edits.

### Feedback Cycles

- Cycle 1: REJECTED — Roborev `branch_final` panel job 2653; surface 9 files/255 LOC vs estimate 9 files/220 LOC (116%); AC unchanged

## Stage Report: implementation

- DONE: Before editing, append an Implementation Intended Surface Declaration naming exact files, estimated changed LOC, each mechanism's value AC, and any difference from the approved 10-file/260-line cap; stop for re-gate on prohibited or over-tolerance drift.
  The declaration names the approved nine files, maps every mechanism/proof to AC-1 through AC-5, and declares 190–220 estimated lines under the 10/260 cap.
- DONE: Add the focused RED proof first: a real source binary carrying a misleading future-minor Version ldflag but no release marker must still emit the embedded manifest plus +dev and pass doctor; include the exact marked-release control.
  Before production edits, `TestSourceBuildCompatibilityIdentity/misleading` red with future-minor line 1 and doctor exit 1; final three-build proof requires source `0.26.0+dev`/doctor OK and marked `0.27.0-pre0`/mismatch.
- DONE: Implement the single-purpose releaseBuild trust boundary in internal/cli without a generalized build profile, provenance subsystem, compatibility carve-out, internal/contract change, or First Officer prose change.
  Commits `a9180d43` and `2992f430` add only unexported `releaseBuild = "false"` plus the exact `releaseBuild == "true" && Version != "dev"` trust check.
- DONE: Stamp the exact release marker in both GoReleaser build entries and add an independent parsed configuration guard whose stable-only and edge-only marker-removal twins fail.
  Stable and edge carry the exact linker marker; parsed build-ID guard tests green, and both channel-specific removal mutants red.
- DONE: Apply the approved source-build, release-guide, CLI, GoReleaser, and release-workflow documentation/comment changes without changing automatic pre0 or edge-advance behavior.
  The documented plain `go build -o spacedock ./cmd/spacedock` emitted `0.26.0+dev` and doctor OK; no pre0/edge implementation file changed.
- DONE: Run focused CLI/contract/release tests, both stable and edge GoReleaser snapshot single-target dry-runs, gofmt -w ./cmd ./internal, go test ./..., go test ./... -race, and git diff --check; record exact outcomes and any infrastructure failure honestly.
  Focused packages, final full suite, and final race suite passed; both GoReleaser 2.16 dry-runs emitted exact `0.27.0-snapshot-070fb828`; gofmt and diff check were clean, including existing strict-minor/auto-pre0/edge-advance tests.
- DONE: Run the three-edit detached adversarial audit from the approved design and show each claim-breaking edit makes the intended proof fail.
  Three throwaway worktrees proved Version-only trust reds source line/doctor, edge-only marker removal reds the parsed guard, and bypassing marked-release identity reds exact release line/doctor; worktrees were removed.
- DONE: Measure the final file/LOC surface against the approved baseline, inspect for unrelated drift, and stop for design reset rather than rationalizing any prohibited or over-tolerance expansion.
  `070fb828..2992f430` is exactly the approved nine files, 208 additions/51 deletions = 259 changed lines; no unrelated file, `internal/contract`, FO prose, pre0/edge logic, profile abstraction, or prototype compatibility drift exists.
- DONE: Request one final Roborev over the completed implementation, triage every finding by released workflow/value AC and materiality, and fix only material findings.
  Roborev `branch_final` job 2653: one material AC-1/AC-5 doctor-verdict proof gap fixed in `2992f430`; snapshot risk deferred until a supported adjacent-plugin workflow exists; four polish observations declined without AC change.
- FAILED: Record any advisory round through gate record --round.
  `${SPACEDOCK_BIN} gate record source-build-compatibility-identity --workflow-dir docs/dev --round implementation/1 --briefing …/briefing.json --log …/briefing.review.jsonl --feedback-cycle …/feedback-cycle.txt` exited 1: `Error: gate record --round requires folder-form entity <slug>/index.md because review artifacts accumulate beside the entity`; the First Officer directed no flat-task migration.
- SKIPPED: Run all-host live CI.
  Claude/Codex/Pi live lanes are required at PR time; this stage opened no PR and claims no local live-host evidence.
- DONE: Append a complete implementation Stage Report with every checklist item DONE, SKIPPED, or FAILED, commit and push code plus state, and stop for fresh independent validation without opening a PR.
  Code commits are pushed on `spacedock-ensign/source-build-compatibility-identity`; this report is path-scoped for the split-root state commit/push, and no PR was opened.

### Summary

Source builds now ignore provenance-only Version stamps and use embedded checkout manifest identity plus `+dev`, while both release channels explicitly opt into exact release identity. Real binaries, parsed release wiring, dry-runs, full/race suites, and three detached mutants prove the boundary; the sole failed ceremony is advisory-round recording, structurally unavailable for this flat entity and explicitly left unmigrated for independent validation.

## Stage Report: validation

- DONE: Reproduce AC-1 through AC-3 with real source, misleading-Version, and explicitly marked-release binaries: assert exact version line, doctor exit, and compatibility outcome independently of the implementation report.
  Eight independent builds made plain/tag/pre0/SHA/dirty/wrong-marker/dev-sentinel shapes emit exact `spacedock 0.26.0+dev (contract 3)` with doctor exit 0; only marked `0.27.0-pre0` emitted exactly and exited 1 against manifest 0.26.
- DONE: Verify AC-1 (VALUE) — A source build remains compatible with its adjacent checkout even when its Git-derived version candidate is a future-minor automatic `pre0`.
  Real `Version=v0.27.0-pre0-17-gabcdef0` without the marker emitted manifest-derived `0.26.0+dev`, doctor exit 0, and the exact compatible message.
- DONE: Verify AC-2 — Only an explicitly marked release pipeline artifact claims its exact release version.
  Independent YAML parsing found both exact stamps once in stable and edge; GoReleaser 2.16 snapshots both emitted exact `0.27.0-snapshot-2992f430`, while wrong-case or missing markers remained source identity.
- DONE: Verify AC-3 — Source-build compatibility is invariant under tag, revision, distance, and dirty provenance.
  Clean tag, future-minor describe, bare SHA, and dirty describe candidates produced byte-identical version lines, doctor exits, and compatibility messages against the independent manifest baseline.
- DONE: Verify AC-4 — Published source-build guidance has one canonical unstamped command and explains the identity boundary accurately.
  Human review found one plain build command and no provenance stamp; executing exact `go build -o spacedock ./cmd/spacedock` emitted `0.26.0+dev` and doctor exit 0.
- DONE: Verify AC-5 — Automatic `pre0` release behavior and the strict First Officer minor gate remain unchanged.
  `TestCompare`, `TestAutoPre0MinorEqualsRequiredMinor`, edge decision/wiring/always-cut guards, and marked Y+1 doctor refusal passed; the exact diff touches no contract, pre0, edge-decision, or First Officer file.
- DONE: Verify AC-2, AC-4, and AC-5 from independent sources: parse both release configurations, run stable/edge guard and snapshot proofs, inspect the published source-build guidance, and confirm automatic pre0 plus strict-minor behavior did not change.
  Independent Ruby parsing, `goreleaser check`, both single-target snapshots, the executed published command, and focused policy tests all passed with exact outcomes above.
- DONE: Audit the exact 070fb828..HEAD surface against the approved 10-file/260-line boundary, run focused/full/race/format/diff checks and the required detached adversarial proof; classify the flat-task round-recorder failure separately and reject v218 only if it invalidates a value AC or non-negotiable boundary.
  `070fb828..2992f430` is nine files/259 lines with no prohibited drift; focused, full, race, gofmt, and range diff checks passed, and the code worktree remained clean.
- DONE: Perform the required detached adversarial proof.
  Three separate throwaway checkouts red exactly when Version-only trust broke AC-1, edge marker removal broke AC-2, or marked-release bypass broke the exact line and strict doctor verdict; those temporary artifacts were removed.
- DONE: Classify the flat-task round-recorder failure separately.
  Observed trigger is `--round` against flat `<slug>.md` without `<slug>/index.md`; the byte-clean refusal affects no AC or compatibility boundary, so it is separate workflow-mechanism friction, not a v218 material or deferred product finding.
- SKIPPED: Run all-host live CI.
  Claude/Codex/Pi lanes remain required at PR time; no PR-time lanes were run or claimed by this local validation.
- DONE: Issue a validation recommendation.
  PASSED: AC-1 through AC-5 have independent valid evidence, no material finding remains, and there are no v218 deferred risks; all-host live CI remains a later PR gate.
- DONE: Append a complete validation Stage Report with every checklist item marked DONE, SKIPPED, or FAILED; commit and push split-root state, then stop at the validation gate.
  This report is appended without frontmatter edits and is path-scoped for the required state commit and push.

### Summary

Independent binaries, parsed release wiring, canonical GoReleaser validation, snapshots, policy guards, full/race suites, and detached mutants establish every acceptance criterion. Recommendation is PASSED with no v218 finding; the flat-form round refusal is separate workflow friction, and all-host live CI remains required at PR time.
