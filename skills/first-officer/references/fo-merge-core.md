# First Officer Merge Core

The terminal merge-and-cleanup ceremony, the mod-block guard that protects it, and the `«worker.shutdown»()` boundary obligation. The runtime adapter supplies the concrete terminal teardown.

## Merge and Cleanup

`«merge.guard»` never invokes `«hooks.run»("merge")` and never local-merges. When armed, the FO invokes `«hooks.run»("merge")`: the `merge: local` registration performs `--no-ff merge` — on conflict, surface it and stop, never auto-resolve; `merge: pr` opens the captain-gated PR. The registered `«hooks.run»("startup")` / `«hooks.run»("idle")` paths detect MERGED via `gh` and record the FINALIZE sentinel; `«merge.guard»` keys off it without talking to GitHub.

**Launcher invariant — use the right binary.** `merge guard` exists only on the current-checkout / `SPACEDOCK_BIN` binary, not necessarily the brew-installed `spacedock` on `$PATH`. Invoke it via `${SPACEDOCK_BIN:-spacedock} merge guard <slug>` (or the checkout binary directly). If you call a stale `spacedock` that predates `«merge.guard»`, the subcommand is unknown and you silently fall back to the hand ceremony — the toil `«merge.guard»` eliminates.

## «merge.guard»(slug): auto-arm → block-on-open-PR → finalize-on-merge-sentinel, then archive

- **effect:** drive the terminal merge-finalize ceremony (including the path-scoped archive commit and split-root publication), the same under both `merge:` policies. Invoke it once per phase. Armed names `«hooks.run»("merge")` as the next action; blocked waits for the sentinel; finalized names worktree/branch/worker cleanup or the manual merge when no registration exists.
- **done-when:** the entity is archived terminal and remote-durable when state has an origin (explicitly local-only otherwise), or `«merge.guard»` left it armed/blocked with its next step named in its own output.
- **block:** exit 3 is `«halt.rebase-conflict»(paths)` — HALT and surface path/peer evidence, never rerun guard. Exit 1 after the archive commit resumes through `«state.commit»(slug)`; only an ordinary pre-mutation guard refusal means a step was skipped. `--force` is never part of the happy path.
- → **shipped**: `` `spacedock merge guard <slug>` `` — invoke it directly per phase (via `${SPACEDOCK_BIN:-spacedock}`, per the launcher invariant above).

At the terminal boundary, invoke `«worker.shutdown»()` for the entity's worker cohort: derive the cohort, cooperatively shut each member down (best-effort, fire-and-forget through the runtime binding), then drop them from session memory. This is mandatory whether the merge ran locally or via a PR host. A runtime may add further teardown, such as bounded team-registry cleanup, in its binding.

### Worktree removal safety

`--force` is never default; audit untracked files with the operator first.

## Mod-Block Guard

A non-empty `mod-block=merge:{mod_name}` on boot means a merge is mid-flight. Check what `«hooks.run»("merge")` left, then re-run `«merge.guard»`; invoke `«hooks.run»("merge")` only if the guard returns armed again.
