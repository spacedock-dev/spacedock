---
id: jhazq8v9prphgs0bk5pn0xfw
title: Dispatch-build flag/file runtime docs must not outpace accepted binaries
status: backlog
source: "FO dogfood (2026-06-06) - plugin 0.19.5 Codex runtime instructs `dispatch build --checklist-file`, but installed spacedock 0.19.4 still accepts only stdin JSON while satisfying contract 1."
score: "0.27"
worktree: ""
issue:
sprint: 019x-pre-flip-cleanups
group: dispatch-hygiene
sprint-readiness: ready
---

The Codex first-officer runtime now instructs the FO to build dispatch prompts
with the flag/file form:

```bash
spacedock dispatch build --workflow-dir DIR --entity-path FILE --stage STAGE --checklist-file FILE
```

That works in current source, but the installed `spacedock 0.19.4 (contract 1)`
still treats the request as stdin JSON and fails with `invalid JSON on stdin`.
The startup contract gate accepts both binaries because they both advertise
`contract 1`, so the skill can put an FO into a command form the accepted binary
does not support.

This is a runtime/source skew problem, not a Codex host-dispatch leak. Current
source with `--host codex` emits a Codex-shaped `Read /tmp/...` prompt; the old
binary simply does not support the flag/file input mode.

## Acceptance criteria

**AC-1 - Runtime instructions are compatible with accepted binaries.**
Verified by a fixture or version-compatibility test that fails if shipped
first-officer runtime text requires a dispatch-build flag form unsupported by
the minimum accepted binary/version.

**AC-2 - Feature detection or fallback exists for older accepted binaries.**
Verified by a command-level test or runtime-text test showing the FO either uses
stdin JSON when flag/file mode is unavailable, or the startup gate rejects the
older binary with a clear upgrade instruction before dispatch.

**AC-3 - Current Codex host dispatch remains Codex-shaped.**
Verified by a dispatch-build test using `--host codex` that asserts the emitted
prompt is `Read /tmp/...` and does not contain `Skill(skill="spacedock:ensign")`.

## Stage test gates

- Ideation should decide whether this is solved by raising the binary/plugin
  compatibility gate, probing `dispatch build --print-schema` or `--help`, or
  documenting a stdin-JSON fallback in the runtime adapter.
- Implementation should include a regression test that models `0.19.4`-style
  stdin-only behavior or otherwise binds the minimum supported command surface.
- Validation should run focused dispatch-build tests plus `go test ./...`.
