---
id: cnh3nk0yfhy5er186dm1g2h0
title: "`--version` prints every host's runtime line and reports safehouse availability instead of the session's actual runtime and sandbox state"
status: backlog
source: "Captain observation 2026-07-27 on bootstrap output; confirmed in printVersion (internal/cli/cli.go:752-759). The binary already carries the runtime detector it does not use."
started:
completed:
verdict:
score: 0.4
worktree:
issue:
---

Make the first command every First Officer session runs describe the session it is actually in: the detected runtime, and whether this session is sandboxed.

## Problem

`printVersion` (`internal/cli/cli.go:752-759`) emits, unconditionally:

```
spacedock <version> (contract 3)
Sandbox: <safehouse state>
claude: ...
codex: ...
pi: ...
```

Two separate defects sit in those seven lines.

**1. Every host's runtime line is printed regardless of the running runtime.** The loop is a hardcoded `[]string{"claude", "codex", "pi"}` with no detection, so a session running under Claude reads two lines about runtimes it is not in. The binary already owns the detector this needs — `internal/dispatch/build.go:254-281` resolves the host from the `CLAUDECODE` / `CODEX_THREAD_ID` / `PI_CODING_AGENT_DIR` markers and refuses on ambiguity. `printVersion` simply does not call it.

**2. `Sandbox:` answers a different question than it appears to.** It renders `safehouse.State(safehouse.Present(dir), available)` — whether the safehouse binary is on PATH and whether a sandbox is configured for this directory. It does not report whether the current session is executing inside a sandbox. "Sandbox: unavailable (safehouse not on PATH)" is read as "you are not sandboxed", which is not what was measured.

This is small but it compounds: `--version` is the first command every FO session runs, at the version gate, before discovery or boot. Every session pays the noise and every session reads a sandbox line that may not mean what it says.

## Constraints the fix must respect

- **Line 1 is a parsed contract.** The FO version gate parses line 1 as `spacedock <version>` and aborts on anything else. Any change must leave line 1 byte-compatible; only the lines below it may move.
- **`--version` must not fail on ambiguous runtime markers.** The dispatch detector refuses when two markers are set, which is correct there. Refusing here would break the version gate and therefore every boot — including the nested-runtime case that already occurs in practice. Ambiguity needs a reported state, not an error.
- Cross-host probing must remain reachable somehow; some diagnostics legitimately want all three. Whether that becomes a flag or moves to `doctor` is a design question, not a foregone conclusion.

## Also worth checking

`internal/cli/frontdoor.go` and `internal/cli/launch_banner_sandbox_test.go` carry the same `Sandbox:` rendering on the launch banner. Whatever the sandbox line comes to mean should mean the same thing in both places.

## Out of scope

- Changing what the version gate requires of line 1.
- Building a new runtime-detection mechanism; the existing detector is the one to reuse.
- Sandbox enforcement behaviour — this is about reporting only.

## Acceptance criteria

Ideation fills these in. The end state is that `--version` names the runtime the session is actually in and reports this session's sandbox state, with line 1 unchanged, and does not fail when runtime markers are ambiguous.

## Test plan

Ideation fills this in. Golden-output fixtures per runtime marker set, including the ambiguous case, are the likely substrate; existing version golden tests are the starting point.
