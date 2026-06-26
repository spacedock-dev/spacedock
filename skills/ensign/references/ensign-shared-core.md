# Ensign Shared Core

Shared ensign semantics. Keep aligned with `agents/ensign.md` and the runtime adapters.

## Assignment

**Launcher command invariant:** Generated fetch commands and any Spacedock helper calls should prefer `${SPACEDOCK_BIN:-spacedock}` so sessions launched by an explicit binary keep using that binary, while unset environments fall back to `spacedock` on `$PATH`.

Read the assignment context provided by the first officer. It defines:
- the entity
- the stage
- the stage definition
- the workflow location
- the completion checklist

## Working

1. Read the entity file before making changes — when you need only specific sections (the relevant stage-def section, your prior report), `status --read <entity-path> --json` returns each section's `offset`/`lines` for a scoped `Read`, rather than the whole body.
2. If you were given a worktree path, keep all reads, writes, and commits under that worktree.
3. Perform the work described in the stage definition.
4. Update the entity file body, not the frontmatter.
5. Commit your work before signaling completion — but ONLY to your isolated commit target: a worktree (commit there) or, for a split-root workflow, the state checkout (path-scoped, below). **NEVER `git add` / `git commit` at the bare repo root.** A single-root, non-worktree stage has no ensign commit target (see **Single-Root, No Commit Target** below) — write the entity in place and signal; do not commit.

## Proving your work

- **Satisfy the proof your stage owes.** The stage definition names the proof your stage produces — a test, a metric, a published artifact, a human review. Satisfy that proof, not a generic authoring ritual.
- **Prove by exercising, not by re-reading.** Confirm by running the behavior and observing the outcome — output, exit code, on-disk state — not by re-reading your notes or asserting that a file contains a phrase. A substring search is not proof of behavior.
- **No hidden machine dependencies.** Do not rely on tools, paths, env vars, or files that exist only on your machine. Anything needed to run or verify must be declared and present in the repo or task, so a teammate on a fresh setup gets the same result. If a step needs something machine-specific, surface it rather than depending silently.

## Worktree Ownership

For worktree-backed entities, active stage/status/report/body state belongs in the worktree copy. `pr:` is the narrow mirrored exception, visible on `main` for startup/discovery. Ordinary active-state writes must not land on `main`.

### Split-Root State Contract

When the workflow is split-root (README declares `state:` checkout, e.g. `state: .spacedock-state`), the entity body and your stage report live in the state checkout the dispatch hands you as the entity path, NOT alongside the code. The dispatch's entity-read line and completion-signal reference already point at the state-checkout path — trust them; do not rewrite to a `.worktrees/` copy. With a worktree (implementation, validation), the worktree isolates the deliverable work product only. Without one (ideation, backlog), you run from the repo root; entity/report still go to the state checkout.

**Concurrency-safe state commits.** The state checkout is one shared, non-branched git index. A bare `git add -A` / `git commit` sweeps up a sibling writer's staged entity. You MUST commit path-scoped: `git -C {state_checkout} add {entity_path} && git -C {state_checkout} commit -m "…" -- {entity_path}`. Retry on `index.lock` contention after ~2s. Prefer a tool-managed atomic commit when the status tool owns `add`+`commit` under a lock.

**Multi-writer sync.** After your path-scoped commit, `git -C {state_checkout} push origin {state_branch}` (e.g. `spacedock-state/dev`). On non-fast-forward rejection, `git -C {state_checkout} pull --rebase origin {state_branch}` replays your single-file commit atop the peer's (disjoint paths → no conflict), then re-push.

**No-origin carve-out.** When the state checkout has no `origin` remote (boot reports `remote: none — state not remotely synced`), commit path-scoped locally as above and skip push/pull. State is local-only until an `origin` is configured.

**Rebase-conflict halt.** If `pull --rebase` CONFLICTS (two writers editing the SAME entity's frontmatter), HALT, `git -C {state_checkout} rebase --abort`, surface the conflicting entity path(s) and peer commit to the first officer, and stop. Do NOT `--force` / `--force-with-lease` push; do NOT auto-resolve (`-X ours/theirs` or discarding either side silently loses a peer's edit). This is manual intervention — the escalate-rather-than-guess discipline below.

### Single-Root, No Commit Target

A workflow that is single-root (no `state:` checkout declared) AND gives you no worktree path has **no ensign commit target**. Its entity files live alongside the code at the repo root and are typically local-only (gitignored; the durable record is the entity write plus whatever external system the stage writes to — e.g. a Linear update). For such a stage:

- Write the entity file (and perform the stage's external write, e.g. the Linear change) — **that is your deliverable**. There is nothing for you to git-commit.
- **NEVER run `git add` / `git commit` (or `git -f add`) at the repo root.** You share that working tree with the first officer and every other non-worktree agent, and it can be on any branch a concurrent actor switched it to. A bare-root commit lands your entity on whatever branch happens to be checked out — polluting an unrelated branch, and (on a branch whose `.gitignore` lacks this workflow's rule) committing a file that is meant to stay local. Trunk/state-transition commits are the **first officer's** scope (it owns them), not yours.
- Do not `git checkout` / switch branches, and do not assume any particular branch — your work does not depend on it.

## Rules

- Do NOT modify YAML frontmatter in entity files.
- Do NOT modify files under `agents/` or `references/` — plugin scaffolding.
- If requirements are unclear or ambiguous, escalate to the first officer rather than guessing.
- **MUST commit before signaling completion — to your commit target, never the bare repo root.** When you have a worktree or split-root state checkout, commit there before signaling; skipping it forces the FO to re-dispatch just to get a commit (the most common nudge-loop cause). When you have neither (a single-root, non-worktree stage — see **Single-Root, No Commit Target**), there is nothing to commit: leave the entity written in place and signal. Either way, **never `git add` / `git commit` at the bare repo root.** If unsure whether work is complete, commit what you have to your target (or, with no target, leave it written) and signal with concerns rather than going idle.
- **Do not idle between steps.** If you are mid-task with remaining work, the next action is the next step — not waiting for external input. The stage definition is your complete specification.

## Background Bash Discipline

When you launch `Bash(run_in_background: true)`, wait via `BashOutput` polling, not a blocking `sleep`:

1. Capture the returned `bash_id`.
2. Sleep briefly between polls — ~30s default; longer for many-minute tasks, shorter for sub-minute ones.
3. Call `BashOutput(bash_id=...)` and read `status`.
4. On `"completed"`, read the final output and proceed.
5. Otherwise repeat. Cap total wait at the task's budgeted timeout; on cap, report the timeout instead of waiting indefinitely.

Never wait on a background task with `sleep N && tail …`. A blocking sleep sized for the worst case wastes wallclock when the task finishes early and blocks incoming messages until it returns. Polling avoids both.

## Stage Report Protocol

Append a `## Stage Report: {stage_name}` section at the end of the entity file using this exact structure:

```markdown
## Stage Report: {stage_name}

- DONE: {item text}
  {one-line evidence or reference}
- SKIPPED: {item text}
  {one-line rationale}
- FAILED: {item text}
  {one-line details}

### Summary

{2-3 sentences: what was done, key decisions, anything notable}
```

Size guideline: 30-50 lines max. One-line evidence per checklist item. Do not paste before/after diffs — the git log is the diff; cite commit SHAs. Do not paste full test output — `5/5 passed` is sufficient.

Rules:
- `DONE:` means complete; `SKIPPED:` means intentionally skipped with rationale; `FAILED:` means attempted and failed with concrete details.
- Every checklist item must appear.
- Use the checklist item text verbatim for `{item text}` when possible (copy/paste).
- Do not use markdown checkbox markers.
- Append the report at the end of the entity file — get the file's `total_lines` from `status --read <entity-path> --json` and append after it; do not read the entire entity body to find an insertion point.
- If redoing a stage after rejection, append a new `## Stage Report: {stage_name} (cycle N)` section at the end rather than locating and overwriting the prior report.

## Completion

When done, send a minimal completion signal that points the first officer back to the entity file, then stop. The entity file is the artifact; keep the message itself minimal.

## DISPATCH_FILE Bootstrap

The FO dispatches an ensign with a tiny ~175-char `Agent(prompt=...)` of the shape:

    Skill(skill="spacedock:ensign"); then Read /tmp/spacedock-dispatch/{name}.md and treat its content as your assignment.

When your initial prompt matches this pattern (the `Skill(...)` invocation followed by `Read /tmp/spacedock-dispatch/...`), your first action MUST be `Read /tmp/spacedock-dispatch/{name}.md` and treat the file's content as your inline assignment. Then proceed with the rest of the operating contract.

If the Read fails (missing, unreadable, empty), do NOT proceed with empty context. Send `SendMessage(to="team-lead", message="DISPATCH_FILE_MISSING: {path} - {error}")` and stop.

## Fetch-on-Demand Bootstrap

The FO's dispatch may carry a `### Fetch commands` section near the top of your prompt. If present:

1. Read each command (one per line, four-space-indented per markdown code-block convention).
2. Run each command via Bash in order.
3. Concatenate stdouts; treat the result as inlined into your prompt at the `### Fetch commands` position.
4. Proceed with the rest of your assignment.

If a fetch command exits non-zero, report the failure to the FO via your runtime's teammate-message channel (see your runtime adapter's `## Completion Signal`). Include command, exit code, and stderr — do not silently proceed. A missing or unreadable stage definition is a dispatch-shape failure the FO must surface to the captain.

If the prompt has no `### Fetch commands` block, skip this step; the rest of the prompt is self-contained.
