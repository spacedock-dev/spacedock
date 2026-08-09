# Test-behavior-completeness staff review

Date: 2026-08-09

Verdict: **READY FOR IDEATION GATES**. The folds close all seven Material
findings. No Material finding remains.

This review excludes `g3`, `47g`, and all durable-decisions work. It does not
approve product changes or edit task state.

## Evidence read

The review reread the corrected bodies for `ts`, `zh`, `6x5`, `9a`, and `xp6`.
It also read all three new task bodies and their complete ideation reports:

- `2e4` — commit the Pi gate before presentation
- `dvd` — continue Codex rejection after the first validation
- `fh6` — hold the Pi default headless validation gate

The review checked the source bindings, shared assertions, expected surfaces,
line budgets, dependencies, and serial landing order.

The workflow query remains the membership authority:

```bash
spacedock status --workflow-dir docs/dev \
  --where sprint=test-behavior-completeness \
  --where 'slug != define-fo-moving-target-conflict-ownership' \
  --where 'sprint-readiness != defer'
```

## Prior Material findings

### M1 — Strict XFAIL multi-error safety: closed

`ts` now runs every durable semantic assertion after infrastructure succeeds.
XFAIL requires the unique code set to equal `{expected}`. An empty set is XPASS.
An additional or different code is FAIL.

The first landing also changes two real journey results. Sonnet and Codex
`default-headless-gate-stop` change from TODO to executed XFAIL. A helper-only,
parser-only, or metrics-only landing is banned.

### M2 — Rejection mechanisms: closed

`zh` now owns only stable recorder publication. Pi has exact evidence for
`rejection-round-incomplete`. Sonnet and Opus need matching exact evidence before
they join that repair.

`dvd` owns the separate Codex continuation value. Its target is one fresh open
gate after correction and validation/2. It does not change recorder bytes.

### M3 — Two Pi mechanisms: closed

`2e4` owns prepare, commit, reread, and presentation order for Pi
`gate-guardrail`.

`fh6` owns the separate terminal-field failure in Pi
`default-headless-gate-stop`. The artifact identifies `completed` and `verdict`
as the exact failed clause. The task does not assume the `98a` mechanism.

### M4 — Owner transfer and shared-file order: closed

`6x5` now transfers only the three `smallest-sufficient-mechanism` rows from
`9a` to `6x5`. The transfer and XFAIL baseline precede product bytes.

`9a` keeps only the three `keep-moving-posture` rows. It rebases after `6x5` and
reruns the exact candidate.

The sprint now declares one serial merge order for all shared files.

### M5 — Evidence capstone boundary: closed

`xp6` remains evidence-only and net zero. It removes four passing bindings after
exact runs. Two Pi repairs stay with `2e4` and `fh6`.

The passed Codex withdrawn-gate row can lose its TODO without adding `47g` to
the sprint. The two Opus TODO rows remain only while authenticated execution is
unavailable.

## Final delta findings

### M6 — Codex baseline order: closed

The `dvd` spike found that the shared runner uses the Claude round extractor for
Codex JSONL. That path reports a false missing-recorder error. The task states
that the Codex extractor is required before the baseline can isolate
`rejection-flow-not-completed`.

The corrected task now uses this exact order:

1. Rebase after `ts` and `zh`.
2. Add the existing Codex extractor selection and strict final-gate oracle.
3. Add the Codex XFAIL binding in the same non-product baseline commit.
4. Run the real Codex journey and require the sole code
   `rejection-flow-not-completed`.
5. Change the feedback-flow and Codex runtime behavior only after that XFAIL.
6. Require XPASS with the binding, then PASS after removal.

The extractor commit cannot land alone. The final task landing changes the real
Codex journey from XFAIL to PASS with one fresh open gate.

### M7 — Pi headless XFAIL ownership: closed

`fh6` now uses this source row:

```text
liveXFail("pi", "fh6rv0k6wr25zty0jjan4jp7", "gate-hold-terminal-fields-set")
```

The body, acceptance criteria, evidence, and baseline plan use that same owner.
`xp6` remains evidence-only. The binding cannot name `xp6` at any repair stage.

## End-user value refutation

| Member | Required journey result | Review result |
| --- | --- | --- |
| `0a` | A manual Pi cadence runs common journeys and retains evidence. | Pass. A workflow-only change is insufficient without the exact manual run. |
| `ts` | Two real headless cells run as strict XFAIL. | Pass after M1. A classifier-only landing is banned. |
| `98a` | Sonnet and Codex complete implementation before validation. | Pass. Both exact live cells must change from XFAIL to PASS. |
| `6x5` | Each initial-stage worker leaves a ready report before terminal state. | Pass. The durable journey counts two reports and two archives. |
| `9a` | A consumed gate dispatches once and reaches non-forced terminal state. | Pass. The live value includes the exact dispatch commit and terminal result. |
| `zh` | Pi retains one complete rejected round before re-gating. | Pass. A recorder-text change without a real Pi result is banned. |
| `dvd` | Codex reaches one fresh open gate after validation/2. | Pass after M6. The extractor cannot land without the journey repair. |
| `2e4` | Pi presents a review only after the gate is committed and reread. | Pass. The exact Pi cell must change from XFAIL to PASS. |
| `fh6` | Pi reaches an open gate with zero terminal fields. | Pass after M7. The binding names `fh6` as repair owner. |
| `xp6` | Four passing cells stop skipping and expose current results. | Pass. It can remove bindings only, with no repair bytes. |

No retained task ends with a helper, parser, binding, contract, metrics, fixture,
or test-only result. `dvd` can include parser work only with its final journey
repair. `ts` can include common infrastructure only with its two first XFAIL
cells.

## Strict-XFAIL and exact-candidate audit

- `98a` starts from the two `ts` XFAIL records.
- `6x5` starts from three owner-correct XFAIL records.
- `9a` starts from each executable `keep-moving-posture` cell with one sole code.
- `zh` starts from Pi `rejection-round-incomplete`.
- `dvd` puts the extractor, oracle, controls, and binding in one baseline commit.
- `2e4` starts from Pi `gate-prepare-state-commit-missing`.
- `fh6` starts from Pi `gate-hold-terminal-fields-set` under its own task ID.

Each repair keeps its binding through one XPASS run. Binding removal then uses
the same target and exact candidate for PASS.

## `0a` recovery audit

`0a` remains in implementation. Its registered checkout exists at the expected
path. The registered branch tip is `e838fba69`.

The entity has no implementation Stage Report. The Commander must not treat the
candidate as finished. Cold boot first checks durable state, the registered
checkout, branch tip, and active owner availability. Normal ownership rules then
select reuse or fresh implementation dispatch.

`0a` validates and merges before `ts`.

## Shared surfaces and landing order

The required serial merge order is:

```text
0a -> ts -> 98a -> 6x5 -> 9a -> zh -> dvd -> 2e4 -> fh6 -> xp6
```

This order covers these collisions:

- `docs/runtime-live-ci.md`: `0a`, `ts`, and `dvd`
- `internal/ensigncycle/shared_live_runner_test.go`: all evidence and repair tasks
- `skills/first-officer/references/fo-dispatch-core.md`: `98a`, `6x5`, and `9a`
- `skills/feedback-rejection-flow/SKILL.md`: `zh` and `dvd`
- `skills/first-officer/references/pi-first-officer-runtime.md`: `2e4` and `fh6`

Workers can run independent investigation and focused tests in parallel. The
Commander owns all shared-file conflict resolution and serial merges.

## Line and capstone audit

- `0a`: net must stay within its approved 8-file, 380-insertion reset.
- `ts`: hard cap `+210` net.
- `98a`: about `+6` net.
- `6x5`: about `+12` net.
- `9a`: about `+228` net, with 25% tolerance.
- `zh`: about `+2` net.
- `dvd`: about `+26` net.
- `2e4`: about `+14` net.
- `fh6`: about `+12` net.
- `xp6`: product net `0`.

`xp6` cannot absorb a product fix, helper, parser, fixture, or metrics change.
If no exact passing binding remains, it lands no product source change.

## Readiness result

The delta read found M6 and M7 in durable task state. Every retained task has a
visible journey result. Every product repair starts from strict XFAIL. The
shared surfaces have one serial landing order. `xp6` remains evidence-only.

No Material finding remains. The Shaping FO can record and present the ideation
gates. Commander execution still waits for those approvals, the target-train
decision, and captain activation.
