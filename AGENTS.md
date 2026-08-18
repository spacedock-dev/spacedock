# Agent Instructions

This repo builds the Go-based `spacedock` launcher and the project-side skill integration behind its stable command surface.

## Priorities

- Prefer small Go packages with clear boundaries. Avoid a single large CLI file.
- Use the standard library unless a dependency removes real complexity.
- Keep command output stable and test it with fixtures.

## Expected Commands

Run these before claiming work is complete:

```bash
go test ./...
go test ./... -race
gofmt -w ./cmd ./internal
```

Add focused tests for each stage before implementing that stage.

## Releasing

Cut stable releases from `main` via an annotated `vX.Y.Z` tag — see `docs/releasing.md`. `next` is a dev-only source-build convenience branch (`go install …@next`, `--plugin-dir`) — the edge marketplace entry resolves `main` directly, so `next` is not a re-pull source for any installer. Do not cut a stable `vX.Y.Z` tag from it.

## Project Shape

- `cmd/spacedock/`: process entry point only.
- `internal/cli/`: command routing, usage text, exit-code behavior.
- `internal/status/`: status implementation.
- `docs/specs/`: design contracts, including the state behavior extension.
- `docs/roadmap/`: bootstrap and migration roadmap.
- `docs/dev/README.md`: development workflow definition.
- `docs/dev/.spacedock-state/`: development workflow entities in a separate state checkout.
- `skills/`: first-officer, ensign, and satellite skills plus integration fixtures.

## Skill Development

- Skills should call `spacedock`, not plugin-private script paths.
- Keep skill instructions declarative. Let the binary own path resolution and mutation guards.
- Add skill smoke tests before changing first-officer or ensign command text.
- Preserve current FO/ensign write-scope rules: the first officer mutates entity state; ensigns write assigned code, reports, and artifacts.

## Runtime Support

- When adding a new runtime host or debugging first-contact runtime friction, read `docs/runtime-support.md` first.
- Use the documented "assume it already works" operating prompt before declaring a host impossible due to auth setup, extension/package discovery, or tool-shape mismatch.
- Prove runtime claims with live or fixture-backed durable state evidence - see `docs/runtime-support.md`.
