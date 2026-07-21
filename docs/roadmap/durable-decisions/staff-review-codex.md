# Durable-decisions — sprint-wide preflight staff review (Codex seat)

## Verdict: NOT READY

The sprint has a sound decomposition around one gates writer, but the assembled package
does not yet deliver its stated promise: a decision recorded by machinery that cannot
mis-file it, presented by machinery that cannot lose it, and consumed exactly once.
Six material findings survive. Three require captain choices because the available
repairs change either the sprint boundary or the meaning of “exactly once.”

Finding count: **6 material, 3 needs-decision, 4 recorded declines/non-blockers**.

## Scope and method

This is an independent sprint-wide refutation, not a per-task design review. I read the
roadmap staff-review contract, the sprint index, the fable seat, the Commander package,
all four member bodies, the shared contract, the frozen current gate packages, the
production evidence, and the 0260 shaping debrief. I checked DoD ownership, dependency
order, shared-region composition, write boundaries, gate provenance, digest
reproducibility, state-checkout durability, and cold-boot instructions.

Cheap falsifiable checks used:

```sh
spacedock status --workflow-dir docs/dev --where sprint=durable-decisions
git -C docs/dev/.spacedock-state status --short
git -C docs/dev/.spacedock-state ls-files \
  gate-review-presentation-command/review/ideation/briefing-3
shasum -a 256 <current briefing and contract snapshots>
jq -cS . <briefing.json> | tr -d '\n' | shasum -a 256
jq '{resolution, status}' <current float-result.json>
rg -n 'gate review|implemented by xb|hand-authored|DEFERRED|dispatch follows' \
  <sprint index, package, member bodies, shared contract>
```

No workflow state or source was mutated. “Confirmed” below means directly observed in
the committed tree or current state checkout. “Impact” is the sprint-wide inference.

## Material findings

### Material 1 — two current gate approvals are not durable, binding captain approvals

**Confirmed.** The state checkout reports these current xb room artifacts as untracked:

```text
?? gate-review-presentation-command/review/ideation/briefing-3/briefing.review.jsonl
?? gate-review-presentation-command/review/ideation/briefing-3/float-result.json
```

`git -C docs/dev/.spacedock-state ls-files` confirms that the briefing, entity snapshot,
contract snapshot, and summary are tracked, but the review log and provider result are
not. A cold checkout therefore cannot recover the provider evidence for xb's current
approval. This directly refutes the index's claim that every briefing is verifiable in
its room and the fable seat's “closure complete” state-status check.

The tracked 02av provider result and the current untracked xb provider result both say:

```json
{"status":"advisory","binding":false,"actor":"person:reviewer","approver":"agent:invoking-session","decision":"approve"}
```

Their entity records copy the resulting `person:reviewer` Resolutions into closed gate
attempts with pending `advance` applications. The Commander package then calls every
such record “the captain's approval.” No durable field or note shows who externally
authorized either advisory result to become binding, and neither record identifies the
captain as the renderer under the sprint's recording-identity ruling.

**Impact.** A cold Commander is told to consume two applications whose durable evidence
proves only an advisory reviewer approval, not a binding captain approval. Tracking xb's
room files is necessary but insufficient; both gates need a durable authorization and
identity repair. Re-record or supersede them under the actor who actually rendered the
binding decision, preserving the provider result as evidence rather than silently
promoting `binding:false`.

### Material 2 — the recorded digest convention contradicts the one shared contract

**Confirmed.** The contract and 3k body define `briefing.digest` as SHA-256 over RFC 8785
canonical Briefing bytes (`gate-resolution-frontmatter-contract.md`, “Round records”; 3k
“Resolved storage decisions”). The current 02av, h1, and xb gate fields instead equal the
raw, pretty-printed `briefing.json` file hashes. Re-serializing the same JSON with sorted,
compact keys changes every hash. For example, h1's recorded/raw digest is
`f98f7ac3...dbc78`, while the normalized candidate is `3bb7e3ed...e603`. The current 3k
attempt is farther from the contract: its `briefing.digest` (`fd95df2a...e5f1a`) is the
raw Markdown hash of `entity-snapshot.md`, not a canonical Review & Gate Briefing object.

The byte checks in the fable seat proved that each recorded string matches *some file*;
they did not prove the normative digest algorithm. The records are reproducible as raw
blob pins, but they are not the JCS Briefing digests that the proposed recorder says it
will validate and replay.

**Impact.** A strict implementation must either reject the sprint's own replay fixtures
or weaken the approved contract. The advisory-digest hole remains open at the semantic
level even though the raw bytes are committed. Before implementation, define a versioned
digest domain (for example, raw snapshot blob versus canonical Briefing), migrate or
grandfather historical records explicitly, and add one fixture that changes only JSON
formatting: a JCS digest must remain stable while a raw-blob digest must not.

### Material 3 — the “one spec” still specifies the architecture xb's approved gate forbids

**Confirmed.** The live shared contract (`b4c96236...17b`) and the contract snapshot in
xb's approved briefing-3 have the same hash. That spec:

- instructs the ensign to run `${SPACEDOCK_BIN:-spacedock} gate review ...`;
- assigns provider-id normalization to xb;
- puts `gate review` in the Go helper boundary and behavioral proof.

The other artifact in the same approved package, xb's entity snapshot, says the opposite:

- no `spacedock gate review` verb;
- zero Subspace coupling in the Spacedock binary;
- result validation and id normalization implemented by the recorder;
- a provider-owned override script supplies presentation.

The sprint index adopts xb's later architecture, and the Commander package says the
shared contract is the design authority while mentioning only an owner-tag amendment.
The contract landing pass removes shaping tokens; it does not reconcile these operative
sections.

**Impact.** The current approval package approved two mutually exclusive implementations,
and a cold Commander cannot know which text wins. This is a shared-region collision and a
failure of the one-spec boundary, not stale historical prose. Amend every operative
contract section to the chosen architecture, route the change to its owners, and create
a superseding gate attempt for the affected owner as the sprint's own change protocol
requires.

### Material 4 — xb's load-bearing implementation has no in-sprint owner

**Confirmed.** Sprint Goal bullet 3 requires one command that retains the result, log,
and diagnostics on success and failure. xb says the hardened override script and its
committed CI drive suite are load-bearing and must live in the sibling Subspace repo.
The sprint index makes direct sibling-repo edits and Subspace product work out of scope.
The Commander package says to route the suite as a cross-repo ask and “do not absorb it
here.” No member, pinned sibling commit, dependency gate, or release criterion owns that
ask through completion.

**Impact.** xb can land prose and a recorder request in this repo while the only mechanism
that satisfies its value AC remains absent. The sprint would then cut 0.27.0 while its
presentation guarantee exists only as an external request and a scratch spike. Either
bring a pinned, tested provider artifact into the sprint's owned dependency graph, add an
explicit cross-repo member/gate, or narrow the sprint DoD and defer xb. “Route the ask” is
not a cold-boot execution plan.

### Material 5 — 02av's release line is hand-authored because its recorder is deferred

**Confirmed.** Goal bullet 5 promises a recorded advisory decline in the sprint whose
headline says decisions are recorded by machinery that cannot mis-file them. 02av's body
and the shared contract explicitly defer the room append, frontmatter pointer, and body
projection beyond the first implementation. Its AC-1 therefore substitutes a
hand-authored room record and hand-authored `### Feedback Cycles` projection. The
Commander package repeats that interim and states that no in-sprint member builds the
plumbing.

The fable seat found the missing owner, then treated “hand-authored interim” as the fold
that closed it. That changes the promised mechanism rather than closing the finding.

**Impact.** The fifth DoD bullet can prove triage semantics and a zero-line diff, but it
cannot prove the sprint's durability claim. Hand authoring is the failure mode 3k exists
to remove. Add the minimum advisory-round recorder surface to an in-scope owner, or state
plainly that bullet 5 is a convention-only follow-up outside this sprint and remove it
from the machine-durability success criterion.

### Material 6 — h1 proves one authorization consume, not exactly-once dispatch

**Confirmed.** h1 correctly observes that current dispatch has no transactional
idempotency guard. Its repair atomically writes `status: <target>` and
`application.state: consumed`, then dispatches as a separate step. On dispatch failure,
the retry path dispatches again even though the application is already consumed. The
design deliberately carries no dispatch identity, receipt, or recovery state.

This conflicts with the shared contract, which says the application becomes consumed
only in the durable state change that records the transition/dispatch machinery's
successful outcome. It also leaves two untested crash windows:

1. crash after `consumed` is committed but before dispatch starts;
2. crash after dispatch succeeds but before the caller durably observes success.

h1's three-pass fake proves that a consumed application is not admitted again. It cannot
prove that the external dispatch effect occurs exactly once across either window.

**Impact.** The sprint may honestly claim “the gate authorization is consumed once,” or
it may claim “the worker dispatch happens exactly once,” but the current design cannot
claim both. Resolve the semantic fork, align the shared contract and Goal, and add
response-loss/crash fixtures. If dispatch exactly-once remains promised, the design needs
a durable effect identity, receipt, or an already-proven idempotent downstream operation.

## Needs-decision

### Needs-decision 1 — what does “consumed exactly once” cover?

Choose one contract before h1 starts:

- **Authorization-only:** atomically consume the gate application and advance status;
  document dispatch as separately retryable and potentially at-least-once.
- **End-to-end effect:** retain the current wording, but add durable dispatch identity and
  reconciliation sufficient for response loss.

The first is smaller and consistent with h1's implementation; the second preserves the
headline promise but expands scope.

### Needs-decision 2 — is the Subspace provider artifact part of the 0.27.0 release gate?

If yes, name its owner, repository, pinned revision, required CI evidence, landing order,
and what blocks the Spacedock tag. If no, defer xb or narrow Goal bullet 3 to the chat
channel. The present package attempts to exclude the work and depend on it simultaneously.

### Needs-decision 3 — how are historical raw digests and advisory approvals migrated?

Choose whether the recorder accepts a versioned legacy raw-blob digest and explicit
advisory-to-binding adoption record, or whether the four current gate histories are
superseded under the final JCS/identity contract. Silent reinterpretation would make the
records look cleaner while losing honest history.

## Recorded declines and non-blockers

### Decline 1 — blocker evaluation and hold authoring need not ship now

h1's decline is evidence-based: no dry-run gate carried a blocker or execution hold, and
the empty-blocker eligibility read is fail-closed. **Promote when:** a real gate must
survive a session boundary while a machine, rather than FO judgment, waits on a named
predicate. At promotion, authoring and satisfaction evaluation must land together.

### Decline 2 — the four-member carve and leading dependency are correct

The membership query returns exactly 3k, 02av, h1, and xb. 3k must lead because every
other member consumes its schema or verbs; h1 and the in-repo half of xb may then proceed
in parallel. **Promote when:** the repaired shared contract changes a consumer-visible
field or command after a dependent gate; that owner must re-anchor through a superseding
attempt rather than rely on the current drift waiver.

### Decline 3 — the one-writer boundary is conceptually sound

No member proposes an independent `gates:` writer: 3k owns recording, h1 extends the same
binary, xb hands results to it, and 02av stays advisory. **Promote when:** an override
script, hand-authored round path, or convenience command begins mutating entity gates
directly. That would be a second writer and a release blocker.

### Decline 4 — landing-pass ownership and raw artifact commitment mostly held

The Commander now owns removal of contract shaping scaffolding, and the current frozen
briefing/snapshot files needed to reproduce raw digest strings are tracked except for
xb's current review log and provider result identified above. **Promote when:** the
pre-cut audit still finds task ids/owner tags in the landed spec, a render regression, or
any current room artifact outside `git ls-files`.

## What held under refutation

- The sprint has one coherent product center: durable gate records, not a general
  scheduler, provenance system, panel, or event ledger.
- The four entity bodies declare useful expected surfaces, tolerance, red fixtures, and
  hard self-checks. Their riskiest local mechanisms were exercised early.
- Current gate pointers and cross-attempt pending-application state are internally
  consistent after the fable fold; each current gate has one pending advance.
- Raw briefing and snapshot hashes are byte-reproducible from present files. The defect is
  the mismatch between that raw convention and the normative JCS domain, not missing
  bytes across the board.
- h1's blocker decline has a concrete promotion condition, and fail-closed handling is a
  safe interim.
- The scope exclusions for panels, provenance enforcement, route-context projection, and
  the broader event design remain proportionate.

## Closure and cold-boot readiness condition

The sprint becomes ready only when all six material findings close in durable state:

1. commit xb's current room result/log and repair or supersede xb and 02av approvals with
   binding authorization and honest renderer identity;
2. settle the digest domain and prove replay plus formatting-only drift with fixtures;
3. reconcile the shared contract with xb's selected architecture and re-present affected
   owner gates under the sprint's change protocol;
4. give the Subspace script/suite a named, pinned owner and release gate, or remove that
   external mechanism from this sprint's DoD;
5. give advisory round recording an in-sprint machine owner, or narrow/defer the 02av
   machine-durability promise; and
6. align h1's exactly-once claim with a crash-safe proof and the shared lifecycle text.

After those folds, update the Commander package to read **both** staff-review seats and
the reconciled spec. Verify from clean main and state checkouts: the four-member query,
zero state-checkout status, every current room artifact in `git ls-files`, normative
digest recomputation, binding approval provenance, the eight-history replay, the two h1
crash windows, and the xb retention suite at its pinned provider revision. Until that
evidence exists, do not consume the pending advances or cut 0.27.0.
