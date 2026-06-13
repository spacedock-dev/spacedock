# j9 Phases 2-3 — Lazy-TeamCreate + shallow-boot-then-greet

**Milestone:** 0.20.3 (0203 FO efficiency)
**Stage:** ideation
**Backbone task:** j9 (one task, three phases). Phase 1 (the contract structural split) is shaped in `T1-ideation.md` (labeled T1 there). This document shapes **Phase 2 (lazy-TeamCreate)** and **Phase 3 (shallow-boot-then-greet)**.
**Depends on:** Phase 1 makes the dispatch/merge content lazily loadable; Phase 3's greet-first sequencing rests on that split. Phase 2 (lazy-TeamCreate) needs no split.

---

## Problem

Boot forensics on a live FO session (`/tmp/boot-analysis-spacedock-v1.md`) measured **~160k peak context and ~13.6 min wall-clock to reach an interactive greet — with no team created and no worker dispatched.** A 100% pre-dispatch session paid full deep-boot cost. Two structural levers dominate that picture and are this document's targets:

1. **The 89k team-mode prefix re-cache.** The milestone README ranks lazy-TeamCreate as the single biggest lever at **~89k cache-creation removed** — far larger than the contract reads (~16k), the status-table render (~8.7k), or the mod reads (~6.5k). Creating a Claude team re-caches the whole conversation prefix under the new team context. The measured session never created a team (it never dispatched), so on a no-dispatch boot the entire 89k is avoidable — *if* the contract stops telling the FO to create a team at startup. Today the runtime adapter's Team Creation section reads "At startup (after reading the README, before dispatch)," which couples the team into the boot path.

2. **Everything read before the greet that the greet does not need.** The forensics timeline shows the FO, before answering the captain, reading both FO references (~16.2k), the full `docs/dev/README.md` (~7.9k), both mod files (~6.5k), and rendering the human status table (~4.3k) — then thinking at 100k+ context, where the two slowest turns (128.6s, 100.1s) fired. The greet itself runs off `status --boot --json`, already executed at t=27s (forensics row 8), long before any of those heavy reads. The cost is structural: the loader and Startup procedure front-load all of it, and generation latency scales with loaded context, so the wall-clock is dominated by thinking at a context the greet never required.

The goal (milestone success criterion): an FO reaches interactive readiness — greet + state summary + *able to present a gate* — in seconds at **<~60k** context, deferring everything not needed for an accurate greet to the moment it is needed.

---

## Spike (riskiest-first) — can TeamCreate and the team-dependent startup steps all defer past the greet without breaking a correctness guarantee?

**The riskiest unknown:** the shallow-boot bet is that the FO can greet *before* creating a team, *before* reading the mods, and *before* loading the dispatch/merge modules — yet still honor every correctness guarantee those steps carry (a merged PR advances, an orphan surfaces, a mod-block resumes, a superseded agent is cleaned, a split-root state tree is fresh). If any of those guarantees genuinely must fire before the greet, the shallow-boot sequence is wrong. This spike enumerates every boot step in execution order, classifies each as must-run-before-greet vs defer-to-first-action, and hunts for a correctness break.

**Honesty on spike depth:** this is STATIC analysis — step enumeration plus correctness tracing over the SKILL loader, the shared-core Startup procedure, the runtime Team Creation / standing-teammate / Event-Loop sections, and the two mod files. It establishes that the deferral *can* be clean. It does NOT prove the FO actually greets correctly with the team and modules unloaded — that is the live `shallow-boot` drive in `internal/ensigncycle` (AC-1), at implementation/validation. The spike says the sequence is sound on paper; the live scenario proves it in flight.

### Spike step 1 — every boot/startup step in execution order

Traced from the SKILL loader through the first event-loop pass. Sources: `SKILL.md` (lines 7-27), `first-officer-shared-core.md` (`## Startup` steps 1-7), `claude-first-officer-runtime.md` (`## Team Creation`, `### Standing teammate discovery pass`, `### Standing teammate lazy-spawn`, `## Event Loop` step 0), the startup mod hooks in `docs/dev/_mods/comm-officer.md` and `docs/dev/_mods/pr-merge.md`.

| # | Step | Source | What it does at boot |
|---|------|--------|----------------------|
| L1 | Skill load: inline shared core (`@references/first-officer-shared-core.md`), read runtime adapter | `SKILL.md` 18, 23-25 | Loads the contract. Phase-1 split slims this to boot-resident core only. |
| S0 | Single-entity-mode check | `SKILL.md` 7-14 | Non-interactive + named entity → bounded mode (out of scope for shallow-boot, which is the interactive greet path). |
| S1 | Contract-version gate (`spacedock --version`, parse `contract <N>`, range check) | shared-core Startup 1 | Aborts on missing binary / out-of-range contract. |
| S2 | `git rev-parse --show-toplevel` | shared-core Startup 2 | Project root. |
| S3 | Workflow discovery (`status --discover` or explicit path) | shared-core Startup 3 | Resolves `{workflow_dir}`. |
| S4 | Read `{workflow_dir}/README.md` | shared-core Startup 4 | Mission, entity labels, stage ordering/defaults, stage properties (initial/terminal/gate/worktree/feedback-to/agent). |
| S5 | `status --boot` (FO consumes `--boot --json`) | shared-core Startup 5 | MODS, ID_STYLE, NEXT_ID, ORPHANS, PR_STATE, DISPATCHABLE, STATE_BACKEND, team_state, sandbox — one call. |
| S6 | Split-root state halt-gate | shared-core Startup 6 | If `state_backend==split-root && entity_dir_present==false`, HALT — state checkout not initialized (would render EMPTY + VALID, a silent failure). |
| S7 | Split-root pull-on-boot (`git -C <state> pull --rebase`) | shared-core Startup 7 | Integrate peers' state once at boot; rebase-conflict → HALT. |
| T1 | Team Creation (`Skill(using-claude-team)` → TeamCreate) | runtime Team Creation | **The ~89k re-cache.** Currently framed "at startup … before dispatch." |
| T2 | Standing-teammate discovery pass (`dispatch list-standing`) | runtime discovery pass | Records standing-teammate mod paths. "No spawn calls at boot." |
| M1 | Startup mod hooks: pr-merge `## Hook: startup` | `pr-merge.md` 13-26 | Scan entities with non-empty `pr` + non-terminal; `gh pr view`; advance MERGED, report CLOSED, no-op OPEN. |
| M2 | Startup mod hooks: comm-officer `## Hook: startup` | `comm-officer.md` 12-27 | If team lacks `comm-officer`, spawn it (fire-and-forget). Needs a team. |
| E0 | Event Loop step 0: reconcile sweep (`dispatch reconcile --team-name {team_name}`) | runtime Event Loop 0 | "(a) at boot, AFTER the split-root pull, BEFORE the first dispatch." Roster-derived A/B/C drift + git-only D/E. |
| E1+ | Event Loop steps 1-4: PR-pending check, mod-block check, `--next` dispatch, idle | runtime Event Loop 1-4 | The normal dispatch loop. |

**Note on the greet:** there is no explicit "greet" step in the contract today — the FO greets when it has enough state to report. The forensics show `--boot --json` (S5) already executed at t=27s, but the FO did not greet until t=511s because it ran S4 (whole README), M1, M2, and the orphan reconciliation reads first, all at rising context. The shallow-boot change is to **insert the greet immediately after the accuracy-critical boot steps and stop for input**, pushing T1, T2, M1-as-spawn, M2, E0, and the deferred contract modules past it.

### Spike step 2 — classify each step: must-run-before-greet vs defer

The discriminator is: *does the greet's accuracy or a correctness guarantee depend on this step having run before the FO speaks to the captain?* An accurate greet = correct workflow state summary + the ability to present a gate. A guarantee that the deferred timing still catches on first-action is safe to defer.

| Step | Verdict | Reason |
|------|---------|--------|
| L1 (boot-resident core load) | **before-greet** | The greet/gate/status-viewer instructions live here. Phase-1-slimmed. |
| S1 (contract gate) | **before-greet** | A version mismatch must abort before any state read; greeting against a wrong-contract binary is unsafe. |
| S2 (git root) | **before-greet** | Every subsequent path resolves from it. |
| S3 (workflow discovery) | **before-greet** | The greet names the workflow; can't greet without `{workflow_dir}`. |
| S4 (README read) | **before-greet for the greet-relevant fields; the full ~7.9k read DEFERS** | The greet needs entity-label, stage names/ordering, gate flags — a small slice. The full body (proof-policy prose, add-a-scenario procedure, CI setup, PR-body template, task template) is dispatch/validation material the greet never uses. See spike step 5. |
| S5 (`status --boot --json`) | **before-greet** | This IS the greet's data source. ORPHANS, PR_STATE, DISPATCHABLE, STATE_BACKEND, team_state all feed the state summary. |
| S6 (split-root halt-gate) | **before-greet** | A silent-empty state checkout would make the greet REPORT A FALSE "no work" state. The halt is an accuracy guarantee for the greet itself — it must fire before the FO summarizes state. |
| S7 (split-root pull-on-boot) | **before-greet** | State freshness: the greet must report peers' committed state, not a stale local tree. A rebase-conflict halt also pre-empts a dispatch against an unmerged tree. The greet's accuracy depends on this. |
| T1 (TeamCreate) | **defer-to-first-dispatch** | The ~89k re-cache. No team is needed to greet, report state, or present a gate (gates render as captain-facing text, not team messages). Needed only when the FO dispatches a worker. |
| T2 (standing-teammate discovery) | **defer-to-first-dispatch** | Records mod paths for later spawn; "no spawn at boot." Not greet-relevant. Travels with the dispatch module. |
| M1 (pr-merge startup hook) | **before-greet — but as a STATE READ, not via the mod file** | This carries a real guarantee: a PR merged while the FO was away must be advanced, and the greet should report it. BUT `status --boot --json` already carries `pr_state` (PR-pending entities + current merge state). The accuracy obligation is "the greet reports merged PRs," satisfied by reading `pr_state` from the boot JSON — which is before-greet anyway (S5). Reading the **pr-merge mod FILE** (~3.3k) and running its `gh pr view` advancement can ride the first event-loop pass. See the correctness hunt (step 3, break #1). |
| M2 (comm-officer startup hook = spawn) | **defer-to-first-dispatch** | Spawning the prose-polisher needs a team (T1). comm-officer polishes deliberate drafts (PR bodies, gate summaries) — none of which exist at the greet. The greet is a live captain reply, explicitly OUT of comm-officer's scope (`comm-officer.md` 41-42). Reading the mod file and spawning both defer to first dispatch alongside T1. |
| E0 (reconcile sweep) | **defer-to-first-dispatch** | Phase-1 finding C3: reconcile fires "before the first dispatch," NOT before-greet, and needs a `team_name` (A/B/C are roster-derived). It travels with the dispatch module that loads at first dispatch. See the correctness hunt (step 3, breaks #2/#3/#4). |
| E1+ (dispatch loop) | **defer-to-first-action** | The loop runs when there is work to dispatch; the shallow boot greets and stops for input first. |

### Spike step 3 — correctness hunt: does deferring past the greet drop any guarantee?

For each guarantee carried by a deferred step, the question is: does the deferred timing (first-action / first-dispatch) still catch it, or MUST it run pre-greet?

**Break #1 — a merged PR not advanced (pr-merge startup hook M1).** The guarantee: a PR that merged while the FO was offline must advance its entity to terminal and archive it. If the FO greets without running M1, does the merged PR get lost?
- **Caught on first-action, with one accuracy caveat.** The pr-merge hook's advancement logic is duplicated in the Event Loop: step 1 ("Check PR-pending entities … advance merged PRs") runs on the first loop pass, AND pr-merge declares an `## Hook: idle` that re-checks (`pr-merge.md` 28-30, "defense in depth"). So the *advancement* is caught the moment the FO acts. The caveat is the *greet's accuracy*: if the FO greets reporting "PR #347 pending" when it actually merged, the greet is stale. **Resolution:** the greet reads `pr_state` from `status --boot --json` (S5, before-greet), which reports each PR-pending entity's current merge state. The boot probe itself queries merge state, so the greet can say "PR #347 (now MERGED — will advance)" accurately WITHOUT reading the mod file or running the mod's `gh pr view`. The mod-file read and the advancement action defer to first event-loop pass. **No guarantee dropped; the greet stays accurate off the boot JSON.** (Verify during implementation that `--boot --json`'s `pr_state.entries[].state` reflects live merge state, not just the stored `pr:` field — the boot probe parity test `internal/status/boot_probe_parity_test.go` is the place this is pinned.)
- **Genuine constraint:** this only holds if `pr_state` in the boot JSON carries live merge state. If it carries only the stored `pr:` value, the greet cannot report a freshly-merged PR accurately and M1's `gh pr view` would need to run pre-greet. This is the one item to confirm at implementation; the spike flags it as the single accuracy dependency to pin.

**Break #2 — an orphan not surfaced (reconcile sweep E0).** The guarantee: a lingering/superseded agent, an un-advanced PR, a stale branch, or stale local main is detected and acted on.
- **Caught: ORPHANS already surface at the greet via `status --boot --json`.** The boot JSON's `orphans` array (worktree fields cross-referenced against filesystem + git state) is computed by the binary, needs no team, and is before-greet (S5). The greet reports anomalies from it (Startup step 5: "Report anomalies; do not auto-redispatch"). The reconcile sweep's roster-derived classes (A lingering, B superseded, C un-advanced PR) are a SUPERSET that needs a `team_name`, and the contract already says it fires "before the first dispatch," not before-greet. The git-only classes (D stale branch, E stale local main) are session-independent but are pre-dispatch hygiene, not greet-accuracy. **No guarantee dropped: the greet surfaces filesystem/git orphans off the boot JSON; the roster-derived reconcile rides the first-dispatch path where its `team_name` exists.** This is exactly Phase-1's C3 finding ("the reconcile sweep travels with the dispatch module — fires before-first-dispatch, not before-greet").

**Break #3 — a mod-block not resumed.** The guarantee: an entity left `mod-block`-pending across a session boundary resumes its pending merge action.
- **Caught on first-action.** Event Loop step 2 ("Check mod-blocked entities … resume its pending action") runs on the first loop pass; the runtime Mod-Block-at-Terminal section ("On session resume, scan entities with non-empty `mod-block` and resume") is merge-module content. A mod-blocked entity is by definition NOT dispatchable (the loop refuses new work for it), so resumption belongs in the loop, not the greet. The greet reports the pending state from the boot JSON's `mods`/`pr_state` view. **No guarantee dropped: resumption is a first-action obligation, and the merge-module read defers to terminalization where it's needed.**

**Break #4 — a superseded agent not cleaned (supersede-shutdown / terminal teardown).** The guarantee: a stale cohort from a prior dispatch is shut down.
- **Vacuous at a shallow boot.** A freshly-booted FO that has not created a team has no live roster to clean. Supersede-shutdown fires "on fresh dispatch from a -cycleN increment" (a dispatch action) and terminal teardown fires at the terminal boundary (a merge action). Reconcile Class A/B (the resume-time backstop for a missed teardown/supersede) needs a team and rides first-dispatch (break #2). **No guarantee dropped: there is nothing to clean before the first dispatch creates a team.**

**Break #5 — split-root state freshness / silent-empty (S6, S7).** These are the one place the hunt says **MUST run pre-greet** — and they already do (classified before-greet in step 2). S6's halt prevents the greet from reporting a false "no work" state off an uninitialized checkout; S7's pull prevents the greet from reporting a stale tree. Neither needs a team. They stay in the before-greet set. **This is the step that genuinely cannot defer** — but it's a shared-core boot-resident step, not a deferred one, so it is no obstacle to the shallow-boot design; it is part of it.

**Summary of the hunt:** no guarantee carried by a deferred step (TeamCreate, reconcile, standing-teammate spawn, the mod-file reads, the dispatch/merge modules) is dropped by deferring past the greet. Every one is either (a) re-asserted as a state READ in `status --boot --json` so the greet stays accurate (merged-PR state, orphans, mod-block/pending), or (b) an ACTION that the first event-loop pass / first dispatch catches before it matters (PR advancement, reconcile A/B/C, mod-block resume, supersede/teardown). The single accuracy dependency to pin at implementation is that `pr_state` in the boot JSON reflects live merge state (break #1). The only genuinely-cannot-defer steps are the split-root halt-gate and pull-on-boot (S6/S7) — and those are already before-greet shared-core steps.

### Spike step 4 — lazy-TeamCreate mechanics

**Where the contract currently mandates team creation at startup:**

`claude-first-officer-runtime.md` `## Team Creation`, line 7 (verbatim):

> "At startup (after reading the README, before dispatch), invoke the generic Claude-team-harness discipline: `Skill(skill="spacedock:using-claude-team")`"

and line 11 (the truthful trigger already present in the next clause):

> "Invoke it before the first team-mode tool call in the session."

Phase-1's C2 tweak (from `T1-ideation.md`) already identifies this: retire "at startup" in favor of the first-dispatch trigger. The two clauses contradict — "at startup" vs "before the first team-mode tool call" — and the forensics confirm the truthful one (the measured session never created a team because it never dispatched).

**Minimal wording change to make it genuinely first-dispatch-triggered:**

- Replace line 7's "At startup (after reading the README, before dispatch), invoke …" with: **"Before the first team-mode dispatch (the first `Agent()` call that uses a `team_name`), invoke the generic Claude-team-harness discipline:"**. Drop the "at startup" clause entirely; keep line 11's "before the first team-mode tool call" as the now-consistent trigger.
- **Companion timing changes (already aligned by Phase 1, confirmed here):**
  - The **standing-teammate discovery pass** (line 19: "After team creation succeeds … and BEFORE entering the normal dispatch event loop") moves with team creation to first-dispatch. It already says "No spawn calls at boot"; the discovery `list-standing` call (cheap) can either ride first-dispatch with the rest of the dispatch module or stay a one-line boot probe — recommend riding first-dispatch so the boot read drops the whole section.
  - The **reconcile sweep** (Event Loop step 0): already "(a) at boot, AFTER the split-root pull --rebase and BEFORE the first dispatch." Its "at boot" timing is the before-FIRST-DISPATCH boot moment, not before-greet. With shallow-boot the FO greets before any dispatch, so reconcile fires when the first dispatch arrives (or at the first idle/explicit-action pass), where its `team_name` exists. **No wording change needed beyond making clear "boot reconcile" means before-first-dispatch, which the contract already says — but its *placement* in the deferred dispatch module (Phase 1) is what realizes the deferral.**
  - The **comm-officer / standing-teammate spawn** (lazy-spawn, line 28: "Before the first `Agent()` call that uses a `team_name`, spawn all declared standing teammates") is ALREADY first-dispatch-triggered in wording. The only change is that the comm-officer mod FILE read (~3.3k) defers with it instead of being read at boot.

The lazy-TeamCreate change is therefore **one clause** in the runtime adapter (line 7), plus the Phase-1 placement of the team/dispatch sections behind the first-dispatch load point so the FO never reads them — or creates a team — at boot.

### Spike step 5 — README-read question (Startup step 4)

**Is the full ~7.9k README read greet-blocking?** No — only a small slice is.

`docs/dev/README.md` is **31,456 chars (~7.9k tokens)**. Its frontmatter (lines 1-27) carries everything the greet needs: `entity-label`/`entity-label-plural` (the FO speaks the workflow's declared noun, per Working Principles), `id-style`, `state`, and the full `stages` block — stage names, ordering, and per-stage `initial`/`terminal`/`gate`/`worktree`/`fresh`/`feedback-to` flags. **That frontmatter is ~700 chars (~175 tokens) of the 7.9k.** The remaining ~7.7k is prose the greet never touches: the field-reference table, the long per-stage Good/Bad/proof-policy narratives (the ~6k anti-prose-grep and detached-audit prose), Workflow State, Runtime Live CI / shared-scenario add procedure, the PR-body template, the task template, Testing Resources, Commit Discipline. All of that is dispatch-time, validation-time, or merge-time material.

**Does `status --boot --json` already carry the greet-relevant slice?** Partly, and not the parts that matter. Confirmed by reading `internal/status/json_commands.go` (`bootJSON`): the boot envelope carries `command, mods, id_style, next_id, [min_prefix], orphans, pr_state, dispatchable, team_state, state_backend, definition_dir, entity_dir, entity_dir_present, sandbox`. The `dispatchable` array carries each ready entity's `current`/`next` stage NAMES, but the boot JSON does **NOT** carry: the entity-label, the full stage ordering/taxonomy, or the per-stage gate/terminal flags. Those are parsed from the README frontmatter by `internal/status/stages.go` (`mappingValue(doc.Content[0], "stages")`), never emitted by any status `--json` command (confirmed: no `"stages"`/`"labels"`/`"gates"` key in `json_commands.go`). So the greet's need for the entity-label and the gate/stage taxonomy is **not** met by the current boot JSON — it is met only by reading the README frontmatter.

**Recommendation (two viable shapes; recommend the first):**

1. **Read the README frontmatter only at Startup step 4, defer the body.** The FO already parses YAML frontmatter elsewhere; reading just the `---`-delimited head (lines 1-27, ~700 chars) gives the greet the entity-label, stage names/ordering, and gate/terminal flags. The body (the per-stage prose, proof policy, templates, CI docs) defers — it is read when its phase begins: stage Good/Bad and proof policy at dispatch/gate adjudication, the PR-body template at merge, the add-a-scenario procedure at validation. This is a clean, behavior-preserving slim: the same `## stage` subsections are copied verbatim into dispatch messages (Dispatch step 8, "the full stage definition") at first dispatch, so the body is genuinely not needed before then. **~7.7k cut from every boot.** Implementation note: Startup step 4's wording changes from "Read `{workflow_dir}/README.md` for mission, entity labels, stage ordering …" to "Read the README frontmatter for entity labels, stage ordering and per-stage flags; defer the body (per-stage prose, proof policy, templates) to the phase that needs it."

2. **Extend `status --boot --json` to carry the greet slice (heavier, defers to a later task).** Add a `labels` + `stages` projection to the boot JSON so the greet needs zero README read. This removes the README read from boot entirely but is a binary change with its own test surface (boot JSON schema, golden fixtures, the FO's key-order parse) — larger than Phase 3 warrants. **Recommend deferring this to p2/vc or a follow-on; Phase 3 takes shape 1 (frontmatter-only read).**

**Verdict on the README read:** the full ~7.9k read is NOT greet-blocking; only the ~175-token frontmatter is. Phase 3 slims Startup step 4 to a frontmatter read and defers the body. The boot JSON cannot today substitute for the frontmatter (it omits labels + stage taxonomy), so the frontmatter read stays in the before-greet set; only the body defers.

### Spike step 6 — VERDICT and the greet-blocking step set

**Verdict: VIABLE with the adjustments named below.** The shallow-boot + lazy-TeamCreate sequence is structurally sound. No correctness guarantee is dropped by deferring TeamCreate, the reconcile sweep, standing-teammate spawn, the mod-file reads, and the dispatch/merge modules past the greet (spike step 3). Every guarantee is either re-asserted as a state READ in `status --boot --json` (so the greet stays accurate) or caught as an ACTION on the first event-loop pass / first dispatch. The lazy-TeamCreate change is one clause (spike step 4). The README read slims to its frontmatter (spike step 5).

**Adjustments required (none is a redesign):**
- **A1 — lazy-TeamCreate wording:** runtime adapter line 7, "at startup … before dispatch" → "before the first team-mode dispatch." (Phase-1 C2; restated as Phase 2's core change.)
- **A2 — README frontmatter-only read:** Startup step 4 reads the frontmatter, defers the body.
- **A3 — defer the mod-file reads:** pr-merge and comm-officer mod files are read when their hook fires (pr-merge advancement on first event-loop pass; comm-officer spawn at first dispatch), not at boot. The greet reports PR/orphan/mod state off the boot JSON.
- **A4 — insert an explicit greet-and-stop after the before-greet set.** The contract gains a "greet and stop for input" step after S7, ahead of T1/T2/M1-action/M2/E0.
- **A5 — pin the boot-JSON `pr_state` live-merge-state dependency** (spike step 3, break #1) — confirm at implementation that the greet can report a freshly-merged PR off the boot JSON without running the mod.

**The greet-blocking step set (the new shallow-boot sequence) — the deliverable:**

```
contract-gate (S1)
  → git root (S2)
  → workflow discovery (S3)
  → README FRONTMATTER read (S4-slim: entity-label, stage names/ordering, gate/terminal flags)
  → status --boot --json (S5: state summary source — orphans, pr_state w/ live merge state, dispatchable, team_state, state_backend)
  → split-root halt-gate (S6: prevents a false-empty greet)
  → split-root pull-on-boot (S7: state freshness; rebase-conflict HALT)
  → GREET the captain (state summary off the boot JSON; able to present a gate) and STOP for input
```

**Everything deferred past the greet:** TeamCreate (T1, the ~89k re-cache), standing-teammate discovery + spawn (T2 + lazy-spawn), the comm-officer mod read+spawn (M2), the pr-merge mod read + `gh pr view` advancement action (M1-action), the reconcile sweep (E0), the dispatch/merge contract modules (Phase 1), the README body, and the human status-table render (rendered to the captain only on explicit request, once).

**The one step that genuinely cannot defer:** the split-root halt-gate + pull-on-boot (S6/S7) — but it is already a before-greet shared-core step, so it is part of the shallow-boot sequence, not an obstacle to it.

**Honesty on depth:** this VERDICT rests on static step-enumeration and correctness tracing. The live `shallow-boot` drive in `internal/ensigncycle` — observing the FO greet correctly with NO team created and the dispatch/merge modules unloaded — is the implementation/validation proof (AC-1), not this spike's. The spike establishes the sequence *can* be correct; the live scenario establishes that it *is*.

---

## Proposed approach

### The new shallow-boot sequence

Reshape the Startup procedure and the runtime Team Creation trigger so the FO reaches the greet through exactly the before-greet step set (spike step 6), then stops for input. Concretely:

1. **Lazy-TeamCreate (Phase 2).** Change the runtime adapter's Team Creation trigger from "at startup … before dispatch" to "before the first team-mode dispatch" (A1). The FO creates no team at boot; the ~89k re-cache happens only when the first worker is dispatched. The standing-teammate discovery/spawn (already "no spawn at boot" / "before the first `Agent()` with a `team_name`") and the reconcile sweep (already "before the first dispatch") travel with the dispatch module so the FO never reads them — or creates a team — at boot.

2. **Greet off `status --boot --json` (Phase 3).** After S1-S7 (the before-greet set), the FO greets the captain with a state summary built from the boot JSON (orphans, pr_state with live merge state, dispatchable, team_state, state_backend) and the README frontmatter (entity-label, stage taxonomy, gate flags), then stops for input (A4). It can present a gate from this state without a team (gates are captain-facing text).

3. **README frontmatter-only read (Phase 3).** Startup step 4 reads the README frontmatter for the greet-relevant fields and defers the ~7.7k body to the phase that needs it (A2).

4. **Defer mod-file reads (Phase 3).** The comm-officer and pr-merge mod FILES are read when their hooks fire — comm-officer at first dispatch (spawn needs a team), pr-merge advancement on the first event-loop pass (A3). The greet reports PR/mod/orphan state off the boot JSON, never the mod files.

5. **Defer the human status-table render (Phase 3).** The FO never renders the human-formatted table for its own reasoning (it has the boot JSON); it renders it to the captain at most once, on explicit request, per the Status Viewer's existing captain-facing-display rules.

### Greet-blocking classification table

| Step | Classification | Reason |
|------|----------------|--------|
| contract-gate (S1) | **before-greet** | A version mismatch must abort before any state read. |
| git root (S2) | **before-greet** | All paths resolve from it. |
| workflow discovery (S3) | **before-greet** | The greet names the workflow. |
| README **frontmatter** (S4-slim) | **before-greet** | Entity-label, stage names/ordering, gate/terminal flags — the greet's vocabulary and the gate taxonomy. |
| README **body** | defer | Per-stage prose, proof policy, templates, CI docs — dispatch/validation/merge material. |
| `status --boot --json` (S5) | **before-greet** | The state-summary data source (orphans, pr_state, dispatchable, team_state, state_backend). |
| split-root halt-gate (S6) | **before-greet** | Prevents a false-empty greet off an uninitialized checkout. |
| split-root pull-on-boot (S7) | **before-greet** | State freshness; rebase-conflict HALT pre-empts a stale-tree greet/dispatch. |
| **GREET + stop for input** | **the boundary** | State summary + able to present a gate. |
| TeamCreate (T1) | defer (first dispatch) | The ~89k re-cache. No team needed to greet/report/present-a-gate. |
| standing-teammate discovery + spawn (T2) | defer (first dispatch) | Spawn needs a team; not greet-relevant. |
| comm-officer mod read + spawn (M2) | defer (first dispatch) | Polishes deliberate drafts, not the live greet; spawn needs a team. |
| pr-merge mod read + `gh pr view` advance (M1-action) | defer (first event-loop pass) | Advancement caught by Event-Loop step 1 + the pr-merge idle hook; the greet reports merge state off the boot JSON. |
| reconcile sweep (E0) | defer (first dispatch) | Needs a `team_name`; fires before-first-dispatch, not before-greet (Phase-1 C3). |
| dispatch/merge contract modules | defer (Phase 1 load points) | First dispatch / terminalization. |
| human status-table render | defer | Captain-facing, on explicit request, once. |

### Lazy-TeamCreate mechanics (restated for implementation)

Single clause: runtime adapter `## Team Creation` line 7. Current: "At startup (after reading the README, before dispatch), invoke …". New: "Before the first team-mode dispatch (the first `Agent()` call that uses a `team_name`), invoke …". Drop "at startup"; the existing line-11 clause ("before the first team-mode tool call in the session") becomes the consistent trigger. Companion: discovery/spawn and reconcile travel with the dispatch module (Phase 1 placement), so the boot read drops them and the FO never creates a team at boot. No binary change.

### README-read decision

Take spike-step-5 shape 1: Startup step 4 reads the README **frontmatter only** (entity-label, stage taxonomy, gate/terminal flags); defer the body to the phase that consumes it. The boot JSON cannot substitute (it omits labels + stage taxonomy — confirmed in `json_commands.go`), so the frontmatter read stays before-greet; only the body defers. Extending the boot JSON to carry labels+stages (shape 2) is a heavier binary change deferred to a follow-on, not Phase 3.

---

## Out of scope

- **Phase 1 — the contract structural split** (extract dispatch/merge modules behind lazy load points, slim the boot-resident core). Shaped in `T1-ideation.md`. Phase 3's deferred-module unload rests on it, but the split itself is Phase 1's deliverable, not this body's.
- **p2 / vc** — `spacedock pr complete` + `reconcile --act` (the binary-simplification line). Parked to 0.20.4.
- **Extending `status --boot --json` to carry labels + stage taxonomy** (README-read shape 2). A heavier binary change; a follow-on, not Phase 3.
- **T3 — residual-prose audit + comm-officer polish.** Files along post-Phase-1; the cut-list does not exist until the split lands.
- **The Codex and Pi runtime adapters' lazy-TeamCreate.** Phase 2 changes the Claude adapter (the bulk file the forensics measured). The Codex/Pi team-creation timing is a follow-on if their boot cost warrants it.

---

## Acceptance criteria

Each AC names an end-state property of finished Phases 2-3, verified by something outside this task body that can fail. No AC is proven by a string/substring/regex match over an instruction file the model reads — that is banned by this workflow (a passing match only asserts the implementer's own text is present). The behavioral ACs are live drives; the structural AC tests a relationship between independent values.

**AC-1 — A freshly-booted FO greets the captain and reports accurate workflow state with NO team created and the dispatch/merge modules NOT loaded.**
Verified by: a new live shared-runtime scenario `shallow-boot` in `internal/ensigncycle` (added to `sharedRuntimeScenarios()` with Claude + Codex runners per the README's 4-step add-a-scenario procedure). The runner launches the real host front door against a fixture with at least one dispatchable entity, and the host-neutral assertion over `(before, after, observed)` confirms: (a) the FO produced a greet with a state summary in its final message; (b) durable state shows NO team artifact created and NO worker dispatched (no entity advanced past its boot stage, no worktree created); (c) the FO stopped for input rather than auto-dispatching. The team-not-created observation is the lazy-TeamCreate proof; the greet-without-dispatch is the shallow-boot proof. Behavioral and live, not a contract grep. An offline negative in `shared_scenarios_negative_test.go` builds the broken end-state (a team artifact present / an entity dispatched at boot) and proves the assertion goes red.

**AC-2 — The shallow boot greets at materially lower loaded context than the deep boot, with the team-mode re-cache absent.**
Verified by: the `shallow-boot` live scenario's captured host artifacts (stream jsonl / session transcript) show the FO reached its greet without a `TeamCreate` tool call in the pre-greet window, and the pre-greet context never incurred the team-mode prefix re-cache. The check parses the live transcript for the presence/absence and ordering of the `TeamCreate` (or team-tool) call relative to the greet message — a behavioral observation over the real run's tool-call sequence (the team call is absent before the greet), NOT a grep over the contract. The negative control: a deep-boot run (or a fixture forcing eager team creation) shows the `TeamCreate` call before the greet, proving the assertion distinguishes the two.

**AC-3 — The deferred startup steps still run correctly on first action: a merged PR advances and the reconcile sweep fires before the first dispatch.**
Verified by: the existing live scenarios `merge-hook-guardrail` and `rejection-flow`/`gate-guardrail` in `internal/ensigncycle` still pass after Phases 2-3 — they exercise the merge module loading at terminalization (mod-block enforcement, merged-PR advancement) and the dispatch module loading at first dispatch (reconcile sweep, reuse/feedback routing). A green run of all existing live scenarios is the proof that deferring the startup steps past the greet did not drop their first-action obligations. Run via `go test -tags live -run TestLiveClaudeSharedScenarios ./internal/ensigncycle`. (Regression, zero new authoring beyond AC-1's scenario.)

**AC-4 — The greet reads only the README frontmatter, not the full body, for its greet-relevant fields.**
Verified by: a structural check in `internal/contractlint` (the allowed quarantine) confirming the boot-resident Startup step-4 instruction targets the README frontmatter slice (entity-label, stage taxonomy, gate flags) and that the README-body-only anchors (the proof-policy prose, the add-a-scenario procedure, the PR-body template, the task template) are reachable only from the deferred phase modules, not from the boot-resident core. The expected value (which README sections are greet-relevant vs phase-deferred) comes from the deferred-module manifest, an independent source the boot files can diverge from — so the check can fail if a future edit points the boot core at the README body. A control test plants a boot-resident reference into a body-only section and proves the guard goes red. (Structural, not a prose-grep; pairs with AC-1, which is the behavioral proof the greet is accurate off the slim read.)

---

## Test plan

- **AC-1 (`shallow-boot` live scenario):** the costly item. Following the README's 4-step add-a-scenario procedure — host-neutral entry in `sharedRuntimeScenarios()`, fixture + prompt in `shared_fixtures_test.go`, host-neutral assertion (greet present + no-team + no-dispatch + stopped-for-input) and offline negative in `shared_scenarios_negative_test.go`, runner entries in BOTH `claudeScenarioRunners()` and `codexScenarioRunners()`. Live, real host, durable-state + final-message assertion. Cost: one model-spend scenario added to the serial suite (~5-7 min Claude opus). **Spot-check first:** run the parity/definition guards (`TestSharedScenarioRunnerCoverage|TestSharedRuntimeScenarioDefinitions`) at zero spend before paying for the live run.
- **AC-2 (no team-mode re-cache before greet):** rides AC-1's live run — adds a transcript assertion over the tool-call sequence (no `TeamCreate` before the greet message). Cost: zero additional model spend; it reads AC-1's captured artifacts. The negative control (a deep-boot / eager-team fixture showing the call before the greet) proves the assertion distinguishes the two.
- **AC-3 (regression):** zero new authoring — run the existing live scenarios (`gate-guardrail`, `rejection-flow`, `merge-hook-guardrail`, `feedback-3-cycle-escalation`) after Phases 2-3. Cost: the existing serial-suite wall-time, already budgeted in CI. Green proves the deferred startup steps still fire correctly on first action.
- **AC-4 (structural guard):** a Go test in `internal/contractlint`, no model spend, runs in the offline gate job (`go test ./...`). Extends the existing reference-closure pattern. Ships with a control test (planted violation goes red) so the guard is proven able to fail, not vacuous.
- **Implementation-time spike to pin (spike step 3, break #1):** before authoring AC-1, confirm `status --boot --json`'s `pr_state.entries[].state` carries live merge state (so the greet can report a freshly-merged PR without reading the mod). The cheap check is a `internal/status/boot_probe_parity_test.go`-adjacent fixture; pin it before paying for the live scenario. If it does NOT carry live merge state, the pr-merge `gh pr view` would need to stay before-greet for accuracy — flag and re-scope.
- **Fixture vs live:** AC-1, AC-2, AC-3 are live (the runtime integration — does the FO greet correctly with the team and modules deferred — is the claim). AC-4 is a structural fixture over the shipped surface. No AC leans on a prose-grep over the contract.

---

## Spike result

**Verdict: VIABLE with the named adjustments (A1-A5).** The shallow-boot + lazy-TeamCreate sequence is structurally sound. Evidence:

- **The boot steps enumerate cleanly into before-greet vs defer** (spike steps 1-2). The before-greet set is exactly: contract-gate, git root, workflow discovery, README frontmatter, `status --boot --json`, split-root halt-gate, split-root pull-on-boot — then greet and stop. Everything else (TeamCreate, standing-teammate discovery/spawn, the mod reads, the reconcile sweep, the dispatch/merge modules, the README body, the human table) defers.
- **No correctness guarantee is dropped by deferring past the greet** (spike step 3). Merged-PR advancement, orphan surfacing, mod-block resume, and supersede/teardown are each either re-asserted as a state READ in `status --boot --json` (so the greet stays accurate) or caught as an ACTION on the first event-loop pass / first dispatch. The single accuracy dependency to pin is that the boot JSON's `pr_state` reflects live merge state (break #1).
- **The only genuinely-cannot-defer step is the split-root halt-gate + pull-on-boot** (break #5) — and it is already a before-greet shared-core step, so it is part of the shallow-boot sequence, not an obstacle.
- **Lazy-TeamCreate is one clause** (spike step 4): runtime adapter line 7, "at startup … before dispatch" → "before the first team-mode dispatch." The companion timings (discovery/spawn, reconcile) are already first-dispatch / before-first-dispatch in wording; Phase 1's placement realizes the deferral. No binary change.
- **The full README read is not greet-blocking** (spike step 5): only the ~175-token frontmatter is (entity-label + stage taxonomy + gate flags). The boot JSON omits labels and stage taxonomy (confirmed in `internal/status/json_commands.go` — the boot envelope carries no `"stages"`/`"labels"`/`"gates"` key), so the frontmatter read stays before-greet and the ~7.7k body defers.

**Honesty on spike depth:** this is STATIC analysis — step enumeration plus correctness tracing over the SKILL loader, the shared-core Startup procedure, the runtime Team Creation / standing-teammate / Event-Loop sections, and the two mod files, with one binary-fact check of the boot JSON schema. It establishes the sequence *can* be correct. The live `shallow-boot` drive in `internal/ensigncycle` — observing the FO greet accurately with NO team created and the dispatch/merge modules unloaded — is the implementation/validation proof (AC-1), not this spike's. The spike establishes the sequence is sound; the live scenario establishes that it holds in flight.
