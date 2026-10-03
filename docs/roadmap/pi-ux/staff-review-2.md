# Pi-UX sprint preflight staff review — second pass

Date: 2026-10-03

**Verdict: materially improved, but not yet cold-boot drivable.** M1, M3 and
M4 are closed as sprint-shaping findings. M5 and M6 remain open, with partial
progress. M2 is **superseded in its recorded-gate premise and remains open as a
re-identified delivery dependency**: two different, unbound journeys fail on the
published family, and their historical repair owner is terminal/PASSED. The
seven-member sprint has no authorized repair owner for that current red lane.

This is an independent sprint-wide review, not a second ideation gate, task
validation, scope approval, or permission to merge. “Closed” below means the
missing sprint delivery contract now exists, not that all implementation proof
has passed. AC citations establish accounting, not runtime acceptance.

## Evidence and snapshot

Read [the sprint discipline](../README.md), [the seven-outcome index](index.md),
[the first review](staff-review.md), all seven current member bodies, and the
named code seams. Also read `docs/dev/README.md`'s Proof policy,
`docs/runtime-support.md`, nta's historical repair/validation record, and mc's
retained remedy transactions.

- Review code checkout: `main`, `dc7d9a0cd` (the first review commit).
- State checkout: `cdd4154fd5a7bd9b8e5f1f853bca1c91ae4378d3`.
- mc: fold `1935cb19a`, implementation `1f41f289faad73fa1ce14684b80ee920e502f89e`,
  [PR 817](https://github.com/spacedock-dev/spacedock/pull/817), base `main`.
- mh: implementation `094e06ff7f195d5dac59e93a1760c813e62076c2`,
  [PR 816](https://github.com/spacedock-dev/spacedock/pull/816), base mc's branch.
  Both PRs remain open; mc is an ancestor of mh. Stack identifier **818** is
  supplied by the handoff, not a GitHub issue number.
- Other supplied folds: 3w1 `44622cdc3` / `9576dfe01`; z6e `fa14e3d66`;
  9w `d4d41c301`; ekw `5dc6d9a3d` / `42f40cc33`; d52 `4f657242f`.
- Independently retrieved run metadata, logs and artifact
  `runtime-live-e2e-pi-live` (artifact ID `11266907788`) for
  [run 37101046846](https://github.com/spacedock-dev/spacedock/actions/runs/37101046846).
  Its exact head is mh's `094e06ff7`, not this review's main checkout.

Both canonical membership queries now return **exactly seven members**, with
five in ideation and mc/mh in implementation. The stale w5/w5s membership is
fixed. The query abbreviates d52 as `d5`; its full ID remains
`d525n1p5zgnz99hmtjq16z57`. Membership authority is the query, not this snapshot:

```bash
spacedock status --workflow-dir docs/dev --where sprint=pi-ux
spacedock status --workflow-dir docs/dev \
  --where sprint=pi-ux --where 'sprint-readiness != defer'
```

## Measured matrix: what changed, and what did not

| Surface at `094e06ff7` | Observed result | Delivery meaning |
| --- | --- | --- |
| Offline job | Success | A real CI baseline, not a substitute for affected live lanes or race proof. |
| Install / compatibility guard / current-checkout setup | All success; installed Pi 1.0.0, subagents 0.75.0, intercom 0.16.0 | The published-family install and four replacement assertions execute successfully. |
| Pi front-door smoke | PASS, 132.27s | mc removes the launch-readiness barrier on this family. Does not prove 3w1's no-override discovery. |
| Pi common package | 17 tests, 2 failures, 1429.151s | Not a green pi-live lane and not mh AC-1 completion. |
| Recorded gate | `TestLiveCommonRecordedGateLifecycle` PASS, 323.62s | The previous recorded-gate FAIL does not reproduce here. |
| Auto-continue, single-root | FAIL: `[validation-worker-not-dispatched]` | Ordinary unbound failure. |
| Default headless gate stop | FAIL: `[gate-hold-violation gate-not-held implementation-worker-not-dispatched]` | Ordinary unbound failure. |

The downloaded `pi-coverage-detail.jsonl` directly confirms recorded-gate PASS;
this is stronger than inferring it from absence in the failure summary. However,
`piSharedLiveDriver.prepareRecordedGate` still returns `noLiveGrade` at this
exact tip. One existing-oracle pass is **not** 9w's restored semantic grade or
its declared 3/3 full-transaction batch. It neither closes gc globally nor proves
all historical presentation/continuation defects gone.

“Fifteen pass” means fifteen green Go tests, not fifteen unbound product passes.
The same artifact reports:

- owner handoff: **XPASS**, owner fe7, `observed=[]`;
- smallest mechanism: **XPASS**, owner h30, `observed=[]`;
- keep moving: **XFAIL**, owner x0, `[keep-moving-violation]`;
- rejection: **XFAIL**, still archived p17, `[rejection-worker-topology]`.

This updates d52's evidence inventory. In particular, a blanket claim that the
smallest-mechanism binding cannot emit XPASS is refuted by the measured output;
it does not establish that its extractor is trustworthy. Neither new XPASS is
permission to clear a binding without the required exact unbound normal-PASS
pair and approved disposition. Historical keep-moving XPASS is not today's
result.

## First-review finding dispositions

### M1 — Closed by the folded discovery contract

Read `pi-default-extension-discovery.md`, **“M1 staff-review fold”**, specifically
**“One isolated-home setup contract”**, **“Non-live seam, old test contract, and
bounded surface”**, and **“Proposed proof refinements.”** The fold now:

- captures the real HOME/agentDir before isolation, aligns child piHome with
  clean HOME, and shares that setup between both fixtures;
- registers BOTH npm packages and preserves one supported absolute Spacedock
  entry, rather than treating symlinks as registration;
- separates discovered roots from independent explicit overrides and forbids
  either root env variable or substrate extension argv in default proof;
- names the non-live helper file and replaces the sibling-root test contract;
- leaves entry resolution with mc and requires both loaded tools plus the
  existing durable smoke outcome.

The real-loader spike also establishes the unsupported `file:` prefix and the
working absolute-path replacement. This closes the first review's composition
and first-day build-tag decisions. Captain approval of the proposed fold and
implementation/live proof remain due; the already-green stack smoke is not a
substitute.

### M2 — Superseded premise; open, re-identified external prerequisite

Read 9w's **“External dependency handoff (M2, proposed for captain's gate)”**,
**“Proposed scope and authority restatements”**, and **AC-1/AC-4**; mh's
**AC-1** and **“Risk evidence”** dependency paragraph; and
`repair-pi-default-headless-gate-stop.md`'s **“Live-root-cause finding”**,
subsequent **“Stage Report: validation”**, and **cycle 2** reports.

The 9w fold correctly forbids repairing conduct inside a grade-only task and
requires an accepted external handoff before green-lane merges. It does **not**
commission a repair. Its assertion that gc's recorded-gate repair is a necessary
current prerequisite is now stale: that journey passed on the measured family.
Do not require an arbitrary gc code change before attempting 9w's stronger proof.
Keep gc as the routing owner if that stronger grade exposes recorded-gate conduct.
The compact-prefix/full-canonical-authority distinction also resolves the w5
contract collision without adding w5 back to scope.

**The external prerequisite is not void. It re-identifies as auto-continue and
default-headless-gate-stop.** At both main and the tested stack tip,
`shared_live_runner_test.go:120,130` pass `nil` gaps. Neither failure can be
credited as expected, and neither belongs to 9w's recorded-gate callback.

nta (`ntarrp8jp5h34g6528d66kbe`), the historical/archived PASSED repair of those
same two journeys, is not an active commissioned dependency. Its readable state
body is `status: done`, `verdict: PASSED`. Its past repair recognized
`subagent_wait` completion and later a legitimate redispatch; that historical
success is not current Pi 1.0 closure evidence.

The new log gives a useful diagnostic boundary, not a root cause: both failures
show `spawns=1 completed=-1`; auto-continue also describes a stale entity lookup,
and default-headless reports no successful `gate prepare` and zero Briefings.
An error named “worker-not-dispatched” does not by itself establish that no worker
ran. nta's history specifically warns against inferring conduct from an obsolete
completion extractor. The external owner must inspect current root/tool/durable
artifacts and distinguish observation defects from actual gate/dispatch defects;
this review does not prescribe a conduct patch or an assertion relaxation.

**Required disposition before delivery:** Shaping FO/captain commissions an
active repair/reopened owner in the appropriate external scope, or explicitly
revises sprint scope/acceptance. Record responsible owner, availability, accepted
candidate SHA(s), current failure evidence and exact-tip unbound closure for both
journeys, followed by the complete required lane. No new XFAIL, narrowed selector,
removed assertion or silent repair inside mc/mh/9w/d52 is authorized. 3w1 may
change the measured matrix through legitimate setup work; test that composition,
but do not assume it repairs these journeys or expand its scope if it does not.

### M3 — Closed by the artifact interface fold

Read z6e's **“Operator artifact contract (M3 fold; proposed for the captain's
gate)”**, **“Expected surface and tolerance”**, and **“Test plan”** M3 placement.
It supplies concrete JSON, field domains, immutable input snapshots with explicit
read/list errors, conservative unique-root/linear-history rules, all-record
parsing, and warning-before-original-timeout write failure. The deterministic
wiring check is a subtest of the existing registered front-door owner, with its
fifth file explicitly accounted for. Those are the missing implementation-day
interface choices; no new standing framework is necessary. Approval and proof
are still stage obligations, not preflight omissions.

### M4 — Closed: remedy acceptance has an owner and retained candidate evidence

Read mc's **AC-4**, **“M4 remedy acceptance and scope proposal”**, **“Integrated
proof procedure”**, and **“Stage Report: implementation.”** Implementation owns
the disposable-home missing-target → printed command → same-line OK transaction;
validation independently checks it. The folded bridge-only executable remedy
exception is covered by the recorded ideation approval.

Also read `_evidence/pi-doctor-probes-stale-subagents-layout/remedy-transactions.txt`:
both independently removed manifest targets move from MISSING to OK after the
printed `pi install npm:pi-subagents`, with Pi 1.0.0/subagents 0.75.0 recorded.
The unaffected probe remains OK. This is not merely the old unfixed-binary spike
or an install-exit-only claim. Independent task validation still owns transaction
provenance/exits and final acceptance; the sprint-wide owner gap is closed.

### M5 — Still open, membership portion closed

Read 9w's **“Retained evidence and portability”**, d52's **“Test plan”** explicit
state-join requirement, and the sprint discipline's **“Package”** step.
Membership is reconciled. `dispatch-sprint-execution.md` deliberately does not
exist yet: that is an expected pending shaping step, not a missing product file
or a reason to re-ideate the five tasks. It nevertheless prevents cold-boot
Commander activation.

The two old-worktree root files still exist and match the body's SHA-256 values;
`internal/ensigncycle/testdata/pi-recorded-gate/` does not yet exist on main.
The package must supply retrievable, provenance-preserving inputs before that
worktree disappears, not only instruct a remote Commander to use its absolute
path. d52 also needs an initialized state checkout, immutable state SHA and
`SPACEDOCK_LIVE_STATE_DIR`; otherwise its owner join skips. The updated M2 handoff,
stack inheritance and validation map remain packaging work.

### M6 — Still open, partially evidenced on the stack

Read mc's **“Expected surface and tolerance”** / **“Test plan”**, ekw's
**“Expected surface”** / **“Primary proof owners and implementation order”**,
d52's **AC-4** / **“Test plan”**, and mh's **“Test plan.”** The local proof plans
do not themselves compose the final-diff lane map. The actual mc diff reaches
launcher readiness/fallback, ekw reaches shared loaded install instructions, and
d52 reaches the shared live declaration file.

The measured Pi front-door success now supplies real evidence for mc/mh's
composition, but their common lane is red. This run skipped Claude/Codex and
cannot discharge those lanes for later shared changes. Package the following
map, re-evaluated on each final diff:

| Change | Required validation composition |
| --- | --- |
| mc + mh stack | Pi front door and complete Pi common lane on an accepted composite tip; launcher and CI detached audits; no green acceptance from install-only results. |
| 3w1 / z6e / 9w Pi harness edits | Complete affected Pi lane plus each distinct proof: no-env discovery, intentionally capped diagnostic run, and declared 3/3 restored-grade batch respectively. |
| ekw extension + shared `fo-install.md` | Pi identity/actual FO offer traces and Pi lane; Claude/Codex lanes for shared loaded instructions, not just TS mocks or text checks. Shipped-surface detached audit. |
| d52 shared declarations | Explicit-state registry join and attribution evidence; under the file-based merge rule, shared live-file changes require affected host lanes, conservatively Pi/Claude/Codex. Any narrower determination must be recorded from the actual diff, not its Pi-only title. |

Retain full/race/format gates and the independent assembled-sprint audit. z6e's
expected-red cap is diagnostic evidence, not the required green regression run.
Policy permits attributable serial flake reruns; it does not permit replacing
failed/skipped samples in 9w's declared batch or ignoring the current two failures.

## Seven-outcome reachability

| # | Outcome | Reachable owner and remaining boundary |
| --- | --- | --- |
| 1 | No-export/no-hand-wired isolated discovery | **3w1**, now composed with mc. Both loaders/tools and durable smoke remain required; CI's retained exports do not prove this. |
| 2 | Truthful doctor and effective remedy | **mc**, implemented candidate with transaction evidence; validation/merge still due. |
| 3 | Classified stall with limits | **z6e**, now has a shaped interface; live confirmation uses mh's published family and an actual capped Pi process. |
| 4 | Repeated presentation through successor dispatch | **9w**, conditionally reachable by restoring measurement and observing the existing conduct. Current recorded-gate PASS removes the demonstrated gc prerequisite, but does not satisfy the stronger 3/3 requirement. Route any newly exposed conduct failure externally. |
| 5 | Session-scoped identity/install offer | **ekw**, current-session API and A/B/A actual-offer proof, not parent env. |
| 6 | Published-family pi-live installation and validation | **mh** owns pins/assertions and has measured setup success, but **no in-scope member owns the two failing journeys needed for its complete-green-lane acceptance**. This outcome cannot currently be delivered by the bounded seven-member plan alone. |
| 7 | Active XFAIL owners plus classified evidence | **d52**, reachable without repairing bound journeys. Re-anchor p17 to active 6h3 and refresh dispositions using this run; XPASS alone does not clear. |

All seven have a named member; that is not the same as a closed delivery plan.
**Outcome 6 has the current repair-ownership hole.** Outcome 4 remains an honest
measurement risk, not a presently demonstrated absent repair owner. More broadly,
M2's red lane blocks affected merges and therefore integrated delivery of the
whole sprint. The additional format/full/race/exact-live bullet is an integration
gate, not an eighth product outcome.

## Composition, sequencing and scope after the two-layer stack

Do not dispatch mc or mh implementation again, reorder mh behind 3w1 merely to
match the first review, or treat either open PR as already merged. The existing
mc → mh order is now a tested setup prerequisite. New candidates must explicitly
inherit the approved stack tip until normal bottom-up landing/rebase makes those
changes available on main. Stack-tip evidence is evidence for that composite,
not proof that an isolated lower PR or a later rebased tip passed the same lane.
Neither layer is ready for a red-lane merge under the current policy; close M2
and validate the final candidate/merge plan, without creating a circular demand
that the lower layer pass against the obsolete pins before testing the composite.

Remaining collision seams:

- **3w1 → z6e → 9w:** serial Pi-driver construction, timeout and success-grade
  edits. z6e also restructures the existing smoke into sibling subtests: preserve
  3w1's paid-smoke/default-discovery assertions and keep deterministic setup out
  of the paid sibling. Timeout classification must not claim to run 9w's success
  callback. Their JSONL semantics differ; do not grow a generic parser framework.
- **mh → z6e:** retain pin documentation while adding diagnostic limits in
  `docs/runtime-live-ci.md`. **mc → 9w:** retain manifest-path documentation while
  adding the proof-boundary paragraph in `docs/runtime-support.md`.
- **mc / 3w1 / mh:** manifest entries versus root discovery versus CI package
  assertions remain separate. mc resolves the first declared extension; mh
  checks all declared entries. 3w1 registers roots rather than duplicating either
  resolver. The separate launcher local-source reader mismatch noted in 3w1 is
  not an unapproved reason to change production resolution here.
- **ekw / external repair:** identity changes touch `.pi/extensions/spacedock.ts`
  and the Pi adapter; a newly commissioned journey repair might touch either.
  Agree on exact candidate bases and serialize/rebase if it does. Do not assume
  the new failure has the same repair surface or cause as archived nta.
- **d52:** reconcile the actual assembled registry/state last. It has no mandatory
  dependency on z6e's timeout vocabulary: semantic/infrastructure/evidence
  dispositions are different from a capped-run diagnostic label.

The scope boundary stays intact: Pi journey repairs remain outside pi-ux unless
the captain explicitly changes it. w5, 6h3 implementation, the flat-entity/
merge-guard follow-up and an activity feed are not new dispatch slots. Active
backlog 6h3 can satisfy d52's owner join without being implementation-ready;
record its scheduling disposition rather than making d52 fix topology.

### Recommended remaining-five sequence

After captain approval of the remaining folds, use this serial implementation /
landing preference on the inherited mc+mh candidate:

```text
3w1 -> z6e -> ekw -> 9w -> d52
```

1. **3w1:** settle common/smoke isolation and real package registration first,
   using mc's now-working compiled entry resolution. Preserve separate no-env
   proof even though the CI setup is already green.
2. **z6e:** rebase onto the final setup and mh docs; add timeout diagnostics before
   further Pi-driver grading edits. Run deterministic wiring and retain the
   deliberately capped real-family evidence separately from regression greens.
3. **ekw:** prove session-scoped offers on the established family; its TS work is
   largely independent, but shared instructions require broader lane validation
   and coordination with any external repair in the adapter/extension.
4. **9w:** restore the successful-run grade after setup/timeout edits, retain
   authentic negatives and a real positive, then run the declared exact-tip
   three-run batch. Do not wait for a presumed gc fix solely because the old
   handoff says so; do route any stronger-grade failure before claiming delivery.
5. **d52:** reconcile final bindings against explicit current state and current
   artifacts, preserving distinctions between XFAIL, XPASS and unbound PASS.

This is not permission to merge sequentially through red required lanes. External
triage/commissioning for the two observed failures starts **now**, before the first
affected merge, including mc/mh. Deterministic preparation may proceed within
approved gates; acceptance waits for a real green dependency closure. If no
external owner is available, stop for the captain's scope/acceptance decision,
not a sixth undisclosed repair assigned to one of these five.

## Commander cold-boot checklist still owed

Shaping FO should package, rather than ask the Commander to rediscover:

- The reconciled query, captain-approved fold/restatement precedence, exact code
  and state initialization/SHAs, and which mc/mh implementation and validation
  work is already banked. Historical “proposed” and old-stage reports are not
  independent instructions to restart completed work.
- PR 817 → PR 816 / stack 818 bases, accepted continuation tip, bottom-up merge /
  rebase ownership, and renewed evidence rules after changes.
- The **updated** M2 external owner/availability/SHA/closure handoff, stop rules,
  and retained red artifacts. Remove the stale unconditional gc prerequisite
  from dispatch guidance without pretending gc or 9w is fully validated.
- Portable historical 9w inputs with checksums/redaction provenance; current CI
  artifact retrieval and preservation before retention expiry. For example:
  `gh run download 37101046846 -n runtime-live-e2e-pi-live -D <evidence-dir>`.
  Preserve `pi-coverage-detail.jsonl`, smoke details, session/durable artifacts
  and run/head identity, not only the job summary.
- Explicit state checkout/env for d52, current classified binding dispositions
  (including the two new XPASS observations), and 6h3 ownership disposition.
- The final-diff required-lane map, detached audit ownership, authorized auth /
  model isolation and CI approvals, runtime budgets and artifact destinations.
  Keep old local-suite failures/timeouts distinct from the green CI offline job.
- One integrated seven-outcome evidence checklist, final format/full/race and
  exact runtime gates, and the independent assembled-sprint audit. Target remains
  the `next` development line; **no stable-release tag from this sprint**.

**Readiness conclusion:** the task folds are no longer the principal obstacle.
Commission/discharge the re-identified lane dependency, complete M6's validation
composition, approve the remaining gates and write the deliberate cold-boot
package. Until then, this is a bounded preparation plan, not a drivable promise
that these seven members alone can deliver the sprint.

## Review validation and limits

Executed both membership queries, inspected both PR bases/heads, verified mc's
ancestry in mh, retrieved the exact CI run/log/artifact, parsed the emitted Go
results and grade lines, and rechecked the two historical negative-root hashes.
No new model-backed live run or product/state mutation was performed.

Executed:

```bash
SPACEDOCK_LIVE_STATE_DIR="$PWD/docs/dev/.spacedock-state" \
  go test ./internal/contractlint \
  -run '^(TestRuntimeLiveRegistryReconciliation|TestRuntimeLiveTODOOwnersAreActive)$' \
  -count=1 -v -timeout=60s
```

Registry reconciliation passed (`xfail/pi=4`); active-owner join failed **only**
for rejection-flow's archived p17. That is d52's expected unimplemented baseline,
not a new member or clearance authorization.

`go test ./...` and `go test ./... -race` each exceeded this review's 120-second
budget and were terminated; neither is claimed green. `gofmt -w ./cmd ./internal`
completed; its pre-existing formatting-only delta in
`internal/release/runtime_live_evidence_workflow_test.go` was restored. No product
code changes or tests were added. `git diff --check` passes. Only this review is
committed; pre-existing unrelated untracked files are untouched.
