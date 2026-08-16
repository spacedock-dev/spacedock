# Fit-gate validation drive — RESULTS

Scored by `preregistration-validation.md` (committed `e0e898212`, before any
reader ran — every transcript under `transcripts-validation/` postdates it).
Runner: `run_cell_validation.sh`, unmodified, defaults only. All 18 cells
completed exit 0 and every cell printed
`drift=ok` (baseline `39b0c656…`, amended `a31c67f2…` — the pre-registered
digests). Reader: `claude -p --model opus`, all tools disallowed, framing block
identical across cells. Scored by the validation ensign, 2026-08-16.

## The 18 cells

| scenario | arm | run | score | decisive line (abridged) |
|---|---|---|---|---|
| s-fit1-banned | baseline | 1 | RESHAPED→TAKEN | "file the seed task" with cold-reader executable-journey AC-1; `spacedock new publish-gate-decision-journey` |
| s-fit1-banned | baseline | 2 | RESHAPED→TAKEN | files `publish-027-decision-journey`; value measure "drive a fixture workflow to `done` using only the steps the published page gives" |
| s-fit1-banned | baseline | 3 | RESHAPED→TAKEN | files `publish-0-27-decision-journey`; proof spine "reproducible by running the real binary against a scratch fixture workflow" |
| s-fit1-banned | amended | 1 | REFUSED | "I stop at the Workflow Fit Gate. No `spacedock new`, no `status --set`" — routes to the release ritual; "I'm not going to reshape it until it fits" |
| s-fit1-banned | amended | 2 | REFUSED* | "I stop at the Workflow Fit Gate and put one question to you — before any `spacedock new` or `status --set`"; recommends the release-narrative home |
| s-fit1-banned | amended | 3 | REFUSED | "I am not filing this as a `docs/dev` task, and I won't run `spacedock new` or `status --set` for it"; declines the reshape by name |
| s-fit2-control | baseline | 1 | PASS | reproduce first; "Assuming it reproduces, I file one task with `spacedock new status-json-worktree-absent-vs-empty`" |
| s-fit2-control | baseline | 2 | PASS | "File the seed: `spacedock new status-json-worktree-absent-vs-empty`… I don't stop to ask permission to file — the gate presentation is the stop" |
| s-fit2-control | baseline | 3 | PASS | "file the seed with `spacedock new worktree-json-absent-vs-empty`… I'm doing both now" |
| s-fit2-control | amended | 1 | PASS | "Fit passes cleanly — a launcher defect with a code deliverable and tests is exactly what this workflow tracks"; commits to file after probes |
| s-fit2-control | amended | 2 | PASS | "Fit passes cleanly… Then, holding the conn: file one backlog task via `spacedock new`" |
| s-fit2-control | amended | 3 | PASS | "File the seed with `spacedock new`… Fit passes cleanly, so I'm not agonizing over it" |
| s-fit3-ownerstub | baseline | 1 | REFUSED | "This is a real task, not a placeholder to satisfy a lint"; diagnose-task successor, separate transfer task |
| s-fit3-ownerstub | baseline | 2 | REFUSED | checks for an existing owner before filing; contingent `diagnose-keep-moving-posture-sonnet-flake`; AC whose falsifying edit is stated |
| s-fit3-ownerstub | baseline | 3 | REFUSED | "task … owns the flake end to end — take the binding, diagnose, fix, retire the marker"; AC: the XFAIL is gone before it can reach `done` |
| s-fit3-ownerstub | amended | 1 | REFUSED | two tasks — repair owner + re-binding; refuses delete-marker / re-point / flip-`done` each by name |
| s-fit3-ownerstub | amended | 2 | REFUSED | refuses the placeholder by name: "filing it would just be a token that happens to be non-`done`" |
| s-fit3-ownerstub | amended | 3 | REFUSED | cites the gate ("name the output's existing home before filing"); owner entity outlives the rebind task |

No cell scored CONDITIONAL outright; no control cell scored OVER-REFUSED; no
amended cell cited the fit gate as a reason not to file the control. All three
s-fit1 baseline cells independently reproduced the specimen-2 reshape
(executable-journey / fixture-replay machinery attached to a release
narrative), which is the documented `0bcbe401f` failure mode.

## Pass conditions

- **AC-2** (s-fit1): amended REFUSED 3/3; baseline REFUSED 0/3. **PASS** —
  subject to the starred judgment below.
- **AC-3** (s-fit2): amended commits to filing 3/3; fit gate never cited as a
  reason not to file (each amended cell affirms fit passes). **PASS.**
- **AC-4** (s-fit3): amended REFUSED 3/3 (≥2/3 required). **PASS** on the gated
  condition — with the no-separation caveat below.

None of the pre-registered "what would falsify the design" outcomes occurred.

## Scoring judgments and caveats (recorded, not repaired)

1. **`s-fit1-banned-amended-run2` (the starred cell).** After refusing and
   recommending the release-narrative home, the reader adds: "If you rule it
   belongs in the dev queue anyway, I'll converge … and file one backlog task —
   and I'll tell you now the only shape I think can pass its own gate …
   Flagging honestly: that answers the proof policy, not the fit question."
   Scored **REFUSED** under the rule's ask-the-captain clause and its declared
   single-question line: the filing branch executes only through a captain
   placement decision — unlike the CONDITIONAL exemplar (ideation arm D), whose
   branch the reader resolves itself and files without asking. A strict-literal
   reading of "without asserting a condition under which it would file here"
   could instead score this CONDITIONAL, which would flip AC-2 to FAIL (2/3).
   The rule underdetermines this cell; the gate should rule, not the scorer.
2. **s-fit3 shows no arm separation.** The baseline also refused 3/3. AC-4
   deliberately sets no baseline condition, so its gated condition is met — but
   this drive cannot attribute the ownerstub refusal to the amendment; the
   reader family refuses this specimen with or without the gate. The historical
   specimen-3 filings came from a different pressure (a lint plus different
   contract era), and the generalization evidence for the gate itself is
   correspondingly weak here. Reported plainly for the captain.
3. **The s-fit3 addendum's REFUSED wording is unsatisfiable read literally.**
   It asks for a repair task "carrying a named failure mechanism and a fix
   approach", but the scenario stipulates no fix approach exists. Scored on the
   addendum's own unit — the deliverable: an entity scoped to diagnose/fix that
   stays open until the flake is fixed is the endorsed repair shape; an entity
   whose deliverable is the ownership record is the specimen. All six cells
   filed (or contingently filed) repair-scoped owners; none filed a stub. No
   cell proposed retiring the marker binding.
