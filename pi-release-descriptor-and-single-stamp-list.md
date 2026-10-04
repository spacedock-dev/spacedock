---
title: Give the Pi extension its own release descriptor, and make the stamp list exist once
status: implementation
score: 0.8
source: "Captain constraint, 2026-10-04: the Pi extension's canonical metadata must not be the top-level package.json, because every other host has its own. Astra design review, run b3928dbb."
id: 8vahpzjd358ygqce5etw6fr2
started: 2026-10-04T05:28:21Z
worktree: .worktrees/spacedock-ensign-pi-release-descriptor-and-single-stamp-list
---

## Problem

The Pi extension has no release descriptor of its own. Claude and Codex each carry
one, `.claude-plugin/plugin.json` and `.codex-plugin/plugin.json`, and each carries
its `version`. Pi has neither, so the only candidate for its version is the root
`package.json` — the same file every other host would want, and the file Pi reads
for resource routing. Claiming it for one host is wrong, and leaving Pi unversioned
means the release cannot stamp it at all.

Two release defects sit in the same family. `StampVersion` returns a named manifest
unchanged when it has no top-level `version`, and the release CLI then writes those
unchanged bytes and prints success. And `replaceFirstVersion` rewrites the first
textual `version` occurrence rather than the top-level key, so a nested version can
capture the write.

The stamp list is written out in three places in `docs/releasing.md`, and repeated
again in the release workflow, so adding a fourth manifest would add a fifth
authority.

## Design

**1. New file `.pi/plugin.json`**, beside the extension it describes. Spacedock-owned
release metadata, not something Pi natively discovers:

    {
      "name": "spacedock",
      "version": "0.28.0-pre3",
      "description": "Turn directories of markdown files into structured workflows operated by AI agents",
      "author": { "name": "CL Kao" },
      "homepage": "https://github.com/spacedock-dev/spacedock",
      "repository": "https://github.com/spacedock-dev/spacedock",
      "license": "Apache-2.0",
      "keywords": ["workflow", "pipeline", "agents", "markdown", "automation"],
      "interface": {
        "displayName": "Spacedock",
        "shortDescription": "Plain text workflows for AI agents",
        "longDescription": "Turn directories of markdown files into structured workflows operated by AI agents.",
        "developerName": "CL Kao",
        "category": "workflow",
        "capabilities": ["Interactive", "Write"],
        "websiteURL": "https://github.com/spacedock-dev/spacedock"
      }
    }

It carries identity, description and version. It carries NO resource routing: Pi
continues to read `pi.extensions` and `pi.skills` from the root `package.json`, which
stays unchanged and unversioned.

**2. One authoritative target list** in `internal/release`:

    func StampTargets() []string {
        return []string{
            ".claude-plugin/plugin.json",
            ".codex-plugin/plugin.json",
            ".pi/plugin.json",
            "skills/first-officer/references/first-officer-shared-core.md",
        }
    }

- `stamp-version <version>` and `manifest-tag-gate <tag>` use it when given no
  explicit targets.
- A new `stamp-paths` subcommand prints it, so the release workflow's gate, stamp,
  diff and commit steps stop repeating it, and the `FO_PROSE` variable goes away.
- `docs/releasing.md` keeps one invocation per ritual step and drops its target lists.

**3. Two correctness fixes in the release CLI.** A named JSON manifest with no usable
top-level `version` must fail with a clear error and a non-zero exit instead of
writing unchanged bytes and reporting success; the low-level marketplace no-op
contract keeps its behaviour, because it is explicitly tested. And the top-level
`version` must be the field replaced, even when a nested `version` appears earlier.

**4. Root `package.json` stays as it is.** Its narrower responsibility is documented
where the release ritual is described, not encoded in the file.

## Out of scope

- Moving the Pi package into a subdirectory, and any `git:.../subdirectory` syntax.
  Pi's git source model has no subdirectory selector.
- Splitting the extension into its own repository. The extension is coupled to the
  shared skills tree, so that is a distribution migration, not a metadata move.
- Making Pi display or enforce the descriptor's version.
- Prerelease coherence, where CI skips the manifest and tag gate for prereleases.

## Expected surface and tolerance

New `.pi/plugin.json`; `internal/release/release.go` plus its test;
`cmd/spacedock-release/main.go`; `docs/releasing.md`; `.github/workflows/release.yml`.
Estimate net +120, across 6 files. Tolerance net +50 to +200, at most 8 files.

## Acceptance criteria

**AC-1 (VALUE) - One invocation stamps every host, including Pi.**
Verified by: `stamp-version X.Y.Z` with no target arguments rewrites the version in
`.claude-plugin/plugin.json`, `.codex-plugin/plugin.json`, `.pi/plugin.json` and the
shared-core prose. Falsifier: omit the Pi descriptor from the list and the stamp
leaves it at its old version.

**AC-2 - The list has one authority.**
Verified by: `stamp-paths` prints exactly the list the stamping and the tag gate use,
and the test fails when the two diverge. Falsifier: give the gate a different list.

**AC-3 - A named JSON manifest with no version fails loudly.**
Verified by: naming a JSON manifest without a usable top-level version exits
non-zero with a clear error and prints no success line. Falsifier: remove the check
and the command exits zero on unchanged bytes.

**AC-4 - The top-level field is the one replaced.**
Verified by: a manifest carrying a nested `version` before the top-level one; the
top-level value changes and the nested one does not. Falsifier: keep
first-occurrence replacement and the nested value changes instead.

**AC-5 (no-regression) -** The marketplace no-op contract, `dev-preversion`,
`manifest-tag-gate`, and the existing release tests keep their behaviour; the Pi
extension still loads through the real installed loader with the root manifest
untouched; the offline suite passes.

## Test plan

Behavioural tests only, with the version supplied independently. No test may assert a
file's own text against a copy of it, and no expected value may come from the thing
it tests. Pi discovery is checked through the real installed loader, not a copied
manifest string.

## Stage Report: implementation

- DONE: Create .pi/plugin.json with the identity, description and version shown in the design, and no resource routing. Leave the root package.json unchanged and unversioned.
  New `.pi/plugin.json` (30 lines) carries name/version/description/author/homepage/repository/license/keywords/interface only. Commit `d9842ab5a` touches no `package.json` (`git show --stat` grep for it is empty).
- DONE: Put the authoritative stamp target list in internal/release, including the new Pi descriptor, and have stamp-version and manifest-tag-gate use it by default. Add a stamp-paths subcommand that prints it, then remove the repeated target lists from the release workflow's gate, stamp, diff and commit steps and drop its FO_PROSE variable, and from the three invocations in docs/releasing.md.
  `release.StampTargets()` is the sole list (`.pi/plugin.json` included); `stampVersion`/`runManifestTagGate` default to it; `stampPaths` prints it. `grep -n 'plugin.json\|FO_PROSE' .github/workflows/release.yml docs/releasing.md` returns nothing.
- DONE: Make a named JSON manifest with no usable top-level version fail with a clear error and a non-zero exit instead of writing unchanged bytes and printing success. Keep the low-level marketplace no-op contract.
  Manual: `stamp-version 9.9.9 marketplace.json` -> exit 1, `stamp marketplace.json: no usable top-level "version" field...`, bytes unchanged. `StampVersion` still returns marketplace bytes unchanged (TestStampVersionLeavesMarketplaceCalendarUntouched, unchanged).
- DONE: Replace the top-level version field, not the first textual occurrence, and cover a nested version appearing earlier in the file.
  `replaceTopLevelVersion` walks top-level keys with json.Decoder. Manual nested-first run left `"metadata":{"version":"schema-7"}` and wrote only top-level `3.4.5`; covered by TestStampVersionReplacesTopLevelVersionWhenNestedVersionPrecedes.
- DONE: Add behavioural tests for all four: the default list stamps the Pi descriptor, stamp-paths and the gate share one list, a version-less named manifest fails, and a nested version is not the one written. Verify Pi discovery through the real installed loader with the root manifest untouched, and say what that check observes.
  Tests: TestStampVersionCommandDefaultsToStampTargets (fails if `.pi/plugin.json` missing from the default list -> stays 0.0.0); TestStampPathsMatchesDefaultStampList + TestManifestTagGateDefaultsToStampPathsList (fail if stamp-paths and the default stamp/gate lists diverge); TestStampVersionCommandFailsOnVersionlessNamedManifest (fails if exit 0 or a `stamped` success line on unchanged bytes); TestStampVersionReplacesTopLevelVersionWhenNestedVersionPrecedes (fails if the nested value changes). All in `go test ./internal/release/ ./cmd/spacedock-release/`.
- DONE: Pi discovery through the real installed loader (pi SDK `DefaultResourceLoader`, v1.0.2) with the worktree as the registered package, root manifest untouched.
  Observed: `ROOT_MANIFEST.pi = {"extensions":["./.pi/extensions/spacedock.ts"],"skills":["./skills"]}`; `ROOT_MANIFEST_UNCHANGED = true`; `DISCOVERED_EXTENSIONS = [".pi/extensions/spacedock.ts"]`; `EXTENSION_ERRORS = 0`; repo skills include `first-officer` and `ensign`. `.pi/plugin.json` was present during the load and caused no resource or error: Pi reads `pi.extensions`/`pi.skills` from the root manifest, not the descriptor.

### Summary

Gave Pi its own release descriptor and made the stamp target list exist once in `internal/release`, consumed by `stamp-version`, `manifest-tag-gate`, and the new `stamp-paths` subcommand; removed the repeated lists from `release.yml` (incl. `FO_PROSE`) and `docs/releasing.md`. Fixed the two release defects: a version-less named JSON manifest now exits non-zero with a clear error and no success line (low-level marketplace no-op kept), and the top-level `version` field is replaced even when a nested `version` precedes it. Verified Pi discovery through the real installed pi loader with the root manifest untouched.

Deviation: diff is net +353 across 9 files vs the design's estimate +120 / tolerance +50..+200 / at most 8 files. The extra surface is forced: `manifest_tag_gate.go` must change for the gate default, and the four mandated CLI behavioural tests land in the pre-existing `cmd/spacedock-release` test files. Pre-existing gofmt deviation in `internal/release/runtime_live_evidence_workflow_test.go` (unformatted at HEAD) was left untouched to keep scope narrow.



## Review-finding disposition

### V1 — case-insensitive read bypasses exact-key stamp guard

- Observation: `internal/release/release.go:120-127` decodes into a struct (case-insensitive field names); `:137-144` checks an exact map key. `cmd/spacedock-release/main.go:198-217` consequently writes unchanged bytes and announces a version it never installed.
- Released user and normal workflow: release cutter explicitly names a JSON descriptor; rejecting descriptors without a usable top-level `version` is a promised input-validation boundary, not limited to the four default files.
- Observable harm: for `uppercase.json` containing `{"name":"fixture","Version":"1.0.0"}\n`, `stamp-version 7.8.9 uppercase.json` exits 0; stdout exactly `stamped uppercase.json version=7.8.9\n`; stderr empty; bytes identical, mtime changed.
- Affected authority: value-ac[AC-3] A named JSON manifest without the exact top-level version must fail rather than silently claim a successful stamp.
- Trigger evidence: independently constructed fixture executed against the built candidate CLI; `Version` is not JSON's case-sensitive `version`. Ordinary lowercase missing/empty/null cases correctly fail, so this is specifically reader/writer identity disagreement.
- Advisory classification: **outcome defect / material**. Ownership: this task's named-manifest guard. Recommend a narrow fix aligning exact-key validation with stamping and an independent regression fixture; no new controller or design reset is needed for this bug. **REJECTED; candidate unchanged; FO authorization required before any repair.**

### V2 — duplicate top-level keys disagree across writer and reader

- Released user and workflow: explicitly naming a descriptor with duplicate `version` members; this non-interoperable JSON shape is outside the promised ordinary unique-field host descriptors.
- Observable harm and trigger: `{"version":"7.8.9","version":"1.0.0"}\n` stamped to `7.8.9` exits 0, stdout exactly `stamped duplicate.json version=7.8.9\n`, stderr empty, bytes unchanged and mtime changed; subsequent `manifest-tag-gate v7.8.9 duplicate.json` exits 1, stderr exactly `spacedock-release manifest-tag-gate: tag v7.8.9 does not match tagged commit's plugin.json version 1.0.0; stamp the manifest to 7.8.9 and tag THAT commit (duplicate.json)\n`, stdout empty. Writer chooses first; canonical Go/Python readers choose last.
- Affected authority: value-ac[AC-4] The effective top-level version must be the one replaced; duplicate-key semantics are not specified for this release.
- Advisory classification: **outcome defect / deferred risk**. Ownership: release JSON validation; recommend FO disposition, not automatic scope expansion. Unique-key defaults and nested-first fixtures pass AC-1/AC-4. Promote to material if duplicate-key manifests are accepted as supported input or appear in a released host descriptor.

### V3 — root package responsibility is not documented as designed

- Released user and workflow: maintainer reading `docs/releasing.md`; Design 4 promised an explanation that root `package.json` stays unversioned and owns Pi resource routing. The document has no such explanation.
- Observable harm: missing maintenance guidance only; root manifest actually remains unchanged/unversioned and real Pi discovery passes. Affected authority: none: no current value AC or runtime boundary fails from this omission.
- Advisory classification: **outcome defect / deferred risk (documentation polish)**. Ownership: this task's release documentation. Recommend documenting the distinction when authorized; promote if a release process begins stamping/routing through the wrong manifest.

## Stage Report: validation

- DONE: Run the release CLI yourself and report exact output: stamp-paths; stamp-version with no target arguments on a scratch copy; manifest-tag-gate with no target arguments; and a named JSON manifest with no usable top-level version, which must exit non-zero with a clear error, leave the bytes unchanged, and print no success line.
  Built candidate `d9842ab5a` using `go build -o .validation-release/spacedock-release ./cmd/spacedock-release`; ran the binary from `.validation-release/copy` with the actual four target files copied there. Outputs below use JSON string notation, preserving exact newlines; stderr empty unless stated.
  `stamp-paths`: exit 0, stdout `.claude-plugin/plugin.json\n.codex-plugin/plugin.json\n.pi/plugin.json\nskills/first-officer/references/first-officer-shared-core.md\n`.
  `stamp-version 7.8.9`: exit 0, stdout `stamped .claude-plugin/plugin.json version=7.8.9\nstamped .codex-plugin/plugin.json version=7.8.9\nstamped .pi/plugin.json version=7.8.9\nstamped skills/first-officer/references/first-officer-shared-core.md version=7.8.9\n`.
  `manifest-tag-gate v7.8.9`: exit 0, stdout `tag v7.8.9 matches tagged commit's plugin.json version 7.8.9 (.claude-plugin/plugin.json)\ntag v7.8.9 matches tagged commit's plugin.json version 7.8.9 (.codex-plugin/plugin.json)\ntag v7.8.9 matches tagged commit's plugin.json version 7.8.9 (.pi/plugin.json)\ntag v7.8.9 matches prose-stamped minor 7.8 (skills/first-officer/references/first-officer-shared-core.md)\n`.
  `stamp-version 7.8.9 missing.json`, input `{"name":"fixture","metadata":{"version":"nested"}}\n`: exit 1, stdout empty, stderr `stamp missing.json: no usable top-level "version" field to replace (marketplace/no-op shape is not a stamp target)\n`; bytes and mtime unchanged. V1 falsifies the universal claim despite this ordinary case passing.
- DONE: Confirm the single-authority claim by inspection: no stamp target list may remain in docs/releasing.md or in .github/workflows/release.yml, the FO_PROSE variable must be gone, and every remaining site must consume stamp-paths or the default list. Name each site you checked.
  Checked `internal/release/release.go:StampTargets`, `main.go:stampVersion/stampPaths`, `manifest_tag_gate.go:runManifestTagGate`; docs stable ritual stamp/commit/gate, prerelease stamp/commit, patch stamp and waiver reference; workflow e2e-gate and publish stamp/diff/commit. All active sites consume defaults or `STAMP_PATHS`; no target list or `FO_PROSE` remains in those docs/workflow. Historical roadmap reports are not active authorities; test fixtures independently enumerate expected hosts (not production routing).
- DONE: Independently construct a manifest with a nested version before the top-level one and require that only the top-level field changes.
  Input `{ "meta": {"version":"schema-Ω"}, "version" : "1.0.0", "tail":"stay" }\n`; output exactly `{ "meta": {"version":"schema-Ω"}, "version" : "7.8.9", "tail":"stay" }\n`; CLI exit 0, stdout `stamped nested.json version=7.8.9\n`, stderr empty. Exact-byte assertion, not only a substring; escaped `ver\u0073ion` also succeeds.
- DONE: Verify the Pi discovery claim through the real installed loader, and confirm that .pi/plugin.json has no effect on discovery. State what the check observed, and treat that boundary as expected, not a defect.
  Used installed `@earendil-works/pi-coding-agent` 1.0.2 at `/Users/clkao/.local/share/fnm/node-versions/v24.13.1/installation/lib/node_modules/@earendil-works/pi-coding-agent/dist/index.js`, importing `DefaultResourceLoader`/`SettingsManager`; registered an actual package copy with `SettingsManager.inMemory({packages:[pkg]},{projectTrusted:true})`, isolated cwd/agentDir, `reload()`. Compared descriptor present, absent, and invalid; no candidate/global settings changed.
  Each run's observed resource set was exactly `{"extensions":[".pi/extensions/spacedock.ts"],"errors":[],"skills":["commission","debrief","ensign","feedback-rejection-flow","first-officer","fo-dispatch-recovery","fo-gate-lifecycle","fo-status-viewer","present-gate","refit","survey"]}`. Final output `ROOT_MANIFEST_UNCHANGED=true; DISCOVERY_UNAFFECTED=true; LOADER_VERSION=1.0.2`. Root bytes equal original, candidate `package.json` diff empty, descriptor has identity/description/version and no routing. Ignoring this Spacedock-owned descriptor is the expected Pi boundary.
- DONE: Audit for tautological tests: no test may assert a file's own text against a copy of it, and no expected value may come from the thing it tests. Reject any test that does, and name it.
  Inspected all changed tests in `release_test.go`, `stamp_version_test.go`, `manifest_tag_gate_test.go`: none is tautological. Expected versions are independent literals; target fixtures pin host membership independently; relational stamp-paths tests exercise distinct writer/gate operations, not text copied from their implementation. No-op byte comparisons use authored input fixtures. Omitting Pi from stamping fails `TestStampVersionCommandDefaultsToStampTargets`; gate omission fails `TestManifestTagGateDefaultsToStampPathsList`; old first-text replacement fails `TestStampVersionReplacesTopLevelVersionWhenNestedVersionPrecedes`. Existing missing-version test simply lacks V1's case variant, not an invalid oracle.
- DONE: Look for any remaining path that writes unchanged bytes and reports success.
  V1 and V2 reproduced with changed mtimes. Already-correct `{"version":"7.8.9"}\n` also writes unchanged bytes and prints `stamped already-stamped.json version=7.8.9\n` (exit 0, stderr empty), but the claimed version is true: expected idempotent post-tag behavior, not the false-success defect. Prose at the same major.minor is similarly intentionally idempotent.
- DONE: Run applicable tests from the Testing Resources section and report results.
  `go test ./internal/release ./cmd/spacedock-release -count=1`: both pass (11.897s / 0.541s); `go test ./internal/release ./cmd/spacedock-release -race -count=1`: both pass (14.257s / 1.910s); `git diff --check main...HEAD`: pass. Includes unchanged marketplace no-op, dev-preversion, manifest gate, stamp/prose and workflow guards. `dev-preversion 7.8.9`: exit 0, stdout exactly `7.9.0-pre1\n`, stderr empty.
- DONE: Semantic adversarial pass and acceptance-criterion evidence.
  **AC-1 (VALUE)** passes: independent JSON decode finds all three host versions `7.8.9`, real prose pin `7.8`. **AC-2** passes: independently diverged each of the four expected files to `2.3.4` / `2.3`; each default gate exits 1 naming that file. **AC-3** FAILS V1, while absent/empty/whitespace/null/non-string/array/top-level-null/truncated EOF all fail without success or mutation. **AC-4** passes ordinary unique-key nested-first plus Unicode/escaped-key fixtures; V2 deferred. **AC-5 (no-regression)** focused regressions and real loader pass; full offline-suite evidence is intentionally not claimed.
  Linear parser scans/read-all allocations inspected; 4,000,072-byte UTF-8 nested fixture stamped in 0.061s with exact bytes preserved except top-level value (exit 0, stdout `stamped large.json version=7.8.9\n`, stderr empty). No new implicit size limit observed. Preliminary harness assertion assumed a space after the real prose minor; corrected the assertion to the independently observed sentence punctuation and reran. One inspection guessed nonexistent `prose_version.go`; actual `prose_stamp.go` inspected next.
- SKIPPED: Repository-wide suite, CI lane run, and candidate formatting/repair.
  Explicit FO verification-only scope; only focused packages and isolated fixtures exercised. No product/test files added or modified; scratch resources removed after verification.
- DONE: Report each finding with the defect kind and release scope, and whether the change can proceed.
  Recommendation **REJECTED: do not proceed** for material outcome defect V1; V2/V3 are separately deferred. See advisory dispositions above; no fixes authorized or attempted. Implementation also exceeds declared tolerance (9 files, +427/-74, net +353 vs max 8 files / +200); retain captain-visible scope/design-reset acknowledgement before another cycle rather than silently accepting the implementation's explanation.

### Summary

Verified the default stamp/gate/list behavior and real Pi discovery, but falsified the promised named-manifest fail-loud boundary: a case-variant key still produces the original unchanged-byte false success. The intended value remains reachable by a narrow exact-key guard correction, not a new mechanism; validation recommends REJECTED pending FO disposition, with duplicate-key handling and missing root-manifest documentation recorded separately. Only this state report was changed; candidate `d9842ab5a` remains untouched.

## Stage Report: implementation (cycle 2)

- DONE: Fix the false success the validation reproduced.
  Commit cc20c4c59. `release.VerifyStampedVersion` re-reads each target after the write and the CLI prints success only once it passes. Live run on `{"name":"fixture","Version":"1.0.0"}`: `stamp-version 7.8.9 uppercase.json` -> exit 1, stdout empty, stderr `stamp …/uppercase.json: wrote bytes but the target did not take version 7.8.9: target still reports version "1.0.0", want "7.8.9"`; shasum identical before/after.
- DONE: Do not patch by widening the key check; verify the command's own result, for every named target not only JSON manifests.
  The exact-lowercase-key writer and the pre-check are untouched; the new guard is a post-write read-back. Dispatch by extension: JSON via `ManifestVersion`, prose via `ProseMinor` (both readers the gate already uses), so a prose target whose minor did not change also fails.
- DONE: Keep the low-level marketplace no-op contract.
  `TestStampVersionLeavesMarketplaceCalendarUntouched` unchanged and passing; `StampVersion` still returns marketplace bytes unchanged.
- DONE: Add the negative case for the exact reproduction, and state why each new test could fail.
  `TestStampVersionCommandFailsOnCaseVariantVersionKey` (cmd): asserts exit != 0, no success line, bytes == authored `src`. Falsifier demonstrated: stashing only `main.go` makes it FAIL with "exit = 0".
  `TestVerifyStampedVersionRejectsUnchangedProse` (release): prose at minor 1.2 for requested 9.9.0 must error; fails if the guard returns nil without reading `ProseMinor`.
  `TestVerifyStampedVersionAcceptsStampedTargets` (release): stamped manifest + prose must pass; fails if the prose expectation is the full "9.9.0" instead of major.minor.
- DONE: Record the surface overrun for the captain; do not reduce it.
  Branch vs main: 9 files, +559/-74, net +485 (design tolerance +50..+200, at most 8 files). This cycle alone: 4 files, +132/-0, no new files.

### Summary

Replaced the defeated pre-check with a read-back verification of the command's own result: after writing each named target, `stamp-version` re-reads it and requires the requested version, exiting non-zero with an error naming the target and printing no success line otherwise. The exact case-variant reproduction now fails loudly over unchanged bytes, `git stash` of the guard makes the new test fail, and the marketplace no-op contract is unchanged. No expected value in the new tests is copied from the file under test.

## Stage Report: validation (cycle 2)

- DONE: Reproduce the exact case the first round found: a named manifest whose top-level version key is capitalized. It must exit non-zero, print nothing to stdout, name the failing target in the error, and leave the bytes unchanged. Report the exact output.
  Built candidate `cc20c4c59`; independently authored `uppercase.json` as `{"name":"fixture","Version":"1.0.0"}\n`, then ran `stamp-version 7.8.9 uppercase.json` in a worktree-local scratch directory. Exit **1**, stdout exactly `""`; stderr exactly `stamp uppercase.json: wrote bytes but the target did not take version 7.8.9: target still reports version "1.0.0", want "7.8.9"\n`. Here `\n` denotes the actual terminal newline.
  Bytes exactly equal the authored input; SHA-256 `d872f659d7736d619ea594036cd6efb746e619f424164722a715754c31af8103`. **V1: outcome defect / material — repaired**, AC-3's false-success reproduction no longer occurs. This is a post-write check, not an atomic/no-write guarantee.
- DONE: Independently confirm the falsifiability claim rather than trusting it: revert only the command-line change and require the new regression test to fail. Then restore it.
  Replaced only `cmd/spacedock-release/main.go` with `HEAD^` bytes; `git diff --name-only` confirmed it was the sole changed file. `go test ./cmd/spacedock-release -run '^TestStampVersionCommandFailsOnCaseVariantVersionKey$' -count=1` exited **1**, failing at `stamp_version_test.go:73`: `stamp-version exit = 0 for a manifest whose top-level version key is not lowercase `version`; want non-zero`. The new test and verifier implementation stayed in place.
  Restored original main.go in a `finally` block (SHA-256 `70c3446579db1f0f7ab53f16accd2279b1d7b1a14bbd48b46b5ede4f9f9f83a6`); reran that exact regression with `-v`: PASS, exit 0. Candidate tracked diff empty. **Evidence defect / material: none**; the claimed falsifier was independently observed.
- DONE: Confirm the low-level marketplace no-op contract still holds and that its test is unchanged.
  `TestStampVersionLeavesMarketplaceCalendarUntouched` passes; compared its complete function against `main`: identical. Its oracle is authored marketplace input, not copied repository text; changing the nested calendar version falsifies it. No regression finding (AC-5).
- DONE: Re-confirm briefly that the earlier passing checks still hold: stamp-paths prints the list; stamp-version and manifest-tag-gate use it with no arguments; no target list remains in docs/releasing.md or .github/workflows/release.yml; FO_PROSE is gone; the top-level version is the field replaced when a nested version appears earlier; no test is tautological, meaning no test asserts a file's own text against a copy and no expected value comes from the thing it tests.
  `stamp-paths` exit 0, stderr empty, stdout exactly `.claude-plugin/plugin.json\n.codex-plugin/plugin.json\n.pi/plugin.json\nskills/first-officer/references/first-officer-shared-core.md\n`. Independent four-path expectation, not derived from this output.
  Copied actual targets into isolated scratch; `stamp-version 7.8.9` and `manifest-tag-gate v7.8.9` with **no target arguments** both exit 0. Independently decoded all three JSON versions as `7.8.9` and checked prose minor `7.8`; separately diverging each of the four files makes the default gate exit 1 naming that file. AC-1/AC-2 pass; omitting Pi from defaults or diverging the gate list falsifies this proof.
  Inspected `StampTargets`, CLI `stampVersion`/`stampPaths`, `runManifestTagGate`; docs stable/prerelease/patch rituals and waiver reference; workflow gate and publish stamp/diff/commit. All use defaults or `STAMP_PATHS`; neither doc/workflow retains literal targets or `FO_PROSE`.
  Nested-first input `{ "meta": {"version":"schema-7"}, "version" : "1.0.0", "tail":"stay" }\n` becomes exactly `{ "meta": {"version":"schema-7"}, "version" : "7.8.9", "tail":"stay" }\n`; exit 0. AC-4 passes; first-text replacement would fail the exact-byte assertion and `TestStampVersionReplacesTopLevelVersionWhenNestedVersionPrecedes`.
  Audited all changed tests in `release_test.go`, `stamp_version_test.go`, `manifest_tag_gate_test.go`: no tautological assertions. Expected versions and host membership are independently authored. Relational stamp-paths tests compare separate writer/gate behavior against literal versions; they do not use a file's own bytes as its correctness oracle. Prose verifier mismatch/acceptance tests together reject both unconditional acceptance and full-semver/minor confusion.
- DONE: Run applicable tests and a focused adversarial pass.
  `go test ./internal/release ./cmd/spacedock-release -count=1` PASS (21.040s / 1.289s); same packages with `-race -count=1` PASS (22.419s / 2.118s); explicit marketplace/nested/read-back tests PASS. `dev-preversion 7.8.9` exit 0, stdout `7.9.0-pre1\n`, stderr empty. AC-5 focused regressions pass; root `package.json` diff remains empty.
  Missing/empty/null top-level versions each exit 1, stdout empty, authored bytes unchanged. Previous duplicate-key V2 fixture now exits 1 with empty stdout and target-naming read-back error instead of false success. **V2: outcome defect / deferred risk — false-success instance repaired; duplicate-key writer/reader ambiguity remains outside unique-key host descriptors**; promote only if duplicate-key manifests become supported or shipped. Ordinary/nested unique-key paths pass.
- SKIPPED: Complete repository-wide suite evidence and live Pi rerun.
  Attempted `go test ./...`: exceeded 240s; `go test ./... -race`: exceeded 180s; both inconclusive, not passing evidence or an observed assertion failure. Focused suites completed. No live Pi run started, as explicitly instructed; cycle-1 installed-loader evidence is retained rather than re-claimed. AC-5 full offline-suite evidence remains incomplete.
  Ran `gofmt -w ./cmd ./internal`; only pre-existing `internal/release/runtime_live_evidence_workflow_test.go` differed, restored its original bytes to preserve verification-only scope. `git diff --check main...HEAD` passes; all scratch resources removed; candidate HEAD unchanged and index empty.
- DONE: Judge the surface overrun and record it for the captain, naming the estimate and the actual.
  Independently measured `git diff --numstat main...HEAD`: **9 files, +559/-74, net +485**, versus estimate **6 files / net +120**, tolerance **at most 8 files / net +50..+200**. Actual is 150% of estimated file count and 404% of estimated net lines; exceeds ceilings by 1 file and 285 net lines. Scope variance, not a newly demonstrated runtime or evidence defect; release disposition **needs captain decision**, not silently accepted. Most added surface is parser/read-back behavior and behavioral tests; that rationale does not waive the declared tolerance.
- DONE: Report each finding with the defect kind and release scope, and whether the change can proceed.
  **V3: outcome defect / deferred documentation polish** remains: root package's unversioned routing responsibility is not explained in releasing.md; no current runtime/value loss (unchanged root and prior loader evidence). Revisit if release instructions cause stamping/routing through the wrong file. No new material product or tautological-evidence finding in this cycle; V1 is closed, V2 mitigated, V3 unchanged.
  Recommendation **PASSED for the requested repair and focused release checks**; can proceed to captain review, but final delivery still needs explicit acknowledgement of the above-tolerance surface and incomplete whole-repository suite evidence. No candidate repair, code commit, new test, or scope waiver was made by validation.

### Summary

The exact rejected uppercase-key case now fails loudly over unchanged bytes, and removing only the CLI repair independently makes the new regression fail; restoration makes it pass. Default stamping/gating, nested-field identity, marketplace compatibility and non-tautological test oracles remain sound under focused tests and race checks. The repaired behavior can proceed to captain review, with the 9-file/net +485 overrun and inconclusive full suites explicitly visible rather than waived.
