# Prototype: prose-functions for the state-repo operating principle

A quick prototype of the restructuring idea: anything that can later become a binary module is
**first** declared as a *prose-function* — `«fn»(arg)` — invoked in the main flow and defined in a
`## «fn»(arg): natural language` block. The flow reads as layer-3 intention; the layer-2 mechanics
live in the function body; each body names the binary verb it will become.

**Notation carries the migration status:**

- `«fn»(arg)` (guillemets) — a **prose-function**: the FO follows its body by hand. Layer 2, not yet code.
- `` `fn arg` `` (backticks) — a **shipped binary**: the FO calls it.

So migrating a step is literally `«state.commit»(slug)` → `` `spacedock state commit <slug>` ``. The flow
that invokes it never changes; only the function body flips from a hand-followed recipe to "call the verb."

Scoped here to the boot/state-repo slice (the 9-step Startup, read every session).

---

## BEFORE (today: inline procedure, mechanics inline)

Startup steps 5-8 + the State-Management commit discipline, the FO hand-executes:

> 6. **Split-root state halt-gate.** If `state_backend == split-root` AND `entity_dir_present == false` [...]
>    HALT dispatch, report "state not initialized," run `spacedock state init` [...] Re-read `--boot`.
> 7. **Split-root pull-on-boot.** `git -C <state-path> pull --rebase origin <state-branch>` [...] On CONFLICT, halt.
> 8. **Merged-PR sweep.** For each `pr_state` entry `MERGED` + non-terminal, run the pr-merge startup advancement.
>
> **Concurrency-safe state commits.** Fallback: `git -C {state} add {path} && commit -m "…" -- {path}`.
> Never bare `git add -A`. After commit → push. On push rejection → `pull --rebase` then re-push.

The reader extracts the intention ("get the state repo ready, then greet"; "land one entity's change
safely") from the git mechanics. The FO re-derives the git dance every time.

---

## AFTER (prose-functions: intention in the flow, mechanics in the function body)

### State repo (operating principle)

The state repo is the session's source of truth. **The FO declares intent against it by invoking the
prose-functions below; their bodies own the mechanics.** Each function is idempotent — re-invoking
checks its *done-when* and is a no-op if already satisfied — so a sequence converges rather than runs
as a script, and each function can become a binary independently without touching the flow.

**Boot converges the state repo to ready, then greets:**

1. `«state.boot»()`
2. `«state.ensure-ready»()`
3. `«state.sweep-merged»()`
4. greet from the boot read, then stop.

**Every state write is one call:** `«state.commit»(slug)`.

### Declarations

## «state.boot»(): read all startup state in one call
Yields the boot record (mods, id style, orphans, `pr_state`, dispatchable, team, `state_backend`).
→ **shipped**: `` `spacedock status --boot --json` `` — already a binary; invoke it directly.

## «state.ensure-ready»(): the split-root checkout is linked and integrated with peers before any dispatch
- **guard:** `state_backend == split-root`
- **effect:** `entity_dir` absent → halt + `«state.init»()`; else `git -C {state} pull --rebase origin {branch}`
- **done-when:** `entity_dir_present` ∧ rebased clean
- **block:** rebase conflict → halt, surface the conflict, stop (no dispatch on an unmerged tree)
- → **prose**, becomes `` `spacedock state ready` ``

## «state.sweep-merged»(): merged PRs reach their terminal stage at boot
- **guard:** `pr_state` has an entry `MERGED` with a non-terminal entity
- **effect:** run the pr-merge startup advancement (clear `mod-block`, terminalize PASSED, archive, drop worktree)
- **done-when:** no `MERGED` + non-terminal entity remains
- **block:** `gh` unavailable → skip, report merge state UNKNOWN
- → **prose**, becomes `` `spacedock state sweep` ``

## «state.commit»(slug): record one entity's change durably and concurrency-safe
- **effect:** path-scoped `git -C {state} add {path} && commit -- {path}`; push; on reject `«state.sync»()` + re-push; retry on `index.lock`
- **done-when:** the entity's change is pushed to `origin`
- **block:** rebase conflict on sync → halt per the rebase-conflict rule
- → **prose**, becomes `` `spacedock state commit <slug>` ``

---

## The migration property

The flow (1-4 above) is final prose. Nothing in it changes as the binary grows:

- `«state.boot»` is already shipped; the other three are prose (the FO follows their bodies by hand, as today).
- Ship `` `spacedock state commit <slug>` `` → the `«state.commit»` body collapses to "call the verb" and the
  notation flips to backticks. Every `«state.commit»(slug)` in the contract is now one binary call; the
  ~15-per-session git dance is gone. **No other prose moves.**
- Later fold `«state.ensure-ready»` + `«state.sweep-merged»` into `` `spacedock state ready` `` → two more flips.

The contract shrinks one function at a time, never a rewrite, and the intention (boot = converge-then-greet;
write = declare-intent) stays legible throughout. Same treatment applies to the dispatch 9-step and the
merge ceremony — same imperative-prose smell, same prose-function cure.
