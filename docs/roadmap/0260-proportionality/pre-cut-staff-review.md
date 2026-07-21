# 0260 proportionality — independent pre-cut staff review

## Verdict

**NOT READY.** The assembled product at `4f669e6f32081bc68d210ab35f8d57dab21cc75e`
passes the repository-wide deterministic, race, formatting, and merged-PR live-lane
checks I inspected. The three post-replay additions are narrow, compose coherently,
and close the gaps their own targeted evidence found. The release gate is nevertheless
open on three proof failures: the feedback-cycle DoD has no live replay, the template
member has no independent validation record or required human review, and the
captain-assigned final-candidate lure drive is still externally owned and pending.

Finding count: **3 ship blockers, 0 additional Material findings, 0 Needs-decision
items, and 6 recorded non-blockers/declines**.

## Audit identity, range, and tag state

- Audited product tip: `4f669e6f` (`Host-neutral contract core: containment lint +
  per-host budget + «fn» relocation`, PR #546).
- Audited assembled range: `5dac2d6a..4f669e6f`. `5dac2d6a` is the closed preflight
  package immediately before the first 0260 product merge; the member/product merge
  commits in scope are `bdf39f01`, `c240d49e`, `2b3e0287`, `45f54678`, `70daeb65`,
  `a3cce638`, `41880e36`, `effbc769`, `f4188af2`, and `4f669e6f`. I also inspected
  the intervening Pi support merges because they affect host-lane interpretation.
- The local `main` tip is later than the audit tip and carries unrelated
  durable-decisions roadmap work. Commits after `4f669e6f` are excluded from the
  0.26.0 product verdict.
- Stable-tag state, checked both locally and at `origin`: **`v0.26.0` does not exist**.
  The only matching remote tag is `v0.26.0-pre0` (tag object `76f6f613`, dereferenced
  commit `601c3f53`). This is still a pre-cut review.
- The release candidate is expected to receive the normal manifest/minor stamp and an
  exact-SHA Runtime Live E2E green before tagging, per `docs/releasing.md`; their absence
  at this pre-stamp audit tip is not a defect.

## Scope and method

I treated this as an assembled-system refutation, not a second validation pass over
each member. I read the sprint index, both preflight seats, Commander package and
debrief, shaping debrief, all scoped entity bodies, the latest relevant implementation
and validation reports, merged diffs, PR check state, and the committed lure evidence
needed to understand the three additions. I inspected shared contract sections after
all merges rather than accepting clean rebases as composition proof.

The authoritative archived done set covered here is `85`, `ht`, `bw`, `z7`, `az`,
`841`, and `2ae`. `02av` is accounted for separately because the captain removed it
from 0260 and parked/reframed it onto the durable-decisions recorder group. The
post-replay additions outside the archived membership query are `f6yg`, `v4dm`, and
`j8s4`.

I did not rerun member-owned green checks. The assembly-level checks I did run were
from a temporary `git archive 4f669e6f` outside the repository:

```text
go test ./...          PASS
go test ./... -race    PASS
gofmt -d ./cmd ./internal
                       CLEAN
```

The temporary tree was removed after the run. GitHub PR evidence shows all four live
lanes green on PRs #540, #541, #543, #545, and #546. PR #536's Pi failure occurred
under the then-active Pi-only waiver; PR #538 repaired that substrate skew, and the
later all-host greens prevent the waiver from silently carrying the assembled result.

## Ship blockers

### 1. The central feedback-cycle DoD has arithmetic evidence, but no live replay of the decision barrier

**Confirmed fact.** `docs/roadmap/0260-proportionality/index.md`, Definition of Done,
requires a live replay that records reconfirm/re-scope/park/escalate before a third
repair dispatch. `bw`'s final validation report in
`docs/dev/.spacedock-state/_archive/feedback-cycle-record-command/index.md` records a
one-off reconstruction of e6j's history: the deviation is visible at round 2, with
negative controls for missing fields and prior-round arithmetic. That validly proves
the numbers and the fixed-baseline rule. It does not exercise a live FO taking the
next-action decision, persisting the decision, and withholding the next repair
dispatch until that record exists.

The same report explicitly records that the AC's proposed checked-in fixture was not
committed. The captain's no-new-check ruling explains that choice and I do not require
a committed test. It does not erase the sprint index's separate live-replay promise.
No committed report or evidence file supplies that outcome.

**Why this blocks the cut.** The value claim is behavioral: narration must change the
loop's next action. The existing evidence proves only that a reader can compute the
threshold. Shipping on that proof would repeat the sprint's own rejected pattern of
proving the mechanism one seam before the promised behavior.

**Closure.** Run one live FO replay against the assembled candidate (or its release-
stamped descendant with no contract change), using the archived e6j round shape. Keep
durable evidence of the resulting entity/body or event record showing: round 2 exceeds
tolerance; a named reconfirm/re-scope/park/escalate decision is persisted; and a third
repair dispatch does not occur before that record. A report assertion without the
observable order is insufficient.

### 2. `2ae` merged and archived done without its validation gate or the captain-required human review

**Confirmed fact.** The archived entity
`docs/dev/.spacedock-state/_archive/template-rigor-propagation/index.md` ends with the
implementation reports and pre-validation roborev triage. It has no `Stage Report:
validation`, no validation gate record, and its current gate pointer still names the
ideation gate. Its attempt-1 approval carries a binding condition: at validation,
present the refitted workflow-README delta for human review. PR #542 merged as
`41880e36`; GitHub records no PR review, review comment, or issue comment, and the host
live checks remained waiting. The implementation reports contain useful dispatched
commission/refit drives, but they are not the independent validation or human gate the
captain required.

**Why this blocks the cut.** This is both Sprint DoD evidence and a direct captain
condition. The committed fixture proves there is something reproducible to validate;
it does not prove that a fresh validator drove Phase 3b or that the captain reviewed the
emitted delta. Archive status cannot substitute for the missing gate record.

**Closure.** On the assembled candidate, dispatch a fresh validator to drive both the
commission path and refit Phase 3b against `fixtures/refit-content-propagation/`, record
the branch/control content delta and falsifying change, and present the actual refitted
README delta to the captain through the recorded review path. Append the validation
report and durable approval to `2ae`'s workflow record. Because the product is already
merged, this may close as evidence-only if it passes; any product correction must land
and re-enter the affected pre-cut checks.

### 3. The mandatory final-candidate lure evidence is pending outside this review seat

**Confirmed fact.** The Commander package requires a final assembled lure drive. By
captain direction, this staff-review seat skipped duplicate execution of the six-
scenario/two-runtime drive because a separate Claude reviewer owns it. Hegel left no
durable report, transcript, or outcome recoverable after the session crash; I make no
claim that Hegel completed any lure cell.

The committed z7 matrix and the targeted `f6yg`/`v4dm` replays are valid evidence for
their respective branch claims. They are not the missing final-candidate rerun after
`f6yg`, `v4dm`, and `j8s4` compose on `4f669e6f`. I inspected them only to understand
the fixes and did not relabel them as assembled-main proof.

**Closure.** The externally owned Claude review must leave a durable assembled-tip
six-scenario/two-runtime result. Any material regression routes to a product fix and a
focused recheck of the affected seam; an acceptable result closes this external gate.

## Material findings

None beyond the three proof failures above. I found no additional assembled product
defect that independently blocks a value AC.

## Needs-decision items

None. The decisions most likely to be mistaken for open questions are already settled:
`02av` moved; `az` ships no new prose gate; `bw` remains prose-only; `j8s4` keeps the
standing-teammate idempotency conflict separately owned; and no stable tag is authorized
by this report.

## Recorded non-blockers and declines

### 1. `02av` is moved, not silently missing

The captain parked the 0260 implementation, withdrew every shipped `findings` field and
standing triage block, changed the sprint DoD line to MOVED, and re-opened `02av` under
`sprint: durable-decisions`, group `recorder`. Its current body designs the decline as
an advisory resolution and keeps generic recorder plumbing out of its scope. This is
not a 0.26.0 failure. Promote to material if 0.26.0 release notes claim the ensign-side
decline mechanism shipped, or if any withdrawn `findings`/standing-block text reappears
in the candidate.

### 2. `85` deliberately ships a different, narrower value than its original title

The original “armed is not a stopping point” payload is explicitly parked after its
careful-reader probe failed to discriminate and it could not fund its bytes. PR #537
ships the `--no-ff` conflict surface-and-stop rule instead, at net negative prompt
surface, with all host lanes green. The body preserves the abandoned claim honestly.
Promote if a context-pressure drive again parks at `armed`; the named follow-up then
owns a better instrument.

### 3. `ht` leaves top-level help wording intentionally unguarded

The ten-mutation validation proves the retained behavior checks still red on behavior
breaks; it also proves five help-content edits survive. That is the approved behavior-
versus-wording boundary, not hidden coverage. Promote on an observed user-facing help
regression or if a machine consumer begins parsing top-level help.

### 4. The commissioned template outruns its parent on finding triage

`2ae` carries the three-class materiality taxonomy even though `02av`'s parent-workflow
delivery moved. The sprint DoD independently requires that taxonomy, and no withdrawn
Feedback Cycles entry format ships with it, so I do not block on the asymmetry. Promote
if a commissioned workflow needs to persist a decline before the recorder-backed
advisory-resolution path exists, or if the template claims that record path already
ships.

### 5. `841` records runtime-token and declaration-versus-wire limits instead of inventing owners

The final mechanical enumeration leaves nine tokens uncovered and two merely named,
not asserted. Its bindings prove Go declarations, not runtime wire behavior; later
reports correct that language. This is honest debt and the exact-tip suite stays green.
Promote when a supported claim depends on one of those tokens, or when a runtime change
shows a declaration binding can green while the wire behavior regresses.

### 6. `j8s4` defers host-key parity until a fourth host exists

The per-host ratchet has no exact-key-parity guard across its three maps. All three
current hosts are present; the union/topology guard and exact-tip tests cover today's
set. Adding the parity mechanism now would guard only a future host. Promote exactly
when adding another runtime host, as the validation report records.

## Post-replay additions and proportionality

### `f6yg` — PR #543, merge `effbc769`

The gap was real and introduced by z7's first wording: an unowned claim that a direct
read settles still licensed a verifier. The one-clause carve-out routes that case to the
read itself while keeping judgment, runtime behavior, and the mandated detached audit
eligible for adversarial work. Its targeted Claude replay moved from 3/8 verifier
climbs to 0/8, and all four PR live lanes are green. It adds no check, schema, or
standing process. At the assembled tip, the clause remains distinct from the fan-out
checkpoint and j8s4's host-binding rule.

### `v4dm` — PR #545, merge `f4188af2`

The gap was ordering, not missing dedupe: Claude spent verifiers while results streamed
and collapsed duplicates only afterward. The clause now requires a barrier before the
per-finding verifier spawn and binds streaming language to `«async-dispatch»` being
async. Targeted evidence is 8/8 Claude dedupe-first and 4/4 Codex non-regression, with
all PR live lanes green. The later report correctly retracts “self-funded”: the change
uses 276 bytes of existing margin. That accounting correction is durable and j8s4 then
replaces the old summed ratchet, so no hidden 87-byte-margin claim survives as the
current accounting model.

### `j8s4` — PR #546, merge `4f669e6f`

The final core reads coherently in order: declare fan-out count/tolerance; collapse
duplicates before verifiers; author the plan against this host's async/addressability
bindings; and use a second verifier only for an unowned claim a direct read cannot
settle. No clause reopens the earlier lure gaps.

Its standing machinery is proportionate to the observed class. The containment lint is
a structural location check over five shared files, with a discriminator, replacing
flat host claims that had already produced multiple contradictions; it does not assert
runtime behavior from prose. The per-host byte ratchet replaces the inaccurate
sum-of-all-adapters metric rather than adding a second budget. The dispatch change and
thirteen goldens are a mechanical `BashOutput` to `background-task` ownership move. No
binary command, schema, runtime policy, or recurring live lane was added. PR #546 has
all four live lanes green, and the exact-tip deterministic/race/format gates pass.

## Per-member and addition coverage

| Member | Durable outcome inspected | Assembled assessment |
|---|---|---|
| `85` | archived done; PR #537 / `70daeb65`; validation cycle 2 | Narrow substituted payload is honest and green; original value parked, non-blocking. |
| `ht` | archived done; PR #535 / `bdf39f01`; 10-mutation validation | Eight named tautologies fixed; retained behavior coverage and known help-wording gap recorded. |
| `bw` | archived done; PR #541 / `a3cce638`; final implementation and validation | Prose convention composes and lanes are green; missing live next-action replay is blocker 1. |
| `z7` | archived done; PR #540 / `45f54678`; 30-drive validation | Ordering/consent/minting/fan-out core landed; three targeted follow-ons inspected below. |
| `az` | archived done; PR #536 / `2b3e0287`; validation | Falsifiable-evidence rule and existing-audit trigger land without a new prose gate; later Pi fix and host greens close the waived-lane concern. |
| `841` | archived done; PR #539 / `c240d49e`; validation cycle 2 | Runtime phrase assertions retire into declaration/structural bindings with uncovered claims recorded honestly. |
| `2ae` | archived done; PR #542 / `41880e36`; implementation reports only | Product diff is coherent, but validation and captain human-review condition are absent: blocker 2. |
| `02av` | active under durable-decisions, ideation; 0260 park record | Explicitly moved/parked; no 0260 product bytes or decline DoD claimed. |
| `f6yg` | archived done; PR #543 / `effbc769`; targeted replay report | Direct-read carve-out closes the observed s4 climb without new machinery. |
| `v4dm` | archived done; PR #545 / `f4188af2`; targeted replay report | Dedupe-before-verify ordering closes the observed Claude default; funding wording corrected. |
| `j8s4` | archived done; PR #546 / `4f669e6f`; detached audit report | Host-neutral core, containment, and per-host budget compose cleanly; all live lanes green. |

## What held under refutation

- The assembled contract does not contain the retired “prefer a code gate over a
  prose-only rule” or “5/5 passed is sufficient” instructions on shipped surfaces.
- Shared edits did not overwrite one another: z7 and bw occupy separate
  `first-officer-shared-core.md` regions; f6yg, v4dm, and j8s4 form distinct clauses in
  `fo-dispatch-core.md`; the exact-tip preservation and topology suites pass.
- The prompt-surface accounting is no longer a misleading sum of mutually exclusive
  adapters. Every current host has a measured load and a ratchet; a single-host plant
  is discriminated; the union remains tied to the thirteen-file address-lint set.
- Host-specific behavior moved to adapters without a stranded core `«fn»`; the new
  containment check and existing adapter set-equality checks cover opposite failure
  directions, and PR #546's Claude/Codex/Pi lanes all passed.
- The template/refit product diff carries the required taxonomy and fixed Verified-by
  text, and the fixture makes its behavior reproducible. The blocker is missing
  independent/human proof, not a contradiction found in the assembled text.
- No new binary command, schema, daemon, CI lane, or recurring verifier was smuggled in
  under proportionality language. The only new standing check in the additions is the
  structural host-containment lint, justified by a recurrent cross-host claim class.

## Exact closure condition

This verdict becomes **SHIP-CLEAR** only when all three conditions hold on the final
candidate lineage:

1. The live e6j feedback replay leaves durable evidence of the decision record before
   any third repair dispatch.
2. `2ae` has an independent validation report and a durable captain approval that
   includes the human-reviewed refitted README delta.
3. The separately owned Claude lure review leaves an acceptable durable assembled-tip
   six-scenario/two-runtime result; Hegel is not used as evidence.

After any product change made to close a failed condition, rerun only the checks mapped
to that changed seam plus the repository completion gates. With those conditions met,
the current assembled product findings do not require another broad staff-review pass;
the normal release-stamp, exact-SHA Runtime Live E2E, manifest gate, and captain tag
authorization remain the release procedure.
