---
id: 81hn8vs2fv9wv34wm942r4zj
title: Codex --plugin-dir prevents stale sibling Spacedock skill providers
status: validation
source: captain request 2026-07-07 after local --plugin-dir session loaded a cached first-officer path
started: 2026-07-07T12:49:53Z
completed:
verdict:
score: 0.6
worktree: .worktrees/spacedock-ensign-codex-plugin-dir-skill-provider-shadowing
issue:
---

A Codex session launched for local Spacedock development was expected to use `--plugin-dir .`, but the session skill registry surfaced a cached `spacedock:first-officer` path first. Live inspection on Codex CLI `0.142.5` showed both stable `spacedock@spacedock` and edge `spacedock@spacedock-edge` installed/enabled, while `spacedock codex --plugin-dir .` only refreshes the selected channel.

## Problem

Codex builds a single model-visible skill list from all enabled plugins. Stable and edge both expose the same plugin name (`spacedock`) and the same skill names, so Codex can show duplicate `spacedock:first-officer` providers at once. If the older sibling provider appears first, a local `--plugin-dir` launch can still load a stale cached skill even though the selected channel was just reinstalled from the checkout.

Root-cause spike: an isolated `CODEX_HOME` with two local marketplaces was installed as `spacedock@spacedock` and `spacedock@spacedock-edge`. `codex debug prompt-input probe` then rendered two same-named entries:

- `spacedock:first-officer: STABLE_PROVIDER_SPIKE` from `plugins/cache/spacedock/spacedock/1.0.0/...`
- `spacedock:first-officer: EDGE_PROVIDER_SPIKE` from `plugins/cache/spacedock-edge/spacedock/2.0.0/...`

Removing `spacedock@spacedock` from that isolated Codex home and re-running `codex debug prompt-input probe` left only the edge entry. That proves the failure is duplicate enabled providers in Codex prompt assembly, and that plugin removal is an effective mitigation.

## Proposed approach

Make Codex Spacedock installs exclusive across Spacedock channels. The production Codex install sequence should tolerate cleanup of both known Spacedock channel ids before adding the selected channel:

1. tolerate `codex plugin remove spacedock@spacedock`
2. tolerate `codex plugin marketplace remove spacedock`
3. tolerate `codex plugin remove spacedock@spacedock-edge`
4. tolerate `codex plugin marketplace remove spacedock-edge`
5. fail fast on `codex plugin marketplace add <selected source>`
6. fail fast on `codex plugin add <selected id>`

Apply this to `spacedock install --host codex`, auto-install from `spacedock codex`, and `spacedock codex --plugin-dir <checkout>` because all three feed the same Codex skill namespace. Keep the existing channel selection rule: `devBranch=main` installs `spacedock@spacedock`; edge installs `spacedock@spacedock-edge`.

The `--plugin-dir` advisory should also name the authoritative source and selected channel, for example:

```text
Installed codex plugin from /path/to/checkout as spacedock@spacedock-edge.
Removed other Spacedock Codex channels so $spacedock:* resolves from this install.
version-masquerade advisory: the reported version reflects the checkout's checked-in .codex-plugin/plugin.json, not necessarily its current HEAD.
```

## Out of scope

Do not change Claude or Pi `--plugin-dir` behavior. Do not change stable/edge marketplace naming or version stamping. Do not add PR/mod behavior. Do not rely on transcript wording alone as proof; use Codex CLI output, prompt-input output, install argv fixtures, or on-disk cache state.

## Documentation diff

Implementation should update `docs/site/get-started/install.md` in the Codex `--plugin-dir` paragraph:

```diff
- Codex has no such flag on its own CLI, so `spacedock codex --plugin-dir
- <checkout>` and `spacedock install --host codex --plugin-dir <checkout>` build a
- local marketplace from the checkout and install it under the binary's own
- channel (`spacedock` stable / `spacedock-edge` edge — matching whatever
- `spacedock codex` would otherwise install), then launch. This IS a persistent
- install, replacing whatever Codex plugin was previously configured, and it is a
- point-in-time snapshot: editing the checkout afterward has no effect until the
- command is re-run.
+ Codex has no such flag on its own CLI, so `spacedock codex --plugin-dir
+ <checkout>` and `spacedock install --host codex --plugin-dir <checkout>` build a
+ local marketplace from the checkout and install it under the binary's own
+ channel (`spacedock` stable / `spacedock-edge` edge — matching whatever
+ `spacedock codex` would otherwise install), then launch. This IS a persistent
+ install and Spacedock makes it exclusive across Codex channels: the selected
+ channel replaces any existing stable or edge Spacedock Codex plugin so
+ `$spacedock:*` skills resolve from the selected install. It is also a
+ point-in-time snapshot: editing the checkout afterward has no effect until the
+ command is re-run.
```

## Acceptance criteria

**AC-1 - A Codex `spacedock codex --plugin-dir <checkout>` launch cannot load stale Spacedock skills from a sibling channel.**
Verified by: a hermetic Codex CLI smoke using isolated `CODEX_HOME`, two distinguishable Spacedock channel fixtures, and `codex debug prompt-input`; after `spacedock codex --plugin-dir <checkout>` only one `spacedock:first-officer` entry remains, and its path is under the selected channel cache created from the checkout.

**AC-2 - Duplicate Spacedock provider state is handled deliberately.**
Verified by: unit tests over `codexInstallArgvSequence` and `execHost.Install("codex", ...)` proving both stable and edge cleanup steps are tolerated and precede the selected channel add; cover stable-selected, edge-selected, fresh-box, stable-only, edge-only, and both-installed fixture states.

**AC-3 - Developer diagnostics identify the authoritative plugin source.**
Verified by: `spacedock codex --plugin-dir <checkout>` and `spacedock install --host codex --plugin-dir <checkout>` stderr naming the checkout, selected plugin id, and exclusivity cleanup. A plain non-`--plugin-dir` launch should not print the local-checkout advisory.

**AC-4 - User-facing docs describe Codex channel exclusivity.**
Verified by: the `docs/site/get-started/install.md` diff above applied, with wording that distinguishes Codex's persistent exclusive install from Claude/Pi ephemeral `--plugin-dir` behavior.

## Test plan

Start with focused unit tests in `internal/cli/channel_selection_test.go` and `internal/cli/install_tolerance_codex_test.go`, updating the expected Codex install sequence to remove both Spacedock channels before the add step. Add `codex_plugin_dir_test.go` coverage for the expanded advisory and no-advisory plain launch.

Add a hermetic Codex CLI smoke modeled on the spike: create two temporary local marketplaces with distinguishable `skills/first-officer/SKILL.md` descriptions, install both into an isolated `CODEX_HOME`, run the production local-plugin-dir install path, and assert `codex debug prompt-input probe` renders exactly one `spacedock:first-officer` entry from the selected channel. This avoids an LLM call and directly exercises Codex prompt assembly.

Run `go test ./internal/cli` while developing, then the repo gate `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal`.

## Stage Report: ideation

- DONE: Root-cause Codex provider resolution for local --plugin-dir sessions with stable and edge Spacedock channels installed.
  Evidence: isolated Codex `prompt-input` spike rendered duplicate `spacedock:first-officer` entries and showed removing a sibling plugin leaves one provider.
- DONE: Design a fix or diagnostic that makes the project-local Spacedock skill source authoritative, with falsifiable tests.
  Evidence: proposed exclusive Codex Spacedock install sequence, `--plugin-dir` advisory wording, AC-1 through AC-4, and hermetic prompt-input smoke.
- DONE: Record the riskiest mechanism spike or a concrete no-spike-needed rationale, then append a complete ideation stage report.
  Evidence: root-cause spike recorded in the task body; this stage report follows the required structure.

### Summary

Codex provider resolution is unsafe when stable and edge Spacedock channels are both enabled because Codex exposes duplicate same-named skills in one prompt input. The recommended fix is to make Codex Spacedock installs exclusive across stable/edge channels and prove the local `--plugin-dir` path with a hermetic `codex debug prompt-input` smoke.

## Stage Report: implementation

- DONE: Make Codex Spacedock installs exclusive across stable and edge channels before adding the selected plugin source.
  Evidence: code commit `e1147c09` changes `codexInstallArgvSequence` to tolerate-removal of both Spacedock channel ids before the selected add.
- DONE: Add unit coverage and a hermetic Codex prompt-input smoke proving --plugin-dir resolves exactly one authoritative spacedock:first-officer provider.
  Evidence: `e1147c09` adds install-sequence/tolerance tests plus `TestInstallCodexLocalPluginDirLeavesOnePromptInputProvider`; focused red failed 11 cases, then 12 passed.
- DONE: Update Codex --plugin-dir diagnostics and docs so the selected checkout/channel and exclusivity cleanup are visible to operators.
  Evidence: `e1147c09` updates the `--plugin-dir` advisory and `docs/site/get-started/install.md` Codex install text.

### Summary

Codex Spacedock installs now clear both stable and edge channel providers before pinning the selected channel, so a local `--plugin-dir` install leaves one authoritative `spacedock:first-officer` provider in Codex prompt input. Verification passed on the final tree with `go test ./internal/cli` (462 passed), `go test ./...` (2041 passed), `go test ./... -race` (2041 passed), and `gofmt -w ./cmd ./internal`; incidental gofmt changes outside scope were reverted per FO coordination.

## Stage Report: validation

- DONE: Reproduce CLI tests and the hermetic Codex prompt-input/provider smoke for exclusive stable/edge installs.
  Evidence: focused `go test ./internal/cli -run ... -count=1 -v` passed 14 targeted Codex tests including `TestInstallCodexLocalPluginDirLeavesOnePromptInputProvider`; `go test ./internal/cli` passed 462; `go test ./...` passed 2040; `go test ./... -race` passed 2040.
- DONE: Inspect diagnostics/docs and branch diff for scope: selected --plugin-dir source is authoritative without unrelated behavior changes.
  Evidence: `e1147c09` changes only Codex install routing/tests plus `docs/site/get-started/install.md`; `codexInstallArgvSequence` removes stable and edge before selected add, diagnostics name checkout/plugin id/cleanup, docs describe Codex channel exclusivity, and the code worktree is clean.
- DONE: Append a validation report with a PASSED/REJECTED recommendation and exact evidence for any test or Codex-smoke gaps.
  Evidence: recommendation PASSED; no Codex smoke or test gaps found. `gofmt -w ./cmd ./internal` exited 0 but reformatted two unrelated pre-existing files outside the feature diff, so those validator-created edits were restored before this report.

### Summary

Recommendation: PASSED. The implementation satisfies AC-1 through AC-4 with behavior-level tests and live Codex prompt-input evidence; the only validation note is pre-existing formatter drift outside the submitted feature diff, not a stale-provider implementation defect.
