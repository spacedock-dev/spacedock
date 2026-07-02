---
title: "Debrief collects a first-person agent testimonial: the driving agent's experience using spacedock vs not"
status: ideation
group: tooling
source: "Captain request 2026-07-02 (Claude Commander session): mid-session, the captain asked the FO 'forget that we are developing spacedock for a moment — as the agent driving this session, how would you describe the experience using spacedock vs not using it?' and the answer (drift resistance, auditability under context pressure, honest friction list) was valuable enough to want collected EVERY session. Add the prompt to the debrief flow so testimonials accumulate from the agents' perspective."
id: qdb1w5r7k9nvbvkf8qetcd5m
started: 2026-07-02T07:35:57Z
---

## Problem
The debrief skill records commits, task state changes, decisions, issues, observations, and next-session routing, but it does not capture the driving agent's first-person experience of operating Spacedock. The missing signal is not a generic satisfaction quote: the value is in what the workflow caught under context pressure, what would likely have been lost without it, and where the machinery added friction. Without an explicit end-of-session prompt, that evidence evaporates at session end and Spacedock loses both product feedback and honest testimonial raw material.

The load-bearing constraint is that the prompt must ask for friction, not praise. A praise-seeking prompt would produce marketing-colored prose and make the corpus less useful; the desired record is an agent-authored comparison between using Spacedock and driving the same work without it, including costs.

## Proposed approach
Update `skills/debrief/SKILL.md` so the debrief flow collects a first-person testimonial before writing the debrief file and includes it in the produced debrief template.

Implementation should add a new extraction/draft step near Phase 3, before the final markdown is written, using wording near this prompt:

> Setting the project's subject matter aside: as the agent driving this session, how would you describe the experience using Spacedock versus driving the same work without it? Be honest about friction, costs, or places where the workflow got in your way; this is not a request for praise.

The answer should be stored in a new `## Agent Testimonial` section in the debrief file, together with provenance fields:

- `Date`: the debrief `session-date`.
- `Host/runtime`: the runtime host in use when known (for example Claude Code, Codex, Pi); if unknown, record `unknown` rather than inventing.
- `Model`: the driving model when known; if unknown, record `unknown`.
- `Session scale`: counts for entities touched, workers dispatched, and PRs touched/merged, derived from the debrief's session data where possible and marked `unknown` only when the data is not available.
- `Testimonial`: the first-person answer, preserving the agent's voice.

The template change should be concrete rather than a loose instruction. The intended before/after for `skills/debrief/SKILL.md` is:

```diff
@@ Phase 3: Draft and Review
-### Step 1 — Present the draft
+### Step 1 — Collect the agent testimonial
+
+Ask the driving agent:
+
+> Setting the project's subject matter aside: as the agent driving this session, how would you describe the experience using Spacedock versus driving the same work without it? Be honest about friction, costs, or places where the workflow got in your way; this is not a request for praise.
+
+Record the answer as `{agent_testimonial}`. Also record `{host_runtime}`, `{model}`, and `{session_scale}` (`entities touched`, `workers dispatched`, `PRs touched/merged`) from the current session context when known; use `unknown` only for fields that cannot be determined.
+
+### Step 2 — Present the draft
@@ Write the debrief to `{debrief_root}/_debriefs/{date}-{sequence:02d}.md`:
 ## Observations
 {captain-contributed content, or "_(none recorded)_" if none provided}
+
+## Agent Testimonial
+- Date: {YYYY-MM-DD}
+- Host/runtime: {host_runtime}
+- Model: {model}
+- Session scale: {entities_touched} entities touched; {workers_dispatched} workers dispatched; {prs_touched_or_merged} PRs touched/merged
+
+{agent_testimonial}
 
 ## What's Next
```

No separate docs-site change is proposed for the first pass: this is a user-invocable skill behavior/template change, and the changed user-visible text lives in `skills/debrief/SKILL.md` itself. If a later docs page describes the debrief schema, that page should be updated in the same implementation, but ideation found no existing docs page that duplicates this template.

No spike needed: the task relies on already-proven mechanisms in the debrief skill — asking the captain/agent for content during the debrief flow, writing markdown sections into `{debrief_root}/_debriefs/`, and committing the generated debrief file in the appropriate single-root or split-root checkout. The only risky part is behavioral compliance by the driving model, which is covered by a live debrief smoke rather than a static prose grep.

Path-to-lane rule: the implementation touches `skills/debrief/SKILL.md` under `skills/**`, a shipped skill surface loaded by live runtime sessions. The merge gate is the Runtime Live E2E skill/contract lane for the affected host(s), with the first-pass required lane named `claude-live` when the change is validated through Claude Code. If the implementation or validation also exercises Codex or Pi-specific skill loading, the corresponding host live lane must be green too; deterministic tests alone are not sufficient for merge.

## Acceptance criteria
- **AC-1 (VALUE)** A real debrief run produces an on-disk debrief record containing exactly one `## Agent Testimonial` section with a first-person comparison of using Spacedock versus not using it, and the answer includes at least one explicit friction/cost/negative observation rather than only praise. Verified by running the debrief flow in a live or fixture-backed workflow and inspecting the generated `_debriefs/{date}-{sequence}.md` file; the check must read the produced debrief artifact, not just `skills/debrief/SKILL.md`.
- **AC-2** The debrief flow asks the testimonial prompt before writing the final debrief and the prompt carries the near-verbatim honesty clause: `Be honest about friction, costs, or places where the workflow got in your way; this is not a request for praise.` Verified by observing a debrief run transcript or live-run artifact where the driving agent is asked the prompt; a static string check over the skill file is not sufficient evidence.
- **AC-3** The produced debrief's `## Agent Testimonial` section includes provenance fields for date, host/runtime, model, and session scale (`entities touched`, `workers dispatched`, `PRs touched/merged`), with unknown values explicitly written as `unknown` instead of omitted or fabricated. Verified by inspecting the generated debrief artifact from the same run used for AC-1.
- **AC-4** The implementation remains compatible with split-root debrief storage: when the workflow uses a state checkout, the testimonial-bearing debrief is written and committed under `{state_checkout}/_debriefs/`, not the definition worktree. Verified with a split-root fixture or live dev workflow run by checking the resulting file path and `git -C {state_checkout} status/log`.
- **AC-5** The implementation's merge evidence names and runs the live lane required by the path-to-lane rule for `skills/**`; for the first pass this is `claude-live`, because the changed debrief skill behavior must be observed through a real runtime lane before merge. Verified by the validation report citing the live lane run URL or local live command, result, and artifact path.

## Test plan
1. **Focused skill smoke (medium cost, fixture-backed or local live):** Run the debrief skill against a small workflow fixture or the dev workflow with a bounded session range, answer the new testimonial prompt with a first-person response that includes both value and friction, and confirm the generated debrief contains one `## Agent Testimonial` section with the preserved answer.
2. **Provenance artifact check (low cost, artifact inspection):** Inspect the generated debrief file from the smoke run and verify date, host/runtime, model, and session-scale fields are present. Use `unknown` only where the run truly cannot provide the value.
3. **Split-root path check (low cost if using `docs/dev`, otherwise fixture):** For a split-root workflow, verify the debrief lands under the state checkout `_debriefs/` directory and the state checkout git log/status reflect the new debrief commit path, not a write in the definition dir.
4. **Regression gate (standard cost):** Run `go test ./...` to catch unrelated integration or contract regressions. If code paths under `internal/` are changed, also run the repo's normal `go test ./... -race` gate.
5. **Required live lane (high cost, merge-gating):** Run the Runtime Live E2E lane that loads the changed skill surface. First-pass lane: `claude-live`. Evidence must include the lane result and the produced debrief artifact or transcript showing the prompt and generated testimonial section.

## Stage Report

### Summary
Refined the debrief testimonial idea into a concrete skill/template change: add an honesty-framed first-person prompt, store the answer in a new debrief section with provenance, require a produced debrief artifact as proof, and name the `claude-live` lane as the first-pass merge gate for the `skills/**` change.

### Completion checklist
- DONE — Prompt lands verbatim-or-near in the debrief skill flow WITH the honesty/friction clause (not praise) — captured as AC-2 and the proposed `skills/debrief/SKILL.md` diff.
- DONE — Debrief output template gains a testimonial section + provenance fields (date, host/model, session scale: entities/workers/PRs touched); a produced debrief record demonstrates it end-to-end — captured as AC-1, AC-3, and the test plan artifact checks.
- DONE — A concrete first-pass path-to-lane rule is identified (touches skills/**, so the live lane gates the merge) — named `claude-live` as the first-pass required live lane in Proposed approach and AC-5.

### Validation performed
- Ran `${SPACEDOCK_BIN:-spacedock} dispatch show-stage-def --workflow-dir /Users/clkao/git/spacedock-research/spacedock-v1/docs/dev --stage ideation` and used it as the stage contract.
- Read the entity file and `skills/debrief/SKILL.md` to anchor the proposal in the current debrief flow/template.
- Read the dev workflow proof policy/path-to-lane rule in `docs/dev/README.md` to identify the live lane requirement.

### Residual risk
The exact mechanism for deriving host/runtime, model, worker count, and PR count may need small implementation judgment from available session context. The acceptance criteria require explicit `unknown` rather than fabricated values when a field cannot be determined.

### Feedback Cycles

**Cycle 1 (pre-gate, captain-derived from live testimonial drafting, 2026-07-02):** The captain asked the FO to draft the proposed testimonial from its own session perspective, then revised the draft across three passes. The learnings below are guidance for the implementation stage — they constrain *how the prompt and section should read*, not the ACs (which stand).

1. **Length and register.** A good testimonial reads like a person, not a report: roughly 4–8 sentences, no internal workflow jargon ("ensign", "gate", "AC cross-check", run ids, literal commands, tool names). The reader is an operator running multiple agents who is not yet using spacedock; internal terms are uninformative to them.
2. **The honesty clause is load-bearing in practice, not just in spec.** Once the prompt forced naming friction, the genuine costs surfaced (a flat boot tax paid every cold session; a hard stop at approval that adds latency) instead of defaulting to praise. The prompt must keep the "not a request for praise" framing verbatim-or-near; a softer phrasing degrades the corpus.
3. **What the machinery caught that mattered.** The single load-bearing catch was *not letting the driving agent self-approve the stage transition* — the asymmetry between "the agent thought it was fine" and "an operator signed off" is where drift enters. Implementation/test plan should treat the testimonial as evidence of that asymmetry, not generic satisfaction.
4. **Provenance format that read cleanly** at the end of the drafts: `date · runtime host · model · {N} task(s), {N} worker(s), {N} gate(s), {N} PR(s)`. Plain nouns ("task"/"worker"/"gate"), not workflow-internal ones ("entity"/"ensign"). AC-3's field list is correct; the label naming in the produced section should follow this plain-noun form.
5. **A genuine runtime-adapter gap surfaced during this session** (not part of qd's scope, recorded for awareness): the FO contract assumes an `ensign` agent type the Pi adapter binds, but Pi exposes `worker`/`delegate`/`oracle`; the FO fell back to `worker` and hand-assembled assignment prose that `dispatch build` is meant to own. This does not change qd's ACs but is worth not fabricating away if the testimonial mechanism ever asks the driving agent about dispatch friction.

These do not change the five ACs or the test plan; they tighten how implementation renders the prompt, the section labels, and the provenance nouns.
