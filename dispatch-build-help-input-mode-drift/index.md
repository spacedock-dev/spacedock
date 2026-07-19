---
id: 2690fpqe9pkn917am6bt6eqs
title: Make dispatch-build help match its input-mode parser
status: ideation
source: "FO dogfood, 2026-07-19: 0.26.0-pre0 help advertised stdin JSON plus --advance but the same invocation selected flag/file mode and rejected stdin."
started: 2026-07-19T05:04:28Z
completed:
verdict:
score: "0.65"
worktree:
issue:
milestone: 0.26.0
group: binary-ux
---

Make `spacedock dispatch build --help` describe every supported input form, its
mode-selection rules, and a complete reuse-advance invocation that the parser
actually accepts.

## Problem

The current help prints only the stdin-JSON usage, omits the accepted
`--entity-path`, `--stage`, `--checklist-file`, `--scope-notes-file`, and
`--feedback-context-file` flags, then lists `--advance` without explaining that
it participates in request-flag detection. A caller following that surface can
pipe a complete JSON request and add `--advance`, only to receive:

```text
error: flag/file input requires --entity-path, --stage, and --checklist-file
```

The implementation selects flag/file parsing whenever `hasRequestFlags()` is
true, so stdin is ignored in that case. The help neither exposes that boundary
nor supplies the complete flag/file form required by the First Officer
contract. This is a recurring dispatch-boundary round trip.

The related active task `dispatch-build-flag-form-version-skew` and closed
GitHub issue #313 concern an older accepted binary versus newer plugin
instructions. This task concerns one current binary contradicting its own help.

## Acceptance criteria

- **AC-1:** `dispatch build --help` documents the complete stdin-JSON and
  flag/file forms, including all accepted request flags and the exact rule that
  selects between them. Verified by running the rendered examples against the
  real parser, not by grepping source prose.
- **AC-2:** The supported reuse-advance form is unambiguous: either stdin JSON
  plus `--advance` works, or help explicitly rejects that combination and gives
  a complete accepted flag/file example. Verified by paired success/failure
  command tests asserting stdout, stderr, and exit status.
- **AC-3:** A behavioral help-example test fails whenever an example printed by
  `--help` no longer parses successfully with a minimal workflow/entity fixture.

## Scope

Keep the fix to command help, input-mode selection if necessary, and
load-bearing behavioral tests. Do not redesign the dispatch envelope or add a
third request form.
