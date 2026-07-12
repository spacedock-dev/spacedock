---
title: Move prose-polish routing policy into the workflow
status: ideation
score: "0.70"
source: "Captain correction after c6 ideation review."
id: csb4c89dteavbq1htdac7fwm
started: 2026-07-12T23:19:51Z
---

# Move prose-polish routing policy into the workflow

## Problem

The shared First Officer dispatch contract names a standing prose-polisher
convention and describes when drafts may route through it. That policy is not
universal: `comm-officer` exists because this workflow declares it in its mod.
The shared layer should own only generic standing-teammate discovery, lifecycle,
addressing, the turn-starting `teammate.request` operation, and mod projection.

The split caused the c6 incident. Git commit `8f40cf60` durably records a
9,229-word task body before late cleanup. Session testimony says no polish route
reached the writer before that point, but no durable tool-event artifact proves
the zero-route claim. Treat that claim as incident testimony, not as an existing
fixture. Implementation must first reproduce the failure by running the new
registered scenario against the pre-change revision and retaining its tool-event
and state artifacts as the RED baseline.

## Proposed direction

The ideation stage owns qualification and commit ordering; the comm-officer mod
owns caller behavior; shared dispatch owns only generic projection and the
turn-starting `teammate.request` operation. Use those existing seams. Do not add
a routing registry, policy schema, or routing engine.

### Ownership and trigger

Before an ideation dispatch, the FO applies the stage's existing staff-review
complexity rule. When it decides the task is complex, it writes the exact line
`Prose-polish: required (staff-review-complex)` into the existing dispatch scope
notes. Scope notes are already carried into the worker's dispatch file, so the
decision reaches the writer without a new field or policy format.

The writer independently counts the final draft. The counted narrative region
starts after the closing YAML-frontmatter `---` and ends before the first
`## Stage Report` heading. YAML frontmatter and every Stage Report are excluded.
Within that region, the threshold is the whitespace-delimited word count
reported by `wc -w`: 1,499 does not qualify; 1,500 does.

A body qualifies when either:

- its scope notes carry `Prose-polish: required (staff-review-complex)`; or
- its counted narrative region is at least 1,500 words.

For a qualifying body, the writer MUST make one polish request at the final
draft boundary and before the first commit that records that final draft or its
stage-state transition. Meeting both triggers still produces exactly one
request. "MUST" applies only to the attempt and its ordering. Receipt and
acceptance remain best-effort; the writer retains ownership of technical meaning,
qualifiers, criteria, and evidence.

The comm-officer mod owns the prose-polish purpose and caller contract: the
target teammate, modes, two-minute bound, exclusions, and fallback. A response
is non-load-bearing. The timeout clock begins when the runtime records a
successful turn-starting route attempt. An absent target or explicit rejection
falls back immediately. A successfully routed but silent target falls back at
two real wall-clock minutes; no fake clock is assumed. The writer then commits
the original narrative and records the exact Stage Report line
`comm-officer unavailable; proceeded unpolished`. No retry or replacement agent
is required. If a timely response changes protected semantics, the writer rejects
the edit and immediately commits the original narrative. That is a reviewed
response disposition, not teammate unavailability, so it records no fallback
line.

The exclusions are direct captain chat, a non-complex ideation body below 1,500
words, short operational status, tool output, logs, commit messages, and prose
that cannot be separated from secrets or sensitive context. Excluded content
causes no polish request and needs no fallback note. Exclusions are evaluated
before either qualification trigger, so a long log or sensitive body never
routes merely because it crosses 1,500 words.

### Propagation seam

`dispatch build` already composes two fetches into a worker assignment:
`show-stage-def` supplies the workflow stage and `show-standing` projects each
declared mod's `## Routing Usage`. Preserve that shape and close its narrow host
gap: every non-bare named worker for a workflow with a standing mod receives the
standing-mod fetch, whether or not a legacy Claude `team_name` exists.

The projected caller contract invokes `teammate.request(name, payload)`: start or
resume one turn for the addressed declared teammate and report success only when
the runtime accepts a turn-starting call. It does not say `SendMessage`, "team",
or "mirror Claude". The host bindings are explicit:

- Codex uses `followup_task` for running, idle, or completed-but-addressable
  targets. `send_message` is preservation-only and never satisfies a request.
- Claude legacy and merged use `SendMessage`, whose standing-target call starts
  or resumes work.
- Pi uses its named turn-starting teammate route.
- With no turn-starting binding, the target is immediately unavailable. The
  caller does not start a two-minute clock on an unwakeable queued message.

The Codex runtime already defines `followup_task` as turn-starting and
`send_message` as non-triggering under addressable-worker reuse. Implementation
reuses that proven distinction for standing-teammate requests rather than
inventing a second Codex lifecycle. Bare mode projects nothing.

This is the smallest existing mechanism. The stage definition reaches every
ideation writer already, the mod parser and `show-standing` renderer already
exist, and the dispatch file already carries fetch commands. No new policy
representation or lookup path is needed.

`teammate.request` is generic standing-teammate mechanics, not a prose-policy
abstraction. Declared standing roles already have multiple consumers and every
host must distinguish a turn-starting request from discovery, queue-only
preservation, and teardown. The semantic names that existing compatibility
boundary once; it adds no registry, configuration format, or policy evaluator.

### Implementation seam

- Make standing-mod enumeration team-name agnostic. `dispatch build` emits the
  existing `show-standing` fetch when `bare_mode` is false and a standing mod is
  declared; legacy team identity remains relevant only to spawn/dedup.
- Move the rendered standing-usage header out of the Claude-specific package (or
  otherwise make that renderer host-neutral). Its exact generic header is:
  `Declared standing teammates are addressable through teammate.request. Each
  projected mod below owns purpose, triggers, timeout, exclusions, and fallback.
  A queue-only preservation message is not a request.` Remove the current
  `SendMessage`, `team`, `2-minute`, `un-polished`, and shared-core "full routing
  contract" claims.
- Add the `teammate.request` semantic to the existing dispatch core and bind it
  in the three runtime adapters as specified below. Do not add a request registry
  or a new runtime state machine.
- Consolidate the comm-officer mod's duplicated `Routing guidance` and `Routing
  Usage` into the single projected runtime-neutral caller contract below. Keep
  its four payload modes, but remove caller-side host tool names.

### Exact process-documentation changes

`docs/dev/README.md`, inside the existing `### ideation` Outputs list:

Before:

```markdown
  - When captain feedback changes the target behavior, update the task body, acceptance criteria, and test plan together before re-validating.
  - For template or skill text changes: specific before/after wording, not just "change X".
```

After:

```markdown
  - When captain feedback changes the target behavior, update the task body, acceptance criteria, and test plan together before re-validating.
  - Prose polish for this workflow:
    - Before dispatch, the FO applies the Staff review rule. When it decides the ideation is complex, scope notes include `Prose-polish: required (staff-review-complex)`.
    - Apply the comm-officer exclusions before either qualification trigger. Excluded content never routes because of complexity or word count.
    - The writer counts the final narrative after YAML frontmatter and before the first `## Stage Report` with `wc -w`; YAML and all Stage Reports are excluded. A staff-review-complex draft or a narrative of at least 1,500 words requires exactly one route using the projected comm-officer `## Routing Usage`.
    - The route starts before the first commit that records the final draft or its stage-state transition. The attempt and ordering are required; reply and acceptance are best-effort, and technical meaning, qualifiers, criteria, and evidence remain writer-owned.
  - For template or skill text changes: specific before/after wording, not just "change X".
```

`docs/dev/_mods/comm-officer.md`, replace the complete current `## Routing Usage`
body. Before:

```markdown
Four caller patterns (mirror Claude's Read/Edit/Write tool shapes). Pick the pattern first, then format the SendMessage body to match.

1. **Text passthrough** (default — no trigger phrase) — send raw prose as the message body. Reply: polished text first, then `---` + `**Polish notes**` block. Caller places the result.
2. **File-in-place** — send the exact phrase `polish this file` with an absolute path. Teammate Edits/Writes the file in place. Reply: one-line receipt + `---` + `**Polish notes**`.
3. **Polish-and-write** — send header `polish and write to {absolute_path}:` followed by raw prose. Teammate Writes the polished prose to that path (create-or-overwrite). Reply: one-line receipt + `---` + `**Polish notes**`.
4. **Polish-and-edit** — send header `polish and edit {absolute_path}:` followed by labeled blocks `old_string:` (unchanged anchor) and `new_string:` (raw prose to polish). Teammate polishes `new_string` and Edits the file at that anchor. Reply: one-line receipt + `---` + `**Polish notes**`.

Notes block fields: `Mode`, `Guide applied`, `Changes`, `Flagged for review`. Absolute paths required for patterns 2-4; no inferred targets. Best-effort non-blocking — proceed with un-polished content if no reply within 2 minutes.
```

Rename `## Routing guidance (for FO and ensigns)` to `## Purpose and scope` and
retain its two existing Scope lists. Delete this exact duplicated caller block
from that section because the new `## Routing Usage` below is its single home:

```markdown
**Four usage patterns (mirrors Claude Code's read/Edit/Write tool shapes):**

1. **Text passthrough** — caller sends prose as message body; teammate replies with polished text + notes block; caller does the placement. Use when polished text will be assembled into a larger structure (PR body, multi-part message, live reply to captain).
2. **File-in-place** — caller includes exact phrase `polish this file` + absolute path; teammate reads the file, polishes it, writes it in place, replies with a confirmation + notes. Use when a file already exists on disk with unpolished prose to tighten.
3. **Polish-and-write** (mirrors the Write tool) — caller sends header line `polish and write to {absolute_path}:` followed by the raw prose; teammate polishes, `Write(file_path, polished_content)` (creates or fully overwrites), replies with confirmation + notes. Use when creating a new file whose content IS polished prose (e.g., a draft narrative block).
4. **Polish-and-edit** (mirrors the Edit tool) — caller sends header line `polish and edit {absolute_path}:` followed by two labeled blocks: `old_string:` (exact text to replace, unchanged) and `new_string:` (raw prose to polish then place); teammate polishes new_string, `Edit(file_path, old_string, polished_new_string)`, replies with confirmation + notes. Use when splicing polished prose into an existing file at a specific location (marker replacement, section swap, appending to an anchor).

Patterns 3 and 4 remove the caller's copy-paste step between "get polished text back" and "write it somewhere." Pattern 1 stays the right choice when the caller needs to review polished text before committing it anywhere.

**Hard rules:**

- MUST NOT block on `comm-officer` reply. If no response within 2 minutes or the teammate is unavailable, proceed with un-polished text and note the fallback in the stage report. Polish is best-effort, not load-bearing.
- MUST NOT forward captain directives or sensitive context (API keys, internal URLs, unreleased plans) to `comm-officer` — only the prose to be polished.
```

After:

```markdown
Call `teammate.request("comm-officer", prose)` to start or resume the comm officer. A queue-only preservation message is not a successful request.

For this workflow's ideation stage, route exactly once when scope notes contain `Prose-polish: required (staff-review-complex)` or the final narrative region is at least 1,500 whitespace-delimited words. The narrative starts after YAML frontmatter and ends before the first `## Stage Report`; YAML and all Stage Reports are excluded. Route before the first commit that records the final draft or its stage-state transition.

Send only the prose to polish. Do not send captain directives, secrets, internal URLs, or unreleased plans. Direct captain chat, a non-complex narrative below 1,500 words, short operational status, tool output, logs, and commit messages are excluded and produce no route or fallback note. Evaluate exclusions before complexity or word count.

Payload modes remain text passthrough; `polish this file {absolute_path}`; `polish and write to {absolute_path}:`; and `polish and edit {absolute_path}:` with exact `old_string:` and `new_string:` blocks. The caller chooses the target and reviews the returned change.

The required obligation is the route attempt and its ordering. A reply and its acceptance are best-effort. Preserve technical meaning, qualifiers, criteria, and evidence; the caller rejects any returned change that alters them, immediately commits the original narrative, and records no unavailability fallback.

The two-minute clock starts at the runtime's recorded successful turn-starting attempt. If the target is absent or explicitly rejects the request, proceed immediately. If a successfully routed target is silent, proceed after two real wall-clock minutes. Commit the original narrative and add this exact Stage Report line: `comm-officer unavailable; proceeded unpolished`. Do not retry or spawn a replacement.
```

`skills/first-officer/references/fo-dispatch-core.md` replaces this complete
workflow-specific paragraph. Before:

```markdown
**Routing through a standing prose-polisher.** When composing drafts for captain review (PR bodies, gate-review summaries, long narrative entity-body sections, debrief content), the FO MAY route through a live standing prose-polisher (convention: `comm-officer`). Best-effort, non-blocking regardless of duration; if absent, proceed un-polished. Read the polisher's usage (when to polish, the polish modes) from its mod.
```

After:

```markdown
**`teammate.request(name, payload)`.** Start or resume one turn for an addressed declared standing teammate. The runtime adapter binds the turn-starting operation and reports unavailable when it has no such binding or target. A queue-only preservation message is not a successful request. Purpose, triggers, required/best-effort semantics, timeout, exclusions, and fallback come only from the workflow and the teammate's projected mod usage.
```

The host adapters add these exact bindings:

```markdown
- Claude: `teammate.request` → `SendMessage` to the declared standing target.
- Codex: `teammate.request` → `followup_task` for a running or completed-but-addressable target; `send_message` is preservation-only and never satisfies the operation.
- Pi: `teammate.request` → the named turn-starting teammate route.
```

The comm-officer Agent Prompt also sheds caller-tool and residency assumptions.
Replace `Each SendMessage is a discrete standalone message` with `Each response
is a discrete standalone delivery`, and replace `Stay live. Go idle between
tasks` with `After each response, return to the runtime's addressable idle or
completed state; do not self-shutdown. The dispatcher may resume you through
teammate.request.` Online readiness and replies use the runtime's normal teammate
reply channel, not a named host tool.

No public CLI syntax changes, so the docs site needs no command-reference edit.

### Mechanism check

No implementation spike is needed for parsing: this live Codex dispatch
demonstrated that `show-stage-def` reaches the writer while `show-standing` is
missing, and `spacedock dispatch show-standing --workflow-dir docs/dev`
demonstrated that the existing parser/renderer emits the mod body. Runtime
routing behavior is not treated as proven; the registered live scenario below
must first RED on the pre-change revision, then pass after implementation.

## Acceptance criteria

**AC-1 (VALUE — early route): A c6-shaped complex ideation improves from zero successful pre-commit turn-starting routes on the reproducible pre-change RED run to exactly one before the commit that records the final draft or stage transition.**

Verified by: add `prose-polish-routing` to the registered shared runtime
scenarios and run it first from a detached pre-change checkout, retaining the
raw tool-event stream, entity body, state git log, and clean status. Re-run after
the change. The host adapter records the successful route event and the first
state commit in one ordered scenario log; the route sequence must precede the
commit sequence. The 9,229-word c6 commit remains durable incident context, while
its zero-route claim remains labeled testimony.

**AC-2 (runtime-neutral projection): Every non-bare named dispatch projects declared mod routing without Claude-specific caller prose; bare and no-mod dispatches do not.**

Verified by: a table-driven dispatch test builds Claude legacy/team-name,
Claude merged/no-team-name, Codex/no-team-name, and Pi/named requests, executes
their emitted fetch commands, and finds independently sourced stage and mod
fixture markers in composed input. It separately asserts zero standing fetches
for bare mode and a workflow with no standing mod. The projected generic header
contains no `SendMessage`, `team`, or `mirror Claude` wording.

**AC-3 (qualification and single route): The FO's recorded complexity decision and the exact 1,500-word boundary reach the writer and produce zero or one route as specified.**

Verified by: the registered live scenario drives four inert-sentinel bodies:
non-complex 1,499 words (zero), non-complex 1,500 words (one), scope-noted complex
under 1,500 (one), and both triggers (one, not two). Each count uses only the
defined narrative region; YAML and Stage Reports contain extra decoy words that
must not affect the result.

**AC-4 (bounded fallback and commit gate): An absent or unwakeable target falls back immediately and a successfully routed silent target falls back at two real wall-clock minutes; a commit occurs only after response disposition or genuine fallback, fallback commits preserve the original narrative plus the exact durable note, and a semantically rejected response commits the original without a false unavailability note.**

Verified by: registered live arms exercise absent, no-turn-starting-binding, and
deliberately silent targets plus accepted and semantic-drift responses. The runner times from the
recorded successful request event; the silent arm uses real wall time and
observes no commit before 120 seconds, followed by the bounded commit. The
accepted arm commits only after the response event and semantic acceptance. The
semantic-drift arm rejects the edit, then commits the original immediately with
no unavailability line.
Fallback commits contain `comm-officer unavailable; proceeded unpolished`. The
before/after SHA-256 covers only the defined narrative region, so the Stage
Report note does not invalidate the unchanged-draft proof.

**AC-5 (turn-starting behavior and exclusions): Actual runtime routes start the target's turn, while excluded content produces no route or fallback note.**

Verified by: live adapters observe the actual request tool event. Codex covers
active, completed-but-addressable, and absent targets: the first two wake through
`followup_task`; queue-only `send_message` fails the route assertion; absence
falls back immediately. Claude merged and Pi named run the same registered
scenario; the legacy Claude `SendMessage` event shape is covered by a recorded
adapter fixture while that runtime is unavailable live. A no-mod arm produces no
request. Exclusion arms for direct chat, operational status, tool output, log,
commit message, and sensitive content use inert sentinels only and observe zero
request events.

**AC-6 (no new abstraction): The implementation changes the existing stage definition, mod projection, dispatch conditional, and tests without adding a routing registry, policy schema, or generalized engine.**

Verified by: validation reviews the product diff and pairs that structural fact
with AC-2's executable dispatch proof; adding a new policy package, registry
type, or config format fails this criterion even if tests pass.

**AC-7 (narrow ownership): Identical generic dispatch machinery does not invent prose policy for another workflow.**

Verified by: a control live/recorded fixture declares a different standing
reviewer with unrelated `Routing Usage` and an ideation stage with no prose rule.
The same complex draft causes no comm-officer route and no prose fallback there,
while the dev workflow routes according to its own stage and mod. A structural
contractlint check may additionally keep the deleted named shared paragraph
absent, but it is not the behavior proof.

## Test plan

First add the host/mode dispatch table beside the existing standing/build tests.
It executes emitted fetches and proves projection for Claude legacy, Claude
merged, Codex, and Pi plus bare/no-mod negatives. Estimated cost: small and
offline.

Then add `prose-polish-routing` to `sharedRuntimeScenarios`, with Claude, Codex,
and Pi runners and CI registration. The scenario uses temporary split-root state,
inert generated word-count bodies, raw runtime request/response events, an
ordered commit event, state git log, entity body, and clean status. It covers the
four trigger arms, exclusions (evaluated before triggers), control workflow,
accepted and semantic-drift responses, Codex active
and completed targets, absent/no-binding/no-mod targets, and one real two-minute
silent target. This is high-cost live integration work; do not replace it with a
Go fake clock or prose inspection. Run it once against the pre-change revision
for RED evidence and once after the fix.

Run focused tests, `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test
./... -race`. The host-neutral dispatch-core and shared-scenario changes require
every host lane: `claude-live` (both configured matrix legs, including merged
mode), `codex-live`, and `pi-live`; the legacy Claude adapter fixture also runs
offline. All required deployments must be approved and green.

Before merge, run the workflow's detached adversarial audit on a throwaway
checkout. At minimum, mutate away the non-team-name projection, replace Codex
`followup_task` with queue-only `send_message`, treat a completed target as dead,
shift the threshold to 1,501, permit two routes when both triggers fire, move the
commit before response disposition/fallback, falsely label semantic rejection as
unavailability, and remove the fallback note. Each
claim-breaking mutation must turn its owning test red. Structural contractlint
may guard the deleted workflow-specific shared prose, but prose-grep is never the
sole proof.

## Out of scope

Changing standing-teammate spawn timing, making prose polish load-bearing beyond this workflow, redesigning the comm officer's prose rules, or implementing a general routing-policy registry.

## Stage Report: ideation

- DONE: Choose the smallest existing mechanism that moves prose-polish policy from the shared contract to this workflow and exposes it to the producing writer.
  Chose the existing ideation stage definition plus comm-officer `Routing Usage`, projected through the existing dispatch fetch composition; live `show-standing` output confirmed the parser/renderer seam.
- DONE: Define qualifying content, required or best-effort routing, bounded fallback, and exclusions without inventing a registry or policy engine.
  Defined staff-review-complex or 1,500-word qualification, a required pre-commit attempt, best-effort response, immediate-unavailable/two-minute-silent fallback, and explicit exclusions.
- DONE: Produce behavior-first acceptance criteria and tests from the c6 RED case, then route the captain-facing draft through the comm officer before commit.
  ACs compare c6's zero pre-commit routes with one ordered route and exercise exclusions/fallback; `/root/comm_officer` remained pending for two minutes, so the required route used the declared fallback and the draft proceeded unchanged.

### Summary

Ideation assigns qualification and commit ordering to this workflow, detailed caller policy to its comm-officer mod, and only generic mod projection to dispatch composition. The design reuses current seams, supplies executable ordering and fallback proofs, and records the bounded comm-officer fallback without changing technical meaning.

## Stage Report: ideation (cycle 2)

- DONE: Choose the smallest existing mechanism that moves prose-polish policy from the shared contract to this workflow and exposes it to the producing writer.
  Repaired the design around the existing stage scope-notes, mod projection, and a generic `teammate.request` runtime binding; no registry, schema, or policy engine is proposed.
- DONE: Define qualifying content, required or best-effort routing, bounded fallback, and exclusions without inventing a registry or policy engine.
  Defined the FO complexity decision path, exact narrative/word boundary, single-request rule, response/commit gate, runtime lifecycle, exclusions, and immediate-versus-two-minute fallback.
- DONE: Produce behavior-first acceptance criteria and tests from the c6 RED case, then route the captain-facing draft through the comm officer before commit.
  Replaced testimony-as-proof and fake-clock fixtures with a reproducible RED plus registered live tool-event/state evidence; both comm-officer reviews were incorporated after the initial bounded fallback.

comm-officer unavailable; proceeded unpolished

### Summary

Cycle 2 resolves every independent-review finding with exact process wording, runtime-neutral turn-starting semantics, full host/mode coverage, real-wall-time fallback, and a detached adversarial audit. The c6 zero-route claim is now labeled testimony until the registered pre-change run creates durable RED evidence; commits are proven to follow either an accepted response or a genuine fallback.
