---
title: "Reuse-advance messages go through a dispatch-built file pointer instead of FO-hand-assembled verbatim stage sections"
status: ideation
group: tooling
source: "fable-token-trim-scout analysis 2026-07-02 (captain-ordered fresh-angle token review): the reuse path explicitly does NOT route through dispatch build (fo-dispatch-core.md:40) — the FO hand-assembles a SendMessage embedding the full README stage subsection verbatim (claude-fo-dispatch.md:40), costing a README section read PLUS a verbatim echo per advance, ~400-800 tok x ~10 advances in a 5-entity/3-stage session = 4-8k tok. The biggest recurring dispatch-path cost found; fresh dispatch already proves the file-pointer mechanism (~175-char prompt pointing at a written dispatch file)."
id: jhe1244c8cjdymfbnnpnsvsw
started: 2026-07-02T03:02:47Z
---

## Problem
Advancing a reused worker to its next stage is the only dispatch-shaped message the FO still assembles by hand: the contract requires copying the next stage's full `### Stage definition` subsection from the README verbatim into a SendMessage, plus the checklist and continuation instruction. Every advance pays the section read + echo twice over; initial dispatch already solved this shape with `dispatch build` writing the assignment to a file and emitting a short pointer prompt.

## Proposed approach

**Verb shape: a `--advance` flag on `dispatch build`, not a new verb.** Advance mode shares nearly all of build's machinery: flag/file input plumbing (`--entity-path`, `--stage`, `--checklist-file`, `--scope-notes-file`, `--feedback-context-file`, `--feedback-reflow`), the validation rules (stage exists, checklist non-empty, entity readable, worktree stickiness, feedback-reflow requires context), effective-model resolution, split-root state-commit guidance, fetch-command emission, and the collision-keyed dispatch-file write. The deltas are mode-conditional branches in `runBuildFields` (internal/dispatch/build.go), not a parallel assembly path.

**Advance-mode deltas from initial-dispatch build:**
1. **No first-action block.** The reused worker already holds its operating contract; the file opens with an advance header instead: `## Advancing to next stage: {stage}` + `You are continuing work on: {entity title}`.
2. **Continue-on-entity line** replaces the "Read the entity file ... for the full spec" wording: `Continue working on the entity at {entity_path}.` (split-root/worktree path resolution identical to build).
3. **Completion Signal block names the next stage** — `Done: {title} completed {next_stage}. Report written to {path}.` — and keeps the "after all commits and stage report writes are done" clause, so commit-before-signal rides the same block. (The current hand template leaves the reused ensign to derive the new Done: wording itself; the built file pins it.)
4. **Dispatch filename gets an `-advance` suffix** (`{key}-{workerKey}-{slug}-{next_stage}-advance.md`) so an advance file can never alias a fresh-dispatch file for the same slug+stage (a fresh dispatch after a failed advance would otherwise collide with the stale advance body).
5. **Envelope drops the spawn-only fields.** Output: `schema_version`, `description`, `fetch_commands`, `dispatch_file_path`, `prompt`, `model`. No `subagent_type`/`name`/`team_name`/`run_in_background` — nothing is spawned. `model` stays so reuse-condition-4's comparator can read `next_stage.effective_model` from the build output instead of a separate README read.
6. **`--advance` + `--bare-mode` is a usage error (exit 2).** A reuse advance presupposes an addressable worker; bare mode has none.
7. **Pointer message (the `prompt` field), per host:** Claude: `Advancing to next stage: {stage}.\n\nRead {dispatch_file_path} and treat its content as your next-stage assignment.` Codex: same wording (no skill-wrapper clause, matching the codex pointer today). Pi: n/a — reuse-advance is deferred on Pi (pi adapter: fresh redispatch is the default first slice).

**The target handle stays FO-side.** The envelope does NOT emit a `to`/worker-name field: the live handle may carry a `-cycleN` suffix the helper cannot know, and the FO's session roster is already the authoritative dead-vs-alive tracker (claude-fo-dispatch.md `## Context Budget and Dead Ensign Handling`). The FO sends `output.prompt` through the reuse-advance handle it already holds.

**Ensign-side bootstrap clause.** ensign-shared-core.md's `## DISPATCH_FILE Bootstrap` covers the initial prompt only; it gains an advance clause: when a mid-session message matches `Advancing to next stage: {stage}. ... Read /tmp/spacedock-dispatch/{name}.md and treat its content as your next-stage assignment.`, Read the file and treat its content as the next-stage assignment (fetch-commands bootstrap included); on Read failure, send `DISPATCH_FILE_MISSING: {path} - {error}` to team-lead and stop — same failure shape as the initial bootstrap.

**Unchanged:** the reuse conditions 0–4, the frontmatter advance (`status --set ... status={next_stage}`, `advance:` commit), supersede-shutdown, and the break-glass rule (`«dispatch.build»` block clause: non-zero exit → manual template; the current verbatim SendMessage template survives demoted to break-glass).

## Content contract (pinned before the helper design)

Every element the current reuse-advance template (claude-fo-dispatch.md:40) carries, and its carrier in the file/pointer shape:

| Current template element | Carrier in new shape |
|---|---|
| Next stage name (`Advancing to next stage: {next_stage_name}`) | Pointer message AND file header |
| `### Stage definition:` — full README subsection verbatim | Fetch line (`show-stage-def`) in the file's `### Fetch commands`; the ensign materializes the verbatim section via the proven Fetch-on-Demand bootstrap — identical mechanism to initial dispatch, which never inlines the section either |
| `### Completion checklist` (Dispatch step 2) | File `### Completion checklist` (via `--checklist-file`, same one-item-per-line discipline) |
| `Continue working on {entity title} at {entity_file_path}` | File header (title) + continue-on-entity line (path); split-root guidance carries the state-checkout commit discipline alongside |
| `Commit before sending your completion message` | File `### Completion Signal` block ("after all commits and stage report writes are done"), now also pinning the exact next-stage `Done:` wording |
| Feedback context when re-entering a feedback-to stage (feedback-rejection-flow step 5: the routed message "must carry the concrete next-stage assignment and fix work" — today that path has NO template) | `--feedback-context-file` → file `### Feedback from prior review`; `--feedback-reflow` keeps build's rule-5 validation (missing context = error) |

Content parity is at the "what the ensign ends up with" level: the stage definition arrives verbatim via `show-stage-def`, exactly as every initial dispatch already delivers it.

## Spike record (riskiest mechanism)

Riskiest path: a live, already-completed worker receiving a mid-session pointer-shaped message and acting on the pointed file instead of an inline body. Spiked 2026-07-02 in this session:

1. Spawned a throwaway worker (haiku) with a file-pointer prompt (`Read {scratchpad}/spike-stage1.md and treat its content as your assignment`); it wrote `spike-out-stage1.txt` = `stage1 done` and completed.
2. Sent the completed worker (resumed by agent handle) the production-shaped advance: `Advancing to next stage: stage2.\n\nRead {scratchpad}/spike-advance-stage2.md and treat its content as your next-stage assignment.`
3. Observed on-disk proof: the resumed worker read the advance file and wrote `spike-out-stage2.txt` = `stage2 done`.

Composite evidence: pointer-following at spawn is proven by every shipped dispatch (the ~175-char v2 prompt); `SendMessage(to=name)`-to-a-live-worker delivery is the production reuse-advance channel today (with inline bodies); the spike closes the remaining gap — pointer-following mid-session, post-completion. (Harness note: teammates cannot spawn named background workers, so the spike resumed a completed synchronous subagent by agentId; the message-content mechanism under test is identical.)

## Contract-prose touch points (before → after)

1. **skills/first-officer/references/fo-dispatch-core.md `## Reuse and Fresh Dispatch`, "If reuse" (line ~40).** Before: "Send the next assignment through the runtime adapter's reuse-advance handle (its live-worker messaging call) — carrying the next stage name, the full `### Stage definition` subsection copied from the README verbatim, the `### Completion checklist` from Dispatch step 2, and an instruction to keep working on the entity at its path and commit before signaling. The reuse path does NOT route through `«dispatch.build»` — assemble the advancement message directly." After: "Build the advancement with `«dispatch.build»` in advance mode (`--advance`, same checklist-file discipline; `--feedback-context-file` when routing rejection findings) and send the emitted `prompt` through the runtime adapter's reuse-advance handle (its live-worker messaging call). On non-zero helper exit only, fall back to the adapter's manual advance template (the break-glass rule)."
2. **fo-dispatch-core.md `«dispatch.build»` closing line (~136).** Before: "`«dispatch.build»` serves initial dispatch only; the reuse-advance path assembles its message directly (`## Reuse and Fresh Dispatch`)." After: "`«dispatch.build»` serves initial dispatch (spawn envelope) and reuse advance (`--advance`: a pointer message for the reuse-advance handle, no spawn fields)."
3. **fo-dispatch-core.md `«dispatch.build»` effect block (~120-131).** Add `[--advance]` to the flag listing and one sentence on the advance-mode envelope (no `subagent_type`/`name`/`team_name`/`run_in_background`; `prompt` is the advance pointer).
4. **skills/first-officer/references/claude-fo-dispatch.md `## Spawn Call (Agent)` reuse-advance handle (lines ~38-40).** Before: the hand-assembled `SendMessage(to="{agent}-{slug}-{completed_stage}", message="Advancing to next stage: ...[STAGE_DEFINITION — copy the full ### stage subsection from the README verbatim]...")` template. After: "On a zero-exit `spacedock dispatch build --advance ...`, send `SendMessage(to="{live worker handle from session roster}", message=output.prompt)`." The current verbatim template moves under the break-glass heading (non-zero exit only), gaining the completion-signal line for the next stage.
5. **skills/ensign/references/ensign-shared-core.md `## DISPATCH_FILE Bootstrap`.** Add the advance clause (mid-session pointer → Read → next assignment; Read failure → `DISPATCH_FILE_MISSING`, stop).
6. **skills/first-officer/references/codex-first-officer-runtime.md `«addressable-worker»` line (~20).** Half-line addition: the `followup_task(target,message)` advance payload is `dispatch build --advance` `output.prompt`.
7. **Pi: no change.** pi-first-officer-runtime.md defers reuse-advance (fresh redispatch is the default first Pi slice); recorded as out of scope.

## Acceptance criteria

- **AC-1 (VALUE, measured): the per-advance FO-side message payload is O(pointer), invariant to stage-definition size.** For the dev workflow's stages, the emitted advance `prompt` is ≤ 300 bytes for every stage, while the hand-assembled template scales with the README section. Measured baseline (2026-07-02, `spacedock dispatch show-stage-def --workflow-dir docs/dev`): ideation section 4,866 bytes, validation 3,139, implementation 1,259 — vs a 239-byte pointer, i.e. ≥ 5–20x per advance before adding checklist/framing, and the FO additionally stops paying the README-section read it performed only to echo it. Test: a Go unit test builds `--advance` over fixture workflows with small and large stage sections and asserts (a) `prompt` ≤ 300 bytes and byte-identical length across section sizes (modulo stage-name/path length), and (b) the old-template assembly (section + checklist + framing, computed in-test from the same fixture) exceeds the emitted prompt by ≥ 5x on the smallest fixture stage. The baseline moves the wrong way if the pointer ever embeds the section.
- **AC-2 (content contract): the advance file + pointer carry every element in the pinned content-contract table** — next stage name, stage definition via `show-stage-def` fetch line, completion checklist, entity title/path + continue-on-entity instruction, next-stage `Done:` completion block with the commits-first clause, and feedback context + reflow validation when re-entering. Test: golden fixtures for (i) plain advance, (ii) split-root advance (state-commit guidance present), (iii) feedback-reflow advance, plus the rule-5 error case (`--feedback-reflow` without context file); a parity assertion walks the table's elements against the emitted file body + prompt.
- **AC-3 (mechanism shipped, paired with AC-1/AC-4): the shipped contract routes the reuse path through the helper.** fo-dispatch-core.md's "If reuse" names `--advance` and the "does NOT route through `«dispatch.build»`" sentence is gone; claude-fo-dispatch.md's reuse-advance handle is `SendMessage(to={handle}, message=output.prompt)` with the verbatim template demoted to break-glass; ensign-shared-core.md carries the advance-pointer bootstrap clause. Test: static contract checks (grep-anchored test beside the existing contract tests) for the three anchors, enforced alongside the live proof in AC-4.
- **AC-4 (live behavior): a live reused worker advanced with the pointer message performs the next stage** — reads the advance file, runs the fetch commands, does the next-stage work, and emits the next-stage `Done:` signal. Test: the mechanism spike above seeds this; implementation verifies whether the runtime-live-e2e claude scenario exercises a reuse advance and, if it does not, extends the live scenario (or records a manual live drive in this entity) so the advanced worker's next-stage report is observed live. The merge gate below enforces the lanes.

## Test plan

- **Go unit + golden fixtures (internal/dispatch, low cost):** advance-file body and envelope goldens (plain / split-root / feedback-reflow / codex host), flag validation (`--advance --bare-mode` exit 2; shared required-input errors), filename `-advance` suffix and collision-keying, envelope field absence (no spawn fields), AC-1 measurement test.
- **Static contract checks (low cost):** AC-3 anchors.
- **Live lanes (merge gate):** this diff touches `skills/first-officer/references/fo-dispatch-core.md` and `skills/ensign/references/ensign-shared-core.md` — host-neutral contract files under `skills/**/references/**` — so per docs/dev/README.md ("Required CI lanes are a function of the diff"), EVERY host live lane is REQUIRED green before merge: claude-live (both matrix variants: sonnet/CI-E2E and claude-opus-4-8/CI-E2E-OPUS), codex-live, and pi-live. claude-fo-dispatch.md additionally pins claude-live; codex-first-officer-runtime.md pins codex-live. A red or flaky lane re-runs to green — never skipped or waved off.
- **No spike needed beyond the recorded one:** the remaining mechanisms (`dispatch build` assembly, `show-stage-def` fetch, split-root guidance emission, SendMessage-to-live-worker delivery) are shipped and live-proven.

## Documentation

No docs-site diff: no page under docs/ enumerates `dispatch build` flags (docs/*.md are releasing / runtime-live-ci / runtime-support methodology docs; runtime-support.md references `dispatch build` only as a probe method). The FO-facing command surface is the skill prose enumerated in the touch points above; implementation applies those before/afters.

## Out of scope

- Pi reuse-advance (deferred by the Pi adapter; fresh redispatch remains its default).
- Emitting the target worker handle from the helper (cycle suffixes make derived names unreliable; the FO roster stays authoritative).
- Any change to reuse conditions 0–4, supersede-shutdown, or the frontmatter advance commit.

## Stage Report: ideation

- DONE: The advancement-message content contract is pinned before the helper is designed: enumerate exactly what the current reuse-advance template carries (stage definition verbatim, checklist, entity path, commit-before-signal, feedback context when re-entering) and design the file/pointer shape to carry all of it — the riskiest mechanism (a reused ensign acting correctly on a pointer instead of an inline body) is spiked or evidenced from the proven initial-dispatch pattern.
  `## Content contract` table maps every current-template element (claude-fo-dispatch.md:40) to its carrier; `## Spike record` — a live completed worker acted on a production-shaped mid-session pointer, on-disk proof observed.
- DONE: ACs include a VALUE measure of the recurring saving: per-advance FO-side tokens O(pointer) vs the current O(stage-section) baseline, with how it is measured.
  AC-1: prompt ≤ 300 bytes invariant to section size vs measured baselines (ideation 4,866 B / validation 3,139 B / implementation 1,259 B vs 239 B pointer), with the in-test old-template comparison spelled out.
- DONE: The contract-prose touch points (fo-contract-core.md reuse section, claude-fo-dispatch.md reuse-advance template, any codex/pi adapter parallel) are enumerated with concrete before/after, and the claude-live merge-gate consequence of touching skills/** is recorded in the test plan.
  `## Contract-prose touch points` lists 7 items with before → after (codex half-line; pi explicitly no-change); test plan records the all-lane merge gate (host-neutral core ⇒ claude-live both variants + codex-live + pi-live).

### Summary

Fleshed out the task: `--advance` mode on `dispatch build` (shared machinery, 7 mode deltas including the `-advance` filename suffix and a spawn-field-free envelope), content contract pinned in a parity table, riskiest mechanism spiked live (completed worker acted on a pointer-shaped advance; on-disk proof), measured AC-1 baselines from the dev workflow's own stage sections, 7 contract-prose before/afters enumerated, and the all-host live-lane merge gate recorded. Target handle stays FO-side; Pi reuse-advance stays out of scope.
