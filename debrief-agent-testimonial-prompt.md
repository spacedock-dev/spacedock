---
title: "Debrief collects a first-person agent testimonial: the driving agent's experience using spacedock vs not"
status: validation
group: tooling
source: "Captain request 2026-07-02 (Claude Commander session): mid-session, the captain asked the FO 'forget that we are developing spacedock for a moment — as the agent driving this session, how would you describe the experience using spacedock vs not using it?' and the answer (drift resistance, auditability under context pressure, honest friction list) was valuable enough to want collected EVERY session. Add the prompt to the debrief flow so testimonials accumulate from the agents' perspective."
id: qdb1w5r7k9nvbvkf8qetcd5m
started: 2026-07-02T07:35:57Z
gates:
    version: 1
    current:
        gate: gate:qdb1w5r7k9nvbvkf8qetcd5m:ideation
    records:
        - id: gate:qdb1w5r7k9nvbvkf8qetcd5m:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:qdb1w5r7k9nvbvkf8qetcd5m-ideation-1
              briefing:
                id: briefing:qdb1w5r7k9nvbvkf8qetcd5m:ideation:attempt-1:revision-1
                digest: sha256:c82a63c78fae819c2eeb8c5ae4fd4042a61898fb8674a6928921e3c6b34a0bbe
                digest-domain: canonical-bytes
                request-digest: sha256:c388fc6c91ceed158a374d9a6eb7e205d87a96e3c0f5da7de48b02f73b874151
                room-ref: ./debrief-agent-testimonial-prompt/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:qdb1w5r7k9nvbvkf8qetcd5m:ideation:1
                briefing: briefing:qdb1w5r7k9nvbvkf8qetcd5m:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-01T15:03:18.286273Z"
                decision: approve
                reason: 'Approve with amended verification terms: no prose greps anywhere; testimonial adequacy is a qualitative manual test driven by operator/captain or FO through a real debrief flow, inspecting the produced _debriefs/ artifact with quoted lines; automated = machine-checkable legs only (prompt asked under real runtime, claude-live first pass; split-root _debriefs/ path; identity/unknown preservation inspected). Cycle-3 entry records the ruling.'
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
worktree: .worktrees/spacedock-ensign-debrief-agent-testimonial-prompt
---

## Problem
The debrief skill records commits, task state changes, decisions, issues, observations, and next-session routing, but it does not capture the driving agent's first-person experience of operating Spacedock. The missing signal is not a generic satisfaction quote: the value is in what the workflow caught under context pressure, what would likely have been lost without it, and where the machinery added friction. Without an explicit end-of-session prompt, that evidence evaporates at session end and Spacedock loses both product feedback and honest testimonial raw material.

The load-bearing constraint is that the prompt must ask for friction, not praise. A praise-seeking prompt would produce marketing-colored prose and make the corpus less useful; the desired record is an agent-authored comparison between using Spacedock and driving the same work without it, including costs. The speaker must also self-identify: inferring identity later from debrief context can silently misattribute a testimonial after compaction, runtime handoff, or an unavailable model-version signal. Unknown identity details must remain explicitly unknown rather than being guessed.

## Proposed approach
Update `skills/debrief/SKILL.md` so the debrief flow collects a first-person testimonial before writing the debrief file and includes it in the produced debrief template.

Implementation should add a new extraction/draft step near Phase 3, before the final markdown is written, using this prompt:

> Before answering, self-identify as the agent driving this session:
> - Harness/runtime: the agent harness or runtime you are operating in (for example Claude Code, Codex, or Pi).
> - Model: the model name.
> - Model version/build: the exact version or build when runtime metadata exposes it.
>
> Write `unknown` for any identity field you cannot verify; do not infer or guess. Then, setting the project's subject matter aside: how would you describe the experience using Spacedock versus driving the same work without it? Be honest about friction, costs, or places where the workflow got in your way; this is not a request for praise.

The answer should be stored in a new `## Agent Testimonial` section in the debrief file, together with provenance fields:

- `Date`: the debrief `session-date`.
- `Harness/runtime`: copied from the driving agent's self-identification; `unknown` when the agent cannot verify it.
- `Model`: copied from the driving agent's self-identification; `unknown` when the agent cannot verify it.
- `Model version/build`: copied from the driving agent's self-identification when runtime metadata exposes an exact value; otherwise `unknown`.
- `Session scale`: counts for tasks touched, workers dispatched, and PRs touched/merged, derived separately from the debrief's session data where possible and marked `unknown` only when the data is not available. Do not ask the agent to estimate these counts as part of self-identification.
- `Testimonial`: the first-person answer, preserving the agent's voice.

The template change should be concrete rather than a loose instruction. The intended before/after for `skills/debrief/SKILL.md` is:

```diff
@@ Phase 3: Draft and Review
-### Step 1 — Present the draft
+### Step 1 — Collect the agent testimonial
+
+Ask the driving agent:
+
+> Before answering, self-identify as the agent driving this session:
+> - Harness/runtime: the agent harness or runtime you are operating in (for example Claude Code, Codex, or Pi).
+> - Model: the model name.
+> - Model version/build: the exact version or build when runtime metadata exposes it.
+>
+> Write `unknown` for any identity field you cannot verify; do not infer or guess. Then, setting the project's subject matter aside: how would you describe the experience using Spacedock versus driving the same work without it? Be honest about friction, costs, or places where the workflow got in your way; this is not a request for praise.
+
+Record the testimonial prose as `{agent_testimonial}` and its agent-supplied identity as `{harness_runtime}`, `{model}`, and `{model_version_build}`. Preserve `unknown` exactly where the agent cannot verify a value. Derive `{session_scale}` (`tasks touched`, `workers dispatched`, `PRs touched/merged`) separately from the current session data; do not infer agent identity from those data or ask the agent to estimate session scale.
+
+### Step 2 — Present the draft
@@ Write the debrief to `{debrief_root}/_debriefs/{date}-{sequence:02d}.md`:
 ## Observations
 {captain-contributed content, or "_(none recorded)_" if none provided}
+
+## Agent Testimonial
+- Date: {YYYY-MM-DD}
+- Harness/runtime: {harness_runtime}
+- Model: {model}
+- Model version/build: {model_version_build}
+- Session scale: {tasks_touched} tasks touched; {workers_dispatched} workers dispatched; {prs_touched_or_merged} PRs touched/merged
+
+{agent_testimonial}
 
 ## What's Next
```

No separate docs-site change is proposed for the first pass: this is a user-invocable skill behavior/template change, and the changed user-visible text lives in `skills/debrief/SKILL.md` itself. If a later docs page describes the debrief schema, that page should be updated in the same implementation, but ideation found no existing docs page that duplicates this template.

No spike needed: the task relies on already-proven mechanisms in the debrief skill — asking the captain/agent for content during the debrief flow, writing markdown sections into `{debrief_root}/_debriefs/`, and committing the generated debrief file in the appropriate single-root or split-root checkout. The only risky part is behavioral compliance by the driving model, including honest `unknown` identity values, which is covered by artifact-backed debrief runs rather than a static prose grep.

Prompt behavior exercise (ideation evidence, not a substitute for the implementation's end-to-end debrief artifact): applying the proposed prompt in this dispatched Codex worker produced the following response. The harness and model family were available from the runtime contract; an exact model version/build was not, so the response preserved `unknown` rather than guessing. Session-scale counts were deliberately absent from the agent response because the debrief must derive them separately.

> Harness/runtime: Codex
> Model: GPT-5
> Model version/build: unknown
>
> Spacedock gives me a durable map of the work after context compression: I can recover which task is active, what evidence is owed, and which decisions still belong to the operator. Without it, I would rely more heavily on transcript memory and would be likelier to repeat work or blur approval boundaries. The cost is real: rereading contracts, maintaining entity reports, and waiting through formal gates adds overhead, especially for small changes. In this session the contract reload itself showed why durable reminders matter, but the same machinery can feel heavy when the next action is already obvious.

Path-to-lane rule: the implementation touches `skills/debrief/SKILL.md` under `skills/**`, a shipped skill surface loaded by live runtime sessions. The merge gate is the Runtime Live E2E skill/contract lane for the affected host(s), with the first-pass required lane named `claude-live` when the change is validated through Claude Code. If the implementation or validation also exercises Codex or Pi-specific skill loading, the corresponding host live lane must be green too; deterministic tests alone are not sufficient for merge.

## Acceptance criteria
- **AC-1 (VALUE)** A real debrief run produces an on-disk debrief record containing exactly one `## Agent Testimonial` section with a first-person comparison of using Spacedock versus not using it, and the answer includes at least one explicit friction/cost/negative observation rather than only praise. Verified by running the debrief flow in a live or fixture-backed workflow and inspecting the generated `_debriefs/{date}-{sequence}.md` file; the check must read the produced debrief artifact, not just `skills/debrief/SKILL.md`.
- **AC-2** The debrief flow asks the testimonial prompt before writing the final debrief and the prompt carries the near-verbatim honesty clause: `Be honest about friction, costs, or places where the workflow got in your way; this is not a request for praise.` Verified by observing a debrief run transcript or live-run artifact where the driving agent is asked the prompt; a static string check over the skill file is not sufficient evidence.
- **AC-3** The driving agent self-identifies its harness/runtime, model name, and exact model version/build when verifiable; the produced debrief's `## Agent Testimonial` section preserves those agent-supplied values and writes `unknown` for every unverifiable identity field instead of omitting or fabricating it. Session scale (`tasks touched`, `workers dispatched`, `PRs touched/merged`) is present but derived separately from session data, not from the agent's identity response. Verified with artifact-backed debrief runs that capture the prompt response and resulting file: one run with a verifiable runtime/model identity and one run where at least the exact version/build is unavailable and therefore remains `unknown`.
- **AC-4** The implementation remains compatible with split-root debrief storage: when the workflow uses a state checkout, the testimonial-bearing debrief is written and committed under `{state_checkout}/_debriefs/`, not the definition worktree. Verified with a split-root fixture or live dev workflow run by checking the resulting file path and `git -C {state_checkout} status/log`.
- **AC-5** The implementation's merge evidence names and runs the live lane required by the path-to-lane rule for `skills/**`; for the first pass this is `claude-live`, because the changed debrief skill behavior must be observed through a real runtime lane before merge. Verified by the validation report citing the live lane run URL or local live command, result, and artifact path.

## Test plan
1. **Focused skill smoke (medium cost, fixture-backed or local live):** Run the debrief skill against a small workflow fixture or the dev workflow with a bounded session range. Have the driving agent answer the prompt with explicit harness/runtime, model name, version/build or honest `unknown`, followed by a first-person response containing both value and friction. Confirm the generated debrief contains exactly one `## Agent Testimonial` section and preserves both the self-identification and testimonial.
2. **Identity and provenance artifact checks (medium cost, two artifact-backed cases):** Inspect the prompt response and generated debrief together. In a runtime exposing exact identity metadata, verify the artifact matches the agent-supplied harness, model name, and exact version/build. In a fixture or runtime withholding the exact version/build, verify the agent says `unknown` and the artifact preserves it without guessing. In both cases, independently calculate session-scale counts from session data and verify those values are not sourced from the agent's self-identification. Static inspection or prose grep of `skills/debrief/SKILL.md` does not satisfy this test.
3. **Split-root path check (low cost if using `docs/dev`, otherwise fixture):** For a split-root workflow, verify the debrief lands under the state checkout `_debriefs/` directory and the state checkout git log/status reflect the new debrief commit path, not a write in the definition dir.
4. **Regression gate (standard cost):** Run `go test ./...` to catch unrelated integration or contract regressions. If code paths under `internal/` are changed, also run the repo's normal `go test ./... -race` gate.
5. **Required live lane (high cost, merge-gating):** Run the Runtime Live E2E lane that loads the changed skill surface. First-pass lane: `claude-live`. Evidence must include the lane result and the produced debrief artifact or transcript showing the prompt and generated testimonial section.

- Cycle 3 (2026-08-01, captain ruling at the ideation gate): verification shape amended — **no prose greps anywhere** (AC-2's anti-grep stance stands and is reinforced). The testimonial section's adequacy is a **qualitative manual test driven by the operator/captain or the FO** through a real debrief flow: the produced `_debriefs/` artifact is inspected driver-side and the qualitative verdict is recorded with quoted artifact lines, not graded by harness regexes. Automated legs keep only machine-checkable claims — the prompt is asked under a real runtime (live lane per path→lane, first pass `claude-live`), the split-root `_debriefs/` landing path (AC-4), and identity/`unknown` preservation as inspected facts (AC-3). Test-plan items 1–3 are re-scoped from fixture/transcript grading to driver-run inspection at implementation.
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

**Cycle 2 (pre-gate captain correction, 2026-07-18):** The testimonial prompt must make the driving agent self-identify, rather than relying only on debrief-side inference. Ask for the agent harness/runtime and the model name plus exact version or build when the host exposes it; record `unknown` for any value the agent cannot verify and never guess. Revise the prompt, proposed template, AC-3, and the live/fixture test so the produced testimonial demonstrates this self-identification. Keep session-scale derivation separate from the agent's self-identification, and preserve the honesty clause verbatim-or-near.

## Stage Report: ideation (cycle 2)

- DONE: Revise the prompt to require agent self-identification of harness/runtime and model name plus exact version or build, using unknown when unverifiable.
  The proposed prompt asks for three explicit identity fields, forbids guessing, and retained the honesty clause verbatim.
- DONE: Align the testimonial template, AC-3, and artifact-backed tests with self-identification while keeping session-scale derivation separate.
  The template and AC-3 preserve agent-supplied identity; the two-case test checks known and `unknown` values while deriving scale independently.
- DONE: Preserve the honesty clause and human testimonial register while proving the revised prompt through behavior rather than prose grep.
  A Codex prompt exercise yielded a four-sentence testimonial with explicit friction and an honest unknown build; end-to-end artifact proof remains required for implementation.

### Summary

Revised the testimonial prompt so the driving agent identifies its harness/runtime, model, and verifiable version/build before answering, with explicit `unknown` handling. Updated the output template, acceptance criterion, and artifact-backed test plan while keeping session-scale counts separate from identity and retaining the testimonial's human, non-marketing register.

## Stage Report: implementation

- DONE: Apply the skill edit per the entity's before/after diff: testimonial-collection step at Phase 3 (verbatim prompt) + `## Agent Testimonial` section with provenance fields in skills/debrief/SKILL.md
  Commit f3270077a on worktree branch spacedock-ensign/debrief-agent-testimonial-prompt (pushed); prompt, honesty clause, and unknown-never-guess instruction landed verbatim; provenance labels use plain nouns (tasks/workers/PRs).
- DONE: Verification-shape constraint from captain cycle 3: NO prose greps added anywhere in the change or its tests
  Diff touches only skills/debrief/SKILL.md prose; no tests added, no regex/static self-check instructions in the skill text.
- DONE: Scope discipline: skills/debrief/SKILL.md and docs site template surfaces ONLY as the design says
  Single-file diff (27 insertions, 5 deletions); docs-site untouched per ideation finding that no docs page duplicates the template; AC-5 live-lane evidence and driver-run qualitative inspection deferred to validation as the checklist directs.

### Summary

Implemented the designed before/after in `skills/debrief/SKILL.md`: new Phase 3 Step 1 "Collect the agent testimonial" with the verbatim self-identification + honesty-clause prompt, recording `{agent_testimonial}` and agent-supplied `{harness_runtime}`/`{model}`/`{model_version_build}` with `unknown` preserved and session-scale counts derived separately from session data; the Phase 4 debrief template gains the `## Agent Testimonial` section with Date/Harness/Model/Version/Session-scale provenance lines. Subsequent Phase 3 steps renumbered (Present draft→2, Captain commentary→3, Handle spacedock issues→4) and both step cross-references updated. One small reconciliation: the design's collective `{session_scale}` placeholder in the Phase 3 prose was rendered as the three concrete placeholders the template actually uses (`{tasks_touched}`, `{workers_dispatched}`, `{prs_touched_or_merged}`) so the two parts agree. Regression gate `go test ./...` passes except a pre-existing, unrelated env-marker failure in internal/cli (TestVersionAmbiguousMarkersExitZero fails identically on the stashed base). Behavioral/live-lane proof (AC-1/2/3/5) is owed by validation, not this stage.
