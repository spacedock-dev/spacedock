# 0230 — stable v0.23.0 finalization (lean cut)

> **Commissioned 2026-06-20 (Shaping FO).** The first sprint whose deliverable is a
> STABLE release, not a feature. `v0.22.0` is the last stable tag; two pre-releases
> (`v0.23.0-pre.1`, `v0.23.0-pre.2`) shipped edge content but accumulated boot/dispatch-
> resident contract bloat. Stable was held precisely because that bloat must be clawed back
> before it ships to every stable user. This sprint cuts a lean `v0.23.0` from `main`.

## Theme

Cut a LEAN stable `v0.23.0` whose resident contract — FO **and** ensign — is back at or
below the `v0.22.0` baseline, recovering the token bloat that accreted through `pre.1`/
`pre.2`, with the `#420` cask channel-switch fix in, codex edge install fixed, a fresh
`brew install` surviving first launch, a reliably-green e2e-gate, and release notes covering
everything the two pre-releases shipped.

## Membership is the query, never a hard-coded list

```bash
spacedock status --workflow-dir docs/dev --where sprint=0230-stable-finalization
spacedock status --workflow-dir docs/dev --where sprint=0230-stable-finalization --where 'sprint-readiness != defer'
```

> **The query returns the 6 tagged entities** (members 1, 2, 4, 5, 6, 7 below). `#420`
> (member 3) is already merged — not an entity. **Member 4 is the promoted single
> consolidated lead:** `releasing-doc-pre-stamp-drift` (`7yd3mbsy2am5qggc17sxvz2v`) carries
> the whole `docs/releasing.md` reconciliation; its two former co-targets —
> `stamp-then-tag-release-ritual` (`ezn308z0chwc2zvmyny9ry8w`) and
> `steady-state-stable-release-runbook` (`qpfmdxy6438fsndp9nw4c89e`) — are **untagged and
> folded into the lead** (their in-0230 substance lives in the lead's ACs; any out-of-scope
> remainder stays backlog). So the Commander drives member 4 as ONE serialized pass over
> `docs/releasing.md`, not three.

## Goal (success criterion)

A stable `v0.23.0` is tagged from `main` and a fresh `brew install` of it launches, where:

1. **The resident contract is lean** — every FO boot-resident contract file AND the ensign
   runtime adapters are `wc -c` at or below their `v0.22.0` baseline (the VALUE GATE below).
   This is the WHY the stable was held; it is the deliverable's gate, not side polish.
2. **The `#420` cask channel-switch fix is in** — a user switching between the stable
   `spacedock` cask and the `spacedock@next` edge cask gets the right binary, the two casks
   declaring `conflicts_with`. (Already merged; in because the tag regenerates both casks.)
3. **Codex edge installs** — edge resolves as `spacedock@spacedock-edge` (member 5); the
   channel-via-marketplace-name fix makes codex's entry-name == `plugin.json`-name check pass.
4. **A fresh `brew install` survives first launch** — the cask postflight strips
   `com.apple.quarantine`, so a Gatekeeper-clean first run works without Developer-ID
   notarization (deferred — it needs an Apple cert).
5. **The e2e-gate is reliably green** — the `rejection-flow` live lane is proven-rare or
   hardened (member 7) so the tag's required green Runtime Live E2E is not a coin-flip.
6. **Release notes cover everything shipped** — the `v0.23.0` changelog names the `pre.1`
   and `pre.2` content PLUS the token clawback, so a stable user upgrading from `0.22.0`
   sees the whole story in one set of notes.

## The VALUE GATE — the definition of done (begin-with-the-end)

**The `v0.23.0` tag does NOT fire until this gate is met.** This is the WHY the stable
release was held: the token clawback is the deliverable. Stop at this gate and check it
BEFORE doing release-cut work — a green test suite and a clean audit do not authorize the
tag if the byte counts still exceed the baseline.

**FO boot-resident contract** — `wc -c` MUST be `≤` the `v0.22.0` baseline:

| resident file (`skills/first-officer/references/`) | v0.22.0 (B) | HEAD at commission (B) | delta | met? |
| --- | ---: | ---: | ---: | --- |
| `first-officer-shared-core.md` | 28586 | 30640 | +2054 | NO — clawback owed |
| `fo-dispatch-core.md` | 17488 | 22929 | +5441 | NO — clawback owed |
| `fo-merge-core.md` | 8059 | 8597 | +538 | NO — clawback owed |
| `claude-first-officer-runtime.md` | 4575 | 4575 | 0 | met (at the line — net-zero) |
| `codex-first-officer-runtime.md` | 6004 | 6043 | +39 | NO — small clawback owed |
| `pi-first-officer-runtime.md` | 3754 | 3754 | 0 | met (#418 over-delivered) |

**Ensign runtime adapters** (member 6) — `wc -c` MUST be `≤` the `v0.22.0` baseline:

| ensign file (`skills/ensign/references/`) | v0.22.0 (B) | HEAD (B) | delta | met? |
| --- | ---: | ---: | ---: | --- |
| `codex-ensign-runtime.md` | 2390 | 2847 | +457 | NO — clawback owed |
| `pi-ensign-runtime.md` | 1768 | 2838 | +1070 | NO — clawback owed |
| `claude-ensign-runtime.md` | 2556 | 2556 | 0 | met (net-zero guardrail) |
| `ensign-shared-core.md` | 8829 | 8829 | 0 | met (net-zero guardrail) |

The release is NOT cut while ANY byte count — FO or ensign — exceeds its `v0.22.0` number.
Re-measure with `wc -c` on these exact paths against the assembled `main` tip immediately
before the cut; if any file is over, the tag does not fire — bounce back to the owning
member. (The gate is human-measured `wc -c` — no in-repo byte-budget guard this sprint.)

> **Measurement de-risk (riskiest-first).** The byte counts above are the gate's authority.
> Before any cut is applied, pin the baseline numbers and the measurement command (`wc -c`
> on the exact file paths) so the gate is checked against a fixed, agreed oracle, not a
> re-derived estimate. Member 1 leads this de-risk.

## Members (riskiest-first)

| # | slug | id | deliverable | the WHY |
| --- | --- | --- | --- | --- |
| 1 | `fo-contract-token-cut` | `y2r7ew51xqs6q3avsb6mcaka` | Boot-resident clawback: apply the verified cuts (`first-officer-shared-core.md` + the FO runtime adapters) from the `fo-contract-token-cleanup.md` proposal; fold in the zero-discover "no filesystem sweep" prohibition. #418 landed → a rebase + the remaining cuts. | The gate's core; recovers the FO resident-contract bloat. |
| 2 | `trim-dispatch-adapter-prose` | `adk755xqeb4a9dxhhgtjwawh` | The `+5441 B` dispatch/merge-core half: trim `fo-dispatch-core.md` (+ the merge-core overage) to capability bindings. SB2-class risk. | The single largest delta on the gate. |
| 3 | `#420` cask channel-switch fix | — (merged `b7ecd04a`) | Edge cask installs the `spacedock` binary, declares `conflicts_with` stable. | In the tag because goreleaser regenerates both casks at cut time. |
| 4 | `docs/releasing.md` reconciliation (consolidated lead) | `7yd3mbsy2am5qggc17sxvz2v` | ONE pass over `docs/releasing.md` (folds the stamp-then-tag + steady-state-runbook substance). Fixes the stale step-3 manual pre-stamp that would block a cutter at the exact-SHA e2e-gate. | The cutter must follow a doc that tags the e2e-green commit. |
| 5 | `spacedock-marketplace-source-env` | `z2tjv3570ahjxewv1c309rbc` | Codex EDGE install fix: channel-via-marketplace-name. Edge installs as `spacedock@spacedock-edge` from a distinct `spacedock-edge` MARKETPLACE NAME so codex's entry-name == `plugin.json`-name check passes. | Codex edge install is rejected without it. |
| 6 | `ensign-runtime-binding-block-cleanup` | `x1khmz0e80fyhe7vnjg8w59y` | Ensign-side clawback: reduce the codex/pi ensign runtime adapters (+457 / +1070 over `v0.22.0`) to host-specific binding blocks (or absorb into the FO adapter); claude-ensign + ensign-shared-core stay net-zero. The ensign analog of t0g/#418. | The ensign half of the resident-contract clawback (gate legs above). |
| 7 | `opus-rejection-flow-reviewer-routing-flake` (7hc) | `7hczkc0c6ezgwy1p627ejp6x` | Establish the opus `rejection-flow` flake frequency (N≥3); if recurring, harden the feedback-routing prose so opus reliably dispatches a SEPARATE reviewer (sonnet already gets it right). | De-risks the green Runtime Live E2E the tag requires; a flaky lane makes the cut a coin-flip. |

### Riskiest-first rationale

Members 1, 2, and 6 carry the value-gate clawback (1/2 the FO contract, 6 the ensign
adapters) and the SB2-class risk that a trim drops a load-bearing guarantee — they lead.
Member 3 is already merged. Members 4 (release-machinery), 5 (codex-edge install), and 7
(e2e-gate reliability) must be in before the tag but do not threaten the byte gate; member 7
is the one that gates the *tag's required green e2e run*, so it lands before the cut.

### Member 1 — `fo-contract-token-cut` (the WHY, FO boot-resident)

`#418` already landed (the pi adapter over-delivered → the pi FO-leg is met). What remains:
rebase the `fo-contract-token-cleanup.md` proposal and apply the verified cuts to
`first-officer-shared-core.md` and the FO runtime adapters (incl. clawing the codex-FO +39),
and fold in the zero-discover "no filesystem sweep" prohibition (bound to the
`detectBroadSearchAtBoot` detector / `TestLiveZeroDiscoverReportsAndStops`, not a prose-grep;
a stochastic zero-discover red is re-run-grounds, never a merge blocker). Leads the
value-measurement de-risk: pin the `wc -c` baselines + the no-guidance-control micro-test on
the cut clauses before applying, then re-test the keeps. Also carries a `docs/runtime-support.md`
alignment AC (the authority doc reflects the final binding-block shape). `claude-FO-runtime`
is exactly at baseline — net-zero guardrail (must not add).

### Member 2 — `trim-dispatch-adapter-prose` (the +5441 B half, SB2-class)

The largest single delta. SB2-class: the in-surface reuse/await text in `fo-dispatch-core.md`
(`## Reuse and Fresh Dispatch` conditions 0-4, the `does not match next stage effective_model`
anchor, supersede-shutdown) and the `«merge.guard»` phase semantics in `fo-merge-core.md` stay
BYTE-INTACT; and `claude-fo-dispatch.md` (the `## Awaiting Completion` / premature-reap ban /
idle guardrail home, OUT of the cut surface) gets a git-diff-ZERO assertion. AC-1 gates on
absolute `wc -c` ≤ v0.22.0 (not a delta-vs-main); `contractlint`-guarded; + a runtime-support
alignment AC. Lands AFTER member 1 (both touch resident prose; serialize — and M2 must not
edit shared-core's `→ shipped (this sprint)` markers, which are M1's).

### Member 3 — `#420` cask channel-switch (ALREADY MERGED)

Squash `b7ecd04a` on `main`. The edge cask installs the `spacedock` binary and declares
`conflicts_with` stable. In because the stable tag regenerates BOTH casks via goreleaser. No
drive; verified by the post-cut fresh-install check.

### Member 4 — `docs/releasing.md` reconciliation (consolidated lead, serialize)

ONE pass over `docs/releasing.md` (+ its byte-identical mirror `docs/site/contributing/
releasing.md`). The bug: step 3 today does a manual pre-stamp commit, creating a fresh SHA
the exact-SHA e2e-gate has never run green on — so a cutter following the doc verbatim tags a
commit the gate BLOCKS. The reconciliation makes the doc tag the GREEN e2e commit, describes
the moving-`stable`-branch auto-advance (not a manual repoint), and carries a divergeable
guard asserting the tagged commit's `plugin.json` version == the tag semver. Verify-first:
`releasing.md` already says "cut from main," so the steady-state fold is likely a DELETION of
stale next-framing, not an addition. The `next` marketplace version-field cleanup is OUT
(standalone repo, unreachable here).

### Member 5 — `spacedock-marketplace-source-env` (z2, codex EDGE install)

The standalone-marketplace decouple is a SHIPPED PREREQUISITE (PR #352), NOT a live member.
z2 carries the remaining codex-edge fix: today the standalone marketplace is one marketplace
NAMED `spacedock` carrying two ENTRIES; codex's entry-name == `plugin.json`-name check REJECTS
the `spacedock-edge` entry. z2 exposes a distinct `spacedock-edge` MARKETPLACE NAME so edge
installs as `spacedock@spacedock-edge`. AC-1 = a real `codex plugin add` of the edge channel
on the live lane (the failing-today baseline flipped to success); the claude half is a
construction-guarantee + resolver test (a real `claude plugin install` mutates global state
→ out-of-band, not a CI gate). Independent of 1/2/4 — drive in parallel.

### Member 6 — `ensign-runtime-binding-block-cleanup` (ensign-side clawback)

The ensign analog of t0g/#418. The codex/pi ensign runtime adapters grew (+457 / +1070 over
`v0.22.0`) from duplicated shared-core narration; reduce them to host-specific binding blocks
(clarification tool, completion signal, captain comms, shutdown) — or absorb fully into the FO
adapter if the ensign-specific content is FO-side. `claude-ensign-runtime.md` + `ensign-shared-
core.md` stay net-zero (guardrail). AC-3 is the value gate (codex-ensign ≤2390, pi-ensign
≤1768); + a `docs/runtime-support.md` alignment AC (the ensign binding-block shape documented,
aligned with whichever path lands). Independent of the FO clawback — drive in parallel.

### Member 7 — `opus-rejection-flow-reviewer-routing-flake` (7hc, e2e-gate de-risk)

A live-lane reliability member, not a byte-gate member. The `v0.23.0` tag fires only on a
green Runtime Live E2E (the e2e-gate); the opus `rejection-flow` lane intermittently mis-routes
the cycle-2 re-review to the impl worker (fix≠reviewer violation) where sonnet gets it right —
a red there makes the cut a coin-flip. AC-1: establish the opus frequency (N≥3 isolated). AC-2:
if recurring, harden the `feedback-rejection-flow` / `«feedback.route»` prose so opus reliably
dispatches a separate reviewer, proven N≥3; if rare, record the re-run-to-green policy and
close. Land before the cut so the e2e-gate is not a coin-flip.

## Deferred (NOT in 0230) — and why

| deferred | why it is out of 0230 |
| --- | --- |
| `notarize-macos-release` | The cask postflight already strips `com.apple.quarantine`, so fresh installs launch. Developer-ID notarization needs an Apple cert — out of a leanness cut. |
| `gate-on-end-value` | Adds resident tokens during a leanness sprint — directly opposes the value gate. |
| `fo-self-evidence-bar` | Adds contract text — same opposition to the gate. |
| `next-independent-release-line` | Edge versioning, not the stable cut. |
| `live-test-stream-tee-ci-stdout-noise` | CI log hygiene → `0203`. |
| `haiku-drive-validation` (0221) | The layered-FO premise proof; separable from the cut (member 1/2 AC-2 already covers contract-trim safety). |
| `fo-tier-delegation` (0221) | In validation; not pulled into the cut (captain: defer for now). |
| `dispatch-build-flag-form-version-skew` (0221), `dispatch-build-request-file` (0222) | Dispatch-build enhancements; not cut blockers. |

## Dependencies & sequencing

- **Members 1, 2, 6 are the clawback** — member 1 leads (FO boot-resident + the value-measurement
  de-risk); member 2 lands AFTER member 1 (both touch FO resident prose — serialize); member 6
  (ensign) is independent — drive in parallel.
- **Member 3 is already merged** — no drive.
- **Member 4 is a single serialized pass** over `docs/releasing.md` (the consolidated lead).
- **Member 5** (codex-edge) is independent — drive in parallel.
- **Member 7** (rejection-flow e2e-gate de-risk) lands before the cut so the tag's required green
  Runtime Live E2E is reliable.
- **The value gate is the funnel** — the cut does not fire until every FO and ensign resident byte
  count is met, re-measured against the assembled `main` tip.

## Definition of Done

`v0.23.0` is cut from `main` when:

- **The VALUE GATE is met** — every FO boot-resident contract file AND the ensign runtime adapters
  are `wc -c` `≤` their `v0.22.0` baseline (the tables above), re-measured against the assembled
  `main` tip.
- Members 1, 2, 4, 5, 6 merged to `main`; member 3 already in; member 7 resolved (lane proven-rare
  or hardened) so the e2e-gate is reliably green.
- The pre-cut antipattern audit (independent reviewer, tag not yet fired) is clean.
- `go test ./...` green from the root.
- The `v0.23.0` release notes cover the `pre.1`/`pre.2` content AND the token clawback.
- The plugin manifests (`.claude-plugin/plugin.json`, `.codex-plugin/plugin.json`, currently
  `0.22.0`) are stamped `0.23.0` on the tagged commit, the tag is on an e2e-green commit, and the
  tagged commit's manifest reads `0.23.0`.
- The tag push fires goreleaser (regenerating both casks with `#420`'s fix), and a fresh
  `brew install` of the released cask launches.

## Sprint lifecycle checklist (owner-tagged)

**Shape — Shaping FO (this session)**
- [x] **Scope-lock** with the captain — members in / deferred ✓
- [x] **Carve** — the 6 members stamped `sprint: 0230-stable-finalization`; the 2 ritual co-targets
  untagged/folded; this `index.md` written ✓
- [x] **Ideate / gate-review** — staff readiness review run; members refined to value-anchored ACs
  (absolute `wc -c` vs `v0.22.0`), SB2-protection split, runtime-support alignment AC ✓
- [ ] **⚠️ Preflight staff review (sprint-wide)** — an independent reviewer refutes the assembled
  sprint as a whole before the cut (the pre-cut antipattern audit covers this in the drive phase)
- [x] **Package** — `dispatch-sprint-execution.md` (the cold-boot Commander package) ✓

**Drive — Commander (a separate, cold-booted session)**
- [ ] Implementation → validation → done per member; detached adversarial audit at validation for
  the high-stakes surfaces (the shipped FO + ensign contract trims; the release machinery)
- [ ] Merge each to `main` (PR-merge); state commits concurrency-safe
- [ ] **VALUE-GATE CHECK** — re-measure all FO + ensign resident byte counts vs `v0.22.0`; if any
  exceeds, the tag does NOT fire — bounce back
- [ ] **⚠️ Pre-cut antipattern audit** — independent reviewer over the assembled sprint before the
  tag fires
- [ ] **Cut the release** — `go test ./...` green, stamp `0.23.0`, tag the e2e-green commit, push to
  fire goreleaser, verify the tagged manifest reads `0.23.0` + fresh install launches *(captain
  authorizes)*

**Close — Shaping FO**
- [ ] Seed the next sprint — fold the pre-cut audit's non-blockers + the deferred set (notarize,
  gate-on-end-value, fo-self-evidence-bar, next-independent-release-line, live-test-stream-tee,
  haiku-drive-validation, fo-tier-delegation, the dispatch-build pair) into the next backlog + a
  light post-cut release verification.
