---
id: qcwrzza9xkr5kfwmdmvbqhkv
title: Pin the installed Pi package to the launcher release
status: validation
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
        - id: gate:qcwrzza9xkr5kfwmdmvbqhkv:validation
          stage: validation
          attempts:
            - id: gate-attempt:qcwrzza9xkr5kfwmdmvbqhkv-validation-1
              briefing:
                id: briefing:qcwrzza9xkr5kfwmdmvbqhkv:validation:attempt-1:revision-1
                digest: sha256:736b7d17f20e1a7195d3e02f36c8dd8b35592de855349e75e3fd2aa1251191f9
                room-ref: ./review/validation/briefing-1
              withdrawal:
                by: agent:first-officer
                at: "2026-08-28T21:42:13.356773Z"
                reason: 'Captain re-scope (2026-08-28 chat): proof may come from executed use of the code with behavior observation; the 335-line fake-seam flow suite, 96-line live journey, and 18-line registry are cut. AC verifiers re-scoped to hand-executed scenario runs (retained artifacts) + pure-function tables. Routing to ideation for AC re-scope, then implementation shrink.'
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

Estimate net LOC change: +295, across 3 files. Tolerance: +/-80 net LOC and +/-1 file (re-baselined per the captain's 2026-08-28 AC re-scope ruling: proof comes from executed use of the code with behavior observation, so the fake-seam flow suite, the dedicated live journey, and the registry wiring leave this task's surface; the automated surface is the pure-function tables only).

Post-shrink files: `internal/cli/pi.go` (+54/−1; pinned-source derivation call site, one front-door repair attempt with recheck and refusal, repair suppression for non-git entries and dev builds), `internal/cli/pi_package_source.go` (+150 new, the source-derivation helper: linker stamp, then build-info tag for proxy builds, then the dev sentinel), one new pure-function test file (`internal/cli/pi_package_source_test.go`, ~+90: derivation table, source/ref classifier table, repair-trigger decision table — no fake-runtime scaffolding). Cut from the delivered 178661933 by the implementation correction round: the +335 fake-seam flow suite (`pi_package_repair_test.go`), the +9 `pi_frontdoor_test.go` seam extension, the +96 live journey (`internal/ensigncycle/pi_frontdoor_pin_live_test.go`), and the +18 registry registration. Observable semantics that may change: stored format (the `packages` entry for spacedock in `~/.pi/agent/settings.json` gains an `@<ref>` suffix for release-stamped binaries) and front-door runtime behavior (one repair attempt with launch refusal on remaining mismatch). Command grammar is unchanged; no new flags. Authority is unchanged. Doctor remedy lines may gain wording pointing at the pinned repair.

**AC-1 (VALUE) - A released launcher repairs and boots the Pi skill suite from the same release line.**
Verified by: one hand-executed front-door run at validation — the ordinary installed front door with no development override (no `--plugin-dir`, no `SPACEDOCK_REPO_ROOT`), against a fixture pi-home seeded with the incident state (Spacedock package absent, unpinned, or wrong-line), observing exactly one `pi install` of the pinned source, the resulting `settings.json` carrying the binary's own `@<ref>` entry, and a launch-ready report. Run transcript + resulting `settings.json` are retained in the entity as durable artifacts. Independent baseline that must move the wrong way without the fix: the same start state with the repair removed leaves the package missing/unpinned and the FO binary gate aborts before workflow work (the observed v0.27.2 failure).

**AC-2 - The default Pi install source is release-pinned for every release-shaped identity, and the development override remains local.**
Verified by: the kept pure-function tables (`piReleaseRefFrom` derivation rows for linker-stamped, proxy-tagged (`Version=="dev"` + build-info `v0.27.2`), prerelease, dev-sentinel, and pseudo-version identities) plus the executed front-door runs as the seam proof — the run transcript records the source `pi install` was invoked with. Falsifier: changing the release source to omit its ref, letting a proxy-tagged build float, or making the dev-override path use the release source, flips a table row or an observed run.

**AC-3 - The release fix is safe for v0.27.3, the v0.28 prerelease line, and proxy-installed launchers.**
Verified by: the derivation table rows over stable (`v0.27.3`), prerelease (`v0.28.0-pre1`), unstamped (`dev`), and proxy-tagged identities asserting literal expected refs (stable/prerelease/proxy pin to their own tag; dev floats); the executed front-door run adds the seam observation for the stamped identity. Falsifier: changing the source derivation for any covered identity.

**AC-4 - A failed or ineffective repair cannot launch Pi.**
Verified by: hand-executed scenario runs observing behavior — (a) install failure (a PATH-shimmed `pi` recorder that fails the install, the same technique as the CI PATH shim): exit 1, no launch, actionable stderr naming the remedy; (b) remaining mismatch after the repair: exit 1, no second install, no launch. Observed exit code, stderr, and (for no-launch) the pi sessions dir / recorder log. Falsifier: a second install attempt, a launch, or exit 0 after remaining mismatch.

**AC-5 - The repair never rewrites a user-managed or dev-owned package surface.**
Verified by: hand-executed suppression runs — a `file:` settings entry, a `dev`-sentinel binary, and a `--plugin-dir` run: the entry string is byte-unchanged afterward, zero install calls in the recorder log, and launch/exit behavior per shape (the first-officer binary gate remains the version backstop). Falsifier: a repair install attempted for any of those shapes, or an entry string rewritten, visible as a changed `settings.json`.

## Test plan

Automated: one new pure-function test file `internal/cli/pi_package_source_test.go` (~90 lines, no scaffolding) — the `piReleaseRefFrom` derivation table (10 identity rows incl. proxy-tagged and pseudo-version float), the `piPinnedSource`/`piGitSourceRef` classifier table (git/file/npm/path/empty), and the `piPackageNeedsRepair` decision table (missing/unpinned/wrong-line/own-ref/non-git/dev). Falsifier: changing any derivation arm, classifier class, or trigger state flips a table row.

Hand-executed scenario matrix at validation (executed use of the code, behavior observation; never prose grep): with a PATH-shimmed `pi` recorder (the CI-shim technique — logs install/status/launch invocations and can fail on demand) or the real Pi against a fixture `piHome`/`HOME`, run the stamped binary's ordinary front door and observe: (1) incident repair — unpinned entry → exactly one install of `@<own ref>`, `settings.json` entry rewritten to the pinned source, launch proceeds; (2) wrong-line entry `@v0.28.0-pre1` under a 0.27.2-identity binary → repairs to the binary's own ref; (3) AC-4 install-failure arm → exit 1, no launch, actionable stderr naming `spacedock doctor --host pi`; (4) AC-4 remaining-mismatch arm → exit 1, exactly one install, no launch; (5) AC-5 suppression: `file:` entry byte-unchanged, dev-sentinel build no install, `--plugin-dir` run suppressed; (6) healthy pinned entry → zero installs, launch proceeds. Each run retains the recorder log, stderr, and resulting `settings.json` in the entity's artifacts. Then gofmt, the full suite, and race. Documentation impact: none on the docs site (no command-surface change); the doctor remedy wording change, if any, is covered by the executed runs' output observations.

Accepted residuals: Pi versions older than the locally spiked one may not honor `@ref` pinning — the hand-executed incident run at validation exercises the locally installed Pi version, and the remedy wording covers the fallback; concurrent `spacedock pi` invocations could race the repair's settings.json rewrite — the front door is a single-user launch path, so the race window is accepted rather than locked.

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

## Stage Report: validation

- DONE: independently verify the deliverable at commit 178661933 against each AC: run the AC-2/AC-3/AC-4/AC-5 command-behavior tests and confirm each named falsifier would flip, and confirm gofmt, full go test ./..., and go test ./... -race results (the pre-existing env-marker failure verified on origin/main before attribution)
  All 10 tests in internal/cli/pi_package_repair_test.go PASS verbosely (derivation table, pinned-source parsing, needs-repair table, per-identity literals, unpinned/wrong-line repair, install-fail refusal, ineffective-repair refusal, suppression, healthy-pinned no-repair). Falsifier audit by reading the assertions: AC-2/AC-3 pin literal sources per identity (`git:github.com/spacedock-dev/spacedock@v0.27.2` for stamped AND proxy-tagged via withPiReleaseIdentity("0.27.2","(devel)") / ("dev","v0.27.2"); dev sentinel asserts the bare floating source; plugin-dir asserts "/checkout" and never the release source) — omitting a ref, floating a proxy build, or routing the override to the release source flips a t.Fatalf literal. AC-4 asserts len(piInstalls)==1, launched==0, ops.statusCalls==2 (initial+one recheck), exit 1, and stderr naming the remedy (`spacedock doctor --host pi`) for both install-failure and ineffective-repair shapes. AC-5 asserts zero installs for file:/npm: entries, dev-sentinel builds (exit 1 with the pre-existing refusal stands, zero installs), and --plugin-dir runs; the healthy-pinned test adds the no-op positive. AC-1's live journey TestLivePiFrontDoorInstallsPinnedPackage builds under -tags live and is registered at docs/runtime-live-ci-registry.md:336 (pi-front-door-pinned-package) — the CI pi-live lane is its proof run, open by design. gofmt clean on changed dirs; full suite with ambient markers unset: ok across all packages (rc=0, -json verified); race: ok (cli 181s, ensigncycle 251s). TestVersionAmbiguousMarkersExitZero with PI_CODING_AGENT set FAILs identically on origin/main (verified at the repo root) — pre-existing environment-marker failure, not this branch's regression.
- DONE: audit the repair-path code in internal/cli/pi_package_source.go and internal/cli/pi.go against the approved design: ordered derivation (linker stamp, build-info semver tag, dev sentinel), one repair attempt, recheck, launch refusal on remaining mismatch, AC-5 suppression for non-git entries/dev builds/--plugin-dir/SPACEDOCK_REPO_ROOT, and that no non-git entry or dev-build clobber path exists
  piReleaseRefFrom (pi_package_source.go:47) is the ordered chain: piSemverTag(linkerStamp) → build-info main version gated to semver-not-pseudo (piPseudoVersion rejects go-install @next pseudo-versions, so @next resolutions float per design) → "" sentinel; piPinnedSource("") returns the bare floating source. piPackageNeedsRepair (pi_package_source.go:126) triggers only for a release-shaped binary with missing/unpinned/wrong-ref git entry; non-git entries return ok=false (never rewritten); dev sentinel returns false. piRepairAttempt (pi.go:632) suppresses on cfg.repoRoot != "" (covers both --plugin-dir and SPACEDOCK_REPO_ROOT per piRuntimeConfigFromEnv:528) and !check.piBinOK, spends exactly one PiInstall, returns true so the caller rechecks once; refusal on remaining mismatch is the existing launch gate (test-pinned: exit 1, zero launches, one recheck, actionable stderr naming doctor). No clobber path exists for dev builds (no install call at all) or non-git entries (piGitSourceRef returns ok=false before any rewrite path).

### Summary

Independent validation of pin-pi-package-to-binary-release's deliverable (178661933) against the re-baselined design: all AC-2..AC-5 behavior tests pass with their named falsifiers audited as literal-pinning, the repair-path code matches the approved design (ordered derivation, one repair attempt, recheck, refusal, AC-5 suppression), gofmt/full-suite/race are clean (the sole failure is the pre-existing environment-marker failure, verified identical on origin/main), and AC-1's live journey is wired, registry-registered, and left to the CI pi-live lane as designed. The entity remains nonterminal.

## Stage Report: ideation (cycle 4)

- DONE: re-scope the acceptance criteria verifiers per the captain's ruling (proof from executed use of the code with behavior observation, never prose grep): AC-1's verifier is now one hand-executed front-door run against a fixture pi-home seeded with the incident state, retaining run artifacts and the resulting settings.json showing the binary's own @ref registered and a launch-ready report; AC-4's verifier is hand-executed refusal runs (install-failure and remaining-mismatch arms) observing exit 1, no launch, actionable stderr; AC-5's verifier is hand-executed suppression runs (file: entry, dev build, --plugin-dir) observing the byte-unchanged entry and zero installs; AC-2/AC-3 consolidate into the kept pure-function tables (derivation, classifier, repair-trigger) with the executed runs as the seam proof; the dedicated live journey and its registry registration are out of this task's surface
  AC-1..AC-5 verifier paragraphs and the Test plan rewritten in place (index.md lines 118-141 span); AC-1 names the fixture-pi-home hand run with retained artifacts; AC-4/AC-5 name observed exit/stderr/on-disk state and the recorder log as evidence; the pure-function tables (`piReleaseRefFrom` derivation, `piPinnedSource`/`piGitSourceRef` classifier, `piPackageNeedsRepair` decision) carry AC-2/AC-3 with falsifier rows unchanged in meaning (omitting a ref, floating a proxy build, routing the override to the release source still flips a literal)
- DONE: re-baseline the Expected surface and tolerance to the post-shrink estimate (~300 net / 3 files: production ~204, tests ~90) with real slack (+/-80 net, +/-1 file), replacing the self-referential +661/6 baseline, and update the test plan to name the hand-executed scenario matrix (commands, fixture setup, observed exit/stderr/on-disk state, retained artifacts) that validation runs
  Expected surface now reads +295 net / 3 files, tolerance ±80/±1, with the post-shrink file list (pi.go +54/−1, pi_package_source.go +150, pi_package_source_test.go ~+90) and the named cut set (flow suite 335, seam extension 9, live journey 96, registry 18 — removed by the implementation correction round); the test plan names the hand-executed matrix (fixture piHome, stamped binary via -ldflags, PATH-shimmed `pi` recorder for install/status/launch observation, observed exit codes/stderr/settings.json, artifacts retained under the entity).
- DONE: append ## Stage Report: ideation (cycle 4) per the ensign discipline (evidence on separate indented lines, frontmatter off-limits), commit path-scoped to the state checkout
  This section; frontmatter untouched (status: validation, id, pr: unchanged — verified by scoped read before commit).

### Summary

Re-scoped the verification design per the captain's ruling that proof comes from executed use of the code with behavior observation: AC-1/AC-4/AC-5 move to hand-executed scenario runs with retained artifacts at validation, AC-2/AC-3 consolidate into the pure-function tables, and the fake-seam flow suite, live journey, and registry wiring leave the surface. Expected surface re-baselined to ~+295 net / 3 files (±80/±1) — a prediction again, not a mirror of actuals.
