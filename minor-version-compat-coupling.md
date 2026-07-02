---
title: Replace the contract-integer compatibility gate with minor-version coupling (self-maintaining skill↔binary compat)
status: ideation
score: 0.6
source: "captain proposal during the 0.23.0 cut (2026-06-22). The contract integer stayed dead at 1 across 0.20/0.22/0.23 and was never bumped when the interface broke (the v0.23.0 skill-skew bug). contract-version-bump-v2 (7h3) is the tactical integer bump this would replace."
id: kr7s9efxas4fqhtcj94fs4br
started: 2026-07-02T09:55:07Z
---

Replace the manual contract-integer gate (binary reports `contract <N>`; plugin declares `requires-contract: >=X,<Y`) with minor-version coupling: 0.X.* skills require a >=0.X.0 binary (same-minor compatible; patches interchangeable). The integer is a hand-maintained value that got forgotten across three minors — minor-coupling moves the floor automatically each release, so the skill-skew class (old binary + new skills booting clean, then breaking on missing verbs) cannot recur across minors (within-minor prerelease skew is an accepted residual — see D1).

## Problem

The contract integer decouples interface-version from release-version IN THEORY, but in PRACTICE it stayed `1` across 0.20/0.22/0.23 even when the skill↔binary interface broke (v0.23.0 added hard deps on `state commit/ready/sweep` + `merge guard`, absent in contract-1 binaries). The gate gave a false-green and the FO broke cryptically. `contract-version-bump-v2` (7h3) ships the integer bump to 2 for v0.23.0 as a tactical fix; this task replaces the mechanism so the bump is never forgotten again.

## Design (decided at ideation, 2026-07-02)

The release train already couples both sides: goreleaser stamps `internal/cli.Version` from the git tag, and `spacedock-release stamp-version` stamps both vendored `plugin.json` manifests at THREE call sites — the manual pre-tag stamp (releasing.md step 3), the post-tag reconcile of `main` to the stable `X.Y.Z` (release.yml:228), and the `next` edge-branch bump to the dev pre-version `X.(Y+1).0-pre1` (release.yml:306) — with `manifest-tag-gate` (release.yml:138) blocking a stable tag whose manifests disagree. Prerelease tags are deliberately carved out of both the gate and the main stamp (a `vX.Y.Z-pre.N` tag legitimately differs from the manifests). The gate therefore needs NO new declared field: the plugin manifest's existing `version` IS the compatibility declaration, and the floor moves on every release with zero hand edits. D5's prose stamp must ride ALL THREE stamp call sites (plus the tag gate's argument list), or a branch's skills would contradict their own manifest — after a stable cut, `next` carries manifests at `X.(Y+1).0-pre1` while un-restamped prose would still say `0.X`, turning the sync test red on that branch and mis-gating dev-channel sessions.

### D1 — Gate semantics: minor-exact, both directions

The binary and the plugin must share `(major, minor)`. `binary < plugin` on that pair → `too-old-binary` (upgrade the binary); `binary > plugin` → `too-old-plugin` (refresh the plugin); equal → `compatible`. Patch skew (0.24.0 ↔ 0.24.3) and prerelease skew (0.24.0-pre1 ↔ 0.24.0) are interchangeable. Rationale: this reproduces the half-open `>=N,<N+1` bracket with the minor as the auto-advancing integer, and keeps BOTH rejection directions the integer scheme had — a version-floor-only rule would silently pass old skills against a newer binary, the mirror image of the v0.23.0 skew bug (output-shape changes break old skills just as missing verbs break old binaries). A newer binary is NOT assumed forward-compatible; its remedy is the free `spacedock install`, auto-run by the front door (D6).

**Accepted residual: within-minor prerelease skew stays ungated.** Verbs added between prerelease snapshots of the SAME minor — between `0.23.0-pre.2` and `0.23.0-pre.4`, or between a preN and the final cut, the very window where the v0.23.0 verbs actually landed — are not caught: same-minor artifacts are presumed compatible. This is a decision, not an oversight. The strict alternative (exact full-version match whenever either side carries a prerelease suffix) is broken by the actual release mechanics: prerelease tags deliberately stamp NO manifests (the release.yml carve-outs above), so a `v0.X.0-pre.N` binary reports `0.X.0-pre.N` while its matching next-channel skills say `0.X.0-pre1` — strict matching would false-red the entire next channel unless prerelease cuts grew their own stamping machinery. The exposure is confined to next/edge-channel users, who opt into instability and whose remedy is refreshing both sides, whereas the v0.23.0 incident's victims were STABLE users crossing a minor boundary — exactly the class minor-coupling kills. The same granularity limit covers mid-cycle dev skew (a day-1 dev binary vs day-30 skills, both claiming `0.X.0-pre1`): ungated, known, accepted.

### D2 — Comparison inputs

Binary side: `cli.Version` in-process (the binary never parses its own `--version` text). Plugin side: the manifest `version` field `readManifest` already returns. `requires-contract` is no longer read by the new binary (a manifest that still carries it is simply ignored — an installed 0.23 plugin classifies as `too-old-plugin` by its version, the correct verdict). Prerelease/build suffixes are cut at the first `-` or `+` before extracting major.minor; both published suffix styles (`0.24.0-pre1`, `0.23.0-pre.4`) parse identically (spiked below). A manifest with a missing or unparseable `version` is `malformed-version` (packaging bug). Verdict set: `compatible`, `too-old-binary`, `too-old-plugin`, `malformed-version` (replaces `malformed-range`), `no-plugin-found`; `plugin-predates-contract` is deleted — every manifest since 0.x has a `version`, so an ancient plugin classifies as `too-old-plugin` with the same install remedy. Message shapes stay the existing version-bearing ones ("Spacedock version mismatch: binary X, plugin Y…"); the `upgradeHint` keeps its current conservative parser unchanged.

### D3 — dev/source builds embed the checkout's version

Source builds (`go build`, no ldflags) currently report `spacedock dev (contract 2)` — no semver, and historically a false-green skew source. Decision: the module root gains a `go:embed` of `.claude-plugin/plugin.json` (the proven `schema_embed.go` pattern; verified cycle-free — the root package has zero `internal/` deps, so `internal/cli` may import it; the go:embed dot-directory hazard was empirically closed by the staff reviewer — explicit file paths embed fine, only wildcard patterns skip dot-directories). When `cli.Version == "dev"` (unstamped), the display version becomes `<embedded-manifest-version>+dev`, e.g. `spacedock 0.24.0-pre1+dev`. A dev binary thus always claims its checkout's minor: built from an 0.23 checkout it is rejected by 0.24 skills with the rebuild remedy; built fresh it passes. This also fixes `go install …@vX.Y.Z` builds (module proxy, no ldflags — the tagged manifest equals the tag, so they self-report `X.Y.Z+dev` and gate correctly; the `+dev` suffix on a proxy build is a cosmetic oddity noted in the doc diff). A version token with no parseable major.minor (`dev`) can then only come from an integer-era source build, and the gate treats it as too-old-binary with the rebuild hint — no always-pass carve-out anywhere.

### D4 — Cross-era sentinels (the migration off `>=2,<3`)

Removing the mechanism outright would strand integer-era installs with WRONG remedies: a 0.23 binary reading a manifest with no `requires-contract` reports `plugin-predates-contract` → "reinstall the plugin" → reinstall loops forever; a 0.23 FO prose gate seeing no `contract <N>` token classifies the binary as absent → "spacedock is not on PATH" nonsense. Two frozen literals fix both directions:

- The vendored manifests keep `"requires-contract": ">=3,<4"` as a tombstone, never edited again. An integer-era binary (contract 1 or 2) reads it and aborts `too-old-binary` — "Upgrade the binary", the correct one-step remedy.
- `--version` line 1 keeps a frozen `(contract 3)` token. Integer-era FO prose (range `>=2,<3`) sees 3 at/above the upper bound and aborts "update the plugin" — again the correct direction. The value must be 3: keeping 2 would false-green old skills against every future binary.

Both are frozen strings with comments, not maintained constants; they carry no compare math in the new binary and can be dropped at 1.0.

### D5 — FO prose gate (shared-core Startup step 1) with a release-stamped minor

The requirement must travel WITH the skill text in context (exactly as `>=2,<3` does today): a doctor- or `--version`-plugin-block-based prose gate would compare the binary against the INSTALLED plugin, which is the wrong pair in a `--plugin-dir` dev session running checkout skills. So the prose carries a release-stamped minor literal. The mechanism, in full:

- **Exactly ONE literal in the file.** The stamped phrase appears once (`These skills require binary minor 0.24 …`); every other clause in the step references "the required minor" rather than repeating the number, so no partial rewrite can yield self-contradictory instructions.
- **Stamper:** `release.StampVersion` is strictly JSON today (json.Unmarshal, then a first-`"version"` regex rewrite) and errors on markdown, so `stamp-version` grows a second handler selected by file extension: `.json` args take the existing JSON stamp; `.md` args take a prose stamp that rewrites the single line matching a pinned pattern (`These skills require binary minor <maj.min>`) to the major.minor derived from the release version, ERRORING unless exactly one match is found. One literal + one whole-file rewrite = atomic by construction.
- **All three stamp call sites carry the prose file:** releasing.md step 3's manual command, release.yml:228's main reconcile (and its `git diff`/`git commit` path lists), and release.yml:306's next edge bump (and its commit path list). Missing any one leaves a branch whose skills contradict their own manifest (the next-branch failure named in the design intro).
- **Tag gate:** `manifest-tag-gate` also takes the prose file, verifying its minor against a stable tag — same prerelease carve-out as today (release.yml:138's `if: !contains(github.ref, '-')`).
- **Drift check between stamps:** a test in `internal/contractlint` pins prose-minor == `.claude-plugin/plugin.json` minor. contractlint is the ONLY sanctioned location for tests that read instruction files (the docs/dev README's instruction-file read quarantine), and this is the legitimate shape the proof policy itself names: two independent values that can diverge — the prose literal comes from the stamped skill file, the expectation from the manifest the binary actually parses. It runs on every branch, so a hand edit to either side goes red immediately.

Before/after for Startup step 1 (the binary-absent asymmetry is load-bearing per the token-cleanup pass's SC-1 keep — the "do NOT run doctor on an absent binary" guard survives verbatim):

Before (current):
> 1. **Contract version gate.** Before discovery or boot read, run `${SPACEDOCK_BIN:-spacedock} --version` and parse `contract <N>`. Confirm `<N>` satisfies this contract's range `>=2,<3`. Abort by class:
>    - **Binary absent or non-executable** — … emits no parseable `contract <N>` token. …
>    - **Binary present but contract out of range** — `<N>` is below the lower bound (binary too old) or at/above the upper bound (plugin too old). ABORT with the mismatch message and run `${SPACEDOCK_BIN:-spacedock} doctor` for the per-class remedy.

After (the single `0.24` literal is release-stamped; every other clause says "the required minor"):
> 1. **Binary version gate.** Before discovery or boot read, run `${SPACEDOCK_BIN:-spacedock} --version` and parse line 1: `spacedock <version>`. These skills require binary minor 0.24 (same major.minor; patch and prerelease skew are fine). Abort by class:
>    - **Binary absent or non-executable** — `${SPACEDOCK_BIN:-spacedock} --version` is not found or line 1 is not `spacedock <version>`. [retry-once / ABORT / install-hint / do-NOT-run-doctor text unchanged]
>    - **Binary present but wrong version** — the version's major.minor is below the required minor (binary too old) or above it (these skills are too old — update the plugin), or the version token carries no major.minor at all (`dev` — an integer-era source build; rebuild it). ABORT with the mismatch message and run `${SPACEDOCK_BIN:-spacedock}` doctor for the per-class remedy.

The ensign shared core has no version gate; the FO shared core is the only skill text carrying the range (verified by grep across `skills/`).

### D6 — Front doors auto-heal `too-old-plugin` (claude AND codex)

Under minor-exact, EVERY binary upgrade makes the installed plugin too old at the next front-door launch. To keep the single-command UX, BOTH `runClaude` and `runCodex` (which carries an identical NoPluginFound auto-install arm today) extend that auto-install to the `too-old-plugin` verdict: refresh the plugin, then launch. `--no-install` opts out (refuse-and-instruct, as today). `too-old-binary` stays fail-fast — the front door cannot safely self-upgrade a brew binary. Two behaviors specified explicitly:

- **No scary print before a silent heal.** `gateHost` today prints the mismatch remedy for every non-Compatible verdict before returning; on the auto-heal arm that would show the operator a "version mismatch — update the plugin" remedy right before we silently fix it. Restructure: `gateHost` returns the verdict+message without printing for `too-old-plugin`; the caller prints either "Refreshing the {host} plugin…" (auto arm) or the remedy (the `--no-install` refuse arm).
- **Re-gate after install, refuse on second miss.** The auto-heal arms (both `no-plugin-found` and `too-old-plugin`) re-run the gate ONCE after the install; if the refreshed plugin still mismatches (a release-window race, a stale channel), refuse with the remedy rather than launching blind. One retry, no loop.

Codex's `--plugin-dir` differs from claude's and is treated explicitly: it is NOT a gate bypass — the checkout is installed as a real plugin and then gated (frontdoor.go's takeCodexPluginDir path). Under minor coupling that stays correct: a checkout NEWER than the binary gates `too-old-binary` there — the right verdict, since the just-installed skills genuinely outrun the binary — and the remedy is building the binary from the same checkout (which, with D3's embed, then matches exactly). Claude's `--plugin-dir` remains an ephemeral bypass, covered in-session by D5's prose gate.

The `--skip-contract-check` flag renames to `--skip-compat-check` (same behavior; "contract" no longer names anything). No alias. In-repo consumers updated: the flag definition and gate comments (frontdoor.go:780, :296, :312, :508), the bootstrap hint in `noPluginRemedy` (frontdoor.go:162), help text references (help.go:42/:83), the pi.go:408 comment, the runtime-live consumer `internal/ensigncycle/live_test.go:280`, and `docs/runtime-live-ci.md`.

### Compatibility matrix (migration record)

| | 0.23/0.24-pre plugin (integer-era) | 0.24.0+ plugin (minor-era) |
|---|---|---|
| **0.23 binary** | integer gate, unchanged | reads tombstone `>=3,<4` → too-old-binary, "upgrade the binary" ✔ |
| **0.24.0+ binary** | version compare → too-old-plugin → front door auto-refreshes ✔ | minor-exact compare ✔ |

In-session prose: 0.23 skills + new binary → old prose reads frozen `(contract 3)` → "update the plugin" ✔; new skills + 0.23 binary → new prose reads `0.23` < stamped `0.24` → "upgrade the binary" ✔; new skills + integer-era source build → `dev` token → rebuild hint ✔.

Front doors: claude/codex `too-old-plugin` → auto-refresh, re-gate, launch ✔ (or refuse under `--no-install`); codex `--plugin-dir <newer-checkout>` + older binary → real install, then `too-old-binary` with the build-from-source remedy ✔ (not a bypass; claude's `--plugin-dir` bypass stays covered by the in-session prose gate).

## Spike (ideation, 2026-07-02)

Riskiest mechanism exercised first: the real `--version` line-1 shapes across published minors and current builds, then the minor-extraction parse against them.

- Downloaded and ran the published darwin_arm64 release binaries: `v0.22.0` → `spacedock 0.22.0 (contract 1)`; `v0.23.0` → `spacedock 0.23.0 (contract 2)`; `v0.24.0-pre1` → `spacedock 0.24.0-pre1 (contract 2)`. Local source build → `spacedock dev (contract 2)`. Line 1 is uniformly `spacedock <version> (contract <N>)`; the release history carries two prerelease suffix styles (`-pre1`, `-pre.4`).
- A throwaway Go exercise of the proposed parse (cut at first `-`/`+`, split dots, take first two ints) extracted the correct major.minor from every published shape including both prerelease styles, and cleanly classified `dev` as no-semver. This seeds the implementation's first test table.
- Embed feasibility for D3: `go list -deps` confirms the module-root package (`schema_embed.go` only) imports zero `internal/` packages, so `internal/cli` importing a root-package embed is cycle-free; the `go:embed`-at-root pattern is already proven by `EntityMDSchema`.

## Acceptance criteria

- **AC-1 (VALUE)** — an old-minor binary is rejected by new-minor skills with the clean "binary too old, upgrade" abort, WITHOUT any hand-maintained integer to bump: verified by a test where restamping ONLY the release version (`stamp-version` on fixture copies of the manifest and prose; no other edit) moves the compatibility floor — the same binary version flips from compatible to too-old-binary. Tested by unit tests over the new Compare plus the stamp-then-gate fixture test.
- **AC-2** — patch- and prerelease-level skew stays compatible (0.X.0 skills ↔ 0.X.1 binary; 0.X.0-pre1 ↔ 0.X.0 — within-minor prerelease skew is the accepted residual recorded in D1), minor-level skew is rejected in BOTH directions with the correct per-direction remedy. Tested by a Compare table over the real captured version shapes, and by the real published binaries of two minors (v0.22.0, v0.23.0) in the cross-era check below.
- **AC-3** — the contract-integer MECHANISM is removed — no `CONTRACT_VERSION` constant, no `ParseRange`/integer Compare math, no bump discipline, and the FO contract text describes the minor-version gate — with exactly two frozen cross-era sentinels remaining (the manifest tombstone `requires-contract: ">=3,<4"` and the `(contract 3)` version-line token), each comment-documented as frozen and pinned by the sync test so neither can be "helpfully" bumped. (Amends the original "no requires-contract range": full removal strands integer-era installs in a wrong-remedy reinstall loop — see D4.) Tested by grep-level absence of the old mechanism plus the `internal/contractlint` sync test pinning the sentinels and prose-minor == manifest-minor (contractlint is the quarantine-sanctioned home for instruction-file reads; the check binds two independent values, the proof policy's legitimate shape).
- **AC-4** — `brew upgrade` then `spacedock claude` OR `spacedock codex` is a working single command: a `too-old-plugin` verdict at either front door auto-refreshes the plugin, re-gates once, and proceeds to launch (no mismatch-remedy print on the heal arm); a second miss after refresh refuses with the remedy rather than launching blind; `--no-install` refuses with the remedy; `too-old-binary` still fails fast. Tested via the existing `hostOps`-stub front-door tests for both hosts.
- **AC-5** — an unstamped source build reports `spacedock <checkout-manifest-version>+dev` and gates like any binary of that minor (a stale-checkout dev build is rejected by newer skills; a fresh build passes). Tested by a behavior fixture that `go build`s and runs `--version`, plus a unit test on the embed fallback.

## Test plan

- **Unit (internal/contract):** Compare table over the spiked real shapes (`0.22.0`, `0.23.0`, `0.24.0-pre1`, `0.23.0-pre.4`, `dev`, `0.24.0-pre1+dev`) × both directions × patch/prerelease interchange; malformed/missing manifest version; verdict-to-message/exit mapping (doctor). Low cost.
- **Sync test (successor of the #410 bracketing test) in `internal/contractlint`:** `.claude-plugin` and `.codex-plugin` manifest minor == the shared-core prose stamped minor (the prose read lives here per the instruction-file read quarantine); tombstone field frozen at `">=3,<4"`; `(contract 3)` token present on line 1 of `--version` output. The manifest-only halves of `skills/integration/plugin_manifest_test.go` / `marketplace_manifest_test.go` / `contract_skew_test.go` migrate to the new compare. Low cost.
- **Release tooling:** the `.md` prose-stamp handler (rewrites the pinned single-literal line; errors on zero or multiple matches; leaves `.json` handling untouched); `stamp-version` stamps manifests + prose in one invocation (fixture round-trip); `manifest-tag-gate` rejects a stable tag whose prose minor disagrees (prerelease carve-out preserved); AC-1's divergeable stamp-then-gate test. Low cost.
- **Startup-gate behavior fixture (gate_test.go successor):** stub binary printing each real shape drives the gate; out-of-range aborts before `status --discover` with the pinned per-direction remedy. Moderate cost, pattern exists.
- **Front door (frontdoor_test.go pattern, both hosts):** AC-4 cases — auto-refresh + re-gate-once + second-miss refusal + `--no-install` + no-remedy-print-on-heal, for claude AND codex; the codex `--plugin-dir` install-then-gate path with a newer checkout. Low cost.
- **Dev-build fixture:** AC-5 `go build` + `--version` shape + gate pass against the checkout. Moderate cost.
- **Real-binary cross-era check (local/e2e tier, commands recorded):** the downloaded published v0.22.0/v0.23.0 binaries' `doctor` against a tombstone-manifest fixture → "Upgrade the binary" exit 1; the new binary against a real 0.23 manifest → too-old-plugin. Proves the migration against actual shipped artifacts, not stubs.
- **Detached adversarial audit** at validation (contract/release high-stakes surface) — anticipated: the audit should attack the sentinel freeze (can anything bump them?), the prose/manifest drift window, and the dev-build minor claim.

## Documentation diff (user-visible changes)

- `docs/releasing.md` step 3: `go run ./cmd/spacedock-release stamp-version X.Y.Z .claude-plugin/plugin.json .codex-plugin/plugin.json` → `go run ./cmd/spacedock-release stamp-version X.Y.Z .claude-plugin/plugin.json .codex-plugin/plugin.json skills/first-officer/references/first-officer-shared-core.md` (and the same file added to the release commit's path list on the next line).
- `.github/workflows/release.yml`: the prose file joins the `stamp-version` argument lists AND the `git diff`/`git commit` path lists at BOTH CI stamp sites (the main reconcile at :228 and the next edge bump at :306), plus the `manifest-tag-gate` argument list at :138.
- `docs/runtime-live-ci.md` (plugin-install row): `spacedock claude --plugin-dir <checkout> --skip-contract-check` → `spacedock claude --plugin-dir <checkout> --skip-compat-check`.
- Proof-policy worked example (`docs/site/contributing/proof-policy.md:13` and `docs/dev/README.md:78`): "such as the plugin manifest's contract range bracketing the binary's contract version" → "such as the plugin manifest's version sharing a major.minor with the binary's version" (the old example names a deleted mechanism).
- `--version` line 1 for source builds: `spacedock dev (contract 2)` → `spacedock 0.24.0-pre1+dev (contract 3)`; release builds keep their shape with the frozen token. `go install …@vX.Y.Z` proxy builds report `X.Y.Z+dev` (correct gating; the `+dev` marks the unstamped build path) — one line in the install docs.
- Doctor verdict tokens: `malformed-range` → `malformed-version`; `plugin-predates-contract` removed (classifies as `too-old-plugin`). Message bodies already speak in display versions, unchanged in shape.

## Stage Report: ideation

- DONE: The exact gate semantics are DECIDED and recorded with rationale (minor-exact vs version-floor, forward-compat for a newer binary, patch-skew interchange, prerelease suffixes like 0.24.0-pre1, and the source-build `spacedock dev` version string that carries no semver) plus the migration path off `requires-contract: >=2,<3` — decisions, not an options list.
  D1–D6 in the body: minor-exact both directions (rationale vs floor recorded in D1); no forward-compat assumption, front door auto-heals instead (D6); prerelease cut at first `-`/`+` (D2); dev builds embed the checkout manifest version so `dev` only ever means an integer-era build with a rebuild remedy (D3); migration = two frozen cross-era sentinels with a full compatibility matrix (D4).
- DONE: AC-1's divergeable measure holds in the design: a test where bumping ONLY the release version (no separate hand-maintained integer edit) already moves the compatibility floor the right way.
  AC-1 rewritten as the stamp-then-gate fixture test: `stamp-version` on manifest+prose copies flips the same binary version from compatible to too-old-binary with no other edit; the sync test makes hand-drift between stamps impossible.
- DONE: The riskiest unverified mechanism is spiked first against real binaries — parse the actual `--version` output shapes across two published minors and the current prerelease/dev builds — and the result recorded in the body (or "no spike needed" with proven mechanisms named).
  `## Spike` section: ran published v0.22.0/v0.23.0/v0.24.0-pre1 darwin_arm64 binaries plus the local source build; parse exercise handled every shape incl. both prerelease styles (`-pre1`, `-pre.4`); embed cycle-check for D3 verified via `go list -deps`.

### Summary

Replaced the options-list approach section with a decided design: minor-exact coupling derived from the existing release stamps (no new declared field — the manifest `version` is the declaration), two frozen cross-era sentinels (`requires-contract: ">=3,<4"` tombstone + `(contract 3)` version token) so integer-era installs abort with correct remedies in both directions, a release-stamped minor literal in the FO prose gate, front-door auto-refresh on too-old-plugin, and embed-versioned dev builds. AC-3 was amended (sentinels vs full removal) and AC-4/AC-5 added; the gate review should look at the AC-3 amendment, the D6 auto-refresh behavior change, and the `--skip-contract-check` → `--skip-compat-check` rename.

### Repairs after staff review (2026-07-02)

All five material findings verified against the code and repaired in the body; the design/AC/test-plan/doc-diff sections were amended together.

- DONE: Release-train premise error + prose-stamp coverage gap.
  Verified release.yml — :228 reconcile-stamps MAIN to stable X.Y.Z, :306 stamps the `next` edge branch to X.(Y+1).0-pre1; design intro corrected, all THREE stamp call sites (+ the :138 tag-gate arg list) enumerated in the intro/D5/doc diff, each carrying the prose file, with the next-branch contradiction named as the consequence.
- DONE: runCodex omitted from D6/AC-4.
  Verified frontdoor.go:496-560; D6 extended to both hosts, codex `--plugin-dir` install-then-gate treated explicitly (correct too-old-binary on a newer checkout, build-from-source remedy), matrix gains a front-doors row, AC-4/test plan cover both hosts.
- DONE: Prerelease-interchange residual false-green.
  Decided: accepted residual, recorded in D1 with rationale — strict prerelease matching false-reds the whole next channel because prerelease tags deliberately stamp no manifests; exposure is next/edge-only while the v0.23.0 class hit stable users at a minor boundary. Seed intro's "cannot recur" narrowed to "across minors". Mid-cycle dev-skew granularity limit folded in.
- DONE: Sync-test placement violates the Proof policy quarantine.
  Moved to `internal/contractlint` (D5, AC-3, test plan) — the sanctioned instruction-file-read home; the check binds two independent values (prose literal vs manifest), the policy's legitimate shape.
- DONE: Prose-stamping mechanism undesigned.
  Verified release.go:74-84 is JSON-only; D5 now specifies an extension-selected `.md` handler rewriting a single pinned-pattern literal (error unless exactly one match), the after-text was rewritten to carry the `0.24` literal exactly once with other clauses saying "the required minor", and the tooling tests cover zero/multi-match errors.
- DONE: Minor findings folded in.
  Full `--skip-compat-check` consumer list in D6 (frontdoor.go:780/:296/:312/:508/:162, help.go:42/:83, pi.go:408, live_test.go:280, runtime-live-ci.md); gateHost no-remedy-print-on-heal + re-gate-once-then-refuse in D6/AC-4; proof-policy worked-example rewording (proof-policy.md:13, docs/dev/README.md:78) and the `go install` `+dev` cosmetic in the doc diff; go:embed dot-path closure cited in D3.
