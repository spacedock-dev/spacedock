# Paseo Ensign Runtime

How the shared ensign core executes when dispatched through Paseo. The shared core owns the ensign discipline (assignment reading, worktree, split-root commit, proof, stage report); this adapter binds only the Paseo-specific concerns.

## Runtime implementation

- `Fresh dispatch` -> Load `ensign-shared-core.md` and `paseo-ensign-runtime.md`, Read the `/tmp/spacedock-dispatch/...md` pointer file, run `### Fetch commands` verbatim, and follow the assignment together with any ordinary FO scope notes appended to the prompt.
- `Git sync in headless/sandboxed workers` -> Commit code to the assigned worktree branch and entity/report updates path-scoped in the split-root state checkout. When pushing over SSH without an unlocked agent, use `GIT_SSH_COMMAND='ssh -o BatchMode=yes'` so git fails fast rather than hanging on a passphrase prompt; if push fails due to sandbox or missing credentials, keep the local path-scoped commit clean and note the push failure in `### Summary`.
- `Clarification` -> Ordinary Paseo ensigns run without a daemon send grant; if genuinely blocked, commit safe partial progress (or leave worktree clean), record the blocker in `## Stage Report` (`FAILED:` + `### Summary`), and end your turn so the FO observes `idle` and intervenes.
- `Completion signal` -> After all commits and stage-report writes succeed, end your final turn with one single-line message: `Done: {entity title} completed {stage}. Report written to {entity_file_path}.` Do not call host-specific mailbox tools; Git commits and the stage report are the completion gate.
