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

**2. `Sandbox:` reports the inverse of the truth in the sandboxed case.** It renders `safehouse.State(safehouse.Present(dir), available)`, where `selected` is "a `.safehouse` profile exists here" and `available` is "the `safehouse` binary resolves on PATH". Both are questions about *launching a command into* a sandbox. Neither asks whether this process is *already inside* one.

That is not merely a different question — in the common case the two are anti-correlated. Reproduced 2026-07-27 in a live session:

- `APP_SANDBOX_CONTAINER_ID=agent-safehouse` — the session IS executing inside a safehouse sandbox.
- `.safehouse` profile present in the workdir, so `selected == true`.
- `command -v safehouse` finds nothing, so `available == false` — precisely BECAUSE the wrap already happened; you do not invoke safehouse from within safehouse.
- `State(true, false)` renders `unavailable (safehouse not on PATH)`, and `state.go`'s own comment pins that precedence: "This dominates even when selected."

So the state in which the answer to "am I sandboxed" is most emphatically *yes* is the state the line renders as *unavailable*. The doc comment's reasoning — "nothing can wrap the launch, so a present profile cannot take effect" — is correct about a launch and exactly backwards about the running session.

Meanwhile `APP_SANDBOX_CONTAINER_ID` is a live, present signal that nothing reads.

**Scope is wider than `--version`.** `state.go`'s comment states that the launcher banner, `status --boot`, and `--version` all source these strings "so the posture reads identically across surfaces" — so all three carry it. A boot record captured this session contained `"sandbox":"unavailable (safehouse not on PATH)"` while running inside the sandbox, which means the First Officer's own durable boot evidence is wrong on this field.

This compounds: `--version` is the first command every FO session runs, at the version gate, before discovery or boot. Every session pays the runtime-line noise, and every session reads — and records — a sandbox posture that inverts under the very condition it most needs to report.

Worth checking during ideation, not asserted here: `Available()` returns an `installHint` advising `brew install eugene1g/safehouse/agent-safehouse`. `printVersion` discards it, so it does not fire there — but if any surface does emit it, inside a sandbox it advises installing the thing the session is already running in.

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
