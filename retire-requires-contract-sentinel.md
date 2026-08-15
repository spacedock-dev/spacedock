---
id: 6qhgsezz7v4g4h76t0jf98b0
title: Retire the extinct requires-contract manifest sentinel
status: validation
source: "0.27 audit (2026-08-14) Priority 2; pre-ship cleanup companion to remove-startup-capability-probe"
started: 2026-08-15T02:55:26Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-retire-requires-contract-sentinel
issue:
gates:
    version: 1
    records:
        - id: gate:6qhgsezz7v4g4h76t0jf98b0:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:6qhgsezz7v4g4h76t0jf98b0-backlog-1
              briefing:
                id: briefing:6qhgsezz7v4g4h76t0jf98b0:backlog:attempt-1:revision-1
                digest: sha256:886472682b4b65555d37be5b39151f04f478b618bebafce013e3c41dea56d935
                request-digest: sha256:e3dfe8c776eebd1adfb7660b57d36d1886b392d9e26c2d6f3cc392d374eab563
                room-ref: ./retire-requires-contract-sentinel/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:6qhgsezz7v4g4h76t0jf98b0:backlog:1
                briefing: briefing:6qhgsezz7v4g4h76t0jf98b0:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T02:53:19.566854Z"
                decision: approve
                reason: 'Captain ruling 2026-08-14 (dispatch them): approved into ideation'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:6qhgsezz7v4g4h76t0jf98b0:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:6qhgsezz7v4g4h76t0jf98b0-ideation-1
              briefing:
                id: briefing:6qhgsezz7v4g4h76t0jf98b0:ideation:attempt-1:revision-1
                digest: sha256:b1b75c4742f4d15d4dddd66d1ca1c1c970a0d9947cbc477f543a9340638e58bf
                request-digest: sha256:9f349e885539ce7a8a27ef58a086c51203f911d183dc931dd530aaf49aaed13b
                room-ref: ./retire-requires-contract-sentinel/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:6qhgsezz7v4g4h76t0jf98b0:ideation:1
                briefing: briefing:6qhgsezz7v4g4h76t0jf98b0:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T03:55:08.408275Z"
                decision: approve
                reason: 'Captain ruling 2026-08-15 (approve all except x8): approved into implementation'
              application:
                target-stage: implementation
                state: consumed
        - id: gate:6qhgsezz7v4g4h76t0jf98b0:validation
          stage: validation
          attempts:
            - id: gate-attempt:6qhgsezz7v4g4h76t0jf98b0-validation-1
              briefing:
                id: briefing:6qhgsezz7v4g4h76t0jf98b0:validation:attempt-1:revision-1
                digest: sha256:0607c0748dda2c474d48c8607abb6dbb06562ea2ce49cece551ef10eacffc7a2
                request-digest: sha256:13fa83a060fcff5df3403aa1315827095db27b88a1a2933c563f74d3ad6d049e
                room-ref: ./retire-requires-contract-sentinel/review/validation/briefing-1
---

Retire the pre-0.19 `requires-contract` sentinel; its audience is extinct and production code explicitly ignores it. Pure deletion, no behavior change.

The audit scoped this as "testdata fixtures only", but a live grep (2026-08-14) finds 8 source/test files referencing it: internal/contract/doctor.go (a comment stating the field "is not read"), internal/release/release_test.go, internal/cli/codex_plugin_dir_test.go, internal/cli/codex_name_match_test.go, internal/cli/upgrade_from_stale_test.go, internal/cli/decoupling_behavior_test.go, internal/cli/install_behavior_codex_test.go, skills/integration/plugin_manifest_test.go. Historical mentions in docs/roadmap and archived state entities stay — they are records, not live surface.

## Problem

`requires-contract` is a pre-0.19 plugin-manifest compatibility sentinel whose value chain is broken at all three links. Confirmed at HEAD `4d1912a69`:

- **No writer.** Neither shipped manifest carries the field. `.claude-plugin/plugin.json` and `.codex-plugin/plugin.json` declare `version` and nothing else compatibility-related.
- **No reader.** `readManifest` (`internal/contract/doctor.go:21`) unmarshals into a one-field struct `{ Version string }`. No Go symbol in the repo references the key — and for a JSON key that grep is complete, since any reader would have to name the literal in a struct tag or a map lookup.
- **No verifier.** `internal/contractlint/prose_manifest_minor_sync_test.go` reads only `version`.

What survives is 15 matching lines across 8 live files that keep the dead field alive as fixture noise and, worse, as false documentation. Two comments assert a guard that does not exist: `skills/integration/plugin_manifest_test.go:29-32` and `skills/integration/marketplace_manifest_test.go:38-40` both state that the "D4 cross-era tombstone" is "pinned by the internal/contractlint sync test". That test pins the manifest `version` against the FO shared-core's stamped prose minor and never looks at a tombstone. Anyone auditing the compatibility surface is told a check exists that does not.

The ongoing cost is copy-paste: 5 of the 8 sites are inline manifest literals in test helpers, so every new host fixture inherits the field. Note the audit's "testdata fixtures only" framing was wrong in both directions — the scope is larger than testdata, and there are no testdata files at all. Every site is either an inline Go string literal or a comment.

### The other D4 sentinel stays

D4 shipped **two** cross-era sentinels. Only the manifest field is extinct. `frozenContractToken = "contract 3"` (`internal/cli/cli.go:859`, emitted at `:901`) has a fully intact value chain — writer `printVersion`, consumers being integer-era FO prose plus `docs/site/reference/command-reference.md:38`, and verifier `internal/cli/version_session_test.go` (6 assertions including `TestVersionContractTokenPlacement`). It is untouched here. Its own retirement condition is documented in place at `cli.go:855-858` and is a separate future call.

## Proposed approach

Pure deletion plus comment correction across 10 files. No behavior change, no new test, no shipped-artifact change.

**Delete the field from 5 inline fixture manifests** (the code under test reads only `version`, so the field is inert in every one):

- `internal/cli/codex_plugin_dir_test.go:332` (`">=2,<3"`) and `:343` (`">=3,<4"`) — the two values differ arbitrarily; nothing reads either.
- `internal/cli/codex_name_match_test.go:103`, `internal/cli/decoupling_behavior_test.go:126`, `internal/cli/install_behavior_codex_test.go:187` (all `">=2,<3"`).

Each edit is in-place on one line: `{ "name": "spacedock", "version": …, "requires-contract": ">=N,<M", "skills": "./skills/" }` becomes `{ "name": "spacedock", "version": …, "skills": "./skills/" }`.

**Delete the field and its assertion from `internal/release/release_test.go`.** Drop `"requires-contract": ">=1,<2"` from the `TestStampVersionRewritesPluginVersion` fixture (moving the trailing comma up to the `"skills"` line) and remove the 3-line assertion at `:39-41`. The test's "untouched fields survive" property keeps real coverage through the surviving `name` and `skills` assertions — those two must stay.

**Correct 5 comments.** These are the load-bearing edits; three of them currently mislead:

1. `internal/contract/doctor.go:18-19` — replace "requires-contract, if present, is not read — it is a frozen cross-era sentinel for integer-era binaries only (D4)." with "Every other manifest field is ignored."
2. `internal/cli/upgrade_from_stale_test.go:2` and `:128-130` — the fixture's staleness is driven by its *version* (`0.12.1`), not by a missing field; the test body already says so (D2 note at `:18-19`). Re-anchor both comments: "install of a stale (no requires-contract) plugin" becomes "install of an ancient-minor (0.12.1) plugin", and "whose plugin manifest has NO requires-contract field — the 0.12.1 shape that predates the contract mechanism" becomes "whose plugin manifest declares version 0.12.1 — an ancient minor that resolves to too-old-plugin". The fixture itself is unchanged.
3. `internal/cli/decoupling_behavior_test.go:102` and `internal/cli/install_behavior_codex_test.go:162` — "The plugin manifest carries a requires-contract bracketing CONTRACT_VERSION." becomes "The plugin manifest carries the display version the doctor verdict reads."
4. `internal/cli/codex_plugin_dir_test.go:325` — "a .codex-plugin/plugin.json bracketing CONTRACT_VERSION + one skill" becomes "a .codex-plugin/plugin.json with a version + one skill". (Same stale framing, no literal — found only by widening the grep.)
5. `skills/integration/plugin_manifest_test.go:25-32` — drop the false "D4 cross-era tombstone … pinned by the internal/contractlint sync test" sentence and state what that test actually pins: the manifest `version`'s binding to the FO shared-core's stamped minor.

**Two adjacent files, deliberately included.** `skills/integration/marketplace_manifest_test.go:38-40` and `internal/contractlint/prose_manifest_minor_sync_test.go:2` describe the extinct tombstone *without using the literal string*, so a grep-scoped sweep would leave them behind — still claiming a verifier for a field that no longer appears anywhere. Leaving them is leaving the exact confusion this entity exists to remove, at a cost of two comment lines. Same edit as above: say that the sync test pins the manifest minor against the stamped prose minor, and drop the tombstone clause.

**Simplest alternative considered:** delete only the 8 literal-matching files and stop. Rejected — it satisfies a grep while leaving two live comments asserting a guard that does not exist, which is the failure mode the entity is named after. No new mechanism is introduced by either option.

### No-spike determination

**No spike needed: pure deletion.** The design rests on one load-bearing claim — that removing the field changes nothing observable — and on three proven mechanisms: `readManifest`'s single-field unmarshal (read directly), the completeness of a literal grep for a JSON key, and `frozenContractToken`'s independent pinning in `version_session_test.go`.

The claim was exercised anyway rather than asserted, since it was nearly free. A throwaway differential test (written, run, deleted — not committed) called `ManifestVerdict` and `ManifestVersion` on two manifests identical but for the field, across 3 hosts x 3 binary versions spanning all three verdict classes. Verdict and user-facing message were byte-identical in all 9 combinations: `compatible` at 0.26.0, `too-old-plugin` at 0.27.0, `too-old-binary` at 0.25.0. This is why the entity ships **no** regression test — see the keep-boundaries below.

### Keep-boundaries

1. `frozenContractToken` / `"contract 3"` (`internal/cli/cli.go`, `version_session_test.go`, `dev_version_test.go`, `gate_test.go`, `docs/site/reference/command-reference.md`) — the live second sentinel. Untouched.
2. `docs/roadmap/0199-pre-flip-mechanics/debrief.md` (2 mentions) — the historical record of the session that deleted the false doctor note. Untouched.
3. Archived state entities under `docs/dev/.spacedock-state/` — records, not live surface. Untouched.
4. `release_test.go`'s `name` and `skills` assertions — the fixture must retain at least two non-version fields or "untouched fields survive" stops proving anything.
5. The `buildStaleMarketplace` fixture manifest — unchanged; only its describing comment moves.
6. `internal/contract/contract.go`'s integer-era handling (`:22`, `:72`, `:112`) — that is the `dev`-sentinel path in *binary version* parsing, unrelated to the manifest field. Untouched.
7. **Do not add a regression test asserting the field is ignored.** It would reintroduce the literal and fail AC-1. The unchanged test inventory plus a green suite is the proof.

### Documentation

No doc diff needed: no user-visible surface changes. `docs/site/reference/command-reference.md:38` describes the *live* `contract 3` sentinel and stays as-is. No docs-site page mentions `requires-contract`.

### Coordination

No file overlap with the sibling entity `remove-startup-capability-probe`, which touches `skills/first-officer/references/first-officer-shared-core.md`, `skills/integration/testdata/version_gate_flow.sh`, `skills/integration/version_gate_fixture_test.go`, and `docs/site/get-started/install.md`. Both land edits in the `skills/integration` package but in disjoint files, and this entity adds and removes no symbols there, so neither ordering produces a conflict or a compile break.

## Out of scope

Priority-3 audit items (`gate validate` demotion, `provider-evidence` schema strip, `hold` decision) — each needs its own captain call.

## Expected surface and tolerance

**Files: 10.** Eight carry the literal; two carry the same stale claim without it.

| File | Edit | Net lines |
| --- | --- | --- |
| `internal/contract/doctor.go` | comment | -1 |
| `internal/release/release_test.go` | fixture field + assertion | -4 |
| `internal/cli/codex_plugin_dir_test.go` | 2 fixtures + 1 comment | 0 |
| `internal/cli/codex_name_match_test.go` | fixture | 0 |
| `internal/cli/decoupling_behavior_test.go` | fixture + comment | 0 |
| `internal/cli/install_behavior_codex_test.go` | fixture + comment | 0 |
| `internal/cli/upgrade_from_stale_test.go` | 2 comments | 0 |
| `skills/integration/plugin_manifest_test.go` | comment | -1 |
| `skills/integration/marketplace_manifest_test.go` | comment (adjacent) | -1 |
| `internal/contractlint/prose_manifest_minor_sync_test.go` | ABOUTME (adjacent) | 0 |

**Net LOC: about -7.** Tolerance -4 to -12; comment rewrapping can move a line either way. A positive net change is a breach.

**Tolerance on file count: +1**, and only for a site in one of the two classes above (a literal-matching file, or a comment claiming the tombstone is verified). Any file outside those classes is a boundary breach the implementation report must name rather than absorb.

**Semantic changes: none.** Declared explicitly, since a small diff is exactly where an undeclared semantic hides:

- Command grammar: unchanged.
- Stored formats: unchanged. The shipped `.claude-plugin/plugin.json` and `.codex-plugin/plugin.json` are not touched — they never carried the field.
- Authority: unchanged.
- Runtime behavior: unchanged, exercised across all three verdict classes (see the no-spike determination).
- Public Go API: unchanged. No symbol added, removed, or renamed.
- Test inventory: unchanged at 1377 names.

## Acceptance criteria

**AC-1 (value) — No live source, test, or fixture file carries or describes the sentinel as real.**
The live reference count falls from its measured baseline to zero, while the historical record is preserved intact.
Verified by: `git grep -c "requires-contract" -- internal skills cmd docs/site` returns nothing (baseline: 15 matching lines across 8 files), and `git grep -c "requires-contract"` over the whole tracked tree returns exactly `docs/roadmap/0199-pre-flip-mechanics/debrief.md:2` (baseline: 17 lines across 9 files). Additionally `git grep -n "tombstone" -- internal skills` returns nothing, closing the two adjacent sites that make the same claim without the literal.

**AC-2 (value guard) — The reference count reaches zero by deletion, not by dropping coverage.**
The test inventory is byte-identical before and after: 1377 names, same set. This is the criterion that makes AC-1 non-gameable — deleting the fixtures' owning tests would satisfy AC-1 just as well and is the cheapest wrong way to do this task.
Verified by: `go test ./... -list '.*' | grep -E '^(Test|Example|Fuzz)' | sort` diffed against the baseline captured at `4d1912a69` produces no output. (`-list` compiles without running, so it is unaffected by the live-credential caveat below.)

**AC-2b (suite) — No package's test result changes from its `4d1912a69` baseline.**
Stated as a differential rather than an absolute green, because an absolute green is not achievable on a developer box — see the baseline caveat in the test plan.
Verified by: `go test ./...` and `go test ./... -race` run in the credential-free shape CI's `offline` job uses, with the four affected packages (`internal/contract`, `internal/release`, `internal/contractlint`, `skills/integration`) green — all four confirmed green at `4d1912a69` during ideation.

**AC-3 (no-behavior-change) — The shipped artifacts and the compatibility verdict are untouched.**
Verified by: `git diff --stat -- .claude-plugin .codex-plugin` is empty, and the `internal/contract` package tests — which exercise `ManifestVerdict` across the compatible / too-old-plugin / too-old-binary classes — stay green under AC-2's run.

**AC-4 (keep-boundary) — The live `contract 3` sentinel survives the sweep.**
An over-eager D4 cleanup that also removed `frozenContractToken` would break every integer-era reader's abort path, and AC-1 would not notice.
Verified by: `go test ./internal/cli/ -run 'TestVersion'` passes, and `git grep -c '"contract 3"' -- internal/cli/cli.go` still reports the constant.

## Test plan

**No new tests.** The deletion's correctness is carried by the existing suite staying green with an unchanged test inventory. A new test asserting "the field is ignored" would have to name the literal and would fail AC-1 — the throwaway spike already exercised that property (9 host x version combinations, all three verdict classes, byte-identical messages) and was deleted rather than committed.

What verifies each half of the change:

- **The 6 fixture-field deletions** are behavioral: the affected fixtures feed `internal/cli` and `internal/release` tests that already assert on manifest-derived outcomes, so if the field were load-bearing anywhere those tests go red. Cost: one full-suite run.
- **The 7 comment corrections** have no behavior to exercise; a textual check is the appropriate proof, not a stand-in for a missing behavioral one. AC-1's grep is that check.

Cost and complexity: low. No fixture, CLI, or live-workflow tests are needed. Nothing here touches the live-host paths — `upgrade_from_stale_test.go` receives comment-only edits, so the implementer needs no `claude` on PATH (those tests skip without it, in CI as well).

Baseline artifacts an implementer needs, both reproducible from `4d1912a69`: the 1377-name test inventory (`go test ./... -list '.*'`) and the 15-line live reference count. Both are captured from a clean tracked tree — untracked scratch under `.worktrees/` and `.pi-subagents/` carries vendored copies of these files and must be excluded, which `git grep` does by construction and `grep -r` does not.

### Baseline caveat: `go test ./...` does not finish on a credentialed box

Measured at `4d1912a69` during ideation, before any edit: a bare `go test ./...` **fails**, with `internal/ensigncycle` hitting the default 10-minute timeout at 600.5s. Every other package passes (13 `ok`). This is pre-existing and unrelated to this entity, which touches no file in that package.

The cause is environmental, not a red repo. `internal/ensigncycle` is the live-agent harness; its tests self-skip when no host credential is present, which is why CI's `offline` job runs a bare `go test ./...` green — it runs **without secrets**. On a developer box inside a Claude session the ambient credential is present, so the live legs actually execute and blow the timeout. Confirmed by running CI's own deterministic control subset from `runtime-live-e2e.yml:78`, which passes in 1.0s.

Consequence for whoever implements this: do not chase that failure, and do not treat a bare `go test ./...` as this entity's gate. Prove the change against the four affected packages plus the inventory diff. The four were confirmed green at `4d1912a69` (`internal/contract` cached, `internal/release` 272s, `internal/contractlint` 6.3s, `skills/integration` 39.3s). The same caveat applies to the sibling entity `remove-startup-capability-probe` running on this box.

## Stage Report: ideation

- DONE: Design confirms the 8-file reference list at HEAD; historical roadmap and archive records stay untouched
  `git grep -l "requires-contract" HEAD` at `4d1912a69` returns exactly the 8 named files plus `docs/roadmap/0199-pre-flip-mechanics/debrief.md`; the roadmap record and the archived state entities are named as keep-boundaries 2 and 3.
- DONE: Value AC: zero live requires-contract references and a green suite
  AC-1 measures the live count 15 lines/8 files to 0; AC-2 guards it with an unchanged 1377-name test inventory so the count cannot be reached by deleting tests. The "green suite" half became AC-2b, restated as a per-package differential — see the FAILED item.
- DONE: No-spike determination recorded: pure deletion
  Recorded under "No-spike determination" with the three mechanisms it rests on, and the one load-bearing claim was exercised anyway rather than asserted.
- FAILED: Confirm a green `go test ./...` baseline at HEAD
  A bare `go test ./...` at `4d1912a69`, before any edit, fails: `internal/ensigncycle` hits the default 10m timeout at 600.5s while 13 other packages pass. Environmental, not a red repo — that package is the live-agent harness and its legs fire on a credentialed box, where CI's `offline` job runs without secrets. CI's own deterministic subset from that package passes in 1.0s. Recorded as a baseline caveat and folded into AC-2b; the sibling ensign was notified.

### Summary

Scope confirmed exactly as the entity described it — 8 files, 15 live lines — then widened by 2 files. `skills/integration/marketplace_manifest_test.go` and `internal/contractlint/prose_manifest_minor_sync_test.go` describe the same extinct tombstone without using the literal string, so a grep-scoped sweep would leave them behind still claiming a verifier for a field that no longer exists anywhere; including them costs two comment lines. Three of the comment edits correct claims that are false at HEAD, the sharpest being that the contractlint sync test "pins the D4 cross-era tombstone" when it reads only `version`.

The critical keep-boundary is that D4 shipped **two** cross-era sentinels and only the manifest field is extinct: `frozenContractToken = "contract 3"` has a writer, consumers, and 6 assertions pinning it, so AC-4 exists specifically to catch a sweep that takes both. The "not read" claim was proven by exercising rather than by grep — a throwaway differential over `ManifestVerdict` across 3 hosts x 3 binary versions returned byte-identical verdicts and messages in all 9 combinations, spanning compatible, too-old-plugin, and too-old-binary. It was deleted, not committed, and the entity ships no regression test because one would have to name the literal and fail AC-1.

## Stage Report: implementation

- DONE: Ten-file scope including the two comment-only files; the live contract-3 sentinel untouched per AC-4
  Commit `571ac9ec7` touches exactly the 10 declared files, +22/-29 = net -7 (declared estimate -7, tolerance -4..-12); no file outside the two declared classes. AC-4: `git grep -c '"contract 3"' -- internal/cli/cli.go` still reports the constant, and `go test ./internal/cli/ -run 'TestVersion'` passes 6 tests including `TestVersionContractTokenPlacement`, which goes red if `frozenContractToken` is deleted or stops being emitted in the version block.
- DONE: Zero live requires-contract references and the 1377-name test inventory unchanged
  AC-1 both halves: `git grep -c "requires-contract" -- internal skills cmd docs/site` and `git grep -n "tombstone" -- internal skills` each exit 1 with no output; the whole-tree grep returns only `docs/roadmap/0199-pre-flip-mechanics/debrief.md:2`. AC-2: the `go test ./... -list` inventory diff against the pre-edit tree is empty, and `git grep -hE '^func (Test|Example|Fuzz)'` is 1416 at both `4d1912a69` and HEAD. The inventory number is corrected to **1375** — see the correction below.
- DONE: Suite verdicts per AC-2b differential; reproduce any failure on a clean control before attributing
  All four AC-2b packages green plain and under `-race` (`internal/contract` 4.7s, `internal/release` 207.2s, `internal/contractlint` 31.8s, `skills/integration` 37.7s). `internal/cli` — which AC-2b omits despite owning 5 of the 10 edited files — was run too: it fails on this credentialed box, and a clean control at `4d1912a69` in a throwaway checkout reproduces the same failures with no edits present. In CI's credential-free shape it is green: `env PATH=/usr/bin:/bin:/usr/sbin:/sbin go test ./internal/cli/` returns `ok ... 442.7s`.

### Correction: the inventory baseline is 1375, not 1377

The entity records 1377 names; the measured count at `4d1912a69` is 1375, from both a clean worktree and the main checkout (identical sets, empty diff). 1377 does not reproduce by `go test ./... -list` (1375) or by declaration grep (1416), so it appears to be an ideation mis-measurement rather than a change in the tree. AC-2's load-bearing property is unaffected: the property is that the inventory set is byte-identical before and after, and that diff is empty. What would falsify it is deleting or renaming any test function — including the cheapest wrong way to satisfy AC-1, deleting the fixtures' owning tests.

### Correction: the credentialed-box caveat covers `internal/cli`, not just `internal/ensigncycle`

Ideation recorded `internal/cli` as `ok` at baseline and scoped the caveat to `internal/ensigncycle`. On this box `internal/cli` fails two independent ways, both reproduced on the clean `4d1912a69` control before attributing:

1. `TestCodexResolveManifestAgainstInstalledHost` fails in ~0.4s because the real `codex plugin list` cannot read `/Users/clkao/.codex/config.toml` (`Operation not permitted`) — a sandbox permission condition, not a manifest condition; the test's own guard is `exec.LookPath("codex")`, so it skips where codex is absent.
2. The package hits the test timeout on a live-host leg. The test still running at the panic differs between runs — `TestClaudePluginInstallIsHostNative` on the candidate, `TestGateRecordHaltsOnSameEntityConflict` on the control — so the hang tracks live-subprocess scheduling, not any edited file.

Removing `claude` and `codex` from PATH turns both off and the package passes, which is why CI's secretless `offline` job is green. Validation should not treat a bare credentialed `go test ./internal/cli/` as this entity's gate.

### Summary

Pure deletion landed as designed: the field is gone from the 6 inline fixture manifests and from `release_test.go`'s assertion (its `name` and `skills` assertions stay, so "untouched fields survive" still proves something), and 7 comments were corrected — including the three that falsely claimed the `internal/contractlint` sync test pins a D4 tombstone it never reads. The two adjacent comment-only files were included as ideation directed; without them a grep-clean tree would still have asserted a verifier for a field that no longer exists.

Two ideation figures needed correcting and both are recorded above: the test-inventory baseline is 1375 rather than 1377, and the credentialed-box test caveat extends to `internal/cli` rather than `internal/ensigncycle` alone. Neither changes the outcome — AC-1 through AC-4 all hold, with `internal/cli` green in the credential-free shape AC-2b specifies. No behavior change: `git diff --stat -- .claude-plugin .codex-plugin` is empty, no Go symbol moved, and `gofmt -l` is clean over `cmd`, `internal`, and `skills`.

## Stage Report: validation

- DONE: Independently re-exercise the ACs against worktree commit 571ac9ec7, never by reading the report: repo grep for requires-contract empty over internal, skills, cmd; test-name inventory byte-identical before/after (corrected baseline 1375); net delta negative; the live contract-3 sentinel and its six version_session_test assertions untouched
  AC-1: `git grep -in "requires-contract" -- internal skills cmd docs/site` and `git grep -n "tombstone" -- internal skills` both exit 1 empty; the whole-tree grep returns only `docs/roadmap/0199-pre-flip-mechanics/debrief.md:2`; variant spellings (`requires_contract`, `requiresContract`) absent. AC-2: inventory measured on a git-archive control at `4d1912a69` vs the worktree — 1375 names both sides, diff empty; this goes non-empty if any test is deleted or renamed, including the cheapest wrong path of deleting the fixtures' owning tests. AC-4: `git grep -c '"contract 3"' -- internal/cli/cli.go` reports 1; `go test ./internal/cli/ -run 'TestVersion'` passes 6/6 including `TestVersionContractTokenPlacement`, which reds if `frozenContractToken` stops being emitted; `version_session_test.go` is absent from the commit's file list. Surface: exactly the 10 declared files, +22/-29 = net -7, negative and inside the declared -4..-12.
- DONE: Suite verdicts per the AC-2b differential: reproduce any failure on a clean 4d1912a69 control before attributing; the credentialed-box caveat covers internal/cli AND internal/ensigncycle - never gate on a bare credentialed run of either
  Credential-free shape (PATH shim exposing only `go`; `command -v claude`/`codex` resolve nothing): the four AC-2b packages green plain (contract 6.8s, release 33.2s, contractlint 5.0s, skills/integration 30.1s) and under -race (release 102.8s), both exit 0 — no failure arose, so no control attribution was needed. `internal/cli`, which owns 5 of the 10 edited files but sits outside AC-2b's package list, was additionally run credential-free uncached: ok 186.3s, independently confirming the implementation's credential-free green. No bare credentialed run of internal/cli or internal/ensigncycle was used as a gate.
- DONE: Verdict PASSED or REJECTED with per-AC evidence citations; the two recorded figure corrections (1375 inventory, cli caveat widening) are report facts to confirm, not defects
  PASSED. AC-3: `git diff --stat 4d1912a69..HEAD -- .claude-plugin .codex-plugin` is empty and both shipped manifests were inspected — neither ever carried the field; `internal/contract` (all three verdict classes) green under the AC-2b run. The 1375 figure is confirmed by independent measurement on the clean control (1377 does not reproduce). The cli-caveat widening is confirmed structurally: the five fixture-consuming cli tests flip from SKIP (credential-free) to RUN (hosts on PATH), which is exactly the live-leg activation the correction describes.

### Adversarial pass beyond the ACs

The credential-free green of `internal/cli` does not execute the five edited fixture consumers — all guard on `exec.LookPath("claude"/"codex")` and SKIP without hosts, in CI's offline job too, so a green offline run cannot falsify the fixture edits. Closed by exercising one consumer per fixture class individually with real hosts on PATH: `TestCodexEntryNameMustMatchPluginName` PASS 0.63s (real codex installs from the fieldless `.codex-plugin` manifest) and `TestStableChannelDecoupledFromBranchHead` PASS 33.93s (real claude runs the full channel-decoupling flow over the fieldless `.claude-plugin` manifest). Each fails if its host rejects or misreads the edited manifest. Every corrected comment was verified true at HEAD by reading the code it describes: `readManifest` unmarshals only `{Version string}`, the stale fixture declares `0.12.1`, and the contractlint sync test pins the manifest minor against `release.ProseMinor` with no tombstone read. `gofmt -l` clean over cmd, internal, skills; worktree clean at `571ac9ec7`.

Deferred risk (pre-existing harness shape, not introduced by this change): the five fixture-consuming `internal/cli` tests execute only where live hosts are on PATH, so offline CI never exercises those fixture manifests. No AC fails, and both fixture classes passed live on this box. Promote to material only if the fixture manifest schema changes again without a live-host run in the loop.

### Summary

PASSED. All four ACs re-exercised independently at worktree commit `571ac9ec7` against a clean `4d1912a69` control, with no material findings and one pre-existing deferred risk recorded above. The live reference count is zero with the historical debrief record intact, the 1375-name inventory is byte-identical, the four AC-2b packages are green plain and under -race in the credential-free shape, the shipped manifests are untouched, and the `contract 3` sentinel survives with all six version assertions passing. Both of the implementation's figure corrections check out as facts: 1375 reproduces on the control and 1377 does not, and the credentialed-box caveat demonstrably covers `internal/cli`. Beyond the ACs, both edited fixture classes were exercised against real `claude` and `codex` hosts rather than trusting skip-guarded suites.
