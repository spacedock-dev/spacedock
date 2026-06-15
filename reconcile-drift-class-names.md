---
title: Reconcile drift output carries descriptive class names instead of A-E letters
status: validation
source: captain (2026-06-14, this session) — `spacedock dispatch reconcile` emits opaque single-letter drift classes (A-E) when the descriptive names already exist in the FO dispatch contract event-loop step-0 mapping (A=lingering, B=superseded, C=un-advanced-pr, D=stale-branch, E=local-main-drift). The letter is indirection over a name the `reason` field already states in English. Decided this session — emit the descriptive names and drop the letters; a string enum is as stable to branch on, with no machine-stability cost.
started: 2026-06-15T05:28:02Z
completed: 2026-06-15T15:45:51Z
verdict: PASSED
score: 0.3
worktree:
issue:
id: pd7fqh4f8yzf9dacbbbamfg7
mod-block:
pr: "#382"
---

The `spacedock dispatch reconcile` output tags each drift entry with a bare letter (`"class": "A"` .. `"E"`). The descriptive names already exist — the FO dispatch contract event-loop step-0 spells every class out parenthetically — but the helper discards them and emits the letter, so a reader (FO or human) must carry the A-E mapping to read the output. Decided this session: emit the descriptive name and drop the letter.

## Problem

`spacedock dispatch reconcile` emits `drift[].class` as a single letter A-E whose meaning lives only in the FO dispatch contract event-loop step-0 mapping:
- A = lingering (roster member, no live work)
- B = superseded
- C = un-advanced-pr
- D = stale-branch
- E = local-main-drift

The `reason` field already states the condition in English, so the letter adds an indirection the reader must resolve against the contract. The FO one-line drift summary (`A={N} B={N} C={N} D={N} E={N}`) is equally opaque.

## Proposed approach

Behavior decided this session — ideation formalizes the ACs, test plan, and doc-diff, not the design:
- Emit the descriptive name as `drift[].class`: `lingering` | `superseded` | `un-advanced-pr` | `stale-branch` | `local-main-drift`. Remove the letters.
- Update the FO dispatch contract event-loop step-0 — the per-class action mapping and the one-line summary format — to reference the names; the FO branches on the string name.
- A string enum is as stable to branch on as a letter, so there is no machine-stability cost; the gain is self-documenting JSON and a readable summary line (`stale-branch=3` instead of `D=3`).

## Out of scope

The drift-class set itself (no new classes, no semantics change); the reconcile detection logic; any other `dispatch` subcommand's output. The `--include` flag is in scope (AC-3): it is the same `reconcile` surface's input vocabulary and shares the identical opaque-letter problem, so its accepted tokens rename in lockstep with the output — keeping input and output vocabulary equal. If the gate prefers to scope `--include` out and leave it on A-E, drop AC-3; the output rename (AC-1/AC-2) stands independently.

## Acceptance criteria

**AC-1 — `dispatch reconcile` output carries the descriptive class name, never a letter.** Each `drift[].class` in real reconcile output is one of `lingering`/`superseded`/`un-advanced-pr`/`stale-branch`/`local-main-drift`, and no `drift[].class` is a bare single A-E letter.
Verified by: a behavioral Go test that drives `Reconcile()` over the existing five-class fixture (`internal/dispatch/reconcile_test.go`) and asserts every emitted `class` is in the descriptive-name set and matches none of `{A,B,C,D,E}`. The existing `TestReconcileFiveClasses` / `TestReconcileFlipReclassifies` assertions (which group by `"A"`..`"E"`) are rewritten to the descriptive names — the same fixture, re-pinned — so the behavior is proven by running the helper, not by reading prose.

**AC-2 — The FO dispatch contract event-loop step-0 names exactly the helper's emitted class set, bound as two independent sources that red on drift in either direction.** The contract event-loop step-0 (in `skills/first-officer/references/claude-fo-dispatch.md`) and the helper (`internal/dispatch/reconcile.go`) enumerate the identical five-name set; neither side can rename, add, or drop a class without the other.
Verified by: a `contractlint` quarantine-package structural lint (the single sanctioned instruction-read path; the boundary guard exempts `internal/contractlint`). It extracts the class set from BOTH sources independently and asserts set equality: (a) helper side — `go/ast` scan of `reconcile.go` for the canonical declared class set (the `var` the rename introduces; see test plan); (b) contract side — a regex over the step-0 `"class":"…"` JSON-shape token, splitting its `|`-alternation into members. This is the same doc⟷code lock-test shape `scenario-testing-principles.md` already uses ("red on drift in either direction"), a DEDUP/cross-source structural check — NOT a prose-grep (it never asserts the doc contains a given word; it compares two extracted enumerations) and NOT a behavior substitute (AC-1 is the behavior test that runs the helper). Spiked: see "Spike" below.

**AC-3 — The `--include` flag accepts the descriptive class names, not A-E letters, with input and output vocabulary equal.** `dispatch reconcile --include stale-branch,local-main-drift` scopes the sweep to those two classes; an unknown token (including a bare `A`) is a usage error (exit 2) whose message names the descriptive-name set; the help line advertises the descriptive names.
Verified by: a behavioral Go test driving `runReconcile`/`Reconcile` with `--include` set to descriptive names (assert only those classes emit), with a bad token (assert exit 2 and the error names the descriptive set), reusing the existing `TestReconcileIncludeScope` / `TestReconcileUsageError` surface re-pinned to the new vocabulary.

## Spike

Spiked the AC-2 dual-extraction mechanism against the CURRENT A-E state (proving the mechanism, not the post-rename values), in a throwaway `internal/contractlint` test:
- Contract side: regex `"class":"([A-Za-z|\-]+)"` over `claude-fo-dispatch.md` yielded `[A B C D E]` from the step-0 JSON-shape token. The `\-` in the char class already tolerates the hyphenated names (`un-advanced-pr`, `stale-branch`, `local-main-drift`).
- Helper side: `go/ast` scan of `reconcile.go` for `Class: "…"` composite-literal values yielded `[A B C D E]`.
- The two extracted sets were equal today, and the test compares them as sets — so it reds when either source renames a class without the other. Both extractions are structural (a delimited-token parse and an AST literal scan), neither is a prose-grep. Mechanism proven; no further unverified path remains.

## Test plan

Three behavioral surfaces plus one structural lint, all offline Go tests on the existing reconcile/contractlint test surface — no live workflow run needed (the claim is helper output bytes and a doc⟷code enumeration, both exercisable offline). Estimated cost: low — re-pins existing fixtures and adds one ~40-line contractlint test; the reconcile fixture already builds a real git tree, so the descriptive-name assertions are a find/replace over assertion literals plus the new set-membership check.

Recommended helper-side shape (makes AC-2's "helper emitted enum" a single source rather than five scattered `Class:` literals): introduce one canonical declared set in `reconcile.go` — e.g.

    const (
        classLingering     = "lingering"
        classSuperseded    = "superseded"
        classUnadvancedPR  = "un-advanced-pr"
        classStaleBranch   = "stale-branch"
        classLocalMainDrift = "local-main-drift"
    )
    var driftClasses = []string{classLingering, classSuperseded, classUnadvancedPR, classStaleBranch, classLocalMainDrift}

`classA`..`classE` emit these constants as `Class:`; `parseInclude` validates against `driftClasses` (replacing the `{A,B,C,D,E}` map and the `--include` error string); `sortDrift`'s `a.Class < b.Class` lexical ordering still holds (it sorts whatever the strings are — deterministic regardless of the rename). The AC-2 lint then reads `driftClasses` (one source) via AST. Acceptable minimal fallback: leave the five `Class:` literals scattered and have the lint scan all of them; the canonical var is preferred because it also removes the A-E coupling from `parseInclude` and gives the lint one symbol to read.

Tests:
1. **AC-1 behavioral** (`reconcile_test.go`, re-pinned): `TestReconcileFiveClasses` and `TestReconcileFlipReclassifies` assert the descriptive class names (the class-iteration list `{"A".."E"}` → the descriptive set; `d.Class == "A"` → `d.Class == "lingering"`, etc.) plus a new assertion that no emitted `class` matches `^[A-E]$`.
2. **AC-2 structural lint** (new `internal/contractlint/reconcile_class_binding_test.go`): the spiked dual-extraction, asserting set equality between the helper `driftClasses` (AST) and the contract step-0 JSON-shape token (regex). Empty-set guard on both sides so it cannot pass vacuously.
3. **AC-3 behavioral** (`reconcile_test.go`, re-pinned): `TestReconcileIncludeScope` drives `--include` with descriptive names; `TestReconcileUsageError`'s bad-`--include` case asserts exit 2 and the error message names the descriptive set.

Run all four offline: `go test ./internal/dispatch/ ./internal/contractlint/ -run 'Reconcile|ClassBinding' -count=1`.

## Documentation changes

User-visible surfaces the rename touches: the `dispatch reconcile` JSON output (`drift[].class`), the `--include` help line, and the FO dispatch contract event-loop step-0. The contract step-0 doc-diff (the assignment's required exact wording, in `skills/first-officer/references/claude-fo-dispatch.md`):

**(1) Step-0 JSON-shape line — the lede sentence's `Stdout:` clause:**

Before:
> Stdout: `{"command":"reconcile","team_name":…,"drift":[{"class":"A|B|C|D|E",…}]}`. Empty `drift[]` is green. Act per drift class:

After:
> Stdout: `{"command":"reconcile","team_name":…,"drift":[{"class":"lingering|superseded|un-advanced-pr|stale-branch|local-main-drift",…}]}`. Empty `drift[]` is green. Act per drift class:

**(2) Step-0 per-class action mapping — the five bullets:**

Before:
> - **A (lingering)** / **B (superseded)** → `SendMessage({"type":"shutdown_request"})` to `name`; drop from session memory.
> - **C (un-advanced PR)** → enter Merge-and-Cleanup for the named slug.
> - **D (stale branch)** → only when `drift.owned == true`: `git -C {worktree} pull --rebase origin {drift.trunk}`; halt on conflict per the rebase-conflict halt rule. When `drift.owned` is false the item is report-only — surface it to the captain; do NOT rebase a worktree the current session does not own.
> - **E (local main drift)** → behind only (`drift.behind > 0 && drift.ahead == 0`): `git -C {repo} fetch origin {drift.trunk} && git -C {repo} merge --ff-only origin/{drift.trunk} && cd {repo} && go build -o spacedock ./cmd/spacedock`. Ahead/unpushed or diverged (`drift.ahead > 0`): report-only — surface `drift.reason` to the captain and NEVER `reset --hard`; the captain decides push vs. manual reconcile.

After:
> - **lingering** / **superseded** → `SendMessage({"type":"shutdown_request"})` to `name`; drop from session memory.
> - **un-advanced-pr** → enter Merge-and-Cleanup for the named slug.
> - **stale-branch** → only when `drift.owned == true`: `git -C {worktree} pull --rebase origin {drift.trunk}`; halt on conflict per the rebase-conflict halt rule. When `drift.owned` is false the item is report-only — surface it to the captain; do NOT rebase a worktree the current session does not own.
> - **local-main-drift** → behind only (`drift.behind > 0 && drift.ahead == 0`): `git -C {repo} fetch origin {drift.trunk} && git -C {repo} merge --ff-only origin/{drift.trunk} && cd {repo} && go build -o spacedock ./cmd/spacedock`. Ahead/unpushed or diverged (`drift.ahead > 0`): report-only — surface `drift.reason` to the captain and NEVER `reset --hard`; the captain decides push vs. manual reconcile.

**(3) Step-0 one-line drift summary — the closing sentence:**

Before:
> On drift, report one line: `reconcile: {N} entries: A={N_A} B={N_B} C={N_C} D={N_D} E={N_E} — acting`.

After:
> On drift, report one line: `reconcile: {N} entries: lingering={N} superseded={N} un-advanced-pr={N} stale-branch={N} local-main-drift={N} — acting`.

**(4) CLI help line** (`internal/dispatch/dispatch.go`, the reconcile usage line):

Before:
> `  spacedock dispatch reconcile --workflow-dir DIR [--team-name NAME] [--repo-root DIR] [--include A,B,C,D,E]`

After:
> `  spacedock dispatch reconcile --workflow-dir DIR [--team-name NAME] [--repo-root DIR] [--include lingering,superseded,un-advanced-pr,stale-branch,local-main-drift]`

## Stage Report: ideation

- DONE: The design is DECIDED (emit lingering/superseded/un-advanced-pr/stale-branch/local-main-drift as the class value, drop the A-E letters) — formalize it; do not re-open the design or weigh alternatives.
  Formalized AC-1/AC-2/AC-3 and a test plan from the decided design; did not re-open or weigh alternatives. Located the three emitting helper sites (classA..E `Class:` literals) and `parseInclude`/`sortDrift` ripples in `internal/dispatch/reconcile.go`.
- DONE: Behavior-first ACs plus a test plan that proves the rename over real `dispatch reconcile` output (the existing reconcile test surface), and an independent-source check binding the contract event-loop step-0 class set to the helper emitted enum — never a prose-grep over the contract.
  AC-1/AC-3 are behavioral over `internal/dispatch/reconcile_test.go` (re-pinning `TestReconcileFiveClasses`/`FlipReclassifies`/`IncludeScope`/`UsageError`). AC-2 is a `contractlint` quarantine lint binding the contract step-0 `"class":"…"` JSON-shape token (regex) to the helper `driftClasses` (AST) as two independent sources — the sanctioned doc⟷code lock-test shape, not a prose-grep. Spiked both extractions against the live A-E state (both yielded `[A B C D E]`, compared as sets) — mechanism proven.
- DONE: Record the exact contract doc-diff: the FO dispatch contract event-loop step-0 per-class action mapping and the one-line drift summary, before/after wording.
  Recorded under "## Documentation changes": (1) step-0 `Stdout:` JSON-shape line, (2) the five per-class action bullets, (3) the one-line drift summary, (4) the CLI `--include` help line in `internal/dispatch/dispatch.go:274` — each before/after, verbatim against `skills/first-officer/references/claude-fo-dispatch.md` lines 186/187-190/192.

### Summary

Formalized the decided rename: three behavior-first ACs (output `drift[].class`, contract⟷helper enum binding, `--include` input vocabulary) plus a test plan over the existing reconcile/contractlint surfaces. Two open decisions for the gate: (a) AC-3 puts `--include` in scope — renaming the input vocabulary in lockstep so input and output stay equal; the gate can scope it out and keep `--include` on A-E without touching AC-1/AC-2. (b) The test plan recommends a single canonical `driftClasses` var in `reconcile.go` (vs. five scattered `Class:` literals) so AC-2's lint reads one symbol; this also removes the A-E coupling from `parseInclude`. Baseline `go test ./internal/dispatch/ -run TestReconcile` is green; the AC-2 dual-extraction mechanism was spiked and proven.

## Stage Report: implementation

- DONE: Rename the emitted `drift[].class` to the descriptive names via a single canonical `driftClasses` var in reconcile.go (classA..E emit the constants), and rename `parseInclude`'s accepted vocabulary + the usage/help line in lockstep (AC-3) so a bad token is an exit-2 usage error naming the descriptive set.
  Added `const classLingering..classLocalMainDrift` + `var driftClasses` in reconcile.go; classA..E `Class:` literals and the `Reconcile` include gating now reference the constants; `parseInclude` validates against `driftClasses` and its error names the descriptive set. Built binary confirms `--include A` → exit 2, error `(expected lingering,superseded,un-advanced-pr,stale-branch,local-main-drift)`. Commit ef8ed3f5.
- DONE: Apply the recorded contract doc-diff to the FO dispatch contract event-loop step-0 in skills/first-officer/references/claude-fo-dispatch.md (the Stdout JSON-shape clause, the five per-class action bullets, the one-line drift summary) and the CLI --include help line — match by content, not blind line numbers.
  All three step-0 edits applied verbatim (JSON-shape token, five bullets, summary line); CLI help line in dispatch.go updated. Also updated two same-paragraph/backstop prose references that named the classes by letter ("roster-derived classes (A/B/C)"/"(D/E)" → descriptive; "Class A/Class B" backstop → "lingering/superseded class") for internal consistency — see Summary note, this was beyond the three recorded wordings.
- DONE: Add the three tests: AC-1 behavioral (re-pin reconcile_test.go to the descriptive names + assert no emitted class is a bare A-E letter), AC-2 the contractlint dual-extraction lint binding the contract step-0 class set to the helper `driftClasses` as two independent sources, AC-3 the --include vocabulary + bad-token exit-2 behavior.
  AC-1: reconcile_test.go re-pinned + new `isBareLetterClass` assertion over real output. AC-2: new `internal/contractlint/reconcile_class_binding_test.go` (AST var-resolve + regex token), verified to RED on drift in both directions (drop-a-class mutation on each source) and green at rest. AC-3: new `TestReconcileIncludeVocabulary` + `TestReconcileUsageError` re-pinned. Sibling reconcile_*_test.go files (session/trunk/de_safety/decompose/namecap/e) also re-pinned to keep the package green. `go test ./...` all green; `go vet`/`gofmt` clean.

### Summary

Renamed the reconcile drift-class vocabulary from A-E letters to descriptive names through one canonical `driftClasses` var (the AC-2 lint's single helper-side symbol), with `parseInclude` and the contract step-0 doc-diff in lockstep. The AC-2 dual-extraction lint was proven to red on drift in either direction. Scope note: the recorded doc-diff covered only the step-0 JSON-shape clause, the five bullets, and the summary line — but applying only those left "roster-derived classes (A/B/C)"/"(D/E)" in the same paragraph's lede and "Class A/Class B" in the Claude backstop as the only A-E letters left in the entire surface, so I updated those parentheticals/labels to the descriptive names for internal consistency. This is a mechanical same-paragraph consistency fix, not a design change; flagging it since it exceeds the three recorded wordings. Validation owes a re-run of `go test ./internal/dispatch/ ./internal/contractlint/ -run 'Reconcile|ClassBinding|ParseInclude' -count=1`.

## Feedback Cycles

### Cycle 1 — validation (detached adversarial audit)

Audit: AC-2 dodge-construction on a throwaway detached checkout of `spacedock-ensign/reconcile-drift-class-names` (never the implementation worktree). Finding: the AC-2 contractlint binding is NARROWER than AC-2's stated claim.

- **AC-2 claims** "the contract event-loop step-0 … neither side can rename, add, or drop a class without the other," and the entity's "Documentation changes" + AC-2 text name THREE step-0 surfaces as bound: the JSON-shape `"class":"…"` token, the five per-class action bullets, and the one-line drift summary.
- **The lint binds only ONE of the three.** `contractClassToken = "class":"([A-Za-z|\-]+)"` with `FindSubmatch` reads only the FIRST `"class":"…"` JSON-shape token. The per-class **action bullets** (`- **stale-branch** → …`) and the **one-line summary** (`stale-branch={N} …`) are never extracted.
- **Demonstrated dodges (both passed AC-2 AND the AC-1/AC-3 behavioral tests):**
  - Adversarial edit A: renamed the action bullet `- **stale-branch**` → `- **stale-branch-TYPO**` (JSON token untouched) → `TestReconcileClassBinding` PASSED (should RED).
  - Adversarial edit B: renamed the summary token `stale-branch={N}` → `stale-branch-OOPS={N}` → `TestReconcileClassBinding` PASSED (should RED).
- **Why it matters:** the FO branches on `drift.class` (JSON), then looks up the matching action bullet to know what to DO. A drifted bullet/summary leaves the FO unable to resolve the action for a real class — and nothing reds. This is the test-strength hole class the validation stage's adversarial audit exists to catch (cf. #262's two `contract_gate_test.go` holes named in the stage def).
- **Bounded fix:** extend the contract-side extractor so all THREE step-0 class-naming surfaces feed the set-equality (or assert each emitted class appears as both a bullet and a summary token). Keep it structural (parse the bullet/summary tokens), not a prose-grep — the helper `driftClasses` stays the single source; the contract side must enumerate from every step-0 site the AC claims to bind, not just the JSON token. The helper-side extractor (AST `driftClasses` var) is sound and needs no change.

Note — what the audit did NOT find (these are clean): residual bare A-E class letters in the output/contract surface (`reconcile.go`, `dispatch.go`, `claude-fo-dispatch.md` step-0) — none remain; the only A/B/C-D/E mentions left are Go-function shorthand in code comments and a historical debrief, both out of scope. The helper-emitted enum, the `driftClasses` var, the JSON-shape token, the bullets, and the summary all agree at rest (5 classes, identical members).

## Stage Report: validation

- DONE: Reproduce the AC evidence, do not trust the report: run full `go test ./...` green, and independently MUTATE one source of the AC-2 dual-extraction lint and confirm the contractlint test goes RED in both directions, then restore.
  `go test ./...` all green (uncached contractlint 0.360s, all packages ok). MUTATED helper side (`classStaleBranch` const → `stale-branch-RENAMED`) → `TestReconcileClassBinding` RED with correct set-mismatch diff; restored. MUTATED contract side (step-0 JSON token `stale-branch` → `stale-branch-RENAMED`) → RED; restored. Worktree `git status` clean + `git diff` empty after restore; lint green again. Binding catches drift in BOTH directions — proven, not trusted.
- FAILED: Detached adversarial audit on a THROWAWAY checkout: construct a class-vocabulary drift that DODGES the AC-2 lint; grep the whole reconcile + contract surface for any residual bare A-E class letter.
  Ran on a throwaway detached worktree (`/tmp/reconcile-audit-*`, torn down). FOUND TWO dodges that pass AC-2 AND the behavioral tests: an action-bullet rename and a summary-line rename both slip past the lint, which only binds the first `"class":"…"` JSON token. See `### Feedback Cycles` Cycle 1. Residual-letter grep half is clean: no bare A-E class letter remains in `reconcile.go`/`dispatch.go`/`claude-fo-dispatch.md` step-0 (only Go-func-shorthand comments and a historical debrief, both out of scope).
- DONE: AC cross-check: AC-1 (output names + no bare-letter assertion), AC-2 (the binding lint), AC-3 (--include descriptive vocabulary + bad-token exit-2) each have reproduced evidence; the contract step-0 doc-diff matches the emitted `driftClasses` enum; PASSED/REJECTED recommendation.
  AC-1 reproduced: `TestReconcileFiveClasses`/`FlipReclassifies` drive real `Reconcile()` over a real git tree; `isBareLetterClass` asserts no emitted class is bare A-E — PASS. AC-3 reproduced: `TestReconcileIncludeVocabulary` accepts each descriptive name + rejects bare letters; `TestReconcileUsageError` confirms `--include A` → exit 2 naming the descriptive set — PASS. AC-2 lint passes at rest and reds on JSON-token drift, BUT under-delivers its own stated scope (bullets + summary unbound). Step-0 doc-diff matches `driftClasses` at rest (JSON token, var, bullets, summary all 5, identical). Recommendation: REJECTED.

### Summary

The rename behavior is correct and well-tested: AC-1 and AC-3 are genuine behavioral tests over real `Reconcile()` output and pass, no bare A-E letter survives in any user-visible surface, and the helper's emitted enum, the `driftClasses` var, and all contract step-0 sites agree at rest. The AC-2 binding works for the JSON-shape token (proven RED on drift in both directions). REJECTED on one material gap surfaced by the detached adversarial audit: AC-2 claims to bind all of step-0 ("neither side can rename a class"), and the entity explicitly names the action bullets and one-line summary as bound surfaces, but the lint only reads the first `"class":"…"` JSON token — a class rename in an action bullet or the summary line passes AC-2 and all behavioral tests (two demonstrated dodges). The fix is bounded: extend the contract-side extractor to enumerate from the bullets and summary too, so every step-0 site the AC claims is actually bound. Helper-side extractor and all behavioral tests need no change.

### Feedback Cycles

#### Cycle 1 (validation REJECTED, 2026-06-15)

- **Stage:** validation → routed back to implementation (same worktree, reused implementation ensign).
- **Finding (detached adversarial audit, material):** the AC-2 contractlint binding under-delivers its own stated scope. The entity claims AC-2 binds the full contract step-0 class set ("neither side can rename a class"), and names the per-class action bullets and the one-line summary as bound surfaces, but the lint's contract-side extractor reads only the first `"class":"…"` JSON-shape token. Two demonstrated dodges pass AC-2 AND every behavioral test: a class rename in an action bullet, and a class rename in the summary line.
- **Not in scope of the gap:** the rename behavior itself is correct (AC-1/AC-3 are genuine behavioral tests over real `Reconcile()` output and pass; no bare A-E letter survives any user-visible surface; `go test ./...` green; the JSON-token binding reds on drift in both directions, reproduced by mutation).
- **Fix routed:** extend the AC-2 contract-side extractor to enumerate class names from all three step-0 surfaces (the JSON-shape token, the per-class action bullets, and the one-line summary), so every surface the AC claims is actually bound. Helper-side extractor and behavioral tests unchanged. Re-verify by reproducing both dodges and confirming each now REDs the lint.

## Stage Report: implementation (cycle 2)

- DONE: Extend the AC-2 contract-side extractor to bind all three step-0 surfaces (JSON-shape token, per-class action bullets, one-line drift summary) against the helper `driftClasses` set.
  Rewrote the contract side of `internal/contractlint/reconcile_class_binding_test.go`: `step0Block` bounds reads to the `0. **Reconcile sweep.**`→`1. **` slice; three structural extractors — `contractClassesFromToken` (`"class":"…"` split on `|`), `contractClassesFromBullets` (bold name(s) before each `   - … →` action arrow), `contractClassesFromSummary` (`name={N}` pairs) — each asserted equal to the helper set, with the failure message naming the offending surface. Helper-side AST extractor and the AC-1/AC-3 behavioral tests unchanged. Commit 3a216fe2.
- DONE: Re-verify by reproducing the audit's two dodges plus the prior drift mutations; each must RED; green at rest.
  Reproduced both dodges and confirmed RED: action-bullet rename (`un-advanced-pr`→`un-advanced-pr-X`) → RED naming "per-class action bullets"; summary-line rename (`stale-branch`→`stale-branch-X`) → RED naming "one-line drift summary". Prior mutations still RED: JSON-token rename → RED naming "JSON-shape token"; helper-side var drop (`classLocalMainDrift`) → RED. Each restored clean (`git diff` empty); lint green at rest. `go test ./internal/dispatch/ ./internal/contractlint/ -count=1` and full `go test ./...` green; gofmt/vet clean.

### Summary (cycle 2)

Closed the cycle-1 gap: the AC-2 contract-side extractor now binds all three step-0 class-naming surfaces — the JSON-shape token, the five per-class action bullets, and the one-line drift summary — each independently set-equal to the helper `driftClasses` enum, scoped to the bounded step-0 block. A class renamed/added/dropped in the helper or in ANY single surface now reds, with the error naming the offending surface. Both audit dodges (bullet-rename, summary-rename) were reproduced and confirmed to RED, alongside the prior JSON-token and helper-side drift mutations; green at rest, full suite green. Helper-side extractor and the AC-1/AC-3 behavioral tests were untouched, as scoped.

## Stage Report: validation (cycle 2)

- DONE: Reproduce the cycle-1 fix, do not trust the report: confirm BOTH audit dodges (a class rename in a per-class action bullet, and in the one-line summary) now RED the AC-2 lint naming the offending surface, AND the prior JSON-token-drift + helper-side-var-drop mutations still RED; full `go test ./...` green at rest.
  Reproduced all four mutations on the impl worktree, each restored clean: action-bullet rename (`un-advanced-pr`) → RED naming "per-class action bullets"; summary rename (`stale-branch={N}`) → RED naming "one-line drift summary"; JSON-token rename (`local-main-drift`) → RED naming "JSON-shape token"; helper var drop (`classLocalMainDrift`) → RED against all three surfaces. `go test ./...` green at rest; vet/gofmt clean.
- DONE: Detached adversarial audit on a THROWAWAY checkout (re-run against the now-three-surface binding): try to find a NEW dodge the extended extractor still misses; confirm the binding is a genuine independent-source set-equality, not a prose-grep; "refuted nothing material" is a valid recorded outcome.
  Ran on a throwaway detached worktree (`/tmp/reconcile-audit-cycle2-*`, torn down). Refuted nothing material. Two dodges pass the lint but are out of AC-2's binding-correctness scope: (1) renaming a class in the line-186 PROSE parentheticals (`roster-derived classes (…)` / `git/filesystem classes (…)`), and (2) a hypothetical SECOND contradictory `"class":"…"` example token in step-0 (the token extractor reads the first match). Neither breaks the FO action-resolution path (JSON class → action bullet → summary), which is fully bound; at rest only one JSON token exists and both parentheticals agree with the bound set. Combined-bullet rename (`superseded` inside `**lingering** / **superseded**`) and a malformed summary placeholder (`={count}`) both correctly RED — well-covered. Independence proven NOT a prose-grep: a coordinated `superseded`→`reclassed` rename across the helper const AND all three contract surfaces stays GREEN (compares extracted values, not a baked-in word), and diverging only one surface back REDs.
- DONE: AC cross-check: AC-2 now binds all three named step-0 surfaces with reproduced evidence; AC-1/AC-3 unaffected (no regression); the contract step-0 doc-diff still matches the emitted driftClasses enum; PASSED/REJECTED recommendation.
  AC-1 reproduced: `TestReconcileFiveClasses`/`FlipReclassifies` run real `Reconcile()` over a real git tree, assert every emitted class is a descriptive name, `isBareLetterClass` rejects any bare A-E — PASS. AC-3 reproduced: `TestReconcileIncludeVocabulary` + `TestReconcileUsageError` drive the CLI, `--include A` → exit 2 naming the descriptive set; help line carries descriptive vocab, no A-E survives — PASS. AC-2 binds JSON token + bullets + summary, reds on drift in any direction (four mutations, each naming its surface) — PASS. Step-0 doc-diff matches `driftClasses` at rest (JSON token, var, bullets, summary, help line all 5 identical members; no bare A-E class shape in step-0). Recommendation: PASSED.

### Summary (cycle 2)

PASSED. The cycle-1 gap is genuinely closed: the AC-2 contractlint binding now set-equals the helper `driftClasses` enum against all three step-0 class-naming surfaces — the JSON-shape token, the five per-class action bullets, and the one-line drift summary — each independently, with the failure message naming the offending surface. I reproduced all four drift mutations (both audit dodges plus the prior JSON-token and helper-var-drop) and confirmed each REDs; the binding is a genuine independent-source set-equality, not a prose-grep (proven by a coordinated rename that stays green while a single-surface divergence reds). The detached adversarial audit on a throwaway checkout refuted nothing material: the two dodges that slip the lint (line-186 prose parentheticals, a hypothetical second JSON example token) are doc-consistency nits outside AC-2's binding-correctness purpose — neither breaks the FO's JSON→bullet→summary action-resolution path, and both surfaces agree with the bound set at rest. AC-1/AC-3 remain genuine behavioral tests over real `Reconcile()` output and pass with no regression; full `go test ./...` is green.
