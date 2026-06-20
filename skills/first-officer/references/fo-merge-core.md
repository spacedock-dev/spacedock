# First Officer Merge Core (host-neutral)

The terminal merge-and-cleanup ceremony, the mod-block guard that protects it, and step 10's boundary obligation. Lazily loaded at the terminal boundary (named by the boot-resident core); a boot, dispatch, or gate that never terminalizes never reads it. The runtime adapter supplies the host's concrete terminal teardown (step 10's host-specific part), read alongside this file.

## Merge and Cleanup

When an entity reaches its terminal stage, `«merge.guard»(slug)` drives the terminal merge-finalize ceremony as a re-entrant envelope: the FO invokes `spacedock merge guard` once per phase, and `«merge.guard»` owns arm + clear + terminalize + archive (including the path-scoped archive commit). The FO owns only the steps `«merge.guard»` cannot: invoking the hook, pushing, opening the captain-gated PR, and removing the worktree+branch. The role is the same under both `merge:` policies — `«merge.guard»` reads the `pr`/`mod-block`/`verdict` state delta and signals the phase:

- **armed** — entering terminal with an empty mod-block and a merge hook registered: `«merge.guard»` sets `mod-block=merge:{hook}` and signals the FO to invoke the hook. This is the same under `merge: local` (the FO's hook does the local `--no-ff` merge and records a `local-merge:{sha}` sentinel) and `merge: pr` (the FO's hook opens the captain-gated PR and sets `pr`).
- **blocked** — a bare/open PR reference is set (`pr: #42`) and the verdict is not `rejected`: the PR has not merged. `«merge.guard»` leaves the mod-block + pr intact and waits. It NEVER finalizes on `pr`-presence alone — archiving an entity before its PR landed is the premature-finalize bug.
- **finalized** — the verdict is `rejected`, OR the `pr` field carries a MERGE SENTINEL (`pr-merge:{number}` for a landed PR, `local-merge:{sha}` for a local merge) — the local signal that the merge LANDED. `«merge.guard»` clears the mod-block in its own standalone `--set`, terminalizes (`completed verdict={verdict} worktree=`), archives, and commits the archive move path-scoped. This finalizes EVEN from a non-armed (empty mod-block) state — the stranded case a re-validation bounce leaves behind.

`«merge.guard»` never invokes the hook and never local-merges, and the clear+terminalize is always two separate `--set` calls, not one. The merge sentinel is the FINALIZE key: the FO's `pr-merge` startup/idle/sweep hook detects MERGED via `gh` and records the sentinel; `«merge.guard»` keys off it without ever talking to GitHub.

**Launcher invariant — use the right binary.** `merge guard` exists only on the current-checkout / `SPACEDOCK_BIN` binary, not necessarily the brew-installed `spacedock` on `$PATH`. Invoke it via `${SPACEDOCK_BIN:-spacedock} merge guard <slug>` (or the checkout binary directly). If you call a stale `spacedock` that predates `«merge.guard»`, the subcommand is unknown and you silently fall back to the hand ceremony — the toil `«merge.guard»` eliminates.

## «merge.guard»(slug): auto-arm → block-on-open-PR → finalize-on-merge-sentinel, then archive

- **effect:** drive the terminal merge-finalize ceremony. **auto-arm (both policies):** with an empty mod-block and a merge hook registered, set `mod-block=merge:{mod_name}` and signal the FO to invoke the hook (`«merge.guard»` does NOT invoke it). **finalize:** read the `pr`/`mod-block`/`verdict` delta — if the verdict is `rejected` OR `pr` carries a merge sentinel (`pr-merge:`/`local-merge:`), clear the mod-block in its own standalone `--set`, terminalize (`completed verdict={verdict} worktree=`), archive, and commit the archive move path-scoped; finalize works EVEN from a non-armed state. **block:** if `pr` is a bare/open reference and the verdict is not `rejected`, signal `blocked` and leave state intact — never finalize on `pr`-presence alone. `«merge.guard»` never invokes the hook and never local-merges; the mechanism-level mod-block / merge-hook guards (below) refuse any out-of-order or hook-skipping terminal transition without `--force`, and `«merge.guard»` propagates — never bypasses — that refusal.
- **done-when:** the entity is archived terminal with its mod-block cleared and `pr`/sentinel recorded, the archive move committed path-scoped (or `«merge.guard»` left it blocked on an open PR, mod-block intact and the pending state reported).
- **block:** if the hook blocks (open PR pending, captain approval pending, external wait), leave `mod-block` set and do not local-merge. `--force` is never part of the happy path — if the guard refuses, a step was skipped, not a flag forgotten.
- → **shipped**: `` `spacedock merge guard <slug>` `` — invoke it directly per phase (via `${SPACEDOCK_BIN:-spacedock}`, per the launcher invariant above). The hand-followed sequence it automates is the steps below.

`«merge.guard»` performs the arm/clear/terminalize/archive sequence and commits the archive move path-scoped. The FO owns only what `«merge.guard»` does NOT:

1. **Invoke the merge hook.** When `«merge.guard»` signals `armed`, run the merge hooks (`merge: local` does the local `--no-ff` merge and records a `local-merge:{sha}` sentinel; `merge: pr` opens the captain-gated PR and sets `pr`), then re-run `merge guard`. `«merge.guard»` signals; the FO invokes. When the PR is detected MERGED (the `pr-merge` startup/idle/sweep hook's `gh` check), record the `pr-merge:{number}` sentinel in `pr` and re-run `merge guard` — the sentinel is what unlocks finalize.
2. **Default local merge** when no merge hook is registered: merge `{branch}` onto the trunk (`spacedock dispatch trunk --workflow-dir {workflow_dir}`, default `main`) from the stage worktree branch. `«merge.guard»` never local-merges under any policy.
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
