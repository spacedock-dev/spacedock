# Sprint 0200-flip — pre-cut antipattern audit

**Verdict: SHIP-CLEAR — 0 blockers.**

Independent pre-cut antipattern audit of the assembled `next` (sprint
0200-flip, the PRE-FLIP work) BEFORE the outward 0.20.0 flip+cut. Four merged
members (nzb #337, cmx #338, k6d #339, fix #340), assembled base
`85404bdd` → head `4082b90a` (the only non-code commit between is the ignored
roadmap doc `929cac74`). Read-only over code; tests run; adversarial edits done
only in throwaway worktrees (all removed — main checkout confirmed at
`4082b90a` with no tracked modifications). The "tag" here is the FLIP itself,
owned by entity `pj` and HELD for a later, captain-gated step — so a
ship-blocker would split into `fix-now-on-next` (the Commander fixes it on
`next` now) vs `pj-at-flip` (the flip entity handles it). The audit found no
test weakness, no proof-policy violation, no broken cross-change integration,
and no half-done work the flip would freeze. `go test ./...` from the repo root
is GREEN (15 packages `ok`, 0 FAIL, 2 SKIP). One real cross-cutting landmine —
the stamp step retargeted `next`→`main` and armed live on `next` before the
flip — was filed as a candidate ship-blocker, runtime-reproduced, and on
verification of bounded impact DOWNGRADED to record-for-next-sprint (R1).
Nothing must change before the flip.

## Ship-blockers

None.

The integrated-pipeline-landmine dimension filed the live-armed stamp retarget
as a candidate ship-blocker. It was runtime-reproduced (see §3 and R1) and then
downgraded: a pre-flip `v0.19.x` patch cut the documented way WOULD push a wrong
version-stamp commit onto the legacy `origin/main`, but (a) binaries build
correctly from the `next` tip, (b) the marketplace ref is still `next` so
default installs do not resolve from `main`, and (c) the flip's archive+force-
replace absorbs the junk commit — so it freezes nothing and breaks no
user-facing default-install path. Real, but non-blocking. Recorded as R1.

## How each of the 5 dimensions was checked

### 1. Shipped test weakness / proof-policy violations — CLEAR

Every merged test in the four members was read at HEAD `4082b90a` against the
validation section of `docs/dev/README.md` (the "one test that settles it: does
the expected value come from somewhere OTHER than the file under test?",
`README:84`) and the "Instruction-file read quarantine" (`README:123`). The
changed packages ran green and five adversarial edits ran in a throwaway
worktree (`git worktree add --detach /tmp/audit0200-twpp 4082b90a`, since
removed). None is a tautology, a prose-grep over an instruction file, or a
self-referential assertion of the implementer's own text.

- **nzb — predicate + workflow tests.**
  `internal/release/e2egate_test.go` feeds CONSTRUCTED `gh run list` JSON
  fixtures (`greenForCommitJSON:11`, `parkedRunJSON:15`, `greenWrongCommitJSON:19`)
  into the real `EvaluateE2EGate` and asserts the pass/block `Decision` — the
  oracle is the fixture, not the predicate's prose. Adversarial: dropping
  `&& run.HeadSha == releaseCommit` from `e2egate.go:60` reds
  `TestE2EGateBlocksGreenRunOnWrongCommit` (`e2egate_test.go:67`) and
  `TestE2EGatePicksMatchingRunAmongMany` (`:157`) — non-vacuous.
  `cmd/spacedock-release/e2e_gate_test.go` drives the real `runE2EGate` with a
  `fakeRunLister` supplying the EXTERNAL `gh` boundary (not a mock of the
  logic), asserts exit code + that independent inputs (gateCommit fed as argv,
  reason fed via `SPACEDOCK_E2E_GATE_WAIVER`) appear in the
  `$GITHUB_STEP_SUMMARY` tempfile; the `strings.Contains(ToUpper(got),"WAIV")`
  check (`e2e_gate_test.go:79`) tests rendered runtime output of the waived
  path, not an instruction-file grep.
  `internal/release/e2egate_workflow_test.go` +
  the re-aligned `journey_workflow_test.go` anchors parse
  `.github/workflows/release.yml` into a real JOB GRAPH (`parseWorkflowJobs`) and
  prove the goreleaser carrier `needs:` the e2e-gate job that resolves the SHA and
  runs `spacedock-release e2e-gate "$RELEASE_COMMIT"`. release.yml is a CI/build
  artifact parsed in CODE — the README-blessed "parses real artifacts in code"
  case (`README:84`), NOT a prompt/skill/contract the model ingests. The anchors
  exist in the real file (confirmed: `release.yml:99` `e2e-gate:`, `:125` the
  `go run ./cmd/spacedock-release e2e-gate "$RELEASE_COMMIT"` invocation, `:128`
  `needs: e2e-gate`). Adversarial: deleting the `needs: e2e-gate` line from the
  real release.yml reds `TestReleaseWorkflowGatesGoreleaserOnE2E`
  (`e2egate_workflow_test.go:20`) AND the realigned
  `TestReleaseWorkflowJobGraphMatchesGitHubActions`
  (`journey_workflow_test.go:425`) — the realignment is consistent with the new
  edge, non-vacuous.
- **cmx — banner tests.** They render the REAL `launchBanner` and assert on
  rendered bytes; expected phrases are authored independently in the test
  (e.g. `"launching " + host + " as your first officer"`,
  `launch_banner_wording_test.go:44`), never grepped from the `frontdoor.go`
  format string. The AC-5 golden (`TestLaunchBannerSingleWorkflowGolden`, `:122`)
  is byte-exact but RIGHT-grained (the 3-line happy path), built from the package
  `Version` var (an independent runtime value) plus independently-written line
  text. Adversarial: rewording `"first officer"→"firstofficer"` in
  `frontdoor.go` reds both the golden (`:130`) and the metaphor test (`:46`).
  `detectedWorkflow` is now `(label,value)` and tests assert label/content
  agreement (`:104`).
- **k6d — channel-agreement guard.**
  `internal/release/channel_agreement_guard_test.go` parses THREE real artifacts
  (release.yml stamp target, `.goreleaser.yaml` stable/edge `devBranch` ldflags
  via `yaml.Unmarshal`, `.claude-plugin/marketplace.json` `source.ref`) and tests
  the relationship `== main` / `!= next`. The expected value is the independent
  invariant; each surface is authored by a different change. Adversarial:
  stamping `cli.devBranch=next` on the stable build reds
  `TestStableChannelBinaryPairAgreesOnMain` (`:183`,`:186`) and the
  channel-collapse guard in `TestEdgeChannelStampsNext` (`:203`).
  `codex_channel_smoke_test.go` drives the real `runCodex`/`runClaude`
  no-plugin auto-install with `devBranch` set per channel, OBSERVES the branch
  off the recorded install seam, and confirms it threads into `--ref <branch>` /
  `source@branch` — the observed value IS the production argv, the expected
  (`main`/`next`) comes from the case table.
- **fix — `internal/dispatch/launcher_command_test.go`** is test-only hardening
  (`environWithoutSpacedockBin` strips ambient `SPACEDOCK_BIN`, `:63`). Confirmed
  REAL: reverting the fix and running with `SPACEDOCK_BIN=/usr/bin/false` reds the
  `/unset` subcase (`launcher_command_test.go:41`) — the exact ambient condition
  inside a spacedock claude/codex session; the fixed version passes.

**Quarantine boundary holds.** `git diff --stat 85404bdd..4082b90a --
internal/contractlint/**` is empty (no member touched contractlint). No member
added any read of a skill/contract/agent/README instruction file — the only file
reads added are release.yml, runtime-live-e2e.yml, `.goreleaser.yaml`,
marketplace.json (CI/build artifacts) and the `$GITHUB_STEP_SUMMARY` runtime
tempfile. No TODO/FIXME/HACK/XXX in any changed non-md file. Exactly one
`t.Skip` (`channel_agreement_guard_test.go:220`) — the documented pre-flip/flip
split: `TestTriSurfaceChannelAgreement` skips RED-by-construction because
`marketplace.json` `source.ref` is still `next` (`marketplace.json:12`) and pj
owns flipping surface 3 at the flip; the binary-side pair (surfaces 1+2 ==
`main`) is asserted UNCONDITIONALLY by `TestStableChannelBinaryPairAgreesOnMain`,
so the skip masks no binary-side hole. Runtime:
`go test ./internal/cli/ ./internal/release/ ./cmd/spacedock-release/ ./internal/dispatch/`
all `ok`.

### 2. Cross-change integration (release.yml × internal/release × Version='dev') — CLEAR

The three cross-change seams the brief names were audited.

**release.yml (nzb e2e-gate + k6d stamp-swap, merged sequentially).** The merged
job graph is `e2e-gate` (entry, no needs) → `goreleaser` (`needs: e2e-gate`,
`release.yml:127-128`) → `journey-ledger` (`needs: goreleaser`,
`release.yml:24-25`) — confirmed by grep of the head file. The two merges edited
DISJOINT regions: `git show 4862c2ab -- .github/workflows/release.yml` shows k6d
touched ONLY the "Stamp plugin manifests" step (comment rewrite + next→main on
the git fetch/switch/push at `release.yml:186-217`); a grep of k6d's diff for
`e2e-gate`/`needs:`/`goreleaser:` returns nothing, so k6d did not disturb nzb's
e2e-gate wiring (added at `cd4fd457`). The document-order builder<goreleaser
guard (`assertReleaseWorkflowPublishesJourneyCosts`,
`workflow_exec_guard_test.go:95-145`) still holds because it parses steps
job-unaware in text order — the journey-cost builder sits in journey-ledger,
e2e-gate's step before the goreleaser action, so builderStep<goreleaserStep
regardless of the e2e-gate step now between them. `isJourneyCostBuilder`
(`workflow_exec_guard_test.go:250-254`) keys on `journey-costs`, so it does NOT
false-match the e2e-gate command (`spacedock-release e2e-gate`).

**internal/release (nzb e2egate*.go + journey_workflow anchor re-align; k6d
channel_agreement_guard_test.go).** `go test ./internal/release/ -count=1` → `ok`
(75 PASS, 1 SKIP at runtime). No test-name collision across the four
changed/added files. The journey_workflow fixtures (re-aligned by nzb in the SAME
commit that added e2e-gate) anchor on the literal goreleaser-header
`needs: e2e-gate`, and `TestReleaseWorkflowJobGraphMatchesGitHubActions` asserts
goreleaser needs EXACTLY `[e2e-gate]` and journey-ledger needs EXACTLY
`[goreleaser]` — both PASS. The separation guards
(`TestReleaseWorkflowGuardRejectsGoreleaserNeedsJourneyLedger` + ViaJobIdentity)
and the e2e-gate wiring guard coexist and PASS. Shared parser helpers
(`parseWorkflowSteps`/`parseWorkflowJobs`/`executableShellCommands`,
`readGoreleaserConfig`) each have exactly ONE definition site — no helper-
redefinition collision; `go build ./...` OK.

**Version='dev' interaction (cmx vs nzb/k6d).** cmx changed `internal/cli/cli.go:30`
`var Version = "dev"` (from `"0.19.0"`). The `.goreleaser.yaml` stable+edge builds
both stamp `-X …cli.Version={{ .Version }}` (git-describe) at `.goreleaser.yaml`
stable/edge ldflags — so `"dev"` is the unstamped sentinel overwritten at
release-build. nzb's release.yml version-stamp uses
`RELEASE_VERSION="${GITHUB_REF_NAME#v}"` (a shell var feeding `stamp-version` on
plugin.json), wholly separate from `cli.Version`. No collision. cmx's
`launch_banner_wording_test.go:126` renders the banner with the LIVE `Version` var
(not a hard-coded literal) and `TestUnstampedVersionIsNotARelease` asserts
`Version=="dev"` — both PASS. k6d's `codex_channel_smoke_test.go` mutates
`devBranch` (not Version) with save/restore; neither cmx's nor k6d's mutating
tests call `t.Parallel`, so no shared-state race on `Version`/`devBranch`. Full
`internal/cli` package: `ok`. `cmd/spacedock-release/main.go` carries
stamp-version + journey-costs + e2e-gate subcommands together.

### 3. Integrated-pipeline landmine — FINDINGS (recorded R1/R2; not blocking)

The full post-sprint pipeline was read as one system at head `4082b90a`
(`.github/workflows/release.yml`, `.goreleaser.yaml`), tracing a `v*` tag
end-to-end. Two sub-questions came back CLEAN; one real landmine surfaced.

**Job graph / e2e-gate scope — CLEAR.** release.yml triggers only on
`push: tags: ['v*']`. `goreleaser` is ONE job that emits BOTH channels —
`.goreleaser.yaml` has two builds `spacedock-stable` (devBranch=main) and
`spacedock-edge` (devBranch=next), two archives, two casks (`spacedock` +
`spacedock@next`) — all produced by the single `goreleaser release` invocation.
So the single `needs: e2e-gate` (`release.yml:128`) gates the WHOLE two-channel
run, not one build. `EvaluateE2EGate` (`e2egate.go:43-58`) matches
`conclusion==success && headSha==releaseCommit`, channel-agnostic, bound to the
tagged commit. One gate, whole run — correct, not a hole.

**Tri-surface skip — CLEAR / acceptable.** `TestStableChannelBinaryPairAgreesOnMain`
and `TestEdgeChannelStampsNext` assert the binary-side pair (release.yml stamp
target == `.goreleaser.yaml` stable devBranch == `main`) and channel-separation
UNCONDITIONALLY; `TestTriSurfaceChannelAgreement` `t.Skipf`s while
marketplace.json `source.ref` is `next` (pj owns surface 3). Runtime:
`go test -run TestStableChannelBinaryPairAgreesOnMain|TestEdgeChannelStampsNext|TestTriSurfaceChannelAgreement -v`
→ PASS / PASS / SKIP with the exact message `marketplace source.ref = "next", not
yet flipped to "main" (pj owns surface 3); binary-side pair is covered by
TestStableChannelBinaryPairAgreesOnMain`. Three independently-authored real
artifacts, binary pair covered without the skip — a legitimate, documented
disposition, NOT an illegitimate skip hiding unproven behavior.

**THE PRE-FLIP-CUT LANDMINE (the finding, recorded as R1).** The stamp step
changed from `git switch next`/`git push origin next` (confirmed at
`git show 85404bdd:.github/workflows/release.yml:160-172`) to `git fetch origin
main`/`git switch main`/`git push origin main` (`release.yml:203-217` at head).
This is LIVE on `next` NOW, before the flip. The live, documented cut procedure
still says cut-from-next: `AGENTS.md:28` — "Cut releases from `next` … Never
release from `origin/main`." Tag history confirms nine consecutive patches
v0.19.0–v0.19.9 cut from `next`; v0.19.9 (`2ffe9c9b`) was cut BEFORE k6d
(`4862c2ab`) landed, so its stamp correctly targeted `next`; the NEXT patch hits
the retargeted step. `origin/main` is the legacy 0.12.1 tip (`8c069d95`),
confirmed NOT an ancestor of `next` (`git merge-base --is-ancestor origin/main
4082b90a` = NO); main's plugin.json is `0.12.1`, next's is `0.19.9`. I
reproduced the exact CI step in a detached worktree
(`/tmp/audit0200-final` at `4082b90a`, since removed): built the stamp tool from
the tag, then `git switch main`, then `stamp-version 0.19.10
.claude-plugin/plugin.json .codex-plugin/plugin.json` rewrote main's
`"version": "0.12.1"` → `"0.19.10"` (DIFF PRESENT → the step commits + `git push
origin main`). Severity is bounded: binaries build from the correct next tip; the
marketplace ref pre-flip is `next` so default installs don't resolve from the
corrupted main; the flip's archive+force-replace absorbs it. It does NOT freeze
the flip and does NOT break a user-facing default-install path — but it WILL push
a wrong, published version-stamp commit (0.12.1 code labeled 0.19.10) onto the
legacy stable branch on any pre-flip patch cut done the documented way. This is
the catalog "premature behavior live on next before the flip" antipattern, on a
release branch, armed. k6d's body
(`docs/dev/.spacedock-state/_archive/two-channel-release-devbranch-stamp.md`)
reasons ENTIRELY about the post-flip agreement invariant and never addresses the
pre-flip-cut window; the per-task audit could not see this cross-cutting hole
because it only emerges from the AGENTS.md cut-from-next instruction × the
retargeted stamp × the divergent legacy main. The
AGENTS.md/`docs/releasing.md` contradiction itself is pre-existing and explicitly
pj-owned (flip index AC-4), so it is NOT charged here — only the live armed stamp.

**Secondary (recorded R2).** `next-publish.yml` stamps ONLY the marketplace
CALENDAR key (`bump-calendar`, confirmed `next-publish.yml:15-29`), NOT plugin.json
`version`. Pre-sprint the stamp step updated next's plugin.json `version` on each
cut. After k6d the stamp targets `main`, so next's plugin.json `version` freezes
at 0.19.9 — no workflow updates it. The `release.yml:195` comment ("the edge
channel's version-display rides the unchanged next-publish.yml") is imprecise:
next-publish.yml does not touch plugin.json `version`. Cosmetic version-panel
staleness on the edge channel during the pre-flip window; no install/resolve
impact.

`go test ./internal/release/ ./internal/cli/ ./skills/integration/` is GREEN.

### 4. `go test ./...` green from the repo root + skip/TODO hygiene (DoD#2) — CLEAR

`go test ./... -count=1` from the repo root at HEAD `4082b90a` (clean): 15
packages `ok`, 0 FAIL, 2 SKIP (uncached). Re-ran the four members' packages
explicitly — `internal/release` ok, `internal/cli` ok, `cmd/spacedock-release`
ok, `internal/dispatch` ok.

**#340 (the internal/dispatch hermeticity fix) validated as a REAL fix, not
cosmetic.** The breakage condition is an executable ambient `SPACEDOCK_BIN`
leaking into the `unset` subcase of
`TestLauncherCommandFallsBackToPathWhenSpacedockBinUnsetEmptyOrUnusable`
(`launcher_command_test.go:35,41-44`). Reproduced the FAIL at pre-#340 commit
`4862c2ab` with an executable ambient bin (`--- FAIL: …/unset … launcher command
output = "ambient-bin:dispatch", want "path-bin:dispatch"`); at HEAD the same
condition PASSES (all four subcases unset/empty/non-executable/missing green).
The fix is `environWithoutSpacedockBin()` (`launcher_command_test.go:63-77`),
used at `:52`. The full suite from repo root ran green under three env
conditions — default, `env -u SPACEDOCK_BIN`, and `SPACEDOCK_BIN=<executable>`.

**Skip/TODO hygiene.** Scanned all 17 changed code/test files (the 19-file diff
minus the two ignored roadmap docs) for TODO/FIXME/XXX/HACK — zero hits. Exactly
ONE new `t.Skip` across all changed test files:
`channel_agreement_guard_test.go:220`. It is conditioned on REAL on-disk state —
reads the live marketplace.json `source.ref` and skips only while that ref !=
`main`; confirmed the live ref at HEAD is `next`, so the skip is firing because
surface 3 genuinely is not yet flipped. In a throwaway worktree, flipping the
ref to `main` UN-SKIPS the test (which then PASSES, because k6d already put the
binary-side surfaces on `main`); introducing a real drift (stable build
`cli.devBranch=main→next`) reds it loudly (`channel surface ".goreleaser.yaml
stable devBranch" = "next", want "main"`). It parses three independent real
artifacts, the binary-side pair is asserted UNCONDITIONALLY by
`TestStableChannelBinaryPairAgreesOnMain` — a legitimate state-absent guard, not
hiding unproven behavior. The second repo-wide skip
(`TestUpdateFrontmatterBlockScalarsRewrapped`,
`internal/status/node_roundtrip_test.go:270`) is OUT OF SCOPE: not in the audited
diff, pre-existing at base `85404bdd` with a documented reactivation condition.

### 5. Flip-readiness / scope integrity — CLEAR (provenance nit recorded R3)

Audited the three sub-dimensions against base `85404bdd` → head `4082b90a`.

**(a) SCOPE INTEGRITY — airtight.** `git diff --name-only 85404bdd..4082b90a`
yields 19 files; grep for `marketplace.json|marketplace_manifest_test.go|next-publish.yml`
returns NOTHING — "NO LEAK - clean". The two `docs/` files in range
(`docs/roadmap/0200-flip/dispatch-sprint-execution.md`, `docs/roadmap/README.md`)
are EXACTLY the two the ignored roadmap commit `929cac74` touched. Positively
verified the pj-owned surfaces are pristine: `.claude-plugin/marketplace.json:12`
still `"ref": "next"`; `skills/integration/marketplace_manifest_test.go:67` still
`if src.Ref != "next"`; `next-publish.yml` UNTOUCHED.

**(b) ENTITY STATE — confirmed all four.**
`_archive/{gate-release-on-e2e,two-channel-release-devbranch-stamp,frontdoor-launch-banner-ux}.md`
frontmatter all show `status: done` / `verdict: PASSED` and are archived.
`gh pr view` confirms 337/338/339/340 all `MERGED`. The fourth entity (fix, #340)
is a test-only hermetic change with no standalone archived body — appropriate for
its scope. Each high-stakes surface's detached adversarial audit is on the record:
nzb (7 adversarial edits all caught RED), k6d (stamp mismatch / switch-push split
/ channel collapse / dropped --ref all caught, tri-surface skip shown honest),
cmx (5 reverted-wording edits all caught over rendered bytes).

**(c) FLIP HAND-OFF — green-by-construction, PROVEN at runtime.** On disk at
head: `release.yml:204/215` `git switch main`/`git push origin main` (stamp
target = main); `.goreleaser.yaml` stable `cli.devBranch=main`, edge
`cli.devBranch=next`. Runtime: binary-pair PASS, edge-stamps-next PASS,
tri-surface SKIP with the exact `pj owns surface 3` message. In a throwaway
worktree at `4082b90a` I applied pj's flip (marketplace.json ref `next→main` +
the paired `marketplace_manifest_test.go:67` edit `next→main`): tri-surface
un-skipped to PASS and the manifest test stayed PASS — confirming pj's ONE ref
edit (plus paired test) is sufficient. Inverse honesty check: with ref=main but a
drifted stable devBranch, BOTH `TestStableChannelBinaryPairAgreesOnMain` and the
un-skipped `TestTriSurfaceChannelAgreement` FAIL — the skip is bounded to the
ref-still-next case, not a blanket escape. nzb's e2e-gate is wired
(`release.yml:99` `e2e-gate:` job, `:128` goreleaser `needs: e2e-gate`, `:125`
`go run ./cmd/spacedock-release e2e-gate "$RELEASE_COMMIT"`, `:120`
`SPACEDOCK_E2E_GATE_WAIVER`) — the exact path the pj runbook relies on. No
premature flip behavior: `cli.go:30` `Version = "dev"` sentinel, `frontdoor.go:50`
`devBranch = "next"` unstamped default (NOT prematurely main, not changed by the
diff). pj's outstanding obligations (ref flip + paired test, `AGENTS.md:28` +
`docs/releasing.md` reconciliation, calendar-bump-on-main, archive/flip runbook,
upgrade journeys) are all correctly assigned to pj-at-flip per the dispatch doc
and flip index AC-4 — no gap the sprint should have delivered.

## Record for next sprint

Non-blocking observations seeded for the Shaping FO's Close step — none affects
the flip.

- **R1 — stamp step retargeted to `main` is LIVE on `next` while the documented
  cut procedure still cuts from `next`; a pre-flip 0.19.x patch pushes a wrong
  version-stamp commit to the legacy `origin/main`.**
  `@ .github/workflows/release.yml:203-217`. k6d swapped the "Stamp plugin
  manifests" step from `git switch next`/`git push origin next` to `git fetch
  origin main`/`git switch main`/`git push origin main` and merged it to `next`
  BEFORE pj flips — the step is armed now. `AGENTS.md:28` (the live, pre-flip cut
  procedure) still says "Cut releases from `next` … Never release from
  `origin/main`", and nine patches v0.19.0–v0.19.9 were all cut from `next`.
  `origin/main` is the divergent legacy 0.12.1 tip (`8c069d95`), not an ancestor
  of `next`. So a normal pre-flip v0.19.10 patch cut from `next` runs goreleaser
  (publishes correct binaries from the next tip), then `git switch main` (lands
  on the 0.12.1 tip), stamps main's plugin.json `0.12.1 → 0.19.10` (runtime-
  reproduced: DIFF present), commits, and `git push origin main` — publishing a
  misleading stamp (0.12.1 code labeled 0.19.10) onto the legacy stable branch.
  The flip's archive+force-replace absorbs the junk commit so the flip is not
  frozen, and the marketplace ref pre-flip is still `next` so default installs are
  unaffected — but the step writes wrong published state to the shared `main`
  branch the moment anyone cuts a patch the documented way. Fix it (guard the
  stamp or freeze 0.19.x cuts) before any pre-flip cut can fire. *(Filed as a
  candidate ship-blocker; downgraded on verification of bounded impact.)*

- **R2 — after the stamp moves to `main`, no workflow stamps `next`'s
  plugin.json `version`; the edge plugin-panel version freezes at 0.19.9; the
  release.yml:195 comment is imprecise.**
  `@ .github/workflows/release.yml:195` and `.github/workflows/next-publish.yml:15-29`.
  AC-4's purpose was that the host plugin panel displays a version tracking the
  binary; pre-sprint, the stamp step updated `next`'s plugin.json `version` on
  each `v*` cut. With the stamp retargeted to `main`, next's plugin.json `version`
  is no longer updated by any workflow — `next-publish.yml` bumps ONLY the
  marketplace CALENDAR key (`0.0.YYYYMMDDNN`), not plugin.json `version`. So the
  edge/`next` channel's displayed plugin version freezes at 0.19.9 until the flip.
  The `release.yml:195` comment "the edge channel's version-display rides the
  unchanged next-publish.yml" is imprecise — next-publish.yml does not touch
  plugin.json `version`. Cosmetic version-panel staleness on the edge channel
  during the pre-flip window; no install/resolve impact.

- **R3 — archived entity bodies have empty `issue:` frontmatter; PR provenance
  lives only in commit subjects/GitHub, not the queryable state.**
  `@ docs/dev/.spacedock-state/_archive/gate-release-on-e2e.md`,
  `two-channel-release-devbranch-stamp.md`, `frontdoor-launch-banner-ux.md`
  (all `issue:` blank). The PR refs (#337/#338/#339) that bind each entity to its
  merged change exist only in the squash-merge commit subjects and on GitHub —
  not in the on-disk Spacedock state. A reviewer reconstructing the flip's
  provenance from `spacedock status` alone cannot map an archived entity to its PR
  without leaving the state store. This does not block the flip (the binding is
  recoverable from git/gh, and all four PRs are MERGED), but it is a
  provenance-trail gap. The fix entity (#340) additionally has no standalone
  archived body at all — acceptable for a test-only change, but worth a one-line
  note in whichever sprint record closes 0200.
