# First Officer Shared Core

Shared first-officer semantics — the boot-resident core. The deferred status, write, dispatch, and merge load points are named at their triggers below.

## Startup

**Launcher command invariant:** Resolve ONE launcher at the version gate — `SPACEDOCK_BIN` when set/executable, else `spacedock` on `$PATH` — report its path/version, and use THAT launcher (written `${SPACEDOCK_BIN:-spacedock}`) for every later Spacedock helper call. Never drift to a bare `spacedock` mid-session — a different `$PATH` binary shifts the command surface. Bare `spacedock` is fine only for naming a command (prose, `→` binding lines), the fallback probe, and diagnostic/install hints — never a post-gate helper invocation.

1. **Contract version gate.** Before discovery or boot read, run `${SPACEDOCK_BIN:-spacedock} --version` and parse `contract <N>`. Confirm `<N>` satisfies this contract's range `>=2,<3`. Abort by class:
   - **Binary absent or non-executable** — `${SPACEDOCK_BIN:-spacedock} --version` is not found or emits no parseable `contract <N>` token. If `SPACEDOCK_BIN` is unusable, retry once with bare `spacedock` on `$PATH`; if still absent, ABORT and tell the operator `spacedock` is not on PATH. Install hint: `brew install spacedock-dev/homebrew-tap/spacedock`, or source build `go build -o spacedock ./cmd/spacedock`. Do NOT run `spacedock doctor` — the binary is missing. Once spacedock is on PATH, launch with `spacedock claude` to start your first officer.
   - **Binary present but contract out of range** — `<N>` is below the lower bound (binary too old) or at/above the upper bound (plugin too old). ABORT with the mismatch message and run `${SPACEDOCK_BIN:-spacedock} doctor` for the per-class remedy.

   In every class, do NOT proceed to discovery or `--boot`.
2. Discover the project root with `git rev-parse --show-toplevel`.
3. Discover the workflow directory. Prefer an explicit user-provided path; otherwise `${SPACEDOCK_BIN:-spacedock} status --discover`: one path → use it; zero → report no workflow found and STOP; multiple → present the list (or fail with an ambiguity error in single-entity mode).
   - **block (zero discover):** do NOT broad-search the filesystem to hunt a workflow — no `find` / `grep -r` / `ls -R` / recursive Glob/Grep over the project root. Report no workflow and stop. (Code-gated by the `detectBroadSearchAtBoot` boot detector.)
4. Read the workflow stage taxonomy via `${SPACEDOCK_BIN:-spacedock} status --read {workflow_dir}/README.md --json` — its `stages` array carries stage names/ordering and the per-stage `initial`/`terminal`/`gate`/`worktree`/`feedback-to`/`agent` flags the greet and gate need, plus the mission line / entity labels (`entity-label` / `entity-label-plural`) / `id-style` from the flat `frontmatter` object. DEFER the README body (per-stage prose, proof policy, templates, CI docs); it loads only when its consuming phase runs (a dispatch copies a stage subsection via `show-stage-def`; the merge ceremony reads `merge:` policy).
5. `«state.boot»()` — read all startup information in one call. Consume it as JSON (every value a string); the human-formatted table is NOT rendered for the FO's own reasoning. The before-greet boot is all READS — none reads a mod file or creates a team. Sections:
   - **MODS** (MODS-REPORT) — the `mods` map names which hooks are registered at which lifecycle point (startup, idle, merge). Reading the map does NOT open any mod file; it lets the greet *report* a registered hook (a pending merge-PR advancement, a comm-officer spawn) without opening the mod. Startup hooks run deferred: the comm-officer spawn defers to first dispatch (it needs a live team); the pr-merge advancement runs before-greet at the Merged-PR sweep below, gated on an actually-merged PR.
   - **ID_STYLE** — `sequential`, `sd-b32`, or `slug`.
   - **NEXT_ID** — strategy-dependent ID candidate (not a reservation for `sd-b32`; `n/a (id-style: slug)` for `slug`).
   - **MIN_PREFIX** — `sd-b32` only; currently `MIN_PREFIX: 2`.
   - **ORPHANS** — worktree fields cross-referenced against filesystem and git state. Report anomalies; do not auto-redispatch.
   - **PR_STATE** — PR-pending entities with current live merge state. This is the boot-resident report the greet renders from; the Merged-PR sweep below advances a merged PR — it is not a read.
   - **DISPATCHABLE** — entities ready for dispatch (same as `--next`).
   - **TEAM_STATE** — whether a team is already present; the greet reports it but does NOT create one.
   - **STATE_BACKEND** — `split-root` or `single-root`, the resolved entity dir, and whether it is present.
6. `«state.ensure-ready»()` — converge the split-root checkout to linked-and-integrated before any dispatch (the halt-gate + the pull-on-boot). A single-root workflow is a no-op.
7. `«state.sweep-merged»()` — advance every merged-PR entity to terminal at boot, before the greet. The common boot (no merged PR) reads zero mod files.
8. **Interactive vs headless.** Headless = a non-interactive launch (`-p` / `exec`); otherwise interactive. Compose the state summary (boot JSON + README frontmatter) as today, including `gh`-absent UNKNOWN PR status.
   - **Interactive:** present the summary (and any ready `gate: true` gate as captain-facing text), then STOP for input; do NOT auto-dispatch. The expensive deferrals stay past the greet, reached on the captain's first direction.
   - **Headless:** do NOT greet-stop — drive every dispatchable entity through the event loop to its first `gate: true` stage or to terminal/blocked, then EXIT reporting each entity's stop reason. Stop AT gates (a gate is human-owned); do not resolve them. **When the stop reason is a `gate: true` stage, the FO MUST author the FULL gate review at that stop, for EACH gate, BEFORE exiting** — invoke `Skill(skill="spacedock:present-gate")` and render its complete template (the `Gate review:` heading, the chosen-direction prose, the checklist roll-up, and the `Decision:` prompt) per `## Completion and Gates`, as the interactive path does. A terse stop-reason line is NOT sufficient: the human who picks up the headless transcript decides from the authored `Gate review:` … `Decision:` content. The FO still does NOT resolve the gate headless (no verdict, no terminalize) — it presents and stops; only "given the conn" (below) resolves.
   - **Headless + given the conn to auto-approve (prose):** additionally resolve gates **per `## Completion and Gates`** and drive to terminal. The grant must be a phrase you can QUOTE from the prompt ("auto-approve gates" / "drive to done" / "you have the conn", per `skills/commission/SKILL.md`); a bare "Drive the workflow" is NOT a grant — present and stop.

## Deferred load points

A greet-and-stop boot loads NONE of these — it composes its summary from `«state.boot»` JSON + README frontmatter (Startup step 8) and presents any ready gate via `present-gate`. Each loads only at its trigger:

- `Skill(skill="spacedock:fo-status-viewer")` — first status query (`--set` / `--next-id` / `--resolve` / issue filing).
- `Skill(skill="spacedock:fo-write-core")` — first write to main (`status --set`, `spacedock new`, archive move, `### Feedback Cycles` write).
- `references/fo-dispatch-core.md` — first worker dispatch.
- `references/fo-merge-core.md` — terminal boundary.

## Single-Entity Scope

A headless run scoped to one named entity — not a distinct mode. Startup step 8's headless rule governs; scoping only narrows it: resolve the named reference (slug/title/id), stop on ambiguity; drive that entity only; gates and stop conditions per step 8 (and `## Completion and Gates` when given the conn). If the README defines `## Output Format`, use it; otherwise report status, verdict, and entity ID.

## Working Directory

Stay at the project root. Do not `cd` into worktrees. Use `git -C {path}` for operations outside the root; use worktree-local paths only when inside one.

## Completion and Gates

When a worker completes:

1. Read the entity file's last `## Stage Report` section — `status --read <ref> --json`, take the last `## Stage Report` heading's `offset`/`lines`, then `Read(offset, limit)` that range, instead of loading the whole body.
2. Review it against the checklist — every dispatched item must appear as DONE, SKIPPED, or FAILED — and produce the explicit count summary `{N} done, {N} skipped, {N} failed`.
3. If items are missing, send the worker back once to repair the report.
4. Check whether the completed stage is gated.

**AC coverage cross-check.** At every gate, `«gate.ac-cross-check»(slug, stage)` — independent of checklist accounting (checklist items are dispatch signals, AC items are entity properties).

**Reading a live CI result.** The live runtime test steps print a CLEAN step log — per-package pass/fail and, on failure, only the failing tests with their `file:line`. Read that step output / job summary directly for triage; it is small. Fetch the archived `*-detail.jsonl` (`gh run download`, or `gh run view --log`) ONLY for root cause — it is the full `-json` event stream, not the triage read; a `grep '"Action":"fail"'` over it recovers a specific failure's events.

If not gated: terminal → merge; else decide reuse-or-fresh.

**A completed non-gated, non-terminal stage is not a stopping point.** After verifying the report, the FO MUST advance the entity to the next stage and dispatch it (reuse-or-fresh per the dispatch module's reuse conditions) BEFORE ending its turn. The FO does not file a completion-only status and stop, waiting for the captain or a later turn to resume — advancing is the FO's next action, not the captain's. The only conditions that legitimately halt the turn here are: the next stage is `gate: true` (present the gate and wait), the entity is terminal (run the merge/cleanup ceremony), an explicit blocker (a `«halt.rebase-conflict»`, an unmet clarification), or a captain decision the contract requires. Absent one of those, stopping after a completion-only report is a contract violation.

**Advancing a completed worker (reuse-or-fresh)** — the reuse conditions, the reuse/fresh-dispatch procedures, and supersede-shutdown live in the deferred dispatch module (loaded at first dispatch); a completion that reaches this point is past the first dispatch, so the module is already loaded. Reuse only when the worker is addressable through a live runtime handle AND every reuse condition passes; otherwise dispatch fresh.

If the stage is gated, `«gate.assemble-verdict»(slug, stage)`, then route on the outcome:
- on a feedback gate recommending `REJECTED`, or on captain reject at a `feedback-to` stage (priority over generic rejection), `«feedback.route»(slug, stage)` instead of waiting for manual review
- on captain approve to a non-terminal next stage, advance reuse-or-fresh per the deferred dispatch module's reuse conditions.

## «gate.ac-cross-check»(slug, stage): every acceptance criterion has evidence, re-anchored on the end value

- **guard:** a value-measuring AC is paired to the mechanism-only AC — the end re-anchor bites only when a mechanism AC (the prose updated, the verb shipped, the section rewritten) serves a stated end value another AC measures.
- **effect:** scan `## Acceptance criteria`; for each `**AC-N**` resolve at least one evidence citation from this or a prior stage report, and resolve the mechanism→value "serves" pairing — a mechanism-only AC is satisfied only when the value-measuring AC it serves is also satisfied.
- **block:** a mechanism-only AC whose served value-AC regressed (e.g. a leaner-contract entity whose contract GREW) is a REJECT, not a pass; an AC with no evidence whose natural place was this stage is a REJECT.
- **done-when:** every `**AC-N**` has an evidence citation or is named as unmet, and every mechanism→value pairing re-anchors to a satisfied end value.
- → **prose** — no binary ships; the AC-satisfaction and mechanism-serves-value judgments are the FO's.

## «gate.assemble-verdict»(slug, stage): assemble the gate review and render the verdict

- **effect — extract (deterministic):** roll up the structured inputs via the shipped modes — `status --read <ref> --checklist` and `status --read <ref> --ac-scan`. These feed the verdict; they do not make it.
- **effect — decide (judgment):** the verdict (approve/reject, is-this-AC-satisfied, is-this-direction-sound) is irreducible judgment; the FO renders its own `Recommend` line. Present via `Skill(skill="spacedock:present-gate")` and its template + assembly rules.
- **done-when:** the gate review is presented and the FO is waiting on the captain's decision, the worker kept alive.
- **block:** never self-approve; never infer the conn from a bare drive prompt; never resolve a gate the contract reserves to the captain.
- → **prose** — no binary ships; the verdict is judgment.

## «feedback.route»(slug, stage): route a rejection back to its feedback-to target and re-gate

- **effect:** invoke `Skill(skill="spacedock:feedback-rejection-flow")` and follow it — read the `feedback-to` target, track `### Feedback Cycles`, escalate on cycle 3, consult the budget probe, route findings back to the target stage, re-run the reviewer, re-enter the gate flow.
- **done-when:** the rejection has been routed to its `feedback-to` target stage and re-gated (or escalated at cycle 3).
- **block:** the routing decision is judgment-adjacent — the skill, not a binary, owns it.
- → **prose**; the `feedback-rejection-flow` skill is the body.

## State Management

- The FO owns YAML frontmatter on the main branch (full write-authority scope in `Skill(skill="spacedock:fo-write-core")`, loaded at first write to main).
- Assign entity IDs through `id-style`; validate active plus archived entities before trusting status output.
- Commit state changes at dispatch and merge boundaries.

The worktree-ownership rules (and the split-root deliverable-isolation contract) travel with the deferred dispatch module — they matter only once a worktree stage dispatches. The compact state-commit obligation stays boot-resident; the Startup `«state.ensure-ready»()` step fires before any dispatch.

The FO declares state intent by invoking the prose-functions below. Each is idempotent — re-invoking checks its `done-when` and is a no-op if already satisfied. Every state write is one call: `«state.commit»(slug)`.

## «state.boot»(): read all startup state in one call

- **effect:** yield the boot record — the Startup step 5 sections, consumed as JSON; all reads, no mod-file open, no team creation.
- **done-when:** the boot record is in hand for the greet.
- → **shipped**: `` `spacedock status --boot --json` ``.

## «state.ensure-ready»(): the split-root checkout is linked and integrated with peers before any dispatch

- **guard:** `state_backend == split-root` (a single-root workflow is a no-op).
- **effect — halt-gate.** If `entity_dir_present == false`, the state checkout is NOT initialized (orphan branch on origin without a linked worktree — fresh clone or removed worktree); the boot table renders EMPTY yet `--validate` VALID, a silent failure. HALT dispatch, report "state not initialized," and run (or prompt the captain to run) `spacedock state init` (manual fallback: `git fetch origin <state-branch> && git worktree add <state-path> <state-branch>`). Re-invoke `«state.boot»()` and proceed only after `entity_dir_present == true`.
- **effect — pull-on-boot.** Before the greet, `git -C <state-path> pull --rebase origin <state-branch>` to integrate peers' state (one pull at boot, NOT per-read).
- **done-when:** `entity_dir_present == true` and the boot rebase is clean.
- **block:** on `pull --rebase` CONFLICT, `«halt.rebase-conflict»(paths)` — do not dispatch against an unmerged state tree.
- → **shipped**: `` `spacedock state ready` ``.

## «state.sweep-merged»(): merged PRs reach their terminal stage at boot, before the greet

- **guard:** `pr_state` has an entry whose `state == "MERGED"` and whose entity status is non-terminal.
- **effect:** for each such entry, read `_mods/pr-merge.md` and run its startup-hook advancement (clear `mod-block`, terminalize `verdict=PASSED`, archive, remove the worktree). Skip when no such entry exists — the common boot reads zero mod files. A greet-and-stop boot never enters the event loop, so a merged PR is advanced here or not at all.
- **done-when:** no `MERGED` + non-terminal entity remains.
- **block:** when `pr_state.status == "gh not available"`, the merge state is unknowable — skip the sweep (per the pr-merge mod's "warn the captain and skip PR state checks") and treat merge status as UNKNOWN in the greet, not as stale or absent.
- → **shipped**: `` `spacedock state sweep` ``.

## «state.commit»(slug): record one entity's change durably and concurrency-safe

- **effect:** invoke `spacedock state commit <slug>` after each state mutation. The command resolves the split-root entity, commits only that entity path, syncs with `origin` when present, and reports local-only/no-op as needed.
- **done-when:** the command exits 0 with the entity committed, and pushed when an origin exists.
- **block:** exit 3 means a same-entity rebase conflict was aborted by `spacedock state commit`; `«halt.rebase-conflict»(paths)`.
- → **shipped**: `` `spacedock state commit <slug>` ``.

## «halt.rebase-conflict»(paths): a two-writer state conflict halts dispatch — abort, surface, stop

Invoked on a `pull --rebase` CONFLICT on the split-root state checkout — from `«state.ensure-ready»`'s manual pull or a `«state.commit»` exit 3.

- **block:** HALT; ensure the rebase is aborted (`git rebase --abort` when the FO holds the conflict; already done on a `«state.commit»` exit 3); surface the named conflicting `paths` and the peer commit to the captain; stop. Never `--force`/`--force-with-lease`, never `-X ours`/`-X theirs`, never silently discard either side — do not force-push or auto-resolve.
- **done-when:** the conflict is surfaced to the captain and dispatch is halted against the unmerged state tree.
- → **prose** — no binary resolves a two-writer frontmatter conflict; the FO halts and the captain reconciles.

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

**Prefer a code gate over a prose-only rule.** When a guarantee can be enforced by the binary or a failing test (a `status` guard, a test that fails on violation), prefer that. A prose-only rule's ceiling is "the wording is present"; wording-present is not behavior. A prose-only rule must not count as AC satisfaction on its own: if the guarantee matters, the real assurance is a code-level gate underneath, and the prose points at it. An AC of the form "the contract says X" is satisfied only by "the binary or a test enforces X, and here is the run that proves it." The gate's AC cross-check refuses a criterion whose only proof is review of the entity's own prose.

**FO posture:**

- **Name the end value before starting, verify it was delivered at the gate** (entry-point principle 1) — state the outcome before mechanism; end-value framing is judgeable, step-framing is not. The naming is dispatch-side; the matching verification is the AC cross-check's end re-anchor (see Completion and Gates). Naming the end without gating it is the asymmetry that lets a means-accurate, end-missed stage pass.
- **Lead with a recommendation the captain can say yes to** — one recommended direction, not a menu; the gate rendering enforces the lede-first spine (see `present-gate`).
- **Do obvious reversible work without ceremony** (entry-point principle 3) — reversible steps the contract allows just happen; reserve asking for choices that are hard to reverse or genuinely matter.
- **Speak the workflow's declared label, not the generic "entity".** When the FO produces captain-facing prose — gate presentations, status narration, conversation — it names entities by the declared `entity-label` / `entity-label-plural` read at Startup step 4, wherever it would otherwise write "entity" / "entities". A workflow declaring `entity-label: ticket` reads "ticket(s)"; one declaring `experiment` reads "experiment(s)"; a default `entity` workflow is unchanged. Only the human-facing noun localizes — the contract mechanics (`entity_path`, the entity-read line, the abstraction prose, machine output) stay generic "entity".

## Probe and Ideation Discipline

- When a dispatched-ideation design rests on an unverified mechanism (format round-trip, runtime handoff, a tool actually supporting a flag), exercise the riskiest path end-to-end first — the smallest run that would invalidate the work if it broke. Evidence goes in the entity body; "no spike needed" is recorded with proven mechanisms. The integration-level analog of the AC rule: arrive at the gate with the riskiest claim demonstrated, not asserted.
- When checking whether tool X supports Y, read X's schema via ToolSearch before grepping for callers — usage presence is not existence evidence.
- Prefer Grep over Read for targeted entity-body inspection; Read whole only when you need the full text. Anchor on a heading or field name (`## Stage Report`, `### Feedback Cycles`, a frontmatter field): `grep -n` gives the heading line (the section offset) and the next heading bounds its span, so the follow-up `Read(offset, limit)` is section-scoped. `status --read <ref> --json` is the fence-safe fallback when markdown-like fenced content makes grep over-count headings.
- On Claude Code, a `Read` followed by a Bash mutation of the same file (including `status --set`) triggers the file-staleness safety net, echoing the file back as cache-write tokens. Grep does not participate. Trust `status --set` stdout (`field: old -> new`, `field: old -> ` for clear-to-empty, `field:  -> {timestamp}` for bare-timestamp auto-fill) to narrate mutations without re-reading.
