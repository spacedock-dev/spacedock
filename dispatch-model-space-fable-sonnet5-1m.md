---
title: "Dispatch model space: fable joins the model enum; sonnet-5 recognized as a 1M-context model"
status: ideation
group: tooling
source: "Captain request 2026-07-02 (Claude Commander session), while probing per-stage model routing (fable ideation ensigns, sonnet-5 implementers): dispatch build validates declared models against {sonnet, opus, haiku} (internal/dispatch/build.go:59), so 'model: fable' on a stage errors; and the context-budget probe's window mapping (internal/claudeteam/contextbudget.go) knows the [1m] suffix and the claude-opus-4-{minor>=7} forward family but not sonnet-5, so a sonnet-5 member resolves to the 200k default."
id: wcex4yjx4mvecybxjb43gwtw
started: 2026-07-02T01:45:25Z
---

## Problem

Two gaps block captain-directed per-stage model routing:
1. `dispatch build` rejects `model: fable` declared on a stage or in workflow defaults — the enum at `internal/dispatch/build.go:59` (error string at build.go:23) is `{sonnet, opus, haiku}`, while the Claude Code Agent tool already accepts `fable`. A second copy of the same enum lives at `internal/dispatch/standing.go:15` (`spawnModelEnum`, error string standing.go:17), so a standing mod declaring `model: fable` in its `## Hook: startup` errors the same way. Golden fixtures (`build-error-invalid-stage-model.txt`, `build-error-invalid-defaults-model.txt`, `spawn-standing-bad-model-enum.txt`, `spawn-standing-missing-model.txt`) and the canonical-enum prose in `skills/first-officer/references/claude-fo-dispatch.md` pin the old enum. The README schema (`docs/schema/workflow-readme.mdschema.yml`) declares the enum twice (defaults.model and states[].model).
2. `spacedock dispatch context-budget` resolves a sonnet-5 member model to the 200k default window: `internal/claudeteam/contextbudget.go` grants 1M only to the `[1m]` suffix and the `claude-opus-4-{minor}` family with minor >= 7. Sonnet-5 is a 1M-context model. The same gap hits fable: a `claude-fable-5` member also resolves to 200k, and fable-5's window is 1M by default — directly relevant since fable ideation ensigns are the captain's stated use case (see evidence below; scope addition to flag at the gate).

Captain constraint (2026-07-02, verbatim intent): the per-stage model declaration needs some design so it can be OVERLAYED for different providers/hosts — a workflow README declaring claude-space models must not break codex or pi sessions running the same workflow — but do not overcomplicate it.

## Verified evidence (riskiest unknown: the model-id string shapes)

Inspected real team configs (`~/.claude/teams/*/config.json`) and session/subagent transcripts (`~/.claude/projects/**/*.jsonl`) on the development machine, 2026-07-02:

- **Runtime jsonl assistant `message.model` values** for the Claude-5 family: `claude-fable-5` (83,546 entries), `claude-sonnet-5` (1,524 entries), `claude-fable-5[1m]` (81 entries). **No date suffix and no minor-version token** — the 5-family drops the `claude-sonnet-4-5-20250929` dated shape and the `claude-sonnet-4-6` major-minor shape. No `claude-sonnet-5[1m]` observed.
- **Team-config `members[].model` values** mix full ids (`claude-fable-5[1m]`, `claude-opus-4-8[1m]`, `claude-sonnet-4-6`) and Agent-schema short aliases (`sonnet`, `opus[1m]`, `fable` ×4, `inherit`). So `Agent(model=fable)` has already succeeded live on this machine, and this session's Agent tool schema declares the enum `["sonnet", "opus", "haiku", "fable"]`.
- **Window evidence:** a bare-id `claude-fable-5` transcript (no `[1m]`) reached **561,372 resident tokens** — past the 200k default, proving the harness runs fable-5 on an extended window without the suffix. Max observed `claude-sonnet-5` resident is 176,910 (locally inconclusive); the platform model catalog (claude-api reference, cached 2026-06-24) lists sonnet-5 context = 1M and fable-5 context = 1M ("the maximum is also the default"). Haiku 4.5 remains 200K.
- **Forward-family-rule decision (regex shape):** one generation rule `^claude-(sonnet|fable|opus)-(\d+)` with major >= 5 → 1M, layered alongside the existing `^claude-opus-4-(\d+)` minor >= 7 rule and the `[1m]` suffix override. Guards checked: `claude-sonnet-4-6` captures major 4 → 200k (pinned); dated old shape `claude-sonnet-4-5-20250929` captures 4 → 200k; hypothetical dated/minor 5-family shapes (`claude-sonnet-5-20260301`, `claude-sonnet-5-1`, `claude-sonnet-6`) capture >= 5 → 1M, so later releases stay correct without an edit. **Haiku is deliberately excluded**: haiku-4-5 is 200K, a future haiku-5's window is unverified, and under-granting is the safe failure mode (the probe under-reports the limit and refuses reuse early; it never overflows a member).

**No code spike needed**: the mechanisms relied on are all proven — the Agent schema's fable acceptance is observed live (schema + 4 stamped members), the model-id shapes come from real config/transcript data above, and the change sites (enum map, regex table, golden fixtures) are existing, tested code paths.

## Proposed approach

### 1. Enum: fable joins, at both sites

- `internal/dispatch/build.go:59` `modelEnum` gains `"fable"`; `modelEnumList` (build.go:23) becomes `"must be one of: sonnet, opus, haiku, fable"`.
- `internal/dispatch/standing.go:15` `spawnModelEnum` + `spawnModelEnumList` get the same change — it is the same Agent-schema enum; implementation should deduplicate the two definitions into one shared package-level enum/list so the next model lands with one edit.
- `docs/schema/workflow-readme.mdschema.yml` adds `"fable"` to both `model` enums (defaults and states[]). (Grep confirms nothing in `internal/` consumes this schema — `field_conformance.go` loads only `entity.mdschema.yml` — so this is a doc-truthfulness edit.)

### 2. Host overlay: the model declaration's value space is host-scoped (the "do not overcomplicate" shape)

Today `runBuildFields` validates the declared model against the claude enum unconditionally, before/independent of the resolved `host`. The chosen shape: **keep the single flat `model:` key; interpret its value space per resolved host at validation time.**

- **host=claude**: validate against the enum (now incl. fable); emit the model into the envelope; unknown model still errors with the updated enum list. Unchanged behavior apart from fable.
- **host=codex / host=pi**: any declared model is **ignored-with-note** — build exits 0, emits `model: null`, and prints one stderr note: `[build] declared model '{m}' ignored on host {host}: outside {host}'s dispatch-settable model space; emitting model=null`. It replaces the `[build] effective_model=…` line for that build.

**Why ignore-with-note, not error:** portability is the point of the constraint — the same README must build on every host; an error would force per-host README forks. `model: null` already has defined semantics on every host (claude: harness default; codex: the worker runs the thread's model — `«worker-identity»` says codex's model space IS the thread's model, so there is nothing dispatch-settable to honor; pi: the `pi-dispatch-model-stamping` adapter stamps the parent's live model on null). The stderr note keeps the drop observable to the FO. Consequence stated plainly: on codex/pi a model typo (`model: frobnicate`) is also ignored, not caught — it is caught where it binds, on a claude-host build.

**Rejected alternatives (recorded per captain's non-binding list):** per-host keys (`model@codex:`) add declaration syntax and parser surface no workflow needs yet; abstract tiers require a tier→model mapping per host — the most machinery. Both remain compatible later extensions because the flat key is untouched. When pi grows a positive per-stage model space (owned by `pi-dispatch-model-stamping`, per fo-dispatch-core.md:93), the pi branch changes from "ignore all" to "pass through pi-space values"; claude-enum names stay not-honorable on pi and keep the ignore-with-note path.

Standing-mod spawn (`spawn-standing`) renders Claude Agent envelopes only, so `spawnModelEnum` stays claude-scoped with no host branch.

### 3. Context window: the 5-generation family rule

`internal/claudeteam/contextbudget.go` `contextLimitForModel` gains the generation rule from the evidence section: `^claude-(sonnet|fable|opus)-(\d+)` with major >= 5 → `extendedContextLimit`. Existing rules unchanged: `[1m]` suffix → 1M; `claude-opus-4-{minor>=7}` → 1M; everything else (pre-5 sonnet, opus 4-6 bare, haiku, short aliases, unknown) → 200k default.

Known advisory consequence, pinned in the test plan: a member whose team config stamps the short alias (`fable`, `sonnet`) while the runtime jsonl reports the full 5-family id now trips `config_drift_warning` (config resolves 200k, runtime 1M). Reuse decisions are unaffected — `reuse_ok` computes from the runtime model's limit — and the warning is truthful, so it is accepted rather than special-cased with alias mapping (aliases are harness-version-ambiguous; `sonnet` meant a 200k model on older harnesses).

## Acceptance criteria

Each criterion states how it is tested.

- **AC-1** — **A README declaring `model: fable` on a stage builds a dispatch artifact whose JSON `model` field is `"fable"`; same for `defaults.model: fable`.** RED today: exit 1, `invalid model … must be one of: sonnet, opus, haiku` (build.go:379/:384). Tested: new cases in `TestBuildModelPrecedence` (internal/dispatch/build_hazards_test.go) with golden fixtures, asserting the emitted `"model"` JSON value and the `[build] effective_model=fable` stderr line.
- **AC-2 (VALUE)** — **`context-budget` on a member whose runtime model is `claude-sonnet-5` reports `"context_limit": 1000000`, and the reuse verdict flips with it**: a 250k-resident sonnet-5 member reads usage_pct 125.0 / reuse_ok false today, 25.0 / reuse_ok true after. RED today: 200000. Tested: in-package `ContextBudget` fixture test (the `writeBudgetFixtureModels` harness in contextbudget_test.go) asserting `context_limit`, `usage_pct`, `reuse_ok` in the emitted JSON.
- **AC-3 (VALUE)** — **A `claude-fable-5` member likewise resolves the 1M window.** RED today: 200k — live evidence shows a bare fable-5 transcript at 561,372 resident, which under the current rule reads usage_pct 280.7 / reuse_ok permanently false, so the FO can never reuse a healthy fable ensign. Tested: same fixture harness. (Scope addition beyond the seed's title; flag at the ideation gate.)
- **AC-4** — **Pre-5 and non-family behavior is pinned unchanged**: `claude-sonnet-4-6` → 200k, `claude-sonnet-4-5-20250929` → 200k, `claude-opus-4-6` → 200k, `claude-opus-4-7`/`4-8` → 1M, any `[1m]` suffix → 1M, `claude-haiku-4-5` and hypothetical `claude-haiku-5` → 200k, short aliases and unknown models → 200k. Tested: `TestContextLimitForModelBoundary` table extended with the new rows; existing rows stay green.
- **AC-5** — **`dispatch build --host codex` (and `--host pi`) over a fable-declaring README exits 0, emits `model: null`, and prints the ignore-note on stderr.** RED today: codex-host build over that README exits 1 with the enum error. Tested: new cases beside `build_codex_host_test.go` / `build_pi_host_test.go`, asserting exit code, the JSON `model` null literal, and the stderr note.
- **AC-6** — **An unknown model on a claude-host build still refuses, with the updated enum list in the error.** Tested: existing error goldens (`build-error-invalid-stage-model.txt`, `build-error-invalid-defaults-model.txt`) regenerated to the new `… sonnet, opus, haiku, fable` string; the tests that drive them stay red-on-regression.
- **AC-7** — **A standing mod declaring `model: fable` spawns; a bad model still errors with the updated list.** Tested: spawn-standing tests + regenerated `spawn-standing-bad-model-enum.txt` / `spawn-standing-missing-model.txt` goldens.
- **AC-8** — **Contract prose names the shipped enum and the shipped window rules** (concrete diffs below). Counts only as paired with ACs 1–7 — proof is anchored in the binary's tests; prose is updated to match, never grepped as evidence.

AC-2/AC-3 are the value-measuring ACs: the entity exists for captain-directed per-stage model routing, and its end value is measured as the built artifact's model field (AC-1, a behavior) plus the probe's numeric limit and the reuse_ok flip against the 200k baseline that can move the wrong way (AC-2/3).

## Test plan

- **Go unit tests** (internal/claudeteam): extend `TestContextLimitForModelBoundary` with the AC-4 rows plus `claude-sonnet-5`, `claude-sonnet-5[1m]`, `claude-sonnet-5-20260301`, `claude-sonnet-6`, `claude-fable-5`, `claude-opus-5`; one `ContextBudget` envelope test per AC-2/AC-3 shape, including the alias-config drift case (config `fable`, runtime `claude-fable-5` → limit 1M, `config_drift_warning` present) so the advisory consequence is pinned deliberately.
- **Golden fixtures** (internal/dispatch/testdata/golden): new `build-model-stage-fable` (and defaults variant) goldens; regenerate the four enum-error goldens; new codex/pi host-ignore goldens. Existing model-precedence goldens (`build-model-stage-wins-opus`, `build-model-defaults-haiku`, `build-model-null`) must stay byte-identical.
- **Command-level tests** (internal/dispatch): `TestBuildModelPrecedence` fable cases; host-ignore cases driving `runBuild` with `--host codex` / `--host pi` over a fable-declaring README fixture.
- No live workflow smoke needed: every claim is parser/command behavior the fixture harness already exercises; the one runtime-shaped claim (model-id shapes) is settled by the recorded evidence, not by a live run.
- Estimated cost: small — one focused implementation pass; all test infrastructure (golden harness, budget fixture writers, host test files) already exists.

## Documentation / contract-prose diff (applied by implementation, reviewed at the gate)

`skills/first-officer/references/claude-fo-dispatch.md:52` — break-glass conditional model slot:
- Before: `The canonical enum the conditional slot draws from is \`sonnet | opus | haiku\`.`
- After: `The canonical enum the conditional slot draws from is \`sonnet | opus | haiku | fable\`.`

`skills/first-officer/references/claude-fo-dispatch.md:140` — context-budget paragraph, two sentences:
- Before: `The opus context window follows a forward family rule (\`claude-opus-4-{minor}\` with minor ≥ 7 → 1M; the \`[1m]\` suffix → 1M; else 200k), so a new opus release stays correct without an edit.`
- After: `The context window follows forward family rules (\`claude-opus-4-{minor}\` with minor ≥ 7 → 1M; \`claude-{sonnet|fable|opus}-{major}\` with major ≥ 5 → 1M; the \`[1m]\` suffix → 1M; else 200k), so a new release in these families stays correct without an edit.`
- Before: `The canonical model enum reuse-condition-4 compares against is \`sonnet\`, \`opus\`, \`haiku\` (the \`dispatch build\` effective_model values).`
- After: `The canonical model enum reuse-condition-4 compares against is \`sonnet\`, \`opus\`, \`haiku\`, \`fable\` (the \`dispatch build\` effective_model values).`

`docs/schema/workflow-readme.mdschema.yml` — both `model` enums (defaults ~line 72, states[] ~line 118) gain `"fable"` after `"haiku"`.

`skills/first-officer/references/fo-dispatch-core.md:93` — checked: it defers to "the Claude enum in `claude-fo-dispatch.md`" with no literal enum, so **no edit needed** there; the codex line ("the thread's model") already matches the host-overlay design.

## Stage Report: ideation

- DONE: ACs are measured RED/GREEN behavioral tests: a README declaring model: fable on a stage builds an artifact with model=fable (RED today: enum error at build.go:59), context-budget on a sonnet-5 member reports the 1M limit (RED today: 200k default), and pre-5 sonnet / opus-family / [1m]-suffix behavior is pinned unchanged.
  ACs 1–8 above; AC-2 measures the reuse_ok flip (125.0/false → 25.0/true at 250k resident) against the 200k baseline; AC-4 pins pre-5/opus/[1m]/haiku rows in the boundary table.
- DONE: The riskiest unknown is verified, not guessed: the actual sonnet-5 model-id string shape the runtime stamps into team config (family prefix, date suffix, [1m] interplay) — inspect real config/transcript data or spike it, then record the forward-family-rule decision (regex shape) in the task body.
  Verified from real data: runtime jsonl stamps bare `claude-sonnet-5` / `claude-fable-5` (no date, no minor; 1,524 / 83,546 entries), configs carry aliases incl. `fable`; regex decision `^claude-(sonnet|fable|opus)-(\d+)` major ≥ 5 → 1M recorded in the body with guard cases; bare fable-5 observed at 561,372 resident tokens (proves >200k window).
- DONE: The contract-prose touch points are enumerated with concrete before/after wording (claude-fo-dispatch.md's break-glass canonical enum line and context-budget canonical-enum paragraph; any fo-dispatch-core.md mention) — with proof anchored in the binary's tests per the proof policy, prose updated to match, never prose-grep as evidence.
  Both claude-fo-dispatch.md lines (52, 140) have before/after diffs in the body; fo-dispatch-core.md:93 checked and needs no edit (deferral, no literal enum); schema yml enums included.

### Summary

Fleshed out the task: fable joins the Agent-schema enum at both dispatch enum sites (build.go and the previously-unlisted standing.go duplicate, to be deduplicated), the host-overlay question is resolved as host-scoped validation with ignore-with-note on codex/pi (rejected per-host keys and abstract tiers as premature), and the context window gets one 5-generation family rule (sonnet/fable/opus major ≥ 5 → 1M, haiku deliberately excluded). The riskiest unknown was settled with real transcript/config data rather than a spike; fable-5's 1M window (observed 561k resident on a bare id) is a flagged scope addition since the current rule makes healthy fable ensigns permanently non-reusable.
