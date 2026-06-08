# 0199 pre-flip-mechanics — preflight staff review

Independent staff review run BEFORE the Commander drive. Mandate: refute the shaping
FO's gate decisions, not ratify them. Read-only over the package + the three entity
bodies + the code/config each claim rests on. The verdicts below are mine, not the FO's.

## Verdict per task

- **`v3` ship-linux-binaries** — `sound`. The riskiest mechanism (goreleaser linux
  cross-compile) was genuinely exercised; the config guard and install.sh tamper-case
  are real independent checks; the one weak AC is honestly self-flagged as non-load-bearing.
- **`th` safehouse-preserves-spacedock-bin** — `sound-with-conditions`: the argv-oracle
  half is testable and honestly scoped as "shape, not survival," and the captain
  real-safehouse run is a bounded plan — BUT the gate semantics must be enforced (no
  `done`/PASSED on shape+fake-smoke alone) and the captain check must cover a SECOND
  untested assumption the body underplays (bare `env` resolves on the sandbox PATH).
- **`jm` entity-label-localization** — `MATERIAL HOLE`: AC-1 and AC-2 localize prose
  positions that do not exist in the artifacts as the body describes, and reduce to a
  near-tautology (the FO substitutes a value it was handed into prose it authors itself,
  then observes the value). Nothing independent of the FO's authoring act can fail.

---

## Material

### M1 (`jm`, AC-1 + AC-2) — the localization targets do not exist; the ACs are near-tautological

**This is the highest-risk item in the sprint.** jm's two behavioral ACs name surfaces
to localize that, as described, are not in the artifacts under test.

**present-gate (AC-1).** jm's body (entity line 33, AC-1 line 62) claims the present-gate
`Decision:` prose currently says generic "entity" — e.g. "reject to bounce this **task**
back to …" — and the fix swaps it for `{entity-label}`. I checked `skills/present-gate/SKILL.md`.
The word "entity" appears in exactly three places, and **all three are positions jm's own
scope (entity line 33) declares OUT**:
- line 16 `{entity title}` — a field placeholder ("structure, not operating voice" per jm)
- line 20 `{entity_file_path}` — a field placeholder (same)
- line 37 "…or open the entity file" — FO *assembly-rule guidance*, never rendered to the captain

The actual `Decision:` template (SKILL.md line 29) is **noun-free**: its example reads
"reject to bounce back to `{feedback-to target}`" — there is no "entity"/"task"/"member"
to localize. So AC-1's claimed target ("the `Decision:` bounce-back prose reads 'ticket',
not 'entity'") **does not exist**. To make AC-1 pass, the implementer must *invent* a new
noun-bearing clause + `{entity-label}` placeholder, then render it and watch the FO echo
the label it was just handed. That is not "localize existing generic prose"; it is "add a
placeholder and prove a model can substitute." The render then carries almost no behavioral
content — it cannot fail unless the FO ignores a value already in hand.

**commander-dispatch (AC-2).** AC-2 says "the FO filling `dispatch-sprint-execution.md`
**from the localized template**." There is no template. I verified: no `skills/roadmap/`
exists (jm's own Out-of-scope, entity line 55, concedes "There is no `roadmap` skill in
`skills/` yet"), and `docs/roadmap/README.md` states the roadmap construct is a
"convention-only dry run … prose + frontmatter … **no new binary code**." Each
`dispatch-sprint-execution.md` is a per-sprint artifact the shaping FO **authors by hand
from scratch** (the 0198 / 019x ones are instances, not a fillable template). So "render
the declared plural from the localized template" reduces to: the FO writes "tickets" into
freshly-authored prose, then a check observes the prose says "tickets." The expected value
("tickets") nominally comes from the fixture README, but the *production* of the text is an
unconstrained authoring act with no template binding it — there is no artifact a future
edit could regress, so the check can never fail. This is the exact trap README `## ideation`
names: a check whose only proof is the deliverable's own prose, where "a passing check just
confirms the text the implementer wrote is present."

**What survives in jm:** only AC-3 (the shared contract stays generic) is a real,
independent, fail-able check — it rides the existing `internal/dispatch` `dispatch build`
golden, whose frozen expected text is independent of the skill edits and goes red on a leak.
That is sound. But AC-3 is a *guardrail on the change*, not the change's own positive proof.
A task whose only genuinely-checkable AC is "the thing I changed did NOT leak into the thing
I didn't change" is, by README `## done` policy, "a contract or skill change PASSED only when
a live drive observed the behavior it claims" — and there is no behavior here for a drive to
observe beyond an FO echoing a label.

**Fix (one of):**
1. **Reshape AC-1/AC-2 against a real, regressable artifact.** If localization is to ship,
   it needs a fillable artifact whose *baseline* says a generic noun, so a render can show
   the noun changed AND a guard can show it regressing. As scoped, present-gate has no such
   baseline position and commander-dispatch has no template at all. Either (a) define the
   exact generic-noun clause being replaced and confirm it exists in the SKILL.md baseline
   (it does not today), or (b) downscope jm to AC-3 only and record AC-1/AC-2 as deferred
   until the `spacedock:roadmap` skill (the real template generator) exists — which jm's own
   Out-of-scope already says is the prerequisite for the commander-dispatch half.
2. **At minimum**, validation must NOT accept a live-render of an FO echoing a just-handed
   label as proof of AC-1/AC-2 — that is the "FO writes X, observes X" tautology, not an
   independent check. If the gate is presented as clean on that evidence, it is the #262/`1x`
   self-referential-AC class the detached-audit policy exists to catch.

This should change jm's gate (or its AC set) **before** the Commander drives it. As written,
validation will either (a) flag AC-1/AC-2 as unprovable, stalling the drive, or (b) rubber-
stamp a tautological render — both bad outcomes the FO can prevent now.

### M2 (`th`, DoD#4 / gate semantics) — enforce the "not done on fake-smoke" boundary, and broaden the captain check

The FO's split — argv-oracles testable here, real-safehouse survival deferred to the
captain off-box — is honest, and I verified the testable half is real (`executablePath` is a
stubbable var; `withExecutablePath` exists; the AC-5 smoke correction from env-preserving
`exec "$@"` to env-scrubbing `unset SPACEDOCK_BIN; exec "$@"` genuinely turns a test that
"passes for the wrong reason" into a real regression oracle). The argv-oracle scope is NOT
overclaimed: the body's Spike + test plan state plainly it proves shape, not survival.

Two conditions the package must hold, or th ships unproven:

**(a) The gate boundary must be load-bearing, not narrative.** th's argv-oracles + corrected
fake-smoke ALL pass in this environment without the fix ever facing real safehouse. A drive
that sees green here could route th to `validation PASSED` and then to `done`. DoD#4 says the
captain's real-safehouse run is the proof — but nothing in the *entity* (frontmatter, AC set)
prevents a `verdict: PASSED` being stamped on the fake-smoke alone. The dispatch package
(line 52) states this correctly ("DO NOT present th as validated on the fake smoke alone");
ensure that lands as an actual hold on the verdict, not just prose the Commander may skim.
Concretely: th's validation can recommend "argv-shape PASSED, survival PENDING captain
real-safehouse" but th must NOT reach `done`/PASSED until DoD#4 is recorded. This matches
README `## done`: PASSED only "when a live drive observed the behavior it claims" — and the
claimed behavior is *survival across real safehouse*, which the fake-smoke does not observe.

**(b) The captain check has TWO untested assumptions, not one.** The body and package both
frame the captain run around "does the re-assert survive real safehouse" + "does `/usr/bin/env`
exist in the sandbox." But the proposed inner argv inserts a **bare `env` token**, not
`/usr/bin/env` (entity line 31: `… -- env SPACEDOCK_BIN=<bin> claude …`; AC-1 line 45 asserts
the literal `env` token). Bare `env` resolves only if `env` is on `$PATH` *inside the
sandbox*. The spike narration ("`/usr/bin/env` is present and portable") quietly substitutes
the absolute path for the bare token that actually ships — these are not the same assumption.
So the captain's real-safehouse run must confirm THREE things, and the package should say so:
1. real safehouse passes argv verbatim past `--` (untested here — only AC-6's non-scrubbing
   path and the controlled fake-stub exercise it);
2. the inner `env` token resolves on the sandbox `$PATH` (PATH survives per the spike, but
   `env`'s presence on it is untested);
3. `env VAR=val <cmd>` actually sets the var for the inner process under the real sandbox's
   exec semantics (not just the fake stub's).
If the captain's check covers only #1, a bare-`env`-not-on-PATH failure (which would make the
*launch itself* fail, not just lose the var) ships unproven. **Fix:** add #2 to DoD#4's
captain checklist, OR change the inner token to `/usr/bin/env` (the absolute path the spike
narration actually validated) so #2 collapses — the body should pick one, not leave the bare
`env` token while citing the absolute path as the proof.

---

## Polish

- **(`v3`, AC-3) live GitHub API in CI is a real flake source, not "fine."** AC-3 hits the
  unauthenticated GitHub "latest release" API on every CI run. The FO isolated the OS-logic
  test offline (good — AC-2's local-dist override), so a flake is contained to AC-3 alone, but
  AC-3 will still occasionally red on GitHub 5xx / rate-limit (unauthenticated = 60 req/hr/IP,
  shared across CI runners). Recommend: AC-3 asserts only URL *construction* against a pinned
  expected shape and treats a network failure as `t.Skip`, not `t.Fatal` — the production fetch
  is exercised for real by the first actual release anyway. Non-blocking; note it so a flake
  isn't misread as a regression.

- **(`v3`, AC-4) the doc grep is honestly self-flagged** as weak ("a grep over a doc is a weak
  proof on its own … the honesty of the safehouse wording is a human/reviewer judgment at the
  gate"). I confirmed `docs/install-journey.md` has no Linux/`curl|sh` line today, so the grep
  will genuinely fail-then-pass on the real edit — it is not a pure tautology (the expected
  string is a known-absent line being added). Acceptable as a presence check; the safehouse-
  honesty wording is correctly punted to gate judgment. No change needed.

- **(`v3`, AC-1) the config guard is genuinely independent.** I verified `internal/release`
  already parses YAML via `yaml.v3` (the existing workflow guards), the spike proved the
  cross-compile from `macos-latest` (which is what `release.yml`'s goreleaser job runs on, so
  the real release WILL produce the linux tarballs — pure-Go `CGO_ENABLED=0`, no separate linux
  runner needed), and the guard's expected value (the `{linux/amd64, linux/arm64}` target set)
  is independent of the YAML it reads (a config dropping linux fails it). This AC is sound.

- **(`v3`, AC-2) the tamper case is correctly load-bearing.** AC-2's "corrupt a tarball byte
  or its checksum line → install.sh exits non-zero, installs nothing" is the right shape: it
  proves the checksum gate, not just the happy path. The detached audit's mandate (package
  line 47, "tamper case must stay load-bearing") is the right thing to point the auditor at.

- **(`jm`, AC-4 deferral) the scope-record is clean for validation.** I confirmed jm's
  `## Feedback Cycles` cycle 1 records AC-4 deferred and the package (line 57) says validation
  requires AC-1/AC-2/AC-3 only. So validation will NOT flag AC-4 as missing evidence — the
  downscope is bounded correctly. (The problem is AC-1/AC-2 themselves, per M1, not AC-4.)

- **(`m1` deferral) correct.** rtk-only, already-caught, no shipped artifact — deferring it is
  right; building a guard for a one-off already-caught issue would be the YAGNI violation the
  FO avoided. No quarrel.

- **(detached-audit assignments) sound.** All three earn an audit (v3 = release machinery,
  th = front-door, jm = shipped scaffolding) per the four high-stakes surfaces in README
  `## validation`. The per-task audit targets in the package (v3 = checksum gate tamper-case;
  th = unwrapped-path unchanged + no-blank-value-on-unresolved; jm = no-leak-into-shared-
  contract + non-default-label render genuinely driven) are the right adversarial edges. Note
  for the jm auditor: per M1, "a non-default-label render is genuinely driven (not a prose
  assertion)" is precisely the edge that will expose the tautology — the auditor should try to
  pass AC-1 with an FO that echoes a hardcoded label and confirm the test cannot tell the
  difference. If it can't, that is a Material audit finding.

---

## Summary

Refuted one material hole (jm) and attached two binding conditions to th; v3 is sound. The
single highest-risk item is **M1: jm's AC-1/AC-2 localize prose that doesn't exist and reduce
to "the FO writes the label, then observes the label"** — a self-referential check the README
proof-policy bans and the detached-audit policy exists to catch. It should change jm's gate or
AC set (downscope to the real, code-gated AC-3; defer AC-1/AC-2 to the `spacedock:roadmap`
skill that is their actual prerequisite) before the Commander drives the sprint.

---

# Final review (post-correction)

Second pass after the sprint settled: `m1` DEFERRED (out); `jm` RE-IDEATED (cycle 2 reframe);
`47rx` ADDED and grown into a survey-skill body-rendering pass (cycles 1–5). Re-reviewed the
four ideated members against the actual code/config each claim rests on. (Out of this pass:
`frontdoor-launch-ux` = unideated backlog; `migration-check-share-walk-helper` = a later
fast-track refactor.) Read-only; I verified the claims by exercising the riskiest ones, not
by trusting the bodies.

## Verdict per member

- **`jm` entity-label-localization** — `M1 CLOSED`. The cycle-2 reframe from "edit template
  files" to "a generic operating-voice convention in the shared contract, proven by the FO's
  *generated* prose tracking the README label" removes the tautology I flagged, and the
  `ticket`-vs-`experiment` two-fixture differential makes a hardcoded/parroted label
  un-passable. The convention's home is real. Closed.
- **`47rx` survey body-rendering pass** — `sound`. All four deliverables (D1-depth-ii, C, E, F)
  have real targets I confirmed in the live SKILL.md / queries.sql; the gate-able query halves
  rest on the genuine sqlite3 query-smoke with non-tautological fixture-derived oracles (I
  reproduced the clustering extraction in pure sqlite3 — it yields the exact AC-3 labels); the
  rendered-body halves bottom out on a live drive per the survey discipline. No material hole.
- **`v3` ship-linux-binaries** — `held` (sound, unchanged). Body byte-identical to first review
  (no diff in the state repo); the soundness findings stand.
- **`th` safehouse-preserves-spacedock-bin** — `held` (M2 resolution accepted, sound). Body
  unchanged; the captain accepted both M2 conditions (verdict-hold + `/usr/bin/env` absolute +
  3-assumption captain check) at the package/scope level. The resolution is sound and
  sufficient — see the M2-resolution sanity-check below.

## Material

**Nothing material.** I tried to break all four and could not find a correctness, test-strength,
or proof-policy hole that should change a gate or the package. The two earlier material findings
are resolved: M1 (jm) is genuinely closed by the reframe; M2 (th) is accepted with a sound
resolution. Detail on the two verifications, since they were the load-bearing ones:

### M1 (jm) — verified CLOSED

The cycle-2 body (Problem ¶3, AC-1/AC-2, Test-plan "anti-tautology design") now targets a
generic convention, not a placeholder edit, and I confirmed every load-bearing claim:

- **The home exists.** jm puts the convention in `first-officer-shared-core.md`'s
  `## Working Principles` → `**FO posture:**` list. I confirmed that list exists (line ~304)
  with exactly the sibling bullets jm names ("Name the end value before starting", "Lead with a
  recommendation the captain can say yes to"). The present-gate reinforcement lands in the
  assembly rules, which exist. Unlike cycle 1, the edit has a real target.
- **The proof is non-tautological.** AC-1's expected noun ("ticket") lives ONLY in the fixture
  README's `entity-label` field — an independent artifact the FO must READ (Startup step 4) and
  RESOLVE, never handed to the FO by the test. The assertion is on the FO's *generated* prose
  vs that field — two artifacts that can diverge, with a real observable failure (a
  non-conforming FO emits "entity" and the drive reds). This is exactly the "expected value from
  OUTSIDE the file under test" the README `## ideation` policy demands.
- **The differential kills the residual loophole.** AC-2 runs the same drive over a SECOND
  fixture (`entity-label: experiment`). A hardcoded or test-handed label could fool one fixture
  but cannot match both READMEs — so two declarations yielding two different generated nouns
  proves field-resolution, not parroting. This directly answers the "FO writes X, observes X"
  objection from M1.
- **AC-3 (the no-leak guardrail) is still a real code gate** — the existing `dispatch build`
  golden (`go test ./internal/dispatch -run TestBuild`), whose frozen text is independent of the
  skill edits. Baseline green; the golden suite exists and runs.

The proof now rests on observing a behavior the FO can genuinely fail, with an independent
oracle and an anti-parroting differential. M1 is closed — not papered over.

### M2 (th) — resolution sanity-checked (accepted, not re-opened)

Per the brief the captain accepted both M2 conditions at the package/scope level; I confirm the
resolution is sound and sufficient:
- **(a) verdict-hold** — th must not reach `done`/PASSED on the in-env argv-oracles + corrected
  fake-smoke alone; survival is PENDING the captain's real-safehouse run. Correct: the in-env
  tests prove argv *shape*, and the README `## done` rule requires a live drive observing the
  *claimed* behavior (survival across real safehouse), which only the off-box run does.
- **(b) `/usr/bin/env` absolute + 3-assumption captain check** — switching the inner token from
  bare `env` to `/usr/bin/env` collapses the "is `env` on the sandbox PATH" assumption I raised
  (the absolute path is what the spike narration actually validated), and the captain check
  covering {verbatim-argv past `--`, the token resolves/execs under the real sandbox, the var
  reaches the inner process} closes the gap. Sound and sufficient. Not re-flagged.

## Polish

- **(`jm`, AC-1 negative control) make it an explicit assertion, not just a note.** The test
  plan mentions "the default-`entity`/dev workflow drive still says 'entity' (the convention is
  a no-op when label == 'entity')" as a negative control. Recommend the implementer assert it,
  not just observe it: a convention that accidentally hardcodes "ticket" or always-localizes
  would pass AC-1/AC-2 yet wrongly rewrite the default-`entity` workflow. The negative control is
  the cheap guard against an over-broad convention. Non-blocking — the `ticket`/`experiment`
  differential already catches the worst case; this catches the "label == entity" edge.

- **(`47rx`, AC-3) confirm the clustering is SQL in `queries.sql`, not model-applied prose.**
  AC-3's whole gateability claim is that the query-smoke runs the clustering. I verified it CAN
  be pure sqlite3: I reproduced the three-case extraction (`instr`/`substr`/`rtrim`, no regex)
  over the exact AC-3 fixture messages and it returned `journey-cost-ledger`, `codex-live-ci`,
  `orient-workflow-discovery`, `(unlabeled)` with the stage suffix stripped and the two dispatch
  tasks not merged. The implementer MUST land the clustering as a labeled `-- name:` query block
  in `queries.sql` (so the smoke extracts and runs the artifact), NOT as a prose rule only the
  model applies — otherwise AC-3's "query-smoke" proof evaporates into a SKILL.md-grep-equivalent
  the survey discipline bans. The cycle-4 note ("validated the clustering SQL … corrected a
  backtick off-by-one") implies it is SQL; pin it. Non-blocking but load-bearing for AC-3's
  proof tier.

- **(`47rx`, AC-7 / F) the `<external>` mechanics tag spans two branches — tag both.** F's
  mechanics set is `<external>` + `.worktrees` + `.claude*`. In the live `work-by-area` query,
  `<external>` is the ELSE branch (paths not under repo root) while `.worktrees`/`.claude*` are
  first-segment buckets UNDER the root. The `kind` partition must tag the `<external>` ELSE
  branch AND the in-root mechanics segments as `mechanics`, leaving every other in-root segment
  `product`. AC-7a's non-vacuousness flip (re-segment a `.worktrees` row to `internal` → it moves
  into the product lead) correctly exercises the in-root case; add a parallel check that an
  `<external>` row also sorts below product, so both mechanics branches are proven demoted. Minor
  — the fixture already carries an `<external>` row, so it is one assertion, not new fixture data.

- **(`47rx`, AC-4 live-drive flakiness) the rendered-body labels depend on a live synced DB.**
  AC-4/5/6/7b are proven by ONE combined live drive over the real agentsview-synced DB. That is
  the right proof bar (the survey discipline forbids a SKILL.md grep), but the rendered Codex
  workstream labels (`journey-cost-ledger`, `orient-workflow-discovery`) come from THIS repo's
  live Codex history — which drifts as new sessions land. The drive should assert the *shape*
  (a Codex section with workdir-attributed count + ≥1 named workstream cluster + activity), not a
  frozen label list, so a future session that shifts the top workstreams doesn't red the drive.
  The body already says "labels come from the real synced DB" — just keep the assertion on
  presence-of-clusters, not a pinned roster. Non-blocking; the body permits a captain-run drive
  escalation (as vh's AC-4/AC-5 did), which is the right pressure valve.

- **(`47rx`, scope) depth (ii) is honestly bounded; the deferrals are clean.** I confirmed the
  body keeps depth (iii) (per-file work-by-area, which needs raw-rollout `apply_patch` parsing —
  genuinely not in the agentsview DB) and D3 (source-health, no agentsview signal) as
  upstream-gated follow-ups with a concrete recorded `kenn-io/agentsview` dependency. No
  over-reach: D1(ii) reads only DB signals (`first_message`, `exec_command`/`update_plan`
  tallies), no rollout parsing. The attribution spike (0 sibling cross-contamination across 160
  sibling sessions, 93% recall on a 95%-present `$.workdir` field) is a real exercised mechanism,
  not an assertion. m1's deferral remains correct. No quarrel with any scope call.

## Summary (final review)

Re-reviewed the four ideated 0199 members. **jm's M1 is genuinely CLOSED** — the reframe to a
shared-contract operating-voice convention proven against the FO's generated prose, with a
`ticket`-vs-`experiment` differential, removes the tautology and binds the expected noun to an
independent source (the fixture README). **47rx is sound** — I confirmed C/E/F's targets exist
verbatim in the live SKILL.md/queries.sql, the query-smoke is a real sqlite3 oracle, and I
reproduced AC-3's clustering extraction in pure sqlite3 to confirm it is gateable and
non-tautological. **v3 and th are held** — bodies unchanged; th's accepted M2 resolution is
sound and sufficient. **No material findings this pass**; five polish notes, the load-bearing
one being "land 47rx's clustering as a `queries.sql` query block, not model-applied prose, so
AC-3's query-smoke proof is real."

---

# yq review (frontdoor-launch-ux)

Late-added 0199 member, reviewed after the comprehensive pass. The highest-stakes surface (the
front door every user hits first) and the one whose defect A shipped in 0.19.8 *because nobody
ran the front door end-to-end* — so I was hard on it and verified every line-anchored claim
against the real `internal/cli/frontdoor.go` (and against `th`'s actual landed worktree, since
`th` is now in validation).

**Verdict — `sound`.** A is a correct message-ownership fix that preserves the fail-fast
verdicts and derives the `--no-install` remedy from `host` (no new hardcoded `claude`); the
ACs each pair an independent-source oracle with a genuinely load-bearing (not skippable) live
drive and none degrade to a frontdoor.go grep; the `th` overlap is substantively disjoint
(verified against th's real code, not the body's word); "no spike needed" is justified. Two
polish notes, no material hole.

## Material

**Nothing material.** I tried to break A's verdict-preservation, the AC proof tiers, and the
th-disjointness claim against the actual code, and found no correctness/test-strength/proof
hole that should change yq's gate or the package. The verifications, since A is the load-bearing
one:

### A (message ownership) — verified correct, fail-fast preserved

- **The defect is real and the diagnosis is right.** I confirmed `gateHost` (frontdoor.go
  116-142) prints a remedy for *every* non-Compatible verdict, while the caller's NoPluginFound
  response is non-uniform: `runCodex`'s arm is `if fd.noInstall { return 1 }` then
  `ops.Install("codex", …)` — and its own comment says "gateHost already printed the instruct
  remedy, so just fail fast." So today the `--no-install` refuse path has NO print of its own; it
  rides gateHost's. That is exactly why the remedy fires on the auto-install path too, and why
  the hardcoded `spacedock claude` hint (gateHost's NoPluginFound text) is wrong in a codex run.
- **Deleting the two NoPluginFound prints does NOT regress the fail-fast verdicts.** A keeps
  gateHost's resolve-error print (119-121 → `MalformedRange`) and the mismatch
  `Fprintln(res.Message)` (139 → not-Compatible). Both are verdicts the caller ALWAYS fails fast
  on, so their message must stay at gateHost — and A keeps them. The mismatch still prints its
  `res.Message`. Correct.
- **A MUST add the `--no-install` print, and the body specifies it.** Because the refuse arm is
  silent today (it depended on gateHost), deleting gateHost's NoPluginFound print would leave
  `spacedock codex --no-install` emitting nothing unless the caller prints. The body's §A (line
  46) requires exactly this: a host-correct `noPluginRemedy(host)` before `return 1` in the
  `--no-install` branch. If the implementer deletes the gateHost prints but forgets the caller
  print, `TestGateRemedyNamesLiveInstallCommand`'s "no plugin"/"missing manifest" cases go RED
  (they assert the remedy via gateHost today) — so the existing suite catches that miss, and
  AC-A(2) correctly moves those two cases to assert on the launcher's `--no-install` output while
  keeping the resolve-error case on gateHost. The test reconfiguration is sound.
- **The remedy is genuinely host-derived.** `noPluginRemedy(host)` interpolates `host` into both
  `spacedock install --host {host}` and `spacedock {host} --skip-contract-check` — no residual
  hardcoded `claude`. AC-A asserts the codex path does NOT contain `spacedock claude`, which is
  the real anti-regression for the original defect.
- **Deleting the phantom-manifest message loses no needed signal.** Both NoPluginFound
  sub-states (empty manifestPath, resolved-but-missing manifest) mean the same thing to the
  operator: no usable plugin, installing now. Collapsing them to one caller-owned message is the
  right call — the "manifest path missing vs no entry" distinction is internal, not actionable.

### AC-A/B/D — non-grep, live-drive load-bearing

Each AC pairs (i) an independent-source oracle and (ii) a real front-door live drive, and the
live leg is explicitly framed as load-bearing, not decorative:
- **Independent sources are real.** AC-D's `wantBootstrapPrompt`/`wantCodexBootstrapPrompt` are a
  SECOND hand-written copy in `safehouse_frontdoor_test.go` (line 18/156) — production drift fails
  the argv-equality oracle and a residual "I love you" fails the absence assertion; the expected
  value is not in frontdoor.go. AC-B's discovered-workflow expectation comes from the fixture
  README the test writes (`commissioned-by: spacedock@…`), an independent artifact. AC-A's
  expected stderr strings live in the test file / live terminal.
- **The live drive is not quietly skippable.** A shipped precisely from a coverage gap (tests
  passed, nobody ran the front door), so the body makes the live drive the named proof and the
  detached adversarial audit at validation re-runs it over the `--no-install`/resume/outside-
  workflow edges. Host CLIs are present on the box (`codex`/`claude` resolve), so the leg is
  runnable — not a deferred off-box hand-wave like th's real-safehouse. No AC rests on a
  frontdoor.go grep; the body bans it explicitly and the bans are real (the oracles read
  `fakeHost.launchedArg`/stderr buffers and test-file copies, never the source).

### th overlap — verified disjoint against th's REAL code (not the body's word)

`th` is in validation with its changes in `.worktrees/spacedock-ensign-safehouse-preserves-
spacedock-bin`, NOT yet merged. I diffed th's actual `frontdoor.go` against the current one. th
adds a new `envBinary` const + `launcherBinArgvPrefix()` helper (new top-level decls) and
changes ONE statement in each launcher: `argv = safehouse.Wrap(inner, extra)` →
`argv = safehouse.Wrap(append(launcherBinArgvPrefix(), inner...), extra)` at the `if wrap {`
block (~216 / ~362). yq's edits — `gateHost` body (116-142), the `case contract.NoPluginFound:`
caller arms, a banner print sited BEFORE `inner`/`argv` assembly (after `warnStrayPromptAfterDash`,
~197/344), and the two consts (24/289) — touch no statement th touches. The banner sits *above*
th's wrap-site edit in the same function; textual adjacency at worst, never a semantic collision.
yq's disjointness conclusion is correct, and th's worktree `internal/cli` suite is green (I ran
it), so the "th lands first, yq rebases cleanly" sequencing is realistic. (One body imprecision
noted under Polish.) I also confirmed th's landed code uses `/usr/bin/env` (absolute) — the M2
resolution is in the actual deliverable, not just the package note.

### "No spike needed" — justified

The banner's only non-trivial dependency is `status.DiscoverWorkflowDir(dir)`. I confirmed it is
public (`internal/status/discover_walkup.go:16`) and is the same walk-up native `status` already
uses (walk up from `dir`, match the nearest ancestor README with `commissioned-by: spacedock@`).
B claims no new behavior — it reuses a proven walk-up — and AC-B tests both arms (a fixture dir
with a commissioned-README ancestor → names the rel path; a bare temp dir → `none detected`)
with independent fixture expectations. There is no unexercised mechanism: no parser round-trip,
no on-disk format, no runtime handoff, no tool-flag assumption. The real risk was the coverage
gap, and the live drive closes it. Justified. (Note: from inside a git worktree the walk-up
finds a workflow only when an ancestor of `dir` carries the commissioned README — same as status
today; the banner honestly renders `none detected` otherwise, which is correct, not a bug.)

## Polish

- **(yq, th-overlap body wording) one inaccuracy in how the body describes `th`.** The Overlap
  section (line 111) says th "factors `resolvedLauncherBin()`". It does not — `resolvedLauncherBin()`
  already existed (frontdoor.go:75, pre-th baseline). th only ADDS `launcherBinArgvPrefix()`,
  which CALLS the pre-existing `resolvedLauncherBin()`. The disjointness conclusion is unaffected
  (the new helper is a fresh top-level decl), so this is wording-only — but if a Commander reads
  the body to predict the merge, the description should say "th adds `launcherBinArgvPrefix()`
  calling the existing `resolvedLauncherBin()`", not "factors" it. Non-blocking.

- **(yq, AC-B banner-vs-resume) pin the resume-suppression decision before implementation.** B
  says the banner is "suppressed on `--resume`/`codexResume` … though if simpler to always print,
  that is acceptable; the AC pins the no-resume case." AC-B only tests the no-resume case, so an
  implementer could ship either behavior and pass. That is a real under-specification on a
  user-facing surface: a banner printed on every `--resume` adds noise to the exact path
  (continuing a session) where the user least wants orientation chrome. Recommend the AC pin the
  resume case too (suppressed, matching the bootstrap-prompt's existing resume-suppression via
  `containsResume`), so validation can't pass a noisier-than-intended banner. Minor, but it's the
  front door — and the body already leaves the door open to drift. The detached audit at
  validation should probe the resume edge regardless (the body lists it), which mitigates this.

## Summary (yq review)

yq is `sound`. A correctly relocates NoPluginFound message ownership from `gateHost` to the
caller — I verified it preserves the always-fail-fast verdict prints (resolve-error, mismatch),
adds the host-derived `--no-install` remedy the silent refuse arm needs today, and kills the
hardcoded `spacedock claude` hint; the existing `TestGateRemedyNamesLiveInstallCommand` will
catch a half-done edit, and AC-A correctly moves the NoPluginFound cases to the launcher. The
ACs are non-grep with a load-bearing, runnable live drive (the front-door run the 0.19.8 gap
missed). The `th` overlap is genuinely disjoint — verified against th's real landed worktree
code, not the body's claim. "No spike needed" holds (`DiscoverWorkflowDir` is a proven walk-up).
No material findings; two polish notes (a wording inaccuracy in the th-overlap description, and
an under-pinned banner-on-resume case).
