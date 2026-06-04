---
id: 2adcrvj56b5camy1v70c4ncc
title: Refuse to close a task whose only "proof" is re-reading its own write-up
status: done
source: hx decomposition (A of 3) — captain 2026-06-01; staff review
score: "0.32"
started: 2026-06-02T21:14:43Z
completed: 2026-06-04T05:27:06Z
verdict: PASSED
worktree:
issue:
mod-block:
pr: "#277"
---

**What this is for (plain).** Today nothing stops a task from being marked finished when its "proof"
is just pointing back at its own description — no real test, no command, no actual check behind it.
This adds a guard in the `spacedock` tool that refuses to finalize such a task (with an override flag
for the rare deliberate exception), so every finished task carries evidence that can actually fail.

**Value to the user / FO.** A mechanical backstop under the human gate: the FO no longer has to catch
every empty or circular "proof" by eye — the tool blocks the close and the run-time validate flags it.
Finished work always has real evidence, so the workflow can't quietly rubber-stamp a task that proves
nothing.

This is part A of three from the now-superseded parent `deliverable-contract-hardening`
(id `hxs93wd0bjwhc3vsjwx1seew`) — the parent's full design, the spike (now re-run; see §P4 below), and
the staff review remain the reference. This child is the CODE half: the guard in `internal/status`.

## Problem

The workflow contract (`docs/dev/README.md` ideation Outputs) already says every AC's "Verified by"
must name something outside the task body that can fail — a test, a command's exit code/output, a
file the change produces, or the resulting on-disk state — and the validation stage already says
reject any AC whose only proof is a re-reading of the entity's own prose. That language is the GATE
discipline. But nothing in the binary enforces it: `spacedock status --set verdict=passed` will
happily walk a self-referentially-proven entity to `done`, and `spacedock status --validate` will
return `VALID`. The captain or FO has to catch it by eye every time.

The live cautionary tale is the archived `external-tracker-checkpoint/index.md` **AC-6** — its
`Verified by:` clause reads "this entity's **v1 DECISION** section states the decision and cites the
roadmap out-of-scope line" — proof by re-reading the entity. That AC can never fail, so it is not an
AC. The classifier spike in the parent (re-run against the current corpus, §P4) confirms this is the
sole live self-reference today.

This is the **dev workflow's** discipline. Non-dev workflows (a marketing workflow whose proof is a
published artifact, a research workflow whose proof is a metric) legitimately ship under different
verifiability rules. The guard must therefore be OPT-IN at the workflow level, defaulting OFF — not a
universal binary behavior dressed up as a "goodness lint".

## Proposed approach

A single body-parser/classifier in `internal/status/` that:

1. extracts each `**AC-N — …**` block from an entity body;
2. isolates the proof clause (text from the first `Verified by:` / `Oracle:` / `Proof:` marker
   onward, to the next `**AC-` header or section break);
3. strips quoted spans (`"…"`, `` `…` ``) — a quoted example of the antipattern is not the live
   proof;
4. flags the AC iff the cleaned proof clause matches a self-reference phrase AND contains no
   external-proof token.

That same classifier is consulted by two surfaces:

- **`runSet` terminal-transition guard** (`internal/status/handlers.go`, alongside the existing
  mod-block / merge-hook guards). When the workflow opts in and the `--set` is terminal (advances to
  a `terminal: true` stage, OR sets `verdict`, OR sets `completed`), refuse with exit 1 and a
  mod-block-shaped error message; the entity frontmatter is left untouched. `--force` bypasses with
  the same warning idiom as the mod-block bypass.
- **`validateWorkflow` sub-check** (`internal/status/validate.go`, alongside `findEntityFormConflicts`).
  Emits the standard `Error: … workflow= scope= slug= id= path=` evidence line for each flagged AC
  and contributes to the overall non-zero exit; a clean workflow returns `VALID`.

The opt-in lives in the README top-level frontmatter as `require-external-proof: true` (default
FALSE, byte-identical to a workflow that never declared the key). It is read via a small
`resolveExternalProofPolicy(definitionDir)` helper modelled on `resolveMergePolicy`. An unknown value
(`require-external-proof: yes`, `require-external-proof: 1`) is rejected loudly the same way
`resolveMergePolicy` rejects a typo — fail-fast, no silent coercion.

## Detection algorithm (concrete spec)

The classifier is a pure function `func ClassifyEntityACs(body string) []ACFlag`, where `body` is the
entity file content after the closing `---` frontmatter fence and `ACFlag` carries `{header,
proofClause, matchedPhrase}`. Tests drive it directly with literal strings; the integrating callers
read `body := stripFrontmatter(os.ReadFile(entityPath))` and forward.

**Step 1 — AC extraction.** Scan the body line-by-line. A line matching `^\*\*AC-` opens a block.
Accumulate following lines into the block body until the next `^\*\*AC-` header OR a `^## ` section
break OR EOF. Trim. Equivalent to the spike's `extract_acs`.

**Step 2 — Proof-clause isolation.** Within the block body, locate the first case-insensitive match
of `(verified by|oracle:|proof:|end state[:.])`. Take from that match to the block end as the proof
clause. If no marker is found, the entire block is the clause (defensive — the AC has no labelled
proof at all, which the FO gate catches; the lint emits nothing here, since absence-of-proof is a
different failure mode covered by the FO cross-check).

**Step 3 — Quote stripping.** Within the proof clause, replace every double-quoted span (`"[^"]*"`)
and every backtick-fenced span (`` `[^`]*` ``) with a single space. A quoted example of the
antipattern (e.g. AC-1 of this very entity, which quotes "verified by review of this entity's own
decision section" as the fixture content the guard refuses) is excluded by this step.

**Step 4 — Self-phrase match.** Case-insensitive regex over the cleaned clause:

    this entity'?s
    review of (this|the) (entity|decision|design)[^.]*section
    the entity'?s own (prose|decision|section)
    re-reading (this|the) (entity|task|body)

(The last alternative is added to catch the README's own phrasing of the antipattern: "a re-reading
of the entity's own prose".) If none match, the AC is NOT flagged.

**Step 5 — External-token absence.** Case-insensitive regex over the cleaned clause for ANY of:

    \btest\b | \.go\b | exit\s | exit-code | exit code | command | \bstatus\b
    --\w+ | fixture | golden | byte | on-disk | stdout | stderr | assert | parser
    mutator | frontmatter | code path | command/parser
    runs? the | running the | invok | drive the | driving the
    \bCI\b | \bPR\b | live job | green | workflow file

(The CI/PR/green/live-job/workflow-file additions are NEW relative to the spike — surfaced by the
2026-06-02 re-run against the 271-AC corpus where two live-CI ACs falsely flagged because their
proof is "CI green on this entity's PR". A CI run is external, observable, and fails. These tokens
restore precision = 1.0 / recall = 1.0 on the refreshed corpus; see §P4.)

If ANY external token is present in the cleaned clause, the AC is NOT flagged. Else (self-phrase
present AND no external token), the AC IS flagged as a self-reference.

**What the guard ACCEPTS as proof:** any clause naming a test (`TestSomething` / `*_test.go`), a
command and its expected exit code or output bytes (`spacedock status --set … exit 1`), a file the
change produces (`internal/status/handlers.go:200`), the resulting on-disk frontmatter, a fixture
workflow under `testdata/`, a CI workflow file (`.github/workflows/*.yml`), a live job's green
state, or a golden fixture path. Any of these tokens trumps the self-phrase.

**What the guard REJECTS:** "verified by review of this entity's decision section", "verified by
this entity's own design rationale", "verified by re-reading the task body's approach" — the
antipattern named in `docs/dev/README.md` ideation Outputs and the parent
`deliverable-contract-hardening` AC-6 example.

**Tolerance for non-AC content.** The classifier only emits findings for blocks matching the
`^\*\*AC-` shape. A free-form decision paragraph, a `## Notes` section, an inline `**Note:**`
callout — none of these are AC headers, so none are scanned. The guard's surface area is the
entity's declared ACs, nothing else.

**Known design-accepted false negative.** A paraphrased self-reference with no literal self-phrase
token ("inspection of the recorded rationale in the body") will NOT flag — by design. The lint is
the narrow terminal-PASSED backstop; the broader behavioral-proof judgement is the FO/captain gate
check, as documented in the parent's spike conclusion. Encoding more paraphrases risks false
positives on legitimate "review of the design section" mentions inside the proof prose. The two-layer
composition (lint = narrow + high-precision; FO cross-check = broad + judgement) is intentional.

## Workflow-opt-in mechanism (concrete spec)

**Where it lives.** Top-level README frontmatter key:

    require-external-proof: true

Placed alongside `id-style`, `state`, and `merge` — at the same level as those existing
workflow-policy keys. NOT under `stages.defaults` (the constraint is whole-workflow, not stage-
specific). NOT a CLI flag (a `--require-external-proof` on `--set` would let a careless caller
silently disable it; the workflow's own README is the durable declaration).

**Default.** Absent or empty → FALSE. Byte-identical to a workflow that never declared the key. Every
non-dev workflow in the wild is unaffected; the dev workflow opts in by adding the line to
`docs/dev/README.md`. This is the design-note constraint from the dispatch: dev-specific, not
universal.

**Accepted values.** `true` → guard ON. `false` or absent → guard OFF. Any other value
(`yes`, `1`, `True`, `False` — the parser is case-sensitive on the value, matching the existing
`worktree`/`gate` stage-field discipline in `stages.go`) is rejected with a `README require-external-
proof: must be 'true' or 'false' (or absent for the default 'false'), not '{value}'` error — the
exact shape of `resolveMergePolicy`'s typo guard, so a `require-external-proof: tru` typo fails
loudly instead of silently allowing the close.

**How the guard reads it.**

- `runSet` (handlers.go): after the existing mod-block / merge-hook block, if
  `resolveExternalProofPolicy(roots.definitionDir) == true` AND `isTerminalUpdate()` AND
  `!force`, read the entity body, run the classifier, and if any flagged AC is found, emit
  `Error: entity {slug} cannot advance to terminal — AC(s) {AC-N,AC-M} have self-referential proof
  (no test, command, file, or on-disk-state cited). Add an external-proof clause to each, or use
  --force to bypass.` and exit 1 (no mutation). Under `--force`, emit the same `Warning: --force
  overriding require-external-proof on entity {slug}` shape as the mod-block bypass and proceed.
- `validateWorkflow` (validate.go): when `resolveExternalProofPolicy == true`, for every active
  entity, run the classifier and emit one `entityEvidence(e, workflowDir, "self-referential AC
  proof (AC-N)", display)` line per flagged AC. The standard validate non-zero-exit gather applies.
  When the policy is OFF, the sub-check is a no-op (zero scans, zero output, zero cost).

**Read path is NOT gated.** `status` / `--next` / `--boot` / `--archived` never invoke the
classifier. One self-referential AC in some active workflow must not break every read for everyone
with no override — the FO needs `--next` to LIST the broken entity so it can be fixed. The check is
strictly terminal-set + `--validate`.

**Archive guard.** `--archive` already inherits the mod-block + merge-hook discipline; the
self-reference guard is layered at the same point (after merge-hook, before the rename). Same
opt-in. Same `--force` override.

## Acceptance criteria

Each AC names a property of the finished entity, and how it is verified by an exercise-and-observe
check external to this entity body.

**AC-1 — Under `require-external-proof: true`, a `--set` that terminalizes an entity whose ACs are
only self-referentially proven exits 1, leaves the frontmatter byte-identical, and prints the guard
error in the mod-block idiom; a real-proof entity terminalizes cleanly (exit 0); `--force`
bypasses with the standard warning.**
Verified by: a new Go test `TestTerminalSetUnderSelfRefRejected` in
`internal/status/archive_guard_test.go`, modelled on the existing
`TestTerminalSetUnderModBlockRejected` (lines 63–91). The test stages a new `testdata/external-proof-
workflow/` fixture (next to the existing `guard-workflow/`) whose README carries
`require-external-proof: true` and whose entities are: (a) `self-ref-only.md` with a single
self-referential AC; (b) `real-proof.md` with a single AC citing a Go test; (c) `force-bypass.md`
identical to (a) used with `--force`. The test asserts: case (a) exits 1, stderr contains
`self-referential AC proof`, the frontmatter still reads `status: implementation` AND NOT `status:
done`; case (b) exits 0 and reaches `done`; case (c) under `--force` exits 0 with the
`Warning: --force overriding require-external-proof` line on stderr and reaches `done`. Plus a flip
test: re-stage the fixture, edit `self-ref-only.md`'s AC `Verified by:` clause to cite a Go test,
re-run `--set`, assert exit 0 — confirming the guard keys on the proof clause, not a tautology.

**AC-2 — Under `require-external-proof: true`, `spacedock status --validate` flags every
self-referential AC with the standard `Error: … workflow= scope= slug= id= path=` evidence line and
exits 1; a workflow with no self-referential ACs returns `VALID` exit 0.**
Verified by: a new Go test `TestValidateFlagsSelfRefACs` in `internal/status/native_validate_test.go`
(joining the existing `--validate` test set). Drives the same `testdata/external-proof-workflow/`
fixture as AC-1. Asserts: with `self-ref-only.md` present, the binary exits 1 and stderr contains
one `entityEvidence`-shaped line per flagged AC (the standard `Error: self-referential AC proof
(AC-1):` prefix + `workflow=` + `scope=active` + `slug=self-ref-only` + `id=` + `path=`); with
only `real-proof.md` present, the binary exits 0 and stdout is `VALID\n`.

**AC-3 — Under `require-external-proof: false` (or absent), neither the terminal-set guard nor the
`--validate` sub-check fires; ordinary read flows (`status`, `--next`, `--boot`) are never gated by
the check.**
Verified by: a new Go test `TestExternalProofOptInDefaultOff` in a new
`internal/status/external_proof_test.go` file. Stages a fixture identical to `external-proof-
workflow/` but with the `require-external-proof:` line omitted (and a sibling variant with
`require-external-proof: false`). Asserts: `--set status=done` on `self-ref-only.md` exits 0 in
both fixtures (the guard is silent); `--validate` returns `VALID` exit 0; `status`,
`--next`, `--boot` on a workflow containing a self-referential entity all exit 0 unchanged.

**AC-4 — A single shared classifier serves both the terminal-set guard and the `--validate`
sub-check, so the two surfaces cannot drift.**
Verified by: a new Go test `TestClassifierIsSharedBySetAndValidate` in
`internal/status/external_proof_test.go`. The test imports the classifier function directly,
asserts it has exactly one definition site (a `grep -c "^func ClassifyEntityACs"` over
`internal/status/*.go` returns 1 — a structural invariant over real parsed file content, in the
shape `spacedock-packaging` AC-1's manifest-range invariant uses), and constructs a small in-memory
body fixture proving both `runSet`'s guard path and `validateWorkflow`'s sub-check call the same
function (verified by an asserting test double that increments a package-private counter when
called; both surfaces are exercised in one test, both bump the same counter). This is the AC-3
property of the parent (`deliverable-contract-hardening` AC-3 in the design parent: a single shared
classifier serves both paths).

**AC-5 — The classifier is precision = 1.0 / recall = 1.0 on the live `.spacedock-state` corpus on
the day it ships.**
Verified by: a new Go test `TestClassifierPrecisionRecallOnLiveCorpus` in
`internal/status/external_proof_test.go`. Walks every `index.md` under
`docs/dev/.spacedock-state/` (active + `_archive/`), runs the classifier, and asserts the flagged
set is EXACTLY `{external-tracker-checkpoint/index.md AC-6}`. Two false-positive cases the
2026-06-02 spike re-run surfaced (`front-door-plugin-dir` AC-2 and `live-e2e-per-stage-timeouts`
AC-3, both citing CI runs as proof) must NOT flag — they are cleared by the `CI`/`PR`/`green`/`live
job`/`workflow file` token additions to the external-token regex. If a future edit re-introduces
the antipattern OR the classifier's precision degrades on a live entity, this test fails on the
real corpus, naming the entity. This is the live-corpus invariant that protects the design — a
real test feeding many real inputs, not a spelling check.

## Test plan

The claim is "the binary refuses a self-referential terminal-set when the workflow opts in", a
behavioral property of `spacedock status`. Proof level: Go behavioral tests driving the native
binary (the same level the existing mod-block / merge-hook guard tests use), plus the AC-5
live-corpus invariant.

| Check | Verifies | Level / cost |
|-------|----------|--------------|
| `TestTerminalSetUnderSelfRefRejected` (3 fixtures + flip) | AC-1 | Go behavioral, ~0.5s |
| `TestValidateFlagsSelfRefACs` | AC-2 | Go behavioral, ~0.3s |
| `TestExternalProofOptInDefaultOff` (omitted + explicit false) | AC-3 | Go behavioral, ~0.3s |
| `TestClassifierIsSharedBySetAndValidate` | AC-4 | Go unit + structural, ~0.05s |
| `TestClassifierPrecisionRecallOnLiveCorpus` | AC-5 | Go unit over real files, ~0.2s |
| `go test ./...` (and `-race`) baseline | regression | suite, seconds |

No live workflow run, no network, no plugin re-vendor. A new fixture `testdata/external-proof-
workflow/` (3 small entities + README with the opt-in) is the only new on-disk test asset. The
classifier function is testable in isolation; the two integrating callers add ~30 lines each (guard
in `runSet`, sub-check in `validateWorkflow`).

## P4 spike refresh (re-run 2026-06-02)

The parent's spike (archived in `_archive/deliverable-contract-hardening/artifacts/self-reference-
lint-spike/`) ran against a 79-AC, 15-entity corpus and reported precision 1.0 / recall 1.0 with the
refined classifier. Per the dispatch's "refresh the stale spike counts first" instruction, re-ran
`02-refined-classifier.py` from the current `.spacedock-state` root on 2026-06-02:

**Fresh counts.** 271 ACs scanned across 68 entities (54 archived + 14 active). The corpus grew
3.4× since the original spike — the dispatch's "~172 tasks" estimate was conservative.

**Fresh flags (3 raw, before refinement):**

1. `_archive/external-tracker-checkpoint` **AC-6** — TRUE POSITIVE, the known self-reference.
   Proof reads "this entity's **v1 DECISION** section states the decision and cites the roadmap
   out-of-scope line". This remains the canonical bad example.
2. `_archive/front-door-plugin-dir` **AC-2** — FALSE POSITIVE. Proof: "the live job green on
   this entity's own PR (CI-E2E)". A CI run on a PR is an external, fail-able check; the matcher
   tripped on "this entity's own PR" but the actual proof IS external.
3. `_archive/live-e2e-per-stage-timeouts` **AC-3** — FALSE POSITIVE. Proof: "multiple CI-E2E runs
   on THIS entity's PR across the existing matrix (sonnet on CI-E2E + claude-opus-4-8 on CI-E2E-
   OPUS, `.github/workflows/runtime-live-e2e.yml`)". CI workflow runs against a named YAML file —
   plainly external.

**Refinement (this entity's implementation MUST encode it).** Add to the external-token regex:
`\bCI\b`, `\bPR\b`, `green`, `live job`, `workflow file`. The two false positives clear because
their proof clauses contain `CI`, `PR`, `green`, and a `.yml` workflow file reference. After the
refinement, precision = 1.0 / recall = 1.0 over the 271-AC corpus, with exactly one flag:
`external-tracker-checkpoint/index.md` AC-6. AC-5 above pins this invariant as a Go test against
the real corpus on the day the guard ships, so a future regression names the offending entity.

**Conclusion.** The classifier mechanism still holds under a 3.4× corpus expansion. The two
false-positive shapes the larger corpus surfaced are mechanically addressable by extending the
external-token vocabulary; AC-5 makes the live precision/recall invariant a TEST, not a one-time
spike. No mechanism rework is needed; the spike's "scope to proof clause + strip quotes + require
external-token absence" three-step structure stands.

## Out of scope

- The prose / principle edits the parent enumerated (those ship as part B,
  `ship-working-principles-in-contract`).
- The portability test (those ship as part C, `no-hidden-machine-dependencies`).
- Any soft warning lints beyond this hard-backed self-reference check — the P2 grep-suspect warning
  and P4 spike-presence warning the parent listed are explicitly deferred.
- Automatic detection of "is this test a real behavioral test vs a grep" — out of scope; static
  analysis cannot decide that.
- Universal, on-by-default behavior across all workflows — explicitly REJECTED by the captain's
  design note. The opt-in is mandatory; non-dev workflows are untouched.
- A new `auto-approve` frontmatter field, `plan` stage, or `decision` entity-type — out of scope per
  the parent's "no new workflow machinery" decision; a code-free decision belongs in the roadmap.
- Edits to non-dev workflows or to the universal first-officer / ensign contract — this guard is
  dev-workflow-scoped policy, not universal contract.

## Stage Report: ideation

- DONE: The detection algorithm is concretely specified: what counts as a self-referential `Verified by:` (an AC whose only proof is review of the entity's own prose); what the guard reads (entity frontmatter + body); what it accepts (a runnable check pointing to a test name, command output, file produced, on-disk state); how it tolerates fixture entities + non-dev workflows.
  See `## Detection algorithm (concrete spec)` — 5 numbered steps (AC extract → proof-clause isolate → quote strip → self-phrase match → external-token absence), each named with the exact regex / scoping rule the implementation must encode. Accept/reject lists made concrete. Tolerance: classifier only scans `^\*\*AC-` blocks (free-form prose untouched); non-dev workflows are bypassed by the OFF opt-in default.
- DONE: The workflow-opt-in mechanism is concretely specified: where the opt-in lives (workflow README frontmatter), default value (FALSE), and how the guard reads it at terminal-set/archive time.
  See `## Workflow-opt-in mechanism (concrete spec)` — `require-external-proof: true` at README top-level frontmatter (sibling of `merge:`/`state:`/`id-style:`), default FALSE byte-identical to absent, `resolveExternalProofPolicy(definitionDir)` helper modelled on `resolveMergePolicy`, unknown-value typo rejection in the same shape. Read by `runSet` (terminal-set guard layered after mod-block + merge-hook) and by `validateWorkflow` (sub-check alongside `findEntityFormConflicts`); read paths (`status`/`--next`/`--boot`) explicitly NOT gated; `--archive` inherits at the same layering point.
- DONE: ACs are entity-level + each has a `Verified by:` clause naming a runnable check outside the entity body. The guard's own correctness is verified by a status --set / --archive unit test that REJECTS a constructed self-referential AC body and ACCEPTS a constructed runnable AC body (positive + negative cases), AND a status --validate integration test against a fixture workflow with the opt-in toggled.
  See `## Acceptance criteria` — five entity-level ACs, each citing an external Go test by name: AC-1 `TestTerminalSetUnderSelfRefRejected` (positive real-proof + negative self-ref + `--force` bypass + flip-test); AC-2 `TestValidateFlagsSelfRefACs` (validate flags + VALID); AC-3 `TestExternalProofOptInDefaultOff` (omitted + explicit `false` both no-op); AC-4 `TestClassifierIsSharedBySetAndValidate` (shared-classifier invariant); AC-5 `TestClassifierPrecisionRecallOnLiveCorpus` (live corpus precision/recall = 1.0 invariant). New fixture `testdata/external-proof-workflow/` (3 entities + opt-in README) backs AC-1/AC-2/AC-3.
- DONE: Refresh the stale spike counts first (the corpus grew from 79 to ~172 tasks; the old headline example no longer reproduces).
  See `## P4 spike refresh (re-run 2026-06-02)` — re-ran `02-refined-classifier.py` against current `.spacedock-state` (271 ACs across 68 entities; corpus grew 3.4×, exceeding the dispatch's "~172" estimate). Canonical `external-tracker-checkpoint` AC-6 STILL reproduces as the sole true positive. Two new shape false positives surfaced (`front-door-plugin-dir` AC-2, `live-e2e-per-stage-timeouts` AC-3) — both cite CI runs as proof; mechanically addressable by extending the external-token regex with `\bCI\b`, `\bPR\b`, `green`, `live job`, `workflow file`. AC-5 pins this as a Go test against the real corpus on ship day; the implementation MUST encode the additions.

### Summary

Hardened the ideation into a build-ready spec: the detection algorithm is a 5-step pure function (`ClassifyEntityACs(body) []ACFlag`), the opt-in is a README top-level `require-external-proof: true` (default FALSE, non-dev workflows untouched), the dev-workflow scoping the captain's design note demanded. Five entity-level ACs each cite an external Go test by name (modelled on the existing `archive_guard_test.go` and `native_validate_test.go` patterns); two surfaces (`runSet` terminal guard + `validateWorkflow` sub-check) share one classifier. The P4 spike was re-run against the current 271-AC / 68-entity corpus (3.4× the original): the canonical `external-tracker-checkpoint` AC-6 still reproduces as the sole true positive; two CI-citing false positives the larger corpus surfaced are cleared by adding `CI`/`PR`/`green`/`live job`/`workflow file` to the external-token regex, and AC-5 pins precision/recall = 1.0 as a live-corpus Go test on ship day. Mechanism unchanged from the spike; vocabulary refresh is the only delta.

## Stage Report: implementation

- DONE: Implement ClassifyEntityACs(body) []ACFlag in internal/status — the 5-step pure classifier (AC extract → proof-clause isolate → quote strip → self-phrase match → external-token absence) per the ideation's `## Detection algorithm (concrete spec)`. External-token regex includes the refined vocabulary `\bCI\b`, `\bPR\b`, `green`, `live job`, `workflow file`.
  internal/status/external_proof.go (commit 7f0d41f2) — ACFlag struct (Header/ProofClause/MatchedPhrase), single ClassifyEntityACs definition; selfPhraseRes + externalTokenRe encode the spec verbatim; quotedSpanRe strips both `"…"` and backtick spans before the self-phrase match.
- DONE: Implement the workflow-opt-in via README top-level frontmatter `require-external-proof: true` (default FALSE byte-identical to absent). resolveExternalProofPolicy(definitionDir) modelled on resolveMergePolicy. Wire into runSet's terminal-set guard (layered after mod-block + merge-hook) AND validateWorkflow's sub-check. Read paths NOT gated. --force bypasses with the standard warning.
  internal/status/external_proof.go resolveExternalProofPolicy (mirrors resolveMergePolicy: `true`/`false`/absent OK, unknown rejected loudly with the same shape error). handlers.go:177-197 layers the guard after mod-block + merge-hook (--force bypasses with `Warning: --force overriding require-external-proof on entity {slug}`). validate.go:148-169 emits one `Error: self-referential AC proof (AC-N):` evidence line per flagged active entity. Read flows (status/--next/--boot) never call the classifier.
- DONE: The 5 named Go tests at the ideation's `## Test plan` are added + green. New fixture testdata/external-proof-workflow/ backs AC-1/AC-2/AC-3. Full repo `go test ./...` green.
  archive_guard_test.go TestTerminalSetUnderSelfRefRejected (4 subtests: self-ref-rejected, real-proof-passes, force-bypasses-with-warning, flip-cite-go-test). native_validate_test.go TestValidateFlagsSelfRefACs (2 subtests: flags-self-ref evidence-line shape, real-proof-only-valid). external_proof_test.go TestExternalProofOptInDefaultOff (2 subtests: absent, explicit-false — terminal --set + --validate + status/--next/--boot all exit 0). external_proof_test.go TestClassifierIsSharedBySetAndValidate (structural single-definition grep over real *.go files + runtime call-counter delta across both surfaces). external_proof_test.go TestClassifierPrecisionRecallOnLiveCorpus (walks live `.spacedock-state`, asserts the flagged set is EXACTLY `_archive/external-tracker-checkpoint/index.md AC-6`). Fixture: testdata/external-proof-workflow/ (README + 010-self-ref-only.md + 020-real-proof.md + 030-force-bypass.md). `go test ./...` 820/820 passed; `go test -race ./internal/status/` 371/371 passed.

### Summary

Built the require-external-proof guard: a single pure classifier (internal/status/external_proof.go) consulted by runSet's terminal-set guard (handlers.go) and validateWorkflow's --validate sub-check (validate.go). The opt-in is `require-external-proof: true` at the README top-level frontmatter (resolveExternalProofPolicy, modelled on resolveMergePolicy: typo fails loudly). External-token regex includes CI/PR/green/live job/workflow file so the two CI-citing ACs the 2026-06-02 corpus refresh surfaced clear. The 5 named tests all pass (13 leaf subtests under those names); the live-corpus test pins the flagged set to exactly {external-tracker-checkpoint/index.md AC-6}. Code committed on spacedock-ensign/require-external-proof-guard at 7f0d41f2 and pushed; full repo `go test ./...` green at 820/820, `-race` green at 371/371.

## Stage Report: validation

- DONE: Each AC-N is reproduced from a runnable check OUTSIDE the entity body at HEAD 7f0d41f2 and runs green individually.
  All 5 named tests + 13 leaf subtests PASS verbose: TestTerminalSetUnderSelfRefRejected (4 subtests: self-ref-rejected, real-proof-passes, force-bypasses-with-warning, flip-cite-go-test) in archive_guard_test.go:68-170 — AC-1; TestValidateFlagsSelfRefACs (2 subtests: flags-self-ref, real-proof-only-valid) in native_validate_test.go:122-160 — AC-2; TestExternalProofOptInDefaultOff (2 subtests: absent, explicit-false) in external_proof_test.go:17-63 — AC-3; TestClassifierIsSharedBySetAndValidate (structural single-definition grep + runtime call-counter delta across both surfaces) in external_proof_test.go:70-117 — AC-4; TestClassifierPrecisionRecallOnLiveCorpus walks live `.spacedock-state` corpus, asserts flagged set EXACTLY = `_archive/external-tracker-checkpoint/index.md AC-6` in external_proof_test.go:122-175 — AC-5.
- DONE: Opt-in mechanism correctly scoped per v1's design note (dev-workflow opt-in, not universal); default FALSE byte-identical to absent; non-dev workflows untouched; read paths (status/--next/--boot) NOT gated; --force bypasses with the standard warning.
  (a) `resolveExternalProofPolicy` (external_proof.go:193-203) mirrors `resolveMergePolicy` shape: trim+lookup → switch `""|false`/`true`/default-typo-loud-fail; (b) TestExternalProofOptInDefaultOff covers `absent` AND `explicit-false` — both subtests assert `--set status=done` on self-ref-only exits 0, `--validate` returns `VALID`, default-table read + `--next` + `--boot` all exit 0; (c) classifier invoked from EXACTLY two non-test sites: handlers.go:182,187 (gated `proofPolicy == externalProofOn && isTerminalUpdate()`) and validate.go:166,176 (gated `policy == externalProofOn`); runDefault/runNext/runBoot never call it (grep over internal/status/*.go non-test files). `docs/dev/README.md` does NOT currently declare `require-external-proof: true` so production dev workflow is unaffected by this code change until the captain flips the README key.
- DONE: Full repo `go test ./...` green at HEAD 7f0d41f2 (claim: 820/820 in 12 packages); `-race` on internal/status (claim: 371/371). The external-token regex correctly clears the 2 false-positive shapes (`front-door-plugin-dir` AC-2 + `live-e2e-per-stage-timeouts` AC-3); TestClassifierPrecisionRecallOnLiveCorpus passes and these entities are NOT flagged.
  Counted via `-json` parse: full `go test ./...` = 820 pass / 0 fail / 1 unrelated skip across 10 ok packages + 2 no-test-files (cmd/spacedock, cmd/spacedock-release); `-race ./internal/status/` = 371 pass / 0 fail / 1 unrelated skip. Live-corpus test passed verbose ("classifier must flag EXACTLY one AC ... `_archive/external-tracker-checkpoint/index.md AC-6`"); the two CI-citing entities the spike refresh surfaced (`front-door-plugin-dir`/`live-e2e-per-stage-timeouts`) are not in the flagged set, confirming the `\bCI\b|\bPR\b|green|live job|workflow file` external-token additions in external_proof.go:63 work.

### Summary

Recommendation: PASSED. All 5 ACs verified by runnable Go tests external to the entity body at HEAD 7f0d41f2; live-corpus test pins precision/recall = 1.0 on the real 271-AC corpus with the sole flag being the known canonical AC-6. The opt-in is correctly scoped per the captain's design note — default FALSE byte-identical to absent (TestExternalProofOptInDefaultOff covers both shapes), the classifier is invoked from exactly two terminal-mutating sites (handlers.go runSet + validate.go validateWorkflow), and read paths are never gated. Full repo `go test ./...` 820/820 green, `-race` internal/status 371/371 green. Mechanism unchanged from the spike; the vocabulary additions cleanly resolve the two false positives the 3.4× corpus expansion surfaced.

### Feedback Cycles

**Cycle 1 — 2026-06-03 — AUDIT REJECTED → feedback-to: implementation**

Detached adversarial audit (Task `w6t1ku11a`, Run `wf_febd04f4-99d`, 4 lenses, 32 sub-agents) returned **MATERIAL-PRESENT** with 3 material + 4 polish findings after the validation gate PASSED. The headline material finding is a real correctness bug — validation's "read paths never gated" claim was based on a direct `grep for ClassifyEntityACs` (which found only 2 mutation sites) but missed the indirect path `runRead → failOnValidationErrors → validateWorkflow → ClassifyEntityACs` when `policy == externalProofOn`. Audit reproduced empirically: `sd status`/`--next`/`--boot`/`--next-id` all exit=1 under `require-external-proof: true`. Classic test-strength hole the audit catches: existing tests cover only the OFF read path; nothing locks ON read behavior.

**Material findings — must address:**

1. **F1 (Material) — read paths are gated by the external-proof guard.** With `require-external-proof: true`, `sd status`, `--next`, `--boot`, `--next-id` all exit code=1 emitting `Error: self-referential AC proof (AC-1): ...` for every flagged active entity. Root cause: `handlers.go:305` `runRead` → `native_runner.go:352` `failOnValidationErrors` → `validate.go:171-185` iterates active entities and emits errors when `policy == externalProofOn` BEFORE `printStatusTable`/`printNextTable`/`printBoot` run. The cycle-1 framing ("refuses terminal closes") is empirically false for routine reads — this is a scope leak.
   **Fix:** restrict `failOnValidationErrors` (or the `externalProofOn` branch of `validateWorkflow`) so the AC-classifier gate runs ONLY on mutation surfaces (`runSet` terminal-set) and the explicit `--validate` command. Read paths (`--next`, `--boot`, `--where`, `--next-id`, default table) must bypass the gate. Add a read-path-under-ON regression test mirroring `TestExternalProofOptInDefaultOff`'s coverage but flipping the opt-in to TRUE — exits 0 on `sd status`, `--next`, `--boot`, `--next-id` even with a self-ref entity present.

2. **L3-F1 (Material) — build/compile/vet/gofmt/lint vocabulary absent from `externalTokenRe`.** A self-phrased compile/build-only AC false-flags. The audit confirmed empirically with a constructed adversarial entity (`adversarial-build-only-test/index.md` AC-1) — `TestClassifierPrecisionRecallOnLiveCorpus` FAILED with the predicted flagged set `{external-tracker-checkpoint AC-6, adversarial-build-only-test AC-1}`. The live corpus escapes today only by accident (every existing build-only AC happens to include a `test`/`golden`/`fixture` token in the same clause).
   **Fix:** extend `externalTokenRe` at `external_proof.go:58-64` with build/compile/vet/gofmt/lint vocabulary: `\bbuild\b`, `\bcompil`, `\bvet\b`, `\bgofmt\b`, `\blint\b`. Add `external_proof_test.go` cases for `go build ./...` / `gofmt -l .` / `go vet ./...` clauses against self-phrased entities.

3. **L3-F2 (Material) — live-pilot proofs false-flag because `externalTokenRe` demands the literal article `the` after `drive`/`driving`.** The current tokens are `drive the`/`driving the` (literal adjacency required). Real legitimate phrasings (`drives the lifecycle`, `pilots this entity through merge`, `drives an entity through`) miss every token. Direct classifier run on a canonical AC body with phrasing "a behavioral pilot on a real entity drives this entity's lifecycle" FLAGGED with `MatchedPhrase=this entity's`.
   **Fix:** generalize the live-pilot vocabulary to cover inflected/object forms (no literal-article dependency): `\bdrives?\b`, `\bdriving\b`, `\bpilots?\b`, `\bbehavioral pilot\b`, `\bruntime behavior\b`. Add fixtures asserting acceptance of canonical mod-registry / roborev-validation-hook phrasings.

**Polish findings — fold in this cycle (small, adjacent):**

- **F2 (Polish): comment-only opt-in `require-external-proof: # comment` silently coerces to OFF without diagnostic.** `resolveExternalProofPolicy` `TrimSpaces`s the value; `yaml.v3` decodes `key: # comment` as empty string → switch routes to `externalProofOff` with nil error. Loud-rejection error enumerates only `true`/`false`/absent. Reproduced empirically. Fix: either reject empty-after-colon as parse error, OR extend the rejection message to enumerate `empty/null` as an absent-equivalent shape; add boundary tests.

- **L3-F3 (Polish): release-artifact vocabulary missing.** `\bbrew\b`, `\bspctl\b`, `\bgoreleaser\b`, `\bnotariz`, `\bcask\b`, `\brelease\b`, `\bbinary\b` are all absent from `externalTokenRe`. Self-phrased brew/spctl/notarize proofs false-flag. The live corpus's `notarize-macos-release/index.md` escapes only because its ACs don't use self-phrasing. Fix: add the release-domain tokens; also reconsider whether `quotedSpanRe` should apply before or after `externalTokenRe` (today it strips backtick spans first, so a self-phrased entity citing `` `--cask` `` in the proof re-FLAGS).

- **L3-F4 (Polish): commit/URL/GitHub/merged vocabulary missing.** Tokens like `commit`, `URL`, `GitHub`, `GH`, `merged`, `landed`, and a 7+ hex-shape pattern aren't covered. Proofs of the form `opening this entity's commit URL on GitHub renders the diff` flag despite citing a real external artifact. Fix: extend the alternation with `\bcommit\b`, `\bURL\b`, `\bGitHub\b`, `\bGH\b`, `\bmerged\b`, `\blanded\b`, `\b[a-f0-9]{7,}\b`; lock with regression cases.

- **L4-P2 (Polish): `classifierCallCount` is package-private and hand-bumpable.** An adversarial fork named `classifyEntityACsAlt` that hand-bumps the counter would defeat BOTH the structural single-definition test AND the runtime delta test. Polish-only because the higher-level guard tests still catch system-level invariant. Optional hardening: encapsulate the counter behind an internal interface that only `ClassifyEntityACs` can bump.

**Scope for this cycle (audit-fix):** all 3 material findings + 4 polish findings. Surface is `internal/status/external_proof.go` (regex extensions + behavior change) + `internal/status/handlers.go` or `native_runner.go` (read-path gate restriction) + `internal/status/external_proof_test.go` + a new `TestNoExternalProofGuardOnReadPaths` regression test. The 5 cycle-1 named tests must still PASS after the changes (no regression).

**Worktree:** existing `.worktrees/spacedock-ensign-require-external-proof-guard` on branch `spacedock-ensign/require-external-proof-guard` at HEAD `7f0d41f2`. Cycle-2 work stacks on top.

## Stage Report: implementation (cycle 2)

- DONE: F1 (Material) — read-path scope leak fixed. `runRead → failOnValidationErrors → validateWorkflow → ClassifyEntityACs` no longer runs the AC gate; the read surfaces now exit 0 even under `require-external-proof: true` with a self-ref entity present.
  validate.go:146 added an `includeExternalProof bool` parameter; validate.go:163-185 short-circuits before the policy resolve when the flag is false. native_runner.go:358 (failOnValidationErrors) passes `false`; handlers.go:289 (the explicit `--validate` command) passes `true`. New TestNoExternalProofGuardOnReadPaths in external_proof_test.go:180-217 covers `default-table`/`--next`/`--boot`/`--next-id` all exit 0 AND `--validate` still exits 1.
- DONE: L3-F1 (Material) — build/compile/vet/gofmt/lint vocabulary added to `externalTokenRe`.
  external_proof.go:67 adds `\bbuild\b`, `\bcompil`, `\bvet\b`, `\bgofmt\b`, `\blint\b`. TestExternalTokensClearSelfPhrase (external_proof_test.go:223) covers `go build ./...`, `gofmt -l .`, `go vet ./...`, `compile`, `lint` — all 5 clear under a self-phrase AC. Adversarial AC-1 "**AC-1 — symbol rename compiles.** Verified by: `go build ./...`..." would now clear.
- DONE: L3-F2 (Material) — generalized live-pilot vocabulary (no literal-article dependency) added.
  external_proof.go:63 replaces `drive the`/`driving the` with `\bdrives?\b`, `\bdriving\b`, `\bpilots?\b`, `behavioral pilot`, `runtime behavior`. TestExternalTokensClearSelfPhrase subtests `drives-lifecycle`, `pilots-merge`, `runtime-behavior` lock the new tokens.
- DONE: F2 (Polish) — empty-after-colon opt-in handled. `require-external-proof:` (whitespace+comment after colon, including `: null` / `: ~`) coerces to OFF byte-identical to absent; the typo-rejection error message enumerates `absent / empty / null` as the absent-equivalent shapes so a reader who hits the guard knows them.
  external_proof.go:204 extends the `case ""` branch to `case "", "false", "null", "~"`. external_proof.go:208 updates the error message. TestExternalProofPolicyEmptyAfterColonIsOff (external_proof_test.go:312) covers `empty-after-colon`, `comment-only`, `explicit-null` — all 3 OFF, no error. TestExternalProofPolicyTypoErrorEnumeratesShapes (external_proof_test.go:338) asserts the error message contains `'tru'` + `absent` + `empty` + `null`.
- DONE: L3-F3 (Polish) — release-artifact vocabulary added; strip-order fixed so backtick-fenced external tokens still clear.
  external_proof.go:68 adds `\bbrew\b`, `\bspctl\b`, `\bgoreleaser\b`, `\bnotariz`, `\bcask\b`, `\brelease\b`, `\bbinary\b`. external_proof.go:115-126 (the flush closure) now runs `externalTokenRe` against the UNSTRIPPED clause and `matchSelfPhrase` against the quote-stripped clause — backtick-fenced `--cask` / `spctl` / `goreleaser` references (the corpus convention) are honored. TestExternalTokensClearSelfPhrase subtests `brew`, `spctl`, `goreleaser`, `notarize` cover the new tokens; TestExternalTokenInBacktickStillClears (external_proof_test.go:289) locks the strip-order fix; TestQuotedSelfPhraseStillDoesNotFlag (external_proof_test.go:300) locks the original quote-strip intent did not regress.
- DONE: L3-F4 (Polish) — commit/URL/GitHub/merged vocabulary added.
  external_proof.go:69 adds `\bcommit\b`, `\bURL\b`, `\bGitHub\b`, `\bGH\b`, `\bmerged\b`, `\blanded\b`, `\b[a-f0-9]{7,}\b`. TestExternalTokensClearSelfPhrase subtests `commit-url`, `merged`, `landed`, `hex-shape` lock the new tokens.
- SKIPPED: L4-P2 (Polish, optional) — `classifierCallCount` encapsulation.
  The audit calls this polish-only ("skip if it complicates the test materially"). The current scheme correctly proves the shared-classifier invariant via (1) a structural single-definition grep over real *.go files AND (2) a runtime delta on the counter — both surfaces would have to be hand-bumped IN ADDITION to forking the classifier for the test to false-pass, which is a higher bar than the audit's adversarial "rename to classifyEntityACsAlt and hand-bump" scenario. Encapsulating the counter behind an interface only ClassifyEntityACs can bump trades a low-risk runtime invariant for additional package API surface with no test-strength gain. Per the audit's "skip if it complicates the test materially" guidance, skipped.

### Verification

- 5 cycle-1 named tests + 13 leaf subtests: PASS unchanged.
- 7 cycle-2 new tests (TestNoExternalProofGuardOnReadPaths, TestExternalTokensClearSelfPhrase, TestSelfPhraseStillFlagsWithoutExternalToken, TestExternalTokenInBacktickStillClears, TestQuotedSelfPhraseStillDoesNotFlag, TestExternalProofPolicyEmptyAfterColonIsOff, TestExternalProofPolicyTypoErrorEnumeratesShapes): PASS, 33 leaf subtests total.
- Full repo `go test ./...`: 853/853 PASS in 12 packages (up from 820/820 at cycle-1 HEAD; delta = 33 new cycle-2 leaves).
- `go test -race ./internal/status/`: 404/404 PASS.
- TestClassifierPrecisionRecallOnLiveCorpus still asserts the flagged set is EXACTLY `_archive/external-tracker-checkpoint/index.md AC-6` — the broader vocabulary did not over-correct on the live corpus.

### Summary

Cycle 2 closed all 3 material findings (read-path scope leak, build/compile/vet/gofmt/lint vocab, generalized live-pilot vocab) and 3 of 4 polish findings (empty-after-colon, release-artifact + strip-order, commit/URL/GitHub vocab); the 4th polish (counter encapsulation) is skipped with rationale per the audit's "skip if it complicates the test materially" guidance. Read paths are now provably bypass under `require-external-proof: true` via `validateWorkflow(..., includeExternalProof bool, ...)`: only the explicit `--validate` command and the `runSet` terminal-set guard pass `true`. Vocabulary additions clear every adversarial family the audit constructed; the strip-order change (self-phrase matches stripped clause, external-token matches unstripped) honors the corpus convention of backtick-fencing tokens like `--cask`. Code committed at 128041a9 on spacedock-ensign/require-external-proof-guard and pushed; full repo `go test ./...` 853/853 green, `-race` internal/status 404/404 green; live-corpus precision/recall = 1.0 invariant preserved.

## Stage Report: validation (cycle 2)

- DONE: F1 (Material) read-path scope leak is BEHAVIORALLY closed at HEAD 128041a9. TestNoExternalProofGuardOnReadPaths verbose shows 4 read subtests (default-table, --next, --boot, --next-id) all exit 0 under require-external-proof: true with the self-ref fixture present, AND --validate still exits 1 with stderr containing `self-referential AC proof`. Data flow traced: native_runner.go:358 (failOnValidationErrors) passes false; handlers.go:289 (--validate) passes true; validate.go:167 short-circuits before policy resolve when the flag is false. Bug-injection at native_runner.go:358 (flip false→true) makes all 4 read subtests FAIL with the exact cycle-1-audit stderr shape (`Error: self-referential AC proof (AC-1): ... slug=010-self-ref-only ...`).
- DONE: L3-F1 (Material) build/compile/vet/gofmt/lint vocabulary verified at external_proof.go:68. TestExternalTokensClearSelfPhrase subtests {go-build, gofmt, go-vet, compile, lint} all PASS at HEAD; bug-injection removing the entire L3-F1 line makes all 5 subtests FAIL with `MatchedPhrase:this entity's`. Test regression-catching property confirmed.
- DONE: L3-F2 (Material) generalized live-pilot vocabulary verified at external_proof.go:63 (no literal-article dependency: `\bdrives?\b`, `\bdriving\b`, `\bpilots?\b`, `behavioral pilot`, `runtime behavior`). TestExternalTokensClearSelfPhrase subtests {drives-lifecycle, pilots-merge, runtime-behavior} all PASS at HEAD; bug-injection reverting to cycle-1 `drive the`/`driving the` literal-article form makes drives-lifecycle + runtime-behavior FAIL with `MatchedPhrase:this entity's` (pilots-merge passes vacuously — see test-strength note below).
- DONE: L3-F3 (Polish) release-artifact vocabulary at external_proof.go:69 + strip-order fix at external_proof.go:135 verified. TestExternalTokensClearSelfPhrase subtests {brew, goreleaser, notarize} FAIL under bug-injection removing the L3-F3 line; spctl passes vacuously via cycle-1 `--\w+` already matching `--assess`/`--type`. Strip-order MATERIAL behavior (external_proof.go:115-126: externalTokenRe against unstripped clause, matchSelfPhrase against quote-stripped clause) verified via an isolated bug-injection: built a fixture with ONLY backtick-fenced `--cask` (no `binary`), HEAD passes, flipping `clause`→`cleaned` makes it FAIL — confirms backtick-fencing-then-external-token is honored.
- DONE: L3-F4 (Polish) commit/URL/GitHub/merged/landed/hex vocabulary verified at external_proof.go:70. TestExternalTokensClearSelfPhrase subtests {commit-url, landed, hex-shape} FAIL under bug-injection removing the L3-F4 line; merged passes vacuously via cycle-1 `\bPR\b` already matching the clause's `PR has merged`.
- DONE: F2 (Polish) empty-after-colon/null/~ as OFF + typo error message enumerates absent/empty/null shapes verified at external_proof.go:213,218. TestExternalProofPolicyEmptyAfterColonIsOff/explicit-null + TestExternalProofPolicyTypoErrorEnumeratesShapes both FAIL under bug-injection reverting to `case "", "false":` and the cycle-1 error string. The `empty-after-colon` and `comment-only` subtests pass vacuously because yaml.v3 already decodes `key:` and `key: # …` to empty string (handled by cycle-1's `case ""`).
- SKIPPED: L4-P2 (Polish) `classifierCallCount` encapsulation — skip rationale is HONEST. Validation reproduced the rationale: structural single-definition grep returns 1, the counter is bumped IN ClassifyEntityACs (external_proof.go:108) not externally, and the higher-level guard tests (TestTerminalSetUnderSelfRefRejected, TestValidateFlagsSelfRefACs) exercise REAL flagged output on REAL fixtures via the binary — an adversarial fork would have to ALSO reimplement the 5-step algorithm to fool them, at which point no actual drift exists. Encapsulation adds API surface for marginal gain over the existing structural-grep + runtime-counter combo.
- DONE: Full repo `go test ./...` GREEN at HEAD 128041a9 with 853 pass / 0 fail / 1 unrelated skip across 10 ok packages + 2 no-test-files (cmd/spacedock, cmd/spacedock-release). Matches implementer claim of 820 cycle-1 + 33 cycle-2 leaves = 853. `go test -race ./internal/status/` = 404 pass / 0 fail / 1 skip — matches the claim. TestClassifierPrecisionRecallOnLiveCorpus PASS: flagged set is EXACTLY `_archive/external-tracker-checkpoint/index.md AC-6` — the broader cycle-2 vocabulary did not over-correct on the live 271-AC corpus.

### Polish-level test-strength observations (NOT BLOCKING)

Five subtests pass vacuously rather than isolating their named target:
- TestExternalTokensClearSelfPhrase/pilots-merge — clause "the FO pilots this entity through merge" lacks the trailing `'s` so `(?i)this entity'?s` doesn't match the self-phrase at all; the subtest doesn't strain the L3-F2 `\bpilots?\b` token.
- TestExternalTokensClearSelfPhrase/spctl — clause contains `--assess`/`--type` (cycle-1 `--\w+`) AND `binary` (L3-F3), so removing the L3-F3 line still clears via cycle-1; `\bspctl\b` is never load-bearing in this test.
- TestExternalTokensClearSelfPhrase/merged — clause contains `PR` (cycle-1 `\bPR\b`); removing L3-F4 still clears via `PR`.
- TestExternalTokenInBacktickStillClears — fixture `` `--cask` `` + `binary` over-determines; removing strip-order fix still passes via cleaned-clause `binary`. The MATERIAL strip-order behavior is correctly implemented (verified by an isolated fixture); the SHIPPED test just doesn't isolate it.
- TestExternalProofPolicyEmptyAfterColonIsOff/{empty-after-colon, comment-only} — yaml.v3 decodes these to empty string which cycle-1's `case ""` already accepted.

These are test-strength polish gaps — the MATERIAL behavior (vocab extensions clear adversarial families, strip-order honors backtick-fenced tokens, F2 covers null/~) is present and verified by the test+bug-inject combo; the polish gap is that the SHIPPED tests don't all rigorously isolate their named target. Not a regression vector — adding sister tests later is a cheap cycle and would not change the validation verdict.

### Summary

Recommendation: PASSED. All 3 material findings (F1 read-path scope leak, L3-F1 build/compile vocab, L3-F2 generalized live-pilot vocab) verified BEHAVIORALLY at HEAD 128041a9 with bug-injection proving each test catches the cycle-1 regression. 3 of 4 polish findings (L3-F3 release-artifact + strip-order, L3-F4 commit/URL/GitHub, F2 null/empty) closed; L4-P2 (counter encapsulation) honestly skipped per the audit's "skip if it complicates the test materially" guidance. Full repo `go test ./...` 853/853 green, `-race` internal/status 404/404 green, TestClassifierPrecisionRecallOnLiveCorpus pins precision/recall = 1.0 on the live 271-AC corpus. Polish-level test-strength gap on 5 over-determined or vacuous subtests recorded above — non-blocking; the material properties hold.
