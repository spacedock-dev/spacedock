---
title: Contractlint runtime-semantics retirement — codex and pi phrase checks become behavior tests
status: ideation
source: "Contractlint antipattern sweep, 2026-07-11: codex_multi_agent_v2_contract_test.go and Codex portions of runtime_binding_block_test.go assert runtime meaning from host-adapter prose."
score: 0.34
id: 8413fc05vpp8116k54x8br15
sprint: 0260-proportionality
group: contract-cleanups
started: 2026-07-20T05:04:09Z
---

## Problem

The Codex and Pi runtime-adapter contract tests encode runtime meaning — capability binding, spawn/wait/reuse behavior, completion-signal routing, host-neutrality, model-space handling — as expected or banned prose in the adapter markdown. They pass after a behavior regression when the wording survives, and red after a safe paraphrase. The two families are the codex set (`codex_multi_agent_v2_contract_test.go` plus the Codex halves of `runtime_binding_block_test.go`) and the pi set (`pi_runtime_negative_contrast_test.go` plus the Pi halves of `runtime_binding_block_test.go`). Across the three files, ~121 literal-in-adapter-prose entries drive 24 `strings.Contains` / `reflect.DeepEqual` / `strings.Count` assertions — the committed-prose-test anti-pattern this sprint retires.

## Proposed approach

Route every runtime-meaning claim to the narrowest independent source that can diverge from the adapter prose, then delete the prose assertion. Three routes exist, and each token maps to exactly one:

1. **spacedock-emitted Go source (bind, reconcile pattern).** Where spacedock's own Go code emits the tool arg-shape or action token the doc names, bind the doc surface to the Go surface and red on drift in either. Confirmed independent sources: `internal/dispatch/codex_v2_adapter.go` `CodexMultiAgentV2Spawn.ToolArgs()` emits `{task_name, message, fork_turns}` (matches the doc's `spawn_agent(task_name,message,fork_turns="none")`); `internal/piruntime/teams.go` `TeamsAction` declares `delegate` / `message_dm` / `member_shutdown` / `team_done`; `piruntime.SubagentStageDispatch` sets `context:"fresh"` and omits `acceptance`. The 8 runtime `«capability»` tokens bind as a set equal across the codex block, the pi block, and the core capability-`«fn»` headings in `fo-dispatch-core.md` — three independent surfaces, set-equality, red on any add/drop/rename.
2. **dispatch-build fixture (already live).** Host-neutrality (no `Agent(` / `SendMessage` / `TeamCreate` / `Skill(` leak), the pi completion form, model:null+ignore-note, and the split-root path are already proven by `internal/dispatch/build_pi_host_test.go` driving `dispatch build --host pi`. Cite it; do not duplicate as prose.
3. **live host lane.** Substrate-native tokens spacedock does not emit — codex `wait_agent` / `followup_task` / `list_agents` / `send_message`, pi `subagent(...)` / `intercom(...)` / `member_spawn` / `contact_supervisor` — have their runtime meaning owned by the gated live lanes `internal/ensigncycle/codex_live_runner_test.go` and `pi_live_runner_test.go`, which record the actual tool call and resulting workflow state. Delete the prose; cite the lane.

Contractlint retains only structural properties: reference closure, the binding-block token set, host-neutral file placement, and section-scoping (tool tokens contained to the runtime-binding/harness section). Each retained structural-absence guard (host-neutrality, "no absolute step ordinal") keeps its discriminator control — the non-vacuity partner the repo already mandates.

Per-check disposition (the two families, ten test functions):

| Test function | Claim | Disposition |
| --- | --- | --- |
| `TestCodexToolNamesStayInRuntimeBindingSection` | codex spawn signature present + tool tokens contained | BIND signature ↔ `ToolArgs()` keys; RETAIN containment/absence-outside as structural |
| `TestCodexEnsignRuntimeAvoidsNegativeHostContrast` | no negative host-contrast; positive budget/completion prose | RETAIN negative-contrast as structural-absence + discriminator; DELETE positive prose (budget/completion owned by codex live lane + capability-set binding) |
| `TestFeedbackRejectionFlowStaysHostNeutral` | shared skill names no host/tool; uses neutral `«tokens»` | RETAIN host/tool ban as structural-absence + discriminator; BIND `«token»` refs ⊆ core capability set (DELETE literal-present asserts) |
| `TestCodexAndPiFirstOfficerRuntimeBindingBlocks` | both blocks bind exactly the 8-capability list | BIND codex block == pi block == core `«fn»` headings (retires the hardcoded slice) |
| `TestCodexAndPiFirstOfficerRuntimeLifecycleHeadingsRemoved` | old lifecycle headings gone | DELETE (subsumed: a re-introduced divergent block reds the capability-set set-equality) |
| `TestPiFirstOfficerRuntimeRejectMutableStepAndNegativeContrast` | no mutable step-ordinal coupling; no negative contrast | RETAIN "no absolute step ordinal" + negative-contrast as structural-absence + discriminator |
| `TestPiToolNamesStayInRuntimeBindingOrHarnessSections` | pi tool tokens present + contained | BIND emitted tokens (`delegate`/`message_dm`/`member_shutdown`/`team_done`, `context:"fresh"`/no-`acceptance`) ↔ `teams.go`/`SubagentStageDispatch`; RETAIN containment; substrate-native tokens routed to live lane |
| `TestPiFirstOfficerRuntimeSemanticsPreserved` | fresh/acceptance/completion semantics present | DELETE (every claim owned by `build_pi_host_test.go` behavior + the Go bindings above) |
| `TestPiFirstOfficerRuntimeAvoidsNegativeHostContrast` | no negative contrast; positive `«worker.shutdown»`/model-space | RETAIN negative-contrast as structural-absence + discriminator; DELETE positive prose (owned by capability-set binding + `TestBuildPiHostIgnoresModelWithNote`) |
| `TestPiEnsignRuntimeAvoidsNegativeHostContrast` | no negative contrast; positive completion/clarification | RETAIN negative-contrast as structural-absence + discriminator; DELETE positive prose (completion/clarification owned by pi live lane) |

No net-new check purpose is added: every committed test maps 1:1 to a retired phrase check's claim or is a deletion; discriminators are the required non-vacuity partners, not new surface. Adapter markdown is not modified.

## Out of scope

Changing Codex or Pi multi-agent policy, model defaults, or the host tool sets. Editing any adapter markdown body (`codex-first-officer-runtime.md`, `pi-first-officer-runtime.md`, `codex-ensign-runtime.md`, `pi-ensign-runtime.md`, feedback-rejection-flow `SKILL.md`) — the tests change, the docs stay. Adding coverage for a token not already asserted by the retired families (minting).

## Acceptance criteria

**AC-1 (VALUE) - Runtime-meaning claims in the two families are proven by a value computed against an independent source, not by a literal-in-adapter-prose assertion.**
Baseline that can move the wrong way: today the three files run ~121 literal-in-adapter-prose entries through 24 `strings.Contains`/`reflect.DeepEqual`/`strings.Count` sites. End-state: the count of committed assertions that assert runtime MEANING via an adapter-prose literal is 0 (tolerance 0). The only surviving `strings.Contains` residue is structural — a host-neutrality or no-absolute-step-ordinal absence guard (each with a discriminator control) or a section-scoping containment check.
Verified by: a validation-time inventory diff vs `origin/main` showing runtime-meaning literal-in-adapter-prose assertions at 0, plus `go test ./internal/contractlint/... ./internal/dispatch/...` green. (The inventory grep is validation-time evidence only; the committed tests are the bindings and the cited behavior lanes.)

**AC-2 (VALUE) - Each spacedock-emitted adapter token binds to its Go source and reds under a planted divergence.**
The codex spawn signature binds to `codex_v2_adapter.go` `ToolArgs()` keys; the pi-agent-teams action tokens bind to `teams.go` `TeamsAction`; `context:"fresh"`/no-`acceptance` binds to `SubagentStageDispatch`; the 8 `«capability»` tokens bind equal across the codex block, the pi block, and the core `«fn»` headings.
Verified by: each binding passes today and reds when its independent source is mutated away from the doc (rename a `ToolArgs` key, drop a `teams.go` const, remove a capability from the core `«fn»` set) — the mutation/discriminator control the repo already requires; 0 false greens across the planted-divergence suite.

**AC-3 - Substrate-native tokens with no spacedock Go source are owned by a live host lane, and the retired prose is deleted with the lane cited.**
Codex `wait_agent`/`followup_task`/`list_agents`/`send_message` and pi `subagent(...)`/`intercom(...)`/`member_spawn`/`contact_supervisor` carry no spacedock-emitted Go source; their runtime meaning is owned by `codex_live_runner_test.go` / `pi_live_runner_test.go`.
Verified by: each deleted prose assertion's claim maps to a named live-lane (or build-fixture) assertion; no new live assertion is minted.

**AC-4 - The adapter docs are unchanged and any doc/code divergence a new binding surfaces is filed, not silently patched.**
Verified by: `git diff --stat` for this task touches only `internal/contractlint/**` (and test support), no `skills/**` doc bodies; the known `member_spawn` doc-only reference (present in `pi-first-officer-runtime.md`, absent from `teams.go`) is recorded as a finding and routed to the live pi lane rather than folded into a struct binding.

## Test plan

Inventory the ten phrase-check functions by claim (done in ideation, recorded in the disposition table). Apply each disposition: three Go-source bindings (codex spawn args, pi teams actions + subagent shape, capability-set cross-source), the retained structural-absence guards with discriminators, the containment/section-scoping retentions, and the deletions with cited owning lanes. Add a planted-divergence control per new binding (mutate the Go source, assert red). Run `go test ./internal/contractlint/... ./internal/dispatch/...` first (fast, fixture/unit only), then `go test ./...` and `go test ./... -race`. The live lanes are cited, not modified: run each once (gated by codex auth / pi login+PATH, serial) to confirm it still owns the routed claim; no new live assertion is added.

Estimate + tolerance: ~3 binding test functions + ~3 structural-absence functions (each +1 discriminator) replace the 10 phrase-check functions; net committed test-function count roughly flat (±2); runtime-meaning literal-in-adapter-prose assertions 121 → 0 (tolerance 0); line delta expected near-flat to negative. Complexity: medium — no new parser; every extraction mechanism is already proven in-repo.

## Spike / mechanism record

No spike needed — every mechanism the design rests on is proven in-repo, and the one risky assumption (that each doc token's claimed independent source actually matches today) was verified inline during ideation:

- AST var/const string-set extraction: proven by `reconcile_class_binding_test.go` (`helperDriftClasses`).
- Markdown `## section` extraction: proven by `capability_binding_test.go` / `runtime_binding_block_test.go` (`extractMarkdownSection`).
- Build-output JSON parse + host-neutrality/completion assertions: proven by `build_pi_host_test.go` (green today).
- Live tool-call recording: proven by `codex_live_runner_test.go` / `pi_live_runner_test.go` (gated lanes).

Inline verification (the spike), each binding passes today: codex `spawn_agent(task_name,message,fork_turns="none")` == `ToolArgs()` keys; pi `delegate`/`message_dm`/`member_shutdown`/`team_done` == `teams.go` consts; pi `context:"fresh"`/no-`acceptance` == `SubagentStageDispatch` (build fixture green); the 8 `«capability»` tokens equal across the codex block, the pi block, and the core `«fn»` headings.

One divergence found and routed, not silently patched: `member_spawn` appears in `pi-first-officer-runtime.md` but has no `teams.go` const — it is a pi-agent-teams substrate-native action spacedock does not emit, so it binds to the live pi lane (AC-3), not the struct binding. A naive `doc-tokens ⊆ teams.go` set-membership would red on it; implementation scopes the struct binding to spacedock-emitted tokens only.

## Merged scope (adopted cross-review re-lock, 2026-07-20)

Absorbs `contractlint-pi-runtime-semantics-retirement` — the pi split is mechanically identical to the codex split; one member retires both runtime-semantics phrase checks into live/fixture behavior tests.

## Stage Report: ideation

- DONE: retire BOTH runtime-semantics phrase-check sets (codex + pi) from contractlint into live or fixture behavior tests; no scope beyond the two named check families.
  Per-check disposition table covers all ten functions across the three in-scope files; out-of-scope explicitly excludes adapter-doc edits and any token not already asserted by these families.
- DONE: Every replacement committed test computes a real value against an independent source that can diverge — name each check's replacement shape or its deletion-with-reason; no new check surface without consent.
  Three Go-source bindings (codex `ToolArgs()` keys, pi `teams.go` `TeamsAction` consts + `SubagentStageDispatch`, 8-`«capability»` set across codex/pi blocks + core `«fn»` headings); structural-absence + containment retentions with discriminators; deletions routed to `build_pi_host_test.go` / `codex_live_runner_test.go` / `pi_live_runner_test.go`. No net-new check purpose (AC-2/AC-3).
- DONE: Written expected surface + tolerance declared; "no spike needed" recorded with proven mechanisms.
  AC-1 value: runtime-meaning literal-in-adapter-prose assertions 121 → 0 (tolerance 0) vs origin/main; estimate ~3 bindings + ~3 absence guards, test-function count flat ±2. Spike/mechanism record proves all four mechanisms in-repo and verifies each binding passes today; one divergence (`member_spawn` doc-only) found and routed to the live lane.

### Summary

Designed the retirement of the codex and pi runtime-semantics phrase checks (three files, ten functions, ~121 prose literals) into value-computing tests. Each runtime-meaning claim routes to one of three independent sources — a spacedock-emitted Go struct/const (bind, reconcile pattern), the already-live `dispatch build --host pi` fixture, or a gated live host lane — with pure prose-greps deleted and only structural-absence/containment residue retained (each keeping its discriminator). Verified inline that every proposed binding matches its source and passes today; the one divergence found (`member_spawn` in the pi doc but not in `teams.go`, a substrate-native token) is routed to the live lane rather than silently patched.
