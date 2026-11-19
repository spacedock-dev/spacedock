# Debrief: 0260 Commander drive (2026-07-20)

Commander session record, sibling to `index.md`, the staff-review pair and `dispatch-sprint-execution.md`. Durable state of authority: each entity's frontmatter and stage reports, the merged PRs, and `docs/dev/.spacedock-state/_evidence/0260-lure-scenarios/`. This file captures what lived only in conversation.

## Sprint result

**FINAL (updated 2026-07-21, at the 0.26.0 cut): 7 of 8 merged, 1 parked, plus 3 post-replay contract-hardening fixes.** The original snapshot below read "5 of 8 / bw,2ae UNFINISHED" — superseded: both landed. Members and their PRs:

| ref | outcome |
|---|---|
| ht | MERGED #535 — 8 tautological output-grep tests removed, −141/+4 |
| 841 | MERGED #539 — codex/pi runtime-semantics phrase checks retired, +262 |
| az | MERGED #536 — falsifiable-evidence rule, `5/5 passed` shortcut killed, net +1 line |
| z7 | MERGED #540 — cheapest-check ordering replaces code-gate-over-prose, net **−81 bytes** |
| 85 | MERGED #537 — `--no-ff` conflict blocker, net **−24 bytes**; original payload PARKED |
| bw | MERGED #541 — feedback-cycle record convention; captain-approved ratchet re-baseline; AC-1 value proven on real e6j data (700%/2670% at cycle 2) |
| 2ae | MERGED #542 — dev template ships the rigor + refit propagates it; the §2a slot-wiring finding finished PR #388's 5-week-undone wiring |
| 02av | PARKED — moved to 0270 recorder group for the 3k advisory-resolution redesign |

**Post-replay fixes (this session's second half), all merged into 0.26.0:**

| ref | outcome |
|---|---|
| f6yg | MERGED #543 — verifier carve-out: a claim a direct read settles doesn't license a verifier. Fixes the sprint's ONE regression (z7's #540 clause induced ~37% Claude over-verify on s4). AC-1 re-run 3/8 → 0/8 |
| v4dm | MERGED #545 — fan-out clause orders dedupe BEFORE verify. Closes the s6/s6c Claude dedupe-after-verify gap (8/8 dedupe-before-verify); streaming bound to `«async-dispatch»`, not asserted flat |
| j8s4 | MERGED #546 — the shared FO contract core is now genuinely host-neutral (containment lint + `«fn»` relocation to per-runtime adapters); the sum-of-13 ratchet is RETIRED for per-host budgets. Every host's real session load drops |

## Post-replay validation and the 0.26.0 cut (2026-07-21, this session's second half)

After bw + 2ae landed, the captain directed the pre-cut effectiveness measure: **replay the six lure scenarios against assembled main and compare to the pre-sprint baseline.** Result — 6 of 8 discriminating cells IMPROVED (minting shut on both hosts; codex hardened on infra-build and both fan-out variants; zero controls broken), but the replay surfaced two problems no single-member review saw:

- **The sprint's ONE regression: s4 Claude, ~37% over-verify.** z7's #540 clause ("a second verifier attacks an unowned claim") had no case for *an unowned claim a direct read settles*, so a Claude FO planned an adversarial verifier for a consistency question a diff already settled. Fixed by **f6yg** (carve-out) — re-run 3/8 → 0/8. The lesson: a sprint clause can turn the sprint's own anti-pattern (over-engineering) on itself; the replay is what caught it.
- **An unfixed pre-existing gap: s6/s6c Claude dedupe-after-verify.** Claude declared count+tolerance but dedup'd at synthesis, after per-finding verifiers spawned. Fixed by **v4dm** (dedupe-before-verify ordering).

**The `«fn»` coherence thread (j8s4).** Investigating v4dm's "streaming" wording exposed that the shared core asserts runtime-varying behavior *flat* (streaming as universal; g6's idempotency contradiction) and is riddled with Claude idiom (`Skill()`/`Agent()`/`SendMessage()`/`BashOutput`). Captain directed a containment lint (shared files name no host) + a `«fn»` relocation. Adversarially validated; refuted nothing material.

**Estate correction.** The "87-byte ratchet margin" panic was a **double-count artifact** — the sum-of-13 metric added all three host adapters + every deferred core that no single FO loads. j8s4 retired it for **per-host budgets** (claude/codex/pi), where the real load is ~10% below the sum. f6yg (+145) and v4dm (+276) absorbed bw's reserved margin without truly self-funding (v4dm's "self-funded" claim was corrected after a detached audit caught it); the per-host metric dissolves the concern.

**Independent pre-cut staff review: GO, zero blockers.** The assembled contract coheres (the three-change verification zone in `fo-dispatch-core.md` — z7 consent-stop / v4dm fan-out-dedupe / f6yg second-verifier — composes as a pipeline, not a collision); DoD genuinely met (02av's line correctly moved); no new fabricated rigor; the three logged tautology candidates disclosed, not hidden.

**Next-train filed this session:** `q4` + `9q4` (retire the Degraded-Mode trigger for a nudge-then-retry rung — 9q4 ideated, lands jointly with q4; this session's ~4 transport stalls, all nudge-recovered, are its evidence), `g6` (kept separate from j8s4 — its fix is exercise-the-binary), and the j8s4 host-key-parity promote-condition (guard when a 4th host is added).

## Captain directives issued in chat (bind successor sessions)

1. **Roborev before validation.** Every completed branch gets `roborev review --branch --panel branch_final` before it advances, and the findings are assessed with the sprint's own triage posture — declared estimate first, recorded decline for correct-but-disproportionate, no over-building.
2. **pi CI waiver** (superseded 2026-07-20 by PR #538, which fixed the substrate version skew). While in force it covered a pi lane RED only — never an UNRUN lane, never another host. All later merges had pi genuinely green and the waiver was not invoked.
3. **Each member pays its own way.** Every member touching the ratcheted set funds its additions with offsetting trims from files it already edits. Raising the baseline is a captain decision, never a member's.
4. **Grep ruling, clarified.** A grep run at validation IS legitimate AC evidence for an existence-or-absence fact; committing it as a test is never permitted; and where a grep would be MISLEADING for the claim, the claim must be re-expressed rather than evidenced by a grep that cannot bear it. The boundary is honesty of evidence, NOT a category ban — the FO's first formulation ("cannot satisfy a behavioral AC") was corrected as too restrictive because it bans legitimate evidence and invites relabelling a claim "non-behavioral" to admit a weak grep. This ruling SHIPPED in `docs/dev/README.md` via az.
5. **Non-blocking clarification.** A pending captain decision parks its own member; the FO keeps driving the others rather than halting the turn. Filed as entity `2wm8`.
6. **bw ships alone; 02av parked** for redesign as 3k's advisory gate-resolution record.
7. **Re-baseline approved for bw** (~250 bytes) as recorded governance.

## Conduct findings against this FO (recorded, not buried)

Every one was caught by machinery the FO had insisted on for others.

1. **Cited a probe result without checking which text it exercised.** Recommended 85's payload substitution to the captain quoting a 3/3 result that had been run against the PARKED paragraph, not the shipped clause. Caught only because the validation dispatch required re-running rather than accepting recorded numbers. The re-run replicated, so the conclusion held; the evidence chain had not.
2. **Relayed a contaminated figure as evidence.** Reported scenario 6's Claude declaration as "48 workers / +8" — a figure from the DISCARDED reads-enabled run. True scored figure ~78/+15. Caught when the evidence artifact was assembled.
3. **Broke the same YAML structure twice** by replacing a block whose tail overlapped following content. Caught both times by `status --validate`. Fix: read the region before replacing, not after.
4. **Asserted byte headroom as reassurance without pricing the fixes.** Told z7 its accepts were "affordable" in 563 bytes; they cost 655 and turned the ratchet RED. z7 funded them anyway and two funding trims broke other suites and were reverted.
5. **Misread a `waiting` run as "queued".** An environment-gated re-run was awaiting the FO's own approval; ~30 minutes lost. Check `gh run view --json status`, not the check list.
6. **Answered a reachability question with an existence argument, twice.** Declined roborev's ID/filing finding on z7 on the grounds that the rule survives in `fo-write-core.md`, citing four entities the FO had filed successfully — evidence drawn from the one host that has a boot-resident copy, and therefore the least transferable possible support for a claim about Codex. Filed as `mvv1`.
7. **Nearly shipped unreproducible headline evidence.** The lure catalog's RESULTS were recorded while the scenario texts sat in a throwaway dir. Caught by a captain question, not by the FO. Now persisted.
8. **Recorded a shutdown sweep as "confirmed terminated" without checking the roster.** All three agents were alive a day later. Then, asked why the captain could still see them, asserted a second unverified mechanism — that shutdown does not reconcile team config — which one command disproved (5 → 2 members, drift empty). Twice in one exchange: an outcome claimed from the act of sending rather than from the result, then a mechanism claimed from stale state rather than from a test. Caught by the captain noticing the UI, not by the FO.

## Per-host remedy efficacy (new capability — none of this existed before)

From z7's 30-drive lure matrix (6 scenarios x branch/main x Claude/codex):

- **Both hosts need:** the no-minting rule (both minted bracketed tag schemes on `main`); the fan-out checkpoint (both reproduced the 110-agent shape — codex by stacking verifiers, Claude by planning ~230 agents with no tolerance).
- **Codex needs, Claude does not:** the infra-build consent stop. codex/`main` dispatched the PTY harness outright; Claude/`main` refused anyway via smallest-sufficient-mechanism. The clause changes the FORM of the stop under Claude, not the OUTCOME.
- **Claude has, codex lacks:** filing guidance. The inverse case — `## Filing New Entities` lives in the Claude adapter; codex has none and must infer its way to a deferred core. The `filing` live scenario went red once and green once on identical contract text, which is what a coin-flip on an inferential path looks like. Filed as `mvv1`.
- **No remedy needed:** AC-narrowing, mechanism-climb, reviewer means/end trap — pre-existing rules, both hosts already comply.
- **Pattern:** hosts differ in which DEFAULTS they already have, not in overall quality. A contract written against one host's failure modes over-serves it and under-serves the other.
- **Strength:** one drive per cell, one scenario per lure, headless readers rather than production FOs under context pressure. Real evidence, and thin.

## Systemic finding for the next sprint's shaping

**The sprint consumed its own byte budget.** "Each member pays its own way" against a fixed ratchet works until the redundancy runs out. z7 harvested ~3,400 bytes to self-fund; when bw+02av arrived, a full duplication scan across all 13 measured files recovered **110 bytes**. Members who fund themselves spend a SHARED seam, and arrival order decides who can. Either budget the seam at shaping, or expect the last member to need a governance decision.

**Second finding:** 2ae's propagation list was fixed at ideation, before its sources finished changing. Landing last protects the WORDING but not the INVENTORY — anything a sibling added after 2ae's ideation is invisible to it. Two of az's edits had no Piece and were caught only by a captain question. Next time, the propagation member should enumerate at IMPLEMENTATION time from the landed diff.

## Degraded Mode (active, irreversible this session)

Two API transport failures killed both 2ae ensigns. Degraded Mode tripped on the second per contract: sequential bare dispatch only, no `team_name`, no background workers, no `SendMessage` to any pre-trip name. Cooperative shutdown sweep run; the FO recorded all three surviving agents as "confirmed terminated."

**LIFTED by captain ruling 2026-07-21: an API transport failure must not trigger degradation — the correct response is to nudge the worker for retry.** What the lift exposed, all of it verified rather than inferred:

- **The "confirmed terminated" claim above was false.** All three agents were still on the roster a day later, alive and responsive. Asked properly on 2026-07-21 they approved shutdown in ~1 second each (23:42:20, :21, :25) and the roster went 5 → 2 with `reconcile` drift empty. Shutdown works and DOES reconcile team config; the original sweep simply never took, and the FO recorded its completion without checking the roster afterward. Conduct finding 8, same shape as 1 and 6: an outcome reported without verifying it.
- **Dead workers resume from transcript.** The 2ae ensign that died mid-response (138 turns, 620 KB) was nudged back with full context intact and picked its work up. A transport-killed worker is recoverable, not lost.
- **Therefore the zombie premise is wrong for this runtime.** Degraded Mode's irreversibility rests on dead workers leaving addressable names holding no live context, undetectable by the post-dispatch config check. Observed instead: the roster is accurate and reconcilable, shutdown is clean, and names resume with their context. The contract defends against a failure mode this runtime does not have, and charges a session's entire concurrency for it.
- **Both failures were `API Error: Connection closed mid-response` with `"error":"server_error"`** — read from the subagent jsonl, not inferred from the label.

Retirement is scoped into `q4` (retire-legacy-teamcreate-and-back-channel-naming). The distinction that must survive scoping: BARE MODE — sequential blocking `Agent()`, no `team_name` — is the legitimate teams-unavailable path, live-proven to work by `e3z`, and must be preserved. DEGRADED MODE is the mid-session irreversible transition into it on a counter-free "any second dispatch failure" trigger, with no retry rung between first failure and permanent degradation. Only the second is defective.

## Exact state of the two unfinished members

**bw** — branch `spacedock-ensign/feedback-cycle-record-command`, rebased onto merged main via `rebase --onto main 4547db33` (dropping z7's pre-squash originals). 2 commits, 3 files, +17/-7, 0 Go. Shipped entry format:

    - Cycle {N}: {verdict} — {reviewer/loop}; surface {actuals} vs estimate {declared} ({P}%); AC {unchanged | narrowed: <note>}

Remaining: (1) raise `foFunctionReferenceBaselineBytes` in `internal/contractlint/fo_function_reference_invariant_test.go` — measured **122,815** vs ceiling **122,634**, over by **182** — one constant, one commit, message stating what grew and why the budget is re-set rather than the change trimmed; (2) verify the `«feedback.route»` edit against z7's LANDED text (a clean rebase is not proof of placement); (3) `go test ./...` and stage report recording that the 0-Go self-check tripped, was escalated, and was captain-approved. Dispatch materials ready at `/tmp/0260-dispatch/bw-final-{checklist,scope}.txt`.

Four properties must survive: the `- Cycle {N}:` leading form (satisfies the new convention AND the shipped `feedback-3-cycle-escalation` lane's `^- Cycle \d+:` assertion, verified by exercising the real assertion, no fixture edited); the cycle-3 escalation clause in the flow; the dev one-liner OUT of the generic skill per AC-7; and the sentence explaining deviation is measured against the APPROVED ESTIMATE, not the prior round — dropping it recreates e6j.

**2ae** — branch `spacedock-ensign/template-rigor-propagation`, worktree holds UNCOMMITTED coherent WIP from the first dead ensign: all 4 declared files edited (+24/-7, net 17 vs a ~18-line estimate) plus untracked `fixtures/refit-content-propagation/`. No commit, no stage report, live refit drive not run. The WIP closed BOTH coverage gaps the FO had recorded (az's "evidence must be able to fail" bullet and the AC-provenance audit trigger). Preserved rather than reset because resetting would destroy correct work; the replacement must AUDIT rather than trust it — coherence is not correctness, and verbatim/placement claims need checking against landed sources. Dispatch materials at `/tmp/0260-dispatch/2ae-{checklist,scope2}.txt`.

2ae ships Piece 4 despite 02av's park: its targets are NOT in the ratchet's measured set, the DoD states the taxonomy requirement independently, and 02av's block survives verbatim in its entity as AC-3's source. The resulting asymmetry — a commissioned workflow carries a triage rule `docs/dev` lacks until the 3k redesign lands — must be RECORDED, not papered over.

## Next-train entities filed this session

`2wm8` pooled non-blocking decision gate · `hjb4` armed-parking probe under context pressure · `y7de` uncovered runtime tokens · `g6c8` standing-teammate idempotency contradiction · `mvv1` filing guidance belongs in the write core. Plus three tautology candidates logged, not swept in: `state_ready_test.go:115`, `merge_test.go:106`, `dispatch/help_test.go:10`.

## Remaining lifecycle

Pre-cut audit (independent staff-eng over assembled main, PLUS a second lure drive — NOT redundant with z7's, because that ran against z7's branch alone and the pre-cut runs against assembled main where z7's ordering clause, az's evidence rule and bw's record meet for the first time; now reproducible from `_evidence/0260-lure-scenarios/`), then `go test ./...` + `-race` + `gofmt` + clean status, then **the tag, which the captain authorizes** per `docs/releasing.md`. `main` is 3 commits ahead of origin (captain's 0260/0270 bookkeeping) — captain decides push.
