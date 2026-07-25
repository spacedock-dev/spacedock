# Durable-decisions — pre-cut review of the recorded gate lifecycle

**Target:** `deac7f8af9b0f56ceea6231976d7ffaa8ab2fc51` (short `deac7f8a`), the tip of `origin/main`.

**Pinning discipline.** The audit ran on a detached throwaway checkout at that SHA, never a working
tree. `git status --porcelain` printed zero lines before and after every command; nothing was staged,
edited, or committed. `git rev-parse HEAD` was re-checked at the end of the audit and still returned
`deac7f8af9b0f56ceea6231976d7ffaa8ab2fc51` — the target did not move. **Any movement of `main`
invalidates this audit**; every `file:line`, byte count, and suite result below is a claim about that
one commit and nothing else. All `file:line` citations were verified against the pinned checkout.

## Scope limits

This is **not** the sprint's completed pre-cut audit. Two limits, both material to how much weight
this document can carry:

1. **Four sprint members are still at `ideation`** — `rq` (git-root review v1 materialization), `s4`
   (provider-neutral gate rooms), `rd` (merge-guard archive durability), `v3v` (foreign runtime marker
   scrub). The roadmap's Drive checklist runs the pre-cut audit *after* every member is merged. This
   review therefore covers only the recorded gate lifecycle surface that has landed. A further audit
   is owed once the remaining members merge.
2. **Live-lane coverage is incomplete on the audited commit.** The final tag candidate still needs
   a fresh exercise, but the last Pi result does not establish a new product defect. See the
   release-prerequisite section.

## Verdict: NOT READY

The recorded gate lifecycle itself is sound: one write surface holds, tests that run the real binary
prove the frozen-closure and one-use-application invariants, and every checkable digest is byte-exact.
The *release* is not ready. The Contract landing pass has not run, the final tag candidate lacks a
fresh live-lane exercise, and the schema carries a dead `application` subtree that becomes expensive
to remove after a strict-decoding v1 tag. The latest Pi lane completed the durable lifecycle before a
captain-deferred review-text diagnostic failed; that result requires an explicit disposition, not
another workflow correction.

Findings: 2 material, 4 minor, 1 nit.

## Release prerequisites

### 1. The Contract landing pass has not run; the landed spec still carries task-id owner tags

The Drive checklist item at `docs/roadmap/durable-decisions/index.md:57` is unchecked. Its scope is
explicit: "strip the spec's shaping scaffolding — owner-tag lines, diagram task-id prefixes …, example
ids genericized — so the landed spec speaks only component terms."

On `deac7f8a`, `docs/specs/gate-resolution-frontmatter-contract.md:126` still reads:

```
## Round records and triage dispositions (advisory; owner: 02av)
```

with three further `02av` occurrences at lines 144, 152 and 153 (the example Annotation and Resolution
ids). `02av` is a workflow task-id prefix, not a component term: it resolves to
`_archive/ensign-finding-triage-disposition` (entity id `02avdajaz0q3hnjwycm5fq45`) in the state
checkout. A landed, shipped product spec therefore points a reader at an archived internal task
reference. This is precisely the class the landing pass exists to remove, and the pass is the recorder
member's declared final step.

**Blocks the cut.** Run the pass before the tag. The edit is four lines in one spec file.

### 2. The final tag candidate needs a fresh live-lane exercise

`gh api repos/spacedock-dev/spacedock/commits/deac7f8a.../check-runs` returns `total_count: 1` — a
single `build` check, `completed / success`. No live lane has run on `main`'s tip.

On the branch merged here, the final lane results were:

| lane | status | conclusion |
|---|---|---|
| `pi-live` | completed | **failure** |
| `claude-live (sonnet, CI-E2E)` | waiting | never ran |
| `claude-live (claude-opus-4-8, CI-E2E-OPUS)` | waiting | never ran |
| `codex-live` | waiting | never ran |
| `offline` | completed | success |
| `build` | completed | success |
| `install (ubuntu-latest)` | completed | success |
| `install (macos-latest)` | completed | success |

Three lanes produced no result in either direction. The Pi lane produced a red check conclusion, but
the retained evidence separates the lifecycle outcome from the later diagnostic:

- Commit `e617f947` already landed the two-file, `+2/-2` selector correction that made the canonical
  Pi job invoke `TestLivePiRecordedGateLifecycle`. It is an ancestor of the audited commit, so no such
  workflow correction remains to land.
- Commit `cc5b27dc` removed the forced provider/model pair. At that exact tip, Pi bound the Briefing,
  recorded approval, consumed the gate, advanced to handoff, dispatched the successor, and retained
  successor commit `d630e06`.
- The check failed only after those outcomes, when review-text extraction reported
  `gate review omits its decision facts`. The final 6y validation classified that diagnostic as a
  **captain-deferred Pi evidence gap**. Its promotion condition is a supported Pi journey failing to
  bind, record, consume, advance, dispatch, or commit its durable effect, or a later decision that
  restores decision-facts extraction as an explicit release requirement.

**Required before the tag.** Exercise every applicable live lane on the exact tag candidate. Count
each lane by its retained lifecycle evidence, not by absence of a red check. A lane must either pass
or carry an explicit captain disposition for the same narrow, post-outcome diagnostic. A different
lifecycle failure remains release-blocking.

---

## Findings

### 1. The `application` subtree ships 13 dead leaf fields into a schema with no migration path

**Material. DELETE NOW.**

**Mechanism.** `internal/gates/model.go:45-51` declares six `Application` fields. Three are
load-bearing; three are dead subtrees:

```go
Blockers      *[]Blocker     `yaml:"blockers,omitempty" ...`       // model.go:49
ExecutionHold *ExecutionHold `yaml:"execution-hold,omitempty" ...` // model.go:50
Feedback      *Feedback      `yaml:"feedback,omitempty" ...`       // model.go:51
```

Their leaf types are `Blocker` (`model.go:54-61`, 6 fields), `ExecutionHold` (`:63-69`, 5 fields) and
`Feedback` (`:71-75`, 3 fields). Counting production writers and readers, excluding `_test.go`: 13 of
the 17 scalar leaves in the subtree have **zero writers and zero readers**. `action`, `target-stage`
and `state` carry the whole live lifecycle — 12 production writers, 25+ production readers, three
entities using them today.

`blockers` is a distinct pathology: proof that cannot fail. Its only production writer is
`internal/gates/operation.go:483-484`, which unconditionally writes the empty literal:

```go
blockers := []Blocker{}
return &Application{Action: "advance", TargetStage: stages[i+1].Name, State: "pending", Blockers: &blockers}, nil
```

Under the sprint's one-write-surface rule (`index.md:74`, rule 3), no reachable state makes
`len(*app.Blockers) != 0` true — so the `blocked` condition at `model.go:212` and
`application.go:76-78` is dead by construction. `execution-hold` has no writer anywhere in the
repository, so the `held` branches at `model.go:203-204` and `application.go:66-67` are unreachable
too. `feedback` has neither writer nor reader.

**Evidence.**
- Full-tree `grep -rIn`: `expected-revision` — 1 hit, its own declaration at `model.go:58`.
  `expected-state` — 1 hit, `model.go:59`. `finding-ref` / `finding-digest` — 2 hits each, the
  declaration plus one round-trip fixture string. `execution-hold` — 4 hits: declaration, two test
  fixtures, and the spec.
- `grep -rn "ExecutionHold{" / "Feedback{" --include=*.go` — zero production hits. `grep -rn
  "Application{"` — 3 production sites (`operation.go:457,478,484`); only `:484` sets `Blockers`.
- `grep -rIn "blockers" hooks/ mods/ agents/` — zero. No out-of-band writer, no JSON/YAML schema file.
- The spec concedes it: `docs/specs/gate-resolution-frontmatter-contract.md:232` lists
  "Blocker-satisfaction evaluation, execution-hold authoring, dispatch identities, or effect receipts"
  under **Explicitly outside v1** — while `:52` puts `blockers: []` in the canonical frontmatter
  example.
- The governing ruling is already recorded: `index.md:17` — "Blockers/eligibility ship only against a
  demonstrated live consumer; absent one, h1 closes as a recorded decline, not a build."
  `staff-review-codex.md:217` records the evidence for that decline: "no dry-run gate carried a
  blocker or execution hold." The decline was taken on the evaluator; the schema footprint shipped.
- Deletion experiment on a scratchpad export of the pinned SHA (pinned tree read-only): removed the
  three types, three fields, four reader branches, the writer, and the fixtures that exist only to
  exercise them. `go build ./...` OK, `go vet ./...` clean, `gofmt` clean, and
  `go test ./internal/gates/ ./internal/status/ ./internal/cli/ ./internal/ensigncycle/` all `ok`.
  Nothing load-bearing failed.

**Live-state corroboration** (from `docs/dev/.spacedock-state/`, a separate checkout — live-state
evidence, not pinned-tree evidence). Three sprint members carry a closed approve resolution with
`application: state: pending` and `blockers: []`: `prepare-provider-neutral-gate-room`,
`sync-merge-guard-archive-state`, and `live-claude-runner-scrub-foreign-runtime-markers`. All three
have a real, named, identical dependency, and all three record it only in the resolution's free-text
`reason` ("implementation remains pending until 6y lands", "dependency-held behind 6y", "must wait for
landed 6y"). The boot readiness projection emitted `approved-awaiting-advance` with `dispatchable: []`
both while that dependency was unmet and after it was satisfied — identical output in both states, so
the projection carries no information about the dependency. A tree-wide scan of live state for a
non-empty `blockers`, any `execution-hold`, or any `feedback:` under an application returned zero
hits. The machine-readable field is present at the exact moment a genuine blocker exists and does not
hold it.

**Why this is time-critical.** `internal/gates/io.go:38-40` sets `decoder.KnownFields(true)`, and the
spec states there is no migration or compatibility rewrite. Deleting these fields *after* the tag
makes every entity carrying `blockers: []` fail to decode, on data the team does not control, in a v1
that by design has no migration path. Before the tag it is a mechanical line-strip in a state checkout
the team owns.

**Simpler alternative.** Delete the `Blocker`, `ExecutionHold` and `Feedback` types and their three
`Application` fields. Keep `action`, `target-stage`, `state`. Simplify the readiness guard at
`model.go:199-214` to `if attempt.Resolution.Decision != "approve" || app.Action != "advance" ||
strings.TrimSpace(app.TargetStage) == "" || app.TargetStage == status { return "invalid" }`. Drop the
`app.Blockers == nil` clause at `application.go:63` and the hold/blocker branches at `:66-79`. The
eligibility vocabulary loses `blocked` and `held`, which no production writer can produce; every
condition a real gate reaches is untouched. The promotion condition for a real blocker layer is
already written down at `index.md:17`, and whatever eventually replaces it should be designed from the
observed usage — a free-text dependency `reason` that has so far been sufficient — not from the shape
guessed here.

**Deletion cost.** Go: -121/+15 across 9 files (production -51, tests -55). 16 struct fields and 3
types removed. Prose: 4 spots in 2 files — spec `:52` (drop `blockers: []` from the canonical
example), spec `:232` (the out-of-v1 bullet reduces to dispatch identities and effect receipts, which
makes the section honest), `skills/fo-gate-lifecycle/SKILL.md:13` (drop `blocked/held` from the
omitted-readiness list) and `:50` (drop `blockers` from "Consume itself rechecks currency, successor,
blockers, and one-use state"). No section is deleted; no section changes ownership. Plus the
mechanical `blockers: []` line-strip in the state checkout.

### 2. The Codex live runner's named archival assertion is tautological

**Minor. NARROW.**

**Mechanism.** `internal/ensigncycle/codex_live_runner_test.go:210` loads `after` from the *active*
entity path, and `:217` then asserts:

```go
after := readFile(t, fixture.entity)               // :210
...
if resolveRecordedGateEntity(fixture) != after {   // :217
    t.Fatal("recorded-gate-task was archived while waiting at the gate")
}
```

`resolveRecordedGateEntity` (`internal/ensigncycle/recorded_gate_lifecycle_test.go:873-884`) tries
`fixture.entity` **first**, then `_archive/recorded-gate-task/index.md`, then
`_archive/recorded-gate-task.md`, returning file *contents*. The two paths behave as follows:

- Not archived — it re-reads the same file the test just read. The comparison is a tautology.
- Archived — the initial `readFile(t, fixture.entity)` calls `t.Fatal` before the comparison, so the
  test fails with a generic read error. The named assertion never runs.

The test therefore detects archival, but the assertion that claims to detect it cannot fire. This is
a diagnostic and proof-clarity defect, not missing outcome coverage.

**Evidence.** `git grep -n resolveRecordedGateEntity -- internal/` at `deac7f8a` returns exactly three
lines: `codex_live_runner_test.go:217` (the defect), `recorded_gate_lifecycle_test.go:851` (correct
usage — `after := resolveRecordedGateEntity(fixture)`, assigned rather than compared, where the
archive fallback is load-bearing), and the definition at `:873`. This is the only live-runner call
site. The sibling merge-guard scenario in the same file gets it right at
`codex_live_runner_test.go:293`:

```go
if _, err := os.Stat(filepath.Join(workflowRoot, "_archive", "merge-check.md")); !os.IsNotExist(err) {
```

which is falsifiable. The Claude gate-guardrail scenario
(`internal/ensigncycle/claude_live_runner_test.go:253-275`) carries no archive check at all.

**Simpler alternative.** Before reading the active entity, replace `:217` with the explicit
archive-absence check already used at `:293`. Add the same pre-read diagnostic to the Claude runner
for parity. No new harness or detector expansion is needed.

**Deletion cost.** Two test files, a few lines each. Nothing in production changes.

**Why it matters.** The current test fails on archival, but it attributes the failure to a file read
rather than the guardrail under test. A direct check makes the proof legible and gives both live
runners the same diagnostic.

### 3. A lifecycle test greps the shipped skill's prose from outside the quarantine, and the grep is a tautology

**Material. DELETE NOW.**

**Mechanism.** `internal/ensigncycle/recorded_gate_lifecycle_test.go:568-570`:

```go
func TestRecordedGateLifecycleCommandTextMutants(t *testing.T) {
	root := recordedGateRepoRoot(t)
	original := readFile(t, filepath.Join(root, "skills", "fo-gate-lifecycle", "SKILL.md"))
```

`recordedGateRepoRoot` (`:1130`) resolves the real repo root, so this reads the **shipped** skill body,
not a fixture. `internal/contractlint/boundary_guard_test.go:49` states the rule without
qualification: "A read outside the quarantine is a failure, full stop."

The assertion is a tautology. The mutant literals at `:576-578` each strictly *contain* the needle the
test-local parser `procedureEvents` (`:650-667`) greps for:

| deleted by the test | grepped by `procedureEvents` |
|---|---|
| `gate record ENTITY --briefing BRIEFING` | `gate record ENTITY --briefing` |
| `gate record ENTITY --decision approve\|revise\|hold` | `gate record ENTITY --decision` |
| `gate consume ENTITY` | `gate consume ENTITY` (identical) |

Deleting a superstring of the needle cannot leave the needle. The loop then asserts the parser misses
it (`:592-593`). No product code runs. The baseline at `:571` compares the shipped prose against three
string literals hardcoded at `:576-578` — neither side is the product. Renaming `gate consume` in
`internal/cli` would leave this test green.

**Evidence.**
- Demonstrated non-discrimination, on a scratch copy (the audit tree untouched): rewrote the shipped
  skill to invert its meaning while keeping every grepped token — `gate consume ENTITY …` became
  `NEVER run: … gate consume ENTITY … # dispatch WITHOUT consuming`, and the briefing record became
  `Skip this step entirely; do NOT run … gate record ENTITY --briefing BRIEFING …`. Result:
  `ok github.com/spacedock-dev/spacedock/internal/ensigncycle` — still green. That is verbatim the
  failure mode the guard doc names at `boundary_guard_test.go:31-34`: "a meaning-inverting paraphrase
  keeps every grepped token."
- The guard is not defective. It is documented not to see this shape:
  `internal/contractlint/instruction_read_detector_test.go:42-46` — "no data-flow tracking. A read
  reached only through a param/local/field/method/closure flow, a transitive helper chain, a
  range-element flow, a cross-package reader … is NOT detected here; that reader-shape axis is the
  detached-audit-backstopped boundary", repeated at `boundary_guard_test.go:42-44`. Confirmed by
  running the pinned detector on both shapes: as written, through the package-local `readFile` helper,
  `detected=false`; the identical read written as `os.ReadFile`, `detected=true`.
- `procedureEvents` has exactly two call sites, `:571` and `:592`, both inside the test under review.
  `grep -rn "CommandTextMutants\|procedureEvents" docs/` — no matches; no spec depends on it.
- Deletion verified clean on a scratch copy: `go vet ./internal/ensigncycle/` exit 0 with no output;
  `go test ./internal/ensigncycle/ -run TestRecordedGateLifecycle -count=1` → `ok` (32.6s).

**What already proves this behavior, by running it.** All in the same package on the pinned tree:
`TestRecordedGateLifecycleRealCLIReplay` (`:144`) drives the real binary through record/record/consume
with a durable `gates.Read` snapshot, a command-log trace, and five negative controls;
`TestRecordedGateLifecycleMissingEventControls` (`:718`) omits each of the three events and requires
dispatch authorization to refuse; `TestRecordedGateLifecycleAC5RefusalMatrix` (`:277`) and
`AC7ResumeMatrix` (`:433`) cover the refusal and resume shapes. Three live journeys grade the command
log an agent actually produces from the shipped skill: `claude_live_runner_test.go:173`,
`codex_live_runner_test.go:88`, and `recorded_gate_lifecycle_pi_live_test.go:16` (which loads
`skills/fo-gate-lifecycle` explicitly at `:50`). Structural checks on this exact file already live
legally in quarantine: `internal/contractlint/structural_checks_test.go:146`,
`internal/contractlint/fo_function_reference_invariant_test.go:137` (byte cap) and `:277`.

**Simpler alternative.** Delete the test (`:568-597`) and the now-orphaned `procedureEvents`
(`:650-667`). No replacement is needed. `recordedGateRepoRoot` stays (used at `:902`, `:957`);
`readFile` stays.

**Deletion cost.** 48 lines, one file, no cross-file edits.

**Two adjacent observations, stated once and not pursued further.** First, the guard's claim at
`boundary_guard_test.go:59-60` that "the count of out-of-quarantine reader files must be zero" is not
literally true of the tree it polices: `skills/integration/survey_probe_test.go:30` and `:123` read
`skills/survey/SKILL.md`, and `internal/cli/prose_function_routing_test.go:129` reads the three FO
cores declared at `:20-24`. Both evade through different documented blind spots (local-variable and
range-element flow), and both **extract and execute** the shipped artifact, so they are defensible in
spirit where the one above is not. Either move them or reword the guard's claim — do not widen the
detector, which is an explicit design choice. Second,
`internal/contractlint/layering_restore_test.go:188` maps `skills/fo-gate-lifecycle/SKILL.md` to the
prose string "terminal current status resumes the existing merge ceremony" — an in-quarantine
prose-grep, the form the package doc bans "here and everywhere" (`boundary_guard_test.go:29-34`).
Worth a separate look; not verified further here.

### 4. Boot-resident FO prose restates the deferred file it points at, and the largest presenter file is uncapped

**Minor. NARROW.**

**Mechanism.** `skills/first-officer/references/first-officer-shared-core.md` is boot-resident, so
every byte is a standing per-session cost. Line 46 correctly declares the deferred trigger for
`Skill(skill="spacedock:fo-gate-lifecycle")`. Line 81 then restates the routed semantics the deferred
file itself carries — "nonterminal → ordinary dispatch, terminal → existing merge ceremony; revise
invokes `«feedback.route»`, while hold/ineligibility stops" — against
`skills/fo-gate-lifecycle/SKILL.md:50`, "A nonterminal target then enters ordinary reuse-or-fresh
dispatch; a terminal target enters the existing merge guard/hook and has no successor dispatch."
Boot-resident bytes paying to summarize a file that is loaded before any gate action anyway.

The new deferred file does earn its place: `internal/contractlint/structural_checks_test.go:146`
asserts it is named, non-user-invocable, absent from user command discovery and adapter-less, and its
trigger set is a strict superset of the presenter's, so pure resume and consume paths avoid loading
the presenter.

**Evidence.** Merged sizes on `deac7f8a` (`wc -c`): `skills/present-gate/SKILL.md` 9,285;
`skills/fo-gate-lifecycle/SKILL.md` 6,596; `first-officer-shared-core.md` 26,706;
`claude-first-officer-runtime.md` 4,796. Across the change that introduced the deferred file (base
`c9633279`), `present-gate` grew 8,616 → 9,285 (+669) and `claude-first-officer-runtime` grew
4,611 → 4,796 (+185).

`internal/contractlint/fo_function_reference_invariant_test.go:137-140` caps exactly two files:
`first-officer-shared-core.md` at 26,754 (48 bytes of headroom) and `fo-gate-lifecycle/SKILL.md` at
6,600 (4 bytes of headroom). `git grep "present-gate" internal/contractlint/` shows no cap entry — the
largest FO instruction file on disk, and the one that grew most in this change, is uncapped.
`foHostLoadBytes` survives at `:85` but is consumed only by a `Logf` argument at `:172`, so the
per-host aggregate is reported and not asserted.

**Recommendation.** Fold the restated routing at `first-officer-shared-core.md:81` down to the
trigger, letting `:46` and the deferred file carry the semantics. Do **not** add a presenter cap in
this cleanup. The captain closed the aggregate prompt-load ratchet by replacing it with selected
component caps; extending that set would create a new standing obligation. A separate, evidence-backed
proposal may request that authority. The prose trim is small and non-release-blocking.

**Honesty note on the measurement.** At the pinned baseline (`c9633279`, the parent of the commit
that added the deferred file), `first-officer-shared-core.md` net *shrank* 32,289 → 26,706. This
audit could
not reproduce a net boot-resident increase across this change. The finding is the duplication and the
uncapped presenter, both directly observable on the pinned tree; only the duplication has an approved
cleanup. The finding is not aggregate growth.

### 5. The recorder derives a Result association and then re-verifies it against itself

**Minor. DELETE NOW.**

**Mechanism.** `recordRoomLocked` builds a struct and immediately re-reads it,
`internal/gates/operation.go:328-334`:

```go
association, err := deriveAssociation(resultBytes, &result, &presented, request.Approver, attempt.Briefing, canonicalItems)
if err != nil { return err }
if err := verifyAssociation(resultBytes, &result, association, request.Approver, attempt.Briefing, canonicalItems); err != nil { return err }
```

The two calls share five of six arguments verbatim; the sixth differs only in that derive takes
`&presented` and verify takes derive's own return value. `verifyAssociation` (`:534-582`) has no input
derive did not either receive or produce. The struct it re-reads (`resultAssociation`, `:74-91`) is
never persisted — what the recorder actually freezes is two raw digests at `:337-340`. The spec agrees
it is ephemeral (`docs/specs/gate-resolution-frontmatter-contract.md:124`).

Of the nine error returns in `verifyAssociation`, **seven are unreachable** and one more is
half-unreachable, each guaranteed by the derive line that precedes it:

- `:538` (Type/Version) — derive sets both as literals at `:595`.
- `:541` (Result.Digest/Briefing) — derive sets both from the same `resultBytes`/`result` at
  `:596-597`.
- `:544` — `association.Actor != actor` is a tautology (derive assigns it at `:595`), and
  `result.Resolution.By != actor` is already enforced at `:309` on the same value.
- `:550` (Canonical.Briefing/Revision) — derive sets both from the same `binding` at `:598-599`.
- `:553` (three length equalities) — derive errors at `:585` on the same condition and appends exactly
  one entry per element.
- `:559` (canonical artifact loop) — byte-for-byte the predicates derive already ran at `:590`.
- `:569` (presentation loop) — guaranteed by derive's equality checks and its dedupe.
- `:578` first clause — derive already forced a bijection.

Three checks are live and genuinely external: `:535` (the `review-v1-result` envelope; derive never
looks at `result.Type` or `result.Artifact`), `:547` (`verifyProviderResolution`, which touches only
`result` and has nothing to do with the association), and `:578`'s `!resultArtifactPresent`.

**Evidence.**
- `grep -rn "resultAssociation\|deriveAssociation\|verifyAssociation" --include=*.go` → 6 hits, all in
  `operation.go`. `grep -n Marshal internal/gates/operation.go` → only `yaml.Marshal` in `nodesEqual`.
  Nothing serializes the struct.
- Runtime probe: instrumented all nine `return fmt.Errorf` statements plus the
  `verifyProviderResolution` delegation and ran `go test ./...` on a pinned export. All 16 packages
  pass; exactly two probes fire, both on live checks. The seven never fired once.
- Property test: 300,000 randomized input tuples into `deriveAssociation`; 9,399 returned nil error;
  for every one, each association-referencing clause of the seven held. Zero fires.
- Deletion experiment: deleted `verifyAssociation` and the struct, collapsed `deriveAssociation` into
  a `checkPresentationMapping` validator, hoisted the live checks to the call site. `gofmt` clean,
  `go build ./...` OK, `go test ./...` all 16 packages `ok`, including every `TestGateRoom*` rejection
  subtest. `diff -rq` over the full tree: exactly one source file differs.
- No test asserts on any of the seven dead error strings.

**Simpler alternative.** Delete `verifyAssociation` (`:534-582`) and `resultAssociation` (`:74-91`).
Turn `deriveAssociation` into `checkPresentationMapping(result, presented, inventory) error` — the
same loop with no struct built — folding the one live association requirement into the loop it already
runs:

```go
if expected.Type == "Artifact" && item.ID == result.Artifact.ID && item.Rev == result.Artifact.Rev {
    resultArtifactPresent = true
}
```

Hoist the envelope check and `verifyProviderResolution` to the call site alongside the other
`result.json` checks at `:306-311`, and drop the tautological and duplicate actor clauses — `:309`
already owns that predicate. One caveat: this changes error *precedence* (the envelope message fires
before the presented-inventory message). No test depends on the ordering; keep the mapping check first
if byte-identical precedence is wanted.

**Deletion cost.** Net -75 LOC in one file (`operation.go` 925 → 850; `1 file changed, 17 insertions,
92 deletions`). Zero test files touched. Zero spec change for the verify pass alone — the spec's
requirements at `:110-111` and `:246` are both preserved by the folded check. Deleting the struct as
well makes two spec sentences stale (`:109`, `:124`) and needs rewording; that is a normative edit, so
it does not ride the landing pass.

**Severity rationale.** Nothing is persisted, so no durable state, migration, or frontmatter byte can
be affected, and the seven dead checks cannot wrongly reject either. The cost is comprehension: 75
lines that read like the authorization proof of the sprint's most security-relevant path and provably
cannot reject anything a reviewer would hope they catch. Do it before the tag, while the v1 schema is
unreleased and the change is one file with no test or spec churn.

### 6. `gates.Summary` declares two eligibility fields no producer ever sets

**Minor. DELETE NOW.**

**Mechanism.** `internal/gates/model.go:114-126` declares eleven `Summary` fields. Two have no
producer anywhere in the tree:

```go
Condition        string   // model.go:123
Eligible         bool     // model.go:124
```

The sole constructor of a non-empty `Summary` is `CurrentSummary` (`model.go:341-360`), which writes
Gate/Attempt/State/Briefing (`:348`), Resolution/Decision (`:350`), and
Application/ApplicationState/TargetStage (`:353-355`) — never Condition or Eligible. The only other
`Summary` literals are the zero values at `model.go:359` and `io.go:52,55`.

Both are read exactly once each, at `internal/status/discover.go:227-228`:

```go
fields["gate-condition"] = summary.Condition
fields["gate-eligible"] = fmt.Sprintf("%t", summary.Eligible)
```

Line 227 assigns the constant `""`. Line 228 `fmt.Sprintf`s a value that is provably always `false`.
The real values come from a different path — `gates.EligibilityFileAt`
(`internal/gates/application.go:105`) via `EvaluateEligibility`, populating the *different* struct
`gates.Eligibility` (`model.go:127-135`, whose `Condition`/`Eligible` at `:134-135` are the real ones)
— and are written into the same two map keys 75 lines later at `discover.go:303-304`.

**Evidence.**
- `grep -rn 'Condition:\|Eligible:' --include=*.go .` → exactly one keyed init tree-wide:
  `internal/gates/application.go:19`, `result := Eligibility{Condition: "ineligible"}`. Not a
  `Summary`. `grep -rn 'Summary{'` → the four sites above, none keyed with either field.
  `grep -rn 'gate-eligible\|gate-condition' . --exclude-dir=.git` → 9 hits, all Go; no skill, doc, or
  fixture consumer.
- Fail-closed check, run live. Built the pinned tree to scratchpad and ran
  `spacedock status --fields id,gate-state,gate-condition,gate-eligible --json` on three fixtures:

  ```
  {"id":"plain","gate-state":"","gate-condition":"","gate-eligible":""}
  {"id":"withstatus","gate-state":"closed","gate-condition":"stale","gate-eligible":"false"}
  {"id":"nostatus","gate-state":"closed","gate-condition":"","gate-eligible":"false"}
  ```

  `nostatus` (valid gates, `status:` key removed) proves the seed is *observable*: `gates.Read`
  succeeds so `:227-228` run, then `EligibilityFileAt` fails at `entityStatus` (`io.go:97-99`), so
  `materializeGateEligibility` `continue`s at `:301` and the seed survives. `plain` proves `:227` is a
  pure no-op — no key set at all, identical `""` output, because every consumer reads
  `e.fields[name]` by map index.
- Deletion proven on a scratchpad copy: removed `model.go:123-124`, deleted `discover.go:227`, rewrote
  `:228` to the literal. `gofmt` clean, `go build ./...` OK, `go test ./internal/gates/...
  ./internal/status/... ./internal/cli/...` all `ok`, and the rebuilt binary's JSON on all three
  fixtures is byte-identical to the baseline.

**Simpler alternative.** Delete `Condition` and `Eligible` from `Summary` (`model.go:123-124`). Delete
`discover.go:227` outright — a verified no-op. Replace `:228` with the literal it always evaluates to:

```go
fields["gate-eligible"] = "false"
```

That one line must stay: it is the fail-closed default for the reachable case of an entity whose gates
record parses but whose `status:` key is missing. Deleting it too would flip that path from "not
eligible" to "unknown", against the fail-closed posture declared at `internal/gates/application.go:1`.
A literal states that intent far more plainly than routing a constant through a struct field no
producer sets.

**Deletion cost.** 2 files, net -3 LOC. No import churn (`fmt` still used at `discover.go:304`). No
test edits.

**Related, out of scope here.** The fail-closed seed is not universal: the whole block sits inside
`if gateErr == nil` (`discover.go:217`), so an entity with no gates record renders
`"gate-eligible":""` (verified on fixture `plain`). A uniform floor would put `gate-eligible` in
`defaultEntityKeys` (`discover.go:149`) with a `"false"` default — a behavior change, not a cleanup.

### 7. Two recorder entry points have no production caller

**Nit. ACCEPT-WITH-REASON.**

`internal/gates/operation.go:149-151` is a one-line delegation:

```go
func RecordBriefing(entityPath, briefingPath string) error {
	return RecordSemantic(entityPath, RecordInput{BriefingPath: briefingPath})
}
```

`git grep -n "RecordBriefing(" -- '*.go' | grep -v _test.go` returns only the declaration; 13 test
call sites use it. Likewise `internal/gates/application.go:85` `func EligibilityFile` has exactly one
caller in the whole tree, `EligibilityFileAt` at `:106`.

Both are thin, both are cheap, and both are genuinely useful as test-facing seams. Accept them, but do
not let the pattern grow: an exported entry point whose only callers are tests is a maintenance
surface, not an API.

---

## Disposition and ownership

Keep each correction within the component that owns its invariant:

1. **6y proof hygiene.** Delete `TestRecordedGateLifecycleCommandTextMutants` and
   `procedureEvents`; improve the Codex archive diagnostic and add the same pre-read diagnostic to
   Claude. Inspect the adjacent in-quarantine prose grep separately. These changes correct proof
   quality without changing the lifecycle.
2. **Application schema cleanup.** Remove the dead `Blocker`, `ExecutionHold`, and `Feedback`
   subtrees before the first v1 tag, then strip `blockers: []` mechanically from the owned state
   checkout. This is an h1/schema follow-up, not 6y implementation work. The two dead `Summary`
   fields may join that cleanup if doing so does not broaden its behavioral boundary.
3. **Recorder cleanup.** Route the ephemeral association self-verification to the recorder owner.
   It is a separate, non-blocking simplification.
4. **Contract landing pass.** Genericize the four remaining `02av` references in the product spec.
   This existing checklist item remains a release prerequisite.
5. **Durability probe.** Drive a multi-record, multi-attempt approval history through the real
   `status --set` binary during the final audit. Fix code only if that falsifiable probe shows a
   history mutation.
6. **Final audit.** Merge the four remaining sprint members, exercise the exact tag candidate on all
   applicable live lanes, and repeat the detached adversarial audit before recording sprint closure.

---

## Suite and lane evidence

All commands run on the pinned checkout. Toolchain: go1.26.1 darwin/arm64. Worktree clean before and
after; `git rev-parse HEAD` re-checked after the last run and unchanged.

| command | result |
|---|---|
| `go test ./internal/gates ./internal/status ./internal/ensigncycle ./internal/contractlint` | PASS, exit 0. 4/4 `ok`. |
| `go test ./...` | PASS, exit 0. 16/16 test packages `ok`; 2 packages have no test files. |
| `go test ./... -race` | **Completed** (~5.5 min wall). PASS, exit 0. 16/16 `ok`, 0 data races. |
| `gofmt -l ./cmd ./internal` | PASS, exit 0, zero lines of output. |

Test-level tally from `go test ./... -json -count=1` (cache bypassed): **2401 pass, 0 fail, 8 skip.**
No failing test in any package. Every skip adjudicated individually from its emitted message:

- Seven are environment gates that would run in CI: five `internal/release` CI-log tests
  (`cilog_clean_output_test.go:91,111,134,162,187` — "gotestsum not on PATH"), `internal/cli`
  `TestCodexResolveManifestAgainstInstalledHost` (`codex_resolve_test.go:46`), and
  `skills/integration` `TestSurveyCodexPresenceThroughSync` (`survey_sync_codex_test.go:55` —
  "agentsview not on PATH").
- One is parked unconditionally: `internal/status` `TestUpdateFrontmatterBlockScalarsRewrapped`
  (`node_roundtrip_test.go:270-272`) is a function whose entire body is
  `t.Skip("documented divergence #L2-b — no live block-scalar frontmatter to pin; reactivate when a
  fixture exists")` — zero assertions, so it can never fail. Its comment states the intent, and it is
  outside this sprint's surface; noted, not filed.

No test errored on a missing runtime or credential. The credential-gated live runners did not skip —
their guards resolved without needing to.

**Lane diagnosability.** The Pi live lane's provider and model are not pinned:
`git grep SPACEDOCK_PI_LIVE_PROVIDER` and `SPACEDOCK_PI_LIVE_MODEL` both return nothing, and
`internal/ensigncycle/recorded_gate_lifecycle_pi_live_test.go` carries no `--provider` or `--model`
flag — the only knob is `SPACEDOCK_PI_LIVE_CHILD_MODEL` at `:50`, which stamps the *successor*
dispatch. Commit `cc5b27dc` made that inheritance intentional because the forced pair could not use
the lane's available authentication. Do not restore the pin as an inferred obligation. Instead,
retain the resolved provider and model in the run evidence when the host exposes them; absence of that
diagnostic does not invalidate a durable lifecycle result.

## Boundaries that held

- **One write surface.** Every recorder write in `internal/gates/operation.go` (`:184`, `:209`,
  `:230`, `:344`, `:385`) goes through `writeDocument`, which is `internal/gates/io.go:103-105` →
  `writeEntityDocument(path, expected, nil, doc, nil)` — status always nil. The only `status:` writer
  in the package is `internal/gates/application.go:279`, the consume/apply path. No second `gates:` or
  `status:` writer exists.
- **Unrelated frontmatter mutation does not corrupt gates.** `internal/status/gates_coexist_test.go:159`,
  `TestUnrelatedSetPreservesGatesAndStatusProjectsResolution`, writes an entity with a full nested
  `gates` block, snapshots the entire parsed value, mutates an unrelated field through
  `updateFrontmatter`, and asserts `reflect.DeepEqual` on the whole value. Covered. One residual,
  stated honestly: it exercises `updateFrontmatter` rather than the `status --set` CLI end to end, on
  a single-record single-attempt `hold` shape rather than a multi-record approval history.
- **The provider seam is honestly cheap.** `git grep "interface {" -- internal/gates/` returns
  nothing — zero Go interfaces in the seam.
- **`fo-gate-lifecycle` earns its file.** `internal/contractlint/structural_checks_test.go:146`
  asserts it is named, non-user-invocable, absent from user command discovery and adapter-less;
  `internal/contractlint/fo_function_reference_invariant_test.go:277` reads lifecycle and presenter
  together and asserts the lifecycle carries the fail-closed and presentation contract clauses the
  presenter does not.
- **No hidden schema versioning.** No revision counter surfaced anywhere in the v1 schema.

## Examined and not substantiated

- **Whether `status --set` can destroy an approval history at exit 0.** Not established either way.
  The covering test named above is real but narrower than the claim: `updateFrontmatter` rather than
  the CLI, and a single-record single-attempt shape rather than an approval history. Settling this
  needs its own falsifiable probe — a multi-record, multi-attempt approval history driven through the
  real `status --set` binary — because it is the only open question in this set that would be a
  durability defect rather than an over-build.
- **The per-host aggregate FO prompt load on the audited commit.** `foHostLoadBytes`
  (`fo_function_reference_invariant_test.go:85`) is computed but flows only into a `Logf` at `:172`.
  No check asserts it. See the closed ruling below.

## Closed by prior ruling — not reopened

- **The writer-less `raw-file-pin` digest domain is deliberate.** Recorded at `index.md:26` as a
  captain ruling, and stated normatively in the spec at
  `docs/specs/gate-resolution-frontmatter-contract.md:70-73`: "`canonical-bytes`: … New recorder binds
  always use this domain" / "`raw-file-pin`: an explicitly labelled raw-byte pin … never silently
  reinterpreted as a canonical digest." It is accepted at `internal/gates/model.go:254` and evaluated
  at `internal/gates/application.go:231`. A read-only domain with no writer is the ordered mechanism,
  not dead capacity. **Closed.**
- **The aggregate prompt-load ratchet, replaced by per-component caps.** Reported as a captain ruling
and treated here as closed rather than a defect. One line of honesty: the ruling's citation was not
located during this audit, and as a consequence the per-host aggregate load on `deac7f8a` is not
asserted by any check (see above). That absence does not authorize a new `present-gate` cap.
**Closed.**

## What a completed pre-cut audit still owes

Once `rq`, `s4`, `rd` and `v3v` merge, the roadmap's pre-cut audit (`index.md:58`) still owes, at a
minimum:

1. The same detached adversarial pass over each of the four members' landed surfaces, at the SHA that
   will actually be tagged.
2. Live-lane evidence on that SHA — all applicable lanes run, with each retained lifecycle outcome
   green or an exact captain disposition for the already deferred post-outcome Pi diagnostic.
3. Re-verification that the landing pass is complete and that no different lifecycle failure hides
   behind a check conclusion on the tagged tree.
4. Confirmation that the `application` subtree decision in finding 1 was taken *before* the tag, since
   `KnownFields(true)` makes it substantially more expensive afterwards.
5. A settled answer on the `status --set` approval-history question above, or an explicit decision to
   ship without one.

*Audited on a detached throwaway checkout at `deac7f8af9b0f56ceea6231976d7ffaa8ab2fc51`, re-checked
unchanged at the end of the audit; every `file:line` verified there; nothing mutated.*
