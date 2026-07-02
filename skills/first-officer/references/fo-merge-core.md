# First Officer Merge Core

The terminal merge-and-cleanup ceremony, the mod-block guard that protects it, and step 10's boundary obligation. The runtime adapter supplies the host's concrete terminal teardown (step 10's host-specific part), read alongside this file.

## Merge and Cleanup

`«merge.guard»` never invokes the merge hook and never local-merges — `merge: local`'s hook itself performs the `--no-ff merge`; `merge: pr` opens the captain-gated PR. The merge sentinel is the FINALIZE key: the FO's `pr-merge` startup/idle/sweep hook detects MERGED via `gh` and records the sentinel; `«merge.guard»` keys off it without ever talking to GitHub. Each phase's next action for the FO — which hook to invoke, when to wait, what to clean up — is named in `merge guard`'s own output; see `«merge.guard»` below.

**Launcher invariant — use the right binary.** `merge guard` exists only on the current-checkout / `SPACEDOCK_BIN` binary, not necessarily the brew-installed `spacedock` on `$PATH`. Invoke it via `${SPACEDOCK_BIN:-spacedock} merge guard <slug>` (or the checkout binary directly). If you call a stale `spacedock` that predates `«merge.guard»`, the subcommand is unknown and you silently fall back to the hand ceremony — the toil `«merge.guard»` eliminates.

## «merge.guard»(slug): auto-arm → block-on-open-PR → finalize-on-merge-sentinel, then archive

- **effect:** drive the terminal merge-finalize ceremony — auto-arm, block on an open PR, finalize on a merge sentinel, then archive (including the path-scoped archive commit) — the same under both `merge:` policies. Invoke it once per phase; its own stdout/stderr name the FO's next action (armed: which hook to invoke and where; blocked: wait for the sentinel, never finalize on `pr`-presence alone; finalized: worktree/branch/worker cleanup, or the manual merge when no hook is registered).
- **done-when:** the entity is archived terminal, or `«merge.guard»` left it armed/blocked with its next step named in its own output.
- **block:** `--force` is never part of the happy path — if the guard refuses, a step was skipped, not a flag forgotten.
- → **shipped**: `` `spacedock merge guard <slug>` `` — invoke it directly per phase (via `${SPACEDOCK_BIN:-spacedock}`, per the launcher invariant above).

**Step 10: teardown workers at terminal.** At the terminal boundary, derive the entity's worker cohort and cooperatively shut each one down (best-effort, fire-and-forget — the runtime adapter supplies the shutdown call), then drop them from session memory. Mandatory whether the merge ran locally or via a PR host — `«merge.guard»`'s own finalized output only points here ("tear down the entity's workers per your runtime adapter"); the cohort-derivation rule and the cooperative-shutdown call are the adapter's. A runtime MAY add a further teardown step (e.g. a bounded team-registry teardown); the adapter declares it where it applies.

### Worktree removal safety

`--force` is never default; audit untracked files with the operator first.

## Mod-Block Guard

A non-empty `mod-block=merge:{mod_name}` on boot means a merge is mid-flight — check what the hook left and re-run `merge guard` to continue.
