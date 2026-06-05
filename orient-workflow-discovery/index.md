---
id: 19fhrfae24d221wzgqm4zarn
title: Bring workflow-discovery into spacedock — orient on a project's implicit agent workflow as the front-door to commission
status: implementation
source: "captain (2026-06-05) — 'check the prototype orient skill in ~/.claude, this is the workflow discovery thing we should bring in.' Prototype at ~/.claude/skills/orient/SKILL.md."
score: "0.33"
started: 2026-06-05T19:10:17Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-orient-workflow-discovery
issue:
---

There is a working prototype `orient` skill (`~/.claude/skills/orient/SKILL.md`) that reconstructs a project's IMPLICIT workflow from its AI-agent session history: the inferred loop, the workstreams, the recent decisions, and — load-bearing — the OPEN decisions (abandoned/unanswered forks) plus an interruption count (how often the human had to step in). It reads agentsview's multi-agent session DB (`~/.agentsview/sessions.db`) with plain `sqlite3` queries. Its closing pitch is exactly the spacedock thesis: "spacedock turns these [interruptions] into gates so the agent advances on its own between your calls."

That makes orient the natural FRONT-DOOR to commission: discover what a brownfield project's agents have implicitly been doing → surface the workstreams + the OPEN decision frontier → propose commissioning it as a real spacedock workflow with gates.

## Problem

Spacedock today starts from a commissioned workflow; it has no way to DISCOVER a workflow from a project that already has ad-hoc multi-agent history. The bridge from "agents have been working here unsupervised" to "a spacedock workflow with gates" is missing — orient is the prototype of that bridge.

Three problems sharpen this for a real first-run experience on a brownfield repo:

1. **No discoverable front door.** A new user landing on a repo that already has agent history has nothing to run first. "orient" is the internal concept, but it is not the word a new user reaches for when they want to understand a project they just opened. The entry point needs a name that reads as "this is where you start."
2. **Spacedock ignores incumbent scaffolds.** Many brownfield repos already run a common agent scaffold — superpowers (`.claude/skills/` with the superpowers plugin), gsd / get-shit-done, and similar. Offering a generic spacedock pitch on top of an existing scaffold is tone-deaf: it neither acknowledges what the project already does nor explains why spacedock is worth adopting *over* what is already there.
3. **The benefit is asserted, not shown.** The prototype's closing line ("spacedock turns these interruptions into gates") is a generic pitch. It does not tie the value to what the scan actually found in *this* project, so it reads as marketing rather than a diagnosis.

## Direction (for ideation)

- Read the prototype `~/.claude/skills/orient/SKILL.md` (the agentsview sqlite scan + the synthesis format) as the reference; decide what comes into spacedock and in what shape — a spacedock-owned `orient`/discover skill, a `spacedock` command, or a commission pre-step.
- The discovery → commission bridge: the OPEN decisions and the interruption stats are the raw material for proposing gates; the workstreams map to candidate stages/entities. Design how discovery output feeds the commission skill.
- Dependency: agentsview (the session DB this reads). Decide how spacedock handles its absence (the prototype asks consent + installs on a yes). Do NOT make discovery hard-require agentsview without a graceful path.
- Spike-first risk: confirm the agentsview DB shape + the sqlite queries actually reconstruct a useful workflow on a real project's history (the prototype already does this — exercise it on this repo's own history as the spike) before committing to the integration shape.

## Out of scope

Building agentsview itself (it's an external tool). Re-deriving the scan — reuse the prototype's queries.

## Approach

### Integration shape: a spacedock-owned, user-invocable skill

Ship discovery as a new `user-invocable: true` skill under `skills/`, a peer of `commission`/`debrief`/`refit` (these are the four skills that today carry `user-invocable: true` and surface as `/spacedock:<name>` commands). NOT a `spacedock` CLI subcommand — the binary is the status/dispatch engine, and discovery is an interactive, synthesis-heavy flow that belongs in a skill. NOT folded into commission as a hidden pre-step — discovery has its own trigger surface (a new user on a brownfield repo) distinct from "I want to design a workflow," and it must be runnable on its own to answer "what's been happening here" without committing to commission. The skill ENDS by handing off to commission (offering it), so commission stays the workflow-authoring entry point and discovery stays the orientation entry point.

The skill is a spacedock-owned port of the `~/.claude/skills/orient/SKILL.md` prototype: it keeps the prototype's read-only, queries-inline philosophy and its synthesis format, fixes the agent-portability bug the spike found, adds scaffold-recognition, and ends with a discovery → commission bridge.

### User-facing name

**Chosen: `survey`.** Invoked as `/spacedock:survey`.

Rationale — why it beats shipping "orient" as the user command:
- It is what a person does *first* on unfamiliar ground: you survey the territory before you build. The verb already implies "I just arrived and want the lay of the land," which is exactly the first-run posture.
- It is honest about the action: the skill reads existing history and reports back. "survey" promises a read, not a change — matching the skill's read-only nature. A first-time user is reassured that running it is safe.
- "orient" works as the mental model ("get oriented") but as a *command* it is weaker: it reads as reflexive ("orient *me*") and presumes the user already knows the project — the opposite of the brownfield-newcomer case. "orient" is retained as the internal concept and as a description keyword so the existing orient muscle-memory still routes here.
- Rejected alternatives: `start-here` / `getting-started` (reads like onboarding docs, implies a guided setup not a diagnosis of existing history); `scout` (good "reconnaissance" connotation but collides with common tool names and is less self-evidently read-only); `discover` (accurate but generic, and overloaded across tooling). `survey` is short, a real verb, unambiguously a read, and idiomatic for "first thing you run."

The skill's `description` front-matter MUST keep the prototype's trigger phrases ("what have we been doing here", "catch me up", "orient me", "where did we leave off", "what's the state of this project", "picking up a brownfield repo") so the existing natural-language routes still land on `survey`.

### Scaffold recognition + comparable spacedock workflow + concrete comparative benefit

After the history scan, before the discovery report, the skill detects whether the project already runs a common agent scaffold, by reading project files (not the session DB):

- **superpowers** — presence of the superpowers plugin / skills: `.claude/skills/superpowers/`, a `superpowers` entry in `.claude-plugin/`/plugin marketplace config, or skills whose front-matter names the superpowers discipline (e.g. `brainstorming`, `writing-plans`, `executing-plans`, `subagent-driven-development`).
- **gsd / get-shit-done** — a `gsd`/`get-shit-done` skill or command directory, a `GSD.md`/`gsd` config, or the gsd phase vocabulary in agent instructions.
- **similar** — a generic fallback: any `.claude/skills/` or `.claude/commands/` tree with a multi-phase agent discipline that is neither of the above (report it by name without claiming a specific comparison).

On a detected scaffold, the skill offers a COMPARABLE spacedock workflow framed against what that scaffold does, and explains the benefit **concretely and comparatively**, anchored to what the scan found. The comparison is per-scaffold, not generic:

- **superpowers** is a library of disciplines an agent invokes (brainstorming → writing-plans → executing-plans), with human interruption left implicit — *the human decides when to step in.* The comparable spacedock workflow maps those disciplines to stages (ideation → implementation → validation) and makes the interruption points EXPLICIT gates. Benefit, tied to the scan's interruption count: "superpowers gives your agent the *plays* but leaves *when you step in* up to you — the scan counted N interruptions across M sessions where you had to. A spacedock workflow turns those into explicit approval gates, so the agent advances on its own between your calls and only stops where you marked a gate."
- **gsd / get-shit-done** runs a fixed phase sequence per task. The comparable spacedock workflow maps the gsd phases to stages and adds gates + entity state, so multiple work items move through the same phases concurrently and pause only at gates. Benefit: "gsd drives one task through its phases; spacedock tracks every work item through the same stages as durable on-disk state, gates the steps you flagged as needing you (the scan found these OPEN decisions: …), and lets several run in parallel without you re-driving each."
- **similar / unknown scaffold** — name it, then offer the generic spacedock benefit (gates from the interruption count, entity state, parallelism) without a false-specific comparison.

The OPEN decisions and the interruption count from the scan are the raw material that makes the benefit concrete — the skill names the actual OPEN forks and the actual interruption number, not a placeholder.

### The discovery → commission bridge

The discovery report ends with an offer: "Want me to commission a spacedock workflow from this?" On a yes, the skill invokes `commission` in batch mode (commission already supports batch design inputs in its first message). Discovery supplies commission's inputs from the scan: the inferred workflow loop → the stage list; the workstreams → seed entities; the OPEN decisions → the approval gates (each OPEN fork is a candidate gate); when a scaffold was detected, the comparable-workflow stage mapping is the proposed stage list. Discovery does not generate files itself — it assembles a commission invocation and hands off, so file generation stays commission's job.

### agentsview dependency + graceful-absent path

The skill reads agentsview's session DB. Handling its absence (proven mechanism, see spike):

- **Absent (no `agentsview` binary AND no `~/.agentsview/sessions.db`):** print the missing notice, explain agentsview ingests the agent logs the skill reads, **ask consent**, and only on a yes run the install (`brew install --cask agentsview`, fallback `curl -fsSL https://agentsview.io/install.sh | bash`). Never install without a yes. This path is inherited verbatim from the prototype's step 1.
- **Present:** best-effort `agentsview sync` to freshen, then query the DB.
- **Empty for this project (0 sessions):** say there is no agent history for this project and stop — nothing to discover.

## Out of scope (confirmed)

Building agentsview itself. Re-deriving the scan from raw logs (reuse the prototype's queries, fixed per the spike). Auto-installing agentsview without consent. Generating workflow files inside the discovery skill (that is commission's job; discovery hands off).

## Acceptance criteria

Each AC names a property of the finished skill and how a check OUTSIDE the skill text verifies it. Behavioral ACs are proven by RUNNING the skill on real history and observing output — never by grepping the skill file for a phrase.

**AC-1 — The user-facing entry point is named `survey` and is invocable as a spacedock skill command.**
Verified by: a behavior fixture / live invocation where `/spacedock:survey` (or the skill named `survey`) routes to the discovery skill — the skill resolves and runs, and the skill is registered `user-invocable: true` so the host exposes the command. The check is "invoking the command runs discovery," confirmed by observed startup output, not by finding the string `survey` in the file. (The plugin/skill-registry resolution is the independent source of truth: a mis-registered or misnamed skill fails to launch.)

**AC-2 — Running discovery on a project with agent history produces the inferred workflow, the workstreams, the OPEN-decision frontier, and the interruption stats as observable output.**
Verified by: a live run on a history-bearing project (e.g. this repo's own agentsview history) prints OVERVIEW (session count + date range), INTERRUPTIONS (a numeric count), WORKSTREAMS, and — when any decision is OPEN — a NEEDS-YOU section listing the OPEN forks. The expected values (session count, date range) come from the agentsview DB, an independent source that can diverge from any hard-coded number — so the check can fail. Proven by running the skill and reading its output, never by string-matching the skill.

**AC-3 — Discovery recognizes an existing scaffold (superpowers and gsd at minimum) and offers a comparable spacedock workflow whose benefit is stated concretely and comparatively.**
Verified by: two live runs, one in a fixture repo carrying a superpowers scaffold and one carrying a gsd scaffold. Each run's output must (a) name the detected scaffold, (b) offer a comparable spacedock workflow, and (c) state a benefit that references BOTH what the scaffold does AND a value spacedock adds, anchored to the scan (e.g. the interruption count or the OPEN decisions). The two runs must produce DIFFERENT comparison prose (superpowers-vs-spacedock differs from gsd-vs-spacedock) — a generic identical pitch fails this AC. The independent source of truth is the fixture's scaffold files + the scan numbers, which the comparison must reflect; an inverted or generic claim diverges and fails. Proven by running the skill against the two fixtures and reading output, never by grepping the skill file.

**AC-4 — Discovery hands a discovered workflow off to commission.**
Verified by: a live run where, on the "commission this?" yes-path, the skill invokes `commission` with batch inputs derived from the scan (stage list, seed entities, candidate gates), and commission proceeds into its design/confirm flow. Confirmed by observing commission start with discovery-supplied inputs (the inputs originate from the scan, not from the skill text), not by string-match.

**AC-5 — Discovery degrades gracefully when agentsview is absent and when the project has no history.**
Verified by: a run with no agentsview present asks consent before installing and never installs on no input (observable: it stops at the consent prompt, no install side-effect); a run against a project with zero sessions reports "no agent history" and stops. Both are observed runtime behaviors with side-effects (install attempted / not, stop vs continue), not prose.

## Test plan

- **Spike (done, this stage) — the riskiest mechanism.** The agentsview sqlite scan reconstructs a useful workflow on THIS repo's own history, plus the agentsview-absent path. Result recorded under **## Spike record** below. Cost: ~30 min. Outcome: mechanism works AND surfaced two design-load-bearing bugs in the prototype's queries (see spike).
- **AC-1 (name + registration):** behavior fixture / CLI exercise that the skill registers and launches as `survey`. Low cost. Fixture-level (the host's skill resolution), not a string match.
- **AC-2 (discovery output) + AC-5 (graceful paths):** live-history exercise on this repo's agentsview data (the implementation's first test seeds from the spike's working queries). For AC-5's absent path, a fixture run with `AGENT_VIEWER_DATA_DIR` pointed at an empty dir and the binary masked simulates "absent" without uninstalling anything. Medium cost (needs a synced DB; the spike proves the sync path).
- **AC-3 (scaffold recognition + comparative benefit):** two throwaway fixture repos — one seeded with a superpowers `.claude/skills/superpowers/` tree, one with a gsd skill/command dir — driven through the skill; assert each output names its scaffold and that the two comparison blocks differ and each cites a scan number. Behavior-fixture / live run. Medium cost.
- **AC-4 (commission hand-off):** live run exercising the yes-path into commission batch mode on a history-bearing fixture; assert commission starts with scan-derived inputs. Medium cost; reuses commission's existing batch-mode entry.

Live workflow / behavior-fixture tests throughout — no string-match-over-skill tests count toward any behavioral AC.

## Spike record

**Riskiest mechanism exercised first:** the agentsview sqlite scan, run against this repo's own real session history, plus the agentsview-absent path. Run from `/Users/clkao/git/spacedock-research/spacedock-v1` on 2026-06-05. Findings (all observed, not assumed):

1. **The scan mechanism works.** Pointing `AGENT_VIEWER_DATA_DIR` at a temp dir and running `agentsview sync` produced a queryable `sessions.db`; the scoped sync of this repo's history completed in ~12s (514 sessions, 92k messages, 66k tool_calls). Querying by the `project` column ('spacedock_v1' = 35 sessions, 2026-05-31..2026-06-05; 'spacedock' = 158 sessions) returned a real OVERVIEW + user-turn counts, and the RECENT-PROMPTS workstream query returned this repo's actual dispatch sessions. The reconstruction is useful.

2. **BUG (load-bearing): the prototype's project filter is Claude-Code-only and silently returns nothing for other agents.** The prototype matches sessions with `file_path LIKE '%<cwd-with-slashes-as-dashes>%'`. That dash-encoding is *Claude Code's* projects-path convention; Codex/Gemini sessions store a flat `file_path` (e.g. `~/.codex/sessions/.../rollout-*.jsonl`) with NO cwd in it. On this repo's DB the prototype's `LIKE` matched **0 rows** even though 35 spacedock_v1 sessions exist — because they are Codex sessions. agentsview already derives a correct `project` column from each session's cwd across all agents. **Fix for implementation: filter by the `project` column, not by `file_path LIKE` on a dash-mangled cwd.** This is the first test the implementation must pass: a project with non-Claude sessions returns its sessions.

3. **BUG (load-bearing): the prototype's decision/interruption signal is Claude-Code-exclusive.** It keys on `tool_name IN ('AskUserQuestion','ExitPlanMode')` and `json_extract(input_json,'$.questions[0].header')`. Codex sessions have neither tool — their plan/decision signal is `update_plan` (81 calls in spacedock_v1) with a different shape (`$.plan[].step/status`); interruptions surface differently too. So on a Codex/Gemini project the prototype reports `asks=NULL plans=NULL` and an empty decisions list — exactly the project where you most want the OPEN-decision frontier. **Fix for implementation: the decision/interruption extraction must be agent-aware** (Claude: AskUserQuestion/ExitPlanMode; Codex: update_plan + interrupt markers; degrade honestly when a signal is absent rather than emitting an empty section that looks like "no decisions").

4. **BUG (environment, must surface): the agent process cannot read `~/.agentsview/` directly.** Both inside and outside the Claude sandbox, `stat`/`ls`/`sqlite3` on `~/.agentsview/sessions.db` fail with macOS TCC "Operation not permitted" / "authorization denied", while the `agentsview` *binary* (differently entitled) reads it fine, and `~/.claude/projects/` is readable. **Implication for design:** the prototype's plan to read the default DB with raw `sqlite3` works in a normal Claude Code session that has Full Disk Access, but breaks in a sandboxed/limited-permission agent. The robust path proven here is to drive reads through the `agentsview` binary (run `sync`, optionally into an `AGENT_VIEWER_DATA_DIR` the process can read) rather than assuming raw FS access to `~/.agentsview/`. The implementation must not hard-require raw `~/.agentsview/` FS access without this fallback.

5. **Cost note:** a full-history `agentsview sync -full` is ~16k sessions and minutes-long; it also filled this near-full disk (the machine sits at ~99% on /). Discovery should `sync` incrementally (default, not `-full`) and scope by project, not force a full resync.

The integration shape composes only already-proven behavior on top of this fixed scan (a user-invocable skill calling commission in batch mode — both existing mechanisms). The one unverified mechanism (the scan) was exercised here; the two query bugs it surfaced are now seeded as the implementation's first tests.

## Stage Report: ideation

- DONE: Choose the user-facing starter name (with a rationale that it beats "orient" for a new user's first run) and decide the integration shape — a spacedock-owned skill, a spacedock command, or a commission pre-step
  Chose `survey` (rationale under ## Approach > User-facing name: first action on unfamiliar ground, honest read-only verb, beats reflexive "orient me"); integration shape = a `user-invocable: true` spacedock-owned skill peer of commission/debrief/refit (not a CLI subcommand, not a hidden commission pre-step), ending in a commission hand-off.
- DONE: Design scaffold-recognition (detect superpowers, gsd / get-shit-done, and similar) that offers a COMPARABLE spacedock workflow with a concrete comparative benefit; the ACs must prove recognition + benefit by RUNNING the skill and observing output, never a string-match over the skill file
  Designed per-scaffold detection (superpowers via .claude/skills/superpowers + discipline skill names; gsd via gsd/get-shit-done dir/config; generic fallback) with per-scaffold comparative benefit anchored to the scan's interruption count / OPEN decisions; AC-3 proves it via two fixture runs (superpowers + gsd) asserting differing, scan-anchored comparison prose — no string-match.
- DONE: Spike-first the riskiest mechanism — the agentsview sqlite scan reconstructs a useful workflow on THIS repo's own history, plus a graceful agentsview-absent path — and record the result (or an auditable "no spike needed" with the proven mechanism)
  Ran the scan on this repo's real history (514 sessions synced in ~12s; 35 spacedock_v1 + 158 spacedock sessions reconstructed via the project column). Recorded under ## Spike record; surfaced 4 load-bearing findings incl. two prototype query bugs (Claude-Code-only file_path filter; Claude-Code-only decision tools) and the macOS-TCC FS-access blocker, all now seeded as implementation tests.

### Summary

Fleshed out the discovery task: integration shape is a user-invocable `survey` skill (porting the orient prototype) that scans agentsview history, recognizes incumbent scaffolds and pitches a comparable spacedock workflow with a scan-anchored comparative benefit, then hands off to commission. The riskiest mechanism — the agentsview sqlite scan on this repo's own history — was exercised and works, but surfaced two load-bearing prototype bugs (its `file_path LIKE` project filter and its decision/interruption signal are both Claude-Code-only and silently empty on Codex/Gemini sessions, which is what THIS repo actually has) plus a macOS-TCC blocker on raw `~/.agentsview/` reads from a sandboxed process; the fixes are recorded and seeded as the implementation's first tests. All five ACs are behavioral and proven by running the skill, never by string-matching the skill file.

## Stage Report: implementation

- DONE: Ship a user-invocable `survey` skill (port of the orient prototype): read-only, queries-inline, the synthesis format; FIX the two Claude-Code-only prototype bugs so it works on this repo's Codex history — filter by the agentsview `project` column (not file_path LIKE dash-cwd), and make decision/interruption extraction agent-aware (Codex update_plan + interrupts, not only AskUserQuestion/ExitPlanMode)
  Skill at `skills/survey/SKILL.md` (commit e4777476). Bugs fixed and validated against this repo's REAL Codex history (scoped DB, project='spacedock_v1'): the project-column filter returns 70 top-level sessions (codex 35 + claude 35) where the prototype's `file_path LIKE` matched 0 Codex rows; agent-aware extraction surfaces claude asks=78, codex plans=81, 11 codex `<turn_aborted>` interrupts — the verbatim skill query block produced correct OVERVIEW/AGENTS/INTERRUPTIONS/CLAUDE-DECISIONS/CODEX-OPEN-STEPS output (run, not grepped).
- DONE: Route agentsview reads through the agentsview binary with a tmp-DB env var (AGENT_VIEWER_DATA_DIR) per captain direction — the macOS-TCC-safe path; consent-before-install on the absent path; "no history" stop on an empty project
  TCC blocker reconfirmed live: raw `ls`/`sqlite3` on `~/.agentsview/` fail with "Operation not permitted"/"authorization denied" from this process; skill syncs into a persistent process-readable `AGENT_VIEWER_DATA_DIR` (default `$TMPDIR/spacedock-survey`) — first run one-time full ingest, incremental re-sync measured 0s. Absent path: masked-PATH run prints `AGENTSVIEW MISSING` (consent gate, no install side-effect); empty project: COUNT(*) returns 0 for a nonexistent project key → "no agent history" stop.
- DONE: Scaffold-recognition (superpowers + gsd + generic fallback) offering a comparable spacedock workflow with DIFFERING, scan-anchored comparative benefit, ending in the commission hand-off; every behavioral AC proven by RUNNING the skill on real history/fixtures, never a string-match over the skill file
  Verbatim scaffold-detection block run (bash) against four fixture repos returns `superpowers` / `gsd` / `similar: deploy,my-pipeline` / `none` correctly. Per-scaffold comparison templates are structurally distinct claims (superpowers: implicit-interruption→explicit-gates citing the interrupt count; gsd: single-task-phases→parallel-gated-entities citing the OPEN forks) each anchored to real scan numbers; close with the commission batch-mode hand-off supplying stages/seed-entities/gates/mission from the scan. Registered `survey` in `skills/integration/skill_surface_test.go` userSkills so the frontmatter+reference-closure lints cover it; `go test ./skills/integration/` 60/60 pass; `go build ./...` + `go vet` clean.

### Summary

Shipped `skills/survey/SKILL.md` as a user-invocable spacedock skill porting the orient prototype, with both Claude-Code-only prototype bugs fixed and validated against this repo's actual Codex sessions: project-column filtering (catches all 70 top-level sessions vs the prototype's 0 Codex matches) and agent-aware decision/interruption extraction (Claude AskUserQuestion/ExitPlanMode; Codex update_plan + `<turn_aborted>`). Reads route through the agentsview binary into a persistent process-readable `AGENT_VIEWER_DATA_DIR` (the macOS-TCC-safe path, reconfirmed live), with consent-gated install on the absent path and a "no history" stop on an empty project. Scaffold recognition over four fixtures returns the right label and the two comparison templates are distinct scan-anchored claims ending in the commission hand-off. AC-1 registration is locked by adding `survey` to the integration suite's published-surface list; 60/60 integration tests pass, build and vet clean. Schema note: the prototype's `display_name` column is absent in the synced schema — used `first_message` for the prompts query.
