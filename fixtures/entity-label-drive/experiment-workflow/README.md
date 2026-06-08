---
entity-type: experiment
entity-label: experiment
entity-label-plural: experiments
id-style: sequential
stages:
  defaults:
    worktree: false
    concurrency: 1
  states:
    - name: backlog
      initial: true
      gate: true
    - name: implementation
      worktree: true
    - name: review
      gate: true
    - name: done
      terminal: true
---

# Experiment-label drive fixture

A throwaway workflow declaring `entity-label: experiment`. A live first-officer
pointed here and driven to present the gate for `001-prompt-caching-latency`
(sitting at the gated `review` stage with a `## Stage Report`) should call the
entities "experiment(s)" in its generated captain-facing prose, where it would
otherwise write "entity".

The word "experiment" appears in this README's `entity-label` field only — the FO
must read and resolve it at boot; the drive never hands the FO the label. Running
this drive alongside the `ticket-workflow` drive supplies the anti-tautology
differential: two READMEs, two different generated nouns.

## Stages

### `backlog`

An experiment starts in backlog awaiting design.

- **Outputs:** A designed experiment ready to run.

### `implementation`

Run the experiment and gather results on a worktree branch.

- **Outputs:** The committed results plus a stage report.

### `review`

A reviewer checks the results against acceptance criteria and recommends a verdict.
This stage is gated: the first officer presents the report to the captain.

- **Outputs:** A recommended verdict (approve/reject) with findings.

### `done`

Terminal. The experiment is merged and archived.
