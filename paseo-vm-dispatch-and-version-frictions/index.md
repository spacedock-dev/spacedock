---
id: kdpzm9pckwj4pm1djvaf01ce
title: "Paseo stack Layer 3: dev build SHA in --version, collision-free immutable dispatch files, and --dispatch-dir"
status: backlog
source: "CL via MetaFO (2026-11-23, Studfleet VM dogfooding /tmp/paseo-frictions-studfleet-vm-2026-11-24.md)"
started:
completed:
verdict:
score: 0.90
worktree:
issue:
pr:
mod-block:
---

Fix the Spacedock-owned frictions surfaced by real Paseo backend dogfooding on the Studfleet VM (`studfleet-paseo`, PR #897 commit `5210d36d0`), stacked cleanly as Layer 3 atop PR #897 (`paseo-and-generic-harness-fo-support`):

## Problem

1. **Friction #1 — Dev builds hide the running PR commit SHA:** `internal/cli/dev_version.go` emits `0.29.0-pre0+dev` for unstamped `go build`/`go install` binaries (`Version == "dev"`), hiding which PR commit is installed when testing a stack build.
2. **Friction #4 — Concurrent briefs for the same `(parent, entity, stage)` overwrite each other (reported on both Paseo and Claude):** `internal/dispatch/build.go:627-655` keys the filename only on optional `CLAUDE_CODE_SESSION_ID`/`PASEO_AGENT_ID` + `derivedName` (+ `-advance`), not on brief content — so two simultaneous briefs from one parent for the same task overwrite the same mutable file.
3. **Friction #5 — Hardcoded `/tmp/spacedock-dispatch` breaks private `/tmp` namespaces:** `internal/dispatch/build.go:25` hardcodes `dispatchFileDir = "/tmp/spacedock-dispatch"`, so under a private scratch `/tmp` the emitted pointer still spells the private namespace's `/tmp` instead of a shared project/scratch directory.
4. **Operational notes (#3, #6):** Document why `--isolation worktree` writes to `~/.paseo/worktrees/` (reinforcing `--isolation local` adoption in `paseo-fo-dispatch.md`), document `--dispatch-dir`/`SPACEDOCK_DISPATCH_DIR`, and note verifying `inspect.Thinking` after spawn rather than trusting `paseo models` thinking metadata alone.

Non-Spacedock / already-handled items (out of code scope):
- Friction #2 (`state-checkout-missing` on bundle clone): already intentionally fail-closed with remediation `"run spacedock state init (or state new for a newly declared workflow)"` at `internal/status/roots.go:154-155`.
- Friction #6 (`gpt-6.2-*` `thinkingOptionIds: ["none"]`): upstream Paseo / provider gateway metadata bug.
- Friction #7 (`host_unknown` VM → Mac send): Paseo host registration / deployment topology, addressed in `_roadmap/paseo-backend/workspace-role-topology.md`.

## Proposed approach

1. **Dev build SHA (`internal/cli/dev_version.go` + tests):** When `Version == "dev"`, read Go VCS buildinfo (`debug.ReadBuildInfo` via an injectable seam: 7-char `vcs.revision` + `.dirty` if `vcs.modified == "true"`) and append inside semver build metadata (`0.29.0-pre0+dev.<sha7>` / `.dirty`), falling back to `+dev` when VCS buildinfo is absent. Preserves line-1 `^spacedock \S+$`, 3-line/6-line output, minor-version gate, and linker-stamped releases.
2. **Immutable collision-free dispatch files (`internal/dispatch/build.go`):** Append a deterministic 8-hex-char SHA-256 digest of `dispatchBody` (`-<hash8>`) to `dispatchFileName` (respecting `dispatchFileNameMaxLen` = 251) across ALL hosts (Claude, Codex, Pi, Paseo) so distinct concurrent briefs never collide while identical retries stay deterministic/idempotent.
3. **Configurable dispatch directory (`internal/dispatch/build.go`, `dispatch.go`):** Add `--dispatch-dir <dir>` flag + `SPACEDOCK_DISPATCH_DIR` env var fallback (`--dispatch-dir` > `SPACEDOCK_DISPATCH_DIR` > `/tmp/spacedock-dispatch`), resolved cleanly to an absolute path, created via `os.MkdirAll`, and emitted verbatim in `dispatch_file_path` and `prompt`.
4. **Docs/contract polish (`skills/first-officer/references/paseo-fo-dispatch.md`, `docs/site/reference/command-reference.md` if applicable):** Add concise notes for `--dispatch-dir`/`SPACEDOCK_DISPATCH_DIR`, `--isolation worktree` vs `--isolation local`, and `inspect.Thinking` verification.
