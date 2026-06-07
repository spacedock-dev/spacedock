---
name: survey
description: Use when arriving at or returning to a project that already has AI-agent history and you want the lay of the land before doing anything else — "survey this project", "what have we been doing here", "catch me up", "orient me", "where did we leave off", "what's the state of this project", or picking up a brownfield repo with several in-flight agent tracks. Reads existing agent session history (read-only), reports the implicit workflow, the open decisions, and how often you had to step in, then offers to commission a spacedock workflow from it.
user-invocable: true
---

# Survey a Project

## Overview

Survey is the first thing you run on unfamiliar ground: it reconstructs what the AI agents in this project have implicitly been doing, from their session history. It reports the inferred workflow, the workstreams, the recent decisions, and — load-bearing — the OPEN decisions (the abandoned or unanswered forks) plus how often the human had to step in. Then it offers to commission a real spacedock workflow with explicit gates from what it found.

It reads **agentsview**'s session DB and is strictly read-only — the recommended queries live in `references/queries.sql` (one labeled query per concern) so nothing is a black box. For now it surveys **Claude Code** history (the decision and interruption signals below are Claude's); agentsview also ingests Codex, Gemini, and more, and surfacing those agents' decision signals is a deferred follow-up. The closing move is the discovery → commission bridge: the OPEN decisions become candidate gates, the workstreams become candidate entities, the inferred loop becomes the stage list.

Run the four steps in order: **check agentsview → scan → recognize scaffold → report and offer**.

---

## 1. Check agentsview, then sync THIS project only (scoped)

This skill may run in a sandboxed agent that **cannot read `~/.agentsview/` directly** (macOS TCC denies raw FS access to a limited-permission process, even though the `agentsview` binary itself reads it). So do NOT `sqlite3 ~/.agentsview/sessions.db` blindly. Instead, drive the read through the `agentsview` binary into a process-readable data directory under `AGENTSVIEW_DATA_DIR`, then query that copy:

```bash
SURVEY_DB_DIR="${SPACEDOCK_SURVEY_DB_DIR:-${TMPDIR:-/tmp}/spacedock-survey}"
DB="$SURVEY_DB_DIR/sessions.db"

if ! command -v agentsview >/dev/null; then echo "AGENTSVIEW MISSING"; fi
```

If it prints `AGENTSVIEW MISSING`: tell the user agentsview is needed (it ingests the agent logs this skill reads), **ask consent**, and only on a yes run the install (`brew install --cask agentsview`; fallback `curl -fsSL https://agentsview.io/install.sh | bash`). NEVER install without an explicit yes — stop at the consent prompt otherwise.

With the binary present, sync — but **scope the sync to this project**. A bare `agentsview sync` enumerates the ENTIRE `~/.claude/projects` history (16k+ sessions, growing): on a real machine that walk exhausts any sane timeout, so the survey data dir ends up empty or partial. The fix is to narrow Claude's source root to just this repo's session directories before syncing, so the walk is bounded and this project's Claude sessions land in seconds:

```bash
mkdir -p "$SURVEY_DB_DIR"

# Claude Code stores each project's sessions in ~/.claude/projects/<cwd-with-/-as->,
# so this repo's sessions live under dirs that begin with the dash-encoded cwd. Point
# CLAUDE_PROJECTS_DIR at a symlink farm of just those dirs — the sync then walks only
# this project's Claude sessions (this cwd plus its worktrees), not the whole backlog.
CLAUDE_ROOT="${CLAUDE_PROJECTS_DIR:-$HOME/.claude/projects}"
DASH_CWD=$(pwd | sed 's#/#-#g')                      # ~/.claude/projects dir-name convention
NARROW="$SURVEY_DB_DIR/claude-narrow"
rm -rf "$NARROW"; mkdir -p "$NARROW"
shopt -s nullglob 2>/dev/null
for d in "$CLAUDE_ROOT/$DASH_CWD" "$CLAUDE_ROOT/$DASH_CWD"-*; do
  [ -d "$d" ] && ln -s "$d" "$NARROW/$(basename "$d")"
done

AGENTSVIEW_DATA_DIR="$SURVEY_DB_DIR" CLAUDE_PROJECTS_DIR="$NARROW" timeout 300 agentsview sync
```

The survey data dir persists between runs, so a re-survey of the same project is incremental (seconds). Do not pass `--full` — a full resync re-ingests everything and can fill the disk. If the symlink farm is empty (this project has no Claude sessions under `~/.claude/projects`), the synced DB has no history for it; step 2 reports "no agent history" and stops.

If `agentsview sync` fails (network, disk, permissions), report the exact failure and stop — do not fall back to raw `~/.agentsview/` reads (they fail under TCC).

## 2. Scan the project

**Scope by repo IDENTITY, not by `basename(pwd)`.** agentsview keys each session by the basename of its working directory, so a subdir checkout, a worktree, or the split-root state dir each get a DIFFERENT `project` key even though they are ONE repo — run from any of those and a `basename(pwd)` filter finds zero sessions while the real history sits under a sibling key. Resolve the repo root once, then scope every query to the cwds UNDER that root so the divergent keys coalesce:

```bash
# Repo-root identity: the parent of the common .git dir resolves to the SAME absolute
# path from the repo root, a subdir, the state checkout, or a linked worktree — so it
# coalesces every checkout of one repo. --path-format=absolute needs git >= 2.31; if
# this is not a git repo or git is older, fall back to the cwd itself (today's behavior).
if REPO_GIT=$(git rev-parse --path-format=absolute --git-common-dir 2>/dev/null); then
  REPO_ROOT=$(dirname "$REPO_GIT")
else
  REPO_ROOT=$(pwd)
fi
```

Bind `REPO_ROOT` as the `:repo_root` parameter of the queries in `references/queries.sql`. The session-scoping subquery there matches `cwd = :repo_root OR cwd LIKE :repo_root || '/%'` — the cwd AT the root and everything strictly under it, excluding a sibling like `…/proj-other`. Run a labeled query by name (the SQL stays in the reference file; this only extracts and runs it):

```bash
SURVEY_SKILL_DIR="${SPACEDOCK_SURVEY_SKILL_DIR:-skills/survey}"
QUERIES="$SURVEY_SKILL_DIR/references/queries.sql"

run_query() {  # run_query <name> — runs the labeled query, :repo_root bound to REPO_ROOT
  local q
  q=$(awk -v n="$1" '$0=="-- name: "n{f=1;next} /^-- name: /{f=0} f && $0 !~ /^--/{print}' "$QUERIES")
  printf ".param set :repo_root '%s'\n%s\n" "$REPO_ROOT" "$q" | sqlite3 "$DB"
}

run_query scoping        # #318 — sessions|folded_keys|blank_cwd|span over the coalesced repo
run_query scaffold-usage # #319 — behavioral skill_name family tally (spacedock self EXCLUDED)
run_query work-by-area   # #317.2 — Edit/Write file_path bucketed by package (external = reference)
run_query decision-open  # #320 — AskUserQuestion/ExitPlanMode frontier; OPEN sorts first
```

`scoping` returns `sessions=0` → there is no Claude agent history for this repo; say so and stop. Nothing to discover. (Survey reads Claude history only for now; a repo whose only agent history is Codex/Gemini will report "no agent history" here — surfacing those agents is a deferred follow-up.) When `folded_keys > 1`, the repo had its history split across several agentsview keys (subdir / worktree / state checkout) and the scoping coalesced them — note in the report how many keys were folded in, and the `blank_cwd` count if non-zero (sessions agentsview never captured a cwd for, which the repo-root scope cannot place).

**Honest signal accounting.** The `decision-open` rows are the human-decision points; `OPEN` = still needs the human, and you lead the report with those. For the interruption total, count the AskUserQuestion / ExitPlanMode decisions plus the hard-veto markers Claude sessions retain (`[Request interrupted` / `Request interrupted by user` / `doesn't want to proceed` in the message stream), over the same repo-scoped session set; `pct = total*100/user_turns`. Never dress an empty section up as "no decisions" — if a section is empty, say the run found none of that signal.

## 3. Recognize an incumbent scaffold

Recognize the scaffold from TWO signals and reconcile them — a file probe (what is installed on disk) and the behavioral tally (what actually ran), which is the `scaffold-usage` query you already ran in step 2. A file-only probe misses a scaffold that was invoked but isn't checked in; a tally-only read misses one installed but never used. Join them.

**File probe — multi-label, not a single winner.** Probe each scaffold INDEPENDENTLY and name EVERY match; report `none` only when no probe matched (the old single-winner if-ladder hid co-installed scaffolds):

- **superpowers** — `.claude/skills/superpowers` exists, `superpowers` appears in `.claude-plugin/`, or a superpowers discipline skill dir is present (`.claude/skills/{brainstorming,writing-plans,executing-plans,subagent-driven-development,…}`);
- **gsd / get-shit-done** — a `.claude/skills/gsd` or `.claude/skills/get-shit-done` dir, a `.claude/commands/gsd` dir, or a `GSD.md` / `gsd.md` / `.gsd` file;
- **similar / unknown** — any other `.claude/skills` or `.claude/commands` tree (name the dirs you found);
- **none** — none of the above is present on disk.

**Behavioral tally.** The `scaffold-usage` rows are a `family → invocations` tally normalized from `tool_calls.skill_name` (`superpowers:brainstorming` and the bare `running-research-spikes` both fold to family `superpowers`); the `spacedock` family is excluded because survey/ensign self-invocation otherwise dominates and would make every repo read as "uses spacedock".

**Join into a 3-bucket classification.** For each family that appears in EITHER signal:

- **file-present + invoked** → **active** — report it plainly;
- **file-present + never-invoked** → **installed-but-unused** — flag it (installed, not yet adopted);
- **not-file-present + invoked** → **recovered** — report it; this is the case the file-only probe misses entirely, recovered by the behavioral signal.

The classification drives the comparative benefit in the report (step 4). The probe reads files; the comparison's *numbers* come from the scan (step 2).

## 4. Confirm, then report and offer

Every `{slot}` below is a FILL slot: substitute the real value from the step-2 scan before you show the user. A literal `{slot}` (or a `<…>` angle token) left in what you present is a bug — never show the user an unfilled slot. If a slot has no data (e.g. zero OPEN decisions), drop that line rather than printing an empty slot.

**Cross-check the OPEN frontier against the repo (before you present it).** The `decision-open` query is a TRANSCRIPT-only scan — a fork that read OPEN there may already be shipped (a merged PR / a commit) and over-reports. For each transcript-OPEN fork, cross-reference the repo (`git log`, merged PRs via `gh pr list --state merged` if available, the working tree) and split it:

- **shipped** → **DROP** from the frontier. Evidence: a merged PR or a git-log commit whose subject/body CONFIDENTLY references the fork (its decision header or branch — an exact-ish token match).
- **decided-not-shipped** → move to a **backlog** line (decided, no artifact yet).
- **never-decided** → **true open**, stays on the `NEEDS YOU` frontier.

**Conservative-match rule.** DROP only on a CONFIDENT repo match; anything less than confident → KEEP on the frontier. A false "still open" is a cheap nudge; a false "shipped" silently hides a real open fork — so the asymmetry favors keeping.

**Mandatory degrade.** When NO repo signal is available (not a git repo, or `git log` / PR lookup fails or is empty), the frontier degrades to transcript-only and EVERY OPEN fork is flagged **`unverified`** in the report — never silently presented as authoritative. The degrade is the default behavior, not an error.

Tell the user what you found and wait for a yes:

> Found **{N} sessions** in `{project}` (`{date range}`), with **{D} decision points** and **{V} interruptions**. Want me to lay it out?

Then synthesize this, one screen:

```
PROJECT: {basename}     {sessions} Claude sessions · {date range}
  {if folded_keys>1: coalesced from {folded_keys} agentsview keys}{if blank_cwd>0: · {blank_cwd} uncaptured-cwd sessions}

SCAFFOLD
  {the multi-label classification: each family + its bucket — active / installed-but-unused / recovered; or "none"}

INFERRED WORKFLOW
  {the implicit loop across the decisions + prompts, as an arrow chain} — {one honest line}

WORKSTREAMS
  {cluster the decisions + prompts into tracks; one line each, status glyph}

WORK BY AREA   (what this is — where edits land)
  {the work-by-area buckets: package — {edits}; an <external> bucket is edits to a sibling repo (a reference, not this project)}

NEEDS YOU   (only if any decision is still OPEN after the repo cross-check)
  ⚠ {the true-open forks — never-decided questions; lead with them}{if degraded: each flagged unverified (no repo signal)}

BACKLOG   (only if any fork was decided-not-shipped)
  {decided forks with no shipped artifact yet}

RECENT DECISIONS  (answered or shipped)
  {the rest: header — short question}

INTERRUPTIONS  (where spacedock can help)
  {total} times you stepped in across {sessions} sessions — {decisions} decision points
  + {vetoes} course-corrections, {pct}% of your turns.
```

### The discovery → commission bridge (close every report with this)

After the synthesis, recognize the scaffold (step 3) and offer a COMPARABLE spacedock workflow, with a benefit stated **concretely and comparatively**, anchored to the actual scan numbers — never a placeholder, never a generic pitch. As in the synthesis above, every `{slot}` is a FILL slot: substitute the real step-2 number/forks before you show the user; a literal `{slot}` in your output is a bug. Use the per-scaffold framing:

- **superpowers** is a library of disciplines an agent invokes (brainstorming → writing-plans → executing-plans → subagent-driven-development), with human interruption left implicit — *the human decides when to step in.* Offer a spacedock workflow that maps those disciplines to stages (ideation → implementation → validation) and makes the interruption points EXPLICIT gates. State it tied to the scan's interruption count:
  > superpowers gives your agent the *plays* but leaves *when you step in* up to you — this scan counted **{V} interruptions across {N} sessions** where you had to. A spacedock workflow turns those into explicit approval gates, so the agent advances on its own between your calls and only stops where you marked a gate.

- **gsd / get-shit-done** runs a fixed phase sequence per task, one task at a time. Offer a spacedock workflow that maps the gsd phases to stages and adds gates + durable entity state, so multiple work items move through the same phases concurrently and pause only at gates. State it tied to the OPEN forks:
  > gsd drives one task through its phases; spacedock tracks every work item through the same stages as durable on-disk state, gates the steps you flagged as needing you (this scan found these OPEN forks: **{the actual OPEN decisions}**), and lets several run in parallel without you re-driving each.

- **similar / unknown scaffold** — name it (use the names the step-3 detection emitted), then offer the generic spacedock benefit (gates from the interruption count, entity state, parallelism) without inventing a false-specific comparison.

- **none** — offer the generic spacedock benefit anchored to the interruption count and OPEN forks.

The two comparisons MUST differ — superpowers-vs-spacedock (implicit-interruption → explicit-gates) is a different claim from gsd-vs-spacedock (single-task-phases → parallel-gated-entities). Each must cite a real scan number (the filled `{V}`/`{N}` or the filled OPEN forks), not a placeholder.

Then make the offer:

> Want me to commission a spacedock workflow from this?

On a **yes**, invoke commission in batch mode, supplying inputs derived from the scan (commission already accepts batch design inputs in its first message — see its Batch Mode). Assemble:

- **stages** ← the inferred workflow loop (for a detected scaffold, the comparable-workflow stage mapping above is the proposed stage list);
- **seed entities** ← the workstreams;
- **approval gates** ← the true-open forks (each OPEN decision that survived the repo cross-check is a candidate gate);
- **mission / entity** ← inferred from the workstreams and the project.

Hand off by invoking `commission` with those assembled inputs. Survey does NOT generate workflow files itself — file generation stays commission's job; survey only assembles the invocation and hands off.

On a **no**, stop — the survey stands on its own as an orientation.

## Synthesis guidance

- **Project name** = path basename.
- **Workflow + workstreams: infer them**, primarily from the decisions (the `PROMPTS` are sparse/noisy — secondary). Be honest when a track is one-off or stalled.
- **Decisions + stats are data, not invention.** `OPEN` = still needs the human; lead the report with the true-open forks. The transcript scan can't tell shipped from open — that's what the step-4 repo cross-check is for; drop a fork to "shipped" only on a confident match, and flag the whole frontier `unverified` when there's no repo signal.
- **Work-by-area is identity, not a to-do list.** The `work-by-area` buckets say WHAT this project is (where edits land); an `<external>` bucket is edits to a sibling repo — a reference, not this project's identity. Report it separately from the decision frontier (where you stop).
- **Fill every slot, never invent.** Every `{slot}` in the report and the comparison comes from the step-2 numbers; a literal `{slot}` shown to the user is a bug. If a section's signal is empty (no OPEN decisions, no interruptions, no edits), say the run found none — never dress an empty section up as "no decisions."
- **Claude-only for now.** Survey reads Claude history; Codex/Gemini decision signals are a deferred follow-up. Don't imply a Codex/Gemini-only project has "no history" in the user-facing wording beyond what step 2 reports.
