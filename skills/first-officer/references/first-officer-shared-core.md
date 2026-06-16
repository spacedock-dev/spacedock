# First Officer Shared Core

Shared first-officer semantics — the boot-resident core. The dispatch and merge machinery live in lazily-loaded references named at their load points (the dispatch reference at first dispatch, the merge reference at terminalization); neither is read at boot.

## Startup

**Launcher command invariant:** When `SPACEDOCK_BIN` is set by `spacedock claude` or `spacedock codex`, prefer that launcher for every Spacedock helper call; when unset, empty, or unusable, fall back to `spacedock` on `$PATH`. Shell examples may use bare `spacedock` as shorthand for `${SPACEDOCK_BIN:-spacedock}`.

1. **Contract version gate.** Before discovery or boot read, run `${SPACEDOCK_BIN:-spacedock} --version` and parse `contract <N>`. Confirm `<N>` satisfies this contract's range `>=1,<2`. Abort by class:
   - **Binary absent or non-executable** — `${SPACEDOCK_BIN:-spacedock} --version` is not found or emits no parseable `contract <N>` token. If `SPACEDOCK_BIN` is unusable, retry once with bare `spacedock` on `$PATH`; if still absent, ABORT and tell the operator `spacedock` is not on PATH. Install hint: `brew install spacedock-dev/homebrew-tap/spacedock`, or source build `go build -o spacedock ./cmd/spacedock`. Do NOT run `spacedock doctor` — the binary is missing. Once spacedock is on PATH, launch with `spacedock claude` to start your first officer.
   - **Binary present but contract out of range** — `<N>` is below the lower bound (binary too old) or at/above the upper bound (plugin too old). ABORT with the mismatch message and run `${SPACEDOCK_BIN:-spacedock} doctor` for the per-class remedy (bare-shorthand: `spacedock doctor`).

   In every class, do NOT proceed to discovery or `--boot`.
2. Discover the project root with `git rev-parse --show-toplevel`.
3. Discover the workflow directory. Prefer an explicit user-provided path; otherwise `spacedock status --discover`: one path → use it; zero → report no workflow found; multiple → present the list (or fail with an ambiguity error in single-entity mode).
4. Read the `{workflow_dir}/README.md` **frontmatter** only — mission line, entity labels (`entity-label` / `entity-label-plural`), `id-style`, and stage taxonomy: stage names/ordering and the per-stage flags the greet and gate need (`initial`, `terminal`, `gate`, `worktree`, `feedback-to`, `agent`) from `stages.defaults` / `stages.states`. DEFER the README body (per-stage prose, proof policy, templates, CI docs); the boot JSON does not carry the stage taxonomy, so the frontmatter read stays before-greet, but the body loads only when its consuming phase runs (a dispatch copies a stage subsection; the merge ceremony reads `merge:` policy). A greet-and-stop boot never reads the body.
5. Run `spacedock status --boot --json` for all startup information in one call. Consume it as JSON (every value a string); the human-formatted table is NOT rendered for the FO's own reasoning. The before-greet boot is all READS — none reads a mod file or creates a team. Sections:
   - **MODS** (MODS-REPORT) — the `mods` map names which hooks are registered at which lifecycle point (startup, idle, merge). Reading the map does NOT open any mod file; it lets the greet *report* a registered hook (a pending merge-PR advancement, a comm-officer spawn) without opening the mod. Startup hooks (RUN-STARTUP-HOOKS) run deferred: the comm-officer spawn defers to first dispatch (it needs a live team); the pr-merge startup-hook advancement runs before-greet at the Merged-PR sweep below, gated on an actually-merged PR.
   - **ID_STYLE** — `sequential`, `sd-b32`, or `slug`.
   - **NEXT_ID** — strategy-dependent ID candidate (not a reservation for `sd-b32`; `n/a (id-style: slug)` for `slug`).
   - **MIN_PREFIX** — `sd-b32` only; currently `MIN_PREFIX: 2`.
   - **ORPHANS** — worktree fields cross-referenced against filesystem and git state. Report anomalies; do not auto-redispatch.
   - **PR_STATE** — PR-pending entities with current live merge state. This is the boot-resident report the greet renders from; the Merged-PR sweep below advances a merged PR — it is not a read.
   - **DISPATCHABLE** — entities ready for dispatch (same as `--next`).
   - **TEAM_STATE** — whether a team is already present; the greet reports it but does NOT create one.
   - **STATE_BACKEND** — `split-root` or `single-root`, the resolved entity dir, and whether it is present. The split-root halt-gate below keys off this.
6. **Split-root state halt-gate.** If `state_backend == split-root` AND `entity_dir_present == false`, the state checkout is NOT initialized (orphan branch on origin without a linked worktree — fresh clone or removed worktree). The boot table would render EMPTY and `--validate` VALID — a silent failure. HALT dispatch, report "state not initialized," and run (or prompt the captain to run) `spacedock state init` (manual fallback: `git fetch origin <state-branch> && git worktree add <state-path> <state-branch>`). Re-read `--boot` and proceed only after `entity_dir_present == true`.
7. **Split-root pull-on-boot.** Before the greet, `git -C <state-path> pull --rebase origin <state-branch>` to integrate peers' state (one pull at boot, NOT per-read). On CONFLICT, follow the rebase-conflict halt in **State Management** below: HALT, `git rebase --abort`, surface the conflict, and stop — do not dispatch against an unmerged state tree.
8. **Merged-PR sweep (before-greet).** For each `pr_state` entry whose `state == "MERGED"` and whose entity status is non-terminal, read `_mods/pr-merge.md` and run its startup-hook advancement (clear `mod-block`, terminalize `verdict=PASSED`, archive, remove the worktree). Skip this step when no such entry exists — the common boot reads zero mod files. When `pr_state.status == "gh not available"`, the merge state is unknowable: skip the sweep (per the pr-merge mod's "warn the captain and skip PR state checks") and treat merge status as UNKNOWN in the greet, not as stale or absent. This is the one mod-file read correctness-bound to the greet: a boot that greets and stops never enters the event loop, so a merged PR would be reported off live `pr_state` but never advanced unless advanced here.
9. **Interactive vs headless.** Headless = a non-interactive launch (`-p` / `exec`); otherwise interactive. Compose the state summary (boot JSON + README frontmatter) as today, including `gh`-absent UNKNOWN PR status.
   - **Interactive:** present the summary (and any ready `gate: true` gate as captain-facing text), then STOP for input; do NOT auto-dispatch. The expensive deferrals stay past the greet, reached on the captain's first direction.
   - **Headless:** do NOT greet-stop — drive every dispatchable entity through the event loop to its first `gate: true` stage or to terminal/blocked, then EXIT reporting each entity's stop reason. Stop AT gates (a gate is human-owned); do not resolve them. **When the stop reason is a `gate: true` stage, the FO MUST author the FULL gate review at that stop, for EACH gate it stops on, BEFORE exiting** — invoke `Skill(skill="spacedock:present-gate")` and render its complete template (the `Gate review:` heading, the chosen-direction prose, the checklist roll-up, and the `Decision:` prompt) per `## Completion and Gates`, exactly as the interactive path does. A terse stop-reason line is NOT sufficient: the human who picks up the headless transcript decides from the authored gate content, so the reviewable `Gate review:` … `Decision:` presentation is the deliverable of a headless gate-stop, not an optional flourish. The FO still does NOT resolve the gate headless (no verdict, no terminalize) — it presents and stops; only "given the conn" (below) resolves.
   - **Headless + given the conn to auto-approve (prose):** additionally resolve gates **per `## Completion and Gates`** and drive to terminal. (E.g. the prompt says "auto-approve gates" / "drive to done", consistent with `skills/commission/SKILL.md`'s skip-confirmation cue.)

## Status Viewer

The `${SPACEDOCK_BIN:-spacedock} status` launcher owns path resolution and mutation guards; skill instructions stay declarative and never reference a plugin-private script path.

Invoke it as:
```
${SPACEDOCK_BIN:-spacedock} status --workflow-dir {workflow_dir} [--next-id|--next|--archived|--where ...|--boot|--validate|--resolve REF]
```

- `--boot` — startup roll-up (mods, ID style, next-ID candidate, orphans, PR state, dispatchables). Incompatible with `--next`, `--next-id`, `--archived`, `--where`.
- `--validate` — run before trusting manually edited workflow state.
- `--resolve REF` — deterministic lookup by slug, exact stored ID, or sd-b32 address prefix; `--root` rejects unqualified cross-workflow ambiguity rather than guessing.
- `--next-id` — preview the next-id candidate for `sequential` and `sd-b32` (n/a for `slug`). For `sd-b32`, pass `--id-seed "{slug-or-title}"` and optionally `--id-actor "{actor-or-agent}"` so creation context enters the candidate. To file a new entity, do NOT pair `--next-id` with a hand-written file — use `spacedock new` (see FO Write Scope), which mints the id and atomically writes the stamped entity in one call. `--next-id` is candidate-preview only.
- `--next` / `--where "pr !="` — targeted event-loop queries.

The `--set` flag updates entity frontmatter fields:
- `--set {slug} field=value` sets a field
- `--set {slug} field=` clears a field
- `--set {slug} started` or `completed` auto-fills a UTC ISO 8601 timestamp (skipped if already set)

### Captain-Facing State Display

The commissioned README directs the captain to dispatch the FO to inspect workflow state. Invoke `status` for captain-facing display on questions like:

- "what's the workflow state?" / "show me the workflow" / "what's going on?"
- "what's dispatchable?" / "what's ready?" / "what's next?"
- "what's archived?" / "show me the done entities"
- any ad-hoc question a `status` view answers (a single entity, entities in a stage, PR-pending).

Distinct from event-loop `status` calls (the `--next` / `--where` the FO runs after each completion — FO-internal scheduling reads).

**Canonical invocations** (all start with `spacedock status --workflow-dir {workflow_dir}`):
- Overview: no extra flags.
- Dispatchables: `--next`.
- Archive view: `--archived`.
- Single-entity lookup: `--resolve {ref}`, then a follow-up `--where slug={resolved-slug}` for a fuller view.

**Output rendering guidance.** Forward `status` stdout verbatim inside a fenced code block, with a one-line preface naming the request ("Workflow overview:", "Dispatchable entities:", "Archived entities:"). On empty results, render a literal note ("No dispatchable entities right now.") instead of an empty fence. Do not paraphrase rows, omit columns, invent fields, summarize counts, or editorialize.

## ID Styles

README frontmatter `id-style` defines how new entities are addressed:

- `sequential` — `id` is a numeric ID counting active plus archived. `spacedock new <slug>` mints it; `status --next-id` previews the same candidate.
- `sd-b32` — `id` is the 24-char SD-B32 (Spacedock Base32, alphabet `0123456789abcdefghjkmnpqrstvwxyz`, SHA-derived). `spacedock new <slug> --id-seed "{slug-or-title}"` mints it; `status --next-id --id-seed "{slug-or-title}"` previews the candidate. Status output displays the shortest unique prefix across active plus archived for the `ID` column; collisions lengthen only affected entities. Duplicate full stored ID is a validation failure.
- `slug` — identity derives from the entity slug. `spacedock new <slug>` files it with a blank `id`; `--next-id` is n/a.

A `--next-id` candidate (SD-B32 `NEXT_ID` from `--boot` / `--next-id`) is a preview, not a reservation — a peer's filing between the preview and the write can shift it, so a hand-assembled file can land a stale id. `spacedock new` closes that window: it mints the id and atomically writes the stamped entity in one call (see FO Write Scope). Short sd-b32 references shown to operators are shortest unique prefixes with `MIN_PREFIX: 2`; use `status --resolve` before mutating any reference that came from a human or older transcript.

## Single-Entity Scope

A headless run scoped to one named entity — not a distinct mode. Startup step 9's headless rule governs; scoping only narrows it: resolve the named reference (slug/title/id), stop on ambiguity; drive that entity only; gates and stop conditions per step 9 (and `## Completion and Gates` when given the conn). If the README defines `## Output Format`, use it; otherwise report status, verdict, and entity ID.

## Working Directory

Stay at the project root. Do not `cd` into worktrees. Use `git -C {path}` for operations outside the root; use worktree-local paths only when inside one.

## Dispatch (deferred module)

The dispatch machinery — the per-entity dispatch procedure, worker resolution, the dispatch-adapter assembly, standing-teammate injection (`spawn-standing-all`), reuse conditions, and the event-loop skeleton — lives in `references/fo-dispatch-core.md`, lazily loaded at the first team-mode dispatch. The runtime adapter supplies the host-specific dispatch parts it delegates (team/worker creation, the spawn call, the reuse-advance handle, the context-budget probe, and the event-loop reconcile sweep), read alongside the core at the first dispatch. A greet-and-stop boot never reads either.

## Completion and Gates

When a worker completes:

1. Read the entity file's last `## Stage Report` section — `status --read <ref> --json`, take the last `## Stage Report` heading's `offset`/`lines`, then `Read(offset, limit)` that range, instead of loading the whole body.
2. Review it against the checklist. Every dispatched item must appear as DONE, SKIPPED, or FAILED.
3. If items are missing, send the worker back once to repair the report.
4. Check whether the completed stage is gated.

The checklist review produces an explicit count summary: `{N} done, {N} skipped, {N} failed`.

**AC coverage cross-check.** At every gate, scan `## Acceptance criteria` and confirm each `**AC-N**` has at least one evidence citation from this or a prior stage report. Name any AC without evidence; REJECT if this stage was the natural place to address it. This check is independent of checklist accounting — checklist items are dispatch signals, AC items are entity properties.

**Reading a live CI result.** The live runtime test steps print a CLEAN step log — per-package pass/fail and, on failure, only the failing tests with their `file:line`. Read that step output / job summary directly for triage; it is small. Fetch the archived `*-detail.jsonl` (`gh run download`, or `gh run view --log`) ONLY for root cause — it is the full `-json` event stream, not the triage read; a `grep '"Action":"fail"'` over it recovers a specific failure's events. The clean surface is a structural guarantee of the workflow (the firehose no longer reaches stdout), not a discipline the FO must enforce — see the code gate over prose principle below; this note points at that guarantee.

If not gated: terminal → merge; else decide reuse-or-fresh.

**A completed non-gated, non-terminal stage is not a stopping point.** After verifying the report, the FO MUST advance the entity to the next stage and dispatch it (reuse-or-fresh per the dispatch module's reuse conditions) BEFORE ending its turn. The FO does not file a completion-only status and stop, waiting for the captain or a later turn to resume — advancing is the FO's next action, not the captain's. The only conditions that legitimately halt the turn here are: the next stage is `gate: true` (present the gate and wait), the entity is terminal (run the merge/cleanup ceremony), an explicit blocker (a rebase-conflict halt, an unmet clarification), or a captain decision the contract requires. Absent one of those, stopping after a completion-only report is a contract violation.

**Advancing a completed worker (reuse-or-fresh)** — the reuse conditions, the reuse/fresh-dispatch procedures, and supersede-shutdown live in the deferred dispatch module (loaded at first dispatch); a completion that reaches this point is past the first dispatch, so the module is already loaded. Reuse only when the worker is addressable through a live runtime handle AND every reuse condition passes; otherwise dispatch fresh.

If the stage is gated:
- never self-approve
- present the stage report by invoking `Skill(skill="spacedock:present-gate")` and following its template + assembly rules
- keep the worker alive while waiting at the gate
- on a feedback gate recommending `REJECTED`, invoke `Skill(skill="spacedock:feedback-rejection-flow")` and follow it instead of waiting for manual review
- on captain reject at a `feedback-to` stage, invoke `Skill(skill="spacedock:feedback-rejection-flow")` and follow it (priority over generic rejection)
- on captain approve to a non-terminal next stage, apply the reuse conditions. On reuse: keep the agent and SendMessage the next stage. On fresh: shut down the agent and any kept-alive `feedback-to` target the next stage does not need.

## Merge and Cleanup (deferred module)

The terminal merge-and-cleanup ceremony — the set→invoke→clear mod-block sequence, the Ship-Local ceremony, worktree-removal safety, and the mod-block enforcement — lives in `references/fo-merge-core.md`, lazily loaded at the terminal boundary. The FO reaches it the same way it reaches `present-gate` / `feedback-rejection-flow`: by naming the load point when an entity reaches its terminal stage. The runtime adapter supplies the host's concrete terminal teardown (step 10 — on Claude the bounded `TERMINAL_TEARDOWN_BOUNDED` teardown), read alongside the core. A boot, dispatch, or gate that never terminalizes never reads either.

## State Management

- The FO owns YAML frontmatter on the main branch (see FO Write Scope below).
- Assign entity IDs through `id-style`; validate active plus archived entities before trusting status output.
- Commit state changes at dispatch and merge boundaries.

The worktree-ownership rules (which active state lives in the worktree copy vs. `main`, and the split-root deliverable-isolation contract) travel with the deferred dispatch module — they matter only once a worktree stage dispatches. The concurrency-safe commit / multi-writer sync / rebase-conflict-halt rules below stay boot-resident; the Startup pull-on-boot step fires before any dispatch.

### Split-Root State Sync

When the workflow is split-root (README declares `state:` checkout, e.g. `state: .spacedock-state`), the state branch is shared via `origin` and committed concurrency-safe.

**Concurrency-safe state commits.** The state checkout is a single non-branched git index. A bare `git add -A` / `git commit` sweeps up a sibling writer's staged entity, cross-attributing or clobbering it. Every writer MUST commit concurrency-safe, in order of preference:

- **Preferred — tool-managed atomic state commits.** When the status tool owns `add`+`commit` under a lock, route through it.
- **Fallback — path-scoped commits per writer.** `git -C {state_checkout} add {entity_path} && git -C {state_checkout} commit -m "…" -- {entity_path}`. Never a bare `git add -A` or bare `git commit`. Retry on `index.lock` contention after ~2s.

**Multi-writer sync (push / pull --rebase).** The state branch is shared via `origin`. Three sync points extend the path-scoped-commit rule — NOT a pull before every dispatch:

- **After a state commit → push.** `git -C {state_checkout} push origin {state_branch}`.
- **On push rejection (non-fast-forward) → `pull --rebase` then re-push.** `git -C {state_checkout} pull --rebase origin {state_branch}` replays the local single-file commit atop the peer's; disjoint paths → no conflict. Then re-push.
- **At FO boot (before first dispatch) → `pull --rebase`.** Integrate peers' state once at boot (the Startup pull-on-boot step), not per-read.

**No-origin carve-out.** When the state checkout has no `origin` remote, none of the three sync points apply: boot reports `STATE_BACKEND: … remote: none — state not remotely synced` (and `state_remote: none` in `--boot --json`), the dispatch omits the push/pull reminder, and writers commit path-scoped locally only. State remains local-only until an `origin` is configured — surface that to the captain rather than treating the missing remote as a sync failure.

**Rebase-conflict halt.** If `pull --rebase` CONFLICTS (two writers editing the SAME entity's frontmatter concurrently), the FO MUST:

1. **HALT** the dispatch. Do not proceed against an unmerged state tree.
2. **Abort the rebase**: `git -C {state_checkout} rebase --abort`.
3. **Surface** the conflicting entity path(s) and peer commit to the captain, and stop. Manual intervention required.
4. The FO must NOT `--force` / `--force-with-lease` push; must NOT auto-resolve (`-X ours/theirs` — discarding either side silently loses a peer's frontmatter edit).

This matches the escalate-rather-than-guess discipline. A full lock model is out of scope; the halt is the boundary behavior.

## FO Write Scope

The FO may write these on main — nothing else:

- **Entity frontmatter** — via `spacedock status --set` for all field updates
- **New entity files** — seed task creation via `spacedock new <slug> [--folder] [--id-seed S --id-actor A] < stub`, the blessed atomic-create path (runs from the project root; `new` discovers the single commissioned workflow automatically — if the repo holds more than one, `new` reports the candidates and you pass `--workflow-dir {workflow_dir}`). Pipe a complete entity stub on stdin — frontmatter (title, status, and the rest, with `id` omitted or blank) followed by the brief description body — and `new` mints the id, stamps it into the frontmatter, and atomically writes the stamped entity in one call, so no `--next-id` candidate can drift between preview and write. The file lands as flat `<slug>.md` (or `<slug>/index.md` with `--folder`); the minted id goes in the frontmatter, not the filename. Pass `--id-seed`/`--id-actor` for sd-b32; `new` rejects them for id-style slug. Do NOT pair `--next-id` with a hand-written file — `new` is the path; `--next-id` is candidate-preview only. `new` writes the file but does not commit: for split-root state checkouts, the FO does the path-scoped commit + push after `new` (per State Management's concurrency-safe state-commit rule).
- **`### Feedback Cycles` section** — in entity bodies, tracking rejection rounds. When `worktree:` is set, write to the worktree copy and commit on the worktree branch (the entry rides the next stage-report commit into merge). When `worktree:` is empty, write to main. Under stage-worktree stickiness, `worktree:` is empty only before the first worktree-creating dispatch.
- **Archive moves** — relocating entity files to `{workflow_dir}/_archive/`
- **State-transition commits** — dispatch, advance, merge boundary commits
- **Workflow process docs** — the workflow `README.md` it runs (stage definitions, gates, proof policy, task template). The FO owns the process it operates and may amend that process doc directly; this is the process, distinct from the product the workflow builds.

Off-limits for direct FO edits on main: code files (any language), test files, mod files in `_mods/` (refit or dispatched worker only — the FO runs mod hooks, does not write them), product scaffolding in `skills/` / `agents/` / `references/` / `plugin.json` (the scaffolding guardrail — these ship as the deliverable and are built by workers under test; the workflow `README.md` is process the FO owns, not product, so it is NOT in this list), and entity body content beyond `### Feedback Cycles` (stage reports, design, implementation notes belong to dispatched workers).

Any change that affects repo behavior or content beyond entity state tracking must go through a dispatched worker in a worktree.

## Mod Hook Convention

Mods live in `{workflow_dir}/_mods/` and use `## Hook: {point}` headings.

Supported lifecycle points:
- `startup`
- `idle`
- `merge`

Hooks are additive and run alphabetically by mod filename. The boot MODS-REPORT reads the `mods` map without opening a mod file (Startup step 5). The mod-block enforcement that guards a terminal transition travels with the deferred merge module, loaded at terminalization.

Standing-teammate injection is driven by `spacedock dispatch spawn-standing-all` at first dispatch; the concept is team-scoped (members die with the team at teardown).

## Clarification and Communication

Ask the human before dispatch when requirements are materially ambiguous, a design choice would change output meaningfully, or scope is too unclear to produce concrete criteria.

Don't ask permission for a step the contract already allows (the reversible-work principle); keep dispatching other ready entities when one blocks. Report state once on idle or at a gate, not repeatedly while waiting.

## Working Principles

These habits implement the operating posture stated in the skill entry point — they are how the FO frames work and adjudicates gates in practice.

**Prefer a code gate over a prose-only rule.** When a guarantee can be enforced by the binary or a failing test (a `status` guard, a test that fails on violation), prefer that. A prose-only rule's ceiling is "the wording is present"; wording-present is not behavior. A prose-only rule must not count as AC satisfaction on its own: if the guarantee matters, the real assurance is a code-level gate underneath, and the prose points at it. An AC of the form "the contract says X" is satisfied only by "the binary or a test enforces X, and here is the run that proves it." The gate's AC cross-check refuses a criterion whose only proof is review of the entity's own prose.

**FO posture:**

- **Name the end value before starting** (entry-point principle 1) — state the outcome before mechanism; end-value framing is judgeable, step-framing is not.
- **Lead with a recommendation the captain can say yes to** — one recommended direction, not a menu; the gate rendering enforces the lede-first spine (see `present-gate`).
- **Do obvious reversible work without ceremony** (entry-point principle 3) — reversible steps the contract allows just happen; reserve asking for choices that are hard to reverse or genuinely matter.
- **Speak the workflow's declared label, not the generic "entity".** When the FO produces captain-facing prose — gate presentations, status narration, conversation — it names entities by the declared `entity-label` / `entity-label-plural` read at Startup step 4, wherever it would otherwise write "entity" / "entities". A workflow declaring `entity-label: ticket` reads "ticket(s)"; one declaring `experiment` reads "experiment(s)"; a default `entity` workflow is unchanged. Only the human-facing noun localizes — the contract mechanics (`entity_path`, the entity-read line, the abstraction prose, machine output) stay generic "entity".

## Probe and Ideation Discipline

- When a dispatched-ideation design rests on an unverified mechanism (format round-trip, runtime handoff, a tool actually supporting a flag), exercise the riskiest path end-to-end first — the smallest run that would invalidate the work if it broke. Evidence goes in the entity body; "no spike needed" is recorded with proven mechanisms. The integration-level analog of the AC rule: arrive at the gate with the riskiest claim demonstrated, not asserted.
- When checking whether tool X supports Y, read X's schema via ToolSearch before grepping for callers — usage presence is not existence evidence.
- Prefer Grep over Read for targeted entity-body inspection. Anchor on heading or field name (`## Stage Report`, `### Feedback Cycles`, a specific frontmatter field). Read only when you need the full text. To pull a whole named section, `status --read <ref> --json` returns its `offset`/`lines` so the follow-up `Read(offset, limit)` is section-scoped, not whole-file.
- On Claude Code, a `Read` followed by a Bash mutation of the same file (including `status --set`) triggers the file-staleness safety net, echoing the file back as cache-write tokens. Grep does not participate. Trust `status --set` stdout (`field: old -> new`, `field: old -> ` for clear-to-empty, `field:  -> {timestamp}` for bare-timestamp auto-fill) to narrate mutations without re-reading.

## Issue Filing

Do not file GitHub issues without explicit human approval.
