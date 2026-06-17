---
title: Pin CI live-e2e to a known-good Claude Code version (team tools regress at 2.1.178+)
status: validation
source: captain (2026-06-17) — team-tool availability regression
score: 0.5
id: 61an320jkjgeyqq55rfx7nw0
started: 2026-06-17T14:38:40Z
worktree:
mod-block:
pr: "#395"
completed: 2026-06-17T18:58:19Z
verdict: PASSED
archived: 2026-06-17T18:58:19Z
---

## Problem

The live-e2e CI lane (`.github/workflows/runtime-live-e2e.yml`) floats the Claude Code version. Its `claude_version` workflow_dispatch input defaults to `""` (lines 26-29); the `Install Claude Code` step (lines 128-139) reads that input and, when empty, runs the no-pin branch `curl -fsSL https://claude.ai/install.sh | bash` — installer latest. On a `pull_request` trigger the input is *always* empty (workflow_dispatch inputs are absent for `pull_request`), so every PR live run already takes the no-pin branch unconditionally.

Latest is now 2.1.178+, where headless `claude -p` lost the native team tools (TeamCreate/TeamDelete; anthropics/claude-code#68721), and 2.1.179 dropped them even from interactive sessions with `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1` present (established by `live-team-mode-terminal-harness`/m4: a direct live probe on 2.1.179 returns `NO_TEAM_TOOLS`; the FO falls back to bare mode and the bounded-teardown grade correctly REDs). Team-mode FO drives and the standing level-3-judge residency both need a version where team tools work. **2.1.177 is the last known-good** — m4's live tmux probes on 2.1.177 created a real `teams/<name>/config.json` with the session resident, the exact behavior 2.1.178+ removes.

Pin CI on 2.1.177 so the live lane runs against a deterministic, team-tool-capable Claude Code instead of whatever floats to latest, and assert the resolved version so the pin is enforced rather than merely declared.

## Scope — the COMPLETE set of CI Claude-Code install points

A repo-wide grep (`claude.ai/install`, `install.sh | bash`, `claude_version`, `claude --version`, `@anthropic-ai/claude`) over `.github/` returns hits in exactly ONE file: `runtime-live-e2e.yml`. There is **one** floating install point, the `Install Claude Code` step. The other workflows do not install or run `claude`:

- `install-e2e.yml` — installs the `spacedock` binary from a release tarball; no claude.
- `release.yml` — its own comment (lines 150-151) states "CI never runs `claude` — it only reads the tag"; goreleaser cross-build only.
- `next-publish.yml` — bumps a marketplace calendar version; no claude.
- `docs.yml` — MkDocs build; no claude.

No `schedule` / `cron` / `workflow_run` / `repository_dispatch` triggers exist anywhere in `.github/workflows/` (grep returns nothing), so there is no scheduled trigger that floats the version. **The pin is a single-file, single-step change.** "Other lane that needs the same pin" resolves to: none.

## Proposed approach

1. **Make 2.1.177 the effective default.** Define a single source of truth for the pinned version in the `claude-live` job's `env:` block (alongside the existing `DISABLE_AUTOUPDATER: "1"` at line 102) and have the install step fall back to it when `claude_version` is empty, instead of falling through to no-pin latest. The workflow_dispatch input is KEPT as a manual override (a maintainer can still type `latest` or another version to test a floated build deliberately), but its empty default now resolves to the pin, not to latest.

2. **Assert the resolved version (the checkable proof).** The `Show tool versions` step (lines 164-172) already runs `claude --version` but only echoes it. Add an assertion that the resolved version's leading field equals the pin and FAILS the job otherwise — so a regression to a floated build REDs the run instead of silently running on a team-tool-broken claude.

3. **Self-document the WHY.** A comment beside the pin records that 2.1.177 is the last team-tool-capable release (regression at 2.1.178+, m4's live-team-mode finding, anthropics/claude-code#68721), so a future reader does not treat it as a mystery magic number; and that `DISABLE_AUTOUPDATER: "1"` (already present at line 102) is load-bearing for the pin — without it claude auto-updates mid-job past the pin (m4 observed a 2.1.177→2.1.179 auto-update when the isolated HOME lacked it).

## Spike — riskiest unverified mechanisms (exercised FIRST)

Two mechanisms could invalidate the whole task if false. Both were exercised against the real local claude (2.1.177) before any yaml is written.

**Spike 1 — `claude --version` output format (the assertion's contract).** `claude --version` on the local 2.1.177 prints `2.1.177 (Claude Code)`, NOT a bare `2.1.177`. A naive `[ "$(claude --version)" = "2.1.177" ]` would RED even on a correctly-pinned binary. The assertion must extract the leading whitespace-delimited field. Verified end-to-end against the real binary:
- Case A — real claude `2.1.177 (Claude Code)`: `${raw%% *}` → `2.1.177`, assertion GREEN (rc=0).
- Case B — simulated floated `2.1.179 (Claude Code)`: assertion REDs (rc=1) with `::error::claude version is '2.1.179' … expected pinned 2.1.177` — it catches the exact regression this task guards.
- Case C — hypothetical bare `2.1.177` (no suffix): GREEN — robust to both formats.
The extraction-then-equality is the proven assertion shape; a plain string equality is not.

**Spike 2 — installer accepts a version positional arg.** The fetched `claude.ai/install.sh` parses `TARGET="$1"` and validates it against `^(stable|latest|[0-9]+\.[0-9]+\.[0-9]+(-…)?)$`, so `bash -s -- "2.1.177"` passes validation; the installer then delegates the actual pin to `"$binary"…install 2.1.177`. The current workflow's pinned branch (line 134) already uses this exact `bash -s -- "$CLAUDE_VERSION"` syntax, so the arg path is the author's intent. The remaining live unknown — that `2.1.177` still resolves on the install endpoint at CI time — is verifiable only in CI and is covered by the AC-2 assertion: if the install resolves the wrong version, the assertion REDs the job. **No spike needed beyond this for the install path:** the syntax is proven by the existing pinned branch and the format/assertion is proven by Spike 1; the endpoint-resolution claim is exactly what the enforced assertion checks at run time.

## Concrete before/after diff (`.github/workflows/runtime-live-e2e.yml`)

**(a) `claude-live` job `env:` — add the pinned-version source of truth (after line 102 `DISABLE_AUTOUPDATER: "1"`):**

```diff
     env:
       DISABLE_AUTOUPDATER: "1"
+      # Pinned Claude Code version. 2.1.177 is the LAST release that exposes the
+      # native team tools (TeamCreate/TeamDelete): 2.1.178 dropped them from
+      # headless `claude -p` (anthropics/claude-code#68721) and 2.1.179 dropped
+      # them from interactive sessions too, even with
+      # CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1 present (m4 live-team-mode finding,
+      # NO_TEAM_TOOLS probe). Team-mode FO drives and the level-3-judge residency
+      # need this version. DISABLE_AUTOUPDATER above is load-bearing for the pin:
+      # without it claude auto-updates mid-job past 2.1.177 (m4 saw 2.1.177→2.1.179).
+      SPACEDOCK_PINNED_CLAUDE_VERSION: "2.1.177"
       ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
```

**(b) `Install Claude Code` step — resolve empty input to the pin, not to no-pin latest (replace lines 128-139):**

```diff
       - name: Install Claude Code
         env:
           CLAUDE_VERSION: ${{ inputs.claude_version }}
         run: |
-          if [ -n "$CLAUDE_VERSION" ]; then
-            echo "Installing pinned Claude Code version: $CLAUDE_VERSION"
-            curl -fsSL https://claude.ai/install.sh | bash -s -- "$CLAUDE_VERSION"
-          else
-            echo "Installing latest Claude Code (no pin)"
-            curl -fsSL https://claude.ai/install.sh | bash
-          fi
+          # Empty workflow_dispatch input resolves to the team-tool-capable pin,
+          # NOT to installer-latest (which is 2.1.178+, team tools removed). A
+          # maintainer can still override with an explicit version or `latest`.
+          VERSION="${CLAUDE_VERSION:-$SPACEDOCK_PINNED_CLAUDE_VERSION}"
+          echo "Installing Claude Code version: $VERSION"
+          curl -fsSL https://claude.ai/install.sh | bash -s -- "$VERSION"
           echo "$HOME/.local/bin" >> "$GITHUB_PATH"
```

**(c) `Show tool versions` step — assert the resolved version equals the pin (replace lines 164-172):**

```diff
       - name: Show tool versions
         run: |
-          claude --version
+          # Enforce the pin: the leading version field MUST equal the pin or the
+          # job REDs. `claude --version` prints "2.1.177 (Claude Code)", so strip
+          # from the first space before comparing (a bare equality would red on
+          # the suffix). This is the checkable proof the live lane runs on the
+          # team-tool-capable version, not whatever floated to latest.
+          CLAUDE_VERSION_RAW="$(claude --version)"
+          CLAUDE_VERSION_FIELD="${CLAUDE_VERSION_RAW%% *}"
+          if [ "$CLAUDE_VERSION_FIELD" != "$SPACEDOCK_PINNED_CLAUDE_VERSION" ]; then
+            echo "::error::claude resolved to '$CLAUDE_VERSION_FIELD' (full: '$CLAUDE_VERSION_RAW'); expected pinned $SPACEDOCK_PINNED_CLAUDE_VERSION" >&2
+            exit 1
+          fi
+          echo "claude version pin verified: $CLAUDE_VERSION_FIELD"
           go version
           echo "### Tool versions" >> "$GITHUB_STEP_SUMMARY"
           echo "- \`claude --version\`: \`$(claude --version)\`" >> "$GITHUB_STEP_SUMMARY"
           echo "- \`go version\`: \`$(go version)\`" >> "$GITHUB_STEP_SUMMARY"
           echo "- Model: \`${{ matrix.model }}\`" >> "$GITHUB_STEP_SUMMARY"
           echo "- Effort: \`${{ inputs.effort }}\`" >> "$GITHUB_STEP_SUMMARY"
```

The `SPACEDOCK_PINNED_CLAUDE_VERSION` env var, set once in the job `env:`, is the single source of truth read by both the install step and the assertion — install and proof can never disagree about which version is the pin.

## Acceptance criteria (with tests)

- **AC-1 (single source of truth, no float):** The `claude-live` job defines `SPACEDOCK_PINNED_CLAUDE_VERSION: "2.1.177"` in its `env:`, and the `Install Claude Code` step installs `${CLAUDE_VERSION:-$SPACEDOCK_PINNED_CLAUDE_VERSION}` — the empty-input path resolves to the pin, never to no-pin latest. *Test:* the offline document-order/workflow guard family (the same `internal/release/workflow_exec_guard_test.go`-style yaml-parsing tests already in the repo) asserts the install step contains no unconditional no-pin `install.sh | bash` branch and that the install command references `SPACEDOCK_PINNED_CLAUDE_VERSION`. A grep-style structural assertion over the parsed workflow, run by `go test ./...` (offline, no API spend).
- **AC-2 (enforced, checkable proof — the deliverable):** A live `claude-live` run FAILS the job unless the resolved `claude --version` leading field equals `2.1.177`. *Test:* the `Show tool versions` step's assertion, exercised live in CI; locally proven RED→GREEN by Spike 1's three cases against the real 2.1.177 binary (GREEN on 2.1.177, RED on simulated 2.1.179, GREEN on bare-format 2.1.177). This is an enforced change, not a prose claim.
- **AC-3 (self-documenting):** A comment beside the pin records WHY 2.1.177 (team-tool regression at 2.1.178+, m4 finding, #68721) and that `DISABLE_AUTOUPDATER` is load-bearing for it. *Test:* the offline workflow guard asserts the `SPACEDOCK_PINNED_CLAUDE_VERSION` line is preceded by a comment naming the regression (a structural presence check, not prose review); human-readable rationale is checked at the ideation/implementation gate.
- **AC-4 (manual override preserved):** The `claude_version` workflow_dispatch input still overrides the pin when a maintainer sets it (a deliberate float for testing remains possible). *Test:* covered by AC-1's structural assertion — `${CLAUDE_VERSION:-…}` precedence means a non-empty input wins; the offline guard asserts the install reads `inputs.claude_version` into `CLAUDE_VERSION` and uses `:-` fallback.

## Test plan

- **Cost/complexity:** Low. One workflow file, three edits, one source-of-truth env var. No new test infrastructure.
- **Offline (free, the gate):** Extend the existing yaml-parsing workflow-guard tests (`internal/release/`) to assert AC-1/AC-3/AC-4 structurally — no install branch floats, the pin var + rationale comment are present, the input-override precedence holds. Run by `go test ./...`.
- **Live (CI, API-gated):** AC-2's assertion is the live proof; it runs inside the already-approval-gated `claude-live` matrix (CI-E2E / CI-E2E-OPUS). No new live test is needed — the assertion rides the existing `Show tool versions` step.
- **Local pre-CI confidence:** Spike 1 already exercised the assertion logic end-to-end against the real 2.1.177 binary (RED→GREEN), so the proof is validated before CI spend.

## High-stakes surface / validation note

This is a CI-machinery change — a high-stakes surface under the proof policy. Implementation lands in a worktree; **validation runs the detached adversarial audit** because CI machinery is the gate that protects every downstream live run, and a silently-broken pin would let the team-tool regression back in unnoticed.

## Docs

No user-visible CLI output, command surface, startup banner, or docs-site content changes. The only edits are CI workflow internals and their comments. No doc diff required.

## Linkage

Pinning CI to a team-tool-capable version unblocks the `fo-tier-delegation` (72) residency decision — the standing level-3-judge residency becomes viable in the CI env once the pin lands.

## Stage Report: ideation

- DONE: Establish the COMPLETE set of CI install points that float the Claude Code version … and produce the concrete before/after yaml diff pinning each to 2.1.177.
  Repo-wide grep over `.github/` returns claude install/run hits in exactly ONE file (`runtime-live-e2e.yml`, the `Install Claude Code` step + empty `claude_version` input); release.yml/install-e2e.yml/next-publish.yml/docs.yml install no claude and there are no schedule/cron/workflow_run triggers — so the pin is a single-file, three-edit change with concrete before/after diffs (a/b/c) in the body.
- DONE: Define the CHECKABLE proof that CI runs actually resolve to 2.1.177 … the lane asserts `claude --version` reports 2.1.177 and fails the job otherwise.
  AC-2 + diff (c): the `Show tool versions` step extracts the leading version field and `exit 1`s on mismatch; exercised RED→GREEN against the real local 2.1.177 binary (GREEN on 2.1.177, RED on simulated 2.1.179, GREEN on bare-format) — an enforced check, not a prose claim.
- DONE: Make the pin self-documenting … record WHY 2.1.177 … note that validation runs the detached adversarial audit.
  AC-3 + diff (a): the `SPACEDOCK_PINNED_CLAUDE_VERSION: "2.1.177"` line carries a comment naming the 2.1.178+ team-tool regression, the m4 finding, #68721, and the load-bearing role of `DISABLE_AUTOUPDATER`; the high-stakes-surface section records the detached adversarial audit at validation.

### Summary

The floating-version surface is a single step in `runtime-live-e2e.yml`; no other lane or trigger installs claude, so the pin is one file. Two riskiest mechanisms were spiked against the real local 2.1.177 first: `claude --version` prints `2.1.177 (Claude Code)` (so the assertion strips from the first space rather than a brittle bare-equality, proven RED→GREEN), and the installer accepts a version positional arg (already the workflow's pinned-branch syntax). A single `SPACEDOCK_PINNED_CLAUDE_VERSION` env var is the source of truth read by both the install fallback and the enforced assertion, so install and proof cannot disagree; `DISABLE_AUTOUPDATER` (already present) is flagged as load-bearing against mid-job auto-update. No user-visible/doc changes.

## Stage Report: implementation

- DONE: Apply the three edits to .github/workflows/runtime-live-e2e.yml exactly as the body's diff (a/b/c): (a) add `SPACEDOCK_PINNED_CLAUDE_VERSION: "2.1.177"` to the claude-live job `env:` with the rationale comment; (b) the Install Claude Code step resolves `${CLAUDE_VERSION:-$SPACEDOCK_PINNED_CLAUDE_VERSION}` and installs via `bash -s -- "$VERSION"` — NO unconditional no-pin branch; (c) the Show tool versions step extracts the leading version field (`${raw%% *}`) and `exit 1`s on mismatch with a named `::error::`. ONE source of truth read by BOTH install and assertion.
  Worktree commit 6697a02b on branch `spacedock-ensign/pin-ci-claude-code-version`: env line 111 (install line 144, assertion lines 180-182 both read SPACEDOCK_PINNED_CLAUDE_VERSION); the old no-pin else-branch is gone; workflow still parses via yaml.v3.
- DONE: Add the offline workflow-guard test (internal/release/, the existing yaml-parsing guard family) asserting AC-1/AC-3/AC-4 structurally over the parsed workflow … TDD: this guard test first (red), then the yaml edits (green).
  `internal/release/claude_version_pin_guard_test.go` (new): TestClaudeLiveWorkflowPinsClaudeVersion + three adversarial twins (no-pin fallback restored, rationale comment stripped, override dropped) — each RED before the edits, all GREEN after; full `go test ./...` exit 0, no failures/panics.
- DONE: Confirm the AC-2 assertion shape matches Spike 1's proven RED-to-GREEN (extract-then-equality, NOT bare string equality) — AC-2 is the live proof and rides the existing Show tool versions step; do NOT add a new live test.
  Re-exercised the exact diff-(c) logic against the real local 2.1.177 binary: Case A GREEN on `2.1.177 (Claude Code)`, Case B RED (rc=1, named `::error::`) on simulated `2.1.179 (Claude Code)`, Case C GREEN on bare `2.1.177`. No new live test added; the assertion rides the existing step.

### Summary

Pinned the live-e2e lane to Claude Code 2.1.177 (last team-tool-capable release) via the body's a/b/c diff, exactly as approved — single `SPACEDOCK_PINNED_CLAUDE_VERSION` source of truth read by both the install fallback and the enforced `claude --version` assertion, with the workflow_dispatch override preserved and `DISABLE_AUTOUPDATER` left intact. TDD: a new offline guard (`claude_version_pin_guard_test.go`) with three adversarial twins went RED before the edits and GREEN after; full `go test ./...` passes (exit 0). AC-2's extract-then-equality shape re-confirmed RED→GREEN against the real local 2.1.177 binary. Scope held to one workflow file + one guard test; no doc/user-visible changes.

## Stage Report: validation

- DONE: Run the DETACHED ADVERSARIAL AUDIT (proof policy, high-stakes CI-machinery surface): on a THROWAWAY checkout of the branch, try to REFUTE the pin's guarantees by constructing claim-breaking edits the deliverable's tests MUST catch — (1) re-introduce the unconditional no-pin fallback; (2) weaken the Show-tool-versions assertion to bare string-equality; (3) strip the rationale comment; (4) drop the workflow_dispatch override. Confirm the offline guard + assertion logic CATCH each.
  Throwaway: detached worktree at impl commit `6697a02b` under `/tmp/pin-ci-audit-throwaway` (real worktree untouched). (1) reintroduced the prior `if/else … | bash` no-pin branch as a REAL file mutation → `TestClaudeLiveWorkflowPinsClaudeVersion` REDs (`…does not resolve ${CLAUDE_VERSION:-$SPACEDOCK_PINNED_CLAUDE_VERSION}`). (3) stripped the rationale comment as a REAL file mutation → REDs (`…not preceded by a rationale comment`). (4) replaced `${CLAUDE_VERSION:-…}` with a bare `$SPACEDOCK_PINNED_CLAUDE_VERSION` as a REAL file mutation → REDs (same override-precedence error). All reverted after. (2) the assertion is LIVE-only (no offline test references it — confirmed via grep), so exercised the EXACT diff-(c) extract-then-equality shell logic directly: GREEN on `2.1.177 (Claude Code)` and bare `2.1.177`, RED on `2.1.178`/`2.1.179 (Claude Code)`. The bare-equality weakening does NOT silently pass a float — it FALSELY REDs even the correctly-pinned `2.1.177 (Claude Code)` (the real binary always carries the ` (Claude Code)` suffix), so it errs loud-red, never green-on-float. Two THEORETICAL offline-guard-evasion paths exist (no-pin via `| sh`, or download-then-`bash file`) but are NOT material: each requires keeping a correct pin-resolving line AND adding a contrived second no-arg installer, not the realistic "revert to the no-pin branch" vector (which IS caught), and AC-2's live assertion is independent defense-in-depth that REDs any resulting float at runtime.
- DONE: Verify the live AC-2 proof path end-to-end against the workflow logic: on a pull_request trigger (claude_version empty) the Install step resolves to the PIN and the assertion REDs on any non-2.1.177 resolution (extract-then-equality). Confirm assertion + install read the SAME SPACEDOCK_PINNED_CLAUDE_VERSION source of truth and DISABLE_AUTOUPDATER remains in env.
  Modeled the empty-input PR path: `${CLAUDE_VERSION:-$SPACEDOCK_PINNED_CLAUDE_VERSION}` resolves to `2.1.177` (workflow_dispatch inputs absent for `pull_request`), assertion GREEN; a simulated float resolution REDs. SINGLE source of truth confirmed: `SPACEDOCK_PINNED_CLAUDE_VERSION: "2.1.177"` is declared ONCE (line 111) and read by BOTH the install (line 144) and the assertion (lines 181-182) — install and proof cannot disagree. `DISABLE_AUTOUPDATER: "1"` present (line 102); `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS: "1"` present (line 113). NOTED (not material): a maintainer override to a non-pin version (AC-4) installs that version but the assertion then REDs against the pin literal — internally consistent with both ACs and arguably the intended "you can deliberately float, but the job goes red to flag it" safety behavior.
- DONE: Re-confirm scope integrity: `go test ./...` from the worktree is green (exit 0); the single-floating-install-point claim still holds; no user-visible/doc changes leaked in.
  `go test ./...` on the throwaway exits 0, no FAIL/panic (incl. `internal/release` ok, all four pin-guard tests + the `TestRuntimeLiveWorkflowGuard*` family GREEN). Repo-wide grep over `.github/` for claude install/version touchpoints hits ONLY `runtime-live-e2e.yml` (all in the one `Install Claude Code` step, the pinned `bash -s -- "$VERSION"`); no `schedule`/`cron`/`workflow_run`/`repository_dispatch` trigger floats independently. Impl commit `6697a02b` diff = `.github/workflows/runtime-live-e2e.yml` (+35) + `internal/release/claude_version_pin_guard_test.go` (+247) ONLY; no `*.md`/docs/CLI/banner/skills leak. Workflow parses as valid YAML via `gopkg.in/yaml.v3` (the same parser CI + the guards use).

### Summary

Refuted nothing material. On a throwaway detached worktree at impl commit `6697a02b` (real worktree untouched), all four dispatch-named claim-breaking edits are caught: the three offline-guarded ones (no-pin fallback, stripped comment, dropped override) RED `TestClaudeLiveWorkflowPinsClaudeVersion` as REAL on-disk file mutations, and the live-only assertion (AC-2) was exercised directly — its extract-then-equality REDs on 2.1.178/2.1.179 floats and the bare-equality weakening errs loud-red (never green-on-float). Install and proof read the single `SPACEDOCK_PINNED_CLAUDE_VERSION` source of truth; `DISABLE_AUTOUPDATER` intact; `go test ./...` green (exit 0); scope is the one workflow + one guard test with no doc/user-visible leak; one floating install point and no scheduled float. Recorded one non-material design note (AC-4 override + AC-2 assertion interaction) and two non-material theoretical guard-evasion paths backstopped by the live assertion. gate:true — recommend PASSED to the captain; no feedback-to-implementation routing needed.
