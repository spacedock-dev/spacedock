# 0230 — stable v0.23.0 finalization (lean cut)

> **Commissioned 2026-06-20 (Shaping FO).** The first sprint whose deliverable is a
> STABLE release, not a feature. `v0.22.0` is the last stable tag; two pre-releases
> (`v0.23.0-pre.1`, `v0.23.0-pre.2`) shipped edge content but accumulated boot-resident
> FO-contract bloat. Stable was held precisely because that bloat must be clawed back
> before it ships to every stable user. This sprint cuts a lean `v0.23.0` from `main`.

## Theme

Cut a LEAN stable `v0.23.0` whose boot-resident FO contract is back at or below the
`v0.22.0` baseline — recovering the token bloat that accreted through `pre.1`/`pre.2` —
with the `#420` cask channel-switch fix in, a fresh `brew install` surviving first
launch, and release notes covering everything the two pre-releases shipped.

## Membership is the query, never a hard-coded list

```bash
spacedock status --workflow-dir docs/dev --where sprint=0230-stable-finalization
spacedock status --workflow-dir docs/dev --where sprint=0230-stable-finalization --where 'sprint-readiness != defer'
```

> **Member 4 is ONE consolidated pass over three query members.** All six members carry
> `sprint: 0230-stable-finalization` (state commit `5b0dcbc8`) and DO surface in the query.
> Member 4's three reconciliation targets — `releasing-doc-pre-stamp-drift`
> (`7yd3mbsy2am5qggc17sxvz2v`), `stamp-then-tag-release-ritual` (`ezn308z0chwc2zvmyny9ry8w`),
> and `steady-state-stable-release-runbook` (`qpfmdxy6438fsndp9nw4c89e`) — all collide on
> the same file (`docs/releasing.md`), so the Commander reconciles all three in ONE
> serialized pass rather than three parallel dispatches. The query returns them as three
> entities; member 4 is the single consolidated unit of work over them.

## Goal (success criterion)

A stable `v0.23.0` is tagged from `main` and a fresh `brew install` of it launches, where:

1. **The boot-resident FO contract is lean** — every resident contract file is measured
   `wc -c` at or below its `v0.22.0` baseline (the VALUE GATE below). This is the WHY the
   stable was held; it is the deliverable's gate, not side polish.
2. **The `#420` cask channel-switch fix is in** — a user who switches between the stable
   `spacedock` cask and the `spacedock@next` edge cask gets the right binary, with the two
   casks declaring `conflicts_with` each other. (Already merged; in because the tag
   regenerates both casks via goreleaser.)
3. **A fresh `brew install` survives first launch** — the cask postflight strips
   `com.apple.quarantine`, so a Gatekeeper-clean first run works without Developer-ID
   notarization (which is deferred — it needs an Apple cert).
4. **Release notes cover everything shipped** — the `v0.23.0` changelog names the `pre.1`
   and `pre.2` content PLUS the token clawback, so a stable user upgrading from `0.22.0`
   sees the whole story in one set of notes.

## The VALUE GATE — the definition of done (begin-with-the-end)

**The `v0.23.0` tag does NOT fire until this gate is met.** This is the WHY the stable
release was held: the token clawback is the deliverable. Stop at this gate and check it
BEFORE doing release-cut work — a green test suite and a clean audit do not authorize the
tag if the byte counts still exceed the baseline.

Measured `wc -c` on the boot-resident contract files MUST be `≤` the `v0.22.0` baseline
for EVERY runtime:

| resident file | v0.22.0 baseline (B) | HEAD at commission (B) | delta | met? |
| --- | ---: | ---: | ---: | --- |
| `skills/first-officer/references/first-officer-shared-core.md` | 28586 | 30640 | +2054 | NO — clawback owed |
| `skills/first-officer/references/fo-dispatch-core.md` | 17488 | 22929 | +5441 | NO — clawback owed |
| `skills/first-officer/references/fo-merge-core.md` | 8059 | 8597 | +538 | NO — clawback owed |
| `skills/first-officer/references/claude-first-officer-runtime.md` | 4575 | 4575 | 0 | met (at the line) |
| `skills/first-officer/references/codex-first-officer-runtime.md` | 6004 | 6043 | +39 | NO — small clawback owed |
| `skills/first-officer/references/pi-first-officer-runtime.md` | 3754 | 3754 | 0 | met (#418 over-delivered) |

The release is NOT cut while ANY byte count exceeds its `v0.22.0` number. Re-measure with
`wc -c` against the assembled `main` tip immediately before the cut; if any file is over,
the tag does not fire — bounce back to the owning member.

> **Measurement de-risk (riskiest-first).** The byte counts above are the gate's authority.
> Before any cut is applied, pin the baseline numbers and the measurement command (`wc -c`
> on the exact file paths) so the gate is checked against a fixed, agreed oracle, not a
> re-derived estimate. Member 1 leads this de-risk (it owns the boot-resident clawback and
> the no-guidance-control micro-test that proves a cut preserves FO behavior).

## Members (riskiest-first)

| # | slug | id | deliverable | the WHY |
| --- | --- | --- | --- | --- |
| 1 | `fo-contract-token-cut` | `y2r7ew51xqs6q3avsb6mcaka` | Boot-resident clawback: apply the verified cuts (`first-officer-shared-core.md` + the runtime adapters) from the `fo-contract-token-cleanup.md` proposal; fold in the zero-discover "no filesystem sweep" prohibition. #418 already landed → this is a rebase + the remaining cuts. | The gate's core; recovers the resident-contract bloat. |
| 2 | `trim-dispatch-adapter-prose` | `adk755xqeb4a9dxhhgtjwawh` | The `+5441 B` dispatch/merge-core half: trim `fo-dispatch-core.md` (and the merge-core overage) to capability bindings. SB2-class risk. | The single largest delta on the gate. |
| 3 | `#420` cask channel-switch fix | — (merged: `b7ecd04a`) | Edge cask installs the `spacedock` binary, declares `conflicts_with` stable. | In the tag because goreleaser regenerates both casks at cut time. |
| 4 | `docs/releasing.md` reconciliation (consolidated) | `7yd3mbsy2am5qggc17sxvz2v` + `ezn308z0chwc2zvmyny9ry8w` + `qpfmdxy6438fsndp9nw4c89e` | ONE pass folding `releasing-doc-pre-stamp-drift` + `stamp-then-tag-release-ritual` + `steady-state-stable-release-runbook`. Fixes the stale step-3 manual pre-stamp that would block a cutter at the exact-SHA e2e-gate. | All three collide on `docs/releasing.md`; serialize them. |
| 5 | `spacedock-marketplace-source-env` | `z2tjv3570ahjxewv1c309rbc` | Codex EDGE install fix: channel-via-marketplace-name. Edge installs as `spacedock@spacedock-edge` from a distinct `spacedock-edge` MARKETPLACE NAME so codex's entry-name == `plugin.json`-name check passes. Carries the marketplace restructure to expose that distinct name. | Codex edge install is rejected without it. |

### Riskiest-first rationale

Members 1 and 2 carry the value-gate clawback and the SB2-class risk that the trim drops
a load-bearing guarantee — they go first. Member 3 is already merged. Members 4 and 5 are
release-machinery and install-path correctness that must be in before the tag but do not
threaten the gate.

### Member 1 — `fo-contract-token-cut` (the WHY, boot-resident)

The boot-resident clawback. `#418` already landed (the pi adapter over-delivered → the pi
gate is already met). What remains: rebase the `fo-contract-token-cleanup.md` proposal and
apply the verified cuts to `first-officer-shared-core.md` and the runtime adapters, and
fold in the zero-discover "no filesystem sweep" prohibition. Leads the value-measurement
de-risk: pin the `wc -c` baseline numbers and run the no-guidance-control micro-test on the
cut clauses before applying, then re-test the keeps (the control can overturn a keep).

### Member 2 — `trim-dispatch-adapter-prose` (the +5441 B half)

The largest single delta on the gate. SB2-class risk: the cuts MUST keep the
premature-reap ban, the dispatch idle guardrail, and the reuse/await semantics
**byte-intact** — those are load-bearing FO-safety guarantees. The acceptance criterion
gates on a NEGATIVE cumulative line delta and is `contractlint`-guarded. Lands AFTER
member 1 (they both touch resident contract prose; serialize to avoid a collision).

### Member 3 — `#420` cask channel-switch (ALREADY MERGED)

Squash `b7ecd04a` on `main` (`fix(release): edge cask installs spacedock command,
conflicts_with stable`). Listed as a member because the stable tag regenerates BOTH casks
via goreleaser — it just needed to be in before the tag fires. No drive work; verified by
the post-cut fresh-install check.

### Member 4 — `docs/releasing.md` reconciliation (consolidated, serialize)

ONE pass folding three entities that collide on `docs/releasing.md`. The bug: step 3
today does a manual pre-stamp commit (`stamp-version` → commit), creating a fresh SHA
the e2e-gate has never exercised green — so a cutter who follows the doc verbatim tags a
commit the `e2e-gate` job will block (it publishes only after confirming the tagged commit
has a green Runtime Live E2E or a recorded waiver). The reconciliation makes the doc tag
the GREEN e2e commit, not a freshly-stamped one. Because all three entities edit the same
file, they MUST be serialized into a single pass; do not dispatch them in parallel.

### Member 5 — `spacedock-marketplace-source-env` (z2)

Codex EDGE install fix. The standalone-marketplace decouple is a SHIPPED PREREQUISITE (PR
#352 — the standalone `spacedock-dev/marketplace` repo with the `spacedock`/`spacedock-edge`
entries already exists), NOT a live member. z2 carries the remaining codex-edge fix:
channel-via-marketplace-name. Today the standalone marketplace is one marketplace NAMED
`spacedock` carrying two ENTRIES (`spacedock` stable, `spacedock-edge` edge); codex's
entry-name == `plugin.json`-name check REJECTS the `spacedock-edge` entry. z2's fix exposes
a distinct `spacedock-edge` MARKETPLACE NAME so edge installs as `spacedock@spacedock-edge`
and the name-match passes. That marketplace restructure is within z2's scope.

## Deferred (NOT in 0230) — and why

| deferred | why it is out of 0230 |
| --- | --- |
| `notarize-macos-release` | The cask postflight already strips `com.apple.quarantine`, so fresh installs launch. Developer-ID notarization needs an Apple cert — out of scope for a leanness cut. |
| `gate-on-end-value` | Adds resident tokens during a leanness sprint — directly opposes the value gate. |
| `fo-self-evidence-bar` | Adds contract text — same opposition to the gate. |
| `next-independent-release-line` | Edge versioning, not the stable cut. |
| `live-test-stream-tee-ci-stdout-noise` | CI log hygiene → belongs to `0203`. |

## Dependencies & sequencing

- **Member 1 leads** (boot-resident clawback + the value-measurement de-risk). Member 2
  lands AFTER member 1 (both touch resident contract prose — serialize).
- **Member 3 is already merged** — no drive; it is in the tree.
- **Member 4 is a single serialized pass** over `docs/releasing.md` (three colliding
  entities).
- **Member 5** (codex-edge install) is independent of 1/2/4 and can run in parallel.
- **The value gate is the funnel** — every member feeds the cut, but the cut does not fire
  until the gate's byte counts are all met, re-measured against the assembled `main` tip.

## Definition of Done

`v0.23.0` is cut from `main` when:

- **The VALUE GATE is met** — every resident contract file is `wc -c` `≤` its `v0.22.0`
  baseline (the table above), re-measured against the assembled `main` tip.
- Members 1, 2, 4, 5 merged to `main`; member 3 already in.
- The pre-cut antipattern audit (independent reviewer, tag not yet fired) is clean —
  ship-blockers fixed, non-blockers seeded.
- `go test ./...` green from the root.
- The `v0.23.0` release notes cover the `pre.1`/`pre.2` content AND the token clawback.
- The plugin manifests (`.claude-plugin/plugin.json`, `.codex-plugin/plugin.json`,
  currently `0.22.0`) are stamped `0.23.0` on the tagged commit, the tag is on an
  e2e-green commit, and the tagged commit's manifest reads `0.23.0`.
- The tag push fires goreleaser (regenerating both casks with `#420`'s fix), and a fresh
  `brew install` of the released cask launches.

## Sprint lifecycle checklist (owner-tagged)

**Shape — Shaping FO (this session)**
- [x] **Scope-lock** with the captain — members in / deferred ✓
- [x] **Carve** — members stamped `sprint: 0230-stable-finalization` (member 1, 2, and the
  releasing-doc-pre-stamp-drift entity); write this `index.md` ✓
- [ ] **Ideate** any gated member that is not already designed (the token-cut + trim cuts
  carry their proposals; the releasing reconciliation + codex-edge fix carry their entity
  bodies)
- [ ] **⚠️ Preflight staff review (sprint-wide)** — independent reviewer refutes the sprint
  as a whole → `staff-review.md`
- [ ] **Package** — `dispatch-sprint-execution.md` (the cold-boot Commander package — below) ✓

**Drive — Commander (a separate, cold-booted session)**
- [ ] Implementation → validation → done per member; detached adversarial audit at
  validation for the high-stakes surfaces (the shipped FO contract trims; the release
  machinery)
- [ ] Merge each to `main` (PR-merge); state commits concurrency-safe
- [ ] **VALUE-GATE CHECK** — re-measure all resident byte counts vs `v0.22.0`; if any
  exceeds, the tag does NOT fire — bounce back
- [ ] **⚠️ Pre-cut antipattern audit** — independent reviewer over the assembled sprint
  before the tag fires
- [ ] **Cut the release** — `go test ./...` green, stamp `0.23.0`, tag the e2e-green
  commit, push to fire goreleaser, verify the tagged manifest reads `0.23.0` + fresh
  install launches *(captain authorizes)*

**Close — Shaping FO**
- [ ] Seed the next sprint — fold the pre-cut audit's non-blockers + the deferred set
  (notarize, gate-on-end-value, fo-self-evidence-bar, next-independent-release-line,
  live-test-stream-tee → 0203) into the next backlog + a light post-cut release
  verification.
