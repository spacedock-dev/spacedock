---
title: Bound direct-launcher filing evidence to one shell command segment
status: backlog
score: 0.2
source: "Split from 88t by captain direction on 2026-07-11: the separator counterexample is low-realism test-oracle hardening and does not block the named-function release."
completed:
verdict:
worktree:
issue:
id: 6p740rw6era1xz9wmg29qfay
---

Fresh validation of 88t found that the test-only filing evidence detector can combine a direct launcher token from one shell command with `new <slug>` from a later command after a separator. The captain accepted 88t because this artificial stream does not represent the user-facing failure mode and directed that the matcher hardening be tracked separately.

## Scope

Unify direct literal and captured-variable launcher recognition through one bounded simple-command segment path. Retain `;`, `&&`, `||`, and `|` counterexamples as negative controls. Preserve existing positive filing evidence and make no product-runtime or first-officer contract change unless a real behavior case is demonstrated.

## Acceptance criteria

- Direct and captured launchers use the same bounded command-segment matcher.
- Separator counterexamples for `;`, `&&`, `||`, and `|` remain uncredited.
- Existing legitimate direct and captured `spacedock new` evidence remains credited.
- Focused detector tests and repository gates pass without changing shipped runtime behavior.
