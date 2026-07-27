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
gates:
    version: 1
    current:
        gate: gate:docs-dev:cn:backlog
    records:
        - id: gate:docs-dev:cn:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:cn-backlog-1
              briefing:
                id: briefing:docs-dev:cn:backlog:attempt-1:revision-1
                digest: sha256:7adb8b917e29b7e52dacb9e330ae55d3f76edd72a05619f9d72821e0a1c5a6aa
                digest-domain: canonical-bytes
                room-ref: ./version-output-runtime-and-sandbox-state/review/backlog/briefing-1
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

## Proposed output shape (captain-approved, 2026-07-27)

The organising rule: **inside a session, report the session; outside one, report the version.** Anything about what is *installed* belongs to `doctor`, which is already the gate's own named remedy surface.

Outside any runtime — one line:

```
spacedock 0.26.0+dev
```

Inside a session:

```
spacedock 0.26.0+dev
Runtime: claude (CLAUDECODE)
Sandbox: inside (agent-safehouse)
contract 3
```

Ambiguous markers — reports the ambiguity, does not guess, stays exit 0:

```
spacedock 0.26.0+dev
Runtime: ambiguous (CLAUDECODE, CODEX_THREAD_ID) — pass --host
Sandbox: inside (agent-safehouse)
contract 3
```

Outside variants of the sandbox line, which must carry current state and launch capability without conflating them again:

```
Sandbox: not wrapped — a launch from here would wrap (safehouse, .safehouse profile)
Sandbox: not wrapped — a launch from here cannot wrap (safehouse not on PATH; .safehouse profile present)
Sandbox: not wrapped — no .safehouse profile
```

The second is the only place the `brew install` hint belongs, and it can now never fire at a session already running inside the sandbox.

The per-host installed-plugin lines are dropped entirely. The version gate runs in one direction — the skills parse the binary's line 1 — so the binary echoing the skills' version back has no programmatic consumer. That is an install fact; `doctor` owns it.

## The two contract tombstones

`(contract 3)` and `requires-contract: ">=3,<4"` are mirrored D4 sentinels. Neither is read by any current reader: the binary carries no contract-integer mechanism, and both hosts ignore unknown manifest fields. Each exists only so the OLD side of a version mismatch fails correctly. The captain's decision splits them, and the split is technical rather than aesthetic:

**Keep `(contract 3)`, moved off line 1.** It guards old skills against a new binary. New skills cannot protect old skills — the old prose is the only thing checking, and it checks for a contract integer. Nothing else reaches that reader. Note the gate is prose executed by a model, so an absent token is not a reliable abort: a deterministic parser errors on a missing field, but a model may reason "no token, nothing to check, proceed" — a silent false-green. Emitting the literal `3` keeps the abort correct in the old prose's own terms. It moves below line 1 (the old prose says "run `--version` and parse `contract <N>`" and never pins it to line 1) and prints only inside a session, since every integer-era reader is itself a session.

**Remove `requires-contract` from both manifests.** It guards a new plugin against an old binary — a direction the skills already cover, because a new plugin's own shared core carries the minor gate and an old binary meets the FO's abort at boot regardless. It is belt-and-braces for a case the prose already catches, so removing it loses no coverage.

It is also ours, not a host contract: introduced by `080ec3ef`, undocumented in `docs/`, and placed outside Claude Code's schema by the guide check recorded in the 0199 debrief. Delete the key from both manifests. Do not rename or namespace it — a renamed key nothing reads is the same dead weight under a safer name, and the field's only readers are pre-#468 binaries a rename would break anyway.

**Write a retirement condition into the surviving constant.** `frozenContractToken` currently reads "Frozen; pinned by the internal/contractlint sync test. Do not edit." with no exit, which is how a migration aid becomes permanent furniture. Both reader populations are ours — old spacedock skills and old spacedock binaries, things we ship and tag — so "removable once no plugin or binary predating #468 can still be running" is a query against brew and the marketplace, not a guess.

## No skill change is required

The FO shared core's gate (`skills/first-officer/references/first-officer-shared-core.md:9`) is a **binary version gate** parsing `spacedock <version>` against minor `0.26`. It does not check a contract integer; #468 replaced that with minor-version coupling. Nothing in `skills/` needs editing for this task.

## What must not be disturbed

`plugin.json`'s `version` field is the live compatibility declaration that replaced the contract integer. It is bound in two directions — `internal/contractlint/prose_manifest_minor_sync_test.go` pins it against the FO prose's stamped minor, and `release.yml`'s `manifest-tag-gate` binds tag to both manifests to prose on every release. The two manifests differ only by Codex's `"hooks": "./hooks.json"`, which is correct. Removing `requires-contract` must leave all of that untouched.

## Constraints the fix must respect

- **Line 1 stays parseable.** The FO version gate parses line 1 as `spacedock <version>` and aborts on anything else. Dropping the trailing `(contract 3)` from line 1 is the point of this change; the version token itself must remain first and unchanged in shape.
- **`--version` must not fail on ambiguous runtime markers.** The dispatch detector refuses when two markers are set, which is correct there. Refusing here would break the version gate and therefore every boot — including the nested-runtime case that already occurs in practice. Ambiguity needs a reported state, not an error.
- **`Runtime: none detected` is a normal state**, not a fault — a human at a terminal is outside every runtime.
- The "inside" signal is sandbox-specific: `APP_SANDBOX_CONTAINER_ID=agent-safehouse` is what is observably set. A second sandbox implementation needs its own signal, exactly like the per-host runtime markers, so this is a small registry rather than a boolean and should be built as one.

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
