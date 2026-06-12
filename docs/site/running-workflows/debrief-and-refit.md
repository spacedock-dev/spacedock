# Debrief & refit

Two maintenance commands keep a workflow durable across sessions and releases. `/spacedock:debrief` captures what happened in a session into a record the next session reads to start with context. `/spacedock:refit` brings an existing workflow's scaffolding up to the current Spacedock version while leaving your local edits in place. You run both as the captain; each pauses for your confirmation before it writes anything.

## Debrief: capture a session

`/spacedock:debrief` writes a structured record of a session (shipped entities, newly filed backlog seeds, workflow-only commits, gate decisions, issues, and what's next) to `{dir}/_debriefs/{date}-{sequence}.md` and commits it. The next session's first officer reads the most recent debrief instead of starting cold.

Run it at the end of a working session, or whenever you want a checkpoint:

```bash
spacedock claude "/spacedock:debrief"
```

The skill works in four phases. You make the decisions at the boundaries; everything else is git and local-file reads, with no external services until you ask it to file an issue.

1. **Discovery.** It finds the workflow, anchors the session start where the previous debrief ended (or at the workflow's recent history when there is none), shows you the boundary (since-commit and commit count), and waits for your confirmation or a corrected starting commit.

2. **Extract.** It buckets every commit in range: shipped PRs roll up into a **Shipped** section as links, routine state churn is suppressed, and only workflow-only commits that never flowed through a PR are listed. It reads entity frontmatter for what reached `done`, scans for gate approvals and rejections, and fills **What's next** from the dispatchable queue.

3. **Draft and review.** It presents the draft with **Decisions** and **Observations** left as placeholders for you to fill. Add why a gate was approved or rejected, scope changes, design insights, or confirm as-is. Issues are split into **Workflow** (quirks in your pipeline, kept local) and **Spacedock** (framework bugs). For each Spacedock issue it offers to file an **anonymized** GitHub issue: the body carries the bug, repro steps, and scale, but never your mission, entity titles, or domain. You approve, edit, or decline each one before any `gh issue create` runs.

4. **Write and commit.** It writes the debrief to `{dir}/_debriefs/{date}-{sequence:02d}.md`, commits it, and reports the path:

   ```
   Debrief written to {dir}/_debriefs/2026-06-09-01.md and committed.
   ```

The debrief's frontmatter records where the session ended, which is how the next debrief knows where to start.

## Refit: upgrade scaffolding to the current release

`/spacedock:refit` upgrades a workflow's scaffolding files (the README and any installed mods in `_mods/`) to match the current Spacedock version, and migrates entity frontmatter when a schema change requires it. Agent files and the status viewer ship with the plugin, so they are never refit locally. The skill never auto-replaces a file you may have customized; it shows you a diff and you decide.

You must give it the workflow directory:

```bash
spacedock claude "/spacedock:refit path/to/workflow"
```

It reads the version stamp from the README frontmatter (`commissioned-by: spacedock@X.Y.Z`) and each mod's `version` field, then compares them against the current version from the plugin manifest. If everything matches it stops with "Workflow is already up to date." Otherwise it presents an upgrade plan and proceeds per file by strategy:

- **`README.md`: show diff, never auto-replace.** Because you customize stages, schema fields, and quality criteria here, the skill generates what the current template would produce, diffs it against your README, and leaves it to you to apply the changes you want. It modifies only the version stamp at the end, not the README body.
- **`_mods/{name}.md`: version diff.** For each installed mod it compares your `version` against the canonical mod at `mods/{name}.md`. Matching versions are skipped; differing versions get a diff and a y/n. A mod with no canonical match is treated as custom: acknowledged, no action. Canonical mods you don't have installed are offered for install.
- **`status` (legacy): remove.** A workflow-local `status` script predates the launcher. The status viewer is now the `spacedock status` command, so refit removes the local copy with `git rm`.

### Schema migration and ID style

After scaffolding, refit compares the old and new README `## Schema` and `### Field Reference` sections for changed types or ranges, renamed fields, removed fields, or new required fields. If a change affects entity data, it lists the affected entities, proposes the migration (for example, "Convert score from /25 to 0.0–1.0 by dividing by 25"), and waits for your y/n. On approval it edits **only** the named frontmatter fields with the Edit tool, never an entity body, never a whole-file rewrite.

Refit preserves the README's `id-style` (`sequential`, `sd-b32`, or `slug`) and never changes it silently. It recommends `sd-b32` only under collaboration pressure (worktree stages, PR/merge mods, multiple creators, branches, offline work) and requires your explicit approval. Before any approved style change it runs `spacedock status --validate` against the workflow and reports failures; the actual ID rewrite is manual in this release.

### When there is no version stamp

If the README has no `commissioned-by` stamp, refit cannot tell what the original scaffolding looked like, so it enters **degraded mode** and offers two choices: **stamp only** (add stamps without changing anything, to establish a baseline) or **full refit with review** (show a full diff for every file and require your approval before replacing each). It never auto-replaces an unstamped file.

When the refit finishes it updates the README stamp to the current version, prints a per-file summary, and suggests the commit:

```bash
git commit -m "refit: upgrade workflow scaffolding to spacedock@{current_version}"
```

Git is the safety net throughout: `git diff` and `git checkout` recover anything you didn't mean to keep.

## Where these fit

Debrief and refit bracket the working loop described in [Operating a workflow](operating.md): you commission once, operate session by session, debrief at the end of a session, and refit when you upgrade Spacedock. For the commands these skills call, see the [Command reference](../reference/command-reference.md).
