---
entity-type: ticket
entity-label: ticket
entity-label-plural: tickets
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

# Ticket-label drive fixture

A throwaway workflow declaring `entity-label: ticket`. A live first-officer pointed
here and driven to present the gate for `001-add-login-rate-limit` (sitting at the
gated `review` stage with a `## Stage Report`) should call the entities "ticket(s)"
in its generated captain-facing prose, where it would otherwise write "entity".

The word "ticket" appears in this README's `entity-label` field only — the FO must
read and resolve it at boot; the drive never hands the FO the label.

## Stages

### `backlog`

A ticket starts in backlog awaiting triage.

- **Outputs:** A triaged ticket ready for implementation.

### `implementation`

The work to satisfy the ticket. Produces the deliverable on a worktree branch.

- **Outputs:** The committed deliverable plus a stage report.

### `review`

A reviewer checks the implementation against acceptance criteria and recommends a
verdict. This stage is gated: the first officer presents the report to the captain.

- **Outputs:** A recommended verdict (approve/reject) with findings.

### `done`

Terminal. The ticket is merged and archived.
