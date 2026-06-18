---
title: Command-surface gate — enforce the anti-sprawl routing principles in the binary (allowlist, not prose)
status: backlog
source: captain (2026-06-17) — the command-surface-and-routing-principles.md proposal is PROSE; a prose rule's ceiling is "the wording is present" (the same failure mode as the Shaping-FO/Commander role-boundary slip). Make it a CODE gate.
score: 0.4
sprint: 0205-layered-fo
group: surface
sprint-readiness: ready
id: dwd8hhn898apwgay08mdj7p6
---

The command-surface + routing principles (`docs/dev/_proposals/command-surface-and-routing-principles.md`) are currently PROSE — a rubric reviewers are supposed to apply. Per the project's own code-gate-over-prose discipline, a prose rule does not hold under pressure. As 0205+ ADDS verbs (`state`, `merge guard`, `gate-extract` modes, later `next-action`), the surface sprawls by default unless a code gate enforces the principle.

Build the binary/contractlint gate that makes the anti-sprawl rule HOLD:
- An EXPLICIT ALLOWLIST of top-level commands (the `root.AddCommand` set at internal/cli/cli.go:141), pinned by a test in `internal/contractlint` (the precedent package for structural gates — boot-resident closure, boundary guards, reconcile-class binding). Adding a top-level command REQUIRES a deliberate allowlist edit → forces the decision-tree justification at review + fails the build if a command is added without the deliberate step.
- BIND the allowlist to the completion `verbs=` string (cli.go:549) — these are two INDEPENDENT values that can drift (a legitimate cross-check, NOT a prose-grep tautology per the Proof policy: the expected value comes from outside the file under test). The test fails if the registered commands and the advertised completion verbs disagree.
- Reference `command-surface-and-routing-principles.md` as the rubric the gate enforces (the prose explains WHY; the test enforces THAT).

This sprint is the proof case: `merge guard` introduces a NEW top-level `merge` command group — the gate would FORCE that into the allowlist deliberately (the merge-finalize gate confirms "merge earns a top-level group: distinct loop phase"), while `state ready`/`commit` are SUBcommands under the existing `state` noun (no new top-level). So the gate catches exactly the "should this be a new command?" decision the principles describe.

Out of scope: enforcing the FULL decision tree (flag-vs-mode-vs-mod) in code — that stays human judgment at the gate. This gate enforces the one mechanical invariant: the top-level surface is an explicit, deliberate allowlist that AGREES with the advertised completion verbs.

Acceptance criteria (ideation fleshes; behavior-first / oracle-based, NO prose-grep): a contractlint test that RED on (a) a top-level command added without an allowlist entry, and (b) the allowlist and completion verbs disagreeing — proven by exercising a stray-command fixture (RED) vs the real surface (GREEN), asserting exit/failure. The allowlist value is external to the cli.go file under test (the test owns the expected set), so it is not a tautology.
