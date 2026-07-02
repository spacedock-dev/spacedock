---
title: "Dispatch model space: fable joins the model enum; sonnet-5 recognized as a 1M-context model"
status: backlog
group: tooling
source: "Captain request 2026-07-02 (Claude Commander session), while probing per-stage model routing (fable ideation ensigns, sonnet-5 implementers): dispatch build validates declared models against {sonnet, opus, haiku} (internal/dispatch/build.go:59), so 'model: fable' on a stage errors; and the context-budget probe's window mapping (internal/claudeteam/contextbudget.go) knows the [1m] suffix and the claude-opus-4-{minor>=7} forward family but not sonnet-5, so a sonnet-5 member resolves to the 200k default."
id: wcex4yjx4mvecybxjb43gwtw
---

## Problem
Two gaps block captain-directed per-stage model routing:
1. `dispatch build` rejects `model: fable` declared on a stage or in workflow defaults — the enum at `internal/dispatch/build.go:59` (and its error string at build.go:23) is `{sonnet, opus, haiku}`, while the Claude Code Agent tool already accepts `fable`. Golden fixtures (`build-model-stage-wins-opus.txt`, `build-model-defaults-haiku.txt`) and the canonical-enum prose in `skills/first-officer/references/claude-fo-dispatch.md` (break-glass conditional model slot; context-budget canonical enum for reuse-condition-4) pin the old enum.
2. `spacedock dispatch context-budget` resolves a sonnet-5 member model to the 200k default window: `internal/claudeteam/contextbudget.go` grants 1M only to the `[1m]` suffix and the `claude-opus-4-{minor}` family with minor >= 7. Sonnet-5 is a 1M-context model and should resolve to `extendedContextLimit`, ideally via a forward family rule so later sonnet-5 releases stay correct without an edit (the same shape the opus family got).

## Desired direction (for ideation to refine)
- `model: fable` passes dispatch build validation at both precedence sites (stage and defaults) and is emitted as the artifact's `model`; an unknown model still refuses with the updated enum list in the error.
- The context-budget probe resolves sonnet-5 model ids to the 1M window (forward family rule), observable in the probe's JSON limit; the 200k default for pre-5 sonnets and the existing opus/[1m] rules unchanged.
- The contract prose naming the canonical model enum is updated with the shipped enum so reuse-condition-4's comparator and the break-glass template stay truthful.

## Rough acceptance sketch (ideation tightens into measured ACs + a test plan)
- A README declaring `model: fable` on a stage builds a dispatch artifact whose model field is fable (RED today: build exits non-zero with the enum error).
- `context-budget` on a sonnet-5 member reports the 1M limit (RED today: 200k), with pre-5 sonnet and opus behavior pinned unchanged.
- Golden fixtures and enum error strings updated; existing model-precedence tests stay green.
