---
id: rrrhd7e79w41w1p39r0268e8
title: Unpin CI Claude Code version — run live-e2e on the merged-team floor (retire the #395 pin)
status: validation
source: 'Captain 2026-06-18 — 9243/#396 (using-claude-team merged-model support) merged + green on 2.1.181. The #395 keystone pinned live-e2e to 2.1.177 (last native-TeamCreate release) ONLY to keep the legacy team contract alive. With the merged contract shipped, the pin should be retired so CI runs the current (merged-team) Claude. Ships in 0.20.5 alongside m4''s merged lane.'
started: 2026-06-19T00:33:52Z
completed:
verdict:
score: 0.6
worktree: .worktrees/spacedock-ensign-ci-unpin-claude-version
issue:
---

Retire the 2.1.177 CI pin so the live-e2e lane runs current/unpinned Claude Code (the merged-team floor), now that #396 shipped the merged-team contract.

## Problem

`.github/workflows/runtime-live-e2e.yml` pins `SPACEDOCK_PINNED_CLAUDE_VERSION: "2.1.177"` (with `DISABLE_AUTOUPDATER: "1"`) because 2.1.177 was the last release exposing native `TeamCreate`/`TeamDelete`, and the legacy team contract required them (#395 keystone). #396 (`using-claude-team`) re-architected around the merged-team model (named background `Agent` + `SendMessage`, no `TeamCreate`), validated green on 2.1.181. The pin now holds CI on a deprecated team API while real users run the merged floor — CI no longer tests what ships. `internal/release/claude_version_pin_guard_test.go` actively ENFORCES the 2.1.177 pin (it REDs on any other pin / on no-pin install), so an unpin must update or retire that guard in the same change.

## Proposed approach

### Version policy (RESOLVED): float to installer-`latest`, freeze it for the run with `DISABLE_AUTOUPDATER`

The install step resolves the empty-input default to `latest` (the installer's own newest), and `DISABLE_AUTOUPDATER: "1"` stays — so each run installs current reality at install time, then the resolved binary is frozen for the multi-minute test run (no mid-run drift). The `workflow_dispatch` `claude_version` override survives as the escape hatch.

**Why float, not a `>=2.1.178` floor (or `stable`):**

1. **The defect #395 created is "CI tests a deprecated API, not what ships."** Users who `curl install.sh | bash` get installer-latest; the merged floor (#396) is what real sessions run. Float-to-latest makes CI test exactly the version users run — the fix's whole purpose.
2. **`stable` is a trap (rejected).** Live-fetched today: `latest` = 2.1.181, `stable` = 2.1.170. `stable` (2.1.170) is BELOW the merged floor AND still exposes native `TeamCreate` — pinning `stable` would re-create #395 inverted (CI back on a legacy-native version), the opposite of this task's goal.
3. **A `>=2.1.178` floor is a magic number that drifts and dulls the sensor (rejected).** A pinned floor lags reality, re-introduces the exact "mystery version a future reader rips out" problem the old guard's rationale comment fought, and — being frozen — would NOT have sensed the upstream regressions CI is meant to catch. The captain's recorded intent (m4 debrief: "CI stays on latest as the regression sensor that caught both this and the zero-discover flake") is float.
4. **The surprise-regression risk is bounded, not ignored.** A bad upstream release mid-release-window is the real cost of float. It is bounded by (a) `DISABLE_AUTOUPDATER` freezing the install-time-resolved version per run, so the "Show tool versions" step records the exact version that ran (per-run reproducibility — what failed is always known); and (b) the surviving `workflow_dispatch` `claude_version` input, which lets a maintainer pin a known-good version for a release cut if latest breaks.

**`DISABLE_AUTOUPDATER` coherence under float:** kept `"1"`. Its role shifts from "don't update past the pin" to "don't drift mid-run from the version the install step resolved." The install step resolves `latest` to a concrete version once; `DISABLE_AUTOUPDATER` then prevents that binary from auto-updating during the test run. Within-run reproducibility is preserved; cross-run currency is gained.

### The change, in three coupled edits

- **Workflow install step** — drop the `SPACEDOCK_PINNED_CLAUDE_VERSION` env var and the `${CLAUDE_VERSION:-$SPACEDOCK_PINNED_CLAUDE_VERSION}` fallback; resolve `latest` when the `workflow_dispatch` input is empty. **Drop the "Show tool versions" pin-enforcement block** (the `!= $SPACEDOCK_PINNED_CLAUDE_VERSION` RED) — there is no pin to enforce; the step still RECORDS the resolved version to the step summary. Keep `DISABLE_AUTOUPDATER: "1"` with a rewritten rationale comment (freeze-for-run, not hold-at-2.1.177). See the before/after diff below.
- **`internal/release/claude_version_pin_guard_test.go`** — its four tests all assert "the pin is exactly 2.1.177 with a team-tool rationale and an override-precedence fallback." Under float there is no pin. **Disposition: DELETE the file** (see the pin-guard disposition section); a new minimal `assertClaudeLiveWorkflowFloatsAndRecords` guard replaces the one still-relevant invariant.
- **`internal/contractlint/legacy_teamcreate_layering_test.go::TestLegacyRemovalTriggerIsExternallyAnchored`** — this guard the dispatch did NOT name but the unpin BREAKS: it asserts the workflow still pins 2.1.177, and is DESIGNED to RED when the pin moves ("the signal that the legacy branch has lost its last live consumer"). **Disposition: REWRITE** (see that disposition section) — the unpin satisfies only ONE of the legacy removal trigger's two conditions, so the guard must re-anchor without falsely declaring legacy dead.

### m4 coordination (AC-3)

The legacy interactive pty tests (`TestLivePtyStandingResidencyInjectsCommOfficer`, `TestLivePtyEnsignCycleTeamTeardown`, on m4's branch) require native `TeamCreate`, which the merged floor (≥2.1.178) does not expose. m4's implementation task owns adding a capability-aware skip (`ToolSearch(select:TeamCreate)` empty ⇒ no `TeamCreate` ⇒ `t.Skip`, mirroring AC-6's auth-skip; per `docs/roadmap/0205-layered-fo/m4-readiness.md` §"Next steps" step 2). This task's AC-3 is the *coherence* claim: on the unpinned host the legacy lane SKIPs (not REDs), green via m4's merged lane. The skip mechanism is m4's deliverable; this task does not author it, but the unpin must not land before it (sequencing recorded in the test plan).

## Out of scope

- Deleting the legacy `using-legacy-claude-team` path — the removal trigger has TWO conditions ("no runtime the FO targets still exposes `TeamCreate`" AND "no live lane drives the legacy branch"). Unpinning satisfies only the second; `stable` (2.1.170) still exposes `TeamCreate`, so a user on `stable` still hits the legacy path — it is NOT yet dead code. Legacy deletion retires later when STABLE Claude catches up to the merged floor (a separate trigger).
- Authoring m4's merged lane and m4's legacy-lane skip-gate (separate m4 work; this task changes the version + the two guards only).

## Concrete workflow diff (`.github/workflows/runtime-live-e2e.yml`)

### Edit 1 — the `claude-live` job `env:` block (drop the pin var, rewrite the `DISABLE_AUTOUPDATER` comment)

BEFORE (lines ~102-111):

```yaml
    env:
      DISABLE_AUTOUPDATER: "1"
      # Pinned Claude Code version. 2.1.177 is the LAST release that exposes the
      # native team tools (TeamCreate/TeamDelete): 2.1.178 dropped them from
      # headless `claude -p` (anthropics/claude-code#68721) and 2.1.179 dropped
      # them from interactive sessions too, even with
      # CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1 present (m4 live-team-mode finding,
      # NO_TEAM_TOOLS probe). Team-mode FO drives and the level-3-judge residency
      # need this version. DISABLE_AUTOUPDATER above is load-bearing for the pin:
      # without it claude auto-updates mid-job past 2.1.177 (m4 saw 2.1.177→2.1.179).
      SPACEDOCK_PINNED_CLAUDE_VERSION: "2.1.177"
      ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
```

AFTER:

```yaml
    env:
      # Freeze the installed version for the run's duration: the Install step below
      # resolves installer-`latest` ONCE, and this keeps that resolved binary from
      # auto-updating mid-run (a multi-minute test must not drift to a second
      # version partway through — the "Show tool versions" step records the exact
      # version that ran). It no longer holds a PIN: the live lane floats to current
      # Claude (the merged-team floor #396 ships on), so CI tests what users run
      # rather than the retired team-tool-capable 2.1.177 (#395, now retired).
      DISABLE_AUTOUPDATER: "1"
      ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
```

### Edit 2 — the "Install Claude Code" step (resolve `latest` on empty input)

BEFORE (lines ~137-147):

```yaml
      - name: Install Claude Code
        env:
          CLAUDE_VERSION: ${{ inputs.claude_version }}
        run: |
          # Empty workflow_dispatch input resolves to the team-tool-capable pin,
          # NOT to installer-latest (which is 2.1.178+, team tools removed). A
          # maintainer can still override with an explicit version or `latest`.
          VERSION="${CLAUDE_VERSION:-$SPACEDOCK_PINNED_CLAUDE_VERSION}"
          echo "Installing Claude Code version: $VERSION"
          curl -fsSL https://claude.ai/install.sh | bash -s -- "$VERSION"
          echo "$HOME/.local/bin" >> "$GITHUB_PATH"
```

AFTER:

```yaml
      - name: Install Claude Code
        env:
          CLAUDE_VERSION: ${{ inputs.claude_version }}
        run: |
          # Empty workflow_dispatch input floats to installer-`latest` — the live
          # lane tests the version users actually run (the merged-team floor #396
          # ships on), not a frozen pin. A maintainer can still override with an
          # explicit version or `stable` via the claude_version dispatch input
          # (the escape hatch to pin a known-good version for a release cut).
          VERSION="${CLAUDE_VERSION:-latest}"
          echo "Installing Claude Code version: $VERSION"
          curl -fsSL https://claude.ai/install.sh | bash -s -- "$VERSION"
          echo "$HOME/.local/bin" >> "$GITHUB_PATH"
```

### Edit 3 — the "Show tool versions" step (drop the pin-enforcement RED; keep recording)

BEFORE (lines ~172-191): the step computes `CLAUDE_VERSION_FIELD` and `exit 1`s if it `!= $SPACEDOCK_PINNED_CLAUDE_VERSION`, then records to the step summary.

AFTER:

```yaml
      - name: Show tool versions
        run: |
          # No pin to enforce: the live lane floats to current Claude. Record the
          # resolved version (the install step froze it for the run) so a failure
          # is always attributable to a known version, and surface it in the step
          # summary for the regression-sensor read.
          echo "claude version (frozen for this run): $(claude --version)"
          go version
          echo "### Tool versions" >> "$GITHUB_STEP_SUMMARY"
          echo "- \`claude --version\`: \`$(claude --version)\`" >> "$GITHUB_STEP_SUMMARY"
          echo "- \`go version\`: \`$(go version)\`" >> "$GITHUB_STEP_SUMMARY"
          echo "- Model: \`${{ matrix.model }}\`" >> "$GITHUB_STEP_SUMMARY"
          echo "- Effort: \`${{ inputs.effort }}\`" >> "$GITHUB_STEP_SUMMARY"
```

(The `workflow_dispatch` `claude_version` input description at lines ~26-28 already says "Empty = installer default" — that is now accurate; no edit needed there.)

## Pin-guard test disposition

### `internal/release/claude_version_pin_guard_test.go` — DELETE, replace with a minimal float guard

All four tests in this file (`TestClaudeLiveWorkflowPinsClaudeVersion`, `…RejectsNoPinInstall`, `…RejectsUndocumentedPin`, `…RejectsDroppedOverride`) and the shared predicate `assertClaudeLiveWorkflowPinsClaudeVersion` / `assertPinVarDocumented` assert the RETIRED policy: pin == 2.1.177, a team-tool rationale comment, and the `${CLAUDE_VERSION:-$SPACEDOCK_PINNED_CLAUDE_VERSION}` fallback. Under float every one of these is false by design — keeping any of them rewritten-to-2.1.178 would just re-encode a pin we decided not to have.

**Why delete, not rewrite in place:** three of the four tests guard a pin's *properties* (exact value, rationale comment, override-precedence over the pin) that no longer exist. The one invariant still worth a guard is narrower: the empty-input install path must float to `latest` and must NOT silently re-introduce a pin. Carrying the old file's scaffolding to assert that would be more retired-policy code than the new invariant needs.

**Replacement (new `claude_version_float_guard_test.go` in the same package):** one test `TestClaudeLiveWorkflowFloatsToLatest` + its adversarial twin, over the parsed workflow (reusing the existing `readWorkflow`/`parseWorkflowSteps`/`findStepByName`/`executableShellCommands` helpers):

- **POSITIVE:** the "Install Claude Code" step resolves `${CLAUDE_VERSION:-latest}` (empty input ⇒ `latest`), still reads `inputs.claude_version` into `CLAUDE_VERSION` (the override escape hatch survives), and the workflow declares NO `SPACEDOCK_PINNED_CLAUDE_VERSION` var anywhere (the pin is gone, not hidden).
- **ADVERSARIAL twin** (the real discriminator): re-introduce a `SPACEDOCK_PINNED_CLAUDE_VERSION: "2.1.x"` var + a `${CLAUDE_VERSION:-$SPACEDOCK_PINNED_CLAUDE_VERSION}` fallback into the fixture and assert the new guard REDs — so a future PR that silently re-pins the float lane is caught.

This guard's PROOF is `go test ./internal/release/...` (offline, deterministic), exercising both the positive and the mutated-fixture adversarial path — not a prose-grep.

### `internal/contractlint/legacy_teamcreate_layering_test.go::TestLegacyRemovalTriggerIsExternallyAnchored` — REWRITE the anchor

This test (the dispatch did not name it; the unpin breaks it) asserts the workflow contains `SPACEDOCK_PINNED_CLAUDE_VERSION: "2.1.177"` and is DESIGNED to RED when the pin moves — its comment calls that RED "the signal that the legacy branch has lost its last live consumer and the skill may be deleted." Unpinning IS that move, but the dispatch keeps legacy-deletion out of scope. The removal trigger (per `using-legacy-claude-team/SKILL.md`) has TWO conditions; unpinning satisfies only "no live lane drives the legacy branch." A user on `stable` (2.1.170, still `TeamCreate`-capable) still hits the legacy path, so it is NOT dead code — the guard must re-anchor without falsely declaring legacy dead.

**Disposition: REWRITE (not delete).** The renamed `TestLegacyConsumerRetiredButPathLives` asserts the now-true layering fact:
- the workflow NO LONGER pins a `TeamCreate`-capable version (`SPACEDOCK_PINNED_CLAUDE_VERSION` absent) — the live consumer is retired; AND
- the legacy skill `skills/using-legacy-claude-team/SKILL.md` AND its load line in `claude-fo-dispatch.md` (the `spacedock:using-legacy-claude-team` reference) STILL EXIST — the best-effort path is retained, not deleted.

It binds two independent sources (the live CI yaml and the on-disk skill/contract files), so it stays a real structural guard, not a self-referential prose-grep. The comment is updated to record that the externally-checkable proxy for legacy *deletion* is no longer the CI pin (now retired) but the future condition "even `stable` Claude is ≥ the merged floor" — a separate trigger. PROOF: `go test ./internal/contractlint/...` green.

(`TestNormalPathContractInlinesNoLegacyMachinery` and `TestNormalPathContractNamesLegacySkill` in the same file are UNAFFECTED — they assert the dispatch-contract layering, not the pin.)

## Acceptance criteria (proof = behavior, never prose-grep)

**AC-1 — the live-e2e lane installs and runs current/unpinned Claude, frozen for the run.**
End-state: the `claude-live` job resolves `${CLAUDE_VERSION:-latest}` (no `SPACEDOCK_PINNED_CLAUDE_VERSION` var), `DISABLE_AUTOUPDATER` stays set, and the "Show tool versions" step records the resolved version without a pin-enforcement RED.
Verified by: a green live-e2e CI run whose "Show tool versions" step summary reports a ≥2.1.178 (merged-floor) version, AND whose Claude shared-scenarios + merged-lane scenario green on that version (the merged path, no `TeamCreate`). The CI run is the proof — not a grep over the yaml.

**AC-2 — the version guard reflects the float policy and catches a silent re-pin.**
End-state: `internal/release/claude_version_pin_guard_test.go` is deleted; the new `claude_version_float_guard_test.go` asserts the float-resolution + no-pin invariant and REDs on a re-introduced pin; `internal/contractlint/legacy_teamcreate_layering_test.go` no longer asserts the 2.1.177 pin and instead asserts "consumer retired, path retained."
Verified by: `go test ./internal/release/... ./internal/contractlint/...` green under the new policy, AND the new release guard's adversarial twin demonstrably RED when a `SPACEDOCK_PINNED_CLAUDE_VERSION` + fallback is mutated back into the fixture (run the mutated case; observe the RED). Deterministic, offline.

**AC-3 — the legacy interactive lane skips (not REDs) on the unpinned host.**
End-state: on the floated live-e2e run the m4 pty team-mode tests SKIP (capability probe finds no `TeamCreate` ⇒ merged host), and the `claude-live` job is green via the merged-lane + shared scenarios.
Verified by: the unpinned live-e2e CI run shows the legacy pty tests SKIPPED (in the gotestsum step log / archived `*-detail.jsonl`) and the job green. Depends on m4's capability-skip landing first (test-plan sequencing); this task's diff does not author the skip.

## Test plan

- **AC-2 (offline, do first — the cheapest checkable proof):** apply the guard dispositions, run `go test ./internal/release/... ./internal/contractlint/...`. Confirm green, then mutate the float-guard fixture to re-introduce a pin and confirm the new adversarial twin REDs (the discriminator). Cost: minutes, deterministic, no credentials. This is the riskiest-mechanism-first check — if the guards can't be made to reflect float coherently, the workflow edit is moot.
- **AC-1 + AC-3 (live, one CI run):** apply the three workflow edits, trigger `Runtime Live E2E` (or land on a `main` PR), approve the `CI-E2E`/`CI-E2E-OPUS` environments. Confirm the "Show tool versions" summary reports ≥2.1.178, the merged-lane + shared scenarios green, and the legacy pty tests SKIP (not RED). Cost: one approval-gated live run (same cost class as the existing live lane).
- **Sequencing (load-bearing):** AC-3 depends on m4's capability-skip being on `main` (or co-landing) — otherwise the unpinned legacy pty tests RED instead of SKIP. The unpin PR must NOT merge ahead of m4's skip-gate. Recorded as a merge-order dependency, not a code dependency in this diff.
- **No spike needed:** the mechanisms this rests on are all proven. Installer `latest`/`stable`/explicit-version resolution is the documented `install.sh` contract (verified live: `latest`=2.1.181, `stable`=2.1.170). The merged floor running flag-free on ≥2.1.178 is proven by #396 (merged, green on 2.1.181). The capability-probe skip is m4's proven mechanism (its branch already drives `ToolSearch(select:TeamCreate)`). The guard helpers (`readWorkflow`/`parseWorkflowSteps`/`executableShellCommands`) are the existing, exercised release-test machinery. The only judgment call — float vs floor — is a policy decision (resolved above), not an unverified mechanism.

## Related

- `m4` live-team-mode-terminal-harness — ships its merged lane + the legacy-lane capability-skip in 0.20.5 alongside this unpin; the two are the 0.20.5 cut. AC-3 depends on m4's skip landing first.
- `docs/roadmap/0205-layered-fo/m4-readiness.md` — the FO package that scopes the both-semantics seam, names the legacy-lane skip-on-merged-host plan (§"Next steps" step 2), and frames this unpin (§step 3).
- #395 pin keystone (`ea03d094`) — what this retires.
- #396 `using-claude-team` merged-model support (`1d691b45`, merged) — the contract that makes the unpin safe; validated green on 2.1.181.
- `internal/contractlint/legacy_teamcreate_layering_test.go::TestLegacyRemovalTriggerIsExternallyAnchored` — the SECOND guard (not named in the dispatch) the unpin breaks; rewritten here, NOT deleted (legacy path stays best-effort).
- `skills/using-legacy-claude-team/SKILL.md` — its removal-trigger prose names the CI pin as the legacy live-consumer proxy; the rewrite keeps that coherent (consumer retired, path retained).
- `bare-mode-coverage` — the `-p`-assumes-bare audit that pin-retirement opens (sequenced after).

## Stage Report: ideation

- DONE: Version policy resolved with rationale — float-to-latest vs pin-a-merged-floor-minimum (>=2.1.178) — and the concrete runtime-live-e2e.yml before/after diff recorded in the task body.
  RESOLVED to float-to-installer-`latest` with `DISABLE_AUTOUPDATER` retained to freeze the resolved version for the run (4-point rationale in §"Version policy"; `stable`=2.1.170 rejected as a legacy-native trap, `>=2.1.178` floor rejected as a drifting magic number that dulls the regression sensor). Three coupled before/after yaml edits recorded in §"Concrete workflow diff" (drop the pin var + rewrite the `DISABLE_AUTOUPDATER` comment; resolve `${CLAUDE_VERSION:-latest}`; drop the "Show tool versions" pin-enforcement RED, keep recording).
- DONE: The internal/release/claude_version_pin_guard_test.go disposition decided (delete, or rewrite to assert the new floor/rationale) and recorded as a concrete before/after, so the guard reflects the shipped policy rather than the retired 2.1.177 pin.
  DECIDED: DELETE the file (all four tests + the two predicates assert pin==2.1.177 properties that no longer exist under float) and replace with a minimal `claude_version_float_guard_test.go` — a positive float-resolution/no-pin assertion plus an adversarial twin that REDs on a re-introduced pin (proof = `go test ./internal/release/...` over both paths, not a prose-grep). Recorded with the concrete before/after in §"Pin-guard test disposition". Also surfaced and dispositioned a SECOND guard the dispatch did not name — `internal/contractlint/legacy_teamcreate_layering_test.go::TestLegacyRemovalTriggerIsExternallyAnchored`, which binds the literal 2.1.177 pin and is DESIGNED to RED on a pin move — REWRITTEN (not deleted) to "consumer retired, path retained," coherent with the legacy removal trigger's two conditions (out-of-scope deletion preserved).

### Summary

Resolved the load-bearing float-vs-floor decision in favor of floating to installer-`latest` (with `DISABLE_AUTOUPDATER` reframed from hold-at-pin to freeze-for-run), because the merged floor #396 is what users run and the captain's recorded intent is to keep CI as the latest-version regression sensor; `stable`(2.1.170) and a `>=2.1.178` floor were both considered and rejected with rationale. Recorded three concrete before/after yaml edits and the guard dispositions: delete `claude_version_pin_guard_test.go` for a minimal float guard with an adversarial re-pin twin, and — a guard the dispatch did not name but the unpin breaks — rewrite `legacy_teamcreate_layering_test.go`'s removal-trigger assertion to "live consumer retired, best-effort legacy path retained" so it stays coherent with the two-condition legacy removal trigger without falsely declaring legacy dead. No spike needed (installer version resolution, the #396 merged floor, m4's capability-skip, and the release-test helpers are all proven mechanisms); AC-3 carries a load-bearing merge-order dependency on m4's legacy-lane skip landing first.

## Stage Report: implementation

- DONE: Apply the three coupled runtime-live-e2e.yml edits from the ideation body: drop the SPACEDOCK_PINNED_CLAUDE_VERSION env var + the ${CLAUDE_VERSION:-$SPACEDOCK_PINNED_CLAUDE_VERSION} fallback and resolve ${CLAUDE_VERSION:-latest}; rewrite the DISABLE_AUTOUPDATER comment (freeze-for-run, keep "1"); drop the "Show tool versions" pin-enforcement RED while keeping it RECORDING the resolved version.
  Commit 647b5ef4 on `spacedock-ensign/ci-unpin-claude-version`; `grep SPACEDOCK_PINNED_CLAUDE_VERSION` in the yaml returns nothing, resolution is `${CLAUDE_VERSION:-latest}`, `DISABLE_AUTOUPDATER: "1"` retained with the freeze-for-run comment, the Show-tool-versions step records `claude --version` to the summary with no `exit 1`. YAML re-parsed clean via `yq eval`.
- DONE: Apply both guard dispositions from the ideation body: DELETE internal/release/claude_version_pin_guard_test.go and add the new claude_version_float_guard_test.go (positive float-resolution/no-pin assertion + the adversarial twin that REDs on a re-introduced pin); REWRITE internal/contractlint/legacy_teamcreate_layering_test.go's removal-trigger assertion to "live consumer retired, best-effort legacy path retained" (binds the now-no-pin yaml AND the still-present using-legacy-claude-team skill + its load line).
  Same commit. `git rm` of the pin guard; new float guard carries the two helpers (`findStepByName`/`stepEnvBinds`) that were only used in the deleted file. `TestLegacyRemovalTriggerIsExternallyAnchored` → `TestLegacyConsumerRetiredButPathLives`, now binding three sources (no-pin yaml + the `spacedock:using-legacy-claude-team` load line in claude-fo-dispatch.md + the non-empty SKILL.md). Also fixed two now-false comments in that file that justified themselves by the retired 2.1.177 pin.
- DONE: Confirm offline green + the adversarial discriminator: `go test ./internal/release/... ./internal/contractlint/...` passes, AND the new float guard demonstrably REDs when a SPACEDOCK_PINNED_CLAUDE_VERSION + fallback is mutated back into the fixture (run the mutated case, observe the RED).
  `go test ./internal/release/... ./internal/contractlint/...` green; full `go test ./...` green. Discriminator proven TWO ways: the in-test adversarial twin passes, AND mutating the REAL fixture back to a pin made both `TestClaudeLiveWorkflowFloatsToLatest` and `TestLegacyConsumerRetiredButPathLives` RED with their discriminating messages (then reverted, re-confirmed green).

### Summary

Applied the three coupled `runtime-live-e2e.yml` edits (drop the pin var + fallback → `${CLAUDE_VERSION:-latest}`; reframe `DISABLE_AUTOUPDATER` to freeze-for-run; drop the Show-tool-versions enforcement RED, keep recording) and both guard dispositions (delete the pin guard, add `claude_version_float_guard_test.go` with the no-pin/float-resolution assertion + the re-pin adversarial twin; rewrite the contractlint removal-trigger guard to "consumer retired, path retained" binding the no-pin yaml + the still-present legacy skill and its dispatch load line). Carried the two helpers (`findStepByName`/`stepEnvBinds`) that lived only in the deleted file into the new float guard, and fixed two now-false comments in the contractlint file that referenced the retired pin. AC-2 fully proven offline (suites green; the discriminator RED demonstrated by mutating the real fixture and reverting). AC-1/AC-3 are the live CI proof, gated on m4's capability-skip landing first (recorded merge-order dependency) — not authored here. All on branch `spacedock-ensign/ci-unpin-claude-version`, commit 647b5ef4.
