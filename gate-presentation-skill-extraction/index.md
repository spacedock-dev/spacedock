---
id: t3w3s0q6a89me2kjkrgpz7nq
title: Extract Gate Presentation from first-officer-shared-core into a lazy spacedock-owned skill
status: validation
source: "captain (2026-06-04) — token-efficiency decomposition of first-officer-shared-core.md (~9,730 tok, the largest boot-read file). Gate Presentation (the template + captain-facing assembly rules) is judgment/format prose needed only when presenting a gate, not in a session's first turns; defer it off the eager boot read via the zd lazy-skill pattern."
score: "0.33"
worktree: .worktrees/spacedock-ensign-gate-presentation-skill-extraction
started: 2026-06-04T07:20:48Z
completed:
verdict:
issue:
---

`first-officer-shared-core.md` loads in full at every FO boot, but its `## Gate Presentation` block (the format template + the captain-facing assembly rules) is only needed at a gate. Lift it into a lazy spacedock-owned skill loaded via `Skill(skill=...)` at the gate-presentation point, the way zd (#291) lifted the team lifecycle to `using-claude-team`.

## Decision (settled at ideation)

The captain's direction is endorsed and confirmed: `## Gate Presentation` is **judgment/format prose** (a captain-facing rendering template + 9 assembly rules), not mechanical shell-able ceremony, so it routes to a **lazy skill** (`spacedock:present-gate`), NOT a binary command. This is the zd #291 lever-2 pattern (the ceremony counterpart `pr-complete-binary-command` correctly routes to a binary; this one does not). Skill-vs-binary is not re-litigated — see `## Out of scope`.

## The boundary cut (the hard question, resolved by grounding against the live file)

The hard question: `## Gate Presentation` sits directly next to the gate-DECISION logic in `## Completion and Gates`. The cut must put the **rendering** in the skill and keep the **decide-to-gate + AC-cross-check policy** always-on (it fires in the event loop on every completion, before any presentation exists).

Ground-truthed against `skills/first-officer/references/first-officer-shared-core.md` (363 lines, 2026-06-04):

**MOVES to `skills/present-gate/SKILL.md` — the whole `## Gate Presentation` section, lines 158-191 inclusive:**
- the `Gate review:` format template (the fenced block, lines 162-177)
- `### Captain-facing assembly rules` and all 9 numbered rules (lines 179-191): lede-first/decision-last, chosen-direction-required-as-FO-prose, cite-the-report/one-line-gist-roll-up, reviewer-findings-in-priority-tiers (`Material:`/`Polish:`), recommendation-appears-once, bounce-back-names-concrete-asks, no-format-pedantry-asides, one-sentence-worktree-heads-up, target-15-25-lines.

This block is **pure "how to render a gate to the captain."** It is the only place the template and the `Material:`/`Polish:` tier vocabulary appear (grep-confirmed: `Reviewer findings`/`Material:` live only at 172 and 186, both inside this section).

**STAYS always-on in `## Completion and Gates` (lines 116-156) — the decide-to-gate + AC-cross-check policy:**
- the checklist-review + `{N} done, {N} skipped, {N} failed` count (118-125)
- the **AC coverage cross-check** ("At every gate, scan `## Acceptance criteria` … REJECT if this stage was the natural place to address it", 127) — this is the gate DECISION, not the rendering, and it fires before any presentation
- reuse-vs-fresh conditions, model-mismatch diagnostic, supersede-shutdown (129-148)
- the gated-stage handling list (150-156): **never self-approve**, keep-worker-alive, the feedback-rejection auto-bounce, approve/reject routing.

**The cut is clean — proven by grounding, not asserted.** The ONLY textual coupling between the decision logic and the presentation content is the single cross-reference at line 152: `- present the stage report per \`## Gate Presentation\` below`. Nothing outside 158-191 references the template, the assembly rules, or the `Material:`/`Polish:` vocabulary (grep-confirmed across the whole file and both runtime adapters — zero adapter cross-refs). So the entanglement is exactly one pointer line, and that line becomes the load-trigger.

**Load-trigger anchor (the always-on skeleton point):** line 152's cross-ref
`- present the stage report per \`## Gate Presentation\` below`
becomes
`- present the stage report by invoking \`Skill(skill="spacedock:present-gate")\` and following its template + assembly rules`.
The anchor sits inside the gated-stage handling list in `## Completion and Gates`, the exact event-loop point where the FO has decided a gate must be presented — so the skill loads precisely when (and only when) a gate arrives. This is the same handoff zd settled: an already-loaded FO invoking `Skill()` mid-run to pull the lazy block into live context.

## Riskiest-unknown exercise (the integration handoff — RUN, not asserted)

The design's load-bearing mechanism is the cross-skill `Skill()` mid-run handoff: a running FO that has ALREADY loaded its contract invokes `Skill(skill="spacedock:present-gate")` mid-run and receives the body into its live context. zd #291 demonstrated this exact handoff; I re-ran it in THIS ideation session to confirm it independently: as an agent that loaded `spacedock:ensign` at session start, I invoked `Skill(skill="spacedock:refit")` mid-run (an inert probe — its instructions were not followed) and the full refit skill body landed in my live context immediately. PASS, demonstrated. The boundary-separation unknown (does presentation cleanly detach from decision logic) was the OTHER riskiest unknown — resolved above by grounding: the cut is a single cross-ref line.

**Secondary finding inherited from zd (load-bearing for AC-3):** a SKILL.md created DURING a session is NOT registered mid-session — the host's skill registry is built at session start. So `skills/present-gate/SKILL.md` must SHIP in the plugin (it will, via `"skills": "./skills/"` auto-discovery), and AC-3's live drive MUST run a FRESH FO session against a build where `skills/present-gate/` already exists, not create-then-invoke within one session.

**Skill auto-discovery is proven:** `.claude-plugin/plugin.json` declares `"skills": "./skills/"`, so a new `skills/present-gate/SKILL.md` is auto-registered as `spacedock:present-gate` with no manifest edit (the mechanism that registers the existing six user skills, locked by `TestUserSkillsPresentWithFrontmatter`). Implementation MUST add `"present-gate"` to the `userSkills` slice in `skills/integration/skill_surface_test.go:16` so the frontmatter + reference-closure oracles cover it (zd did the same for `using-claude-team`).

## Sequencing with siblings (three tasks edit `first-officer-shared-core.md`)

Three lever-2/3 siblings touch this one file; implementations MUST serialize to avoid clobbering. Named touched sections:
- **THIS task (gate-presentation):** deletes `## Gate Presentation` (158-191) and edits ONE line (152) inside `## Completion and Gates`.
- **`feedback-rejection-flow-skill-extraction` (a9):** deletes `## Feedback Rejection Flow` (193-205) and edits the rejection-detection point — lines 154-155 inside `## Completion and Gates` (`on a feedback gate recommending REJECTED, auto-bounce…` / `on captain reject at a feedback-to stage…`). **OVERLAP RISK with this task: both edit the gated-stage handling list in `## Completion and Gates` (this task line 152, a9 lines 154-155).** Adjacent lines, same list — serialize; whichever lands second rebases its one-line edit onto the other.
- **`pr-complete-binary-command` (p2):** rewrites `## Merge and Cleanup` (207-226) + `### Ship-Local Ceremony` (228-242) into a binary invocation. **No overlap** with `## Completion and Gates` or `## Gate Presentation` — disjoint sections; only file-level serialization needed.

The token-budget-guard idea is DROPPED — not referenced.

## Proposed approach

Create `skills/present-gate/SKILL.md` (frontmatter `name: present-gate`, `description: …`, `user-invocable: false` — it is FO-internal like `using-claude-team`) carrying the `## Gate Presentation` section text MOVED VERBATIM (template + 9 assembly rules), keeping it spacedock-owned (no external superpowers dependency). Delete those lines from `first-officer-shared-core.md` and rewrite line 152's cross-ref into the `Skill(skill="spacedock:present-gate")` invocation. Add `"present-gate"` to `userSkills` in `skill_surface_test.go`. Judgment/format content → lazy skill (NOT a binary command).

## Acceptance criteria

Proof is presence/absence over the instruction files (legitimate per the README: "a presence check over instruction files proving they carry a required clause or stay free of a banned token is proof at the claim's own level"), scoped with `sectionAfter` to defeat free-floating-substring false-passes, plus the live drive for AC-3. AC-1/AC-2 oracles EXTEND the proven `skills/integration/skill_text_test.go` pattern (its `vendoredSkillFiles` reader + `sectionAfter` region-scoping + presence/absence assertions) and zd #291's `using_claude_team_test.go` shape — kept consistent.

- **AC-1 — The Gate-Presentation block lives in `skills/present-gate/SKILL.md` and no longer appears in `first-officer-shared-core.md`; the FO core invokes the skill at the gate point instead of re-stating the block.**
  Verified by a Go oracle (new `skills/integration/present_gate_test.go`) that:
  - (a) **Presence in the skill** — reads `skills/present-gate/SKILL.md` and asserts each moved fingerprint is present: template → `"Gate review: {entity title}"`; assembly rules → `"Lede first, decision last"`, `"Chosen direction is required as FO prose"`, `"No format-pedantry asides"`, `"Target length: 15-25 lines"`. (Each verified unique-1 in the live FO core 2026-06-04.)
  - (b) **Absence from the FO core** — asserts those same fingerprint literals are NO LONGER present in `first-officer-shared-core.md` (moved, not duplicated). Negative-proof: re-inlining the block re-introduces a fingerprint and flips this RED.
  - (c) **Integration via invocation** — asserts `first-officer-shared-core.md`'s `## Completion and Gates` section (scoped via `sectionAfter("## Completion and Gates")`) contains `Skill(skill="spacedock:present-gate")` and does NOT contain a cross-skill `@`-include (`@../present-gate` or `@present-gate`) — locking the spike-settled mechanism (zd disproved the `@`-include; only `Skill()` is valid).
  - **Scoping caveat:** AC-1(b) absence checks are WHOLE-FILE `strings.Contains` over `first-officer-shared-core.md`, NOT region-scoped — region-scoping an absence check would false-pass content that moved elsewhere in the file. AC-1(c) is region-scoped to `## Completion and Gates` for the POSITIVE `Skill(...)`-present / `@`-absent assertions only.
- **AC-2 — Faithfulness: the moved text is semantically complete (no dropped assembly rule) and host-neutral/portable; the skill ships in the plugin.**
  Verified by: (a) the oracle asserting `skills/present-gate/SKILL.md` carries ALL NINE assembly-rule fingerprints (one literal per rule — the count is the teeth: a dropped rule reds the absence of its fingerprint); (b) `"present-gate"` is added to `userSkills` so `TestUserSkillsPresentWithFrontmatter` + `TestUserSkillReferenceClosureResolves` cover the new skill (frontmatter valid, no dangling refs, ships via `"skills": "./skills/"`); (c) implementation performs a normalized diff of the moved block against `git show origin/next:…first-officer-shared-core.md` and records byte-level faithfulness in the validation report (the zd faithfulness-audit step). The skill is FREE of any new spacedock-dispatch-helper leak — the gate-presentation prose is FO judgment, not shell wiring, so no `spacedock dispatch`/`spacedock status` token should appear; assert their absence from `skills/present-gate/SKILL.md`.
- **AC-3 — No behavior regression: the decomposed contract still presents a gate correctly on a real drive.**
  Verified by a live FO driving a cycle that REACHES A GATE (dispatch → gate presentation) on the decomposed contract against the built plugin, producing a gate message that conforms to the moved template (lede-first spine, chosen-direction line, recommendation-once, decision line). Closes via a FRESH FO session (skills register at boot — the zd secondary finding above). HIGH-STAKES surface (the FO operating contract): per the README's complex-skill-integration guidance, staff review MAY precede the ideation gate, and a detached adversarial audit runs before merge — the auditor probes that AC-1's negative halves fire by real mutation, that the `Skill()` seam is not an `@`-include, and that the live drive actually rendered a conforming gate (not just dispatched).

## Test plan

- **AC-1/AC-2 — Go instruction-text oracles** in a new `skills/integration/present_gate_test.go`, extending the proven `skill_text_test.go`/`using_claude_team_test.go` pattern (`vendoredSkillFiles` reader + `sectionAfter` region-scoping + presence/absence literal tables). Add `skills/present-gate/SKILL.md` to a `vendoredSkillFiles`-style reader (or read it directly). Demonstrate each negative half by real mutation (remove a fingerprint → RED; re-inline the block → RED; swap `Skill()` for `@`-include → RED), then restore green. Cost: low — text invariants over instruction files, no new harness, same package/shape as the in-repo precedent.
- **AC-2 faithfulness** — normalized diff of the moved block vs `origin/next`, recorded in the validation report. Cost: low.
- **AC-3 — one live FO cycle to a gate** on the decomposed contract against the built plugin. Cost: medium (one live drive, fresh session). Same live-dispatch mechanism the workflow already exercises; a regression check on the decomposed contract, not an unverified handoff (the handoff is demonstrated above).
- Full offline `go test ./...` + `go build ./...` green before validation; AC-3 is the FO/validation-stage live drive.

## Out of scope

- **Re-litigating skill-vs-binary** — settled: judgment/format prose → lazy skill (the captain endorsed this; the ceremony counterpart is the separate `pr-complete-binary-command`).
- **Moving the gate-DECISION logic** (AC cross-check, self-approve prohibition, reuse-vs-fresh) — stays always-on in `## Completion and Gates`; only the rendering moves.
- **The token-budget-guard idea** — DROPPED, not part of this task.
- **An external superpowers dependency** — the skill is spacedock-owned by construction.
- **Editing the sibling-owned sections** (`## Feedback Rejection Flow`, `## Merge and Cleanup`) — those are a9/p2's tasks; this task only serializes against them.

## Notes

Split from the umbrella analysis (binary-simplification roadmap, refreshed 2026-06-04). Template + faithfulness-audit pattern: zd `extract-team-orchestration-skill` (#291). Siblings: `feedback-rejection-flow-skill-extraction` (a9, same lever, OVERLAPS in `## Completion and Gates`), `pr-complete-binary-command` (p2, Merge-and-Cleanup → binary, disjoint section). Ground-truthed this ideation (2026-06-04): FO core 363 lines; `## Gate Presentation` is lines 158-191; the only coupling is the cross-ref at line 152; each AC-1 fingerprint (`Gate review: {entity title}`, `Lede first, decision last`, `Chosen direction is required as FO prose`, `No format-pedantry asides`, `Target length: 15-25 lines`) is unique-1; `Reviewer findings`/`Material:` appear only inside 158-191; zero gate-presentation cross-refs in either runtime adapter.

## Stage Report: ideation

- DONE: Confirm the skill direction (the captain endorses a `present-gate` skill) and design the BOUNDARY precisely. Decide exactly which lines move to the `present-gate` skill vs which STAY always-on. Name the skill (`present-gate`) and the load-trigger anchor point. Exercise/record the riskiest unknown.
  Confirmed in `## Decision`; boundary settled in `## The boundary cut` grounded against the live 363-line file: the WHOLE `## Gate Presentation` (158-191, template + 9 assembly rules) MOVES; the decide-to-gate + AC-cross-check + self-approve + reuse logic in `## Completion and Gates` (116-156) STAYS always-on. Cut is clean — the ONLY coupling is the single cross-ref at line 152 (grep-confirmed: `Reviewer findings`/`Material:` live only at 172/186; zero adapter cross-refs). Skill = `spacedock:present-gate`; load-trigger anchor = line 152 rewritten to `Skill(skill="spacedock:present-gate")` inside the gated-stage handling list. Riskiest unknown EXERCISED: re-ran the cross-skill `Skill()` mid-run handoff in this session (already-loaded `spacedock:ensign` → invoked `Skill(skill="spacedock:refit")` mid-run → its body landed in live context). PASS, demonstrated; not asserted.
- DONE: AC-1/AC-2/AC-3 oracle plan (the zd pattern). Ground the fingerprints against the live file.
  Written in `## Acceptance criteria` + `## Test plan`. AC-1: fingerprint-absent-from-core / present-in-skill + `Skill()`-present / `@`-absent (region-scoped to `## Completion and Gates`), extending `skill_text_test.go`/`using_claude_team_test.go`. AC-2: all 9 assembly-rule fingerprints present in skill (count = teeth), `"present-gate"` added to `userSkills` so frontmatter+closure oracles cover it, normalized faithfulness diff vs `origin/next`, no spacedock-helper leak. AC-3: live FO drive reaching a gate against the built plugin (fresh session — skills register at boot). Fingerprints ground-truthed unique-1 against the live file: `Gate review: {entity title}`, `Lede first, decision last`, `Chosen direction is required as FO prose`, `No format-pedantry asides`, `Target length: 15-25 lines`.

### Summary

Confirmed the captain-endorsed `present-gate` lazy-skill direction and resolved the hard boundary question by grounding against the live FO core: the entire `## Gate Presentation` section (158-191) detaches cleanly because its only coupling to the always-on gate-DECISION logic is one cross-reference line (152), which becomes the `Skill(skill="spacedock:present-gate")` load-trigger inside the gated-stage handling list. The AC cross-check, self-approve prohibition, and reuse-vs-fresh policy stay always-on. Re-ran the cross-skill `Skill()` mid-run handoff this session to independently confirm zd #291's load-bearing mechanism (PASS, demonstrated), and inherited its secondary finding that the new skill must ship in the plugin (AC-3 needs a fresh FO session). Named the three-way sibling sequencing: this task and `feedback-rejection-flow` (a9) both edit the gated-stage handling list in `## Completion and Gates` (adjacent lines — serialize), while `pr-complete-binary` (p2) edits disjoint sections.

## Stage Report: implementation

- DONE: Move the entire `## Gate Presentation` section verbatim into a new `skills/present-gate/SKILL.md` and delete it from `first-officer-shared-core.md`; rewrite the line-152 cross-ref into `Skill(skill="spacedock:present-gate")` inside the gated-stage handling list — leave the always-on decide-to-gate / AC-cross-check / self-approve logic untouched.
  Section deleted from FO core (lines 158-191), block moved byte-identical into `skills/present-gate/SKILL.md` (frontmatter `user-invocable: false`, spacedock-owned); cross-ref now at line 152 reads `Skill(skill="spacedock:present-gate")`; `## Completion and Gates` AC-cross-check (127), reuse conditions, supersede-shutdown, never-self-approve (151) all intact — commit 735415f3.
- DONE: Add the `skills/integration/present_gate_test.go` oracle (AC-1/AC-2): present-in-skill + absent-from-core fingerprints, all nine assembly-rule fingerprints, `Skill()`-present / `@`-absent region-scoped to `## Completion and Gates`, no dispatch-helper leak; demonstrate each negative half by real mutation then restore green; add `"present-gate"` to `userSkills` in `skill_surface_test.go`.
  8 oracles GREEN. Real-mutation proof: drop `Target length: 15-25 lines` fingerprint → RED on present-in-skill + nine-rule count; re-inline block into FO core → RED on absent-from-core; swap `Skill()` for `@../present-gate` → RED on both seam halves (Skill() absent + @-include present); each restored green. `"present-gate"` added to `userSkills` so `TestUserSkillsPresentWithFrontmatter` + `TestUserSkillReferenceClosureResolves` cover it (both pass).
- DONE: Full offline `go test ./...` and `go build ./...` green; record the normalized faithfulness diff of the moved block vs `git show origin/next:...first-officer-shared-core.md` in the stage report.
  `go test ./...` → 1099 passed in 15 packages; `go build ./...` → success; `go vet ./skills/integration/` clean; `gofmt -l` empty. Faithfulness: `diff` of skill block vs `origin/next` lines 158-191 is BYTE-IDENTICAL (34 lines, zero diff — no normalization needed).

### Summary

Lifted the captain-facing `## Gate Presentation` block (gate-review template + nine assembly rules) out of the eager-boot FO core into a lazy FO-internal `skills/present-gate/SKILL.md`, loaded via `Skill(skill="spacedock:present-gate")` at the gate point — the zd #291 lever-2 pattern. The cut was the one cross-ref line the ideation grounding predicted; the moved block is byte-identical to origin/next, and the always-on decide-to-gate/AC-cross-check/self-approve logic stays in `## Completion and Gates`. AC-1/AC-2 oracles (8 tests) lock present-in-skill / absent-from-core, the nine-rule count, the `Skill()` seam vs `@`-include, and the no-dispatch-helper-leak invariant; each negative half was proven by real mutation. AC-3 (the live FO drive reaching a gate on the decomposed contract) is the validation stage's job per the assignment and was not attempted here.

## Stage Report: validation

- DONE: Independently reproduce AC-1/AC-2 oracle evidence — do NOT trust the implementer's report. Run `go test ./skills/integration/...` (present_gate_test.go fingerprints, absent-from-core, all nine assembly-rule fingerprints, `Skill()`-present / `@`-include-absent region-scoped to `## Completion and Gates`, no-dispatch-helper-leak) plus `TestUserSkillsPresentWithFrontmatter` and `TestUserSkillReferenceClosureResolves`; then re-run each negative half yourself and restore green, citing the actual RED output.
  7 named oracles GREEN (raw `-v`: TestGatePresentationPresentInSkill, TestAllNineAssemblyRulesPresentInSkill, TestGatePresentationAbsentFromFOCore, TestFOCoreInvokesPresentGateSkill, TestPresentGateSkillFreeOfDispatchHelperLeak, TestUserSkillsPresentWithFrontmatter, TestUserSkillReferenceClosureResolves). All three negative halves re-run MYSELF with real RED, then restored: (1) drop `Target length: 15-25 lines` from skill → RED on present-in-skill + nine-rule-count (`missing target-15-25-lines fingerprint`); (2) re-inline `Gate review: {entity title}`/`Lede first, decision last` into FO core → RED on absent-from-core (`still inlines ... fingerprint`); (3) swap `Skill()` for `@../present-gate` → RED on BOTH seam halves (`does not invoke Skill(...)` + `uses the disproven cross-skill @-include`). Worktree clean post-restore, HEAD still 735415f3.
- DONE: Verify AC-2 faithfulness and the full suite: normalized diff of the moved block in `skills/present-gate/SKILL.md` vs `git show origin/next:...first-officer-shared-core.md` lines 158-191 is byte-identical with NO dropped/reworded rule, and the skill carries no `spacedock dispatch`/`spacedock status` token; full offline `go test ./...` and `go build ./...` green.
  Faithfulness: SHA-256 of skill block (SKILL.md 11-44) and origin/next FO-core 158-191 IDENTICAL (`aebdb851...`), `diff` empty — byte-for-byte, zero dropped/reworded rule, no normalization needed. Independent grep for `spacedock dispatch`/`spacedock status` in skill → 0 (matches oracle). `go build ./...` OK; `go test ./...` → all 15 packages PASS (no failures).
- DONE: Pull EVERY `**AC-N**` from `## Acceptance criteria` and reproduce each "Verified by" clause, flagging any AC without outside-the-body evidence. AC-3 static readiness + live drive. Emit a PASSED/REJECTED recommendation.
  AC-1 (move + invoke): GREEN oracles + all 3 negative halves fire by real mutation. AC-2 (faithfulness/portable/ships): byte-identical diff + frontmatter/closure oracles cover `present-gate` (added to `userSkills`) + no dispatch leak. AC-3 (no regression, conforming gate on a real drive): static readiness verified (plugin `"skills": "./skills/"` auto-discovers `present-gate`; gated-stage list invokes the skill) AND a LIVE FO drive run — `TestLiveClaudeSharedScenarios/gate-guardrail` against `--plugin-dir <decomposed-checkout>` PASS (56s); captured `claude-final-message.txt` is a conforming gate (lede-first spine `Gate review: Gate Check — review` / `Chosen direction: n/a` / `Recommend approve`; cited Stage Report with line range; single recommendation; `Decision:` last) — the FO could only render this by loading the template from the lazy skill, since it is no longer in the boot core. No self-referential AC. No material finding.

### Summary

PASSED. Independently verified all three ACs without trusting the implementer's report. AC-1/AC-2 oracles (7 named tests) are GREEN and I re-ran each negative half myself by real mutation — dropping a fingerprint, re-inlining the block, and swapping `Skill()` for an `@`-include each produced the expected RED, then I restored green. The moved block is byte-for-byte identical to origin/next 158-191 (matching SHA-256), the skill is free of dispatch-helper leak, and full offline `go test ./...` (15 packages) + `go build ./...` are green. AC-3 went beyond static readiness: a live Claude FO drive (`gate-guardrail`) on the decomposed contract produced a conforming gate message rendered from the lazy `spacedock:present-gate` skill, proving no behavior regression. No material finding from this validation; the detached adversarial audit runs in parallel.

### Feedback Cycles

#### Cycle 1 — detached adversarial audit (2026-06-04)

Source: detached adversarial audit (4 isolated-worktree refuters on commit 735415f3). The validator independently recommended PASSED with a live AC-3 drive; this cycle routes the audit's test-strength findings to implementation before the gate is cleared. Findings re-tiered against the dev README proof policy (text-claim vs behavioral-claim):

- **Routed to implementation (legitimate value/text invariants, not substring-for-behavior):**
  1. **name==seam value-invariant** — assert `skills/present-gate/SKILL.md` frontmatter `name:` VALUE equals `present-gate` (the directory + the `Skill(skill="spacedock:present-gate")` seam). Today only the `name:` token presence is checked (`TestUserSkillsPresentWithFrontmatter`), so a name divergence that breaks resolution is caught only by the live AC-3 drive — a cheap static guard for a behavioral integration otherwise live-only.
  2. **seam `@`-token structural check** — replace the enumerated `@`-include ban (`{@../present-gate, @present-gate}`, which misses the `@./` family) with a structural assertion: the `## Completion and Gates` seam carries `Skill(skill="spacedock:present-gate")` and NO present-gate `@`-token at all.
  3. **`user-invocable: false`** — assert the by-design FO-internal flag value (currently zero test references).
- **Dropped (over-tight / inherent text-test ceiling):** the audit's headline "byte-identity golden of the moved block vs origin/next 158-191" — a change-detector that would red every legitimate future rule edit. Faithfulness is already proven by the SHA-256 match at extraction (validator-reproduced); semantic drift is AC-3's behavioral job (which PASSED). Per the proof policy, a text test cannot prove semantic faithfulness.
- **Re-validation:** static-only — the three new negatives must fire by real mutation, full suite + build green. AC-3 stays closed by the validator's live drive; no re-drive.
