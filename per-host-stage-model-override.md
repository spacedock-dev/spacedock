---
id: e3g7s1jtr05fp6n1w0z89w4v
title: Per-stage model declaration must support divergent per-host values (and Codex's effort axis), not just ignore-elsewhere
status: ideation
source: Captain conversation 2026-07-03, following dispatch-model-space-fable-sonnet5-1m (archived, id wcex4yjx4mvecybxjb43gwtw, PR #463) — that task's "host overlay" design deliberately rejected per-host keys ("no workflow needs it yet") in favor of validate-on-claude / ignore-with-note-and-null-elsewhere. A concrete need surfaced immediately after: declaring `fable` for the `ideation` stage on Claude while wanting a different, specific model (e.g. a Codex model id) plus a reasoning-effort tier (e.g. "xhigh") when the same workflow runs under Codex. The current mechanism cannot express this — Codex/Pi dispatch unconditionally gets `model: null`, with zero per-host substitution capability.
started: 2026-07-03T04:12:53Z
completed:
verdict:
score:
worktree:
issue:
---

The shipped per-stage `model:` mechanism (`internal/dispatch/build.go`'s `runBuildFields`) only routes a declared model on Claude; on Codex and Pi the declared value is discarded entirely (an "ignored-with-note" stderr line, `model` always emitted as `null`). There is no way today for a workflow README to say "this stage runs fable on Claude, but a specific Codex model at a specific reasoning effort on Codex" — the mechanism was scoped to graceful degradation, not actual cross-host model routing, despite that being the stated motivation for the original task.

## Problem

A workflow author cannot declare divergent per-host models for one stage. The concrete case (captain, 2026-07-03): the `ideation` stage should run `fable` when the workflow runs under Claude, and `gpt-5.5` at reasoning effort `xhigh` when the same workflow runs under Codex. The captain's own Codex config (`~/.codex/config.toml`: `model = "gpt-5.5"`, `model_reasoning_effort = "xhigh"`) is exactly that pairing — but at session scope only, with no per-stage control. The shipped mechanism (PR #463) cannot express any of it:

- `model:`'s value space is Claude-only. On host=codex/pi the declared value is discarded (ignore-with-note; `model` emitted as the JSON literal `null`). Graceful degradation was that task's entire scope; per-host substitution was explicitly deferred as "no workflow needs it yet". The need arrived the next day.
- Codex's effort axis has no representation anywhere: no README key to write `xhigh` into, no envelope field to carry it, no adapter mapping to `spawn_agent`.
- The consumption side already exists: Codex 0.142.5's `spawn_agent` accepts per-spawn `model` AND `reasoning_effort` (evidence below). Only spacedock's declaration/emission side is missing.

**Baseline (RED), driven 2026-07-03** against the current binary (built from HEAD) over a fixture README whose `ideation` stage declares `model: fable`, `model-codex: gpt-5.5`, `effort-codex: xhigh`:

- `--host codex`: envelope `"model": null`, no effort field exists; stderr `[build] declared model 'fable' ignored on host codex: outside codex's dispatch-settable model space; emitting model=null`. The suffix keys are silently dropped (unknown optional keys never reach the parser output).
- `--host claude`: envelope `"model": "fable"`, stderr `[build] effective_model=fable (from stage) → Agent model=fable`. The suffix keys do not interfere — the chosen syntax is purely additive, no collision.
- `--host pi`: same shape as codex (null + ignore-note).

## Verified evidence (riskiest unknown: does Codex expose a per-spawn model/effort surface at all?)

The design is worthless if Codex workers can only inherit the thread's model — which was true when the codex adapter was written (live probe 2026-06-20, `codex-multi-agent-v2-runtime-support`: "omit unsupported `description`, `subagent_type`, and `model` arguments"). Verified 2026-07-03 that this has changed, two independent ways:

1. **Installed binary strings** (codex-cli 0.142.5, `/opt/homebrew/Caskroom/codex/0.142.5/codex-aarch64-apple-darwin`): the spawn_agent tool instruction reads "sub-agents that inherit your current model by default. Do not set the `model` field unless the user explicitly asks for a different model or there is a clear task-specific reason"; runtime validation errors "Unknown model `%s` for spawn_agent. Available models: ..." and "Reasoning effort `%s` is not supported for model `%s`. Supported reasoning efforts: ..."; per-spawn telemetry fields `requested_model`, `requested_reasoning_effort`.
2. **Version-pinned source** (`openai/codex` tag rust-v0.142.5, `core/src/tools/handlers/multi_agents_v2/spawn.rs`): `SpawnAgentArgs { message, task_name, agent_type?, model?, reasoning_effort?, service_tier?, fork_turns?, fork_context? }` — `model` and `reasoning_effort` are independently optional (effort-only spawns are structurally valid).

Two constraints that shape the design fell out of the same evidence:

- **Fork-mode interaction:** spawn overrides (`model`/`reasoning_effort`/role) are REJECTED when `fork_turns="all"` (full-history fork) — and `"all"` is the default when `fork_turns` is absent. Overrides are permitted with `"none"` or an integer. So any spawn forwarding these overrides must explicitly pass `fork_turns="none"` (which is what spacedock's self-contained file-pointer dispatch wants anyway, per the codex runtime prose).
- **Both value spaces are dynamic:** the available-models list is composed at runtime (account/catalog/version dependent), and the effort enum itself grows across Codex releases (0.142.5 carries `none|minimal|low|medium|high|xhigh|ultra`, with per-model supported subsets — the "not supported for model" error proves the subsetting). This drives the validation decision below.

**No further spike needed** — the proven mechanisms the rest relies on: yaml.v3 flat-scalar parsing of stage optional keys (existing, tested; the baseline drive proves suffix keys are additive today), the golden-fixture build harness (existing), and the host branch in `runBuildFields` (existing, shipped by #463).

## Proposed approach

**Decision: per-host suffix keys — `model-codex:` and `effort-codex:` alongside the existing flat `model:`.** (Central design decision; reasoning and rejected alternatives recorded below.)

```yaml
- name: ideation
  gate: true
  model: fable          # Claude value space, unchanged semantics
  model-codex: gpt-5.5  # Codex model id, passthrough
  effort-codex: xhigh   # Codex reasoning effort, passthrough
```

Both keys are also legal under `stages.defaults`, resolved per axis with the existing precedence: stage > defaults > null. The axes resolve independently — a stage `effort-codex` combines with a defaults `model-codex`, and effort-only (`effort-codex: xhigh` with no `model-codex`) is meaningful: inherit the thread's model, request the effort. The pinned source proves both spawn params are independently optional, so the schema mirrors the target surface.

Per-host semantics of the resolved declaration:

- **host=claude**: unchanged — flat `model:` validates against the Agent enum and becomes `effective_model`; `model-codex`/`effort-codex` are ignored silently (they name their host; this is routing, not degradation, so no stderr note). Envelope `effort` is null.
- **host=codex**: `model-codex`/`effort-codex` become the envelope's `model`/`effort` (passthrough, no spacedock-side validation). Stderr reports what resolved, mirroring the claude line: `[build] effective_model=gpt-5.5 effort=xhigh (from stage)` (a null axis renders as `null`). The flat `model:` ignore-with-note fires ONLY when no codex-space key is declared for the build (unchanged degradation path); when the author explicitly covered codex, the note would be noise and is suppressed. A declared-but-empty value (`model-codex: ""`) is a build error naming the key — same catch-it-where-it-binds rule as claude's enum error.
- **host=pi**: unchanged — codex-suffixed keys ignored silently, flat model ignore-with-note, model/effort null. (`pi-dispatch-model-stamping` owns any future pi space.)

**Envelope**: `buildOutput` and `buildAdvanceOutput` gain `Effort *string \`json:"effort"\`` beside `Model`, always present and null unless a codex-host build resolves one — mirroring `model`'s null-literal treatment so FO adapters read both uniformly. Additive nullable field: no `schema_version` bump (flag at the gate). The advance envelope carries it for the same reason it carries model: reuse-condition-4's comparator reads `next_stage.effective_model`/`effective_effort` from `--advance` output, and a reused Codex worker's effort cannot be changed on a live thread, so an effort mismatch forces fresh dispatch exactly like a model mismatch.

**Codex adapter** (`internal/dispatch/codex_v2_adapter.go`, the executable spec of the FO prose mapping, consumed by its own tests): `CodexMultiAgentV2SpawnInput` additionally reads `model`/`effort` from the build envelope; `ToolArgs` includes `model` and `reasoning_effort` when non-empty and forces `fork_turns: "none"` whenever overrides are forwarded and no explicit fork_turns was set (required — Codex rejects overrides on the default full-history fork). Envelope name `effort` is the host-neutral domain name; the adapter owns the mapping to Codex's wire name `reasoning_effort`.

**Value-space validation decision** (checklist item): **unvalidated passthrough for both `model-codex` and `effort-codex`, with Codex as the authoritative validator at spawn time.** This is forced, not merely chosen: spacedock cannot enumerate Codex's model catalog (runtime-composed, account- and version-dependent) and must not hard-code the effort enum (it grew `xhigh` then `ultra` across recent releases, and the supported subset varies per model — a hard-coded enum would wrongly reject valid values on newer Codex and cannot catch invalid combos anyway). Codex's own spawn errors ("Unknown model ... Available models: ...", "Reasoning effort ... not supported for model ...") are precise and surface through the FO's spawn-failure path. The asymmetry with Claude's validated enum is principled: the Claude enum is the Agent tool schema spacedock itself renders against, fixed and known at release time; Codex's spaces are dynamic and owned upstream. Only structural validation lands spacedock-side (non-empty values).

**Rejected alternatives:**

1. **Structured per-host block** (`model: {claude: fable, codex: {model: ..., effort: xhigh}}`): rejected. The "one field, atomic pairing" argument dissolves on the evidence — model and effort are independently optional axes on Codex, not a forced pair. Against no expressiveness gain it costs: `scalarMap` flattens nested mappings to `""` today, so the parser needs a special-cased nested decode and a richer carrier than `Stage.optional map[string]string`; both verbatim emitters (`stagesJSONArr` at json_commands.go:304, `formatReadText` at section_read.go:181) need invented nested renderings (the `--read` mirror is line-oriented `key=value` with no nested form); the flat `model: fable` shape must keep parsing forever, so every consumer handles two shapes; the schema needs oneOf(string, object). Bigger blast radius across three packages for the same expressible outcomes, against the standing "do not overcomplicate" constraint.
2. **Combined micro-syntax** (`model-codex: gpt-5.5@xhigh` or space-joined): rejected — invents a DSL to re-join two axes the target surface keeps separate, and hyphen-joining is genuinely ambiguous against real model ids (`gpt-5.1-codex-mini`).
3. **Abstract tiers** (rejected already in #463): still the most machinery; unchanged.

The suffix-key shape is the minimal parser surface: three flat-key list extensions (`stages.go:90`, `json_commands.go:304`, `section_read.go:181`), the build.go resolution branch, the envelope field, and the adapter mapping. The `<axis>-<host>` pattern extends naturally if pi grows a model space or Claude grows an effort axis — without committing to either now.

## Out of scope

- Validating that a given Codex model id / effort combination is acceptable — Codex validates at spawn time against its live catalog; that error path already reaches the FO.
- Pi's model space (`pi-dispatch-model-stamping` owns it; no positive per-stage space exists).
- A Claude effort axis (the Agent tool has no effort parameter).
- Codex's `service_tier` spawn override (exists in `SpawnAgentArgs`; no named need — YAGNI).
- Changing `fork_turns` defaults for ALL codex dispatches. Observed pre-existing gap: the adapter omits `fork_turns`, so plain dispatches full-history-fork by default although the runtime prose says file-pointer dispatch "should normally use no inherited turn context" — follow-up candidate, not this task. Here only override-carrying spawns force `"none"` (required for the spawn to succeed at all).

## Acceptance criteria

Each criterion states how it is tested.

- **AC-1 (VALUE)** — **For one stage declaring `model: fable`, `model-codex: gpt-5.5`, `effort-codex: xhigh`, `dispatch build` emits the host-correct envelope on each host**: claude → `"model": "fable"`, `"effort": null`; codex → `"model": "gpt-5.5"`, `"effort": "xhigh"` (+ the effective_model/effort stderr line); pi → both null. Measured against the recorded 2026-07-03 baseline where the codex/pi branch emits `model: null` unconditionally and no effort field exists — the envelope bytes are the independent baseline that can move the wrong way. Verified by: command-level Go tests with golden fixtures beside `build_codex_host_test.go` / `build_pi_host_test.go` / `TestBuildModelPrecedence`, asserting envelope JSON and stderr.
- **AC-2** — **Degradation paths are byte-stable**: a codex-host build over a README declaring ONLY flat `model: fable` still exits 0 with `model: null`, `effort: null`, and the verbatim ignore-note; the existing codex/pi host-ignore and claude model-precedence goldens stay byte-identical. When `model-codex` or `effort-codex` is declared, the flat-model ignore-note is suppressed. Verified by: existing goldens unchanged + a new golden for the flat+suffix case asserting note absence.
- **AC-3** — **Axes resolve independently with stage > defaults > null per key**: stage `effort-codex` combines with defaults `model-codex`; effort-only declaration emits `model: null` + the effort; model-only emits the model + `effort: null`. Verified by: table-driven build tests + goldens.
- **AC-4** — **Passthrough is pinned as passthrough**: arbitrary values (`model-codex: some-future-model`, `effort-codex: ultra`) build exit 0 on host=codex with values emitted verbatim; a declared-but-empty value errors naming the key and stage index; claude/pi hosts ignore the same keys silently. Verified by: build tests incl. an error golden.
- **AC-5** — **The status surface carries the keys**: `status --json`'s stages array and `--read`'s text mirror include `model-codex`/`effort-codex` verbatim when declared, absent when not (presence semantics matching the other optional keys). Verified by: internal/status unit tests over `stagesJSONArr` / `formatReadText`.
- **AC-6** — **The spawn mapping is safe by construction**: the codex adapter's `ToolArgs` carry `model`/`reasoning_effort` when the envelope resolves them and always pair forwarded overrides with `fork_turns: "none"` (absent `fork_turns` means full-history fork, which Codex rejects for override spawns). Verified by: `codex_v2_adapter_test.go` cases for override-present, override-absent, and explicit-fork_turns shapes.
- **AC-7** — **`dispatch build --advance` on host=codex carries the resolved model+effort** so reuse-condition-4 can compare both against the stamped worker. Verified by: advance-mode build test + golden.
- **AC-8** — **Contract/schema prose names the shipped behavior** (concrete diffs below). Counts only paired with AC-1..7; prose is updated to match, never grepped as evidence. The diff touches `skills/**` (codex adapter prose), so per the proof policy the `codex-live` lane is REQUIRED green before merge.

## Test plan

- **Command-level Go tests + goldens** (internal/dispatch): the AC-1 three-host matrix, AC-2 byte-stability + note suppression, AC-3 precedence table, AC-4 passthrough/empty-error, AC-7 advance envelope. All on the existing golden harness; regenerate nothing existing except where AC-2 proves byte-identity.
- **Parser/status unit tests** (internal/status): AC-5 round-trip of the two new optional keys through `ParseStagesWithDefaults`, `stagesJSONArr`, `formatReadText`.
- **Adapter unit tests** (internal/dispatch): AC-6 `ToolArgs` shapes.
- **CI lanes**: deterministic lanes plus `codex-live` (REQUIRED — the diff touches the codex runtime adapter prose under `skills/**`). No new live scenario is strictly owed by our claim boundary (the envelope/adapter/prose are spacedock's; honoring the overrides at spawn is Codex's own validated behavior, cited from the pinned source), but the existing codex-live scenarios must stay green over the prose change.
- Estimated cost: small-medium — one implementation pass; every harness (goldens, host tests, adapter tests, status fixtures) already exists.

## Documentation / contract-prose diff (applied by implementation, reviewed at the gate)

`docs/schema/workflow-readme.mdschema.yml` — both `fields` blocks that carry `model` (defaults ~line 70, states[] ~line 119) gain two plain-string fields, no enum:

```yaml
"model-codex": { "type": "string" }
"effort-codex": { "type": "string" }
```

`skills/first-officer/references/codex-first-officer-runtime.md`:
- Line 9 (probe): `- «worker.spawn»: spawn_agent(task_name,message,fork_turns)` → `- «worker.spawn»: spawn_agent(task_name,message,fork_turns,model,reasoning_effort)`
- Line 19 («worker.spawn» binding), after "with the helper-emitted prompt unchanged as `message`": add `Forward the helper-emitted model/effort as spawn_agent's model/reasoning_effort when non-null, and pair forwarded overrides with fork_turns="none" — Codex rejects spawn overrides on a full-history fork, the default when fork_turns is absent.`
- Line 22 («worker-identity»): `and thread-inherited model when the helper emits null` → `and the helper-emitted model and effort when present, thread-inherited model when the helper emits null`

`skills/first-officer/references/fo-dispatch-core.md`:
- Line 48 (`«reuse.model-match»` guard): `skip when next_stage.effective_model is null` → `skip when next_stage.effective_model and next_stage.effort are both null`
- Line 49 (effect): append `On a host with an effort axis (Codex), the stamped effort is part of the comparand alongside the model.`
- Line 51 (done-when): `the models match (or the declared model is null)` → `the models and efforts match (or the declared model and effort are both null)`; the mismatch diagnostic gains the effort form: `reused worker {name} effort {X} does not match next stage effort {Y} — fresh-dispatching`
- Line 93 (host model spaces): `· **Codex:** the thread's model.` → `· **Codex:** the Codex model catalog plus its per-model effort axis (helper passthrough, validated by Codex at spawn); the thread's model when the helper emits null.`
- Line 133 (`«dispatch.build»` done-when): the forwarded-fields list gains `effort` with the per-host clause `(effort binds only on hosts with an effort axis; Claude's adapter omits it)`.

`skills/first-officer/references/claude-fo-dispatch.md` — checked: the Agent() field mapping (lines 23-31) forwards `model` only and needs no edit; `effort` is never forwarded on Claude. No edit.

## Stage Report: ideation

- DONE: Firm up the Problem statement with concrete real-world grounding (the fable-on-Claude / gpt-5.5-xhigh-on-Codex example, and why the current ignore-with-note mechanism falls short of it)
  Problem section grounds it in the captain's live Codex config (gpt-5.5 at xhigh, session-scoped only) and a recorded RED baseline drive of the current binary: codex/pi emit model null with no effort field, suffix keys silently dropped.
- DONE: Evaluate both candidate approach shapes from the seed (per-host suffix keys vs. a structured per-host block) plus any better alternative, and pick one with reasoning recorded — this is the task's central design decision, do not leave it open
  Decided: per-host suffix keys `model-codex:`/`effort-codex:`. The structured block's "atomic pairing" argument dissolves on the pinned-source evidence (model and effort are independently optional spawn axes) and it costs nested parsing, two invented emitter renderings, and permanent two-shape consumers; a combined micro-syntax DSL also rejected (hyphen-ambiguity with real ids like gpt-5.1-codex-mini).
- DONE: Determine whether Codex's effort axis needs its own validated enum the way Claude's model: does, or should be an unvalidated passthrough string
  Passthrough, and forced rather than chosen: 0.142.5's enum is {none,minimal,low,medium,high,xhigh,ultra} with per-model supported subsets, grown twice in recent releases (xhigh, ultra) — a spacedock enum would wrongly reject valid future values and cannot catch invalid combos; Codex's own spawn errors are the authoritative validator. Only non-empty structural checks land spacedock-side.
- DONE: Produce concrete, testable Acceptance criteria with Verified-by evidence, per the README's ideation rules
  ACs 1-8; AC-1 is the value-measuring criterion (host-correct envelope bytes vs the recorded always-null baseline); riskiest mechanism (does Codex expose per-spawn model/effort at all?) verified FIRST two independent ways — installed 0.142.5 binary strings and version-pinned spawn.rs source — surfacing the fork_turns="all" override-rejection constraint that shaped AC-6.

### Summary

Fleshed out the task around a verified mechanism reversal: the codex adapter's June evidence ("spawn_agent has no model arg") is obsolete — Codex 0.142.5 accepts per-spawn model and reasoning_effort, independently optional, but rejects both on the default full-history fork, so forwarded overrides must pair with fork_turns="none". The central design decision is settled as per-host suffix keys (minimal parser surface, matches the target's independent axes) with unvalidated passthrough for both codex values, an always-present nullable `effort` envelope field mirroring `model`, and concrete contract-prose diffs for the codex runtime adapter and reuse-condition-4's effort comparand.
