---
title: One-command gate review presentation with atomic result retention
status: ideation
source: "Split from the gate-recorder task (3k), captain-approved 2026-07-21. The subspace-coupled presentation half; 3k cycles 11-12 are its banked design history."
id: xbatj4hxtxw9t83vvmfem27f
gates:
    version: 1
    current:
        gate: gate:docs-dev:xb:ideation
        attempt: gate-attempt:xb-ideation-1
    records:
        - id: gate:docs-dev:xb:ideation
          stage: ideation
          current-attempt: gate-attempt:xb-ideation-1
          attempts:
            - id: gate-attempt:xb-ideation-1
              sequence: 1
              state: open
              current-briefing:
                id: briefing:docs-dev:xb:ideation:briefing-1
                digest: sha256:bfb87cff3c021c17af9a9d8a999cb682cde972d73d292a5838f21baec7a240de
                room-ref: "./review/ideation/briefing-1"
                note: "Multi-artifact package: gate summary, frozen entity snapshot, frozen recorder-contract snapshot — each digest-pinned inside the briefing; the digest above binds briefing.json itself."
sprint: durable-decisions
group: recorder
started: 2026-07-21T01:43:36Z
---

One binary command presents a gate: it validates an explicit briefing package (gate summary, frozen design snapshot, frozen probe input/history), derives the canonical title, launches the Subspace TUI as a blocking child through the caller's terminal transport, and atomically validates and retains the review log, resolution, and diagnostics on success or failure — the ensign presenter stays addressable until TUI exit plus retention, and pane creation or timeout is never completion. Includes the provider id-mapping adapter (the provider resolution binds its own envelope briefing id; normalize to the attempt briefing id after digest validation) as SPECIFIED in 3k's gate-resolution-frontmatter-contract.md. Scope moved from 3k at the split: ACs 7 and 15, the presentation-side AC-8 mutants (early-completion, detached-worker, live-Reference append, controller/child/validation/retention), and the probes companion (gate-review-probes.md rides as provider-owned convention; probes.jsonl proved out in one dry run). Evidence base: the 0260 shaping float findings 1-15 (blank-float EOF defect, launcher repair, probe-first ritual, retention-deleted-on-failure incidents, and finding 15 — the captain's own approval destroyed when a dead launcher unlinked its scratch result; finding 14 is a deliberate live-session numbering skip) in the shaping debrief and 3k's attempt-7 resolution provenance note. Cross-repo sequencing: depends on the subspace-tui briefing-package and result surfaces; the working-copy-skill ritual recorded in the debrief is the interim. Land after the recorder (3k).

## Expected surface + tolerance

Go product code, spacedock-side only: ~2-3 new files under `internal/` (briefing-package validation + title derivation + the launch/retain/validate/id-map lifecycle) plus 1 new `spacedock gate review` verb entry in `cmd/` and a small provider-probe helper. ~350-550 production LOC, roughly equal test LOC (a fixture table driving the real command against a fake subspace-tui child per exit path — resolved approve/revise, leave-open/hold, blank-float EOF, mid-run crash, post-exit validation failure, missing binary). No contract doc change (the spec is the already-banked `gate-resolution-frontmatter-contract.md`, owned by 3k); ~1 doc-reference addition for the new verb. **Tolerance 2×.** Hard self-check: any entity-frontmatter write (that is the recorder's, 3k), any blocker/eligibility/hold/application computation (that is h1's), any subspace-tui product change (xb consumes its surface, does not build it), or any schema change to the gates records trips a reconfirm.

## Problem

A captain's gate decision presented through Subspace can be lost. The 0260 shaping session ran five successful Subspace floats but collected the presentation-side failure evidence in the same two days (shaping debrief float findings 1-15; finding 14 is a deliberate live-session numbering skip):

- **A dead launcher destroyed the captain's own approval (the primary red fixture).** Finding 15, recorded in 3k's attempt-7 resolution provenance note: the captain resolved a gate in a float pane whose launcher had already died; the resolution was written to a launcher-owned scratch file that was then unlinked, destroying the decision. Only a chat re-affirmation recovered it. A retention target owned by a dying launcher is exactly the failure this command removes — the result lands at a caller-owned `--result` path outside any launcher-owned scratch, so the launcher's death cannot unlink it. The interim launcher's later room-resident-result fix is the same lesson reached by hand.
- **The result is destroyed on the failure and hold paths.** The interim `review-local-zellij` launcher captures the child's stdout into a `mktemp -d` scratch and carries `trap 'rm -rf "$scratch"' EXIT` (findings 3, 12). When result validation fails, or the reviewer leaves the review open/hold, the scratch — and the only copy of the result — is deleted. The leave-open/hold result shape is exactly the branch that was never retained.
- **The blank-float EOF leaves no recoverable state.** The `subspace-tui present` child protocol opened an empty float and ended with `present-child protocol ended early: EOF`, exit 1, no recovery path emitted (finding 12). Reproduced here against the exact-tip binary: on the can't-render path the provider writes **nothing** to its `--result` path and removes the zero-entry package, so the only trace is one line of stderr — which the interim launcher then throws away.
- **Single-file concatenation.** The interim launcher accepts exactly one Markdown file, forcing the gate summary and the frozen entity snapshot to be concatenated; the multi-artifact briefing package is designed (3k) but the interim path cannot present it as separate artifacts (finding 5, captain-flagged).

The decision therefore does not reliably become durable, recoverable state, and the ritual that works today depends on a skill-local shell script with a destroy-on-exit defect. This is the presentation half split from the gate recorder (3k); the recorder owns the `gates:` writes, this command owns getting a validated, retained result into the First Officer's hands without losing it.

## Required capability

One binary command — `spacedock gate review <slug> --workflow-dir DIR --stage STAGE --briefing FILE` — presents a gate and cannot lose the result:

1. **Validate the explicit briefing package before launch.** Parse the briefing JSON, verify each Artifact and supporting Reference revision/digest, and derive the canonical pane title in the provider's title grammar (`<short-id> · <stage> · (approve|revise|hold)?`, ≤40 columns). Probe for `subspace-tui` and the review skill; if either is missing, print the exact install commands and do not launch.
2. **Announce that workflow state will not advance**, then launch `subspace-tui` as **one blocking child process** through the caller's terminal transport, binding a caller-owned `--result` path.
3. **Own a retention directory the command never auto-deletes.** It holds the result the child writes to `--result`, the review log when one exists, and a diagnostics record (the exact launch argv, the child exit code, and captured stderr). Retention is atomic and happens on **success AND failure** — a launch, controller, child-exit, validation, or retention-write failure preserves the package, diagnostics, and any partial result for retry.
4. **Completion is exit + validated result + retention, never pane creation or a wait timeout.** The command returns success only after the child exits, the returned result validates against the bound briefing, and retention completes. The gate-attempt ensign that invokes the command stays unresolved and addressable for the whole blocking call; the First Officer waits with `wait_agent({timeout_ms:300000})` and re-waits on timeout rather than treating a timeout as completion.
5. **Normalize the provider envelope briefing id to the attempt briefing id after digest validation**, as specified in `gate-resolution-frontmatter-contract.md` — see below.
6. **Never mutate entity frontmatter.** Recording the resolution into `gates:` is the recorder's job (3k); computing blockers, holds, eligibility, and applying the action is h1's. This command hands the First Officer a validated, retained result and stops there.

## What ships spacedock-side now vs. the subspace-tui dependency

**Ships spacedock-side now (this task, fully testable with a fake subspace-tui child):** briefing-package validation, canonical title derivation, provider probe with exact install text, retention-directory allocation, the blocking child launch, stderr→diagnostics capture, launch-metadata recording, post-exit result validation, provider id normalization, and atomic retention on every exit path. None of this needs the real TUI to be exercised — the launch is a subprocess with a `--result` path and an exit code, so the whole lifecycle is driven in fixtures against a stub child (proven by the spike below).

**Needs subspace-tui surfaces (declared cross-repo dependency, NOT built here):**
- A **blocking review invocation that accepts a briefing package and writes to a caller-owned `--result` path.** Today the provider's `--result` durable write exists only in the single-file one-file mode; `--review-v1 <briefing.json>` (the multi-artifact briefing package) parses exactly six args and ignores `--result`. So briefing-package presentation and durable `--result` retention do not yet coexist in subspace-tui. This command binds to a narrow launch contract (argv + `--result <path>` + exit code + stderr); the concrete binding that works today is the single-file advisory invocation, and the multi-artifact briefing-package + `--result` coexistence is the subspace-tui gap that unlocks the full journey.
- A **non-EOF blocking transport.** `subspace-tui present --transport zellij` (the present-child protocol) carries the blank-float EOF defect; the proven working path is the direct blocking launch of the TUI in a Zellij float. Repairing/owning the transport is subspace-tui's, per scope ("xb consumes its surface, does not build it").

**The interim baseline this command must beat:** the `subspace-r-working-copy` skill's `review-local-zellij` (shaping debrief finding 13) — proven across five floats, but single-file only and destroying the result via `trap rm -rf` on the failure and hold paths. This command beats it on the measurable axis (below): it accepts a multi-artifact briefing package and retains the result/log/diagnostics on the exit paths where the interim path retains nothing.

## Provider id-mapping (as specified in the contract)

The provider mints its own briefing id: subspace-tui's one-file mode derives `briefing:single-file:<hex16>` from the result path (`BriefingIDForInvocation`), and stamps it into `result.briefing` and `result.resolution.briefing`. The recorder (3k) verifies the current attempt Briefing id/digest, so a provider-minted id would never bind. Per `gate-resolution-frontmatter-contract.md` (the frozen Briefing binds an immutable snapshot; the provider result is keyed by Briefing id and joined after digest validation), this command:

1. validates the artifact digest in the returned result against the digest bound by the attempt briefing;
2. **only on digest match**, rewrites the provider envelope briefing id in the retained result to the attempt briefing id, so the recorder receives a result already keyed to the attempt it verifies;
3. on digest mismatch, rejects and retains diagnostics without normalizing (an unverified result is never laundered into the attempt id).

Normalization touches the retained result file (presentation state this command owns), never entity frontmatter.

## Spike: the riskiest unverified mechanism (atomic retention on success AND failure)

Exercised end-to-end before design lock; harness at `scratchpad/retention-spike.sh`, plus a probe of the exact-tip binary.

- **Real-binary grounding.** Built the tip `subspace-tui` (`go build ./cmd/subspace-tui`, version `dev`) and ran the advisory one-file mode with `--result` on the can't-render path (no TTY, the blank-float EOF analogue). Result: **exit 1, nothing written to `--result`, zero-entry package removed, one line of stderr.** Confirms the command cannot rely on the provider to retain anything on the failure path — it must own retention.
- **Retention-contract spike (9/9 assertions pass).** A fake TUI child reproduces the three red fixtures — (A) blank-float EOF: exit non-zero, empty result; (B) leave-open/hold: writes a `"status":"open"` result, exit 0; (C) launcher/controller death after the result is written (finding 15): a real launcher process writes the approval then is `kill -TERM`ed before returning — under two wrappers: the interim `trap rm -rf` scratch pattern (launcher owns the result path), and the command's contract (caller-owned dir never deleted, `--result` binding, stderr→diagnostics, argv+exit recorded). The interim baseline retains **nothing** on all three (reproducing finding 3, finding 12, and finding 15 — the dying launcher's trap unlinks the captain's approval); the command's contract retains the result-or-diagnostics, the log when present, and the launch record on **all three**, and the approval written to the caller-owned `--result` path survives the launcher's death because the launcher owns no cleanup over it.

Conclusion: the retention mechanism is proven at the contract level with real files, real exit codes, and a real killed launcher process, reproducing the exact debrief failure signatures — including finding 15, 3k's named primary red fixture. The spike table is the first implementation test (AC-1). The interactive TUI transport is out of scope (subspace-tui's) and is not what this command's value rests on.

## Acceptance criteria

**AC-1 (VALUE) — No presented decision is lost on any exit path.** Driving `gate review` through each retention fixture — resolved approve, resolved revise, leave-open/hold, blank-float EOF (child exits non-zero, empty result), mid-run child crash, and post-exit result-validation failure — leaves a caller-owned retention directory that still holds the result (or, when the child produced none, a diagnostics record naming the launch argv, exit code, and captured stderr) plus the review log when one exists. Retained count is N/N across the fixtures; the same fixtures driven through the interim `review-local-zellij` baseline retain 0/N on the failure and hold paths. *Test:* Go fixture table driving the real command with a fake TUI child per fixture (the spike harness promoted to a committed test), asserting on-disk retention-directory contents after each run and contrasting the baseline's `trap rm -rf` destruction. This measures the end-value (result-retention rate) against an independent baseline that can regress to 0.

**AC-2 — Retention survives every failure class, including launcher/controller death.** Launch, controller, child-exit, validation, and retention-write failures each leave the briefing package, diagnostics, and any partial result recoverable; none deletes the retention directory, and each returns non-zero. In particular, a launcher/controller that dies after the result is written (finding 15 — the red fixture where a dying launcher unlinked the captain's approval) leaves the result intact, because it lands at the caller-owned `--result` path this command allocates, not in a launcher-owned scratch the launcher's cleanup can unlink. *Test:* table test injecting each failure class, including a launcher process that writes the result then is killed before returning; asserts directory survival, the result surviving the launcher's death, diagnostics presence, and exit code.

**AC-3 — Pane/session creation and wait-timeout are never completion.** The command returns success only after the child exits, the result validates, and retention completes; a launched-but-unexited child yields neither a success exit nor a retained validated result. A mutant that returns on pane creation, and a mutant that lets the presenter resolve before child exit, both fail. *Test:* fixture with a child that emits a pane marker then blocks; assert no success and no validated result until the child exits.

**AC-4 — The provider envelope briefing id is normalized to the attempt briefing id, only after digest validation.** Given a result whose `briefing`/`resolution.briefing` carry the provider-minted `briefing:single-file:<hex>` id: on a digest that matches the attempt briefing, the retained result's briefing ids equal the attempt briefing id; on a digest mismatch, the command rejects and retains diagnostics with the ids un-rewritten. *Test:* fixture feeding a provider result with the minted id plus a matching and a mismatching artifact digest; assert normalized ids on match and rejection (no normalization) on mismatch.

**AC-5 — The briefing package is validated and the title derived before launch; missing prerequisites emit the exact install action.** A valid multi-artifact briefing yields a canonical pane title in the provider's grammar and confirms each Artifact and supporting Reference revision is present; an absent `subspace-tui` or review skill prints the exact install commands and does not launch. *Test:* a valid-briefing fixture (assert title derivation + Reference presence) and a missing-binary fixture (assert exact install text, zero launch).

**AC-6 — The command never mutates entity frontmatter.** Across every success and failure fixture, the entity file's frontmatter bytes (including `gates:` and `status`) are unchanged; only the caller-owned retention directory and the provider result change. *Test:* byte-compare the entity file before and after each fixture run.

## New mechanisms (value AC served / simplest alternative / why insufficient)

- **Caller-owned retention directory the command never auto-deletes** — serves AC-1, AC-2. *Alternative:* rely on subspace-tui's own `--result` write and package retention. *Insufficient:* the real-binary probe shows the provider writes nothing and removes the zero-entry package on the blank-float EOF path, and the interim skill compounds it with `trap rm -rf`. The command must own a directory it never deletes.
- **Diagnostics capture (argv + exit code + stderr) when the child yields no result** — serves AC-1, AC-2. *Alternative:* capture only the result/stdout. *Insufficient:* the EOF path produces no result and only one stderr line; without argv/exit/stderr there is nothing to retry from ("no recovery path emitted", finding 12).
- **Provider id normalization gated on digest validation** — serves AC-4. *Alternative:* let the recorder consume the provider-minted id, or normalize before validating the digest. *Insufficient:* the recorder binds the attempt Briefing id/digest, so a `briefing:single-file:<hex>` id never binds; normalizing before validation would launder an unverified result into the attempt id.
- **Blocking single-invocation launch contract (argv + `--result` + exit code)** — serves AC-1, AC-3. *Alternative:* the present-child streaming protocol (result over stdout). *Insufficient:* present-child carries the blank-float EOF defect; a `--result` write is atomic (temp + rename) and decouples retention from a fragile stdout stream.

## Behavioral test plan

Behavior-first Go fixtures driving the real command against a fake subspace-tui child; no live TUI, no host smoke, no committed prose-greps. Estimated cost: medium — a child-stub table plus on-disk assertions and one byte-identity check.

1. **Retention table (AC-1, AC-2).** One fixture per exit path (approve, revise, hold/leave-open, EOF, crash, validation-failure, retention-write-failure, and launcher/controller death after the result is written — finding 15); assert retention-directory contents, the launcher-death approval surviving, and the N/N-vs-0/N baseline contrast. Promotes `scratchpad/retention-spike.sh` (9/9) to a committed Go test.
2. **Completion boundary (AC-3).** Child emits a pane marker then blocks; assert no success/validated result until exit. Mutants (return-on-pane, resolve-before-exit) fail.
3. **Id normalization (AC-4).** Provider-minted id with matching/mismatching digest; assert normalize-on-match, reject-on-mismatch.
4. **Pre-launch validation (AC-5).** Valid multi-artifact briefing → canonical title + Reference presence; missing binary → exact install text, zero launch.
5. **Frontmatter immutability (AC-6).** Byte-compare the entity file across all fixtures.

Riskiest mechanism first: item 1 (atomic retention on success AND failure) — already exercised end-to-end in the spike.

## Documentation change proposal

Add to the gate-review command reference in `docs/site` (the presentation-side entry; the recorder's frontmatter/gates doc changes are 3k's):

```diff
+### Review a complete gate Briefing
+
+When the First Officer offers Subspace and you answer yes, the gate-attempt ensign runs
+`spacedock gate review <slug> --workflow-dir DIR --stage STAGE --briefing FILE`. The command
+validates the explicit briefing package, derives the pane title, and launches Subspace as one
+blocking child. It atomically retains the review log, the result, and diagnostics on success
+**and** on failure — a launch, child-exit, validation, or retention failure preserves the
+package and diagnostics for retry — and never changes workflow state. The First Officer records
+and applies any binding decision separately.
```

## Out of scope

- Writing the resolution into `gates:` frontmatter, or any entity-frontmatter mutation (the recorder, 3k).
- Computing blockers, execution holds, eligibility, or applying the gate action (h1, `gate-blockers-and-eligibility`).
- Building or repairing the subspace-tui transport, the `--review-v1` + `--result` coexistence, or any subspace-tui product change (xb consumes the surface).
- Rendering ProbeResult/comparison UI (a recorded Subspace product gap; this command carries the frozen semantic-delta summary as a supporting Reference, it does not render it).
- The probes companion (`gate-review-probes.md`) as a provider-owned convention — it rides along, not built here.

## Stage Report: ideation

- DONE: A concrete command design against the banked 3k cycles 11-12 lifecycle: blocking presenter, atomic result/log/diagnostics retention on success AND failure (the destroyed hold-path result and the blank-float EOF are the red fixtures), and the provider id-mapping implemented as SPECIFIED in 3k's gate-resolution-frontmatter-contract.md.
  `## Required capability` designs the six-step lifecycle (validate → announce → blocking child → owned retention dir → completion-is-exit+validated+retained → no frontmatter write) matching cycle 11-12's addressable-blocking-presenter and non-terminal-pane/timeout discipline; `## Provider id-mapping` implements the contract's normalize-after-digest-validation rule; AC-1/AC-2 carry both red fixtures as fixtures.
- DONE: The cross-repo dependency declared honestly: what needs subspace-tui surfaces vs what ships spacedock-side now; the working-copy-skill ritual as the measured interim baseline it must beat.
  `## What ships spacedock-side now vs. the subspace-tui dependency` splits the two: the full validate/launch/retain/id-map lifecycle ships spacedock-side (testable against a stub child); the briefing-package + `--result` coexistence and the non-EOF transport are the declared subspace-tui gaps. The `review-local-zellij` interim is named as the baseline, and AC-1 measures the command beating it (N/N vs 0/N retention).
- DONE: Expected surface + tolerance declared; riskiest unverified mechanism spiked end-to-end first or an auditable no-spike-needed recorded.
  `## Expected surface + tolerance` declares ~350-550 prod LOC ≈ equal test LOC, ~2-3 `internal/` files + 1 `cmd/` verb, tolerance 2×, with a hard self-check fencing off 3k/h1/subspace scope. `## Spike` records the end-to-end retention spike (7/7, `scratchpad/retention-spike.sh`) plus the exact-tip binary probe proving the provider retains nothing on the blank-float EOF path.

### Summary

Designed `spacedock gate review` as the presentation half split from 3k: a blocking presenter that owns a never-deleted retention directory and atomically retains result/log/diagnostics on every exit path, reproducing and beating the interim ritual's two destruction defects. The riskiest mechanism (atomic retention on success AND failure) was spiked end-to-end against the real tip binary before design lock — the provider writes nothing to `--result` on the blank-float EOF path, so the command owns retention. Scope holds strictly to the resolution/briefing side: no frontmatter writes (3k), no application/eligibility (h1), no subspace-tui product work.

## Stage Report: ideation (cycle 2)

Preflight fold applied (the first decline: the destroyed-approval fixture).

- DONE: Extend the evidence base to findings 1-15.
  Bumped both evidence lines (overview + Problem) to findings 1-15, cited 3k's attempt-7 resolution provenance note (verified at the source: index.md line 164), and marked finding 14 as the deliberate live-session numbering skip.
- DONE: Add finding 15 as the controller/launcher-death red fixture, with the AC-2 clause.
  Added a Problem bullet naming finding 15 (a dying launcher unlinked the captain's own approval) as the primary red fixture; extended AC-2 with the launcher-death clause (the caller-owned `--result` path survives because the launcher owns no cleanup over it) and its test. Exercised it: spike fixture C kills a real launcher process after it writes the result — interim scratch is unlinked (finding 15 reproduced), caller-owned result survives. Spike now 9/9.

### Summary

The fold strengthens the evidence base with finding 15 — the strongest red fixture for this exact mechanism, since it is the caller-owned retention directory's whole reason to exist. Verified the provenance at 3k's attempt-7 note rather than taking it on faith, then proved the fixture by exercising a real killed launcher process (spike 9/9), not by asserting it.
