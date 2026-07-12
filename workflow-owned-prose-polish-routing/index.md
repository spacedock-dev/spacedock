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
The shared layer should own only standing-teammate discovery, lifecycle,
addressing, and projection.

The split caused a live failure in c6. The FO loaded the generic dispatch
contract but did not pass the workflow/mod guidance to the ideation writer. The
writer made no polish route before the task body reached roughly 9,200 words;
`comm-officer` arrived only for late cleanup. The RED baseline is therefore an
observable ordering failure: zero polish requests before the first state commit,
despite a qualifying complex ideation.

## Proposed direction

Use the existing stage-definition and mod-projection seams. Do not add a routing
registry, policy schema, or routing engine.

### Ownership and trigger

The `docs/dev/README.md` `ideation` stage owns this workflow's qualification and
write ordering. A task body qualifies when either:

- the ideation meets the stage's existing staff-review complexity rule; or
- its body, excluding YAML frontmatter and stage reports, reaches 1,500 words.

For a qualifying body, the writer MUST make one polish request at the final
draft boundary and before the first entity-body/state commit. "MUST" applies to
the attempt and its ordering, not to receiving or accepting a reply. Any returned
edit remains behavior-preserving: the writer rejects changes to technical
meaning, qualifiers, criteria, or evidence.

The comm-officer mod owns the prose-polish purpose and caller contract: the
target teammate, modes, two-minute bound, exclusions, and fallback. A response
is best-effort and non-load-bearing. If the teammate is absent, rejects the
request, or has not replied within two minutes, the writer commits the original
draft and records `comm-officer unavailable; proceeded unpolished` in the stage
report. No retry or replacement agent is required.

The exclusions are direct captain chat, a non-complex ideation body below 1,500
words, short operational status, tool output, logs, commit messages, and prose
that cannot be separated from secrets or sensitive context. Excluded content
causes no polish request and needs no fallback note.

### Propagation seam

`dispatch build` already composes two fetches into a worker assignment:
`show-stage-def` supplies the workflow stage and `show-standing` projects each
declared mod's `## Routing Usage`. Preserve that shape and close its narrow host
gap: every non-bare named worker for a workflow with a standing mod receives the
standing-mod fetch, including Codex workers without a legacy Claude `team_name`.
The emitted section stays generic and tells the runtime to use its teammate
message channel; all prose-specific semantics come verbatim from this workflow's
comm-officer mod.

This is the smallest existing mechanism. The stage definition reaches every
ideation writer already, the mod parser and `show-standing` renderer already
exist, and the dispatch file already carries fetch commands. No new policy
representation or lookup path is needed.

### Concrete edits

- In `docs/dev/README.md` under `### ideation`, add the qualification, required
  pre-commit attempt, and reference to the projected comm-officer routing usage
  stated above.
- In `docs/dev/_mods/comm-officer.md` `## Routing Usage`, make the two-minute
  best-effort boundary, exclusions, exact fallback note, and behavior-preserving
  acceptance rule the single caller contract.
- In `internal/dispatch/build.go`, gate the existing `show-standing` fetch on a
  non-bare named dispatch plus declared standing mods, not on legacy
  `team_name` alone. Reuse `EnumerateDeclaredStandingTeammates`; add no registry.
- In `internal/claudeteam/standing.go`, replace the shared header's
  prose-specific `2-minute` / `un-polished` defaults with generic direction to
  follow each projected mod's routing usage through the runtime teammate-message
  channel.
- Delete `Routing through a standing prose-polisher` from
  `skills/first-officer/references/fo-dispatch-core.md`. Keep its surrounding
  generic standing-injection and dispatch mechanics unchanged.

No public CLI syntax changes, so the docs site needs no command-reference edit.
The workflow README and mod diff above are the user-visible process
documentation change.

### Mechanism check

No implementation spike is needed. This live Codex dispatch demonstrated that
`show-stage-def` reaches the writer while `show-standing` is missing; running
`spacedock dispatch show-standing --workflow-dir docs/dev` demonstrates that the
existing parser/renderer emits the comm-officer routing usage. The unproven part
is only the current build conditional, which the first host-neutral dispatch
fixture below exercises before implementation.

## Acceptance criteria

**AC-1 (VALUE — early route): A c6-shaped complex ideation produces one comm-officer request before its first state commit, improving the RED baseline from zero pre-commit requests to one.**

Verified by: a registered host-neutral workflow fixture records ordered
`route_requested` and `state_commit` events and asserts exactly one route whose
sequence number is lower than the first commit's; the original c6 fixture
records zero pre-commit routes and remains the negative baseline.

**AC-2 (writer receives policy): Every non-bare named ideation dispatch for this workflow supplies both its stage definition and the declared comm-officer mod's routing usage, including a Codex request with no legacy `team_name`.**

Verified by: a Go dispatch-package test builds Claude-team and Codex/no-team
requests from the same fixture workflow, executes every emitted fetch command,
and asserts the composed worker input contains the fixture stage marker and
fixture mod-routing marker sourced independently from the two fixture files.

**AC-3 (bounded fallback): An unavailable comm officer cannot block the workflow; a silent comm officer cannot block it past two minutes; the original draft commits and the stage report records the declared fallback.**

Verified by: the workflow fixture observes immediate continuation for an
unavailable target and uses a fake clock for a silent target; both arms produce
a successful commit, unchanged body hash, and the durable fallback line, while
the silent arm advances exactly two virtual minutes.

**AC-4 (no unnecessary routing): A non-complex ideation below 1,500 words and each explicitly excluded content class produce zero comm-officer requests and no fallback note.**

Verified by: a table-driven workflow fixture feeds the short-body boundary
(1,499 words), direct chat, operational status, tool output, log, commit message,
and sensitive-context cases; its message-channel recorder stays empty while the
permitted output or commit succeeds.

**AC-5 (narrow ownership): Shared contract output remains generic while this workflow and its comm-officer mod determine prose purpose, qualification, timeout, exclusions, and fallback.**

Verified by: the AC-1 through AC-4 behavior fixtures are run once with the dev
workflow/mod and once with a workflow that declares no prose-polish mod; only the
dev fixture routes. `internal/contractlint` may additionally enforce structural
absence of the obsolete named shared subsection, but that check is not the
behavior proof.

**AC-6 (no new abstraction): The implementation changes the existing stage definition, mod projection, dispatch conditional, and tests without adding a routing registry, policy schema, or generalized engine.**

Verified by: validation reviews the product diff and pairs that structural fact
with AC-2's executable dispatch proof; adding a new policy package, registry
type, or config format fails this criterion even if tests pass.

## Test plan

Start with the host-neutral dispatch-package test because it exercises the only
composition seam that failed in this session. Use fixture markers rather than
production wording, and run the emitted fetch commands instead of grepping the
dispatch file. Estimated cost is small: one table-driven Go test beside the
existing standing/build tests.

Add a registered workflow behavior fixture with an event recorder, fake clock,
and temporary split-root state checkout. Its qualifying arm proves route-before-
commit ordering; its unavailable arm proves the bounded fallback and committed
note; its short/excluded table proves zero routes. This is medium complexity and
must drive the assignment/write boundary, not merely inspect instruction prose.

Run focused dispatch and workflow-fixture tests first, then `gofmt -w ./cmd
./internal`, `go test ./...`, and `go test ./... -race`. Because the change
touches shipped host-neutral dispatch and skill-contract surfaces, validation
also runs the registered host lanes required by the workflow proof policy. A
structural contractlint assertion may ensure the obsolete shared subsection
does not return, but prose-grep is never the sole proof.

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
