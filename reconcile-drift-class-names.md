---
title: Reconcile drift output carries descriptive class names instead of A-E letters
status: implementation
source: captain (2026-06-14, this session) — `spacedock dispatch reconcile` emits opaque single-letter drift classes (A-E) when the descriptive names already exist in the FO dispatch contract event-loop step-0 mapping (A=lingering, B=superseded, C=un-advanced-pr, D=stale-branch, E=local-main-drift). The letter is indirection over a name the `reason` field already states in English. Decided this session — emit the descriptive names and drop the letters; a string enum is as stable to branch on, with no machine-stability cost.
started: 2026-06-15T05:28:02Z
completed:
verdict:
score: 0.3
worktree: .worktrees/spacedock-ensign-reconcile-drift-class-names
issue:
id: pd7fqh4f8yzf9dacbbbamfg7
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
