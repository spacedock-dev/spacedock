---
title: Walk the assembled durable-decision journey before the pre-release
status: ideation
source: "Captain request on 2026-07-26 to exercise the sprint Definition of Done end to end and surface cross-member seams before release."
started: 2026-07-26T10:51:13Z
completed:
verdict:
score: 1.0
worktree:
issue:
sprint: durable-decisions
id: ph0zv6azcrhcxmg57wwnxah7
gates:
    version: 1
    current:
        gate: gate:ph0zv6azcrhcxmg57wwnxah7:backlog
    records:
        - id: gate:ph0zv6azcrhcxmg57wwnxah7:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:ph0zv6azcrhcxmg57wwnxah7-backlog-1
              briefing:
                id: briefing:ph0zv6azcrhcxmg57wwnxah7:backlog:attempt-1:revision-1
                digest: sha256:0782c65c06c7ee9378226b3a7ef88d92939a54c05d916fe3690cc7d99804278f
                digest-domain: canonical-bytes
                request-digest: sha256:77aabae5f9e5af378e377bc1eaefccde931c8932e5e6023a661f5eae4a22e438
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:ph0zv6azcrhcxmg57wwnxah7:backlog:1
                briefing: briefing:ph0zv6azcrhcxmg57wwnxah7:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-26T10:49:10.054082Z"
                decision: approve
                reason: The sprint's end value is an operable durable-decision journey; repeated late integration seams make one real release-candidate walking skeleton necessary before the pre-release.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
---

The sprint has strong component proofs but repeatedly discovered integration seams only
when a First Officer operated the assembled lifecycle. Before the pre-release, drive
one real walking skeleton through the actual launcher, binary, First Officer skills,
default chat presentation, and Subspace room-only entry. Preserve both correctness
evidence and honest agent-friction evidence.

This is a one-time release-candidate acceptance journey, not a standing harness, a
second implementation, or a place to repair findings inline. It starts only after
`s4`, `rq`, `v21`, and `withdraw-stale-open-gate-attempt` land. The x9
`git-root-uri-joined-as-relative-path` reproduction is already rq's exact owned
evidence; it is not another implementation member.

## Problem and end value

Component tests do not establish that a Captain can express intent in chat, a Shaping
First Officer can retain the exact decision without applying it, and a separately
cold-booted Commander can discover and spend that authority once. They also do not
establish that the same prepared-room authority survives the real Subspace presentation
boundary after main and state checkouts move.

The end value is one durable, inspectable release-candidate run in which:

- the Captain authors no JSON/YAML and names no gate, attempt, Briefing, room,
  association, provider, or state-commit coordinates;
- the Shaping First Officer records approvals but never calls the Commander's
  `gate consume` verb;
- a cold Commander alone applies the designated pending approval and dispatches its
  successor;
- default chat and `/subspace:r gate <room>` both close through the same recorder
  authority; and
- every intended state boundary is reconstructible from committed ticket and room
  bytes after another cold boot.

Success is not a command-count result. It is the final on-disk state plus the absence
of five forbidden recoveries: caller-authored metadata, reconstructed coordinates,
direct `gates:` edits, hand-selected Git package commits, and reverted accidental
transitions.

## Exact release-candidate boundary and hard stops

Ideation was inspected at Spacedock main
`50f8d1fb7b0cbc40747622d9f9d95467a0bec6c0`. The ambient launcher reported
`spacedock 0.27.0-pre0`, the adjacent skills reported `0.27.0-pre1`, and its top-level
help exposed no `gate` command. The installed `subspace-tui` reported
`0.10.0-beta.6`. This composition intentionally does **not** satisfy the run
preconditions and is not release evidence.

Immediately before the walk, record these full 40-hex values in the run table:

| Pin | Required identity |
|---|---|
| `spacedock_rc` | Captain-selected clean Spacedock commit containing the landed merge commits for s4, rq's Spacedock half, v21, and withdrawal |
| `spacedock_state_start` | Clean, pushed `spacedock-state/dev` commit before the first journey mutation |
| `subspace_rc` | Pinned clean Subspace commit containing rq's provider half and its real moved-root E2E |
| `fo_gate_skill` / `present_gate_skill` | Git blob ids from `spacedock_rc`, loaded by the actual launcher rather than an ambient plugin |
| `subspace_r_skill` | Git blob id from `subspace_rc`, loaded by the host that handles `/subspace:r` |
| `spacedock_bin` / `subspace_bin` | SHA-256 of binaries built from those exact clean commits |

Build the Spacedock binary from a detached checkout at `spacedock_rc`, pass its
absolute path as `SPACEDOCK_BIN`, and launch fresh Shaping and Commander sessions
through that same binary with the adjacent checkout's plugin. Build/select the
Subspace binary and skill from `subspace_rc`. An installed binary, skill, or TUI with
only a matching version string is not an acceptable substitute for the recorded Git
and file identities.

Stop before a Captain sees a review when any condition below is false:

1. The four prerequisite task states are terminal `PASSED`; every named merge commit
   is an ancestor of `spacedock_rc`; rq records `subspace_rc`; and both repositories'
   required tests are green on those exact commits.
2. One fresh `${SPACEDOCK_BIN} gate --help` exposes `prepare`, `record`, `validate`,
   `eligibility`, `consume`, `--question`, `--artifact`, `--summary`, `--reference`,
   `--room`, `--decision`, `--actor`, `--reason`, and `--withdraw`. The withdrawal
   grammar is exactly:

   ```text
   ${SPACEDOCK_BIN} gate record ENTITY --withdraw \
     --actor agent:first-officer --reason TEXT --workflow-dir docs/dev
   ```

3. `spacedock state ready` succeeds from the clean acceptance checkout, and its state
   head equals `spacedock_state_start`; a rebase conflict is an immediate halt.
4. A real, folder-form second task has a current prepared request-backed attempt,
   a concrete scope change that makes its Briefing stale before decision, and a
   committed replacement Artifact. If no such review exists, ask the Captain to name
   real work later; never fabricate a decision or a synthetic task to satisfy the run.
5. The exact Subspace skill and binary pass their landed availability/capability
   preflight before provider side effects. After a selected Subspace launch, failure
   never falls back to chat.
6. Both selected Git-root commits and paths exist in the local main/state object
   stores. Missing objects are a closed failure, not permission to fetch, deepen,
   reconstruct a path, or use working-tree bytes.

Any product defect ends the walking skeleton at the last durable boundary and routes
to its owner. This task records it; it changes no product file.

## Minimum role-separated journey

The actual acceptance checkout is a clean clone/worktree at a different absolute main
and state location from the authoring checkout. The Shaping First Officer receives
only ordinary-language Captain intent:

> Present this task's ideation review in chat. If I approve, record my decision for the
> Commander and stop without applying it. The separate review on `{rescope-target}` is
> stale because `{truthful scope change}`; withdraw it, restart cleanly, prepare the
> replacement, and present that replacement through Subspace.

The FO uses the real `fo-gate-lifecycle` and `present-gate` skills. The command examples
below are the exact product entries the evidence table records; the Captain does not
type or assemble them.

### Track A — chat records authority and stops

1. Cold-boot the Shaping FO through the pinned launcher. Run `state ready`, then
   `status --boot --identify --json`; load `fo-gate-lifecycle` at gate engagement and
   `present-gate` at presentation.
2. Prepare this folder-form task with:

   ```text
   ${SPACEDOCK_BIN} gate prepare durable-decisions-release-walking-skeleton \
     --question "Is this release-candidate walking skeleton ready for Commander execution?" \
     --artifact docs/dev/.spacedock-state/durable-decisions-release-walking-skeleton/index.md \
     --summary "One real role-separated durable-decision acceptance journey." \
     --reference docs/roadmap/durable-decisions/index.md \
     --workflow-dir docs/dev
   ```

   Use the emitted absolute `room=` verbatim. `spacedock state commit
   durable-decisions-release-walking-skeleton` must commit the binding and exactly two
   prepare-time metadata files before presentation.
3. Default chat emits one root review naming the bound Briefing id/digest. After the
   Captain personally approves, record:

   ```text
   ${SPACEDOCK_BIN} gate record durable-decisions-release-walking-skeleton \
     --decision approve --actor person:captain \
     --reason "Captain approved the presented release-candidate journey." \
     --workflow-dir docs/dev
   ```

   Commit through `spacedock state commit`. Require the task to remain at `ideation`
   with `approved-awaiting-advance`. The Shaping FO stops. Calling `gate consume` here
   is agent-role friction, not evidence that the product lacks an application verb.

### Track B — truthful re-scope and the real Subspace entry

1. On the separate real task, retain the old room tree digest, record the truthful
   withdrawal with the exact `--withdraw` entry above, and run `spacedock state commit
   {rescope-target}`. A repeated withdrawal must exit nonzero with unchanged task and
   room bytes.
2. End that Shaping session. From a newly cold-booted Shaping session in the relocated
   clean checkout, `status --boot --identify --json` must surface
   `withdrawn-awaiting-prepare`. Invoking `/subspace:r gate {old-room}` at this point
   is the stale-attempt mutant: it must refuse before display or provider mutation.
3. Prepare attempt N+1 with the same `gate prepare` surface, selecting the committed
   replacement task Markdown as the one Artifact and
   `docs/roadmap/durable-decisions/index.md` as the one Reference. Use a new truthful
   question and exact summary. Commit the emitted room with `state commit`; do not
   infer its attempt number or pathname.
4. Invoke exactly:

   ```text
   /subspace:r gate <emitted-room>
   ```

   The real TUI must display the canonical question, the exact Artifact summary and
   bytes, and the Reference bytes resolved from the pinned Git objects. After the
   Captain renders the truthful Result, record only through:

   ```text
   ${SPACEDOCK_BIN} gate record {rescope-target} \
     --room <emitted-room> --workflow-dir docs/dev
   ```

   Commit through `state commit`. The Shaping FO stops after recording. If the Result
   is approve, its application remains pending for a Commander; the Shaping FO does
   not consume it. A prior Subspace-side consume is classified as operator-role
   friction and routed to the sprint-role owner, not as a missing recorder capability.

### Track C — cold Commander applies once

Start a new top-level Commander session through the same pinned launcher; do not resume
the Shaping session. Boot must discover this task as `approved-awaiting-advance`.
The Commander invokes:

```text
${SPACEDOCK_BIN} gate consume durable-decisions-release-walking-skeleton \
  --workflow-dir docs/dev
${SPACEDOCK_BIN} state commit durable-decisions-release-walking-skeleton \
  --workflow-dir docs/dev
```

Require `condition=approved-pending`, `eligible=true`, `consumed=true`, and successor
`implementation`. Immediately repeat the consume once: it must exit nonzero, keep the
post-consume task bytes and state commit unchanged, and grant no second dispatch
authorization. The Commander then uses the ordinary first-officer dispatch path once
to dispatch this task's implementation worker. The resulting worker handle, dispatch
file digest, implementation status, and later stage report—not a shell invocation
count—prove the one successor effect.

## Evidence table and end-value measurements

Implementation appends one `## Release-candidate run` section to this file. Its first
table has exactly these columns:

| UTC timestamp | Actor intent | Exact command or skill entry | Pre-state | Exit/output | Post-state | Retained artifact/digest | Elapsed | Friction |
|---|---|---|---|---|---|---|---|---|

It contains rows for pinning/preflight; both cold boots; chat prepare/commit/present/
record/commit; withdrawal/commit/repeat; stale-room refusal; replacement prepare/
commit; `/subspace:r`; room record/commit; Commander consume/commit/repeat; successor
dispatch; clean-checkout reboot; and every focused owner test below. Raw transcripts
may be cited by immutable digest, but no new event collector or log protocol is added.

The report then records these independent outcome comparisons:

- **Operator recovery baseline:** all five forbidden recovery counters remain zero.
  Any hand-authored metadata, coordinate reconstruction, direct gate edit, manual
  package commit, or transition revert makes AC-1 fail even if every command exits 0.
- **Identity baseline:** compare canonical Briefing id/digest, request digest,
  Artifact/Reference URI plus raw `rev`, Result digest, and presented-inventory digest
  from prepared bytes, recorder projection, and `gate validate`. One divergent value
  fails AC-2/AC-3.
- **Lifecycle baseline:** the chat close changes no status; only the cold Commander's
  consume changes `ideation` to `implementation`; repeat consume changes no bytes and
  grants no successor; one dispatched worker reaches the declared successor.
- **Durability baseline:** for each lifecycle boundary, compare the parent and child
  state Git trees. The target task/room paths are complete, a deliberately dirty
  sibling sentinel is absent from every commit, historical rooms retain their
  pre-withdraw tree digests, and a fresh checkout at the final pushed state commit
  reproduces `gate validate` plus boot readiness.
- **Provider surface baseline:** prepare-time room inventory is exactly
  `gate-briefing.json` and `request.json`; after presentation, durable additions are
  under `provider/`; `association.json`, copied Artifact/Reference payloads, and
  retained `resolved-sources` are absent on success and every catchable exit.

These resulting bytes, transitions, and clean cold-boot observations are the release
claim. The number of commands or green tests alone is not.

## Focused failure mutants

Run the three lifecycle mutants in the live journey: stale old-room presentation,
repeated withdrawal, and repeated consume. Capture exit code, actionable stderr,
before/after task hash, before/after room tree hash, provider-path existence, and
state Git head.

Run digest, missing-object, and catchable-provider mutants from the exact landed
owner suites at `spacedock_rc` and `subspace_rc`, inline in the same report:

- s4's focused prepared-authority tests must cover tampered request/Briefing/source
  digest and unavailable Git object before mutation. The expected landed anchors are
  `TestPreparedAuthorityIsRecomputedDuringReadOnlyValidation` and
  `TestRecordBriefingRejectsUnavailablePreparedGitSourceWithoutChangingTree`;
  record the final landed names rather than assuming these branch names survived.
- rq's real room-only suite must cover canonical Briefing/manifest mismatch, missing
  local Git object, catchable provider failure/interrupt diagnostics retention, and
  deletion of ephemeral `resolved-sources` before completion. x9's relative-path
  reproduction remains evidence owned by rq; do not reproduce it as a new task.
- withdrawal's focused test must prove stale identity, actor/reason mutations, and
  repeat withdrawal are byte-clean; the live repeat is the composition check.
- the existing consume tests must prove repeat/stale consume cannot transition or
  authorize again; the live repeat is the role-separated composition check.

If a prerequisite lands without a runnable owner test for its assigned mutant, stop
and route that evidence gap to the owner. Do not create a local fake provider, wrapper,
fixture framework, retry loop, or test script in this task.

## Friction disposition and ownership

Every nonempty `Friction` cell receives one class, one owner, and one next action:

| Class | Meaning in this run | Required disposition |
|---|---|---|
| `address now` | Supported release journey cannot reach an AC, including invalid or lossy durable state | Stop at the durable boundary; route to the owning product task and require a corrected exact RC before rerun |
| `deferred with promotion condition` | Real risk outside the promised journey | Record trigger, supported-path evidence, and the condition that promotes it to material |
| `polish` | End value remains intact but wording, discoverability, or latency is needlessly rough | Route to the owning skill/command task; it does not block by itself |
| `needs decision` | Fix would change authority, scope, retained evidence, or release promise | Stop and ask the Captain; do not choose or repair inline |

Owner routing is singular:

- room preparation, source identity, canonical Briefing, recorder validation, and
  two-file/association rules → s4;
- Git-root materialization, TUI/provider lifecycle, summaries, Result/inventory
  evidence, cleanup, and x9 → rq and its pinned Subspace counterpart;
- source-build/plugin compatibility identity → v21;
- truthful withdrawal/current-attempt boot projection → withdrawal;
- Shaping-versus-Commander role misuse, including a Shaping FO consuming a recorded
  Subspace approval → `first-officer-sprint-shaping-commander-lifecycle`;
- pending-approval cold-boot discovery/consume handoff →
  `cold-commander-pilot-gate-consumption-seam`;
- `state commit` path scope/sibling leakage →
  `state-commit-folder-entity-scope` (`vn`); and
- an ambiguous product-policy tradeoff → Captain as `needs decision`.

The final report distinguishes command defect, provider defect, agent-role error,
Captain choice, and test-infrastructure failure. “The agent used the wrong verb” is not
automatically a product capability gap; “the only documented role path instructed the
wrong verb” is routed to the role-contract owner.

## Mechanism-to-value ledger

| Mechanism | Serves | Simpler alternative considered | Why insufficient |
|---|---|---|---|
| Full commit, skill-blob, and binary pins | AC-1, AC-2, AC-3 | Record version strings | Equal strings do not prove the launcher, skills, and binaries came from the tested candidate bytes |
| One table appended to this task | AC-5 | Add an event logger or sidecar report | A logger is a new harness and a sidecar splits the durable decision from its evidence |
| Real current task plus one legitimately stale review | AC-1, AC-3 | Create synthetic acceptance tasks | Synthetic work can hide discovery, authority, and role friction present on real workflow state |
| Live stale/repeat mutants plus landed owner suites | AC-2 | Reimplement every mutant in a new E2E script | A second suite duplicates dependency semantics and can disagree with the product owners' canonical validators |
| Dirty sibling sentinel plus clean final checkout | AC-4 | Trust `state commit` exit zero | Exit zero alone cannot prove full room inclusion, sibling exclusion, remote durability, or cold-boot readability |

## Expected surface and tolerance

This is an artifact-only acceptance task.

- **Product code, tests, skills, and docs/site:** 0 files, 0 LOC.
- **Human-authored retained surface:** this
  `durable-decisions-release-walking-skeleton/index.md` only, with the implementation
  run section, friction ledger, and stage reports; up to 350 added lines after the
  approved ideation body.
- **Product-generated state:** gate frontmatter changes on this task and the one real
  re-scope target; at most two newly prepared rooms with exactly two metadata files
  each at prepare time; one provider-owned evidence subtree on the Subspace room; and
  path-scoped state commits for each mutation boundary.
- **Transient validation state:** one dirty sibling sentinel and disposable clean
  checkouts only; neither remains in a retained commit or final working tree.

No report sidecar, orchestration script, fixture, fake provider, generic harness,
event schema, or product fix is permitted. A third durable task, more than two new
rooms, any product-file diff, or more than 350 post-ideation report lines requires a
design reset rather than silent scope growth.

No documentation diff is proposed: the task changes no user-visible behavior. If the
walk discovers a documentation defect, its owner receives the concrete before/after
wording as a routed finding.

## Acceptance criteria

**AC-1 (VALUE)** — The Captain, Shaping First Officer, and cold Commander complete the
supported natural-language journey with zero forbidden recoveries.

Verified by the timestamped live table plus final state: zero caller-authored JSON/YAML,
zero reconstructed room coordinates, zero direct `gates:` edits, zero manual package
commits, and zero reverted transitions; this task reaches `implementation` with one
real successor worker. Any one forbidden recovery or absence of the successor outcome
fails, even when command/test counts are green.

**AC-2 (CORRECTNESS)** — Retained state proves exact identity, truthful retirement,
record-versus-apply separation, and one-use authorization.

Verified by recomputing the Briefing/request/Artifact/Reference/Result/inventory
identity table, old-room tree digest, and task state from the final clean checkout.
Changing any retained byte, allowing a withdrawn/consumed attempt to act again, or
letting chat/room record advance status fails.

**AC-3 (CHANNEL PARITY)** — Default chat and real Subspace reach the same
provider-neutral recorder boundary without sharing presentation mechanics.

Verified by both tracks producing recorder-valid current-attempt bindings and durable
`gate record` closures, while only the Subspace room contains provider evidence.
Adding provider authority to chat, changing lifecycle semantics under the override,
or requiring an association/payload copy fails.

**AC-4 (DURABILITY)** — Every lifecycle boundary is one path-scoped state commit that
survives a clean cold boot without sibling leakage.

Verified by parent/child Git-tree comparisons and a fresh checkout at the pushed final
state commit. Omitting any changed target/room path, including the seeded sibling,
mutating a historical room, or failing boot/validation from the clean checkout fails.

**AC-5 (OPERABILITY)** — The retained report makes every step, delay, and friction
actionable without hiding role errors as product defects.

Verified by the required nine-column table and friction ledger: each nonempty friction
has one class, owner, next action, and promotion condition when deferred. A missing
owner, absent elapsed time, ambiguous command/skill entry, or command-count-only
conclusion fails.

## Test plan and proof order

1. **Dependency/pin preflight (low cost, CLI/read-only):** resolve exact commits,
   binaries, skill blobs, merge ancestry, task states, `state ready`, and one fresh
   gate help. Stop on any mismatch.
2. **Smallest live boundary first (medium, live workflow):** prepare this task, commit
   its exact two-file room, present in default chat, record Captain approval, and
   observe `approved-awaiting-advance` with unchanged status. If no visible retained
   proof appears within 90 minutes, stop for architecture review; add no controller.
3. **Re-scope/Subspace boundary (high, live TUI):** withdraw the real stale attempt,
   run repeat/stale-room mutants, cold-boot, prepare N+1 from relocated clean roots,
   invoke the actual room-only skill, record its binding Result, and stop the Shaping
   FO after commit.
4. **Commander boundary (medium, separate live workflow):** cold-boot, consume this
   task once, commit, prove repeat consume byte-clean, and dispatch the actual
   implementation successor once.
5. **Owner mutants (medium, existing focused tests):** run the exact s4, rq,
   withdrawal, and consume tests on the pinned commits. No new test source.
6. **Durability/end-value verification (medium, Git/CLI):** seed sibling dirt, compare
   each boundary tree, clone the final state cleanly, rerun boot and `gate validate`,
   inspect room inventories, and calculate the five forbidden-recovery outcomes.
7. **Repository release checks (existing release gate):** on `spacedock_rc`, run
   `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`; run rq's
   declared Subspace repository gates on `subspace_rc`. Formatting must leave both
   exact candidate trees unchanged.

The composition-level runtime handoff is deliberately the deliverable being tested,
so there is no smaller throwaway composition spike and no second journey. Its
component mechanisms already carry owner tests and s4/rq moved-root spikes; x9 is
folded into rq. The live default-chat boundary is implementation's first falsifiable
exercise, before Subspace or Captain application.

## Out of scope

- Fixing, documenting, or working around any friction in this task.
- Creating synthetic gate work merely to make the walk pass.
- Changing gate authority, adding dispatch receipts, or claiming effect-level
  exactly-once semantics; only the authorization is one-use.
- Fetching missing Git objects, adding compatibility flags, retrying a provider room,
  or falling back to chat after provider launch.
- Cutting or tagging the release. This report is evidence for the later pre-release
  decision.

## Stage Report: ideation

- DONE: Define the minimum executable role-separated chat and Subspace journey against exact release-candidate revisions.
  AC-1 and AC-3 design evidence: the body pins code/state/skill/binary identities, uses the real launcher and command surfaces, stops Shaping FOs after record, and reserves consume plus successor dispatch for a cold Commander.
- DONE: Specify the evidence table, focused failure mutants, dependency stop conditions, and friction disposition without a new harness.
  AC-2, AC-4, and AC-5 design evidence: the nine-column run table, live repeat/stale mutants, owner-suite digest/object/provider mutants, hard stops, and singular owner map are executable without new code or scripts.
- DONE: Declare the expected artifact-only surface and prove the run checks end value rather than command counts.
  AC-1, AC-4, and AC-5 design evidence: the surface is one human-authored state file plus product-generated rooms/state; outcomes measure forbidden recoveries, identities, transitions, Git trees, clean boot, and routed friction.
- SKIPPED: Execute the release-candidate journey during ideation.
  Current main `50f8d1fb7b0cbc40747622d9f9d95467a0bec6c0` lacks the complete gate surface and the four declared dependencies have not all landed; executing now would be false release evidence.

### Summary

The task is now a one-time, artifact-only release acceptance runbook. It composes
default chat, truthful withdrawal, real room-only Subspace, and Commander-only
consumption against exact post-land commits, with byte/state evidence and explicit
owner routing instead of inline fixes or a new harness.
