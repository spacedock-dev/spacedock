---
title: Pi FO runtime — apply runtime-support.md principles (positive bindings, no step coupling, capability names)
status: validated
source: "Captain (2026-06-20): sweeping runtime-contract verbosity post-#414. Three spots in skills/first-officer/references/pi-first-officer-runtime.md violate docs/runtime-support.md's own 'Runtime contract principles' — negative host contrast, mutable-step-number coupling, Claude-centric enum contrast. The pi adapter was heavily rewritten by 0223 concurrently, so the prior prose sweeps didn't reach it."
score:
started: 2026-06-20T18:09:00Z
completed: 2026-06-20T19:30:00Z
verdict: PASSED
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

## Stage Report: validation

- **AC-1 (positive bindings — structural review, not prose-grep):** PASS. Verified against the committed files at 15509fe6 (not the report). FIX 1 (pi-first-officer-runtime.md:7) — the negative Claude sentences ("does not expose Claude Code team-tool signatures" / "Do not call or ask workers to call Claude team tools") are GONE; the section opens with "Pi is a first-class runtime target. Pi dispatch uses a Pi-native substrate selected by the launch/test harness:"; the pi-subagents/pi-agent-teams substrate bullets are unchanged (that IS the positive binding). FIX 2 (line 62) — "fo-merge-core.md's Merge-and-Cleanup step 10" is GONE, replaced with the capability binding "This is the Pi realization of `«worker.shutdown»` — the terminal-boundary teardown the shared core names. It runs at the terminal boundary whether the merge ran locally or via a PR host." No step number; no duplicated teardown sequence. FIX 3 (line 40) — reframed positively: "**Pi's model-space binding is provider/model strings**: `z-ai/glm-5.2`, `~openai/gpt-mini-latest`, ... — any valid pi-subagents model string ... Reuse-condition-4's comparator operates on these Pi-native strings." The "Claude-centric" / "Pi has no such enum" / "no Claude-centric enum" sentences are GONE from this paragraph. FIX 4 (pi-ensign-runtime.md:9) — "Do not assume Claude team tools exist in Pi." is GONE; the positive completion statement ("Completion is reported by the worker's final result in the Pi turn or by the active Pi adapter's task-completion notification.") is retained. All four edited regions are positive Pi bindings; no negative Claude contrast remains inside the edited regions. The review bound two independent values (the runtime-support.md positive-binding principle + the adapter text); it was NOT a "phrase is absent" tautology.
- **AC-2 (substance preserved):** PASS. Pi substrate list (`pi-subagents` default + `pi-agent-teams` optional, with the subagent(...)/teams action mapping) intact (lines 8-10). Shutdown specifics intact: "For `pi-subagents`, a completed child invocation needs no mailbox shutdown. Mark the worker complete/closed in first-officer memory and continue." + "For `pi-agent-teams`, use the adapter's lifecycle mapping to request member shutdown or end the team run. Do not emulate Claude team deletion." Model-space declaration intact: provider/model strings + reuse-condition-4 comparator on Pi-native strings. The rewrites are shape-only.
- **AC-3 (gates):** PASS. `go test ./...` green; `go test ./... -race` green; `gofmt -l ./cmd ./internal` clean. (The two task files are `.md`; gofmt does not apply to them — the AGENTS.md `gofmt -w ./cmd ./internal` invocation covers only Go, and is clean.)
- **Adversarial audit (a) — re-introduce negative contrast in ensign:** I restored "Do not assume Claude team tools exist in Pi." in pi-ensign-runtime.md and re-ran the full contractlint suite: it stayed GREEN. No automated test catches the regression — a human structural review would, but the guard is absent. This is the determinant for the contractlint-test recommendation below.
- **Adversarial audit (b) — re-introduce step-number coupling in Shutdown:** I restored "fo-merge-core.md's Merge-and-Cleanup step 10" in pi-first-officer-runtime.md and re-ran contractlint: it stayed GREEN. The codex FO test (TestCodexFirstOfficerRuntimeAvoidsNegativeHostContrast) bans "Merge-and-Cleanup step" but scans ONLY codex-first-officer-runtime.md, so it does not protect the pi adapter. Same gap as (a). (Both adversarial edits were reverted; the worktree task files match 15509fe6 exactly — confirmed by `git diff 15509fe6 -- <both files>` = empty.)
- **Residual assessment (follow-on paragraph after FIX 3, line 42):** FINDING — should be cleaned (captain gate; not self-fixed). The paragraph reads: "Reuse-condition-4's model-match comparator MUST operate on these Pi-native strings, not the Claude enum. Otherwise every Pi reuse would be a 'captain-session fallback value outside the enum' and force fresh dispatch — defeating reuse entirely." Two residual negative-contrast tokens: (1) "not the Claude enum" — pure negative contrast; the positive assertion "MUST operate on these Pi-native strings" already stands, so the negation of another runtime is redundant and is the same smell FIX 3 removed. (2) "captain-session fallback value outside the enum" — this is NOT a verbatim quote of the core's term. `fo-dispatch-core.md` condition 4 (line 40) says "A member stamped with a captain-session fallback value — one outside the host's canonical model space — never matches..."; the core's neutral term is "canonical model space", NOT "enum". The pi adapter substituted "enum" (a Claude-centric term, per `«worker-identity»`'s Claude realization "canonical model space = the Claude enum") for the core's neutral "canonical model space". So the scare-quoted phrase re-imports the Claude-centric vocabulary FIX 3 stripped, rather than legitimately quoting the core. Recommended clean form: "Reuse-condition-4's model-match comparator MUST operate on these Pi-native strings. Otherwise every Pi reuse would match a captain-session fallback value (one outside the canonical model space) and force fresh dispatch — defeating reuse entirely." Note (related, same smell, out of the assessed paragraph but observed): line 34 ("Stage-declared models on Pi MUST use Pi-valid model strings ..., not the Claude-centric enum.") carries the same "not the Claude-centric enum" residual; flag for the same sweep if the captain opens it. (Pre-existing, non-blocking Claude mentions at lines 16/30/36/66 are NOT findings: "without rewriting it into Claude syntax" is a transport instruction; "Claude/Codex adapters instead omit..." and "Claude inherits the team's model" are comparative technical contrast explaining model-resolution; "Do not emulate Claude team deletion" is runtime-support.md-conformant — the guide itself says "Do not emulate Claude ... unless the host really provides those tools".)
- **Contractlint-test recommendation:** FINDING/RECOMMENDATION — a pi-equivalent negative-contrast test SHOULD be added (captain gate; not self-added). Codex ships TestCodexFirstOfficerRuntimeAvoidsNegativeHostContrast and TestCodexEnsignRuntimeAvoidsNegativeHostContrast; pi has none. The adversarial audit (a)+(b) proves the gap: re-introducing the exact regressions this task removed leaves the entire `go test ./...` suite green. Recommended shape (mirrors codex; legitimate structural test binding the runtime-support.md principle to the pi adapter text, NOT prose-grep): TestPiEnsignRuntimeAvoidsNegativeHostContrast — ban the literal negative phrases ("Do not assume Claude team tools", "Claude team tools exist in Pi"); the pi ensign currently has ZERO "Claude" mentions so this passes today and would have caught the FIX-4 regression. TestPiFirstOfficerRuntimeAvoidsNegativeHostContrast — IMPORTANT: a blanket "Claude" substring ban is TOO BROAD for the pi FO file (it has legitimate, guide-conformant Claude mentions at lines 16/30/36/66); the pi FO test must ban the SPECIFIC smell phrases ("does not expose Claude Code team-tool signatures", "Do not call or ask workers to call Claude team tools", "Pi has no such enum", "Claude-centric enum", "no Claude-centric", "Merge-and-Cleanup step 10", "Merge-and-Cleanup step"), plus assert the positive bindings (`«worker.shutdown»` realization, "Pi's model-space binding is provider/model strings"). Adding these would also pin the residual finding above as a future regression guard once the residual is cleaned.
- **Worktree-hygiene observation (not a finding on the task):** the worktree has unrelated pre-existing dirty modifications to 3 internal files (internal/cli/prose_function_routing_test.go, internal/ensigncycle/haiku_loop_spike_live_test.go, internal/status/section_read.go) — smart-quote corruption in comments (`` ` `` → `"`) and whitespace alignment — NOT introduced by 15509fe6 and not part of this task. They do not affect the 4-spot fix (committed and clean: `git diff 15509fe6 -- <both skill files>` = empty) and the gates are green despite them, but the FO/captain should be aware the worktree is not pristine. The two task files themselves are exactly as committed.

### Summary

VERDICT: PASSED (with two non-blocking findings for the captain gate). All four fixes (15509fe6) are correct, shape-only, and substance-preserving; AC-1/AC-2/AC-3 all pass; gates green. The adversarial audit confirmed a real regression-guard gap: no pi-equivalent contractlint test exists, so re-introducing the removed negative contrast (ensign) or step-number coupling (FO shutdown) leaves `go test ./...` green — caught only by human review. Two findings surfaced for the captain: (1) residual negative-contrast phrasing in the follow-on paragraph after FIX 3 (line 42: "not the Claude enum" + "captain-session fallback value outside the enum" — same smell, not a verbatim core quote; should be cleaned; line 34 has the same residual); (2) recommend adding TestPiFirstOfficerRuntimeAvoidsNegativeHostContrast / TestPiEnsignRuntimeAvoidsNegativeHostContrast mirroring codex (with a specific-phrase ban, not a blanket "Claude" ban, for the FO file). Neither blocks merging the 4-spot fix; both are forward-looking hardening.
