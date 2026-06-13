---
id: 7d72hqzg9tpnr6tdjfa34nxt
title: dispatch context-budget — suppress spurious model warnings on healthy team members
status: implementation
source: github#344 (captain intake 2026-06-13)
started: 2026-06-13T17:31:54Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-context-budget-spurious-warnings
issue: "#344"
---

`spacedock dispatch context-budget` emits two warnings that read as faults but are environmental noise on a healthy reused team member, eroding trust in the reuse-condition-0 budget signal. Intook to the 0.20.3 (0203-fo-efficiency) sprint as FO dispatch-path quality.

## Problem

(github#344) On a healthy reused team member, `spacedock dispatch context-budget --name {member}` emits:
- `config_drift_warning` — team config records the captain session model string (e.g. `claude-fable-5[1m]`) while the runtime jsonl reports the canonical id (`claude-fable-5`); these never match, so it fires every probe.
- `mixed_models_warning` — jsonl `<synthetic>` (harness-injected) entries mix with the real model → "multiple models seen — using smallest context window"; with `<synthetic>` unrecognized the fallback window is opaque.

Expected: `[1m]`-suffixed config values normalize before comparison; `<synthetic>` entries excluded from the model census. Version: binary 0.20.0 (contract 1), plugin 0.19.9.

## Root cause (reproduced empirically)

A throwaway in-package test drove the real `ContextBudget` (the function `spacedock dispatch context-budget` calls) over two fixtures and captured the actual JSON:

**config_drift on a healthy member** — config `claude-fable-5[1m]`, runtime jsonl `claude-fable-5`:
```
"model": "claude-fable-5", "context_limit": 200000,
"config_declared_model": "claude-fable-5[1m]",
"config_drift_warning": "team config requested claude-fable-5[1m] but runtime is claude-fable-5"
```
The drift check (`internal/claudeteam/contextbudget.go:144`) is `driftSeen := len(runtimeModels) > 0 && configLimit != contextLimit`. `contextLimitForModel("claude-fable-5[1m]")` is 1M (the `[1m]` substring rule, line 40); `contextLimitForModel("claude-fable-5")` is 200k (non-opus, no suffix). The limits differ → drift fires every probe. **Second harm:** the chosen window comes from the runtime model only (line 140), so a captain session genuinely running the 1M window is measured against a 200k denominator — the wrong denominator, the same class of error as the archived opus-4-8 false-negative.

**mixed_models on a healthy member** — config `claude-fable-5`, jsonl carrying a real `claude-fable-5` entry plus a harness-injected `<synthetic>` entry:
```
"model": "<synthetic>", "context_limit": 200000,
"mixed_models_warning": "multiple models seen in jsonl: ['<synthetic>', 'claude-fable-5'] — using smallest context window"
```
`extractRuntimeModels` (line 301) collects every distinct `message.model`, including `<synthetic>`. With two entries the `default` census branch (line 121) fires the mixed warning and picks the smallest window; `<synthetic>` is unrecognized → 200k, and it sorts first (`<` < `c`) so it becomes the emitted `model`. The chosen model is a placeholder.

## Approach

Two localized fixes in `internal/claudeteam/contextbudget.go`. No vendored Python oracle exists — context-budget parity is native-vs-frozen-golden under `internal/dispatch/testdata/golden/`, regenerated with `go test ./internal/dispatch -run TestContextBudget... -update`. The riskiest-mechanism record below establishes the design needs no spike.

**(a) Exclude `<synthetic>` from the model census.** In `extractRuntimeModels` (the single source of `runtimeModels`), skip assistant entries whose `message.model == "<synthetic>"` — add a guard alongside the existing `entry.Message.Model == ""` skip at line 317. Smallest possible point of the fix: every downstream consumer (the `len(runtimeModels)` switch, the mixed warning, the chosen model, the window) reads from this census, so excluding the placeholder there fixes the warning, the model, and the window in one place.

Before (line ~317):
```go
if entry.Type != "assistant" || entry.Message == nil || entry.Message.Model == "" {
    continue
}
```
After:
```go
if entry.Type != "assistant" || entry.Message == nil || entry.Message.Model == "" {
    continue
}
if entry.Message.Model == "<synthetic>" {
    continue
}
```

**(b) Normalize the `[1m]` suffix before the config-vs-runtime comparison.** When the single runtime model equals the config model modulo a trailing `[1m]` suffix the runtime dropped (e.g. config `claude-fable-5[1m]`, runtime `claude-fable-5`), the config model is authoritative — it carries the suffix that names the window the captain session genuinely runs. Promote the chosen `model` to `configModel` in the single-model branch when they are the same base differing only by the suffix.

Add a helper:
```go
// sameModelModuloExtended reports whether a and b name the same base model
// differing only by an [1m] suffix (the captain-session model string carries
// the suffix the spawned runtime drops). Distinguishes a healthy suffix-dropped
// member from a genuinely different runtime model.
func sameModelModuloExtended(a, b string) bool {
    return a != b && stripExtended(a) == stripExtended(b)
}

// stripExtended drops a trailing [..] suffix (e.g. [1m]) from a model string.
func stripExtended(m string) string {
    if i := strings.Index(m, "["); i >= 0 {
        return m[:i]
    }
    return m
}
```
In `case 1:` (line 119-120), after `model = runtimeModels[0]`:
```go
case 1:
    model = runtimeModels[0]
    if sameModelModuloExtended(model, configModel) {
        model = configModel
    }
```
After promotion `model == configModel`, so `contextLimit == configLimit` → the existing drift check (line 144) no longer fires, the emitted `model` is `claude-fable-5[1m]`, and the window is 1M (correct denominator). The drift line and the `default` (mixed) branch are unchanged; genuine cross-family drift (config `claude-opus-4-8` / runtime `claude-sonnet-4-6`) has different bases, so no promotion, drift still fires.

Note `contextLimitForModel` already strips `[..]` via `strings.Index(base, "[")` at line 44, so its behavior on the promoted `claude-fable-5[1m]` is the explicit 1M path (line 40), consistent with `stripExtended`.

## Acceptance criteria

- **AC-1 — a healthy member with config `model[1m]` + runtime canonical `model` emits no `config_drift_warning` and reads the `[1m]` (1M) window.** End state: `ContextBudget` over a fixture with config `claude-fable-5[1m]` and a runtime jsonl entry `claude-fable-5` produces JSON with no `config_drift_warning` key, `context_limit` 1000000, and `model` `claude-fable-5[1m]`.
  - Verified by: a Go test in `internal/claudeteam/contextbudget_test.go` driving real `ContextBudget` over a `t.TempDir` fixture, parsing stdout JSON, asserting the `config_drift_warning` key is absent, `context_limit == 1000000`, and `reuse_ok` correct for the resident count. (Fails today: drift fires + window is 200k.)

- **AC-2 — a healthy member whose jsonl carries `<synthetic>` alongside the real model emits no `mixed_models_warning`, picks the real model, and reads the real model's window.** End state: `ContextBudget` over a fixture with one `claude-fable-5` entry and one `<synthetic>` entry produces JSON with no `mixed_models_warning` key, `model` `claude-fable-5`, and `context_limit` 200000.
  - Verified by: a Go test driving real `ContextBudget`, asserting `mixed_models_warning` absent, `model == "claude-fable-5"`, `context_limit == 200000`. (Fails today: mixed warning fires, `model` is `<synthetic>`.)

- **AC-3 (over-suppression guard) — drift normalization does not silence a genuinely drifted member.** End state: a member whose config and runtime are genuinely different models (config `claude-opus-4-8`, runtime `claude-sonnet-4-6`, distinct bases) still emits `config_drift_warning`.
  - Verified by: a Go test asserting `config_drift_warning` is present for the cross-family fixture. (This is the existing `ensign-mm` shape; the new test asserts the warning survives the fix.)

- **AC-4 (over-suppression guard) — synthetic exclusion does not silence a genuinely mixed-model member.** End state: a member whose jsonl carries two genuinely different REAL models (`claude-opus-4-8` + `claude-sonnet-4-6`, no `<synthetic>`) still emits `mixed_models_warning` and picks the smallest window.
  - Verified by: a Go test asserting `mixed_models_warning` present and `context_limit == 200000` for the two-real-model fixture. The existing `TestContextBudgetWarningNonASCIIParity/mixed-models-warning` parity case (two real models) already covers this and must stay green.

- **AC-5 — the existing context-budget golden parity stays green.** End state: the eight `internal/dispatch/testdata/golden/context-budget-*.txt` goldens still match native output; only the goldens whose fixtures the fix legitimately changes are regenerated, and `ensign-mm` (two real models) is NOT among them.
  - Verified by: `go test ./internal/dispatch -run TestContextBudget` passing without `-update` after any necessary golden regeneration; the diff to `testdata/golden/` reviewed at the gate. The current fixtures use no `[1m]`-config and no `<synthetic>`, so the existing goldens are expected to be byte-unchanged — a regeneration that touches them flags an unintended behavior change.

## Test plan

- **Location / kind:** Go unit tests in `internal/claudeteam/contextbudget_test.go` driving the real `ContextBudget` over `t.TempDir` `~/.claude` fixtures (the `writeBudgetFixture` / `writeTeamConfigWithSession` helpers already in that file), parsing the emitted JSON and asserting key presence/absence + scalar fields. These are the AC-1..AC-4 oracles — they read the probe's ACTUAL output, never an instruction file.
- **New fixtures:** (1) config `model[1m]` + runtime canonical `model` (AC-1); (2) real model + `<synthetic>` entry (AC-2); (3) cross-family config/runtime (AC-3); (4) two-real-model jsonl (AC-4). The current `writeBudgetFixture` writes a single entry — add a small multi-line-jsonl variant for the `<synthetic>`/two-model cases (one assistant line per model, mirroring `assistantLine` in the dispatch parity test).
- **Parity goldens (AC-5):** run `go test ./internal/dispatch -run TestContextBudget` first to confirm the existing goldens stay green unchanged; regenerate with `-update` only if a golden's fixture legitimately exercises the new paths (none of the current eight do, so expect no golden churn).
- **Cost/complexity:** low — pure Go unit tests over in-memory `t.TempDir` fixtures, no live workflow, no network. ~4 new test functions plus the multi-entry fixture helper.

## Riskiest-mechanism on record

No spike needed. The probe is a deterministic Go function over on-disk JSON; both fixes are pure string operations on already-parsed model strings, composing already-proven behavior (`contextLimitForModel`, the census map, the drift comparison). The one unknown — does the `[1m]` normalization distinguish a suffix-dropped healthy member from genuine drift without over-suppressing — was exercised with a throwaway `sameModelModuloExtended` probe over five model pairs: it promotes suffix-only differences (`claude-fable-5` vs `claude-fable-5[1m]`, `opus` vs `opus[1m]`, `claude-opus-4-6` vs `claude-opus-4-6[1m]`) and rejects genuine drift (`claude-sonnet-4-6` vs `claude-opus-4-8`) and identical strings. Proven mechanisms relied on: `contextLimitForModel`'s existing `[..]`-strip + `[1m]`→1M rule, the census map, the limit-based drift comparison.

## Documentation

No doc diff. The `config_drift_warning` / `mixed_models_warning` strings are FO-internal probe output; the docs site (`docs/site/`, `docs/runtime-support.md`) does not describe them (grep: zero matches for `config_drift` / `mixed_models` / `context-budget` outside roadmap and `.spacedock-state`). The FO runtime skill (`skills/first-officer/references/claude-first-officer-runtime.md:136`) documents the model-to-context family rule, which this change does not alter (it composes the existing rule). The healthy-member output simply stops carrying the two spurious keys.

## Stage Report: ideation

- DONE: Flesh the approach with specific before/after: locate where `spacedock dispatch context-budget` emits `config_drift_warning` and `mixed_models_warning`, then specify (a) normalizing the `[1m]`-suffixed captain-session config model string before the config-vs-jsonl model comparison, and (b) excluding `<synthetic>` jsonl entries from the model census.
  Approach section gives before/after for both: synthetic exclusion in `extractRuntimeModels` (line ~317) and a `sameModelModuloExtended`/`stripExtended` promotion in the `case 1` branch (line 119) that fixes both the spurious drift warning and the wrong 200k denominator in `internal/claudeteam/contextbudget.go`.
- DONE: Acceptance criteria each proven by a Go test over the probe's ACTUAL output (warnings present/absent + reuse_ok + chosen context window) with fixtures — a healthy member yields zero spurious warnings and the correct window; never a string-match over an instruction file.
  AC-1..AC-5: each names a Go test driving real `ContextBudget` over a `t.TempDir` fixture, parsing emitted JSON, asserting warning-key absence/presence + `context_limit` + `model` + `reuse_ok`; AC-5 binds the parity goldens.
- DONE: Riskiest-mechanism on record: the fix must NOT over-suppress ... an over-suppression regression-guard test. If the probe is a deterministic Go function, record "no spike needed" naming the proven mechanism.
  Over-suppression guards are AC-3 (genuine cross-family drift still warns) and AC-4 (two real models still warn); "no spike needed" recorded with the `sameModelModuloExtended` five-pair probe result and the proven mechanisms (`contextLimitForModel`, census map, limit-based drift).

### Summary

Reproduced both bugs empirically by driving the real `ContextBudget` over fixtures: config_drift fires because the drift check compares context limits and `claude-fable-5[1m]` (1M) never matches runtime `claude-fable-5` (200k) — which also picks the wrong 200k denominator; mixed_models fires because `<synthetic>` pollutes the census and even becomes the emitted `model`. Fix (a) excludes `<synthetic>` in `extractRuntimeModels`; fix (b) adds a `[1m]`-suffix-normalizing promotion in the single-runtime-model branch so the config's declared window is authoritative. Both are pure deterministic Go string ops over already-parsed models — no spike needed — with two over-suppression guards and the native-vs-frozen-golden parity (no Python oracle) recorded.

## Stage Report: implementation

- DONE: Implement both fixes in internal/claudeteam/contextbudget.go per the gate-approved design: (a) skip `<synthetic>` entries in extractRuntimeModels (beside the existing empty-model skip); (b) a stripExtended/sameModelModuloExtended promotion in the case-1 branch making the config's `[1m]`-declared window authoritative when the runtime is the same base modulo the suffix — fixing both the spurious warning AND the wrong-window denominator.
  Commit 46224f5f: +22 lines in contextbudget.go — `<synthetic>` skip after the empty-model skip in extractRuntimeModels (line ~339), and `sameModelModuloExtended`/`stripExtended` helpers + the case-1 promotion (`if sameModelModuloExtended(model, configModel) { model = configModel }`). Matches the entity body's before/after verbatim.
- DONE: Land all 5 ACs as Go tests over the probe's real JSON output (AC-1/AC-2 healthy cases; AC-3/AC-4 over-suppression guards; AC-5 golden parity — regenerate only goldens the fix legitimately changes, and ensign-mm MUST stay byte-unchanged).
  AC-1..AC-4 added to internal/claudeteam/contextbudget_test.go driving real `ContextBudget` over `t.TempDir` fixtures, parsing stdout JSON, asserting warning-key absence/presence + model + context_limit + reuse_ok (4/4 pass; AC-1/AC-2 confirmed red before the fix, AC-3/AC-4 held pre-fix). AC-5: `go test ./internal/dispatch -run TestContextBudget -update` produced ZERO golden churn (`git status` on testdata/golden empty after `-update`), so ensign-mm and all eight goldens are byte-unchanged.
- DONE: Whole package green (`go test ./internal/claudeteam/ ./internal/dispatch/`), build clean; commit the deliverable on the worktree branch.
  `go test ./internal/claudeteam/ ./internal/dispatch/` → 213 passed (209 baseline + 4 new). `go build ./...` Success; `go vet` both packages clean. Committed as 46224f5f on branch spacedock-ensign/context-budget-spurious-warnings.

### Summary

Implemented the two gate-approved fixes exactly as specified — no re-design. (a) `extractRuntimeModels` now skips `<synthetic>` placeholders, removing them from the census so the mixed warning, chosen model, and window all read from real models only. (b) A `sameModelModuloExtended` promotion in the single-runtime-model branch makes the config's `[1m]`-suffixed model authoritative when the runtime is the same base sans suffix, clearing the spurious drift warning and the wrong 200k denominator (now correctly 1M). Five ACs land as Go tests over real probe JSON; both over-suppression guards confirmed; golden parity byte-unchanged including ensign-mm.
