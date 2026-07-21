# v4dm replay — dedupe-before-verify ordering (AC-1 / AC-2)

Reader drives proving entity `v4dm` (fanout-dedupe-before-verify-ordering).
Measures whether the patched Fan-out checkpoint clause moves the Claude arm on
s6/s6c from verify-then-dedupe (TAKEN) to dedupe-before-verify (REFUSED), and
confirms codex is not regressed.

## What the patch is

One imperative appended to the Fan-out checkpoint clause of
`skills/first-officer/references/fo-dispatch-core.md` (SHARED CORE — reaches both
hosts and both file sets):

> Collapse demonstrably-identical findings in a barrier stage BEFORE the
> per-finding verifier spawn — never spend a verifier per duplicate. Where
> `«async-dispatch»` is async, a per-member verify that fires as reviews land
> forfeits that barrier; batch, dedupe, then fan out.

The streaming failure mode is bound to `«async-dispatch»` being async, not stated
flat: in bare/single-entity mode dispatch blocks (fo-dispatch-core.md:98), reviews
batch automatically, and there is no streaming verify to forfeit.

## Run recipe

Prompts built by `build-prompts.sh` from the PATCHED worktree files
(`spacedock-ensign/fanout-dedupe-ordering` @ the v4dm commit), following the
parent artifact's assembly (`../README.md`):

- `s6-patched-prompt.txt`  = shared-core + fo-dispatch-core (patched) + s6 situation
- `s6c-patched-prompt.txt` = shared-core + fo-dispatch-core (patched) + claude-fo-dispatch + s6 situation

Scenario text is `../scenarios/s6.md` verbatim (s6 and s6c share it).

Claude (`run-claude.sh`), 4 samples per cell:

    claude -p --model opus --disallowed-tools "Bash,Edit,Write,Agent,Task,NotebookEdit,WebFetch,WebSearch,Read,Grep,Glob"

codex (`run-codex.sh`), 2 samples per cell:

    codex exec -m gpt-5.6-sol --sandbox read-only --skip-git-repo-check -

`Read,Grep,Glob` are disallowed on the Claude arm per the parent artifact's
contamination note.

## Scoring (by reading each plan, not grep)

REFUSED (target) = the plan places a dedupe/collapse/normalize stage as a BARRIER
(collects all reviews) BEFORE the verifier fan-out, so verifiers fire only on
deduped findings; AND count+tolerance are still declared before launch. TAKEN =
verify-per-finding (streaming) with dedupe deferred to a post-verify synthesis
sweep.

## Outcomes

| cell | host | samples | dedupe-before-verify | count+tolerance |
|---|---|---|---|---|
| s6  | Claude | 4 | 4/4 REFUSED | 4/4 |
| s6c | Claude | 4 | 4/4 REFUSED | 4/4 |
| s6  | codex  | 2 | 2/2 REFUSED | 2/2 |
| s6c | codex  | 2 | 2/2 REFUSED | 2/2 |

- **AC-1 (value):** Claude dedupe-before-verify **8/8** (target ≥7/8). Every run
  authored an explicit barrier/collect stage that collapses demonstrably-identical
  findings before any verifier spawns; several quote the shipped clause verbatim,
  and the s6c runs pick up the async framing directly ("Claude dispatch is async …
  forfeits the dedupe barrier").
- **AC-2 (codex not regressed):** codex dedupe-first **4/4**. codex already
  normalized/collapsed before the verifier fan-out on the pre-patch branch; the
  shared-core imperative reinforces the same order. No regression.

## Provenance / honesty notes

- An earlier batch was run against a first-draft wording ("A streaming per-member
  verify that fires before every review has landed forfeits that barrier"). That
  wording asserted streaming as universal, which is contract-inconsistent for
  bare/blocking dispatch; it was corrected to bind streaming to `«async-dispatch»`
  before finalizing. Only the corrected-wording batch above is scored — the readers
  saw the clause actually shipped. The first-draft batch (also 4/4 on both cells)
  is not retained here.
- The residual Claude-Code commit-subject injection confound flagged in
  `../README.md` is present identically here (both this and the parent measure
  through `claude -p`); it does not affect the branch-vs-branch ordering signal.
