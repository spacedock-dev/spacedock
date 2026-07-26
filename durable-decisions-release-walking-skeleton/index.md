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
second implementation, or a place to repair findings inline. It starts only after the
six prerequisite sprint members `gqs`, `mn`, `0m6`, `rq`, `s4`, and `v21` are
terminal and landed. The x9 `git-root-uri-joined-as-relative-path` reproduction is
already rq's exact owned evidence; it is not another implementation member.

## Problem and end value

Component tests do not establish that a Captain can express intent in chat, a Shaping
First Officer can retain the exact decision without applying it, and a separately
cold-booted Commander can discover and spend that authority once. They also do not
establish that the same prepared-room authority survives the real Subspace presentation
boundary after main and state checkouts move.

The end value is one durable, inspectable release-candidate run in which:

- the Captain authors no JSON/YAML and names no gate, attempt, Briefing, room,
  association, provider, or state-commit coordinates;
- one acceptance worker is dispatched before the journey, receives exact logs and
  digests through its normal assignment channel, and alone retains the immutable
  evidence plus final Stage Report;
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
direct `gates:` edits, manual Git commits of generated gate/source packages, and
reverted accidental transitions.

This role split is a candidate requirement, not current supported behavior. The
walking skeleton must not call it supported until `mn`
`respect-gate-application-authority` lands and proves the narrow First Officer
contract: a Shaping FO commits the recorded decision and stops; a fresh Commander
consumes, commits, and dispatches.

This Cycle 2 correction governs where it conflicts with the retained first ideation
report: it replaces the four-prerequisite readiness claim, post-journey evidence
worker, authoritative-state mutants, five-member package assumption, expected surface,
acceptance criteria, and proof order. The first report remains historical evidence of
the rejected design.

## Exact release-candidate boundary and hard stops

Ideation was inspected at Spacedock main
`50f8d1fb7b0cbc40747622d9f9d95467a0bec6c0`. The ambient launcher reported
`spacedock 0.27.0-pre0`, the adjacent skills reported `0.27.0-pre1`, and its top-level
help exposed no `gate` command. The installed `subspace-tui` reported
`0.10.0-beta.6`. This composition intentionally does **not** satisfy the run
preconditions and is not release evidence.

At that main commit, the roadmap index blob is
`de338245b8a5bb91e9dad4d7ecacbfc20724ff7c`; it correctly delegates membership to the
state query. The Commander-package blob
`b851f7600439286aa246198da74dcb8b17f98f74` still describes the retired five-member
drive. The authoritative state at
`006b6512e8d3c185949bd3666e2ecfe7add4e435` instead returns the seven current members
`gqs`, `ph`, `rq`, `mn`, `0m6`, `s4`, and `v21`. The old package blob is a stale
baseline, not an acceptable run pin.

Immediately before the walk, record these full 40-hex values in the run table:

| Pin | Required identity |
|---|---|
| `spacedock_rc` | Captain-selected clean Spacedock commit containing gqs, s4, rq's Spacedock half, v21, 0m6 withdrawal, and mn's narrow role-boundary contract |
| `spacedock_state_start` | Clean, pushed `spacedock-state/dev` commit before the first journey mutation |
| `subspace_rc` | Pinned clean Subspace commit containing rq's provider half and its real moved-root E2E |
| `fo_gate_skill` / `present_gate_skill` / `role_contract` | Git blob ids from `spacedock_rc`, loaded by the actual launcher rather than an ambient plugin |
| `role_contract_owner` | `mn` id/slug, terminal state commit, and merged product commit |
| `subspace_r_skill` | Git blob id from `subspace_rc`, loaded by the host that handles `/subspace:r` |
| `spacedock_bin` / `subspace_bin` | SHA-256 of binaries built from those exact clean commits |
| `commander_package` | Git blob id of the reconciled seven-member `docs/roadmap/durable-decisions/dispatch-sprint-execution.md` at `spacedock_rc` |
| `acceptance_state` | Captain-named state branch/checkout plus exact starting commit and the two exclusively owned real task slugs |

Build the Spacedock binary from a detached checkout at `spacedock_rc`, pass its
absolute path as `SPACEDOCK_BIN`, and launch fresh Shaping and Commander sessions
through that same binary with the adjacent checkout's plugin. Build/select the
Subspace binary and skill from `subspace_rc`. An installed binary, skill, or TUI with
only a matching version string is not an acceptable substitute for the recorded Git
and file identities.

Stop before a Captain sees a review when any condition below is false:

1. All six prerequisite sprint members `gqs`, `mn`, `0m6`, `rq`, `s4`, and `v21`
   are terminal `PASSED`; every named product merge is an ancestor of `spacedock_rc`;
   rq records `subspace_rc`; and no prerequisite remains only approved, in validation,
   or present in an unmerged worktree.
2. `mn`'s product commit is merged, and its proof shows Shaping
   record/commit/stop followed by cold Commander consume/commit/dispatch.
   Current `fo-gate-lifecycle` behavior that consumes immediately after approval is
   insufficient; until this prerequisite lands, the role-separated journey is not a
   supported release claim.
3. The Commander package is reconciled to the roadmap's authoritative current
   seven-member query and its exact blob is pinned. At ideation correction time the
   live query is exactly `gqs`, `ph`, `rq`, `mn`, `0m6`, `s4`, and `v21`, while the
   package still says five and enumerates the retired `3k`, `vn`, `h1`, `02av`, and
   `xb` cohort. A stale member count, landing order, responsibility table, or close-out
   condition halts.
4. Candidate identity, all prerequisite owner suites, and both repositories' required
   full/race/format checks have run first on the exact pins and are green. These checks
   precede every live gate mutation.
5. One fresh `${SPACEDOCK_BIN} gate --help` exposes `prepare`, `record`, `validate`,
   `eligibility`, `consume`, `--question`, `--artifact`, `--summary`, `--reference`,
   `--room`, `--decision`, `--actor`, `--reason`, and `--withdraw`. The withdrawal
   grammar is exactly:

   ```text
   ${SPACEDOCK_BIN} gate record ENTITY --withdraw \
     --actor agent:first-officer --reason TEXT --workflow-dir docs/dev
   ```

6. `spacedock state ready` succeeds from the Captain-named clean acceptance checkout,
   no other FO/Commander/worker owns either selected task, and its state
   head equals `spacedock_state_start`; a rebase conflict is an immediate halt.
7. The two Captain-named tasks are real and folder-form. The chat target is ready for
   a nonterminal approval; the re-scope target has a current prepared request-backed
   attempt,
   a concrete scope change that makes its Briefing stale before decision, and a
   committed replacement Artifact. If no such review exists, ask the Captain to name
   real work later; never fabricate a decision or a synthetic task to satisfy the run.
8. The exact Subspace skill and binary pass their landed availability/capability
   preflight before provider side effects. After a selected Subspace launch, failure
   never falls back to chat.
9. Both selected Git-root commits and paths exist in the local main/state object
   stores. Missing objects are a closed failure, not permission to fetch, deepen,
   reconstruct a path, or use working-tree bytes.
10. The acceptance worker has been dispatched through the ordinary entered-stage
    worker path. Its initial assignment names the six terminal prerequisites, exact
    starting commits, one-file package scope, evidence directory, and final Stage
    Report checklist. After it commits the package correction and the proof-first
    checks pass, an ordinary follow-up binds the exact candidate pins, exclusive task
    scope, command/session log paths, and SHA-256 digests. It is addressable before
    the first `gate prepare` and has no authority to call `gate record`, `gate
    consume`, or edit gate frontmatter.

Any product defect ends the walking skeleton at the last durable boundary and routes
to its owner. This task records it and changes no product code, tests, or skills. Its
sole repo-authored implementation change is the Commander-package reconciliation.

### Seven-member Commander-package reconciliation

The roadmap remains strategy and its query remains membership authority; the package
must nevertheless stop making false cardinality and retired-cohort claims. Its
release-candidate revision must say that the query returns seven and order the current
work as:

1. s4 and v21 land their independent Spacedock prerequisites;
2. 0m6 and gqs land the truthful re-scope and entered-stage dispatch boundaries;
3. rq lands after s4 with its pinned Subspace counterpart and x9 folded into its
   evidence;
4. mn lands the narrow Shaping-records/stops,
   Commander-consumes/dispatches contract; and
5. ph is the final assembled acceptance member, dispatched only after all preceding
   conditions are terminal.

The reconciliation removes no history, revives no retired member, and does not make
the package another tracker. During implementation, the ph acceptance worker first
applies this correction to exactly
`docs/roadmap/durable-decisions/dispatch-sprint-execution.md`, commits it, and pins its
resulting blob before any journey mutation.

## Minimum role-separated journey

After the six prerequisite members pass, a cold Commander consumes ph's approved
ideation gate, commits entry into implementation, and—under gqs—immediately dispatches
the ph acceptance worker. This setup transition is not counted as journey proof. The
worker applies the one-file Commander-package reconciliation, commits it, pins the
resulting candidate, and runs the owner/repository checks before any live gate
mutation. Its normal follow-up assignment then binds the exact candidate pins,
Captain-named exclusive state/tasks, intended positive steps, command/session log
paths and digests, evidence directory, and final Stage Report checklist. Beyond that
one package artifact it retains evidence only; FO/Captain actors retain all authority
to prepare, record, consume, commit lifecycle state, and dispatch target work.

The actual acceptance checkout is clean and at a different absolute main/state
location from the authoring checkout. No other session mutates the named chat or
re-scope task until handback. The Shaping First Officer receives only ordinary-language
Captain intent:

> Present `{chat-target}` in chat. If I approve, record my decision for the Commander
> and stop without applying it. The separate review on `{rescope-target}` is stale
> because `{truthful scope change}`; withdraw it, restart cleanly, prepare the
> replacement, and present that replacement through Subspace.

The FO uses the real `fo-gate-lifecycle` and `present-gate` skills. The command
examples below are the exact product entries sent to the already-running acceptance
worker with command/session log paths and SHA-256 digests. The Captain does not type
or assemble them, and the worker does not execute their authoritative mutations.

### Track A — chat records authority and stops

1. Cold-boot the Shaping FO through the pinned launcher. Run `state ready`, then
   `status --boot --identify --json`; load `fo-gate-lifecycle` at gate engagement and
   `present-gate` at presentation.
2. Prepare the Captain-named folder-form chat target with:

   ```text
   ${SPACEDOCK_BIN} gate prepare {chat-target} \
     --question "{Captain-approved real decision question}" \
     --artifact {committed-chat-target-review.md} \
     --summary "{exact concise Artifact summary}" \
     --reference docs/roadmap/durable-decisions/index.md \
     --workflow-dir docs/dev
   ```

   Use the emitted absolute `room=` verbatim. `spacedock state commit
   {chat-target}` must commit the binding and exactly two prepare-time metadata files
   before presentation.
3. Default chat emits one root review naming the bound Briefing id/digest. After the
   Captain personally approves, record:

   ```text
   ${SPACEDOCK_BIN} gate record {chat-target} \
     --decision approve --actor person:captain \
     --reason "{Captain's truthful approval reason}" \
     --workflow-dir docs/dev
   ```

   Commit through `spacedock state commit`. Require the task to remain at its actual
   gated stage with `approved-awaiting-advance`. The Shaping FO stops. Calling
   `gate consume` here is a failed role-contract proof and agent-role friction, not
   evidence that the product lacks an application verb.
4. The Shaping FO sends the acceptance worker the exact root-session log, command log,
   emitted room, resulting state commit, and SHA-256 for each file through ordinary
   worker messaging. The worker verifies the bytes before retaining them and does not
   infer missing coordinates from prose.

### Track B — truthful re-scope and the real Subspace entry

1. On the separate real task, retain the old room tree digest, record the truthful
   withdrawal with the exact `--withdraw` entry above, and run `spacedock state commit
   {rescope-target}`. The authoritative state receives no repeated-withdrawal probe.
2. End that Shaping session. From a newly cold-booted Shaping session in the relocated
   clean checkout, `status --boot --identify --json` must surface
   `withdrawn-awaiting-prepare`.
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
5. The Shaping FO passes the exact command log, complete Subspace session log,
   provider-result paths, state commit, and all SHA-256 digests to the acceptance
   worker through its existing assignment channel. The worker verifies and retains
   those bytes without a watcher, collector, wrapper, or retry protocol.

### Track C — cold Commander applies once

Start a new top-level Commander session through the same pinned launcher; do not resume
the Shaping session. Boot must discover `{chat-target}` as
`approved-awaiting-advance`.
The Commander invokes:

```text
${SPACEDOCK_BIN} gate consume {chat-target} \
  --workflow-dir docs/dev
${SPACEDOCK_BIN} state commit {chat-target} \
  --workflow-dir docs/dev
```

Require `condition=approved-pending`, `eligible=true`, `consumed=true`, and successor
equal to the target's declared nonterminal next stage. The Commander then uses the
ordinary first-officer dispatch path to dispatch that entered stage before any later
transition. The target worker handle, dispatch-file digest, entered-stage status, and
eventual matching Stage Report—not a shell invocation count—prove the successor
effect and gqs composition.

The Commander sends its exact boot/consume/commit/dispatch command log, root-session
log, dispatch assignment, and all digests to the acceptance worker. The acceptance
worker authors no further repo surface beyond immutable evidence and appends ph's
final Stage Report after all positive state, disposable negative probes, clean-boot
validation, and handback checks finish.

### Exclusive-state handback

Positive state is never reset to its pre-run form. After evidence is committed, the
acceptance worker reports the exact final pushed state commit and target statuses to
the Captain and owning FO. Handback succeeds only when:

- both authoritative task/room trees validate and are clean at the reported commit;
- the chat target's entered-stage worker and the re-scope target's truthful
  post-Result owner are named;
- no disposable negative branch/clone was pushed or merged;
- every disposable path has been removed after its exact logs were retained and
  hashed; and
- the Captain/owning FO explicitly releases the temporary exclusivity hold and resumes
  normal workflow ownership from the reported state, without deleting rooms or
  rewriting history.

## Evidence table and end-value measurements

The pre-dispatched acceptance worker appends one `## Release-candidate run` section and
the final implementation Stage Report to this file. Its first table has exactly these
columns:

| UTC timestamp | Actor intent | Exact command or skill entry | Pre-state | Exit/output | Post-state | Retained artifact/digest | Elapsed | Friction |
|---|---|---|---|---|---|---|---|---|

It contains rows for the seven-member query and reconciled Commander-package blob;
candidate/skill/binary pins; prerequisite terminal states; role-contract proof; owner
suites and repository checks; acceptance-worker dispatch/assignment; both cold boots;
chat prepare/commit/present/record/commit; withdrawal/commit; replacement prepare/
commit; `/subspace:r`; room record/commit; Commander consume/commit/entered-stage
dispatch; disposable negative probes; clean-checkout reboot; and handback.

Each FO/Commander sends the worker its exact command log and host session log path plus
SHA-256 through ordinary worker messaging. The worker verifies the supplied digest,
retains the exact bytes under this task's evidence directory, and records their new
Git blob id and SHA-256. Missing bytes or a digest mismatch is an evidence defect and
a stop. There is no watcher, collector, event schema, wrapper, daemon, or polling
protocol; normal assignment/follow-up delivery is the handoff.

The report then records these independent outcome comparisons:

- **Operator recovery baseline:** all five forbidden recovery counters remain zero.
  Any hand-authored metadata, coordinate reconstruction, direct gate edit, manual
  Git commit of a generated gate/source package, or transition revert makes AC-1 fail
  even if every command exits 0.
- **Identity baseline:** compare canonical Briefing id/digest, request digest,
  Artifact/Reference URI plus raw `rev`, Result digest, and presented-inventory digest
  from prepared bytes, recorder projection, and `gate validate`. One divergent value
  fails AC-3/AC-4.
- **Lifecycle baseline:** the chat close changes no status; only the cold Commander's
  consume enters the declared working stage; gqs then exposes and dispatches that
  stage before any later transition. The worker assignment and matching Stage Report
  prove the successor effect; repeat behavior is tested only on a disposable branch.
- **Durability baseline:** for each lifecycle boundary, compare the parent and child
  state Git trees. The target task/room paths are complete, a deliberately dirty
  sibling in a disposable clone is excluded, authoritative historical rooms retain
  their pre-withdraw tree digests, and a fresh checkout at the final pushed state
  commit reproduces `gate validate` plus boot readiness.
- **Provider surface baseline:** prepare-time room inventory is exactly
  `gate-briefing.json` and `request.json`; after presentation, durable additions are
  under `provider/`; `association.json`, copied Artifact/Reference payloads, and
  retained `resolved-sources` are absent on success and every catchable exit.

These resulting bytes, transitions, and clean cold-boot observations are the release
claim. The number of commands or green tests alone is not.

## Focused failure mutants

First, after the worker establishes exact candidate identity and before live state
mutation, run the exact landed owner suites at `spacedock_rc` and `subspace_rc`:

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
  repeat withdrawal are byte-clean.
- the existing consume tests must prove repeat/stale consume cannot transition or
  authorize again; gqs tests must prove consumed entry dispatches before transition.

If a prerequisite lands without a runnable owner test for its assigned mutant, stop
and route that evidence gap to the owner.

Second, after each positive boundary supplies exact real bytes, create a disposable
local clone or branch at that exact `spacedock_rc` plus state commit. Never point it at
the authoritative state remote; never push or merge it. The role-authorized FO or
Commander—not the evidence worker—runs the matching public command there, and passes
the exact logs/digests to the worker:

| Disposable clone source | Mutant | Required result |
|---|---|---|
| post-prepare chat room | tampered request/Briefing/source digest; missing local Git object | fail before presentation/mutation; task/room authority byte-clean |
| post-withdraw re-scope task | repeated withdrawal | nonzero actionable refusal; no byte or state-head change |
| post-N+1 re-scope room | old-room/stale-attempt presentation | fail before provider allocation/display; historical room unchanged |
| post-Subspace room | catchable provider failure/interrupt | diagnostics retained; `resolved-sources` removed; no fabricated Result |
| post-consume chat target | repeated consume and dirty sibling | nonzero/no second authorization; target bytes unchanged; sibling excluded from any target commit |

The clone must contain the exact retained real room/task bytes and local objects, not a
hand-built approximation. The worker verifies clone base commits before accepting a
result. After evidence bytes are retained and hashed, remove the disposable clone.
Do not create a local fake provider, wrapper, fixture framework, retry loop, or test
script in this task.

## Friction disposition and ownership

Every nonempty `Friction` cell receives one class, one owner, and one next action:

| Class | Meaning in this run | Required disposition |
|---|---|---|
| `address now` | Pinned candidate journey cannot reach an AC, including invalid or lossy durable state | Stop at the durable boundary; route to the owning product task and require a corrected exact RC before rerun |
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
  Subspace approval → mn's narrow role-boundary First Officer contract;
- consumed-stage projection or failure to dispatch before transition → gqs;
- stale seven-member count, landing order, or responsibility/close-out text → ph's
  one-file Commander-package reconciliation;
- `state commit` path scope/sibling leakage →
  `state-commit-folder-entity-scope` (`vn`); and
- an ambiguous product-policy tradeoff → Captain as `needs decision`.

The final report distinguishes command defect, provider defect, agent-role error,
Captain choice, and test-infrastructure failure. “The agent used the wrong verb” is not
automatically a product capability gap; “the only documented role path instructed the
wrong verb” is routed to mn.

## Mechanism-to-value ledger

| Mechanism | Serves | Simpler alternative considered | Why insufficient |
|---|---|---|---|
| Full commit, skill-blob, and binary pins | AC-1, AC-2, AC-3 | Record version strings | Equal strings do not prove the launcher, skills, and binaries came from the tested candidate bytes |
| Reconciled seven-member Commander-package blob pin | AC-1, AC-5 | Trust the live query alone | The Commander package currently contradicts the query and would drive a retired cohort/order |
| Pre-dispatched evidence worker using normal assignment | AC-2, AC-5 | Reconstruct evidence after the journey | Retrospective prose cannot bind exact session/command bytes or prove who held mutation authority |
| Exact retained logs plus SHA-256, no collector | AC-2, AC-5 | Add an event logger | A logger is a new harness and changes the runtime being accepted |
| Captain-named exclusively owned real tasks | AC-1, AC-3, AC-4 | Create synthetic acceptance tasks | Synthetic work can hide discovery, authority, concurrency, and handback friction present on real workflow state |
| Owner suites first, then disposable exact-byte clone mutants | AC-2, AC-4 | Mutate the positive journey or build a new E2E script | Positive state must remain truthful; a second suite can diverge from owner validators |
| Clean final checkout and explicit handback | AC-4, AC-5 | Trust `state commit` exit zero | Exit zero cannot prove remote durability, ownership release, or cold-boot readability |

## Expected surface and tolerance

This remains an artifact-only acceptance task, with one explicit expansion from the
rejected design: ph owns the one-file Commander-package correction needed to make the
seven-member release walk truthful.

- **Product code, tests, skills, user-facing docs, and site changed by ph:** 0 files,
  0 LOC. The mn role contract is an independently owned landed prerequisite.
- **Roadmap package changed by ph:** exactly
  `docs/roadmap/durable-decisions/dispatch-sprint-execution.md`; expected at most 100
  touched lines and 40 net added lines to replace retired membership, landing order,
  responsibility, and close-out prose with the seven-member composition. The roadmap
  index and every other roadmap/package file remain byte-identical.
- **Human-authored retained surface:** this
  `durable-decisions-release-walking-skeleton/index.md` plus one
  `evidence/<run-id>/` directory containing the exact in-scope command/session logs
  delivered to the acceptance worker and one digest inventory. Expected maximum:
  12 evidence files, 25 MiB total, and 250 added report/Stage Report lines after the
  approved correction.
- **Product-generated state:** ph's setup consume/dispatch state plus gate frontmatter
  changes on the two Captain-named real targets; at most two newly prepared rooms with
  exactly two metadata files each at prepare time; one provider-owned evidence subtree
  on the Subspace room; and path-scoped state commits for each mutation boundary.
- **Transient validation state:** disposable exact-byte clones/branches only; no
  negative clone, dirty sibling, or mutation survives in an authoritative branch,
  remote, or final working tree.

No orchestration script, fixture, fake provider, generic harness, event schema,
collector, watcher, or product fix is permitted. A third positive target, more than
two new authoritative rooms, more than 12 evidence files/25 MiB, any ph-owned product
diff, a second roadmap/package file, more than the declared package-line tolerance, or
more than 250 post-correction report lines requires a design reset rather than silent
scope growth.

The one Commander-package diff changes execution guidance, not user-visible behavior.
It lands under ph before the walk; mn's narrow role contract lands under mn. Any new
documentation defect discovered by the walk is routed with concrete before/after
wording.

## Acceptance criteria

**AC-1 (VALUE)** — After the narrow role contract lands, the Captain, Shaping First
Officer, and cold Commander complete the natural-language journey on exclusively owned
real state with a pre-dispatched evidence worker and zero forbidden recoveries.

Verified by exact retained actor logs, the worker dispatch/assignment timestamp, final
task/room state, and handback: zero caller-authored JSON/YAML, reconstructed room
coordinates, direct `gates:` edits, manual Git commits of generated gate/source
packages, or reverted transitions.
Dispatching the evidence worker after the first journey mutation, mutating outside the
exclusive scope, or lacking the entered-stage successor worker fails even if all
commands/tests are green.

**AC-2 (READINESS)** — The walked candidate is exactly the reconciled seven-member
release composition, with every owner prerequisite and repository gate complete before
live mutation.

Verified by equality between the seven-item status query and pinned Commander-package
meaning; terminal `PASSED` plus merge ancestry for gqs, mn, 0m6, rq, s4, and v21;
mn's pinned role-contract proof; candidate/skill/binary digests; owner-suite logs; and
full/race/format results. A stale five-member package, nonterminal prerequisite,
unmerged owner commit, version-only identity, or check timestamp after the first
journey mutation fails.

**AC-3 (CORRECTNESS)** — Retained state proves exact identity, truthful retirement,
record-versus-apply separation, and one-use authorization.

Verified by recomputing the Briefing/request/Artifact/Reference/Result/inventory
identity table, old-room tree digest, and task state from the final clean checkout.
Changing any retained byte, allowing a withdrawn/consumed attempt to act again, or
letting chat/room record advance status fails.

**AC-4 (PARITY AND DURABILITY)** — Default chat and real Subspace reach the same
provider-neutral recorder boundary, and every positive lifecycle boundary survives a
clean cold boot without negative-probe contamination or sibling leakage.

Verified by both tracks producing recorder-valid current-attempt bindings and durable
`gate record` closures, while only Subspace retains provider evidence; by parent/child
Git-tree comparisons and a fresh checkout at the pushed final state; and by clone-base
proof that every negative result ran off-authority. Adding provider authority to chat,
requiring association/payload copies, omitting a target path, mutating a historical
room, or placing a mutant/sibling in authoritative history fails.

**AC-5 (OPERABILITY)** — The retained report makes every step, delay, and friction
actionable, binds evidence to its actor, and hands exclusive state back without hiding
role errors as product defects.

Verified by the required nine-column table, exact retained command/session log plus
digest inventory, friction ledger, worker-authored final Stage Report, and explicit
handback record. Each friction has one class, owner, next action, and promotion
condition when deferred. A digest mismatch, FO-authored substitute report, missing
owner/elapsed time, ambiguous command/skill entry, unreleased exclusivity hold, or
command-count-only conclusion fails.

## Test plan and proof order

1. **Dispatch and candidate identity (low, package edit plus CLI/read-only):** verify
   the six terminal prerequisite states and mn proof, consume ph into implementation,
   commit, and let gqs dispatch the acceptance worker. The worker reconciles and
   commits the one-file seven-member Commander package, then pins its blob with the
   exact query, commits/binaries/skill blobs, merge ancestry, `state ready`, exclusive
   task ownership, and one fresh gate help. Stop on any mismatch.
2. **Owner suites and repository gates first (medium):** run exact gqs, s4, rq,
   withdrawal, consume, v21, and mn role-contract owner suites on the pinned commits. On
   `spacedock_rc`, run
   `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`; run rq's
   declared Subspace repository gates on `subspace_rc`. Formatting must leave both
   candidate trees unchanged. Finish this step before any live mutation.
3. **Bind immutable evidence (low, ordinary follow-up):** give the already-dispatched
   worker exact candidate pins, exclusive scope, log paths, and digests through its
   normal assignment channel. Verify it is addressable, owns no gate authority, and
   will author only evidence plus the final Stage Report from this point before the
   first positive `gate prepare`.
4. **Positive chat boundary (medium, live workflow):** on the exclusive chat target,
   prepare/commit, present in default chat, record/commit Captain approval, pass exact
   logs/digests to the worker, and stop the Shaping FO.
5. **Positive re-scope/Subspace boundary (high, live TUI):** on the exclusive re-scope
   target, withdraw/commit, cold-boot, prepare/commit N+1 from relocated clean roots,
   invoke the actual room-only skill, record/commit its Result, pass exact logs/digests,
   and stop the Shaping FO.
6. **Positive Commander boundary (medium, separate live workflow):** cold-boot,
   consume/commit the chat target, dispatch its entered stage under gqs, and pass exact
   boot/command/session/dispatch logs and digests to the worker.
7. **Disposable exact-byte mutants (medium):** from each corresponding positive state
   commit, run the digest/object, repeated-withdrawal, stale-room, provider-failure,
   repeated-consume, and dirty-sibling probes only in unpushed disposable clones/
   branches. Retain logs/digests, then remove them.
8. **Durability, report, and handback (medium, Git/CLI):** compare boundary trees,
   clone final state cleanly, rerun boot and `gate validate`, inspect inventories,
   compute forbidden-recovery outcomes, have the same worker append the final Stage
   Report, commit evidence, and complete explicit ownership handback.

The composition-level runtime handoff is deliberately the deliverable being tested,
so there is no smaller throwaway composition spike and no second journey. Its
component mechanisms already carry owner tests and s4/rq moved-root spikes; x9 is
folded into rq. Candidate identity, owner suites, and repository gates run first; the
live default-chat boundary is the first falsifiable journey exercise only after the
evidence worker is bound.

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

## Stage Report: ideation (cycle 2)

- DONE: Add terminal prerequisites for gqs, mn, 0m6, rq, s4, and v21.
  AC-1/AC-2 evidence: every owner must be terminal PASSED and merged; the role split remains explicitly unsupported until mn proves Shaping-record/stop and cold-Commander consume/dispatch.
- DONE: Reconcile and pin the current seven-member durable-decisions Commander package.
  AC-2 evidence: the stale five-member package blob and exact seven-member state-query baseline are recorded, with corrected landing order and the future release-candidate package blob required before the walk.
- DONE: Put candidate identity, owner suites, and repository checks before live mutation.
  AC-2 evidence: after entered-stage dispatch lets the worker create the one-file package candidate, proof order stops on pin/ancestry/suite/full/race/format failure before any acceptance state changes.
- DONE: Dispatch the acceptance worker before the journey and bind immutable evidence through normal assignment.
  AC-1/AC-5 evidence: after its one package correction, the worker verifies exact command/session logs plus SHA-256, retains them, and alone writes ph's final Stage Report; FO/Captain keep gate mutation authority.
- DONE: Keep the positive journey on Captain-named exclusively owned state and define handback.
  AC-1/AC-4/AC-5 evidence: two real tasks remain authoritative and unmutated by probes; final state, downstream owners, cleanup, exclusivity release, and no-history-rewrite handback are explicit.
- DONE: Run negative mutants only in disposable exact-byte clones or branches.
  AC-3/AC-4 evidence: each mutant is based on a corresponding real positive commit, remains unpushed/unmerged, and is removed after logs/digests are retained.
- DONE: Preserve x9 folded into rq and declare the small artifact-only expansion without admitting product code or a harness.
  AC-2/AC-5 evidence: x9 stays rq evidence; ph permits exactly one Commander-package file plus its state report/evidence artifacts and product-generated rooms, with zero ph-owned product LOC or harness.
- SKIPPED: Execute the release-candidate journey during corrected ideation.
  The seven-member package, six prerequisite members, and exact release-candidate pins are not yet all terminal and reconciled; execution now would violate AC-2.

### Summary

Cycle 2 makes readiness and evidence ownership binding rather than retrospective. The
positive run now starts only after the seven-member candidate, including mn's role
contract, is proven, runs under exclusive state ownership with a pre-dispatched
evidence worker, and confines every negative mutation to disposable copies before
explicit handback.
