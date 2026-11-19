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
`v0.22.0` stable floor. The plugin manifests read `0.22.0` (`.claude-plugin/plugin.json`,
`.codex-plugin/plugin.json`).

> **The query returns 6 tagged entities (members 1, 2, 4, 5, 6, 7).** `#420` (member 3) is
> already merged — not an entity. **Member 4 is the promoted single consolidated lead:**
> `releasing-doc-pre-stamp-drift` (`7yd3mbsy2am5qggc17sxvz2v`) carries the whole
> `docs/releasing.md` reconciliation; its two former co-targets `stamp-then-tag-release-ritual`
> (`ezn308z0chwc2zvmyny9ry8w`) and `steady-state-stable-release-runbook`
> (`qpfmdxy6438fsndp9nw4c89e`) are UNTAGGED and folded into the lead — drive member 4 as ONE
> serialized pass over `docs/releasing.md`, not three.

## Deliverable & DoD — begin with the end

**`v0.23.0`** = a LEAN stable release cut from `main`. The token clawback is the WHY the
stable was held; it is the deliverable's gate, not side polish. Done when the VALUE GATE is
met AND the e2e-gate is reliably green AND the tag fires on an e2e-green commit AND a fresh
`brew install` launches.

### THE VALUE GATE (the DoD; begin-with-the-end STOP)

**The `v0.23.0` tag does NOT fire until this gate is met.** Stop here and check it BEFORE
doing release-cut work. A green `go test`, a clean pre-cut audit, and finished members do
NOT authorize the tag if the byte counts still exceed the baseline.

Measured `wc -c` MUST be `≤` the `v0.22.0` baseline for EVERY file — FO **and** ensign:

```
# FO boot-resident (skills/first-officer/references/)
first-officer-shared-core.md    ≤ 28586   (30640 at commission, +2054 owed)
fo-dispatch-core.md             ≤ 17488   (22929, +5441 owed — the big one)
fo-merge-core.md                ≤ 8059    (8597, +538 owed)
claude-first-officer-runtime.md ≤ 4575    (4575 — met at the line, net-zero)
codex-first-officer-runtime.md  ≤ 6004    (6043, +39 owed)
pi-first-officer-runtime.md     ≤ 3754    (3754 — met; #418 over-delivered)

# Ensign runtime adapters (skills/ensign/references/) — member 6
codex-ensign-runtime.md         ≤ 2390    (2847, +457 owed)
pi-ensign-runtime.md            ≤ 1768    (2838, +1070 owed)
claude-ensign-runtime.md        ≤ 2556    (2556 — met, net-zero guardrail)
ensign-shared-core.md           ≤ 8829    (8829 — met, net-zero guardrail)
```

Re-measure with `wc -c` on these exact paths against the assembled `main` tip immediately
before the cut. **The release is NOT cut while ANY byte count (FO or ensign) exceeds its
`v0.22.0` number.** If any is over, the tag does not fire — bounce back to the owning member
(1 owns the FO shared-core + adapters; 2 owns `fo-dispatch-core.md` + the merge-core overage;
6 owns the ensign adapters). The gate is human-measured `wc -c` (no in-repo byte-budget guard).

## The members (riskiest-first)

| # | slug | id | deliverable |
| --- | --- | --- | --- |
| 1 | `fo-contract-token-cut` | `y2r7ew51xqs6q3avsb6mcaka` | FO boot-resident clawback (`first-officer-shared-core.md` + the FO adapters incl. codex-FO +39); fold in the zero-discover "no filesystem sweep" prohibition. #418 landed → rebase + remaining cuts. Leads the value-measurement de-risk. |
| 2 | `trim-dispatch-adapter-prose` | `adk755xqeb4a9dxhhgtjwawh` | The `+5441 B` dispatch/merge-core half. SB2-class: in-surface reuse/await text + `«merge.guard»` semantics BYTE-INTACT; git-diff-ZERO to `claude-fo-dispatch.md`. AC-1 = absolute `wc -c` ≤ v0.22.0; `contractlint`-guarded. After member 1. |
| 3 | `#420` cask channel-switch | merged `b7ecd04a` | ALREADY MERGED. In the tag because goreleaser regenerates both casks. No drive. |
| 4 | `docs/releasing.md` reconciliation | `7yd3mbsy2am5qggc17sxvz2v` | ONE serialized pass (the promoted lead; stamp-then-tag + steady-state folded in). Fixes the stale step-3 manual pre-stamp that blocks the cutter at the exact-SHA e2e-gate. |
| 5 | `spacedock-marketplace-source-env` | `z2tjv3570ahjxewv1c309rbc` | Codex EDGE install via channel-via-marketplace-name: edge installs as `spacedock@spacedock-edge` from a distinct `spacedock-edge` MARKETPLACE NAME so codex's entry-name == `plugin.json`-name check passes. (Decouple is a SHIPPED prereq, PR #352 — not a member.) |
| 6 | `ensign-runtime-binding-block-cleanup` | `x1khmz0e80fyhe7vnjg8w59y` | Ensign-side clawback: reduce codex/pi ensign adapters (+457 / +1070) to host-specific binding blocks (or absorb into the FO adapter); claude-ensign + ensign-shared-core net-zero. AC-3 is the gate (codex-ensign ≤2390, pi-ensign ≤1768). |
| 7 | `opus-rejection-flow-reviewer-routing-flake` (7hc) | `7hczkc0c6ezgwy1p627ejp6x` | e2e-gate de-risk: establish the opus `rejection-flow` flake frequency (N≥3); harden the feedback-routing prose if recurring. NOT a byte-gate member. |

## Drive order

1. **Member 1 (`fo-contract-token-cut`) FIRST** — leads the value-measurement de-risk (pin
   the `wc -c` baselines + the measurement command) AND owns the FO boot-resident clawback.
   Run the no-guidance-control micro-test on the cut clauses, apply the cuts, re-test the
   keeps (the control can overturn a keep), confirm the 4 live shared scenarios green on
   Claude AND Codex. Bind the no-sweep prohibition to `detectBroadSearchAtBoot`, not a grep.
2. **(Parallel) Member 6 (`ensign-runtime-binding-block-cleanup`)** — independent of the FO
   clawback (different files: `skills/ensign/references/`); drive in parallel with 1/2.
3. **Member 2 (`trim-dispatch-adapter-prose`) AFTER member 1** — both touch FO resident prose;
   serialize. SB2-protection: the in-surface reuse/await text + `«merge.guard»` semantics stay
   byte-intact, and `claude-fo-dispatch.md` (out of surface) gets git-diff-ZERO. M2 must NOT
   touch shared-core's `→ shipped (this sprint)` markers (those are M1's).
4. **Member 3 (`#420`) is already merged** — confirm `b7ecd04a` is in `main`'s log; no drive.
5. **(Parallel) Member 5 (codex-edge)** — independent of 1/2/4/6; drive whenever.
6. **(Parallel/early) Member 7 (rejection-flow e2e de-risk)** — establish the opus frequency
   early; if recurring, harden before the cut so the tag's required green Runtime Live E2E is
   not a coin-flip.
7. **VALUE-GATE CHECK** — re-measure ALL FO + ensign resident byte counts vs `v0.22.0`; if any
   exceeds, the tag does NOT fire — bounce back to the owning member.
8. **Member 4 (`docs/releasing.md` reconciliation)** — one serialized pass; land before the cut
   so the cutter follows a doc that tags the e2e-green commit, not a freshly-stamped one.
9. **Pre-cut antipattern audit** — independent reviewer (tag not yet fired).
10. **Write the `v0.23.0` release notes** — MUST cover `pre.1`/`pre.2` content + the clawback.
11. **Cut.**

## Drive procedure (per member)

For each drivable member (1, 2, 4, 5, 6, 7), the Commander runs the standard FO dispatch cycle:

1. **Advance `→ implementation`** (worktree stage). `status --set {slug} status=implementation
   worktree=.worktrees/spacedock-ensign-{slug} started`. Commit the state transition path-scoped,
   push.
2. **Create the worktree** on first dispatch.
3. **Build the dispatch** via `spacedock dispatch build` with a checklist file (≤3 items — the
   dispatch-core cap).
4. **Dispatch** the worker (the FO follows its runtime adapter's dispatch capability).
5. **On completion, verify the stage report** against the entity file (`status --read {ref}
   --json` → last `## Stage Report` → `Read(offset,limit)`). Never advance on a cheerful summary.
6. **Advance `→ validation`** (fresh validator; `feedback-to: implementation`; `gate: true`).
7. **At the validation gate**, present it (`present-gate`); the captain decides. Approve →
   terminal ceremony (PR-merge mod). Reject → `feedback-rejection-flow`.
8. **Detached adversarial audit at validation** for the high-stakes surfaces — the shipped FO
   contract trims (members 1, 2), the ensign contract trim (member 6), and the release machinery
   (member 4). For the contract trims, the audit refutes that the diff dropped a load-bearing
   guarantee while the byte count dropped.

## In-drive gates (captain-owned)

- **Validation gates: captain decides.** The Commander presents each (`present-gate`); never
  self-approve.
- **Merge gate:** each member merges to `main` via the `pr-merge` mod; the mod-block enforcement
  runs at terminalization.
- **VALUE-GATE CHECK (the DoD funnel):** after members merge and before the pre-cut audit,
  re-measure ALL FO + ensign resident byte counts vs `v0.22.0`. If any exceeds → the tag does NOT
  fire; bounce back. This is the begin-with-the-end stop.
- **Pre-cut antipattern audit:** with all members merged and the tag NOT yet fired, dispatch an
  INDEPENDENT reviewer (staff-eng persona; not the Commander, not the implementers) over the
  assembled sprint. The contract trims are the highest-stakes surface: confirm no load-bearing
  guarantee was dropped to make the byte count. Ship-blockers fixed before the cut.
- **Release cut:** `go test ./...` green from the root, then `docs/releasing.md`. Captain
  authorizes the cut.

## Per-member build notes

### Member 1 — `fo-contract-token-cut` · FO boot-resident clawback (the WHY) · leads the de-risk

`#418` landed, so the proposal (`docs/dev/_proposals/fo-contract-token-cleanup.md`) is a rebase +
the remaining cuts to `first-officer-shared-core.md` and the FO runtime adapters (incl. clawing
codex-FO +39; `claude-FO` is at the line — net-zero guardrail). The proposal carries the per-clause
`safe-cut`/`cut-with-care`/`keep` verdicts and the no-guidance-control micro-test method. Drive:
validate the control on sample clauses across the verdict space, apply the default-path cuts,
RE-TEST the keeps (the control can overturn a keep), confirm the four live shared scenarios
(`gate-guardrail`, `rejection-flow`, `feedback-3-cycle-escalation`, `merge-hook-guardrail`) green on
Claude AND Codex. Folds in the zero-discover "no filesystem sweep" prohibition bound to
`detectBroadSearchAtBoot` / `TestLiveZeroDiscoverReportsAndStops` (a stochastic zero-discover red is
re-run-grounds, never a merge blocker). + a `docs/runtime-support.md` alignment AC. Contract files
are shipped scaffolding — edits ship through a dispatched worker in a worktree under test. **Pins
the value-gate measurement** so the gate has a fixed oracle.

### Member 2 — `trim-dispatch-adapter-prose` · the +5441 B half · SB2-class · after 1

The largest single delta (`fo-dispatch-core.md` +5441; `fo-merge-core.md` +538). SB2-class: keep
BYTE-INTACT the in-surface reuse/await text (`## Reuse and Fresh Dispatch` conditions 0-4, the
verbatim anchor `does not match next stage effective_model`, supersede-shutdown) and the
`«merge.guard»` armed/blocked/finalized phase semantics + never-finalize-on-pr-presence rule; and
assert git-diff-ZERO to `claude-fo-dispatch.md` (the `## Awaiting Completion` / premature-reap ban /
idle guardrail home, OUT of the cut surface). AC-1 = absolute `wc -c` ≤ v0.22.0 (fo-dispatch-core
≤17488 AND fo-merge-core ≤8059), report the signed cumulative delta; `contractlint`-guarded; + a
runtime-support alignment AC. After member 1 (resident-prose collision); does NOT touch shared-core's
`→ shipped (this sprint)` markers (M1's).

### Member 3 — `#420` cask channel-switch · ALREADY MERGED · no drive

Squash `b7ecd04a` on `main`. Edge cask installs the `spacedock` binary, declares `conflicts_with`
stable. In because the stable tag regenerates BOTH casks via goreleaser. Verified by the post-cut
fresh-install check.

### Member 4 — `docs/releasing.md` reconciliation · ONE serialized pass · before the cut

The promoted lead (`releasing-doc-pre-stamp-drift`); `stamp-then-tag-release-ritual` +
`steady-state-stable-release-runbook` are folded in (untagged). ONE pass over `docs/releasing.md` +
its byte-identical mirror `docs/site/contributing/releasing.md`. The bug: step 3 today does a manual
pre-stamp commit, creating a fresh SHA the e2e-gate has never run green on; `docs/releasing.md`'s
`## What the Tag Push Does` already says goreleaser "publishes only after the `e2e-gate` job confirms
the tagged commit has a green Runtime Live E2E run (or a recorded `SPACEDOCK_E2E_GATE_WAIVER`)" — so a
cutter following step 3 verbatim tags a commit the gate BLOCKS. Reconcile so the doc tags the GREEN
e2e commit, describes the moving-`stable`-branch auto-advance (not a manual repoint-to-tag), and adds
a divergeable guard (tagged commit's `plugin.json` version == tag semver). VERIFY FIRST: the live doc
already says "cut from main", so the steady-state fold is likely a DELETION of stale next-framing, not
an addition. The `next` marketplace version-field cleanup is OUT (standalone repo, unreachable).

### Member 5 — `spacedock-marketplace-source-env` (z2) · codex EDGE install · independent

The standalone-marketplace decouple is a SHIPPED PREREQUISITE (PR #352). z2 carries the remaining
codex-edge fix: today the standalone marketplace is ONE marketplace NAMED `spacedock` carrying two
ENTRIES; codex's entry-name == `plugin.json`-name check REJECTS the `spacedock-edge` entry, so edge
cannot install. z2 exposes a distinct `spacedock-edge` MARKETPLACE NAME so edge installs as
`spacedock@spacedock-edge`. Also: make the codex fresh-install programmatic, make both resolvers
channel-aware, add `SPACEDOCK_MARKETPLACE_SOURCE`. AC-1 = a real `codex plugin add` of the edge
channel on the live lane (the failing-today baseline flipped to success); the claude half is a
construction-guarantee + resolver test, real install OUT-OF-BAND (mutates global state, not a CI
gate). Flip `channel_selection_test.go` from `spacedock-edge@spacedock` to `spacedock@spacedock-edge`.
Independent — drive in parallel.

### Member 6 — `ensign-runtime-binding-block-cleanup` · ensign-side clawback · independent

The ensign analog of t0g/#418. The codex/pi ensign runtime adapters grew (+457 / +1070 over
`v0.22.0`) from ~70-80% duplication of `ensign-shared-core.md`; reduce each to a compact host-specific
binding block (clarification tool, completion signal, captain comms, shutdown) — or, if the
ensign-specific content is fully FO-side, absorb into the FO adapter and delete the per-host ensign
files. `claude-ensign-runtime.md` + `ensign-shared-core.md` stay net-zero (guardrail; the shared core
is the authority — do NOT trim it). AC-3 is the value gate (codex-ensign ≤2390, pi-ensign ≤1768); + a
`docs/runtime-support.md` alignment AC (document the ensign binding-block shape, aligned with whichever
path lands); update the `#417` pi-ensign contractlint guard. Independent of the FO clawback — parallel.

### Member 7 — `opus-rejection-flow-reviewer-routing-flake` (7hc) · e2e-gate de-risk · early

A live-lane reliability member, not a byte-gate member. The `v0.23.0` tag fires only on a green
Runtime Live E2E (the e2e-gate). The opus `rejection-flow` lane intermittently routes the cycle-2
re-review to the impl worker (fix≠reviewer violation) where sonnet gets it right. AC-1: re-run the opus
`rejection-flow` scenario N≥3 isolated, record the fail rate — clean ⇒ rare flake (accept re-run,
close); recurring ⇒ real gap. AC-2 if recurring: identify the contract prose opus mis-reads
(`feedback-rejection-flow` / `«feedback.route»`), harden it, prove opus passes N≥3. Land before the cut
so the e2e-gate is not a coin-flip.

## Release-cut recipe (the finalization sequence — DoD path)

1. **Members 1, 2, 4, 5, 6 merged to `main`** via PR-merge; member 3 already in (`b7ecd04a`); member 7
   resolved (lane proven-rare or hardened).
2. **VALUE-GATE CHECK** — re-measure ALL FO + ensign resident byte counts (`wc -c` on the
   `skills/first-officer/references/` + `skills/ensign/references/` paths) vs `v0.22.0`. **If any
   exceeds, the tag does NOT fire — bounce back.** This is the begin-with-the-end stop.
3. **Pre-cut antipattern audit** — independent reviewer over the assembled sprint, tag NOT yet fired.
   The contract trims are the highest-stakes surface (refute that a load-bearing guarantee was dropped
   to make the byte count). Ship-blockers fixed before the cut.
4. **Write the `v0.23.0` release notes** — MUST cover the `pre.1`/`pre.2` content AND the token clawback.
5. **`go test ./...` green** from the root.
6. **Cut** per the (now-reconciled) `docs/releasing.md`:
   - Stamp `0.23.0` into the plugin manifests (`.claude-plugin/plugin.json`, `.codex-plugin/plugin.json`
     — currently `0.22.0`).
   - Tag the e2e-GREEN commit `v0.23.0` (annotated, the release notes as the changelog).
   - Push to fire goreleaser — it regenerates BOTH casks (with `#420`'s fix), publishes the GitHub
     Release, advances the moving `stable` branch.
7. **Verify** the tagged commit's manifest reads `0.23.0` AND a fresh `brew install` of the released
   cask launches (the postflight strips `com.apple.quarantine` → Gatekeeper-clean).
8. **Cut from `main`** (NOT `next` — `next` is the dev/edge line). *Captain authorizes the cut.*

## Escalation (Commander → captain)

Escalate only on: a 3rd feedback cycle on any member, a budget blowout, an irrecoverable block, or a
genuine scope fork. Specific triggers for THIS sprint:

- **The value gate cannot be met without dropping a load-bearing guarantee** — if a trim (FO or ensign)
  cannot reach the byte target without touching a named keep (premature-reap ban, idle guardrail,
  reuse/await semantics, the ensign shared-core authority), escalate; do not ship a leaner-but-unsafe
  contract.
- **The e2e-gate blocks the tag** — if the tagged commit lacks a green Runtime Live E2E and no waiver
  applies, escalate (member 4's stale-step-3 bug, or member 7's rejection-flow flake, surfacing live).
- **A fresh `brew install` does not launch** post-cut — escalate; the `#420`/postflight chain is the
  first-launch guarantee.

Otherwise the Commander drives to the DoD and presents the validation gates + the value-gate check +
the pre-cut audit + the release cut for captain authorization.

## Close (Shaping FO, post-drive)

- Seed the next sprint — fold the pre-cut audit's non-blockers + the deferred set
  (`notarize-macos-release`, `gate-on-end-value`, `fo-self-evidence-bar`, `next-independent-release-line`,
  `live-test-stream-tee-ci-stdout-noise`, `haiku-drive-validation`, `fo-tier-delegation`, the
  dispatch-build pair) into the next backlog.
- Light post-cut release verification — some release-machinery issues only manifest when the tag fires
  (cask regeneration, the `stable`-branch advance, the fresh-install launch).
