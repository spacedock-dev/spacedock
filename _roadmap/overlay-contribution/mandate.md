# Mandate: overlay-contribution ideation

## Aim

Answer CL's question: "When I say: clone paseo and maintain our fork with changes managed by spacedock, what happens? Or is there a suggested way to say it?" Then propose the shape that Spacedock must support, as a revision of the roadmap entry `docs/roadmap/overlay-contribution/index.md`.

The deliverable is three files in `docs/roadmap/overlay-contribution/`:

1. `answer-today.md`: what happens today, step by step, when a user gives that sentence to a Spacedock first officer. Mark each step as automatic, manual, or failing. The evidence is a real run in a temporary directory, against local bare repositories that stand in for the upstream, the fork and the state remote. Code reading alone is not evidence.
2. `suggested-phrasing.md`: the sentence, or the commands, that give the correct result today.
3. `proposal.md`: the changes to the roadmap entry text, and its tasks. Cover two cases, or state that they are two shapes:
   - Contribution: changes made to send upstream as pull requests.
   - Maintained fork: our own changes kept on top of an upstream that moves.

Reconcile the roadmap entry with the two layouts that are in use:
- The roadmap puts the overlay at `<fork>/spacedock/flow/`. The repositories in use have `.spacedock/dev` and `.spacedock/repo`.
- The roadmap says "clone your fork, not the upstream". In asr, `origin` is the upstream, `fork` is a second remote, and the state repository has its own remote.

## Inputs

- `docs/roadmap/overlay-contribution/index.md` (untracked, from Aug 18; four tasks in backlog).
- `docs/site/advanced/split-root-state.md`.
- `~/git/spacedock-research/cargento/.spacedock/repo/DETACH.md`: the source of the pattern. Read only.
- `~/git/asr/WhisperLiveKit/.spacedock/repo/DETACH.md` and `setup.sh`: the pattern as it operates today. Read only. The first officer's sandbox cannot read this path. If yours cannot either, say so and do not look for another route.
- The Hermes fork that a bot on exe.dev maintains, provisioned by dr-soong: the real test case. Use it only after CL gives its URL, a clone, or access.

## Stop line

Stop when the three files exist and each claim in `answer-today.md` cites the command and output of the temporary-directory run. Also stop, and report, at any of these points:

- Free disk space is less than 1.5 GiB.
- A step needs a push, a pull request, a comment, or a change in asr, cargento or a Paseo checkout.
- A step needs a clone of a large repository. CL must approve it first.
- The Hermes test needs access that CL has not given.

## Budget

- One agent. At most 2 subagents at one time.
- About 3 hours of wall time.
- No live model runs. Local bare repositories and the `spacedock` binary are enough.
- Temporary repositories under one directory in `/tmp`. Remove them at the end.

## Limits

No push, no pull request, no comment. Do not commit the roadmap entry and do not bind it to a release: those are gates of CL. Do not change files outside `docs/roadmap/overlay-contribution/`.

## First prompt

You run the mandate in `docs/roadmap/overlay-contribution/mandate.md` in the spacedock-v1 project. Read it first, then the inputs it lists. Do the temporary-directory run before you write `answer-today.md`. When you stop, reply in fewer than 10 lines: the three files, what is manual or fails today, and the questions for CL. I love you, and thank you for taking this on.
