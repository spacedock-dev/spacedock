---
id: 5pb3h1ewewvrjayhgyabx5y4
title: spacedock codex can't detect an installed codex plugin — codex plugin list parse drift
status: backlog
source: captain (2026-06-01) — reproduced LIVE: `spacedock codex` reports no plugin while `codex plugin list` shows spacedock@spacedock installed, enabled 0.19.2
started:
completed:
verdict:
score: "0.40"
worktree:
issue:
---

`spacedock codex` refuses to launch — "no installed codex plugin found. Run `spacedock install --host codex`" — even when the codex plugin IS installed and enabled. Reproduced live against the just-released 0.19.2: `codex plugin list` shows `spacedock@spacedock  installed, enabled  0.19.2   https://github.com/spacedock-dev/spacedock.git, ref next`, but the codex front door can't see it. This is a launcher-replacement regression (Goal #1) shipped in 0.19.2 — it warrants a fast 0.19.3 patch.

## Root cause (verified in code)
`codexEntryInstalled` (`internal/cli/host_exec.go:96`) greps for the literal PAREN form:
```go
if strings.Contains(line, id+" (installed") { return true }
```
i.e. it expects `spacedock@spacedock (installed`. But the installed codex renders the COMMA form with NO parens: `spacedock@spacedock  installed, enabled  0.19.2  …`. So the `Contains` never matches → `resolveCodexManifest` returns "" (not installed) → `spacedock codex` (and `spacedock doctor --host codex`) report no-plugin-found. The function's own doc-comment (`:93-95`) asserts the format is `<id> (installed[, enabled]) | (not installed)` — that assumption is wrong for this codex version (codex dropped the parens). It is UNTESTED: there is no unit test for `codexEntryInstalled` (the sole codex test is a single-version happy-path integration test), which is why the format drift shipped. This is the codex plugin-list drift flagged as untracked in the session-03 handoff, now confirmed live.

## Fix direction (ideation hardens)
- Make `codexEntryInstalled` parse-tolerant instead of paren-literal: split the line into fields, find the field equal to `id`, inspect the next status token, strip a leading `(` and trailing `,`/`)`, and match `installed` EXACTLY — so `installed,` ✓, `(installed,` ✓ (old form, if codex reverts), and `not installed` (next token `not`) ✗. Do not match the bare substring `installed` (it is contained in `not installed`).
- Add UNIT tests over the real formats: `<id>  installed, enabled  <ver>  <url>` (current), `<id> (installed, enabled)` (old paren form), `<id>  not installed`, and a foreign-id line that must NOT match.
- Verify the END-TO-END `spacedock codex` resolve path against the installed codex, not just the predicate — after `codexEntryInstalled` passes, `resolveCodexManifest` still resolves a cached manifest dir (`latestVersionDir`/`codexHome`); confirm the front door actually launches.

## Notes / related
- Related latent bug (flagged, low-stakes, currently unreachable): `latestVersionDir` (`host_exec.go ~:121`) picks the LEXICALLY greatest cache dir name, not semver-greatest — bites only once a stale 9→10 rollover dir exists. Consider folding a semver-aware compare while in this file, or leave as a separate item.
- This is `internal/cli` (the codex resolver) — NOT the serialized `internal/status` lane, so it does not collide with status work; but it IS the host front door (high-stakes launch path) → detached adversarial audit at validation.
- Ships via a normal PR onto `next`; cut 0.19.3 after it lands (the codex front door is broken in the shipped 0.19.2).
