---
title: Gate-extract verbs — structured extraction of stage report, AC coverage, and reviewer findings (the verdict stays level-3)
status: ideation
source: 'FO shaping 0205 (2026-06-15, this session) — the FO gate prep (checklist DONE/SKIPPED/FAILED review, the AC cross-check, the Material/Polish reviewer tiers) is structured extraction over a stage-report file: deterministic, but done by the model today. A weak FO can do it reliably only as a binary call returning structured data; the verdict (judgment) stays level-3. No 2y dependency. 0205 layered-FO, no-2y parallel track.'
started: 2026-06-16T02:21:21Z
completed:
verdict:
score: 0.4
worktree:
issue:
sprint: 0205-layered-fo
sprint-readiness: defer
id: 6reqad9gff9wk544det3x4fj
---

Three deterministic extraction verbs so a weak FO assembles a gate from structured data instead of re-parsing prose. The verdict is NOT computed here; it routes to level-3.

## Problem

Gate preparation today is the model reading a stage report and a `## Acceptance criteria` section and extracting: the checklist accounting, each AC's evidence citation, and the reviewer findings tiered Material vs Polish. This is deterministic extraction a binary should own, so a Haiku FO gets the same structured input every time instead of re-deriving it (and mis-deriving it).

## Proposed approach

Three deterministic extraction verbs under a new `spacedock gate` command, each reading a stage-report-bearing entity file and emitting the all-strings ordered JSON envelope the workflow already uses (`internal/status/json.go`: every leaf a string, insertion-order keys, compact + newline-terminated). The text default mirrors the JSON for parity, as `--read` does. None of the three computes a verdict.

**Surface — a new `gate` command, not a `dispatch` subcommand.** `spacedock gate checklist-extract --entity-path FILE`, `spacedock gate ac-scan --entity-path FILE`, `spacedock gate reviewer-parse --entity-path FILE`, each `--json` for the envelope. Rationale: `dispatch` builds the worker artifact (the prompt the ensign runs); these verbs read back the completed worker's report at the gate. Different phase of the loop, different noun — a sibling command keeps the surfaces honest. The verbs share one stage-report parser internally.

**Verb 1 — `checklist-extract`.** Locates the LAST `## Stage Report` section (the latest cycle) and emits each `- DONE:`/`- SKIPPED:`/`- FAILED:` item with its status, the bullet text, and the 1-based `start`/`end` line range it owns (the bullet line through its trailing evidence line, stopping at the next bullet or the `### Summary`/next-section heading). `start`/`end` use the same 1-based line semantics as `status --read` so the FO can `Read(offset, limit)` an item's evidence. Also emits a `chosen_direction_required` flag: true when the stage definition declares a chosen-direction concept (ideation picks an approach, validation picks PASS/REJECTED) so the present-gate `Chosen direction:` line is required rather than `n/a`. The flag is read from the workflow stage definition (the same `--workflow-dir`/`--stage` inputs `show-stage-def` takes), NOT guessed from the report prose.

**Verb 2 — `ac-scan`.** For each `**AC-N**` heading in the entity's `## Acceptance criteria` section, emits the AC id, its 1-based heading line, and a citation list: every `AC-N` token mention **found within a checklist item's line range** (the ranges verb 1 computes), each with the citing line number and line text. An AC with zero such citations is flagged `unevidenced: true`. **The citation search is scoped to checklist evidence lines, never the `### Summary` or `**Reviewer findings:**` prose** — this is the decisive boundary the spike below settled: a reviewer *complaining* "AC-2 has no evidence" must not be counted as AC-2 *being* evidenced. Also emits a `natural_place` flag per AC: true when this stage's `Outputs:` make it the stage where that AC is expected to gain evidence (read from the stage definition), so level-3 knows which unevidenced ACs are a reject signal here versus expected-later. **The verb emits the citation map and the flags; it does NOT decide whether an AC is *satisfied*.** Whether a red-test citation is evidence-of-failure or evidence-of-progress is judgment that routes to level-3.

**Verb 3 — `reviewer-parse`.** Parses a `**Reviewer findings:**` block (the present-gate two-tier convention) into `Material:` and `Polish:` arrays of finding text. Emits `{material: [...], polish: [...]}`; an absent block yields both empty. The tier names are the present-gate contract's fixed labels — Material (fact-corrections, contract violations, missing AC evidence, broken claims) vs Polish (wording, format drift) — so the FO renders the gate's `Reviewer findings` tiers from structured data instead of re-tiering prose.

**Why deterministic extraction, not a verdict.** The present-gate assembly rules already define exactly these three structures (checklist roll-up with a line-range citation, AC cross-check, Material/Polish reviewer tiers). A capable FO derives them by reading; a Haiku FO re-derives them unreliably. A binary that returns the same structured input every time removes that variance. The verdict — approve/reject, is-this-AC-satisfied, is-this-chosen-direction-sound — stays judgment and routes to level-3 via `fo-tier-delegation`. These verbs feed L3 a clean structured input; they never adjudicate.

## Riskiest-mechanism spike (DONE — exercised on a known-structure fixture, not asserted)

The riskiest unverified path was whether the three extractions are actually deterministic over a real stage-report shape AND compose cleanly — specifically whether `ac-scan`'s "is this AC evidenced" can be answered without sliding into the verdict it must not compute. Spiked with a throwaway Go parser (`/tmp/gate-extract-spike/spike4.go`) over a fixture with a KNOWN structure (`/tmp/gate-extract-spike/fixture.md`; a stage report with a DONE, a SKIPPED, and a FAILED item, two ACs where AC-1 is cited in a DONE evidence line and AC-2 is mentioned only in the reviewer findings, and a Material/Polish reviewer block — line numbers fixed by `cat -n`).

The spike found and settled three design decisions the prose alone would have missed:

1. **AC-id token boundary.** A naive split on `-` truncated `AC-1` to `AC`. The id must be tokenized on the heading regex `^\*\*(AC-[0-9A-Za-z]+)\b`, not split on `-`. Confirmed: the corrected parse emits `AC-1`/`AC-2`.
2. **Citation scope.** Scanning the whole `## Stage Report` umbrella swept the `### Summary` and `**Reviewer findings:**` text, double-counting AC-2 (the reviewer's "AC-2 is unsatisfied" complaint registered as an AC-2 citation — a perverse inversion). Scoping the citation search to the checklist items' own line ranges fixed it: AC-1 cited at line 24 (its DONE evidence), AC-2 cited at `[]` → `unevidenced: true`. This is why the three verbs compose: `checklist-extract`'s line ranges bound `ac-scan`'s search.
3. **Verdict boundary.** AC-2 *is* mentioned in the report, but only as red/unsatisfied. A token-match "evidenced" boolean would mislabel it. The verb therefore emits the raw citation map (with surrounding line text) and the `unevidenced` flag; whether a citation is evidence-of-satisfaction is L3's call. The spike proved the extraction stops exactly at the judgment boundary.

Final spike output against the oracle: checklist `DONE [23-24] / SKIPPED [25-26] / FAILED [27-28]`; `AC-1 cited at [24] (unevidenced=false)`, `AC-2 cited at [] (unevidenced=true)`; reviewer `Material` ×2, `Polish` ×1. The line ranges match the fixture's `cat -n` exactly. The implementation's first tests seed directly from this exercise.

## Out of scope

- **The verdict decision (approve/reject), the satisfied/unsatisfied AC determination, and chosen-direction soundness** — these are judgment and route to level-3 via `fo-tier-delegation`. The verbs feed L3 structured input; they never adjudicate.
- **`spacedock gate assemble-verdict` as a binary** — the verdict decision tree may stay a prose-function this sprint (sprint index open question). These three extract verbs ship; assembly stays prose with the L3 route.
- **The rest of the FO loop** — `next-action`, state-verbs, merge-finalize, the prose-function restructure are separate sprint members.
- **2y (`shared-merge-dispatch-contract`)** — no dependency; this member is the no-2y-startable slice and proceeds in parallel.

## Acceptance criteria

All criteria are entity-level properties of the finished verbs, each proven by EXERCISING real output against a committed fixture's KNOWN structure (the external oracle — the fixture's `cat -n` line numbers and enumerated contents), never a prose/substring/regex match. The verdict is explicitly NOT asserted by any AC (it is not computed).

**AC-1 — `checklist-extract` emits each stage-report item's status and exact line range for the latest cycle.**
Given a committed fixture with a DONE, a SKIPPED, and a FAILED item under the last of two `## Stage Report` sections, `gate checklist-extract --json` emits exactly those three items with their `status` and 1-based `start`/`end` line ranges, and the ranges, when sliced from the fixture (`lines[start-1:end]`), equal the item's known bullet+evidence text. Items from the EARLIER cycle's report are absent (latest-cycle only).
Verified by: Go test slicing the fixture by each emitted range and asserting byte-equality against the known item text (table-driven over the fixture's three items), plus an assertion that the cycle-1 item is absent.

**AC-2 — `ac-scan` cites AC evidence from checklist lines only and flags the unevidenced AC.**
Given a fixture whose AC-1 is cited in a DONE evidence line and whose AC-2 is mentioned ONLY in the reviewer-findings prose, `gate ac-scan --json` emits AC-1 with a citation at the DONE line and `unevidenced: false`, and AC-2 with an empty citation list and `unevidenced: true` (the reviewer-prose mention is NOT counted). No `satisfied` field is emitted.
Verified by: Go test asserting the citation line numbers against the fixture's known DONE-line position and asserting AC-2's empty citation list / `unevidenced: true`; a mutation test that widening the citation scope to the whole stage report (counting the reviewer prose) flips AC-2 to evidenced and FAILS the test — proving the scope boundary is load-bearing, not incidental.

**AC-3 — `reviewer-parse` tiers the findings into Material and Polish; an absent block yields empty tiers.**
Given a fixture with a `**Reviewer findings:**` block carrying two Material items and one Polish item, `gate reviewer-parse --json` emits `material` with those two items and `polish` with the one, in order. Given a fixture with NO reviewer block, both arrays are empty (not an error).
Verified by: Go test asserting the parsed `material`/`polish` arrays equal the fixture's enumerated finding texts; a second test over a no-reviewer fixture asserting both empty and exit 0.

**AC-4 — `checklist-extract` and `ac-scan` read `chosen_direction_required` and `natural_place` from the stage definition, not the report prose.**
Given `--workflow-dir`/`--stage` resolving to a stage definition that declares a chosen-direction concept (ideation), `checklist-extract` emits `chosen_direction_required: true`; for a stage without one, `false`. `ac-scan` emits `natural_place: true` for an AC the stage's `Outputs:` make it the natural place to evidence, `false` otherwise. Flipping only the stage definition (same report) flips the flags.
Verified by: Go test driving both verbs against two stage definitions (one with, one without a chosen-direction/natural-place mapping) over the SAME report fixture, asserting the flags track the stage definition, not the report.

**AC-5 — A missing or malformed input fails loudly with a non-zero exit and a named diagnostic; never a partial/silent emit.**
A missing `--entity-path`, an unreadable file, or an absent `## Stage Report` / `## Acceptance criteria` section exits non-zero with a stderr diagnostic naming the missing input — matching the `show-stage-def` loud-failure idiom (exit 1 / exit 2 for a usage error).
Verified by: Go test asserting the exit code and the stderr diagnostic for each missing-input case (no `--entity-path`; nonexistent file; a report-less fixture).

## Test plan

- **Mechanism (already paid):** the spike above — throwaway parser over a known-structure fixture — settled the AC-id boundary, the citation-scope boundary, and the verdict boundary. Seeds AC-1/AC-2's first tests. Cost: done.
- **Unit (Go, `internal/gate` or the chosen package):** AC-1–AC-3 + AC-5 as table-driven tests over THREE committed stage-report fixtures — a passing report (all DONE, every AC cited), a report with a FAILED item, and a report with an unevidenced AC (AC mentioned only in reviewer prose). Fixtures live under the package's `testdata/`. The oracle is each fixture's known structure (line numbers, enumerated items/ACs/findings), computed in-test — never a golden byte file that re-breaks on every fixture edit. AC-2 carries a mutation test (widen the citation scope → AC flips → test fails) so the scope boundary is proven non-tautological. Cost: low — pure-function parser, reads only the fixtures; minutes.
- **Behavior (Go, drives the binary):** AC-4 — the two verbs against two stage definitions over one report, asserting the stage-definition-sourced flags. Cost: low; reuses the `show-stage-def` workflow-dir/stage resolution.
- **No live workflow test:** these are pure extraction over committed files; there is no runtime handoff to exercise. The verdict (the only judgment) is out of scope and routes to L3, so no live drive is owed by THIS member — the live Haiku drive that exercises the L3 routing is `haiku-drive-validation`'s proof, not this verb member's.
- **The verdict is explicitly NOT asserted** by any test — no AC computes or checks approve/reject or satisfied/unsatisfied. A test that asserted a verdict would be testing a behavior this member does not own.

## Documentation changes

This adds a user-visible CLI surface (a new `spacedock gate` command), so ideation proposes the doc diff; implementation applies it.

**`docs/site/reference/command-reference.md`** — under `## Workflow`, add one row to the command table (addition only, no rewording of existing rows):

> **Before** (last data row, unchanged): `| \`spacedock completion\` | Print a bash or zsh completion script |`
> **After** — add a new row above `completion` (after the `state` row):
> `| \`spacedock gate\` | Extract structured gate input from a completed stage report (\`gate checklist-extract\`, \`gate ac-scan\`, \`gate reviewer-parse\`) — the deterministic structure the first officer assembles a gate from; the verdict stays a judgment call |`

The bash/zsh completion verb lists in `internal/cli/cli.go` (`bashCompletion`/`zshCompletion`) gain `gate` in the `verbs` list and a `gate) compadd checklist-extract ac-scan reviewer-parse` case — implementation applies it under the completion tests, not a separate doc.

No FO/ensign contract wording changes are proposed at ideation: the verbs are read by the FO at the gate, and the present-gate skill's three structures (checklist roll-up, AC cross-check, reviewer tiers) ALREADY name what these verbs return. The adoption wiring (pointing present-gate's assembly at the verbs) is the `prose-function-restructure` member's concern (`«gate.assemble-verdict»` body), not this extraction-verb member's — flagged here so the gate review can confirm the boundary.

## Stage Report: ideation

- DONE: Design the three extraction verbs (checklist-extract / ac-scan / reviewer-parse) as DETERMINISTIC structured output over a stage-report file — checklist DONE/SKIPPED/FAILED + line ranges + chosen-direction-required flag; per-AC evidence citations + a natural-place flag for level-3 routing; Material/Polish reviewer tiers. The verdict is NOT computed here (routes to level-3 via fo-tier-delegation).
  Body "Proposed approach": three verbs under a new `spacedock gate` command, all-strings JSON envelope (the `internal/status/json.go` idiom). checklist-extract → status + 1-based start/end ranges + `chosen_direction_required` (read from the stage def). ac-scan → per-AC citation map scoped to checklist evidence lines + `unevidenced`/`natural_place` flags, no `satisfied` field. reviewer-parse → `{material, polish}` from the present-gate two-tier block. "Out of scope" routes the verdict / satisfied-determination / chosen-direction-soundness to level-3 via fo-tier-delegation.
- DONE: Behavior-first ACs over committed stage-report fixtures (a passing report, a report with a FAILED item, a report with an unevidenced AC); assert the extracted structure against the fixture's known contents, never a prose/substring match. No 2y dependency.
  Five ACs, all oracle-based (fixture `cat -n` line numbers + enumerated contents): AC-1 slices the fixture by each emitted range and asserts byte-equality (+ latest-cycle-only); AC-2 asserts checklist-scoped citations with a MUTATION test (widen scope → AC-2 flips evidenced → test fails) proving the boundary is load-bearing; AC-3 tiers Material/Polish + empty-on-absent; AC-4 flips the stage def to prove flags are stage-def-sourced; AC-5 loud-failure on missing input. Test plan names the three fixtures (passing / FAILED-item / unevidenced-AC), explicitly NOT asserting the verdict. "Out of scope" + source line confirm no 2y dependency.
- DONE: Spike the riskiest mechanism FIRST (per the proof policy) — prove the three extractions are deterministic over a real stage-report shape and that ac-scan stops at the verdict boundary it must not cross.
  Body "Riskiest-mechanism spike (DONE)": throwaway Go parser (`/tmp/gate-extract-spike/spike4.go`) over a known-structure fixture settled three design decisions the prose would have missed — the AC-id token boundary (regex, not split-on-`-`), the citation-scope boundary (checklist lines only, NOT the reviewer prose, or the reviewer's "AC-2 unevidenced" complaint inverts into an AC-2 citation), and the verdict boundary (emit the raw citation map + `unevidenced` flag, never a satisfied boolean). Output matched the fixture's `cat -n` exactly: checklist DONE[23-24]/SKIPPED[25-26]/FAILED[27-28], AC-1 cited[24]/AC-2 cited[] unevidenced.
- DONE: Propose the documentation diff for the new user-visible CLI surface (per the ideation rule for behavior-changing tasks).
  Body "Documentation changes": a concrete before/after adding one `spacedock gate` row to the `## Workflow` command table in `docs/site/reference/command-reference.md` (addition only); the `internal/cli/cli.go` completion verb-list/case update flagged for implementation under the completion tests; and an explicit note that present-gate adoption wiring belongs to the prose-function-restructure member, not this one.

### Summary

Designed `gate-extract-verbs` as three deterministic extraction verbs (`spacedock gate checklist-extract` / `ac-scan` / `reviewer-parse`) emitting the workflow's established all-strings JSON envelope over a completed stage-report file, so a weak (Haiku) FO assembles a gate from the same structured input every time instead of re-parsing prose. The riskiest path was spiked FIRST with a throwaway parser over a known-structure fixture, which settled three non-obvious boundaries — the AC-id token regex, scoping AC citations to checklist evidence lines only (so a reviewer's "AC-N is unevidenced" complaint is not perversely counted as an AC-N citation), and stopping `ac-scan` at the verdict boundary (it emits a raw citation map and an `unevidenced` flag, never a satisfied/unsatisfied decision). The verdict, the satisfied-determination, and chosen-direction soundness are explicitly out of scope and route to level-3 via fo-tier-delegation; the five ACs are oracle-based over three committed fixtures (passing / FAILED-item / unevidenced-AC) with a mutation test pinning the citation-scope boundary, and no AC asserts a verdict. No 2y dependency — this is the no-2y-startable slice. A CLI doc diff for the new `spacedock gate` surface is recorded for implementation.

