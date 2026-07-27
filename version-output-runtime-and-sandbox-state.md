---
id: cnh3nk0yfhy5er186dm1g2h0
title: "`--version` prints every host's runtime line and reports safehouse availability instead of the session's actual runtime and sandbox state"
status: ideation
source: "Captain observation 2026-07-27 on bootstrap output; confirmed in printVersion (internal/cli/cli.go:752-759). The binary already carries the runtime detector it does not use."
started: 2026-07-27T08:05:18Z
completed:
verdict:
score: 0.4
worktree:
issue:
gates:
    version: 1
    current:
        gate: gate:docs-dev:cn:ideation
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
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:cn:backlog:1
                briefing: briefing:docs-dev:cn:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-27T08:04:47.847614Z"
                decision: approve
                reason: 'Captain directive in the active conversation: ''dispatch cn''. Direction accepted on evidence: the sandbox posture was reproduced inverting live (APP_SANDBOX_CONTAINER_ID=agent-safehouse set, .safehouse profile present, safehouse absent from PATH, so State(true,false) renders ''unavailable'' with state.go documenting that precedence as dominating), the same strings feed the launcher banner and status --boot so First Officer boot evidence carries the false posture, and the runtime detector the fix needs already exists at internal/dispatch/build.go:254-281 and is simply not called. The output shape, the tombstone split, and the requires-contract deletion are captain-approved and recorded in the entity body.'
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
        - id: gate:docs-dev:cn:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:cn-ideation-1
              briefing:
                id: briefing:docs-dev:cn:ideation:attempt-1:revision-2
                digest: sha256:784787fef86b6b4caf18557c61dddc839dfe37209bab72d0d725cdd90763e9ac
                digest-domain: canonical-bytes
                room-ref: ./version-output-runtime-and-sandbox-state/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:cn:ideation:1
                briefing: briefing:docs-dev:cn:ideation:attempt-1:revision-2
                by: person:captain
                at: "2026-07-27T08:53:23.476676Z"
                decision: revise
                reason: 'Captain resolved the float in the review TUI (resolution:captain-1785142370707298000) with two annotations, direction accepted. 1: the Runtime line should also carry the session identifier, not only the host and the marker that proved it. 2: the outside-a-session sandbox strings are wrong in shape — outside a session it is never wrapped, so ''not wrapped'' is a constant carrying no information; report availability instead. Design direction and the rest of the six acceptance criteria stand.'
              application:
                action: feedback
                target-stage: ideation
                state: superseded
            - id: gate-attempt:cn-ideation-2
              briefing:
                id: briefing:docs-dev:cn:ideation:attempt-2:revision-1
                digest: sha256:8d8348fe9e1f1083c56532dd3d4d96286bd649d6c36c02205a48297276b02650
                digest-domain: canonical-bytes
                room-ref: ./version-output-runtime-and-sandbox-state/review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:cn:ideation:2
                briefing: briefing:docs-dev:cn:ideation:attempt-2:revision-1
                by: person:captain
                at: "2026-07-27T09:06:27.862968Z"
                decision: approve
                reason: 'Captain resolved the float in the review TUI (resolution:captain-1785143163921397000) with approve and no annotations, against briefing attempt-2 revision-1 whose artifact is the frozen room copy at sha256:74628c92. Both prior annotations are folded: the Runtime line carries a session identifier resolved from a per-host identity field rather than the detection marker (checking caught that CLAUDECODE is the literal 1, so the assumed form would have rendered ''session 1''), and the outside sandbox strings reduce to two session-state renders after the launch question moved to the surface that asks it. Six acceptance criteria stand, two value-measuring, each with a falsifier; three mechanisms were exercised rather than assumed, including performing the requires-contract deletion against both binding gates and reverting clean; expected surface is net negative at roughly -60 lines.'
              application:
                action: advance
                target-stage: implementation
                state: pending
                blockers: []
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
Runtime: claude (CLAUDECODE, session afd74765)
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

Not sandboxed — the sandbox line reports availability, nothing more:

```
Sandbox: not sandboxed (safehouse available)
Sandbox: not sandboxed (safehouse not installed)
```

### Captain revision, gate cycle 1

Two annotations at the ideation gate changed the shape above; both are folded in and both narrow it.

**The Runtime line carries the session identifier.** `Runtime: claude (CLAUDECODE, session afd74765)` — see "Session identity" below for how it resolves and why eight characters.

**The three outside sandbox strings are replaced by two.** The captain: *"this makes no sense. outside it is never wrapped. we should just say it is available."* He is right, and the old strings were wrong twice over. `not wrapped` is constant across all three, so it carried no information while reading as though it did. Worse, all three answered *"would a launch from here be wrapped"* — a launch question — in a line whose whole purpose after this task is to report the running session. That is the same category error as the original defect, surviving into the fix. The `.safehouse` profile is dropped from this line entirely: a profile determines whether a *launch* would wrap, and a session that is already running is not about to launch itself.

The launch question does not disappear; it moves to the surface that actually asks it. See "Two questions, two renderers" below.

The `brew install` hint therefore belongs to the launch path alone, and can never fire at a session already running inside the sandbox.

## Session identity

The Runtime line names the host, the marker that proved it, and the session's own identifier. The third of those needed checking rather than assuming, and the assumption would have been wrong.

**`CLAUDECODE` is `1`.** It is a boolean flag, not an identifier. Claude's session identity lives in a *separate* variable, `CLAUDE_CODE_SESSION_ID` — confirmed live at `afd74765-9000-4e63-acf4-3b1f4645a8f3`, and already read by this repo for the merged-team lead-session reconcile (`internal/claudeteam/reconcile.go:50`). So detection marker and identity marker are two different columns, not one. `Detect` gains an identity field per host rather than reusing the marker value.

| Host | Detection marker | Identity variable |
| --- | --- | --- |
| codex | `CODEX_THREAD_ID` | `CODEX_THREAD_ID` — marker and identity coincide here, and only here |
| claude | `CLAUDECODE` | `CLAUDE_CODE_SESSION_ID` |
| pi | `PI_CODING_AGENT`, `PI_CODING_AGENT_DIR` | none |

**Pi has no session identifier, and this is a finding rather than a decision.** `PI_CODING_AGENT_DIR` is the agent installation directory. The one variable that looks like a session is `PI_CODING_AGENT_SESSION_DIR`, but `internal/cli/pi.go:524-528` shows it resolving to `~/.pi/agent/sessions` — the *collection* of all sessions, not one session. Nothing in pi's environment distinguishes one pi session from another. So pi renders no identifier, and the rule that produces that is the general one, not a pi special case: **a host with no identity variable, and a host whose identity variable is unset, take the same path — the segment is omitted.** That answers the "must not print an empty parenthetical" requirement without a branch dedicated to it, and it means the day pi gains a session variable, the change is one table cell.

**Eight characters, and the number is not arbitrary.** These identifiers are UUID-shaped and would otherwise dominate a four-line output. The prefix length matches a convention already on disk: Claude Code names its own team directory `~/.claude/teams/session-<first 8 chars>`, verified live this session — `CLAUDE_CODE_SESSION_ID=afd74765-…` against the existing directory `session-afd74765`. Rendering the same eight characters makes the printed token directly greppable against the session state already on the filesystem, which is worth more than an arbitrary truncation of the same length. An identifier of eight characters or fewer renders whole rather than padded.

Eight hex characters distinguish concurrent sessions on one host with room to spare — the point of showing it at all. A human runs a handful of sessions at once, against a 4.3-billion-value space.

**Ambiguous markers render no identifier.** Printing one would assert that a host was resolved, which is exactly what the ambiguous branch declines to do.

## Two questions, two renderers

The original body observed that `--version`, the launch banner, and `status --boot` "all source the same three-way strings so the posture reads identically across surfaces", and treated that sharing as the thing to preserve. Following the captain's second annotation to its root, the sharing is itself part of the bug: these surfaces ask two different questions, and one string set cannot answer both without the conflation this task exists to remove.

**`SessionState(insideName, inside, available)` — "is this process sandboxed?"** Used by `--version` and `status --boot`. Three strings: `inside (agent-safehouse)`, `not sandboxed (safehouse available)`, `not sandboxed (safehouse not installed)`. No `.safehouse` mention.

**`LaunchState(insideName, inside, selected, available)` — "will this launch be wrapped?"** Used by the launch banner only, where the question is genuinely about a launch that has not happened yet, and where the `.safehouse` profile is the load-bearing input:

```
Sandbox: inside (agent-safehouse) — launching without re-wrapping
Sandbox: wrapping this launch (safehouse, .safehouse profile)
Sandbox: not wrapping this launch (no .safehouse profile)
Sandbox: not wrapped (safehouse not installed; .safehouse profile present)
```

The banner still needs the `inside` arm: launched from within the sandbox today it renders `unavailable (safehouse not on PATH)`, the same inversion. So the three-surface REACH decision stands unchanged — all three stop lying. What changes is that two of them now answer the question they were being asked.

*Simplest alternative considered:* keep one renderer and let the banner print the session-state strings. Rejected because the banner would then stop reporting whether the launch it is about to perform will be sandboxed, which is the only reason that line is on a pre-launch banner. Two ten-line functions are cheaper than deleting a fact the operator uses.

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

**Resolved in ideation, and NOT the way this paragraph anticipated.** Following the captain's cycle-1 annotation to its root showed the shared rendering to be part of the defect: the banner asks whether a launch will be wrapped, `--version` and `--boot` ask whether this process is sandboxed, and one string set cannot answer both honestly. The banner keeps launch semantics under `LaunchState` and gains the `inside` arm so it stops inverting; `--version` and `--boot` take `SessionState`. What must mean the same thing in all three is the *detection* — `Inside()` is the single registry — not the rendered sentence. See "Two questions, two renderers".

## Out of scope

- Changing what the version gate requires of line 1.
- Building a new runtime-detection mechanism; the existing detector is the one to reuse.
- Sandbox enforcement behaviour — this is about reporting only.

## Spike record (2026-07-27)

Two mechanisms the design rests on were unverified, so both were exercised before the design was written.

**1. The live inversion, reproduced in this ideation session.** `APP_SANDBOX_CONTAINER_ID=agent-safehouse`, `CLAUDECODE=1`, `.safehouse` present in the workdir, `command -v safehouse` empty. `go run ./cmd/spacedock --version` in that process emits:

```
spacedock 0.26.0+dev (contract 3)
Sandbox: unavailable (safehouse not on PATH)
claude: spacedock 0.27.0-pre1
codex: spacedock 0.26.0
pi: spacedock ready
```

and `status --boot` emits `SANDBOX: unavailable (safehouse not on PATH)`. Both wrong, from inside the sandbox. This is the baseline AC-1 measures against.

**2. The detector extraction and the new renderers, exercised as a throwaway program.** A standalone Go program implementing the proposed `Detect` (env-only, never errors), the proposed sandbox registry, and the proposed render functions was run against the live env and against five pinned marker sets. It produced the captain-approved shape byte-for-byte in every case, exit 0 throughout:

- live env → `Runtime: claude (CLAUDECODE)` / `Sandbox: inside (agent-safehouse)` / `contract 3`
- no markers → the single line `spacedock 0.26.0+dev`, nothing else
- `CODEX_THREAD_ID` + `CLAUDECODE` → `Runtime: ambiguous (CODEX_THREAD_ID, CLAUDECODE) — pass --host`, still exit 0
- `PI_CODING_AGENT` + `PI_CODING_AGENT_DIR` → `Runtime: pi (PI_CODING_AGENT, PI_CODING_AGENT_DIR)` — two markers for the SAME host is not ambiguity, which the naive "more than one marker set" rule would have got wrong
- the three outside-sandbox variants each rendered as approved (since replaced — see spike 4)

The riskiest path — an ambiguity rule that stays exit 0 while still distinguishing real ambiguity from the two-pi-markers case — is therefore proven, not assumed. The throwaway's marker table and its five cases seed the implementation's first test.

**3. The `installHint` question the body left open, answered: it is real, and it is a launch bug, not a reporting bug.** `Available()`'s hint is discarded by `printVersion` (confirmed — it does not fire there), but it IS emitted at three launch sites: `frontdoor.go:418`, `frontdoor.go:631`, `pi.go:301`. Each is gated on `wrap && !Available`. Inside this sandbox, `.safehouse` is present so `wrap` is true and `safehouse` is off PATH, so a nested `spacedock claude` prints "install safehouse (brew install …)" and exits 1 — advising installation of the thing the session is already running in, and refusing the launch. This is the launch decision, which the body puts out of scope ("sandbox enforcement behaviour — this is about reporting only"), so this task does not change it. Recorded as a follow-up: **the launch wrap gate should consult the same `Inside()` registry this task builds, and either skip the wrap or refuse with an accurate nested-sandbox message.** Filing it as its own task is the recommendation; the registry this task ships is the prerequisite it was blocked on.

### Spike 4 — the captain's two revisions (gate cycle 1)

Both revisions were exercised as a second throwaway program before the design above was written, and the first of them overturned an assumption that would otherwise have shipped.

**The identity assumption was wrong, and cheap to catch.** The gate note said the per-host session identity was already in hand from `Detect`. It is not: `CLAUDECODE` is `1`, a flag. Reading the live environment found the identity in a *separate* variable, `CLAUDE_CODE_SESSION_ID=afd74765-9000-4e63-acf4-3b1f4645a8f3`. Had the identity been taken from the detection marker, Claude sessions would have rendered `session 1` — every session identical, the exact failure the annotation asked to avoid. AC-2a now pins this with a row asserting the identity is not `"1"`.

**The eight-character prefix was verified against on-disk state, not chosen.** `ls ~/.claude/teams/session-afd74765` succeeded while `CLAUDE_CODE_SESSION_ID` was `afd74765-…`, confirming the host's own directory-naming convention uses that prefix.

**Pi's identity was checked rather than assumed absent.** `PI_CODING_AGENT_SESSION_DIR` looked like the answer; `internal/cli/pi.go:524-528` shows it resolving to `~/.pi/agent/sessions`, the collection directory. Pi exposes no per-session identifier.

The program rendered every revised case, exit 0 throughout:

- live env → `Runtime: claude (CLAUDECODE, session afd74765)` / `Sandbox: inside (agent-safehouse)`
- codex → `Runtime: codex (CODEX_THREAD_ID, session 01937f2a)` — marker and identity coincide
- pi, and claude with `CLAUDE_CODE_SESSION_ID` unset → `Runtime: pi (PI_CODING_AGENT, PI_CODING_AGENT_DIR)` and `Runtime: claude (CLAUDECODE)`, no empty parenthetical in either
- an identifier shorter than eight characters rendered whole, not padded
- two `CLAUDE_CODE_SESSION_ID` values rendered two distinguishable lines
- ambiguous markers rendered no identifier
- the revised sandbox strings, and the four launch-banner strings, each as designed

## Proposed approach

**A new `internal/runtimehost` package owns the env-only detection.** It exports `Detect(getenv func(string) string) (host string, markers []string, identity string, ambiguous bool)` over the table lifted from `internal/dispatch/build.go:250-281`, extended with the identity column from the Session identity section above. It never returns an error: the two callers need opposite dispositions of the same facts — dispatch must refuse on ambiguity, `--version` must report it — so the shared unit is the detection, and the policy stays with each caller. It also exports `ShortID(id string) string`, the eight-character truncation, so the prefix rule has one definition rather than one per caller.

`identity` is the raw value; truncation happens at render. Keeping the full value in the returned struct means a later consumer that needs to match a session id exactly — the reconcile path already does — is not forced to un-truncate.

`resolveBuildHost` is rewritten to call `Detect` and keep its own error strings verbatim (`ambiguous runtime host sources: … (%s); pass --host …` from the returned marker list, and the missing-host message unchanged). This is the body's "reuse the existing detector" made literal: one marker table, two policies.

*Simplest alternative considered:* copy the four-line marker table into `printVersion` and leave `build.go` alone. Rejected because a duplicated table is exactly how the two surfaces drift — a future host added to one and not the other reports a runtime the dispatcher cannot resolve, silently. The extraction is ~50 lines and removes the drift class outright. It serves AC-2.

**The sandbox registry lives beside the launch seam, in `internal/safehouse/state.go`.** `Inside(getenv) (name string, ok bool)` scans a table of `{env, wantValue, displayName}` entries, today the single `{APP_SANDBOX_CONTAINER_ID, agent-safehouse, agent-safehouse}`. It matches on the VALUE, not mere presence, because `APP_SANDBOX_CONTAINER_ID` is a generic macOS app-sandbox variable that other containers also set; matching presence would claim safehouse for any of them. A second sandbox implementation adds one table row.

`State(selected, available bool)` is deleted and replaced by the two renderers described in "Two questions, two renderers" above: `SessionState(insideName string, inside, available bool)` and `LaunchState(insideName string, inside, selected, available bool)`. The old three strings go with it; nothing keeps them, since the whole defect is that they answer a launch question in a session-reporting slot.

*Simplest alternative considered:* a bare `Inside() bool` plus a hardcoded `"agent-safehouse"` in the render string. Rejected on the body's own constraint — the display name must come from the same row as the signal, or the second sandbox reports itself as safehouse. One struct field is cheaper than that bug. It serves AC-1.

**REACH: all three surfaces, deliberately.** `printVersion` (`cli.go:755`) and `status --boot` (`boot.go:256`) take `SessionState`; `launchBanner` (`frontdoor.go:179`) takes `LaunchState`. Fixing only `--version` would leave every boot record the First Officer writes carrying `"sandbox":"unavailable (safehouse not on PATH)"` from inside the sandbox, which is the sharpest cost named in the problem statement: durable, machine-read evidence that is wrong. `boot.go` already has an env seam (`e.get`) and `newRootCommand` already threads `env []string`, so no new plumbing is needed at any of the three.

The `Runtime:` line is `--version`-only. The banner already names the host it is launching, and `--boot` is not the version surface; adding it to either would be a line nothing asked for.

**`printVersion` is rewritten around the organising rule.** Signature becomes `printVersion(w io.Writer, dir string, getenv func(string) string, lookPath func(string) (string, error))` — the `runtimeProbe` parameter is gone. It prints `spacedock <version>` unconditionally; then, only when `Detect` reports a host or ambiguity, the `Runtime:` line, the `Sandbox:` line, and a bare `contract 3`. Outside every runtime it prints one line and stops.

**The per-host block is deleted, and with it `internal/cli/host_runtime.go` entirely** (207 lines: `runtimeStatus`, `runtimeProbe`, `runtimeLine`, `enabledMarker`, `claudeMarker`, `codexEntryDisabled`, `execRuntimeProbe`). Verified: nothing outside that file and its two test files uses any of it. `pluginListEntry` and `codexEntryInstalled` live in `host_exec.go` and are used elsewhere — they stay.

*Honest cost, recorded rather than glossed:* the body says "that is an install fact; `doctor` owns it." `doctor` is single-host — it resolves ONE host's manifest and returns a verdict (`internal/contract/doctor.go:57-90`); it does not enumerate claude/codex/pi. So the multi-host inventory disappears with no replacement. The recommendation is to accept that: the organising rule is "inside a session, report the session", and in a session only your own host's plugin can matter — which is exactly what `doctor` already reports. If the captain wants the inventory back, `spacedock doctor --all-hosts` is the place for it, as a separate task; it is not `--version`'s job and it is not in this task's approved shape.

**`(contract 3)` moves.** `frozenContractToken` changes from `"(contract 3)"` to `"contract 3"` (the parens were line-1 punctuation) and prints on its own line below the sandbox line, inside a session only. The value stays 3. Its doc comment gains the retirement condition the body asks for, replacing "Frozen … Do not edit." with the removable-once condition: no plugin or binary predating #468 still in circulation — a query against the Homebrew formula and the marketplace, not a guess.

**`requires-contract` is deleted from `.claude-plugin/plugin.json:19` and `.codex-plugin/plugin.json`.** `TestVendoredManifestTombstoneFrozen` (`prose_manifest_minor_sync_test.go:98-108`) asserts the field's presence and value, so it is deleted with the field, along with the `frozenTombstone` constant and the `RequiresContract` struct field. `TestProseMinorMatchesVendoredManifestMinor` reads only `version` and is untouched. `release.yml`'s `manifest-tag-gate` (line 166) binds tag→`version`→prose and never reads `requires-contract`. The `requires-contract` strings in `internal/release/release_test.go` and the `internal/cli/*_test.go` fixtures are synthetic HOST manifests exercising round-trip preservation and old-shape install paths — they are not the vendored manifests and stay as they are.

## Acceptance criteria

**AC-1 (value). A process running inside the safehouse sandbox is reported as inside, on all three surfaces.** In an environment with `APP_SANDBOX_CONTAINER_ID=agent-safehouse`, a `.safehouse` profile present, and `safehouse` absent from PATH — the exact configuration this ideation session ran in — `--version`, the launch banner, and `status --boot --json` each report `inside (agent-safehouse)`. The baseline moves the wrong way and is independently observable: today all three report `unavailable (safehouse not on PATH)` under identical inputs, and `status --boot --json` writes that string into `sandbox`, so a durable boot record is the artifact that changes. *Tested by:* three table tests — `internal/cli/version_session_test.go`, `internal/status/boot_sandbox_test.go` (both over `SessionState`'s three `(inside, available)` outcomes) and `internal/cli/launch_banner_sandbox_test.go` (over `LaunchState`'s four `(inside, selected, available)` outcomes) — each naming the inside-and-safehouse-absent row as the regression case, and each asserting the rendered string against a test-supplied literal rather than the production constant. All three fail today.

**AC-2 (value). `--version` names the runtime the session is in — and which session that is — and names no other runtime.** Under a single host marker the output carries exactly one `Runtime:` line, naming that host, the marker that proved it, and the session's own identifier where the host exposes one; and it carries zero lines about any other runtime. Two concurrent sessions on one host render distinguishable identifiers, which is the point of showing one. The measurable baseline: today's output contains three host lines regardless of markers and no session identity at all, and `--version` execs three host CLIs to produce them — measured at 0.86s wallclock against 0.31s for `status --help` in the same process and directory, so ~64% of the first command every First Officer session runs is spent reporting runtimes the session is not in. After the change no host CLI is executed. *Tested by:* a PATH-shim test that puts executable `claude`/`codex`/`pi` shims on PATH which append to a witness file, runs `--version` through `cli.Run`, and asserts the witness file was never created, plus an assertion that no output line begins `claude:`, `codex:`, or `pi:`; and a two-session case rendering `CLAUDE_CODE_SESSION_ID=afd74765-…` and `afd74799-…` and asserting the two `Runtime:` lines differ. Both fail today (the witness gets three entries; neither line carries an identifier).

**AC-2a. The session identifier renders from the right variable, truncates to a greppable prefix, and never renders empty.** `claude` reads `CLAUDE_CODE_SESSION_ID` (not `CLAUDECODE`, which is the flag `1`); `codex` reads `CODEX_THREAD_ID`; `pi`, which exposes no session variable, renders no identifier segment — as does any host whose identity variable is unset. No output ever contains an empty or dangling parenthetical. An identifier longer than eight characters renders as its first eight; one of eight or fewer renders whole. *Tested by:* `internal/runtimehost/runtimehost_test.go` asserting `identity` per host over set and unset identity variables, including the `CLAUDECODE=1` row that asserts the identity is NOT `1` — the falsifying change being an implementation that reuses the detection marker's value as the identifier; a `ShortID` table covering longer-than-8, exactly-8, shorter-than-8, and empty; and render cases asserting the pi and unset-identity outputs equal `Runtime: pi (PI_CODING_AGENT, PI_CODING_AGENT_DIR)` and `Runtime: claude (CLAUDECODE)` exactly, so a stray `, session ` or `()` fails.

**AC-2b. The sandbox line reports the session, and the launch banner reports the launch.** `--version` and `status --boot` render one of exactly three strings — `inside (agent-safehouse)`, `not sandboxed (safehouse available)`, `not sandboxed (safehouse not installed)` — and none of them mentions a `.safehouse` profile, because a profile is a launch fact and neither surface is about a launch. The launch banner renders the launch question, and its wrapping-this-launch string is the only place the `brew install` hint can be reached from. *Tested by:* the AC-1 tables, extended with an assertion that no `--version` or `--boot` sandbox string contains the substring `.safehouse`. The falsifying change: collapsing the two renderers back into one puts the profile text on `--version` and turns it red.

**AC-3. Ambiguous runtime markers are reported, not guessed at, and never fail.** With `CODEX_THREAD_ID` and `CLAUDECODE` both set — the recorded nested-Claude-under-Codex marker leak, not a hypothetical — `--version` exits 0 and emits `Runtime: ambiguous (CODEX_THREAD_ID, CLAUDECODE) — pass --host`. Two markers belonging to the SAME host (`PI_CODING_AGENT` + `PI_CODING_AGENT_DIR`) are not ambiguous and report `Runtime: pi (PI_CODING_AGENT, PI_CODING_AGENT_DIR)`. *Tested by:* `internal/runtimehost/runtimehost_test.go` over the marker-set matrix asserting `(host, markers, ambiguous)`, and a `--version` case asserting both the exact line and exit code 0. The falsifying change: an ambiguity rule of "more than one marker set" turns the two-pi-markers row red; a rule that returns an error turns the exit-code assertion red.

**AC-4. Line 1 stays exactly `spacedock <version>` and the version gate still passes.** Line 1 carries the version token first and nothing after it. `contract 3` appears below the sandbox line inside a session, and does not appear at all outside every runtime. *Tested by:* an updated `TestVersionFirstLineUnchanged` asserting line 1 matches `^spacedock \S+$` and equals `"spacedock " + displayVersion()`; `internal/cli/dev_version_test.go:77`'s expectation updated to the same shape; and a case asserting `contract 3` is present in the in-session render and absent from the outside-every-runtime render. The falsifying change: leaving the token on line 1 turns the regex red.

**AC-5. Deleting `requires-contract` leaves the live minor-version binding intact.** Neither vendored manifest contains a `requires-contract` key; both retain their `version` field unchanged at `0.26.0`. `go test ./internal/contractlint/` is green including `TestProseMinorMatchesVendoredManifestMinor`, and `go run ./cmd/spacedock-release manifest-tag-gate v0.26.0 .claude-plugin/plugin.json .codex-plugin/plugin.json skills/first-officer/references/first-officer-shared-core.md` exits 0. *Tested by:* running both directly — the sync test in the normal suite, and the manifest-tag-gate invoked exactly as `release.yml:166` invokes it, with its exit code recorded. Nothing under `skills/` is edited; a `git diff --stat -- skills/` showing no entries is the check.

**AC-6. Documentation matches the shipped output.** `docs/site/reference/command-reference.md`'s `--version` section shows the in-session and outside-a-session shapes, the session identifier, and the three session sandbox strings, with no surviving mention of the per-runtime block or of a contract level on line 1. *Tested by:* the doc diff below applied verbatim, and a manual read-back of the rendered section against real `--version` output captured in both a session and a bare shell.

## Expected surface

Roughly 20 files, about +630 / −600 lines, for a net of about +30.

**This moved, and the direction is worth naming.** Before the captain's revision this was +540 / −600, a net NEGATIVE of about −60. The session identifier adds a table column, a `ShortID` rule, and the tests that pin both; the two-renderer split adds a second function and its cases. Both are additive, so the net crosses from −60 to +30 — a swing of about +90, inside the ±250-line tolerance and leaving the file count unchanged. The tolerance below is NOT restated to absorb it: the numbers are the revised estimate, and the original net-negative claim no longer holds. The large deletions still dominate the shape of the change; they no longer quite outweigh it.

| Path | Change |
| --- | --- |
| `internal/runtimehost/runtimehost.go` | new, ~70 |
| `internal/runtimehost/runtimehost_test.go` | new, ~140 |
| `internal/dispatch/build.go` | ~+15 / −25 |
| `internal/safehouse/state.go` | ~+65 / −20 |
| `internal/safehouse/state_test.go` | ~+85 / −30 |
| `internal/cli/cli.go` | ~+25 / −12 |
| `internal/cli/host_runtime.go` | DELETED, −207 |
| `internal/cli/version_runtime_test.go` | DELETED, −143 |
| `internal/cli/version_claude_enabled_test.go` | DELETED, ~−50 |
| `internal/cli/version_session_test.go` | new, ~175 |
| `internal/cli/frontdoor.go` | ~+6 / −3 |
| `internal/cli/launch_banner_sandbox_test.go` | ~+30 / −12 |
| `internal/cli/launch_banner_wording_test.go` | ~+3 / −2 |
| `internal/cli/dev_version_test.go` | ~+4 / −3 |
| `internal/status/boot.go` | ~+5 / −2 |
| `internal/status/boot_sandbox_test.go` | ~+35 / −15 |
| `internal/contractlint/prose_manifest_minor_sync_test.go` | ~−25 |
| `.claude-plugin/plugin.json` | −1 |
| `.codex-plugin/plugin.json` | −1 |
| `docs/site/reference/command-reference.md` | ~+22 / −10 |

**Tolerance: ±6 files and ±250 lines.** The soft spot is the blast radius of replacing `State` — a grep for the three old strings found five source/test files plus the doc, and the recorded live-session `.jsonl` fixtures under `internal/ensigncycle/testdata/` which contain the strings as captured transcript content and are inputs, not oracles. If any of those turns out to be asserted against, the file count rises. Exceeding the tolerance means the surface was mis-scoped; report it rather than absorbing it.

## Test plan

All tests are in-process Go tests over injected seams. No live host CLI is executed in the test path, and no test reads the running machine's real environment.

| Test | What it proves | Cost |
| --- | --- | --- |
| `internal/runtimehost/runtimehost_test.go` | `Detect` over the marker-set matrix: none, each host alone, two markers same host (not ambiguous), two markers different hosts (ambiguous), all four set. Asserts `(host, markers, identity, ambiguous)`. | low, ~90 lines, pure function |
| `internal/runtimehost/runtimehost_test.go` identity rows | AC-2a: identity comes from `CLAUDE_CODE_SESSION_ID` and never from `CLAUDECODE` (the row asserting the identity is not `"1"`); codex's marker doubles as its identity; pi resolves none; an unset identity variable resolves empty. Plus a `ShortID` table: longer than 8, exactly 8, shorter than 8, empty. | low, ~50 lines, pure function |
| `internal/dispatch` existing build tests | `resolveBuildHost`'s error strings and refusal policy are unchanged after the extraction. Existing tests are the oracle; if none covers the ambiguity message, add one. | low, mostly re-running what exists |
| `internal/safehouse/state_test.go` | `Inside` matches on value not presence (a row asserting `APP_SANDBOX_CONTAINER_ID=something-else` is NOT inside); `SessionState` renders its three strings, none containing `.safehouse`; `LaunchState` renders its four, with the inside arm dominating in both. | low, ~85 lines |
| `internal/cli/version_session_test.go` | The full `--version` render per marker set through a fake getenv and pinned lookPath: outside (one line, no contract token), in-session with and without an identifier, ambiguous (exit 0, no identifier), and both not-sandboxed variants. Line 1 shape, `contract 3` placement, and no empty parenthetical anywhere. | medium, ~175 lines, golden strings |
| `internal/cli/version_session_test.go` PATH-shim case | No host CLI is executed by `--version` (AC-2). Shims on a temp PATH append to a witness file; the assertion is that the file does not exist. | medium — needs a temp dir, executable shims, and PATH control; the one test with real filesystem setup |
| `internal/cli/version_session_test.go` two-session case | AC-2: two `CLAUDE_CODE_SESSION_ID` values render two different `Runtime:` lines — the identifier actually distinguishes concurrent sessions rather than merely appearing. | low |
| `internal/cli/launch_banner_sandbox_test.go` | The banner's `LaunchState` line, four combinations. Existing test, expectations replaced. | low |
| `internal/status/boot_sandbox_test.go` | `status --boot` and `--boot --json`'s `sandbox` value, `SessionState`'s three outcomes. Existing test, expectations replaced. | low |
| `internal/contractlint/` suite | The minor-version binding survives the `requires-contract` deletion. | free — existing suite |
| `manifest-tag-gate` invocation | The release gate reads only `version` and still exits 0. Run once by hand, exit code recorded in the stage report. | low, one command |

No live workflow test is needed: nothing here is runtime handoff behavior. No new golden FIXTURE FILES are needed either — the strings are short enough that in-test literals are clearer than a testdata round-trip, and an in-test literal is an independent expected value in a way a regenerated golden file is not.

## Documentation diff

`docs/site/reference/command-reference.md`. Line 3, the intro sentence:

```diff
-The `spacedock` binary groups its subcommands into Launch, Setup, and Workflow, plus a top-level `spacedock --version` (the binary version and contract level). For the exact flags of any command, run `spacedock <command> --help`, the always-current source of truth; `spacedock` with no arguments prints the grouped help.
+The `spacedock` binary groups its subcommands into Launch, Setup, and Workflow, plus a top-level `spacedock --version` (the binary version, and — inside an agent session — that session's runtime and sandbox state). For the exact flags of any command, run `spacedock <command> --help`, the always-current source of truth; `spacedock` with no arguments prints the grouped help.
```

Lines 5-17, the `--version` section, replaced in full:

```markdown
## --version

`spacedock --version` reports the binary version, and — when it is running inside an agent session — that session's runtime and sandbox state. Outside any session it prints one line:

    spacedock 0.26.0

Inside a session it also names the runtime it detected, the marker that proved it, which session this is, and whether this process is running inside a sandbox:

    spacedock 0.26.0
    Runtime: claude (CLAUDECODE, session afd74765)
    Sandbox: inside (agent-safehouse)
    contract 3

The session identifier is the first eight characters of the host's own session id — the same prefix Claude Code uses to name `~/.claude/teams/session-afd74765` — so you can tell two concurrent sessions apart and match one against its state on disk. Hosts that do not expose a session id, such as pi, omit it:

    Runtime: pi (PI_CODING_AGENT, PI_CODING_AGENT_DIR)

When markers for more than one runtime are set — a nested session can leak them — it reports the ambiguity rather than guessing, and still exits 0:

    Runtime: ambiguous (CLAUDECODE, CODEX_THREAD_ID) — pass --host

`Runtime: none detected` is a normal state: it means a human at a terminal, outside every runtime.

The `Sandbox:` line answers one question — is this process sandboxed? Inside a sandbox it names it; otherwise it reports whether safehouse is available to sandbox future launches:

    Sandbox: inside (agent-safehouse)
    Sandbox: not sandboxed (safehouse available)
    Sandbox: not sandboxed (safehouse not installed)

`spacedock status --boot` reports the same three. The pre-launch banner answers the neighbouring but different question — whether the launch it is about to perform will be wrapped — so its `Sandbox:` line reads in terms of that launch.

The trailing `contract 3` is a frozen compatibility sentinel read only by skill versions predating the current version gate. It prints inside a session only. For what is installed for each host — plugin versions and enablement — use `spacedock doctor`.
```

The first line stays `spacedock <version>` because the First Officer and ensign skills parse it; that invariant is unchanged and is why the contract token moved off it.

## Follow-ups this task deliberately does not do

1. **The nested-launch `installHint` bug.** Confirmed live, described in the spike record above. Requires changing the launch wrap gate, which the body puts out of scope. Blocked on the `Inside()` registry this task ships; file it once this lands.
2. **`spacedock doctor --all-hosts`**, if the captain wants the multi-host plugin inventory that the per-host block currently provides. Not `--version`'s job under the approved organising rule.

## Stage Report: ideation

- DONE: Value AC measures the inversion actually being fixed: a session running inside the sandbox reports itself as inside, against a baseline that can move the wrong way — today that exact session reports "unavailable". Assert the reported posture against the real session state, never that a string changed.
  AC-1 asserts `inside (agent-safehouse)` on all three surfaces against the live baseline reproduced this session — `--version` and `status --boot` both printed `unavailable (safehouse not on PATH)` while `APP_SANDBOX_CONTAINER_ID=agent-safehouse` was set.
- DONE: The ambiguous-marker branch stays exit 0 and reports rather than guesses. This is the one branch that can break every boot, and its shape is the recorded nested-Claude-under-Codex marker leak, not a hypothetical.
  AC-3, proven by a throwaway program run over five marker sets: the ambiguous pair rendered `Runtime: ambiguous (CODEX_THREAD_ID, CLAUDECODE) — pass --host` at exit 0, and it also caught that two pi markers are NOT ambiguity.
- DONE: Deleting `requires-contract` leaves the live minor-version binding untouched — `internal/contractlint/prose_manifest_minor_sync_test.go` and `release.yml`'s `manifest-tag-gate` still green, and `plugin.json`'s `version` field unchanged.
  Exercised, not reasoned: the key was deleted from both manifests, `TestProseMinorMatchesVendoredManifestMinor` PASSed and `manifest-tag-gate v0.26.0 …` exited 0, only `TestVendoredManifestTombstoneFrozen` went red (the test that ships with the field); manifests restored, `git diff --stat` clean.
- DONE: REACH decision recorded (dispatch instruction, not a checklist item)
  All three surfaces — `--version`, launch banner, `status --boot` — get the session-aware sandbox render, so the FO's durable boot record stops carrying the false posture; the `Runtime:` line stays `--version`-only.
- DONE: The `Available()` installHint question the dispatch flagged as worth checking
  Confirmed real and confirmed NOT a reporting bug: the hint fires at `frontdoor.go:418`, `frontdoor.go:631`, `pi.go:301`, all gated on `wrap && !Available`, so a nested launch inside the sandbox advises installing safehouse and exits 1. Left out as sandbox enforcement (explicitly out of scope) and recorded as follow-up 1.

### Summary

The entity body now carries a proposed approach, six acceptance criteria, an expected surface with tolerance, a test plan, and the full documentation diff for `docs/site/reference/command-reference.md`. Three mechanisms were exercised rather than assumed: the live inversion was reproduced in this session, a throwaway Go program produced the captain-approved output shape byte-for-byte across five marker sets (catching that two same-host markers are not ambiguity), and the `requires-contract` deletion was actually performed against both gates and reverted.

Two things are recorded as costs rather than glossed. Dropping the per-host block deletes `internal/cli/host_runtime.go` entirely and takes the multi-host plugin inventory with it — `doctor` is single-host and does not replace it, so `doctor --all-hosts` is offered as a separate task if the captain wants it back. And the `installHint` nested-launch bug is real but is a launch-gate change, so it is filed as a follow-up blocked on the `Inside()` registry this task ships. Expected surface is net negative (~−60 lines across ~20 files), and AC-2 carries a measured baseline: `--version` costs 0.86s against 0.31s for a non-probing command, ~64% of it spent reporting runtimes the session is not in.

## Stage Report: ideation (cycle 2)

- DONE: Value AC measures the inversion actually being fixed: a session running inside the sandbox reports itself as inside, against a baseline that can move the wrong way — today that exact session reports "unavailable". Assert the reported posture against the real session state, never that a string changed.
  Unchanged and still held by AC-1; its test rows re-split as `SessionState` (three outcomes, `--version` and `--boot`) and `LaunchState` (four, banner), each naming the inside-and-safehouse-absent row as the regression case.
- DONE: The ambiguous-marker branch stays exit 0 and reports rather than guesses. This is the one branch that can break every boot, and its shape is the recorded nested-Claude-under-Codex marker leak, not a hypothetical.
  AC-3 unchanged; spike 4 re-confirmed exit 0 and added that the ambiguous branch renders NO session identifier, since printing one would assert a host had been resolved.
- DONE: Deleting `requires-contract` leaves the live minor-version binding untouched — `internal/contractlint/prose_manifest_minor_sync_test.go` and `release.yml`'s `manifest-tag-gate` still green, and `plugin.json`'s `version` field unchanged.
  Unchanged from cycle 1, where it was proven by actually deleting the key: sync test PASSed, `manifest-tag-gate v0.26.0 …` exited 0, only `TestVendoredManifestTombstoneFrozen` went red; manifests restored.
- DONE: Captain revision 1 — the Runtime line carries the session identifier
  New AC-2a plus an identity column in `Detect`. Spike 4 overturned the premise: `CLAUDECODE` is `1`, so identity comes from the separate `CLAUDE_CODE_SESSION_ID`; taking it from the detection marker would have rendered `session 1` for every Claude session.
- DONE: Captain revision 2 — the outside sandbox strings report availability
  Three strings became two (`not sandboxed (safehouse available|not installed)`), the `.safehouse` profile dropped from the session line, and the launch question moved to a second renderer used by the banner alone. AC-2b pins it with an assertion that no `--version` or `--boot` sandbox string contains `.safehouse`.

### Summary

Both annotations are folded in, and both narrowed the design. The identifier work found the gate note's premise wrong on the facts — `CLAUDECODE` is a flag, not an id — so identity resolves from a separate per-host variable, truncated to the eight characters Claude Code itself uses to name `~/.claude/teams/session-afd74765` (verified against the live directory), with pi rendering none because `PI_CODING_AGENT_SESSION_DIR` is the sessions collection, not a session. The empty-identity case is the general rule rather than a special case, so no empty parenthetical is reachable.

Following the second annotation to its root showed the three-surface string sharing to be part of the defect rather than a property to preserve: the banner asks a launch question, `--version` and `--boot` ask a session question. Two renderers now answer them separately, `Inside()` stays the single detection registry, and all three surfaces still stop inverting — the REACH decision is unchanged.

The surface moved and is reported rather than absorbed: +630/−600 against the previous +540/−600, so the net crosses from about −60 to about +30. That is inside the ±250-line tolerance with the file count unchanged, but the cycle-1 claim that this task is net line-negative no longer holds and is retracted here.
