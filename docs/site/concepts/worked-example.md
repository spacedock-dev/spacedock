# A worked example

This page traces one real entity, `z9` `codex-plugin-auto-install`, from
backlog through to `done` / PASSED, using artifacts that live in this repo. It is
a concrete read of the abstract stage machine: backlog → ideation →
implementation → validation → done, the gates between them, and what the captain
actually decides at each one.

The workflow is `docs/dev` (the Spacedock v1 dev workflow); its stages, gates,
and entity schema are defined in
[the development workflow reference](../contributing/development-workflow.md).
Runtime entity state lives in a separate `.spacedock-state` checkout, so the
finished entity itself is not in the main tree, but its full trajectory is on
the record in the `0198-pre-flip-hardening` sprint directory: `index.md`,
`dispatch-sprint-execution.md`, `debrief.md`, and `post-sprint-audit.md`.

`z9` delivered front-door Codex plugin auto-install: `spacedock codex` now
installs a missing plugin then launches, the Codex analog of the Claude path. It
shipped as [PR #329](https://github.com/spacedock-dev/spacedock/pull/329).

## See where it sits

Each first-officer loop starts by asking the workflow what is ready to move. You
run the same query:

```bash
spacedock status --workflow-dir docs/dev --next
```

Lists the entities ready to dispatch. To scope to one sprint, filter with
`--where`:

```bash
spacedock status --workflow-dir docs/dev \
  --where sprint=0198-pre-flip-hardening --where 'sprint-readiness != defer'
```

This is the sprint membership query, the source of truth, not a hand-kept list.
At the start of `0198`, `z9` shows up here in the `binary-ux` group.

## backlog → ideation: shape the work

`z9` enters backlog as a seed: a title, a source, and a brief description of the
problem. It carries no design yet. Ideation is where a worker fleshes it out into
a problem statement, a proposed approach, entity-level acceptance criteria, and a
test plan, with each AC naming a check outside the entity body that can fail.

For `z9` that meant a concrete approach: install through the shared `devBranch`
rather than a hardcoded `"next"`, so the install channel tracks the release
channel (`next` today, `main` after the flip). The design also resolved the
riskiest unknown up front (that the install branch was the channel variable, not
a literal) and recorded it before committing to the rest of the plan.

Ideation ends at a **gate**. Because the dev workflow marks `ideation` with
`gate: true`, the first officer does not advance on its own: it presents the
design to the captain. The captain approved `z9`'s ideation on 2026-06-08. That
approval is the entry condition for implementation, and it is captured in the
sprint package so a later Commander session drives implementation directly
without re-presenting the gate.

## implementation: produce the deliverable

Once the design is approved, `z9` moves to implementation. The dev workflow runs
this stage in a `worktree` (`worktree: true`), so the dispatched ensign works in
an isolated checkout, not the shared tree.

The dispatch package gave the implementer three concrete build notes, all of
which the work honored:

- **Install on the shared `devBranch`**, via `ops.Install("codex",
  marketplaceSource, devBranch)` and `--ref <devBranch>`, not a hardcoded
  `next`, so the channel tracks the flip's later `devBranch` retarget.
- **Fix the now-false comments and error strings** it builds around
  (`host_exec.go`, `frontdoor.go`), rather than adding new code around stale text.
- **Invert the obsolete test** `TestCodexFrontDoorNoPluginFailsFastWithoutInstalling`,
  whose old assertion contradicted the new auto-install behavior.

Implementation is complete when the deliverable is committed and the stage report
is filed. It is not a parking spot: a completed implementation routes straight to
`validation` dispatch.

## validation: verify against the criteria

`z9` moves to validation to be checked, not finished. The dev workflow marks
`validation` with `fresh: true` and `worktree: true`, so a fresh validator (one
that did not write the code) runs in its own worktree. It pulls every `AC-N`
from the entity's acceptance-criteria section, reproduces the evidence each one
cites, and produces a PASSED or REJECTED recommendation.

`z9` is a front-door change, a high-stakes surface, so validation alone was not
sufficient. The dev workflow requires a **detached adversarial audit** for the
launcher front door: a read-only pass on a throwaway checkout of the merge result
that tries to refute the validation. The `z9` audit ran at commit `0b714fac`
and exercised five mandated probes, each reddening the test suite then reverting.
It refuted nothing material. Channel-tracking was confirmed clean. That clean
audit satisfied the sprint's DoD item for the high-stakes surface.

Validation also has `gate: true` and `feedback-to: implementation`. A REJECTED
recommendation routes the finding back to implementation for another cycle; a
PASSED recommendation goes to the captain for the terminal decision.

## done: the captain's verdict

`z9` reaches `done` when the captain reads the validation report and approves.
The terminal stage records the outcome in frontmatter: `status: done`,
`verdict: PASSED`, and a `completed` timestamp. The sprint's `post-sprint-audit.md`
confirms `z9` finished `status: done`, `verdict: PASSED`, archived, with PR
reference #329, alongside its three sprint siblings (`kb` #327, `qa` #328,
`vh` #330).

That is the whole point of the machine: nothing reached `done` on assertion
alone. `z9` advanced only on an approved design, an isolated implementation, a
fresh validation, a detached audit that failed to break it, and an explicit
captain verdict. Each one a decision, each one on the record.

## Read the trail yourself

The full trajectory is reconstructable from the sprint directory:

- [`index.md`](https://github.com/spacedock-dev/spacedock/blob/next/docs/roadmap/0198-pre-flip-hardening/index.md): goal, members, definition of done, the approved-gate note.
- [`dispatch-sprint-execution.md`](https://github.com/spacedock-dev/spacedock/blob/next/docs/roadmap/0198-pre-flip-hardening/dispatch-sprint-execution.md): the per-member drive plan and `z9`'s build notes.
- [`debrief.md`](https://github.com/spacedock-dev/spacedock/blob/next/docs/roadmap/0198-pre-flip-hardening/debrief.md): what shipped, the PR links, and the decisions made along the way.
- [`post-sprint-audit.md`](https://github.com/spacedock-dev/spacedock/blob/next/docs/roadmap/0198-pre-flip-hardening/post-sprint-audit.md): the final-state confirmation and the detached-audit record.
