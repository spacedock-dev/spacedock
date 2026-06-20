---
title: Pi FO runtime — apply runtime-support.md principles (positive bindings, no step coupling, capability names)
status: implementation
source: "Captain (2026-06-20): sweeping runtime-contract verbosity post-#414. Three spots in skills/first-officer/references/pi-first-officer-runtime.md violate docs/runtime-support.md's own 'Runtime contract principles' — negative host contrast, mutable-step-number coupling, Claude-centric enum contrast. The pi adapter was heavily rewritten by 0223 concurrently, so the prior prose sweeps didn't reach it."
score:
started: 2026-06-20T18:09:00Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-pi-fo-runtime-runtime-support-compliance
issue:
sprint:
sprint-readiness:
id: 2ygdt8xs5wxnpys9vrc8djaq
---

# Pi FO runtime — apply runtime-support.md principles

## End value

`skills/first-officer/references/pi-first-officer-runtime.md` complies with `docs/runtime-support.md`'s "Runtime contract principles" — positive Pi bindings (no negative Claude contrast), `«worker.shutdown»` capability binding (no mutable step-number coupling), and a positive Pi model-space binding (no Claude-centric enum contrast). The pi adapter reads as a self-contained Pi binding, not a diff against Claude.

## Problem — three violations of docs/runtime-support.md

`docs/runtime-support.md` "Runtime contract principles" encodes three rules the pi adapter violates:

### 1. Negative host contrast (line ~7, Runtime Shape)

Current: *"Pi is a first-class runtime target, but it does not expose Claude Code team-tool signatures. Do not call or ask workers to call Claude team tools."*

Violates: *"Write adapters as positive bindings for the host that is running. Prefer 'Pi maps `«worker.shutdown»` to ...' over 'Pi has no Claude TeamDelete.' Avoid negative contrast against another runtime unless the document is deliberately explaining a migration or compatibility hazard."*

### 2. Mutable step-number coupling (line ~62, Shutdown)

Current: *"This is the Pi terminal teardown — fo-merge-core.md's Merge-and-Cleanup step 10, mandatory at the terminal boundary..."*

Violates: *"Do not couple adapter text to mutable procedure step numbers. If a shared procedure says to run `«worker.shutdown»` at the terminal boundary, adapters should bind `«worker.shutdown»`; they should not say 'Merge-and-Cleanup step 10' or duplicate the shared teardown sequence."*

### 4. Negative Claude contrast — ensign (skills/ensign/references/pi-ensign-runtime.md:9)

Current: *"Do not assume Claude team tools exist in Pi. Completion is reported by the worker's final result in the Pi turn or by the active Pi adapter's task-completion notification."*

Violates: the same positive-binding principle (lower risk — it's runtime-specific — but still better as a positive Pi completion/clarification binding). The second sentence (completion = worker's final result / adapter notification) is the positive substance; the first sentence ("Do not assume Claude team tools exist in Pi") is the negative contrast that can drop.

### 3. Claude-centric enum contrast (line ~40, Canonical Model Space)

Current: *"The core's enum assumption is Claude-centric (`sonnet`/`opus`/`haiku`). **Pi has no such enum.** On Pi, any provider/model string is a valid model... There is no Claude-centric `sonnet`/`opus`/`haiku` enum on Pi."*

Violates: the positive-binding principle. The substance (Pi's model space is provider/model strings) is valid; the shape (defining Pi by negating Claude's enum) is the smell.

## Approach — positive Pi bindings, cite runtime-support.md

Three targeted rewrites. Each replaces the negative/contrastive framing with a positive Pi binding. The substance is preserved; only the shape changes.

### Fix 1 (Runtime Shape) — positive Pi substrate binding

Replace the negative Claude contrast with a positive statement of Pi's dispatch substrate. Drop "does not expose Claude Code team-tool signatures. Do not call or ask workers to call Claude team tools." Keep the `pi-subagents` / `pi-agent-teams` substrate list (that IS the positive binding). Shape: "Pi dispatch uses a Pi-native substrate selected by the launch/test harness:" as the opening (the substrate list follows); the negative Claude sentences removed.

### Fix 2 (Shutdown) — bind «worker.shutdown», drop step 10

Replace "fo-merge-core.md's Merge-and-Cleanup step 10" with the capability binding: "This is the Pi realization of `«worker.shutdown»` — the terminal-boundary teardown the shared core names. It runs at the terminal boundary whether the merge ran locally or via a PR host." No step number; no duplicated teardown sequence. The `pi-subagents` / `pi-agent-teams` shutdown specifics (no mailbox shutdown / adapter lifecycle mapping) stay — those are the positive Pi bindings.

### Fix 3 (Canonical Model Space) — positive Pi model-space binding

Reframe from "Pi has no such enum" (negating Claude) to "Pi model-space binding = provider/model strings" (asserting Pi). Drop "The core's enum assumption is Claude-centric" and "There is no Claude-centric `sonnet`/`opus`/`haiku` enum on Pi." Keep the substance: the Pi canonical model space is all valid pi-subagents model strings (provider-qualified or `~`-prefixed); reuse-condition-4's comparator operates on these Pi-native strings. Cite `docs/runtime-support.md`'s positive-binding principle if a rationale sentence is needed.

### Fix 4 (ensign, line ~9) — positive Pi completion binding

Replace "Do not assume Claude team tools exist in Pi. Completion is reported by the worker's final result in the Pi turn or by the active Pi adapter's task-completion notification." with a positive Pi completion binding: "Completion is reported by the worker's final result in the Pi turn or by the active Pi adapter's task-completion notification." (Drop the negative "Do not assume Claude team tools exist in Pi." sentence — the positive completion statement that follows is the binding.) Lower-risk than fixes 1-3 (runtime-specific), but folds into the same sweep for consistency.

## Scope

In scope:
- Four targeted rewrites: three in `skills/first-officer/references/pi-first-officer-runtime.md` + one in `skills/ensign/references/pi-ensign-runtime.md`.
- Cite `docs/runtime-support.md` as the authority where a rationale helps (one sentence max — don't bloat).

Out of scope:
- The rest of the pi adapter — only the four flagged spots.
- Other host adapters (claude/codex) — separate sweep work.
- `docs/runtime-support.md` itself — it's the authority; don't edit the guide.
- Substantive contract changes — this is prose-shape only; the substance (Pi's substrate, shutdown behavior, model space) is unchanged.

## Acceptance criteria (provisional — ideation finalizes; proof = behavior, never prose-grep)

**AC-1 — The four spots are positive Pi bindings, no negative Claude contrast.**
Verified by: a structural review that (a) the Runtime Shape section opens with the Pi substrate (no "does not expose Claude" / "Do not call... Claude team tools"); (b) the Shutdown section binds `«worker.shutdown»` (no "Merge-and-Cleanup step 10"); (c) the Canonical Model Space section asserts Pi's model space positively (no "Pi has no such enum" / "no Claude-centric enum"); (d) the ensign Agent Surface section states Pi completion positively (no "Do not assume Claude team tools exist in Pi"). Note: a contractlint "the phrase is absent" check is the prose-grep tautology the dev-workflow proof policy bans — the structural review is the gate, binding two independent values (the runtime-support.md principle and the adapter text) that can diverge.

**AC-2 — The substance is preserved.**
Verified by: the Pi substrate list (`pi-subagents`/`pi-agent-teams`), the shutdown specifics (no mailbox shutdown / adapter lifecycle mapping), and the model-space declaration (provider/model strings; reuse-condition-4 comparator on Pi-native strings) are all still present and correct. The rewrites are shape-only.

**AC-3 — Gates green.**
Verified by: `go test ./...` (the contractlint/structural tests that pin host-neutrality must stay green); `gofmt -l` clean.

## Test plan

- Structural review (AC-1, AC-2): the four spots rewritten as positive bindings; substance preserved.
- `go test ./...` (AC-3): the host-neutrality contractlint tests (e.g. `prose_inflator_locks_test.go` on the sweep branch, if merged to main) stay green; if any test pins the removed negative phrasing, update it to pin the positive binding instead (record why — the test was pinning the smell, not the contract).
- `gofmt -l` clean.
- This is a shipped-contract change (high-stakes surface per the dev-workflow proof policy). Detached adversarial audit at validation: construct an edit that re-introduces a negative contrast and confirm a reviewer catches it.

## Related

- `docs/runtime-support.md` "Runtime contract principles" — the authority (positive bindings; no step coupling; capability names in core).
- `skills/first-officer/references/pi-first-officer-runtime.md` — the file to edit (lines ~7, ~40, ~62).
- PR #414 (merged) + `shared-contract-prose-sweep` (#276, merged) — the prior sweeps that didn't reach the pi adapter (0223 was rewriting it concurrently).
- The 0223 capstone (`pi-back-channel-dispatch`, #409) — the most recent pi adapter rewrite; these three spots survived it.

## Stage Report: implementation

- DONE: FIX 1 — Runtime Shape (pi-first-officer-runtime.md): Dropped negative Claude contrast ("does not expose Claude Code team-tool signatures" / "Do not call or ask workers to call Claude team tools"). Opens with "Pi is a first-class runtime target. Pi dispatch uses a Pi-native substrate selected by the launch/test harness:" — the pi-subagents / pi-agent-teams substrate bullets follow unchanged.
- DONE: FIX 2 — Shutdown (pi-first-officer-runtime.md): Replaced mutable step-number coupling ("fo-merge-core.md's Merge-and-Cleanup step 10") with the capability binding: "This is the Pi realization of `«worker.shutdown»` — the terminal-boundary teardown the shared core names. It runs at the terminal boundary whether the merge ran locally or via a PR host." No step number; no duplicated teardown sequence. The pi-subagents / pi-agent-teams shutdown specifics unchanged.
- DONE: FIX 3 — Canonical Model Space (pi-first-officer-runtime.md): Reframed from negating Claude's enum to asserting Pi's model space. Dropped "The core's enum assumption is Claude-centric" / "Pi has no such enum" / "There is no Claude-centric sonnet/opus/haiku enum on Pi." Now reads: "The dispatch core's reuse-condition-4 compares a worker's stamped model against the host's canonical model space. **Pi's model-space binding is provider/model strings**: `z-ai/glm-5.2`, `~openai/gpt-mini-latest`, `anthropic/claude-sonnet-4`, `google/gemini-2.5-pro`, etc. — any valid pi-subagents model string (provider-qualified or `~`-prefixed). Reuse-condition-4's comparator operates on these Pi-native strings." Follow-on declaration paragraph preserved unchanged.
- DONE: FIX 4 — ensign Agent Surface (pi-ensign-runtime.md): Dropped negative "Do not assume Claude team tools exist in Pi." sentence. Kept positive completion statement: "Completion is reported by the worker's final result in the Pi turn or by the active Pi adapter's task-completion notification." Surrounding section preserved.
- Gates: `go test ./...` green; `go test ./... -race` green; `gofmt -w ./cmd ./internal` clean. No contractlint/structural test pinned the removed negative phrasing — the existing host-neutrality tests (TestDispatchCoreHasNoClaudeModelToken, TestDispatchCoreHasNoClaudeTeamImperative, TestFOContractCoresHaveNoDeferredTierToken) scan fo-dispatch-core.md or the deferred-tier token set, not the pi adapter's Claude-contrast sentences. The codex negative-contrast tests (TestCodexFirstOfficerRuntimeAvoidsNegativeHostContrast, TestCodexEnsignRuntimeAvoidsNegativeHostContrast) are codex-only and do not scan the pi adapter. No test update was needed.

### Summary

Applied docs/runtime-support.md's Runtime contract principles to 4 spots across 2 Pi runtime adapter files. All rewrites are shape-only — the substance (Pi substrate list, shutdown specifics, model-space declaration, ensign completion mechanism) is preserved. Worktree commit: 15509fe6 on spacedock-ensign/pi-fo-runtime-runtime-support-compliance. No test updates were required; all gates green.
