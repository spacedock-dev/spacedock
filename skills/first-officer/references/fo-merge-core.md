# First Officer Merge Core (host-neutral)

The terminal merge-and-cleanup ceremony, the mod-block enforcement that guards it, and step 10's boundary obligation. Lazily loaded at the terminal boundary (named by the boot-resident core); a boot, dispatch, or gate that never terminalizes never reads it. The runtime adapter supplies the host's concrete terminal teardown (step 10's host-specific part), read alongside this file.

## Merge and Cleanup

When an entity reaches its terminal stage:

1. If merge hooks are registered, set the mod-block before invoking:
   `spacedock status --workflow-dir {workflow_dir} --set {slug} mod-block=merge:{mod_name}`
   Commit: `mod-block: {slug} awaiting merge:{mod_name}`.
   The mechanism enforces this — `status --set` and `status --archive` refuse terminal updates while merge hooks exist with both `pr` and `mod-block` empty, unless `--force`, `merge: local`, or `verdict=rejected` exempts (a rejected entity never ran the merge ceremony). Setting `mod-block` also lets session resume identify which mod is blocking.
2. Run merge hooks before local merge, archival, or status advancement.
3. Detect hook completion via the state delta. A hook blocks if (a) `pr` is now set, (b) its prose requires captain approval and the captain has not responded, or (c) it declares an external wait. Otherwise it completed.
4. If blocked, leave `mod-block` set, report the pending state, and do not local-merge.
5. If completed without blocking, clear the mod-block in its own `--set` call:
   `spacedock status --workflow-dir {workflow_dir} --set {slug} mod-block=`
   Commit: `mod-block: {slug} cleared ({mod_name} completed)`.
   The clear MUST be standalone — `status --set` exits 1 if `mod-block=` is combined with `status={terminal}`, `completed`, `verdict`, or `worktree=` in one call. Use two commits, or `--force` with captain approval only.
6. If no merge hook handled the merge, perform the default local merge from the stage worktree branch.
7. Update frontmatter: `spacedock status --workflow-dir {workflow_dir} --set {slug} completed verdict={verdict} worktree=`.
8. Archive: `spacedock status --workflow-dir {workflow_dir} --archive {slug}`.
9. Remove the worktree (`git worktree remove {path}`) and delete the local branch (`git branch -d {branch}`). Do NOT delete the remote branch while a PR is pending — the reviewer needs it. Remote cleanup belongs to the PR merge.
10. **Teardown agents at terminal.** At the terminal boundary, derive the entity's worker cohort and run the host's terminal-teardown ceremony — cooperative shutdown of the cohort (best-effort, fire-and-forget), drop them from session memory, then the host's bounded team/worker teardown. Teardown is mandatory at the terminal boundary whether the merge ran locally or via a PR host. The cohort-derivation rule, the team/worker-teardown call, the settle interval, the attempt cap, and any terminal-status marker are the host's — the runtime adapter supplies them. This core names none of them; it states only the boundary obligation.

### Ship-Local Ceremony

When the merge boundary has no PR host (README declares `merge: local`, or pr-merge fallback applies — no `gh`, push failed, captain chose local), the FO runs one fixed ceremony per entity. The README's top-level `merge:` key (default `pr`) selects this ceremony or the PR path. The happy path uses no `--force`:

1. Set the merge mod-block: `spacedock status --workflow-dir {workflow_dir} --set {slug} mod-block=merge:{mod_name}` (commit path-scoped).
2. Resolve the integration trunk `BASE=$(spacedock dispatch trunk --workflow-dir {workflow_dir})` (configured trunk, default `main`), then invoke the merge hook (local `--no-ff` merge of `{branch}` onto `{BASE}`).
3. Record the merge so the terminal guard is satisfied without `--force`:
   - If `merge: local`, the policy exempts the pr-requirement — skip to step 4.
   - Otherwise set the post-merge sentinel `spacedock status --workflow-dir {workflow_dir} --set {slug} pr=local-merge:{short-sha}` (the merge commit on `{BASE}`; set only after merge has landed; commit path-scoped). The status table renders as `{short-sha} (local)`.
4. Clear the mod-block in a standalone `--set`: `spacedock status --workflow-dir {workflow_dir} --set {slug} mod-block=` (commit path-scoped). MUST be separate from terminalization — the guard refuses combining `mod-block=` with terminal fields.
5. Terminalize: `spacedock status --workflow-dir {workflow_dir} --set {slug} completed verdict={verdict} worktree=`.
6. Archive: `spacedock status --workflow-dir {workflow_dir} --archive {slug}`.
7. Remove worktree, delete local branch (Merge-and-Cleanup step 9), and run the terminal agent teardown (step 10). Teardown is mandatory at the terminal boundary whether the merge ran locally or via a PR host.

The set→invoke→clear sequence (steps 1, 2, 4) is mandatory whenever a merge hook is registered, regardless of `merge: local`. `--force` is never part of the happy path — if the guard refuses, a step was skipped, not a flag forgotten.

### Worktree removal safety

Use `git worktree remove {path}` (no `--force`). The default refuses to delete a worktree with untracked changes — that refusal is the safety net.

If removal fails on untracked files, the FO MUST:

1. Audit: `git -C {path} status --short` from the parent worktree.
2. Decide per file: commit to the worktree branch (audit-essential per gitignore), move to a persistent location (experiment-output outside the worktree), or explicitly confirm destruction with the captain.
3. ONLY after the audit, `--force` is permitted.

`--force` is never default; it is an explicit captain-confirmed bypass.

## Mod-Block Enforcement

Merge hooks can block (captain approval before pushing, waiting for PR merge). The FO enforces through the entity `mod-block` field and a mechanism-level invariant in `status --set` / `status --archive`:

- **Set** by the FO before invoking a merge hook: `mod-block=merge:{mod_name}`.
- **Cleared** after the blocking action completes or the captain force-overrides. The clear runs in its own `--set` — combining `mod-block=` with terminal fields (`status={terminal}`, `completed`, `verdict`, `worktree=`) is refused without `--force`.
- **Guarded** — `status --set` refuses terminal transitions while `mod-block` is non-empty unless `--force` is passed.
- **Enforced at the mechanism level** — `status --set` and `status --archive` also refuse terminal transitions and archival when merge hooks (`_mods/*.md` with `## Hook: merge`) are registered AND `pr` is empty AND `mod-block` is empty. `--force` bypasses. `merge: local` exempts only the pr-requirement; `verdict=rejected` likewise exempts only the pr-requirement on both surfaces (a rejected entity never ran the merge ceremony, so the requirement is vacuous); the mod-block-pending and combined-clear refusals remain. See the Ship-Local Ceremony.
- **Survives session resume** — the FO reads `mod-block` from frontmatter on boot and resumes the pending action.

In the empty-pr/empty-mod-block state the merge hook has provably not run. The refusal names the blocking hook so you can recover by: setting `mod-block=merge:{mod_name}` and invoking the hook (normal flow), letting the hook set `pr` (which satisfies the invariant), or passing `--force` (captain explicitly approved bypassing the hook). Do NOT pass `--force` merely to clear the guard — it exists to catch exactly the mistake of skipping the hook.

On session resume, scan entities with non-empty `mod-block` and resume the pending action. Do not re-run the hook from scratch — check what it left (PR created? branch pushed?) and continue from there.

If the blocking mod file (`{workflow_dir}/_mods/{mod_name}.md`) is missing or unreadable, report to the captain: "Blocking mod {mod_name} is missing. The entity is stuck. Options: restore the mod file, or use `--force` to clear the block and resume normal flow." Wait for direction.
