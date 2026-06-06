---
name: survey
description: Use when arriving at or returning to a project that already has AI-agent history and you want the lay of the land before doing anything else — "survey this project", "what have we been doing here", "catch me up", "orient me", "where did we leave off", "what's the state of this project", or picking up a brownfield repo with several in-flight agent tracks. Reads existing agent session history (read-only), reports the implicit workflow, the open decisions, and how often you had to step in, then offers to commission a spacedock workflow from it.
user-invocable: true
---

# Survey a Project

## Overview

Survey is the first thing you run on unfamiliar ground: it reconstructs what the AI agents in this project have implicitly been doing, from their session history. It reports the inferred workflow, the workstreams, the recent decisions, and — load-bearing — the OPEN decisions (the abandoned or unanswered forks) plus how often the human had to step in. Then it offers to commission a real spacedock workflow with explicit gates from what it found.

It reads **agentsview**'s session DB and is strictly read-only — every query is shown inline so nothing is a black box. For now it surveys **Claude Code** history (the decision and interruption signals below are Claude's); agentsview also ingests Codex, Gemini, and more, and surfacing those agents' decision signals is a deferred follow-up. The closing move is the discovery → commission bridge: the OPEN decisions become candidate gates, the workstreams become candidate entities, the inferred loop becomes the stage list.

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

agentsview derives a `project` column for every session from that session's working directory — its **basename, with non-alphanumerics replaced by `_`** (so this repo's cwd `…/spacedock-v1` keys as `spacedock_v1`). Filter by that `project` column. Do NOT filter by `file_path LIKE` a dash-mangled cwd: the project key is the stable, agentsview-computed key, and `file_path` matching is brittle across agentsview versions and source layouts.

The runnable scan surface is `bin/scan-project`, resolved relative to this skill directory (in this repo: `skills/survey/bin/scan-project`). It contains the explicit sqlite queries and their comments; keep it paired with this intent list as the source of truth:

- `OVERVIEW` counts this project's top-level Claude sessions by agentsview `project`;
- `INTERRUPTIONS` counts AskUserQuestion/ExitPlanMode stops, hard-veto markers, and user turns;
- `DECISIONS` lists AskUserQuestion/ExitPlanMode decisions, marks only answered-confirmation results as `done`, marks every rejection/error/prompt-echo as `OPEN`, and sorts OPEN before done so the frontier cannot be truncated by the recency `LIMIT`;
- `RECENT PROMPTS` provides secondary workstream signal.

Run the artifact, then read the labelled output:

```bash
SURVEY_SKILL_DIR="${SPACEDOCK_SURVEY_SKILL_DIR:-skills/survey}"
"$SURVEY_SKILL_DIR/bin/scan-project"
```

`OVERVIEW` is empty / `0 sessions` → there is no Claude agent history for this project; say so and stop. Nothing to discover. (Survey reads Claude history only for now; a project whose only agent history is Codex/Gemini will report "no agent history" here — surfacing those agents is a deferred follow-up.)

**Honest signal accounting.** The DECISIONS section lists the human-decision points; `OPEN` = still needs the human, and you lead the report with those. The interruption total is `asks + plans + vetoes` (the AskUserQuestion / ExitPlanMode decision tools plus the hard-veto markers Claude sessions retain); `pct = total*100/user_turns`. Never dress an empty section up as "no decisions" — if a section is empty, say the run found none of that signal.

## 3. Recognize an incumbent scaffold

Before the report, detect whether the project already runs a common agent scaffold — by reading project FILES (not the session DB). The runnable classifier is `bin/detect-scaffold`, resolved relative to this skill directory (in this repo: `skills/survey/bin/detect-scaffold`). It checks:

- superpowers via `.claude/skills/superpowers`, marketplace/plugin config, or superpowers discipline skill names;
- gsd/get-shit-done via a gsd skill/command dir or gsd config;
- similar scaffolds via any other `.claude/skills` or `.claude/commands` tree;
- none when no scaffold is present.

```bash
SURVEY_SKILL_DIR="${SPACEDOCK_SURVEY_SKILL_DIR:-skills/survey}"
"$SURVEY_SKILL_DIR/bin/detect-scaffold"
```

The detected scaffold name drives the comparative benefit in the report (step 4). The detection reads files; the comparison's *numbers* come from the scan (step 2).

## 4. Confirm, then report and offer

Every `{slot}` below is a FILL slot: substitute the real value from the step-2 scan before you show the user. A literal `{slot}` (or a `<…>` angle token) left in what you present is a bug — never show the user an unfilled slot. If a slot has no data (e.g. zero OPEN decisions), drop that line rather than printing an empty slot.

Tell the user what you found and wait for a yes:

> Found **{N} sessions** in `{project}` (`{date range}`), with **{D} decision points** and **{V} interruptions**. Want me to lay it out?

Then synthesize this, one screen:

```
PROJECT: {basename}     {sessions} Claude sessions · {date range}

INFERRED WORKFLOW
  {the implicit loop across the decisions + prompts, as an arrow chain} — {one honest line}

WORKSTREAMS
  {cluster the decisions + prompts into tracks; one line each, status glyph}

NEEDS YOU   (only if any decision is OPEN)
  ⚠ {the OPEN forks — abandoned/unanswered decision questions; lead with them}

RECENT DECISIONS  (answered)
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
- **approval gates** ← the OPEN forks (each OPEN decision is a candidate gate);
- **mission / entity** ← inferred from the workstreams and the project.

Hand off by invoking `commission` with those assembled inputs. Survey does NOT generate workflow files itself — file generation stays commission's job; survey only assembles the invocation and hands off.

On a **no**, stop — the survey stands on its own as an orientation.

## Synthesis guidance

- **Project name** = path basename.
- **Workflow + workstreams: infer them**, primarily from the decisions (the `PROMPTS` are sparse/noisy — secondary). Be honest when a track is one-off or stalled.
- **Decisions + stats are data, not invention.** `OPEN` = still needs the human; lead the report with those. Don't claim a decision was implemented — the scan doesn't know that.
- **Fill every slot, never invent.** Every `{slot}` in the report and the comparison comes from the step-2 numbers; a literal `{slot}` shown to the user is a bug. If a section's signal is empty (no OPEN decisions, no interruptions), say the run found none — never dress an empty section up as "no decisions."
- **Claude-only for now.** Survey reads Claude history; Codex/Gemini decision signals are a deferred follow-up. Don't imply a Codex/Gemini-only project has "no history" in the user-facing wording beyond what step 2 reports.
