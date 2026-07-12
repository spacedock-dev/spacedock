---
title: Preserve Zellij targeting metadata through Safehouse
status: implementation
score: 0.55
source: "Captain request 2026-07-12."
id: 6jk4gverktmthbzkn4vas1kg
worktree: .worktrees/spacedock-ensign-safehouse-zellij-env-pass
started: 2026-07-12T04:17:30Z
---

Safehouse sanitizes the environment passed to a wrapped Spacedock host. Tooling
inside a Zellij-hosted sandbox needs the current Zellij identity in order to
address the correct pane and session.

## Chosen direction

For every Safehouse-wrapped Spacedock front door, forward the inherited values
of `ZELLIJ`, `ZELLIJ_PANE_ID`, and `ZELLIJ_SESSION_NAME` through Safehouse's
existing named environment allowlist. The variable *names* are a built-in,
conditional compatibility allowlist; no values are hardcoded or manufactured.
When the parent process is not in Zellij, do not add those names to the
Safehouse argv.

Keep Safehouse's native `SAFEHOUSE_ENV_PASS` as the operator's global mechanism
for any additional, explicitly chosen variables. Do not add a competing
Spacedock config file, and do not put environment forwarding in the trusted
repository `.safehouse` profile.

## Acceptance criteria

**AC-1 (value): A Safehouse-wrapped Spacedock host launched from Zellij receives
the exact inherited values of `ZELLIJ`, `ZELLIJ_PANE_ID`, and
`ZELLIJ_SESSION_NAME`, so in-sandbox tooling can identify and target the active
pane/session.** Verified by: a scrubbed-environment executable smoke whose host
prints the three values, plus a control that omits the named pass-through.

**AC-2 (scope): A non-Zellij parent neither manufactures Zellij values nor
changes the existing wrapped-launch argv beyond the established
`SPACEDOCK_BIN` behavior.** Verified by: deterministic front-door argv tests
with and without an explicit Zellij parent environment.

**AC-3 (compatibility): The built-in Zellij trio composes with the native
Safehouse global allowlist for unrelated operator-selected variables, without a
new Spacedock configuration surface or broad `--env` inheritance.** Verified
by: a focused Safehouse integration fixture and scoped code review of the
launcher/config surfaces.

## Test plan

1. Refactor the Safehouse env-pass composition into an environment-aware helper
   with explicit parent-environment input so tests do not depend on the
   developer's terminal.
2. Extend the scrubbed Safehouse smoke to prove exact forwarding of all three
   Zellij values and prove the stripped control cannot see them.
3. Cover each Safehouse-wrapped front door that uses the shared helper; preserve
   the unchanged non-Zellij argv control.
4. Run `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal`.

## Out of scope

- Hardcoding `ZELLIJ=0`, a pane ID, or a session name.
- Passing the full host environment or secrets through Safehouse.
- Adding a Spacedock-specific persistent configuration format.

## Stage Report: ideation

- DONE: The captain approved a built-in conditional allowlist for the three
  Zellij targeting variables, retaining Safehouse's native global
  `SAFEHOUSE_ENV_PASS` as the escape hatch for additional variables.
- DONE: The smallest mechanism is the existing Safehouse named environment
  forwarding surface; a new Spacedock config format would duplicate it and make
  security-sensitive allowlisting harder to audit.
