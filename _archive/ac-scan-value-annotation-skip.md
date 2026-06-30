---
title: "status --ac-scan skips ACs annotated inside the bold (e.g. **AC-1 (VALUE)**), hiding the value AC from the gate cross-check"
status: done
score: 0.45
source: "v0.23.0 cut FO session, 2026-06-30. During the 3p validation gate, `status --read <ref> --stage validation --ac-scan` reported only AC-2 and AC-3 — it silently SKIPPED AC-1 because `**AC-1 (VALUE)**` carries the (VALUE) annotation INSIDE the bold markers, which the scan's `**AC-N**` matcher does not catch. The value AC (the most important one) was invisible to the deterministic gate AC cross-check; the FO confirmed AC-1 evidence manually, so 3p was unaffected, but the automated cross-check is weakened for ANY annotated AC."
id: 48gz5715kc4d2j687jbags7v
sprint: 0240-lean-contract
group: tooling
started: 2026-06-30T16:55:24Z
worktree: .worktrees/spacedock-ensign-ac-scan-value-annotation-skip
mod-block:
pr: pr-merge:447
verdict: passed
completed: 2026-06-30T19:39:05Z
archived: 2026-06-30T19:39:05Z
---

`spacedock status --read <ref> --stage <stage> --ac-scan` enumerates `**AC-N**` items and reports each one's evidence/unevidenced status, feeding the gate AC cross-check. Its matcher only recognizes a bare `**AC-N**` token, so an AC whose bold span carries an annotation — `**AC-1 (VALUE)**`, `**AC-2 (no-regression)**`, etc. — is NOT enumerated and is silently dropped from the scan.

## Problem
The README ideation policy explicitly encourages a `(VALUE)`-tagged AC ("At least one AC must MEASURE the end-value"), and the contract's AC cross-check re-anchors on exactly that value AC. So the convention the workflow recommends (`**AC-1 (VALUE)**`) is the convention `--ac-scan` cannot see — the deterministic extraction drops the single most important AC. Live evidence: on `3p`, `--ac-scan` returned only AC-2 and AC-3; AC-1 (VALUE) was absent, even though it was present and evidenced. The gate held only because the FO cross-checked AC-1 by hand.

## Proposed approach
Broaden the `--ac-scan` AC heading matcher so the closing `**` need not immediately follow the id: allow an asterisk-free trailing label inside the bold span, so `**AC-1 (VALUE)**` and `**AC-1**` both enumerate as `AC-1`. The id capture stays exact (`AC-1`, `AC-2`, …); anything after it inside the bold is consumed as a label and discarded — the scan reports id + line only, so no new capture group is needed.

Exact change — `internal/status/gate_extract.go:52`:

    // before
    var acHeadingRe = regexp.MustCompile(`\*\*(AC-[0-9A-Za-z]+)\*\*`)
    // after
    var acHeadingRe = regexp.MustCompile(`\*\*(AC-[0-9A-Za-z]+)[^*]*\*\*`)

`[^*]*` matches the label up to the closing `**`; because it excludes `*` it can never span a `**` boundary, so two headers on one line (`**AC-1** … **AC-2**`) still enumerate separately and never merge. The heading-doc comment (lines 49–52) updates to describe the trailing-label allowance. Only `acHeadingRe` (the heading enumerator) changes; `acTokenRe` (the in-checklist citation scanner, `\bAC-[0-9A-Za-z]+\b`) already matches annotated ACs via its word boundaries and stays as-is.

## Acceptance criteria
- **AC-1 (the value)** — an AC written with the annotation INSIDE the bold (`**AC-1 (VALUE)**`) is enumerated by `status --ac-scan` with id `AC-1` and its evidence/unevidenced status reported, exactly as a bare `**AC-1**` would be — so the value AC the README ideation policy recommends is no longer dropped from the gate cross-check. Verified by: (a) a fixture-driven Go test feeding an entity body whose `## Acceptance criteria` uses `**AC-1 (VALUE)**`, asserting `--ac-scan --json` lists `id=AC-1` — RED on the current matcher (AC-1 absent), GREEN after the fix; AND (b) the independent on-disk baseline that moved the wrong way — the real 0240 body `fn-binding-refinements.md`, whose `--ac-scan` today enumerates 7 ACs and SILENTLY DROPS `**AC-4 (value guardrail)**` (confirmed live, see Spike record), enumerates all 8 (AC-4 present) after the fix.
- **AC-2 (no over-match)** — the broadened matcher does not over-capture: in the same fixture a bare `**AC-2**` enumerates exactly once, and a bare-prose mention in the AC section ("see AC-3 above") is NOT treated as an AC heading (AC-3 absent from the scan). Verified by: the same fixture-driven test asserting the enumerated id set is exactly {AC-1, AC-2} — AC-2 once, no AC-3 — proving the `**…**` delimiter requirement survived the broadening.

## Test plan
- **Fixture:** a new committed fixture `internal/status/testdata/section-reader/annotated-ac-fixture.md` — minimal frontmatter, a `## Acceptance criteria` section with `- **AC-1 (VALUE)** …` (plus a "see AC-3 above" prose mention) and a bare `- **AC-2** …`, and a `## Stage Report: ideation` section whose DONE evidence cites AC-1 (so AC-1's unevidenced flag is exercised). Modeled on `interleaved-fixture.md` and the `cycle3_extract_test.go` harness (`runNative` + `acScanEnvelope`).
- **Test (covers both ACs):** `TestACScanEnumeratesAnnotatedAC` in `internal/status/cycle3_extract_test.go` (or a sibling `*_test.go`) drives `--read <fixture> --stage ideation --ac-scan --json` and asserts the enumerated id set is exactly {AC-1, AC-2}: AC-1 present (value AC → AC-1), AC-2 present exactly once, AC-3 absent (no over-match → AC-2). RED on the bare matcher, GREEN after broadening `acHeadingRe`.
- **Cost/complexity:** trivial — one fixture + one Go unit test, pure binary-exercise, no live host. Reuses the existing `runNative`/`acScanEnvelope` harness.
- **No live workflow / no CLI golden needed:** the `--ac-scan` JSON envelope schema is unchanged (same `id`/`line`/`unevidenced`/`citations` fields); only WHICH ACs appear changes, which the fixture test covers directly.

## Documentation changes
The dev README's validation-stage cross-check prose describes the same "pull every AC" behavior the matcher implements, and currently names only the bare form — the same doc-vs-convention gap this entity exists to close. Concrete diff — `docs/dev/README.md:131`:

    - Pull every `**AC-N**` item from the entity body's `## Acceptance criteria` section; reproduce the evidence cited in each "Verified by" clause; flag any AC without evidence.
    + Pull every `**AC-N**` item (including a value annotation inside the bold, e.g. `**AC-1 (VALUE)**`) from the entity body's `## Acceptance criteria` section; reproduce the evidence cited in each "Verified by" clause; flag any AC without evidence.

`docs/site/reference/command-reference.md:46` describes `--ac-scan` as "per-AC evidence citations" with no header-format detail, so it stays accurate — no change.

## Spike record
Spiked the riskiest path (does the broadened regex catch annotated headers AND not over-capture?) against REAL 0240 bodies and discriminators, before any product edit:

- **Bug confirmed live.** `spacedock status --read fn-binding-refinements.md --stage ideation --ac-scan` enumerates AC-1/2/3/5/6/7/8 but DROPS `**AC-4 (value guardrail)**` — the value AC absent from the deterministic gate cross-check, exactly the reported defect.
- **Broadening confirmed.** A throwaway Go program applied the current (`\*\*(AC-[0-9A-Za-z]+)\*\*`) vs broadened (`\*\*(AC-[0-9A-Za-z]+)[^*]*\*\*`) regex to the real annotated headers (`**AC-4 (value guardrail)**`, `**AC-2 (end value — measured…)**`, an `**AC-2d (calendar)**` trailing-letter id): current captures NONE; broadened captures all with the exact id (AC-4, AC-2, AC-2d). Bare `**AC-1**` still matches under both. Since `scanACs` is unchanged, AC-4's line is the only previously-unmatched AC line in `fn-binding-refinements.md`, so the enumerated set deterministically grows 7→8 with AC-4 added — confirmed end-to-end by AC-1(b)'s RED-then-GREEN test at implementation.
- **No over-match confirmed.** Discriminators: `**AC-2** … see AC-3 above` → broadened captures only `[AC-2]`; `**AC-1** and **AC-2**` on one line → `[AC-1 AC-2]` (no merge across the `**` boundary); `**AC-1 (VALUE)** … plain prose AC-9` → only `[AC-1]`.

This is a localized matcher broadening (one regex literal + its doc comment); the mechanism is fully exercised. Known boundary, out of scope: an annotation containing a literal `*` (e.g. `**AC-1 (*x*)**`) is not a real AC convention and is not matched by `[^*]*`.

## Stage Report: ideation

- DONE: REFINE and COMPLETE the body — the exact matcher change + the test plan
  Proposed approach now pins the one-line change at `gate_extract.go:52` (`\*\*(AC-[0-9A-Za-z]+)\*\*` → `\*\*(AC-[0-9A-Za-z]+)[^*]*\*\*`); Test plan names the fixture + RED-then-GREEN test reusing the `runNative`/`acScanEnvelope` harness.
- DONE: The value AC (AC-1) is a RED-then-GREEN fixture test (`**AC-1 (VALUE)**` enumerated as `AC-1`, exactly as bare `**AC-1**`)
  AC-1 rewritten as the value AC with the fixture RED/GREEN proof plus an independent on-disk baseline (real `fn-binding-refinements.md` 7→8 ACs, AC-4 surfaces).
- DONE: The no-over-match control (AC-2) — bare `**AC-2**` once AND prose "see AC-3 above" not a new AC item
  AC-2 rewritten as the discriminator: enumerated id set is exactly {AC-1, AC-2}; the prose AC-3 mention is not enumerated, proving the `**…**` delimiter requirement survived.
- DONE: Locate the `--ac-scan` matcher and confirm the broadening on the REAL 0240 entity bodies; record this as the spike
  Located `acHeadingRe` at `gate_extract.go:52` (sole definition, no Node twin). Spike record: live bug confirmed on `fn-binding-refinements.md` (AC-4 dropped) + throwaway regex spike over the real annotated headers and discriminators (`scratchpad/acscan_spike.go`).

### Summary

Refined the ideation design for the `--ac-scan` annotated-AC fix to implementation-ready. The matcher is a single regex at `internal/status/gate_extract.go:52`; broadening it to `\*\*(AC-[0-9A-Za-z]+)[^*]*\*\*` allows an asterisk-free trailing label inside the bold so `**AC-1 (VALUE)**` enumerates as `AC-1`, with `[^*]*` guaranteeing no over-match across a `**` boundary. Confirmed the bug live (real `fn-binding-refinements.md` silently drops `**AC-4 (value guardrail)**`) and exercised the broadened regex against the real annotated headers + discriminators before any product edit. Added a value-measuring AC-1 (real-entity 7→8 baseline), a no-over-match AC-2, a fixture-backed test plan, and a one-line `docs/dev/README.md:131` doc diff. No Node twin and the `--ac-scan` JSON schema is unchanged.

## Stage Report: implementation

- DONE: acHeadingRe at internal/status/gate_extract.go:52 broadened to allow an asterisk-free trailing label inside the bold (the [^*]* form) + its heading-doc comment updated, so **AC-1 (VALUE)** enumerates as AC-1.
  `\*\*(AC-[0-9A-Za-z]+)\*\*` → `\*\*(AC-[0-9A-Za-z]+)[^*]*\*\*` (now at gate_extract.go:57, comment lines 49-56 describe the trailing-label allowance + the no-`**`-boundary-span guarantee); commit 3fff9062.
- DONE: TestACScanEnumeratesAnnotatedAC + annotated-ac-fixture.md added: RED on the bare matcher, GREEN after; enumerated id set is exactly {AC-1, AC-2} (no over-match, prose "see AC-3" not enumerated) AND the real fn-binding-refinements body goes 7->8 ACs (AC-4 surfaces).
  RED: bare matcher emits only `{AC-2:1}` (AC-1 dropped); GREEN: exactly {AC-1, AC-2}, AC-1 unevidenced=false. Real `fn-binding-refinements.md --stage ideation --ac-scan`: 7 ACs (AC-1/2/3/5/6/7/8) → 8 with AC-4 at line 37, no other id changed — confirmed via rebuilt binary (state checkout absent from worktree, so a binary exercise per AC-1(b)/Spike, not a brittle committed test). Full `internal/status` suite green.
- DONE: docs/dev/README.md:131 doc diff applied (the validation cross-check line gains the "including a value annotation inside the bold" clause).
  Line 131 now reads "Pull every `**AC-N**` item (including a value annotation inside the bold, e.g. `**AC-1 (VALUE)**`) …"; commit 3fff9062.

### Summary

Broadened the sole `--ac-scan` heading matcher (`acHeadingRe`, gate_extract.go) to `\*\*(AC-[0-9A-Za-z]+)[^*]*\*\*` so an annotation inside the bold (`**AC-1 (VALUE)**`) enumerates as `AC-1`; `[^*]*` excludes `*` so it never spans a `**` boundary (no over-match). Drove the change RED→GREEN with `TestACScanEnumeratesAnnotatedAC` + `annotated-ac-fixture.md` (id set exactly {AC-1, AC-2}, prose "see AC-3" not enumerated) and confirmed the real 0240 `fn-binding-refinements.md` goes 7→8 ACs with the previously-dropped `**AC-4 (value guardrail)**` surfacing. README:131 validation cross-check prose updated; `go vet`/`go build ./...`/full status suite all green.

## Stage Report: validation

**Recommendation: PASSED.**

- DONE: MEASURE AC-1 (the value): reproduce live on a rebuilt binary that real fn-binding-refinements.md --stage ideation --ac-scan enumerates 7 ACs on origin/main's matcher and 8 (AC-4 surfaces) after the fix, AND the fixture is RED on the bare matcher (AC-1 dropped) / GREEN after (AC-1 enumerated) — measured by running the binary, not asserted from prose.
  Built two binaries (fixed=worktree HEAD 3fff9062; baseline=`origin/main` 2ef249e9 throwaway checkout, confirmed bare matcher at gate_extract.go:52). On real `fn-binding-refinements.md`: baseline=7 `[AC-1,2,3,5,6,7,8]` (AC-4 dropped, 36→38), fixed=8 `[AC-1,2,3,4@37,5,6,7,8]` — AC-4 surfaces, no other id moved. Fixture RED/GREEN: baseline emits only `[AC-2]` (AC-1 absent → RED); fixed emits `[AC-1,AC-2]`, AC-1 unevidenced=false (GREEN).
- DONE: MEASURE AC-2 (no over-match): the fixture's enumerated id set is exactly {AC-1, AC-2} — AC-2 once, prose "see AC-3 above" NOT enumerated, and two headers on one line do not merge across the ** boundary (the [^*]* guarantee).
  Fixed binary on `annotated-ac-fixture.md`: id set exactly {AC-1, AC-2}, no dupes, AC-3 absent. Crafted a separate discriminator and ran the fixed binary: `**AC-1** and **AC-2**` on one line → both enumerate separately (AC-1@9, AC-2@9, no merge — the `[^*]*` guarantee holds); `**AC-3 (VALUE)** … plain prose AC-9` → only AC-3 (AC-9 not enumerated); bare prose `see AC-7 above` → not enumerated.
- DONE: Reproduce each AC's "Verified by" clause; run go test ./internal/status/... green; confirm the docs/dev/README.md:131 doc diff is present. Deliver a PASSED/REJECTED recommendation with evidence.
  AC-1 (a) fixture RED→GREEN + (b) real-body 7→8 both reproduced by binary exercise above (proof is command/state, not self-referential). `go test ./internal/status/...` → `ok` (69.4s); `TestACScanEnumeratesAnnotatedAC` PASS. README:131 now carries "(including a value annotation inside the bold, e.g. `**AC-1 (VALUE)**`)" — diff present in commit 3fff9062. Worktree clean; commit touches exactly gate_extract.go + test + fixture + README.

### Summary

PASSED. Every AC was MEASURED by running rebuilt binaries, not by prose-grep. The reported defect reproduces on the baseline matcher (real `fn-binding-refinements.md` silently drops `**AC-4 (value guardrail)**`, 7 ACs; the fixture drops `**AC-1 (VALUE)**` entirely) and is fixed by the broadening (AC-4 surfaces → 8; fixture → exactly {AC-1, AC-2}). The `[^*]*` no-over-match guarantee was verified directly: two headers on one line enumerate separately and prose/annotation `*`-free mentions never become headings. The committed test is non-tautological (genuinely RED on the baseline binary). README:131 doc sync present, suite green, change is localized and low-blast-radius (one regex literal + doc comment + test + fixture).
