---
title: FO contract token-cleanup proposal
date: "2026-06-15"
status: proposal
---

# FO contract token-cleanup proposal (2026-06-15)

This is an evidence-based, adversarially-verified token-cleanup pass over the first-officer contract files. Each candidate below was proposed by an analyst and then independently verified against the live files and cross-referenced load points; only the candidates that survived verification as `safe-cut` or `cut-with-care` are listed here. The verdicts that came back `keep` are recorded at the bottom with the concrete first-officer misbehavior each removal would cause, so they are on the record and not re-proposed. Cuts are weighted by load frequency: a cut in `first-officer-shared-core.md` or `claude-first-officer-runtime.md` is boot-resident and pays back every FO session; `using-claude-team/SKILL.md` loads at first team-mode dispatch (most working sessions); `present-gate/SKILL.md` loads only per gate. Total confirmed savings across the kept candidates: **638 tokens** (282 boot-resident in shared-core, 165 boot-resident in the runtime adapter, 98 at first dispatch, 93 per gate). These are proposals for the captain, not applied edits.

This pass does not overlap the prior `fo-contract-prose-audit` (PR #367), which cut only four dead `S7b`/`(B6)` step-numbering pointers plus a comm-officer polish. None of those four locations recur here. That audit also FLAGGED-not-cut the merge module's two adjacent mod-block sections as a judgment-call restructure out of scope; this pass likewise leaves them untouched.

A note on the two cross-file dedup pairs. The shared-core `## Dispatch (deferred module)` and `## Merge and Cleanup (deferred module)` sections (SC-5, SC-7) restate the dispatch/merge machinery that the boot-resident runtime adapter also carries. The redundant copy is the shared-core one, so the cut is applied THERE and the runtime adapter is kept intact. The mirror-image runtime candidates (RT-1, RT-2) came back `keep` precisely because the runtime adapter is the SOLE place naming the concrete reference filenames (`references/claude-fo-dispatch.md`, `references/claude-fo-merge.md`); cutting the runtime copy would strand the filename. So each concept is cut exactly once, at the copy that does not hold the load-bearing filename.

## first-officer-shared-core.md (boot-resident, every session)

Confirmed savings here: 282 tokens. Highest-value file in the pass.

### SC-5 — `## Dispatch (deferred module)`, cut-with-care, ~85 tokens
Replace the full section ("The dispatch machinery ... A greet-and-stop boot never reads it.") with a pointer to the runtime adapter's dispatch section. The deferral invariant is independently stated at shared-core line 3 ("neither is read at boot") and Startup step 9 line 30; the reuse-conditions deferral is restated at line 118; the concrete load point lives in the boot-resident runtime adapter. Verifier note: no constructible FO misbehavior; the loss is descriptive-completeness only (this is the sole runtime-agnostic single-list enumeration of the deferred dispatch inventory). A pointer that explicitly restates "deferred, loaded at first dispatch, skipped on greet-and-stop" is tighter than a bare "see the runtime adapter."

### SC-7 — `## Merge and Cleanup (deferred module)`, safe-cut, ~60 tokens
Replace the full section ("The terminal merge-and-cleanup ceremony ... reaches `present-gate`") with a pointer to the runtime adapter's merge section. Verifier ran six break attempts: the load point survives at runtime lines 11-13 (boot-resident), the mod-block-enforcement location fact survives at shared-core line 190, the lazy-precedent analogy and the boot-non-read guard survive verbatim in the runtime copy and `claude-fo-merge.md`. The replacement preserves the platform-agnostic indirection (shared core defers to the runtime adapter, not directly to the Claude file), so Codex/Pi runtimes are unaffected. None of the six attempts produced misbehavior.

### SC-2 — Startup step 4, the README-body defer rationale, cut-with-care, ~36 tokens
Replace the defer sentence so it drops the "the boot JSON does not carry the stage taxonomy, so the frontmatter read stays before-greet" rationale while keeping the operative rule. The operative rule (read frontmatter incl. stage taxonomy at boot, before greet) survives via step 4's imperative and its position before step 9's greet. Verifier note: an FO trying to fold the step-4 frontmatter read into the step-5 boot JSON call would find the step-5 JSON section list carries no taxonomy and cannot substitute, so the sequencing holds structurally. Residual is a mild defensive-rationale loss against a future contract editor, not an FO action.

### SC-8 — `## State Management`, worktree-ownership defer sentence, cut-with-care, ~36 tokens
The analyst's proposed one-line replacement over-cuts: it additionally drops the timing invariant ("they matter only once a worktree stage dispatches") and the named split-root deliverable-isolation contract, which is the sole boot-resident signpost to a known footgun (writing the entity body to a worktree copy instead of the state checkout). Recommend the narrower cut: keep the timing clause and the deliverable-isolation contract name, drop only the "which active state lives in the worktree copy vs. `main`" filler. The full worktree-ownership rules genuinely live in `claude-fo-dispatch.md` (lines 68, 74), loaded at first dispatch when needed.

### SC-3 — `## Status Viewer`, "Distinct from event-loop `status` calls" parenthetical, safe-cut, ~28 tokens
Delete the parenthetical distinguishing captain-facing status calls from FO-internal event-loop reads. Both behavioral branches are fully and independently specified elsewhere: captain-facing rendering at line 69 (verbatim fenced block, unconditional), event-loop rendering at `claude-fo-dispatch.md` line 184 (parse as JSON, loaded at first dispatch). The event loop is a dispatch-time concept that does not exist at boot, so no boot-time conflation is possible. Verifier built four break attempts, none succeeded.

### SC-12 — `## Working Principles`, section lead-in, safe-cut, ~24 tokens
Delete the introductory sentence ("These habits implement the operating posture ... how the FO frames work and adjudicates gates in practice."). The bullets below are self-contained; the entry-point linkage they reference is preserved by the explicit "(entry-point principle 1)" / "(entry-point principle 3)" parentheticals inside the bullets, and gate-adjudication behavior lives in `## Completion and Gates`, not in the lead-in. No rule is encoded in the preamble.

### SC-4 — `## Status Viewer`, single-entity-lookup bullet, safe-cut, ~8 tokens
Tighten "Single-entity lookup: `--resolve {ref}`, then a follow-up `--where slug={resolved-slug}` for a fuller view." to "Single-entity: `--resolve {ref}` then `--where slug={resolved-slug}`." The two-call sequence and the data dependency survive via the retained "then" and the `{resolved-slug}` placeholder. This is a display-helper list, not a mutation guard; the load-bearing `--resolve` canonicalization rule at line 79 is untouched.

### SC-6 — `## Completion and Gates`, count-summary sentence, cut-with-care, ~5 tokens
Integrate the explicit count ("{N} done, {N} skipped, {N} failed") into the preceding checklist bullet rather than stating it as a separate sentence. The count MUST stay inline (not be deleted) — for a non-gated completion this is the sole place requiring an explicit tally; the present-gate `Assessment:` line only fires for gated stages. Net saving is modest (~5 tokens) and the concept is not deferrable to a later-loaded file, hence cut-with-care.

## claude-first-officer-runtime.md (boot-resident, every session)

Confirmed savings here: 165 tokens.

### RT-4 — `## Team-mode ensign-chat hint`, lines 21-28, cut-with-care, ~120 tokens
Compress the hint section's tracking and exception prose. The reword MUST retain the placement guard "Append it to the dispatch announcement; do not emit it as a standalone message." — that phrase is the sole statement of where in the turn the hint goes, and it interacts with `using-claude-team` `## Awaiting Completion`, which bans any post-`Agent()`-dispatch turn that emits standalone text (the hallucination-drift footgun). Drop the explanatory parentheticals and per-mode reasons; keep once-per-session, first-dispatch, the gate/terminal/bare/degraded/single-entity skips, and the session-memory tracking. The single-entity gate auto-resolve (lines 29-30) stays untouched. Net saving is lower than a naive count because the placement guard and the exact hint text must survive.

### RT-3 — `## Entity-Body Inspection`, lines 37-39, safe-cut, ~45 tokens
Collapse the body sentence to a thin pointer at the shared core's `## Probe and Ideation Discipline`. The shared-core copy (lines 213-219) is strictly more complete: it carries the Grep-over-Read rule, the Read-then-Bash staleness echo already platform-qualified ("On Claude Code"), and the full `status --set` narration form. Both files are boot-resident and co-loaded, so the defer target is resident before entity-body inspection is needed. Keep the `## Entity-Body Inspection` heading as a thin pointer (the file's line-3 index enumerates it) to avoid a dangling reference; collapse only the body sentence.

## using-claude-team/SKILL.md (loaded at first dispatch)

Confirmed savings here: 98 tokens.

### UCT-4 — `## Terminal Team Teardown`, lines 99-101, cut-with-care, ~62 tokens
Compress the section-head meta-commentary that distinguishes terminal teardown from `## Awaiting Completion`, but retain an explicit conflation guard at the section head. The phase reconciliation survives in the section body (line 109 "this is the one place the FO actively attempts `TeamDelete`") and at line 81's forward cross-ref, but only at the TAIL — the cut paragraph is the only framing BEFORE the line-105 `TeamDelete` instruction. Removing it risks a top-down reader treating the required terminal `TeamDelete` as the banned premature-teardown and stranding the `claude -p` subprocess (#275/#282). Keep "This governs the TERMINAL phase only" plus a one-clause "do not conflate with the pre-completion ban above."

### UCT-1 — `## Deferred Team Tools`, line 13, safe-cut, ~22 tokens
Delete "Once a tool's schema appears in the ToolSearch result, it is callable exactly like a normal tool." The fetch-then-callable mechanism survives in the two preceding sentences ("calling one directly fails until its schema is fetched" + "Before the first call to any team tool, run ToolSearch ..."), which also carry the first-call-only scoping so no re-fetch-every-call regression appears. The harness also delivers this exact fact verbatim to every FO at boot via the deferred-tool system reminder. Verifier ran three break attempts, all blocked by the surviving sentences.

### UCT-5 — `## Degraded Mode`, lines 36-38, safe-cut, ~14 tokens
Drop the overlapping adjectives ("explicit, session-wide mid-session transition") and keep the irreversibility invariant verbatim. Scope is restated at line 50 ("for the remainder of the session") and line 60; "explicit" is in mild tension with the auto-triggers (first/second dispatch failure trip Degraded Mode without an explicit command); "mid-session" is carried by the Triggers subsection. Verifier built three misbehavior scenarios, each blocked by a labeled subsection (Captain Report Template, Effects, Triggers). The lost "mid-session vs startup bare mode" signpost is recoverable from Triggers — small enough to remain safe-cut.

## present-gate/SKILL.md (loaded per gate)

Confirmed savings here: 93 tokens.

### SC... present-gate cuts

### PG-6 — Rule 7, lines 42-43, safe-cut, ~30 tokens
Compress to "No format-pedantry asides. Surface format drift only if it blocks the gate, as a Material finding." The reword preserves the directive heading, the suppression gate ("only if it blocks the gate"), the Material-tier routing ("as a Material finding"), and the not-a-separate-paragraph guard (the word "aside" in the heading is a separate paragraph). Dropped content is the illustrative `1./2./3./4.` vs `**AC-N**` examples plus a rationale phrase whose operative meaning is carried by the retained conditional. The term "format drift" is already used un-illustrated in Rule 4 and the template.

### PG-7 — Rule 8, lines 43-44, cut-with-care, ~20 tokens
Drop the repetitive "One sentence, not a section" tail, but keep Rule 8's bold lead-in header (all ten rules share the `**header.** body` form) and the explicit "opens or closes a worktree" dual-direction cue. The merge module clears `worktree=` at terminalization (a close case), and terminal/merge stages still pass through a gate, so the close direction is real; "changes worktree state" alone risks an FO at a terminal merge gate omitting the worktree-removal heads-up because it pattern-matched the open case. The "not a section" placement guard is subsumed by "the Decision line names it in one sentence."

### PG-5 — Rule 6, lines 41-42, cut-with-care, ~15 tokens
Compress the Good/Bad example pair into one inline `e.g. ... not "address the five notes."` clause. The core guard ("name the specific concerns by content, not by reference") and the count-by-reference footgun both survive. Keep Rule 6's bold headline ("Bounce-back recommendations name the concrete asks.") so it does not become the only rule of ten without a scannable lead. Real saving is ~15 tokens, below a naive estimate, because the compressed examples are retained.

### PG-2 — assembly-rules preamble, line 34, cut-with-care, ~12 tokens
Delete "The template is the floor, not the ceiling." and combine into "The FO MUST hold to this discipline when filling the template:". Each of the ten numbered rules is a self-standing imperative; rule 3 (inline evidence) and rule 8 (worktree sentence) independently mandate additions beyond the bare template, and rule 9 independently caps length. Residual is a stylistic-permission loss (the explicit latitude to exceed the skeleton becomes implicit), not a guard/sequencing/invariant loss — hence cut-with-care.

### PG-8 — Rule 9, line 44, safe-cut, ~11 tokens
Compress to "Target length: 15-25 lines of FO-authored prose. If exceeding 25, trim." The budget (15-25) and the ceiling enforcement (trim if >25) survive; the dropped middle sentence is a belt-and-suspenders restatement that also introduced a scope ambiguity ("full gate message" vs the header's "FO-authored prose"). The directional "cut narration not content" signal is carried in-file by Rules 1, 3, and 7, all loaded at the same gate-time.

### PG-1 — lines 8-9, safe-cut, ~5 tokens
Tighten "to render the decision the FO has already made" to "to render the FO's decision." The render-not-decide guard (the verb "render") is preserved; the always-on decide-to-gate / AC-cross-check clause is untouched; the temporal ordering ("after the FO has decided") is independently stated in the frontmatter description (line 3) and in boot-resident shared-core lines 120-126.

## Top recommendations (highest value first)

Boot-resident cuts rank above equal-size lazy-loaded cuts. Within a file, larger and higher-confidence first.

1. **RT-4** — runtime, team-mode ensign-chat hint compression, ~120 tokens, boot-resident. Largest single saving in the pass; cut-with-care, must retain the placement guard and exact hint text.
2. **SC-5** — shared-core, dispatch deferred-module section to a pointer, ~85 tokens, boot-resident. Cut-with-care; loss is descriptive only, the runtime adapter holds the load point.
3. **SC-7** — shared-core, merge-and-cleanup deferred-module section to a pointer, ~60 tokens, boot-resident, safe-cut. Six break attempts failed.
4. **RT-3** — runtime, entity-body inspection to a pointer, ~45 tokens, boot-resident, safe-cut. Shared-core copy is strictly more complete.
5. **SC-2** — shared-core, README-body defer rationale, ~36 tokens, boot-resident, cut-with-care.
6. **SC-8** — shared-core, worktree-ownership defer sentence (narrower than proposed), ~36 tokens, boot-resident, cut-with-care.
7. **SC-3** — shared-core, event-loop status parenthetical, ~28 tokens, boot-resident, safe-cut.
8. **SC-12** — shared-core, Working Principles lead-in, ~24 tokens, boot-resident, safe-cut.
9. **SC-4** — shared-core, single-entity-lookup bullet, ~8 tokens, boot-resident, safe-cut.
10. **SC-6** — shared-core, count-summary inline-fold, ~5 tokens, boot-resident, cut-with-care.
11. **UCT-4** — team skill, terminal-teardown meta-commentary, ~62 tokens, first-dispatch, cut-with-care. Keep a section-head conflation guard.
12. **UCT-1** — team skill, deferred-tool callability sentence, ~22 tokens, first-dispatch, safe-cut.
13. **UCT-5** — team skill, Degraded Mode adjectives, ~14 tokens, first-dispatch, safe-cut.
14. **PG-6** — present-gate Rule 7, ~30 tokens, per-gate, safe-cut.
15. **PG-7** — present-gate Rule 8, ~20 tokens, per-gate, cut-with-care. Keep the bold header and dual-direction cue.
16. **PG-5** — present-gate Rule 6, ~15 tokens, per-gate, cut-with-care.
17. **PG-2** — present-gate assembly preamble, ~12 tokens, per-gate, cut-with-care.
18. **PG-8** — present-gate Rule 9, ~11 tokens, per-gate, safe-cut.
19. **PG-1** — present-gate lines 8-9, ~5 tokens, per-gate, safe-cut.

## Kept as load-bearing (do not re-propose)

These came back `keep` under adversarial verification. Each line names the FO misbehavior the cut would cause.

- **SC-1** (shared-core, contract version gate consolidation): the two abort classes are NOT symmetric. Binary-absent carries an explicit "Do NOT run `spacedock doctor`" guard (doctor is a subcommand of the missing binary); binary-out-of-range REQUIRES doctor. Consolidating drops the negative guard, so an FO with a stale `SPACEDOCK_BIN` runs `spacedock doctor` on an absent binary, gets command-not-found, and surfaces a misleading diagnostic instead of the install hint. This is a Startup step-1 rule that runs before discovery, so deferral is structurally impossible; it is the sole statement.
- **SC-9** (shared-core, FO Write Scope lead-in "on main — nothing else"): "on main" is the scoping invariant separating direct-on-main writes from worktree-worker writes, and "nothing else" is the closure guard. Dropping them lets the allow-list over-claim as "the only things the FO may ever write anywhere," and an FO writing `### Feedback Cycles` to a worktree copy can no longer tell whether that permitted write is in scope.
- **SC-10** (shared-core, "via `spacedock status --set` for all field updates"): sole statement that `--set` is the EXCLUSIVE path for frontmatter field updates. Without it an FO doing a bulk/unusual edit hand-edits the `.md` frontmatter on main, bypassing the launcher mutation guards, the staleness-safe `--set` narration, and (on split-root) the tool-managed atomic-commit concurrency protection. The proposed "see New entity files" pointer misdirects to `spacedock new` (creation), not `--set` (update).
- **SC-11** (shared-core, "the FO owns the process it operates ... distinct from the product the workflow builds"): the dynamic process/product discriminator. In a meta-workflow that PRODUCES a README, the static reword ("not product scaffolding") lets an FO conclude it owns a produced `output/.../README.md` and edit it directly on main, bypassing the dispatched-worker-in-worktree path. Sole statement.
- **SC-13** (shared-core, "keep dispatching other ready entities when one blocks"): sole statement that a single-entity block must not stall the whole event loop. The per-entity halt language at lines 116/196 says a blocker "legitimately halts the turn"; without this clause an FO reads that as license to stop the whole session, leaving independently-ready entities idle. The event loop lives in a lazily-loaded file that does not restate this posture.
- **RT-1** (runtime, dispatch-reference machinery list + filename): sole statement of the concrete path `references/claude-fo-dispatch.md`; the shared-core deliberately defers the filename here. Cutting it produces a circular pointer with the filename nowhere, so at first dispatch the FO cannot resolve which file to read and dispatches by improvisation.
- **RT-2** (runtime, merge-reference machinery list + filename): the only place in the skill stating `references/claude-fo-merge.md` (repo-wide grep: one hit). Cutting it strands the merge ceremony at terminalization — the FO guesses or skips the set/invoke/clear mod-block sequence and the bounded `TERMINAL_TEARDOWN_BOUNDED` teardown.
- **RT-5** (runtime, Agent Back-off): the reword drops clause 2 ("if you notice the captain messaging an agent without telling you, ask whether to back off"), the sole proactive side-channel-detection guard that makes the Shift+Down ensign-chat feature safe. Without it the FO races the captain's mid-stage steering. It also inverts "until told to resume" into "ask before resuming."
- **UCT-2** (team skill, line 34 "Block all Agent dispatch until team setup resolves"): sole standing dispatch gate for the mid-session recovery window. Without it an FO that hits "Team does not exist" can emit the fresh-suffixed `TeamCreate` and the re-dispatched `Agent()` in the same parallel message, racing the TeamCreate and re-contaminating the desynced slot (the #36806 footgun). The recovery ladder describes steps, not the standing gate.
- **UCT-3** (team skill, lines 88-96 Awaiting Completion guardrails): three sole-statement losses — the `system init`-is-not-a-completion-signal classification, the rationalization anti-patterns ("session ending"/"enough time has passed" + "you cannot measure time from inside a turn"), and the DISPATCH IDLE GUARDRAIL that the boot-resident runtime adapter (`Agent Back-off`) explicitly points into. Cutting them lets an FO act on an idle wake, rationalize a premature `shutdown_request`+`TeamDelete`, or read a between-turns idle as "unresponsive" and tear down the team.
- **PG-3** (present-gate Rule 1 second sentence): the votability test ("if the captain stops reading after line three, they can still vote") is the sole operational acceptance guard for spine self-sufficiency. The first sentence only names which lines are the spine; without the test an FO buries the decision-critical substance below the fold and still satisfies "these four lines are the spine."
- **PG-4** (present-gate Rule 2): the reword drops the behavioral guard ("don't make the captain infer from the Checklist gist or open the entity file") and the `validation picks PASS/REJECTED` disambiguator. Without them an FO renders `Chosen direction: see Checklist` (breaking the line-3 spine) or classifies a REJECTED validation gate as `n/a` (burying the verdict).
- **PG-9** (present-gate Rule 10): sole carve-out that the structural headings (`Gate review:`, `Checklist:`, `Decision:`) and the `{entity title}` placeholder stay generic while only the FO-authored noun localizes. Shared-core line 211 pushes the FO to localize "gate presentations," and `{entity title}` literally contains "entity"; without Rule 10 an FO on a `ticket` workflow rewrites the canonical headings (`Gate review:` to `Ticket review:`) and drifts the gate output off its stable format.

## Cut criteria — empirical verification with a no-guidance control

Every verdict above was reached by *adversarial reasoning*: an analyst constructed a break attempt and judged whether the FO would misbehave without the clause. That is a strong prior, but it is docs-confidence, not measurement — nobody ran an FO with and without the clause and observed the difference. Before applying a cut (and to revisit any `keep`), confirm it with a **no-guidance-control micro-test**, the cheap pre-check that sits *upstream* of the four live-scenario behavior-preservation oracle named in the Closing note.

**Method** (from superpowers v6 `skills/writing-skills/SKILL.md` "Micro-Test Wording Before Full Scenarios"): sample the FO's behavior on the smallest realistic exercise that would expose the clause's job, N≥5 times, in two arms — *with* the clause and *without it* (the control) — in the clause's real surrounding contract, not in isolation. Read every relevant behavior by hand (a string match is not a behavior). Treat run-to-run variance as a binding-failure signal, not noise. Cost is ~$0.15–0.30/sample vs. ~$12 for a full scenario run, so the control is affordable per candidate.

**Decision rule (applies in both directions):**
- *Control already satisfies the behavior* → the clause is **confirmed dead weight**: delete it outright, do not merely compress. Highest-confidence recovery, and it improves the contract. This promotes a `safe-cut` to applied, and it can **overturn a `keep`** — if the FO holds the behavior without the clause, the kept-as-load-bearing verdict was a reasoning error.
- *Removing the clause changes behavior across the N samples* → **confirmed load-bearing**: keep it (or, for a `cut-with-care`, keep the load-bearing fragment and cut only the inert remainder).

**Rewording a `cut-with-care` item** (the prose-efficacy half, same source): prefer a *positive recipe* (what the FO should DO) over a prohibition — superpowers measured a prohibition scoring *worse* than no guidance on an output-shape failure, while the positive recipe was best at zero variance. Drop a prohibition whose positive twin already carries the load. Ties go to the shorter phrasing. Do **not** add nuance clauses ("unless it matters") or exemption clauses ("does not apply to X") — both measurably reopen negotiation and degrade a working rule.

**Why the residency weighting is right.** The load-frequency weighting in this proposal is empirically backed: superpowers measured a long session re-reading a resident skill ~500× (`positive-instruction-redesign-design`, 2026-06-10), so a boot-resident cut compounds across the whole session while a lazy-loaded cut is paid once. Honest ceiling: ~46% of a run is prompt-immune thinking/narration (`strict-cost-sdd-design`) — prose changes govern only the instruction-resident remainder, so expect "delete dead weight + collapse prohibitions," not a flat percentage.

This proposal's cuts are applied by the `fo-contract-token-cut` task (`y2r7ew51xqs6q3avsb6mcaka`): it first validates the no-guidance control on three candidates spanning the verdict space — SC-13 (a `keep`), SC-5 (`cut-with-care`), SC-3 (`safe-cut`), reusing the bare-`claude` launch + durable-state grade the `haiku-loop-spike` proved — then applies the full list (re-testing the 13 `keep`s, since the control can overturn one) and ships the trimmed files through the dispatched-worker path. The deliverable is the trimmed contract files; the control is the method.

## Closing note

These are PROPOSALS for the captain. The FO contract files are shipped scaffolding (`skills/` is product, not FO-direct-editable per the scaffolding guardrail in `## FO Write Scope`), so the actual edits ship through a dispatched worker in a worktree under test, exactly the path the prior `fo-contract-prose-audit` used: the four live shared scenarios (`gate-guardrail`, `rejection-flow`, `feedback-3-cycle-escalation`, `merge-hook-guardrail`) stay green on Claude and Codex as the behavior-preservation oracle, plus a `wc` size-delta and a detached adversarial audit of the diff. The boot-resident cuts (shared-core and the runtime adapter, items 1-10 in Top recommendations) are the highest value: they pay back on every FO session, where a per-gate cut pays back only when a gate renders.
