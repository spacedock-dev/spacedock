# First Officer Merge Module (Claude)

The terminal merge-and-cleanup ceremony, the mod-block enforcement that guards it, and the bounded terminal teardown. Lazily loaded at the terminal boundary — the boot-resident core names this file at the merge load point and reads it only when an entity reaches its terminal stage. A boot, dispatch, or gate that never terminalizes never reads it.

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
10. **Teardown agents at terminal.** Derive the entity's agent cohort from the live team roster — every worker whose handle decomposes to this entity's slug (roster and decomposition are the adapter's). Issue the cooperative-shutdown call (best-effort, fire-and-forget) and drop them from session memory. Then tear down the team itself as a **bounded best-effort**: the cooperative shutdown and the team-teardown call race — the first teardown attempt can fail because a member the FO just signalled is still settling out of the roster ("active member(s)"). Do NOT end the turn on that first failure. Between attempts the FO MUST let the roster settle — re-issue the cooperative shutdown to any still-named active member, then **wait a short settle interval before the next teardown attempt** rather than re-firing in the same instant (an instant retry re-loses the same async registry race — a teardown that "retried but raced every time, then stopped" still hangs). Attempt the settle-then-teardown serially until it succeeds or a small **attempt cap** is reached. In an interactive session the roster clears as the member's session-end propagates, so teardown succeeds on an early attempt and the loop exits naturally. In a non-interactive session (single-entity `-p` mode), an approved-shutdown member can stay listed in the roster indefinitely (an upstream defect), so teardown can never succeed — "retry to success" is unreachable and a fast retry loop only re-hangs the subprocess. On **cap-exhaustion the FO STOPS teardown attempts and emits a defined terminal-status marker — `TERMINAL_TEARDOWN_BOUNDED: best-effort teardown exhausted; member(s) stuck in registry; holding for launcher.` (verbatim).** The PROCESS EXIT is the **launcher's** responsibility: the FO cannot self-exit while the roster is non-empty, so a non-interactive launcher (the live-e2e cycle's `kill()`, or a real automation's timeout) ends the subprocess once the marker has been emitted. The FO emitting the marker is the bounded-teardown terminus a watcher grades; a teardown that gives up silently with no marker, or one that retries past the cap and never reaches the marker, is the failure this step prevents. On a subsequent harness re-invocation with the roster still non-empty, the FO again runs the bounded best-effort and re-emits the marker — acceptable (the launcher ends the subprocess). What this step forbids is an UNBOUNDED retry loop that never reaches the marker. **Mandatory at the boundary; the settle interval, the cap value, and the marker emission are the adapter's.**

### Ship-Local Ceremony

When the merge boundary has no PR host (README declares `merge: local`, or pr-merge fallback applies — no `gh`, push failed, captain chose local), the FO runs one fixed ceremony per entity. The README's top-level `merge:` key (default `pr`) selects this ceremony or the PR path. The happy path uses no `--force`:

1. Set the merge mod-block: `spacedock status --workflow-dir {workflow_dir} --set {slug} mod-block=merge:{mod_name}` (commit path-scoped).
2. Invoke the merge hook (local `--no-ff` merge of `{branch}` onto `next`).
3. Record the merge so the terminal guard is satisfied without `--force`:
   - If `merge: local`, the policy exempts the pr-requirement — skip to step 4.
   - Otherwise set the post-merge sentinel `spacedock status --workflow-dir {workflow_dir} --set {slug} pr=local-merge:{short-sha}` (the merge commit on `next`; set only after merge has landed; commit path-scoped). The status table renders as `{short-sha} (local)`.
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

## Mod-Block Enforcement at Terminal Transitions

Before advancing an entity into Merge and Cleanup, the FO MUST:

1. Check whether merge hooks are registered (from boot-time MODS data).
2. If merge hooks exist, set `mod-block` before invoking the first hook.
3. Invoke merge hooks in order. If a hook blocks (sets `pr`, requires captain approval), leave `mod-block` set and report the pending state.
4. Clear `mod-block` only after the blocking condition is resolved (PR merged, captain chose alternative, hook completed without blocking).
5. Proceed to terminal frontmatter updates (completed, verdict, worktree clear) and archival only after `mod-block` is clear.

**The mechanism enforces this even if you forget.** `status --set` and `status --archive` refuse terminal transitions (status to a terminal stage, completed, verdict, worktree clear) and archival when all of the following hold true:

- the workflow registers at least one merge hook (`_mods/*.md` with `## Hook: merge`),
- the entity's `pr` field is empty,
- the entity's `mod-block` field is empty,
- `--force` was not passed.

In that state the merge hook has provably not run. The refusal names the blocking hook so you can recover by: setting `mod-block=merge:{mod_name}` and invoking the hook (normal flow), letting the hook set `pr` (which satisfies the invariant), or passing `--force` (captain explicitly approved bypassing the hook). Do NOT pass `--force` merely to clear the guard — it exists to catch exactly the mistake of skipping the hook.

On session resume, scan entities with non-empty `mod-block` and resume the pending action. Do not re-run the hook from scratch — check what it left (PR created? branch pushed?) and continue from there.

If the blocking mod file (`{workflow_dir}/_mods/{mod_name}.md`) is missing or unreadable, report to the captain: "Blocking mod {mod_name} is missing. The entity is stuck. Options: restore the mod file, or use `--force` to clear the block and resume normal flow." Wait for direction.
