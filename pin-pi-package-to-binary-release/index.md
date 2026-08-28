---
id: qcwrzza9xkr5kfwmdmvbqhkv
title: Pin the installed Pi package to the launcher release
status: ideation
source: "Released v0.27.2 Pi bootstrap abort reported by the captain on 2026-08-28"
started:
completed:
verdict:
score: "1.0"
worktree:
issue:
pr:
mod-block:
gates:
    version: 1
    records:
        - id: gate:qcwrzza9xkr5kfwmdmvbqhkv:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:qcwrzza9xkr5kfwmdmvbqhkv-backlog-1
              briefing:
                id: briefing:qcwrzza9xkr5kfwmdmvbqhkv:backlog:attempt-1:revision-1
                digest: sha256:b978c947fcd19d4db838077f568d7d231f71fd4bbe7b146b6752a560ad0e579e
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:qcwrzza9xkr5kfwmdmvbqhkv:backlog:1
                briefing: briefing:qcwrzza9xkr5kfwmdmvbqhkv:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-28T18:46:31.966666Z"
                decision: approve
                reason: 'Captain approve (2026-08-28 chat): seed QC-verified, QC amendments folded; advance to ideation.'
              application:
                target-stage: ideation
                state: consumed
---

A stable launcher can install skills from a newer release and then fail its first binary gate.

## Problem

`spacedock install --host pi` uses the unpinned source `git:github.com/spacedock-dev/spacedock`.
Pi updates that source from the repository default branch. A stable v0.27.2 launcher therefore installed the v0.28.0-pre1 first-officer skill from `main`. The skill correctly rejected the older binary before `status --boot`.

The prior sibling-provider fix covered Claude and Codex only. Its approved scope explicitly excluded Pi.

## Proposed approach

Bind the default Pi package source to the running launcher's release identity: a release-stamped binary pins the source to its own release ref (`git:github.com/spacedock-dev/spacedock@<ref>`, ref derived from the stamped `internal/cli.Version`) and never floats across tags. The unstamped development-build source behavior is defined: a `dev`-sentinel build (unstamped `go build`/`go install`) keeps the current unpinned source and floats to the default branch — dev builds exist to track the default branch, and pinning one to a tag that does not exist would break the development loop; the consequence is accepted and stated: a dev-build install can skew, and the first-officer binary gate remains the loud backstop. A release-stamped binary never floats. State the pinning tradeoff as intended semantics: a pinned source does not auto-update, so stable users receive package fixes by upgrading the launcher, matching the Claude/Codex marketplace pin. Preserve `--plugin-dir` as the explicit development override.

Pin detection reads the Pi settings package entry: `~/.pi/agent/settings.json` `packages` entry for the spacedock package (already surfaced as `piPackageStatus.source` by `piSpacedockPackageStatus`, `internal/cli/pi.go`). The entry is a source string; a git source carrying `@<ref>` is pinned (Pi's `parseGitUrl`/`buildGitSource` set `pinned: Boolean(ref)`), and a bare `git:repo` entry is unpinned. Missing = no spacedock entry or an unresolvable package root.

Make the normal `spacedock pi` front door own one repair attempt. If the package is missing, unpinned (release-stamped binary with a ref-less entry), or incompatible (installed package manifest version mismatches the binary), run one `pi install <pinned source>`; Pi reconciles an existing clone to the configured ref and rewrites the settings entry to the pinned source. Recheck the package (the existing `piSpacedockPackageStatus` discovery) before launch. Refuse the launch if the repair fails or the package still does not match.

Exercise the ordinary installed Pi front door. Do not use `--plugin-dir`, `SPACEDOCK_REPO_ROOT`, a prose marker, or a direct skill path as the proof.

## Risk evidence

Pi's installed package manager marks a Git source with `@ref` as pinned. It fetches and checks out that ref. An unpinned source updates from its default branch. The released v0.27.2 source is unpinned, while its embedded first-officer contract requires minor 0.27. The current `main` contract requires minor 0.28.

Spike (2026-08-28, this session): the mechanism is verified against the real shipped Pi package-manager implementation on this machine (`@earendil-works/pi-coding-agent` dist, `core/package-manager.js` + `utils/git.js`, the code that executes `pi install` here). `parseGitUrl`/`buildGitSource` set `pinned: Boolean(ref)`; update-candidate collection skips `local || pinned` sources and updates unpinned git sources from origin HEAD; `pi install git:<repo>@<ref>` clones and checks out `ref`; `addSourceToSettings` rewrites the existing `packages` entry (match key ignores ref) so one install call converts an unpinned entry to its pinned form and reconciles the clone. Pin detection reads the `packages` entry string in `~/.pi/agent/settings.json` — the field `piSpacedockPackageStatus` already surfaces as `packageStatus.source`. No live spike is needed beyond the AC-1 installed-front-door run, which exercises this mechanism end-to-end at validation.

## Out of scope

- Lowering the first-officer binary floor.
- Moving or replacing the published v0.27.2 tag.
- Claude or Codex provider inventory.
- A second package manager or a new Pi discovery protocol.
- Prose-grep tests.

## Expected surface and tolerance

Estimate net LOC change: +100, across 4 files. Tolerance: +/-60 net LOC and +/-2 files.

Expected files: `internal/cli/pi.go` (pinned-source derivation replacing the bare const at the install call site, one front-door repair attempt with recheck and refusal), one new small file for the source-derivation helper (release ref derivation + dev sentinel handling), `internal/cli/pi_frontdoor_test.go` (repair/pin command-behavior tests), and the Pi live test wiring for the no-override proof journey. Observable semantics that may change: stored format (the `packages` entry for spacedock in `~/.pi/agent/settings.json` gains an `@<ref>` suffix for release-stamped binaries) and front-door runtime behavior (one repair attempt with launch refusal on remaining mismatch). Command grammar is unchanged; no new flags. Authority is unchanged. Doctor remedy lines may gain wording pointing at the pinned repair.

## Acceptance criteria

**AC-1 (VALUE) - A released launcher repairs and boots the Pi skill suite from the same release line.**
Verified by: one live Pi journey through the ordinary installed front door with no development override (no `--plugin-dir`, no `SPACEDOCK_REPO_ROOT`), started with the Spacedock package absent or unpinned, observing exactly one `pi install` of the pinned source followed by `status --boot` reaching a launch-ready report. Independent baseline that must move the wrong way without the fix: the same start state with the repair removed leaves the package missing/unpinned and the FO binary gate aborts before workflow work (the observed v0.27.2 failure). Verifier: the live journey run plus its retained artifacts (this is also the anti-tautology baseline for AC-2/AC-4's unit assertions).

**AC-2 - The default Pi install source is release-pinned, and the development override remains local.**
Verified by: `internal/cli/pi_frontdoor_test.go` command-behavior tests asserting the literal source passed to the `PiInstall` op: for a release-stamped binary, `git:github.com/spacedock-dev/spacedock@<that binary's tag>`; for a `dev`-sentinel binary, the bare floating source; for `--plugin-dir`, the plugin dir and never the release source. Falsifier: changing the release source to omit its ref, or making the dev-override path use the release source, flips a literal assertion.

**AC-3 - The release fix is safe for v0.27.3 and the v0.28 prerelease line.**
Verified by: table tests over stable (`v0.27.3`), prerelease (`v0.28.0-pre1`), and unstamped (`dev`) identities asserting literal expected sources (stable/prerelease pin to their own tag; dev floats). Falsifier: changing the source derivation for any covered identity.

**AC-4 - A failed or ineffective repair cannot launch Pi.**
Verified by: command-behavior tests where the `PiInstall` op fails or the post-install `SpacedockPackageStatus` still reports missing/unpinned: exactly one install attempt, one recheck, no launch (`pi` exec never reached), and an actionable error naming the remedy. Falsifier: a second install attempt, a launch, or exit 0 after remaining mismatch.

## Test plan

Use focused command behavior tests first: `internal/cli/pi_frontdoor_test.go` (existing fake `piRuntimeOps` pattern) for AC-2/AC-3/AC-4 — literal expected sources per identity, one-install/one-recheck/no-launch-on-failure, and the dev override isolation; estimated cost: cheap unit/behavior tests, no fixtures beyond the existing fakes. Then one existing Pi live journey through the installed-package front door without a development override, started with no valid Spacedock package (fixture: install a wrong-line or bare `git:` source first), per AC-1. Then gofmt, the full suite, race, and a detached adversarial audit because this changes a front door. Documentation impact: none on the docs site (no command-surface change); the doctor remedy wording change, if any, is covered by the behavior tests' output assertions.

### Feedback Cycles

{First officer appends one `- Cycle {N}: ...` line per correction round; the validation gate reads reviewer findings from here.}

## Stage Report: ideation

- DONE: flesh out the task body at pin-pi-package-to-binary-release/index.md: chosen dev-build source behavior (stated with its consequence), pinned-source mechanism, package-lock detection field the front-door repair reads, expected surface with net LOC/files/tolerance, observable-semantics declaration, and a test plan naming the verifier for each AC
  dev sentinel keeps the floating default-branch source (stated consequence: dev installs can skew; the FO binary gate is the backstop); release-stamped binaries pin to their own tag; detection field named (the `packages` entry string in `~/.pi/agent/settings.json`, already surfaced by `piSpacedockPackageStatus` as `packageStatus.source`); AC verifiers named per AC; semantics declared in Expected surface.
- DONE: ensure at least one AC measures end-value against an independent baseline (a stable launcher with a floating source installs mismatched skills and aborts; after the fix, the same launcher installs same-line skills and reaches status --boot), with the no-override installed front door as the proof run
  AC-1 names the baseline (repair removed → package missing/unpinned → FO gate aborts, the observed v0.27.2 failure) and the no-override installed front door as the proof run with retained artifacts.
- DONE: record the spike decision in the task body: either exercise the riskiest unverified mechanism (Pi package manager's @ref pin + update behavior against the real package root) or write "no spike needed:" naming the proven mechanisms relied on
  Riskiest mechanism verified by reading the real shipped package-manager implementation running on this machine (`pi-coding-agent` dist `core/package-manager.js` + `utils/git.js`): `pinned: Boolean(ref)`, unpinned git sources are update candidates from origin HEAD, `pi install git:<repo>@<ref>` clones+checkouts the ref, `addSourceToSettings` rewrites the existing entry (match key ignores ref), and the settings entry string is the pin record. Recorded under "Spike" in Risk evidence; AC-1's live run remains the end-to-end exercise.

### Summary

Ideation fleshed out the seed into a gated design: pinned-source derivation from the launcher's release identity (dev sentinel floats, stated with its consequence), front-door one-shot repair with recheck and refusal, the settings.json `packages` entry as the detection field, AC verifiers named per criterion (literal-source behavior tests + one no-override live journey), and the package-manager `@ref` mechanics verified against the real installed Pi implementation rather than prose. No spike needed beyond that source-level verification; AC-1's live run remains the end-to-end exercise.
