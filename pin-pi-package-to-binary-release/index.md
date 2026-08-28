---
id: qcwrzza9xkr5kfwmdmvbqhkv
title: Pin the installed Pi package to the launcher release
status: implementation
source: "Released v0.27.2 Pi bootstrap abort reported by the captain on 2026-08-28"
started:
completed:
verdict:
score: "1.0"
worktree: .worktrees/spacedock-ensign-pin-pi-package-to-binary-release
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
        - id: gate:qcwrzza9xkr5kfwmdmvbqhkv:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:qcwrzza9xkr5kfwmdmvbqhkv-ideation-1
              briefing:
                id: briefing:qcwrzza9xkr5kfwmdmvbqhkv:ideation:attempt-1:revision-1
                digest: sha256:5486fa9d172e375ab53f5bdd9e3dadfb7afdc3e0825e72f73b43db0ba56f052a
                room-ref: ./review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:qcwrzza9xkr5kfwmdmvbqhkv:ideation:1
                briefing: briefing:qcwrzza9xkr5kfwmdmvbqhkv:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-28T19:21:25.084766Z"
                decision: approve
                reason: 'Captain approve (2026-08-28 chat): cycle-2 design review-confirmed RESOLVED; advance to implementation.'
              application:
                target-stage: implementation
                state: consumed
review-round:
    id: round:qcwrzza9xkr5kfwmdmvbqhkv:implementation:1
    stage: implementation
    cycle: 1
    briefing:
        id: briefing:qcwrzza9xkr5kfwmdmvbqhkv:implementation:round-1:revision-1
        digest: sha256:9d201f249c99e925759d1c17ac725ec593693f96744722c581a00549ac27c0d2
        room-ref: ./review/implementation/round-1
---

A stable launcher can install skills from a newer release and then fail its first binary gate.

## Problem

`spacedock install --host pi` uses the unpinned source `git:github.com/spacedock-dev/spacedock`.
Pi updates that source from the repository default branch. A stable v0.27.2 launcher therefore installed the v0.28.0-pre1 first-officer skill from `main`. The skill correctly rejected the older binary before `status --boot`.

The prior sibling-provider fix covered Claude and Codex only. Its approved scope explicitly excluded Pi.

## Proposed approach

Bind the default Pi package source to the running launcher's release identity: a release-shaped binary pins the source to its own release ref (`git:github.com/spacedock-dev/spacedock@<ref>`) and never floats across tags. The pin ref is derived from the release identity, resolved in order: (1) the linker-stamped `internal/cli.Version` when the release artifacts set it (`-X ...Version={{ .Version }}`, stable and edge channels alike); (2) when `Version == "dev"`, the Go build-info main module version (`debug.ReadBuildInfo().Main.Version`) when it is a semver tag — this covers `go install github.com/.../spacedock@v0.27.2` proxy builds, which carry no ldflags but ARE release-shaped and must not float (they embed the tagged manifest, which is exactly why `displayVersion()` reports them as `X.Y.Z+dev`); (3) otherwise the `dev` sentinel: a plain `go build`/`go install ./cmd/spacedock` checkout build keeps the current unpinned source and floats to the default branch — dev builds exist to track the default branch, and the consequence is accepted and stated: a dev-build install can skew, and the first-officer binary gate remains the loud backstop. A release-shaped binary (stamped or proxy-tagged) never floats. State the pinning tradeoff as intended semantics: a pinned source does not auto-update, so stable users receive package fixes by upgrading the launcher, matching the Claude/Codex marketplace pin. Preserve `--plugin-dir` as the explicit development override.

Pin detection reads the Pi settings package entry: `~/.pi/agent/settings.json` `packages` entry for the spacedock package (already surfaced as `piPackageStatus.source` by `piSpacedockPackageStatus`, `internal/cli/pi.go`). The entry is a source string; a git source carrying `@<ref>` is pinned (Pi's `parseGitUrl`/`buildGitSource` set `pinned: Boolean(ref)`), and a bare `git:repo` entry is unpinned. Missing = no spacedock entry or an unresolvable package root.

Make the normal `spacedock pi` front door own one repair attempt. The repair triggers for a release-shaped binary (stamped or proxy-tagged) when the discovered `packageStatus.source` is absent, or is a git source whose ref is absent or differs from the binary's release ref — the wrong-ref case covers an installed package pinned to another line (e.g. `@v0.28.0-pre1` under a v0.27.2 binary) without any package-manifest version field: the mismatch IS the ref delta, and the repair reinstalls the binary's own ref. A non-git entry (`file:`, `npm:`, or local path) is user-managed: the repair never runs and never rewrites it (no clobber in either direction), and launch proceeds with the first-officer binary gate as the version backstop. A `dev`-sentinel binary performs no repair and no install: it has no pin target, and a floating reinstall would clobber a pinned entry with an unpinned one — the skew recreated in the other direction. When the repair does run: one `pi install <pinned source>`; Pi reconciles an existing clone to the configured ref and rewrites the settings entry to the pinned source. Recheck the package (the existing `piSpacedockPackageStatus` discovery) before launch. Refuse the launch if the repair fails or the package still does not match. Under `--plugin-dir` or `SPACEDOCK_REPO_ROOT` the repair path is suppressed entirely — a declared development override owns the package surface for that run, and the repair never fights it.

Exercise the ordinary installed Pi front door. Do not use `--plugin-dir`, `SPACEDOCK_REPO_ROOT`, a prose marker, or a direct skill path as the proof.

## Risk evidence

Pi's installed package manager marks a Git source with `@ref` as pinned. It fetches and checks out that ref. An unpinned source updates from its default branch. The released v0.27.2 source is unpinned, while its embedded first-officer contract requires minor 0.27. The current `main` contract requires minor 0.28.

Spike (2026-08-28, this session): the mechanism is verified against the real shipped Pi package-manager implementation on this machine (`@earendil-works/pi-coding-agent` dist, `core/package-manager.js` + `utils/git.js`, the code that executes `pi install` here). `parseGitUrl`/`buildGitSource` set `pinned: Boolean(ref)`; update-candidate collection skips `local || pinned` sources and updates unpinned git sources from origin HEAD; `pi install git:<repo>@<ref>` clones and checks out `ref`; `addSourceToSettings` rewrites the existing `packages` entry (match key ignores ref) so one install call converts an unpinned entry to its pinned form and reconciles the clone. Pin detection reads the `packages` entry string in `~/.pi/agent/settings.json` — the field `piSpacedockPackageStatus` already surfaces as `packageStatus.source`. No live spike is needed beyond the AC-1 installed-front-door run, which exercises this mechanism end-to-end at validation.

Accepted residual risks (stated, not silently held): the spike covered the locally installed Pi version only — an older Pi package manager that does not honor `@ref` pinning would need the remedy line to say `pi upgrade` + reinstall, and AC-1's live run is the check that the pinned install actually reconciles on the CI's pinned Pi version. Concurrent `spacedock pi` invocations could race the repair's settings.json rewrite; the last writer wins and the next front-door run repairs again, so the race is self-healing within the one-repair bound — acceptable for a single-user agent setup.

## Out of scope

- Lowering the first-officer binary floor.
- Moving or replacing the published v0.27.2 tag.
- Claude or Codex provider inventory.
- A second package manager or a new Pi discovery protocol.
- Prose-grep tests.

## Expected surface and tolerance

Estimate net LOC change: +661, across 6 files. Tolerance: +/-40 net LOC and +/-1 file (re-baselined to as-implemented actuals per the captain's design-reset ruling, 2026-08-28: the growth over the original +100/4-file estimate is the cycle-2 reviewer-demanded branches' test matrix (proxy-tagged identity, AC-5 no-clobber, wrong-ref repair) and the previously unbudgeted AC-1 live-journey wiring).

As-implemented files: `internal/cli/pi.go` (+54/−1; pinned-source derivation call site, one front-door repair attempt with recheck and refusal, repair suppression for non-git entries and dev builds), `internal/cli/pi_package_source.go` (+150 new, the source-derivation helper: linker stamp, then build-info tag for proxy builds, then the dev sentinel), `internal/cli/pi_frontdoor_test.go` (+9, seam extension), `internal/cli/pi_package_repair_test.go` (+335, AC-2/3/4/5 command-behavior tests), `internal/ensigncycle/pi_frontdoor_pin_live_test.go` (+96, the no-override live proof journey), `docs/runtime-live-ci-registry.md` (+18, live-journey registration). Observable semantics that may change: stored format (the `packages` entry for spacedock in `~/.pi/agent/settings.json` gains an `@<ref>` suffix for release-stamped binaries) and front-door runtime behavior (one repair attempt with launch refusal on remaining mismatch). Command grammar is unchanged; no new flags. Authority is unchanged. Doctor remedy lines may gain wording pointing at the pinned repair.

## Acceptance criteria

**AC-1 (VALUE) - A released launcher repairs and boots the Pi skill suite from the same release line.**
Verified by: one live Pi journey through the ordinary installed front door with no development override (no `--plugin-dir`, no `SPACEDOCK_REPO_ROOT`), started with the Spacedock package absent or unpinned, observing exactly one `pi install` of the pinned source followed by `status --boot` reaching a launch-ready report. Independent baseline that must move the wrong way without the fix: the same start state with the repair removed leaves the package missing/unpinned and the FO binary gate aborts before workflow work (the observed v0.27.2 failure). Verifier: the live journey run plus its retained artifacts (this is also the anti-tautology baseline for AC-2/AC-4's unit assertions).

**AC-2 - The default Pi install source is release-pinned for every release-shaped identity, and the development override remains local.**
Verified by: `internal/cli/pi_frontdoor_test.go` command-behavior tests asserting the literal source passed to the `PiInstall` op: for a linker-stamped binary, `git:github.com/spacedock-dev/spacedock@<that binary's stamp>`; for a `go install …@vX.Y.Z` proxy build (Version=="dev", build-info main version = the tag), `@<that tag>`; for a `dev`-sentinel build, no install at all; for `--plugin-dir`, the plugin dir and never the release source. Falsifier: changing the release source to omit its ref, letting a proxy-tagged build float, or making the dev-override path use the release source, flips a literal assertion.

**AC-3 - The release fix is safe for v0.27.3, the v0.28 prerelease line, and proxy-installed launchers.**
Verified by: table tests over stable (`v0.27.3`), prerelease (`v0.28.0-pre1`), unstamped (`dev`), and proxy-tagged (`Version=="dev"` + build-info `v0.27.2`) identities asserting literal expected sources (stable/prerelease/proxy pin to their own tag; dev floats). Falsifier: changing the source derivation for any covered identity.

**AC-4 - A failed or ineffective repair cannot launch Pi.**
Verified by: command-behavior tests where the `PiInstall` op fails or the post-install `SpacedockPackageStatus` still reports missing or a wrong-ref git source: exactly one install attempt, one recheck, no launch (`pi` exec never reached), and an actionable error naming the remedy. Falsifier: a second install attempt, a launch, or exit 0 after remaining mismatch.

**AC-5 - The repair never rewrites a user-managed or dev-owned package surface.**
Verified by: command-behavior tests where the settings entry is a non-git source (`file:`, `npm:`) or the run declares `--plugin-dir`/`SPACEDOCK_REPO_ROOT`, or the binary is a `dev`-sentinel build: no `PiInstall` call is made, the entry string is untouched, and launch proceeds (the first-officer binary gate remains the version backstop). Falsifier: a repair install attempted for any of those shapes, or an entry string rewritten by a dev build, flips the assertion.

## Test plan

Use focused command behavior tests first: `internal/cli/pi_frontdoor_test.go` (existing fake `piRuntimeOps` pattern) for AC-2/AC-3/AC-4/AC-5 — literal expected sources per identity (stamped, proxy-tagged, dev, plugin-dir), the wrong-ref repair trigger (entry `@v0.28.0-pre1` under a v0.27.2-identity binary → repair installs `@v0.27.2`), one-install/one-recheck/no-launch-on-failure, and no-clobber/no-repair for non-git entries, dev builds, and dev overrides; estimated cost: cheap unit/behavior tests, no fixtures beyond the existing fakes. Then one existing Pi live journey through the installed-package front door without a development override, started with no valid Spacedock package (fixture: install a wrong-line or bare `git:` source first), per AC-1. Then gofmt, the full suite, race, and a detached adversarial audit because this changes a front door. Documentation impact: none on the docs site (no command-surface change); the doctor remedy wording change, if any, is covered by the behavior tests' output assertions.

Accepted residuals: Pi versions older than the locally spiked one may not honor `@ref` pinning — AC-1's live run exercises the CI-pinned Pi version, and the remedy wording covers the fallback; concurrent `spacedock pi` invocations could race the repair's settings.json rewrite — single-process AC-4 cannot catch this, and the front door is a single-user launch path, so the race window is accepted rather than locked.

### Feedback Cycles

- Cycle 2: revision complete (commit 4a018730e) — all three blockers resolved: (a) pin ref derivation chain linker-stamp → build-info semver tag (proxy builds) → dev sentinel; (b) manifest-version trigger dropped, wrong-ref git-entry detection instead; (c) AC-5 no-repair/no-clobber for non-git entries, dev builds, dev overrides. Independent reviewer re-run verdict: RESOLVED — ready for gate.
- Cycle 3: REJECTED at the implementation point (surface breach, no gate — captain ruling 2026-08-28 chat "2"): implemented ~661 net / 6 files vs approved +100±60 / 4±2 (~4× over ceiling). Captain ruled option 2 (design reset): return through ideation to re-baseline. Correction commits 354f33cb1 + 21730a5d8: Expected surface re-baselined to +661 net / 6 files (tolerance ±40/±1), per-file breakdown + growth attribution added; design/ACs/mechanism byte-unchanged; round implementation/1 recorded (4 entries).

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

## Stage Report: ideation (cycle 2)

- DONE: address the independent review blockers in the task body: (1) resolve the unstamped-proxy-build gap — `go install ...@v0.27.2` carries Version=="dev", so a release-shaped binary would float and recreate the incident; either derive the pin from displayVersion()/the embedded release manifest or explicitly state and accept this class with its consequence; (2) specify the incompatible-mismatch branch — piPackageStatus (pi.go:49-55) has no manifest-version field, so either add its detection to the design and give it a named verifier, or drop the incompatible trigger and say which triggers remain; (3) define repair behavior for non-git sources (file:/npm:/local paths) and for --plugin-dir/SPACEDOCK_REPO_ROOT runs in the repair path (AC-2 covers the install call site only)
  (1) resolved in the source-derivation helper: linker-stamped Version wins; else debug.ReadBuildInfo() main module version when it is a semver tag (covers proxy installs — verified the mechanism exists as stdlib and .goreleaser.yaml:46-60 stamps only goreleaser artifacts); else dev sentinel floats; (2) the manifest-version trigger is dropped — the wrong-ref case is detected as a git-entry ref mismatch against the binary's release ref (settings entry string, already surfaced as packageStatus.source), which needs no manifest field and repairs a pinned-wrong-line package; (3) AC-5 specifies no-repair/no-clobber for file:/npm: entries, dev builds, and --plugin-dir/SPACEDOCK_REPO_ROOT runs.
- DONE: extend the ACs/test plan for the resolved branches: an incompatible-trigger test, a non-git-source no-clobber test, and dev-binary front-door behavior (does a dev build repair at all?), each with a named falsifier
  Wrong-ref repair test (entry @v0.28.0-pre1 under v0.27.2 identity → installs @v0.27.2) in the AC-2/AC-4 behavior set; non-git no-clobber + dev-build no-repair + dev-override suppression are AC-5 with its named falsifier; AC-2/AC-3 tables gained the proxy-tagged identity.

## Stage Report: implementation

- DONE: implement the approved design in the worktree .worktrees/spacedock-ensign-pin-pi-package-to-binary-release: the ordered pin-ref derivation helper (linker-stamped internal/cli.Version, then debug.ReadBuildInfo().Main.Version semver tag for proxy builds, else the dev sentinel), the release-pinned default install source at the install call site, the one-repair-attempt front door (triggers: absent entry, git entry with absent/differing ref) with recheck via the existing discovery, launch refusal on remaining mismatch, and AC-5 suppression for non-git entries, dev-sentinel builds, and --plugin-dir/SPACEDOCK_REPO_ROOT runs
  Commit 178661933 on spacedock-ensign/pin-pi-package-to-binary-release: `internal/cli/pi_package_source.go` (derivation + parse + repair decision), `pi.go` install call site uses `piPinnedSource(piReleaseRef())`, `piRepairAttempt` runs ahead of the readiness gate (proactive on the wrong-ref delta — the v0.27.2 incident shape — because the FO binary gate would abort inside pi), recheck via the existing `checkPiRuntime` discovery, refusal on remaining mismatch, AC-5 suppression for file:/npm:/local entries, dev-sentinel builds, and declared dev overrides.
- DONE: add the command-behavior tests for AC-2/AC-3/AC-4/AC-5 in internal/cli/pi_frontdoor_test.go using the existing fake piRuntimeOps pattern: literal sources per identity (stamped, proxy-tagged, dev, plugin-dir), the wrong-ref repair (@v0.28.0-pre1 under a v0.27.2-identity binary installs @v0.27.2), one-install/one-recheck/no-launch-on-failure, and no-clobber for non-git/dev/override shapes, each with its named falsifier
  `internal/cli/pi_package_repair_test.go` + fake `statusAfterInstall` seam: derivation table (incl. pseudo-version and checkout-build sentinels) flips if any identity's source changes; wrong-ref test pins entry @v0.28.0-pre1 under a 0.27.2 binary and asserts the install of @v0.27.2; failure/ineffective-repair tests assert 1 install, 2 status reads (initial+one recheck), zero Launch calls, and the actionable stderr; suppression tests assert zero installs for file:/npm: entries, dev binaries, and --plugin-dir.
- run gofmt on changed dirs, go test ./..., and go test ./... -race; commit the deliverable to the worktree branch spacedock-ensign/pin-pi-package-to-binary-release with path-scoped adds (never a bare git add -A), and write the stage report as `## Stage Report: implementation` per the ensign discipline with each DONE's evidence on a separate indented line below the marker
  gofmt clean on changed dirs (internal/release/runtime_live_evidence_workflow_test.go was already unformatted on main, untouched); full suite green except the pre-existing env-marker failure TestVersionAmbiguousMarkersExitZero (reproduced on main at the repo root: the sandbox sets PI_CODING_AGENT); go test ./... -race same single pre-existing failure, no races. Commit pushed to spacedock-ensign/pin-pi-package-to-binary-release.

### Summary

Implemented the release-pinned Pi package source and the one-repair-attempt front door per the approved cycle-2 design, with the wrong-ref mismatch (settings entry vs the binary's release ref) as the repair trigger — no package-manifest field needed. The review-found proxy-build gap is closed by the build-info arm of the derivation helper. AC-1's live journey is wired and registered in the runtime live registry; CI's pi-live lane is its proof run (needs CI-E2E-PI approval), so the live run itself remains open as designed.

## Stage Report: ideation (cycle 3)

- DONE: re-baseline the Expected surface and tolerance in the task body's "Expected surface and tolerance" section to the as-implemented figures, per the captain's design-reset ruling (2026-08-28 chat, option 2): actuals are ~661 net across 6 files (662 insertions, 1 deletion) (internal/cli/pi.go +54/-1, internal/cli/pi_package_source.go +150 new, internal/cli/pi_package_repair_test.go +335, internal/cli/pi_frontdoor_test.go +9, internal/ensigncycle/pi_frontdoor_pin_live_test.go +96, docs/runtime-live-ci-registry.md +18) — set the estimate to the as-implemented net and file figures and keep the design, acceptance criteria, and approach unchanged; state in one sentence that the growth over the original estimate is the cycle-2 reviewer-demanded branches' test matrix and the previously unbudgeted live-journey wiring
  "Expected surface and tolerance" now reads +661 net across 6 files, tolerance ±40/±1, with the per-file as-implemented breakdown and one sentence attributing the growth to the cycle-2 reviewer-demanded test matrix plus the previously unbudgeted AC-1 live-journey wiring; the as-implemented figures were independently re-measured from `git diff --numstat origin/main...HEAD` on the worktree branch (18+54+9+335+150+96 insertions, 1 deletion) before writing them, and the design, acceptance criteria, and mechanism sections are byte-unchanged.
- append a `## Stage Report: ideation (cycle 3)` section per the ensign discipline recording only the re-baseline (frontmatter off-limits), with each DONE's evidence on a separate indented line, and commit path-scoped to the state checkout
  This section; frontmatter untouched (`status: implementation`, `id`, `pr:` unchanged — verified by scoped read before commit); commit is path-scoped to pin-pi-package-to-binary-release/index.md only.

### Summary

Measurement-only re-baseline of pin-pi-package-to-binary-release's expected surface to the as-implemented figures (+661 net / 6 files) per the captain's design-reset ruling; design, acceptance criteria, and delivered code unchanged.
