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
