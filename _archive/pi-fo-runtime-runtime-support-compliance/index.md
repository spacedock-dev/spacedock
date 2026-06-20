---
title: Pi FO runtime — apply runtime-support.md principles (positive bindings, no step coupling, capability names)
status: validation
source: "Captain (2026-06-20): sweeping runtime-contract verbosity post-#414. Three spots in skills/first-officer/references/pi-first-officer-runtime.md violate docs/runtime-support.md's own 'Runtime contract principles' — negative host contrast, mutable-step-number coupling, Claude-centric enum contrast. The pi adapter was heavily rewritten by 0223 concurrently, so the prior prose sweeps didn't reach it."
score:
started: 2026-06-20T18:09:00Z
completed: 2026-06-20T19:30:00Z
verdict: PASSED
worktree:
issue:
sprint:
sprint-readiness:
id: 2ygdt8xs5wxnpys9vrc8djaq
mod-block:
pr: "#417"
archived: 2026-06-20T20:26:11Z
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

## Stage Report: implementation (feedback cycle 1)

- DONE: F1a (pi-first-officer-runtime.md:42) — Cleaned the residual negative-contrast phrasing the cycle-0 4-spot fix missed. Dropped "not the Claude enum" (pure negative contrast; the positive "MUST operate on these Pi-native strings" assertion stands alone). Replaced the Claude-centric scare-quote "captain-session fallback value outside the enum" with the dispatch core's neutral term: the cycle-0 validation verified `fo-dispatch-core.md` condition 4 says "canonical model space", NOT "enum" — the adapter had substituted the Claude-centric "enum" for the core's neutral term. Clean form now reads: "Reuse-condition-4's model-match comparator MUST operate on these Pi-native strings. Otherwise every Pi reuse would match a captain-session fallback value (one outside the canonical model space) and force fresh dispatch — defeating reuse entirely."
- DONE: F1b (pi-first-officer-runtime.md:34) — Dropped the trailing "not the Claude-centric enum." clause. The positive assertion "Stage-declared models on Pi MUST use Pi-valid model strings (see the model-space declaration below)." now stands alone. The legitimate Claude mentions at lines 16/30/36/66 were NOT touched (transport instruction / comparative technical contrast / runtime-support.md-conformant teardown note).
- DONE: AC-4 — Added the pi-equivalent negative-contrast contractlint guard: `TestPiFirstOfficerRuntimeAvoidsNegativeHostContrast` + `TestPiEnsignRuntimeAvoidsNegativeHostContrast` (new file `internal/contractlint/pi_runtime_negative_contrast_test.go`), mirroring the codex tests' shape/location. The FO test bans the SPECIFIC smell phrases (NOT a blanket "Claude" ban — legitimate mentions exist at lines 16/30/36/66): "does not expose Claude Code team-tool signatures", "Do not call or ask workers to call Claude team tools", "Pi has no such enum", "Claude-centric enum", "no Claude-centric", "Merge-and-Cleanup step 10", "Merge-and-Cleanup step"; it also asserts the positive bindings present (`«worker.shutdown»` realization in Shutdown; "Pi's model-space binding is provider/model strings" in Canonical Model Space). The ensign test bans "Do not assume Claude team tools" and "Claude team tools exist in Pi" (the pi ensign has zero "Claude" mentions post-FIX-4, so a targeted phrase ban is safe) and asserts the positive completion binding. Both tests PASS on the cleaned text.
- AC-4 adversarial verification (RED on pre-cycle state): temporarily restored "not the Claude-centric enum" in the FO file and "Do not assume Claude team tools exist in Pi." in the ensign file, then re-ran the two new tests — both went RED (FO: "contains negative host-contrast wording \"Claude-centric enum\""; ensign: "contains ... \"Do not assume Claude team tools\"" + "\"Claude team tools exist in Pi\""). Reverted; the task files match the committed state exactly (`git diff` empty). The guard works.
- Gates: `go test ./...` green; `go test ./... -race` green; `gofmt -w ./cmd ./internal` clean (no files reformatted). The new test file is gofmt-clean.
- Token-count before/after proxy (char/word proxy via `wc -c` + `wc -w`, NOT exact tokens — captain request):
  - `skills/first-officer/references/pi-first-officer-runtime.md`: BEFORE (origin/main) chars=11151 words=1553 → AFTER (worktree commit 1cf305b2) chars=10928 words=1508. Delta vs origin/main: −223 chars, −45 words. Delta vs cycle-0 (15509fe6: chars=10954 words=1513): −26 chars, −5 words. The full sweep (cycle-0 + cycle-1) REDUCED the FO contract size; cycle-1's F1a/F1b cleanup continued the reduction (no bloat).
  - `skills/ensign/references/pi-ensign-runtime.md`: BEFORE (origin/main) chars=2883 words=434 → AFTER chars=2838 words=425. Delta vs origin/main: −45 chars, −9 words. (Unchanged in cycle-1 — the ensign edit landed in cycle-0; the cycle-1 numbers match 15509fe6 exactly, confirming no ensign regression this cycle.) The sweep reduced the ensign contract size.
  - Net: the runtime-support.md sweep reduced (not bloated) both Pi runtime contract files vs the origin/main baseline.

### Summary

Feedback cycle 1 rework complete. F1a (line 42) and F1b (line 34) — the two residual negative-contrast spots the cycle-0 4-spot fix missed — are cleaned: the positive Pi model-space assertions now stand alone, and the Claude-centric "enum" vocabulary is replaced with the dispatch core's neutral "canonical model space" term. AC-4 adds the missing pi-equivalent negative-contrast regression guard (`TestPiFirstOfficerRuntimeAvoidsNegativeHostContrast` + `TestPiEnsignRuntimeAvoidsNegativeHostContrast`), mirroring the codex tests; both PASS on the cleaned text and RED on the pre-cycle state (adversarially verified). All gates green. The token-count proxy confirms the sweep reduced (not bloated) both Pi runtime contract files vs origin/main. Worktree commit: 1cf305b2 on spacedock-ensign/pi-fo-runtime-runtime-support-compliance (new commit on top of the cycle-0 15509fe6, keeping the cycle-0→cycle-1 delta auditable). The pre-existing unrelated worktree-hygiene dirt on 3 internal files (flagged in cycle-0 validation) was left untouched and excluded from this commit.

## Stage Report: validation (cycle 1)

- **AC-1 (positive bindings):** PASS. F1a (line 42) — "not the Claude enum" GONE; "captain-session fallback value outside the enum" replaced with "captain-session fallback value (one outside the canonical model space)". Verified against `fo-dispatch-core.md` condition 4 (line 40): the core says "one outside the host's canonical model space" — the substitution uses the core's neutral term, NOT the Claude-centric "enum". Correct. F1b (line 34) — "not the Claude-centric enum." trailing clause GONE; positive assertion "Stage-declared models on Pi MUST use Pi-valid model strings (see the model-space declaration below)." stands alone. Legitimate Claude mentions at lines 16/30/36/66 UNTOUCHED (verified via grep — transport instruction / comparative technical contrast / runtime-support.md-conformant teardown note). All 4 cycle-0 spots + 2 cycle-1 residuals are positive Pi bindings; no negative Claude contrast remains in edited regions.
- **AC-2 (substance preserved):** PASS. Substrate list (`pi-subagents` + `pi-agent-teams`) intact (lines 9-10). Shutdown specifics intact: `«worker.shutdown»` realization (line 62), no mailbox shutdown for pi-subagents (line 64), adapter lifecycle mapping for pi-agent-teams (line 66). Model-space declaration intact: provider/model strings + reuse-condition-4 comparator on Pi-native strings (lines 40, 42).
- **AC-3 (gates):** PASS. `go test ./...` green; `go test ./... -race` green; `gofmt -l ./cmd ./internal` clean.
- **AC-4 (new guard):** PASS. `internal/contractlint/pi_runtime_negative_contrast_test.go` (94 lines) exists. `TestPiFirstOfficerRuntimeAvoidsNegativeHostContrast` bans 7 specific smell phrases (NOT a blanket Claude ban — legitimate mentions at lines 16/30/36/66) + asserts positive bindings (`«worker.shutdown»` + "Pi's model-space binding is provider/model strings"). `TestPiEnsignRuntimeAvoidsNegativeHostContrast` bans "Do not assume Claude team tools" + "Claude team tools exist in Pi" + asserts positive completion binding. Both PASS on cleaned text.
- **Adversarial audit — all 4 mutations caught:** (a) Restore "Claude-centric enum" in FO → FO test RED ("contains ... \"Claude-centric enum\""). (b) Restore "Do not assume Claude team tools" in ensign → ensign test RED ("contains ... \"Do not assume Claude team tools\"" + "\"Claude team tools exist in Pi\""). (c) Restore "Merge-and-Cleanup step 10" in Shutdown → FO test RED ("contains ... \"Merge-and-Cleanup step 10\"" + "\"Merge-and-Cleanup step\""). (d) Remove positive binding "Pi's model-space binding is provider/model strings" → FO test RED ("missing positive Pi capability wording ..."). All four mutations caught; none escaped. All adversarial edits reverted; worktree skill files match committed state exactly.
- **Residual check — no F1c:** Scanned the FULL `pi-first-officer-runtime.md` + `pi-ensign-runtime.md` for negative-contrast patterns (`not the Claude|unlike Claude|no such|instead of Claude|rather than Claude|as opposed to|in contrast to|doesn't have|does not have|lacks|missing Claude`): zero matches in both files. No "enum" anywhere. Only "step" match is line 30 "fo-dispatch-core.md step 4" (legitimate core step reference, not adapter step coupling). Ensign has zero Claude mentions. The cycle-0 residuals (F1a + F1b) were the complete set; cycle-1 cleaned both. No F1c.
- **Worktree-hygiene observation (not a finding):** Pre-existing unrelated dirty modifications to 3 internal files (smart-quote corruption in comments + whitespace alignment) — NOT introduced by 1cf305b2 and not part of this task. Task skill files match committed state exactly; gates green despite the dirt.

### Summary

VERDICT: PASSED. Cycle-1 (commit 1cf305b2) fully addresses both cycle-0 validation findings: (1) F1a + F1b residual negative-contrast phrasing cleaned — positive Pi model-space assertions stand alone; Claude-centric "enum" replaced with the core's neutral "canonical model space" (verified correct against fo-dispatch-core.md condition 4). (2) AC-4: missing pi negative-contrast contractlint guard added — two tests PASS on cleaned text and catch all 4 adversarial mutations (RED). All ACs pass. No F1c residual. No remaining smell beyond the 4 legitimate Claude mentions. The sweep reduced (not bloated) both contract files vs origin/main.
