# 0230 (v0.23.0) — stable finalization (lean cut) — Commander dispatch (cold-boot)

> **Self-contained cold-boot package.** A Commander session boots `spacedock:first-officer`,
> reads this file + the sprint `index.md`, and drives the members to a tagged stable
> `v0.23.0`. The Shaping FO's work is done; this is the handoff to the drive phase. This
> package is runnable from a COLD BOOT — everything load-bearing is here.

## Boot

Sprint = the entities matching `sprint: 0230-stable-finalization` (query, not a list).
Drivable set:
`spacedock status --workflow-dir docs/dev --where sprint=0230-stable-finalization --where 'sprint-readiness != defer'`

Boot the first officer (`spacedock claude`), run `status --boot`, and read each member body
for its design + ACs. Goal/DoD + the VALUE GATE: `index.md`. Foundation: this cut rides the
`v0.22.0` stable floor; `main` is at `b7ecd04a` (the `#420` cask fix). The plugin manifests
read `0.22.0` (`.claude-plugin/plugin.json`, `.codex-plugin/plugin.json`).

> **Member 4 is ONE pass over three query members.** All six members carry
> `sprint: 0230-stable-finalization` (state commit `5b0dcbc8`) and DO appear in the query.
> Member 4's three targets — `releasing-doc-pre-stamp-drift` (`7yd3mbsy2am5qggc17sxvz2v`),
> `stamp-then-tag-release-ritual` (`ezn308z0chwc2zvmyny9ry8w`), and
> `steady-state-stable-release-runbook` (`qpfmdxy6438fsndp9nw4c89e`) — collide on
> `docs/releasing.md`, so reconcile all three in ONE serialized pass — see member 4 below.

## Deliverable & DoD — begin with the end

**`v0.23.0`** = a LEAN stable release cut from `main`. The token clawback is the WHY the
stable was held; it is the deliverable's gate, not side polish. Done when the VALUE GATE is
met AND the tag fires on an e2e-green commit AND a fresh `brew install` launches.

### THE VALUE GATE (the DoD; begin-with-the-end STOP)

**The `v0.23.0` tag does NOT fire until this gate is met.** Stop here and check it BEFORE
doing release-cut work. A green `go test`, a clean pre-cut audit, and finished members do
NOT authorize the tag if the byte counts still exceed the baseline.

Measured `wc -c` on the boot-resident contract files MUST be `≤` the `v0.22.0` baseline for
EVERY runtime:

```
first-officer-shared-core.md   ≤ 28586   (30640 at the commission HEAD, +2054 owed)
fo-dispatch-core.md            ≤ 17488   (22929 at HEAD, +5441 owed — the big one)
fo-merge-core.md               ≤ 8059    (8597 at HEAD, +538 owed)
claude-first-officer-runtime.md ≤ 4575   (4575 — met at the line)
codex-first-officer-runtime.md  ≤ 6004   (6043 at HEAD, +39 owed)
pi-first-officer-runtime.md     ≤ 3754   (3754 — already met; #418 over-delivered)
```

Files live under `skills/first-officer/references/`. Re-measure with `wc -c` on these exact
paths against the assembled `main` tip immediately before the cut. **The release is NOT cut
while ANY byte count exceeds its `v0.22.0` number.** If any file is over, the tag does not
fire — bounce back to the owning member (1 owns the shared-core + adapter clawback; 2 owns
`fo-dispatch-core.md` + the merge-core overage).

## The members (riskiest-first)

| # | slug | id | deliverable |
| --- | --- | --- | --- |
| 1 | `fo-contract-token-cut` | `y2r7ew51xqs6q3avsb6mcaka` | Boot-resident clawback: rebase the `fo-contract-token-cleanup.md` proposal (#418 landed → just a rebase + the remaining cuts), apply the verified cuts to `first-officer-shared-core.md` + the runtime adapters, fold in the zero-discover "no filesystem sweep" prohibition. Leads the value-measurement de-risk. |
| 2 | `trim-dispatch-adapter-prose` | `adk755xqeb4a9dxhhgtjwawh` | The `+5441 B` dispatch/merge-core half. SB2-class: premature-reap ban + idle guardrail + reuse/await semantics stay BYTE-INTACT; AC gates on a NEGATIVE cumulative line delta; `contractlint`-guarded. Lands AFTER member 1. |
| 3 | `#420` cask channel-switch | merged `b7ecd04a` | ALREADY MERGED. In the tag because goreleaser regenerates both casks. No drive. |
| 4 | `docs/releasing.md` reconciliation | `7yd3mbsy2am5qggc17sxvz2v` + `ezn308z0chwc2zvmyny9ry8w` + `qpfmdxy6438fsndp9nw4c89e` | ONE serialized pass folding three colliding entities. Fixes the stale step-3 manual pre-stamp that blocks the cutter at the exact-SHA e2e-gate. |
| 5 | `spacedock-marketplace-source-env` | `z2tjv3570ahjxewv1c309rbc` | Codex EDGE install via channel-via-marketplace-name: edge installs as `spacedock@spacedock-edge` from a distinct `spacedock-edge` MARKETPLACE NAME so codex's entry-name == `plugin.json`-name check passes. Carries the marketplace restructure to expose that name. (The standalone-marketplace decouple is a SHIPPED PREREQUISITE, PR #352 — not a member.) |

## Drive order

1. **Member 1 (`fo-contract-token-cut`) FIRST** — it leads the value-measurement de-risk
   (pin the `wc -c` baseline + the measurement command) AND owns the boot-resident
   clawback. Run the no-guidance-control micro-test on the cut clauses, apply the cuts,
   re-test the keeps (the control can overturn a keep), then verify the 4 live shared
   scenarios green on Claude AND Codex.
2. **Member 2 (`trim-dispatch-adapter-prose`) AFTER member 1** — they both touch resident
   contract prose; serialize. AC gates on a negative cumulative line delta;
   `contractlint` green is the merge condition. The premature-reap ban + idle guardrail +
   reuse/await semantics must remain byte-intact.
3. **Member 3 (`#420`) is already merged** — confirm `b7ecd04a` is in `main`'s log; no drive.
4. **(Parallel) Member 5 (codex-edge)** — independent of 1/2/4; drive whenever.
5. **VALUE-GATE CHECK** — re-measure all six resident byte counts vs `v0.22.0`; if any
   exceeds, the tag does NOT fire — bounce back to the owning member.
6. **Member 4 (`docs/releasing.md` reconciliation)** — one serialized pass over the three
   colliding entities; land before the cut so the cutter follows a doc that tags the
   e2e-green commit, not a freshly-stamped one.
7. **Pre-cut antipattern audit** — independent reviewer (tag not yet fired).
8. **Write the `v0.23.0` release notes** — MUST cover `pre.1`/`pre.2` content + the clawback.
9. **Cut.**

## Drive procedure (per member)

For each drivable member (1, 2, 4, 5), the Commander runs the standard FO dispatch cycle:

1. **Advance `→ implementation`** (worktree stage). `status --set {slug} status=implementation
   worktree=.worktrees/spacedock-ensign-{slug} started`. Commit the state transition
   path-scoped, push.
2. **Create the worktree** on first dispatch.
3. **Build the dispatch** via `spacedock dispatch build` with a checklist file (≤3 items —
   the dispatch-core cap).
4. **Dispatch** the worker (the FO follows its runtime adapter's dispatch capability).
5. **On completion, verify the stage report** against the entity file (`status --read {ref}
   --json` → last `## Stage Report` → `Read(offset,limit)`). Never advance on a cheerful
   summary alone.
6. **Advance `→ validation`** (fresh validator; `feedback-to: implementation`; `gate: true`).
7. **At the validation gate**, present it (`present-gate`); the captain decides. Approve →
   terminal ceremony (PR-merge mod). Reject → `feedback-rejection-flow`.
8. **Detached adversarial audit at validation** for the high-stakes surfaces — the shipped
   FO-contract trims (members 1, 2) and the release machinery (member 4). For the contract
   trims, the audit refutes that the diff dropped a load-bearing FO-safety guarantee while
   the byte count dropped.

## In-drive gates (captain-owned)

- **Validation gates: captain decides.** The Commander presents each (`present-gate`); never
  self-approve.
- **Merge gate:** each member merges to `main` via the `pr-merge` mod; the mod-block
  enforcement runs at terminalization.
- **VALUE-GATE CHECK (the DoD funnel):** after members merge and before the pre-cut audit,
  re-measure the six resident byte counts vs `v0.22.0`. If any exceeds → the tag does NOT
  fire; bounce back. This is the begin-with-the-end stop.
- **Pre-cut antipattern audit:** with all members merged and the tag NOT yet fired, dispatch
  an INDEPENDENT reviewer (staff-eng persona; not the Commander, not the implementers) over
  the assembled sprint. The contract trims are the highest-stakes surface: confirm no
  load-bearing FO guarantee was dropped to make the byte count. Ship-blockers fixed before
  the cut; non-blockers seeded.
- **Release cut:** `go test ./...` green from the root, then `docs/releasing.md`. Captain
  authorizes the cut.

## Per-member build notes

### Member 1 — `fo-contract-token-cut` · boot-resident clawback (the WHY) · leads the de-risk

`#418` already landed, so the proposal (`docs/dev/_proposals/fo-contract-token-cleanup.md`)
is just a rebase + the remaining cuts to `first-officer-shared-core.md` and the runtime
adapters. The proposal carries the per-clause `safe-cut` / `cut-with-care` / `keep` verdicts
and the no-guidance-control micro-test method. The drive: validate the control on the
sample clauses spanning the verdict space (e.g. a `keep`, a `cut-with-care`, a `safe-cut`),
apply the default-path cuts, RE-TEST the remaining keeps (the control can overturn a keep),
then confirm the four live shared scenarios (`gate-guardrail`, `rejection-flow`,
`feedback-3-cycle-escalation`, `merge-hook-guardrail`) stay green on Claude AND Codex.
Folds in the zero-discover "no filesystem sweep" prohibition. The contract files are
shipped scaffolding (`skills/` is product, not FO-direct-editable) — edits ship through a
dispatched worker in a worktree under test. **This member also pins the value-gate
measurement** (the exact `wc -c` paths + baseline numbers from `index.md`) so the gate is
checked against a fixed oracle.

### Member 2 — `trim-dispatch-adapter-prose` · the +5441 B half · SB2-class · lands after 1

The single largest delta on the gate (`fo-dispatch-core.md`, +5441 B; plus the merge-core
+538 B overage). SB2-class risk: the trim MUST keep these BYTE-INTACT — the premature-reap
ban, the dispatch idle guardrail, and the reuse/await semantics. These are the FO-safety
guarantees the `fo-contract-token-cleanup.md` keep-list protects (the premature-teardown
ban and the DISPATCH IDLE GUARDRAIL are named keeps). The AC gates on a NEGATIVE cumulative
line delta and is `contractlint`-guarded (the notation↔command-tree binding, not a grep).
Lands after member 1 to avoid a resident-prose collision.

### Member 3 — `#420` cask channel-switch · ALREADY MERGED · no drive

Squash `b7ecd04a` on `main`. The edge cask installs the `spacedock` binary and declares
`conflicts_with` stable, so a user switching between the `spacedock` stable cask and the
`spacedock@next` edge cask gets the right binary. Listed because the stable tag regenerates
BOTH casks via goreleaser — it needed to be in before the tag. Verified by the post-cut
fresh-install check.

### Member 4 — `docs/releasing.md` reconciliation · ONE serialized pass · before the cut

Three entities (all `sprint: 0230-stable-finalization`) collide on `docs/releasing.md` —
serialize them into a single pass, do NOT dispatch in parallel:
- `releasing-doc-pre-stamp-drift` (`7yd3mbsy2am5qggc17sxvz2v`) — step 3 today does a manual
  `stamp-version` → commit, creating a fresh SHA the e2e-gate has never run green on.
- `stamp-then-tag-release-ritual` (`ezn308z0chwc2zvmyny9ry8w`) — the tagged commit's
  manifest must match its tag.
- `steady-state-stable-release-runbook` (`qpfmdxy6438fsndp9nw4c89e`) — advance `main` from a
  green `next` tip; reconcile `docs/releasing.md`.

The bug they collectively fix: `docs/releasing.md`'s `## What the Tag Push Does` already
says goreleaser "publishes only after the `e2e-gate` job confirms the tagged commit has a
green Runtime Live E2E run (or a recorded `SPACEDOCK_E2E_GATE_WAIVER`)" — but the cut steps
(step 3) tell the cutter to make a fresh manual pre-stamp commit, which the e2e-gate has
never exercised. A cutter following the doc verbatim tags a commit the gate will block.
The reconciliation makes the doc tag the GREEN e2e commit.

### Member 5 — `spacedock-marketplace-source-env` (z2) · codex EDGE install · independent

The standalone-marketplace decouple is a SHIPPED PREREQUISITE (PR #352 — the standalone
`spacedock-dev/marketplace` repo with the `spacedock`/`spacedock-edge` entries already
exists, per `docs/releasing.md`). It is NOT a live 0230 member; reference it only as
already-shipped groundwork. z2 carries the remaining codex-edge fix:
channel-via-marketplace-name. Today the standalone marketplace is ONE marketplace NAMED
`spacedock` carrying two ENTRIES (`spacedock` stable, `spacedock-edge` edge); codex's
entry-name == `plugin.json`-name check REJECTS the `spacedock-edge` entry, so edge cannot
install. z2's fix exposes a distinct `spacedock-edge` MARKETPLACE NAME so edge installs as
`spacedock@spacedock-edge` and the name-match passes — that marketplace restructure is
within z2's scope. Independent of members 1/2/4 — drive in parallel.

## Release-cut recipe (the finalization sequence — DoD path)

1. **Members 1, 2, 4, 5 merged to `main`** via PR-merge; member 3 already in (`b7ecd04a`).
2. **VALUE-GATE CHECK** — re-measure all six resident byte counts (`wc -c` on the
   `skills/first-officer/references/` paths) vs `v0.22.0`. **If any exceeds, the tag does
   NOT fire — bounce back to the owning member.** This is the begin-with-the-end stop.
3. **Pre-cut antipattern audit** — independent reviewer over the assembled sprint, tag NOT
   yet fired. The contract trims are the highest-stakes surface (refute that a load-bearing
   FO guarantee was dropped to make the byte count). Ship-blockers fixed before the cut.
4. **Write the `v0.23.0` release notes** — MUST cover the `pre.1`/`pre.2` content AND the
   token clawback, so a stable user upgrading from `0.22.0` sees the whole story.
5. **`go test ./...` green** from the root.
6. **Cut** per the (now-reconciled) `docs/releasing.md`:
   - Stamp `0.23.0` into the plugin manifests (`.claude-plugin/plugin.json`,
     `.codex-plugin/plugin.json` — currently `0.22.0`) via `go run ./cmd/spacedock-release
     stamp-version 0.23.0 ...`.
   - Tag the e2e-GREEN commit `v0.23.0` (annotated, the release notes as the changelog).
   - Push to fire goreleaser — it regenerates BOTH casks (with `#420`'s fix), publishes the
     GitHub Release, and stamps the manifests.
   - The standalone `spacedock-dev/marketplace` stable `spacedock` entry is repointed to the
     released tag (a commit in that repo, NOT a manifest stamp on `main`).
7. **Verify** the tagged commit's manifest reads `0.23.0` AND a fresh `brew install` of the
   released cask launches (the postflight strips `com.apple.quarantine` → Gatekeeper-clean).
8. **Cut from `main`** (NOT `next` — `next` is the dev/edge line).

   *Captain authorizes the cut.*

## Escalation (Commander → captain)

Escalate only on: a 3rd feedback cycle on any member, a budget blowout, an irrecoverable
block, or a genuine scope fork. Specific triggers for THIS sprint:

- **The value gate cannot be met without dropping a load-bearing FO guarantee** — if a trim
  cannot reach the byte target without touching a named keep (premature-reap ban, idle
  guardrail, reuse/await semantics), escalate; do not ship a leaner-but-unsafe contract.
- **The e2e-gate blocks the tag** — if the tagged commit lacks a green Runtime Live E2E run
  and no waiver applies, escalate (this is member 4's bug surfacing live).
- **A fresh `brew install` does not launch** post-cut — escalate; the `#420`/postflight
  chain is the first-launch guarantee and a failure there blocks the deliverable's promise.

Otherwise the Commander drives to the DoD and presents the validation gates + the value-gate
check + the pre-cut audit + the release cut for captain authorization.

## Close (Shaping FO, post-drive)

- Seed the next sprint — fold the pre-cut audit's non-blockers + the deferred set
  (`notarize-macos-release`, `gate-on-end-value`, `fo-self-evidence-bar`,
  `next-independent-release-line`, `live-test-stream-tee-ci-stdout-noise` → 0203) into the
  next backlog.
- Light post-cut release verification — some release-machinery issues only manifest when the
  tag actually fires (the cask regeneration, the marketplace repoint, the fresh-install
  launch).
