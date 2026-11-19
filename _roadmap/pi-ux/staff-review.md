# Pi-UX sprint preflight staff review

Date: 2026-10-03

Verdict: **Gaps to close, not yet cold-boot drivable.** Six Material findings
remain. Fold them before presenting the ideation gates. This review does not
approve the proposed restatements, change membership, or authorize implementation.

The seven intended members name all seven outcomes. They do **not** yet form a
closed delivery plan. In particular, the repeated gate outcome has an in-scope
proof owner but no in-scope conduct-repair owner. First-contact discovery also
stops short of the promised no-hand-wiring path.

## Evidence read

Read the sprint discipline, sprint index, reference staff review, all seven
member bodies and their stage reports, and the named implementation seams.

Snapshot:

- Code: `cdfa462d1d426febb2391733507c7fde2143fb0a`, branch `main`.
- State: `5436dde75e7bb507b6576505e4c54dc20acf0d5a`.
- `mc`: ideation commits `2fac0d774`, `38a5f040d`.
- `9w`: `a95ecf5dd`, `5436dde75`.
- `ekw`: `5dc6d9a3d`, `42f40cc33`.
- `mh`: `337d09c8a`, `55bf5c386`.
- `z6e`: `9b3d1ec4b`; `d52`: `4f657242f`.
- `3w1`: retained ideation body; latest path commit `1d7d04a10`.

Also read `docs/runtime-support.md`, `docs/dev/README.md`, the active-owner
registry test, the four binding owners, archived `p17`, external `gc` and `w5`,
and archived `gqsw`. Checked the workflow's four stale TypeScript assertions.
The two historical `9w` root transcripts exist in the named old local worktree;
their SHA-256 values match the entity. They are not portable repository fixtures.

The membership authority remains:

```bash
spacedock status --workflow-dir docs/dev --where sprint=pi-ux
spacedock status --workflow-dir docs/dev \
  --where sprint=pi-ux --where 'sprint-readiness != defer'
```

**Both queries currently return eight members, not seven.** `w5` remains
`backlog`, `sprint: pi-ux`, `group: gate`, `sprint-readiness: ready`. This differs
from the supplied exclusion. This review respects the explicitly intended
seven-member scope; it does not silently add `w5` back. The query discrepancy
must be repaired by the Shaping FO (M5). The query abbreviates `d52` as `d5`;
the full entity ID is `d525n1p5zgnz99hmtjq16z57`.

## Definition-of-done coverage

| Sprint outcome | Intended owner | Sprint-wide result |
| --- | --- | --- |
| Isolated-home discovery without package-root exports or hand-wired extension paths | `3w1` | Nominally owned, incomplete composition. Symlink presence is not package registration or proof that both extensions load (M1). Requires `mc` on the published family. |
| Truthful installed-family doctor verdicts and a working remedy | `mc` | Manifest resolution is owned. The remedy-to-healthy transition lacks an assigned acceptance exercise (M4). |
| A classified stall from archived artifacts, with stated limits | `z6e` | Owned and correctly bounded to the common runner's per-run timeout. Artifact interface needs a shape (M3); live confirmation follows `mh`. |
| Repeated Pi presentation through successor dispatch | `9w` | Measurement is owned; conduct repair is excluded and remains with external `gc`/`w5`. No in-scope member can repair a red journey under the proposed scope (M2). |
| Session-scoped Pi identity and a new session's install offer | `ekw` | Owned through current-session context, not parent environment. Requires actual FO offer evidence, not only hashed-path tests. |
| Published-family pi-live installation and manifest-derived substrate assertions | `mh` | Owned, including all FOUR assertions at 669/671/777/779. Requires `mc`; green lane acceptance also inherits unresolved conduct dependencies (M2). |
| Active owners and classified evidence for every Pi XFAIL binding | `d52` | Owned as attribution/evidence work. Four retained bindings, one re-anchor, zero clearances. Backlog eligibility is distinct from repair readiness; see below. |

**Coverage result:** no outcome lacks a named intended member. However, the gate
outcome lacks an authorized in-scope repair owner. The discovery and remedy
outcomes have incomplete delivery/evidence paths. Thus seven named owners do
not establish seven deliverable outcomes.

The additional format/full/race/live checks in the index are integration gates,
not an eighth product outcome. Their execution belongs in the Commander package.

## Material findings

### M1 — Discovery does not compose with the actual isolated-home launch

`3w1` proposes symlinks under `cleanHome/.pi/agent/npm/node_modules`, while
`PI_CODING_AGENT_DIR` points to a separate `piHome`. Both current fixture
constructors write **only the Spacedock checkout** into `piHome/settings.json`.
The proposal does not register either substrate there.

At `internal/cli/pi.go::runPi`, an unregistered subagents package still causes
an explicit `--extension` argument. `mc` makes that fallback manifest-derived;
it does not make it package discovery. Nothing in the proposed seeding registers
`pi-intercom` for Pi's package loader. Passing its directory and skill checks
cannot prove its extension loaded. The sprint promises both extensions without
hand-wired paths, not merely a doctor-ready filesystem.

The first implementation day also exposes two unresolved seams:

- `newPiSharedLiveDriver` currently creates its clean HOME inline. The proposal
  names `cleanHome` seeding but does not bind that same directory to the env.
- `piDefaultExtensionRoots` is proposed in a **live-tagged** file, but the
  non-live `pi_live_controls_test.go::piIntercomPackageRoot` will call it.
  Ordinary Go tests cannot see that definition. The existing sibling-root test
  also conflicts with replacing sibling lookup by operator-home discovery.

**Required fold:** define one isolated-home setup contract, including real home,
agent directory, clean HOME, package registration, and explicit-override behavior.
Name the non-live helper seam and the changed old test contract. Register/load
both packages through Pi's supported discovery path while preserving the
Spacedock package entry. Keep `mc` as entry-resolution owner. Prove both loaded
extensions with both package-root variables absent and no substrate extension
path supplied by the harness. Do not count the retained explicit fallback as
that proof. Any broader launcher change needs separate scope approval.

### M2 — Gate conduct and green-lane delivery depend on uncommissioned external work

`9w` explicitly restores only the Pi callback grade. It cannot change the Pi
adapter's unconditional post-review stop or any other conduct mechanism. Its
3/3 requirement therefore depends on more than its implementation.

- `gc` is still backlog in `pi-live-completeness`. Its retained evidence records
  a real unbound recorded-gate FAIL. It owns the missing-reference conduct repair
  and receives `9w`'s adapter-stop finding.
- `w5` is explicitly excluded by this review's assignment. It owns exact-digest
  presentation reliability. Its retained seed requires a full displayed digest;
  `9w` permits a canonical compact prefix. The package must distinguish canonical
  authority exactness from rendering, not silently import the older contract.
- `gqsw` is archived done and passed. Reuse its entered-stage dispatch coverage;
  do not schedule a new implementation or treat it as an outstanding blocker.

`mh` requires the unchanged complete pi-live lane to be green. That lane includes
the known unbound recorded-gate failure. Neither a successful package install nor
an owner re-anchor can discharge it. There is no honest seven-member-only order
that guarantees both `mh`'s green lane and `9w`'s repeated completed journey.

**Required fold:** commission an external dependency handoff with accepted code
SHAs, responsible owner, availability, and closure evidence before Commander
acceptance. Unbound failures must be repaired before the first merge requiring
a green Pi lane, including earlier launcher changes. Coordinate `gc` against the
restored `9w` grade; record how the
presentation dependency on `w5` is discharged under the approved compact-prefix
contract. If that handoff is unavailable, the captain must revise scope or the
promised outcome. Do not let the Commander repair conduct inside `9w`, hide it
behind a new binding, or accept a grade-only result as reliability.

### M3 — The stall artifact interface is named but not fully shaped

`z6e` defines the filename, labels, evidence rules, and causal limit well. However,
`classification`, `evidence`, and `limits` have no exact JSON field types or sample
record. Evidence is described as paths, lines, fields, IDs, and missing-input
notices without a defined representation. This is a new operator-facing artifact,
not just a private helper. Choosing its schema during implementation would create
an unreviewed interface.

The runner/classifier boundary also needs an explicit input representation for
read errors, ambiguous roots, malformed JSONL, and unsupported branches. The
existing rejection extractors skip malformed records; reusing that permissive
behavior would defeat the promised inconclusive result.

**Required fold:** specify one concrete output example and field domains, and a
minimal input/result contract. State root-selection and unsupported-branch rules,
and what happens if writing the diagnostic fails. Preserve the original timeout
failure; do not allow summary failure to disguise it. Define where the deterministic
live-tagged wiring check resides: the registry rejects every new live-tagged test
without registration. Use an existing registered owner or explicitly account for
that surface. Keep older/missing fields conservative and keep `mh` before live
confirmation. No generic transcript framework or new capture mechanism is needed.

### M4 — The working-remedy promise has no end-to-end acceptance owner

`mc` owns corrected paths, preserves remedy text, and declares deterministic tests
sufficient. `mh` checks a healthy installed package. Neither plan exercises the
sprint's distinct claim that the remedy can clear the line it accompanies.

An already-correct package printing OK is useful, but is not a missing prerequisite
followed by the printed remedy and an OK recheck. The original incident itself
shows why `pi install` saying “up to date” is not that proof.

**Required fold:** assign this integrated acceptance to `mc`, after the published
family is available. In a disposable home, demonstrate the applicable failing
prerequisite, execute its printed remedy, and show the same line clears. Retain
package versions, command exits, and before/after doctor output. No global install
mutation is authorized by this review. If the existing text cannot work, return
the necessary text/scope change to the gate rather than weakening the outcome.

### M5 — The cold-boot membership and evidence package is not self-contained

The declared exclusion of `w5` is not reflected by either canonical query.
`dispatch-sprint-execution.md` does not exist. That is normal before packaging,
but prevents a “Drivable” verdict now. A cold Commander cannot infer the intended
seven-member set from the authoritative query.

The `9w` negative traces are available only through an absolute old-worktree path.
The body correctly requires portable testdata during implementation, but the
handoff must provide retrievable inputs before a different Commander loses access.
`d52` also needs the explicit state checkout and its SHA: without
`SPACEDOCK_LIVE_STATE_DIR`, its owner-join test skips.

**Required fold:** reconcile membership through FO-owned state mutation, then
package the resulting query, approved restatements, dependency handoffs, artifact
retrieval/provenance, and serial order. Preserve the deliberate exclusion of
`w5` unless the captain changes it. Keep `6h3` and the separate flat-entity /
merge-guard follow-up outside implementation dispatch unless commissioned.

### M6 — Validation must be composed from the diff, not each member's local proof claim

`mc` changes both doctor and launch readiness/fallback in `internal/cli/pi.go`;
its “no live lane needed” statement does not satisfy the repository's changed-live-
surface merge rule. `ekw` edits shared `fo-install.md` as well as the Pi adapter.
Its non-Pi preservation claim needs the affected host lanes, not only Pi API
fixtures. `d52` changes the shared live declarations even though only one Pi owner
ID changes. These are broader validation surfaces than their visible Pi effects.

**Required fold:** the Commander package must map each final diff to required
runtime lanes under `docs/dev/README.md`. Include Pi front-door/common proof for
launcher/harness changes and the required non-Pi lanes for shared loaded surfaces.
Keep detached adversarial validation for the front door, shipped extension/skills,
and CI machinery. Resolve any narrower lane determination explicitly from the
actual diff; do not infer it from a task's “Pi-only” title.

An intentionally capped `z6e` run is expected red diagnostic evidence, separate
from required green regression lanes. `9w`'s declared three-run batch counts every
attempt and skip; do not replace failed samples with selective retries.

## XFAIL ownership audit: what 6h3 does and does not satisfy

The existing active-owner test defines active as outside `_archive`, with status
neither `done` nor `rejected`. It does **not** require a sprint label or ideation.
`6h3` is backlog, unarchived, and nonterminal. Therefore **d52 AC-1 is mechanically
satisfiable with 6h3 in its present state** once the binding is re-anchored. Do not
report an AC-1 failure merely because it is unlabelled or unideated.

That is not evidence of an implementation-ready topology repair. `6h3` has a
bounded problem and criteria but no ideated approach. Its own problem asks for an
active product owner carrying a real approach. The FO must record its scheduling
and design disposition; the Commander must not turn `d52` into the missing repair.
The sprint's attribution outcome does not require fixing every bound journey.

| Binding | Staff disposition |
| --- | --- |
| Owner handoff | Keep `fe7`; it is active by the test but explicitly deferred by priority. Its short body points to retained xp6 evidence; carry retrievable evidence into the reconciliation record. |
| Rejection flow | Re-anchor p17 to `6h3`, not clear. Archived p17's validation explicitly retains `rejection-worker-topology`; its 1499.96-second Go PASS is XFAIL-backed. |
| Keep moving | Keep `x0`; CI `31770740214` records XPASS, but there is no paired exact unbound PASS. |
| Smallest mechanism | Keep `h30`; the Pi driver's call to `claudeMechanismTrace` is an extractor mismatch, not a justified product diagnosis or clearance. |
| Recorded gate | Remain unbound. `gc` owns the ordinary FAIL; neither `9w` nor `d52` can mask it. |

`d52`'s classifications are semantic/infrastructure/evidence dispositions. They
are **not** required to use `z6e`'s timeout labels. There is no mandatory dependency
from `d52`'s owner-ID change to the new stall classifier. Any later clearance needs
a revised disposition and the exact XPASS/unbound PASS pair.

## Cross-member composition and blast radius

| Shared region | Collision and boundary |
| --- | --- |
| `pi_shared_live_runner_test.go` | `3w1` changes construction/env; `z6e` changes the timeout exit; `9w` changes the successful-run grading callback. Land serially and rebase each. A timeout never reaches `9w`'s success callback. |
| `shared_live_runner_test.go` | `d52` changes only binding attribution. It must not absorb either Pi-driver work or shared assertion changes. |
| Pi JSONL readers | `z6e` and `9w` both observe root records and tool IDs. Preserve role/order/error information. Reuse compatible small types only; their diagnostic and authority-grade semantics are different. |
| Package resolution | `3w1` finds roots; `mc` resolves launcher/doctor entries; `mh` validates installed declarations. Agree on the manifest contract without inventing a cross-language resolver framework. `mc` uses the first extension; `mh` checks all entries. Record this deliberate distinction. |
| `docs/runtime-support.md` | `mc` changes path examples; `9w` adds a proof-boundary paragraph. Preserve both in serial landing. |
| `docs/runtime-live-ci.md` | `mh` changes versions; `z6e` adds diagnostic limits. Land `mh` first. |
| Pi adapter / shared install instructions | `ekw` owns install scope only. External `gc` can touch the same Pi adapter for conduct. Rebase and check both bindings survive. |

The main risks are false readiness from filesystem-only discovery, false greens
from measurement-only work, parser leniency hiding missing evidence, and changed
shared instructions affecting non-Pi hosts. The sprint does not need a new runtime,
activity feed, global identity registry, shared gate protocol, or lifecycle
supervisor. Keep those held or out of scope.

## First-day implementability walk

| Member | First implementation work | Decision still needed before dispatch |
| --- | --- | --- |
| `mc` | Add compiled-manifest and missing-target fixtures, then wire the one resolver. | Assign remedy-transition proof and required live validation (M4/M6). The selected manifest fields and fail-closed behavior are already defined. |
| `3w1` | Exercise root discovery and isolated-home setup without environment overrides. | Agent-dir versus HOME registration, intercom loading, helper build-tag placement, and changed sibling-root test contract (M1). |
| `mh` | Add structural mutation guards, then update all three version/integrity pairs and both setup checkpoints. | Which accepted external conduct candidate makes complete-lane acceptance possible (M2). The inline manifest check's inputs and failure semantics are otherwise defined. |
| `z6e` | Add literal artifact fixtures before the timeout-path write. | Exact JSON/input contracts, root ambiguity handling, write failure, and registered wiring-test placement (M3). |
| `ekw` | Extend the existing extension fake; test context A/B/A and hashed paths before tool registration. | No unresolved identity-source decision: the seed is disproved. Package the harmless failing installer and real FO trace, plus non-Pi lane evidence (M6). Missing capability remains manual-only. |
| `9w` | Retain authentic negative traces, add adversarial grade controls, wire the existing callback. | Portable input delivery and external conduct/presentation closure (M2/M5). The callback already has a `liveResult` shape; do not invent another protocol. |
| `d52` | Run the explicit-state RED owner join, change only the Pi owner ID/comment, then re-run and reverse-mutate. | Record 6h3's ownership disposition and retrievable classified evidence. No extractor decision belongs in this task. |

These are sprint-delivery seams, not a second review of every member's ACs.
The largest estimate is `9w`'s +350 to +650 net lines and up to 650 KB of fixtures;
retain its ceiling and redaction/provenance requirements. Other surface bounds
remain the ideation gates' responsibility. No scope expansion is approved here.

## Recommended dispatch sequence

After the findings are folded and the captain approves the gates, use this
**serial implementation/landing order** for the seven members:

```text
mc -> 3w1 -> mh -> z6e -> ekw -> 9w -> d52
```

1. **mc:** unlock published-package launcher/doctor checks. Prepare the disposable
   remedy proof as part of this work.
2. **3w1:** compose discovery with the corrected launcher before other harness
   edits. Prove no-override discovery independently of CI's retained exports.
3. **mh:** refresh all three pins and four assertions. This establishes the
   runtime family for subsequent live proof. Do not call installation success
   full acceptance; hold its green-lane gate if external conduct is unresolved.
4. **z6e:** land after `mh` and rebase after `3w1` in the Pi driver. Verify both
   deterministic timeout wiring and one real capped Pi run on the new family.
5. **ekw:** prove session scope on that same family, with actual A/B/A FO behavior.
   Its isolated TS work can be prepared earlier; serialize the adapter landing
   with any external `gc` edit.
6. **9w:** restore the grade after the Pi-driver setup/timeout changes. Coordinate
   its negative controls with external `gc` before spending the declared batch.
   Require accepted conduct and presentation dependencies for the final 3/3 proof.
7. **d52:** reconcile the assembled registry and actual owner state last. Re-read
   evidence after all preceding changes; keep the four bindings unless a newly
   approved disposition has the required passing pair.

This is not an instruction to merge red required lanes. **Satisfy M2's external
handoff before the first affected merge needing a green Pi lane, not merely before
`mh` or `9w`.** The seven-member order sits after that external prerequisite; it
cannot replace it. If the restored grade is needed to validate `gc`, coordinate
exact candidate branches and
retain the proof; do not manufacture a circular requirement that a grade-only
repair must pass before conduct can be fixed. No existing required-lane rule is
waived. Deterministic preparation can run independently, but Commander acceptance
and shared-file landings remain serialized.

`gqsw` is a satisfied prerequisite, not a dispatch slot. `w5` and `6h3` are not
automatically new sprint members. The separate flat-entity / merge-guard follow-up
is not a prerequisite for these seven outcomes.

## Commander cold-boot acceptance package

Before activation, supply:

- Correct membership query and captain-approved scope/restatements; no readiness
  inference from `ready` labels alone.
- Code/state checkout initialization and immutable SHAs, including the active-owner
  join environment and external dependency evidence.
- The sequence above, merge/rebase ownership, and explicit stop/escalation rules
  for red conduct rather than permission to broaden an evidence task.
- Portable trace inputs, isolated home/auth/model setup, published version pins,
  artifact destinations, CI approval, and realistic live/full/race budgets.
- An assembled seven-outcome evidence checklist: no-env discovery of both
  extensions; doctor remedy transition; capped-run summary and retained failure;
  3/3 full gate journeys; A/B/A install offers; published-family lane; four binding
  dispositions with owner/evidence joins.
- Required affected lanes and detached audits, followed by the independent
  assembled-sprint audit. Target remains the `next` development line: merge code
  through the normal `main` workflow and publish only as authorized. **No stable
  release tag from this sprint.**

## Review validation

Executed the two canonical membership queries; both returned the stale eighth
member `w5`.

Executed:

```bash
SPACEDOCK_LIVE_STATE_DIR="$PWD/docs/dev/.spacedock-state" \
  go test ./internal/contractlint \
  -run '^(TestRuntimeLiveRegistryReconciliation|TestRuntimeLiveTODOOwnersAreActive)$' \
  -count=1 -v
```

Registry reconciliation passed and derived `xfail/pi=4`. The owner join failed
only for Pi rejection-flow's archived `p17`. This is the expected unimplemented
baseline, not a claim that re-anchoring has shipped.

`go test ./...` and `go test ./... -race` each exceeded a 120-second review budget.
Treat these as the reported environment limitation, not a design defect or green
acceptance. `gofmt -w ./cmd ./internal` completed; its unrelated existing formatting
delta in `internal/release/runtime_live_evidence_workflow_test.go` was restored.
No product files were changed and no live/model-backed run was made by this review.

**Readiness result:** close M1–M6, obtain the dependency decisions, and package the
cold boot before calling this sprint drivable. Naming all outcomes is not yet a
composed proof that this bounded sprint can deliver them.
