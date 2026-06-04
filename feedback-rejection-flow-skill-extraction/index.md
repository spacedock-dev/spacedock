---
id: a9nte184whfmz8ajzn4n51yr
title: Extract Feedback Rejection Flow from first-officer-shared-core into a lazy spacedock-owned skill
status: implementation
source: "captain (2026-06-04) — token-efficiency decomposition of first-officer-shared-core.md. The Feedback Rejection Flow is needed only when a gate rejects (or a feedback stage recommends REJECTED), not in a session's first turns; defer it off the eager boot read via the zd lazy-skill pattern."
score: "0.31"
worktree: .worktrees/spacedock-ensign-feedback-rejection-flow-skill-extraction
started: 2026-06-04T07:20:48Z
completed:
verdict:
issue:
---

`first-officer-shared-core.md`'s `## Feedback Rejection Flow` (cycle tracking, route-to-feedback-to-target, 3-cycle escalation, budget-probe/reuse, re-run-reviewer) only fires on a rejection. Lift its PROCEDURE BODY into a lazy spacedock-owned skill loaded via `Skill(skill="spacedock:feedback-rejection-flow")` at the rejection-handling point (zd #291 pattern); the rejection DETECTION stays in the always-on skeleton.

## End value

The FO boot read no longer eagerly carries the ~11-line feedback-rejection procedure that only fires when a gate rejects. The procedure loads on demand at the rejection point. Token savings on the common path (a session with no rejection never loads it); faithfulness preserved so a rejection still routes exactly as today.

## Decision: the boundary (what moves vs what stays)

**MOVES → `skills/feedback-rejection-flow/SKILL.md`** (the spacedock-owned lazy skill): the 7-step procedure body, `first-officer-shared-core.md:195-205` — the `When a feedback stage recommends REJECTED:` lead-in, steps 1-7 (read `feedback-to` target; track cycles; 3-cycle escalate; budget-probe consult; route-to-target-in-worktree-else-fresh including the Codex `send_input` non-completion caveat; re-run reviewer; re-enter normal gate flow), and the trailing `The FO owns ### Feedback Cycles. Routing follows FO Write Scope…` paragraph (205).

**STAYS in the always-on skeleton** (the rejection DETECTION + the machinery the flow only references by name, never duplicates):
- The two detection lines in `## Completion and Gates` (154-155): `on a feedback gate recommending REJECTED, auto-bounce…` and `on captain reject at a feedback-to stage, enter the Feedback Rejection Flow (priority over generic rejection)`. These are the load-trigger anchor — the event loop detects the rejection here. They get the `Skill(skill="spacedock:feedback-rejection-flow")` invocation (replacing the bare prose "auto-bounce into the feedback rejection flow" / "enter the Feedback Rejection Flow").
- `### Feedback Cycles` OWNERSHIP — defined by Worktree Ownership (264), FO Write Scope (298, 302), and Probe/Ideation Discipline (358). The flow WRITES to `### Feedback Cycles`; it does not own the section's write-scope rules. Those stay always-on (they govern more than feedback).
- Reuse conditions + budget probe (131-138). The flow says "Consult the budget probe (reuse condition 0)" and "reuse conditions pass" — it REFERENCES them. Reuse conditions fire on every non-feedback stage advance too, so they stay always-on. The skill cites them by name.

**Skill name:** `feedback-rejection-flow` → invoked `Skill(skill="spacedock:feedback-rejection-flow")`. Plugin-owned (lives under `skills/`, auto-discovered by `"skills": "./skills/"` in plugin.json — same mechanism that registers `using-claude-team`). NO external superpowers dependency. Unlike `using-claude-team` (which is generic, zero-spacedock), this skill is INTENTIONALLY spacedock-specific — it names `feedback-to`, `### Feedback Cycles`, reuse conditions. So the zd AC-2 "no-spacedock-leak" banned-token table does NOT apply here; the inverse holds (the skill SHOULD carry the spacedock procedure verbatim).

**Trigger anchor:** the rejection-detection bullets in `## Completion and Gates` (154-155). The single-entity-mode auto-bounce in `claude-first-officer-runtime.md:161` ("REJECTED with feedback-to → auto-bounce … subject to the 3-cycle limit") also detects-and-enters the flow; it references the flow by behavior, not by re-stating steps, so it stays as-is (it points at the same now-lazy procedure).

## Decision: the runtime-adapter bare-mode variant STAYS as an adapter seam

`claude-first-officer-runtime.md:163-167` `## Feedback Rejection Flow (bare mode)` is a Claude-runtime realization seam (sequential dispatch in bare mode; keep-reviewer-alive in teams mode) — it parallels how the zd extraction left runtime seams (`## Degraded Mode (spacedock seams)`) inline in the adapter while the generic body moved to the skill. The Codex side has NO feedback-rejection-specific section; `codex-first-officer-runtime.md:43-48 ## Reuse And Feedback Routing` is generic `send_input` reuse routing, not a feedback-cycle variant. So:
- The Claude bare-mode section STAYS in `claude-first-officer-runtime.md` as an adapter seam. It is the teams-vs-bare execution detail of the now-lazy shared procedure; moving it into the spacedock-owned skill would re-couple a Claude-only execution mode into the cross-runtime procedure body. The shared-core skill stays runtime-neutral (the step-5 routing already names "send_input on Codex, SendMessage on Claude teams" generically).
- The skill's step-5 routing prose moves verbatim (it is already runtime-neutral). The Claude bare-mode sequencing seam stays in the Claude adapter; AC-2 asserts it stays consistent (still present, still says "sequential … dispatch fix agent (wait), then reviewer (wait), then gate" / "keep the reviewer alive").

## Riskiest unknown — exercised (entanglement check)

**Question:** does the flow extract cleanly, or is it welded into the always-on gate/event-loop control flow such that a dropped clause mis-routes a rejection?

**Exercised against the live file (ground-truthed 2026-06-04, `first-officer-shared-core.md` = 363 lines):**
- The procedure body is a self-contained 7-step block (193-205). It REFERENCES three pieces of always-on machinery BY NAME without duplicating any of them: (a) `### Feedback Cycles` — 6 occurrences in the file, of which only 2 (198, 205) are inside the flow; the other 4 (264, 298, 302, 358) are the section's write-scope/ownership rules that stay always-on. (b) reuse condition 0 / budget probe — defined at 131-138, the flow only says "Consult the budget probe (reuse condition 0)". (c) the detection bullets at 154-155 inside `## Completion and Gates`. None of these is COPIED into the flow body, so moving 195-205 leaves the referenced machinery intact and the references resolve to the always-on skeleton.
- **Cut point confirmed clean:** lines 195-205 are pure procedure; lines 154-155 are pure detection. The cut is procedure-vs-detection, exactly the zd "generic body moves, the invocation seam stays" shape. No clause spans both — the heading line 193 + the detection bullets 154-155 carry the `Skill()` invocation; 195-205 is the body. Verified `Skill()` cross-skill mid-run load is the production-proven mechanism (zd #291 demonstrated already-loaded-A → `Skill(B)` → B-in-context; the FO already loads `Skill(skill="spacedock:using-claude-team")` and `Skill(skill="spacedock:first-officer")` this exact way).
- **Conclusion:** extracts cleanly. The flow is NOT entangled with the always-on loop — it references always-on machinery by name; the detection (which IS in the loop) stays. NO spike beyond this ground-truth exercise is needed: the `Skill()` cross-skill load is production-proven (zd #291), and the cut is a text relocation over a verified-clean boundary. Recorded per ideation discipline: relies on the proven `Skill()` cross-skill composition path and the proven `skill_text_test.go` oracle pattern.

## Acceptance criteria

Each AC names a property of the finished entity. Proof is presence/absence over the instruction files (legitimate per the README: a presence/absence check over instruction files is proof at the claim's own level), `sectionAfter`-scoped to defeat free-floating-substring false-passes, extending the in-repo `skills/integration/using_claude_team_test.go` + `skill_text_test.go` pattern (`vendoredSkillFiles` reader + `sectionAfter` + presence/absence literal tables), plus the live drive for AC-3.

**AC-1 — The feedback-rejection PROCEDURE lives in `skills/feedback-rejection-flow/SKILL.md` and no longer appears inlined in `first-officer-shared-core.md`; the FO core's rejection-detection point invokes the skill rather than re-stating the procedure.**
Verified by a new Go oracle (`skills/integration/feedback_rejection_flow_test.go`) that:
- (a) **Presence in the skill** — reads `skills/feedback-rejection-flow/SKILL.md` (added to a `vendoredSkillFiles`-style reader) and asserts the procedure is present by unique-fingerprint literals, each ground-truthed unique-1 in the pre-change shared core: `On cycle 3, escalate to the human instead of another round`, `Re-run the reviewer after fixes`, `Re-enter the normal gate flow with the updated result`, `The FO owns \`### Feedback Cycles\``.
- (b) **Absence from the shared core** — asserts those same fingerprint literals are NO LONGER present in `first-officer-shared-core.md` (the procedure moved, not duplicated). Whole-file `strings.Contains` (NOT region-scoped — region-scoping a generic-absence check false-passes content moved elsewhere). Negative-proof: re-inlining the procedure re-introduces a fingerprint and flips this RED.
- (c) **Integration via invocation at the detection point** — asserts `first-officer-shared-core.md`'s `## Completion and Gates` section (scoped via `sectionAfter`) contains `Skill(skill="spacedock:feedback-rejection-flow")` and does NOT contain a cross-skill `@`-include (`@../feedback-rejection-flow`, `@feedback-rejection-flow`) — locking the proven mechanism, banning the zd-disproven `@`-include.
- (d) **Always-on machinery retained** — asserts `first-officer-shared-core.md` STILL contains the `### Feedback Cycles` write-scope/ownership rules (anchor: the FO Write Scope bullet literal `**\`### Feedback Cycles\` section**`) and the reuse-conditions/budget-probe block (anchor: `Reuse conditions`), so the referenced machinery did not move with the procedure.

**AC-2 — Faithfulness + the adapter bare-mode seam stays consistent.**
- The moved procedure is semantically complete: the skill body carries all seven steps including the load-bearing clauses — the 3-cycle escalation, the budget-probe consult (`reuse condition 0`), the route-to-target-in-worktree-else-fresh decision, the Codex `send_input` non-completion caveat (fingerprint: `do not treat the immediate \`send_input\` response as the new completion result`), and re-run-reviewer / re-enter-gate. Verified by the AC-1(a) fingerprint set PLUS a faithfulness check in the oracle asserting the `send_input` caveat and the `feedback-to` target-read clause (`the stage that receives the fix request, not the reviewer`) are present in the skill — these are the clauses whose loss would silently mis-route.
- The Claude adapter bare-mode seam STAYS consistent: assert `claude-first-officer-runtime.md` STILL contains `## Feedback Rejection Flow (bare mode)` and its sequential-dispatch sentence (fingerprint: `the feedback rejection flow is sequential`) and the keep-reviewer-alive sentence (`Keep the reviewer alive`). Negative-proof: removing the seam reds this.
- A normalized-diff faithfulness audit at implementation (`git show` the pre-change procedure block vs the skill body) confirms the text MOVED, not rewritten — same audit the zd validation stage ran. This is a validation-stage check, recorded here as the faithfulness bar.

**AC-3 — No behavior regression: a live reject → route → re-review cycle routes correctly on the decomposed contract.**
Verified by a live FO driving a feedback cycle that REJECTS at a `feedback-to` stage, routes the findings back to the target stage, and re-runs the reviewer — completing the cycle correctly on the contract where the procedure is lazy-loaded via `Skill()`. Requires a FRESH FO session against the built plugin (skills register at boot — the zd #291 secondary finding: a session-created skill is not registered mid-session). This is the HIGH-STAKES surface (FO operating contract); per the README complex-skill-integration guidance, a detached adversarial audit runs at validation before merge.

## Test plan

- **AC-1/AC-2 — Go instruction-text oracles** in a new `skills/integration/feedback_rejection_flow_test.go`, extending the proven `using_claude_team_test.go` shape (`vendoredSkillFiles` reader + a `feedbackRejectionFlowSkill` reader for the new SKILL.md + `sectionAfter` region-scoping + presence/absence fingerprint tables). Add `feedback-rejection-flow` to the `userSkills` list in `skill_surface_test.go` so it ships in the build (zd boot-registration finding). Negative-proof discipline: insert/remove each fingerprint, watch red/green. Cost: low — text invariants, no new harness, same package/shape as the in-repo precedent.
- **AC-3 — one live FO feedback cycle** (reject → route-to-feedback-to → re-review) on the decomposed contract, fresh FO session against the built plugin. Cost: medium (one live drive). The reject-route-rereview path is the highest-risk relocated behavior, so the live drive MUST reach the re-review, not just the reject.
- **High-stakes surface (the FO operating contract itself).** Per the README complex-skill-integration guidance: the FO requests an **independent staff review** before the ideation gate (boundary cleanliness, oracle sufficiency, the entanglement exercise was actually run) and a **detached adversarial audit before merge**. Auditor probes: (1) AC-1(b) — re-inline the procedure into the shared core, confirm the absence half reds; (2) AC-1(c) — replace the `Skill(...)` invocation with an `@`-include, confirm part (c) reds; (3) the live drive actually reached the re-review (the relocated route + re-run-reviewer path), not just the reject.

## Sequencing (serialize the shared-core edits)

Siblings `pr-complete-binary-command` and `gate-presentation-skill-extraction` also edit `first-officer-shared-core.md`. This task's implementation TOUCHES:
- `## Completion and Gates` (lines 154-155 — the two detection bullets get the `Skill()` invocation)
- `## Feedback Rejection Flow` (lines 193-205 — REMOVED, body relocated to the new skill)
- NEW FILE `skills/feedback-rejection-flow/SKILL.md`
- `claude-first-officer-runtime.md` `## Feedback Rejection Flow (bare mode)` — UNCHANGED (seam stays) but adjacent
- `skills/integration/feedback_rejection_flow_test.go` (new), `skill_surface_test.go` (add to `userSkills`)

`gate-presentation-skill-extraction` touches `## Gate Presentation` (158-191) — ADJACENT to but DISJOINT from this task's `## Completion and Gates` (154-155) and `## Feedback Rejection Flow` (193-205). The FO should serialize the two shared-core implementations to avoid a rebase conflict on the same file, but the edited SECTIONS do not overlap. Name the touched sections (above) so the FO can sequence by section, not just by file.

## Out of scope

- Moving the Claude bare-mode seam into the skill (decided against — it is a Claude-runtime execution mode, stays as an adapter seam).
- A Codex analog of the bare-mode section (Codex has no feedback-cycle-specific section; `## Reuse And Feedback Routing` is generic reuse routing, untouched).
- Changing the reuse-conditions / budget-probe block or the `### Feedback Cycles` write-scope rules (always-on machinery the flow references; stays).
- The token-budget-guard idea (DROPPED — not referenced).

## Notes

Split from the umbrella analysis (binary-simplification roadmap, refreshed 2026-06-04). Template: zd `extract-team-orchestration-skill` (#291) — same `Skill()` cross-skill load mechanism, same `skill_text_test.go`/`using_claude_team_test.go` oracle pattern, same staff-review-then-detached-audit high-stakes discipline. Siblings: `gate-presentation-skill-extraction`, `pr-complete-binary-command`. Ground-truthed 2026-06-04 against `first-officer-shared-core.md` (363 lines): `## Feedback Rejection Flow` at 193-205; detection bullets at 154-155; AC-1 fingerprints each unique-1 in the file; `### Feedback Cycles` ×6 (2 in-flow, 4 always-on); `## Feedback Rejection Flow (bare mode)` at `claude-first-officer-runtime.md:163-167`.

## Stage Report: ideation

- DONE: Design the extraction of `## Feedback Rejection Flow` into a lazy spacedock-owned skill + the load-trigger anchor; decide the boundary against what stays always-on; name the skill + trigger anchor; exercise the riskiest unknown (clean extraction vs entangled).
  Boundary firmed in `## Decision: the boundary`: procedure body (shared-core 195-205) MOVES to new `skills/feedback-rejection-flow/SKILL.md` (invoked `Skill(skill="spacedock:feedback-rejection-flow")`); rejection DETECTION (`## Completion and Gates` bullets 154-155) STAYS as the trigger anchor and carries the `Skill()` invocation; `### Feedback Cycles` write-scope + reuse-conditions/budget-probe stay always-on (the flow references them by name, never copies them). Bare-mode adapter variant decided to STAY as a Claude adapter seam (`## Decision: the runtime-adapter bare-mode variant STAYS`); Codex has no feedback-specific section. Riskiest unknown EXERCISED against the live 363-line file (`## Riskiest unknown`): cut is procedure-vs-detection, clean — extracts cleanly; `Skill()` cross-skill load is production-proven (zd #291); no further spike needed (recorded with the proven mechanisms).
- DONE: AC oracle plan (the zd pattern): AC-1 fingerprints ABSENT from shared-core + PRESENT in the new skill with the `Skill()` invocation at the rejection point; AC-2 faithfulness + adapter bare-mode consistency; AC-3 a live reject→route→re-review cycle. Ground fingerprints against the live file.
  `## Acceptance criteria` written: AC-1(a-d) presence/absence/invocation/retained-machinery oracle in a new `feedback_rejection_flow_test.go` extending `using_claude_team_test.go`; AC-2 faithfulness fingerprints (send_input caveat, feedback-to target-read clause) + the Claude bare-mode seam-stays-consistent checks + a normalized-diff move-not-rewrite audit; AC-3 live fresh-FO reject→route→re-review drive. All AC-1 fingerprints ground-truthed unique-1 in the live `first-officer-shared-core.md` (`On cycle 3, escalate to the human instead of another round`, `Re-run the reviewer after fixes`, `Re-enter the normal gate flow with the updated result`, `The FO owns ### Feedback Cycles`). Sequencing named by touched section so the FO serializes the shared-core edits against siblings.

### Summary

Firmed the spec rather than rewriting it. The boundary cut is procedure-vs-detection: the 7-step procedure body (shared-core 195-205) moves to a spacedock-OWNED lazy skill `feedback-rejection-flow`, the rejection-detection bullets (154-155) stay in the always-on skeleton and carry the `Skill(skill="spacedock:feedback-rejection-flow")` invocation. The entanglement exercise (run against the live 363-line file) showed clean extraction — the flow references `### Feedback Cycles` ownership, reuse conditions, and the budget probe BY NAME (none duplicated; the always-on machinery stays put), so relocating 195-205 leaves every reference resolving to the skeleton. Decided the Claude bare-mode variant stays as an adapter seam (not moved — it is a Claude-runtime execution mode; Codex has no feedback-specific section). The AC oracles extend the proven zd #291 `using_claude_team_test.go` pattern; unlike that generic skill, this one is intentionally spacedock-specific (no no-leak banned-token table). No spike beyond the ground-truth exercise: `Skill()` cross-skill load is production-proven and the cut is a text relocation over a verified-clean boundary.
