# First Officer Merge Core (host-neutral)

The terminal merge-and-cleanup ceremony, the mod-block guard that protects it, and step 10's boundary obligation. Lazily loaded at the terminal boundary (named by the boot-resident core); a boot, dispatch, or gate that never terminalizes never reads it. The runtime adapter supplies the host's concrete terminal teardown (step 10's host-specific part), read alongside this file.

## Merge and Cleanup

When an entity reaches its terminal stage, `«merge.guard»(slug)` drives the terminal merge-finalize ceremony as a re-entrant partial envelope whose role depends on the workflow's `merge:` policy. Under `merge: pr` (the default) the verb is a THIN finalize-helper: the FO arms the mod-block, opens the captain-gated PR, and detects the merge; the verb only reads the `pr`/`mod-block`/`verdict` state delta and signals `blocked` (a PR is pending — leave state intact) or `finalized` (the PR landed — clear the mod-block standalone, terminalize, archive). Under `merge: local` the verb additionally owns the arm step: with no in-flight mod-block it sets `mod-block=merge:{hook}` and signals the FO to invoke the hook, then finalizes on re-run. In neither policy does the verb itself invoke the hook or local-merge, and the clear+terminalize is always two separate `--set` calls, not one.

## «merge.guard»(slug): atomically set→invoke→clear the merge mod-block, then terminalize

- **effect:** drive the terminal merge-finalize ceremony, policy-conditional. **arm (merge: local only):** with an empty mod-block, set `mod-block=merge:{mod_name}` and signal the FO to invoke the hook (the verb does NOT invoke it). **finalize (on re-run):** read the `pr`/`mod-block`/`verdict` delta — if a `pr` is set and the verdict is not `rejected`, signal `blocked` and leave state intact; otherwise clear the mod-block in its own standalone `--set`, terminalize (`completed verdict={verdict} worktree=`), and archive. The verb never invokes the hook and never local-merges; the mechanism-level mod-block guard (below) refuses any out-of-order terminal transition without `--force`.
- **done-when:** the entity is archived terminal with its mod-block cleared and `pr`/sentinel recorded (or the hook left it blocked, mod-block still set and the pending state reported).
- **block:** if the hook blocks (`pr` set, captain approval pending, external wait), leave `mod-block` set and do not local-merge. `--force` is never part of the happy path — if the guard refuses, a step was skipped, not a flag forgotten.
- → **shipped** (this sprint): `` `spacedock merge guard <slug>` `` — invoke it directly. The hand-followed sequence it automates is the steps below.

The verb performs the arm/clear/terminalize/archive sequence. The FO owns only what the verb does NOT:

1. **Invoke the merge hook.** When the verb signals `armed`, run the merge hooks (`merge: local` does the local `--no-ff` merge; `merge: pr` opens the captain-gated PR and sets `pr`), then re-run `merge guard`. The verb signals; the FO invokes.
2. **Default local merge** when no merge hook is registered: merge `{branch}` onto the trunk (`spacedock dispatch trunk --workflow-dir {workflow_dir}`, default `main`) from the stage worktree branch. The verb never local-merges under any policy.
3. **Remove the worktree** (`git worktree remove {path}`, no `--force`) and delete the local branch (`git branch -d {branch}`). Do NOT delete the remote branch while a PR is pending — the reviewer needs it; remote cleanup belongs to the PR merge.
4. **Teardown workers at terminal (step 10).** At the terminal boundary, derive the entity's worker cohort and cooperatively shut each one down (best-effort, fire-and-forget — the runtime adapter supplies the shutdown call), then drop them from session memory. Mandatory whether the merge ran locally or via a PR host. The cohort-derivation rule and the cooperative-shutdown call are the adapter's. A runtime MAY add a further teardown step (e.g. a bounded team-registry teardown); the adapter declares it where it applies. This core states only the boundary obligation: cohort shutdown, then drop from session memory.

### Worktree removal safety

Use `git worktree remove {path}` (no `--force`). The default refuses to delete a worktree with untracked changes — that refusal is the safety net.

If removal fails on untracked files, the FO MUST:

1. Audit: `git -C {path} status --short` from the parent worktree.
2. Decide per file: commit to the worktree branch (audit-essential per gitignore), move to a persistent location (experiment-output outside the worktree), or explicitly confirm destruction with the captain.
3. ONLY after the audit, `--force` is permitted.

`--force` is never default; it is an explicit captain-confirmed bypass.

## Mod-Block Guard

The mod-block guard exists so the FO recovers on session resume and so the merge ceremony cannot skip its hook. `merge guard` owns the set→clear mechanics; the FO's standing concern is recovery:

- **Survives session resume.** Read `mod-block` from frontmatter on boot. A non-empty `mod-block=merge:{mod_name}` means a merge hook is mid-flight — do not re-run it from scratch; check what it left (PR created? branch pushed?) and re-run `merge guard` to continue.
- **The guard refuses, it does not auto-fix.** `status --set`/`status --archive` refuse a terminal transition while `mod-block` is non-empty, refuse combining `mod-block=` with terminal fields, and refuse terminalizing with merge hooks registered while `pr` and `mod-block` are both empty (the hook provably has not run). `--force` bypasses; `merge: local` and `verdict=rejected` exempt only the pr-requirement. Do NOT `--force` merely to clear the guard — it catches exactly the mistake of skipping the hook.
- **A missing blocking mod is a captain escalation.** If `{workflow_dir}/_mods/{mod_name}.md` is missing or unreadable, report: "Blocking mod {mod_name} is missing. The entity is stuck. Options: restore the mod file, or use `--force` to clear the block and resume normal flow." Wait for direction.
