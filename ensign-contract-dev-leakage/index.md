---
id: ep0ra3zjf4hhkhx5rrkwsxbb
title: Universal ensign contract has absorbed dev-workflow assumptions (TDD, code-only deliverables, "CODE only" worktree) — re-home dev policy out of the shared core
status: ideation
source: session-10 detached audit (Task waymqcmru, 2026-06-03) — overall verdict MATERIAL-PRESENT; validated the captain's grilling instinct ("why am I seeing dev-workflow specific ones like TDD in the universal contract?")
score: "0.19"
worktree:
started: 2026-06-03T07:09:59Z
completed:
verdict:
issue:
---

The shared ensign contract (`skills/ensign/references/ensign-shared-core.md`) is loaded by **every** dispatched ensign across all three commission templates (development, experiment, refinement). A detached adversarial audit confirmed it has absorbed dev-workflow-specific discipline — TDD authoring, a code-shaped closed deliverable list, and a "CODE only" worktree assumption — and asserts them as universal. The FO contract was already corrected to deliverable-shape-agnostic phrasing (`first-officer-shared-core.md` lines 342-353); the ensign contract is the **uncorrected twin**.

## Problem

`ensign-shared-core.md` lines 22-27 carry a `## Working Practices` block that duplicates dev-workflow policy verbatim into the portable contract. The project's own docs already tag this same policy as dev-only opt-in:
- `docs/dev/README.md` line ~110: "This is dev-workflow policy: an AC's proof here is code/command/state. A non-development workflow's AC proof may legitimately be a published artifact, a metric, or a human review."
- `skills/commission/references/templates/development.md` lines ~111-117: a "Recommended practices (opt-in)" carve-out stating "the universal first-officer contract does not impose them."

The shared ensign core strips that qualifier and presents the rules as universal ensign discipline, breaking at least two non-dev template classes: refinement entities whose deliverable **is** prose (PRD, outreach reply, content — `refinement.md` says they never touch the repo), and experiment entities whose deliverable is a hypothesis/analysis verdict against pre-registered success criteria.

**Material findings from the audit (Task waymqcmru):**
- **L2-01 (umbrella):** The entire `## Working Practices` section duplicates dev policy into the universal contract. Lift it out; re-home each principle.
- **L1-F1 / L2-02:** TDD bullet (line 24, "Write the failing test first… The test is what the gate judges") asserts test-first as the universal authoring contract — false for refinement (gate judges the artifact body) and experiment (gate judges evidence vs. pre-registered criteria) workflows.
- **L2-03:** "Every task produces a real, checkable change" + its escape clause under-serve prose-only deliverables (refinement/experiment bodies that are checkable via pre-fixed success criteria, not external files).
- **L2-05:** The Split-Root State Contract (line ~35) hard-codes "the worktree isolates **CODE only**" — excludes non-code work product (experiment evidence files, refinement attachments). Same fix needed in **both** `ensign-shared-core.md` and `first-officer-shared-core.md`.
- **L3-1:** Structural twin of the FO's already-corrected four-principle block — the ensign block is the uncorrected version.
- **L4-P2 / runtime adapters (Polish):** `claude-ensign-runtime.md` and `codex-ensign-runtime.md` enumerate "worktree path" as a top-level always-present assignment field, while the shared core already uses the substrate-neutral "workflow location" for the same slot.

## Proposed approach

Mirror the FO-contract correction (`first-officer-shared-core.md` Working Principles, lines 342-353) onto the ensign contract. The exact edits are pinned below against verified current line numbers (ground-truthed during ideation, 2026-06-03).

### Re-homing design — what lifts, and where each principle lands

The `## Working Practices` block is `ensign-shared-core.md:22-27` (four bullets). Re-home, do not delete:

1. **TDD authoring bullet (line 24, "Write the failing test first… The test is what the gate judges").** This is the only principle with **no current dev home**. Lift it into `development.md`'s existing `## Recommended practices (opt-in)` section (line 111) as a new `### Test-first authoring` subsection, phrased identically to how that section already frames the validation-stage disciplines ("recommended, not mandatory — the universal contract does not impose them"). `docs/dev/README.md`'s ideation stage already carries the code-gate-over-prose corollary (line 85); the TDD authoring discipline itself is what's missing and must land in `development.md`.

2. **"Every task produces a real, checkable change" bullet (line 25).** Already homed: `docs/dev/README.md:80` carries this verbatim-equivalent as a dev ideation-stage rule. The lift is a **deletion from the shared core** (its dev home already exists) — no new home to author. Confirm during implementation that the README clause still covers the prose-only escape ("belongs in the roadmap").

3. **"Prove by exercising, not by re-reading" + "No hidden machine dependencies" bullets (lines 26-27).** These are **genuinely universal** (they hold for experiment evidence and refinement artifacts too) — they are not dev-leakage. They stay, but move out of a block named "Working Practices" (a dev-flavored heading) into the existing universal flow. Fold them into the `## Working` numbered list or a neutral `## Proving your work` heading. The implementer decides placement; the invariant is: these two survive in the universal core, the two dev-specific ones do not.

4. **Optional universal restatement.** The shared core MAY retain one discipline-neutral sentence so a non-dev ensign still knows it owes proof: *"The stage definition states the proof your stage owes — a test, a metric, a published artifact, a human review. Satisfy that, not a generic test-first ritual."* This covers TDD (dev), pre-registered criteria (experiment), and acceptance criteria (refinement) without asserting any one tradition. This mirrors the FO core's already-neutral "satisfy the stage definition's proof requirements" framing — adopt the same shape, do not re-invent the four-bullet block.

`experiment.md` / `refinement.md` need **no new homes** — confirmed during ideation: `experiment.md` already encodes "success criteria fixed before evidence is gathered" (hypothesis stage), and `refinement.md` carries draft/review/polish. The re-homing targets are `development.md` (new TDD subsection) and the shared core's own structure (delete two, keep two, optionally add one neutral line).

### "CODE only" → substrate-neutral (AC-3)

- `ensign-shared-core.md:35` — "the worktree isolates **CODE only**" → "the worktree isolates **the deliverable work product only**".
- `first-officer-shared-core.md:270` — "a worktree stage isolates **CODE only**" and "The worktree still owns code: … apply to **code changes** only" → "isolates **the deliverable work product only**" / "apply to **deliverable-artifact changes** only". Keep the rest of the clause (pr-mirror exception, state-checkout path) intact — only the substrate noun changes.

### "worktree path" → neutral location vocabulary (AC-4)

- `claude-ensign-runtime.md:7` and `codex-ensign-runtime.md:7` enumerate "entity, stage, stage definition, **worktree path**, and checklist" as universal assignment fields. Rename the slot to the shared core's neutral term: **"workflow location"** (the term `ensign-shared-core.md:11` already uses). A non-worktree stage (ideation, backlog) has no worktree path but always has a workflow location.
- **In-scope nuance found during ideation:** `ensign-shared-core.md:17` ALSO says "If you were given a **worktree path**…". This usage is **conditional** ("if you were given"), not an asserted-universal field, so it is legitimate and does NOT need changing. The AC-4 oracle must scope to the *universal-field-enumeration* position (the runtime adapters' "fields are: …" lists), not ban the literal "worktree path" everywhere — otherwise it false-fails on the legitimate conditional at line 17.

The shared ensign core should retain only universals: read the assignment, follow the stage definition's proof requirements, commit the entity to the state checkout, signal completion.

## Out of scope

- Rewriting the dev-workflow disciplines themselves — they are correct *for dev*; this is a re-homing, not a policy change.
- Changing the FO contract's already-corrected four-principle block (only its "CODE only" line, shared with the ensign core, is in scope — finding L2-05).
- The experiment/refinement template bodies beyond confirming they already carry the re-homed analogues.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action. Proof is presence/absence over the instruction files (legitimate per the README: "a presence check over instruction files proving they carry a required clause or stay free of a banned token is proof at the claim's own level") plus the lock-in oracle. Ideation refines.

**AC-1 — The universal ensign shared core no longer asserts dev-only authoring discipline.**
Verified by: a Go oracle test (extending the proven `internal/hostneutrality/prose_inflator_locks_test.go` pattern) asserting `ensign-shared-core.md` is free of the dev-only test-first vocabulary — the exact banned literals are "failing test", "feature or bugfix", and "the test is what the gate judges". The test fails if the Working Practices TDD bullet is reintroduced. Negative-proof shape: the implementer inserts the banned phrase, watches the test go red, removes it, watches green (the same lock-in-is-real demonstration `TestNoAuditTrailExposition` uses).

**AC-2 — The dev disciplines survive in their dev-specific homes.**
Verified by: the same (or paired) oracle asserting the *positive* presence of the re-homed guidance — `development.md` carries a "test-first" clause in its `## Recommended practices (opt-in)` section, AND `docs/dev/README.md` retains the "real, checkable change" / deliverable-proof policy. This proves the lift *relocated* rather than *deleted* the guidance; the test fails if a future edit strips the dev home.

**AC-3 — The worktree-isolation boundary is substrate-neutral in both shared cores.**
Verified by: the oracle confirming neither `ensign-shared-core.md` nor `first-officer-shared-core.md` contains the literal "CODE only" (case-sensitive, as it appears today), and that both still carry a worktree-isolation clause (presence of the isolation concept, absence of the code-substrate noun — so the fix neutralizes rather than deletes the boundary).

**AC-4 — Both runtime adapters use the shared core's neutral location vocabulary.**
Verified by: the oracle confirming the `## Agent Surface` / `## Dispatch` assignment-field enumeration in `claude-ensign-runtime.md` and `codex-ensign-runtime.md` lists "workflow location" and NOT "worktree path". Scoped to the field-enumeration sentence ("authoritative for all assignment fields: …"), NOT a blanket ban on the literal "worktree path" — the conditional usage at `ensign-shared-core.md:17` ("if you were given a worktree path") is legitimate and must remain passing.

**AC-5 — The contract still drives a real workflow correctly.**
Verified by: a live ensign dispatch on the swept ensign contract completing a stage cleanly (the strongest signal that load-bearing meaning was preserved — same bar Phase 0.A met by driving on both sonnet and opus).

## Lock-in oracle (the negative-proof mechanism)

The proven pattern lives at `internal/hostneutrality/prose_inflator_locks_test.go` and ALREADY enumerates `ensign-shared-core.md` and both runtime adapters in its `contractProseFiles` / `runtimeAdapterPaths` tables. The lock-in oracle for this entity extends that file (same package, same table-driven shape) rather than introducing a new test harness:

- **Negative half (AC-1, AC-3, AC-4) — dev vocabulary absent from the universal core.** A `devLeakageLiterals` table of phrases that must NOT appear in `ensign-shared-core.md` / `first-officer-shared-core.md`: `"failing test"`, `"feature or bugfix"`, `"the test is what the gate judges"` (AC-1, ensign core only), `"CODE only"` (AC-3, both cores). Mirrors `auditTrailLiterals` exactly — a `strings.Contains` loop that `t.Errorf`s on any hit. The AC-4 worktree-field check is a scoped regex over the runtime adapters' field-enumeration sentence (anchor on "authoritative for all assignment fields:"), asserting "workflow location" present and "worktree path" absent **within that sentence** — NOT a file-wide ban (the `ensign-shared-core.md:17` conditional must stay green).

- **Positive half (AC-2) — dev vocabulary present in the dev homes.** A `devHomePresence` table mapping each dev home to a required clause: `development.md` must contain a test-first clause inside `## Recommended practices (opt-in)`; `docs/dev/README.md` must contain "real, checkable change". A `strings.Contains` that `t.Errorf`s on *absence* — the inverse polarity, proving the lift relocated rather than deleted.

This is a presence/absence property of the text where the text IS the claim (legitimate per the README: "a presence check over instruction files proving they carry a required clause or stay free of a banned token is proof at the claim's own level"). The negative-proof discipline (insert banned phrase → red, remove → green) is what makes it a real lock rather than a substring spelling check.

## Spike determination

**No spike needed.** This change composes only already-proven mechanisms:
- The table-driven prose-invariant oracle is proven in-repo (`prose_inflator_locks_test.go`, shipped 0.19.4 / Phase 0.A AC-4).
- The live-ensign-dispatch path (AC-5) is the same dispatch mechanism 0.19.4 already exercised on both sonnet and opus; re-running it on the swept contract is a regression check, not an unverified handoff.
- No new on-disk format, parser round-trip, or runtime handoff is introduced — the edits are prose re-homing across files the oracle already reads.

The one judgement call (which of the four Working-Practices bullets are universal vs. dev-specific) is a design decision recorded above, not an unverified mechanism — it is settled at the gate by review, not by a spike.

## Test plan

- Go oracle tests for AC-1..AC-4 (extending `internal/hostneutrality/prose_inflator_locks_test.go` — the proven Phase 0.A negative-proof + presence pattern). Cost: low — text invariants over instruction files, same package, no new harness.
- One live ensign dispatch for AC-5. Cost: medium (one live cycle). Drive a real stage on the swept contract and confirm clean completion — the strongest signal load-bearing meaning survived the lift.
- High-stakes surface (shipped contract/scaffolding) → a detached adversarial audit is required before merge per the dev README's validation stage. The auditor should specifically probe the AC-4 scoping (try the adversarial edit of removing "if you were given" from `ensign-shared-core.md:17` and confirm the field-scoped oracle stays correctly green — i.e. it does not over-reach into the conditional usage).

## Notes

- The full audit synthesis (6 material + 2 polish findings, MATERIAL-PRESENT) is Task `waymqcmru` from session 10. Findings are quoted above; re-run the audit if deeper detail is needed.
- This is the ensign-side completion of the contract-simplification arc shipped in 0.19.4 (Phase 0.A swept the FO + ensign cores for *length*; this sweeps the ensign core for *dev-leakage* — a different axis the captain's grilling surfaced).
- Sequence: 0.19.5. Independent of the sonnet-live-CI-flake entity, though both touch the ensign-contract surface.

## Stage Report: ideation

- DONE: Design the re-homing: what lifts out of ensign-shared-core.md Working Practices and where each principle lands (development.md / docs/dev/README.md), mirroring the already-corrected FO block (first-officer-shared-core.md 342-353).
  Added `### Re-homing design` to Proposed approach: per-bullet disposition (TDD → new development.md "Recommended practices (opt-in)" subsection; "real, checkable change" already homed at docs/dev/README.md:80 → delete-only; "prove by exercising" + "no hidden machine deps" are universal → keep, de-dev the heading; optional one neutral universal line mirroring the FO core). Pinned exact before/after wording for AC-3 ("CODE only", both cores) and AC-4 ("worktree path" → "workflow location", both adapters) against verified current line numbers.
- DONE: Specify the lock-in oracle (dev-only vocabulary absent from the universal core, present in the dev homes) — the proven Phase 0.A negative-proof pattern.
  Added `## Lock-in oracle` + `## Spike determination` sections: the oracle EXTENDS the in-repo `internal/hostneutrality/prose_inflator_locks_test.go` (which already lists these files) with a negative-half banned-literal table (mirrors `auditTrailLiterals`) and a positive-half dev-home-presence table. Recorded "no spike needed" with the proven mechanisms (table-driven prose oracle + already-exercised live-dispatch path).

### Summary

Refined the existing spec rather than rewriting it: ground-truthed every cited line number, then sharpened the two checklist items into an implementable design. Key ideation finding: only the TDD authoring bullet genuinely needs a NEW home (development.md) — the deliverable-shape principle is already at docs/dev/README.md:80 (delete-only), and two of the four Working-Practices bullets ("prove by exercising", "no hidden machine deps") are genuinely universal and must STAY. Flagged an AC-4 over-reach trap: the conditional "if you were given a worktree path" at ensign-shared-core.md:17 is legitimate, so the oracle must scope to the runtime adapters' field-enumeration sentence, not ban the literal file-wide. Recorded "no spike needed" — composes only the proven prose-oracle and live-dispatch mechanisms.
