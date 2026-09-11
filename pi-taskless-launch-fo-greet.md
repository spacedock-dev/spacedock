---
title: "A taskless spacedock pi launch never boots the FO: no first request, silent startup"
status: implementation
source: "Captain observation 2026-09-11 in the wrapped dev-override TUI session: silence at startup was the evidence the FO never booted; manual /first-officer produced the greeting and workflow stats. Diagnosis trail: post-s98 the frontdoor no longer appends any launch task (AC-3 removed the inert $spacedock:first-officer sentence), so a taskless launch sends no first model request; the extension bootstrap (session_start-armed, context-hook injected) is request-time-only and once-per-turn (agent_end disarms), so it never fires. Pre-s98 behavior greeted at startup because the always-appended prompt forced a first request and the old ungated extension injected."
id: n315frdw60950kjde4cx017q
gates:
    version: 1
    records:
        - id: gate:n315frdw60950kjde4cx017q:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:n315frdw60950kjde4cx017q-backlog-1
              briefing:
                id: briefing:n315frdw60950kjde4cx017q:backlog:attempt-1:revision-1
                digest: sha256:0bde44c6681bd0f95ed27aa526533c8a4ad99776979c01b4888ed8dae11d7d39
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:n315frdw60950kjde4cx017q:backlog:1
                briefing: briefing:n315frdw60950kjde4cx017q:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-09-11T15:48:17.855233Z"
                decision: approve
                reason: 'Captain approved in chat 2026-09-11 at the seed presentation: fast-track, ideation steered toward mechanism (a) — extension delivers the bootstrap as a real first-turn message via the documented pi.sendUserMessage()/expandPromptTemplates upgrade path; frontdoor stays argv-and-env-only'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:n315frdw60950kjde4cx017q:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:n315frdw60950kjde4cx017q-ideation-1
              briefing:
                id: briefing:n315frdw60950kjde4cx017q:ideation:attempt-1:revision-1
                digest: sha256:21495f93b61ac4aebb40968429cd7373c2163901410cdaa33edbfa999c772a26
                room-ref: '@review/ideation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:n315frdw60950kjde4cx017q:ideation:1
                briefing: briefing:n315frdw60950kjde4cx017q:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-09-11T18:10:20.74954Z"
                decision: approve
                reason: 'Captain approved in chat 2026-09-11 at the ideation gate: mechanism (a) with the live-verified spike accepted — extension-owned real-message bootstrap via sendUserMessage, env-only taskless marker, fallback retained, lean ACs with existing test owners'
              application:
                target-stage: implementation
                state: consumed
started: 2026-09-11T15:48:32Z
---
Problem: a fresh `spacedock pi` launch with no operator task sits silent — the FO never presents the boot summary or workflow stats until the operator types something (and a weak model may never self-boot even then; glm-5.3-flash misreported its own context twice today). The pre-s98 launch greeted at startup; s98's correct argv fix removed the trigger without replacing it.

Deliverable: a taskless launch boots and greets (binary gate, boot identify, session summary, stop for input) with no operator input. Two candidate mechanisms to weigh at ideation: (a) the extension delivers the bootstrap as a REAL first-turn message via pi.sendUserMessage()/expandPromptTemplates — the upgrade path the s98 body already documents — which also makes the bootstrap persist instead of being once-per-turn and request-time-only (the observability trap hit today: multi-turn probes made a working injection look absent); (b) the frontdoor appends a neutral default launch task when no operator task is given (no $spacedock syntax — AC-3 intact), restoring the startup request; the bootstrap stays ephemeral.

Acceptance criteria and test plan to be fleshed out at ideation. Value AC direction: a fresh taskless launch (wrapped and unwrapped, installed and dev-override) presents the interactive greet without operator input, measured against today's silent-startup baseline; plain `pi` sessions still receive nothing (single-owner gate preserved).

## Proposed approach — mechanism (a) selected

The extension delivers the FO bootstrap as a REAL first-turn user message via `pi.sendUserMessage()` on `session_start` (reason `"startup"`), replacing the once-per-turn context-hook injection as the primary delivery path. The frontdoor stays argv-and-env-only: it gains one more env marker, `PI_SPACEDOCK_LAUNCH_TASKLESS=1`, set only when the launch carries no operator task argv and no resume passthrough (the exact condition where today nothing triggers a first request). No frontdoor-owned launch text is introduced or returned — AC-3 (frontdoor argv-and-env-only) stays intact; the greet text remains extension-owned (`FO_BOOTSTRAP_TEXT`, unchanged wording, single owner preserved).

Mechanisms, each tied to the value AC it serves:

1. **Taskless marker (`PI_SPACEDOCK_LAUNCH_TASKLESS=1`)** — serves the value AC: it is the only signal that distinguishes "taskless launch, no first request will ever happen" from "taskful launch / resume, a first request already exists or history is present". The extension cannot derive this itself: at `session_start` the argv prompt is not yet in the session, so a session-content check is indistinguishable between taskless and taskful. Simplest alternative considered: the frontdoor appends a neutral default launch task when no operator task is given (mechanism (b)) — rejected because it reintroduces a frontdoor-owned launch string adjacent to the AC-3 boundary the captain drew at s98, and it leaves the bootstrap ephemeral (request-time-only, invisible in transcripts — the observability trap that cost today's diagnosis). A second alternative — the extension sniffing `process.argv` of the pi process — is rejected as fragile (passthrough reshapes argv) and an undeclared surface.
2. **`pi.sendUserMessage()` at `session_start`** — serves the value AC: it both triggers the first model turn (the thing a taskless launch lacks) and persists the bootstrap as a normal transcript user message, which fixes the once-per-turn request-time-only ephemerality by construction (no disarm needed for the primary path; `agent_end` disarm survives only for the fallback arm). `expandPromptTemplates` stays **false** (the default): `FO_BOOTSTRAP_TEXT` deliberately contains no `/skill:` or template syntax (the s98 resolvable-trigger design), so expansion would be a no-op dependency; the s98-documented `expandPromptTemplates` upgrade path (delivering the skill body inline) remains out of scope. Simplest alternative: keep the context hook and make the frontdoor always send a task — that is mechanism (b), rejected above.
3. **Context-hook bootstrap retained as FALLBACK** — the existing `context`-hook injection arm stays, but only arms when the real-message send FAILS. Failure behavior (named for the gate): if `pi.sendUserMessage()` throws (API drift, sub-floor pi regression, unexpected streaming state at startup), the extension catches the error, leaves `injectBootstrap = true` so the context hook still injects `FO_BOOTSTRAP_TEXT` at the first request exactly as today, and does not crash the session. On a successful send, `injectBootstrap` is disarmed for that turn so the first request carries exactly one bootstrap (the real message, which also lands in the transcript). The compaction boot-record arm is untouched. Fallback test: a fake `pi` whose `sendUserMessage` throws must still yield a context-hook injection on the first `context` event.
4. **Gating composition** — the send fires only when ALL hold: `PI_SPACEDOCK_LAUNCH=1`, `PI_SPACEDOCK_LAUNCH_TASKLESS=1`, not `PI_SUBAGENT_CHILD=1`, and `event.reason === "startup"` (never on resume/new/fork/reload — a resumed or forked session has history and must not be greeted). The single-owner gate is preserved: a plain `pi` session (no markers, even with the package installed) receives nothing.

### Spike result (riskiest mechanism exercised first — 2026-09-11)

The risk was whether `pi.sendUserMessage()` from a `session_start` handler actually delivers at/before the first turn in an interactive TUI session (the exact failing shape) and does NOT fire ungated. Exercised live against installed pi 0.85.1 with a throwaway extension (`/tmp/pi-greet-spike/spike.ts`) that calls `pi.sendUserMessage(...)` on `session_start` reason `"startup"`:

- **Positive arm** (`PI_SPIKE_SEND=1 PI_SPACEDOCK_LAUNCH=1 pi --extension /tmp/pi-greet-spike/spike.ts` in a tmux TUI, temp cwd, no operator input): the TUI showed the agent Working within ~10s of launch with no operator keystroke; the session transcript `~/.pi/agent/sessions/--private-tmp-pi-greet-spike-work--/2026-09-11T16-19-42-402Z_*.jsonl` contains a **persisted role=user message** with the bootstrap text followed by assistant turns (the model began executing the greet: "Simple: greet with one-line boot summary and stop..."). Delivery at/before the first turn: CONFIRMED. Persistence (fixes the request-time-only trap): CONFIRMED.
- **Control arm** (same launch, both markers unset): TUI idle, session directory created but EMPTY — no session file, no user message, no turn. This is today's silent-startup baseline, measured: **0 persisted user messages, 0 model turns before operator input**. Gate preserved: CONFIRMED.
- Spike note: the weak default model (glm-5.3-flash-background) spent a long first turn spelunking for skills instead of greeting crisply — greet QUALITY is the FO skill's business; the mechanism's job (deliver the bootstrap as a real first turn) is what this entity owns.
- A deterministic vitest-harness spike in the pi monorepo checkout was attempted and abandoned (no node_modules in the scratch checkout); the live TUI spike is the stronger evidence and is the shape the ACs reuse.

### Declared pi-behavior dependencies

- `pi.sendUserMessage(content)` — sends a real user message (transcript-persisted) and always triggers a turn when idle; throws when streaming without `deliverAs` (at `session_start` reason `startup` the agent is idle by construction). Exercised live at pi 0.85.1.
- `session_start` event with `reason` (`startup`/`reload`/`new`/`resume`/`fork`) and the `ExtensionUIContext`-free handler form already in use.
- Existing floor `pi >= 0.83.0` (binary-level, ready-gate/doctor) already covers the context-hook and `<available_skills>` behaviors; `sendUserMessage` predates the floor (documented in s98's dependency list). No floor bump required.

## Acceptance criteria

- **AC-1 (value, live, measures end value against the silent-startup baseline).** A fresh taskless `spacedock pi` launch — the 2×2 of wrap/non-wrap × installed/dev-override, at minimum non-wrap installed and non-wrap dev-override live-probed plus one wrap arm — starts the FO greet turn with NO operator input: within 60s of launch the session transcript contains a role=user message carrying the `SPACEDOCK-FO-BOOTSTRAP-v1` marker followed by at least one assistant response. Baseline (measured in the spike control, 2026-09-11): 0 user messages, 0 model turns before operator input. Falsifying change: reverting the extension's send (or the frontdoor's taskless marker) returns the transcript to the empty-until-input baseline.
- **AC-2 (gate, live + deterministic).** A plain `pi` session (extension installed, no `PI_SPACEDOCK_LAUNCH`) receives no bootstrap message and starts no turn before operator input — session transcript empty until the operator types; deterministically, the extension calls `sendUserMessage` zero times and injects nothing into `context` when the markers are absent. Falsifying change: removing the marker gate makes the plain-session unit test fail (send observed).
- **AC-3 (frontdoor authority, deterministic).** The frontdoor's only change is environmental: `PI_SPACEDOCK_LAUNCH_TASKLESS=1` is set exactly when no task positional and no resume passthrough is present, forwarded through safehouse `--env-pass` under wrap, and NEVER appears as appended argv text or a returned launch string. Falsifying change: passing a task (or `--resume`) makes the marker-asserting Go test fail; any greet text in argv makes the argv-shape assertion fail.
- **AC-4 (taskful/resume unchanged, deterministic).** A launch WITH an operator task (or a resume) sends no greet message — no double first turn; the extension's `sendUserMessage` is called zero times when the taskless marker is absent, and the existing context-hook first-request injection behavior for taskful frontdoor launches is preserved. Falsifying change: dropping the taskless-marker check makes this unit test fail (send observed on a taskful fake launch).
- **AC-5 (fallback, deterministic).** When `sendUserMessage` throws, the session does not crash and the first request still carries the bootstrap via the retained context-hook arm. Falsifying change: removing the fallback arming makes the throwing-fake unit test fail (no context injection after the throw).
- **AC-6 (persistence/observability).** The bootstrap on a taskless launch is a real transcript message persisted in the session file, not a request-time-only injection — subsumed by AC-1's transcript assertion and the unit assertion that delivery goes through `sendUserMessage` (not only the `context` hook). Falsifying change: reverting to context-hook-only delivery makes the transcript assertion fail against the baseline.

### Test plan (lean — no new harness, per the captain's standing no-test-infra directive)

- Primary proof owners that exist and are extended: `test/spacedock.test.ts` (extension gate harness — extend with send-path cases: taskless send fires with the marker text; taskful/resume/plain/fork send zero times; throwing send falls back to context-hook arming) and `internal/cli/pi_frontdoor_test.go` (frontdoor env/argv harness — extend with the taskless marker cases: set on taskless non-wrap and wrap (`--env-pass`), absent on taskful and resume, argv shape unchanged).
- Live one-off launches for AC-1/AC-2 (the same shape as the spike, four arms): taskless installed non-wrap, taskless dev-override non-wrap, taskless wrap installed, plain `pi` control. Each launch is observed via its session transcript (user message + assistant turn present/absent) — no harness, one-off manual/live validation only.
- Distinct falsifying edits are named per AC above; no test infrastructure is created for its own sake.

## Expected surface

| File | Change | Rough LOC |
|---|---|---|
| `.pi/extensions/spacedock.ts` | session_start send path + taskless gate + fallback disarm | +40 / −8 |
| `internal/cli/pi.go` | `piLaunchTasklessEnv` const, set in `launchEnvList`, `--env-pass` under wrap | +14 / −0 |
| `internal/cli/pi_frontdoor_test.go` | taskless marker env/argv assertions (wrap + non-wrap + taskful/resume negatives) | +55 / −0 |
| `test/spacedock.test.ts` | send-path unit cases per AC-2/4/5/6 | +70 / −0 |
| `docs/site/contributing/adding-a-runtime.md` | doc diff below | +18 / −6 |

**Estimate: net +183 LOC (+197 insertions / −14 deletions), across 5 files. Tolerance ±60 net.** Ideation runs without a worktree; implementation applies these.

**Observable semantics this task may change:** runtime behavior (a taskless frontdoor launch now boots the FO with a real first-turn message; taskful launches, resumes, plain pi sessions, and subagent children behave exactly as today); stored format adjacency (session transcripts of taskless launches now begin with the bootstrap as a normal role=user message — declared, this is the fix, not a side effect); new env surface `PI_SPACEDOCK_LAUNCH_TASKLESS` (env-only marker, same family as `PI_SPACEDOCK_LAUNCH`). No command grammar, no authority, no stored-format schema change.

### Bounded edges (declared, not open scope)

- An operator who smuggles their own prompt through `--` passthrough on an otherwise taskless call gets the greet queued AND their prompt; the frontdoor's `hasTask` does not see passthrough positionals. Operator-owned bypass; documented in the doc diff.
- A taskless launch that also passes `-p`/print-mode via passthrough fires the greet into a non-interactive run — degenerate input (a taskless print run has nothing to print); accepted, documented.
- Greet quality on weak models (glm-5.3-flash rambling on first boot, observed in the spike) is out of scope: this entity guarantees delivery + turn trigger, not the model's greeting prose.

### Doc diff (ideation-owned, applied at implementation)

`docs/site/contributing/adding-a-runtime.md`, section "First-officer contract bootstrap (Pi)":

Before (paragraph 1, excerpt):
> ...the extension injects the contract through Pi's context hook, gated on the `PI_SPACEDOCK_LAUNCH=1` marker `spacedock pi` sets on every launch, so a plain `pi` session for unrelated work receives no bootstrap.

After:
> ...the extension delivers the contract as a REAL first-turn user message via `pi.sendUserMessage()` on a taskless launch — `spacedock pi` with no operator task sets `PI_SPACEDOCK_LAUNCH_TASKLESS=1` alongside `PI_SPACEDOCK_LAUNCH=1`, and the extension sends the bootstrap at `session_start`, which boots the FO greet with no operator input and persists the bootstrap in the session transcript. On a launch WITH an operator task (or a resume) the extension sends nothing — the task itself is the first request. If the send fails, the extension falls back to the original context-hook injection at the first request. Both markers are env-only; the frontdoor never carries launch text (AC-3), and a plain `pi` session for unrelated work still receives no bootstrap.

Before (bullet):
> - **Request-time only.** Context-hook injections are never persisted to session logs. Absence of the bootstrap message in a session log is NOT evidence that it was absent at request time.

After:
> - **Delivery path.** On a taskless launch the bootstrap is a real transcript message — its presence in the session log IS evidence it was delivered. Only the fallback path (context-hook injection after a failed send, and the compaction boot-record) remains request-time-only and invisible in logs.

Plus one bullet in the same section's fact list:
> - **Taskless greet.** A taskless `spacedock pi` launch greets without operator input; passing a prompt through `--` passthrough on an otherwise taskless call queues it behind the greet (operator-owned bypass).

And the "Skill install and load paths" paragraph's phrase "installs the FO contract through Pi's context hook, gated on the `PI_SPACEDOCK_LAUNCH=1` marker the frontdoor sets on every launch" becomes "installs the FO contract as a real first-turn message on taskless launches (context-hook fallback otherwise), gated on the `PI_SPACEDOCK_LAUNCH=1` / `PI_SPACEDOCK_LAUNCH_TASKLESS=1` markers the frontdoor sets".

## Stage Report: ideation

- DONE: Fleshed-out body with mechanism (a) as the selected design — the extension delivers the bootstrap as a REAL first-turn user message via pi.sendUserMessage() at session_start (reason "startup"), gated on a new frontdoor env marker PI_SPACEDOCK_LAUNCH_TASKLESS=1; frontdoor stays argv-and-env-only (AC-3 intact, greet text remains extension-owned FO_BOOTSTRAP_TEXT); expandPromptTemplates explicitly stays false with rationale.
  Entity body "Proposed approach — mechanism (a) selected" + "## Acceptance criteria"; mechanism (b) rejected in-body (frontdoor-owned launch string re-adjacent to the AC-3 boundary; leaves the bootstrap ephemeral).
- DONE: Existing context-hook bootstrap disposition stated: retained as FALLBACK when sendUserMessage throws (catch, leave injectBootstrap armed, first request still carries the bootstrap via the context hook); named failure behavior and its test (throwing-fake unit case in test/spacedock.test.ts, AC-5). Compaction boot-record arm untouched.
  Entity body mechanism 3 + AC-5.
- DONE: Once-per-turn request-time-only ephemerality fixed by the mechanism (real message persists in the transcript; no disarm needed on the primary path; agent_end disarm survives only for the fallback arm).
  Entity body mechanism 2 + AC-6; spike transcript shows the persisted role=user bootstrap message.
- DONE: Spike first, riskiest mechanism — live TUI exercise against installed pi 0.85.1: positive arm (markers set) delivered a persisted first-turn user message and booted the model turn with zero operator input; control arm (markers unset) stayed silent (empty session dir, 0 turns) — single-owner gate preserved. Result recorded in the body.
  Entity body "Spike result"; evidence /tmp/pi-greet-spike/ and ~/.pi/agent/sessions/--private-tmp-pi-greet-spike-work--/2026-09-11T16-19-42-402Z_01a09144-86c1-7334-94a2-830c18353adb.jsonl (persisted user msg + assistant turns) and --private-tmp-pi-greet-spike-control-- (empty dir = baseline).
- DONE: Lean ACs per the no-test-infra directive — AC-1 measures the end value against the measured silent-startup baseline (0 user messages / 0 turns before input → greet turn within 60s); gate AC-2 = plain pi receives nothing; proof = four one-off live launch arms plus extensions of the two EXISTING test owners (test/spacedock.test.ts, internal/cli/pi_frontdoor_test.go); no new harness.
  Entity body "## Acceptance criteria" + "Test plan".
- DONE: Expected surface declared — 5 files, net +183 (+197/−14), tolerance ±60; observable semantics declared (runtime behavior + env marker + transcript-adjacency; no grammar/authority/schema change); bounded edges listed (passthrough prompt smuggling, print-mode taskless, weak-model greet quality).
  Entity body "Expected surface" + "Bounded edges".
- DONE: Doc diff for the user-visible launch-behavior change proposed against docs/site/contributing/adding-a-runtime.md ("First-officer contract bootstrap (Pi)" paragraph, "Request-time only" bullet replaced by "Delivery path", new "Taskless greet" bullet, "Skill install and load paths" phrase).
  Entity body "Doc diff (ideation-owned, applied at implementation)".

### Summary

Selected and spike-verified mechanism (a): the spacedock pi extension sends the FO bootstrap as a real first-turn user message via pi.sendUserMessage() on taskless frontdoor launches (new PI_SPACEDOCK_LAUNCH_TASKLESS=1 env marker from the frontdoor, env-only so AC-3 holds), which boots the greet without operator input and persists the bootstrap in the transcript — fixing both the silent startup and the request-time-only observability trap. The live TUI spike on pi 0.85.1 confirmed delivery (persisted user message + booted turn, no operator input) and the negative gate (plain pi stays silent). Body now carries exact-heading Acceptance criteria, a per-AC falsifying-edit test plan extending only the two existing test owners, expected surface/LOC, declared semantics, bounded edges, and the concrete doc diff for the ideation gate.
- Cycle 2 (2026-09-11, captain directive at the ideation gate presentation, post-consume): single delivery path — NO context-hook bootstrap fallback. `sendUserMessage` is the only bootstrap delivery; the fallback arm, its `agent_end` disarm retention, and the throwing-fake test case are dropped. Send failure behavior: no bootstrap, no crash, no silent contract loss on taskful launches (taskful launches keep the first-request context-hook bootstrap? NO — single path: send failure means no bootstrap, named in the report). Compaction boot-record arm stays (separate mechanism, rides requests by nature). Surface shrinks accordingly; the implementation worker updates the body's approach/AC-5/surface to match.
