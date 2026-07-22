---
id: mp2jx24h3c92ef1yz9w1tjpz
title: Fix four confirmed tautological tests in internal/ensigncycle and internal/status
status: validation
source: "Found via a per-package audit sweep applying the 'testing-without-tautologies' checklist (github.com/kenn-io/middleman skills/testing-without-tautologies/SKILL.md), each confirmed by actually applying a suggested production-code mutation and observing the test stay green (not just static reasoning). Four findings: (1) internal/ensigncycle/pty_session_test.go TestFOSessionPinning/activeSessionFile_would_flip_to_teammate — the only failure-path action is t.Logf, no t.Error/t.Fatal/t.Fail anywhere in the subtest, so it passes regardless of activeSessionFile's actual behavior. (2) internal/status/native_new_test.go TestNewFolderForm — computes its 'expected' bytes via the exact same unexported stampID function the code under test also calls internally (mirror assertion); a real splice-offset regression in stampID left it green. (3) internal/status/zz_independent_parity_test.go TestIndNewAtomicCreate — identical mirror-assertion defect against the same stampID function, independently reconfirmed. (4) internal/status/boot_probe_parity_test.go TestBootTeamStateProbeConfinement — builds its 'expected' string directly from the same production constant (teamStateNeutralHint) the code renders from; emptying that constant left it green. Scope for ideation: for each, replace the mirror/no-op assertion with one that has an independent oracle (a literal/hand-checked expected value, or an assertion with a real failure path) while preserving whatever real coverage the test does provide elsewhere in the same function."
started: 2026-07-06T15:54:51Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-tautological-test-fixes
issue:
mod-block: merge:pr-merge
pr: pr-merge:548
---

Four tests across internal/ensigncycle and internal/status assert nothing that can actually fail — confirmed by mutation-testing each one. Fix them so they have an independent oracle instead of mirroring the production logic or dropping the assertion entirely.

## Problem

A test whose "expected" value is produced by the same code it tests, or whose failure branch never fails the run, passes no matter how the behavior breaks. Four such tests were confirmed at HEAD `b7f331d5` by applying a mutation and observing the test stay green — not by static reasoning. Two defect shapes:

- **No-op assertion** — the failure branch only logs; nothing calls `t.Error`/`t.Fatal`/`t.Fail`, so the subtest passes regardless of the behavior.
- **Mirror assertion** — the test builds its `expected` value by calling the same production function (`stampID`) or reading the same production constant (`teamStateNeutralHint`) the code under test renders from, so a regression moves both sides together and the equality still holds.

The four confirmed findings (re-verified against current HEAD; line numbers current):

1. `internal/ensigncycle/pty_session_test.go:379` — `TestFOSessionPinning/activeSessionFile_would_flip_to_teammate`. The `active != teammatePath` branch runs only `t.Logf` (`:384`); no `t.Error`/`t.Fatal`/`t.Fail` anywhere in the subtest, so it passes for any return value of `activeSessionFile`. (`activeSessionFile`, `foRootSessionID`, `sessionFileByID`, `newFileLineSourceByID` are test-local offline ports — "ported from spacedock-gym, reference-only, not imported", `:188`. The subtest's real job is to LOCK that the fixture is a genuine adversarial trap: the naive newest-assistant-bearing heuristic DOES pick the teammate, which is why the sibling by-id subtests — that production actually relies on — are meaningful rather than vacuous.)
2. `internal/status/native_new_test.go:143` — `TestNewFolderForm`. `wantBytes := string(stampID([]byte(newBody), want))` computes the expected file bytes by calling `stampID`, the same function the `--new --folder` path uses to write the file. A splice-offset regression in `stampID` (`new.go:131`) moves both sides identically and leaves the byte-identity check green.
3. `internal/status/zz_independent_parity_test.go:552` — `TestIndNewAtomicCreate`. `wantSeed := string(stampID([]byte(body), expectedID))` — the identical mirror defect against `stampID`, independently reconfirmed on the parity path.
4. `internal/status/boot_probe_parity_test.go:112` — `TestBootTeamStateProbeConfinement`. `wantNil := "TEAM_STATE\npresent: false\nhint: " + teamStateNeutralHint` builds the expected nil-probe block from the same constant (`boot.go:23`) the render reads (`boot.go:201`). Emptying `teamStateNeutralHint` empties both sides and the equality holds. (The sibling `wantAbsent` at `:106` has the SAME mirror shape against `claudeteam.PresentFalseHint`; see Notes — flagged, not part of the confirmed four.)

## Approach

For each, replace the tautological check with one that has an INDEPENDENT oracle — a hand-written literal, a hand-known value, or a real failure path — while preserving every other real assertion in the same function.

1. **#1 no-op → real failure path with a deterministic fixture.** Change the `t.Logf` at `:384` to `t.Errorf`, asserting `activeSessionFile(dir) == teammatePath`. The oracle `teammatePath` is the hand-set path the test wrote itself (`:363`), not a value computed by the code under test. Because the assertion now depends on the teammate file being the NEWEST (both fixture files carry assistant entries, so `activeSessionFile`'s tiebreak is mod-time recency, `pty_session_test.go:219`), pin the ordering deterministically with `os.Chtimes` so the teammate is provably newer than the FO file — otherwise coarse-resolution filesystem timestamps could tie and flake the run. Smallest change is `t.Logf`→`t.Errorf`, but that alone is flake-prone; the `os.Chtimes` pin is required for a pristine, deterministic assertion.
2. **#2 mirror → hand-built literal.** Replace the `stampID` call at `:143` with a literal reconstruction of `newBody` (`native_new_test.go:12`) with `id: `+`want` spliced immediately before the closing `---` of the frontmatter — the position `stampID` inserts at, hand-written and human-checkable, NOT by calling `stampID`. `want` remains the minted id from `--next-id` (an independent value). The spelling is the implementer's choice (a full literal or an anchored `strings.Replace(newBody, "source: roadmap\n---", "source: roadmap\nid: "+want+"\n---", 1)`) provided it does not call `stampID` or replicate its line-scan.
3. **#3 mirror → hand-built literal.** Same treatment for the non-slug branch at `:552`: reconstruct `body` (`zz_independent_parity_test.go:480`) with `id: `+`expectedID` spliced before the closing `---`, independent of `stampID`. The slug branch (`:538`), the id-line `Contains` check (`:546`), and the post-create `--validate` (`:558`) are untouched.
4. **#4 mirror → literal hint.** Replace `teamStateNeutralHint` at `:112` with the literal `"no active team runtime detected"` (the constant's current value, hand-copied as the oracle). The confinement assertions (`:91`,`:95`), the present-true/absent blocks (`:100`,`:106`), and the negative Claude-advice check (`:118`) are untouched.

## Acceptance criteria

Each fixed test is a live regression guard: with the fix in place, applying the mutation that originally left it green turns it RED (non-zero exit), and reverting returns it GREEN. The independent baseline that can move the wrong way is each test's pass/fail under its exposing mutation — pre-fix it stayed GREEN (the tautology); the finished value is that it goes RED.

- **AC-1 (value — #1 catches a resolver recency regression).** After the fix, reversing `activeSessionFile`'s recency tiebreak (`pty_session_test.go:219`, `cands[i].mod > cands[j].mod` → `<`) makes the subtest FAIL (it returns the older FO file, so `active != teammatePath`); reverting restores GREEN. Pre-fix baseline that moved wrong: the same reversal leaves the subtest GREEN. Tested by: `go test ./internal/ensigncycle -run 'TestFOSessionPinning/activeSessionFile_would_flip_to_teammate' -count=1` before and after the one-line source mutation.
- **AC-2 (value — #2 catches a stampID splice regression).** After the fix, changing `stampID`'s insert splice (`new.go:131`) to place the id line after the closing `---` instead of before makes `TestNewFolderForm` FAIL on the byte-identity check; reverting restores GREEN. Pre-fix baseline: the same splice mutation leaves it GREEN. Tested by: `go test ./internal/status -run TestNewFolderForm -count=1` before and after the splice mutation.
- **AC-3 (value — #3 catches the same stampID regression on the parity path).** The single `stampID` splice mutation from AC-2 also makes `TestIndNewAtomicCreate` FAIL after the fix (it stayed GREEN before). Tested by: `go test ./internal/status -run TestIndNewAtomicCreate -count=1` under the same mutation — one mutation, two independently-oracle'd tests both go RED.
- **AC-4 (value — #4 catches an emptied hint constant).** After the fix, emptying `teamStateNeutralHint` (`boot.go:23` → `""`) makes `TestBootTeamStateProbeConfinement` FAIL on the nil-probe block (render emits `hint: `, literal oracle still expects `hint: no active team runtime detected`); reverting restores GREEN. Pre-fix baseline: emptying the constant left it GREEN. Tested by: `go test ./internal/status -run TestBootTeamStateProbeConfinement -count=1` before and after emptying the constant.
- **AC-5 (no lost coverage).** With no mutation applied, `go test ./internal/ensigncycle ./internal/status -count=1` is GREEN — the four fixed tests still pass on correct code and every sibling assertion in each function still runs. Tested by: the full package run, plus a diff review confirming only the four cited assertions changed (no sibling assertion removed).

## Implementation dispatch (smallest-sufficient-mechanism gate)

Rung picked: **one worker** — a single implementation-stage ensign in one worktree applies all four edits and runs the five mutation proofs plus the package suite. Justification (one line): the four edits are small and deterministic and share ONE verification pass (`go test` + the mutation RED/GREEN checks) with no cross-file collision, so there is no genuine fan-out, no isolation need, and no parallel-mutation hazard that would justify climbing to a fan-out sub-workflow or parallel workers; and dropping to FO-in-house is below the floor because the deliverable's proof obligation is per-test mutation RED/GREEN verification, which belongs in a worktree-backed stage rather than a gate-time inline edit (the redundant-re-verification direction the gate also forbids).

## Spike determination

**No spike needed.** The proof rests on proven, in-use mechanisms: applying a one-line source mutation and observing `go test` flip RED→GREEN is standard Go unit testing already used throughout these packages; there is no parser round-trip, on-disk format, runtime handoff, or new tool flag to validate. The splice-offset and emptied-constant mutations were already exercised by the audit that found the four defects (each confirmed green-under-mutation), so the exposing mutations are known to compile and run. The implementation's first action seeds directly from these: apply each fix, then re-run its named mutation to confirm the now-real assertion goes RED (the TDD-seed).

## Test plan

- **AC-1–AC-4 — mutation-kill, per test.** For each, apply the named source mutation, run the scoped `go test -run … -count=1`, confirm FAIL; revert, confirm PASS. Fixture/unit-level Go tests (`go test`), no CLI-behavior or live-workflow test needed — the claim is test-internal (does the assertion have teeth), provable by exercising the test under a controlled source mutation. Cost/complexity: low — five single-token/single-line source edits and ten `go test -run` invocations, all sub-second; no CI dependency, no worktree isolation beyond the implementation stage's own.
- **AC-5 — no-lost-coverage.** `go test ./internal/ensigncycle ./internal/status -count=1` green with no mutation, plus a diff review that only the four cited assertions (and #1's `os.Chtimes` fixture pin) changed. Cost: low.
- **Not needed:** live workflow smoke, golden-fixture regeneration, or host-lane runs — no user-visible surface, CLI output, or banner changes.

## Docs

No documentation change owed: this edits test files only — no CLI output, command surface, startup banner, or docs-site-described behavior changes. No doc diff.

## Notes

- Sibling mirror at `boot_probe_parity_test.go:106` (`wantAbsent := "…hint: " + claudeteam.PresentFalseHint`) is the SAME defect shape as finding #4 but was NOT independently mutation-confirmed by the audit, so it is outside the confirmed-four scope. Recommend the implementer fold it into the same three-assertion-block edit for consistency — swap to the literal `"run TeamCreate before first team-mode dispatch (claude runtime supports it)"` (`claudeteam.go:22`), whose exposing mutation is emptying `claudeteam.PresentFalseHint`. Flagged for the gate to confirm rather than silently expanding the AC set.
- The `activeSessionFile`/`foRootSessionID`/`sessionFileByID` test-local ports duplicate logic that lives in the pty/gym integration (reference-only copies). Whether the port pattern itself is worth revisiting is a separate, broader concern than the confirmed-four no-op assertion and is deliberately NOT in scope here.

## Stage Report: ideation

- DONE: Re-confirm each of the four cited file:test:pattern findings still matches current HEAD (the repo has moved since the audit) before designing fixes
  Read all four at HEAD b7f331d5 and confirmed each defect: #1 pty_session_test.go:384 failure branch is `t.Logf` only (no t.Error/Fatal/Fail in the subtest); #2 native_new_test.go:143 and #3 zz_independent_parity_test.go:552 both compute `want` via `stampID` (mirror); #4 boot_probe_parity_test.go:112 builds `wantNil` from the `teamStateNeutralHint` constant (mirror). Baseline suites GREEN today (`go test` both packages ok).
- DONE: Write Problem/Approach/Acceptance criteria/Test plan per docs/dev's ideation bar, with one AC per test pinning its fix to the SAME mutation that originally exposed it as a live-again regression guard
  AC-1 pins #1 to the recency-tiebreak reversal (:219), AC-2/AC-3 pin #2/#3 to the single `stampID` splice-offset mutation (new.go:131), AC-4 pins #4 to emptying `teamStateNeutralHint` (boot.go:23). Each AC measures the value end-state (test flips RED under its exposing mutation, GREEN reverted) against the pre-fix baseline that moved the wrong way (stayed GREEN). AC-5 guards no-lost-coverage. Added #1's `os.Chtimes` fixture pin to remove a coarse-timestamp flake risk the raw t.Errorf would introduce.
- DONE: Apply zm's smallest-sufficient-mechanism gate to the implementation-dispatch call itself (four small independent test edits across two packages) and record the one-line justification for whatever rung you pick
  Recorded in "Implementation dispatch" section: rung = one worker (single implementation ensign, one worktree); no fan-out/isolation/parallel-mutation hazard justifies climbing to a sub-workflow/parallel workers, and FO-in-house is below the floor because the deliverable's per-test mutation RED/GREEN proof belongs in a worktree-backed stage (the redundant-re-verification direction the gate also forbids).

### Summary

Re-confirmed all four tautological tests at HEAD b7f331d5 and designed independent-oracle fixes: #1 becomes a real failure path (t.Errorf) with an os.Chtimes-pinned fixture so the newest-file flip is deterministic; #2/#3 replace `stampID` mirror computations with hand-built literal expected bytes; #4 replaces the `teamStateNeutralHint` mirror with the literal hint string. Each AC pins its test to the exact mutation that originally left it green as a live regression guard (one shared `stampID` splice mutation covers #2 and #3). Flagged the sibling mirror at boot_probe_parity_test.go:106 as an out-of-confirmed-scope consistency fix for the gate to rule on, not a silent AC expansion; the implementation dispatch is gated to a single worker.

## Stage Report: implementation

- DONE: Apply the four CONFIRMED fixes per ## Approach, each with an independent oracle, TDD-style (apply fix -> run its named mutation -> confirm RED -> revert -> confirm GREEN). AC-1..AC-4 each go RED under their exposing mutation and GREEN reverted.
  Commit 9de3e058 (branch spacedock-ensign/tautological-test-fixes), 4 test files, +24/-7. Per-test mutation proof:
  AC-1 #1 (pty_session_test.go): asserts `activeSessionFile(dir) == teammatePath` via t.Errorf, with an os.Chtimes pin making the teammate provably newest. Falsifier: reversing the recency tiebreak (`:219` `>`→`<`) returns the older FO file → RED; reverted → GREEN. Pre-fix the same reversal stayed GREEN (the subtest only t.Logf'd — the tautology).
  AC-2 #2 / AC-3 #3 (native_new_test.go, zz_independent_parity_test.go): byte-identity `wantBytes`/`wantSeed` rebuilt via anchored `strings.Replace` (NOT stampID). Named mutation (stampID splice `new.go:131` to insert after the closing `---`) makes both RED on the byte-identity check (native_new_test.go:148, zz_independent_parity_test.go:556); reverted → GREEN.
  AC-4 #4 (boot_probe_parity_test.go): `wantNil` hint is the literal `"no active team runtime detected"`. Falsifier: emptying `teamStateNeutralHint` (`boot.go:23` → `""`) → RED at boot_probe_parity_test.go:116; reverted → GREEN. Pre-fix emptying it stayed GREEN (mirror).
- DONE: GATE RULING — EXCLUDE the :106 PresentFalseHint sibling; do NOT touch boot_probe_parity_test.go:106; stay strictly to the confirmed-four.
  `git diff` touches only `wantNil` (:112); no change references `PresentFalseHint` or `wantAbsent`. Confirmed via grep over the diff (no hits).
- DONE: AC-5 no-lost-coverage — full suite GREEN with no mutation; diff review confirms only the four cited assertions (+ #1's os.Chtimes pin) changed; gofmt clean; -race green. Base is origin/main (ca136f83).
  `go test ./internal/ensigncycle ./internal/status -count=1` → both ok (9.0s / 29.5s). `gofmt -l` on the 4 files: empty. `go test -race ./internal/ensigncycle ./internal/status` → both ok (11.0s / 35.4s). `git diff --stat`: 4 test files only, no production file changed (all mutations reverted). Base confirmed `git merge-base HEAD origin/main` = ca136f83.

### Summary

Gave all four tautological tests independent oracles and mutation-proved each (RED under its exposing mutation, GREEN reverted; pre-fix each stayed GREEN under the same mutation — the tautology). NOTABLE FINDING on #2/#3: the ideation's named mutation (id after the closing `---`) pushes the id out of the frontmatter, so the sibling `--validate` assertion catches it independently — it does NOT isolate the byte-identity mirror tautology (pre-fix under that mutation the tests fail at `--validate`, not the byte-identity check). I additionally used a within-frontmatter splice (id after the opening `---`, id stays in the frontmatter so `--validate` passes) as the clean isolator: pre-fix it left both mirror tests GREEN, post-fix it makes both RED on the byte-identity check ALONE. That confirms the literal-oracle fix genuinely added teeth the mirror lacked, beyond what `--validate` already guarded. The `:106` PresentFalseHint sibling was left untouched per the gate ruling. No production code changed.

## Stage Report: validation

- DONE: Verify AC-1 (activeSessionFile recency-tiebreak reversal reds #1) with reproduced mutation-kill.
  Post-fix + mutation (`pty_session_test.go:219` `cands[i].mod > cands[j].mod` → `<`) → RED at `:394` (returned the older FO file); reverted → GREEN. Pre-fix (ca136f83) + same mutation → GREEN (subtest only `t.Logf`'d — tautology). The fix's `t.Errorf` against the hand-set `teammatePath` is what fails.
- DONE: Verify AC-2/AC-3 (stampID splice regression reds #2 and #3) with reproduced mutation-kill.
  Post-fix + within-frontmatter splice (`new.go:131` `lines[:fmEnd:fmEnd]` → `lines[:fmStart+1:fmStart+1]`) → RED at `native_new_test.go:148` and `zz_independent_parity_test.go:556` on the byte-identity check ALONE (slug-flat branch correctly PASSES); reverted → GREEN.
- DONE: CRITICAL cross-check on #2/#3 — isolate the byte-identity tautology, not the id-out-of-frontmatter case.
  Reproduced the implementer's stronger proof: the within-frontmatter splice (id stays valid) is NOT caught by `--validate` — the parity test uses `t.Errorf` and ran on to its `:561` `--validate` with NO `:563` failure. Pre-fix mirror + that splice → fully GREEN (tautology); post-fix → RED on byte-identity alone. Also reproduced the ideation's NAMED mutation (id after closing `---`, `lines[:fmEnd+1:fmEnd+1]`): pre-fix mirror → RED but at `--validate` ("missing required id"), NOT byte-identity — confirming it never isolates the tautology. The tautology IS genuinely killed.
- DONE: Verify AC-4 (emptying teamStateNeutralHint reds #4) with reproduced mutation-kill.
  Post-fix + emptied constant (`boot.go:23` → `""`) → RED at `boot_probe_parity_test.go:116` (render emits `hint: `, literal oracle expects `hint: no active team runtime detected`); reverted → GREEN. Pre-fix mirror + same → GREEN (both sides emptied — tautology).
- DONE: Verify AC-5 (full suite GREEN, no lost coverage, diff review).
  `go test ./internal/ensigncycle ./internal/status -count=1` → both ok, exit 0 (9.3s / 29.8s). Diff vs origin/main = exactly the four cited assertions + #1's `os.Chtimes` pin and `time` import; 4 test files (+24/-7); gofmt-clean; no sibling assertion removed.
- DONE: Confirm :106 PresentFalseHint sibling NOT touched and NO production code changed (test-file-only).
  `git diff origin/main...HEAD` has zero hits for `PresentFalseHint`/`wantAbsent`; `git diff --name-only` shows only `*_test.go` (no production file). Commit under review: 9de3e058.
- DONE: Semantic adversarial pass over each fixed assertion.
  Each oracle is independent of the code under test: #1 asserts against the hand-written `teammatePath`; #2/#3 rebuild expected bytes via an anchored `strings.Replace` on `"source: roadmap\n---"` (fixed textual anchor, not a `stampID` line-scan) with the id from `--next-id`, pinning the id-line POSITION which the mirror could not; #4 is a hand-copied literal of the constant. No residual "passes while behavior wrong" path found — each mutation-kill above is the falsifier.
- DONE: Detached adversarial audit — N/A (test-file-only, low blast radius; ACs' oracle is external RED/GREEN-under-mutation).
  Per the stage's audit carve-out and the assignment: no high-stakes production surface changed. The AC-provenance trigger is satisfied — every AC's proof is a reproduced command RED/GREEN, not a string match.

### Summary

Reproduced every AC's mutation-kill from a clean worktree at commit 9de3e058 (base origin/main ca136f83): each fixed test goes RED under its exposing mutation and GREEN reverted, and pre-fix each stayed GREEN under the same mutation — proving the tautology was real and is now closed. The #2/#3 critical cross-check holds: the byte-identity oracle catches a clean within-frontmatter `stampID` splice that `--validate` does not, while the ideation's named (id-out-of-frontmatter) mutation is caught only by `--validate` — so the fix adds genuine teeth. No production code changed, the `:106` sibling is untouched, gofmt-clean, full suite GREEN. No material findings; no deferred risks. Recommendation: **PASSED**.
