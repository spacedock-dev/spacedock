# Mandate: Paseo as a first-class Spacedock backend (probe)

## Aim

Find out what it takes for Spacedock to launch its workers (ensigns) through Paseo, using the Paseo CLI or Paseo MCP, instead of the host's native subagent mechanism (Claude Code `Agent()`, Codex, Pi).

The deliverable answers these questions, with evidence from a real run:

1. **Bindings.** Which Paseo operations implement each dispatch capability the first officer contract names in `skills/first-officer/references/fo-dispatch-core.md`: `«worker.spawn»`, `«addressable-worker»` (mid-run messages), `«async-dispatch»`, `«completion-signal»`, `«worker.shutdown»`, `«context-budget»`, and `«roster-reconcile»`. For each one, say whether Paseo supports it today, partly, or not at all.
2. **Worktrees.** Who owns the worktree: `spacedock dispatch build --stamp` or a Paseo worktree workspace. What breaks when the two disagree. Note that `--stamp` today creates worktrees from the local trunk branch, not from `origin/<trunk>`.
3. **Cross-runtime.** Can a first officer on one host (for example Claude Code) dispatch ensigns that run on another (for example Codex) through Paseo? What does `spacedock dispatch build` need to emit for that?
4. **Access and sandboxing.** What daemon access the first officer needs, which Safehouse labels apply, and how an ensign launched by Paseo gets the correct working directory and write scope.
5. **Completion and recovery.** How the first officer learns that a Paseo-launched ensign finished, failed, or stalled, and how the dispatch-failure retry rung maps onto Paseo.

## Deliverable

Files in `docs/roadmap/paseo-backend/`:

- `probe-report.md`: the answer to questions 1-5. Each claim cites the command and its output from the probe run.
- `run-log.md`: the commands and outputs of the probe run.
- `proposal.md`: the shape of a `paseo` runtime adapter (a new `references/paseo-fo-dispatch.md` binding, or a dispatch-build output mode), what changes in the binary, and the tasks, with the riskiest one first.

## Method

Read `docs/runtime-support.md` first. Follow its "assume it already works" operating prompt before you declare that something is impossible.

The probe run: in a temporary directory, create a small split-root fixture workflow with one task. Run `spacedock dispatch build` for one stage. Launch the ensign with the Paseo CLI or MCP from the emitted envelope, using Codex with model `astra-6.2`. Observe the stage report, the completion signal and the shutdown. Record what you had to do by hand.

## Stop line

Stop when the three files exist and every claim in `probe-report.md` cites run evidence. Also stop, and report, at any of these points:

- Free disk space is less than 1.5 GiB.
- A step needs a push, a pull request, or a comment.
- A step needs a change in another project's files or agents.
- A step needs daemon access that this workspace does not have.

## Budget

- One mandate agent: Codex, model `astra-6.2`. At most 2 Paseo-launched probe ensigns at one time, and at most 6 in total.
- About 3 hours of wall time.
- Temporary repositories under one directory in `/tmp`. Remove them at the end.

## Limits

No push, no pull request, no comment. Do not change product code in this repository. The deliverable is the three files. Start and stop only the probe agents this mandate creates, and stop all of them before you finish.

## Workspace and labels

- Workspace title: `paseo-backend probe`, under the spacedock-v1 Paseo project, with the project root as its directory and local isolation.
- Agent labels: `role=worker`, `repo=spacedock-v1`, `mandate=paseo-backend`, `paseo-safehouse.paseo-access=daemon`. The daemon label is needed because the probe launches ensigns through Paseo.

## First prompt

You run the mandate in `docs/roadmap/paseo-backend/mandate.md` in the spacedock-v1 project. Read it first, then `docs/runtime-support.md` and `skills/first-officer/references/fo-dispatch-core.md`. Do the probe run before you write `probe-report.md`. When you stop, reply in fewer than 10 lines: the three files, which dispatch capabilities Paseo supports today, the riskiest gap, and the questions for CL. I love you, and thank you for taking this on.
