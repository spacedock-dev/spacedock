---
title: Trim the redundant status --read section-read guidance — grep already covers it
status: implementation
source: "FO+captain analysis (2026-06-16), measured on real files: for entity bodies (the actual read target) `status --read --json` and `grep -nE '^#{1,4} '` produce IDENTICAL heading maps (m4 entity 19=19; FO shared-core 18=18). grep over-counts ONLY on fenced markdown-like content (dev README 23 vs 18 — the task-template block). The contract's FO line `first-officer-shared-core.md:214` pairs 'prefer grep, anchor on the heading' WITH 'status --read for offset/lines' — but grep's heading list already yields a section's offset AND its span (the next heading's line); `wc -l` yields the append-point total_lines (ensign `ensign-shared-core.md:92`); `status --resolve`/`--where --json` yield frontmatter. The ensign sites (`:18`, `:92`) name ONLY --read, dropping the grep alternative entirely. So the --read adoption guidance largely re-states the grep-anchor rule the contract already mandates; the sole non-redundant residue is fence-safe heading detection (situational). hf's four FO captures read 0/0 — consistent with re-selling a tool grep already covers, and explains why the trimmed site-6 (4x) was the wrong lever (instruct harder)."
sprint: 0240-lean-contract
sprint-readiness: ready
issue:
id: 82kzghcy3j3cet3hynwa4165
started: 2026-06-30T09:09:18Z
worktree: .worktrees/spacedock-ensign-read-guidance-redundant-with-grep
---

## Problem
The `status --read` section-read guidance spans FO `first-officer-shared-core.md:214` and ensign `ensign-shared-core.md:18,:92`. Measured against grep on real files, for the primary read target (entity bodies) the two produce identical heading maps; grep over-counts only on fenced markdown-like content (the README's task-template block). Everything the --read guidance instructs is already available: grep's heading list gives a section's offset + span; `wc -l` gives the append-point total_lines; `status --resolve`/`--where --json` give frontmatter. The FO bullet (214) at least names grep as primary and tacks --read on; the ensign sites (18, 92) name ONLY --read with no grep alternative — the redundancy is concentrated there. Net: the adoption push re-sells what "prefer grep" already covers, which is why adoption read 0/0 in hf's FO captures and why instructing harder (the trimmed site-6) missed.

## What's needed (evidence-first; composes with hf + f5)
- Reduce the --read adoption instruction to its non-redundant residue: fence-safe heading detection where grep over-matches, plus one-call frontmatter+spans as a *convenience* — not a mandate over grep.
- Make FO and ensign consistent: name grep as the primary section-locator in the ensign sites (18, 92) as FO 214 does, or collapse to one shared rule.
- Keep the `status --read` TOOL (correctness-by-construction is a thin but real rationale); trim the INSTRUCTION, not the binary.
- Gate the trim on measured adoption (hf's metric, ensign-transcript-aware per f5): prove --read+scoped-Read usage does not regress after the trim, so this is evidence-driven rather than another assertion.

## Scope: the redundant sites vs. the load-bearing ones
The `status --read --json` TOOL has two genuinely irreplaceable uses that grep CANNOT cover — these STAY untouched:
- **Structured frontmatter / stages extraction** (`first-officer-shared-core.md:17`, the README taxonomy read): `--json` yields the parsed `stages` array and flat `frontmatter` object. grep cannot parse YAML into a stages array.
- **Structured roll-up modes** (`first-officer-shared-core.md:121`): `status --read <ref> --checklist` and `status --read <ref> --ac-scan`. grep cannot produce a checklist/AC roll-up.

The REDUNDANT instruction — the only thing this task trims — is narrowly *"use `--read --json` to get a HEADING's `offset`/`lines` for a section-scoped `Read`"*, because grep's heading list already yields that. Four sites carry it; the checklist names three, and inspection found a fourth (FO `:100`, the gate-verdict read of the last `## Stage Report`) carrying the identical redundancy with NO grep alternative — see **Open decisions for the gate**.

## Acceptance criteria
- **AC-1 (end value — measured against a baseline that can move the wrong way)** — After the trim, a journeymetrics before/after over REAL FO+ensign transcripts shows `--read`+scoped-`Read` adoption does NOT regress: post-trim per-journey `status_read_calls`+`scoped_read_calls` ≥ the pre-trim baseline. The baseline can move the wrong way — if the trim removed load-bearing guidance, the ensign would fall back to whole-file Reads and the counts would DROP. Measured, not asserted. **Hard dependency:** the metric must be ensign-transcript-aware (f5, `f53zr0ehhzekzbgydpybaq5g`) — `ParseClaudeJSONL` today parses only the FO front-door stream, so without f5 the comparison is the vacuous `0==0` hf already hit (see **Riskiest-mechanism spike**). f5 is backlog in sprint 0204 and must land (or be pulled into 0240) before AC-1 is satisfiable.
- **AC-2 (means — counts only paired with AC-1)** — The redundant section-read instruction is reduced to its residue (fence-safe heading detection where grep over-counts; structured-mode convenience as a non-mandate) and grep is named the primary section-locator consistently across all trimmed sites (FO Probe/Read bullet, ensign `:18`, ensign `:92`, and — pending the gate — FO `:100`). Verified BEHAVIORALLY: (a) `go test ./internal/dispatch/ -run TestBuild` goldens stay BYTE-IDENTICAL — a negative control proving the trim is runtime-loaded skill prose, not dispatch-prompt text (the actual dispatch hint, "site-6", was already removed in #392 `9458a636`); (b) `go test ./internal/contractlint/ ./internal/ensigncycle/` stays green. NOT verified by a prose-grep over the core files — the `instruction_read_detector`/boundary guard forbids a test that reads an instruction file to assert a phrase. Per the ideation rule, this means-AC counts only paired with AC-1's value measurement.

## Before/after wording
Three mandated sites + one recommended (D). Each AFTER names grep primary and keeps `--read --json` only as the fence-safe residue (A, C) or drops it where `wc -l`/grep fully covers (B, D). Line numbers are as-found 2026-06-30; implementation anchors on bullet TEXT (the sibling shifts ensign-core line numbers — see Collision map).

**Site A — `skills/ensign/references/ensign-shared-core.md:18` (`## Working`, list item 1)**

BEFORE:
```text
1. Read the entity file before making changes — when you need only specific sections (the relevant stage-def section, your prior report), `status --read <entity-path> --json` returns each section's `offset`/`lines` for a scoped `Read`, rather than the whole body.
```
AFTER:
```text
1. Read the entity file before making changes — for a specific section (the relevant stage-def section, your prior report), locate its heading with `grep -nE '^#{1,4} '` and scoped-`Read(offset, limit)` the span to the next heading. `status --read <entity-path> --json` is the fence-safe fallback when the body carries fenced markdown-like content that bare grep over-counts.
```

**Site B — `skills/ensign/references/ensign-shared-core.md:92` (`## Stage Report Protocol`, append bullet)**

BEFORE:
```text
- Append the report at the end of the entity file — get the file's `total_lines` from `status --read <entity-path> --json` and append after it; do not read the entire entity body to find an insertion point.
```
AFTER:
```text
- Append the report at the end of the entity file — get the append point from `wc -l <entity-path>` and append after the last line; do not read the entire entity body to find an insertion point.
```

**Site C — `skills/first-officer/references/first-officer-shared-core.md:229` (`## Probe and Ideation Discipline`, "Prefer Grep over Read" bullet)**

BEFORE:
```text
- Prefer Grep over Read for targeted entity-body inspection. Anchor on heading or field name (`## Stage Report`, `### Feedback Cycles`, a specific frontmatter field). Read only when you need the full text. To pull a whole named section, `status --read <ref> --json` returns its `offset`/`lines` so the follow-up `Read(offset, limit)` is section-scoped, not whole-file.
```
AFTER:
```text
- Prefer Grep over Read for targeted entity-body inspection; Read whole only when you need the full text. Anchor on a heading or field name (`## Stage Report`, `### Feedback Cycles`, a frontmatter field): `grep -n` gives the heading line (the section offset) and the next heading bounds its span, so the follow-up `Read(offset, limit)` is section-scoped. `status --read <ref> --json` is the fence-safe fallback when markdown-like fenced content makes grep over-count headings.
```

**Site D — `skills/first-officer/references/first-officer-shared-core.md:100` (`## Completion and Gates`, step 1) — RECOMMENDED, pending gate**

BEFORE:
```text
1. Read the entity file's last `## Stage Report` section — `status --read <ref> --json`, take the last `## Stage Report` heading's `offset`/`lines`, then `Read(offset, limit)` that range, instead of loading the whole body.
```
AFTER:
```text
1. Read the entity file's last `## Stage Report` section — `grep -n '## Stage Report' <ref>` gives every heading line; scoped-`Read(offset, limit)` from the last one to EOF, instead of loading the whole body.
```

Net line/token impact: A and C are length-neutral; B and D shorten. The trim is net slightly-negative, but this item's GATE is adoption (AC-1), not a byte delta (unlike the boot-core deferral siblings).

## Test plan
- **AC-1 (adoption-not-regress) — live workflow, ensign-transcript-aware.** Cost: high (live FO+ensign runs) and BLOCKED on f5. Once f5 folds `subagents/agent-*.jsonl` into `ParseClaudeJSONL`'s counts: capture a pre-trim baseline (run the journey on `origin/main`'s cores; record per-journey `status_read_calls`+`scoped_read_calls`), apply the trim, re-run the same journey, assert post-trim ≥ baseline. The detector primitives (`commandInvokesStatusRead`, `readInputIsScoped`) are already unit-proven (see spike) — what's new and untested is the ensign-transcript fold, which is f5's own AC-1 test. This task consumes f5's mechanism; it does not re-implement it.
- **AC-2 (residue trim + grep-primary consistency) — Go fixture/CLI, fast (<5s).** `go test ./internal/dispatch/ -run TestBuild` must stay GREEN with goldens byte-identical (negative control: trimming runtime-loaded skill prose must NOT move the dispatch prompt). `go test ./internal/contractlint/ ./internal/ensigncycle/` must stay green (the cores feed both). No new prose-grep test is added — the boundary guard forbids reading a contract file to assert a phrase; the behavioral surfaces above ARE the proof.
- **Estimated complexity:** the EDIT is trivial (4 bullets); the RISK and cost are entirely in AC-1's dependency on f5. If f5 cannot land in 0240, AC-1 degrades to a forward-looking gate (ship the trim; the ensign-aware metric confirms no regression in a later live run) — a captain decision (see Open decisions).

## Riskiest-mechanism spike
**Question:** can the journeymetrics `--read` adoption metric actually produce a real before/after from transcripts?

**Spiked end-to-end by source-trace + running the existing suite (2026-06-30):**
1. **Detector primitives PROVEN.** `go test ./internal/journeymetrics/ -run 'StatusRead|ReadAdoption|CommandInvokesStatusRead|CodexCharacterization|DedupAcross'` → all green (incl. multi-delta dedup + codex). `ParseClaudeJSONL` counts `status --read` Bash calls (`StatusReadCalls`) and scoped `Read` calls carrying offset/limit (`ScopedReadCalls`) from a stream, launcher-agnostic and flag-order-independent. So "detect `--read`+scoped-`Read` in a transcript" = YES, proven.
2. **But the ENSIGN surface is NOT yet measured.** `ParseClaudeJSONL` parses ONE stream — the FO front-door `claude -p` stream (`internal/ensigncycle/journey_metrics_live_test.go:19`, `result.stream`). It never opens `subagents/agent-*.jsonl`. The trimmed sites steer principally the ENSIGN's reads; the ensign runs as a separate team-agent session. Empirically (hf) four real FO captures all read `0/0`. So a real before/after of the ENSIGN's adoption is `0==0` until the sub-agent transcript is folded in.
3. **The fold is f5** (`journeymetrics-ensign-read-adoption.md`, `f53zr0ehhzekzbgydpybaq5g`) — backlog, sprint 0204, NOT landed.

**Conclusion:** the metric CAN produce read-adoption counts from a transcript (mechanism proven), but a *meaningful* ENSIGN before/after is GATED on f5. This is recorded as the hard dependency on AC-1, not "no spike needed." **Recommendation:** sequence f5 before this task's implementation (pull it into 0240 as a predecessor), or accept the AC-1 degradation in Open decisions.

## Collision map with sibling `ensign-contract-dev-leakage` (scr, `scr2rx4589p7j6mpgh50hdct`)
Both edit `skills/ensign/references/ensign-shared-core.md`. The sprint sequences Wave 1 STRICTLY SEQUENTIAL (never parallel), so the second implementer rebases on the first's committed core. Anchoring (2026-06-30 line numbers):
- **MY ensign-core sites:** `## Working` list item 1 (`:18`); `## Stage Report Protocol` append bullet (`:92`). My FO-core sites (`:229`, and `:100` if approved) are in `first-officer-shared-core.md`, which the sibling does NOT touch.
- **SIBLING's ensign-core targets (dev-only discipline):** the worktree/deliverable prose — currently surfacing at `## Working` item 2 (`:19`, "keep all reads, writes, and commits under that worktree"), `## Worktree Ownership` (`:30`–`:44`, esp. `:36` "the worktree isolates the deliverable work product only"), and `## Split-Root State Contract` (`:34`–`:44`). The sibling is un-ideated backlog, so its exact removals aren't fixed.
- **Disjoint by section** — EXCEPT one adjacency: my `:18` (Working item 1) sits directly above the sibling's likely `:19` (Working item 2). Same ordered list. No paragraph-level overlap, but a naive line-based merge could touch adjacent lines. Mitigation: sequential Wave-1 execution + anchor edits on bullet TEXT, not line number; the second worker re-reads `## Working` fresh. The FO core is mine alone — zero overlap there.

## Open decisions for the gate
1. **FO `:100` in scope?** The checklist enumerates only the FO Probe/Read bullet (`:229`) + ensign `:18`/`:92`. Inspection found FO `:100` (gate-verdict read of the last `## Stage Report`) carrying the IDENTICAL redundancy with NO grep alternative — leaving it would keep the FO core internally inconsistent, defeating the task's "grep primary consistently across both" thesis. RECOMMEND including Site D. Flagged rather than silently expanding scope.
2. **f5 dependency.** AC-1 is unsatisfiable until f5 makes the metric ensign-transcript-aware. RECOMMEND pulling f5 into 0240 as this task's predecessor. Fallback (captain call): degrade AC-1 to a forward-looking gate.
3. **Stale dispatch-golden clause.** The original AC-1's "dispatch goldens regenerated to the trimmed prompt" was correct when the dispatch prompt carried a `--read` hint (site-6), but #392 already removed it. Reframed in AC-2 as a byte-identity negative control. No golden bytes change from this trim.

## Stage Report: ideation

- DONE: Concrete before/after wording for the redundant `status --read` section-read guidance at the FO core (the Probe/Read bullet) and the ensign core (`:18` and `:92`), naming grep as primary consistently, keeping only the fence-safe-heading residue
  `## Before/after wording` Sites A (ensign `:18`), B (ensign `:92`), C (FO `:229`) — each AFTER names grep primary, keeps `--read --json` only as the fence-safe fallback (A, C) or drops it for `wc -l` (B); FO `:100` added as recommended Site D for consistency.
- DONE: An end-value AC measuring `--read`+scoped-`Read` adoption before/after (ensign-transcript-aware per f5) against a baseline that can move the wrong way, plus the dispatch-golden treatment, with the riskiest mechanism spiked end-to-end
  AC-1 (adoption-not-regress, baseline can drop) + AC-2 (means, paired); spike ran the detector suite green AND source-traced the FO-only stream — metric mechanism PROVEN, ensign before/after GATED on f5; dispatch-golden clause found stale (site-6 removed in #392) and reframed as byte-identity negative control.
- DONE: A recorded note of exactly which ensign-core sections this trim edits, so the design composes with sibling `ensign-contract-dev-leakage` with no implementation-time collision
  `## Collision map` — my ensign sites (`## Working` item 1 `:18`; `## Stage Report Protocol` append bullet `:92`) are disjoint-by-section from the sibling's worktree/deliverable targets; one adjacency flagged (`:18` above the sibling's likely `:19`), mitigated by sequential Wave-1 + text-anchored edits; FO core is mine alone.

### Summary
Refined the design with verbatim before/after wording for four trim sites (three mandated + one recommended), splitting the lone AC into a measured end-value AC (adoption-not-regress) and its paired means-AC (residue trim + grep-primary consistency), both verified behaviorally per the boundary guard's prose-grep ban. The riskiest-mechanism spike is the load-bearing finding: the journeymetrics detector primitives are proven (suite green), but a real ENSIGN before/after is a vacuous `0==0` until f5 folds the dispatched sub-agent transcript in — so AC-1 has a hard, surfaced dependency on f5 (backlog, sprint 0204). Two further findings for the gate: FO `:100` carries the same redundancy the checklist didn't enumerate (recommend including), and the original AC's "dispatch-golden regeneration" is stale (the real dispatch hint, site-6, was removed in #392) so the goldens stay byte-identical as a negative control. No product/contract files were edited (ideation is design-only).

## Stage Report: implementation

- DONE: All four trim sites edited by bullet TEXT not line number (ensign `## Working` item 1; ensign `## Stage Report Protocol` append bullet; FO `## Probe and Ideation Discipline` Prefer-Grep bullet; FO `## Completion and Gates` step 1 / Site D), each naming grep primary and keeping `--read --json` only as the fence-safe fallback (A,C) or dropping it for `wc -l`/grep (B,D)
  Code commit `21675518` on branch `spacedock-ensign/read-guidance-redundant-with-grep`; `git diff` = exactly 4 insertions / 4 deletions matching the `## Before/after wording` AFTER text verbatim. A,C keep `--read --json` as fence-safe fallback; B uses `wc -l`; D greps `## Stage Report` headings + scoped Read to EOF.
- DONE: AC-2 behavioral negative control green — `go test ./internal/dispatch/ -run TestBuild` goldens stay BYTE-IDENTICAL (trim is runtime-loaded skill prose, not dispatch-prompt text) and `go test ./internal/contractlint/ ./internal/ensigncycle/` stays green. No prose-grep test added.
  dispatch TestBuild `ok` (goldens byte-identical → negative control confirms trim did not move the dispatch prompt); contractlint `ok`; ensigncycle `ok`. No new test added (boundary guard forbids reading a contract file to assert a phrase).

### Summary

Applied the four trim edits verbatim from the entity's `## Before/after wording`, anchoring on bullet TEXT (not line number) per the collision map with the sibling `ensign-contract-dev-leakage`; the sibling has not yet committed so the worktree was clean off origin/main. AC-2's behavioral surfaces are all green, with the dispatch-golden byte-identity standing as the negative control that the trim is runtime-loaded skill prose rather than dispatch-prompt text. AC-1 (live adoption-not-regress) is NOT exercised here — it is the captain-gated forward-looking item blocked on f5 (`f53zr0ehhzekzbgydpybaq5g`, ensign-transcript-aware metric, not landed); the trim ships and AC-1 confirms no regression in a later live run per Open decision #2.

## Stage Report: validation

- DONE: Verify AC-2 (the means, offline-provable) on the rebased worktree: `go test ./internal/dispatch/ -run TestBuild` goldens stay BYTE-IDENTICAL (negative control: the trim is runtime-loaded skill prose, not dispatch-prompt text) AND `go test ./internal/contractlint/ ./internal/ensigncycle/` stays green. No prose-grep test over the core files.
  On the rebased tree (trim `6f58941f` atop `ef53f071`, past 48g + f5/#448 + #446): dispatch TestBuild `ok` (goldens byte-identical → negative control holds even with #446's dispatch changes present), contractlint `ok`, ensigncycle `ok`. Diff vs main = the 2 core `.md` files only (4 ins / 4 del); no test file added.
- DONE: MEASURE AC-1 (adoption-not-regress) with f5's now-merged ensign-transcript-aware metric; if the small-count metric is genuinely inconclusive, apply the captain-ratified pre-agreed degradation: ship the trim, record AC-1 as a forward-looking non-regression check, and report the measured numbers explicitly.
  f5 (#448, `a97ed24a`) is merged to main; full `journeymetrics` suite green incl. `TestFoldEnsignReadAdoption` (non-tautological by perturbation). Self-measured six real on-disk ensign transcripts with the merged parser: status_read+scoped_read sums = 4, 5, 7, 0, 1, 4 (status_read=0 in 5/6 even under un-trimmed cores). Instrument is LIVE / non-vacuous (no longer the pre-f5 `0/0`), but single-digit counts with 0→7 per-journey variance make a single paired before/after noise-dominated → genuinely inconclusive. Applied the degradation: ship; AC-1 recorded as forward-looking non-regression check. No fresh paired live journey run — disproportionate to a 4-bullet prose trim and inconclusive by construction at this magnitude.
- DONE: Confirm the rebase onto current main is clean (the four trim sites, anchored on bullet TEXT, survive intact after 48g+f5), then deliver a PASSED/REJECTED recommendation with evidence.
  origin/main advanced to `ef53f071` (#446) mid-validation; #446 touches only dispatch/claudeteam Go, not the core `.md` files. `git merge-tree` conflict-free; rebased clean — all four AFTER bullets (ensign `## Working` item 1, `## Stage Report Protocol` append; FO Prefer-Grep, FO Completion step 1) survive verbatim. Recommendation: PASSED.

### Summary

PASSED. AC-2 (the means) is fully proven behaviorally on the real merge target: the trim is exactly the four `## Before/after wording` AFTER bullets, dispatch goldens stay byte-identical (the negative control that this is runtime-loaded skill prose, not dispatch text), and contractlint+ensigncycle stay green — no prose-grep test added. AC-1's instrument (f5 #448) is now merged and verified live: re-parsing six real ensign transcripts yields real single-digit read-adoption counts (no longer the vacuous pre-f5 `0/0`), and notably status_read=0 in 5/6 even under the UN-trimmed cores — direct corroboration of the task's thesis that the `status --read` adoption guidance was already redundant with grep/scoped-Read. Because those counts are single-digit with 0→7 per-journey variance, a single paired before/after is noise-dominated and genuinely inconclusive, so I applied the captain-ratified degradation: the trim ships and AC-1 stands as a forward-looking non-regression check on accumulated telemetry. The rebase onto current main (`ef53f071`, past 48g + f5 + #446) is conflict-free with all four trim sites intact.

## Feedback Cycles

### Cycle 1 — detached adversarial audit found a MATERIAL Site-D regression (2026-06-30)

- **Reviewer:** `audit-82k-contract`, a detached adversarial audit run on a throwaway checkout (`.worktrees/audit-82k-contract` @ `6f58941f`), binary-built and exercised on real bodies — not a prose read.
- **Verdict: MATERIAL.** The trim regresses the gate-critical FO read at **Site D** (`first-officer-shared-core.md` `## Completion and Gates` step 1 — the last-`## Stage Report` read that feeds the completion cross-check, the done/skipped/failed count, and the captain-facing gate citation).
- **Why Site D specifically:** it is the only trimmed site that BOTH uses the UNANCHORED `grep -n '## Stage Report'` AND drops the fence-safe `status --read` fallback (Sites A & C kept the fallback and use the anchored `^#{1,4} ` form). Two exercised failures the OLD `status --read` handled correctly:
  - *Inline-code, no fence (common):* a report body that merely mentions `` `## Stage Report: validation (cycle 2)` `` in prose → unanchored grep's last match is the prose line → FO reads only that line to EOF → misses every checklist item → computes "0 done, 0 skipped, 0 failed" → falsely bounces a complete report.
  - *Fenced heading:* `status --read` is fence-aware (`internal/status/section_read.go scanHeadings` skips fenced lines); grep is not. A `## Stage Report` inside a ``` ``` ``` fence wins the last-match and the FO reads from there, missing the real items above. The dev README itself over-counts 23 (grep) vs 18 (fence-aware) headings.
  - Trigger class = the entire 0240-lean-contract sprint (entities whose subject IS the stage-report machinery), this entity included once a redo is narrated.
- **Secondary (Site B, `ensign-shared-core.md` `## Stage Report Protocol` append):** `wc -l` undercounts by 1 vs the old `total_lines` on a file with no trailing newline → "append after the last line" lands one line early, splitting the prior report. Minor (tooling bodies usually end in a newline) but a genuine correctness regression vs the old count.
- **Clean:** Sites A & C retain grep-primary AND the fence-safe `status --read --json` fallback (anchored form, immune to inline-code/fence). The load-bearing structured uses (README taxonomy parse at FO:17; `--checklist`/`--ac-scan` roll-ups at FO:121) were correctly retained.
- **Fix routed to implementation (this cycle):** REVERT Site D to the original `status --read <ref> --json` last-`## Stage Report` heading read — it is gate-critical and fence-sensitive, therefore NOT redundant with grep. REVERT Site B to the newline-robust `status --read --json` `total_lines` append point. Keep Sites A & C trimmed (grep-primary + fence-safe fallback). The corrected scope: trim only the plain-section-read guidance where grep + the fence-safe fallback genuinely suffices; keep `status --read` where its fence-awareness / robust count is load-bearing. Re-audit Site D after the fix before re-gating.
