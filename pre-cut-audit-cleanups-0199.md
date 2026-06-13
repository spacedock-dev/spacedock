---
id: 5ar2193yw8sv0rcyrt23wxg9
title: Pre-cut audit cleanups (0.19.9) — checksum-gate test, darwin-only doc drift, gofmt, hasGitEntry comment, node-action bump
status: ideation
source: "0.19.9 pre-cut antipattern audit (Commander staff review, 2026-06-08) — four recorded non-blockers, none gated the 0.19.9 tag. Seeded per the roadmap Close step."
started: 2026-06-13T05:52:37Z
completed:
verdict:
score:
worktree:
issue:
sprint: 0202-survey-improvements
group: cleanup
sprint-readiness: ready
---

Four non-blocking findings the 0.19.9 pre-cut antipattern audit recorded (none blocked the cut; grouped here as the next-sprint seed). Small, independent.

## Items

1. **Checksum-gate fail-closed has no `*_test.go`.** v3's `curl | sh` installer checksum check is proven only in `install-e2e.yml` — no Go test loads it, so a contributor could delete the checksum verification and `go test` stays green. Add a `workflow_exec_guard`-style test (or a Go e2e tamper test) that goes RED if the checksum gate is removed/weakened. *This is the substantive one — a real test-strength hole in shipped install machinery.*

2. **Darwin-only prose is stale** in `docs/releasing.md` + the `release.yml` header — 0.19.9 added Linux binaries, but both still describe a darwin-only release. Doc fix. *(Backstop: `pj`'s AC-4 doc-reconciliation re-verifies `docs/releasing.md` against the real machinery at the flip; this fixes the platform-language drift sooner so the doc isn't inaccurate in the interim.)*

3. **gofmt drift** on `skills/integration/survey_sync_codex_test.go` (pre-existing, cosmetic). `gofmt -w` to clear.

4. **`hasGitEntry`'s only unmasked guard** is the single `TestDiscoverWorkflowsSkipsNestedCheckout` — add a cross-reference comment so a future edit knows that test is what protects it.

5. **Node-20 GitHub Actions deprecation — time-sensitive (≈2026-06-16).** The `v0.19.9` release run warned that `actions/checkout@v4`, `actions/setup-go@v5`, and `goreleaser/goreleaser-action@v6` are forced to Node-24 starting **2026-06-16** ("may not work as expected"). Bump these actions (and any sibling Node-20 actions across the workflows) to Node-24-compatible versions **before** the 0.20.0 flip cut runs `release.yml` — or before 2026-06-16, whichever comes first. Priority alongside #1, because the flip cut depends on a healthy `release.yml`.

## Out of scope

Anything requiring a behavior redesign. Each item is a localized test/doc/format/CI fix.

## Ideation findings (verified against the tree, 2026-06-12)

Each item was checked against the real machinery before fixing the spec. The audit's framing was right on substance but imprecise on two spots; the corrections are recorded here so the implementation does not "fix" already-correct text.

- **#1 (checksum gate).** `install.sh`'s checksum gate is `install.sh:161-169`: pull this asset's expected hash from `checksums.txt`, compute the tarball's actual sha256, `die` on mismatch or a missing checksum line. It is proven ONLY by `install-e2e.yml`'s tamper leg, which needs a goreleaser snapshot and runs on PR — `go test` stays green if lines 164-169 are deleted. The gate is exercisable from Go WITHOUT goreleaser: install.sh's `SPACEDOCK_INSTALL_FROM=<dir>` path (`install.sh:115-134`) consumes a plain `dist/` dir holding a `spacedock_<ver>_<os>_<arch>.tar.gz` (a tar with a bare `spacedock` file at root) plus a `checksums.txt`. A Go test can build that fixture in a temp dir, then drive `sh install.sh` exactly as the sibling `internal/release/install_url_test.go` already does (`exec.Command("sh", "../../install.sh")`, env scrubbed via `scrubInstallEnv`).
- **#2 (darwin-only prose).** Narrower than the audit stated. `docs/releasing.md` is ALREADY reconciled — line 11 reads "cross-builds darwin and linux (arm64 + amd64)"; "macos-latest" (line 9) and "macOS binaries are adhoc-signed" (line 99) are accurate platform facts, not drift. The stale prose is ONLY the `release.yml` file-header comment (lines 1-4): "goreleaser cross-builds the **darwin** arm64+amd64 tarballs … Runs on macOS so the released **darwin** binaries are built natively." `.goreleaser.yaml:25-30` builds `goos: [darwin, linux]`, so the header omits linux and implies darwin-only output. Fix is the release.yml header only.
- **#3 (gofmt drift).** `gofmt -l skills/integration/survey_sync_codex_test.go` lists the file — drift is real. BUT under this host's go1.26.1, a naive `gofmt -w` does TWO things: (a) benign comment-column re-alignment on lines 62-65, and (b) it rewrites the ASCII `''` in line 23's comment (`cwd = ''.`, bytes `27 27`) into a Unicode curly `”` (U+201D, bytes `e2 80 9d`) — the go/doc comment reformatter treating `''` as a typographic quote. That is a documentation CONTENT change, not whitespace. Verified: rewording line 23 to drop the bare `''` (e.g. `cwd = empty-string`) leaves only the benign (a) alignment for `gofmt -w` to apply. Decision recorded in AC-3.
- **#4 (hasGitEntry comment).** `hasGitEntry` is `internal/status/handlers.go:557-564`; its sole unmasked guard is `TestDiscoverWorkflowsSkipsNestedCheckout` (`internal/status/discover_worktree_noise_test.go:103`). Comment-only; rides #1's commit, no standalone behavioral AC.
- **#5 (Node-20 actions).** Verified against GitHub's changelog (revised deadline **2026-06-16**) and each action's releases. The deprecation warning names exactly `actions/checkout@v4`, `actions/setup-go@v5`, `goreleaser/goreleaser-action@v6` — all three sit on the flip-cut path (`release.yml`, `install-e2e.yml`). Node-24-compatible target majors:
  - `actions/checkout@v4 → @v5` (v5.0.0 migrated to node24; v6 also exists)
  - `actions/setup-go@v5 → @v6` (v6 runs node24)
  - `goreleaser/goreleaser-action@v6 → @v7` (v7.0.0 first runs node24; v6 stays node20)
  - `actions/setup-node@v4 → @v6` (node24) — `runtime-live-e2e.yml`
  - `actions/setup-python@v5 → @v6` (node24) — `docs.yml`
  - `actions/upload-artifact@v4` — v4 still node20 by default; bump to the latest node24 major.
  - `actions/deploy-pages@v4` / `actions/upload-pages-artifact@v3` — NO node24 release published yet (upload-pages-artifact internally pins node20 upload-artifact). NOT on the flip-cut path (`docs.yml` only). Record as out-of-reach-this-task; do not invent a tag.
  - Caveat: checkout@v5 needs runner ≥ v2.327.1 — GitHub-hosted `ubuntu-latest`/`macos-latest` already satisfy this.

## Spike determination

**No spike needed** — these are localized test/doc/format/CI fixes over already-proven mechanisms. The one "is this exercisable?" unknown (can a Go test drive install.sh's checksum gate without goreleaser?) and the one "what does the tool actually do?" unknown (does `gofmt -w` change comment content?) were both resolved during ideation above by running them: the `SPACEDOCK_INSTALL_FROM=<dir>` local-dist path is already exercised by `install_url_test.go`'s sibling, and the gofmt curly-quote rewrite was confirmed by a throwaway `gofmt -d` on a temp copy. Nothing in the plan rests on an unexercised mechanism.

## Acceptance criteria

Each AC is an end-state property of the finished task, verified by something outside this body that can fail.

- **AC-1 — A Go test reds when install.sh's checksum gate is removed or weakened.**
  A new test in `internal/release` builds a local `dist/` fixture (a tar holding a bare `spacedock` file + a matching `checksums.txt`) in a temp dir and drives `sh ../../install.sh` via `SPACEDOCK_INSTALL_FROM`. It asserts the happy path installs a runnable binary AND that a byte-appended (tampered) tarball makes install.sh exit non-zero installing nothing. Plus a companion "guard is load-bearing" assertion in the style of `goreleaser_guard_test.go`'s `TestGoreleaserBuildGuardRejectsDroppedLinux`: feed install.sh's text with the checksum-gate lines (`install.sh:164-169`) stripped to a temp copy, run the tamper case against THAT copy, and assert it now wrongly exits 0 — proving the live test would have reded.
  *Verified by:* `go test ./internal/release/ -run Checksum` — the tamper assertion fails if the gate is deleted; the load-bearing assertion fails if stripping the gate does NOT change behavior (i.e. the test wasn't actually exercising the gate). Expected value (the tampered tarball's hash ≠ checksums.txt) comes from the fixture the test builds, not from any file under test.

- **AC-2 — The release.yml header describes the real darwin+linux cross-build.**
  `release.yml`'s file-header comment (lines 1-4) no longer says goreleaser builds only "darwin" tarballs or that "darwin binaries are built natively"; it names the darwin+linux × arm64+amd64 matrix `.goreleaser.yaml` actually produces, while keeping the accurate "runs on macos-latest" note. `docs/releasing.md` is left unchanged (already correct).
  *Verified by:* an extension to the existing `goreleaser_guard_test.go` (or a sibling) that parses `.goreleaser.yaml`'s build targets and asserts `release.yml`'s header does not claim a darwin-only build — i.e. if the header names a build OS set, it must not be missing `linux` while `.goreleaser.yaml` builds linux. The independent oracle is `.goreleaser.yaml`'s parsed `goos`, which can diverge from the header — that is what lets the check fail. (If a header-vs-config text check proves too brittle to bind cleanly, AC-2 downgrades to a reviewer-confirmed prose fix with `TestGoreleaserBuildsLinuxAndDarwin` standing as the existing proof that linux is built; the header is then plain doc-accuracy. Implementation decides which is clean — recorded as an open call.)

- **AC-3 — `gofmt -l` is clean on the formerly-drifting test file, with no comment content silently mangled.**
  `gofmt -l skills/integration/survey_sync_codex_test.go` prints nothing. Because go1.26.1's `gofmt -w` would rewrite the line-23 `''` into a curly `”`, the implementation FIRST rewords line 23 to drop the bare `''` (so the only formatting change is the benign column re-alignment), THEN runs `gofmt -w`. The committed file must contain no U+201D introduced by the formatter.
  *Verified by:* `gofmt -l skills/integration/survey_sync_codex_test.go` exits clean (empty output), and `grep -c $'”' skills/integration/survey_sync_codex_test.go` shows the formatter did not inject a curly quote into the cwd comment.

- **AC-4 — (rides AC-1, not standalone) a cross-reference comment at `hasGitEntry`.**
  A comment at `internal/status/handlers.go:557` names `TestDiscoverWorkflowsSkipsNestedCheckout` as the single test that guards this function, so a future edit knows what protects it. Comment-only; no behavioral proof of its own — it ships in the same change as AC-1 and is covered by the suite staying green.
  *Verified by:* the comment present at the guard site (reviewer check); `go test ./internal/status/` stays green.

- **AC-5 — Every GitHub Actions `uses:` with a published node24 major is pinned to it, enforced by a guard test.** *(time-critical, deadline 2026-06-16)*
  All workflow files pin: `checkout@v5`+ , `setup-go@v6`+ , `goreleaser-action@v7`+ , `setup-node@v6`+ , `setup-python@v6`+ , `upload-artifact` at its node24 major. `deploy-pages`/`upload-pages-artifact` are left at current pins (no node24 release exists) with a one-line comment recording why. A new guard test in `internal/release` (modeled on `goreleaser_guard_test.go`) parses every `.github/workflows/*.yml`, extracts each `uses: owner/action@vN`, and asserts the pinned major ≥ the recorded node24-minimum for that action; it reds if any action regresses below its minimum (e.g. someone re-pins `checkout@v4`). The node24-minimum map is written in the test as the independent oracle (sourced from the deprecation facts above), so the check can fail.
  *Verified by:* `go test ./internal/release/ -run Node24` — fails if any workflow pins below the recorded minimum; a companion adversarial sub-test rewrites one `uses:` line back to `@v4` and asserts the guard reds. Confirmatory (not load-bearing, since the end-of-job annotation is date-gated and known-buggy under FORCE_JAVASCRIPT_ACTIONS_TO_NODE24): the next CI run on the bumped workflows shows no "Node.js 20 actions are deprecated" annotation for the bumped actions.

## Test plan

- **AC-1:** Go unit test, `internal/release` package, ~80 lines. Cost: low — pure filesystem + `sh` exec, no network, no goreleaser. Reuses `install_url_test.go`'s script-path + env-scrub helpers. Risk: building a valid `checksums.txt` line format (`<sha256>␣␣<filename>`) and a tar with the binary at root — both directly mirrored from install.sh's awk parse (`install.sh:164`) and `tar -xzf … spacedock` (`install.sh:173`).
- **AC-2:** Either a YAML-parse extension to `goreleaser_guard_test.go` (low cost, preferred) or a reviewer-confirmed prose edit backed by the existing `TestGoreleaserBuildsLinuxAndDarwin`. Decide at implementation; no live workflow run needed.
- **AC-3:** Mechanical: reword one comment line, `gofmt -w`, commit. Verified by `gofmt -l` (exit-clean) + a curly-quote grep. Near-zero cost.
- **AC-4:** One comment line; covered by the existing `internal/status` suite staying green. Near-zero cost.
- **AC-5:** Edit the `uses:` pins across the five workflow files, add the guard test (~60 lines, mirrors `goreleaser_guard_test.go`'s parse + adversarial pair). Cost: low for the test; the bump itself is mechanical. The confirmatory no-warning CI observation is opportunistic — the guard test is the durable proof. **PRIORITIZE: land AC-5's bumps + guard before the 0.20.0 flip cut runs `release.yml`, and before 2026-06-16.**

Fixture vs CLI vs live: all five proofs are Go unit tests over real files/scripts (`gofmt`, `sh install.sh`, parsed YAML). No live workflow drive is required; the only network touch (`install_url_test.go`'s live-release check) is pre-existing and unaffected.

## Notes

Provenance: 0.19.9 pre-cut audit. `#5` (node20 actions) is time-critical (2026-06-16) and `#1` (checksum test-strength) is the substantive test-strength fix — both are the priorities; #2/#3/#4 are mechanical. Two audit-framing corrections recorded above: #2 is release.yml-header-only (docs/releasing.md already reconciled), and #3 carries a gofmt-mangles-comment-content trap that the implementation must dodge by rewording before formatting.

## Stage Report: ideation

- DONE: Firm #1 checksum-gate tamper test as a concrete fix with EXTERNAL-proof AC (reds when checksum verification removed/weakened)
  AC-1: Go test in internal/release builds a local dist/ fixture, drives `sh install.sh` via SPACEDOCK_INSTALL_FROM (no goreleaser), asserts tamper exits non-zero + load-bearing adversarial strip; mirrors install_url_test.go.
- DONE: Firm #2 darwin->darwin+linux doc reconcile AC
  Corrected scope: docs/releasing.md ALREADY says "darwin and linux" (line 11); only release.yml header (lines 1-4) is stale. AC-2 fixes the header, verified against parsed .goreleaser.yaml goos.
- DONE: Firm #3 gofmt AC (`gofmt -l` clean)
  AC-3: found go1.26.1 `gofmt -w` rewrites line-23 `''`→curly `”` (comment content). Plan rewords first, then formats; verified by `gofmt -l` clean + curly-quote grep.
- DONE: Firm #4 hasGitEntry cross-reference comment (rides #1, not standalone AC)
  AC-4: comment at handlers.go:557 naming TestDiscoverWorkflowsSkipsNestedCheckout; covered by internal/status suite, no standalone behavioral proof.
- DONE: Firm #5 node-action bump (PRIORITIZED, deadline 2026-06-16)
  AC-5: verified the three warned actions (checkout@v4, setup-go@v5, goreleaser-action@v6) + targets checkout@v5/setup-go@v6/goreleaser-action@v7/setup-node@v6/setup-python@v6/upload-artifact node24 major. Guard test (modeled on goreleaser_guard_test.go) asserts pinned major >= recorded node24 minimum; deploy-pages/upload-pages-artifact have no node24 release yet (out of reach, not on flip-cut path). Flagged time-critical.
- DONE: Record spike determination
  "No spike needed" — localized fixes; the two genuine unknowns (Go can drive install.sh's gate sans goreleaser; gofmt mangles comment content) were resolved by running them during ideation.

### Summary

Fleshed all five items into entity-level ACs each with an external, fail-able proof — four are Go guard tests (checksum tamper, header-vs-goreleaser, node24-minimum pins) plus mechanical gofmt/comment fixes. Two audit-framing corrections recorded: #2 is release.yml-header-only (docs/releasing.md already reconciled), and #3 carries a go1.26.1 gofmt trap that rewrites a `''` comment into a curly quote unless reworded first. #5 (node20 actions, deadline 2026-06-16) and #1 (checksum test-strength) are the priorities; verified exact node24-compatible target majors against GitHub's changelog and each action's releases.
